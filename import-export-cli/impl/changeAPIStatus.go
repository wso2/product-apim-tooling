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
	"fmt"
	"net/http"

	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ChangeAPIStatusInEnv function is used with change-status api command
func ChangeAPIStatusInEnv(credential credentials.Credential, accessToken, environment, stateChangeAction, name, version, provider string) (*resty.Response, error) {
	changeAPIStatusEndpoint := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return changeAPIStatus(changeAPIStatusEndpoint, stateChangeAction, name, version, provider, environment, accessToken, credential)
}

// changeAPIStatus
// @param changeAPIStatusEndpoint : API Manager Publisher REST API Endpoint for the environment
// @param stateChangeAction : Action to be performed to change the state of the API
// @param name : Name of the API
// @param version : Version of the API
// @param provider : Provider of the API
// @param environment : Environment where the API resides
// @param accessToken : Access Token for the resource
// @param credential : Credentials of the logged-in user
// @return response Response in the form of *resty.Response
func changeAPIStatus(changeAPIStatusEndpoint, stateChangeAction, name, version, provider, environment, accessToken string, credential credentials.Credential) (*resty.Response, error) {
	changeAPIStatusEndpoint = utils.AppendSlashToString(changeAPIStatusEndpoint)
	apiId, err := getAPIIdForStateChange(accessToken, name, version, provider, environment, credential)
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

// Get the ID of an API if available
// @param accessToken : Token to call the Publisher Rest API
// @param name : Name of the API
// @param version : Version of the API
// @param provider : Provider of the API
// @param environment : Environment where the API resides
// @return apiId, error
func getAPIIdForStateChange(accessToken, name, version, provider, environment string, credential credentials.Credential) (string, error) {
	// Unified Search endpoint from the config file to search APIs
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(environment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryVal string
	queryVal = "name:\"" + name + "\" version:\"" + version + "\""
	if provider != "" {
		queryVal = queryVal + " provider:\"" + provider + "\""
	}
	resp, err := utils.InvokeGETRequestWithQueryParam("query", queryVal, unifiedSearchEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		apiData := &utils.ApiSearch{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &apiData)
		if apiData.Count != 0 {
			apiId := apiData.List[0].ID
			return apiId, err
		}
		if provider != "" {
			return "", errors.New("Requested API is not available in the Publisher. API: " + name +
				" Version: " + version + " Provider: " + provider)
		}
		return "", errors.New("Requested API is not available in the Publisher. API: " + name +
			" Version: " + version)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching API: " + name)
		}
		return "", errors.New("Request didn't respond 200 OK for searching APIs. Status: " + resp.Status())
	}
}
