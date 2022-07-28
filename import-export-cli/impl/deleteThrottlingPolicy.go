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

// DeleteThrottlingPolicy
// @param accessToken : Access Token for the resource
// @param policyName : Name of the Throttling Policy to delete
// @param policyType : Type of the Throttling Policy to delete
// @param environment : Environment where Throttling should be deleted
// @return response Response in the form of *resty.Response
func DeleteThrottlingPolicy(accessToken, policyName, policyType, environment string) (*resty.Response, error) {
	endpoint := utils.GetAdminEndpointOfEnv(environment, utils.MainConfigFilePath)
	endpoint = utils.AppendSlashToString(endpoint)

	resource := "throttling/policies/search"
	searchEndpoint := endpoint + resource

	queryParam := `?name=` + policyName

	url := searchEndpoint + queryParam

	policyId, err := getThrottlingPolicyId(accessToken, environment, url, policyType, policyName)

	if policyId == "" && err == nil {
		errMsg := "Requested Policy with name=" + policyName + " and type=" + policyType + " no found."
		utils.HandleErrorAndExit("Deletion Failed ! ", errors.New(errMsg))
	}

	if err != nil {
		utils.HandleErrorAndExit("Error while getting Throttling Policy Id for deletion ", err)
	}

	resource = "throttling/policies/"

	switch policyType {
	case CmdPolicyTypeSubscription:
		resource += utils.ThrottlingPolicyTypeSub
	case CmdPolicyTypeApplication:
		resource += utils.ThrottlingPolicyTypeApp
	case CmdPolicyTypeAdvanced:
		resource += utils.ThrottlingPolicyTypeAdv
	case CmdPolicyTypeCustom:
		resource += utils.ThrottlingPolicyTypeCus
	}

	resource = utils.AppendSlashToString(resource)

	resource += policyId

	url = endpoint + resource
	fmt.Println(url)
	utils.Logln(utils.LogPrefixInfo+"DeleteThrottlingPolicy: URL:", url)
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

func getThrottlingPolicyId(accessToken, environment, url, policyType, policyName string) (string, error) {
	utils.Logln(utils.LogPrefixInfo+"DeleteThrottlingPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)

	if err != nil {
		return "", err
	}

	var policyList utils.ThrottlingPoliciesDetailsList
	err = json.Unmarshal(resp.Body(), &policyList)

	if err != nil {
		return "", err
	}

	for _, obj := range policyList.List {
		if obj.Type == policyType && obj.PolicyName == policyName {
			return obj.Uuid, nil
		}
	}

	return "", nil
}

func PrintDeleteThrottlingPolicyResponse(policyName, policyType string, err error) {
	if err != nil {
		fmt.Println("Error deleting Throttling Policy:", err)
	} else {
		fmt.Println(policyName + " Throttling Policy with type " + policyType + " deleted successfully!")
	}
}
