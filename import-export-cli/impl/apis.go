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

// GetAPIId Get the ID of an API if available
// @param accessToken : Token to call the Publisher Rest API
// @param environment : Environment where API needs to be located
// @param apiName : Name of the API
// @param apiVersion : Version of the API
// @param apiProvider : Provider of API
// @return apiId, error
func GetAPIId(accessToken, environment, apiName, apiVersion, apiProvider string) (string, error) {
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

// GetAPIDefinition scans filePath and returns APIDefinition or an error
func GetAPIDefinition(filePath string) (*v2.APIDefinition, []byte, error) {
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
	api, err := extractAPIDefinition(buffer)
	if err != nil {
		return nil, nil, err
	}
	return api, buffer, nil
}

// GetAPIList Get the list of APIs available in a particular environment
// @param accessToken : Access Token for the environment
// @param apiListEndpoint : API List endpoint
// @param query : string to be matched against the API names
// @param limit : total # of results to return
// @return count (no. of APIs)
// @return array of API objects
// @return error
func GetAPIList(accessToken, apiListEndpoint, query, limit string) (count int32, apis []utils.API, err error) {
	queryParamAdded := false
	getQueryParamConnector := func() (connector string) {
		if queryParamAdded {
			return "&"
		} else {
			return ""
		}
	}

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryParamSring string
	if query != "" {
		queryParamSring = "query=" + query
		queryParamAdded = true
	}
	if limit != "" {
		queryParamSring += getQueryParamConnector() + "limit=" + limit
	}
	utils.Logln(utils.LogPrefixInfo+"URL:", apiListEndpoint + "?" + queryParamSring)
	resp, err := utils.InvokeGETRequestWithQueryParamsString(apiListEndpoint, queryParamSring, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiListResponse := &utils.APIListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return apiListResponse.Count, apiListResponse.List, nil
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
func GetRevisionsList(accessToken, revisionListEndpoint string) (count int32, revisions []utils.Revisions, err error) {

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
