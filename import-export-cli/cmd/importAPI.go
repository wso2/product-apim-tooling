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
	"archive/zip"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var importAPIFile string
var importEnvironment string
var importAPICmdUsername string
var importAPICmdPassword string
var importAPICmdPreserveProvider bool
var importAPIUpdate bool
var importAPIConfigFile string
var importAPIInject bool

// ImportAPI command related usage info
const importAPICmdLiteral = "import-api"
const importAPICmdShortDesc = "Import API"

const DefaultAPIMParamsFileName = "api_params.yaml"

type ApiInfo struct {
	ID IdInfo `json:"id"`
}

type IdInfo struct {
	Name     string `json:"apiName"`
	Version  string `json:"version"`
	Provider string `json:"providerName"`
}

const importAPICmdLongDesc = "Import an API to an environment"

const importAPICmdExamples = utils.ProjectName + ` ` + importAPICmdLiteral + ` -f qa/TwitterAPI.zip -e dev
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f staging/FacebookAPI.zip -e production -u admin -p admin
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f ~/myapi -e production -u admin -p admin --update
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f ~/myapi -e production -u admin -p admin --update --inject`

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use: importAPICmdLiteral + " --file <Path to API> --environment " +
		"<Environment to be imported>",
	Short:   importAPICmdShortDesc,
	Long:    importAPICmdLongDesc,
	Example: importAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + importAPICmdLiteral + " called")
		var apisExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)
		executeImportAPICmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, apisExportDirectory)
	},
}

// executeImportAPICmd executes the import api command
func executeImportAPICmd(mainConfigFilePath, envKeysAllFilePath, exportDirectory string) {
	b64encodedCredentials, preCommandErr :=
		utils.ExecutePreCommandWithBasicAuth(importEnvironment, importAPICmdUsername, importAPICmdPassword,
			mainConfigFilePath, envKeysAllFilePath)

	if preCommandErr == nil {
		apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(importEnvironment, mainConfigFilePath)
		err := ImportAPI(importAPIFile, apiImportExportEndpoint, b64encodedCredentials, exportDirectory, importAPIConfigFile)
		if err != nil {
			utils.HandleErrorAndExit("Error importing API", err)
			return
		}
	} else {
		// env_endpoints file is not configured properly by the user
		fmt.Println("Error:", preCommandErr)
		utils.Logln(utils.LogPrefixError + preCommandErr.Error())
	}
}

