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

// Export an API from one environment as a super tenant non admin user (who has API Create and API Publish permissions)
// by specifying the provider name
func TestExportApiNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: apiCreator, password: apiCreatorPassword},
		ctlUser:     credentials{username: apiPublisher, password: apiPublisherPassword},
		api:         api,
		srcAPIM:     dev,
	}

	validateAPIExportFailure(t, args)
}

// Export an API from one environment and import to another environment as super tenant admin by specifying the provider name
func TestExportImportApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: apiCreator, password: apiCreatorPassword},
		ctlUser:     credentials{username: adminUsername, password: adminPassword},
		api:         api,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAPIExportImport(t, args)
}

// Export an API from one environment as tenant non admin user (who has API Create and API Publish permissions)
// by specifying the provider name
func TestExportApiNonAdminTenantUser(t *testing.T) {
	tenantApiPublisher := publisher.UserName + "@" + TENANT1
	tenantApiPublisherPassword := publisher.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := addAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: tenantApiCreator, password: tenantApiCreatorPassword},
		ctlUser:     credentials{username: tenantApiPublisher, password: tenantApiPublisherPassword},
		api:         api,
		srcAPIM:     dev,
	}

	validateAPIExportFailure(t, args)
}

// Export an API from one environment and import to another environment as tenant admin by specifying the provider name
func TestExportImportApiAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: tenantApiCreator, password: tenantApiCreatorPassword},
		ctlUser:     credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		api:         api,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAPIExportImport(t, args)
}

// Export an API as super tenant admin without specifying the provider
func TestExportApiAdminSuperTenantUserWithoutProvider(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &apiImportExportTestArgs{
		ctlUser:  credentials{username: adminUsername, password: adminPassword},
		api:      api,
		srcAPIM:  dev,
		destAPIM: prod,
	}

	validateAPIExport(t, args)
}

// Export an API as tenant admin without specifying the provider
func TestExportApiAdminTenantUserWithoutProvider(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &apiImportExportTestArgs{
		ctlUser:  credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		api:      api,
		srcAPIM:  dev,
		destAPIM: prod,
	}

	validateAPIExport(t, args)
}

// Export an API using a tenant user by specifying the provider name - API is in a different tenant
func TestExportApiAdminTenantUserFromAnotherTenant(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: superTenantApiCreator, password: superTenantApiCreatorPassword},
		ctlUser:     credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		api:         api,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAPIExportFailure(t, args)
}

// Export an API using a tenant user without specifying the provider name - API is in a different tenant
func TestExportApiAdminTenantUserFromAnotherTenantWithoutProvider(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &apiImportExportTestArgs{
		ctlUser:  credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		api:      api,
		srcAPIM:  dev,
		destAPIM: prod,
	}

	validateAPIExportFailure(t, args)
}

// Export an API from one environment as super tenant admin and import to another environment as cross tenant admin
// (with preserve-provider=false)
func TestExportImportApiCrossTenantUserWithoutPreserveProvider(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: superTenantApiCreator, password: superTenantApiCreatorPassword},
		ctlUser:     credentials{username: superTenantAdminUsername, password: superTenantAdminPassword},
		api:         api,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAPIExport(t, args)

	// Since --preserve-provider=false both the apiProvider and the ctlUser is tenant admin
	args.apiProvider = credentials{username: tenantAdminUsername, password: tenantAdminPassword}
	args.ctlUser = credentials{username: tenantAdminUsername, password: tenantAdminPassword}

	// Import the API to env2 as tenant admin across domains
	validateAPIImport(t, args)
}

