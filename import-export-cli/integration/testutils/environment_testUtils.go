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
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func InitProjectWithOasFlag(t *testing.T, args *InitTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, err := base.Execute(t, "init", args.InitFlag, "--oas", args.OasFlag, "--verbose")
	return output, err
}

func EnvironmentSetExportDirectory(t *testing.T, args *SetTestArgs) (string, error) {
	apim := args.SrcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--export-directory", args.ExportDirectoryFlag, "-k", "--verbose")
	return output, error
}

func EnvironmentSetHttpRequestTimeout(t *testing.T, args *SetTestArgs) (string, error) {
	apim := args.SrcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--http-request-timeout", strconv.Itoa(args.httpRequestTimeout), "-k", "--verbose")
	return output, error
}

func EnvironmentSetTokenType(t *testing.T, args *SetTestArgs) (string, error) {
	apim := args.SrcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--token-type", args.TokenTypeFlag, "-k", "--verbose")
	return output, error
}

func ValidateThatRecievingTokenTypeIsChanged(t *testing.T, args *ApiGetKeyTestArgs, expectedTokenType string) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	var err error
	_, err = GetKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while getting key")

	tokenType := args.Apim.GetApplication(args.Apim.GetApplicationByName(DefaultApictlTestAppName).ApplicationID).TokenType
	assert.Equal(t, strings.ToUpper(expectedTokenType), tokenType, "Error getting token type of application.")

	UnsubscribeAPI(args.Apim, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)
}

func ValidateExportDirectoryIsChanged(t *testing.T, args *SetTestArgs) {
	t.Helper()
	output, _ := EnvironmentSetExportDirectory(t, args)
	base.Log(output)
	assert.Contains(t, output, "Export Directory is set to", "Export Directory change is not successful")
}

func ValidateExportApisPassed(t *testing.T, args *InitTestArgs, directoryName string) {
	t.Helper()

	output, error := ExportApisWithOneCommand(t, args)
	assert.Nil(t, error, "Error while Exporting APIs")
	assert.Contains(t, output, "export-apis execution completed", "Error while Exporting APIs")

	//Derive exported path from output
	exportedPath := base.GetExportedPathFromOutput(strings.ReplaceAll(output, "Command: export-apis execution completed !", ""))
	count, _ := base.CountFiles(exportedPath)
	assert.Equal(t, 1, count, "Error while Exporting APIs")

	t.Cleanup(func() {
		//Remove Exported apis
		base.RemoveDir(directoryName + TestMigrationDirectorySuffix)
	})
}

func ValidateExportApiPassed(t *testing.T, args *ApiImportExportTestArgs, directoryName string) {
	t.Helper()

	output, error := exportAPI(t, args.Api.Name, args.Api.Version, args.Api.Provider, args.SrcAPIM.EnvName)
	assert.Nil(t, error, "Error while Exporting APIs")
	assert.Contains(t, output, "Successfully exported API!", "Error while Exporting API")

	//Derive exported path from output
	exportedPath := filepath.Dir(base.GetExportedPathFromOutput(output))

	assert.True(t, strings.HasPrefix(exportedPath, directoryName), "API export path "+exportedPath+" is"+
		" not within the expected export location "+directoryName)

	assert.True(t, base.IsAPIArchiveExists(t, exportedPath, args.Api.Name, args.Api.Version), "API archive"+
		" is not correctly exported to "+directoryName)

	t.Cleanup(func() {
		//Remove Exported api
		base.RemoveDir(directoryName)
	})
}

