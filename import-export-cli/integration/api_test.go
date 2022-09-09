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

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
)

const numberOfAPIs = 5 // Number of APIs to be added in a loop

// Export an API from one environment and check the structure of the DTO whether it is similat to what is being
// maintained by APICTL
func TestExportApiCompareStruct(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		Api:         api,
		SrcAPIM:     dev,
	}

	testutils.ValidateExportedAPIStructure(t, args)
}

// Export an API from one environment as a super tenant non admin user (who has Internal/publisher role)
// by specifying the provider name
func TestExportApiNonAdminSuperTenantUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)
}

// Export an API from one environment and import to another environment as super tenant user with
// Internal/devops role by specifying the provider name
func TestExportImportApiDevopsSuperTenantUser(t *testing.T) {
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

	testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)
}

// Export an API from one environment as a super tenant user with Internal/publisher role by specifying the provider name
func TestExportApiSuperTenantPublisherUser(t *testing.T) {
	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: apiPublisher, Password: apiPublisherPassword},
		Api:         api,
		SrcAPIM:     dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API from one environment as a tenant user with Internal/publisher role by specifying the provider name
func TestExportApiTenantPublisherUser(t *testing.T) {
	tenantApiPublisher := publisher.UserName + "@" + TENANT1
	tenantApiPublisherPassword := publisher.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantApiPublisher, Password: tenantApiPublisherPassword},
		Api:         api,
		SrcAPIM:     dev,
	}

	testutils.ValidateAPIExport(t, args)
}

// Export an API using a super tenant user who does not have the required scopes (who has the role Internal/subscriber)
func TestExportApiSuperTenantSubscriberUser(t *testing.T) {
	apiSubscriber := subscriber.UserName
	apiSubscriberPassword := subscriber.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: apiSubscriber, Password: apiSubscriberPassword},
		Api:         api,
		SrcAPIM:     dev,
	}

	testutils.ValidateAPIExportFailureUnauthenticated(t, args)
}

// Export an API using a super tenant user who does not have the required scopes (who has the role Internal/subscriber)
func TestExportApiTenantSubscriberUser(t *testing.T) {
	tenantApiSubscriber := subscriber.UserName + "@" + TENANT1
	tenantApiSubscriberPassword := subscriber.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantApiSubscriber, Password: tenantApiSubscriberPassword},
		Api:         api,
		SrcAPIM:     dev,
	}

	testutils.ValidateAPIExportFailureUnauthenticated(t, args)
}

// Export an API from one environment and import to another environment as tenant admin by specifying the provider name
func TestExportImportApiAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)
}

// Export an API from one environment and import to another environment as tenant user with
// Internal/devops role by specifying the provider name
func TestExportImportApiDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantApiCreator := creator.UserName + "@" + TENANT1
	tenantApiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	api := testutils.AddAPI(t, dev, tenantApiCreator, tenantApiCreatorPassword)

	args := &testutils.ApiImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: tenantApiCreator, Password: tenantApiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		Api:         api,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}

	testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)
}

// Export an API as super tenant admin without specifying the provider
func TestExportApiAdminSuperTenantUserWithoutProvider(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

	dev := GetDevClient()
	prod := GetProdClient()

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

// Export an API with the life cycle status as Blocked and import to another environment
// and import update it
func TestExportImportApiBlocked(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)
			api = testutils.ChangeAPILifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID, "Block")

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			importedApi := testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)

			// Change the lifecycle to Published in the prod environment
			testutils.ChangeAPILifeCycle(prod, user.ApiPublisher.Username, user.ApiPublisher.Password, importedApi.ID, "Re-Publish")

			args.Update = true
			testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)
		})
	}
}

// Export an API with the life cycle status as Deprecated and import to another environment
func TestExportImportApiDeprecated(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)
			api = testutils.ChangeAPILifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID, "Deprecate")

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)
		})
	}
}

// Export an API with the life cycle status as Retired and import to another environment
func TestExportImportApiRetired(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishAPI(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID)
			testutils.ChangeAPILifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID, "Deprecate")
			api = testutils.ChangeAPILifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api.ID, "Retire")

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeREST)
		})
	}
}

func TestListApisAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

// APIs listing with JsonArray format
func TestListApisWithJsonArrayFormat(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
				// Add the API to env1
				testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			}

			args := &testutils.ApiImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}

			testutils.ValidateAPIsListWithJsonArrayFormat(t, args)
		})
	}
}

func TestDeleteApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	// Create and Deploy Revision of the above API
	testutils.CreateAndDeployAPIRevision(t, dev, apiPublisher, apiPublisherPassword, api.ID)

	args := &testutils.ApiGetKeyTestArgs{
		CtlUser: testutils.Credentials{Username: adminUser, Password: adminPassword},
		Api:     api,
		Apim:    dev,
	}
	//Publish created API
	testutils.PublishAPI(dev, apiPublisher, apiPublisherPassword, api.ID)

	testutils.ValidateGetKeysWithoutCleanup(t, args, false)
	//args to delete API
	argsToDelete := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: adminUser, Password: adminPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	//validate Api with active subscriptions delete failure
	testutils.ValidateAPIDeleteFailure(t, argsToDelete)
}

