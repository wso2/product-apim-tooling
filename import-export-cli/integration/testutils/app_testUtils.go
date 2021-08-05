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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func AddAppWithSpaceInAppName(t *testing.T, client *apim.Client, username string, password string) *apim.Application {
	client.Login(username, password)
	app := client.GenerateSampleAppWithNameInSpaceData()
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

func GenerateKeys(t *testing.T, client *apim.Client, username, password, appId, keyType string) apim.ApplicationKey {
	client.Login(username, password)
	generateKeyReq := utils.KeygenRequest{
		KeyType:                 keyType,
		GrantTypesToBeSupported: utils.GrantTypesToBeSupported,
		ValidityTime:            utils.DefaultTokenValidityPeriod,
	}
	keyGenResponse := client.GenerateKeys(t, generateKeyReq, appId)
	return keyGenResponse
}

func getApp(t *testing.T, client *apim.Client, name string, username string, password string) *apim.Application {
	client.Login(username, password)
	appInfo := client.GetApplicationByName(name)
	return client.GetApplication(appInfo.ApplicationID)
}

func getOauthKeys(t *testing.T, client *apim.Client, username, password string,
	application *apim.Application) *apim.ApplicationKeysList {
	client.Login(username, password)
	applicationKeysList := client.GetOauthKeys(t, application)
	return applicationKeysList
}

func ListApps(t *testing.T, env string) []string {
	response, _ := base.Execute(t, "list", "apps", "-e", env, "-k")

	return base.GetRowsFromTableResponse(response)
}

func ListAppsWithOwner(t *testing.T, env string, owner string) []string {
	response, _ := base.Execute(t, "list", "apps", "-e", env, "-k", "--owner", owner)

	return base.GetRowsFromTableResponse(response)
}

func getEnvAppExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedAppsDirName, envName)
}

func exportApp(t *testing.T, args *AppImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "export-app", "-n", args.Application.Name, "-o", args.AppOwner.Username,
		"--withKeys="+strconv.FormatBool(args.WithKeys), "-e", args.SrcAPIM.GetEnvName(), "-k", "--verbose")

	t.Cleanup(func() {
		base.RemoveApplicationArchive(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
			args.Application.Name, args.AppOwner.Username)
	})

	return output, err
}

