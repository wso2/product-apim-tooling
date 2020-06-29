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
	base.Log(response)
	assert.Contains(t,response,apim.GetEnvName(),"TestListEnvironments Failed")
}

//Change Export directory using apictl and assert the change
func TestChangeExportDirectory(t *testing.T) {
	apim := apimClients[0]
	changedExportDirectory := utils.MockTestExportDirectory + utils.DefaultExportDirName
	defaultExportPath := utils.DefaultExportDirPath

	args := &setTestArgs{
		srcAPIM: apim,
		exportDirectoryFlag: changedExportDirectory,
	}
	output, _ := environmentSetExportDirectory(t, args)
	base.Log(output)

	//Change value back to default value
	argsDefault := &setTestArgs{
		srcAPIM: apim,
		exportDirectoryFlag: defaultExportPath,
	}
	environmentSetExportDirectory(t, argsDefault)
	assert.Contains(t,output,"Token type set to:  JWT","Export Directory change is not successful")
}

//Change HTTP request Timeout using apictl and assert the change
func TestChangeHttpRequestTimout(t *testing.T) {
	apim := apimClients[0]
	defaultHttpRequestTimeOut:=utils.DefaultHttpRequestTimeout
	newHttpRequestTimeOut := 20000
	args := &setTestArgs{
		srcAPIM: apim,
		httpRequestTimeout: newHttpRequestTimeOut,
	}
	output, _ := environmentSetHttpRequestTimeout(t, args)
	base.Log(output)

	//Change value back to default value
	argsDefault := &setTestArgs{
		srcAPIM: apim,
		httpRequestTimeout: defaultHttpRequestTimeOut,
	}
	environmentSetHttpRequestTimeout(t, argsDefault)
	assert.Contains(t,output,"Token type set to:  JWT","HTTP Request TimeOut change is not successful")
}

//Change Token type using apictl and assert the change (for both "jwt" and "oauth" token types)
func TestChangeTokenType(t *testing.T) {
	apim := apimClients[0]

	tokenType1 := "Oauth"
	args := &setTestArgs{
		srcAPIM: apim,
		tokenTypeFlag: tokenType1,
	}
	output, _ := environmentSetHttpRequestTimeout(t, args)
	base.Log(output)
	assert.Contains(t,output,"Token type set to:  JWT","1st attempt of Token Type change is not successful")
	tokenType2 := "JWT"

	//Change value back to default value with a test
	argsDefault := &setTestArgs{
		srcAPIM: apim,
		tokenTypeFlag: tokenType2,
	}
	output2, _ := environmentSetHttpRequestTimeout(t, argsDefault)
	base.Log(output2)
	assert.Contains(t,output2,"Token type set to:  JWT","1st attempt of Token Type change is not successful")
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
