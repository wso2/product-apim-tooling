/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package aws

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ghodss/yaml"

	"github.com/wso2/product-apim-tooling/import-export-cli/box"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	yaml2 "gopkg.in/yaml.v2"

	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagApiNameToGet string //name of the api to get from aws gateway
var flagStageName string    //api stage to get
var dir string              //dir where the aws init command is executed from
var path string             //path of the OAS extracted from AWS
var tmpDir string           //temporary directory created to store the OAS extracted from aws till the project is initialized
var getRestAPIsCmdOutput string
var err error
var awsInitCmdForced bool
var initCmdOutputDir string

//common aws cmd flags
var apiGateway string = "apigateway"

//aws cmd Output type
var outputFlag string = "--output"
var outputType string = "json"

//getAWSCLIVersionCmd
var awsCLIVersionFlag string = "--version"

//getRestAPIsCmd
var getAPI string = "get-rest-apis"

//getExportCmd
var getExport string = "get-export"
var apiIdFlag string = "--rest-api-id"
var stageNameFlag string = "--stage-name"
var exportTypeFlag string = "--export-type"
var exportType string = "oas30" //default export type is openapi3. Use "swagger" to request for a swagger 2.
var debugFlag string            //aws cli debug flag for apictl verbose mode

const awsInitCmdLiteral = "init"
const awsInitCmdShortDesc = "Initialize an API project for an AWS API"
const awsInitCmdLongDesc = `Downloading the OpenAPI specification of an API from the AWS API Gateway to initialize a WSO2 API project`
const awsInitCmdExamples = utils.ProjectName + ` ` + awsCmdLiteral + ` ` + awsInitCmdLiteral + ` -n Petstore -s Demo
` + utils.ProjectName + ` ` + awsCmdLiteral + ` ` + awsInitCmdLiteral + ` --name Petstore --stage Demo
` + utils.ProjectName + ` ` + awsCmdLiteral + ` ` + awsInitCmdLiteral + ` --name Shopping --stage Live

NOTE: Both the flags --name (-n) and --stage (-s) are mandatory as both values are needed to get the openAPI from AWS API Gateway.
Make sure the API name and the Stage name are correct.
Also make sure you have AWS CLI installed and configured before executing the aws init command.
(Vist https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html for more information)`

func getPath() {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	dir = pwd
}

//aws init Cmd
var InitCmd = &cobra.Command{
	Use:     awsInitCmdLiteral,
	Short:   awsInitCmdShortDesc,
	Long:    awsInitCmdLongDesc,
	Example: awsInitCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		getPath()
		initCmdOutputDir = dir + string(os.PathSeparator) + flagApiNameToGet

		if stat, err := os.Stat(initCmdOutputDir); !os.IsNotExist(err) {
			fmt.Printf("%s already exists\n", initCmdOutputDir)
			if !stat.IsDir() {
				fmt.Printf("%s is not a directory\n", initCmdOutputDir)
				os.Exit(1)
			}
			if !awsInitCmdForced {
				fmt.Println("Run with -f or --force to overwrite directory and create project")
				os.Exit(1)
			}
			fmt.Println("Running command in forced mode")
		}
		execute()
	},
}

type Apis struct {
	Items []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}

