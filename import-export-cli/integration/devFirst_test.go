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
	"testing"
	"time"
)

//Initialize a project Initialize an API without any flag
func TestInitializeProject(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	projectName := "SampleTestAPI"

	args := &initTestArgs{
		ctlUser:  credentials{username: username, password: password},
		srcAPIM:  apim,
		initFlag: projectName,
	}

	validateInitializeProject(t, args)

}

func validateInitializeProject(t *testing.T, args *initTestArgs) {
	t.Helper()

	output, err := initProject(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Contains(t, output, "Project initialized", "Project initialization Failed")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.initFlag)
	})
}

//Function to initialize a project using API definition
func validateInitializeProjectWithOASFlag(t *testing.T, args *initTestArgs) {
	t.Helper()

	output, err := initProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test initialization Failed with --oas flag")

	//Remove Created project and logout

	t.Cleanup(func() {
		base.RemoveDir(args.initFlag)
	})

}

//Initialize an API from Swagger 2 Specification
func TestInitializeAPIFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "Swagger2APIProject"
	username := superAdminUser
	password := superAdminPassword

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestSwagger2DefinitionPath,
		forceFlag: false,
	}

	validateInitializeProjectWithOASFlag(t, args)

}

//Initialize an API from OpenAPI 3 Specification
func TestInitializeAPIFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	var projectName = "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestOpenAPI3DefinitionPath,
		forceFlag: false,
	}

	validateInitializeProjectWithOASFlag(t, args)

}

//Initialize an API from API Specification URL
func TestInitializeAPIFromAPIDefinitionURL(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	var projectName = "ProjectInitWithURL"

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestOpenAPISpecificationURL,
		forceFlag: false,
	}

	validateInitializeProjectWithOASFlag(t, args)

}

func validateImportInitializedProject(t *testing.T, args *initTestArgs) {
	t.Helper()
	//Initialize a project with API definition
	validateInitializeProjectWithOASFlag(t, args)

	time.Sleep(1 * time.Second)

	result, error := importApiFromProject(t, args.initFlag, args.srcAPIM.GetEnvName())
	assert.Nil(t, error, "Error while importing Project")
	assert.Contains(t, result, "Successfully imported API", "Error while importing Project")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.initFlag)
	})
}

func validateImportFailedWithInitializedProject(t *testing.T, args *initTestArgs) {
	t.Helper()

	time.Sleep(1 * time.Second)

	result, _ := importApiFromProject(t, args.initFlag, args.srcAPIM.GetEnvName())
	assert.Contains(t, result, "Resource Already Exists", "Test failed because API is imported successfully")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.initFlag)
	})
}

func validateImportUpdatePassedWithInitializedProject(t *testing.T, args *initTestArgs) {
	t.Helper()

	time.Sleep(1 * time.Second)

	result, error := importApiFromProjectWithUpdate(t, args.initFlag, args.srcAPIM.GetEnvName())
	assert.Nil(t, error, "Error while generating Project")
	assert.Contains(t, result, "Successfully imported API", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.initFlag)
	})
}

//Import API from initialized project with swagger 2 definition
func TestImportProjectCreatedFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "Swagger2Project"
	username := superAdminUser
	password := superAdminPassword

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestSwagger2DefinitionPath,
		forceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	validateImportInitializedProject(t, args)

}

//Import API from initialized project with openAPI 3 definition
func TestImportProjectCreatedFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestOpenAPI3DefinitionPath,
		forceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	validateImportInitializedProject(t, args)

}

//Import API from initialized project from API definition which is already in publisher without --update flag
func TestImportProjectCreatedFailWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestOpenAPI3DefinitionPath,
		forceFlag: false,
	}

	//Import API for the First time
	validateImportInitializedProject(t, args)

	//Import API for the second time
	validateImportFailedWithInitializedProject(t, args)

}

//Import API from initialized project from API definition which is already in publisher with --update flag
func TestImportProjectCreatedPassWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   utils.TestOpenAPI3DefinitionPath,
		forceFlag: false,
	}

	//Import API for the First time
	validateImportInitializedProject(t, args)

	//Import API for the second time
	validateImportUpdatePassedWithInitializedProject(t, args)

}
