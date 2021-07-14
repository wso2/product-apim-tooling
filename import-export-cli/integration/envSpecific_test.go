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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoints
// using the params file by a super tenant admin user
func TestEnvironmentSpecificParamsEndpoint(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

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

	assert.Equal(t, apiParams.Environments[0].Endpoints.Production["url"], importedAPI.GetProductionURL())

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"
	apiCopy.SetProductionURL(same)
	importedAPICopy.SetProductionURL(same)

	testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

// Add an API to one environment, export it and re-import it to another environment by setting the retry time out for endpoints
// using the params file by a super tenant admin user
func TestEnvironmentSpecificParamsEndpointRetryTimeout(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

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

	assert.Equal(t, apiParams.Environments[0].Endpoints.Production["url"], importedAPI.GetProductionURL())
	paramConfig := apiParams.Environments[0].Endpoints.Production["config"].(map[interface{}]interface{})

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
	apiCopy.SetProductionURL(same)
	importedAPICopy.SetProductionURL(same)
	apiCopy.SetProductionConfig(map[interface{}]interface{}{})
	importedAPICopy.SetProductionConfig(map[interface{}]interface{}{})

	testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

// Add an API to one environment, export it and re-import it to another environment by disabling the endpoint security
// using the params file by a super tenant admin user
func TestEnvironmentSpecificParamsEndpointSecurityFalse(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

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

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoint security
// (with the security type digest), using the params file by a super tenant user with the Internal/devops role
func TestEnvironmentSpecificParamsEndpointSecurityDigestDevopsSuperTenant(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
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

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoint security
// (with the security type digest), using the params file by a tenant user with the Internal/devops role
func TestEnvironmentSpecificParamsEndpointSecurityDigestDevopsTenant(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
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

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoint security
// (with the security type basic), using the params file by a super tenant user with the Internal/devops role
func TestEnvironmentSpecificParamsEndpointSecurityBasicDevopsSuperTenant(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
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

// Add an API to one environment, export it and re-import it to another environment by overriding the endpoint security
// (with the security type basic), using the params file by a tenant user with the Internal/devops role
func TestEnvironmentSpecificParamsEndpointSecurityBasicDevopsTenant(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
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
