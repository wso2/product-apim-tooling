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

	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	reAPIProductName = regexp.MustCompile(`[~!@#;:%^*()+={}|\\<>"',&/$]`)
)

// extractAPIProductDefinition extracts API Product information from jsonContent
func extractAPIProductDefinition(jsonContent []byte) (*v2.APIProductDefinition, error) {
	apiProduct := &v2.APIProductDefinition{}
	err := json.Unmarshal(jsonContent, &apiProduct)
	if err != nil {
		return nil, err
	}

	return apiProduct, nil
}

// resolveImportAPIProductFilePath resolves the archive/directory for import
// First will resolve in given path, if not found will try to load from exported directory
func resolveImportAPIProductFilePath(file, defaultExportDirectory string) (string, error) {
	// Check current path
	utils.Logln(utils.LogPrefixInfo + "Resolving for API Product path...")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// If the file not in given path it might be inside exported directory
		utils.Logln(utils.LogPrefixInfo+"Looking for API Product in", defaultExportDirectory)
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

// importAPIProduct imports an API Product to the API manager
func importAPIProduct(endpoint, filePath, accessToken string, extraParams map[string]string) error {
	resp, err := ExecuteNewFileUploadRequest(endpoint, extraParams, "file",
		filePath, accessToken, true)
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		// 201 Created or 200 OK
		fmt.Println("Successfully imported API Product.")
		return nil
	} else {
		// We have an HTTP error
		fmt.Println("Error importing API Product.")
		fmt.Println("Status: " + resp.Status())
		fmt.Println("Response:", resp)
		return errors.New(resp.Status())
	}
}

// preProcessDependentAPIs pre processes dependent APIs
func preProcessDependentAPIs(apiProductFilePath, importEnvironment string, importAPIProductPreserveProvider bool) error {
	// Check whether the APIs directory exists
	apisDirectoryPath := apiProductFilePath + string(os.PathSeparator) + "APIs"
	_, err := os.Stat(apisDirectoryPath)
	if os.IsNotExist(err) {
		utils.Logln(utils.LogPrefixInfo + "APIs directory does not exists. Ignoring APIs.")
		return nil
	}

	// If APIs directory exists, read the directory
	items, _ := ioutil.ReadDir(apisDirectoryPath)
	// Iterate through the API directories available
	for _, item := range items {
		apiDirectoryPath := apisDirectoryPath + string(os.PathSeparator) + item.Name()

		// Substitutes environment variables in the project files
		err = replaceEnvVariables(apiDirectoryPath)
		if err != nil {
			return err
		}

		utils.Logln(utils.LogPrefixInfo + "Attempting to inject parameters to the API from api_params.yaml (if exists)")
		paramsPath := apiDirectoryPath + string(os.PathSeparator) + utils.ParamFileAPI
		// Check whether api_params.yaml file is available inside the particular API directory
		if utils.IsFileExist(paramsPath) {
			// Reading API params file and populate api.yaml
			err := handleCustomizedParameters(apiDirectoryPath, paramsPath, importEnvironment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ImportAPIProductToEnv function is used with import-api-product command
func ImportAPIProductToEnv(accessOAuthToken, importEnvironment, importPath string, importAPIs, importAPIsUpdate,
	importAPIProductUpdate, importAPIProductPreserveProvider, importAPIProductSkipCleanup, rotateRevision bool) error {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	return ImportAPIProduct(accessOAuthToken, publisherEndpoint, importEnvironment, importPath, importAPIs, importAPIsUpdate,
		importAPIProductUpdate, importAPIProductPreserveProvider, importAPIProductSkipCleanup, rotateRevision)
}

// ImportAPIProduct function is used with import-api-product command
func ImportAPIProduct(accessOAuthToken, publisherEndpoint, importEnvironment, importPath string, importAPIs, importAPIsUpdate,
	importAPIProductUpdate, importAPIProductPreserveProvider, importAPIProductSkipCleanup, rotateRevision bool) error {
	var exportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApiProductsDirName)

	resolvedAPIProductFilePath, err := resolveImportAPIProductFilePath(importPath, exportDirectory)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"API Product Location:", resolvedAPIProductFilePath)

	utils.Logln(utils.LogPrefixInfo + "Creating workspace")
	tmpPath, err := utils.GetTempCloneFromDirOrZip(resolvedAPIProductFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if importAPIProductSkipCleanup {
			utils.Logln(utils.LogPrefixInfo+"Leaving", tmpPath)
			return
		}
		utils.Logln(utils.LogPrefixInfo+"Deleting", tmpPath)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			utils.Logln(utils.LogPrefixError + err.Error())
		}
	}()
	apiProductFilePath := tmpPath

	// Pre Process dependent APIs
	err = preProcessDependentAPIs(apiProductFilePath, importEnvironment, importAPIProductPreserveProvider)
	if err != nil {
		return err
	}

	utils.Logln(utils.LogPrefixInfo + "Substituting environment variables in API Product files...")
	err = replaceEnvVariables(apiProductFilePath)
	if err != nil {
		return err
	}

	// If apiProductFilePath contains a directory, zip it. Otherwise, leave it as it is.
	apiProductFilePath, err, cleanupFunc := utils.CreateZipFileFromProject(apiProductFilePath, importAPIProductSkipCleanup)
	if err != nil {
		return err
	}

	//cleanup the temporary artifacts once consuming the zip file
	if cleanupFunc != nil {
		defer cleanupFunc()
	}

	if err != nil {
		utils.HandleErrorAndExit("Error getting OAuth Tokens", err)
	}
	extraParams := map[string]string{}
	publisherEndpoint += "/api-products/import" + "?preserveProvider=" +
		strconv.FormatBool(importAPIProductPreserveProvider)+ "&rotateRevision=" + strconv.FormatBool(rotateRevision)

	// If the user has specified import-apis flag or update-apis flag, importAPIs parameter should be passed as true
	// because update is also an import task
	if importAPIs || importAPIsUpdate {
		publisherEndpoint += "&importAPIs=" + strconv.FormatBool(true)
	}

	// If the user need to update the APIs and the API Product, overwriteAPIs parameter should be passed as true
	if importAPIsUpdate {
		publisherEndpoint += "&overwriteAPIs=" + strconv.FormatBool(true)
	}

	// If the user need only to update the API Product, overwriteAPIProduct parameter should be passed as true
	if importAPIsUpdate || importAPIProductUpdate {
		publisherEndpoint += "&overwriteAPIProduct=" + strconv.FormatBool(true)
	}

	utils.Logln(utils.LogPrefixInfo + "Import URL: " + publisherEndpoint)
	err = importAPIProduct(publisherEndpoint, apiProductFilePath, accessOAuthToken, extraParams)
	return err
}
