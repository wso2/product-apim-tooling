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

package credentials

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

type MgToken struct {
	// AccessToken of microgateway adapter
	AccessToken string `json:"accessToken"`
}

// GetOAuthAccessTokenForMI returns access token for mi
func GetOAuthAccessTokenForMGAdapter(username, password, tokenEndpoint string) (string, error) {

	b64encodedCredentials := Base64Encode(username + ":" + password)

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials

	resp, err := utils.InvokeGETRequest(tokenEndpoint, headers)
	utils.Logln(utils.LogPrefixInfo + "connecting to " + tokenEndpoint)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("Unable to connect to MI Token endpoint. Status: " + resp.Status())
	}
	responseDataMap := make(map[string]string)
	data := []byte(resp.Body())
	unmarshalError := json.Unmarshal(data, &responseDataMap)

	if unmarshalError != nil {
		return "", unmarshalError
	}
	if accessToken, ok := responseDataMap["AccessToken"]; ok {
		return accessToken, nil
	}
	return "", errors.New("AccessToken not found")
}
