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

// ActivateEndpoint activates an endpoint deployed in the micro integrator in a given environment
func ActivateEndpoint(env, endpointName string) (interface{}, error) {
	return updateEndpointState(env, endpointName, "active")
}

// DeactivateEndpoint deactivates an endpoint deployed in the micro integrator in a given environment
func DeactivateEndpoint(env, endpointName string) (interface{}, error) {
	return updateEndpointState(env, endpointName, "inactive")
}

func updateEndpointState(env, endpointName, state string) (interface{}, error) {
	url := utils.GetMIManagementEndpointOfResource(utils.MiManagementEndpointResource, env, utils.MainConfigFilePath)
	resp, err := updateArtifactState(url, endpointName, state, env)
	if err != nil {
		return nil, createErrorWithResponseBody(resp, err)
	}
	return resp, nil
}
