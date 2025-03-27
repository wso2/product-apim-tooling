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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// InitApictl : Initializes apictl in the local machine and create directories
func InitApictl(t *testing.T) (string, error) {
	return base.Execute(t)
}

// SetApictlWithCustomDirectory : Set custom directory location using environment variable
func SetApictlWithCustomDirectory(t *testing.T, customDirPath string) {

	t.Log("Setting up the environment variable value for " + EnvVariableNameOfCustomCustomDirectoryAtInit)
	os.Setenv(EnvVariableNameOfCustomCustomDirectoryAtInit, customDirPath)
}

// ValidateApictlInit : Check and verify whether the apictl is initialized properly
func ValidateApictlInit(t *testing.T, err error, output string) {

	// Asserting the apictl initializing message to check whether apictl is initialized or not.
	assert.Nil(t, err)
	assert.Contains(t, output, ApictlInitMessage, "Test failed without initializing apictl")
}

// ValidateCustomDirectoryChangeAtInit :  Validate custom directory change at init
func ValidateCustomDirectoryChangeAtInit(t *testing.T, customDirPath string) {

	//Check .wso2apictl and .wso2apictl.local directories on custom directories
	checkConfigDir, _ := utils.IsDirExists(filepath.Join(customDirPath, utils.ConfigDirName))
	checkLocalCredentialsDir, _ := utils.IsDirExists(filepath.Join(customDirPath, utils.LocalCredentialsDirectoryName))

	assert.Equal(t, true, checkConfigDir)
	assert.Equal(t, true, checkLocalCredentialsDir)

	t.Cleanup(func() {
		// Remove created custom directory directory
		os.Unsetenv(EnvVariableNameOfCustomCustomDirectoryAtInit)
		base.RemoveDir(customDirPath)
	})
}

// AddEnvironmentWithTokenFlag : Add new environment with token endpoint url
func AddEnvironmentWithTokenFlag(t *testing.T, envName, apimUrl, tokenUrl string) (string, error) {
	return base.Execute(t, "add", "env", envName, "--apim", apimUrl, "--token", tokenUrl)
}

// AddEnvironmentWithOutTokenFlag : Add new environment without token endpoint url
func AddEnvironmentWithOutTokenFlag(t *testing.T, envName, apimUrl string) (string, error) {
	return base.Execute(t, "add", "env", envName, "--apim", apimUrl)
}

// RemoveEnvironment : Remove an added environment
func RemoveEnvironment(t *testing.T, envName string) (string, error) {
	return base.Execute(t, "remove", "env", envName)
}

// ValidateAddedEnvironments : Check whether the added environments are added as expected when listed
func ValidateAddedEnvironments(t *testing.T, output, envName string, skipCleanup bool) {

	expectedOutput := fmt.Sprintf(`Successfully added environment '%v'`, envName)
	assert.Contains(t, output, expectedOutput)

	//List all the environments and check for availability of the added environment
	ValidateEnvsList(t, envName, true)

	if skipCleanup == false {
		t.Cleanup(func() {
			base.Execute(t, "remove", "env", envName)
		})
	}
}

// ValidateRemoveEnvironments : Check whether the added environments is removed
func ValidateRemoveEnvironments(t *testing.T, output, envName string) {

	expectedOutput := fmt.Sprintf(`Successfully removed environment '%v'`, envName)
	assert.Contains(t, output, expectedOutput)

	//List all the environments and check for availability of the removed environment
	ValidateEnvsList(t, envName, false)
}

// ValidateEnvsList : Check the list and verify the given env is in the list or not
func ValidateEnvsList(t *testing.T, envName string, checkContains bool) {
	// List environments
	response, err := base.Execute(t, "get", "envs")
	assert.Nil(t, err)

	base.GetRowsFromTableResponse(response)
	base.Log(response)

	//Check added environment in the list
	if checkContains == true {
		assert.Contains(t, response, envName, "TestGetEnvironments Failed")
	} else {
		assert.NotContains(t, response, envName, "TestGetEnvironments Failed")
	}
}

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

	count, _ := base.CountFiles(t, exportedPath)
	assert.Equal(t, 1, count, "Error while Exporting APIs")

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

