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
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
)

const numberOfAPIProducts = 5 // Number of API Products to be added in a loop

// Export an API Product with its dependent APIs from one environment as a super tenant non admin user
func TestExportApiProductNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

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

	dev := apimClients[0]
	prod := apimClients[1]

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

// Export an API Product with its dependent APIs from one environment and import it as super tenant admin
// when dependent APIs are already in that environment and you do not want to update those APIs.
func TestExportImportApiProductAdminSuperTenantUserWithoutImportApis(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]
	prod := apimClients[1]

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

// Export an API Product with its dependent APIs from one environment and import it as super tenant admin
// when dependent APIs are already in that environment and you want to update those APIs.
func TestExportImportApiProductAdminSuperTenantUserWithUpdateApis(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]
	prod := apimClients[1]

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

// Export an API Product with its dependent APIs from one environment and import to another environment freshly as super tenant admin
// and try to update that API Product (without updating dependent APIs)
func TestExportImportApiProductAdminSuperTenantUserWithUpdateApiProduct(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]
	prod := apimClients[1]

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
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       true,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

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
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]
	prod := apimClients[1]

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
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       true,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

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
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]
	prod := apimClients[1]

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
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       true,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

	// Export the API Product as super tenant admin
	testutils.ValidateAPIProductExport(t, args)

	// Since --preserve-provider=false both the apiProductProvider and the ctlUser is tenant admin
	args.ApiProductProvider = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	testutils.ValidateAPIProductImport(t, args)
}

// Export an API Product with its dependent APIs from one environment as super tenant admin and import to another environment freshly
// as cross tenant admin and try to update that API Product (without updating dependent APIs)
func TestExportImportApiProductCrossTenantUserWithUpdateApiProduct(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]
	prod := apimClients[1]

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
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       true,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

	// Export the API Product as super tenant admin
	testutils.ValidateAPIProductExport(t, args)

	// Since --preserve-provider=false both the apiProductProvider and the ctlUser is tenant admin
	args.ApiProductProvider = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	testutils.ValidateAPIProductImport(t, args)

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
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]
	prod := apimClients[1]

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
		ApiProductProvider:   testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		CtlUser:              testutils.Credentials{Username: adminUsername, Password: adminPassword},
		ApiProduct:           apiProduct,
		SrcAPIM:              dev,
		DestAPIM:             prod,
		ImportApisFlag:       true,
		UpdateApisFlag:       false,
		UpdateApiProductFlag: false,
	}

	// Export the API Product as super tenant admin
	testutils.ValidateAPIProductExport(t, args)

	// Since --preserve-provider=false both the apiProductProvider and the ctlUser is tenant admin
	args.ApiProductProvider = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	testutils.ValidateAPIProductImport(t, args)

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

	dev := apimClients[0]

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

func TestListApiProductsAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]

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

func TestDeleteApiProductAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]

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

func TestDeleteApiProductAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName + "@" + TENANT1
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]

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

func TestDeleteApiProductSuperTenantUser(t *testing.T) {
	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]

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
