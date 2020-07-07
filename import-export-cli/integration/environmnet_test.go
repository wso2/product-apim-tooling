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

func validateExportDirectoryIsChanged(t *testing.T, args *setTestArgs) {
	t.Helper()
	output, _ := environmentSetExportDirectory(t, args)
	base.Log(output)
	assert.Contains(t, output, "Export Directory is set to", "Export Directory change is not successful")
}

func validateExportApisPassed(t *testing.T, args *initTestArgs, directoryName string) {
	t.Helper()
	time.Sleep(5 * time.Second)

	output, error := exportApisWithOneCommand(t, args)
	assert.Nil(t, error, "Error while Exporting APIs")
	assert.Contains(t, output, "export-apis execution completed", "Error while Exporting APIs")

	//Derive exported path from output
	exportedPath := base.GetExportedPathFromOutput(strings.ReplaceAll(output, "Command: export-apis execution completed !", ""))
	count, _ := base.CountFiles(exportedPath)
	assert.Equal(t, 1, count, "Error while Exporting APIs")

	t.Cleanup(func() {
		argsDefault := &setTestArgs{
			srcAPIM:             args.srcAPIM,
			exportDirectoryFlag: utils.DefaultExportDirPath,
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

	args := &setTestArgs{
		srcAPIM:             dev,
		exportDirectoryFlag: changedExportDirectory,
	}
	validateExportDirectoryIsChanged(t, args)

	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	apiArgs := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestOpenAPI3DefinitionPath,
		forceFlag: false,
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

func validateETokenTypeIsChanged(t *testing.T, args *setTestArgs) {
	t.Helper()
	output, _ := environmentSetTokenType(t, args)
	base.Log(output)
	assert.Contains(t, output, "Token type is set to", "1st attempt of Token Type change is not successful")
}

//Change Token type using apictl and assert the change (for both "jwt" and "oauth" token types)
func TestChangeTokenType(t *testing.T) {
	apim := apimClients[0]

	tokenType1 := "oauth"
	args := &setTestArgs{
		srcAPIM:       apim,
		tokenTypeFlag: tokenType1,
	}

	validateETokenTypeIsChanged(t, args)

	//Create API and get keys for that API
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	dev := apimClients[0]

	api := addAPI(t, dev, adminUser, adminPassword)

	publishAPI(dev, adminUser, adminPassword, api.ID)

	apiArgs := &apiGetKeyTestArgs{
		ctlUser: credentials{username: adminUser, password: adminPassword},
		api:     api,
		apim:    dev,
	}

	validateThatRecievingTokenTypeIsChanged(t, apiArgs, tokenType1)

	tokenType2 := "jwt"

	//Change value back to default value
	argsDefault := &setTestArgs{
		srcAPIM:       apim,
		tokenTypeFlag: tokenType2,
	}

	validateETokenTypeIsChanged(t, argsDefault)
}

func validateThatRecievingTokenTypeIsChanged(t *testing.T, args *apiGetKeyTestArgs, expectedTokenType string) {
	t.Helper()

	base.SetupEnv(t, args.apim.GetEnvName(), args.apim.GetApimURL(), args.apim.GetTokenURL())
	base.Login(t, args.apim.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	var err error
	_, err = getKeys(t, args.api.Provider, args.api.Name, args.api.Version, args.apim.GetEnvName())
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while getting key")

	tokenType := args.apim.GetApplication(args.apim.GetApplicationByName(utils.DefaultApictlTestAppName).ApplicationID).TokenType
	assert.Equal(t, strings.ToUpper(expectedTokenType), tokenType, "Error getting token type of application.")

	unsubscribeAPI(args.apim, args.ctlUser.username, args.ctlUser.password, args.api.ID)
}