func ValidateAPIImportExportWithDeploymentDir(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Move dummy params and certificates files of an API to the created deployment directory
	MoveDummyAPIParamsAndCertificatesToDeploymentDir(args)

	importedAPI := GetImportedAPI(t, args)

	apiRevisions := GetDeployedAPIRevisions(t, args.DestAPIM, args.CtlUser.Username,
		args.CtlUser.Password, importedAPI.ID)
	gatewayEnvironments := GetGatewayEnvironments(apiRevisions)

	apiParams := ReadParams(t, args.ParamsFile+string(os.PathSeparator)+utils.ParamFile)
	validateParamsWithoutCerts(t, apiParams, importedAPI, nil, importedAPI.Policies, gatewayEnvironments)

	args.SrcAPIM = args.DestAPIM // The API should be exported from prod env
	validateExportedAPICerts(t, apiParams, importedAPI, args)
}

func ValidateAPIProductImportExportWithDeploymentDir(t *testing.T, args *ApiProductImportExportTestArgs) {

	// Move dummay params file of an API Product to the created deployment directory
	srcPathForParamsFile, _ := filepath.Abs(APIProductFullParamsFile)
	destPathForParamsFile := args.ParamsFile + string(os.PathSeparator) + utils.ParamFile
	utils.CopyFile(srcPathForParamsFile, destPathForParamsFile)

	srcPathForCertificatesDirectory, _ := filepath.Abs(CertificatesDirectoryPath)
	utils.CopyDirectoryContents(srcPathForCertificatesDirectory,
		args.ParamsFile+string(os.PathSeparator)+utils.DeploymentCertificatesDirectory)

	importedAPIProduct := ValidateAPIProductImport(t, args, true)

	apiRevisions := GetDeployedAPIProductRevisions(t, args.DestAPIM, args.CtlUser.Username,
		args.CtlUser.Password, importedAPIProduct.ID)
	gatewayEnvironments := GetGatewayEnvironments(apiRevisions)

	apiProductParams := ReadParams(t, args.ParamsFile+string(os.PathSeparator)+utils.ParamFile)
	validateParamsWithoutCerts(t, apiProductParams, nil, importedAPIProduct, importedAPIProduct.Policies, gatewayEnvironments)

	args.SrcAPIM = args.DestAPIM // The API Product should be exported from prod env
	validateExportedAPIProductCerts(t, apiProductParams, importedAPIProduct, args)
}

func ValidateAPIImportExportWithDeploymentDirForAdvertiseOnlyAPI(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Move dummy params and certificates files of an API to the created deployment directory
	MoveDummyAPIParamsAndCertificatesToDeploymentDir(args)

	importedAPI := GetImportedAPI(t, args)

	validateAdvertiseOnlyAPIsEqual(t, importedAPI, args)
}

func MoveDummyAPIParamsAndCertificatesToDeploymentDir(args *ApiImportExportTestArgs) {
	// Move dummay params file of an API to the created deployment directory
	srcPathForParamsFile, _ := filepath.Abs(APIFullParamsFile)
	destPathForParamsFile := args.ParamsFile + string(os.PathSeparator) + utils.ParamFile
	utils.CopyFile(srcPathForParamsFile, destPathForParamsFile)

	// Copy dummy certificates to the created deployment directory
	srcPathForCertificatesDirectory, _ := filepath.Abs(CertificatesDirectoryPath)
	utils.CopyDirectoryContents(srcPathForCertificatesDirectory,
		args.ParamsFile+string(os.PathSeparator)+utils.DeploymentCertificatesDirectory)
}

func ValidateDependentAPIWithParams(t *testing.T, dependentAPI *apim.API, client *apim.Client, username, password string) {

	importedDependentAPI := GetAPI(t, client, dependentAPI.Name, username, password)
	srcPathForParamsFile, _ := filepath.Abs(APIFullParamsFile)
	apiParams := ReadParams(t, srcPathForParamsFile)

	validateParamsWithoutCerts(t, apiParams, importedDependentAPI, nil, importedDependentAPI.Policies,
		nil)
}

