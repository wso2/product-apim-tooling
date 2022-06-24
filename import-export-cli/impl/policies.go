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

// GetAPIId Get the ID of an API if available
// @param accessToken : Token to call the Publisher Rest API
// @param environment : Environment where API needs to be located
// @param apiName : Name of the API
// @param apiVersion : Version of the API
// @param apiProvider : Provider of API
// @return apiId, error
func GetOperationPolicyId(accessToken, environment, policyName, policyVersion string) (string, error) {
	// Unified Search endpoint from the config file to search APIs
	operationPolicyEndpoint := utils.GetPublisherEndpointOfEnv(environment, utils.MainConfigFilePath)
	fmt.Println("Get End: ", operationPolicyEndpoint)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	queryParams := `name=` + policyName + `&version=` + policyVersion

	operationPolicyEndpoint = utils.AppendSlashToString(operationPolicyEndpoint)
	// operationPolicyResource := "operation-policies/c86da87e-da70-4977-bed2-57cb089c115f" + "/content"
	operationPolicyResource := "operation-policies?" + queryParams

	url := operationPolicyEndpoint + operationPolicyResource
	utils.Logln(utils.LogPrefixInfo+"DeleteOperationPolicy: URL:", url)

	fmt.Println("URL: ", url)

	resp, err := utils.InvokeGETRequest(url, headers)
	if err != nil {
		return "", err
	}

	fmt.Println("Resp: ", resp)

	if resp.StatusCode() == http.StatusOK {
		policyData := &utils.OperationPoliciesList{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &policyData)
		if policyData.List[0].Id != "" {
			return policyData.List[0].Id, err
		}

		return "", errors.New("Requested Operation Policy is not available in the Publisher. Policy: " + policyName +
			" Version: " + policyVersion)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while getting Operation Policy Id: " + policyName)
		}
		return "", errors.New("Request didn't respond 200 OK for getting Operation Policy Id. Status: " + resp.Status())
	}

}
