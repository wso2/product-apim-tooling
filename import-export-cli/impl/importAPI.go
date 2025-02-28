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

package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/wso2/product-apim-tooling/import-export-cli/specs/params"

	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/Jeffail/gabs"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// JSONData represents the JSON data
type JSONData struct {
	ComplianceCheck ComplianceCheck `json:"compliance-check"`
}

// ComplianceCheck represents the compliance check
type ComplianceCheck struct {
	Result     string      `json:"result"`
	Violations []Violation `json:"violations"`
}

var (
	reAPIName = regexp.MustCompile(`[~!@#;:%^*()+={}|\\<>"',&/$]`)
)

// extractAPIDefinition extracts API information from jsonContent
func extractAPIDefinition(jsonContent []byte) (*v2.APIDefinitionFile, error) {
	api := &v2.APIDefinitionFile{}
	err := json.Unmarshal(jsonContent, &api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

// resolveImportFilePath resolves the archive/directory for import
// First will resolve in given path, if not found will try to load from exported directory
func resolveImportFilePath(file, defaultExportDirectory string) (string, error) {
	// check current path
	utils.Logln(utils.LogPrefixInfo + "Resolving for API path...")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// if the file not in given path it might be inside exported directory
		utils.Logln(utils.LogPrefixInfo+"Looking for API in", defaultExportDirectory)
		file = filepath.Join(defaultExportDirectory, file)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return "", err
		}
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// resolveYamlOrJSON for a given filepath.
// first it will look for the yaml file, if not will fallback for json
// give filename without extension so resolver will resolve for file
// fn is resolved filename, jsonContent is file as a json object, error if anything wrong happen(or both files does not exists)
func resolveYamlOrJSON(filename string) (string, []byte, error) {
	// lookup for yaml
	yamlFp := filename + ".yaml"
	if info, err := os.Stat(yamlFp); err == nil && !info.IsDir() {
		utils.Logln(utils.LogPrefixInfo+"Loading", yamlFp)
		// read it
		fn := yamlFp
		yamlContent, err := ioutil.ReadFile(fn)
		if err != nil {
			return "", nil, err
		}
		// load it as yaml
		jsonContent, err := utils.YamlToJson(yamlContent)
		if err != nil {
			return "", nil, err
		}
		return fn, jsonContent, nil
	}

	jsonFp := filename + ".json"
	if info, err := os.Stat(jsonFp); err == nil && !info.IsDir() {
		utils.Logln(utils.LogPrefixInfo+"Loading", jsonFp)
		// read it
		fn := jsonFp
		jsonContent, err := ioutil.ReadFile(fn)
		if err != nil {
			return "", nil, err
		}
		return fn, jsonContent, nil
	}

	return "", nil, fmt.Errorf("%s was not found as a YAML or JSON", filename)
}

// Substitutes environment variables in the project files.
func replaceEnvVariables(apiFilePath string) error {
	for _, replacePath := range utils.EnvReplaceFilePaths {
		absFile := filepath.Join(apiFilePath, replacePath)
		// check if the path exists. If exists, proceed with processing. Otherwise, continue with the next items
		if fi, err := os.Stat(absFile); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				utils.Logln(utils.LogPrefixInfo+"Substituting env variables of files in folder path: ", absFile)
				if strings.EqualFold(replacePath, utils.InitProjectSequences) {
					err = utils.EnvSubstituteInFolder(absFile, utils.EnvReplacePoliciesFileExtensions)
				} else {
					err = utils.EnvSubstituteInFolder(absFile, nil)
				}
			case mode.IsRegular():
				utils.Logln(utils.LogPrefixInfo+"Substituting env of file: ", absFile)
				err = utils.EnvSubstituteInFile(absFile, nil)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// importAPI imports an API to the API manager
func importAPI(endpoint, filePath, accessToken string, extraParams map[string]string, isOauth bool, dryRun bool,
	apiLoggingCmdFormat string) error {
	resp, err := ExecuteNewFileUploadRequest(endpoint, extraParams, "file",
		filePath, accessToken, isOauth)
	utils.Logf("Response : %v", resp)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}
	if dryRun {
		if resp.StatusCode() == http.StatusOK && resp.String() != "" {
			// 200 OK
			var data JSONData
			err := json.Unmarshal([]byte(resp.String()), &data)
			if err != nil {
				utils.Logln(utils.LogPrefixError, err)
				fmt.Println("Error occurred while validating API")
				return errors.New(resp.Status())
			}
			if data.ComplianceCheck.Result == "fail" {
				PrintViolations(data.ComplianceCheck.Violations, apiLoggingCmdFormat)
			} else if resp.StatusCode() == http.StatusOK {
				fmt.Printf("No violations found for the API")
			}
		} else {
			// We have an HTTP error
			utils.Logln(utils.LogPrefixError, err)
			fmt.Println("Error occurred while validating API")
			return errors.New(resp.Status())
		}
	} else {
		if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
			// 201 Created or 200 OK
			fmt.Println("Successfully imported API.")
		} else {
			// We have an HTTP error
			utils.Logln(utils.LogPrefixError, err)
			fmt.Println("Status: " + resp.Status())
			fmt.Println("Response:", resp)
			return errors.New(resp.Status())
		}
	}
	return nil
}

// ImportAPIToEnv function is used with import-api command
func ImportAPIToEnv(accessOAuthToken, importEnvironment, importPath, apiParamsPath string, importAPIUpdate,
	preserveProvider, importAPISkipCleanup, importAPIRotateRevision, importAPISkipDeployments bool, dryRun bool,
	apiLoggingCmdFormat string) error {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	return ImportAPI(accessOAuthToken, publisherEndpoint, importEnvironment, importPath, apiParamsPath, importAPIUpdate,
		preserveProvider, importAPISkipCleanup, importAPIRotateRevision, importAPISkipDeployments, dryRun, apiLoggingCmdFormat)
}

// ImportAPI function is used with import-api command
func ImportAPI(accessOAuthToken, publisherEndpoint, importEnvironment, importPath, apiParamsPath string, importAPIUpdate,
	preserveProvider, importAPISkipCleanup, importAPIRotateRevision, importAPISkipDeployments bool,
	dryRun bool, apiLoggingCmdFormat string) error {
	exportDirectory := filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)
	resolvedAPIFilePath, err := resolveImportFilePath(importPath, exportDirectory)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"API Location:", resolvedAPIFilePath)

	utils.Logln(utils.LogPrefixInfo + "Creating workspace")
	tmpPath, err := utils.GetTempCloneFromDirOrZip(resolvedAPIFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if importAPISkipCleanup {
			utils.Logln(utils.LogPrefixInfo+"Leaving", tmpPath)
			return
		}
		utils.Logln(utils.LogPrefixInfo+"Deleting", tmpPath)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			utils.Logln(utils.LogPrefixError + err.Error())
		}
	}()
	apiFilePath := tmpPath

	utils.Logln(utils.LogPrefixInfo + "Substituting environment variables in API files...")
	err = replaceEnvVariables(apiFilePath)
	if err != nil {
		return err
	}

	if importAPISkipDeployments {
		//If skip deployments flag used, deployment_environments files will be removed from import artifacts
		loc := filepath.Join(apiFilePath, utils.DeploymentEnvFile)
		utils.Logln(utils.LogPrefixInfo + "Removing the deployment environments file from " + loc)
		err := utils.RemoveFileIfExists(loc)
		if err != nil {
			return err
		}
	}

	if apiParamsPath != "" {
		//Reading params file of the API and add configurations into temp artifact
		err := handleCustomizedParameters(apiFilePath, apiParamsPath, importEnvironment)
		if err != nil {
			return err
		}
	}

	// if apiFilePath contains a directory, zip it. Otherwise, leave it as it is.
	apiFilePath, err, cleanupFunc := utils.CreateZipFileFromProject(apiFilePath, importAPISkipCleanup)
	if err != nil {
		return err
	}

	//cleanup the temporary artifacts once consuming the zip file
	if cleanupFunc != nil {
		defer cleanupFunc()
	}

	extraParams := map[string]string{}
	publisherEndpoint += "/apis/import"
	if importAPIUpdate {
		publisherEndpoint += "?overwrite=" + strconv.FormatBool(true) + "&preserveProvider=" +
			strconv.FormatBool(preserveProvider) + "&rotateRevision=" + strconv.FormatBool(importAPIRotateRevision)
	} else {
		publisherEndpoint += "?preserveProvider=" + strconv.FormatBool(preserveProvider) + "&rotateRevision=" +
			strconv.FormatBool(importAPIRotateRevision)
	}

	if dryRun {
		publisherEndpoint += "&dryRun=" + strconv.FormatBool(true)
	}
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + publisherEndpoint)

	err = importAPI(publisherEndpoint, apiFilePath, accessOAuthToken, extraParams, true, dryRun, apiLoggingCmdFormat)
	return err
}

