/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"fmt"

	"bytes"
	"crypto/tls"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"time"
)

var importAPIName string
var importEnvironment string
var importAPICmdUsername string
var importAPICmdPassword string

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use:   "import-api (--name <name-of-the-api> --environment <environment-to-which-the-api-should-be-imported>)",
	Short: utils.ImportAPICmdShortDesc,
	Long:  utils.ImportAPICmdLongDesc + utils.ImportAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "import-api called")

		accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(importEnvironment, importAPICmdUsername,
			importAPICmdPassword)

		if preCommandErr == nil {

			resp, _ := ImportAPI(importAPIName, apiManagerEndpoint, accessToken)
			utils.Logln(utils.LogPrefixInfo+"Status: %v\n", resp.Status)
			if resp.StatusCode == 200 {
				fmt.Println("Header:", resp.Header)
			} else {
				fmt.Println("Status: ", resp.Status)
				utils.Logln(utils.LogPrefixError + resp.Status)
			}
		} else {
			// env_endpoints_all.yaml file is not configured properly by the user
			log.Fatal("Error:", preCommandErr)
			utils.Logln(utils.LogPrefixError + preCommandErr.Error())
		}
	},
}

func ImportAPI(name string, url string, accessToken string) (*http.Response, error) {
	// append '/' to the end if there isn't one already
	if string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "import/apis"

	filePath := utils.ExportedAPIsDirectoryPath + utils.PathSeparator_ + name
	fmt.Println("filePath:", filePath)

	// check if '.zip' exists in the input 'name'
	hasZipExtension, _ := regexp.MatchString(`^\S+\.zip$`, name)

	if hasZipExtension {
		// import the zip file directly
		//fmt.Println("hasZipExtension: ", true)

	} else {
		//fmt.Println("hasZipExtension: ", false)
		// search for a directory with the given name
		destination := utils.ExportedAPIsDirectoryPath + utils.PathSeparator_ + name + ".zip"
		err := utils.ZipDir(filePath, destination)
		if err != nil {
			utils.HandleErrorAndExit("Error creating zip archive", err)
		}
		filePath += ".zip"
	}

	extraParams := map[string]string{}
	// TODO:: Add extraParams as necessary

	req, err := NewFileUploadRequest(url, extraParams, "file", filePath, accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Error creating request.", err)
	}

	var tr *http.Transport
	if utils.SkipTLSVerification {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}else{
		tr = &http.Transport{}
	}

	client := &http.Client{
		Transport: tr,
		Timeout: time.Duration(utils.HttpRequestTimeout) * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
	} else {
		//var bodyContent []byte
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		//resp.Body.Read(bodyContent)
		//resp.Body.Close()
		//fmt.Println(bodyContent)
	}

	return resp, err
}

// Helper function for forming multi-part form data
// Returns the formed http request and errors
func NewFileUploadRequest(uri string, params map[string]string, paramName, path string,
	accessToken string) (*http.Request,error) {
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

	request, err := http.NewRequest(http.MethodPut, uri, body)
	request.Header.Add(utils.HeaderAuthorization, utils.HeaderValueAuthBearerPrefix+" "+accessToken)
	request.Header.Add(utils.HeaderContentType, writer.FormDataContentType())
	request.Header.Add(utils.HeaderAccept, "*/*")
	request.Header.Add(utils.HeaderConnection, utils.HeaderValueKeepAlive)

	return request, err
}

// Generated with Cobra
func init() {
	RootCmd.AddCommand(ImportAPICmd)
	ImportAPICmd.Flags().StringVarP(&importAPIName, "name", "n", "",
		"Name of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importEnvironment, "environment", "e", "",
		"Environment from the which the API should be imported")
	ImportAPICmd.Flags().StringVarP(&importAPICmdUsername, "username", "u", "", "Username")
	ImportAPICmd.Flags().StringVarP(&importAPICmdPassword, "password", "p", "", "Password")
}
