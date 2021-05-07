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
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
)

const numberOfAPIProducts = 5 // Number of API Products to be added in a loop

// Export an API Product with its dependent APIs from one environment as a super tenant non admin user
func TestExportApiProductNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProductProvider: testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:            testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		ApiProduct:         apiProduct,
		SrcAPIM:            dev,
	}

	testutils.ValidateAPIProductExportFailure(t, args)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as super tenant admin
func TestExportImportApiProductAdminSuperTenantUserWithImportApis(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	apiProviders := map[string]testutils.Credentials{}
	apiProviders[dependentAPI1.Name] = testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword}
	apiProviders[dependentAPI2.Name] = testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword}

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProviders:         apiProviders,
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       true,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

	// Import the API Product with the dependent APIs to env2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as super tenant user
// with Internal/devops role
func TestExportImportApiProductDevopsSuperTenantUserWithImportApis(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	apiProviders := map[string]testutils.Credentials{}
	apiProviders[dependentAPI1.Name] = testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword}
	apiProviders[dependentAPI2.Name] = testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword}

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProviders:         apiProviders,
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       true,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

	// Import the API Product with the dependent APIs to env2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import it as super tenant admin
// when dependent APIs are already in that environment and you do not want to update those APIs.
func TestExportImportApiProductAdminSuperTenantUserWithoutImportApis(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1 and env2
	dependentAPI1ofEnv1, dependentAPI1ofEnv2 := testutils.AddAPIToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv2.ID)

	// Add the second dependent API to env1 and env2
	dependentAPI2ofEnv1, dependentAPI2ofEnv2 := testutils.AddAPIFromOpenAPIDefinitionToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv2.ID)

	// Map the real name of the APIs with the APIs of env1
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1ofEnv1,
		"SwaggerPetstore": dependentAPI2ofEnv1,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       false,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

	// Import the API Product without importing the dependent APIs to env2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import it as super tenant user with Internal/devops
// role when dependent APIs are already in that environment and you do not want to update those APIs.
func TestExportImportApiProductDevopsSuperTenantUserWithoutImportApis(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1 and env2
	dependentAPI1ofEnv1, dependentAPI1ofEnv2 := testutils.AddAPIToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv2.ID)

	// Add the second dependent API to env1 and env2
	dependentAPI2ofEnv1, dependentAPI2ofEnv2 := testutils.AddAPIFromOpenAPIDefinitionToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv2.ID)

	// Map the real name of the APIs with the APIs of env1
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1ofEnv1,
		"SwaggerPetstore": dependentAPI2ofEnv1,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       false,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

	// Import the API Product without importing the dependent APIs to env2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import it as super tenant admin
