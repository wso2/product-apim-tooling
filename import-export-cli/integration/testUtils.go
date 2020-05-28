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
	id := client.AddAPI(t, api, username, password)
	api = client.GetAPI(id)
	return api
}

func addAPIToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.API, *apim.API) {
	client1.Login(username, password)
	api := client1.GenerateSampleAPIData(username)
	id1 := client1.AddAPI(t, api, username, password)
	api1 := client1.GetAPI(id1)

	client2.Login(username, password)
	id2 := client2.AddAPI(t, api, username, password)
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
	assert.Equal(t, "Exit status 1", base.GetValueOfUniformResponse(result))
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

func exportAPI(t *testing.T, name string, version string, provider string, env string) (string, error) {
	output, err := base.Execute(t, "export-api", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")

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

func importAPIPreserveProvider(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import-api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	t.Cleanup(func() {
		client.DeleteAPIByName(api.Name)
	})

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

	apiProductsListBeforeDelete := args.srcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	apiProductsListAfterDelete := args.srcAPIM.GetAPIProducts()

	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count+1, "API Product delete is not successful")
}

func validateAPIProductDeleteFailure(t *testing.T, args *apiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL(), args.srcAPIM.GetTokenURL())

	// Delete an API Product of env 1
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.username, args.ctlUser.password)

	apiProductsListBeforeDelete := args.srcAPIM.GetAPIProducts()

	deleteAPIProductByCtl(t, args)

	apiProductsListAfterDelete := args.srcAPIM.GetAPIProducts()

	assert.Equal(t, apiProductsListBeforeDelete.Count, apiProductsListAfterDelete.Count, "API Product delete is successful")
}