// envParamsFileProcess function is used to process the environment parameters when they are provided as a file
func envParamsFileProcess(importPath, paramsPath, importEnvironment string) error {
	apiParams, err := params.LoadApiParamsFromFile(paramsPath)
	if err != nil {
		return err
	}
	// check whether import environment is included in params configuration
	envParams := apiParams.GetEnv(importEnvironment)
	if envParams == nil {
		return errors.New("Environment '" + importEnvironment + "' does not exist in " + paramsPath)
	} else {

		// Create a source directory and add source content to it and then zip it
		sourceFilePath := filepath.Join(importPath, "SourceArchive")
		err = utils.MoveDirectoryContentsToNewDirectory(importPath, sourceFilePath)
		if err != nil {
			return err
		}

		err, cleanupFunc := utils.CreateZipFile(sourceFilePath, false)
		if err != nil {
			return err
		}
		//cleanup the temporary artifacts once consuming the zip file
		if cleanupFunc != nil {
			defer cleanupFunc()
		}
		//If environment parameters are present in parameter file
		err = handleEnvParams(importPath, importPath, envParams)
		if err != nil {
			return err
		}
	}
	return nil
}

// envParamsDirectoryProcess function is used to process the environment parameters when they are provided as a
// directory
func envParamsDirectoryProcess(importPath, paramsPath, importEnvironment string) error {
	apiParams, err := params.LoadApiParamsFromDirectory(paramsPath)
	if err != nil {
		return err
	}
	// check whether import environment is included in api params configuration
	envParams := apiParams.GetEnv(importEnvironment)
	if envParams == nil {
		return errors.New("Environment '" + importEnvironment + "' does not exist in " + paramsPath)
	} else {

		// Create a source directory and add source content to it and then zip it
		sourceFilePath := filepath.Join(importPath, "SourceArchive")
		err = utils.MoveDirectoryContentsToNewDirectory(importPath, sourceFilePath)
		if err != nil {
			return err
		}

		err, cleanupFunc := utils.CreateZipFile(sourceFilePath, false)
		if err != nil {
			return err
		}
		//cleanup the temporary artifacts once consuming the zip file
		if cleanupFunc != nil {
			defer cleanupFunc()
		}

		//create new directory for deployment configurations
		deploymentDirectoryPath := filepath.Join(importPath, "Deployment")
		err = utils.CreateDirIfNotExist(deploymentDirectoryPath)
		if err != nil {
			return err
		}

		//copy all the content in the params directory into the artifact to be imported
		err = utils.CopyDirectoryContents(paramsPath, deploymentDirectoryPath)
		if err != nil {
			return err
		}
		//If environment parameters are present in parameter file inside the deployment params directory
		err = handleEnvParams(importPath, deploymentDirectoryPath, envParams)
		if err != nil {
			return err
		}
	}
	return nil
}