func TestExportApisWithExportApisCommand(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := GetDevClient()

	var api *apim.API
	var apisAdded = 0
	for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
		api = testutils.AddAPI(t, dev, tenantAdminUsername, tenantAdminPassword)
		testutils.CreateAndDeployAPIRevision(t, dev, tenantAdminUsername, tenantAdminPassword, api.ID)
		apisAdded++
	}

	// This will be the API that will be deleted by apictl, so no need to do cleaning
	api = testutils.AddAPIWithoutCleaning(t, dev, tenantAdminUsername, tenantAdminPassword)
	testutils.CreateAndDeployAPIRevision(t, dev, tenantAdminUsername, tenantAdminPassword, api.ID)

	args := &testutils.ApiImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Api:     api,
		SrcAPIM: dev,
	}

	testutils.ValidateAllApisOfATenantIsExported(t, args, apisAdded)
}

// Export APIs bunch at once with export apis command and then add new APIs and export APIs once again to check whether
// the new APIs exported
func TestExportApisTwiceWithAfterAddingApis(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			var api *apim.API
			var apisAdded = 0
			for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
				api := testutils.AddAPI(t, dev, user.Admin.Username, user.Admin.Password)
				testutils.CreateAndDeployAPIRevision(t, dev, user.Admin.Username, user.Admin.Password, api.ID)
				apisAdded++
			}

			// This will be the API that will be deleted by apictl, so no need to do cleaning
			api = testutils.AddAPIWithoutCleaning(t, dev, user.Admin.Username, user.Admin.Password)
			testutils.CreateAndDeployAPIRevision(t, dev, user.Admin.Username, user.Admin.Password, api.ID)

			args := &testutils.ApiImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				Api:     api,
				SrcAPIM: dev,
			}

			testutils.ValidateAllApisOfATenantIsExported(t, args, apisAdded)

			// Add new API and deploy
			api = testutils.AddAPI(t, dev, user.Admin.Username, user.Admin.Password)
			testutils.CreateAndDeployAPIRevision(t, dev, user.Admin.Username, user.Admin.Password, api.ID)
			newApiCount := apisAdded + 1

			// Validate again to check whether the newly added API exported properly.
			testutils.ValidateAllApisOfATenantIsExported(t, args, newApiCount)
		})
	}
}

func TestChangeLifeCycleStatusOfApiAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()

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

	dev := GetDevClient()
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

	dev := GetDevClient()
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

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiSubscriber := subscriber.UserName
	apiSubscriberPassword := subscriber.Password

	dev := GetDevClient()

	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	// Create and Deploy Revision of the above API
	testutils.CreateAndDeployAPIRevision(t, dev, apiPublisher, apiPublisherPassword, api.ID)

	testutils.PublishAPI(dev, adminUsername, adminPassword, api.ID)

	// Create an App
	app := testutils.AddApp(t, dev, apiSubscriber, apiSubscriberPassword)

	//Create an active subscription for Api
	testutils.AddSubscription(t, dev, api.ID, app.ApplicationID, testutils.UnlimitedPolicy,
		apiSubscriber, apiSubscriberPassword)

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
}

func TestChangeLifeCycleStatusOfApiWithActiveSubscriptionDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	apiPublisher := publisher.UserName
	apiPublisherPassword := publisher.Password

	apiSubscriber := subscriber.UserName
	apiSubscriberPassword := subscriber.Password

	dev := GetDevClient()

	// Add the API to env
	api := testutils.AddAPI(t, dev, apiCreator, apiCreatorPassword)

	// Create and Deploy Revision of the above API
	testutils.CreateAndDeployAPIRevision(t, dev, apiPublisher, apiPublisherPassword, api.ID)

	testutils.PublishAPI(dev, devopsUsername, devopsPassword, api.ID)

	// Create an App
	app := testutils.AddApp(t, dev, apiSubscriber, apiSubscriberPassword)

	//Create an active subscription for Api
	testutils.AddSubscription(t, dev, api.ID, app.ApplicationID, testutils.UnlimitedPolicy,
		apiSubscriber, apiSubscriberPassword)

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
}

// Export a SOAP API from one environment and import to another environment by specifying the provider name
func TestExportImportSoapApi(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddSoapAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password, testutils.APITypeSoap)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeSoap)

		})
	}
}

// Export a SOAPTOREST API from one environment and import to another environment by specifying the provider name
func TestExportImportSoapToRestApi(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddSoapAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password, testutils.APITypeSoapToRest)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeSoapToRest)

		})
	}
}

// Export a GraphQL API from one environment and import to another environment as by specifying the provider name
func TestExportImportGraphQLApi(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddGraphQLAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeGraphQL)

		})
	}
}

// Export a WebSocket API from one environment and import to another environment by specifying the provider name
func TestExportImportWebSocketApi(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddWebStreamingAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password, testutils.APITypeWebScoket)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeWebScoket)
		})
	}
}