// when dependent APIs are already in that environment and you want to update those APIs.
func TestExportImportApiProductAdminSuperTenantUserWithUpdateApis(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1 and env2
	dependentAPI1ofEnv1, dependentAPI1ofEnv2 := testutils.AddAPIToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv2.ID)

	// Add the second dependent API to env1 and env2
	dependentAPI2ofEnv1, dependentAPI2ofEnv2 := testutils.AddAPIFromOpenAPIDefinitionToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv2.ID)

	// Map the real name of the APIs with the APIs of env1
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1ofEnv1,
		"SwaggerPetstore": dependentAPI2ofEnv1,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       false,
		UpdateApisFlag:       true,
		UpdateApiProductFlag: false,
	}

	// Import the API Product without importing the dependent APIs to env2 but updating those
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import it as super tenant user with Internal/devops
// role when dependent APIs are already in that environment and you want to update those APIs.
func TestExportImportApiProductDevopsSuperTenantUserWithUpdateApis(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()
	prod := GetProdClient()

	// Add the first dependent API to env1 and env2
	dependentAPI1ofEnv1, dependentAPI1ofEnv2 := testutils.AddAPIToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv2.ID)

	// Add the second dependent API to env1 and env2
	dependentAPI2ofEnv1, dependentAPI2ofEnv2 := testutils.AddAPIFromOpenAPIDefinitionToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv1.ID)
	testutils.PublishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv2.ID)

	// Map the real name of the APIs with the APIs of env1
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1ofEnv1,
		"SwaggerPetstore": dependentAPI2ofEnv1,
	}

	// Add the API Product to env1
	apiProduct := testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       false,
		UpdateApisFlag:       true,
		UpdateApiProductFlag: false,
	}

	// Import the API Product without importing the dependent APIs to env2 but updating those
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as super tenant admin
// and try to update that API Product (without updating dependent APIs)
func TestExportImportApiProductAdminSuperTenantUserWithUpdateApiProduct(t *testing.T) {
	adminUser := testutils.Credentials{Username: superAdminUser, Password: superAdminPassword}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = adminUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product from env 1 and import it to env 2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApiProductFlag = true

	// Re-import the API Product to env 1 while updating it
	testutils.ValidateAPIProductImportUpdatePreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as super tenant user with
// Internal/devops role and try to update that API Product (without updating dependent APIs)
func TestExportImportApiProductDevopsSuperTenantUserWithUpdateApiProduct(t *testing.T) {
	devopsUser := testutils.Credentials{Username: devops.UserName, Password: devops.Password}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = devopsUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product from env 1 and import it to env 2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApiProductFlag = true

	// Re-import the API Product to env 1 while updating it
	testutils.ValidateAPIProductImportUpdatePreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as super tenant admin
// and try to update that API Product and dependent APIs.
// This same command can be used to update only the dependent APIs as well.
func TestExportImportApiProductAdminSuperTenantUserWithUpdateApisAndApiProduct(t *testing.T) {
	adminUser := testutils.Credentials{Username: superAdminUser, Password: superAdminPassword}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = adminUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product from env 1 and import it to env 2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)

	// Set the --update-apis flag to update the APIs (it will update the API Product_as well)
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApisFlag = true
	// You can make updateApiProductFlag true too - The same behaviour will happen

	// Re-import the API Product to env 1 while updating the API Product and APIs
	testutils.ValidateAPIProductImportUpdatePreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as super tenant user with Internal/devops
// role and try to update that API Product and dependent APIs.
// This same command can be used to update only the dependent APIs as well.
func TestExportImportApiProductDevopsSuperTenantUserWithUpdateApisAndApiProduct(t *testing.T) {
	devopsUser := testutils.Credentials{Username: devops.UserName, Password: devops.Password}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = devopsUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product from env 1 and import it to env 2
	testutils.ValidateAPIProductExportImportPreserveProvider(t, args)

	// Set the --update-apis flag to update the APIs (it will update the API Product_as well)
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApisFlag = true
	// You can make updateApiProductFlag true too - The same behaviour will happen

	// Re-import the API Product to env 1 while updating the API Product and APIs
	testutils.ValidateAPIProductImportUpdatePreserveProvider(t, args)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as cross tenant admin
func TestExportImportApiProductCrossTenantUserWithImportApis(t *testing.T) {
	adminUser := testutils.Credentials{Username: superAdminUser, Password: superAdminPassword}

	tenantAdminUser := testutils.Credentials{Username: superAdminUser + "@" + TENANT1, Password: superAdminPassword}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = adminUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product as super tenant admin
	testutils.ValidateAPIProductExport(t, args)

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	args.CtlUser = tenantAdminUser
	testutils.ValidateAPIProductImport(t, args, false)
}

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as cross tenant user
// with Internal/devops role
func TestExportImportApiProductCrossTenantDevopsWithImportApis(t *testing.T) {
	devopsUser := testutils.Credentials{Username: devops.UserName, Password: devops.Password}

	tenantDevopsUser := testutils.Credentials{Username: devops.UserName + "@" + TENANT1, Password: devops.Password}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = devopsUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product as super tenant user with Internal/devops role
	testutils.ValidateAPIProductExport(t, args)

	// Import the API Product with the dependent APIs to env2 as tenant user with Internal/devops role across domains
	args.CtlUser = tenantDevopsUser
	testutils.ValidateAPIProductImport(t, args, false)
}

// Export an API Product with its dependent APIs from one environment as super tenant admin and import to another environment freshly
// as cross tenant admin and try to update that API Product (without updating dependent APIs)
func TestExportImportApiProductCrossTenantUserWithUpdateApiProduct(t *testing.T) {
	adminUser := testutils.Credentials{Username: superAdminUser, Password: superAdminPassword}

	tenantAdminUser := testutils.Credentials{Username: superAdminUser + "@" + TENANT1, Password: superAdminPassword}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = adminUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product as super tenant admin
	testutils.ValidateAPIProductExport(t, args)

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	args.CtlUser = tenantAdminUser
	testutils.ValidateAPIProductImport(t, args, false)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApiProductFlag = true

	// Re-import the API Product to env 1 while updating it
	testutils.ValidateAPIProductImportUpdate(t, args)
}

// Export an API Product with its dependent APIs from one environment as super tenant user with Internal/devops role
// and import to another environment freshly as cross tenant admin and try to update that API Product (without updating dependent APIs)
func TestExportImportApiProductCrossTenantDevopsWithUpdateApiProduct(t *testing.T) {
	devopsUser := testutils.Credentials{Username: devops.UserName, Password: devops.Password}

	tenantDevopsUser := testutils.Credentials{Username: devops.UserName + "@" + TENANT1, Password: devops.Password}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = devopsUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product as super tenant user with Internal/devops role
	testutils.ValidateAPIProductExport(t, args)

	// Import the API Product with the dependent APIs to env2 as tenant user with Internal/devops role across domains
	args.CtlUser = tenantDevopsUser
	testutils.ValidateAPIProductImport(t, args, false)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApiProductFlag = true

	// Re-import the API Product to env 1 while updating it
	testutils.ValidateAPIProductImportUpdate(t, args)
}

// Export an API Product with its dependent APIs from one environment as super admin and import to another environment freshly as tenant admin
// and try to update that API Product and dependent APIs.
// This same command can be used to update only the dependent APIs as well.
func TestExportImportApiProductCrossTenantUserWithUpdateApisAndApiProduct(t *testing.T) {
	adminUser := testutils.Credentials{Username: superAdminUser, Password: superAdminPassword}

	tenantAdminUser := testutils.Credentials{Username: superAdminUser + "@" + TENANT1, Password: superAdminPassword}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = adminUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product as super tenant admin
	testutils.ValidateAPIProductExport(t, args)

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	args.CtlUser = tenantAdminUser
	testutils.ValidateAPIProductImport(t, args, false)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApisFlag = true

	// Re-import the API Product to env 1 while updating the API Product and APIs
	testutils.ValidateAPIProductImportUpdate(t, args)
}

// Export an API Product with its dependent APIs from one environment as super tenant user with Internal/devops role
//  and import to another environment freshly as tenant admin and try to update that API Product and dependent APIs.
// This same command can be used to update only the dependent APIs as well.
func TestExportImportApiProductCrossTenantDevopsWithUpdateApisAndApiProduct(t *testing.T) {
	devopsUser := testutils.Credentials{Username: devops.UserName, Password: devops.Password}

	tenantDevopsUser := testutils.Credentials{Username: devops.UserName + "@" + TENANT1, Password: devops.Password}

	apiCreator := testutils.Credentials{Username: creator.UserName, Password: creator.Password}

	apiPublisher := testutils.Credentials{Username: publisher.UserName, Password: publisher.Password}

	dev := GetDevClient()
	prod := GetProdClient()

	args := testutils.AddAPIProductWithTwoDependentAPIs(t, dev, &apiCreator, &apiPublisher)

	args.CtlUser = devopsUser
	args.DestAPIM = prod
	args.ImportApisFlag = true
	args.UpdateApisFlag = false
	args.UpdateApiProductFlag = false

	// Export the API Product as super tenant admin
	testutils.ValidateAPIProductExport(t, args)

	// Import the API Product with the dependent APIs to env2 as tenant user with Internal/devops role across domains
	args.CtlUser = tenantDevopsUser
	testutils.ValidateAPIProductImport(t, args, false)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.ImportApisFlag = false
	args.UpdateApisFlag = true

	// Re-import the API Product to env 1 while updating the API Product and APIs
	testutils.ValidateAPIProductImportUpdate(t, args)
}

func TestListApiProductsAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		// Add the API Product to env1
		testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: adminUsername, Password: adminPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIProductsList(t, args)
}

func TestListApiProductsDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		// Add the API Product to env1
		testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIProductsList(t, args)
}

func TestListApiProductsAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		// Add the API Product to env1
		testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIProductsList(t, args)
}

func TestListApiProductsDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		// Add the API Product to env1
		testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIProductsList(t, args)
}

func TestDeleteApiProductAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		apiProduct = testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	// This will be the API Product that will be deleted by apictl, so no need to do cleaning
	apiProduct = testutils.AddAPIProductFromJSONWithoutCleaning(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct: apiProduct,
		SrcAPIM:    dev,
	}

	testutils.ValidateAPIProductDelete(t, args)
}

func TestDeleteApiProductDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		apiProduct = testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	// This will be the API Product that will be deleted by apictl, so no need to do cleaning
	apiProduct = testutils.AddAPIProductFromJSONWithoutCleaning(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser:    testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		ApiProduct: apiProduct,
		SrcAPIM:    dev,
	}

	testutils.ValidateAPIProductDelete(t, args)
}

func TestDeleteApiProductAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		apiProduct = testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	// This will be the API Product that will be deleted by apictl, so no need to do cleaning
	apiProduct = testutils.AddAPIProductFromJSONWithoutCleaning(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser:    testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		ApiProduct: apiProduct,
		SrcAPIM:    dev,
	}

	testutils.ValidateAPIProductDelete(t, args)
}

func TestDeleteApiProductDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		apiProduct = testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	// This will be the API Product that will be deleted by apictl, so no need to do cleaning
	apiProduct = testutils.AddAPIProductFromJSONWithoutCleaning(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser:    testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		ApiProduct: apiProduct,
		SrcAPIM:    dev,
	}

	testutils.ValidateAPIProductDelete(t, args)
}

