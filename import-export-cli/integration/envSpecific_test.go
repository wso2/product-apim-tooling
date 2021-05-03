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
	"os"
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

	apiParams := testutils.ReadParams(t, args.ParamsFile)

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

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Production["url"], importedAPI.GetProductionURL())
	paramConfig := apiParams.Environments[0].Configs.Endpoints.Production["config"].(map[interface{}]interface{})

	apiEndpointConfig := importedAPI.GetProductionConfig()

	for k, v := range paramConfig {
		key := fmt.Sprintf("%v", k)
		value := fmt.Sprintf("%v", v)
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

	assert.Equal(t, false, importedAPI.GetProductionSecurityConfig()["enabled"])
	assert.Equal(t, false, importedAPI.GetSandboxSecurityConfig()["enabled"])

	api.EndpointConfig.(map[string]interface{})["endpoint_security"] = "override_with_the_same_value"
	importedAPI.EndpointConfig.(map[string]interface{})["endpoint_security"] = "override_with_the_same_value"

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

	apiParams := testutils.ReadParams(t, args.ParamsFile)

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

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	testutils.ValidateEndpointSecurityDefinition(t, api, apiParams, importedAPI)
}

//  Import an API with the external params file that has HTTP/REST endpoints without load balancing or failover configs
func TestEnvironmentSpecificParamsHttpRestEndpointWithoutLoadBalancingOrFailover(t *testing.T) {
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
		ParamsFile:  testutils.APIHttpRestEndpointWithoutLoadBalancingOrFailoverParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	testutils.ValidateHttpEndpointWithoutLoadBalancingAndFailover(t, api, apiParams, importedAPI)
}

//  Import an API with the external params file that has HTTP/SOAP endpoints without load balancing or failover configs
func TestEnvironmentSpecificParamsHttpSoapEndpointWithoutLoadBalancingOrFailover(t *testing.T) {
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
		ParamsFile:  testutils.APIHttpSoapEndpointWithoutLoadBalancingOrFailoverParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	testutils.ValidateHttpEndpointWithoutLoadBalancingAndFailover(t, api, apiParams, importedAPI)
}

//  Import an API with the external params file that has HTTP/REST endpoints with load balancing config
func TestEnvironmentSpecificParamsHttpRestEndpointWithLoadBalancing(t *testing.T) {
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
		ParamsFile:  testutils.APIHttpRestEndpointWithLoadBalancingParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	testutils.ValidateHttpEndpointWithLoadBalancing(t, api, apiParams, importedAPI)
}

//  Import an API with the external params file that has HTTP/SOAP endpoints with load balancing config
func TestEnvironmentSpecificParamsHttpSoapEndpointWithLoadBalancing(t *testing.T) {
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
		ParamsFile:  testutils.APIHttpSoapEndpointWithLoadBalancingParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	testutils.ValidateHttpEndpointWithLoadBalancing(t, api, apiParams, importedAPI)
}

//  Import an API with the external params file that has HTTP/REST endpoints with failover config
func TestEnvironmentSpecificParamsHttpRestEndpointWithFailover(t *testing.T) {
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
		ParamsFile:  testutils.APIHttpRestEndpointWithFailoverParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	testutils.ValidateHttpEndpointWithFailover(t, api, apiParams, importedAPI)
}

//  Import an API with the external params file that has HTTP/SOAP endpoints with failover config
func TestEnvironmentSpecificParamsHttpSoapEndpointWithFailover(t *testing.T) {
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
		ParamsFile:  testutils.APIHttpSoapEndpointWithFailoverParamsFile,
	}

	testutils.ValidateAPIExport(t, args)

	importedAPI := testutils.GetImportedAPI(t, args)

	apiParams := testutils.ReadParams(t, args.ParamsFile)

	testutils.ValidateHttpEndpointWithFailover(t, api, apiParams, importedAPI)
}

