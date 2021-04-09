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

// As a super tenant Admin user, initialize an API project and import it to an environment with certificates and
// check whether the advertise only properties (advertised, original devportal URL and API owner) have been set correctly.
// Export the API and check whether certificates have not been exported.
func TestInitDeploymentDirImportExportAdvertiseOnlyAPIAdminSuperTenant(t *testing.T) {
	apim := GetDevClient()
	projectName := base.GenerateRandomName(16)

	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	prod := GetProdClient()

	advertiseOnlyAPIDef, api := testutils.GenerateAdvertiseOnlyAPIDefinition(t)

	args := &testutils.InitTestArgs{
		CtlUser:        testutils.Credentials{Username: adminUsername, Password: adminPassword},
		SrcAPIM:        apim,
		InitFlag:       projectName,
		DefinitionFlag: advertiseOnlyAPIDef,
	}

	testutils.ValidateInitializeProjectWithDefinitionFlag(t, args)

	filePath, _ := filepath.Abs(projectName)
	genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
		Source:      filePath,
		Destination: testutils.GetEnvAPIExportPath(prod.GetEnvName()),
	}

	testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

	importExportArgs := &testutils.ApiImportExportTestArgs{
		ApiProvider:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:        testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Api:            &api,
		SrcAPIM:        prod,
		DestAPIM:       prod,
		ImportFilePath: filePath,
	}

	// Store the deployment directory path to be provided as the params during import
	importExportArgs.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination,
		api.Name, api.Version)
	testutils.ValidateAPIImportExportWithDeploymentDirForAdvertiseOnlyAPI(t, importExportArgs)
}

// As a tenant Admin user, initialize an API project and import it to an environment (cross tenant import since the API definition
// has the provider as the super tenant admin) with certificates and check whether the advertise only properties (advertised,
// original devportal URL and API owner) have been set correctly. Export the API and check whether certificates have not been exported.
func TestInitDeploymentDirImportExportAdvertiseOnlyAPIAdminTenant(t *testing.T) {
	apim := GetDevClient()
	projectName := base.GenerateRandomName(16)

	adminUsername := superAdminUser + "@" + TENANT1
	adminPassword := superAdminPassword

	prod := GetProdClient()

	advertiseOnlyAPIDef, api := testutils.GenerateAdvertiseOnlyAPIDefinition(t)

	args := &testutils.InitTestArgs{
		CtlUser:        testutils.Credentials{Username: adminUsername, Password: adminPassword},
		SrcAPIM:        apim,
		InitFlag:       projectName,
		DefinitionFlag: advertiseOnlyAPIDef,
	}

	testutils.ValidateInitializeProjectWithDefinitionFlag(t, args)

	filePath, _ := filepath.Abs(projectName)
	genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
		Source:      filePath,
		Destination: testutils.GetEnvAPIExportPath(prod.GetEnvName()),
	}

	testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

	importExportArgs := &testutils.ApiImportExportTestArgs{
		// Since this a cross tenant import (The initialized API was super tenant admin's, but now importing to a tenant)
		// the preserve provider will be false (OverrideProvider: true). Hence, the provider will be overridden to the CTL user
		ApiProvider:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:          testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Api:              &api,
		SrcAPIM:          prod,
		DestAPIM:         prod,
		ImportFilePath:   filePath,
		OverrideProvider: true,
	}

	// Store the deployment directory path to be provided as the params during import
	importExportArgs.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination,
		api.Name, api.Version)
	testutils.ValidateAPIImportExportWithDeploymentDirForAdvertiseOnlyAPI(t, importExportArgs)
}

// As a super tenant user with the Internal/devops role, initialize an API project and import it to an environment
// with certificates and check whether the advertise only properties (advertised, original devportal URL and API owner)
// have been set correctly. Export the API and check whether certificates have not been exported.
func TestInitDeploymentDirImportExportAdvertiseOnlyAPIDevopsSuperTenant(t *testing.T) {
	apim := GetDevClient()
	projectName := base.GenerateRandomName(16)

	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	prod := GetProdClient()

	advertiseOnlyAPIDef, api := testutils.GenerateAdvertiseOnlyAPIDefinition(t)

	args := &testutils.InitTestArgs{
		CtlUser:        testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		SrcAPIM:        apim,
		InitFlag:       projectName,
		DefinitionFlag: advertiseOnlyAPIDef,
	}

	testutils.ValidateInitializeProjectWithDefinitionFlag(t, args)

	filePath, _ := filepath.Abs(projectName)
	genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
		Source:      filePath,
		Destination: testutils.GetEnvAPIExportPath(prod.GetEnvName()),
	}

	testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

	importExportArgs := &testutils.ApiImportExportTestArgs{
		ApiProvider:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:        testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:            &api,
		SrcAPIM:        prod,
		DestAPIM:       prod,
		ImportFilePath: filePath,
	}

	// Store the deployment directory path to be provided as the params during import
	importExportArgs.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination,
		api.Name, api.Version)
	testutils.ValidateAPIImportExportWithDeploymentDirForAdvertiseOnlyAPI(t, importExportArgs)
}

// As a tenant user with the Internal/devops role, initialize an API project and import it to an environment (cross tenant
// import since the API definition has the provider as the super tenant admin) with certificates and check whether the
// advertise only properties (advertised, original devportal URL and API owner) have been set correctly. Export the API
// and check whether certificates have not been exported.
func TestInitDeploymentDirImportExportAdvertiseOnlyAPIDevopsTenant(t *testing.T) {
	apim := GetDevClient()
	projectName := base.GenerateRandomName(16)

	devopsUsername := devops.UserName + "@" + TENANT1
	devopsPassword := devops.Password

	prod := GetProdClient()

	advertiseOnlyAPIDef, api := testutils.GenerateAdvertiseOnlyAPIDefinition(t)

	args := &testutils.InitTestArgs{
		CtlUser:        testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		SrcAPIM:        apim,
		InitFlag:       projectName,
		DefinitionFlag: advertiseOnlyAPIDef,
	}

	testutils.ValidateInitializeProjectWithDefinitionFlag(t, args)

	filePath, _ := filepath.Abs(projectName)
	genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
		Source:      filePath,
		Destination: testutils.GetEnvAPIExportPath(prod.GetEnvName()),
	}

	testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

	importExportArgs := &testutils.ApiImportExportTestArgs{
		// Since this a cross tenant import (The initialized API was super tenant admin's, but now importing to a tenant)
		// the preserve provider will be false (OverrideProvider: true). Hence, the provider will be overridden to the CTL user
		ApiProvider:      testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		CtlUser:          testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:              &api,
		SrcAPIM:          prod,
		DestAPIM:         prod,
		ImportFilePath:   filePath,
		OverrideProvider: true,
	}

	// Store the deployment directory path to be provided as the params during import
	importExportArgs.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination,
		api.Name, api.Version)
	testutils.ValidateAPIImportExportWithDeploymentDirForAdvertiseOnlyAPI(t, importExportArgs)
}
