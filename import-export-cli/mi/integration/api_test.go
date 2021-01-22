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

const validAPIName = "HealthcareAPI"
const invalidAPIName = "abcAPI"
const apisCmd = "apis"

func TestGetAPIs(t *testing.T) {
	testutils.ValidateAPIsList(t, apisCmd, config)
}

func TestGetAPIByName(t *testing.T) {
	testutils.ValidateAPI(t, apisCmd, config, validAPIName)
}

func TestGetNonExistingAPIByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, apisCmd, invalidAPIName, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of "+apisCmd+" [ "+invalidAPIName+" ]  404 Not Found")
}

func TestGetAPIsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, apisCmd)
}

func TestGetAPIsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, apisCmd, config)
}

func TestGetAPIsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, apisCmd, config)
}

func TestGetAPIsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, apisCmd, validAPIName, invalidAPIName)
}
