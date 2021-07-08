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

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

const numberOfApps = 5 // Number of Applications to be added in a loop

func TestListApps(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			apim := GetDevClient()
			testutils.AddApp(t, apim, user.CtlUser.Username, user.CtlUser.Password)
			testutils.AddApp(t, apim, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			base.SetupEnv(t, apim.GetEnvName(), apim.GetApimURL(), apim.GetTokenURL())
			base.Login(t, apim.GetEnvName(), user.CtlUser.Username, user.CtlUser.Password)
			testutils.ListApps(t, apim.GetEnvName())
		})
	}
}

// List all the applications in an environment (by specifying the owner)
func TestListAppWithOwner(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword

	apim := GetDevClient()

	for appCount := 0; appCount < 5; appCount++ {
		testutils.AddApp(t, apim, username, password)
	}

	base.SetupEnv(t, apim.GetEnvName(), apim.GetApimURL(), apim.GetTokenURL())
	base.Login(t, apim.GetEnvName(), username, password)

	testutils.ValidateListAppsWithOwner(t, apim.GetEnvName())
}

// Export an application in the same tenant using a non admin user (with Intenal/subscriber role)
func TestExportAppNonAdminUser(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:    testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:     testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				Application: app,
				SrcAPIM:     dev,
			}

			testutils.ValidateAppExport(t, args)
		})
	}

}

// Export an own application in same tenant and imported it to another environment
func TestExportImportOwnApp(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.CtlUser.Username, user.CtlUser.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
			}

			testutils.ValidateAppExportImport(t, args, true)
		})
	}
}

// Export an application belongs to a user with Internal/subscriber role in same tenant and
// import it to another environment with preserve owner
func TestExportImportOtherApp(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
			}

			testutils.ValidateAppExportImport(t, args, true)
		})
	}
}

// Export an application (created by an admin user) and import it to another
// environment while preserving the owner
func TestExportImportAdminApp(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.Admin.Username, user.Admin.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
			}

			testutils.ValidateAppExportImport(t, args, true)
		})
	}
}

// Export an application (created by a subscriber user) with generated keys and import it to another
// environment while preserving the owner by a user with Internal/devops role
func TestExportImportAppWithGeneratedKeys(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
				WithKeys:      true,
			}

			testutils.ValidateAppExportImportGeneratedKeys(t, args, app.ApplicationID, true)
		})
	}
}

// Export an application (created by a subscriber user) with generated keys and import it to another
// environment while preserving the owner by skipping keys
func TestExportImportAppWithGeneratedKeysBySkippingKeys(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
				WithKeys:      true,
				SkipKeys:      true,
			}

			testutils.ValidateAppExportImportGeneratedKeys(t, args, app.ApplicationID, true)

		})
	}
}

// Export an application (created by a subscriber user) with subscriptions and import it to another
// environment while preserving the owner and invoke one API
func TestExportImportAppWithSubscriptions(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			// Add, deploy and publish the first API to env1 and env2
			api1ofEnv1, api1ofEnv2 := testutils.AddAPIToTwoEnvs(t, dev, prod, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.DeployAndPublishAPI(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api1ofEnv1.ID)
			testutils.DeployAndPublishAPI(t, prod, user.ApiPublisher.Username, user.ApiPublisher.Password, api1ofEnv2.ID)

			// Add, deploy and publish the second API to env1 and env2
			api2ofEnv1, api2ofEnv2 := testutils.AddAPIFromOpenAPIDefinitionToTwoEnvs(t, dev,
				prod, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.DeployAndPublishAPI(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api2ofEnv1.ID)
			testutils.DeployAndPublishAPI(t, prod, user.ApiPublisher.Username, user.ApiPublisher.Password, api2ofEnv2.ID)

			// Create active subscriptions for APIs in env1
			testutils.AddSubscription(t, dev, api1ofEnv1.ID, app.ApplicationID, testutils.UnlimitedPolicy,
				user.ApiSubscriber.Username, user.ApiSubscriber.Password)
			testutils.AddSubscription(t, dev, api2ofEnv1.ID, app.ApplicationID, testutils.UnlimitedPolicy,
				user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   testutils.GetApp(t, dev, app.Name, user.ApiSubscriber.Username, user.ApiSubscriber.Password),
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
				WithKeys:      true,
			}

			importedApplication := testutils.ValidateAppExportImportSubscriptions(t, args, app.ApplicationID, false, true)

			// Generate keys for the imported application in env 2
			applicationKey := testutils.GenerateKeys(t, args.DestAPIM, args.AppOwner.Username, args.AppOwner.Password,
				importedApplication.ApplicationID)
			testutils.InvokeAPI(t, testutils.GetResourceURL(args.DestAPIM, api1ofEnv2), applicationKey.Token.AccessToken, 200)

		})
	}
}

