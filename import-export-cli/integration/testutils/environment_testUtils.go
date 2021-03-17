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
	"io/ioutil"
	"log"
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

func InitProjectWithOasFlag(t *testing.T, args *InitTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, err := base.Execute(t, "init", args.InitFlag, "--oas", args.OasFlag, "--verbose", "-f")
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

func genDeploymentDir(t *testing.T, args *GenDeploymentDirTestArgs) (string, error) {
	output, err := base.Execute(t, "gen", "deployment-dir", "-s", args.Source, "-d", args.Destination, "-k", "--verbose")

	t.Cleanup(func() {
		// Remove generated deployment directory
		base.RemoveDir(args.Destination)
	})

	return output, err
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

func ValidateGenDeploymentDir(t *testing.T, args *GenDeploymentDirTestArgs) {
	t.Helper()

	// Execute apictl command to generate the deployment directory for source project
	output, _ := genDeploymentDir(t, args)

	assert.Contains(t, output, "The deployment directory for "+args.Source+" file is generated at "+args.Destination+" directory",
		"Generating deployment directory is not successful")
}

func ValidateAPIImportExportWithDeploymentDir(t *testing.T, args *ApiImportExportTestArgs, api *apim.API) {

	// Move dummay params file of an API to the created deployment directory
	srcPathForParamsFile, _ := filepath.Abs(APIFullParamsFile)
	destPathForParamsFile := args.ParamsFile + string(os.PathSeparator) + utils.ParamFile
	utils.CopyFile(srcPathForParamsFile, destPathForParamsFile)

	// Copy dummy certificates to the created deployment directory
	srcPathForCertificatesDirectory, _ := filepath.Abs(CertificatesDirectoryPath)
	utils.CopyDirectoryContents(srcPathForCertificatesDirectory,
		args.ParamsFile+string(os.PathSeparator)+utils.DeploymentCertificatesDirectory)

	importedAPI := GetImportedAPI(t, args)

	apiParams := ReadParams(t, args.ParamsFile+string(os.PathSeparator)+utils.ParamFile)
	validateParamsWithoutCerts(t, apiParams, importedAPI, nil, importedAPI.Policies,
		importedAPI.GatewayEnvironments)

	args.SrcAPIM = args.DestAPIM // The API should be exported from prod env
	validateExportedAPICerts(t, apiParams, importedAPI, args)
}

func ValidateAPIProductImportExportWithDeploymentDir(t *testing.T, args *ApiProductImportExportTestArgs,
	apiProduct *apim.APIProduct) {

	// Move dummay params file of an API Product to the created deployment directory
	srcPathForParamsFile, _ := filepath.Abs(APIProductFullParamsFile)
	destPathForParamsFile := args.ParamsFile + string(os.PathSeparator) + utils.ParamFile
	utils.CopyFile(srcPathForParamsFile, destPathForParamsFile)

	srcPathForCertificatesDirectory, _ := filepath.Abs(CertificatesDirectoryPath)
	utils.CopyDirectoryContents(srcPathForCertificatesDirectory,
		args.ParamsFile+string(os.PathSeparator)+utils.DeploymentCertificatesDirectory)

	importedAPIProduct := ValidateAPIProductImport(t, args, true)

	apiProductParams := ReadParams(t, args.ParamsFile+string(os.PathSeparator)+utils.ParamFile)
	validateParamsWithoutCerts(t, apiProductParams, nil, importedAPIProduct, importedAPIProduct.Policies, importedAPIProduct.GatewayEnvironments)

	args.SrcAPIM = args.DestAPIM // The API Product should be exported from prod env
	validateExportedAPIProductCerts(t, apiProductParams, importedAPIProduct, args)
}

func ValidateDependentAPIWithParams(t *testing.T, dependentAPI *apim.API, client *apim.Client, username, password string) {

	importedDependentAPI := GetAPI(t, client, dependentAPI.Name, username, password)
	srcPathForParamsFile, _ := filepath.Abs(APIFullParamsFile)
	apiParams := ReadParams(t, srcPathForParamsFile)

	validateParamsWithoutCerts(t, apiParams, importedDependentAPI, nil, importedDependentAPI.Policies,
		importedDependentAPI.GatewayEnvironments)
}

func validateEndpointSecurity(t *testing.T, apiParams *Params, api *apim.API) {
	assert.Equal(t, strings.ToUpper(apiParams.Environments[0].Configs.Security.Type), api.EndpointSecurity.Type)
	assert.Equal(t, apiParams.Environments[0].Configs.Security.Username, api.EndpointSecurity.Username)
	assert.Equal(t, "", api.EndpointSecurity.Password)
}

func ValidateEndpointSecurityDefinition(t *testing.T, api *apim.API, apiParams *Params, importedAPI *apim.API) {
	t.Helper()

	validateEndpointSecurity(t, apiParams, importedAPI)

	assert.Equal(t, strings.ToUpper(apiParams.Environments[0].Configs.Security.Type), importedAPI.EndpointSecurity.Type)
	assert.Equal(t, apiParams.Environments[0].Configs.Security.Username, importedAPI.EndpointSecurity.Username)
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

	ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

func validateParamsWithoutCerts(t *testing.T, params *Params, api *apim.API, apiProduct *apim.APIProduct,
	policies, gatewayEnvironments []string) {
	t.Helper()

	// Endpoints and endpoint security will only be there for APIs, not for API Products
	if api != nil {
		// Validate endpoints
		assert.Equal(t, params.Environments[0].Configs.Endpoints.Production["url"], api.GetProductionURL(),
			"Mismatched productction URL")
		assert.Equal(t, params.Environments[0].Configs.Endpoints.Sandbox["url"], api.GetSandboxURL(),
			"Mismatched sandbox URL")

		// Validate endpoint security
		validateEndpointSecurity(t, params, api)
	}

	// Validate subscription policies
	assert.ElementsMatch(t, params.Environments[0].Configs.Policies, policies, "Mismatched policies")

	// Validate deployment environments
	validateDeploymentEnvironments(t, params, gatewayEnvironments)
}

func validateDeploymentEnvironments(t *testing.T, apiParams *Params, gatewayEnvironments []string) {

	assert.EqualValues(t, len(apiParams.Environments[0].Configs.DeploymentEnvironments), len(gatewayEnvironments),
		"Mismatched number of deployment environments")

	var deploymentEnvironments []string
	for _, deploymentEnvironmentFromParams := range apiParams.Environments[0].Configs.DeploymentEnvironments {
		deploymentEnvironments = append(deploymentEnvironments, deploymentEnvironmentFromParams.DeploymentEnvironment)
	}

	assert.ElementsMatch(t, deploymentEnvironments, gatewayEnvironments, "Mismatched deployment environments")
}

func validateExportedAPICerts(t *testing.T, apiParams *Params, api *apim.API, args *ApiImportExportTestArgs) {
	output, _ := exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	//Unzip exported API and check whether the imported certificates are there
	exportedPath := base.GetExportedPathFromOutput(output)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	pathOfExportedApi := relativePath + string(os.PathSeparator) + api.Name + "-" + api.Version

	validateEndpointCerts(t, apiParams, pathOfExportedApi)
	validateMutualSSLCerts(t, apiParams, pathOfExportedApi)

	t.Cleanup(func() {
		//Remove Created project and logout
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}

func validateExportedAPIProductCerts(t *testing.T, apiProductParams *Params, apiProduct *apim.APIProduct, args *ApiProductImportExportTestArgs) {
	output, _ := exportAPIProduct(t, args.ApiProduct.Name, utils.DefaultApiProductVersion, args.SrcAPIM.GetEnvName())

	//Unzip exported API Product and check whether the imported certificates are there
	exportedPath := base.GetExportedPathFromOutput(output)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	pathOfExportedApiProduct := relativePath + string(os.PathSeparator) + apiProduct.Name + "-" + utils.DefaultApiProductVersion

	validateMutualSSLCerts(t, apiProductParams, pathOfExportedApiProduct)

	t.Cleanup(func() {
		//Remove Created project and logout
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}

func validateEndpointCerts(t *testing.T, apiParams *Params, path string) {
	pathOfExportedEndpointCerts := path + string(os.PathSeparator) + utils.InitProjectEndpointCertificates
	isEndpointCertsDirExists, _ := utils.IsDirExists(pathOfExportedEndpointCerts)

	if isEndpointCertsDirExists {
		files, _ := ioutil.ReadDir(pathOfExportedEndpointCerts)
		for _, endpointCert := range apiParams.Environments[0].Configs.Certs {
			endpointCertExists := false
			for _, file := range files {
				if strings.EqualFold(file.Name(), endpointCert.Path) {
					endpointCertExists = true
				}
			}
			if !endpointCertExists {
				t.Error("Endpoint certificate " + endpointCert.Path + " not exported")
			}
		}
	} else {
		t.Error("Endpoint certificates directory does not exist")
	}
}

func validateMutualSSLCerts(t *testing.T, apiParams *Params, path string) {
	pathOfExportedMsslCerts := path + string(os.PathSeparator) + utils.InitProjectClientCertificates
	isClientCertsDirExists, _ := utils.IsDirExists(pathOfExportedMsslCerts)

	if isClientCertsDirExists {
		files, _ := ioutil.ReadDir(pathOfExportedMsslCerts)
		for _, msslCert := range apiParams.Environments[0].Configs.MsslCerts {
			msslCertExists := false
			for _, file := range files {
				if strings.EqualFold(file.Name(), msslCert.Path) {
					msslCertExists = true
				}
			}
			if !msslCertExists {
				t.Error("Client (MutualSSL) certificate " + msslCert.Path + " not exported")
			}
		}
	} else {
		t.Error("Client (MutualSSL) certificates directory does not exist")
	}
}
