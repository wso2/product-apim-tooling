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
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ChangeMCPServerStatusInEnv function is used with change-status mcp-server command
func ChangeMCPServerStatusInEnv(accessToken, environment, stateChangeAction, name, version, provider string) (*resty.Response, error) {
	changeMCPServerStatusEndpoint := utils.GetMcpServerListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return changeMCPServerStatus(changeMCPServerStatusEndpoint, stateChangeAction, name, version, provider, environment, accessToken)
}

// changeMCPServerStatus
// @param changeMCPServerStatusEndpoint : API Manager Publisher REST API Endpoint for the environment
// @param stateChangeAction : Action to be performed to change the state of the MCP Server
// @param name : Name of the MCP Server
// @param version : Version of the MCP Server
// @param provider : Provider of the MCP Server
// @param environment : Environment where the MCP Server resides
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func changeMCPServerStatus(changeMCPServerStatusEndpoint, stateChangeAction, name, version, provider, environment, accessToken string) (*resty.Response, error) {
	changeMCPServerStatusEndpoint = utils.AppendSlashToString(changeMCPServerStatusEndpoint)
	mcpServerId, err := GetMCPServerId(accessToken, environment, name, version, provider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting MCP Server Id for state change ", err)
	}
	url := changeMCPServerStatusEndpoint + "change-lifecycle"
	utils.Logln(utils.LogPrefixInfo+"MCPServerStateChange: URL:", url)

	queryParams := make(map[string]string)
	queryParams[utils.LifeCycleAction] = stateChangeAction
	queryParams["mcpServerId"] = mcpServerId

	headers := make(map[string]string)
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokePOSTRequestWithQueryParam(queryParams, url, headers, "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
