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
	"os"
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

	//Import API for the second time
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
	//Move doc file to created project
	srcPathForDoc, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseDocPath)
	destPathForDoc := projectPath + testutils.DevFirstUpdatedSampleCaseDestPathSuffix
	base.Copy(srcPathForDoc, destPathForDoc)

	//Move docMetaData file to created project
	srcPathForDocMetadata, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseDocMetaDataPath)
	destPathForDocMetaData := projectPath + testutils.DevFirstUpdatedSampleCaseDestMetaDataPathSuffix
	base.Copy(srcPathForDocMetadata, destPathForDocMetaData)

	//Import the project with Document
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	testutils.ValidateAPIWithDocIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
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
	srcPathForIcon, _ := filepath.Abs(testutils.DevFirstSampleCasePngPath)
	destPathForIcon := projectPath + testutils.DevFirstSampleCasePngPathSuffix
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
	srcPathForImage, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseJpegPath)
	destPathForImage := projectPath + testutils.DevFirstUpdatedSampleCaseDestJpegPathSuffix
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

	//Move doc file to created project
	projectPath, _ := filepath.Abs(projectName)
	srcPathForDoc, _ := filepath.Abs(testutils.DevFirstSampleCaseDocPath)
	destPathForDoc := projectPath + testutils.DevFirstSampleCaseDestPathSuffix
	base.Copy(srcPathForDoc, destPathForDoc)

	//Move Image file to created project
	srcPathForImage, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseJpegPath)
	destPathForImage := projectPath + testutils.DevFirstUpdatedSampleCaseDestJpegPathSuffix
	base.Copy(srcPathForImage, destPathForImage)

	//Import the project with Document and image thumbnail
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	//Update doc file to created project
	srcPathForDocUpdate, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseDocPath)
	destPathForDocUpdate := projectPath + testutils.DevFirstUpdatedSampleCaseDestPathSuffix
	base.Copy(srcPathForDocUpdate, destPathForDocUpdate)

	//Update docMetaData file to created project
	srcPathForDocMetadataUpdate, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseDocMetaDataPath)
	destPathForDocMetaDataUpdate := projectPath + testutils.DevFirstUpdatedSampleCaseDestMetaDataPathSuffix
	base.Copy(srcPathForDocMetadataUpdate, destPathForDocMetaDataUpdate)

	//Update icon file to created project
	err := os.Remove(destPathForImage)
	if err != nil {
		t.Fatal(err)
	}

	srcPathForIcon, _ := filepath.Abs(testutils.DevFirstSampleCasePngPath)
	destPathForIcon := projectPath + testutils.DevFirstSampleCasePngPathSuffix
	base.Copy(srcPathForIcon, destPathForIcon)

	base.WaitForIndexing()
	//Import the project with updated Document and updated image thumbnail
	testutils.ValidateImportUpdateProject(t, args)

	//Validate that image has been updated
	testutils.ValidateAPIWithDocIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)

	//Validate that document has been updated
	testutils.ValidateAPIWithIconIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
}

// Test a verified (syntactically correct) custom sequence update as a super tenant user with Internal/devops role
func TestAPISequenceUpdateWithDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password
	apim := apimClients[1]
	projectName := base.GenerateRandomName(16)

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		SrcAPIM:   apim,
		InitFlag:  projectName,
		OasFlag:   testutils.TestOpenAPI3DefinitionPath,
		APIName:   testutils.DevFirstDefaultAPIName,
		ForceFlag: true,
	}

	// Initialize the project
	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	// Add custom sequence file to created project
	projectPath, _ := filepath.Abs(projectName)
	srcPathForSequence, _ := filepath.Abs(testutils.DevFirstSampleCaseSequencePath)
	destPathForSequence := projectPath + testutils.DevFirstSampleCaseDestSequencePathSuffix
	base.CreateDir(projectPath + testutils.CustomSequenceDirectory)
	base.Copy(srcPathForSequence, destPathForSequence)

	// Update api.yaml file of initialized project with sequence related metadata
	apiMetadataYamlPath := projectPath + testutils.DevFirstSampleCaseApiMetadataPathSuffix
	inSequenceStr := "inSequence: " + testutils.CustomSequenceName
	base.AppendStringToFile(inSequenceStr, apiMetadataYamlPath)

	// Import the project with the verified (syntactically correct) custom sequence
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	// Update custom sequence file of created project
	srcPathForSequenceUpdate, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseSequencePath)
	destPathForSequenceUpdate := projectPath + testutils.DevFirstUpdatedSampleCaseSequencePathSuffix
	err := os.Remove(destPathForSequenceUpdate)
	if err != nil {
		t.Fatal(err)
	}
	base.Copy(srcPathForSequenceUpdate, destPathForSequenceUpdate)

	base.WaitForIndexing()

	// Import the project with updated sequence
	testutils.ValidateImportUpdateProject(t, args)

	// Validate that sequence has been updated
	testutils.ValidateAPIWithUpdatedSequenceIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
}
