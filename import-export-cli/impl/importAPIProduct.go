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
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/Jeffail/gabs"
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

func populateAPIProductWithDefaults(def *v2.APIProductDefinition) (dirty bool) {
	dirty = false
	def.Context = strings.ReplaceAll(def.Context, " ", "")
	if def.ContextTemplate == "" {
		if !strings.Contains(def.Context, "{version}") {
			def.ContextTemplate = path.Clean(def.Context + "/{version}")
			def.Context = strings.ReplaceAll(def.ContextTemplate, "{version}", def.ID.Version)
		} else {
			def.Context = path.Clean(def.Context)
			def.ContextTemplate = def.Context
			def.Context = strings.ReplaceAll(def.Context, "{version}", def.ID.Version)
		}
		dirty = true
	}
	if def.Tags == nil {
		def.Tags = []string{}
		dirty = true
	}
	return
}

// validateApiProductDefinition validates an API Product against basic rules
func validateAPIProductDefinition(def *v2.APIProductDefinition) error {
	utils.Logln(utils.LogPrefixInfo + "Validating API Product")
	if isEmpty(def.ID.APIProductName) {
		return errors.New("apiProductName is required")
	}
	if reAPIProductName.MatchString(def.ID.APIProductName) {
		return errors.New(`apiProductName contains one or more illegal characters (~!@#;:%^*()+={}|\\<>"',&\/$)`)
	}
	if isEmpty(def.ID.Version) {
		return errors.New("version is required")
	}
	if isEmpty(def.Context) {
		return errors.New("context is required")
	}
	if isEmpty(def.ContextTemplate) {
		return errors.New("contextTemplate is required")
	}
	if !strings.HasPrefix(def.Context, "/") {
		return errors.New("context should begin with a /")
	}
	if !strings.HasPrefix(def.ContextTemplate, "/") {
		return errors.New("contextTemplate should begin with a /")
	}
	return nil
}

// importAPIProduct imports an API Product to the API manager
func importAPIProduct(endpoint, filePath, accessToken string, extraParams map[string]string) error {
	resp, err := ExecuteNewFileUploadRequest(endpoint, extraParams, "file",
		filePath, accessToken)
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
	importAPIProductUpdate, importAPIProductPreserveProvider, importAPIProductSkipCleanup bool) error {
	adminEndpoint := utils.GetAdminEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	return ImportAPIProduct(accessOAuthToken, adminEndpoint, importEnvironment, importPath, importAPIs, importAPIsUpdate,
		importAPIProductUpdate, importAPIProductPreserveProvider, importAPIProductSkipCleanup)
}

// ImportAPIProduct function is used with import-api-product command
func ImportAPIProduct(accessOAuthToken, adminEndpoint, importEnvironment, importPath string, importAPIs, importAPIsUpdate,
	importAPIProductUpdate, importAPIProductPreserveProvider, importAPIProductSkipCleanup bool) error {
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

	// Get API Product info
	apiProductInfo, originalContent, err := GetAPIProductDefinition(apiProductFilePath)
	if err != nil {
		return err
	}
	// Fill with defaults
	if populateAPIProductWithDefaults(apiProductInfo) {
		utils.Logln(utils.LogPrefixInfo + "API Product is populated with defaults")
		// API Product is dirty, write it to disk
		buf, err := json.Marshal(apiProductInfo)
		if err != nil {
			return err
		}

		newContent, err := gabs.ParseJSON(buf)
		if err != nil {
			return err
		}
		originalContent, err := gabs.ParseJSON(originalContent)
		if err != nil {
			return err
		}
		result, err := utils.MergeJSON(originalContent.Bytes(), newContent.Bytes())
		if err != nil {
			return err
		}

		yamlContent, err := utils.JsonToYaml(result)
		if err != nil {
			return err
		}
		p := filepath.Join(apiProductFilePath, "Meta-information", "api.yaml")
		utils.Logln(utils.LogPrefixInfo+"Writing", p)

		err = ioutil.WriteFile(p, yamlContent, 0644)
		if err != nil {
			return err
		}
	}
	// Validate definition
	if err = validateAPIProductDefinition(apiProductInfo); err != nil {
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

	updateAPIProduct := false
	if importAPIsUpdate || importAPIProductUpdate {
		// Check for API Product existence
		id, err := GetAPIProductId(accessOAuthToken, importEnvironment, apiProductInfo.ID.APIProductName,
			apiProductInfo.ID.ProviderName)
		if err != nil {
			return err
		}

		if id == "" {
			updateAPIProduct = false
			utils.Logln("The specified API Product was not found.")
			utils.Logln("Creating: %s %s\n", apiProductInfo.ID.APIProductName, apiProductInfo.ID.Version)
		} else {
			utils.Logln("Existing API Product found, attempting to update it...")
			utils.Logln("API Product ID:", id)
			updateAPIProduct = true
		}
	}

	if err != nil {
		utils.HandleErrorAndExit("Error getting OAuth Tokens", err)
	}
	extraParams := map[string]string{}
	adminEndpoint += "/import/api-product" + "?preserveProvider=" + strconv.FormatBool(importAPIProductPreserveProvider)

	// If the user has specified import-apis flag or update-apis flag, importAPIs parameter should be passed as true
	// because update is also an import task
	if importAPIs || importAPIsUpdate {
		adminEndpoint += "&importAPIs=" + strconv.FormatBool(true)
	}

	// If the user need to update the APIs and the API Product, overwriteAPIs parameter should be passed as true
	if importAPIsUpdate {
		adminEndpoint += "&overwriteAPIs=" + strconv.FormatBool(true)
	}

	// If the user need only to update the API Product, overwriteAPIProduct parameter should be passed as true
	if updateAPIProduct {
		adminEndpoint += "&overwriteAPIProduct=" + strconv.FormatBool(true)
	}

	utils.Logln(utils.LogPrefixInfo + "Import URL: " + adminEndpoint)
	err = importAPIProduct(adminEndpoint, apiProductFilePath, accessOAuthToken, extraParams)
	return err
}
