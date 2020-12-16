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

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

//List Environments using apictl
func TestListEnvironmentsDeprecated(t *testing.T) {
	apim := apimClients[0]
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	response, _ := base.Execute(t, "list", "envs")
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, apim.GetEnvName(), "TestListEnvironmentsDeprecated Failed")
}

// Export an API from one environment as a super tenant non admin user (who has API Create and API Publish permissions)
// by specifying the provider name
func TestExportApiNonAdminSuperTenantUserDeprecated(t *testing.T) {
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

	validateAPIExportDeprecated(t, args)
}

// Export an API from one environment and import to another environment as tenant user with
// Internal/devops role by specifying the provider name
func TestExportImportApiDevopsTenantUserDeprecated(t *testing.T) {
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

	validateAPIExportImportDeprecated(t, args)
}

func TestListApisDevopsTenantUserDeprecated(t *testing.T) {
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

	validateAPIsList(t, args)
}

func TestExportApisWithExportApisCommandDeprecated(t *testing.T) {
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

	validateAllApisOfATenantIsExported(t, args, apisAdded)
}

func TestListApiProductsDevopsTenantUserDeprecated(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

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
		CtlUser: testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		SrcAPIM: dev,
	}

	validateAPIProductsList(t, args)
}

func TestExportAppNonAdminSuperTenantDeprecated(t *testing.T) {
	subscriberUserName := subscriber.UserName
	subscriberPassword := subscriber.Password

	dev := apimClients[0]

	app := testutils.AddApp(t, dev, subscriberUserName, subscriberPassword)

	args := &testutils.AppImportExportTestArgs{
		AppOwner:    testutils.Credentials{Username: subscriberUserName, Password: subscriberPassword},
		CtlUser:     testutils.Credentials{Username: subscriberUserName, Password: subscriberPassword},
		Application: app,
		SrcAPIM:     dev,
	}

	validateAppExportFailure(t, args)
}

func TestExportImportAppDevopsTenantDeprecated(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	app := testutils.AddApp(t, dev, tenantAdminUsername, tenantAdminPassword)

	args := &testutils.AppImportExportTestArgs{
		AppOwner:    testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		CtlUser:     testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		Application: app,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	validateAppExportImportWithPreserveOwner(t, args)
}

func TestListAppsDevopsTenantUserDeprecated(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	otherUsername := subscriber.UserName + "@" + TENANT1
	otherPassword := subscriber.Password

	apim := apimClients[0]
	testutils.AddApp(t, apim, tenantAdminUsername, tenantAdminPassword)
	testutils.AddApp(t, apim, otherUsername, otherPassword)

	base.SetupEnv(t, apim.GetEnvName(), apim.GetApimURL(), apim.GetTokenURL())
	base.Login(t, apim.GetEnvName(), tenantDevopsUsername, tenantDevopsPassword)
	listApps(t, apim.GetEnvName())
}

func TestGetKeysNonAdminSuperTenantUserDeprecated(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, api.ID)

	args := &testutils.ApiGetKeyTestArgs{
		CtlUser: testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		Api:     api,
		Apim:    dev,
	}

	validateGetKeysFailure(t, args)
}

func TestGetKeysAdminSuperTenantUserDeprecated(t *testing.T) {
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := apimClients[0]

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, api.ID)

	args := &testutils.ApiGetKeyTestArgs{
		CtlUser: testutils.Credentials{Username: adminUser, Password: adminPassword},
		Api:     api,
		Apim:    dev,
	}

	validateGetKeys(t, args)
}