// extractAPIDefinition extracts API information from jsonContent
func extractAPIDefinition(jsonContent []byte) (*APIDefinition, error) {
	api := &APIDefinition{}
	err := json.Unmarshal(jsonContent, &api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

// getAPIDefinition scans filePath and returns APIDefinition or an error
func getAPIDefinition(filePath string) (*APIDefinition, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	var buffer []byte
	if info.IsDir() {
		filePath = path.Join(filePath, "Meta-information", "api.json")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, err
		}

		// read file
		buffer, err = ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		// try reading zip file
		r, err := zip.OpenReader(filePath)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		for _, file := range r.File {
			// find api.json file inside the archive
			if strings.Contains(file.Name, "api.json") {
				rc, err := file.Open()
				if err != nil {
					return nil, err
				}

				buffer, err = ioutil.ReadAll(rc)
				if err != nil {
					_ = rc.Close()
					return nil, err
				}

				_ = rc.Close()
				break
			}
		}
	}

	api, err := extractAPIDefinition(buffer)
	if err != nil {
		return nil, err
	}
	return api, nil
}

// mergeAPI merges environmentParams to the API given in apiDirectory
// for now only Endpoints are merged
func mergeAPI(apiDirectory string, environmentParams *utils.Environment) error {
	// read api from Meta-information
	apiPath := path.Join(apiDirectory, "Meta-information", "api.json")
	utils.Logln("Reading API definition from: ", apiPath)
	api, err := gabs.ParseJSONFile(apiPath)
	if err != nil {
		return err
	}
	// extract environmentParams from file
	apiEndpointData, err := utils.ExtractAPIEndpointConfig(api.Bytes())
	if err != nil {
		return err
	}

	configData, err := json.Marshal(environmentParams.Endpoints)
	if err != nil {
		return err
	}
	mergedAPIEndpoints, err := utils.MergeJSON([]byte(apiEndpointData), configData)
	if err != nil {
		return err
	}

	utils.Logln("Writing merged API to:", apiPath)
	// replace original endpointConfig with merged version
	_, err = api.SetP(string(mergedAPIEndpoints), "endpointConfig")
	if err != nil {
		return err
	}

	// write this to disk
	err = ioutil.WriteFile(apiPath, api.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
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

// extractArchive extracts the API and give the path.
// In API Manager archive there is a directory in the root which contains the API
// this function returns it appended to the destination path
func extractArchive(src, dest string) (string, error) {
	files, err := utils.Unzip(src, dest)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("invalid API archive")
	}
	r := strings.TrimPrefix(files[0], src)
	fmt.Println(files[0], src)
	return filepath.Join(dest, strings.Split(path.Clean(r), string(os.PathSeparator))[0]), nil
}

// resolveAPIParamsPath resolves api_params.yaml path
// First it will look at BasePath f the import path (the last directory)
// If not found it will look at current working directory
// If a path is provided search ends looking up at that path
func resolveAPIParamsPath(importPath, paramPath string) (string, error) {
	utils.Logln(utils.LogPrefixInfo + "Scanning for " + DefaultAPIMParamsFileName)
	if paramPath == DefaultAPIMParamsFileName {
		// look in the basepath of importPath
		base := filepath.Dir(importPath)
		utils.Logln(utils.LogPrefixInfo+"Scanning in", base)
		fp := filepath.Join(base, DefaultAPIMParamsFileName)
		if info, err := os.Stat(fp); err == nil && !info.IsDir() {
			// found api_params.yml in the basepath
			return fp, nil
		}

		// look in the current working directory
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		utils.Logln(utils.LogPrefixInfo+"Scanning in", wd)
		fp = filepath.Join(wd, DefaultAPIMParamsFileName)
		if info, err := os.Stat(fp); err == nil && !info.IsDir() {
			// found api_params.yml in the cwd
			return fp, nil
		}

		// no luck, it means paramPath is missing
		return "", fmt.Errorf("could not find %s. Please check %s exists in basepath of "+
			"import location or current working directory", DefaultAPIMParamsFileName, DefaultAPIMParamsFileName)
	} else {
		if info, err := os.Stat(paramPath); err == nil && !info.IsDir() {
			return paramPath, nil
		}
		return "", fmt.Errorf("could not find %s", paramPath)
	}
}

// injectParamsToAPI injects ApiParams to API located in importPath using importEnvironment and returns the path to
// injected API location
func injectParamsToAPI(importPath, apiParamsPath, importEnvironment string) (string, error) {
	var dirPath string
	fileIsDir := false
	file := importPath

	paramsPath, err := resolveAPIParamsPath(file, apiParamsPath)
	if err != nil {
		return "", err
	}
	utils.Logln(utils.LogPrefixInfo+"Loading parameters from", paramsPath)
	apiParams, err := utils.LoadApiParamsFromFile(paramsPath)
	if err != nil {
		return "", err
	}
	// check whether import environment is included in api configuration
	envParams := apiParams.GetEnv(importEnvironment)
	if envParams == nil {
		return "", fmt.Errorf("%s does not exists in %s", importEnvironment, paramsPath)
	}

	// create a temp directory
	tmpDir, err := ioutil.TempDir("", "apim")
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}

	if info, err := os.Stat(file); err == nil {
		fileIsDir = info.IsDir()
	} else {
		return "", err
	}
	if fileIsDir {
		// copy dir to a temp location
		utils.Logln(utils.LogPrefixInfo+"Copying from", file, "to", tmpDir)
		dest := path.Join(tmpDir, filepath.Base(file))
		err = utils.CopyDir(file, dest)
		if err != nil {
			return "", err
		}
		dirPath = dest
	} else {
		// try to extract archive
		utils.Logln(utils.LogPrefixInfo+"Extracting", file, "to", tmpDir)
		finalPath, err := extractArchive(file, tmpDir)
		if err != nil {
			return "", err
		}
		dirPath = finalPath
	}

	err = mergeAPI(dirPath, envParams)
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

// getApiID returns id of the API by using apiInfo which contains name, version and provider as info
func getApiID(name, version, provider, environment, accessOAuthToken string) (string, error) {
	apiQuery := fmt.Sprintf("name:%s version:%s", name, version)
	if provider != "" {
		apiQuery += " provider:" + provider
	}
	count, apis, err := GetAPIList(url.QueryEscape(apiQuery), accessOAuthToken,
		utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath))
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", nil
	}
	return apis[0].ID, nil
}

