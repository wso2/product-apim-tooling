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

const validUserName = "admin"
const invalidUserName = "abc-user"
const userCmd = "users"
const newUserName = "capp-tester"

var validAddUserCmd = []string{"mi", "add", "user", newUserName, "-e", "testing"}

func TestGetUsers(t *testing.T) {
	testutils.ValidateUserList(t, userCmd, config)
}

func TestGetUserByName(t *testing.T) {
	testutils.ValidateUser(t, userCmd, config, validUserName)
}

func TestGetNonExistingUserByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, userCmd, invalidUserName, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of users [ "+invalidUserName+" ]  Requested resource not found. User: "+invalidUserName+" cannot be found.")
}

func TestGetUsersWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, userCmd)
}

func TestGetUsersWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, userCmd, config)
}

func TestGetUsersWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, userCmd, config)
}

func TestGetUsersWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, userCmd, validUserName, invalidUserName)
}

func TestAddNewUserWithoutEnvFlag(t *testing.T) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	response, _ := base.Execute(t, "mi", "add", "user", newUserName)
	base.Log(response)
	expected := `required flag(s) "environment" not set`
	assert.Contains(t, response, expected)
}

func TestAddNewUserWithInvalidArgs(t *testing.T) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	response, _ := base.Execute(t, "mi", "add", "user", "-e", "testing")
	base.Log(response)
	expected := "accepts 1 arg(s), received 0"
	assert.Contains(t, response, expected)
}

func TestAddNewUserWithoutSettingUpEnv(t *testing.T) {
	response, _ := base.Execute(t, validAddUserCmd...)
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, "MI does not exists in testing Add it using add env")
}

func TestAddNewUserWithoutLogin(t *testing.T) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	response, _ := base.Execute(t, validAddUserCmd...)
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, "Login to MI")
}

func TestDeleteUserWithoutEnvFlag(t *testing.T) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	response, _ := base.Execute(t, "mi", "delete", "user", newUserName)
	base.Log(response)
	expected := `required flag(s) "environment" not set`
	assert.Contains(t, response, expected)
}

func TestDeleteUserWithInvalidArgs(t *testing.T) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	response, _ := base.Execute(t, "mi", "delete", "user", "-e", "testing")
	base.Log(response)
	expected := "accepts 1 arg(s), received 0"
	assert.Contains(t, response, expected)
}

func TestDeleteUserWithoutSettingUpEnv(t *testing.T) {
	response, _ := base.Execute(t, validAddUserCmd...)
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, "MI does not exists in testing Add it using add env")
}

func TestDeleteUserWithoutLogin(t *testing.T) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	response, _ := base.Execute(t, validAddUserCmd...)
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, "Login to MI")
}

func TestDeleteUserWithInvalidUserName(t *testing.T) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	response, _ := base.Execute(t, "mi", "delete", "user", invalidUserName, "-e", "testing")
	base.Log(response)
	expected := "[ERROR]: deleting user [ " + invalidUserName + " ] Requested resource not found. User: " + invalidUserName + " cannot be found."
	assert.Contains(t, response, expected)
}

func TestDeleteUser(t *testing.T) {
	status := testutils.AddNewUserFromAPI(config, newUserName, "password", "true")
	assert.Contains(t, status, "Added")
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	response, _ := base.Execute(t, "mi", "delete", "user", newUserName, "-e", "testing")
	base.Log(response)
	expected := "Deleting user [ " + newUserName + " ] status: Deleted"
	assert.Contains(t, response, expected)
}
