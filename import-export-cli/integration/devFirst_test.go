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

	output, err := initProject(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Contains(t, output, "Project initialized", "TestInitializeProject Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Function to initialize a project using API definition
func InitializeAPIFromDefinition(t *testing.T, projectName string, definitionFileName string) (string, error) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   definitionFileName,
		forceFlag: false,
	}

	output, err := initProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	return output, err
}

//Initialize an API from Swagger 2 Specification
func TestInitializeAPIFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "Swagger2APIProject"

	output, err := InitializeAPIFromDefinition(t, projectName, utils.TestSwagger2DefinitionPath)
	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Initialize an API from OpenAPI 3 Specification
func TestInitializeAPIFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	var projectName = "OpenAPI3Project"
	output, err := InitializeAPIFromDefinition(t, projectName, utils.TestOpenAPI3DefinitionPath)

	assert.Nil(t, err, "Error while generating Project")
	assert.Contains(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
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

	output, err := initProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Contains(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Import API from initialized project with swagger 2 definition
func TestImportProjectCreatedFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"

	output, err := InitializeAPIFromDefinition(t, projectName, utils.TestSwagger2DefinitionPath)
	assert.Contains(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")
	assert.Nil(t, err, "Error while generating Project")

	result, error := importApiFromProject(t, projectName, apim.GetEnvName())
	assert.Nil(t, error, "Error while generating Project")
	assert.Contains(t, result, "Successfully imported API", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Import API from initialized project with openAPI 3 definition
func TestImportProjectCreatedFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"

	output, err := InitializeAPIFromDefinition(t, projectName, utils.TestOpenAPI3DefinitionPath)
	assert.Contains(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")
	assert.Nil(t, err, "Error while generating Project")

	result, error := importApiFromProject(t, projectName, apim.GetEnvName())
	assert.Nil(t, error, "Error while generating Project")
	assert.Contains(t, result, "Successfully imported API", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Import API from initialized project from API definition which is already in publisher without --update flag
func TestImportProjectCreatedFailWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"

	InitializeAPIFromDefinition(t, projectName, utils.TestOpenAPI3DefinitionPath)

	//Import API for the First time
	importApiFromProject(t, projectName, apim.GetEnvName())

	//Import API for the second time
	result, _ := importApiFromProject(t, projectName, apim.GetEnvName())
	assert.Contains(t, result, "Resource Already Exists", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Import API from initialized project from API definition which is already in publisher with --update flag
func TestImportProjectCreatedPassWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := "OpenAPI3Project"

	InitializeAPIFromDefinition(t, projectName, utils.TestOpenAPI3DefinitionPath)

	//Import API for the First time
	importApiFromProject(t, projectName, apim.GetEnvName())

	//Import API for the second time
	result, error := importApiFromProjectWithUpdate(t, projectName, apim.GetEnvName())
	assert.Nil(t, error, "Error while generating Project")
	assert.Contains(t, result, "Successfully imported API", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	base.RemoveDir(projectName)
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}
