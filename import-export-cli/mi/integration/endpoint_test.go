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

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/integration/testutils"
)

const validEndpointName = "GrandOakEndpoint"
const invalidEndpointName = "abcEndpoint"
const endpointCmd = "endpoints"

func TestGetEndpoints(t *testing.T) {
	testutils.ValidateEndpointList(t, endpointCmd, config)
}

func TestGetEndpointByName(t *testing.T) {
	testutils.ValidateEndpoint(t, endpointCmd, config, validEndpointName)
}

func TestGetNonExistingEndpointByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, endpointCmd, invalidEndpointName, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of endpoints [ "+invalidEndpointName+" ]  404 Not Found")
}

func TestGetEndpointsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, endpointCmd)
}

func TestGetEndpointsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, endpointCmd, config)
}

func TestGetEndpointsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, endpointCmd, config)
}

func TestGetEndpointsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, endpointCmd, validEndpointName, invalidEndpointName)
}