func validateEndpointSecurity(t *testing.T, apiParams *Params, api *apim.API, endpointType string) {
	var endpointSecurityForEndpointType Security
	var endpointSecurityForEndpointTypeInApi map[string]interface{}
	if strings.EqualFold(endpointType, "production") {
		endpointSecurityForEndpointType = apiParams.Environments[0].Configs.Security.Production
		endpointSecurityForEndpointTypeInApi = api.GetProductionSecurityConfig()
	}
	if strings.EqualFold(endpointType, "sandbox") {
		endpointSecurityForEndpointType = apiParams.Environments[0].Configs.Security.Sandbox
		endpointSecurityForEndpointTypeInApi = api.GetSandboxSecurityConfig()
	}

	assert.Equal(t, endpointSecurityForEndpointType.Enabled, endpointSecurityForEndpointTypeInApi["enabled"])

	if endpointSecurityForEndpointType.Enabled {
		assert.Equal(t, strings.ToUpper(endpointSecurityForEndpointType.Type), endpointSecurityForEndpointTypeInApi["type"])

		if strings.EqualFold(strings.ToUpper(endpointSecurityForEndpointType.Type), EndpointSecurityTypeOAuth) {
			// Validate Oauth 2.0 endpoint security related properties
			assert.Equal(t, endpointSecurityForEndpointType.ClientId, endpointSecurityForEndpointTypeInApi["clientId"])
			assert.Equal(t, "", endpointSecurityForEndpointTypeInApi["clientSecret"])
			assert.Equal(t, endpointSecurityForEndpointType.TokenUrl, endpointSecurityForEndpointTypeInApi["tokenUrl"])
			assert.Equal(t, strings.ToUpper(endpointSecurityForEndpointType.GrantType), endpointSecurityForEndpointTypeInApi["grantType"])

			if strings.EqualFold(strings.ToUpper(endpointSecurityForEndpointType.GrantType), PasswordGrantType) {
				validateEndpointSecurityUsernamePassword(t, endpointSecurityForEndpointType,
					endpointSecurityForEndpointTypeInApi)
			}
		}

		if strings.EqualFold(strings.ToUpper(endpointSecurityForEndpointType.Type), EndpointSecurityTypeBasic) ||
			strings.EqualFold(strings.ToUpper(endpointSecurityForEndpointType.Type), EndpointSecurityTypeDigest) {
			// Validate basic or digest endpoint security related properties
			validateEndpointSecurityUsernamePassword(t, endpointSecurityForEndpointType,
				endpointSecurityForEndpointTypeInApi)
		}
	}
}

func validateEndpointSecurityUsernamePassword(t *testing.T, endpointSecurityForEndpointType Security,
	endpointSecurityForEndpointTypeInApi map[string]interface{}) {
	assert.Equal(t, endpointSecurityForEndpointType.Username, endpointSecurityForEndpointTypeInApi["username"])
	assert.Equal(t, "", endpointSecurityForEndpointTypeInApi["password"])
}