// Export a WebSub/WebHook API from one environment and import to another environment by specifying the provider name
func TestExportImportWebSubApi(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddWebStreamingAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password, testutils.APITypeWebSub)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeWebSub)

		})
	}
}

// Export a Server Sent Events API from one environment and import to another environment by specifying the provider name
func TestExportImportSSEApi(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddWebStreamingAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password, testutils.APITypeSSE)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeSSE)
		})
	}
}

// Export a Web Socket API (that was created using an Async API definition) from one environment and
//import to another environment by specifying the provider name
func TestExportImportWebSocketApiFromAsyncApiDef(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			api := testutils.AddWebStreamingAPIFromAsyncAPIDefinition(t, dev, user.ApiCreator.Username, user.ApiCreator.Password,
				testutils.APITypeWebScoket)

			args := &testutils.ApiImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Api:         api,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}

			testutils.ValidateAPIExportImport(t, args, testutils.APITypeWebScoket)
		})
	}
}

// Import an API and then create a new version of that API by updating the context and version only and import again
func TestCreateNewVersionOfApiByUpdatingVersion(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			apim := GetDevClient()
			projectName := base.GenerateRandomName(16)

			args := &testutils.InitTestArgs{
				CtlUser:   testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:   apim,
				InitFlag:  projectName,
				OasFlag:   testutils.TestSwagger2DefinitionPath,
				APIName:   testutils.DevFirstSwagger2APIName,
				ForceFlag: false,
			}

			//Initialize a project with API definition
			testutils.ValidateInitializeProjectWithOASFlag(t, args)

			//Assert that project import to publisher portal is successful
			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))

			// Read the API definition file in the project
			apiDefinitionFilePath := args.InitFlag + string(os.PathSeparator) + utils.APIDefinitionFileYaml
			apiDefinitionFileContent := testutils.ReadAPIDefinition(t, apiDefinitionFilePath)

			//Change the version
			newVersion := base.GenerateRandomString()
			apiDefinitionFileContent.Data.Version = newVersion

			// Write the modified API definition to the directory
			testutils.WriteToAPIDefinition(t, apiDefinitionFileContent, apiDefinitionFilePath)

			// Import and validate new Api with version change
			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))

			testutils.ValidateApisListWithVersions(t, args, newVersion)
		})
	}
}

// API search using query parameters
func TestApiSearchWithQueryParams(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			var searchQuery string

			// Add set of APIs to env and store api details
			var addedApisList [numberOfAPIs + 1]*apim.API
			for apiCount := 0; apiCount <= numberOfAPIs; apiCount++ {
				// Add the API to env1
				api := testutils.AddAPI(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
				addedApisList[apiCount] = api
			}

			// Add custom API
			customAPI := addedApisList[3]
			customAPI.Name = testutils.CustomAPIName
			customAPI.Version = testutils.CustomAPIVersion
			customAPI.Context = testutils.CustomAPIContext
			dev.AddAPI(t, customAPI, user.ApiCreator.Username, user.ApiCreator.Password, true)

			args := &testutils.ApiImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}

			for i := 0; i < len(addedApisList); i++ {
				apiNameToSearch := addedApisList[i].Name
				apiNameNotToSearch := addedApisList[len(addedApisList)-(i+1)].Name
				searchQuery = fmt.Sprintf("name:%v", apiNameToSearch)

				//Search APIs using query
				testutils.ValidateSearchApisList(t, args, searchQuery, apiNameToSearch, apiNameNotToSearch)

				//Select random context from the added APIs
				apiContextToSearch := addedApisList[i].Context
				apiContextNotToSearch := addedApisList[len(addedApisList)-(i+1)].Context
				searchQuery = fmt.Sprintf("context:%v", apiContextToSearch)

				//Search APIs using query
				testutils.ValidateSearchApisList(t, args, searchQuery, apiContextToSearch, apiContextNotToSearch)
			}

			// Search custom API with name
			searchQuery = fmt.Sprintf("name:%v", testutils.CustomAPIName)
			testutils.ValidateSearchApisList(t, args, searchQuery, testutils.CustomAPIName,
				addedApisList[1].Name)

			// Search custom API with context
			searchQuery = fmt.Sprintf("context:%v", testutils.CustomAPIContext)
			testutils.ValidateSearchApisList(t, args, searchQuery, testutils.CustomAPIContext,
				addedApisList[1].Context)

			// Search custom API with version
			searchQuery = fmt.Sprintf("version:%v", testutils.CustomAPIVersion)
			testutils.ValidateSearchApisList(t, args, searchQuery, testutils.CustomAPIVersion,
				addedApisList[1].Version)

			// Search custom API with version and name
			searchQuery = fmt.Sprintf("version:%v --query name:%v", testutils.CustomAPIVersion, testutils.CustomAPIName)
			testutils.ValidateSearchApisList(t, args, searchQuery, testutils.CustomAPIVersion,
				addedApisList[1].Version)
		})
	}
}