func getOAS() error {
	utils.Logln(utils.LogPrefixInfo + "Executing aws version command")
	//check whether aws cli is installed
	getAWSCLIVersionCmd, err := exec.Command(awsCmdLiteral, awsCLIVersionFlag).Output()
	if err != nil {
		utils.HandleErrorAndExit("Error getting AWS CLI version. Make sure AWS CLI is installed and configured.", err)
		return err
	}
	output := string(getAWSCLIVersionCmd[:])
	utils.Logln(utils.LogPrefixInfo + "AWS CLI version :  " + output)

	if utils.VerboseModeEnabled() {
		debugFlag = "--debug" //activating the aws cli debug flag in apictl verbose mode
	}
	utils.Logln(utils.LogPrefixInfo + "Executing aws get-rest-apis command in debug mode")

	getRestAPIsCmd := exec.Command(awsCmdLiteral, apiGateway, getAPI, outputFlag, outputType, debugFlag)
	stderr, err := getRestAPIsCmd.StderrPipe()
	if err != nil {
		utils.HandleErrorAndContinue("Error creating pipe to standard error. (get-rest-apis command)", err)
	}
	stdout, err := getRestAPIsCmd.StdoutPipe()
	if err != nil {
		utils.HandleErrorAndExit("Error creating pipe to standard output (get-rest-apis command).", err)
	}

	err = getRestAPIsCmd.Start()
	if err != nil {
		utils.HandleErrorAndExit("Error starting get-rest-apis command.", err)
	}

	if utils.VerboseModeEnabled() {
		logsScannerGetRestAPIsCmd := bufio.NewScanner(stderr)
		for logsScannerGetRestAPIsCmd.Scan() {
			fmt.Println(logsScannerGetRestAPIsCmd.Text())
		}

		if err := logsScannerGetRestAPIsCmd.Err(); err != nil {
			utils.HandleErrorAndContinue("Error reading debug logs from standard error. (get-rest-apis command)", err)
		}
	}

	outputScannerGetRestAPIsCmd := bufio.NewScanner(stdout)
	for outputScannerGetRestAPIsCmd.Scan() {
		getRestAPIsCmdOutput = getRestAPIsCmdOutput + outputScannerGetRestAPIsCmd.Text()
	}

	if err := outputScannerGetRestAPIsCmd.Err(); err != nil {
		utils.HandleErrorAndExit("Error reading output from standard output.", err)
	}
	err = getRestAPIsCmd.Wait()
	if err != nil {
		utils.HandleErrorAndExit("Could not complete get-rest-apis command successfully.", err)
	}

	//Unmarshal from JSON into Apis struct.
	apis := Apis{}
	err = json.Unmarshal([]byte(getRestAPIsCmdOutput), &apis)
	if err != nil {
		return err
	}
	extractedAPIs := strconv.Itoa(len(apis.Items))
	utils.Logln(utils.LogPrefixInfo + extractedAPIs + " APIs were extracted")

	var found bool
	apiName := flagApiNameToGet
	stageName := flagStageName
	path = tmpDir + string(os.PathSeparator) + apiName + ".json"
	// Searching for API ID:
	utils.Logln(utils.LogPrefixInfo + "Searching for API ID...")
	for _, item := range apis.Items {
		if item.Name == apiName {
			utils.Logln("API ID found : ", item.Id)
			api_id := item.Id
			found = true

			utils.Logln(utils.LogPrefixInfo + "Executing aws get-export command in debug mode")
			getExportCmd := exec.Command(awsCmdLiteral, apiGateway, getExport, apiIdFlag, api_id, stageNameFlag, stageName, exportTypeFlag,
				exportType, path, outputFlag, outputType, debugFlag)

			stderr, err := getExportCmd.StderrPipe()
			if err != nil {
				utils.HandleErrorAndContinue("Error creating pipe to standard error. (get-export command)", err)
			}
			stdout, err := getExportCmd.StdoutPipe()
			if err != nil {
				utils.HandleErrorAndContinue("Error creating pipe to standard output. (get-export command)", err)
			}

			err = getExportCmd.Start()
			if err != nil {
				utils.HandleErrorAndExit("Error starting get-export command.", err)
			}

			if utils.VerboseModeEnabled() {
				logsScannerGetExportCmd := bufio.NewScanner(stderr)
				for logsScannerGetExportCmd.Scan() {
					fmt.Println(logsScannerGetExportCmd.Text())
				}
				if err := logsScannerGetExportCmd.Err(); err != nil {
					utils.HandleErrorAndContinue("Error reading debug logs from standard error. (get-export command)", err)
				}
			}

			if utils.VerboseModeEnabled() {
				outputScannerGetExportCmd := bufio.NewScanner(stdout)
				for outputScannerGetExportCmd.Scan() {
					fmt.Println(outputScannerGetExportCmd.Text())
				}
				if err := outputScannerGetExportCmd.Err(); err != nil {
					utils.HandleErrorAndContinue("Error reading output from standard output. (get-export command)", err)
				}
			}

			err = getExportCmd.Wait()
			if err != nil {
				os.RemoveAll(tmpDir)
				utils.HandleErrorAndExit("Could not complete get-export command successfully.", err)
			}
			break
		}
	}
	if !found {
		os.RemoveAll(tmpDir)
		utils.HandleErrorAndExit("Unable to find an API with the name "+apiName, err)
		return err
	}
	return nil
}

