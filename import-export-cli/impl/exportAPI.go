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
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExportAPIFromEnv function is used with export api command
func ExportAPIFromEnv(accessToken, name, version, revisionNum, provider, format, exportEnvironment string, preserveStatus,
	exportLatestRevision bool) (*resty.Response, error) {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportAPI(name, version, revisionNum, provider, format, publisherEndpoint, accessToken, preserveStatus,
		exportLatestRevision)
}

// exportAPI function is used with export api command
// @param name : Name of the API to be exported
// @param version : Version of the API to be exported
// @param provider : Provider of the API
// @param publisherEndpoint : API Manager Publisher Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func exportAPI(name, version, revisionNum, provider, format, publisherEndpoint, accessToken string, preserveStatus,
	exportLatestRevision bool) (*resty.Response, error) {
	publisherEndpoint = utils.AppendSlashToString(publisherEndpoint)
	query := "apis/export?name=" + url.QueryEscape(name) + "&version=" + version + "&providerName=" + provider +
		"&preserveStatus=" + strconv.FormatBool(preserveStatus)
	if format != "" {
		query += "&format=" + format
	}
	if revisionNum != "" {
		query += "&revisionNumber=" + revisionNum
	}
	if exportLatestRevision {
		query += "&latestRevision=true"
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
// @param exportAPIRevisionNumber: Revision number of the api
// @param zipLocationPath: Path to the export directory
// @param runningExportApiCommand: Whether the export API command is running
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported API will be written to a zip file
func WriteToZip(exportAPIName, exportAPIVersion, exportAPIRevisionNumber, zipLocationPath string,
	runningExportApiCommand bool, resp *resty.Response) {
	zipFilename := exportAPIName + "_" + exportAPIVersion
	if exportAPIRevisionNumber != "" {
		zipFilename += "_" + utils.GetRevisionNamFromRevisionNum(exportAPIRevisionNumber)
	}
	zipFilename += ".zip" // MyAPI_1.0.0_Revision-1.zip
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

	// Add api_meta.yaml file inside the zip and create a new zup file in exportedFinalZip location
	metaData := utils.MetaData{
		Name:    exportAPIName,
		Version: exportAPIVersion,
		DeployConfig: utils.DeployConfig{
			Import: utils.ImportConfig{
				Update:           true,
				PreserveProvider: true,
				RotateRevision:   false,
			},
		},
	}
	err = IncludeMetaFileToZip(tempZipFile, exportedFinalZip, utils.MetaFileAPI, metaData)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive with api_meta.yaml file", err)
	}

	// Output the final zip file location.
	if runningExportApiCommand {
		fmt.Println("Successfully exported API!")
		fmt.Println("Find the exported API at " + exportedFinalZip)
	}
}
