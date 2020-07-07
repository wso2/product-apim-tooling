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
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const defaultExportPath = utils.DefaultExportDirName

//List Environments using apictl
func TestListEnvironments(t *testing.T) {
	apim := apimClients[0]
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	response, _ := base.Execute(t, "list", "envs")
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, apim.GetEnvName(), "TestListEnvironments Failed")
}

func validateExportDirectoryIsChanged(t *testing.T, args *testutils.SetTestArgs) {
	t.Helper()
	output, _ := testutils.EnvironmentSetExportDirectory(t, args)
	base.Log(output)
	assert.Contains(t, output, "Export Directory is set to", "Export Directory change is not successful")
}

func validateExportApisPassed(t *testing.T, args *testutils.InitTestArgs, directoryName string) {
	t.Helper()
	time.Sleep(5 * time.Second)

	output, error := testutils.ExportApisWithOneCommand(t, args)
	assert.Nil(t, error, "Error while Exporting APIs")
	assert.Contains(t, output, "export-apis execution completed", "Error while Exporting APIs")

	//Derive exported path from output
	exportedPath := base.GetExportedPathFromOutput(strings.ReplaceAll(output, "Command: export-apis execution completed !", ""))
	count, _ := base.CountFiles(exportedPath)
	assert.Equal(t, 1, count, "Error while Exporting APIs")

	t.Cleanup(func() {
		argsDefault := &testutils.SetTestArgs{
			SrcAPIM:             args.SrcAPIM,
			ExportDirectoryFlag: utils.DefaultExportDirPath,
		}
		validateExportDirectoryIsChanged(t, argsDefault)
		//Remove Exported apis
		base.RemoveDir(directoryName + utils.TestMigrationDirectorySuffix)
	})
}

//Change Export directory using apictl and assert the change
func TestChangeExportDirectory(t *testing.T) {
	dev := apimClients[0]
	changedExportDirectory, _ := filepath.Abs(utils.CustomTestExportDirectory + utils.DefaultExportDirName)

	args := &testutils.SetTestArgs{
		SrcAPIM:             dev,
		ExportDirectoryFlag: changedExportDirectory,
	}
	validateExportDirectoryIsChanged(t, args)

	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	apiArgs := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestOpenAPI3DefinitionPath,
		ForceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	validateImportInitializedProject(t, apiArgs)

	//Assert that Export directory change is successful by exporting and asserting that
	validateExportApisPassed(t, apiArgs, changedExportDirectory)

}

//TODO  - Need to come up with  a process to make sure that http timeout is actually changed using another fake server
//Change HTTP request Timeout using apictl and assert the change
//func TestChangeHttpRequestTimout(t *testing.T) {
//	apim := apimClients[0]
//	defaultHttpRequestTimeOut := utils.DefaultHttpRequestTimeout
//	newHttpRequestTimeOut := 20000
//	args := &setTestArgs{
//		srcAPIM:            apim,
//		httpRequestTimeout: newHttpRequestTimeOut,
//	}
//	output, _ := environmentSetHttpRequestTimeout(t, args)
//	base.Log(output)
//
//	//Change value back to default value
//	argsDefault := &setTestArgs{
//		srcAPIM:            apim,
//		httpRequestTimeout: defaultHttpRequestTimeOut,
//	}
//	environmentSetHttpRequestTimeout(t, argsDefault)
//	assert.Contains(t, output, "Http Request Timout is set to", "HTTP Request TimeOut change is not successful")
//}

func validateETokenTypeIsChanged(t *testing.T, args *testutils.SetTestArgs) {
	t.Helper()
	output, _ := testutils.EnvironmentSetTokenType(t, args)
	base.Log(output)
	assert.Contains(t, output, "Token type is set to", "1st attempt of Token Type change is not successful")
}

//Change Token type using apictl and assert the change (for both "jwt" and "oauth" token types)
func TestChangeTokenType(t *testing.T) {
	apim := apimClients[0]

	tokenType1 := "oauth"
	args := &testutils.SetTestArgs{
		SrcAPIM:       apim,
		TokenTypeFlag: tokenType1,
	}

	validateETokenTypeIsChanged(t, args)

	//Create API and get keys for that API
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, adminUser, adminPassword)

	testutils.PublishAPI(dev, adminUser, adminPassword, api.ID)

	apiArgs := &testutils.ApiGetKeyTestArgs{
		CtlUser: testutils.Credentials{Username: adminUser, Password: adminPassword},
		Api:     api,
		Apim:    dev,
	}

	validateThatRecievingTokenTypeIsChanged(t, apiArgs, tokenType1)

	tokenType2 := "jwt"

	//Change value back to default value
	argsDefault := &testutils.SetTestArgs{
		SrcAPIM:       apim,
		TokenTypeFlag: tokenType2,
	}

	validateETokenTypeIsChanged(t, argsDefault)
}

func validateThatRecievingTokenTypeIsChanged(t *testing.T, args *testutils.ApiGetKeyTestArgs, expectedTokenType string) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	var err error
	_, err = testutils.GetKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while getting key")

	tokenType := args.Apim.GetApplication(args.Apim.GetApplicationByName(utils.DefaultApictlTestAppName).ApplicationID).TokenType
	assert.Equal(t, strings.ToUpper(expectedTokenType), tokenType, "Error getting token type of application.")

	testutils.UnsubscribeAPI(args.Apim, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)
}