// loadDefaultAWSDocFromDisk loads document.yaml stored in box/init/document.yaml
func loadDefaultAWSDoc() (*v2.Document, error) {
	docData, ok := box.Get(utils.InitDirName + utils.DefaultAWSDocFileName)
	if !ok {
		return nil, errors.New("error while retrieving " + utils.DefaultAWSDocFileName)
	}
	awsDoc := &v2.Document{}
	err = yaml.Unmarshal(docData, &awsDoc)
	if err != nil {
		return nil, err
	}
	return awsDoc, nil
}

func createAWSDocDirectory(docName string) error {
	awsDocDirectoryPath := initCmdOutputDir + string(os.PathSeparator) + utils.InitProjectDocs
	dirPath := filepath.Join(awsDocDirectoryPath, filepath.FromSlash(docName))
	utils.Logln(utils.LogPrefixInfo + "Creating directory " + dirPath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//write document.yaml file
func writeDocumentFile(docName string, summary string) error {
	document, err := loadDefaultAWSDoc()
	docData := &document.Data
	docData.Name = docName
	docData.Summary = summary
	docDataByte, err := yaml2.Marshal(document)
	if err != nil {
		return err
	}
	apiDocFilePath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.InitProjectDocs+string(os.PathSeparator)+
		docName+string(os.PathSeparator)+utils.DefaultAWSDocFileName))

	utils.Logln(utils.LogPrefixInfo + "Writing " + apiDocFilePath)
	err = ioutil.WriteFile(apiDocFilePath, docDataByte, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// write AWS API security documents based on APIs security schemes
func writeAWSSecurityDocs(oas3ByteValue []byte) error {
	securitySchemes := &v2.SecuritySchemes{}
	json.Unmarshal(oas3ByteValue, &securitySchemes)
	schemes := securitySchemes.Components.SecuritySchemes
	if securitySchemes.ResourcePolicy.Version != "" {
		err = createAWSDocDirectory(utils.ResourcePolicyDocDisplayName)
		resourcePolicyDocPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.InitProjectDocs+string(os.PathSeparator)+
			utils.ResourcePolicyDocDisplayName+string(os.PathSeparator)+utils.ResourcePolicyDocDisplayName))

		utils.Logln(utils.LogPrefixInfo + "Writing " + resourcePolicyDocPath)
		resourcePolicyDoc, _ := box.Get(utils.InitDirName + utils.ResourcePolicyDocName)
		err = ioutil.WriteFile(resourcePolicyDocPath, resourcePolicyDoc, os.ModePerm)
		if err != nil {
			return err
		}
		err = writeDocumentFile(utils.ResourcePolicyDocDisplayName, utils.ResourcePolicyDocSummary)
		if err != nil {
			return err
		}
	}
	if schemes.CognitoAuthorizer.AuthType == "cognito_user_pools" {
		err = createAWSDocDirectory(utils.CognitoDocDisplayName)
		cognitoUpDocPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.InitProjectDocs+string(os.PathSeparator)+
			utils.CognitoDocDisplayName+string(os.PathSeparator)+utils.CognitoDocDisplayName))

		utils.Logln(utils.LogPrefixInfo + "Writing " + cognitoUpDocPath)
		cognitoUpDoc, _ := box.Get(utils.InitDirName + utils.CognitoUserPoolDocName)
		err = ioutil.WriteFile(cognitoUpDocPath, cognitoUpDoc, os.ModePerm)
		if err != nil {
			return err
		}
		err = writeDocumentFile(utils.CognitoDocDisplayName, utils.CognitoDocSummary)
		if err != nil {
			return err
		}
	}
	if schemes.APIKey.Type == "apiKey" {
		err = createAWSDocDirectory(utils.ApiKeysDocDisplayName)
		apiKeyDocPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.InitProjectDocs+string(os.PathSeparator)+
			utils.ApiKeysDocDisplayName+string(os.PathSeparator)+utils.ApiKeysDocDisplayName))

		utils.Logln(utils.LogPrefixInfo + "Writing " + apiKeyDocPath)
		apiKeyDoc, _ := box.Get(utils.InitDirName + utils.AWSAPIKeyDocName)
		err = ioutil.WriteFile(apiKeyDocPath, apiKeyDoc, os.ModePerm)
		if err != nil {
			return err
		}
		err = writeDocumentFile(utils.ApiKeysDocDisplayName, utils.ApiKeysDocSummary)
		if err != nil {
			return err
		}
	}
	if schemes.Sigv4.AuthType == "awsSigv4" {
		err = createAWSDocDirectory(utils.AWSSigV4DocDisplayName)
		awsSigv4DocPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.InitProjectDocs+string(os.PathSeparator)+
			utils.AWSSigV4DocDisplayName+string(os.PathSeparator)+utils.AWSSigV4DocDisplayName))

		utils.Logln(utils.LogPrefixInfo + "Writing " + awsSigv4DocPath)
		awsSigv4Doc, _ := box.Get(utils.InitDirName + utils.AWSSigV4DocName)
		err = ioutil.WriteFile(awsSigv4DocPath, awsSigv4Doc, os.ModePerm)
		if err != nil {
			return err
		}
		err = writeDocumentFile(utils.AWSSigV4DocDisplayName, utils.AWSSigV4DocSummary)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadAPISpec loads the API definition from project folder
