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
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"testing"
)
const defaultExportPath = utils.DefaultExportDirName


//List Environments using apictl
func TestListEnvironments (t *testing.T){
	apim := apimClients[0]
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	response, _ := base.Execute(t, "list", "envs")
	base.GetRowsFromTableResponse(response)
	assert.Contains(t,response,apim.GetEnvName(),"TestListEnvironments Failed")
}

//Change Export directory using apictl and assert the change
func TestChangeExportDirectory(t *testing.T) {
	apim := apimClients[0]
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	var changedExportDirectory = "/testingPath/" + utils.DefaultExportDirName
	output, _ := base.Execute(t, "set",  "---export-directory ", changedExportDirectory, "-k", "--verbose")
	base.Log(output)
	assert.Equal(t,changedExportDirectory,utils.DefaultExportDirPath,"Export Directory change is not successful")
}

//Change HTTP request Timeout using apictl and assert the change
func TestChangeHttpRequestTimout(t *testing.T) {
	apim := apimClients[0]
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	var newHttpRequestTimeOut = 20000
	output, _ := base.Execute(t, "set",  "--http-request-timeout ", string(newHttpRequestTimeOut), "-k", "--verbose")
	base.Log(output)
	assert.Equal(t,newHttpRequestTimeOut,configVars.Config.HttpRequestTimeout,"HTTP Request TimeOut change is not successful")
}

//Change Token type using apictl and assert the change (for both "jwt" and "oauth" token types)
func TestChangeTokenType(t *testing.T) {
	apim := apimClients[0]
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	var tokenType1 = "Oauth"
	output1, _ := base.Execute(t, "set",  "--token-type",tokenType1 , "-k", "--verbose")
	base.Log(output1)
	assert.Equal(t,tokenType1,configVars.Config.TokenType,"1st attempt of Token Type change is not successful")
	var tokenType2 = "JWT"
	output2, _ := base.Execute(t, "set",  "--token-type",tokenType2 , "-k", "--verbose")
	base.Log(output2)
	assert.Equal(t,tokenType2,configVars.Config.TokenType,"2nd attempt of Token Type change is not successful")
}

//Login to the environment using email and logout
func TestLoginWithEmail (t *testing.T){

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword
	dev := apimClients[0]

	args := &loginTestArgs{
		ctlUser:     credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		srcAPIM:  dev,
	}
	// Setup apictl envs
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())
	base.Execute(t, "login", args.srcAPIM.GetEnvName(), "-u", tenantAdminUsername, "-p", tenantAdminPassword, "-k", "--verbose")

	t.Cleanup(func() {
		base.Execute(t, "logout", args.srcAPIM.GetEnvName())
	})
}
