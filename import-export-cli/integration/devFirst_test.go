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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
)

// Initialize a API project by getting the OAS of a AWS API and import it as a super tenant user with
// the Internal/devops role
// WARN: To test this you need to have AWS CLI installed and configured
// WARN: Before running this test create an API on AWS with the name "Shopping" and stage name "Live"
// NOTE: If the above prerequisites are met, uncomment the two aws init test functions and run the aws init tests.

//func TestAWSInitImportSuperTenant(t *testing.T) {
//	username := devops.UserName
//	password := devops.Password
//	apim := GetDevClient()
//	apiName := "Shopping"
//	apiStageName := "Live"
//
//	args := &testutils.AWSInitTestArgs{
//		CtlUser:  testutils.Credentials{Username: username, Password: password},
//		SrcAPIM:  apim,
//		ApiNameFlag: apiName,
//		ApiStageNameFlag : apiStageName,
//		InitFlag: apiName,
//	}
//
//	testutils.ValidateAWSInitProject(t, args)
//	testutils.ValidateAWSProjectImport(t, args, true)
//}

// Initialize a API project by getting the OAS of a AWS API and import it as a tenant user with
// the Internal/devops role
// WARN: To test this you need to have AWS CLI installed and configured
// WARN: Before running this test create an API on AWS with the name "PetStore" and stage name "beta"

//func TestAWSInitImportTenant(t *testing.T) {
//	username := devops.UserName + "@" + TENANT1
//	password := devops.Password
//	apim := GetDevClient()
//	apiName := "PetStore"
//	apiStageName := "beta"
//
//	args := &testutils.AWSInitTestArgs{
//		CtlUser:  testutils.Credentials{Username: username, Password: password},
//		SrcAPIM:  apim,
//		ApiNameFlag: apiName,
//		ApiStageNameFlag : apiStageName,
//		InitFlag: apiName,
//	}
//
//	testutils.ValidateAWSInitProject(t, args)
//	// making preserveprovider false since this is a cross tenant import
//	testutils.ValidateAWSProjectImport(t, args, false)
//}

//Initialize a project Initialize an API without any flag
func TestInitializeProject(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			apim := GetDevClient()
			projectName := base.GenerateRandomString()

			args := &testutils.InitTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:  apim,
				InitFlag: projectName,
				APIName:  projectName, // The logic of apictl init has been written to consider the projectName as
				// the API name, if an OAS or a definition is not provided
			}

			testutils.ValidateInitializeProject(t, args)

			projectPath, _ := filepath.Abs(projectName)
			apiYamlPath := projectPath + string(os.PathSeparator) + testutils.APIYamlFilePath

			// Read the api.yaml file in the exported directory
			fileData, _ := ioutil.ReadFile(apiYamlPath)

			fileContent := make(map[string]interface{})
			err := yaml.Unmarshal(fileData, &fileContent)
			if err != nil {
				t.Error(err)
			}
			apiArtifactVersion := fileContent["version"].(string)

			assert.Equal(t, apiArtifactVersion, "v"+yamlConfig.APICTLVersion,
				"Artifact version: "+apiArtifactVersion+
					" does not matches with the APICTL version: v"+yamlConfig.APICTLVersion)

			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))
		})
	}
}

// Initialize an API with --definition flag and import it
func TestInitializeAPIWithDefinitionFlag(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			apim := GetDevClient()
			projectName := base.GenerateRandomString()

			args := &testutils.InitTestArgs{
				CtlUser:        testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:        apim,
				InitFlag:       projectName,
				DefinitionFlag: testutils.SampleAPIYamlFilePath,
				APIName:        testutils.DevFirstDefinitionFlagSampleAPIName,
				ForceFlag:      false,
			}

			testutils.ValidateInitializeProjectWithDefinitionFlag(t, args)

			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))
		})
	}
}

//Initialize an API from Swagger 2 Specification
func TestInitializeAPIFromSwagger2Definition(t *testing.T) {
	apim := GetDevClient()
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
	apim := GetDevClient()
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
	apim := GetDevClient()
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
	apim := GetDevClient()
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

	//Initialize a project with API definition
	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportProject(t, args, "", true)
}

//Import API from initialized project with openAPI 3 definition
func TestImportProjectCreatedFromOpenAPI3Definition(t *testing.T) {
	apim := GetDevClient()
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

	//Initialize a project with API definition
	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	//Assert that project import to publisher portal is successful
	testutils.ValidateImportProject(t, args, "", true)
}

