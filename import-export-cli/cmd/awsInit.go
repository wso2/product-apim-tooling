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
//./build.sh -t apictl.go -v 3.2.0 -f
/*TODO 
* 01 : verbose flag support 
* 02 : Don't duplicate the default api template. The apictl already has a template. Reuse the same template and set advertise only to true.
*/

package cmd 

import (
	"os"
	"fmt"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"github.com/spf13/cobra"

	//"github.com/ghodss/yaml"
	
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
 
	yaml2 "gopkg.in/yaml.v2"
	
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"
 
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	//"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

var flagApiNameToGet string		//name of the api to get from aws gateway
var flagStageName string		//api stage to get  
var dir string					//dir where the aws init command is executed from
var path string					//path of the OAS extracted from AWS
var tmpDir string				//temporary directory created to store the OAS extracted from aws till the project is initialized
var err error 

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
var exportType string = "oas30"	//openapi 3. Use "swagger" to request for a swagger 2.
var debugFlag string	//aws cli debug flag for apictl verbose mode 


const awsInitCmdLiteral = "aws init"
const awsInitCmdShortDesc = "Get the swagger of an API from AWS API Gateway"
const awsInitCmdLongDesc = `Downloading the swagger definition of an API from the AWS API Gateway`
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
		utils.Logln("hello")
		initCmdOutputDir = dir + "/" + flagApiNameToGet

		if stat, err := os.Stat(initCmdOutputDir); !os.IsNotExist(err) {
			fmt.Printf("%s already exists\n", initCmdOutputDir)
			if !stat.IsDir() {
				fmt.Printf("%s is not a directory\n", initCmdOutputDir)
				os.Exit(1)
			}
			var confirmation string 
			fmt.Println("Enter (y) if you wish to continue OR (n) to exit.")
			fmt.Scanln(&confirmation)
			 
			if confirmation == "n" {
				os.Exit(1)
			}
			if confirmation == "y" {
				execute()
			} else {
				fmt.Println ("Invalid input")
			}
		} else {
			execute()
		}
	},
}

type Apis struct {
	Items []struct {
		Id                    string `json:"id"`
		Name                  string `json:"name"`
		Description 		   string `json:"description"`
		CreatedDate           int    `json:"createdDate"`
		APIKeySource          string `json:"apiKeySource"`
		EndpointConfiguration struct {
			Types []string `json:"types"`
		} `json:"endpointConfiguration"`
	} `json:"items"`
}

func getOAS() error {
	if utils.VerboseModeEnabled() {
		utils.Logln("Executing aws version command")
	}
	//check whether aws cli is installed
	cmd_1, err := exec.Command(awsCmdLiteral, awsCLIVersionFlag).Output()
	if err != nil {
		fmt.Println("Error getting AWS CLI version. Make sure aws cli is installed and configured.")
		return err
	}
	if utils.VerboseModeEnabled() {
		fmt.Println("AWS CLI version : ")
		output := string(cmd_1[:])
		fmt.Println(output)
	}
	if utils.VerboseModeEnabled() {
		utils.Logln("Executing aws get-rest-apis command")
	}
	cmd_2, err := exec.Command(awsCmdLiteral, apiGateway, getAPI, outputFlag, outputType).Output()
	
	//Unmarshal from JSON into Apis struct.
	apis := Apis{}
	err = json.Unmarshal([]byte(cmd_2), &apis)
	if err != nil {
		return err
	}
	var found bool
	apiName := flagApiNameToGet
	stageName := flagStageName
	path = tmpDir + "/" + apiName + ".json"
	// Searching for API ID:
	if utils.VerboseModeEnabled() {
		fmt.Println("Searching for API ID...")
	}
	for _, item := range apis.Items {
		if item.Name == apiName {
			fmt.Println("API ID found : ", item.Id)
			api_id := item.Id 
			if utils.VerboseModeEnabled() {
				utils.Logln("Executing aws get-export command")
				debugFlag = "--debug"	//activating the aws cli debug flag in apictl verbose mode 
			}
			cmd_3, err := exec.Command(awsCmdLiteral, apiGateway, getExport, apiIdFlag, api_id, stageNameFlag, stageName, exportTypeFlag, exportType, path, outputFlag, outputType, debugFlag)
			
			if err != nil {
				return err
			}
			if utils.VerboseModeEnabled() {
				output := string(cmd_3[:])
				fmt.Println(output)
			}
			found = true 
			break
		}
	}
	if !found {
		fmt.Println("Unable to find an API with the name " + apiName)
		return err
	}
	return nil
}

func initializeProject() error {
	initCmdOutputDir = flagApiNameToGet 
	swaggerSavePath := filepath.Join(initCmdOutputDir, filepath.FromSlash("Meta-information/swagger.yaml"))
	fmt.Println("Initializing a new WSO2 API Manager project in", dir)
	
	def, err := loadDefaultSpecFromDisk()
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
	err = v2.Swagger2Populate(def, doc)
	if err != nil {
		return err
	}

	err = v2.GetServerUrlFromOAS(def, path)
	if err != nil {
		return err
	}

	v2.SetAdvertiseOnlyToTrue(def)

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

	apiData, err := yaml2.Marshal(def)
	if err != nil {
		return err
	}

	// write to the disk
	apiJSONPath := filepath.Join(initCmdOutputDir, filepath.FromSlash("Meta-information/api.yaml"))
	utils.Logln(utils.LogPrefixInfo + "Writing " + apiJSONPath)
	err = ioutil.WriteFile(apiJSONPath, apiData, os.ModePerm)
	if err != nil {
		return err
	}

	apimProjParamsFilePath := filepath.Join(initCmdOutputDir, utils.ParamFileAPI)
	utils.Logln(utils.LogPrefixInfo + "Writing " + apimProjParamsFilePath)
	err = impl.ScaffoldParams(apimProjParamsFilePath)
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

	fmt.Println("Project initialized")
	fmt.Println("Open README file to learn more")
	return nil
} 

//execute the aws init command 
func execute() {
	tmpDir, err = ioutil.TempDir(dir, "OAS")
		if err != nil {
			fmt.Println("Error creating temporary directory to store OAS")
			return
		}
	if utils.VerboseModeEnabled() {
		fmt.Println("Temporary directory created")
	}
	err = getOAS()
	if err != nil {
		utils.HandleErrorAndExit("Error getting OAS from AWS. ", err)
	}
	err = initializeProject()
	if err != nil {
		utils.HandleErrorAndExit("Error initializing project. ", err)
	}
	defer os.RemoveAll(tmpDir)
	if utils.VerboseModeEnabled() {
		fmt.Println("Temporary directory deleted")
	}
}

func init() { 
		RootCmd.AddCommand(AwsInitCmd)
		RootCmd.AddCommand(DeleteCmd)
		AwsInitCmd.Flags().StringVarP(&flagApiNameToGet, "name", "n", "", "Name of the API to get from AWS Api Gateway")
		AwsInitCmd.Flags().StringVarP(&flagStageName, "stage", "s", "", "Stage name of the API to get from AWS Api Gateway")

		AwsInitCmd.MarkFlagRequired("name")
		AwsInitCmd.MarkFlagRequired("stage")
}