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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

type newRoleRequestBody struct {
	Role   string `json:"role"`
	Domain string `json:"domain"`
}

// AddMIRole adds a new role to the micro integrator in a given environment
func AddMIRole(env, role, domain string) (interface{}, error) {
	body := newRoleRequestBody{
		Role:   role,
		Domain: domain,
	}
	url := utils.GetMIManagementEndpointOfResource(utils.MiManagementRoleResource, env, utils.MainConfigFilePath)
	return addNewMIRole(env, url, body)
}

// DeleteMIRole deletes a role from a micro integrator in a given environment
func DeleteMIRole(env, roleName, domain string) (interface{}, error) {
	url := utils.GetMIManagementEndpointOfResource(utils.MiManagementRoleResource, env, utils.MainConfigFilePath) + "/" + roleName
	params := make(map[string]string)
	putNonEmptyValueToMap(params, "domain", domain)
	return deleteMIRole(url, env, params)
}

func addNewMIRole(env, url string, body interface{}) (string, error) {
	resp, err := invokePOSTRequestWithRetry(env, url, body)
	return handleResponse(resp, err, url, "status", "Error")
}

func deleteMIRole(url, env string, params map[string]string) (string, error) {
	resp, err := invokeDELETERequestWithRetryAndParams(url, env, params)
	return handleResponse(resp, err, url, "status", "Error")
}
