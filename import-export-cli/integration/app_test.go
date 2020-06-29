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
	"github.com/magiconair/properties/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"testing"
	"time"
)

const numberOfApps = 5 // Number of Applications to be added in a loop

func TestListApp(t *testing.T) {
	username := superAdminUser
	password := superAdminPassword
	otherUsername := subscriber.UserName
	otherPassword := subscriber.Password

	apim := apimClients[0]
	addApp(t, apim, username, password)
	addApp(t, apim, otherUsername, otherPassword)

	base.SetupEnv(t, apim.GetEnvName(), apim.GetApimURL(), apim.GetTokenURL())
	base.Login(t, apim.GetEnvName(), username, password)
	listApps(t, apim.GetEnvName())
}

func TestExportAppNonAdminSuperTenant(t *testing.T) {
	subscriberUserName := subscriber.UserName
	subscriberPassword := subscriber.Password

	dev := apimClients[0]

	app := addApp(t, dev, subscriberUserName, subscriberPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: subscriberUserName, password: subscriberPassword},
		ctlUser:     credentials{username: subscriberUserName, password: subscriberPassword},
		application: app,
		srcAPIM:     dev,
	}

	validateAppExportFailure(t, args)
}

func TestExportAppNonAdminTenant(t *testing.T) {
	subscriberUserName := subscriber.UserName + "@" + TENANT1
	subscriberPassword := subscriber.Password

	dev := apimClients[0]

	app := addApp(t, dev, subscriberUserName, subscriberPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: subscriberUserName, password: subscriberPassword},
		ctlUser:     credentials{username: subscriberUserName, password: subscriberPassword},
		application: app,
		srcAPIM:     dev,
	}

	validateAppExportFailure(t, args)
}

func TestExportImportOwnAppAdminSuperTenant(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	app := addApp(t, dev, adminUsername, adminPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: adminUsername, password: adminPassword},
		ctlUser:     credentials{username: adminUsername, password: adminPassword},
		application: app,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAppExportImportWithPreserveOwner(t, args)
}

//Import an already export App with already generated Keys with --update flag
func TestExportImportOwnAppAdminSuperTenantWithUpdate(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	app := addApp(t, dev, adminUsername, adminPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: adminUsername, password: adminPassword},
		ctlUser:     credentials{username: adminUsername, password: adminPassword},
		application: app,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAppExportImportWithUpdate(t, args)
}

func TestExportImportOtherAppAdminSuperTenant(t *testing.T) {
	otherUsername := subscriber.UserName
	otherPassword := subscriber.Password
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	app := addApp(t, dev, otherUsername, otherPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: otherUsername, password: otherPassword},
		ctlUser:     credentials{username: adminUsername, password: adminPassword},
		application: app,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAppExportImportWithPreserveOwner(t, args)
}

func TestExportImportOwnAppAdminTenant(t *testing.T) {
	adminUsername := superAdminUser + "@" + TENANT1
	adminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	app := addApp(t, dev, adminUsername, adminPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: adminUsername, password: adminPassword},
		ctlUser:     credentials{username: adminUsername, password: adminPassword},
		application: app,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAppExportImportWithPreserveOwner(t, args)
}

func TestExportOtherAppAdminTenant(t *testing.T) {
	otherUsername := subscriber.UserName + "@" + TENANT1
	otherPassword := subscriber.Password
	adminUsername := superAdminUser + "@" + TENANT1
	adminPassword := superAdminPassword

	dev := apimClients[0]
	prod := apimClients[1]

	app := addApp(t, dev, otherUsername, otherPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: otherUsername, password: otherPassword},
		ctlUser:     credentials{username: adminUsername, password: adminPassword},
		application: app,
		srcAPIM:     dev,
		destAPIM:    prod,
	}

	validateAppExportImportWithPreserveOwner(t, args)
}

func TestExportCrossTenantAppAdminTenant(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := apimClients[0]

	app := addApp(t, dev, adminUsername, adminPassword)

	args := &appImportExportTestArgs{
		appOwner:    credentials{username: adminUsername, password: adminPassword},
		ctlUser:     credentials{username: tenantAdminUsername, password: tenantAdminPassword},
		application: app,
		srcAPIM:     dev,
	}

	validateAppExportFailure(t, args)
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
validateAppExportImportWithPreserveOwner
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

func validateAppDelete(t *testing.T, args *appImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())

	// Delete an API of env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	time.Sleep(1 * time.Second)
	appsListBeforeDelete := args.srcAPIM.GetApplications()

	deleteAppByCtl(t, args)

	appsListAfterDelete := args.srcAPIM.GetApplications()
	time.Sleep(1 * time.Second)

	// Validate whether the expected number of API count is there
	assert.Equal(t, appsListBeforeDelete.Count, appsListAfterDelete.Count+1, "Expected number of Applications not deleted")

	// Validate that the delete is a success
	validateApplicationIsDeleted(t, args.application, appsListAfterDelete)
}

//Delete an Application as a super tenant admin
func TestDeleteAppSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	dev := apimClients[0]

	var application *apim.Application
	for appCount := 0; appCount <= numberOfApps; appCount++ {
		application = addApp(t, dev, adminUsername, adminPassword)
	}

	// This will be the Application that will be deleted by apictl, so no need to do cleaning
	application = addApplicationWithoutCleaning(t, dev, adminUsername, adminPassword)

	args := &appImportExportTestArgs{
		ctlUser:     credentials{username: superAdminUser, password: superAdminPassword},
		application: application,
		srcAPIM:     dev,
	}

	validateAppDelete(t, args)
}