func importApp(t *testing.T, args *AppImportExportTestArgs, doClean bool) (string, error) {
	var fileName string
	if args.ImportFilePath == "" {
		fileName = base.GetApplicationArchiveFilePath(t, args.SrcAPIM.EnvName, args.Application.Name,
			args.Application.Owner)
	} else {
		fileName = args.ImportFilePath
	}

	output, err := base.Execute(t, "import-app", "-f", fileName, "--preserveOwner="+strconv.FormatBool(args.PreserveOwner),
		"--update="+strconv.FormatBool(args.UpdateFlag), "--skipKeys="+strconv.FormatBool(args.SkipKeys),
		"--skipSubscriptions="+strconv.FormatBool(args.SkipSubscriptions), "-e", args.DestAPIM.EnvName, "-k", "--verbose")

	if doClean {
		t.Cleanup(func() {
			args.DestAPIM.DeleteApplicationByName(args.Application.Name)
		})
	}

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

func ValidateAppExportImportWithPreserveOwner(t *testing.T, args *AppImportExportTestArgs) {
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

	importApp(t, args, true)

	// Get App from env 2
	importedApp := getApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	validateAppsEqual(t, args.Application, importedApp)
}

func ValidateAppExportImportWithUpdate(t *testing.T, args *AppImportExportTestArgs) {
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

	importApp(t, args, true)

	// Get App from env 2
	importedApp := getApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	validateAppsEqual(t, args.Application, importedApp)
}

func ValidateAppExport(t *testing.T, args *AppImportExportTestArgs) string {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting app from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, _ := exportApp(t, args)

	// Validate that export passed
	assert.True(t, base.IsApplicationArchiveExists(t, getEnvAppExportPath(args.SrcAPIM.GetEnvName()),
		args.Application.Name, args.AppOwner.Username))

	return output
}

func ValidateAppImport(t *testing.T, args *AppImportExportTestArgs, doClean bool) *apim.Application {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import app to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importApp(t, args, doClean)

	// Get App from env 2
	importedApp := getApp(t, args.DestAPIM, args.Application.Name, args.AppOwner.Username, args.AppOwner.Password)

	// Validate env 1 and env 2 App is equal
	validateAppsEqual(t, args.Application, importedApp)

	return importedApp
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
	response := ListAppsWithOwner(t, envName, "admin")
	assert.Equal(t, 5, len(response), "Failed when listing Applications with owner as Admin")

	emptyResponse := ListAppsWithOwner(t, envName, "user1")
	assert.Equal(t, 0, len(emptyResponse), "Failed when listing Applications with owner as User1")
}

func ValidateAppAdditionalPropertiesOfKeysUpdateImport(t *testing.T, args *AppImportExportTestArgs, doClean bool) {

	// Construct the exported application path
	mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	exportedAppPath := mainConfig.Config.ExportDirectory + string(os.PathSeparator) +
		utils.ExportedAppsDirName + string(os.PathSeparator) +
		base.GetApplicationArchiveFilePath(t, args.SrcAPIM.EnvName, args.Application.Name, args.Application.Owner)

	// Unzip exported application
	relativePath := strings.ReplaceAll(exportedAppPath, ".zip", "")
	base.Unzip(relativePath, exportedAppPath)

	args.ImportFilePath = relativePath + string(os.PathSeparator) + args.Application.Name

	keyManagerWiseOAuthApp := updateAdditionalPropertiesOfKeys(t, args)

	// Make the update flag true
	args.UpdateFlag = true
	updatedImportedApp := ValidateAppImport(t, args, false)

	// Retrieve oauth keys of the updated application
	updatedApplicationKeysList := getOauthKeys(t, args.DestAPIM, args.AppOwner.Username, args.AppOwner.Password, updatedImportedApp)

	for _, updatedKey := range updatedApplicationKeysList.List {
		additionalProperties := keyManagerWiseOAuthApp[updatedKey.KeyType].(map[string]interface{})[ResidentKeyManager].(map[string]interface{})["parameters"].(map[string]interface{})["additionalProperties"].(map[string]interface{})
		assert.EqualValues(t, additionalProperties["id_token_expiry_time"],
			updatedKey.AdditionalProperties["id_token_expiry_time"], updatedKey.KeyType+" id_token_expiry_time mismatched")
		assert.EqualValues(t, additionalProperties["application_access_token_expiry_time"],
			updatedKey.AdditionalProperties["application_access_token_expiry_time"], updatedKey.KeyType+" application_access_token_expiry_time mismatched")
		assert.EqualValues(t, additionalProperties["user_access_token_expiry_time"],
			updatedKey.AdditionalProperties["user_access_token_expiry_time"], updatedKey.KeyType+" user_access_token_expiry_time mismatched")
		assert.EqualValues(t, additionalProperties["refresh_token_expiry_time"],
			updatedKey.AdditionalProperties["refresh_token_expiry_time"], updatedKey.KeyType+" refresh_token_expiry_time mismatched")
	}

	if doClean {
		t.Cleanup(func() {
			// Remove extracted directory
			base.RemoveDir(relativePath)
		})
	}
}

func updateAdditionalPropertiesOfKeys(t *testing.T, args *AppImportExportTestArgs) map[string]interface{} {
	applicationDefinitionFilePath := args.ImportFilePath + string(os.PathSeparator) + args.Application.Name + ".json"
	// Read the application.yaml file in the exported directory
	applicationData, err := ioutil.ReadFile(applicationDefinitionFilePath)
	if err != nil {
		t.Error(err)
	}

	// Extract the content to a structure
	applicationContent := make(map[string]interface{})
	err = json.Unmarshal(applicationData, &applicationContent)
	if err != nil {
		t.Error(err)
	}

	updatedAdditionalPropertiesProduction := map[string]interface{}{
		"id_token_expiry_time":                 5001,
		"application_access_token_expiry_time": 5002,
		"user_access_token_expiry_time":        5003,
		"refresh_token_expiry_time":            5004,
	}

	updatedAdditionalPropertiesSandbox := map[string]interface{}{
		"id_token_expiry_time":                 5005,
		"application_access_token_expiry_time": 5006,
		"user_access_token_expiry_time":        5007,
		"refresh_token_expiry_time":            5008,
	}

	applicationContent["keyManagerWiseOAuthApp"].(map[string]interface{})[utils.ProductionKeyType].(map[string]interface{})[ResidentKeyManager].(map[string]interface{})["parameters"].(map[string]interface{})["additionalProperties"] = updatedAdditionalPropertiesProduction
	applicationContent["keyManagerWiseOAuthApp"].(map[string]interface{})[utils.SandboxKeyType].(map[string]interface{})[ResidentKeyManager].(map[string]interface{})["parameters"].(map[string]interface{})["additionalProperties"] = updatedAdditionalPropertiesSandbox

	updatedApplicationData, err := json.Marshal(applicationContent)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile(applicationDefinitionFilePath, updatedApplicationData, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	return applicationContent["keyManagerWiseOAuthApp"].(map[string]interface{})
}
