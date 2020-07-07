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
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"path/filepath"
	"testing"
)

func addApp(t *testing.T, client *apim.Client, username string, password string) *apim.Application {
	client.Login(username, password)
	app := client.GenerateSampleAppData()
	doClean := true
	return client.AddApplication(t, app, username, password, doClean)
}

func addApplicationWithoutCleaning(t *testing.T, client *apim.Client, username string, password string) *apim.Application {
	client.Login(username, password)
	application := client.GenerateSampleAppData()
	doClean := false
	app := client.AddApplication(t, application, username, password, doClean)
	application = client.GetApplication(app.ApplicationID)
	return application
}

func getApp(t *testing.T, client *apim.Client, name string, username string, password string) *apim.Application {
	client.Login(username, password)
	appInfo := client.GetApplicationByName(name)
	return client.GetApplication(appInfo.ApplicationID)
}

func listApps(t *testing.T, env string) []string {
	response, _ := base.Execute(t, "list", "apps", "-e", env, "-k")

	return base.GetRowsFromTableResponse(response)
}

func getEnvAppExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedAppsDirName, envName)
}

func exportApp(t *testing.T, name string, owner string, env string) (string, error) {
	output, err := base.Execute(t, "export-app", "-n", name, "-o", owner, "-e", env, "-k", "--verbose")

	t.Cleanup(func() {
		base.RemoveApplicationArchive(t, getEnvAppExportPath(env), name, owner)
	})

	return output, err
}

func importAppPreserveOwner(t *testing.T, sourceEnv string, app *apim.Application, client *apim.Client) (string, error) {
	fileName := base.GetApplicationArchiveFilePath(t, sourceEnv, app.Name, app.Owner)
	output, err := base.Execute(t, "import-app", "--preserveOwner=true", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	t.Cleanup(func() {
		client.DeleteApplicationByName(app.Name)
	})

	return output, err
}

func importAppPreserveOwnerAndUpdate(t *testing.T, sourceEnv string, app *apim.Application, client *apim.Client) (string, error) {
	fileName := base.GetApplicationArchiveFilePath(t, sourceEnv, app.Name, app.Owner)
	output, err := base.Execute(t, "import-app", "--preserveOwner=true", "--update=true", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	return output, err
}

func validateAppExportFailure(t *testing.T, args *appImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Attempt exporting app from env
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	exportApp(t, args.application.Name, args.appOwner.Username, args.srcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.srcAPIM.GetEnvName()),
		args.application.Name, args.appOwner.Username))
}

func validateAppExportImportWithPreserveOwner(t *testing.T, args *appImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Export app from env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	exportApp(t, args.application.Name, args.appOwner.Username, args.srcAPIM.GetEnvName())

	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.srcAPIM.GetEnvName()),
		args.application.Name, args.appOwner.Username))

	// Import app to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	importAppPreserveOwner(t, args.srcAPIM.GetEnvName(), args.application, args.destAPIM)

	// Get App from env 2
	importedApp := getApp(t, args.destAPIM, args.application.Name, args.appOwner.Username, args.appOwner.Password)

	// Validate env 1 and env 2 App is equal
	validateAppsEqual(t, args.application, importedApp)
}

func validateAppExportImportWithUpdate(t *testing.T, args *appImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Export app from env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	exportApp(t, args.application.Name, args.appOwner.Username, args.srcAPIM.GetEnvName())

	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.srcAPIM.GetEnvName()),
		args.application.Name, args.appOwner.Username))

	// Import app to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	importAppPreserveOwnerAndUpdate(t, args.srcAPIM.GetEnvName(), args.application, args.destAPIM)

	// Get App from env 2
	importedApp := getApp(t, args.destAPIM, args.application.Name, args.appOwner.Username, args.appOwner.Password)

	// Validate env 1 and env 2 App is equal
	validateAppsEqual(t, args.application, importedApp)
}

func validateAppsEqual(t *testing.T, app1 *apim.Application, app2 *apim.Application) {
	t.Helper()

	app1Copy := apim.CopyApp(app1)
	app2Copy := apim.CopyApp(app2)

	// Since the Applications are from too different envs, their respective ApplicationID will defer.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	app1Copy.ApplicationID = "override_with_same_value"
	app2Copy.ApplicationID = app1Copy.ApplicationID

	assert.Equal(t, app1Copy, app2Copy, "Application obejcts are not equal")

}

func deleteAppByCtl(t *testing.T, args *appImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "app", "-n", args.application.Name, "-e", args.srcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func validateApplicationIsDeleted(t *testing.T, application *apim.Application, appsListAfterDelete *apim.ApplicationList) {
	for _, existingApplication := range appsListAfterDelete.List {
		assert.NotEqual(t, existingApplication.ApplicationID, application.ApplicationID, "API delete is not successful")
	}
}
