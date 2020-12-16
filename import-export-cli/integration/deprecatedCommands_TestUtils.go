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
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func exportAPI(t *testing.T, name string, version string, provider string, env string) (string, error) {
	var output string
	var err error

	if provider == "" {
		output, err = base.Execute(t, "export-api", "-n", name, "-v", version, "-e", env, "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "export-api", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")
	}

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, testutils.GetEnvAPIExportPath(env), name, version)
	})

	return output, err
}

func importAPI(t *testing.T, args *testutils.ApiImportExportTestArgs) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.Api.Name, args.Api.Version)

	params := []string{"import-api", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose"}

	if args.OverrideProvider {
		params = append(params, "--preserve-provider=false")
	}

	if args.ParamsFile != "" {
		params = append(params, "--params", args.ParamsFile)
	}

	output, err := base.Execute(t, params...)

	t.Cleanup(func() {
		err := args.DestAPIM.DeleteAPIByName(args.Api.Name)

		if err != nil {
			t.Fatal(err)
		}
		base.WaitForIndexing()
	})

	return output, err
}

func listAPIs(t *testing.T, args *testutils.ApiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "list", "apis", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func listAPIProducts(t *testing.T, args *testutils.ApiProductImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "list", "api-products", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func exportAllApisOfATenant(t *testing.T, args *testutils.ApiImportExportTestArgs) (string, error) {
	//Setup environment
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	//Login to the environmeTestImportAndExportAPIWithJpegImagent
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, error := base.Execute(t, "export-apis", "-e", args.SrcAPIM.GetEnvName(), "-k", "--force")
	return output, error
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

func listApps(t *testing.T, env string) []string {
	response, _ := base.Execute(t, "list", "apps", "-e", env, "-k")

	return base.GetRowsFromTableResponse(response)
}

func getKeys(t *testing.T, provider string, name string, version string, env string) (string, error) {
	return base.Execute(t, "get-keys", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")
}

func validateAPIExportDeprecated(t *testing.T, args *testutils.ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, testutils.GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))
}

func validateAPIExportImportDeprecated(t *testing.T, args *testutils.ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.Api.Provider, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, testutils.GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAPI(t, args)

	// Give time for newly imported API to get indexed, or else getAPI by name will fail
	base.WaitForIndexing()

	// Get App from env 2
	importedAPI := testutils.GetAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	// Validate env 1 and env 2 API is equal
	testutils.ValidateAPIsEqual(t, args.Api, importedAPI)
}

func validateAPIsList(t *testing.T, args *testutils.ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List APIs of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIs(t, args)

	apisList := args.SrcAPIM.GetAPIs()

	testutils.ValidateListAPIsEqual(t, output, apisList)
}

func validateAllApisOfATenantIsExported(t *testing.T, args *testutils.ApiImportExportTestArgs, apisAdded int) {
	output, error := exportAllApisOfATenant(t, args)
	assert.Nil(t, error, "Error while exporting APIs")
	assert.Contains(t, output, "export-apis execution completed", "Error while exporting APIs")

	//Derive exported path from output
	exportedPath := base.GetExportedPathFromOutput(strings.ReplaceAll(output, "Command: export-apis execution completed !", ""))
	count, _ := base.CountFiles(exportedPath)
	assert.GreaterOrEqual(t, count, apisAdded, "Error while exporting APIs")

	t.Cleanup(func() {
		//Remove Exported apis and logout
		pathToCleanUp := utils.DefaultExportDirPath + testutils.TestMigrationDirectorySuffix
		base.RemoveDir(pathToCleanUp)
	})
}

func validateAPIProductsList(t *testing.T, args *testutils.ApiProductImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List API Products of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIProducts(t, args)

	apiProductsList := args.SrcAPIM.GetAPIProducts()

	testutils.ValidateListAPIProductsEqual(t, output, apiProductsList)
}

func validateAppExportFailure(t *testing.T, args *testutils.AppImportExportTestArgs) {
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

func validateAppExportImportWithPreserveOwner(t *testing.T, args *testutils.AppImportExportTestArgs) {
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
	importedApp := testutils.GetApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	testutils.ValidateAppsEqual(t, args.Application, importedApp)
}

func validateGetKeysFailure(t *testing.T, args *testutils.ApiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	var err error
	var result string

	result, err = getKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
}

func validateGetKeys(t *testing.T, args *testutils.ApiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	result, err := getKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while getting key")

	testutils.InvokeAPI(t, testutils.GetResourceURL(args.Apim, args.Api), base.GetValueOfUniformResponse(result), 200)
	testutils.UnsubscribeAPI(args.Apim, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)

}