// Import API from initialized project from API definition which is already in publisher with --update flag
func TestImportProjectCreatedPassWhenAPIIsExisted(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			apim := GetDevClient()
			projectName := base.GenerateRandomString()

			args := &testutils.InitTestArgs{
				CtlUser:   testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:   apim,
				InitFlag:  projectName,
				OasFlag:   testutils.TestOpenAPI3DefinitionPath,
				APIName:   testutils.DevFirstDefaultAPIName,
				ForceFlag: false,
			}

			//Initialize a project with API definition
			testutils.ValidateInitializeProjectWithOASFlag(t, args)

			//Assert that project import to publisher portal is successful
			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))

			// Read the API definition file in the project
			apiDefinitionFilePath := args.InitFlag + string(os.PathSeparator) + utils.APIDefinitionFileYaml
			apiDefinitionFileContent := testutils.ReadAPIDefinition(t, apiDefinitionFilePath)

			// Change the description
			apiDefinitionFileContent.Data.Description = "Updated description"

			// Write the modified API definition to the directory
			testutils.WriteToAPIDefinition(t, apiDefinitionFileContent, apiDefinitionFilePath)

			// Import and validate new API with the description change
			importedApi := testutils.ValidateImportUpdateProject(t, args, !isTenantUser(user.CtlUser.Username, TENANT1))

			assert.Equal(t, importedApi.Description, apiDefinitionFileContent.Data.Description, "Description is not updated")
		})
	}
}

//Import API from initialized project from API definition which is already in publisher without --update flag
func TestImportProjectCreatedFailWhenAPIIsExisted(t *testing.T) {
	apim := GetDevClient()
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

	//Initialize a project with API definition
	testutils.ValidateInitializeProjectWithOASFlag(t, args)

	//Import API for the First time
	testutils.ValidateImportProject(t, args, "", true)

	//Import API for the second time
	testutils.ValidateImportProjectFailed(t, args, "")
}

//Import Api with a Document and Export that Api with a Document
func TestImportAndExportAPIWithDocument(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := GetDevClient()
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
	base.CreateTempDir(t, projectPath+testutils.DevFirstUpdatedSampleCaseDocName)

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

	testutils.ValidateAPIWithDocIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion,
		testutils.DevFirstUpdatedSampleCaseDestPathSuffix)
}

//Import Api with an Image and Export that Api with an image (.png Type)
func TestImportAndExportAPIWithPngIcon(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	apim := GetDevClient()
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
	destPathForIcon := projectPath + testutils.DevFirstSampleCaseDestPngPathSuffix
	base.Copy(srcPathForIcon, destPathForIcon)

	//Import the project with icon image(.png)
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	testutils.ValidateAPIWithIconIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
}

//Import Api with an Image and Export that Api with an image (.jpeg Type)
func TestImportAndExportAPIWithJpegImage(t *testing.T) {
	apim := GetDevClient()
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
	apim := GetProdClient()
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
	base.CreateTempDir(t, projectPath+testutils.DevFirstSampleCaseDocName)

	//Move doc file to created project
	srcPathForDoc, _ := filepath.Abs(testutils.DevFirstSampleCaseDocPath)
	destPathForDoc := projectPath + testutils.DevFirstSampleCaseDestPathSuffix
	base.Copy(srcPathForDoc, destPathForDoc)

	//Move docMetaData file to created project
	srcPathForDocMetadata, _ := filepath.Abs(testutils.DevFirstSampleCaseDocMetaDataPath)
	destPathForDocMetaData := projectPath + testutils.DevFirstSampleCaseDestMetaDataPathSuffix
	base.Copy(srcPathForDocMetadata, destPathForDocMetaData)

	//Move icon file to created project
	srcPathForImage, _ := filepath.Abs(testutils.DevFirstSampleCasePngPath)
	destPathForImage := projectPath + testutils.DevFirstSampleCaseDestPngPathSuffix
	base.Copy(srcPathForImage, destPathForImage)

	//Import the project with Document and image thumbnail
	testutils.ValidateImportUpdateProjectNotAlreadyImported(t, args)

	//Update doc file to created project
	srcPathForDocUpdate, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseDocPath)
	destPathForDocUpdate := projectPath + testutils.DevFirstSampleCaseDestPathSuffix
	base.Copy(srcPathForDocUpdate, destPathForDocUpdate)

	//Update image file to created project
	err := os.Remove(destPathForImage)
	if err != nil {
		t.Fatal(err)
	}

	srcPathForIcon, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCaseJpegPath)
	destPathForIcon := projectPath + testutils.DevFirstUpdatedSampleCaseDestJpegPathSuffix
	base.Copy(srcPathForIcon, destPathForIcon)

	base.WaitForIndexing()
	//Import the project with updated Document and updated image thumbnail
	testutils.ValidateImportUpdateProject(t, args, !isTenantUser(username, TENANT1))

	//Validate that image has been updated
	testutils.ValidateAPIWithDocIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion,
		testutils.DevFirstSampleCaseDestPathSuffix)

	//Validate that document has been updated
	testutils.ValidateAPIWithImageIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
}

