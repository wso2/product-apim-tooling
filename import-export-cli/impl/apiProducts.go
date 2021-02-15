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
	"os"
	"path"

	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// GetAPIProductId Get the ID of an API Product if available
// @param accessToken : Access token to call the Publisher Rest API
// @param environment : Environment where API Product needs to be located
// @param apiProductName : Name of the API Product
// @param apiProductProvider : Provider of the API Product
// @return apiId, error
func GetAPIProductId(accessToken, environment, apiProductName, apiProductProvider string) (string, error) {
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

// GetAPIProductDefinition scans filePath and returns APIProductDefinition or an error
func GetAPIProductDefinition(filePath string) (*v2.APIProductDefinition, []byte, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}

	var buffer []byte
	if info.IsDir() {
		_, content, err := resolveYamlOrJSON(path.Join(filePath, "Meta-information", "api"))
		if err != nil {
			return nil, nil, err
		}
		buffer = content
	} else {
		return nil, nil, fmt.Errorf("looking for directory, found %s", info.Name())
	}
	apiProduct, err := extractAPIProductDefinition(buffer)
	if err != nil {
		return nil, nil, err
	}
	return apiProduct, buffer, nil
}

// GetAPIProductList Get the list of API Products available in a particular environment
// @param accessToken : Access Token for the environment
// @param unifiedSearchEndpoint : Unified Search Endpoint for the environment to retreive API Product list
// @param query : String to be matched against the API Product names
// @return count (no. of API Products)
// @return array of API Product objects
// @return error
func GetAPIProductList(accessToken, unifiedSearchEndpoint, query, limit string) (count int32, apiProducts []utils.APIProduct, err error) {
	// Unified Search endpoint from the config file to search API Products
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	// To filter API Products from unified search
	unifiedSearchEndpoint += "?query=type:\"" + utils.DefaultApiProductType + "\""

	// Setting up the query parameter and limit parameter
	if query != "" {
		unifiedSearchEndpoint += " " + query
	}
	if limit != "" {
		unifiedSearchEndpoint += "&limit=" + limit
	}
	utils.Logln(utils.LogPrefixInfo+"URL:", unifiedSearchEndpoint)
	resp, err := utils.InvokeGETRequest(unifiedSearchEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+unifiedSearchEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiProductListResponse := &utils.APIProductListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiProductListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return apiProductListResponse.Count, apiProductListResponse.List, nil
	} else {
		return 0, nil, errors.New(string(resp.Body()))
	}
}

// GetRevisionsList Get the list of Revisions available for the given API
// @param accessToken 			: Access Token for the environment
// @param revisionListEndpoint 	: Revision List endpoint
// @return count (no. of revisions)
// @return array of revision objects
// @return error
func GetAPIProductRevisionsList(accessToken, revisionListEndpoint string) (count int32, revisions []utils.Revisions,
	err error) {

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	utils.Logln(utils.LogPrefixInfo+"URL:", revisionListEndpoint)
	resp, err := utils.InvokeGETRequest(revisionListEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+revisionListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		revisionListResponse := &utils.RevisionListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &revisionListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return revisionListResponse.Count, revisionListResponse.List, nil
	} else {
		return 0, nil, errors.New(string(resp.Body()))
	}
}
