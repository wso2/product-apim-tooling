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
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

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

func invokeAPIProduct(t *testing.T, url string, key string, expectedCode int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)

	assert.Nil(t, err, "Error while generating GET")

	authHeader := "Bearer " + key
	req.Header.Set("Authorization", authHeader)

	t.Log("invokeAPIProduct() url", url)

	response, err := client.Do(req)

	assert.Nil(t, err, "Error while invoking API Product")
	assert.Equal(t, expectedCode, response.StatusCode, "API Product Invocation failed")
}

func addAPI(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username)
	doClean := true
	id := client.AddAPI(t, api, username, password, doClean)
	api = client.GetAPI(id)
	return api
}

func addAPIWithoutCleaning(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username)
	doClean := false
	id := client.AddAPI(t, api, username, password, doClean)
	api = client.GetAPI(id)
	return api
}

func addAPIToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.API, *apim.API) {
	client1.Login(username, password)
	api := client1.GenerateSampleAPIData(username)
	doClean := true
	id1 := client1.AddAPI(t, api, username, password, doClean)
	api1 := client1.GetAPI(id1)

	client2.Login(username, password)
	id2 := client2.AddAPI(t, api, username, password, doClean)
	api2 := client2.GetAPI(id2)

	return api1, api2
}

func addAPIFromOpenAPIDefinition(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	path := "testdata/petstore.yaml"
	client.Login(username, password)
	additionalProperties := client.GenerateAdditionalProperties(username)
	id := client.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api := client.GetAPI(id)
	return api
}

func addAPIFromOpenAPIDefinitionToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.API, *apim.API) {
	path := "testdata/petstore.yaml"
	client1.Login(username, password)
	additionalProperties := client1.GenerateAdditionalProperties(username)
	id1 := client1.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api1 := client1.GetAPI(id1)

	client2.Login(username, password)
	id2 := client2.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api2 := client2.GetAPI(id2)

	return api1, api2
}

func addAPIProductFromJSON(t *testing.T, client *apim.Client, username string, password string, apisList map[string]*apim.API) *apim.APIProduct {
	client.Login(username, password)
	path := "testdata/SampleAPIProduct.json"
	doClean := true
	id := client.AddAPIProductFromJSON(t, path, username, password, apisList, doClean)
	apiProduct := client.GetAPIProduct(id)
	return apiProduct
}

func addAPIProductFromJSONWithoutCleaning(t *testing.T, client *apim.Client, username string, password string, apisList map[string]*apim.API) *apim.APIProduct {
	client.Login(username, password)
	path := "testdata/SampleAPIProduct.json"
	doClean := false
	id := client.AddAPIProductFromJSON(t, path, username, password, apisList, doClean)
	apiProduct := client.GetAPIProduct(id)
	return apiProduct
}

func getAPI(t *testing.T, client *apim.Client, name string, username string, password string) *apim.API {
	client.Login(username, password)
	apiInfo := client.GetAPIByName(name)
	return client.GetAPI(apiInfo.ID)
}

func getAPIs(client *apim.Client, username string, password string) *apim.APIList {
	client.Login(username, password)
	return client.GetAPIs()
}

func deleteAPI(t *testing.T, client *apim.Client, apiID string, username string, password string) {
	time.Sleep(2000 * time.Millisecond)
	client.Login(username, password)
	client.DeleteAPI(apiID)
}

