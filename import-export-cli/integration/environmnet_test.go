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

// Run apictl without any environments or .wso2apictl config folder
func TestRunApictlWithoutAnyEnvironments(t *testing.T) {

	output, err := testutils.InitApictl(t)

	// Validate apictl initialization
	testutils.ValidateApictlInit(t, err, output)
}

// Run apictl by setting a custom wso2apictl directory location by specifying the APICTL_CONFIG_DIR environment variable
func TestRunApictlWithCustomDirectoryLocation(t *testing.T) {

	// Get absolute path of the custom Directory and create temp custom directory
	absolutePathOfCustomDir, _ := filepath.Abs(testutils.CustomDirectoryAtInit)
	_ = utils.CreateDirIfNotExist(absolutePathOfCustomDir)

	testutils.SetApictlWithCustomDirectory(t, absolutePathOfCustomDir)

	// Initializing apictl
	output, err := testutils.InitApictl(t)

	// Validate apictl initialization
	testutils.ValidateApictlInit(t, err, output)

	// Validate custom directory change at init
	testutils.ValidateCustomDirectoryChangeAtInit(t, absolutePathOfCustomDir)
}

// Adding a new Environments with -- token flag and list them and check it
func TestAddEnvironmentWithTokenEndpoint(t *testing.T) {
	apim := GetDevClient()
	output, err := testutils.AddEnvironmentWithTokenFlag(t, apim.GetEnvName(), apim.GetApimURL(), apim.GetTokenURL())

	// Validate added environment
	assert.Nil(t, err)
	testutils.ValidateAddedEnvironments(t, output, apim.GetEnvName(), false)
}

// Adding a new Environments without -- token flag and list them and check it
func TestAddEnvironmentWithoutTokenEndpoint(t *testing.T) {
	apim := GetDevClient()
	output, err := testutils.AddEnvironmentWithOutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())

	// Validate added environment
	assert.Nil(t, err)
	testutils.ValidateAddedEnvironments(t, output, apim.GetEnvName(), false)
}

//Get Environments using apictl
func TestGetEnvironments(t *testing.T) {
	apim := GetDevClient()
	//Adding a new environment
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())

	// Validate added environment list
	testutils.ValidateEnvsList(t, apim.GetEnvName(), true)
}

// Remove an added Environment when an user is logged into the environment
func TestRemoveEnvironmentWithLoggedInUsers(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			devopsUsername := user.CtlUser.Username
			devopsPassword := user.CtlUser.Password

			// Add an environment
			apim := GetDevClient()
			output, err := testutils.AddEnvironmentWithOutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())

			// Validate added environment
			assert.Nil(t, err)
			testutils.ValidateAddedEnvironments(t, output, apim.GetEnvName(), true)

			// Login to the added environment
			base.Execute(t, "login", apim.GetEnvName(), "-u", devopsUsername, "-p", devopsPassword)

			// Remove the added environment
			removeEnvOutput, errr := testutils.RemoveEnvironment(t, apim.GetEnvName())
			assert.Nil(t, errr)

			//Validate removed environment
			testutils.ValidateRemoveEnvironments(t, removeEnvOutput, apim.GetEnvName())
		})
	}
}

// Remove an added Environment when an user is logged into the environment
func TestRemoveEnvironmentWithoutLoggedInUsers(t *testing.T) {

	// Add an environment
	apim := GetDevClient()
	output, err := testutils.AddEnvironmentWithOutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())

	// Validate added environment
	assert.Nil(t, err)
	testutils.ValidateAddedEnvironments(t, output, apim.GetEnvName(), true)

	// Remove the added environment
	removeEnvOutput, errr := testutils.RemoveEnvironment(t, apim.GetEnvName())
	assert.Nil(t, errr)

	//Validate removed environment
	testutils.ValidateRemoveEnvironments(t, removeEnvOutput, apim.GetEnvName())
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

	//Initialize a project with API definition
	testutils.ValidateInitializeProjectWithOASFlag(t, apiArgs)

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportProject(t, apiArgs, "", true)

	//Check exporting as a single API
	exportArgs := &testutils.ApiImportExportTestArgs{
		Api: &apim.API{
			Name:    apiName,
			Version: apiVersion,
		},
		SrcAPIM: apimClient,
	}

	//Assert that Export directory change is successful
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
