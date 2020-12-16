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

package testutils

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

func InitProject(t *testing.T, args *InitTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, err := base.Execute(t, "init", args.InitFlag)
	return output, err
}

func InitProjectWithDefinitionFlag(t *testing.T, args *InitTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, err := base.Execute(t, "init", args.InitFlag, "--definition", args.definitionFlag)
	return output, err
}

func ValidateInitializeProject(t *testing.T, args *InitTestArgs) {
	t.Helper()

	output, err := InitProject(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Contains(t, output, "Project initialized", "Project initialization Failed")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

//Function to initialize a project using API definition
func ValidateInitializeProjectWithOASFlag(t *testing.T, args *InitTestArgs) {
	t.Helper()

	output, err := InitProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test initialization Failed with --oas flag")

	//Remove Created project and logout

	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

//Function to initialize a project using API definition
func ValidateInitializeProjectWithOASFlagWithoutCleaning(t *testing.T, args *InitTestArgs) {
	t.Helper()

	output, err := InitProjectWithOasFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test initialization Failed with --oas flag")

}

//Function to initialize a project using API definition using --definition flag
func ValidateInitializeProjectWithDefinitionFlag(t *testing.T, args *InitTestArgs) {
	t.Helper()

	output, err := InitProjectWithDefinitionFlag(t, args)
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while generating Project")
	assert.Containsf(t, output, "Project initialized", "Test initialization Failed with --oas flag")

	//Remove created project
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

func ValidateImportProject(t *testing.T, args *InitTestArgs) {
	t.Helper()
	//Initialize a project with API definition
	ValidateInitializeProjectWithOASFlag(t, args)

	result, error := ImportApiFromProject(t, args.InitFlag, args.SrcAPIM, args.APIName, &args.CtlUser, true)

	assert.Nil(t, error, "Error while importing Project")
	assert.Contains(t, result, "Successfully imported API", "Error while importing Project")

	base.WaitForIndexing()

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

func ValidateImportProjectFailed(t *testing.T, args *InitTestArgs) {
	t.Helper()

	result, _ := ImportApiFromProject(t, args.InitFlag, args.SrcAPIM, args.APIName, &args.CtlUser, false)

	assert.Contains(t, result, "409", "Test failed because API is imported successfully")

	base.WaitForIndexing()

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

func ValidateImportUpdateProject(t *testing.T, args *InitTestArgs) {
	t.Helper()

	result, error := ImportApiFromProjectWithUpdate(t, args.InitFlag, args.SrcAPIM, args.APIName, &args.CtlUser, false)

	assert.Nil(t, error, "Error while generating Project")
	assert.Contains(t, result, "Successfully imported API", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

func ValidateImportUpdateProjectNotAlreadyImported(t *testing.T, args *InitTestArgs) {
	t.Helper()

	result, error := ImportApiFromProjectWithUpdate(t, args.InitFlag, args.SrcAPIM, args.APIName, &args.CtlUser, true)

	assert.Nil(t, error, "Error while generating Project")
	assert.Contains(t, result, "Successfully imported API", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

func ValidateExportImportedAPI(t *testing.T, args *InitTestArgs, DevFirstDefaultAPIName string, DevFirstDefaultAPIVersion string) string {
	expOutput, expError := exportApiImportedFromProject(t, DevFirstDefaultAPIName, DevFirstDefaultAPIVersion, args.SrcAPIM.GetEnvName())
	//Check whether api is exported or not
	assert.Nil(t, expError, "Error while Exporting API")
	assert.Contains(t, expOutput, "Successfully exported API!", "Error while exporting API")
	return expOutput
}

func ValidateAPIWithDocIsExported(t *testing.T, args *InitTestArgs, DevFirstDefaultAPIName, DevFirstDefaultAPIVersion, TestCaseDestPathSuffix string) {
	expOutput := ValidateExportImportedAPI(t, args, DevFirstDefaultAPIName, DevFirstDefaultAPIVersion)

	//Unzip exported API and check whether the imported doc is in there
	exportedPath := base.GetExportedPathFromOutput(expOutput)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	docPathOfExportedApi := relativePath + TestDefaultExtractedFileName + TestCaseDestPathSuffix

	//Check whether the file is available
	isDocExported := base.IsFileAvailable(docPathOfExportedApi)
	base.Log("Doc is Exported", isDocExported)
	assert.Equal(t, true, isDocExported, "Error while exporting API with document")

	t.Cleanup(func() {
		//Remove Created project and logout
		base.RemoveDir(args.InitFlag)
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}

func ValidateAPIWithIconIsExported(t *testing.T, args *InitTestArgs, DevFirstDefaultAPIName string, DevFirstDefaultAPIVersion string) {
	expOutput := ValidateExportImportedAPI(t, args, DevFirstDefaultAPIName, DevFirstDefaultAPIVersion)

	//Unzip exported API and check whether the imported image(.png) is in there
	exportedPath := base.GetExportedPathFromOutput(expOutput)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	iconPathOfExportedApi := relativePath + TestDefaultExtractedFileName + TestCase2DestPngPathSuffix

	isIconExported := base.IsFileAvailable(iconPathOfExportedApi)
	base.Log("Icon is Exported", isIconExported)
	assert.Equal(t, true, isIconExported, "Error while exporting API with icon")

	t.Cleanup(func() {
		//Remove Created project and logout
		base.RemoveDir(args.InitFlag)
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}

func ValidateAPIWithImageIsExported(t *testing.T, args *InitTestArgs, DevFirstDefaultAPIName string, DevFirstDefaultAPIVersion string) {
	expOutput := ValidateExportImportedAPI(t, args, DevFirstDefaultAPIName, DevFirstDefaultAPIVersion)

	//Unzip exported API and check whethers the imported image(.png) is in there
	exportedPath := base.GetExportedPathFromOutput(expOutput)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	imagePathOfExportedApi := relativePath + TestDefaultExtractedFileName + TestCase2DestJpegPathSuffix
	isIconExported := base.IsFileAvailable(imagePathOfExportedApi)
	base.Log("Image is Exported", isIconExported)
	assert.Equal(t, true, isIconExported, "Error while exporting API with icon")

	t.Cleanup(func() {
		//Remove Created project and logout
		base.RemoveDir(args.InitFlag)
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}
