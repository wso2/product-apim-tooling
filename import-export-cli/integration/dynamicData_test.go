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

func TestImportAPIProjectWithDynamicDataSuperTenantDevops(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	dev := GetDevClient()

	args := &testutils.InitTestArgs{
		CtlUser:   testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		SrcAPIM:   dev,
		InitFlag:  base.GenerateRandomName(16),
		OasFlag:   testutils.TestOpenAPI3DefinitionWithoutEndpointsPath,
		APIName:   testutils.OpenAPI3DefinitionWithoutEndpointsAPIName,
		ForceFlag: false,
	}

	// Set environment variables to be used as dynamic data from the params file
	testutils.SetDynamicDataForAPI(t, dev)

	// Validate whether the project successfully imported
	importedAPI := testutils.ValidateImportProject(t, args, testutils.APIDynamicDataParamsFile)

	// Validate whether the env variables are correctly injected to the imported API
	testutils.ValidateDynamicData(t, importedAPI)
}
