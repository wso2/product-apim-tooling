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
	"path/filepath"

	"github.com/go-resty/resty"
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// GetAppId Get the ID of an Application if available
// @param accessToken : Token to call the Developer Portal Rest API
// @return appId, error
func GetAppId(accessToken, environment, appName, appOwner string) (string, error) {
	// Application REST API endpoint of the environment from the config file
	applicationEndpoint := utils.GetAdminApplicationListEndpointOfEnv(environment, utils.MainConfigFilePath) +
		"?user=" + appOwner + "&name=" + appName

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(applicationEndpoint, headers)

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		appData := &utils.AppList{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &appData)
		appId := ""
		if appData.Count != 0 {
			for _, app := range appData.List {
				if app.Name == appName {
					appId = app.ApplicationID
				}
			}
			return appId, err
		}
		return "", errors.New("Cannot find the application: " + appName + " for owner: " + appOwner)

	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching CLI application: " + appName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching existing applications. " +
			"Status: " + resp.Status())
	}
}

// GetApplicationDefinition scans filePath and returns ApplicationDefinition or an error
func GetApplicationDefinition(filePath string) (*v2.ApplicationDefinition, []byte, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}

	var buffer []byte
	if info.IsDir() {
		_, content, err := resolveYamlOrJSON(path.Join(filePath, filepath.Base(filePath)))
		if err != nil {
			return nil, nil, err
		}
		buffer = content
	} else {
		return nil, nil, fmt.Errorf("looking for directory, found %s", info.Name())
	}
	api, err := extractAppDefinition(buffer)
	if err != nil {
		return nil, nil, err
	}
	return api, buffer, nil
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

	utils.Logln(utils.LogPrefixInfo+"URL:", applicationListEndpoint)

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
