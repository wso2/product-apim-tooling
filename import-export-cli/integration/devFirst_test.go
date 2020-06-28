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
	"log"
	"testing"
)

//Initialize a project Initialize an API without any flag
func TestInitializeProject(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	var projectName = "SampleTestAPI"

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
	assert.Containsf(t, output, "Project initialized", "TestInitializeProject Failed")

	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Initialize an API the the definition.yaml (apictl init SampleAPI --definition definition.yaml)
func TestInitializeProjectWithDefinitionFlag(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	var projectName = "SampleTestAPI2"
	var definitionFilePath = ""

	args := &initTestArgs{
		ctlUser:        credentials{username: username, password: password},
		srcAPIM:        apim,
		initFlag:       projectName,
		definitionFlag: definitionFilePath,
		forceFlag:      true,
	}

	output, err := initProjectWithDefinitionFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Initialize an API from Swagger 2 Specification
func TestInitializeAPIFromSwagger2Definition(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	var projectName = "SampleTestAPI2"
	var swaggerDefinitionFilePath = "testdata/swagger2Definition.yaml"

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   swaggerDefinitionFilePath,
		forceFlag: false,
	}

	output, err := initProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Initialize an API from OpenAPI 3 Specification
func TestInitializeAPIFromOpenAPI3Definition(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	var projectName = "OpenAPI3"
	var openAPISpecificationFilePath = "testdata/openAPI3Definition.yaml"

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   openAPISpecificationFilePath,
		forceFlag: false,
	}

	output, err := initProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}

//Initialize an API from API Specification URL
func TestInitializeAPIFromAPIDefinitionURL(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	var projectName = "PetStore"
	var openAPISpecificationURL = "https://petstore.swagger.io/v2/swagger.json"

	args := &initTestArgs{
		ctlUser:   credentials{username: username, password: password},
		srcAPIM:   apim,
		initFlag:  projectName,
		oasFlag:   openAPISpecificationURL,
		forceFlag: false,
	}

	output, err := initProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test InitializeProjectWithDefinitionFlag Failed")
	t.Cleanup(func() {
		base.Execute(t, "logout", apim.GetEnvName())
	})
}
