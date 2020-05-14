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
	"crypto/tls"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func getKeys(t *testing.T, provider string, name string, version string, env string) (string, error) {
	return base.Execute(t, "get-keys", "-n", name, "-v", version, "-r", provider, "-e", env, "-k")
}

func invokeAPI(t *testing.T, url string, key string, expectedCode int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)

	assert.Nil(t, err, "Error while generating GET")

	authHeader := "Bearer " + key
	req.Header.Set("Authorization", authHeader)

	t.Log("invokeAPI() url", url)

	response, err := client.Do(req)

	assert.Nil(t, err, "Error while invoking API")
	assert.Equal(t, expectedCode, response.StatusCode, "API Invocation failed")
}

func addAPI(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username)
	id := client.AddAPI(t, api, username, password)
	api = client.GetAPI(id)
	return api
}

func getAPI(t *testing.T, client *apim.Client, name string, username string, password string) *apim.API {
	client.Login(username, password)
	apiInfo := client.GetAPIByName(name)
	return client.GetAPI(apiInfo.ID)
}

func publishAPI(client *apim.Client, username string, password string, apiID string) {
	time.Sleep(2000 * time.Millisecond)
	client.Login(username, password)
	client.PublishAPI(apiID)
}

func unsubscribeAPI(client *apim.Client, username string, password string, apiID string) {
	client.Login(username, password)
	client.DeleteSubscriptions(apiID)
}

func getResourceURL(apim *apim.Client, api *apim.API) string {
	port := 8280 + apim.GetPortOffset()
	return "http://" + apim.GetHost() + ":" + strconv.Itoa(port) + api.Context + "/" + api.Version + "/menu"
}

func validateGetKeysFailure(t *testing.T, args *apiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.apim.GetEnvName(), args.apim.GetApimURL(), args.apim.GetTokenURL())
	base.Login(t, args.apim.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	result, err := getKeys(t, args.api.Provider, args.api.Name, args.api.Version, args.apim.GetEnvName())

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Equal(t, "Exit status 1", base.GetValueOfUniformResponse(result))
}

func validateGetKeys(t *testing.T, args *apiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.apim.GetEnvName(), args.apim.GetApimURL(), args.apim.GetTokenURL())
	base.Login(t, args.apim.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	result, err := getKeys(t, args.api.Provider, args.api.Name, args.api.Version, args.apim.GetEnvName())

	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while getting key")

	invokeAPI(t, getResourceURL(args.apim, args.api), base.GetValueOfUniformResponse(result), 200)
	unsubscribeAPI(args.apim, args.ctlUser.username, args.ctlUser.password, args.api.ID)
}

func addApp(t *testing.T, client *apim.Client, username string, password string) *apim.Application {
	client.Login(username, password)
	app := client.GenerateSampleAppData()
	return client.AddApplication(t, app, username, password)
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

func getEnvAPIExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedApisDirName, envName)
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

func exportAPI(t *testing.T, name string, version string, provider string, env string) (string, error) {
	output, err := base.Execute(t, "export-api", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, getEnvAPIExportPath(env), name, version)
	})

	return output, err
}

func importAPIPreserveProvider(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import-api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	t.Cleanup(func() {
		client.DeleteAPIByName(api.Name)
	})

	return output, err
}

func validateAppExportFailure(t *testing.T, args *appImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Attempt exporting app from env
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportApp(t, args.application.Name, args.appOwner.username, args.srcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.srcAPIM.GetEnvName()),
		args.application.Name, args.appOwner.username))
}

func validateAppExportImport(t *testing.T, args *appImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Export app from env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportApp(t, args.application.Name, args.appOwner.username, args.srcAPIM.GetEnvName())

	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.srcAPIM.GetEnvName()),
		args.application.Name, args.appOwner.username))

	// Import app to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	importAppPreserveOwner(t, args.srcAPIM.GetEnvName(), args.application, args.destAPIM)

	// Get App from env 2
	importedApp := getApp(t, args.destAPIM, args.application.Name, args.appOwner.username, args.appOwner.password)

	// Validate env 1 and env 2 App is equal
	validateAppsEqual(t, args.application, importedApp)
}

func validateAppsEqual(t *testing.T, app1 *apim.Application, app2 *apim.Application) {
	t.Helper()

	app1Copy := apim.CopyApp(app1)
	app2Copy := apim.CopyApp(app2)

	// Since the Applications are from too different envs, their respective ApplicationID will defer.
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	app1Copy.ApplicationID = "override_with_same_value"
	app2Copy.ApplicationID = app1Copy.ApplicationID

	assert.Equal(t, app1Copy, app2Copy, "Application obejcts are not equal")

}

func validateAPIExportFailure(t *testing.T, args *apiImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Attempt exporting api from env
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportAPI(t, args.api.Name, args.api.Version, args.api.Provider, args.srcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsAPIArchiveExists(t, getEnvAPIExportPath(args.srcAPIM.GetEnvName()),
		args.api.Name, args.api.Version))
}

func validateAPIExportImport(t *testing.T, args *apiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportAPI(t, args.api.Name, args.api.Version, args.api.Provider, args.srcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIExportPath(args.srcAPIM.GetEnvName()),
		args.api.Name, args.api.Version))

	// Import api to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	importAPIPreserveProvider(t, args.srcAPIM.GetEnvName(), args.api, args.destAPIM)

	// Give time for newly imported API to get indexed, or else getAPI by name will fail
	time.Sleep(1 * time.Second)

	// Get App from env 2
	importedAPI := getAPI(t, args.destAPIM, args.api.Name, args.apiProvider.username, args.apiProvider.password)

	// Validate env 1 and env 2 API is equal
	validateAPIsEqual(t, args.api, importedAPI)
}

func validateAPIsEqual(t *testing.T, api1 *apim.API, api2 *apim.API) {
	t.Helper()

	api1Copy := apim.CopyAPI(api1)
	api2Copy := apim.CopyAPI(api2)

	same := "override_with_same_value"
	// Since the APIs are from too different envs, their respective ID will defer.
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	api1Copy.ID = same
	api2Copy.ID = same

	api1Copy.CreatedTime = same
	api2Copy.CreatedTime = same

	api1Copy.LastUpdatedTime = same
	api2Copy.LastUpdatedTime = same

	// Sort member collections to make equality chack possible
	apim.SortAPIMembers(&api1Copy)
	apim.SortAPIMembers(&api2Copy)

	assert.Equal(t, api1Copy, api2Copy, "API obejcts are not equal")

}