func ValidateEndpointSecurityDefinition(t *testing.T, api *apim.API, apiParams *Params, importedAPI *apim.API) {
	t.Helper()

	validateEndpointSecurity(t, apiParams, importedAPI, "production")
	validateEndpointSecurity(t, apiParams, importedAPI, "sandbox")

	api.EndpointConfig.(map[string]interface{})["endpoint_security"] = "override_with_the_same_value"
	importedAPI.EndpointConfig.(map[string]interface{})["endpoint_security"] = "override_with_the_same_value"

	ValidateAPIsEqualWithEndpointConfigsFromParam(t, api, importedAPI, apiParams)
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
		validateEndpointSecurity(t, params, api, "production")
		validateEndpointSecurity(t, params, api, "sandbox")
	}

	// Validate subscription policies
	assert.ElementsMatch(t, params.Environments[0].Configs.Policies, policies, "Mismatched policies")

	// This will be nil for dependent APIs of an API Product
	if gatewayEnvironments != nil {
		// Validate deployment environments
		validateDeploymentEnvironments(t, params, gatewayEnvironments)
	}
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

	validateEndpointCerts(t, apiParams, args.DestAPIM, args.ApiProvider, pathOfExportedApi)
	validateMutualSSLCerts(t, apiParams, pathOfExportedApi)

	t.Cleanup(func() {
		//Remove Created project and logout
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}

func validateNonExportedAPICerts(t *testing.T, api *apim.API, args *ApiImportExportTestArgs) {
	output, _ := exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	//Unzip exported API and check whether certificates are there
	exportedPath := base.GetExportedPathFromOutput(output)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	pathOfExportedApi := relativePath + string(os.PathSeparator) + api.Name + "-" + api.Version

	pathOfExportedEndpointCerts := pathOfExportedApi + string(os.PathSeparator) + utils.InitProjectEndpointCertificates
	isEndpointCertsDirExists, _ := utils.IsDirExists(pathOfExportedEndpointCerts)
	assert.Equal(t, false, isEndpointCertsDirExists)

	pathOfExportedClientCerts := pathOfExportedApi + string(os.PathSeparator) + utils.InitProjectClientCertificates
	isClientCertsDirExists, _ := utils.IsDirExists(pathOfExportedClientCerts)
	assert.Equal(t, false, isClientCertsDirExists)

	t.Cleanup(func() {
		// Remove Created project
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}

func validateExportedAPIProductCerts(t *testing.T, apiProductParams *Params, apiProduct *apim.APIProduct, args *ApiProductImportExportTestArgs) {
	output, _ := exportAPIProduct(t, args.ApiProduct.Name, args.ApiProduct.Version, args.SrcAPIM.GetEnvName())

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

func validateEndpointCerts(t *testing.T, apiParams *Params, client *apim.Client, credentials Credentials, path string) {
	pathOfExportedEndpointCerts := path + string(os.PathSeparator) + utils.InitProjectEndpointCertificates

	t.Log("validateEndpointCerts() pathOfExportedEndpointCerts = ", pathOfExportedEndpointCerts)
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
			} else {
				t.Cleanup(func() {
					client.Login(credentials.Username, credentials.Password)
					client.RemoveEndpointCert(endpointCert.Alias)
				})
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

func validateEndpointType(t *testing.T, apiParams *Params, api *apim.API) {
	var endPointTypeInParams string
	endPointTypeInParams = apiParams.Environments[0].Configs.EndpointType
	if strings.EqualFold(endPointTypeInParams, "rest") {
		endPointTypeInParams = "http"
	} else if strings.EqualFold(endPointTypeInParams, "soap") {
		endPointTypeInParams = "address"
	}
	assert.Equal(t, endPointTypeInParams, api.GetEndpointType())
}

func validateEndpointUrl(t *testing.T, apiParams *Params, api *apim.API, endpointType string) {
	var endpointTypeRecord map[string]interface{}

	if strings.EqualFold(endpointType, "production") {
		endpointTypeRecord = apiParams.Environments[0].Configs.Endpoints.Production
		assert.Equal(t, endpointTypeRecord["url"], api.GetSandboxURL())
	}
	if strings.EqualFold(endpointType, "sandbox") {
		endpointTypeRecord = apiParams.Environments[0].Configs.Endpoints.Sandbox
		assert.Equal(t, endpointTypeRecord["url"], api.GetSandboxURL())
	}
}

func ValidateHttpEndpointWithoutLoadBalancingAndFailover(t *testing.T, apiParams *Params, api, importedAPI *apim.API) {
	t.Helper()
	//Validate EndPoint Type
	validateEndpointType(t, apiParams, importedAPI)

	// Validate endpoint type URLs
	validateEndpointUrl(t, apiParams, importedAPI, "production")
	validateEndpointUrl(t, apiParams, importedAPI, "sandbox")

	same := "override_with_same_value"
	api.SetEndpointType(same)
	api.SetProductionURL(same)
	api.SetSandboxURL(same)
	importedAPI.SetEndpointType(same)
	importedAPI.SetProductionURL(same)
	importedAPI.SetSandboxURL(same)

	ValidateAPIsEqual(t, api, importedAPI)
}

func ValidateHttpEndpointWithLoadBalancing(t *testing.T, apiParams *Params, api, importedAPI *apim.API) {
	t.Helper()
	//Validate EndPoint Type
	assert.Equal(t, "load_balance", importedAPI.GetEndpointType())

	endPointConfigInApi := importedAPI.EndpointConfig

	// Validate Algorithm class
	algoClassNameInParams := apiParams.Environments[0].Configs.LoadBalanceEndpoints.AlgorithmClassName
	algoClassNameInApi := endPointConfigInApi.(map[string]interface{})["algoClassName"].(string)
	assert.Equal(t, algoClassNameInParams, algoClassNameInApi)

	// Validate Session TimeOut
	sessionTimeOutInParams := apiParams.Environments[0].Configs.LoadBalanceEndpoints.SessionTimeout
	sessionTimeOutInApi := endPointConfigInApi.(map[string]interface{})["sessionTimeOut"].(string)
	assert.Equal(t, strconv.Itoa(sessionTimeOutInParams), sessionTimeOutInApi)

	// Validate Production endpoints
	productionEndpointsInParams := apiParams.Environments[0].Configs.LoadBalanceEndpoints.Production
	productionEndpointsInApi := endPointConfigInApi.(map[string]interface{})["production_endpoints"].([]interface{})
	if strings.EqualFold(apiParams.Environments[0].Configs.LoadBalanceEndpoints.SessionManagement, "transport") {
		for i := 0; i < len(productionEndpointsInParams); i++ {
			singleProductionEndpointInApi := productionEndpointsInApi[i].(map[string]interface{})
			assert.Equal(t, productionEndpointsInParams[i]["url"], singleProductionEndpointInApi["url"])
		}
	}
	if strings.EqualFold(apiParams.Environments[0].Configs.LoadBalanceEndpoints.SessionManagement, "soap") {
		for i := 0; i < len(productionEndpointsInParams); i++ {
			singleProductionEndpointInApi := productionEndpointsInApi[i].(map[string]interface{})
			assert.Equal(t, productionEndpointsInParams[i]["url"], singleProductionEndpointInApi["url"])
			assert.Equal(t, "address", singleProductionEndpointInApi["endpoint_type"])
		}
	}

	// Validate Sandbox endpoints
	sandboxEndpointsInParams := apiParams.Environments[0].Configs.LoadBalanceEndpoints.Sandbox
	sandboxEndpointsInApi := endPointConfigInApi.(map[string]interface{})["sandbox_endpoints"].([]interface{})
	if strings.EqualFold(apiParams.Environments[0].Configs.LoadBalanceEndpoints.SessionManagement, "transport") {
		for i := 0; i < len(sandboxEndpointsInParams); i++ {
			singleSandboxEndpointInApi := sandboxEndpointsInApi[i].(map[string]interface{})
			assert.Equal(t, sandboxEndpointsInParams[i]["url"], singleSandboxEndpointInApi["url"])
		}
	}
	if strings.EqualFold(apiParams.Environments[0].Configs.LoadBalanceEndpoints.SessionManagement, "soap") {
		for i := 0; i < len(sandboxEndpointsInParams); i++ {
			singleSandboxEndpointInApi := sandboxEndpointsInApi[i].(map[string]interface{})
			assert.Equal(t, sandboxEndpointsInParams[i]["url"], singleSandboxEndpointInApi["url"])
			assert.Equal(t, "address", singleSandboxEndpointInApi["endpoint_type"])
		}
	}

	same := "override_with_same_value"
	api.SetEndpointType(same)
	importedAPI.SetEndpointType(same)
	api.SetEndPointConfig(same)
	importedAPI.SetEndPointConfig(same)

	ValidateAPIsEqual(t, api, importedAPI)
}

func ValidateHttpEndpointWithFailover(t *testing.T, apiParams *Params, api, importedAPI *apim.API) {
	t.Helper()

	var isSaopEndpoint bool = false

	//Validate EndPoint Type
	assert.Equal(t, "failover", importedAPI.GetEndpointType())

	endPointConfigInApi := importedAPI.EndpointConfig
	//Check whether the endpoint are SOAP endpoints
	if strings.EqualFold(apiParams.Environments[0].Configs.EndpointType, "soap") {
		isSaopEndpoint = true
	}

	// Validate Production endpoints
	productionEndpointInParams := apiParams.Environments[0].Configs.FailoverEndpoints.Production
	productionEndpointInApi := endPointConfigInApi.(map[string]interface{})["production_endpoints"].(map[string]interface{})
	assert.Equal(t, productionEndpointInParams.URL, productionEndpointInApi["url"])
	if isSaopEndpoint {
		assert.Equal(t, "address", productionEndpointInApi["endpoint_type"])
	}

	// Validate Production failover endpoints
	productionFailoverEndpointsInParams := apiParams.Environments[0].Configs.FailoverEndpoints.ProductionFailovers
	productionFailoverEndpointsInApi := endPointConfigInApi.(map[string]interface{})["production_failovers"].([]interface{})
	for i := 0; i < len(productionFailoverEndpointsInParams); i++ {
		singleProductionFailoverEndpointInApi := productionFailoverEndpointsInApi[i].(map[string]interface{})
		assert.Equal(t, productionFailoverEndpointsInParams[i]["url"], singleProductionFailoverEndpointInApi["url"])
		if isSaopEndpoint {
			assert.Equal(t, "address", productionEndpointInApi["endpoint_type"])
		}
	}

	// Validate Production endpoints
	sandboxEndpointInParams := apiParams.Environments[0].Configs.FailoverEndpoints.Sandbox
	sandboxEndpointInApi := endPointConfigInApi.(map[string]interface{})["sandbox_endpoints"].(map[string]interface{})
	assert.Equal(t, sandboxEndpointInParams.URL, sandboxEndpointInApi["url"])
	if isSaopEndpoint {
		assert.Equal(t, "address", productionEndpointInApi["endpoint_type"])
	}

	// Validate Sandbox failover endpoints
	sandboxFailoverEndpointsInParams := apiParams.Environments[0].Configs.FailoverEndpoints.SandboxFailovers
	sandboxFailoverEndpointsInApi := endPointConfigInApi.(map[string]interface{})["sandbox_failovers"].([]interface{})
	for i := 0; i < len(sandboxFailoverEndpointsInParams); i++ {
		singleSandboxFailoverEndpointInApi := sandboxFailoverEndpointsInApi[i].(map[string]interface{})
		assert.Equal(t, sandboxFailoverEndpointsInParams[i]["url"], singleSandboxFailoverEndpointInApi["url"])
		if isSaopEndpoint {
			assert.Equal(t, "address", productionEndpointInApi["endpoint_type"])
		}
	}

	same := "override_with_same_value"
	api.SetEndpointType(same)
	importedAPI.SetEndpointType(same)
	api.SetEndPointConfig(same)
	importedAPI.SetEndPointConfig(same)

	ValidateAPIsEqual(t, api, importedAPI)
}

func ValidateAwsEndpoint(t *testing.T, apiParams *Params, api, importedAPI *apim.API) {
	t.Helper()

	//Validate EndPoint Type
	assert.Equal(t, "awslambda", importedAPI.GetEndpointType())

	endPointConfigInApi := importedAPI.EndpointConfig
	//Validate access method
	accessMethodInParams := apiParams.Environments[0].Configs.AWSLambdaEndpoints.AccessMethod
	accessMethodInApi := endPointConfigInApi.(map[string]interface{})["access_method"].(string)
	accessMethodInApi = strings.ReplaceAll(accessMethodInApi, "-", "_")
	assert.Equal(t, accessMethodInParams, accessMethodInApi)

	//Validate stored mode configurations
	if strings.EqualFold(accessMethodInApi, "stored") {
		//Validate Amazon access key
		accessKeyInParams := apiParams.Environments[0].Configs.AWSLambdaEndpoints.AmznAccessKey
		accessKeyInApi := endPointConfigInApi.(map[string]interface{})["amznAccessKey"].(string)
		assert.Equal(t, accessKeyInParams, accessKeyInApi)

		//Validate Amazon access secret key
		accessSecretKeyInApi := endPointConfigInApi.(map[string]interface{})["amznSecretKey"].(string)
		assert.Equal(t, "AWS_SECRET_KEY", accessSecretKeyInApi)

		//Validate Amazon region
		regionInParams := apiParams.Environments[0].Configs.AWSLambdaEndpoints.AmznRegion
		regionInApi := endPointConfigInApi.(map[string]interface{})["amznRegion"].(string)
		assert.Equal(t, regionInParams, regionInApi)
	}

	same := "override_with_same_value"
	api.SetEndpointType(same)
	importedAPI.SetEndpointType(same)
	api.SetEndPointConfig(same)
	importedAPI.SetEndPointConfig(same)

	ValidateAPIsEqual(t, api, importedAPI)
}

func ValidateDynamicEndpoint(t *testing.T, apiParams *Params, api, importedAPI *apim.API) {
	t.Helper()

	//Validate EndPoint Type
	assert.Equal(t, "default", importedAPI.GetEndpointType())

	endPointConfigInApi := importedAPI.EndpointConfig

	//Validate default failover config
	failoverConfigEnableInApi := endPointConfigInApi.(map[string]interface{})["failover"].(string)
	failoverConfigEnableInApiBool, _ := strconv.ParseBool(failoverConfigEnableInApi)
	assert.Equal(t, false, failoverConfigEnableInApiBool)

	//Validate default url value for sandbox and production endpoint
	sandboxEndpointsInApi := endPointConfigInApi.(map[string]interface{})["sandbox_endpoints"].(map[string]interface{})
	productionEndpointsInApi := endPointConfigInApi.(map[string]interface{})["production_endpoints"].(map[string]interface{})
	assert.Equal(t, "default", sandboxEndpointsInApi["url"])
	assert.Equal(t, "default", productionEndpointsInApi["url"])

	same := "override_with_same_value"
	api.SetEndpointType(same)
	importedAPI.SetEndpointType(same)
	api.SetEndPointConfig(same)
	importedAPI.SetEndPointConfig(same)

	ValidateAPIsEqual(t, api, importedAPI)
}

// ValidateAPIsEqualWithEndpointConfigsFromParam : Validate endpoint configs from params and the equality of the two APIs
func ValidateAPIsEqualWithEndpointConfigsFromParam(t *testing.T, api *apim.API, importedAPI *apim.API, apiParams *Params) {
	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Production["url"], importedAPI.GetProductionURL())
	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Sandbox["url"], importedAPI.GetSandboxURL())

	apiCopy := apim.CopyAPI(api)
	importedAPICopy := apim.CopyAPI(importedAPI)

	same := "override_with_same_value"
	apiCopy.SetProductionURL(same)
	importedAPICopy.SetProductionURL(same)
	apiCopy.SetSandboxURL(same)
	importedAPICopy.SetSandboxURL(same)

	ValidateAPIsEqual(t, &apiCopy, &importedAPICopy)
}

func ValidateAPIImportUpdateWithParamsEndpointConfig(t *testing.T, args *ApiImportExportTestArgs, apiParams *Params,
	username string) {

	// Update the production and sandbox endpoints in the params file
	apiParams.Environments[0].Configs.Endpoints.Production["url"] = "https://prod-updated.wso2.com"
	apiParams.Environments[0].Configs.Endpoints.Sandbox["url"] = "https://sand-updated.wso2.com"

	apiData, err := yaml.Marshal(apiParams)
	if err != nil {
		t.Error(err)
	}

	// Write the temporary params file
	tempParams := EnvParamsFilesDir + "temp-" + base.GenerateRandomString() + username + ".yaml"
	err = ioutil.WriteFile(tempParams, apiData, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	args.ParamsFile = tempParams
	args.Update = true

	importedAPIUpdated := GetImportedAPI(t, args)

	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Production["url"], importedAPIUpdated.GetProductionURL())
	assert.Equal(t, apiParams.Environments[0].Configs.Endpoints.Sandbox["url"], importedAPIUpdated.GetSandboxURL())

	t.Cleanup(func() {
		base.RemoveDir(args.ParamsFile)
	})
}
