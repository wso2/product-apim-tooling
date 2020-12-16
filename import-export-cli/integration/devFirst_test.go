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

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

//Initialize a project Initialize an API without any flag
func TestInitializeProject(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)

	args := &testutils.InitTestArgs{
		CtlUser:  testutils.Credentials{Username: username, Password: password},
		SrcAPIM:  apim,
		InitFlag: projectName,
	}

	testutils.ValidateInitializeProject(t, args)
}

//Initialize an API with --definition flag
func TestInitializeAPIWithDefinitionFlag(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestApiDefinitionPath,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithDefinitionFlag(t, args)
}

//Initialize an API from Swagger 2 Specification
func TestInitializeAPIFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestSwagger2DefinitionPath,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)
}

//Initialize an API from OpenAPI 3 Specification
func TestInitializeAPIFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)
}

//Initialize an API from API Specification URL
func TestInitializeAPIFromAPIDefinitionURL(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPISpecificationURL,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)
}

//Import API from initialized project with swagger 2 definition
func TestImportProjectCreatedFromSwagger2Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestSwagger2DefinitionPath,
		APIName:   testutils.DevFirstSwagger2APIName,
		ForceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportProject(t, args)
}

//Import API from initialized project with openAPI 3 definition
func TestImportProjectCreatedFromOpenAPI3Definition(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: false,
	}

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportProject(t, args)
}

//Import API from initialized project from API definition which is already in publisher without --update flag
func TestImportProjectCreatedFailWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: false,
	}

	//Import API for the First time
	testutils.ValidateImportProject(t, args)

	//Import API for the second time
	testutils.ValidateImportProjectFailed(t, args)
}

//Import API from initialized project from API definition which is already in publisher with --update flag
func TestImportProjectCreatedPassWhenAPIIsExisted(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: false,
	}

	//Import API for the First time
	testutils.ValidateImportProject(t, args)

	//Import API for the second time with update flag
	testutils.ValidateImportUpdateProject(t, args)
}

//Import Api with a Document and Export that Api with a Document
func TestImportAndExportAPIWithDocument(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	projectPath, _ := filepath.Abs(projectName)
	base.CreateTempDir(t, projectPath+testutils.TestCase1DocName)

	//Move doc file to created project
	srcPathForDoc, _ := filepath.Abs(testutils.TestCase1DocPath)
	destPathForDoc := projectPath + testutils.TestCase1DestPathSuffix
	base.Copy(srcPathForDoc, destPathForDoc)

	//Move docMetaData file to created project
	srcPathForDocMetadata, _ := filepath.Abs(testutils.TestCase1DocMetaDataPath)
	destPathForDocMetaData := projectPath + testutils.TestCase1DestMetaDataPathSuffix
	base.Copy(srcPathForDocMetadata, destPathForDocMetaData)

	//Import the project with Document
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	testutils.ValidateAPIWithDocIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion,
		testutils.TestCase1DestPathSuffix)
}

//Import Api with an Image and Export that Api with an image (.png Type)
func TestImportAndExportAPIWithPngIcon(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	//Move icon file to created project
	projectPath, _ := filepath.Abs(projectName)
	srcPathForIcon, _ := filepath.Abs(testutils.TestCase2PngPath)
	destPathForIcon := projectPath + testutils.TestCase2DestPngPathSuffix
	base.Copy(srcPathForIcon, destPathForIcon)

	//Import the project with icon image(.png)
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	testutils.ValidateAPIWithIconIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
}

//Import Api with an Image and Export that Api with an image (.jpeg Type)
func TestImportAndExportAPIWithJpegImage(t *testing.T) {
	apim := apimClients[0]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: false,
	}

	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	//Move Image file to created project
	projectPath, _ := filepath.Abs(projectName)
	srcPathForImage, _ := filepath.Abs(testutils.TestCase2JpegPath)
	destPathForImage := projectPath + testutils.TestCase2DestJpegPathSuffix
	base.Copy(srcPathForImage, destPathForImage)

	//Import the project with icon image(.jpeg) provided
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	testutils.ValidateAPIWithImageIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
}

//Import and export API with updated thumbnail and document and assert that
func TestUpdateDocAndImageOfAPIOfExistingAPI(t *testing.T) {
	apim := apimClients[1]
	projectName := base.GenerateRandomName(16)
	username := superAdminUser
	password := superAdminPassword

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: username, Password: password},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: true,
	}
	//Initialize the project
	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	projectPath, _ := filepath.Abs(projectName)
	base.CreateTempDir(t, projectPath+testutils.TestCase2DocName)

	//Move doc file to created project
	srcPathForDoc, _ := filepath.Abs(testutils.TestCase2DocPath)
	destPathForDoc := projectPath + testutils.TestCase2DestPathSuffix
	base.Copy(srcPathForDoc, destPathForDoc)

	//Move docMetaData file to created project
	srcPathForDocMetadata, _ := filepath.Abs(testutils.TestCase2DocMetaDataPath)
	destPathForDocMetaData := projectPath + testutils.TestCase2DestMetaDataPathSuffix
	base.Copy(srcPathForDocMetadata, destPathForDocMetaData)

	//Move icon file to created project
	srcPathForImage, _ := filepath.Abs(testutils.TestCase2PngPath)
	destPathForImage := projectPath + testutils.TestCase2DestPngPathSuffix
	base.Copy(srcPathForImage, destPathForImage)

	//Import the project with Document and image thumbnail
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	//Update doc file to created project
	srcPathForDocUpdate, _ := filepath.Abs(testutils.TestCase1DocPath)
	destPathForDocUpdate := projectPath + testutils.TestCase2DestPathSuffix
	base.Copy(srcPathForDocUpdate, destPathForDocUpdate)

	//Update image file to created project
	srcPathForIcon, _ := filepath.Abs(testutils.TestCase2JpegPath)
	destPathForIcon := projectPath + testutils.TestCase2DestJpegPathSuffix
	base.Copy(srcPathForIcon, destPathForIcon)

	base.WaitForIndexing()
	//Import the project with updated Document and updated image thumbnail
	testutils.ValidateImportUpdateProject(t, args)

	//Validate that image has been updated
	testutils.ValidateAPIWithDocIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion,
		testutils.TestCase2DestPathSuffix)

	//Validate that document has been updated
	testutils.ValidateAPIWithImageIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
}
