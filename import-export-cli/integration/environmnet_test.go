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
	"path/filepath"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const defaultExportPath = utils.DefaultExportDirName

//Get Environments using apictl
func TestGetEnvironments(t *testing.T) {
	apim := GetDevClient()
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	response, _ := base.Execute(t, "get", "envs")
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, apim.GetEnvName(), "TestGetEnvironments Failed")
}

//Change Export directory using apictl and assert the change
func TestChangeExportDirectory(t *testing.T) {
	dev := GetDevClient()
	changedExportDirectory, _ := filepath.Abs(testutils.CustomTestExportDirectory)

	// Create directory to act as custom export directory
	base.CreateTempDir(t, changedExportDirectory)

	args := &testutils.SetTestArgs{
		SrcAPIM:             dev,
		ExportDirectoryFlag: changedExportDirectory,
	}
	testutils.ValidateExportDirectoryIsChanged(t, args)

	// reset the export directory change after the test is finished
	t.Cleanup(func() {
		argsDefault := &testutils.SetTestArgs{
			SrcAPIM:             args.SrcAPIM,
			ExportDirectoryFlag: utils.DefaultExportDirPath,
		}
		testutils.ValidateExportDirectoryIsChanged(t, argsDefault)
	})

	apimClient := GetDevClient()
	projectName := base.GenerateRandomName(16)
	apiName := testutils.DevFirstDefaultAPIName
	apiVersion := testutils.DevFirstDefaultAPIVersion
	username := superAdminUser
	password := superAdminPassword

	apiArgs := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apimClient,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		ForceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportProject(t, apiArgs)

	//Assert that Export directory change is successful by exporting and asserting that
	testutils.ValidateExportApisPassed(t, apiArgs, changedExportDirectory)

	//Check exporting as a single API
	exportArgs := &testutils.ApiImportExportTestArgs{
		Api: &apim.API{
			Name:    apiName,
			Version: apiVersion,
		},
		SrcAPIM: apimClient,
	}
	testutils.ValidateExportApiPassed(t, exportArgs, changedExportDirectory)
}

//TODO  - Need to come up with  a process to make sure that http timeout is actually changed using another fake server
//Change HTTP request Timeout using apictl and assert the change
//func TestChangeHttpRequestTimout(t *testing.T) {
//	apim := GetDevClient()
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