// Test a verified (syntactically correct) custom operation policy (sequence) update
func TestAPISequenceUpdate(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			projectName := base.GenerateRandomString()

			args := &testutils.InitTestArgs{
				CtlUser:   testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:   dev,
				InitFlag:  projectName,
				OasFlag:   testutils.TestOpenAPI3DefinitionPath,
				APIName:   testutils.DevFirstDefaultAPIName,
				ForceFlag: true,
			}

			// Initialize the project
			testutils.ValidateInitializeProjectWithOASFlag(t, args)

			// Add custom operation policy (sequence) file to created project
			projectPath, _ := filepath.Abs(projectName)
			base.CreateTempDir(t, projectPath+testutils.PoliciesDirectory)

			srcPathForSequence, _ := filepath.Abs(testutils.DevFirstSampleCasePolicyPath)
			destPathForSequence := projectPath + testutils.DevFirstSampleCaseDestPolicyPathSuffix
			base.Copy(srcPathForSequence, destPathForSequence)

			srcPathForSequenceDefinition, _ := filepath.Abs(testutils.DevFirstSampleCasePolicyDefinitionPath)
			destPathForSequenceDefinition := projectPath + testutils.DevFirstSampleCaseDestPolicyDefinitionPathSuffix
			base.Copy(srcPathForSequenceDefinition, destPathForSequenceDefinition)

			// Update api.yaml file of initialized project with sequence related metadata
			apiYamlFilePath := filepath.Join(projectPath, testutils.DevFirstSampleCaseApiYamlFilePathSuffix)
			apiYaml, err := ioutil.ReadFile(apiYamlFilePath)
			if err != nil {
				t.Error(err)
			}

			var api *apim.APIFile
			err = yaml.Unmarshal(apiYaml, &api)
			if err != nil {
				t.Error(err)
			}

			// Operation policy that will be added
			var requestPolicies []interface{}
			operationPolicies := apim.OperationPolicies{
				Request: append(requestPolicies, map[string]interface{}{
					"policyName": testutils.TestSamplePolicyName,
					"parameters": map[string]string{
						testutils.TestSampleOperationPolicyPropertyNameField:  testutils.TestSampleOperationPolicyPropertyName,
						testutils.TestSampleOperationPolicyPropertyValueField: testutils.TestSampleOperationPolicyPropertyValue,
					}}),
				Response: []string{},
				Fault:    []string{},
			}

			// Assign the above operation policy to a resource
			apiOperationWithPolicy := apim.APIOperations{
				Target:            testutils.TestSampleOperationTarget,
				Verb:              testutils.TestSampleOperationVerb,
				AuthType:          testutils.TestSampleOperationAuthType,
				ThrottlingPolicy:  testutils.TestSampleOperationThrottlingPolicy,
				OperationPolicies: operationPolicies,
			}

			// Add the operation policy added resource to the API
			apiOperations := []apim.APIOperations{apiOperationWithPolicy}
			api.Data.Operations = apiOperations

			// Write the modified api.yaml to initialized project
			apiYamlContent, err := yaml.Marshal(api)
			if err != nil {
				t.Error(err)
			}

			err = ioutil.WriteFile(apiYamlFilePath, apiYamlContent, os.ModePerm)
			if err != nil {
				t.Error(err)
			}

			// Import the project with the verified (syntactically correct) operation policy
			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))

			// Update operation policy file of created project
			srcPathForSequenceUpdate, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCasePolicyPath)
			destPathForSequenceUpdate := projectPath + testutils.DevFirstSampleCaseDestPolicyPathSuffix
			err = os.Remove(destPathForSequenceUpdate)
			if err != nil {
				t.Fatal(err)
			}
			base.Copy(srcPathForSequenceUpdate, destPathForSequenceUpdate)
			base.WaitForIndexing()

			// Update operation policy file of created project
			srcPathForSequenceDefinitionUpdate, _ := filepath.Abs(testutils.DevFirstUpdatedSampleCasePolicyDefinitionPath)
			destPathForSequenceDefinitionUpdate := projectPath + testutils.DevFirstSampleCaseDestPolicyDefinitionPathSuffix
			err = os.Remove(destPathForSequenceDefinitionUpdate)
			if err != nil {
				t.Fatal(err)
			}
			base.Copy(srcPathForSequenceDefinitionUpdate, destPathForSequenceDefinitionUpdate)
			base.WaitForIndexing()

			// Import the project with the updated sequence
			testutils.ValidateImportUpdateProject(t, args, !isTenantUser(user.CtlUser.Username, TENANT1))

			// Validate that sequence has been updated
			testutils.ValidateAPIWithUpdatedSequenceIsExported(t, args, testutils.DevFirstDefaultAPIName, testutils.DevFirstDefaultAPIVersion)
		})
	}
}