// handleCustomizedParameters handles the configurations provided with params file of the API and the resources that needs to
// transfer to server side will bundle with the artifact to be imported.
func handleCustomizedParameters(importPath, paramsPath, importEnvironment string) error {
	utils.Logln(utils.LogPrefixInfo+"Loading parameters from", paramsPath)
	if strings.Contains(paramsPath, ".yaml") {
		utils.Logln(utils.LogPrefixInfo+"Processing Params file", paramsPath)
		err := envParamsFileProcess(importPath, paramsPath, importEnvironment)
		if err != nil {
			return err
		}
	} else {
		utils.Logln(utils.LogPrefixInfo+"Processing Params in the deployment directory", paramsPath)
		err := envParamsDirectoryProcess(importPath, paramsPath, importEnvironment)
		if err != nil {
			return err
		}

	}
	return nil
}

// Process env params and create the intermediate_params.yaml file to pass to the server
func handleEnvParams(tempDirectory string, destDirectory string, environmentParams *params.Environment) error {
	// read api params from external parameters file
	if len(environmentParams.Config) == 0 {
		return errors.New("configs value is empty in the provided parameters")
	}

	envParamsJson, err := jsoniter.Marshal(environmentParams.Config)
	if err != nil {
		return err
	}

	var apiParamsPath string
	apiParams, err := gabs.ParseJSON(envParamsJson)
	paramsContent, err := utils.JsonToYaml(apiParams.Bytes())
	if err != nil {
		return err
	}

	//over-write the api_params.file with latest configurations
	apiParamsPath = filepath.Join(destDirectory, utils.ParamsIntermediateFile)
	utils.Logln(utils.LogPrefixInfo+"Adding the Params file into", apiParamsPath)
	err = ioutil.WriteFile(apiParamsPath, paramsContent, 0644)
	if err != nil {
		return err
	}
	return nil
}