func loadAPISpec() (*v2.APIDefinitionFile, error) {
	pathToAPIDef := initCmdOutputDir + string(os.PathSeparator) + utils.APIDefinitionFileYaml
	apiDef, err := os.Open(pathToAPIDef)
	if err != nil {
		utils.HandleErrorAndExit("Error opening API definition from the project folder", err)
	}
	byteValue, _ := ioutil.ReadAll(apiDef)
	apiDefFile := &v2.APIDefinitionFile{}
	err = yaml.Unmarshal(byteValue, &apiDefFile)
	if err != nil {
		return nil, err
	}
	return apiDefFile, err
}

func loadAPIMetaFile() (*utils.MetaData, error) {
	pathToAPIMetaFile := initCmdOutputDir + string(os.PathSeparator) + utils.MetaFileAPI
	apiMetaFile, err := os.Open(pathToAPIMetaFile)
	if err != nil {
		utils.HandleErrorAndExit("Error opening api_meta.yaml in project folder", err)
	}
	byteValue, _ := ioutil.ReadAll(apiMetaFile)
	apiMetaData := &utils.MetaData{}
	err = yaml.Unmarshal(byteValue, &apiMetaData)
	if err != nil {
		return nil, err
	}
	return apiMetaData, err
}

func initializeProject() error {
	initCmdInitialState := "CREATED"
	initCmdApiDefinitionPath := ""
	advertiseOnly := true
	err := impl.InitAPIProject(initCmdOutputDir, initCmdInitialState, path, initCmdApiDefinitionPath, advertiseOnly)
	if err != nil {
		utils.HandleErrorAndContinue("Error initializing project", err)
		// Remove the already created project with its content since it is partially created and wrong
		dir, err := filepath.Abs(initCmdOutputDir)
		if err != nil {
			utils.HandleErrorAndExit("Error retrieving file path of the project", err)
		}
		fmt.Println("Removing the project directory " + dir + " with its content")
		err = os.RemoveAll(dir)
		if err != nil {
			utils.HandleErrorAndExit("Error removing project directory", err)
		}
	}
	apiDefFile, err := loadAPISpec()
	if err != nil {
		utils.HandleErrorAndExit("Error loading API definition from project forder", err)
	}

	def := &apiDefFile.Data

	v2.AddAwsTag(def)
	version := def.Version
	versionLength := len(version)
	if versionLength > 10 {
		trimmedVersion := version[:10]
		def.Version = trimmedVersion
	} else {
		def.Version = flagStageName
		fmt.Println("[WARN]: Unknown API version. Stage name was assigned as the API version")
	}
	def.Context = string(os.PathSeparator) + flagApiNameToGet

	oas3ByteValue := v2.CreateEpConfigForAwsAPIs(def, path)
	def.AdvertiseInformation.Advertised = true
	def.AdvertiseInformation.Vendor = "AWS"
	err = writeAWSSecurityDocs(oas3ByteValue)
	if err != nil {
		return err
	}

	apiMetaData, err := loadAPIMetaFile()
	if err != nil {
		utils.HandleErrorAndExit("Error loading api_meta.yaml from project folder", err)
	}
	apiMetaData.Version = def.Version

	newAPIMetaData, err := yaml2.Marshal(apiMetaData)
	if err != nil {
		return err
	}
	apiData, err := yaml2.Marshal(apiDefFile)
	if err != nil {
		return err
	}

	//overriding api.yaml file for AWS APIs with AWS API specific details
	apiJSONPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.APIDefinitionFileYaml))
	utils.Logln(utils.LogPrefixInfo + "Overriding " + apiJSONPath)
	err = ioutil.WriteFile(apiJSONPath, apiData, os.ModePerm)
	if err != nil {
		return err
	}
	//overriding api_meta.yaml file for AWS APIs with AWS API specific details
	apiMetaDataPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.MetaFileAPI))
	utils.Logln(utils.LogPrefixInfo + "Overriding " + apiMetaDataPath)
	err = ioutil.WriteFile(apiMetaDataPath, newAPIMetaData, os.ModePerm)
	if err != nil {
		return err
	}

	return err
}

