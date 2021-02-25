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

package cmd 

import (
	"os"
	"fmt"
	"strconv"
	"os/exec"
	"bufio"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"github.com/spf13/cobra"
	"github.com/Jeffail/gabs"
	jsoniter "github.com/json-iterator/go"
	
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
 
	yaml2 "gopkg.in/yaml.v2"
	
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"
 
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagApiNameToGet string		//name of the api to get from aws gateway
var flagStageName string		//api stage to get  
var dir string					//dir where the aws init command is executed from
var path string					//path of the OAS extracted from AWS
var tmpDir string				//temporary directory created to store the OAS extracted from aws till the project is initialized
var cmd2_output string 
var err error 
var awsInitCmdForced bool

//common aws cmd flags 
var awsCmdLiteral string = "aws"
var apiGateway string = "apigateway"

//aws cmd Output type
var outputFlag string = "--output"
var outputType string = "json"

//cmd_1 
var awsCLIVersionFlag string = "--version"

//cmd_2
var getAPI string = "get-rest-apis"

//cmd_3
var getExport string = "get-export"
var apiIdFlag string = "--rest-api-id"
var stageNameFlag string = "--stage-name"
var exportTypeFlag string = "--export-type"
var exportType string = "oas30"	//default export type is openapi3. Use "swagger" to request for a swagger 2.
var debugFlag string			//aws cli debug flag for apictl verbose mode 

const awsInitCmdLiteral = "aws init"
const awsInitCmdShortDesc = "Initialize a API project from a AWS API"
const awsInitCmdLongDesc = `Downloading the OpenAPI specification of an API from the AWS API Gateway to initialize a WSO2 API project`
const awsInitCmdExamples = utils.ProjectName + ` ` + awsInitCmdLiteral  + ` -n PetStore -s Demo

` + utils.ProjectName + ` ` +  awsInitCmdLiteral + ` --name PetStore --stage Demo

` + utils.ProjectName + ` ` +  awsInitCmdLiteral + ` --name Shopping --stage Live

NOTE: Both flags --name (-n) and --stage (-s) are mandatory as both values are needed to get the openAPI from AWS API Gateway.
Make sure the API name and Stage Name are correct.
Also make sure you have AWS CLI installed and configured before executing the aws init command.
Vist https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html for more info`

func getPath() {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	dir = pwd
}

//aws init Cmd
var AwsInitCmd = &cobra.Command{
	Use:   	awsInitCmdLiteral,
	Short: 	awsInitCmdShortDesc,
	Long:  	awsInitCmdLongDesc,
	Example: 	awsInitCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		getPath()
		initCmdOutputDir = dir + "/" + flagApiNameToGet

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
		Id                    string `json:"id"`
		Name                  string `json:"name"`
	} `json:"items"`
}