func TestDeleteApiProductSuperTenantUser(t *testing.T) {
	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
		apiProduct = testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &testutils.ApiProductImportExportTestArgs{
		CtlUser:    testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		ApiProduct: apiProduct, // Choose the last API Product from the loop to delete it
		SrcAPIM:    dev,
	}

	testutils.ValidateAPIProductDeleteFailure(t, args)
}

func TestDeleteApiProductWithActiveSubscriptionsSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := GetDevClient()

	// Add the first dependent API to env1
	dependentAPI1 := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	var apiProduct *apim.APIProduct
	apiProduct = testutils.AddAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	// This will be the API Product that will be deleted by apictl, so no need to do cleaning
	apiProduct = testutils.AddAPIProductFromJSONWithoutCleaning(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &testutils.ApiGetKeyTestArgs{
		CtlUser:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct: apiProduct,
		Apim:       dev,
	}
	base.WaitForIndexing()

	//Get keys for ApiProduct and keep subscription active
	testutils.ValidateGetKeysWithoutCleanup(t, args)

	argsToDelete := &testutils.ApiProductImportExportTestArgs{
		CtlUser:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct: apiProduct,
		SrcAPIM:    dev,
	}

	testutils.ValidateAPIProductDeleteFailureWithExistingEnv(t, argsToDelete)
}

// API products search using query parameters for Super tenant admin user
func TestApiProductSearchWithQueryParamsAdminSuperTenantUser(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			var searchingApiProductIndex, redundantApiProductIndex int
			var searchQuery string

			maxIndexOfTheArray := 5
			minIndexOfTheArray := 0

			// Add the first dependent API to env1
			dependentAPI1 := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, dependentAPI1.ID)

			// Add the second dependent API to env1
			dependentAPI2 := testutils.AddAPIFromOpenAPIDefinition(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, dependentAPI2.ID)

			// Map the real name of the API with the API
			apisList := map[string]*apim.API{
				"PizzaShackAPI":   dependentAPI1,
				"SwaggerPetstore": dependentAPI2,
			}

			// Add set of API Products to env and store api details
			var addedApiProductsList [numberOfAPIProducts + 1]*apim.APIProduct
			for apiProductCount := 0; apiProductCount <= numberOfAPIProducts; apiProductCount++ {
				// Add the API Product to env1
				apiProduct := testutils.AddAPIProductFromJSON(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, apisList)
				addedApiProductsList[apiProductCount] = apiProduct
			}

			args := &testutils.ApiProductImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}

			//Select random name from the added API Products
			searchingApiProductIndex = base.GenerateRandomNumber(minIndexOfTheArray, maxIndexOfTheArray)
			redundantApiProductIndex = maxIndexOfTheArray - searchingApiProductIndex
			apiProductNameToSearch := addedApiProductsList[searchingApiProductIndex].Name
			apiProductNameNotToSearch := addedApiProductsList[redundantApiProductIndex].Name
			searchQuery = fmt.Sprintf("name:%v", apiProductNameToSearch)

			//Search API Products using query with name
			testutils.ValidateSearchApiProductsList(t, args, searchQuery, apiProductNameToSearch, apiProductNameNotToSearch)

			//Select random context from the added API Products
			searchingApiProductIndex = base.GenerateRandomNumber(minIndexOfTheArray, maxIndexOfTheArray)
			redundantApiProductIndex = maxIndexOfTheArray - searchingApiProductIndex
			apiProductContextToSearch := addedApiProductsList[searchingApiProductIndex].Context
			apiProductContextNotToSearch := addedApiProductsList[redundantApiProductIndex].Context
			searchQuery = fmt.Sprintf("context:%v", apiProductContextToSearch)

			//Search API Products using query with context
			testutils.ValidateSearchApiProductsList(t, args, searchQuery, apiProductContextToSearch, apiProductContextNotToSearch)
		})
	}
}
