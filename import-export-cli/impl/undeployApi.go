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
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func UndeployRevisionFromGateways(accessToken, environment, name, version, provider, revisionNum string,
	gateways []utils.Deployment, allGatewayEnvironments bool) (*resty.Response, error) {

	apiId, err := GetAPIId(accessToken, environment, name, version, provider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id for undeploy", err)
	}
	apiRevisionEndpoint := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return undeployRevision(accessToken, apiRevisionEndpoint, apiId, revisionNum, gateways,
		allGatewayEnvironments)
}

// Function is used with undeploy revision command
// @param accessToken : Access Token for the resource
// @param undeployRevisionEndpoint : API resource to undeploy the revisions
// @param apiId : API ID
// @param revisionNum : Revision number of the API
// @param gateways : Gateway environments in which the revision has to be deployed
// @param allGatewayEnvironments : Boolean to specify whether to undeploy in all gateways
// @return response Response in the form of *resty.Response
func undeployRevision(accessToken, undeployRevisionEndpoint, apiId, revisionNum string,
	gateways []utils.Deployment, allGatewayEnvironments bool) (*resty.Response, error) {
	undeployRevisionEndpoint = utils.AppendSlashToString(undeployRevisionEndpoint) + apiId +
		"/undeploy-revision?revisionNumber=" + revisionNum
	if allGatewayEnvironments {
		//This is used to undeploy all the environments at once for a specific revision
		undeployRevisionEndpoint += "&allEnvironments=true"
	}

	utils.Logln(utils.LogPrefixInfo+"Undeploy URL:", undeployRevisionEndpoint)

	headers := make(map[string]string)
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	body, err := json.Marshal(gateways)
	if err != nil {
		utils.HandleErrorAndExit("Error while converting gateways array", err)
	}

	return utils.InvokePOSTRequest(undeployRevisionEndpoint, headers, string(body))
}