func getOAS() error {
	utils.Logln(utils.LogPrefixInfo + "Executing aws version command")
	//check whether aws cli is installed
	cmd_1, err := exec.Command(awsCmdLiteral, awsCLIVersionFlag).Output()
	if err != nil {
		fmt.Println("Error getting AWS CLI version. Make sure aws cli is installed and configured.")
		return err
	}
	output := string(cmd_1[:])
	utils.Logln(utils.LogPrefixInfo + "AWS CLI version :  " + output)

	if utils.VerboseModeEnabled() {
		debugFlag = "--debug"	//activating the aws cli debug flag in apictl verbose mode 
	}
	utils.Logln(utils.LogPrefixInfo + "Executing aws get-rest-apis command in debug mode")
	cmd_2 := exec.Command(awsCmdLiteral, apiGateway, getAPI, outputFlag, outputType, debugFlag)
	stderr, err := cmd_2.StderrPipe()
	if err != nil {
		fmt.Println("Error creating pipe to standard error. (get-rest-apis command)", err)
	}
	stdout, err := cmd_2.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating pipe to standard output (get-rest-apis command).", err)
	}

	err = cmd_2.Start()
	if err != nil {
		fmt.Println("Error starting get-rest-apis command.", err)
	}

	if utils.VerboseModeEnabled() {
		logsScannerCmd2 := bufio.NewScanner(stderr)
		for logsScannerCmd2.Scan() {
			fmt.Println(logsScannerCmd2.Text())
		}

		if err := logsScannerCmd2.Err(); err != nil {
			fmt.Println("Error reading debug logs from standard error. (get-rest-apis command)", err)
		}
	}

	outputScannerCmd2 := bufio.NewScanner(stdout)
	for outputScannerCmd2.Scan() {
		cmd2_output = cmd2_output + outputScannerCmd2.Text()
	}

	if err := outputScannerCmd2.Err(); err != nil {
		fmt.Println("Error reading output from standard output.", err)
	}
	//
	err = cmd_2.Wait()
	if err != nil {
		fmt.Println("Could not complete get-rest-apis command successfully.", err)
	}

	//Unmarshal from JSON into Apis struct.
	apis := Apis{}
	err = json.Unmarshal([]byte(cmd2_output), &apis)
	if err != nil {
		return err
	}
	extractedAPIs := strconv.Itoa(len(apis.Items))
	utils.Logln(utils.LogPrefixInfo + extractedAPIs + " APIs were extracted")

	var found bool
	apiName := flagApiNameToGet
	stageName := flagStageName
	path = tmpDir + "/" + apiName + ".json"
	// Searching for API ID:
	utils.Logln(utils.LogPrefixInfo + "Searching for API ID...")
	for _, item := range apis.Items {
		if item.Name == apiName {
			fmt.Println("API ID found : ", item.Id)
			api_id := item.Id 
			found = true

			utils.Logln(utils.LogPrefixInfo + "Executing aws get-export command in debug mode")
			cmd_3:= exec.Command(awsCmdLiteral, apiGateway, getExport, apiIdFlag, api_id, stageNameFlag, stageName, exportTypeFlag, exportType, path, outputFlag, outputType, debugFlag)
			stderr, err := cmd_3.StderrPipe()
			if err != nil {
				fmt.Println("Error creating pipe to standard error. (get-export command)", err)
			}
			stdout, err := cmd_3.StdoutPipe()
			if err != nil {
				fmt.Println("Error creating pipe to standard output. (get-export command)", err)
			}

			err = cmd_3.Start()
			if err != nil {
				fmt.Println("Error starting get-export command.", err)
			}

			if utils.VerboseModeEnabled() {
				logsScannerCmd3 := bufio.NewScanner(stderr)
				for logsScannerCmd3.Scan() {
					fmt.Println(logsScannerCmd3.Text())
				}
				if err := logsScannerCmd3.Err(); err != nil {
					fmt.Println("Error reading debug logs from standard error. (get-export command)", err)
				}
			}
			
			if utils.VerboseModeEnabled() {
				outputScannerCmd3 := bufio.NewScanner(stdout)
				for outputScannerCmd3.Scan() {
					fmt.Println(outputScannerCmd3.Text())
				}
				if err := outputScannerCmd3.Err(); err != nil {
					fmt.Println("Error reading output from standard output. (get-export command)", err)
				}
			}

			err = cmd_3.Wait()
			if err != nil {
				fmt.Println("Could not complete get-export command successfully.", err)
			} 
			break
		}
	}
	if !found {
		fmt.Println("Unable to find an API with the name", apiName)
		os.RemoveAll(tmpDir)
		os.Exit(1)
		return err
	}
	return nil
}

