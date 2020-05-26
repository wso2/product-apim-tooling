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
	"math/rand"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
)

// Export an API Product with its dependent APIs from one environment as a super tenant non admin user
func TestExportApiProductNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the first dependent API to env1
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider: credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:            credentials{username: apiCreator, password: apiCreatorPassword},
		apiProduct:         apiProduct,
		srcAPIM:            dev,
	}

	validateAPIProductExportFailure(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       true,
		updateApisFlag:       false,
		updateApiProductFlag: false,
	}

	t.Cleanup(func() {
		apiList := getAPIs(prod, apiCreator, apiCreatorPassword)
		for _, api := range apiList.List {
			deleteAPI(t, prod, api.ID, apiCreator, apiCreatorPassword)
		}
	})

	// Import the API Product with the dependent APIs to env2
	validateAPIProductExportImportPreserveProvider(t, args)
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
	dependentAPI1ofEnv1, dependentAPI1ofEnv2 := addAPIToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv1.ID)
	publishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv2.ID)

	// Add the second dependent API to env1 and env2
	dependentAPI2ofEnv1, dependentAPI2ofEnv2 := addAPIFromOpenAPIDefinitionToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv1.ID)
	publishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv2.ID)

	// Map the real name of the APIs with the APIs of env1
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1ofEnv1,
		"SwaggerPetstore": dependentAPI2ofEnv1,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       false,
		updateApisFlag:       false,
		updateApiProductFlag: false,
	}

	// Import the API Product without importing the dependent APIs to env2
	validateAPIProductExportImportPreserveProvider(t, args)
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
	dependentAPI1ofEnv1, dependentAPI1ofEnv2 := addAPIToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv1.ID)
	publishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI1ofEnv2.ID)

	// Add the second dependent API to env1 and env2
	dependentAPI2ofEnv1, dependentAPI2ofEnv2 := addAPIFromOpenAPIDefinitionToTwoEnvs(t, dev, prod, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv1.ID)
	publishAPI(prod, apiPublisher, apiPublisherPassword, dependentAPI2ofEnv2.ID)

	// Map the real name of the APIs with the APIs of env1
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1ofEnv1,
		"SwaggerPetstore": dependentAPI2ofEnv1,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       false,
		updateApisFlag:       true,
		updateApiProductFlag: false,
	}

	// Import the API Product without importing the dependent APIs to env2 but updating those
	validateAPIProductExportImportPreserveProvider(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       true,
		updateApisFlag:       false,
		updateApiProductFlag: false,
	}

	t.Cleanup(func() {
		apiList := getAPIs(prod, apiCreator, apiCreatorPassword)
		for _, api := range apiList.List {
			deleteAPI(t, prod, api.ID, apiCreator, apiCreatorPassword)
		}
	})

	// Export the API Product from env 1 and import it to env 2
	validateAPIProductExportImportPreserveProvider(t, args)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.importApisFlag = false
	args.updateApiProductFlag = true

	// Re-import the API Product to env 1 while updating it
	validateAPIProductImportPreserveProviderWithoutCleaningImportedAPIProduct(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       true,
		updateApisFlag:       false,
		updateApiProductFlag: false,
	}

	t.Cleanup(func() {
		apiList := getAPIs(prod, apiCreator, apiCreatorPassword)
		for _, api := range apiList.List {
			deleteAPI(t, prod, api.ID, apiCreator, apiCreatorPassword)
		}
	})

	// Export the API Product from env 1 and import it to env 2
	validateAPIProductExportImportPreserveProvider(t, args)

	// Set the --update-apis flag to update the APIs (it will update the API Product_as well)
	// and make the importApisFlag false
	args.importApisFlag = false
	args.updateApisFlag = true
	// You can make updateApiProductFlag true too - The same behaviour will happen

	// Re-import the API Product to env 1 while updating the API Product and APIs
	validateAPIProductImportPreserveProviderWithoutCleaningImportedAPIProduct(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       true,
		updateApisFlag:       false,
		updateApiProductFlag: false,
	}

	// Export the API Product as super tenant admin
	validateAPIProductExport(t, args)

	t.Cleanup(func() {
		apiList := getAPIs(prod, tenantAdminUsername, tenantAdminPassword)
		for _, api := range apiList.List {
			deleteAPI(t, prod, api.ID, tenantAdminUsername, tenantAdminPassword)
		}
	})

	// Since --preserve-provider=false both the apiProductProvider and the ctlUser is tenant admin
	args.apiProductProvider = credentials{username: tenantAdminUsername, password: tenantAdminPassword}
	args.ctlUser = credentials{username: tenantAdminUsername, password: tenantAdminPassword}

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	validateAPIProductImport(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       true,
		updateApisFlag:       false,
		updateApiProductFlag: false,
	}

	// Export the API Product as super tenant admin
	validateAPIProductExport(t, args)

	t.Cleanup(func() {
		apiList := getAPIs(prod, apiCreator, apiCreatorPassword)
		for _, api := range apiList.List {
			deleteAPI(t, prod, api.ID, apiCreator, apiCreatorPassword)
		}
	})

	// Since --preserve-provider=false both the apiProductProvider and the ctlUser is tenant admin
	args.apiProductProvider = credentials{username: tenantAdminUsername, password: tenantAdminPassword}
	args.ctlUser = credentials{username: tenantAdminUsername, password: tenantAdminPassword}

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	validateAPIProductImport(t, args)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.importApisFlag = false
	args.updateApiProductFlag = true

	// Re-import the API Product to env 1 while updating it
	validateAPIProductImportWithoutCleaningImportedApiProduct(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add the API Product to env1
	apiProduct := addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)

	args := &apiProductImportExportTestArgs{
		apiProductProvider:   credentials{username: apiPublisher, password: apiPublisherPassword},
		ctlUser:              credentials{username: adminUsername, password: adminPassword},
		apiProduct:           apiProduct,
		srcAPIM:              dev,
		destAPIM:             prod,
		importApisFlag:       true,
		updateApisFlag:       false,
		updateApiProductFlag: false,
	}

	// Export the API Product as super tenant admin
	validateAPIProductExport(t, args)

	t.Cleanup(func() {
		apiList := getAPIs(prod, apiCreator, apiCreatorPassword)
		for _, api := range apiList.List {
			deleteAPI(t, prod, api.ID, apiCreator, apiCreatorPassword)
		}
	})

	// Since --preserve-provider=false both the apiProductProvider and the ctlUser is tenant admin
	args.apiProductProvider = credentials{username: tenantAdminUsername, password: tenantAdminPassword}
	args.ctlUser = credentials{username: tenantAdminUsername, password: tenantAdminPassword}

	// Import the API Product with the dependent APIs to env2 as tenant admin across domains
	validateAPIProductImport(t, args)

	// Set the --update-api-product flag to update the existing API Product while importing
	// and make the importApisFlag false
	args.importApisFlag = false
	args.updateApisFlag = true

	// Re-import the API Product to env 1 while updating the API Product and APIs
	validateAPIProductImportWithoutCleaningImportedApiProduct(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add a random number (< 10) of API Products (same copy)
	randomNumber := rand.Intn(10)
	for apiProductCount := 0; apiProductCount <= randomNumber; apiProductCount++ {
		// Add the API Product to env1
		addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &apiProductImportExportTestArgs{
		ctlUser: credentials{username: adminUsername, password: adminPassword},
		srcAPIM: dev,
	}

	validateAPIProductsList(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add a random number (< 10) of API Products (same copy)
	randomNumber := rand.Intn(10)
	for apiProductCount := 0; apiProductCount <= randomNumber; apiProductCount++ {
		// Add the API Product to env1
		addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &apiProductImportExportTestArgs{
		ctlUser: credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		srcAPIM: dev,
	}

	validateAPIProductsList(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add a random number (< 10) of API Products (same copy)
	randomNumber := rand.Intn(10)

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= randomNumber; apiProductCount++ {
		// Add the API Product to env1 (Since we are deleting an API Product later, if the auto clean is enabled an error may occur when
		// trying to delete the already deleted API. So the auto cleaning will be disabled.)
		apiProduct = addAPIProductFromJSONWithoutCleaning(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &apiProductImportExportTestArgs{
		ctlUser:    credentials{username: adminUsername, password: adminPassword},
		apiProduct: apiProduct, // Choose the last API Product from the loop to delete it
		srcAPIM:    dev,
	}

	t.Cleanup(func() {
		apiProductList := getAPIProducts(dev, adminUsername, adminPassword)
		for _, apiProduct := range apiProductList.List {
			deleteAPIProduct(t, dev, apiProduct.ID, adminUsername, adminPassword)
		}
	})

	validateAPIProductDelete(t, args)
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
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add a random number (< 10) of API Products (same copy)
	randomNumber := rand.Intn(10)

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= randomNumber; apiProductCount++ {
		// Add the API Product to env1 (Since we are deleting an API Product later, if the auto clean is enabled an error may occur when
		// trying to delete the already deleted API. So the auto cleaning will be disabled.)
		apiProduct = addAPIProductFromJSONWithoutCleaning(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &apiProductImportExportTestArgs{
		ctlUser:    credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		apiProduct: apiProduct, // Choose the last API Product from the loop to delete it
		srcAPIM:    dev,
	}

	t.Cleanup(func() {
		apiProductList := getAPIProducts(dev, tenantAdminUsername, tenantAdminPassword)
		for _, apiProduct := range apiProductList.List {
			deleteAPIProduct(t, dev, apiProduct.ID, tenantAdminUsername, tenantAdminPassword)
		}
	})

	validateAPIProductDelete(t, args)
}

func TestDeleteApiProductSuperTenantUser(t *testing.T) {
	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	dev := apimClients[0]

	// Add the first dependent API to env1
	dependentAPI1 := addAPI(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI1.ID)

	// Add the second dependent API to env1
	dependentAPI2 := addAPIFromOpenAPIDefinition(t, dev, apiCreator, apiCreatorPassword)
	publishAPI(dev, apiPublisher, apiPublisherPassword, dependentAPI2.ID)

	// Map the real name of the API with the API
	apisList := map[string]*apim.API{
		"PizzaShackAPI":   dependentAPI1,
		"SwaggerPetstore": dependentAPI2,
	}

	// Add a random number (< 10) of API Products (same copy)
	randomNumber := rand.Intn(10)

	var apiProduct *apim.APIProduct
	for apiProductCount := 0; apiProductCount <= randomNumber; apiProductCount++ {
		apiProduct = addAPIProductFromJSON(t, dev, apiPublisher, apiPublisherPassword, apisList)
	}

	args := &apiProductImportExportTestArgs{
		ctlUser:    credentials{username: apiCreator, password: apiCreatorPassword},
		apiProduct: apiProduct, // Choose the last API Product from the loop to delete it
		srcAPIM:    dev,
	}

	validateAPIProductDeleteFailure(t, args)
}
