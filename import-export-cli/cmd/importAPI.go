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
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"strings"
	"time"

	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var importAPIFile string
var importEnvironment string
var importAPICmdUsername string
var importAPICmdPassword string
var importAPICmdPreserveProvider bool

// ImportAPI command related usage info
const importAPICmdLiteral = "import-api"
const importAPICmdShortDesc = "Import API"

var importAPICmdLongDesc = "Import an API to an environment"

var importAPICmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f qa/TwitterAPI.zip -e dev
		` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f staging/FacebookAPI.zip -e production -u admin -p admin
	`)

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use: importAPICmdLiteral + " (--file <api-zip-file> --environment " +
		"<environment-to-which-the-api-should-be-imported>)",
	Short: importAPICmdShortDesc,
	Long:  importAPICmdLongDesc + importAPICmdExamples,
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

		resp, err := ImportAPI(importAPIFile, apiImportExportEndpoint, b64encodedCredentials, exportDirectory)

		if err != nil {
			utils.HandleErrorAndExit("Error importing API", err)
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			// 200 OK or 201 Created
			utils.Logln(utils.LogPrefixInfo+"Header:", resp.Header)
			fmt.Println("Successfully imported API!")
		} else {
			fmt.Println("Error importing API")
			utils.Logln(utils.LogPrefixError + resp.Status)
		}
	} else {
		// env_endpoints file is not configured properly by the user
		fmt.Println("Error:", preCommandErr)
		utils.Logln(utils.LogPrefixError + preCommandErr.Error())
	}
}

// ImportAPI function is used with import-api command
// @param name: name of the API (zipped file) to be imported
// @param apiManagerEndpoint: API Manager endpoint for the environment
// @param accessToken: OAuth2.0 access token for the resource being accessed
func ImportAPI(query, apiImportExportEndpoint, accessToken, exportDirectory string) (*http.Response, error) {
	apiImportExportEndpoint = utils.AppendSlashToString(apiImportExportEndpoint)

	apiImportExportEndpoint += "import-api"
	apiImportExportEndpoint += "?preserveProvider=" +
		strconv.FormatBool(importAPICmdPreserveProvider)
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + apiImportExportEndpoint)

	sourceEnv := strings.Split(query, "/")[0] // environment from which the API was exported
	utils.Logln(utils.LogPrefixInfo + "Source Environment: " + sourceEnv)

	fileName := query // ex:- fileName = dev/twitterapi_1.0.0.zip

	zipFilePath := filepath.Join(exportDirectory, fileName)
	fmt.Println("ZipFilePath:", zipFilePath)

	// check if '.zip' exists in the input 'fileName'
	//hasZipExtension, _ := regexp.MatchString(`^\S+\.zip$`, fileName)

	//if hasZipExtension {
	//	// import the zip file directly
	//	//fmt.Println("hasZipExtension: ", true)
	//
	//} else {
	//	//fmt.Println("hasZipExtension: ", false)
	//	// search for a directory with the given fileName
	//	destination := filepath.Join(exportDirectory, fileName+".zip")
	//	err := utils.ZipDir(zipFilePath, destination)
	//	if err != nil {
	//		utils.HandleErrorAndExit("Error creating zip archive", err)
	//	}
	//	zipFilePath += ".zip"
	//}

	extraParams := map[string]string{}
	// TODO:: Add extraParams as necessary

	req, err := NewFileUploadRequest(apiImportExportEndpoint, extraParams, "file", zipFilePath, accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Error creating request.", err)
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
	defer resp.Body.Close()

	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
	} else {
		//var bodyContent []byte
		if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
			// 201 Created or 200 OK
			fmt.Println("Successfully imported API '" + fileName + "'")
		} else {
			// We have an HTTP error
			fmt.Println("Error importing API.")
			fmt.Println("Status: " + resp.Status)

			boduBuf, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return nil, err
			}

			strBody := string(boduBuf)
			fmt.Println("Response:", strBody)

			return nil, errors.New(resp.Status)
		}
	}
	
	return resp, err
}

// NewFileUploadRequest form an HTTP Put request
// Helper function for forming multi-part form data
// Returns the formed http request and errors
func NewFileUploadRequest(uri string, params map[string]string, paramName, path,
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

	request, err := http.NewRequest(http.MethodPost, uri, body)
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
}