func initializeProject() error {
	initCmdOutputDir = flagApiNameToGet 
	swaggerSavePath := filepath.Join(initCmdOutputDir, filepath.FromSlash("Definitions/swagger.yaml"))
	fmt.Println("Initializing a new WSO2 API Manager project in", dir)
	
	definitionFile, err := loadDefaultSpecFromDisk()
	
	// Get the API DTO specific details to process
	def := &definitionFile.Data
	if err != nil {
		return err
	}

	err = createDirectories(initCmdOutputDir)
	if err != nil {
		return err
	}

	// use swagger to auto generate
	// load swagger from tmp directory
	doc, err := loadSwagger(path)
	if err != nil {
		return err
	}

	// We use swagger2 loader. It works fine for now
	// Since we don't use 3.0 specific details its ok
	// otherwise please use v2.openAPI3 loaders
	//doc.Spec().Info.Version = doc.Spec().Info.Version[:10]

	err = v2.Swagger2Populate(def, doc)
	if err != nil {
		return err
	}

	//err = v2.GetServerUrlFromOAS(def, path)
	//if err != nil {
	//	return err
	//}

	//v2.SetAdvertiseOnlyToTrue(def)

	// convert and save swagger as yaml
	yamlSwagger, err := utils.JsonToYaml(doc.Raw())
	if err != nil {
		return err
	}

	// write to file
	err = ioutil.WriteFile(swaggerSavePath, yamlSwagger, os.ModePerm)
	if err != nil {
		return err
	}

	utils.Logln(utils.LogPrefixInfo + "Reading API Definition from " + path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	apiDef := &v2.APIDefinitionFile{}
	fmt.Println(string(content))

	// substitute env variables
	utils.Logln(utils.LogPrefixInfo + "Substituting environment variables")
	data, err := utils.EnvSubstitute(string(content))
	if err != nil {
		return err
	}
	content = []byte(data)

	// read from yaml definition
	err = yaml2.Unmarshal(content, &apiDef)
	if err != nil {
		return err
	}

	// marshal original def
	originalDefBytes, err := jsoniter.Marshal(definitionFile)
	if err != nil {
		return err
	}
	// marshal new def
	newDefBytes, err := jsoniter.Marshal(apiDef)
	if err != nil {
		return err
	}

	// merge two definitions
	finalDefBytes, err := utils.MergeJSON(originalDefBytes, newDefBytes)
	if err != nil {
		return err
	}
	tmpDef := &v2.APIDefinitionFile{}
	err = json.Unmarshal(finalDefBytes, &tmpDef)
	if err != nil {
		return err
	}
	definitionFile.Data = tmpDef.Data

	apiData, err := yaml2.Marshal(definitionFile)
	if err != nil {
		return err
	}

	// write to the disk
	apiJSONPath := filepath.Join(initCmdOutputDir, filepath.FromSlash("api.yaml"))
	utils.Logln(utils.LogPrefixInfo + "Writing " + apiJSONPath)
	err = ioutil.WriteFile(apiJSONPath, apiData, os.ModePerm)
	if err != nil {
		return err
	}

	apimProjReadmeFilePath := filepath.Join(initCmdOutputDir, "README.md")
	utils.Logln(utils.LogPrefixInfo + "Writing " + apimProjReadmeFilePath)
	readme, _ := box.Get("/init/README.md")
	err = ioutil.WriteFile(apimProjReadmeFilePath, readme, os.ModePerm)
	if err != nil {
		return err
	}

	// Create metaData struct using details from definition
	metaData := utils.MetaData{
		Name:    definitionFile.Data.Name,
		Version: definitionFile.Data.Version,
	}
	marshaledData, err := jsoniter.Marshal(metaData)
	if err != nil {
		return err
	}

	jsonMetaData, err := gabs.ParseJSON(marshaledData)
	metaDataContent, err := utils.JsonToYaml(jsonMetaData.Bytes())
	if err != nil {
		return err
	}

	// write api_meta.yaml file to the project directory
	apiMetaDataPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.MetaFileAPI))
	utils.Logln(utils.LogPrefixInfo + "Writing " + apiMetaDataPath)
	err = ioutil.WriteFile(apiMetaDataPath, metaDataContent, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("Project initialized")
	fmt.Println("Open README file to learn more")
	return nil
} 

//execute the aws init command 
func execute() {
	tmpDir, err = ioutil.TempDir(dir, "OAS")
		if err != nil {
			os.RemoveAll(tmpDir)
			fmt.Println("Error creating temporary directory to store OAS")
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
		RootCmd.AddCommand(AwsInitCmd)
		RootCmd.AddCommand(DeleteCmd)
		AwsInitCmd.Flags().StringVarP(&flagApiNameToGet, "name", "n", "", "Name of the API to get from AWS Api Gateway")
		AwsInitCmd.Flags().StringVarP(&flagStageName, "stage", "s", "", "Stage name of the API to get from AWS Api Gateway")
		AwsInitCmd.Flags().BoolVarP(&awsInitCmdForced, "force", "f", false, "Force create project")

		AwsInitCmd.MarkFlagRequired("name")
		AwsInitCmd.MarkFlagRequired("stage")
}
