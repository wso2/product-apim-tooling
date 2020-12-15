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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

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

	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Production["url"], importedAPI.GetProductionURL())

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"
	apiCopy.SetProductionURL(same)
	importedAPICopy.SetProductionURL(same)

	testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

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
	apiCopy.SetProductionURL(same)
	importedAPICopy.SetProductionURL(same)
	apiCopy.SetProductionConfig(map[interface{}]interface{}{})
	importedAPICopy.SetProductionConfig(map[interface{}]interface{}{})

	testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

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

func TestEnvironmentSpecificParamsEndpointSecurityDigest(t *testing.T) {
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
		ParamsFile:  testutils.APISecurityDigestParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadAPIParams(t, args.ParamsFile)

	validateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
}

func TestEnvironmentSpecificParamsEndpointSecurityBasic(t *testing.T) {
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
		ParamsFile:  testutils.APISecurityBasicParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadAPIParams(t, args.ParamsFile)

	validateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
}

func validateEndpointSecurityDefinition(t *testing.T, api *apim.API, apiParams *testutils.APIParams, importedAPI *apim.API) {
	t.Helper()

	assert.Equal(t, strings.ToUpper(apiParams.Environments[0].Configs.Security.Type), importedAPI.EndpointSecurity.Type)
	assert.Equal(t, apiParams.Environments[0].Configs.Security.Username, importedAPI.EndpointSecurity.Username)
	assert.Equal(t, "", importedAPI.EndpointSecurity.Password)

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"

	apiCopy.EndpointSecurity.Type = same
	importedAPICopy.EndpointSecurity.Type = same

	apiCopy.EndpointSecurity.Username = same
	importedAPICopy.EndpointSecurity.Username = same

	apiCopy.EndpointSecurity.Password = same
	importedAPICopy.EndpointSecurity.Password = same

	testutils.ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}
