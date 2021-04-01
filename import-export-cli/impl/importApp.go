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
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-resty/resty/v2"

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
	devportalApplicationsEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return ImportApplication(accessToken, devportalApplicationsEndpoint, filename, appOwner, updateApplication, preserveOwner,
		skipSubscriptions, skipKeys, skipCleanup)
}

// ImportApplication function is used with import-app command
// @param accessToken: OAuth2.0 access token for the resource being accessed
// @param devportalApplicationsEndpoint: Dev Portal Applications Endpoint for the environment
// @param filename: name of the application (zipped file) to be imported
// @param appOwner: Owner of the application
// @param updateApplication: Update the application if it already exists
// @param preserveOwner: Preserve the owner after importing the application
// @param skipSubscriptions: Skip importing subscriptions
// @param skipKeys: skip importing keys of application
// @param skipCleanup: skip cleaning up temporary files created during the operation
func ImportApplication(accessToken, devportalApplicationsEndpoint, filename, appOwner string, updateApplication, preserveOwner,
	skipSubscriptions, skipKeys, skipCleanup bool) (*http.Response, error) {

	exportDirectory := filepath.Join(utils.ExportDirectory, utils.ExportedAppsDirName)
	devportalApplicationsEndpoint = utils.AppendSlashToString(devportalApplicationsEndpoint)

	applicationImportEndpoint := devportalApplicationsEndpoint + "import"
	applicationImportUrl := applicationImportEndpoint + "?appOwner=" + appOwner + utils.SearchAndTag + "preserveOwner=" +
		strconv.FormatBool(preserveOwner) + utils.SearchAndTag + "skipSubscriptions=" +
		strconv.FormatBool(skipSubscriptions) + utils.SearchAndTag + "skipApplicationKeys=" + strconv.FormatBool(skipKeys) +
		utils.SearchAndTag + "update=" + strconv.FormatBool(updateApplication)
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + applicationImportEndpoint)

	applicationFilePath, err := resolveApplicationImportFilePath(filename, exportDirectory)
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

	resp, err := NewAppFileUploadRequest(applicationImportUrl, extraParams, "file", applicationFilePath, accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Error executing request.", err)
	}

	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		// 201 Created or 200 OK
		fmt.Println("Successfully imported Application.")
		return nil, nil
	} else {
		// We have an HTTP error
		fmt.Println("Error importing Application.")
		fmt.Println("Status: " + resp.Status())
		fmt.Println("Response:", resp)
		return nil, errors.New(resp.Status())
	}
}

// resolveApplicationImportFilePath resolves the archive/directory for import
// First will resolve in given path, if not found will try to load from exported directory
func resolveApplicationImportFilePath(file, defaultExportDirectory string) (string, error) {
	// check current path
	utils.Logln(utils.LogPrefixInfo + "Resolving for Application path...")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// if the file not in given path it might be inside exported directory
		utils.Logln(utils.LogPrefixInfo+"Looking for Application in", defaultExportDirectory)
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

// NewFileUploadRequest form an HTTP Put request
// Helper function for forming multi-part form data
// Returns the formed http request and errors
func NewAppFileUploadRequest(uri string, params map[string]string, paramName, path,
	accessToken string) (*resty.Response, error) {
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

	// set headers
	headers := make(map[string]string)
	headers[utils.HeaderContentType] = writer.FormDataContentType()
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = "*/*"
	headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive

	resp, err := utils.InvokePOSTRequest(uri, headers, body.Bytes())

	return resp, err
}