// validateApiDefinition validates an API against rules
func validateApiDefinition(def *APIDefinition) error {
	utils.Logln(utils.LogPrefixInfo + "Validating API")
	if def.ID.APIName == "" {
		return errors.New("apiName is required")
	}
	if def.ID.Version == "" {
		return errors.New("version is required")
	}
	if def.Context == "" {
		return errors.New("context is required")
	}
	if def.ContextTemplate == "" {
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

// NewFileUploadRequest forms an HTTP request
// Helper function for forming multi-part form data
// Returns the formed http request and errors
func NewFileUploadRequest(uri string, method string, params map[string]string, paramName, path,
	b64encodedCredentials string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, uri, body)
	request.Header.Add(utils.HeaderAuthorization, utils.HeaderValueAuthBasicPrefix+" "+b64encodedCredentials)
	request.Header.Add(utils.HeaderContentType, writer.FormDataContentType())
	request.Header.Add(utils.HeaderAccept, "*/*")
	request.Header.Add(utils.HeaderConnection, utils.HeaderValueKeepAlive)

	return request, err
}

// importAPI imports an API to the API manager
func importAPI(endpoint, httpMethod, filePath, accessToken string, extraParams map[string]string) error {
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
		tr = &http.Transport{}
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

	//var bodyContent []byte
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		// 201 Created or 200 OK
		_ = resp.Body.Close()
		fmt.Println("Successfully imported API")
		return nil
	} else {
		// We have an HTTP error
		fmt.Println("Error importing API.")
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

// ImportAPI function is used with import-api command
// @param name: name of the API (zipped file) to be imported
// @param apiManagerEndpoint: API Manager endpoint for the environment
// @param accessToken: OAuth2.0 access token for the resource being accessed
func ImportAPI(importPath, apiImportExportEndpoint, accessToken, exportDirectory, apiParamsPath string) error {
	apiFilePath, err := resolveImportFilePath(importPath, exportDirectory)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"API Location:", apiFilePath)

	// inject if required
	if importAPIInject {
		utils.Logln(utils.LogPrefixInfo + "Injecting parameters to the API")
		injectedPath, err := injectParamsToAPI(importPath, apiParamsPath, importEnvironment)
		if err != nil {
			return err
		}
		defer func() {
			utils.Logln(utils.LogPrefixInfo+"Deleting", injectedPath)
			err := os.RemoveAll(injectedPath)
			if err != nil {
				utils.Logln(utils.LogPrefixError + err.Error())
			}
		}()
		apiFilePath = injectedPath
	}

	// Get API info
	apiInfo, err := getAPIDefinition(apiFilePath)
	if err != nil {
		return err
	}
	// validate definition
	if err = validateApiDefinition(apiInfo); err != nil {
		return err
	}

	// if apiFilePath contains a directory, zip it
	if info, err := os.Stat(apiFilePath); err == nil && info.IsDir() {
		tmp, err := ioutil.TempFile("", "api-artifact*.zip")
		if err != nil {
			return err
		}
		utils.Logln(utils.LogPrefixInfo+"Creating API artifact", tmp.Name())
		err = utils.Zip(apiFilePath, tmp.Name())
		if err != nil {
			return err
		}
		defer func() {
			utils.Logln(utils.LogPrefixInfo+"Deleting", tmp.Name())
			err := os.Remove(tmp.Name())
			if err != nil {
				utils.Logln(utils.LogPrefixError + err.Error())
			}
		}()
		apiFilePath = tmp.Name()
	}

	apiID := ""
	updateAPI := false
	if importAPIUpdate {
		accessOAuthToken, err :=
			utils.ExecutePreCommandWithOAuth(importEnvironment, importAPICmdUsername, importAPICmdPassword,
				utils.MainConfigFilePath, utils.EnvKeysAllFilePath)
		if err != nil {
			return err
		}

		// check for API existence
		id, err := getApiID(apiInfo.ID.APIName, apiInfo.ID.Version, apiInfo.ID.ProviderName, importEnvironment, accessOAuthToken)
		if err != nil {
			return err
		}

		if id == "" {
			fmt.Println("The specified API was not found.")
			fmt.Printf("Creating: %s %s\n", apiInfo.ID.APIName, apiInfo.ID.Version)
		} else {
			fmt.Println("Existing API found, attempting to update it...")
			utils.Logln("API ID:", id)
			apiID = id
			updateAPI = true
		}
	}

	extraParams := map[string]string{}
	httpMethod := ""
	if updateAPI {
		httpMethod = http.MethodPut
		apiImportExportEndpoint += "/" + apiID
	} else {
		httpMethod = http.MethodPost
		apiImportExportEndpoint += "/import-api"
	}
	apiImportExportEndpoint += "?preserveProvider=" +
		strconv.FormatBool(importAPICmdPreserveProvider)
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + apiImportExportEndpoint)
	err = importAPI(apiImportExportEndpoint, httpMethod, apiFilePath, accessToken, extraParams)
	return err
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ImportAPICmd)
	ImportAPICmd.Flags().StringVarP(&importAPIFile, "file", "f", "",
		"Name of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importEnvironment, "environment", "e",
		"", "Environment from the which the API should be imported")
	ImportAPICmd.Flags().StringVarP(&importAPICmdUsername, "username", "u", "", "Username")
	ImportAPICmd.Flags().StringVarP(&importAPICmdPassword, "password", "p", "", "Password")
	ImportAPICmd.Flags().BoolVar(&importAPICmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of API after exporting")
	ImportAPICmd.Flags().BoolVarP(&importAPIUpdate, "update", "", false, "Update API "+
		"if exists. Otherwise it will create API")
	ImportAPICmd.Flags().StringVarP(&importAPIConfigFile, "params", "", DefaultAPIMParamsFileName,
		"Provide a API Manager params file")
	ImportAPICmd.Flags().BoolVarP(&importAPIInject, "inject", "", false, "Inject variables defined"+
		"in params file to the given API.")
	// Mark required flags
	_ = ImportAPICmd.MarkFlagRequired("environment")
	_ = ImportAPICmd.MarkFlagRequired("file")
}