// Export an API from one environment and generate the deployment directory for that. Import it to another environment with the params
// and certificates. Validate the imported API with the used params. Again, re-export it to validate the certs.
// As a super tenant user with Internal/devops role.
func TestExportApiGenDeploymentDirImportSuperTenant(t *testing.T) {
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
	testutils.ValidateAPIImportExportWithDeploymentDir(t, args)
}

// Export an API from one environment and generate the deployment directory for that. Import it to another environment with the params
// and certificates. Validate the imported API with the used params. Again, re-export it to validate the certs.
// As a tenant user with Internal/devops role.
func TestExportApiGenDeploymentDirImportTenant(t *testing.T) {
	devopsUsername := devops.UserName + "@" + TENANT1
	devopsPassword := devops.Password

	apiCreator := creator.UserName + "@" + TENANT1
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
	testutils.ValidateAPIImportExportWithDeploymentDir(t, args)
}

// Export an API Product from one environment and generate the deployment directory for that. Import it to another environment with the params
// and certificates. Validate the imported API Product with the used params. Again, re-export it to validate the certs.
// As a super tenant user with Internal/devops role.
func TestExportApiProductGenDeploymentDirImportSuperTenant(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)
	os.Setenv("DEPENDENTAPI_2", dependentAPI2.Name+"-"+dependentAPI2.Version)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProductProvider: testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:            testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		ImportApisFlag:     true,
		ApiProduct:         apiProduct,
		SrcAPIM:            dev,
		DestAPIM:           prod,
	}

	testutils.ValidateAPIProductExport(t, args)

	genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
		Source: base.ConstructAPIFilePath(testutils.GetEnvAPIProductExportPath(dev.GetEnvName()), apiProduct.Name,
			utils.DefaultApiProductVersion),
		Destination: testutils.GetEnvAPIProductExportPath(dev.GetEnvName()),
	}

	testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

	// Store the deployment directory path to be provided as the params during import
	args.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination, apiProduct.Name, utils.DefaultApiProductVersion)
	testutils.ValidateAPIProductImportExportWithDeploymentDir(t, args)

	// Validate the dependent API (SwaggerPetstore will be the only one that is in params file of the product)
	testutils.ValidateDependentAPIWithParams(t, dependentAPI2, prod, devopsUsername, devopsPassword)
}

// Export an API Product from one environment and generate the deployment directory for that. Import it to another environment with the params
// and certificates. Validate the imported API Product with the used params. Again, re-export it to validate the certs.
// As a tenant user with Internal/devops role.

func TestExportApiProductGenDeploymentDirImportTenant(t *testing.T) {
	devopsUsername := devops.UserName + "@" + TENANT1
	devopsPassword := devops.Password

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)
	os.Setenv("DEPENDENTAPI_2", dependentAPI2.Name+"-"+dependentAPI2.Version)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProductProvider: testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:            testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		ImportApisFlag:     true,
		ApiProduct:         apiProduct,
		SrcAPIM:            dev,
		DestAPIM:           prod,
	}

	testutils.ValidateAPIProductExport(t, args)

	genDeploymentDirArgs := &testutils.GenDeploymentDirTestArgs{
		Source: base.ConstructAPIFilePath(testutils.GetEnvAPIProductExportPath(dev.GetEnvName()), apiProduct.Name,
			utils.DefaultApiProductVersion),
		Destination: testutils.GetEnvAPIProductExportPath(dev.GetEnvName()),
	}

	testutils.ValidateGenDeploymentDir(t, genDeploymentDirArgs)

	// Store the deployment directory path to be provided as the params during import
	args.ParamsFile = base.ConstructAPIDeploymentDirectoryPath(genDeploymentDirArgs.Destination, apiProduct.Name, utils.DefaultApiProductVersion)
	testutils.ValidateAPIProductImportExportWithDeploymentDir(t, args)

	// Validate the dependent API (SwaggerPetstore will be the only one that is in params file of the product)
	testutils.ValidateDependentAPIWithParams(t, dependentAPI2, prod, devopsUsername, devopsPassword)
}
