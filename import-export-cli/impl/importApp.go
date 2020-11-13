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
	"github.com/go-resty/resty"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

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
	applicationImportUrl := applicationImportEndpoint + "?appOwner=" + appOwner + utils.SearchAndTag + "preserveOwner=" +
		strconv.FormatBool(preserveOwner) + utils.SearchAndTag + "skipSubscriptions=" +
		strconv.FormatBool(skipSubscriptions) + utils.SearchAndTag + "skipApplicationKeys=" + strconv.FormatBool(skipKeys) +
		utils.SearchAndTag + "update=" + strconv.FormatBool(updateApplication)
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + applicationImportEndpoint)

	applicationFilePath, err := resolveImportFilePath(filename, exportDirectory)
	if err != nil {
		utils.HandleErrorAndExit("Error creating request.", err)
	}

	utils.Logln(utils.LogPrefixInfo + "Pre Processing Application...")
	error := preProcessApplication(applicationFilePath)
	if error != nil {
		utils.HandleErrorAndExit("Error importing Application", error)
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

//This function will check whether the .json file is included in the directory or not
//True is included false otherwise
func checkDirForJson(appDirectory string) bool {
	file, err := os.Open(appDirectory)

	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	//list all files and check for supported type file
	for _, name := range list {
		match, _ := regexp.MatchString(".*\\.json", name)
		if match {
			return true
		}
	}
	return false
}

func preProcessApplication(appDirectory string) error {
	utils.Logln(utils.LogPrefixInfo+"Loading Application definition from: ", appDirectory)
	file, err := os.Open(appDirectory)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	fileStat, error := file.Stat()
	if error != nil {
		log.Fatalf("failed checking directory: %s", err)
	}
	defer file.Close()

	if fileStat.IsDir() {
		isJsonFileExisted := checkDirForJson(appDirectory)
		if isJsonFileExisted {
			return nil
		} else {
			return fmt.Errorf("Supported file type is not found in the %s directory", appDirectory)
		}
	}
	return nil
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

	resp, err := utils.InvokePOSTRequestWithBytes(uri, headers, body.Bytes())

	return resp, err
}
