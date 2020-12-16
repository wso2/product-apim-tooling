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
	"strconv"

	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExportAPIFromEnv function is used with export api command
func ExportAPIFromEnv(accessToken, name, version, provider, format, exportEnvironment string, preserveStatus bool) (*resty.Response, error) {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportAPI(name, version, provider, format, publisherEndpoint, accessToken, preserveStatus)
}

// exportAPI function is used with export api command
// @param name : Name of the API to be exported
// @param version : Version of the API to be exported
// @param provider : Provider of the API
// @param publisherEndpoint : API Manager Publisher Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func exportAPI(name, version, provider, format, publisherEndpoint, accessToken string, preserveStatus bool) (*resty.Response, error) {
	publisherEndpoint = utils.AppendSlashToString(publisherEndpoint)
	query := "apis/export?name=" + name + "&version=" + version + "&providerName=" + provider +
		"&preserveStatus=" + strconv.FormatBool(preserveStatus)
	if format != "" {
		query += "&format=" + format
	}

	url := publisherEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportAPI: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(url, headers)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WriteToZip
// @param exportAPIName : Name of the API to be exported
// @param exportAPIVersion: Version of the API to be exported
// @param zipLocationPath: Path to the export directory
// @param runningExportApiCommand: Whether the export API command is running
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported API will be written to a zip file
func WriteToZip(exportAPIName, exportAPIVersion, zipLocationPath string, runningExportApiCommand bool, resp *resty.Response) {
	zipFilename := exportAPIName + "_" + exportAPIVersion + ".zip" // MyAPI_1.0.0.zip
	// Writes the REST API response to a temporary zip file
	tempZipFile, err := utils.WriteResponseToTempZip(zipFilename, resp)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the temporary zip file to store the exported API", err)
	}

	err = utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archive: "+zipLocationPath, err)
	}
	exportedFinalZip := filepath.Join(zipLocationPath, zipFilename)
	// Add api_params.yaml file inside the zip and create a new zip file in exportedFinalZip location
	err = IncludeParamsFileToZip(tempZipFile, exportedFinalZip, utils.ParamFileAPI)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive", err)
	}

	// Output the final zip file location.
	if runningExportApiCommand {
		fmt.Println("Successfully exported API!")
		fmt.Println("Find the exported API at " + exportedFinalZip)
	}
}
