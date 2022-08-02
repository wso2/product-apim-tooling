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
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// DeleteAPIPolicy
// @param accessToken : Access Token for the resource
// @param policyName : Name of the API Policy to delete
// @param policyVersion : Version of the API Policy to delete
// @param environment : Environment where API Policy should be deleted
// @return response Response in the form of *resty.Response
func DeleteAPIPolicy(accessToken, policyName, policyVersion, environment string) (*resty.Response, error) {
	deleteEndpoint := utils.GetAPIPolicyListEndpointOfEnv(environment, utils.MainConfigFilePath)
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)

	policyId, err := GetAPIPolicyId(accessToken, environment, policyName, policyVersion)

	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Policy Id for deletion ", err)
	}
	url := deleteEndpoint + policyId
	utils.Logln(utils.LogPrefixInfo+"DeleteAPIPolicy: URL:", url)
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

func PrintDeleteAPIPolicyResponse(policyName string, err error) {
	if err != nil {
		fmt.Println("Error deleting API Policy:", err)
	} else {
		fmt.Println(policyName + " API Policy deleted successfully!")
	}
}