func deleteAPIByCtl(t *testing.T, args *apiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "api", "-n", args.api.Name, "-v", args.api.Version, "-e", args.srcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func getAPIProduct(t *testing.T, client *apim.Client, name string, username string, password string) *apim.APIProduct {
	client.Login(username, password)
	apiProductInfo := client.GetAPIProductByName(name)
	return client.GetAPIProduct(apiProductInfo.ID)
}

func getAPIProducts(client *apim.Client, username string, password string) *apim.APIProductList {
	client.Login(username, password)
	return client.GetAPIProducts()
}

func deleteAPIProduct(t *testing.T, client *apim.Client, apiProductID string, username string, password string) {
	time.Sleep(2000 * time.Millisecond)
	client.Login(username, password)
	client.DeleteAPIProduct(apiProductID)
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

func getResourceURLForAPIProduct(apim *apim.Client, apiProduct *apim.APIProduct) string {
	port := 8280 + apim.GetPortOffset()
	return "http://" + apim.GetHost() + ":" + strconv.Itoa(port) + apiProduct.Context + "/menu"
}

func validateGetKeysFailure(t *testing.T, args *apiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.apim.GetEnvName(), args.apim.GetApimURL(), args.apim.GetTokenURL())
	base.Login(t, args.apim.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	var err error
	var result string
	if args.api != nil {
		result, err = getKeys(t, args.api.Provider, args.api.Name, args.api.Version, args.apim.GetEnvName())
	}

	if args.apiProduct != nil {
		result, err = getKeys(t, args.apiProduct.Provider, args.apiProduct.Name, utils.DefaultApiProductVersion, args.apim.GetEnvName())
	}

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
}

func validateGetKeys(t *testing.T, args *apiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.apim.GetEnvName(), args.apim.GetApimURL(), args.apim.GetTokenURL())
	base.Login(t, args.apim.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	var err error
	var result string
	if args.api != nil {
		result, err = getKeys(t, args.api.Provider, args.api.Name, args.api.Version, args.apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		invokeAPI(t, getResourceURL(args.apim, args.api), base.GetValueOfUniformResponse(result), 200)
		unsubscribeAPI(args.apim, args.ctlUser.username, args.ctlUser.password, args.api.ID)
	}

	if args.apiProduct != nil {
		result, err = getKeys(t, args.apiProduct.Provider, args.apiProduct.Name, utils.DefaultApiProductVersion, args.apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		invokeAPIProduct(t, getResourceURLForAPIProduct(args.apim, args.apiProduct), base.GetValueOfUniformResponse(result), 200)
		unsubscribeAPI(args.apim, args.ctlUser.username, args.ctlUser.password, args.apiProduct.ID)
	}
}

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

func getEnvAPIExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedApisDirName, envName)
}

func getEnvAPIProductExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedApiProductsDirName, envName)
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

func exportAPI(t *testing.T, name string, version string, provider string, env string) (string, error) {
	var output string
	var err error

	if provider == "" {
		output, err = base.Execute(t, "export-api", "-n", name, "-v", version, "-e", env, "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "export-api", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")
	}

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, getEnvAPIExportPath(env), name, version)
	})

	return output, err
}

func exportAPIProduct(t *testing.T, name string, version string, env string) (string, error) {
	output, err := base.Execute(t, "export", "api-product", "-n", name, "-e", env, "-k", "--verbose")

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, getEnvAPIProductExportPath(env), name, version)
	})

	return output, err
}

func importAPI(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import-api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose", "--preserve-provider=false")

	t.Cleanup(func() {
		client.DeleteAPIByName(api.Name)
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

func importAPIPreserveProviderFailure(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import-api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")
	return output, err
}

func listAPIs(t *testing.T, args *apiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "list", "apis", "-e", args.srcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func importAPIProductPreserveProvider(t *testing.T, args *apiProductImportExportTestArgs) (string, error) {
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.srcAPIM.GetEnvName(), args.apiProduct.Name, utils.DefaultApiProductVersion)

	if args.importApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--import-apis")
	} else if args.updateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--update-apis")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose")
	}

	t.Cleanup(func() {
		args.destAPIM.DeleteAPIProductByName(args.apiProduct.Name)
	})

	return output, err
}

func importUpdateAPIProductPreserveProvider(t *testing.T, args *apiProductImportExportTestArgs) (string, error) {
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.srcAPIM.GetEnvName(), args.apiProduct.Name, utils.DefaultApiProductVersion)

	if args.updateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--update-apis")
	} else if args.updateApiProductFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--update-api-product")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose")
	}

	return output, err
}

func importAPIProduct(t *testing.T, args *apiProductImportExportTestArgs) (string, error) {
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.srcAPIM.GetEnvName(), args.apiProduct.Name, utils.DefaultApiProductVersion)

	if args.importApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--import-apis", "--preserve-provider=false")
	} else if args.updateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--update-apis", "--preserve-provider=false")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--preserve-provider=false")
	}

	t.Cleanup(func() {
		args.destAPIM.DeleteAPIProductByName(args.apiProduct.Name)
	})

	return output, err
}

