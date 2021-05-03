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

package testutils

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func AddApp(t *testing.T, client *apim.Client, username string, password string) *apim.Application {
	client.Login(username, password)
	app := client.GenerateSampleAppData()
	doClean := true
	return client.AddApplication(t, app, username, password, doClean)
}

func AddApplicationWithoutCleaning(t *testing.T, client *apim.Client, username string, password string) *apim.Application {
	client.Login(username, password)
	application := client.GenerateSampleAppData()
	doClean := false
	app := client.AddApplication(t, application, username, password, doClean)
	application = client.GetApplication(app.ApplicationID)
	return application
}

func GenerateKeys(t *testing.T, client *apim.Client, username, password, appId string) apim.ApplicationKey {
	client.Login(username, password)
	generateKeyReq := utils.KeygenRequest{
		KeyType:                 utils.ProductionKeyType,
		GrantTypesToBeSupported: utils.GrantTypesToBeSupported,
		ValidityTime:            utils.DefaultTokenValidityPeriod,
	}
	keyGenResponse := client.GenerateKeys(t, generateKeyReq, appId)
	return keyGenResponse
}

func GetApp(t *testing.T, client *apim.Client, name string, username string, password string) *apim.Application {
	client.Login(username, password)
	appInfo := client.GetApplicationByName(name)
	return client.GetApplication(appInfo.ApplicationID)
}

func GetOauthKeys(t *testing.T, client *apim.Client, username, password string,
	application *apim.Application) *apim.ApplicationKey {
	client.Login(username, password)
	applicationKey := client.GetOauthKeys(t, application)
	return applicationKey
}

func ListApps(t *testing.T, env string) []string {
	response, _ := base.Execute(t, "get", "apps", "-e", env, "-k")

	return base.GetRowsFromTableResponse(response)
}

func ListAppsWithOwner(t *testing.T, env string, owner string) []string {
	response, _ := base.Execute(t, "get", "apps", "-e", env, "-k", "--owner", owner)

	return base.GetRowsFromTableResponse(response)
}

func getEnvAppExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedAppsDirName, envName)
}

func exportApp(t *testing.T, args *AppImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "export", "app", "-n", args.Application.Name, "-o", args.AppOwner.Username,
		"--with-keys="+strconv.FormatBool(args.WithKeys), "-e", args.SrcAPIM.GetEnvName(), "-k", "--verbose")

	t.Cleanup(func() {
		base.RemoveApplicationArchive(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
			args.Application.Name, args.AppOwner.Username)
	})

	return output, err
}

func importApp(t *testing.T, args *AppImportExportTestArgs, doClean bool) (string, error) {
	fileName := base.GetApplicationArchiveFilePath(t, args.SrcAPIM.EnvName, args.Application.Name, args.Application.Owner)
	output, err := base.Execute(t, "import", "app", "-f", fileName, "--preserve-owner="+strconv.FormatBool(args.PreserveOwner),
		"--update="+strconv.FormatBool(args.UpdateFlag), "--skip-keys="+strconv.FormatBool(args.SkipKeys),
		"--skip-subscriptions="+strconv.FormatBool(args.SkipSubscriptions), "-e", args.DestAPIM.EnvName, "-k", "--verbose")
	if doClean {
		t.Cleanup(func() {
			args.DestAPIM.DeleteApplicationByName(args.Application.Name)
		})
	}

	return output, err
}

func importAppPreserveOwnerAndUpdate(t *testing.T, sourceEnv string, app *apim.Application, client *apim.Client) (string, error) {
	fileName := base.GetApplicationArchiveFilePath(t, sourceEnv, app.Name, app.Owner)
	output, err := base.Execute(t, "import", "app", "--preserve-owner=true", "--update=true", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	return output, err
}

func ValidateAppExportFailure(t *testing.T, args *AppImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting app from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportApp(t, args)

	// Validate that export failed
	assert.False(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
		args.Application.Name, args.AppOwner.Username))
}

func ValidateAppExport(t *testing.T, args *AppImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting app from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportApp(t, args)

	// Validate that export passed
	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
		args.Application.Name, args.AppOwner.Username))
}

func validateAppImport(t *testing.T, args *AppImportExportTestArgs, doClean bool) *apim.Application {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import app to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importApp(t, args, doClean)

	// Get App from env 2
	importedApp := GetApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	ValidateAppsEqual(t, args, importedApp)

	return importedApp
}

func ValidateAppExportImport(t *testing.T, args *AppImportExportTestArgs, doClean bool) *apim.Application {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export app from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportApp(t, args)

	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
		args.Application.Name, args.AppOwner.Username))

	// Import app to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importApp(t, args, doClean)

	// Get App from env 2
	importedApp := GetApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	ValidateAppsEqual(t, args, importedApp)

	return importedApp
}

func ValidateAppExportImportGeneratedKeys(t *testing.T, args *AppImportExportTestArgs, appId string, doClean bool) {

	// Generate keys for the application in env 1
	applicationKey := GenerateKeys(t, args.SrcAPIM, args.AppOwner.Username, args.AppOwner.Password, appId)

	// Export an application from env 1 and import it to env 2
	importedApplication := ValidateAppExportImport(t, args, doClean)
	// Retrieve oauth keys of the imported application to env2
	importedApplicationKey := GetOauthKeys(t, args.DestAPIM, args.AppOwner.Username, args.AppOwner.Password, importedApplication)

	if !args.SkipKeys {
		// Compare consumer key and secret of the application in env 1 and env 2
		assert.Equal(t, applicationKey.ConsumerKey, importedApplicationKey.ConsumerKey, "Consumer key mismatched")
		assert.Equal(t, applicationKey.ConsumerSecret, importedApplicationKey.ConsumerSecret, "Consumer secret mismatched")
	} else {
		// Assert whether the imported application's key is empty
		assert.Equal(t, importedApplicationKey, &apim.ApplicationKey{}, "Application keys are not empty")
	}
}

