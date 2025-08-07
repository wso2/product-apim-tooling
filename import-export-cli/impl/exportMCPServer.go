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
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExportMCPServerFromEnv function is used with export mcp-server command
func ExportMCPServerFromEnv(accessToken, name, version, revisionNum, provider, format, exportEnvironment string, preserveStatus,
	exportLatestRevision, preserveCredentials bool) (*resty.Response, error) {
	publisherEndpoint := utils.GetPublisherEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportMCPServer(name, version, revisionNum, provider, format, publisherEndpoint, accessToken, preserveStatus,
		exportLatestRevision, preserveCredentials)
}

// exportMCPServer function is used with export mcp-server command
// @param name : Name of the MCP Server to be exported
// @param version : Version of the MCP Server to be exported
// @param provider : Provider of the MCP Server
// @param publisherEndpoint : API Manager Publisher Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func exportMCPServer(name, version, revisionNum, provider, format, publisherEndpoint, accessToken string, preserveStatus,
	exportLatestRevision, preserveCredentials bool) (*resty.Response, error) {
	publisherEndpoint = utils.AppendSlashToString(publisherEndpoint)
	query := "mcp-servers/export?name=" + url.QueryEscape(name) + "&version=" + version + "&providerName=" + provider +
		"&preserveStatus=" + strconv.FormatBool(preserveStatus) + "&preserveCredentials=" +
		strconv.FormatBool(preserveCredentials)
	if format != "" {
		query += "&format=" + format
	}
	if revisionNum != "" {
		query += "&revisionNumber=" + revisionNum
	}
	if exportLatestRevision {
		query += "&latestRevision=true"
	}

	requestURL := publisherEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportMCPServer: URL:", requestURL)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(requestURL, headers)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
