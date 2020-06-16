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
	"errors"
	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
)


//GetApplicationListFromEnv
// @param accessToken : Access Token for the environment
// @param environment : Environment to get the list of applications
// @param appOwner : Owner of the applications
// @param limit : Max number of results to return
// @return count (no. of Applications)
// @return array of Application objects
// @return error
func GetApplicationListFromEnv(accessToken, environment, appOwner, limit string) (count int32, apps []utils.Application, err error) {
	applicationListEndpoint := utils.GetAdminApplicationListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return GetApplicationList(accessToken, applicationListEndpoint, appOwner, limit)
}

//Get Application List
// @param accessToken : Access Token for the environment
// @param applicationListEndpoint : Endpoint to use for listing applications
// @param appOwner : Owner of the applications
// @param limit : Max number of results to return
// @return count (no. of Applications)
// @return array of Application objects
// @return error
func GetApplicationList(accessToken, applicationListEndpoint, appOwner, limit string) (count int32, apps []utils.Application,
	err error) {

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	if limit != "" {
		applicationListEndpoint += "?limit=" + limit
	}

	var resp *resty.Response
	if appOwner == "" {
		resp, err = utils.InvokeGETRequest(applicationListEndpoint, headers)
	} else {
		resp, err = utils.InvokeGETRequestWithQueryParam("user", appOwner, applicationListEndpoint, headers)
	}
	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+applicationListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		appListResponse := &utils.ApplicationListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &appListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return appListResponse.Count, appListResponse.List, nil

	} else {
		return 0, nil, errors.New(resp.Status())
	}
}
