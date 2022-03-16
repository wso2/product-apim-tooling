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

// ChangeAPIStatusInEnv function is used with change-status api command
func ChangeAPIProductStatusInEnv(accessToken, environment, stateChangeAction, name, provider string) (*resty.Response, error) {
	changeAPIProductStatusEndpoint := utils.GetApiProductListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return changeAPIProductStatus(changeAPIProductStatusEndpoint, stateChangeAction, name, provider, environment, accessToken)
}

// changeAPIProductStatus
// @param changeAPIProductStatusEndpoint : API Manager Publisher REST API Endpoint for the environment
// @param stateChangeAction : Action to be performed to change the state of the API Product
// @param name : Name of the API Product
// @param provider : Provider of the API Product
// @param environment : Environment where the API Product resides
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func changeAPIProductStatus(changeAPIProductStatusEndpoint, stateChangeAction, name, provider, environment, accessToken string) (*resty.Response, error) {
	changeAPIProductStatusEndpoint = utils.AppendSlashToString(changeAPIProductStatusEndpoint)
	apiProductId, err := GetAPIProductId(accessToken, environment, name, provider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Product Id for state change ", err)
	}
	url := changeAPIProductStatusEndpoint + "change-lifecycle"
	utils.Logln(utils.LogPrefixInfo+"APIProductStateChange: URL:", url)

	queryParams := make(map[string]string)
	queryParams[utils.LifeCycleAction] = stateChangeAction
	queryParams[utils.APIProductId] = apiProductId

	headers := make(map[string]string)
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokePOSTRequestWithQueryParam(queryParams, url, headers, "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
