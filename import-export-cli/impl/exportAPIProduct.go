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

	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExportAPIProductFromEnv function is used with export api command
func ExportAPIProductFromEnv(accessToken, name, version, revisionNum, provider, format,
	exportEnvironment string, exportLatestRevision bool) (*resty.Response, error) {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportAPIProduct(name, version, revisionNum, provider, format, publisherEndpoint, accessToken,
		exportLatestRevision)
}

// exportAPIProduct
// @param name : Name of the API Product to be exported
// @param version : Version of the API Product to be exported
// @param provider : Provider of the API Product
// @param publisherEndpoint : API Manager Publisher Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func exportAPIProduct(name, version, revisionNum, provider, format, publisherEndpoint, accessToken string,
	exportLatestRevision bool) (*resty.Response, error) {
	publisherEndpoint = utils.AppendSlashToString(publisherEndpoint)
	query := "api-products/export?name=" + name + "&version=" + version + "&providerName=" + provider
	if revisionNum != "" {
		query += "&revisionNumber=" + revisionNum
	}
	if exportLatestRevision {
		query += "&latestRevision=true"
	}
	if format != "" {
		query += "&format=" + format
	}

	url := publisherEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportAPIProduct: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(url, headers)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WriteAPIProductToZip
// @param exportAPIProductName : Name of the API Product to be exported
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported API Product will be written to a zip file
func WriteAPIProductToZip(exportAPIProductName, exportAPIProductVersion, zipLocationPath string, runningExportAPIProductCommand bool, resp *resty.Response) {
	zipFilename := exportAPIProductName + "_" + exportAPIProductVersion + ".zip" // MyAPIProduct_1.0.0.zip
	// Writes the REST API response to a temporary zip file
	tempZipFile, err := utils.WriteResponseToTempZip(zipFilename, resp)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the temporary zip file to store the exported API Product", err)
	}

	err = utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archive: "+zipLocationPath, err)
	}
	exportedFinalZip := filepath.Join(zipLocationPath, zipFilename)
	// Add api_product_params.yaml file inside the zip and create a new zip file in exportedFinalZip location
	err = IncludeParamsFileToZip(tempZipFile, exportedFinalZip, utils.ParamFileAPIProduct)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive", err)
	}

	// Add api_product_meta.yaml file inside the zip and create a new zup file in exportedFinalZip location
	metaData := utils.MetaData{
		Name:    exportAPIProductName,
		Version: exportAPIProductVersion,
	}
	err = IncludeMetaFileToZip(tempZipFile, exportedFinalZip, utils.MetaFileAPIProduct, metaData)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive with api_product_meta.yaml file", err)
	}

	if runningExportAPIProductCommand {
		fmt.Println("Successfully exported API Product!")
		fmt.Println("Find the exported API Product at " + exportedFinalZip)
	}
}
