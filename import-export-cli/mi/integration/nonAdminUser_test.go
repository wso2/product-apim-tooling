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

const nonAdminUserName = "user"
const nonAdminUserPassword = "password"
const isAdmin = "false"

func TestGetUsersFromNonAdminUser(t *testing.T) {
	testutils.AddNewUserFromAPI(t, config, nonAdminUserName, nonAdminUserPassword, isAdmin, true)
	testutils.SetupAndLoginToMI(t, nonAdminConfig)
	response, _ := base.Execute(t, "mi", "get", "users", "-e", "testing")
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting List of users 403 Forbidden")
}

func TestGetUserByNameFromNonAdminUser(t *testing.T) {
	testutils.AddNewUserFromAPI(t, config, nonAdminUserName, nonAdminUserPassword, isAdmin, true)
	testutils.SetupAndLoginToMI(t, nonAdminConfig)
	response, _ := base.Execute(t, "mi", "get", "users", validUserName, "-e", "testing")
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of users [ "+validUserName+" ]  403 Forbidden")
}

func TestGetNonExistingUserByNameFromNonAdminUser(t *testing.T) {
	testutils.AddNewUserFromAPI(t, config, nonAdminUserName, nonAdminUserPassword, isAdmin, true)
	testutils.SetupAndLoginToMI(t, nonAdminConfig)
	response, _ := base.Execute(t, "mi", "get", "users", invalidUserName, "-e", "testing")
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of users [ "+invalidUserName+" ]  403 Forbidden")
}

func TestDeleteUserWithInvalidUserNameFromNonAdminUser(t *testing.T) {
	testutils.AddNewUserFromAPI(t, config, nonAdminUserName, nonAdminUserPassword, isAdmin, true)
	testutils.SetupAndLoginToMI(t, nonAdminConfig)
	response, _ := base.Execute(t, "mi", "delete", "user", invalidUserName, "-e", "testing")
	base.Log(response)
	expected := "[ERROR]: deleting user [ " + invalidUserName + " ] 403 Forbidden"
	assert.Contains(t, response, expected)
}

func TestDeleteUserFromNonAdminUser(t *testing.T) {
	testutils.AddNewUserFromAPI(t, config, nonAdminUserName, nonAdminUserPassword, isAdmin, true)
	testutils.SetupAndLoginToMI(t, nonAdminConfig)
	response, _ := base.Execute(t, "mi", "delete", "user", validUserName, "-e", "testing")
	base.Log(response)
	expected := "[ERROR]: deleting user [ " + validUserName + " ] 403 Forbidden"
	assert.Contains(t, response, expected)
}

func TestGetAPIsFromNonAdminUser(t *testing.T) {
	testutils.AddNewUserFromAPI(t, config, nonAdminUserName, nonAdminUserPassword, isAdmin, true)
	testutils.ValidateAPIsList(t, apisCmd, nonAdminConfig)
}

func TestGetAPIByNameFromNonAdminUser(t *testing.T) {
	testutils.AddNewUserFromAPI(t, config, nonAdminUserName, nonAdminUserPassword, isAdmin, true)
	testutils.ValidateAPI(t, apisCmd, nonAdminConfig, validAPIName)
}
