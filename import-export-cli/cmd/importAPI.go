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
	"github.com/Jeffail/gabs"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
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

const DefaultAPIMConfigFileName = ".apim-vars.yml"

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
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f staging/FacebookAPI.zip -e production -u admin -p admin`

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use: importAPICmdLiteral + " (--file <api-zip-file> --environment " +
		"<environment-to-which-the-api-should-be-imported>)",
	Short:   importAPICmdShortDesc,
	Long:    importAPICmdLongDesc,
	Example: importAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + importAPICmdLiteral + " called")
		var apisExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)
		executeImportAPICmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, apisExportDirectory)
	},
}

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

// getAPIInfo scans filePath and returns API or an error
func getAPIInfo(filePath string) (*ApiInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	var buffer []byte
	if info.IsDir() {
		filePath = path.Join(filePath, "Meta-information", "api.json")
		fmt.Println(filePath)
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

	api, err := extractAPIInfo(buffer)
	if err != nil {
		return nil, err
	}

	return api, nil
}

// extractAPIInfo extracts API information from jsonContent
func extractAPIInfo(jsonContent []byte) (*ApiInfo, error) {
	api := &ApiInfo{}
	err := json.Unmarshal(jsonContent, &api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

// mergeAPI merge API in filepath(this points to a directory extracted by tool) with configuration given in configPath
// it returns the directory containing merged files
func mergeAPI(apiFilePath string, endpointConfig *utils.Environment) (string, error) {
	// copy apiFilePath to a temp location for variable injection
	tmpDir, err := ioutil.TempDir("", "apim")
	if err != nil {
		return "", err
	}

	// copy contents to a directory which contains base name of original
	newPath := path.Join(tmpDir, path.Base(apiFilePath))
	err = utils.CopyDir(apiFilePath, newPath)
	if err != nil {
		return "", err
	}

	// read api from Meta-information
	apiPath := path.Join(newPath, "Meta-information", "api.json")
	utils.Logln("Reading API definition from: ", apiPath)
	api, err := gabs.ParseJSONFile(apiPath)
	if err != nil {
		return "", err
	}
	// extract endpointConfig from file
	apiEndpointData, err := utils.ExtractAPIEndpointConfig(api.Bytes())
	if err != nil {
		return "", err
	}

	configData, err := json.Marshal(endpointConfig.Endpoints)
	if err != nil {
		return "", err
	}
	mergedADIEndpoints, err := utils.MergeJSON([]byte(apiEndpointData), configData)
	if err != nil {
		return "", err
	}

	utils.Logln("Writing merged API to:", apiPath)
	// replace original endpointConfig with merged version
	_, err = api.SetP(string(mergedADIEndpoints), "endpointConfig")
	if err != nil {
		return "", err
	}

	// write this to disk
	err = ioutil.WriteFile(apiPath, api.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return newPath, nil
}

// ImportAPI function is used with import-api command
// @param name: name of the API (zipped file) to be imported
// @param apiManagerEndpoint: API Manager endpoint for the environment
// @param accessToken: OAuth2.0 access token for the resource being accessed
func ImportAPI(importPath, apiImportExportEndpoint, accessToken, exportDirectory, configPath string) error {
	apiID := ""
	updateAPI := false
	apiImportExportEndpoint = utils.AppendSlashToString(apiImportExportEndpoint)

	// fileName can be a environment related path like dev/PizzaShackAPI.zip
	fileName := importPath
	zipFilePath := fileName

	// Check whether the given path is a directory
	// If it is a directory, archive it
	if info, err := os.Stat(fileName); err == nil && info.IsDir() {
		// Get base name of the file
		fileBase := filepath.Base(fileName)
		fmt.Println(fileBase + " is a directory")

		if importAPIInject {
			utils.Logln("Looking up for " + configPath)
			if info, err := os.Stat(configPath); err == nil && info.IsDir() {
				// Need to check whether given file is a directory
				return errors.New(configPath + "is a directory")
			} else if os.IsNotExist(err) && configPath != DefaultAPIMConfigFileName {
				// config file is not mandatory. But if the given config file is not the default one need to check for
				// existence and return error if not found
				return err
			}

			utils.Logln(configPath + " found")
			utils.Logln("Scanning...")
			// load configuration from yml file
			apiConfig, err := utils.LoadConfigFromFile(configPath)
			if err != nil {
				return err
			}

			// check whether import environment is included in api configuration
			endpointConfig := apiConfig.GetEnv(importEnvironment)
			if endpointConfig == nil {
				return fmt.Errorf("%s does not exists in configuration file", importEnvironment)
			}

			utils.Logln("Merging...")
			mergedAPIDir, err := mergeAPI(fileName, endpointConfig)
			if err != nil {
				return err
			}

			// delete the temp directory on return
			defer func() {
				utils.Logln("Deleting:", mergedAPIDir)
				err := os.RemoveAll(mergedAPIDir)
				if err != nil {
					utils.HandleErrorAndExit("Error deleting directory:", err)
				}
			}()

			fileName = mergedAPIDir
		}

		fmt.Println("Creating an archive from the directory...")
		// create a temp file in OS temp directory
		tmpZip, err := ioutil.TempFile("", fileBase+"*.zip")
		if err != nil {
			return err
		}
		// delete the temp zip file on return
		defer func() {
			utils.Logln("Deleting:", tmpZip.Name())
			err := os.Remove(tmpZip.Name())
			if err != nil {
				utils.HandleErrorAndExit("Error deleting file:", err)
			}
		}()

		// zip the given directory
		absFilePath, err := filepath.Abs(fileName)
		if err != nil {
			return err
		}
		utils.Logln("Zipping: ", absFilePath)
		err = utils.Zip(absFilePath, tmpZip.Name())
		if err != nil {
			return err
		}
		// change our zip file path to new archive
		zipFilePath = tmpZip.Name()
	}

	// Test if we can find the zip file in the current work directory
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// Doesn't exist... Check if available in the default exportDirectory
		zipFilePath = filepath.Join(exportDirectory, fileName)
		if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
			return err
		}
	}
	utils.Logln("Archive path:", zipFilePath)

	accessOAuthToken, err :=
		utils.ExecutePreCommandWithOAuth(importEnvironment, importAPICmdUsername, importAPICmdPassword,
			utils.MainConfigFilePath, utils.EnvKeysAllFilePath)
	if err != nil {
		return err
	}

	// Get API info
	apiInfo, err := getAPIInfo(zipFilePath)
	if err != nil {
		return err
	}
	if importAPIUpdate {
		utils.Logln("Reading API meta data from: ", zipFilePath)

		if apiInfo.ID.Name == "" || apiInfo.ID.Provider == "" || apiInfo.ID.Version == "" {
			utils.Logln(utils.LogPrefixInfo, "API: ", apiInfo)
			return errors.New("invalid api information")
		}

		// check for API existence
		id, err := getApiID(apiInfo.ID.Name, apiInfo.ID.Version, apiInfo.ID.Provider, importEnvironment, accessOAuthToken)
		if err != nil {
			return err
		}

		if id == "" {
			fmt.Println("The specified API was not found.")
			fmt.Printf("Creating: %s %s\n", apiInfo.ID.Name, apiInfo.ID.Version)
		} else {
			fmt.Println("Existing API found, attempting to update it...")
			utils.Logln("API ID:", id)
			apiID = id
			updateAPI = true
		}
	}

	extraParams := map[string]string{}
	// TODO:: Add extraParams as necessary

	httpMethod := ""
	if updateAPI {
		httpMethod = http.MethodPut
		apiImportExportEndpoint += apiID
	} else {
		httpMethod = http.MethodPost
		apiImportExportEndpoint += "import-api"
	}

	apiImportExportEndpoint += "?preserveProvider=" +
		strconv.FormatBool(importAPICmdPreserveProvider)
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + apiImportExportEndpoint)
	err = importAPI(apiImportExportEndpoint, httpMethod, zipFilePath, accessToken, extraParams)
	return err
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

// init using Cobra
func init() {
	RootCmd.AddCommand(ImportAPICmd)
	ImportAPICmd.Flags().StringVarP(&importAPIFile, "file", "f", "",
		"Name of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importEnvironment, "environment", "e",
		utils.DefaultEnvironmentName, "Environment from the which the API should be imported")
	ImportAPICmd.Flags().StringVarP(&importAPICmdUsername, "username", "u", "", "Username")
	ImportAPICmd.Flags().StringVarP(&importAPICmdPassword, "password", "p", "", "Password")
	ImportAPICmd.Flags().BoolVar(&importAPICmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of API after exporting")
	ImportAPICmd.Flags().BoolVarP(&importAPIUpdate, "update", "", false, "Update API "+
		"if exists. Otherwise it will create API")
	// TODO: finalize a good name for file
	ImportAPICmd.Flags().StringVarP(&importAPIConfigFile, "config", "", DefaultAPIMConfigFileName,
		"Provide a API Manager configuration file.")
	ImportAPICmd.Flags().BoolVarP(&importAPIInject, "inject", "", false, "Inject variables defined"+
		"in config file to the given API.")
}