//execute the aws init command
func execute() {
	tmpDir, err = ioutil.TempDir(dir, "OAS")
	if err != nil {
		os.RemoveAll(tmpDir)
		utils.HandleErrorAndExit("Error creating temporary directory to store OAS", err)
		return
	}
	utils.Logln(utils.LogPrefixInfo + "Temporary directory created")

	err = getOAS()
	if err != nil {
		os.RemoveAll(tmpDir)
		utils.HandleErrorAndExit("Error getting OAS from AWS.", err)
	}
	err = initializeProject()
	if err != nil {
		os.RemoveAll(tmpDir)
		utils.HandleErrorAndExit("Error initializing project.", err)
	}
	defer os.RemoveAll(tmpDir)
	utils.Logln(utils.LogPrefixInfo + "Temporary directory deleted")
}

func init() {
	InitCmd.Flags().StringVarP(&flagApiNameToGet, "name", "n", "", "Name of the API to get from AWS Api Gateway")
	InitCmd.Flags().StringVarP(&flagStageName, "stage", "s", "", "Stage name of the API to get from AWS Api Gateway")
	InitCmd.Flags().BoolVarP(&awsInitCmdForced, "force", "f", false, "Force create project")

	InitCmd.MarkFlagRequired("name")
	InitCmd.MarkFlagRequired("stage")
}