func importUpdateAPIProduct(t *testing.T, args *apiProductImportExportTestArgs) (string, error) {
	var output string
	var err error

	fileName := base.GetAPIArchiveFilePath(t, args.srcAPIM.GetEnvName(), args.apiProduct.Name, utils.DefaultApiProductVersion)

	if args.updateApisFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--update-apis", "--preserve-provider=false")
	} else if args.updateApiProductFlag {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--update-api-product", "--preserve-provider=false")
	} else {
		output, err = base.Execute(t, "import", "api-product", "-f", fileName, "-e", args.destAPIM.EnvName, "-k", "--verbose", "--preserve-provider=false")
	}

	return output, err
}

func listAPIProducts(t *testing.T, args *apiProductImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "list", "api-products", "-e", args.srcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func deleteAPIProductByCtl(t *testing.T, args *apiProductImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "api-product", "-n", args.apiProduct.Name, "-e", args.srcAPIM.EnvName, "-k", "--verbose")
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

func validateAppExportImportWithPreserveOwner(t *testing.T, args *appImportExportTestArgs) {
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

func validateAppExportImportWithUpdate(t *testing.T, args *appImportExportTestArgs) {
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

	importAppPreserveOwnerAndUpdate(t, args.srcAPIM.GetEnvName(), args.application, args.destAPIM)

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
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
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

	exportAPI(t, args.api.Name, args.api.Version, args.apiProvider.username, args.srcAPIM.GetEnvName())

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

func validateAPIExport(t *testing.T, args *apiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportAPI(t, args.api.Name, args.api.Version, args.apiProvider.username, args.srcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIExportPath(args.srcAPIM.GetEnvName()),
		args.api.Name, args.api.Version))
}

func validateAPIImport(t *testing.T, args *apiImportExportTestArgs) {
	t.Helper()

	// Import api to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	importAPI(t, args.srcAPIM.GetEnvName(), args.api, args.destAPIM)

	// Give time for newly imported API to get indexed, or else getAPI by name will fail
	time.Sleep(1 * time.Second)

	// Get App from env 2
	importedAPI := getAPI(t, args.destAPIM, args.api.Name, args.apiProvider.username, args.apiProvider.password)

	// Validate env 1 and env 2 API is equal
	validateAPIsEqualCrossTenant(t, args.api, importedAPI)
}

func validateAPIImportFailure(t *testing.T, args *apiImportExportTestArgs) {
	t.Helper()

	// Import api to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	// importAPIPreserveProviderFailure is used to eleminate cleaning the API after importing
	result, err := importAPIPreserveProviderFailure(t, args.srcAPIM.GetEnvName(), args.api, args.destAPIM)

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
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

func validateAPIsList(t *testing.T, args *apiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// List APIs of env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	output, _ := listAPIs(t, args)

	apisList := args.srcAPIM.GetAPIs()

	validateListAPIsEqual(t, output, apisList)
}

func validateListAPIsEqual(t *testing.T, apisListFromCtl string, apisList *apim.APIList) {

	for _, api := range apisList.List {
		// If the output string contains the same API ID, then decrement the count
		if strings.Contains(apisListFromCtl, api.ID) {
			apisList.Count = apisList.Count - 1
		}
	}

	// Count == 0 means that all the APIs from apisList were in apisListFromCtl
	assert.Equal(t, apisList.Count, 0, "API lists are not equal")
}

func validateAPIsEqualCrossTenant(t *testing.T, api1 *apim.API, api2 *apim.API) {
	t.Helper()

	api1Copy := apim.CopyAPI(api1)
	api2Copy := apim.CopyAPI(api2)

	same := "override_with_same_value"
	// Since the APIs are from too different envs, their respective ID will defer.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	api1Copy.ID = same
	api2Copy.ID = same

	api1Copy.CreatedTime = same
	api2Copy.CreatedTime = same

	api1Copy.LastUpdatedTime = same
	api2Copy.LastUpdatedTime = same

	// The contexts and providers will differ since this is a cross tenant import
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	api1Copy.Context = same
	api2Copy.Context = same

	api1Copy.Provider = same
	api2Copy.Provider = same

	// Sort member collections to make equality check possible
	apim.SortAPIMembers(&api1Copy)
	apim.SortAPIMembers(&api2Copy)

	assert.Equal(t, api1Copy, api2Copy, "API obejcts are not equal")
}

func validateAPIDelete(t *testing.T, args *apiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Delete an API of env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	time.Sleep(1 * time.Second)
	apisListBeforeDelete := args.srcAPIM.GetAPIs()

	deleteAPIByCtl(t, args)

	apisListAfterDelete := args.srcAPIM.GetAPIs()
	time.Sleep(1 * time.Second)

	// Validate whether the expected number of API count is there
	assert.Equal(t, apisListBeforeDelete.Count, apisListAfterDelete.Count+1, "Expected number of APIs not deleted")

	// Validate that the delete is a success
	validateAPIIsDeleted(t, args.api, apisListAfterDelete)
}

func validateAPIIsDeleted(t *testing.T, api *apim.API, apisListAfterDelete *apim.APIList) {
	for _, existingAPI := range apisListAfterDelete.List {
		assert.NotEqual(t, existingAPI.ID, api.ID, "API delete is not successful")
	}
}

func validateAPIProductExportFailure(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Attempt exporting API Product from env
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportAPIProduct(t, args.apiProduct.Name, utils.DefaultApiProductVersion, args.srcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsAPIArchiveExists(t, getEnvAPIProductExportPath(args.srcAPIM.GetEnvName()),
		args.apiProduct.Name, utils.DefaultApiProductVersion))
}

func validateAPIProductExportImportPreserveProvider(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Export API Product from env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportAPIProduct(t, args.apiProduct.Name, utils.DefaultApiProductVersion, args.srcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIProductExportPath(args.srcAPIM.GetEnvName()),
		args.apiProduct.Name, utils.DefaultApiProductVersion))

	// Import API Product to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	importAPIProductPreserveProvider(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	time.Sleep(1 * time.Second)

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.destAPIM, args.apiProduct.Name, args.apiProductProvider.username, args.apiProductProvider.password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqual(t, args.apiProduct, importedAPIProduct)
}

func validateAPIProductImportUpdatePreserveProvider(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Import api to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	// This is used when you have previously imported an API Product (with preserving the provider) and validated it.
	// So when doing the cleaning you do not need to clean twice. For that, importUpdateAPIProductPreserveProvider will not be doing cleaning again.
	importUpdateAPIProductPreserveProvider(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	time.Sleep(1 * time.Second)

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.destAPIM, args.apiProduct.Name, args.apiProductProvider.username, args.apiProductProvider.password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqual(t, args.apiProduct, importedAPIProduct)
}

func validateAPIProductExport(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Export API Product from env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	exportAPIProduct(t, args.apiProduct.Name, utils.DefaultApiProductVersion, args.srcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIProductExportPath(args.srcAPIM.GetEnvName()),
		args.apiProduct.Name, utils.DefaultApiProductVersion))
}

func validateAPIProductImport(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Import API Product to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	importAPIProduct(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	time.Sleep(1 * time.Second)

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.destAPIM, args.apiProduct.Name, args.apiProductProvider.username, args.apiProductProvider.password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqualCrossTenant(t, args.apiProduct, importedAPIProduct)
}

func validateAPIProductImportUpdate(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.destAPIM.GetEnvName(), args.destAPIM.GetApimURL(), args.destAPIM.GetTokenURL())

	// Import API Product to env 2
	base.Login(t, args.destAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	// This is used when you have previously imported an API Product and validated it.
	// So when doing the cleaning you do not need to clean twice. For that, importUpdateAPIProduct will not be doing cleaning again.
	importUpdateAPIProduct(t, args)

	// Give time for newly imported API Product to get indexed, or else getAPIProduct by name will fail
	time.Sleep(1 * time.Second)

	// Get API Product from env 2
	importedAPIProduct := getAPIProduct(t, args.destAPIM, args.apiProduct.Name, args.apiProductProvider.username, args.apiProductProvider.password)

	// Validate env 1 and env 2 API Products are equal
	validateAPIProductsEqualCrossTenant(t, args.apiProduct, importedAPIProduct)
}

func validateAPIProductsEqual(t *testing.T, apiProduct1 *apim.APIProduct, apiProduct2 *apim.APIProduct) {
	t.Helper()

	apiProduct1Copy := apim.CopyAPIProduct(apiProduct1)
	apiProduct2Copy := apim.CopyAPIProduct(apiProduct2)

	same := "override_with_same_value"
	// Since the API Products are from too different envs, their respective ID will defer.
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	apiProduct1Copy.ID = same
	apiProduct2Copy.ID = same

	apiProduct1Copy.CreatedTime = same
	apiProduct2Copy.CreatedTime = same

	apiProduct1Copy.LastUpdatedTime = same
	apiProduct2Copy.LastUpdatedTime = same

	// Check the validity of the operations in the APIs array of the two API Products
	err := validateOperations(&apiProduct1Copy, &apiProduct2Copy)
	if err != nil {
		t.Fatal(err)
	}

	// Sort member collections to make equality check possible
	apim.SortAPIProductMembers(&apiProduct1Copy)
	apim.SortAPIProductMembers(&apiProduct2Copy)

	assert.Equal(t, apiProduct1Copy, apiProduct2Copy, "API Product obejcts are not equal")
}

func validateAPIProductsEqualCrossTenant(t *testing.T, apiProduct1 *apim.APIProduct, apiProduct2 *apim.APIProduct) {
	t.Helper()

	apiProduct1Copy := apim.CopyAPIProduct(apiProduct1)
	apiProduct2Copy := apim.CopyAPIProduct(apiProduct2)

	same := "override_with_same_value"
	// Since the API Products are from too different envs, their respective ID will defer.
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	apiProduct1Copy.ID = same
	apiProduct2Copy.ID = same

	apiProduct1Copy.CreatedTime = same
	apiProduct2Copy.CreatedTime = same

	apiProduct1Copy.LastUpdatedTime = same
	apiProduct2Copy.LastUpdatedTime = same

	// The contexts and providers will differ since this is a cross tenant import
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	apiProduct1Copy.Context = same
	apiProduct2Copy.Context = same

	apiProduct1Copy.Provider = same
	apiProduct2Copy.Provider = same

	// Check the validity of the operations in the APIs array of the two API Products
	err := validateOperations(&apiProduct1Copy, &apiProduct2Copy)
	if err != nil {
		t.Fatal(err)
	}

	// Sort member collections to make equality check possible
	apim.SortAPIProductMembers(&apiProduct1Copy)
	apim.SortAPIProductMembers(&apiProduct2Copy)

	assert.Equal(t, apiProduct1Copy, apiProduct2Copy, "API Product obejcts are not equal")
}

func validateOperations(apiProduct1Copy, apiProduct2Copy *apim.APIProduct) error {

	// To store the validity of each dependent API operation
	var isOperationsValid []bool

	// Iterate thorugh the APIs array of API Product 1
	for index, apiInProduct1 := range apiProduct1Copy.APIs {
		// Iterate thorugh the APIs array of API Product 2
		for _, apiInProduct2 := range apiProduct2Copy.APIs {
			// If the name of the APIs in the two API Products are same, those should be compared
			if apiInProduct1.(map[string]interface{})["name"] == apiInProduct2.(map[string]interface{})["name"] {

				// Convert the maps to APIOperations array structs (so that the structs can be compared easily)
				var operationsList []apim.APIOperations
				operationsInApiInProduct1, _ := json.Marshal(apiInProduct1.(map[string]interface{})["operations"])
				err := json.Unmarshal(operationsInApiInProduct1, &operationsList)
				if err != nil {
					return err
				}
				operationsInApiInProduct2, _ := json.Marshal(apiInProduct2.(map[string]interface{})["operations"])
				err = json.Unmarshal(operationsInApiInProduct2, &operationsList)
				if err != nil {
					return err
				}

				// Compare the two APIOperations array structs, whether they are equal
				if cmp.Equal(operationsInApiInProduct1, operationsInApiInProduct2) {
					// If the operations are equal, it is valid
					isOperationsValid = append(isOperationsValid, true)

					// Since the apiIds of the dependent APIs in the environments differ, those will be assigned integer values (index value in the loop)
					// Same APIs in both apiInProduct2 and apiInProduct2 will be assigned same integer values based on the order they are in APIs array
					apiInProduct2.(map[string]interface{})["apiId"] = index
					apiInProduct1.(map[string]interface{})["apiId"] = index
					break
				}
				// If the operations are not equal, it is not valid
				isOperationsValid = append(isOperationsValid, false)
			}
		}
	}

	// To store the overall result of whether the operations are valid.
	isAllOperationsValid := true
	for _, value := range isOperationsValid {
		// If any of the value in isOperationsValid array is false, the overall result should be false
		if value == false {
			isAllOperationsValid = false
		}
	}

	if isAllOperationsValid {
		// If all the operations in both of the APIs arrays are equal, make those two APIs arrays equals
		apiProduct2Copy.APIs = apiProduct1Copy.APIs
	}
	return nil
}

func validateAPIProductsList(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// List API Products of env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	output, _ := listAPIProducts(t, args)

	apiProductsList := args.srcAPIM.GetAPIProducts()

	validateListAPIProductsEqual(t, output, apiProductsList)
}

func validateListAPIProductsEqual(t *testing.T, apiProductsListFromCtl string, apiProductsList *apim.APIProductList) {

	for _, apiProduct := range apiProductsList.List {
		// If the output string contains the same API Product ID, then decrement the count
		if strings.Contains(apiProductsListFromCtl, apiProduct.ID) {
			apiProductsList.Count = apiProductsList.Count - 1
		}
	}

	// Count == 0 means that all the API Products from apiProductsList were in apiProductsListFromCtl
	assert.Equal(t, apiProductsList.Count, 0, "API Product lists are not equal")
}

func validateAPIProductDelete(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Delete an API Product of env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	time.Sleep(1 * time.Second)
	apiProductsListBeforeDelete := args.srcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	apiProductsListAfterDelete := args.srcAPIM.GetAPIProducts()
	time.Sleep(1 * time.Second)

	// Validate whether the expected number of API Product count is there
	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count+1, "Expected number of API Products not deleted")

	// Validate that the delete is a success
	validateAPIProductIsDeleted(t, args.apiProduct, apiProductsListAfterDelete)
}

func validateAPIProductIsDeleted(t *testing.T, apiProduct *apim.APIProduct, apiProductsListAfterDelete *apim.APIProductList) {
	for _, existingAPIProduct := range apiProductsListAfterDelete.List {
		assert.NotEqual(t, existingAPIProduct.ID, apiProduct.ID, "API Product delete is not successful")
	}
}

func validateAPIProductDeleteFailure(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Delete an API Product of env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	time.Sleep(1 * time.Second)
	apiProductsListBeforeDelete := args.srcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	time.Sleep(1 * time.Second)
	apiProductsListAfterDelete := args.srcAPIM.GetAPIProducts()

	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count, "API Product delete is successful")
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

func initProject(t *testing.T, args *initTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	output, err := base.Execute(t, "init", args.initFlag)
	return output, err
}

func initProjectWithDefinitionFlag(t *testing.T, args *initTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	output, err := base.Execute(t, "init", args.initFlag, "--definition", args.definitionFlag, "--force", strconv.FormatBool(args.forceFlag))
	return output, err
}

func initProjectWithOasFlag(t *testing.T, args *initTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	output, err := base.Execute(t, "init", args.initFlag, "--oas", args.oasFlag)
	return output, err
}

func environmentSetExportDirectory(t *testing.T, args *setTestArgs) (string, error) {
	apim := args.srcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--export-directory", args.exportDirectoryFlag, "-k")
	return output, error
}

func environmentSetHttpRequestTimeout(t *testing.T, args *setTestArgs) (string, error) {
	apim := args.srcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--http-request-timeout", strconv.Itoa(args.httpRequestTimeout), "-k")
	return output, error
}

func environmentSetTokenType(t *testing.T, args *setTestArgs) (string, error) {
	apim := args.srcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--token-type", args.tokenTypeFlag, "-k")
	return output, error
}

func importApiFromProject(t *testing.T, projectName string, envName string) (string, error) {
	projectPath, _ := filepath.Abs(projectName)
	return base.Execute(t, "import-api", "-f", projectPath, "-e", envName, "-k")
}

func importApiFromProjectWithUpdate(t *testing.T, projectName string, envName string) (string, error) {
	projectPath, _ := filepath.Abs(projectName)
	return base.Execute(t, "import-api", "-f", projectPath, "-e", envName, "-k", "--update")
}

func exportApisWithOneCommand(t *testing.T, args *initTestArgs) (string, error) {
	output, error := base.Execute(t, "export-apis", "-e", args.srcAPIM.GetEnvName(), "-k")
	return output, error
}
