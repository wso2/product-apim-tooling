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
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"strconv"
)

// DeleteAPI
// @param accessToken : Access Token for the resource
// @param environment : Environment where API should be deleted
// @param deleteAPIName : Name of the API to delete
// @param deleteAPIVersion : Version of the API to delete
// @param deleteAPIProvider : Provider of API
// @return response Response in the form of *resty.Response
func DeleteAPI(accessToken, environment, deleteAPIName, deleteAPIVersion, deleteAPIProvider string) (*resty.Response, error) {
	deleteEndpoint := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)
	apiId, err := getAPIId(accessToken, environment, deleteAPIName, deleteAPIVersion, deleteAPIProvider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id for deletion ", err)
	}
	url := deleteEndpoint + apiId
	utils.Logln(utils.LogPrefixInfo+"DeleteAPI: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeDELETERequest(url, headers)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return nil, errors.New(strconv.Itoa(resp.StatusCode()) + ":<" + string(resp.Body()) + ">")
	}
	return resp, nil
}

func PrintDeleteAPIResponse(resp *resty.Response, err error) {
	if err != nil {
		fmt.Println("Error deleting API:", err)
	} else {
		fmt.Println("API deleted successfully!. Status: " + strconv.Itoa(resp.StatusCode()))
	}
}

// Get the ID of an API if available
// @param accessToken : Token to call the Publisher Rest API
// @param environment : Environment where API needs to be located
// @param apiName : Name of the API
// @param apiVersion : Version of the API
// @param apiProvider : Provider of API
// @return apiId, error
func getAPIId(accessToken, environment, apiName, apiVersion, apiProvider string) (string, error) {
	// Unified Search endpoint from the config file to search APIs
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(environment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryVal string
	queryVal = "name:\"" + apiName + "\" version:\"" + apiVersion + "\""
	if apiProvider != "" {
		queryVal = queryVal + " provider:\"" + apiProvider + "\""
	}
	resp, err := utils.InvokeGETRequestWithQueryParam("query", queryVal, unifiedSearchEndpoint, headers)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		apiData := &utils.ApiSearch{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &apiData)
		if apiData.Count != 0 {
			apiId := apiData.List[0].ID
			return apiId, err
		}
		if apiProvider != "" {
			return "", errors.New("Requested API is not available in the Publisher. API: " + apiName +
				" Version: " + apiVersion + " Provider: " + apiProvider)
		}
		return "", errors.New("Requested API is not available in the Publisher. API: " + apiName +
			" Version: " + apiVersion)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching API: " + apiName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching APIs. Status: " + resp.Status())
	}
}
