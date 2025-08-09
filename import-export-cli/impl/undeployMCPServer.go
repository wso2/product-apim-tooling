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

func UndeployRevisionFromGatewaysMCPServer(accessToken, environment, name, version, provider, revisionNum string,
	gateways []utils.Deployment, allGatewayEnvironments bool) (*resty.Response, error) {

	apiId, err := GetMCPServerId(accessToken, environment, name, version, provider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id from MCP for undeploy", err)
	}
	mcpApiRevisionEndpoint := utils.GetMcpServerListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return undeployRevision(accessToken, mcpApiRevisionEndpoint, apiId, revisionNum, gateways,
		allGatewayEnvironments)
}
