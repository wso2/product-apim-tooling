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
	"strconv"
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

	output, err := base.Execute(t, "init", args.InitFlag, "--definition", args.definitionFlag, "--force", strconv.FormatBool(args.ForceFlag))
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

func ValidateImportInitializedProject(t *testing.T, args *InitTestArgs) {
	t.Helper()
	//Initialize a project with API definition
	ValidateInitializeProjectWithOASFlag(t, args)

	result, error := ImportApiFromProject(t, args.InitFlag, args.SrcAPIM, args.APIName, &args.CtlUser, true)

	base.WaitForIndexing()

	assert.Nil(t, error, "Error while importing Project")
	assert.Contains(t, result, "Successfully imported API", "Error while importing Project")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

func ValidateImportFailedWithInitializedProject(t *testing.T, args *InitTestArgs) {
	t.Helper()

	result, _ := ImportApiFromProject(t, args.InitFlag, args.SrcAPIM, args.APIName, &args.CtlUser, false)

	base.WaitForIndexing()

	assert.Contains(t, result, "Resource Already Exists", "Test failed because API is imported successfully")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}

func ValidateImportUpdatePassedWithInitializedProject(t *testing.T, args *InitTestArgs) {
	t.Helper()

	result, error := ImportApiFromProjectWithUpdate(t, args.InitFlag, args.SrcAPIM.GetEnvName())

	base.WaitForIndexing()

	assert.Nil(t, error, "Error while generating Project")
	assert.Contains(t, result, "Successfully imported API", "Test InitializeProjectWithDefinitionFlag Failed")

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}
