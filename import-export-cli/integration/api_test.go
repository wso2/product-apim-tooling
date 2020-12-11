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

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
)

const numberOfAPIs = 5 // Number of APIs to be added in a loop

// Export an API from one environment as a super tenant non admin user (who has API Create and API Publish permissions)
// by specifying the provider name
func TestExportApiNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		Api:         api,
		SrcAPIM:     dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API from one environment and import to another environment as super tenant admin by specifying the provider name
func TestExportImportApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportImport(t, args)
}

// Export an API from one environment and import to another environment as super tenant user with
// Internal/devops role by specifying the provider name
func TestExportImportApiDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportImport(t, args)
}

// Export an API from one environment as tenant non admin user (who has API Create and API Publish permissions)
// by specifying the provider name
func TestExportApiNonAdminTenantUser(t *testing.T) {
	tenantApiPublisher := publisher.UserName + "@" + TENANT1
	tenantApiPublisherPassword := publisher.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantApiPublisher, Password: tenantApiPublisherPassword},
		Api:         api,
		SrcAPIM:     dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API from one environment and import to another environment as tenant admin by specifying the provider name
func TestExportImportApiAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportImport(t, args)
}

// Export an API from one environment and import to another environment as tenant user with
// Internal/devops role by specifying the provider name
func TestExportImportApiDevopsTenantUser(t *testing.T) {
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
	}

	testutils.ValidateAPIExportImport(t, args)
}

// Export an API as super tenant admin without specifying the provider
func TestExportApiAdminSuperTenantUserWithoutProvider(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API as super tenant user with Internal/devops role without specifying the provider
func TestExportApiDevopsSuperTenantUserWithoutProvider(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API as tenant admin without specifying the provider
func TestExportApiAdminTenantUserWithoutProvider(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API as tenant user with Internal/devops role without specifying the provider
func TestExportApiDevopsTenantUserWithoutProvider(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API using a tenant user by specifying the provider name - API is in a different tenant
func TestExportApiAdminTenantUserFromAnotherTenant(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportFailure(t, args)
}

// Export an API using a tenant user with Internal/devops role by specifying the provider name - API is in a different tenant
func TestExportApiDevopsTenantUserFromAnotherTenant(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportFailure(t, args)
}

// Export an API using a tenant user without specifying the provider name - API is in a different tenant
func TestExportApiAdminTenantUserFromAnotherTenantWithoutProvider(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser:  testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:      api,
		SrcAPIM:  dev,
		DestAPIM: prod,
	}

	testutils.ValidateAPIExportFailure(t, args)
}

// Export an API using a tenant user with Internal/devops role without specifying the provider name - API is in a different tenant
func TestExportApiDevopsTenantUserFromAnotherTenantWithoutProvider(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser:  testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		Api:      api,
		SrcAPIM:  dev,
		DestAPIM: prod,
	}

	testutils.ValidateAPIExportFailure(t, args)
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

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider:      testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:          testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		Api:              api,
		OverrideProvider: true,
		SrcAPIM:          dev,
		DestAPIM:         prod,
	}

	testutils.ValidateAPIExport(t, args)

	// Since --preserve-provider=false both the apiProvider and the ctlUser is tenant admin
	args.ApiProvider = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the API to env2 as tenant admin across domains
	testutils.ValidateAPIImport(t, args)
}

// Export an API from one environment as super tenant user with Internal/devops role
// and import to another environment as cross tenant user with Internal/devops role (with preserve-provider=false)
func TestExportImportApiCrossTenantDevopsUserWithoutPreserveProvider(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider:      testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:          testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:              api,
		OverrideProvider: true,
		SrcAPIM:          dev,
		DestAPIM:         prod,
	}

	testutils.ValidateAPIExport(t, args)

	// Since --preserve-provider=false both the apiProvider and the ctlUser is tenant user with Internal/devops role
	args.ApiProvider = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}
	args.CtlUser = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}

	// Import the API to env2 as tenant user with Internal/devops role across domains
	testutils.ValidateAPIImport(t, args)
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

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExport(t, args)

	// Since --preserve-provider=false is not specified, the apiProvider remain as it is and the ctlUser is tenant admin
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the API to env2 as tenant admin across domains
	testutils.ValidateAPIImportFailure(t, args)
}

// Export an API from one environment as super tenant user with Internal/devops role
// and import to another environment as cross tenant user with Internal/devops role (without preserve-provider=false)
func TestExportImportApiCrossTenantDevopsUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	superTenantApiCreator := creator.UserName
	superTenantApiCreatorPassword := creator.Password

	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	dev := apimClients[0]
	prod := apimClients[1]

	api := testutils.AddAPI(t, dev, superTenantApiCreator, superTenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: superTenantApiCreator, Password: superTenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExport(t, args)

	// Since --preserve-provider=false is not specified, the apiProvider remain as it is and the ctlUser is tenant user
	// with Internal/devops role
	args.CtlUser = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}

	// Import the API to env2 as tenant admin across domains
	testutils.ValidateAPIImportFailure(t, args)
}

func TestListApisAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		// Add the API to env1
		testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: adminUsername, Password: adminPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIsList(t, args)
}

func TestListApisDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		// Add the API to env1
		testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIsList(t, args)
}

func TestListApisAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		// Add the API to env1
		testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIsList(t, args)
}

func TestListApisDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		// Add the API to env1
		testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateAPIsList(t, args)
}

func TestDeleteApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	var api *apim.API
	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		api = testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIDelete(t, args)
}

func TestDeleteApiDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	var api *apim.API
	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		api = testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIDelete(t, args)
}

func TestDeleteApiAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]

	var api *apim.API
	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		api = testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)
	}

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIDelete(t, args)
}

func TestDeleteApiDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := apimClients[0]

	var api *apim.API
	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		api = testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)
	}

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIDelete(t, args)
}

func TestDeleteApiSuperTenantUser(t *testing.T) {
	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	var api *apim.API
	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		api = testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	}

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAPIDelete(t, args)
}

func TestDeleteApiWithActiveSubscriptionsSuperTenantUser(t *testing.T) {
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	dev := apimClients[0]

	var api *apim.API

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, adminUser, adminPassword)

	args := &testutils.ApiGetKeyTestArgs{
		CtlUser: testutils.Credentials{Username: adminUser, Password: adminPassword},
		Api:     api,
		Apim:    dev,
	}
	//Publish created API
	testutils.PublishAPI(dev, adminUser, adminPassword, api.ID)

	testutils.ValidateGetKeysWithoutCleanup(t, args)
	//args to delete API
	argsToDelete := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: adminUser, Password: adminPassword},
		Api:     api,
		SrcAPIM: dev,
	}
	base.WaitForIndexing()

	//validate Api with active subscriptions delete failure
	testutils.ValidateAPIDeleteFailure(t, argsToDelete)
}

func TestExportApisWithExportApisCommand(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := apimClients[0]

	var api *apim.API
	var apisAdded = 0
	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		api = testutils.AddAPI(t, dev, tenantAdminUsername, tenantAdminPassword)
		apisAdded++
	}

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, tenantAdminUsername, tenantAdminPassword)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAllApisOfATenantIsExported(t, args, apisAdded)
}

func TestChangeLifeCycleStatusOfApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	//Change life cycle state of Api from CREATED to PUBLISHED
	args := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: adminUsername, Password: adminPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, args)

	//Change life cycle state of Api from PUBLISHED to CREATED
	argsToNextChange := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: adminUsername, Password: adminPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfApiDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	//Change life cycle state of Api from CREATED to PUBLISHED
	args := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, args)

	//Change life cycle state of Api from PUBLISHED to CREATED
	argsToNextChange := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfApiAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	//Change life cycle state of Api from CREATED to PUBLISHED
	args := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, args)

	//Change life cycle state of Api from PUBLISHED to CREATED
	argsToNextChange := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfApiDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	//Change life cycle state of Api from CREATED to PUBLISHED
	args := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, args)

	//Change life cycle state of Api from PUBLISHED to CREATED
	argsToNextChange := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfApiFailWithAUserWithoutPermissions(t *testing.T) {
	subscriberUsername := subscriber.UserName
	subscriberDevopsPassword := subscriber.Password

	apiCreator := creator.UserName + "@" + TENANT1
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	//Change life cycle state of Api from CREATED to PUBLISHED
	args := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: subscriberUsername, Password: subscriberDevopsPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPIFailure(t, args)
}

func TestChangeLifeCycleStatusOfApiWithActiveSubscriptionWithAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	testutils.PublishAPI(dev, adminUsername, adminPassword, api.ID)

	args := &testutils.ApiGetKeyTestArgs{
		CtlUser: testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Api:     api,
		Apim:    dev,
	}

	//Create an active subscription for Api
	testutils.ValidateGetKeysWithoutCleanup(t, args)

	base.WaitForIndexing()

	//Change life cycle state of Api from PUBLISHED to CREATED
	argsToLifeCycleStateChange := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: adminUsername, Password: adminPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, argsToLifeCycleStateChange)
	testutils.UnsubscribeAPI(dev, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)
}

func TestChangeLifeCycleStatusOfApiWithActiveSubscriptionDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]
	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	testutils.PublishAPI(dev, devopsUsername, devopsPassword, api.ID)

	args := &testutils.ApiGetKeyTestArgs{
		CtlUser: testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Api:     api,
		Apim:    dev,
	}

	//Create an active subscription for Api
	testutils.ValidateGetKeysWithoutCleanup(t, args)

	base.WaitForIndexing()

	//Change life cycle state of Api from PUBLISHED to CREATED
	argsToLifeCycleStateChange := &testutils.ApiChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		APIM:          dev,
		Api:           api,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfAPI(t, argsToLifeCycleStateChange)
	testutils.UnsubscribeAPI(dev, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)
}
