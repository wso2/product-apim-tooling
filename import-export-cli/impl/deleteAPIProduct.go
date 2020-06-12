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

	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"net/http"
)

// DeleteAPIProduct
// @param accessToken : Access Token for the resource
// @param environment : Environment where API Product needs to be located
// @param apiProductName : Name of the API Product
// @param apiProductProvider : Provider of the API Product
// @return response Response in the form of *resty.Response
func DeleteAPIProduct(accessToken, environment, apiProductName, apiProductProvider string) (*resty.Response, error) {
	deleteEndpoint := utils.GetApiProductListEndpointOfEnv(environment, utils.MainConfigFilePath)
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)
	apiProductId, err := getAPIProductId(accessToken, environment, apiProductName, apiProductProvider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Product Id for deletion ", err)
	}
	url := deleteEndpoint + apiProductId
	utils.Logln(utils.LogPrefixInfo+"DeleteAPIProduct: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeDELETERequest(url, headers)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the ID of an API Product if available
// @param accessToken : Access token to call the Publisher Rest API
// @param environment : Environment where API Product needs to be located
// @param apiProductName : Name of the API Product
// @param apiProductProvider : Provider of the API Product
// @return apiId, error
func getAPIProductId(accessToken, environment, apiProductName, apiProductProvider string) (string, error) {
	// Unified Search endpoint from the config file to search API Products
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(environment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryVal string
	// TODO Search by version as well when the versioning support has been implemented for API Products
	queryVal = "type:\"" + utils.DefaultApiProductType + "\" name:\"" + apiProductName + "\""
	if apiProductProvider != "" {
		queryVal = queryVal + " provider:\"" + apiProductProvider + "\""
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
			apiProductId := apiData.List[0].ID
			return apiProductId, err
		}
		// TODO Print the version as well when the versioning support has been implemented for API Products
		if apiProductProvider != "" {
			return "", errors.New("Requested API Product is not available in the Publisher. API Product: " + apiProductName +
				" Provider: " + apiProductProvider)
		}
		return "", errors.New("Requested API Product is not available in the Publisher. API Product: " + apiProductName)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching API Product: " + apiProductName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching API Products. Status: " + resp.Status())
	}
}
