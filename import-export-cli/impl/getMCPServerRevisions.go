/*
*  Copyright (c) 2025 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 LLC. licenses this file to you under the Apache License,
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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// GetRevisionListFromEnv
// @param accessToken	: Access Token for the environment
// @param environment	: Environment name to use when getting the API List
// @param apiName		: Name of the API
// @param apiVersion	: Version of the API
// @param provider		: Provider of the API
// @param query			: Query param for the filtering the revisions based on the deployed status
// @return count (no. of APIs)
// @return array of revision objects
// @return error
func GetMCPServerRevisionListFromEnv(accessToken, environment, apiName, apiVersion, provider, query string) (count int32, revisions []utils.Revisions, err error) {
	mcpServerId, err := GetMCPServerId(accessToken, environment, apiName, apiVersion, provider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting MCP Server Id to list revisions ", err)
	}
	revisionListEndpoint := utils.GetMcpServerListEndpointOfEnv(environment, utils.MainConfigFilePath)
	revisionListEndpoint = utils.AppendSlashToString(revisionListEndpoint)
	url := revisionListEndpoint + mcpServerId + "/revisions"
	if query != "" {
		url += "?query=" + query
	}
	return GetRevisionsList(accessToken, url)
}
