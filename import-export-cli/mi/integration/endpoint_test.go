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
const endpointsCmd = "endpoints"
const endpointCmd = "endpoint"

func TestGetEndpoints(t *testing.T) {
	testutils.ValidateEndpointList(t, endpointsCmd, config)
}

func TestGetEndpointByName(t *testing.T) {
	testutils.ValidateEndpoint(t, endpointsCmd, config, validEndpointName)
}

func TestGetNonExistingEndpointByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, endpointsCmd, invalidEndpointName, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of endpoints [ "+invalidEndpointName+" ]  404 Not Found")
}

func TestGetEndpointsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, endpointsCmd)
}

func TestGetEndpointsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, endpointsCmd, config)
}

func TestGetEndpointsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, endpointsCmd, config)
}

func TestGetEndpointsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, endpointsCmd, validEndpointName, invalidEndpointName)
}

func TestActivateEndpoint(t *testing.T) {
	expected := validEndpointName + " is switched On"
	testutils.ExecActivateCommand(t, config, endpointCmd, validEndpointName, expected)
}

func TestActivateNonExistingEndpoint(t *testing.T) {
	expected := "[ERROR]: Activating endpoint [ " + invalidEndpointName + " ] Endpoint does not exist"
	testutils.ExecActivateCommand(t, config, endpointCmd, invalidEndpointName, expected)
}

func TestActivateEndpointWithoutEnvFlag(t *testing.T) {
	testutils.ExecActivateCommandWithoutEnvFlag(t, config, endpointCmd, validEndpointName)
}

func TestActivateEndpointWithInvalidArgs(t *testing.T) {
	testutils.ExecActivateCommandWithInvalidArgCount(t, config, 1, 0, endpointCmd)
}

func TestActivateEndpointWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecActivateCommandWithoutSettingEnv(t, endpointCmd, validEndpointName)
}

func TestActivateEndpointWithoutLogin(t *testing.T) {
	testutils.ExecActivateCommandWithoutLogin(t, config, endpointCmd, validEndpointName)
}

func TestDeactivateEndpoint(t *testing.T) {
	expected := validEndpointName + " is switched Off"
	testutils.ExecDeactivateCommand(t, config, endpointCmd, validEndpointName, expected)
}

func TestDeactivateNonExistingEndpoint(t *testing.T) {
	expected := "[ERROR]: Deactivating endpoint [ " + invalidEndpointName + " ] Endpoint does not exist"
	testutils.ExecDeactivateCommand(t, config, endpointCmd, invalidEndpointName, expected)
}

func TestDeactivateEndpointWithoutEnvFlag(t *testing.T) {
	testutils.ExecDeactivateCommandWithoutEnvFlag(t, config, endpointCmd, validEndpointName)
}

func TestDeactivateEndpointWithInvalidArgs(t *testing.T) {
	testutils.ExecDeactivateCommandWithInvalidArgCount(t, config, 1, 0, endpointCmd)
}

func TestDeactivateEndpointWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecDeactivateCommandWithoutSettingEnv(t, endpointCmd, validEndpointName)
}

func TestDeactivateEndpointWithoutLogin(t *testing.T) {
	testutils.ExecDeactivateCommandWithoutLogin(t, config, endpointCmd, validEndpointName)
}