func ValidateAppExportImportSubscriptions(t *testing.T, args *AppImportExportTestArgs, appId string, importOnly bool,
	doClean bool) *apim.Application {
	// Get the subscriptions of the application to be exported from env 1
	subscriptionsListFromEnv1App := args.SrcAPIM.GetApplicationSubscriptions(appId)

	var importedApplication *apim.Application
	if !importOnly {
		// Export an application from env 1 and import it to env 2
		importedApplication = ValidateAppExportImport(t, args, doClean)
	} else {
		importedApplication = validateAppImport(t, args, doClean)
	}

	// Get the subscriptions of the imported application in env 1
	subscriptionsListFromEnv2App := args.DestAPIM.GetApplicationSubscriptions(importedApplication.ApplicationID)

	if !args.SkipSubscriptions {
		validateSubscriptionsOfApp(t, subscriptionsListFromEnv1App, subscriptionsListFromEnv2App)
	} else {
		assert.NotEqual(t, subscriptionsListFromEnv1App.Count, 0, "The subscriptions count of the app in env 1 is incorrect")
		assert.NotEqual(t, len(subscriptionsListFromEnv1App.List), "The subscriptions list of the app in env 1 is incorrect")
		assert.Equal(t, subscriptionsListFromEnv2App.Count, 0, "The subscriptions count of the imported app is incorrect")
		assert.Equal(t, len(subscriptionsListFromEnv2App.List), 0, "The subscriptions list of the imported app is incorrect")
	}
	return importedApplication
}

func ValidateAppsEqual(t *testing.T, args *AppImportExportTestArgs, app2 *apim.Application) {
	t.Helper()

	app1Copy := apim.CopyApp(args.Application)
	app2Copy := apim.CopyApp(app2)

	// Since the Applications are from too different envs, their respective ApplicationID will defer.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	same := "override_with_same_value"
	app1Copy.ApplicationID = same
	app2Copy.ApplicationID = same

	// When the application is imported with skipped subscriptions the subscription count and the scopes
	// will differ in the two applications. Hence those should be overridden with the same value.
	if args.SkipSubscriptions {
		sameInt := 0
		app1Copy.SubscriptionCount = sameInt
		app1Copy.SubscriptionScopes = []string{same}
		app2Copy.SubscriptionCount = sameInt
		app2Copy.SubscriptionScopes = []string{same}
	}

	assert.Equal(t, app1Copy, app2Copy, "Application obejcts are not equal")

}

func validateSubscriptionsOfApp(t *testing.T, subscriptionsOfApp1 *apim.SubscriptionList,
	subscriptionsOfApp2 *apim.SubscriptionList) {
	t.Helper()
	subscriptionsOfApp1 = OverrideDifferedPropertiesOfSubscriptions(subscriptionsOfApp1)
	subscriptionsOfApp2 = OverrideDifferedPropertiesOfSubscriptions(subscriptionsOfApp2)

	assert.Equal(t, subscriptionsOfApp1, subscriptionsOfApp2, "Subscriptions objects are not equal")
}

func DeleteAppByCtl(t *testing.T, args *AppImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "app", "-n", args.Application.Name, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func ValidateApplicationIsDeleted(t *testing.T, application *apim.Application, appsListAfterDelete *apim.ApplicationList) {
	for _, existingApplication := range appsListAfterDelete.List {
		assert.NotEqual(t, existingApplication.ApplicationID, application.ApplicationID, "API delete is not successful")
	}
}

func ValidateAppDelete(t *testing.T, args *AppImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())

	// Delete an App of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()
	appsListBeforeDelete := args.SrcAPIM.GetApplications()

	DeleteAppByCtl(t, args)

	appsListAfterDelete := args.SrcAPIM.GetApplications()
	base.WaitForIndexing()

	// Validate whether the expected number of App count is there
	assert.Equal(t, appsListBeforeDelete.Count, appsListAfterDelete.Count+1, "Expected number of Applications not deleted")

	// Validate that the delete is a success
	ValidateApplicationIsDeleted(t, args.Application, appsListAfterDelete)
}

func ValidateListAppsWithOwner(t *testing.T, envName string) {
	//Clean up existing default apictl app
	base.Execute(t, "delete", "app", "-n", "default-apictl-app", "-e", envName, "-k", "--verbose")
	response := ListAppsWithOwner(t, envName, "admin")
	assert.Equal(t, 5, len(response), "Failed when listing Applications with owner as Admin")

	emptyResponse := ListAppsWithOwner(t, envName, "user1")
	assert.Equal(t, 0, len(emptyResponse), "Failed when listing Applications with owner as User1")
}

func OverrideDifferedPropertiesOfSubscriptions(subscriptionsList1 *apim.SubscriptionList) *apim.SubscriptionList {
	// Since the Applications are from too different envs, their respective Subscription IDs, API IDs and Application IDs will defer.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	same := "override_with_same_value"
	subscriptionsList1Copy := apim.SubscriptionList{
		Count: subscriptionsList1.Count,
		List:  []apim.Subscription{}}
	for _, subscription := range subscriptionsList1.List {
		subscription.SubscriptionID = same
		subscription.APIID = same
		subscription.ApplicationID = same
		subscriptionsList1Copy.List = append(subscriptionsList1Copy.List, subscription)
	}
	return &subscriptionsList1Copy
}