// Export an API from one environment as super tenant admin and import to another environment as cross tenant admin
// (without preserve-provider=false)
func TestExportImportApiCrossTenantUser(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	api := addAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &apiImportExportTestArgs{
		apiProvider: credentials{username: superTenantApiCreator, password: superTenantApiCreatorPassword},
		ctlUser:     credentials{username: superTenantAdminUsername, password: superTenantAdminPassword},
		api:         api,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAPIExport(t, args)

	// Since --preserve-provider=false is not specified, the apiProvider remain as it is and the ctlUser is tenant admin
	args.ctlUser = credentials{username: tenantAdminUsername, password: tenantAdminPassword}

	// Import the API to env2 as tenant admin across domains
	validateAPIImportFailure(t, args)
}

func TestListApisAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add a random number (< 10) of API (same copy)
	randomNumber := rand.Intn(10)
	for apiCount := 0; apiCount <= randomNumber; apiCount++ {
		// Add the API to env1
		addAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &apiImportExportTestArgs{
		ctlUser: credentials{username: adminUsername, password: adminPassword},
		srcAPIM: dev,
	}

	validateAPIsList(t, args)
}

func TestListApisAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add a random number (< 10) of API (same copy)
	randomNumber := rand.Intn(10)
	for apiCount := 0; apiCount <= randomNumber; apiCount++ {
		// Add the API to env1
		addAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &apiImportExportTestArgs{
		ctlUser: credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		srcAPIM: dev,
	}

	validateAPIsList(t, args)
}

func TestDeleteApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add a random number (< 10) of APIs (same copy)
	randomNumber := rand.Intn(10)

	var api *apim.API
	for apiCount := 0; apiCount <= randomNumber; apiCount++ {
		// Add the API to env1 (Since we are deleting an API later, if the auto clean is enabled an error may occur when
		// trying to delete the already deleted API. So the auto cleaning will be disabled.)
		api = addAPIWithoutCleaning(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &apiImportExportTestArgs{
		ctlUser: credentials{username: adminUsername, password: adminPassword},
		api:     api, // Choose the last API from the loop to delete it
		srcAPIM: dev,
	}

	t.Cleanup(func() {
		apiList := getAPIs(dev, adminUsername, adminPassword)
		for _, api := range apiList.List {
			deleteAPI(t, dev, api.ID, adminUsername, adminPassword)
		}
	})

	validateAPIDelete(t, args)
}

func TestDeleteApiAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add a random number (< 10) of APIs (same copy)
	randomNumber := rand.Intn(10)

	var api *apim.API
	for apiCount := 0; apiCount <= randomNumber; apiCount++ {
		// Add the API to env1 (Since we are deleting an API later, if the auto clean is enabled an error may occur when
		// trying to delete the already deleted API. So the auto cleaning will be disabled.)
		api = addAPIWithoutCleaning(t, dev, tenantApiCreator, tenantApiCreatorPassword)
	}

	args := &apiImportExportTestArgs{
		ctlUser: credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		api:     api, // Choose the last API from the loop to delete it
		srcAPIM: dev,
	}

	t.Cleanup(func() {
		apiList := getAPIs(dev, tenantAdminUsername, tenantAdminPassword)
		for _, api := range apiList.List {
			deleteAPI(t, dev, api.ID, tenantAdminUsername, tenantAdminPassword)
		}
	})

	validateAPIDelete(t, args)
}

func TestDeleteApiSuperTenantUser(t *testing.T) {
	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add a random number (< 10) of APIs (same copy)
	randomNumber := rand.Intn(10)

	var api *apim.API
	for apiCount := 0; apiCount <= randomNumber; apiCount++ {
		// Add the API to env1 (Since we are deleting an API later, if the auto clean is enabled an error may occur when
		// trying to delete the already deleted API. So the auto cleaning will be disabled.)
		api = addAPIWithoutCleaning(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &apiImportExportTestArgs{
		ctlUser: credentials{username: apiCreator, password: apiCreatorPassword},
		api:     api, // Choose the last API from the loop to delete it
		srcAPIM: dev,
	}

	t.Cleanup(func() {
		apiList := getAPIs(dev, apiCreator, apiCreatorPassword)
		for _, api := range apiList.List {
			deleteAPI(t, dev, api.ID, apiCreator, apiCreatorPassword)
		}
	})

	validateAPIDelete(t, args)
}
