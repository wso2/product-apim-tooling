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
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func TestEnvironmentSpecificParamsEndpoint(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
		ParamsFile:  testutils.APIEndpointParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadAPIParams(t, args.ParamsFile)

	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Production["url"], importedAPI.GetProductionURL())

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"
	apiCopy.SetProductionURL(same)
	importedAPICopy.SetProductionURL(same)
	apiCopy.SetSandboxURL(same)
	importedAPICopy.SetSandboxURL(same)

	testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

func TestEnvironmentSpecificParamsEndpointRetryTimeout(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
		ParamsFile:  testutils.APIEndpointRetryTimeoutParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadAPIParams(t, args.ParamsFile)

	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Production["url"], importedAPI.GetProductionURL())
	paramConfig := apiParams.Environments[0].Configs.Endpoints.Production["config"].(map[interface{}]interface{})

	apiEndpointConfig := importedAPI.GetProductionConfig()

	for k, v := range paramConfig {
		key := fmt.Sprintf("%v", k)
		value, _ := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)

		assert.Equal(t, value, apiEndpointConfig[key])
	}

	assert.Equal(t, len(paramConfig), len(apiEndpointConfig))

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"
	apiCopy.EndpointConfig = same
	importedAPICopy.EndpointConfig = same

	testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

func TestEnvironmentSpecificParamsEndpointSecurityFalse(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
		ParamsFile:  testutils.APISecurityFalseParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	testutils.ValidateAPIsEqual(t, api, importedAPI)
}

func TestEnvironmentSpecificParamsEndpointSecurityDigest(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
		ParamsFile:  testutils.APISecurityDigestParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadAPIParams(t, args.ParamsFile)

	testutils.ValidateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
}

func TestEnvironmentSpecificParamsEndpointSecurityBasic(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
		ParamsFile:  testutils.APISecurityBasicParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadAPIParams(t, args.ParamsFile)

	testutils.ValidateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
}

func TestExportApiGenDeploymentDirImport(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExport(t, args)

	genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
		Source:      base.ConstructAPIFilePath(testutils.GetEnvAPIExportPath(dev.GetEnvName()), api.Name, api.Version),
		Destination: testutils.GetEnvAPIExportPath(dev.GetEnvName()),
	}

	testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)
	// Store the deployment directory path to be provided as the params during import
	args.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination, api.Name, api.Version)

	// Move dummay params file of an API to the created deployment directory
	srcPathForParamsFile, _ := filepath.Abs(testutils.APIFullParamsFile)
	destPathForParamsFile := args.ParamsFile + "/" + utils.ParamFile
	utils.CopyFile(srcPathForParamsFile, destPathForParamsFile)

	srcPathForCertificatesDirectory, _ := filepath.Abs(testutils.CertificatesDirectoryPath)
	utils.CopyDirectoryContents(srcPathForCertificatesDirectory, args.ParamsFile+"/"+utils.DeploymentCertificatesDirectory)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadAPIParams(t, args.ParamsFile+"/"+utils.ParamFile)
	testutils.ValidateAPIParamsWithoutCerts(t, apiParams, importedAPI)

	args.SrcAPIM = prod // The API should be exported from prod env
	testutils.ValidateExportedAPICerts(t, apiParams, importedAPI, args)
}
