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

const validCAppName = "HealthCareCompositeExporter"
const invalidCAppName = "abcCApp"
const cAppCmd = "composite-apps"

func TestGetCApps(t *testing.T) {
	testutils.ValidateCAppsList(t, cAppCmd, config)
}

func TestGetCAppByName(t *testing.T) {
	testutils.ValidateCApp(t, cAppCmd, config, validCAppName)
}

func TestGetNonExistingCAppByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, cAppCmd, invalidCAppName, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of composite apps [ "+invalidCAppName+" ]  404 Not Found")
}

func TestGetCAppsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, cAppCmd)
}

func TestGetCAppsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, cAppCmd, config)
}

func TestGetCAppsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, cAppCmd, config)
}

func TestGetCAppsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgs(t, config, 1, 2, cAppCmd, validCAppName, invalidCAppName)
}