// Export an application (created by a subscriber user) with subscriptions and import it to another
// environment while preserving the owner by skipping subscriptions.
// Later add the subscriptions using the update flag.
func TestExportImportAppWithSubscriptionsBySkippingSubscriptions(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			// Add, deploy and publish the first API to env1 and env2
			api1ofEnv1, api1ofEnv2 := testutils.AddAPIToTwoEnvs(t, dev, prod, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.DeployAndPublishAPI(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api1ofEnv1.ID)
			testutils.DeployAndPublishAPI(t, prod, user.ApiPublisher.Username, user.ApiPublisher.Password, api1ofEnv2.ID)

			// Add, deploy and publish the second API to env1 and env2
			api2ofEnv1, api2ofEnv2 := testutils.AddAPIFromOpenAPIDefinitionToTwoEnvs(t, dev,
				prod, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.DeployAndPublishAPI(t, dev, user.ApiPublisher.Username, user.ApiPublisher.Password, api2ofEnv1.ID)
			testutils.DeployAndPublishAPI(t, prod, user.ApiPublisher.Username, user.ApiPublisher.Password, api2ofEnv2.ID)

			// Create active subscriptions for APIs in env1
			testutils.AddSubscription(t, dev, api1ofEnv1.ID, app.ApplicationID, testutils.UnlimitedPolicy,
				user.ApiSubscriber.Username, user.ApiSubscriber.Password)
			testutils.AddSubscription(t, dev, api2ofEnv1.ID, app.ApplicationID, testutils.UnlimitedPolicy,
				user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:          testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:           testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:       testutils.GetApp(t, dev, app.Name, user.ApiSubscriber.Username, user.ApiSubscriber.Password),
				SrcAPIM:           dev,
				DestAPIM:          prod,
				PreserveOwner:     true,
				WithKeys:          true,
				SkipSubscriptions: true,
			}

			// Here the imported application without the subscriptions will get validated
			testutils.ValidateAppExportImportSubscriptions(t, args, app.ApplicationID, false, true)

			// Make skip subscriptions false and update true, so that the imported application in env 2 will get
			// updated with the subscriptions
			args.SkipSubscriptions = false
			args.UpdateFlag = true
			// Here the imported application with the subscriptions will get validated
			testutils.ValidateAppExportImportSubscriptions(t, args, app.ApplicationID, true, false)
		})
	}
}

// Export an application in a cross tenant using an admin tenant user
func TestExportCrossTenantAppAdminTenant(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := GetDevClient()

	app := testutils.AddApp(t, dev, adminUsername, adminPassword)

	args := &testutils.AppImportExportTestArgs{
		AppOwner:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:     testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Application: app,
		SrcAPIM:     dev,
	}

	testutils.ValidateAppExportFailure(t, args)
}

// Export an application (created by a tenant user) by a tenant user with Internal/devops role
func TestExportCrossTenantAppDevopsTenant(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	dev := GetDevClient()

	app := testutils.AddApp(t, dev, adminUsername, adminPassword)

	args := &testutils.AppImportExportTestArgs{
		AppOwner:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:     testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		Application: app,
		SrcAPIM:     dev,
	}

	testutils.ValidateAppExportFailure(t, args)
}

// Export an application (created by a tenant user) by a super tenant user with Internal/devops role
func TestExportCrossTenantAppDevopsSuperTenant(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	dev := GetDevClient()

	app := testutils.AddApp(t, dev, tenantAdminUsername, tenantAdminPassword)

	args := &testutils.AppImportExportTestArgs{
		AppOwner:    testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		CtlUser:     testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		Application: app,
		SrcAPIM:     dev,
	}

	testutils.ValidateAppExportFailure(t, args)
}

// Export an application from one environment and import it as a directory to another environment
func TestImportAppAsDirectory(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
			}

			testutils.ValidateExportAppAndDirectoryImport(t, args, true)
		})
	}
}

// Export an application (created by a subscriber user) and import it to another
// environment while preserving the owner.
// Later update the description and throttling tier.
func TestExportImportAppUpdateMetadata(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
			}

			// Export the application from env 1 and import to env 2
			testutils.ValidateAppExportImport(t, args, false)

			// Update the application description and the throttling policy and import it
			testutils.ValidateAppMetaDataUpdateImport(t, args, true)

		})
	}
}

// TODO: Secondary user store test cases, need to enabled when later on when secondary user store creation is automated
/*
func TestExportAppSecondaryUserStoreAdminSuperTenant(t *testing.T) {
	username := "SECOND.COM/super"
	password := "admin"

	name := "DefaultApplication"
	owner := "SECOND.COM/super"

	base.SetupEnv(t, devEnv, devApim, devTokenEP)
	base.Login(t, devEnv, username, password)

	exportApp(t, name, owner, devEnv)

	assert.True(t, base.IsApplicationArchiveExists(devAppExportPath, name, owner))
}

func TestExportAppSecondaryUserStoreAdminSuperTenantLowerCase(t *testing.T) {
	username := "second.com/super"
	password := "admin"

	name := "DefaultApplication"
	owner := "second.com/super"

	base.SetupEnv(t, devEnv, devApim, devTokenEP)
	base.Login(t, devEnv, username, password)

	exportApp(t, name, owner, devEnv)

	assert.True(t, base.IsApplicationArchiveExists(devAppExportPath, name, owner))
}
*/

// Delete an own/other's Applications
func TestDeleteApp(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			var application *apim.Application
			for appCount := 0; appCount <= numberOfApps; appCount++ {
				application = testutils.AddApp(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)
			}

			// This will be the Application that will be deleted by apictl, so no need to do cleaning
			application = testutils.AddApplicationWithoutCleaning(t, dev, user.ApiSubscriber.Username, user.ApiSubscriber.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:    testutils.Credentials{Username: user.ApiSubscriber.Username, Password: user.ApiSubscriber.Password},
				CtlUser:     testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application: application,
				SrcAPIM:     dev,
			}

			testutils.ValidateAppDelete(t, args)
		})
	}
}


// Export an application with space in application name  and import it to another  to check whether the url
// encoding is working properly
func TestExportImportOwnAppWithSpaceInAppName(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			app := testutils.AddAppWithSpaceInAppName(t, dev, user.CtlUser.Username, user.CtlUser.Password)

			args := &testutils.AppImportExportTestArgs{
				AppOwner:      testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				CtlUser:       testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Application:   app,
				SrcAPIM:       dev,
				DestAPIM:      prod,
				PreserveOwner: true,
			}

			testutils.ValidateAppExportImport(t, args, true)
		})
	}
}