func ValidateEndpointSecurityDefinition(t *testing.T, api *apim.API, apiParams *APIParams, importedAPI *apim.API) {
	t.Helper()

	validateEndpointSecurity(t, apiParams, importedAPI, "production")
	validateEndpointSecurity(t, apiParams, importedAPI, "sandbox")

	assert.Equal(t, strings.ToUpper(apiParams.Environments[0].Security.Type), importedAPI.EndpointSecurity.Type)
	assert.Equal(t, apiParams.Environments[0].Security.Username, importedAPI.EndpointSecurity.Username)
	assert.Equal(t, "", importedAPI.EndpointSecurity.Password)

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"

	apiCopy.EndpointSecurity.Type = same
	importedAPICopy.EndpointSecurity.Type = same

	apiCopy.EndpointSecurity.Username = same
	importedAPICopy.EndpointSecurity.Username = same

	apiCopy.EndpointSecurity.Password = same
	importedAPICopy.EndpointSecurity.Password = same

	apiCopy.EndpointConfig = same
	importedAPICopy.EndpointConfig = same
	ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

func validateEndpointSecurity(t *testing.T, apiParams *APIParams, api *apim.API, endpointType string) {
	endpointSecurityForEndpointType := apiParams.Environments[0].Security
	var endpointSecurityForEndpointTypeInApi map[string]interface{}

	if strings.EqualFold(endpointType, "production") {
		endpointSecurityForEndpointTypeInApi = api.GetProductionSecurityConfig()
	}

	if strings.EqualFold(endpointType, "sandbox") {
		endpointSecurityForEndpointTypeInApi = api.GetSandboxSecurityConfig()
	}

	assert.Equal(t, endpointSecurityForEndpointType.Enabled, endpointSecurityForEndpointTypeInApi["enabled"])
	assert.Equal(t, strings.ToUpper(endpointSecurityForEndpointType.Type), endpointSecurityForEndpointTypeInApi["type"])
	assert.Equal(t, endpointSecurityForEndpointType.Username, endpointSecurityForEndpointTypeInApi["username"])
	assert.Equal(t, "", endpointSecurityForEndpointTypeInApi["password"])
}

func ValidateOAuthEndpointSecurityDefinition(t *testing.T, api *apim.API, apiParams *APIParams, importedAPI *apim.API) {
	t.Helper()

	validateOauthEndpointSecurity(t, apiParams.Environments[0].Security.Production, importedAPI, "production")
	validateOauthEndpointSecurity(t, apiParams.Environments[0].Security.Sandbox, importedAPI, "sandbox")

	assert.Equal(t, strings.ToUpper(apiParams.Environments[0].Security.Type), importedAPI.EndpointSecurity.Type)

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"

	apiCopy.EndpointConfig = same
	importedAPICopy.EndpointConfig = same
	ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

func validateOauthEndpointSecurity(t *testing.T, envSecurityPerEndpoint OAuthEndpointSecurity, api *apim.API, endpointType string) {
	var endpointSecurityForEndpointTypeInApi map[string]interface{}

	if strings.EqualFold(endpointType, "production") {
		endpointSecurityForEndpointTypeInApi = api.GetProductionSecurityConfig()
	}

	if strings.EqualFold(endpointType, "sandbox") {
		endpointSecurityForEndpointTypeInApi = api.GetSandboxSecurityConfig()
	}

	assert.Equal(t, envSecurityPerEndpoint.Enabled, endpointSecurityForEndpointTypeInApi["enabled"])

	if envSecurityPerEndpoint.Enabled {
		assert.Equal(t, strings.ToUpper(envSecurityPerEndpoint.Type), endpointSecurityForEndpointTypeInApi["type"])
		assert.Equal(t, envSecurityPerEndpoint.ClientId, endpointSecurityForEndpointTypeInApi["clientId"])
		assert.Equal(t, envSecurityPerEndpoint.ClientSecret, endpointSecurityForEndpointTypeInApi["clientSecret"])
		assert.Equal(t, envSecurityPerEndpoint.TokenUrl, endpointSecurityForEndpointTypeInApi["tokenUrl"])
		assert.Equal(t, strings.ToUpper(envSecurityPerEndpoint.GrantType), endpointSecurityForEndpointTypeInApi["grantType"])

		if strings.EqualFold(strings.ToUpper(envSecurityPerEndpoint.GrantType), utils.PasswordGrantType) {
			assert.Equal(t, envSecurityPerEndpoint.Username, endpointSecurityForEndpointTypeInApi["username"])
			assert.Equal(t, "", endpointSecurityForEndpointTypeInApi["password"])
		}
	} else {
		assert.Equal(t, "NONE", endpointSecurityForEndpointTypeInApi["type"])
	}
}
