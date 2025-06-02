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

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// Calling token endpoint to get access token for the AI related commands
// @param key : Base64 encoded client_id:client_secret
// @return accessToken, error
func GetAIToken(key, env string) (string, error) {

	tokenEndpoint := utils.GetAITokenServiceEndpointOfEnv(env, utils.MainConfigFilePath)
	body := "grant_type=client_credentials"

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + key
	headers[utils.HeaderContentType] = utils.HeaderValueXWWWFormUrlEncoded
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON

	resp, err := utils.InvokePOSTRequest(tokenEndpoint, headers, body)

	if err != nil {
		return "", errors.New("AI Token Endpoint is not valid. " + err.Error())
	}

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		keygenResponse := &utils.TokenResponse{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &keygenResponse)

		accessToken := keygenResponse.AccessToken
		return accessToken, err
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("authorization failed while generating a token for the AI operations.")
		}
		return "", errors.New("Request didn't respond 200 OK for generating a new token for AI operations. Status: " + resp.Status())
	}
}
