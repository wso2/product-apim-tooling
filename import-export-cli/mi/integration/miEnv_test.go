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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/integration/testutils"
)

const undefinedEnv = "env-abc"

func TestAddEnvironmentWithMI(t *testing.T) {
	response, _ := base.Execute(t, "add", "env", miClient.GetEnvName(), "--mi", miClient.GetMiURL())
	base.Log(response)
	expected := fmt.Sprintf("Successfully added environment '%v'", miClient.GetEnvName())
	assert.Contains(t, response, expected)
}

func TestLoginToMI(t *testing.T) {
	response, _ := base.Execute(t, "mi", "login", miClient.GetEnvName(), "-u", testutils.AdminUserName, "-p", testutils.AdminPassword)
	base.Log(response)
	expected := fmt.Sprintf("Logged into MI in %v environment", miClient.GetEnvName())
	assert.Contains(t, response, expected)
}

func TestLogoutFromMI(t *testing.T) {
	response, _ := base.Execute(t, "mi", "logout", miClient.GetEnvName())
	base.Log(response)
	expected := fmt.Sprintf("Logged out from MI in %v environment", miClient.GetEnvName())
	assert.Contains(t, response, expected)
}

func TestLoginToMIWithInvalidCredentials(t *testing.T) {
	response, _ := base.Execute(t, "mi", "login", miClient.GetEnvName(), "-u", testutils.AdminUserName, "-p", "abc123")
	base.Log(response)
	assert.Contains(t, response, "Error occurred while login :  Unable to connect to MI Token endpoint. Status: 401 Unauthorized")
}

func TestRemoveEnvironmentWithMI(t *testing.T) {
	response, _ := base.Execute(t, "remove", "env", miClient.GetEnvName())
	base.Log(response)
	expected := fmt.Sprintf("Successfully removed environment '%v'", miClient.GetEnvName())
	assert.Contains(t, response, expected)
}

func TestLoginToMIInUndefinedEnv(t *testing.T) {
	response, _ := base.Execute(t, "mi", "login", undefinedEnv, "-u", testutils.AdminUserName, "-p", testutils.AdminPassword)
	base.Log(response)
	expected := fmt.Sprintf("MI does not exists in %v Add it using add env", undefinedEnv)
	assert.Contains(t, response, expected)
}

func TestLoginToMIWithoutEnv(t *testing.T) {
	response, _ := base.Execute(t, "mi", "login")
	base.Log(response)
	assert.Contains(t, response, "accepts 1 arg(s), received 0")
}

func TestLoginToMIWithInvalidArgs(t *testing.T) {
	response, _ := base.Execute(t, "mi", "login", "testing", "dev")
	base.Log(response)
	assert.Contains(t, response, "accepts 1 arg(s), received 2")
}
