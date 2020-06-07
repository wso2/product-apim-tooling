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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/Jeffail/gabs"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	importAPIProductFile                string
	importAPIProductEnvironment         string
	importAPIProductCmdPreserveProvider bool
	importAPIs                          bool
	importAPIProductUpdate              bool
	importAPIsUpdate                    bool
	importAPIProductSkipCleanup         bool
	reApiProductName                    = regexp.MustCompile(`[~!@#;:%^*()+={}|\\<>"',&/$]`)
)

const (
	// ImportAPIProduct command related usage info
	importAPIProductCmdLiteral   = "api-product"
	importAPIProductCmdShortDesc = "Import API Product"
	importAPIProductCmdLongDesc  = "Import an API Product to an environment"
)

const importAPIProductCmdExamples = utils.ProjectName + ` ` + importCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f qa/LeasingAPIProduct.zip -e dev
` + utils.ProjectName + ` ` + importCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f staging/CreditAPIProduct.zip -e production --update-api-product
` + utils.ProjectName + ` ` + importCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f ~/myapiproduct -e production
` + utils.ProjectName + ` ` + importCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f ~/myapiproduct -e production --update-api-product --update-apis
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// ImportAPIProductCmd represents the importAPIProduct command
var ImportAPIProductCmd = &cobra.Command{
	Use: importAPIProductCmdLiteral + " (--file <path-to-api-product> --environment " +
		"<environment-to-which-the-api-product-should-be-imported>)",
	Short:   importAPIProductCmdShortDesc,
	Long:    importAPIProductCmdLongDesc,
	Example: importAPIProductCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + importAPIProductCmdLiteral + " called")
		var apiProductsExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApiProductsDirName)

		cred, err := getCredentials(importAPIProductEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeImportAPIProductCmd(cred, apiProductsExportDirectory)
	},
}

// executeImportAPIProductCmd executes the import api product command
func executeImportAPIProductCmd(credential credentials.Credential, exportDirectory string) {
	adminEndpoint := utils.GetAdminEndpointOfEnv(importAPIProductEnvironment, utils.MainConfigFilePath)
	err := ImportAPIProduct(credential, importAPIProductFile, adminEndpoint, exportDirectory)
	if err != nil {
		utils.HandleErrorAndExit("Error importing API Product", err)
		return
	}
}

// extractAPIProductDefinition extracts API Product information from jsonContent
func extractAPIProductDefinition(jsonContent []byte) (*v2.APIProductDefinition, error) {
	apiProduct := &v2.APIProductDefinition{}
	err := json.Unmarshal(jsonContent, &apiProduct)
	if err != nil {
		return nil, err
	}

	return apiProduct, nil
}

// getAPIProductDefinition scans filePath and returns APIProductDefinition or an error
func getAPIProductDefinition(filePath string) (*v2.APIProductDefinition, []byte, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}

	var buffer []byte
	if info.IsDir() {
		_, content, err := resolveYamlOrJson(path.Join(filePath, "Meta-information", "api"))
		if err != nil {
			return nil, nil, err
		}
		buffer = content
	} else {
		return nil, nil, fmt.Errorf("looking for directory, found %s", info.Name())
	}
	apiProduct, err := extractAPIProductDefinition(buffer)
	if err != nil {
		return nil, nil, err
	}
	return apiProduct, buffer, nil
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

// getApiProductID returns id of the API Product by using apiProductInfo which contains name, version and provider as info
func getApiProductID(name, version, environment, accessOAuthToken string) (string, error) {
	apiProductQuery := fmt.Sprintf("name:%s version:%s", name, version)
	apiProductQuery += " type:\"" + utils.DefaultApiProductType + "\""
	// Unified Search endpoint from the config file to search API Products
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(importAPIProductEnvironment, utils.MainConfigFilePath)
	count, apiProducts, err := GetAPIProductList(url.QueryEscape(apiProductQuery), "", accessOAuthToken, unifiedSearchEndpoint)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", nil
	}
	return apiProducts[0].ID, nil
}

func populateApiProductWithDefaults(def *v2.APIProductDefinition) (dirty bool) {
	dirty = false
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
func validateApiProductDefinition(def *v2.APIProductDefinition) error {
	utils.Logln(utils.LogPrefixInfo + "Validating API Product")
	if isEmpty(def.ID.APIProductName) {
		return errors.New("apiProductName is required")
	}
	if reApiProductName.MatchString(def.ID.APIProductName) {
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
func importAPIProduct(endpoint, httpMethod, filePath, accessToken string, extraParams map[string]string) error {
	req, err := NewFileUploadRequest(endpoint, httpMethod, extraParams, "file",
		filePath, accessToken)
	if err != nil {
		return err
	}

	var tr *http.Transport
	if utils.Insecure {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: utils.GetTlsConfigWithCertificate(),
		}
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(utils.HttpRequestTimeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		// 201 Created or 200 OK
		_ = resp.Body.Close()
		fmt.Println("Successfully imported API Product")
		return nil
	} else {
		// We have an HTTP error
		fmt.Println("Error importing API Product.")
		fmt.Println("Status: " + resp.Status)

		bodyBuf, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}

		strBody := string(bodyBuf)
		fmt.Println("Response:", strBody)

		return errors.New(resp.Status)
	}
}

// preProcessDependentAPIs pre processes dependent APIs
func preProcessDependentAPIs(apiProductFilePath string) error {
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
		paramsPath := apiDirectoryPath + string(os.PathSeparator) + DefaultAPIMParamsFileName
		// Check whether api_params.yaml file is available inside the particular API directory
		if utils.IsFileExist(paramsPath) {
			// Reading API params file and populate api.yaml
			err := injectParamsToAPI(apiDirectoryPath, paramsPath, importAPIProductEnvironment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ImportAPIProduct function is used with import-api-product command
func ImportAPIProduct(credential credentials.Credential, importPath, adminEndpoint, exportDirectory string) error {
	resolvedApiProductFilePath, err := resolveImportAPIProductFilePath(importPath, exportDirectory)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"API Product Location:", resolvedApiProductFilePath)

	utils.Logln(utils.LogPrefixInfo + "Creating workspace")
	tmpPath, err := getTempApiDirectory(resolvedApiProductFilePath)
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
	err = preProcessDependentAPIs(apiProductFilePath)
	if err != nil {
		return err
	}

	utils.Logln(utils.LogPrefixInfo + "Substituting environment variables in API Product files...")
	err = replaceEnvVariables(apiProductFilePath)
	if err != nil {
		return err
	}

	// Get API Product info
	apiProductInfo, originalContent, err := getAPIProductDefinition(apiProductFilePath)
	if err != nil {
		return err
	}
	// Fill with defaults
	if populateApiProductWithDefaults(apiProductInfo) {
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
	if err = validateApiProductDefinition(apiProductInfo); err != nil {
		return err
	}

	// If apiProductFilePath contains a directory, zip it
	if info, err := os.Stat(apiProductFilePath); err == nil && info.IsDir() {
		tmp, err := ioutil.TempFile("", "api-artifact*.zip")
		if err != nil {
			return err
		}
		utils.Logln(utils.LogPrefixInfo+"Creating API Product artifact", tmp.Name())
		err = utils.Zip(apiProductFilePath, tmp.Name())
		if err != nil {
			return err
		}
		defer func() {
			if importAPIProductSkipCleanup {
				utils.Logln(utils.LogPrefixInfo+"Leaving", tmp.Name())
				return
			}
			utils.Logln(utils.LogPrefixInfo+"Deleting", tmp.Name())
			err := os.Remove(tmp.Name())
			if err != nil {
				utils.Logln(utils.LogPrefixError + err.Error())
			}
		}()
		apiProductFilePath = tmp.Name()
	}

	updateAPIProduct := false
	if importAPIsUpdate || importAPIProductUpdate {
		accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, importAPIProductEnvironment)
		if err != nil {
			return err
		}

		// Check for API Product existence
		id, err := getApiProductID(apiProductInfo.ID.APIProductName, apiProductInfo.ID.Version, importAPIProductEnvironment, accessOAuthToken)
		if err != nil {
			return err
		}

		if id == "" {
			updateAPIProduct = false
			fmt.Println("The specified API Product was not found.")
			fmt.Printf("Creating: %s %s\n", apiProductInfo.ID.APIProductName, apiProductInfo.ID.Version)
		} else {
			fmt.Println("Existing API Product found, attempting to update it...")
			fmt.Println("API Product ID:", id)
			updateAPIProduct = true
		}
	}

	accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, importAPIProductEnvironment)
	if err != nil {
		utils.HandleErrorAndExit("Error getting OAuth Tokens", err)
	}
	extraParams := map[string]string{}
	httpMethod := http.MethodPost
	adminEndpoint += "/import/api-product" + "?preserveProvider=" + strconv.FormatBool(importAPIProductCmdPreserveProvider)

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
	err = importAPIProduct(adminEndpoint, httpMethod, apiProductFilePath, accessOAuthToken, extraParams)
	return err
}

// init using Cobra
func init() {
	ImportCmd.AddCommand(ImportAPIProductCmd)
	ImportAPIProductCmd.Flags().StringVarP(&importAPIProductFile, "file", "f", "",
		"Name of the API Product to be imported")
	ImportAPIProductCmd.Flags().StringVarP(&importAPIProductEnvironment, "environment", "e",
		"", "Environment from the which the API Product should be imported")
	ImportAPIProductCmd.Flags().BoolVar(&importAPIProductCmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of API Product after importing")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIs, "import-apis", "", false, "Import "+
		"dependent APIs associated with the API Product")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIProductUpdate, "update-api-product", "", false, "Update an "+
		"existing API Product or create a new API Product")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIsUpdate, "update-apis", "", false, "Update existing dependent APIs "+
		"associated with the API Product")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIProductSkipCleanup, "skipCleanup", "", false, "Leave "+
		"all temporary files created during import process")
	// Mark required flags
	_ = ImportAPIProductCmd.MarkFlagRequired("environment")
	_ = ImportAPIProductCmd.MarkFlagRequired("file")
}
