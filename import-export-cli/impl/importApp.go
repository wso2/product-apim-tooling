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
	"time"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ImportApplicationToEnv function is used with import-app command
// @param accessToken: OAuth2.0 access token for the resource being accessed
// @param environment: Environment to import the application
// @param filename: name of the application (zipped file) to be imported
// @param appOwner: Owner of the application
// @param updateApplication: Update the application if it already exists
// @param preserveOwner: Preserve the owner after importing the application
// @param skipSubscriptions: Skip importing subscriptions
// @param skipKeys: skip importing keys of application
// @param skipCleanup: skip cleaning up temporary files created during the operation
func ImportApplicationToEnv(accessToken, environment, filename, appOwner string, updateApplication, preserveOwner,
	skipSubscriptions, skipKeys, skipCleanup bool) (*http.Response, error) {
	adminEndpoint := utils.GetAdminEndpointOfEnv(environment, utils.MainConfigFilePath)
	return ImportApplication(accessToken, adminEndpoint, filename, appOwner, updateApplication, preserveOwner,
		skipSubscriptions, skipKeys, skipCleanup)
}

// ImportApplication function is used with import-app command
// @param accessToken: OAuth2.0 access token for the resource being accessed
// @param adminEndpoint: Admin REST API endpoint to use for importing the application
// @param filename: name of the application (zipped file) to be imported
// @param appOwner: Owner of the application
// @param updateApplication: Update the application if it already exists
// @param preserveOwner: Preserve the owner after importing the application
// @param skipSubscriptions: Skip importing subscriptions
// @param skipKeys: skip importing keys of application
// @param skipCleanup: skip cleaning up temporary files created during the operation
func ImportApplication(accessToken, adminEndpoint, filename, appOwner string, updateApplication, preserveOwner,
	skipSubscriptions, skipKeys, skipCleanup bool) (*http.Response, error) {

	exportDirectory := filepath.Join(utils.ExportDirectory, utils.ExportedAppsDirName)
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)

	applicationImportEndpoint := adminEndpoint + "import/applications"
	url := applicationImportEndpoint + "?appOwner=" + appOwner + utils.SearchAndTag + "preserveOwner=" +
		strconv.FormatBool(preserveOwner) + utils.SearchAndTag + "skipSubscriptions=" +
		strconv.FormatBool(skipSubscriptions) + utils.SearchAndTag + "skipApplicationKeys=" + strconv.FormatBool(skipKeys) +
		utils.SearchAndTag + "update=" + strconv.FormatBool(updateApplication)
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + applicationImportEndpoint)

	applicationFilePath, err := resolveImportFilePath(filename, exportDirectory)
	if err != nil {
		utils.HandleErrorAndExit("Error creating request.", err)
	}

	// If applicationFilePath contains a directory, zip it. Otherwise, leave it as it is.
	applicationFilePath, err, cleanupFunc := utils.CreateZipFileFromProject(applicationFilePath, skipCleanup)
	if err != nil {
		return nil, err
	}

	//cleanup the temporary artifacts once consuming the zip file
	if cleanupFunc != nil {
		defer cleanupFunc()
	}

	extraParams := map[string]string{}

	req, err := NewAppFileUploadRequest(url, extraParams, "file", applicationFilePath, accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Error creating request.", err)
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
		return nil, err
	}
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		// 201 Created or 200 OK
		_ = resp.Body.Close()
		fmt.Println("Successfully imported Application")
		return resp, nil
	} else {
		// We have an HTTP error
		fmt.Println("Error importing Application.")
		fmt.Println("Status: " + resp.Status)

		bodyBuf, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		strBody := string(bodyBuf)
		fmt.Println("Response:", strBody)

		return nil, errors.New(resp.Status)
	}
}

// NewFileUploadRequest form an HTTP Put request
// Helper function for forming multi-part form data
// Returns the formed http request and errors
func NewAppFileUploadRequest(uri string, params map[string]string, paramName, path,
	accessToken string) (*http.Request, error) {
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
	request.Header.Add(utils.HeaderAuthorization, utils.HeaderValueAuthBearerPrefix+" "+accessToken)
	request.Header.Add(utils.HeaderContentType, writer.FormDataContentType())
	request.Header.Add(utils.HeaderAccept, "*/*")
	request.Header.Add(utils.HeaderConnection, utils.HeaderValueKeepAlive)

	return request, err
}
