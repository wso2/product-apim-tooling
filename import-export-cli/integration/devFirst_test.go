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
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"testing"
)

//Initialize a project Initialize an API without any flag
func TestInitializeProject(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	projectName := "SampleTestAPI"

	args := &testutils.InitTestArgs{
		CtlUser:  testutils.Credentials{Username: username, Password: password},
		SrcAPIM:  apim,
		InitFlag: projectName,
	}

	testutils.ValidateInitializeProject(t, args)

}

//Initialize an API from Swagger 2 Specification
func TestInitializeAPIFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "Swagger2APIProject"
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestSwagger2DefinitionPath,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)

}

//Initialize an API from OpenAPI 3 Specification
func TestInitializeAPIFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	var projectName = "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestOpenAPI3DefinitionPath,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)

}

//Initialize an API from API Specification URL
func TestInitializeAPIFromAPIDefinitionURL(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	var projectName = "ProjectInitWithURL"

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestOpenAPISpecificationURL,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)

}

//Import API from initialized project with swagger 2 definition
func TestImportProjectCreatedFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "Swagger2Project"
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestSwagger2DefinitionPath,
		ForceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportInitializedProject(t, args)

}

//Import API from initialized project with openAPI 3 definition
func TestImportProjectCreatedFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestOpenAPI3DefinitionPath,
		ForceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportInitializedProject(t, args)

}

//Import API from initialized project from API definition which is already in publisher without --update flag
func TestImportProjectCreatedFailWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestOpenAPI3DefinitionPath,
		ForceFlag: false,
	}

	//Import API for the First time
	testutils.ValidateImportInitializedProject(t, args)

	//Import API for the second time
	testutils.ValidateImportFailedWithInitializedProject(t, args)

}

//Import API from initialized project from API definition which is already in publisher with --update flag
func TestImportProjectCreatedPassWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   utils.TestOpenAPI3DefinitionPath,
		ForceFlag: false,
	}

	//Import API for the First time
	testutils.ValidateImportInitializedProject(t, args)

	//Import API for the second time
	testutils.ValidateImportUpdatePassedWithInitializedProject(t, args)

}
