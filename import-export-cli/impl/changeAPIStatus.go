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
	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ChangeAPIStatusInEnv function is used with change-status api command
func ChangeAPIStatusInEnv(accessToken, environment, stateChangeAction, name, version, provider string) (*resty.Response, error) {
	changeAPIStatusEndpoint := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return changeAPIStatus(changeAPIStatusEndpoint, stateChangeAction, name, version, provider, environment, accessToken)
}

// changeAPIStatus
// @param changeAPIStatusEndpoint : API Manager Publisher REST API Endpoint for the environment
// @param stateChangeAction : Action to be performed to change the state of the API
// @param name : Name of the API
// @param version : Version of the API
// @param provider : Provider of the API
// @param environment : Environment where the API resides
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func changeAPIStatus(changeAPIStatusEndpoint, stateChangeAction, name, version, provider, environment, accessToken string) (*resty.Response, error) {
	changeAPIStatusEndpoint = utils.AppendSlashToString(changeAPIStatusEndpoint)
	apiId, err := GetAPIId(accessToken, environment, name, version, provider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id for state change ", err)
	}
	url := changeAPIStatusEndpoint + "change-lifecycle?action=" + stateChangeAction + "&apiId=" + apiId
	utils.Logln(utils.LogPrefixInfo+"APIStateChange: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokePOSTRequest(url, headers, "")

	if err != nil {
		return nil, err
	}
	return resp, nil
}
