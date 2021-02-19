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
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// deleteApplication
// @param deleteEndpoint : API Manager Developer Portal REST API Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func DeleteApplication(accessToken, environment, deleteAppName, deleteAppOwner string) (*resty.Response, error) {
	deleteEndpoint := utils.GetAdminApplicationListEndpointOfEnv(environment, utils.MainConfigFilePath)
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)
	appId, err := getAppId(accessToken, environment, deleteAppName, deleteAppOwner)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting App Id for deletion ", err)
	}
	if appId == "" {
		utils.HandleErrorAndExit("Cannot find the application: "+deleteAppName+" for owner: "+deleteAppOwner, err)
	}
	url := deleteEndpoint + appId
	utils.Logln(utils.LogPrefixInfo+"DeleteApplication: URL:", url)
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

func PrintDeleteAppResponse(resp *resty.Response, err error) {
	if err != nil {
		fmt.Println("Error deleting Application:", err)
	} else {
		fmt.Println("Application deleted successfully!. Status: " + strconv.Itoa(resp.StatusCode()))
	}
}

// Get the ID of an Application if available
// @param accessToken : Token to call the Developer Portal Rest API
// @return appId, error
func getAppId(accessToken, environment, appName, appOwner string) (string, error) {
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
