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

func GetApp(t *testing.T, client *apim.Client, name string, username string, password string) *apim.Application {
	client.Login(username, password)
	appInfo := client.GetApplicationByName(name)
	return client.GetApplication(appInfo.ApplicationID)
}

func ListApps(t *testing.T, env string) []string {
	response, _ := base.Execute(t, "get", "apps", "-e", env, "-k")

	return base.GetRowsFromTableResponse(response)
}

func ListAppsWithOwner(t *testing.T, env string, owner string) []string {
	response, _ := base.Execute(t, "gets", "apps", "-e", env, "-k", "--owner", owner)

	return base.GetRowsFromTableResponse(response)
}

func getEnvAppExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedAppsDirName, envName)
}

func exportApp(t *testing.T, name string, owner string, env string) (string, error) {
	output, err := base.Execute(t, "export", "app", "-n", name, "-o", owner, "-e", env, "-k", "--verbose")

	t.Cleanup(func() {
		base.RemoveApplicationArchive(t, getEnvAppExportPath(env), name, owner)
	})

	return output, err
}

func importAppPreserveOwner(t *testing.T, sourceEnv string, app *apim.Application, client *apim.Client) (string, error) {
	fileName := base.GetApplicationArchiveFilePath(t, sourceEnv, app.Name, app.Owner)
	output, err := base.Execute(t, "import", "app", "--preserveOwner=true", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	t.Cleanup(func() {
		client.DeleteApplicationByName(app.Name)
	})

	return output, err
}

func importAppPreserveOwnerAndUpdate(t *testing.T, sourceEnv string, app *apim.Application, client *apim.Client) (string, error) {
	fileName := base.GetApplicationArchiveFilePath(t, sourceEnv, app.Name, app.Owner)
	output, err := base.Execute(t, "import", "app", "--preserveOwner=true", "--update=true", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	return output, err
}

func ValidateAppExportFailure(t *testing.T, args *AppImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting app from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportApp(t, args.Application.Name, args.AppOwner.Username, args.SrcAPIM.GetEnvName())

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

	exportApp(t, args.Application.Name, args.AppOwner.Username, args.SrcAPIM.GetEnvName())

	// Validate that export passed
	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
		args.Application.Name, args.AppOwner.Username))
}

func ValidateAppExportImportWithPreserveOwner(t *testing.T, args *AppImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export app from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportApp(t, args.Application.Name, args.AppOwner.Username, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
		args.Application.Name, args.AppOwner.Username))

	// Import app to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAppPreserveOwner(t, args.SrcAPIM.GetEnvName(), args.Application, args.DestAPIM)

	// Get App from env 2
	importedApp := GetApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	ValidateAppsEqual(t, args.Application, importedApp)
}

func ValidateAppExportImportWithUpdate(t *testing.T, args *AppImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export app from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportApp(t, args.Application.Name, args.AppOwner.Username, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
		args.Application.Name, args.AppOwner.Username))

	// Import app to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAppPreserveOwnerAndUpdate(t, args.SrcAPIM.GetEnvName(), args.Application, args.DestAPIM)

	// Get App from env 2
	importedApp := GetApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	ValidateAppsEqual(t, args.Application, importedApp)
}

func ValidateAppsEqual(t *testing.T, app1 *apim.Application, app2 *apim.Application) {
	t.Helper()

	app1Copy := apim.CopyApp(app1)
	app2Copy := apim.CopyApp(app2)

	// Since the Applications are from too different envs, their respective ApplicationID will defer.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	app1Copy.ApplicationID = "override_with_same_value"
	app2Copy.ApplicationID = app1Copy.ApplicationID

	assert.Equal(t, app1Copy, app2Copy, "Application obejcts are not equal")

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
