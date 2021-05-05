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
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

// Initialize a project and import it with a params file with dynamic data and
// check whether the env variable values have been set correctly in the imported application
func TestImportAPIProjectWithDynamicData(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			preserveProvider := true
			// The default apictl project will contain the provider as super tenant admin
			// Hence, if the CTL user is a tenant user, then this will act like a cross tenant project import.
			// The preserveProvider flag should be false since this is like a cross tenant import.
			if isTenantUser(user.CtlUser.Username, TENANT1) {
				preserveProvider = false
			}

			dev := GetDevClient()

			args := &testutils.InitTestArgs{
				CtlUser:   testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:   dev,
				InitFlag:  base.GenerateRandomName(16),
				OasFlag:   testutils.TestOpenAPI3DefinitionWithoutEndpointsPath,
				APIName:   testutils.OpenAPI3DefinitionWithoutEndpointsAPIName,
				ForceFlag: false,
			}

			// Set environment variables to be used as dynamic data from the params file
			testutils.SetEnvVariablesForAPI(t, dev)

			// Initialize a project with API definition
			testutils.ValidateInitializeProjectWithOASFlag(t, args)

			// Validate whether the project successfully imported
			importedAPI := testutils.ValidateImportProject(t, args, testutils.APIDynamicDataParamsFile, preserveProvider)

			// Validate whether the env variables are correctly injected to the imported API
			testutils.ValidateDynamicData(t, importedAPI)
		})
	}
}

// Initialize a project and import it with a sequence file with dynamic data and export it to
// check whether the env variable values have been set correctly in the exported sequence
func TestImportAPIProjectWithDynamicDataSequence(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			preserveProvider := true
			// The default apictl project will contain the provider as super tenant admin
			// Hence, if the CTL user is a tenant user, then this will act like a cross tenant project import.
			// The preserveProvider flag should be false since this is like a cross tenant import.
			if isTenantUser(user.CtlUser.Username, TENANT1) {
				preserveProvider = false
			}

			dev := GetDevClient()

			args := &testutils.InitTestArgs{
				CtlUser:   testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:   dev,
				InitFlag:  base.GenerateRandomName(16),
				OasFlag:   testutils.TestOpenAPI3DefinitionWithoutEndpointsPath,
				APIName:   testutils.OpenAPI3DefinitionWithoutEndpointsAPIName,
				ForceFlag: false,
			}

			// Set environment variables to be used as dynamic data from the params file
			testutils.SetEnvVariablesForAPI(t, dev)

			// Initialize a project with API definition
			testutils.ValidateInitializeProjectWithOASFlag(t, args)

			// Add the custom IN sequence with dynamic data
			updatedAPIFileContent := testutils.AddSequenceWithDynamicDataToAPIProject(t, args)

			// Validate whether the project successfully imported
			testutils.ValidateImportProject(t, args, "", preserveProvider)

			// Validate whether the exported project contains the correctly env variables substituted content
			testutils.ValidateExportedSequenceWithDynamicData(t, args, updatedAPIFileContent.Data)
		})
	}
}
