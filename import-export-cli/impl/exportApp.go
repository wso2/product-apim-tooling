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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExportAppFromEnv function is used with export app command
func ExportAppFromEnv(accessToken, name, owner, exportEnvironment string, exportAppWithKeys bool) (*resty.Response, error) {
	adminEndpiont := utils.GetAdminEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return ExportApp(name, owner, adminEndpiont, accessToken, exportAppWithKeys)
}

// ExportApp
// @param name : Name of the Application to be exported
// @param owner : Owner of the Application to be exported
// @param adminEndpoint : Admin Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func ExportApp(name, owner, adminEndpoint, accessToken string, exportAppWithKeys bool) (*resty.Response, error) {
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	query := "export/applications?appName=" + name + utils.SearchAndTag + "appOwner=" + owner

	if exportAppWithKeys {
		query += "&withKeys=true"
	}

	url := adminEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportApp: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(url, headers)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WriteApplicationToZip
// @param exportAppName : Name of the Application to be exported
// @param exportAppOwner : Owner of the Application to be exported
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported Application will be written to a zip file
func WriteApplicationToZip(exportAppName, exportAppOwner, zipLocationPath string,
	resp *resty.Response) {
	zipFilename := replaceUserStoreDomainDelimiter(exportAppOwner) + "_" + exportAppName + ".zip" // admin_testApp.zip
	// Writes the REST API response to a temporary zip file
	tempZipFile, err := utils.WriteResponseToTempZip(zipFilename, resp)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the temporary zip file to store the exported application", err)
	}

	err = utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archive: "+zipLocationPath, err)
	}

	exportedFinalZip := filepath.Join(zipLocationPath, zipFilename)
	// Add application_params.yaml file inside the zip and create a new zip file in exportedFinalZip location
	err = IncludeParamsFileToZip(tempZipFile, exportedFinalZip, utils.ParamFileApplication)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive", err)
	}
	fmt.Println("Successfully exported Application!")
	fmt.Println("Find the exported Application at " + exportedFinalZip)
}

// The Application owner name is used to construct a unique name for the app export zip.
// When an app belonging to a user from a secondary user store is exported, the owner name will have
// the format '<Userstore_domain>/<Username>'. The '/' character will be mistakenly considerd as a
// file separator character, resulting in an invalid path being constructed.
// Therefore this function overcomes this issue by replacing the '/' character.
func replaceUserStoreDomainDelimiter(username string) string {
	return strings.ReplaceAll(username, "/", "#")
}
