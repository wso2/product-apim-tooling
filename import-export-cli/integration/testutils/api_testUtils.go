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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/adminservices"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	yaml2 "gopkg.in/yaml.v2"
)

func GetAPIById(t *testing.T, client *apim.Client, username, password, apiId string) *apim.API {
	client.Login(username, password)
	return client.GetAPI(apiId)
}

func AddAPI(t *testing.T, client *apim.Client, username, password string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username, "", DevFirstDefaultAPIVersion, "")
	doClean := true
	id := client.AddAPI(t, api, username, password, doClean)
	api = client.GetAPI(id)
	return api
}

func AddCustomAPI(t *testing.T, client *apim.Client, username, password, name, version, context string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username, name, version, context)
	id := client.AddAPI(t, api, username, password, true)
	api = client.GetAPI(id)
	return api
}

func UpdateAPI(t *testing.T, client *apim.Client, api *apim.API, username, password string) *apim.API {
	client.Login(username, password)
	id := client.UpdateAPI(t, api, username, password)
	api = client.GetAPI(id)
	return api
}

func AddSoapAPI(t *testing.T, client *apim.Client, username, password, apiType string) *apim.API {
	path := "testdata/phoneverify.wsdl"
	client.Login(username, password)
	additionalProperties := client.GenerateAdditionalProperties(username, SoapEndpointURL, apiType, nil)
	id := client.AddSoapAPI(t, path, additionalProperties, username, password, apiType)
	api := client.GetAPI(id)
	return api
}

func AddGraphQLAPI(t *testing.T, client *apim.Client, username, password string) *apim.API {
	client.Login(username, password)
	path := "testdata/products-schema.graphql"
	validationResponse := client.ValidateGraphQLSchema(t, path, username, password)
	if validationResponse.IsValid {
		operations := validationResponse.GraphQLInfo.Operations
		additionalProperties := client.GenerateAdditionalProperties(username, GraphQLEndpoint, APITypeGraphQL, operations)
		id := client.AddGraphQLAPI(t, path, additionalProperties, username, password)
		api := client.GetAPI(id)
		return api
	} else {
		t.Error(t, validationResponse.ErrorMessage)
	}
	return nil
}

func AddWebStreamingAPI(t *testing.T, client *apim.Client, username, password, apiType string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleStreamingAPIData(username, apiType)
	doClean := true
	id := client.AddAPI(t, api, username, password, doClean)
	api = client.GetAPI(id)
	return api
}

func AddWebStreamingAPIFromAsyncAPIDefinition(t *testing.T, client *apim.Client, username, password, apiType string) *apim.API {
	client.Login(username, password)
	path := "testdata/streetlights.yml"
	additionalProperties := client.GenerateAdditionalProperties(username, WebSocketEndpoint, apiType, nil)
	id := client.AddStreamingAPI(t, path, additionalProperties, username, password)
	api := client.GetAPI(id)
	return api
}

func CreateAndDeployAPIRevision(t *testing.T, client *apim.Client, username, password, apiID string) string {
	client.Login(username, password)
	revision := client.CreateAPIRevision(apiID)
	client.DeployAPIRevision(t, apiID, "", "", revision.ID)
	base.WaitForIndexing()
	return revision.ID
}

func DeployAndPublishAPI(t *testing.T, client *apim.Client, username, password, apiID string) {
	CreateAndDeployAPIRevision(t, client, username, password, apiID)
	PublishAPI(client, username, password, apiID)
	base.WaitForIndexing()
}

func GetDeployedAPIRevisions(t *testing.T, client *apim.Client, username, password,
	apiID string) *apim.APIRevisionList {
	client.Login(username, password)
	revisionsList := client.GetAPIRevisions(apiID, "deployed:true")
	return revisionsList
}

func GetDeployedAPIProductRevisions(t *testing.T, client *apim.Client, username, password,
	apiProductID string) *apim.APIRevisionList {
	client.Login(username, password)
	revisionsList := client.GetAPIProductRevisions(apiProductID, "deployed:true")
	return revisionsList
}

func GetGatewayEnvironments(apiRevisions *apim.APIRevisionList) []string {
	var gatewayEnvironments []string
	for _, apiRevision := range apiRevisions.List {
		for _, deployment := range apiRevision.DeploymentInfo {
			gatewayEnvironments = append(gatewayEnvironments, deployment.Name)
		}
	}
	return gatewayEnvironments
}

func AddAPIWithoutCleaning(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username, "", DevFirstDefaultAPIVersion, "")
	doClean := false
	id := client.AddAPI(t, api, username, password, doClean)
	api = client.GetAPI(id)
	return api
}

func AddAPIToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.API, *apim.API) {
	client1.Login(username, password)
	api := client1.GenerateSampleAPIData(username, "", DevFirstDefaultAPIVersion, "")
	doClean := true
	id1 := client1.AddAPI(t, api, username, password, doClean)
	api1 := client1.GetAPI(id1)

	client2.Login(username, password)
	id2 := client2.AddAPI(t, api, username, password, doClean)
	api2 := client2.GetAPI(id2)

	return api1, api2
}

func AddAPIFromOpenAPIDefinition(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	client.Login(username, password)
	path := GetSwaggerPetstoreDefinition(t, username)
	additionalProperties := client.GenerateAdditionalProperties(username, RESTAPIEndpoint, APITypeREST, nil)
	id := client.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api := client.GetAPI(id)
	return api
}

func AddAPIFromOpenAPIDefinitionToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.API, *apim.API) {
	client1.Login(username, password)
	path := GetSwaggerPetstoreDefinition(t, username)
	additionalProperties := client1.GenerateAdditionalProperties(username, RESTAPIEndpoint, APITypeREST, nil)
	id1 := client1.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api1 := client1.GetAPI(id1)

	client2.Login(username, password)
	id2 := client2.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api2 := client2.GetAPI(id2)

	return api1, api2
}

func GenerateAdvertiseOnlyAPIDefinition(t *testing.T) (string, apim.API) {
	projectPath, _ := filepath.Abs(base.GenerateRandomString())
	base.CreateTempDir(t, projectPath)

	// Read the sample-api.yaml file in the testdata directory
	sampleContent := ReadAPIDefinition(t, SampleAPIYamlFilePath)

	// Inject advertise only API specfic parameters
	apim.GenerateAdvertiseOnlyProperties(&sampleContent.Data, "https://localhost:9443/devportal", "https://production-ep:9443",
		"https://sandbox-ep:9443")

	advertiseOnlyAPIDefinitionPath := filepath.Join(projectPath, filepath.FromSlash(utils.APIDefinitionFileYaml))

	// Write the API definition to the temp directory
	WriteToAPIDefinition(t, sampleContent, advertiseOnlyAPIDefinitionPath)
	return advertiseOnlyAPIDefinitionPath, sampleContent.Data
}

func CreateAndDeploySeriesOfAPIRevisions(t *testing.T, client *apim.Client, api *apim.API,
	apiCreator *Credentials, apiPublisher *Credentials) map[int]*apim.API {

	apiRevisions := make(map[int]*apim.API)

	revisionIds := make([]string, 0, 3)

	originalApiId := api.ID

	// Create and Deploy Revision 1 of the above API
	revisionIds = append(revisionIds, CreateAndDeployAPIRevision(t, client, apiPublisher.Username, apiPublisher.Password, originalApiId))

	api.Transport = []string{"https"}

	api = UpdateAPI(t, client, api, apiCreator.Username, apiCreator.Password)

	// Create and Deploy Revision 2 of the above API
	revisionIds = append(revisionIds, CreateAndDeployAPIRevision(t, client, apiPublisher.Username, apiPublisher.Password, originalApiId))

	api.AuthorizationHeader = "AuthorizationNew"

	UpdateAPI(t, client, api, apiCreator.Username, apiCreator.Password)

	// Create and Deploy Revision 3 of the above API
	revisionIds = append(revisionIds, CreateAndDeployAPIRevision(t, client, apiPublisher.Username, apiPublisher.Password, originalApiId))

	for _, rev := range revisionIds {
		api := client.GetAPI(rev)
		apiRevisions[api.RevisionID] = api
		t.Log("CreateAndDeploySeriesOfAPIRevisions api.RevisionID: ", api.RevisionID, ", api ID: ", api.ID)
	}

	return apiRevisions
}

func GetAPI(t *testing.T, client *apim.Client, name string, username string, password string) *apim.API {
	if username == adminservices.DevopsUsername {
		client.Login(adminservices.AdminUsername, adminservices.AdminPassword)
	} else if username == adminservices.DevopsUsername+"@"+adminservices.Tenant1 {
		client.Login(adminservices.AdminUsername+"@"+adminservices.Tenant1, adminservices.AdminPassword)
	} else {
		client.Login(username, password)
	}
	apiInfo, err := client.GetAPIByName(name)

	if err != nil {
		t.Fatal(err)
	}

	return client.GetAPI(apiInfo.ID)
}

func getAPIs(client *apim.Client, username string, password string) *apim.APIList {
	client.Login(username, password)
	return client.GetAPIs()
}

func deleteAPIByCtl(t *testing.T, args *ApiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "api", "-n", args.Api.Name, "-v", args.Api.Version, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func PublishAPI(client *apim.Client, username string, password string, apiID string) {
	base.WaitForIndexing()
	client.Login(username, password)
	client.PublishAPI(apiID)
}

func ChangeAPILifeCycle(client *apim.Client, username, password, apiID, action string) *apim.API {
	base.WaitForIndexing()
	client.Login(username, password)
	client.ChangeAPILifeCycle(apiID, action)
	api := client.GetAPI(apiID)
	return api
}

func UnsubscribeAPI(client *apim.Client, username string, password string, apiID string) {
	client.Login(username, password)
	client.DeleteSubscriptions(apiID)
}

func GetResourceURL(apim *apim.Client, api *apim.API) string {
	port := 8280 + apim.GetPortOffset()
	return "http://" + apim.GetHost() + ":" + strconv.Itoa(port) + api.Context + "/" + api.Version + "/menu"
}

func GetEnvAPIExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedApisDirName, envName)
}

func GetEnvAPIProductExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedApiProductsDirName, envName)
}

func exportAPI(t *testing.T, name, version, provider, env string) (string, error) {
	var output string
	var err error

	if provider == "" {
		output, err = base.Execute(t, "export", "api", "-n", name, "-v", version, "-e", env, "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "export", "api", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")
	}

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, GetEnvAPIExportPath(env), name, version)
	})

	return output, err
}

func exportAPIRevision(t *testing.T, args *ApiImportExportTestArgs) (string, error) {
	var output string
	var err error

	flags := []string{"export", "api", "-n", args.Api.Name, "-v", args.Api.Version, "-e", args.SrcAPIM.GetEnvName(), "-k", "--verbose"}

	if args.ApiProvider.Username != "" {
		flags = append(flags, "-r", args.ApiProvider.Username)
	}

	if args.IsLatest {
		flags = append(flags, "--latest")
	} else {
		flags = append(flags, "--rev", args.Revision)
	}

	output, err = base.Execute(t, flags...)

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()), args.Api.Name, args.Api.Version)
	})

	return output, err
}

func ValidateAllApisOfATenantIsExported(t *testing.T, args *ApiImportExportTestArgs, apisAdded int) {
	output, error := exportAllApisOfATenant(t, args)
	assert.Nil(t, error, "Error while exporting APIs")
	assert.Contains(t, output, "export-apis execution completed", "Error while exporting APIs")

	//Derive exported path from output
	exportedPath := base.GetExportedPathFromOutput(strings.ReplaceAll(output, "Command: export-apis execution completed !", ""))
	count, _ := base.CountFiles(t, exportedPath)
	assert.GreaterOrEqual(t, count, apisAdded, "Error while exporting APIs")

	t.Cleanup(func() {
		//Remove Exported apis and logout
		pathToCleanUp := utils.DefaultExportDirPath + TestMigrationDirectorySuffix
		base.RemoveDir(pathToCleanUp)
	})
}

func importAPI(t *testing.T, args *ApiImportExportTestArgs, doClean bool) (string, error) {
	var fileName string
	if args.ImportFilePath == "" {
		fileName = base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.Api.Name, args.Api.Version)
	} else {
		fileName = args.ImportFilePath
	}

	params := []string{"import", "api", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose"}

	if args.OverrideProvider {
		params = append(params, "--preserve-provider=false")
	}

	if args.ParamsFile != "" {
		params = append(params, "--params", args.ParamsFile)
	}

	if args.Update {
		params = append(params, "--update=true")
	}

	output, err := base.Execute(t, params...)

	if !args.Update && doClean {
		t.Cleanup(func() {
			if strings.EqualFold("DEPRECATED", args.Api.LifeCycleStatus) {
				args.CtlUser.Username, args.CtlUser.Password =
					apim.RetrieveAdminCredentialsInsteadCreator(args.CtlUser.Username, args.CtlUser.Password)
				args.DestAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)
			}
			err := args.DestAPIM.DeleteAPIByName(args.Api.Name)

			if err != nil {
				t.Fatal(err)
			}
			base.WaitForIndexing()
		})

	}
	return output, err
}

func importAPIPreserveProviderFailure(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import", "api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")
	return output, err
}

func listAPIs(t *testing.T, args *ApiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "get", "apis", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func listAPIsWithJsonArrayFormat(t *testing.T, args *ApiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "get", "apis", "-e", args.SrcAPIM.EnvName, "--format", "jsonArray",
		"-k", "--verbose")
	return output, err
}

func changeLifeCycleOfAPI(t *testing.T, args *ApiChangeLifeCycleStatusTestArgs) (string, error) {
	output, err := base.Execute(t, "change-status", "api", "-a", args.Action, "-n", args.Api.Name,
		"-v", args.Api.Version, "-e", args.APIM.EnvName, "-k", "--verbose")
	return output, err
}

func ValidateAPIExportFailure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting api from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsAPIArchiveExists(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version), "Test failed because the API was exported successfully")
}

func ValidateAPIExportFailureUnauthenticated(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting api from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	result, _ := exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())
	assert.Contains(t, result, "401", "Test failed because the response does not contains Unauthenticated request")

	// Validate that export failed
	assert.False(t, base.IsAPIArchiveExists(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version), "Test failed because the API was exported successfully")
}

func ValidateAPIExportImport(t *testing.T, args *ApiImportExportTestArgs, apiType string) *apim.API {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.Api.Provider, args.SrcAPIM.GetEnvName())

	validateAPIProject(t, args, apiType)

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	result, err := importAPI(t, args, true)
	assert.Nil(t, err, "Error while importing the API")
	assert.Contains(t, result, "Successfully imported API", "Error while importing the API")

	// Give time for newly imported API to get indexed, or else GetAPI by name will fail
	base.WaitForIndexing()

	// Get App from env 2
	importedAPI := GetAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	// Validate env 1 and env 2 API is equal
	ValidateAPIsEqual(t, args.Api, importedAPI)

	return importedAPI
}

func ValidateAPIImportExportForAdvertiseOnlyAPI(t *testing.T, args *ApiImportExportTestArgs, apiType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.Api.Provider, args.SrcAPIM.GetEnvName())

	validateAPIProject(t, args, apiType)

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	result, err := importAPI(t, args, true)
	assert.Nil(t, err, "Error while importing the API")
	assert.Contains(t, result, "Successfully imported API", "Error while importing the API")

	// Give time for newly imported API to get indexed, or else GetAPI by name will fail
	base.WaitForIndexing()

	// Get App from env 2
	importedAPI := GetAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	validateAdvertiseOnlyAPIsEqual(t, importedAPI, args)
}

func ValidateAPIRevisionExportImport(t *testing.T, args *ApiImportExportTestArgs, apiType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPIRevision(t, args)

	validateAPIProject(t, args, apiType)

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAPI(t, args, true)

	// Give time for newly imported API to get indexed, or else GetAPI by name will fail
	base.WaitForIndexing()

	// Get App from env 2
	importedAPI := GetAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	// Validate env 1 and env 2 API is equal
	ValidateImportedAPIsEqualToRevision(t, args.Api, importedAPI)
}

func validateAPIProject(t *testing.T, args *ApiImportExportTestArgs, apiType string) {
	assert.True(t, base.IsAPIArchiveExists(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))

	if strings.EqualFold(apiType, APITypeREST) {
		assert.True(t, base.IsFileExistsInAPIArchive(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
			utils.InitProjectDefinitionsSwagger, args.Api.Name, args.Api.Version))
	}

	if strings.EqualFold(apiType, APITypeSoap) {
		wsdlFilePathInProject := utils.InitProjectWSDL + string(os.PathSeparator) + args.Api.Name + "-" + args.Api.Version + ".wsdl"
		assert.True(t, base.IsFileExistsInAPIArchive(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()), wsdlFilePathInProject,
			args.Api.Name, args.Api.Version))
	}

	if strings.EqualFold(apiType, APITypeGraphQL) {
		assert.True(t, base.IsFileExistsInAPIArchive(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
			utils.InitProjectDefinitionsGraphQLSchema, args.Api.Name, args.Api.Version))
	}

	if strings.EqualFold(apiType, APITypeWebScoket) || strings.EqualFold(apiType, APITypeWebSub) ||
		strings.EqualFold(apiType, APITypeSSE) || strings.EqualFold(apiType, APITypeAsync) {
		assert.True(t, base.IsFileExistsInAPIArchive(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
			utils.InitProjectDefinitionsAsyncAPI, args.Api.Name, args.Api.Version))
	}
}

func ValidateAPIExport(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))
}

func ValidateExportedAPIStructure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, _ := exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	validateAPI(t, args.Api, output, args.IsDeployed, SampleAPIYamlFilePath)

	assert.True(t, base.IsAPIArchiveExists(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))
}

func ValidateExportedAPIRevisionStructure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, _ := exportAPIRevision(t, args)

	validateAPI(t, args.Api, output, args.IsDeployed, SampleRevisionedAPIYamlFilePath)

	assert.True(t, base.IsAPIArchiveExists(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))
}

func ValidateExportedAPIRevisionFailure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, _ := exportAPIRevision(t, args)

	assert.Contains(t, output, "404", "Test failed because the response does not contains a not found request")

	// Validate that export failed
	assert.False(t, base.IsAPIArchiveExists(t, GetEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version), "Test failed because the API Revision was exported successfully")
}

func validateAPI(t *testing.T, api *apim.API, exportedOutput string, isDeployed bool, sampleFile string) {
	// Unzip exported API
	exportedPath := base.GetExportedPathFromOutput(exportedOutput)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	unzipedProjectPath := relativePath + string(os.PathSeparator) + api.Name + "-" + api.Version

	// Read the api.yaml file in the exported directory
	fileData, err := ioutil.ReadFile(filepath.Join(unzipedProjectPath, APIYamlFilePath))

	if err != nil {
		t.Error(err)
	}

	validateAPIStructure(t, &fileData, sampleFile)

	if isDeployed {
		assert.True(t, base.IsFileAvailable(t, filepath.Join(unzipedProjectPath, DeploymentEnvYamlFilePath)), "Expected deployment_environments.yaml not found")
	} else {
		assert.False(t, base.IsFileAvailable(t, filepath.Join(unzipedProjectPath, DeploymentEnvYamlFilePath)), "Non required deployment_environments.yaml found")
	}

	t.Cleanup(func() {
		// Remove the extracted project
		base.RemoveDir(exportedPath)
		base.RemoveDir(relativePath)
	})
}

func validateAPIStructure(t *testing.T, fileData *[]byte, sampleFile string) {
	// Extract the "data" field to an interface
	fileContent := make(map[string]interface{})
	err := yaml.Unmarshal(*fileData, &fileContent)
	if err != nil {
		t.Error(err)
	}
	apiData := fileContent["data"].(map[interface{}]interface{})

	// Read the sample-api.yaml file in the testdata directory
	sampleData, err := ioutil.ReadFile(sampleFile)
	if err != nil {
		t.Error(err)
	}

	// Extract the "data" field to an interface
	sampleDataContent := make(map[string]interface{})
	err = yaml.Unmarshal(sampleData, &sampleDataContent)
	if err != nil {
		t.Error(err)
	}
	sampleAPIData := sampleDataContent["data"].(map[interface{}]interface{})

	// Compare the artifact versions
	exportedAPIArtifactVersion := fileContent["version"].(string)
	sampleAPIArtifactVersion := sampleDataContent["version"].(string)
	assert.Equal(t, exportedAPIArtifactVersion, sampleAPIArtifactVersion,
		"Exported artifact version: "+exportedAPIArtifactVersion+
			" does not matches with the sample artifact version: "+sampleAPIArtifactVersion)

	// Check whether the fields of the API DTO structure in APICTL has all the fields in API DTO structure from APIM
	base.Log("\n-----------------------------------------------------------------------------------------")
	base.Log("Checking whether the fields of APICTL API DTO struct has all the fields from APIM API DTO struct")
	for key, _ := range apiData {
		keyValue := key.(string)
		_, ok := sampleAPIData[key]
		base.Log("\"" + keyValue + "\" is in both the structures")
		if !ok {
			t.Error("Missing \"" + keyValue + "\" in the API DTO structure from APICTL")
		}
	}

	// Check whether the fields of the API DTO structure in APIM has all the fields in API DTO structure from APICTL
	base.Log("\n-----------------------------------------------------------------------------------------")
	base.Log("Checking whether the fields of APIM API DTO struct has all the fields from APICTL API DTO struct")
	for key, _ := range sampleAPIData {
		keyValue := key.(string)
		_, ok := apiData[key]
		base.Log("\"" + keyValue + "\" is in both the structures")
		if !ok {
			t.Error("Missing \"" + keyValue + "\" in the API DTO structure from APIM")
		}
	}
}

func GetImportedAPI(t *testing.T, args *ApiImportExportTestArgs) *apim.API {
	t.Helper()

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())
	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	_, err := importAPI(t, args, true)

	if err != nil {
		t.Fatal(err)
	}

	// Give time for newly imported API to get indexed, or else GetAPI by name will fail
	base.WaitForIndexing()

	// Get App from env 2
	importedAPI := GetAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	return importedAPI
}

func ReadParams(t *testing.T, apiParamsPath string) *Params {
	reader, err := os.Open(apiParamsPath)

	if err != nil {
		base.Fatal(err)
	}
	defer reader.Close()

	apiParams := Params{}
	yaml.NewDecoder(reader).Decode(&apiParams)

	return &apiParams
}

func ReadAPIDefinition(t *testing.T, path string) apim.APIFile {

	// Read the file in the path
	sampleData, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	// Extract the content to a structure
	sampleContent := apim.APIFile{}
	err = yaml.Unmarshal(sampleData, &sampleContent)
	if err != nil {
		t.Error(err)
	}

	return sampleContent
}

func WriteToAPIDefinition(t *testing.T, content apim.APIFile, path string) {
	apiData, err := yaml2.Marshal(content)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile(path, apiData, os.ModePerm)
	if err != nil {
		t.Error(err)
	}
}

func ValidateAPIImport(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAPI(t, args, true)

	// Give time for newly imported API to get indexed, or else GetAPI by name will fail
	base.WaitForIndexing()

	// Get App from env 2
	importedAPI := GetAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	// Validate env 1 and env 2 API is equal
	validateAPIsEqualCrossTenant(t, args.Api, importedAPI)
}

func ValidateAPIImportForMultipleVersions(t *testing.T, args *ApiImportExportTestArgs, firstImportedAPIId string) *apim.API {

	t.Helper()

	isFirstImport := false
	if strings.EqualFold(firstImportedAPIId, "") {
		isFirstImport = true
	}

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAPI(t, args, isFirstImport)

	// Give time for newly imported API to get indexed, or else getAPI by name will fail
	base.WaitForIndexing()

	if !isFirstImport {
		args.DestAPIM.DeleteAPI(firstImportedAPIId)
		base.WaitForIndexing()
	}

	// Get App from env 2
	importedAPI := GetAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	// Validate env 1 and env 2 API is equal
	validateAPIsEqualCrossTenant(t, args.Api, importedAPI)

	return importedAPI
}

func ValidateAPIImportFailure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// importAPIPreserveProviderFailure is used to eleminate cleaning the API after importing
	result, err := importAPIPreserveProviderFailure(t, args.SrcAPIM.GetEnvName(), args.Api, args.DestAPIM)

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
}

// ValidateAPIsEqual : Validate if two APIs are equal while ignoring unique fields
func ValidateAPIsEqual(t *testing.T, api1 *apim.API, api2 *apim.API) {
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

	// If an API is not advertise only, the API owner will be changed during export and import to the current provider
	if (api1Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		api1Copy.AdvertiseInformation.ApiOwner = same
	}
	if (api2Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		api2Copy.AdvertiseInformation.ApiOwner = same
	}

	// Sort member collections to make equality check possible
	apim.SortAPIMembers(&api1Copy)
	apim.SortAPIMembers(&api2Copy)

	assert.Equal(t, api1Copy, api2Copy, "API obejcts are not equal")
}

func validateAdvertiseOnlyAPIsEqual(t *testing.T, importedAPI *apim.API, args *ApiImportExportTestArgs) {
	t.Helper()

	assert.Equal(t, args.Api.AdvertiseInformation.Advertised, importedAPI.AdvertiseInformation.Advertised)
	assert.Equal(t, args.Api.Provider, importedAPI.AdvertiseInformation.ApiOwner)
	assert.Equal(t, args.Api.AdvertiseInformation.ApiExternalProductionEndpoint, importedAPI.AdvertiseInformation.ApiExternalProductionEndpoint)
	assert.Equal(t, args.Api.AdvertiseInformation.ApiExternalSandboxEndpoint, importedAPI.AdvertiseInformation.ApiExternalSandboxEndpoint)

	if (args.CtlUser.Username == adminservices.AdminUsername) ||
		(args.CtlUser.Username == adminservices.AdminUsername+"@"+adminservices.Tenant1) {
		// Only the users who has admin privileges (apim:admin scope) were allowed to set the original devportal URL.
		assert.Equal(t, args.Api.AdvertiseInformation.OriginalDevPortalUrl,
			importedAPI.AdvertiseInformation.OriginalDevPortalUrl)
	} else {
		assert.Equal(t, "", importedAPI.AdvertiseInformation.OriginalDevPortalUrl)
	}

	// Certificates should not get exported for advertise only APIs
	validateNonExportedAPICerts(t, importedAPI, args)
}

// ValidateImportedAPIsEqualToRevision : Validate if the imported API and exported revision is the same by ignoring
// the unique details and revision specific details.
func ValidateImportedAPIsEqualToRevision(t *testing.T, api1 *apim.API, api2 *apim.API) {
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

	// When imported the revision as API, the "IsRevision" property will be false for the imported API. Hence,
	//the property of the imported API should be changed
	api2Copy.IsRevision = true

	// When imported revision as API, the "RevisionID" property will be 0 for the imported API. Hence, the property of
	// the imported API be changed
	api2Copy.RevisionID = 1

	// If an API is not advertise only, the API owner will be changed during export and import to the current provider
	if (api1Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		api1Copy.AdvertiseInformation.ApiOwner = same
	}
	if (api2Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		api2Copy.AdvertiseInformation.ApiOwner = same
	}

	// Sort member collections to make equality check possible
	apim.SortAPIMembers(&api1Copy)
	apim.SortAPIMembers(&api2Copy)

	assert.Equal(t, api1Copy, api2Copy, "API obejcts are not equal")
}

func ValidateAPIsList(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List APIs of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIs(t, args)

	apisList := args.SrcAPIM.GetAPIs()

	ValidateListAPIsEqual(t, output, apisList)
}

func ValidateListAPIsEqual(t *testing.T, apisListFromCtl string, apisList *apim.APIList) {
	unmatchedCount := apisList.Count
	for _, api := range apisList.List {
		// If the output string contains the same API ID, then decrement the count
		assert.Truef(t, strings.Contains(apisListFromCtl, api.ID), "apisListFromCtl: "+apisListFromCtl+
			" , does not contain api.ID: "+api.ID)
		unmatchedCount--
	}

	// Count == 0 means that all the APIs from apisList were in apisListFromCtl
	assert.Equal(t, 0, unmatchedCount, "API lists are not equal")
}

// ValidateAPIsListWithJsonArrayFormat : Validate the received list of APIs are in JsonArray format and verify only
// the required ones are there and others are not in the command line output
func ValidateAPIsListWithJsonArrayFormat(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List APIs of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIsWithJsonArrayFormat(t, args)

	apisList := args.SrcAPIM.GetAPIs()

	// Validate APIs list with added APIs
	ValidateListAPIsEqual(t, output, apisList)

	// Validate JsonArray format
	assert.Contains(t, output, "[\n {\n", "Error while listing APIs in JsonArray format")

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

	// If an API is not advertise only, the API owner will be changed during export and import to the current provider
	if (api1Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		api1Copy.AdvertiseInformation.ApiOwner = same
	}
	if (api2Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		api2Copy.AdvertiseInformation.ApiOwner = same
	}

	// Sort member collections to make equality check possible
	apim.SortAPIMembers(&api1Copy)
	apim.SortAPIMembers(&api2Copy)

	assert.Equal(t, api1Copy, api2Copy, "API obejcts are not equal")
}

func ValidateAPIDelete(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Delete an API of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()
	apisListBeforeDelete := args.SrcAPIM.GetAPIs()

	deleteAPIByCtl(t, args)

	apisListAfterDelete := args.SrcAPIM.GetAPIs()
	base.WaitForIndexing()

	// Validate whether the expected number of API count is there
	assert.Equal(t, apisListBeforeDelete.Count, apisListAfterDelete.Count+1, "Expected number of APIs not deleted")

	// Validate that the delete is a success
	validateAPIIsDeleted(t, args.Api, apisListAfterDelete)
}

func ValidateAPIDeleteFailure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	apisListBeforeDelete := args.SrcAPIM.GetAPIs()

	output, _ := deleteAPIByCtl(t, args)

	apisListAfterDelete := args.SrcAPIM.GetAPIs()
	base.WaitForIndexing()

	// Validate whether the expected number of API count is there
	assert.NotContains(t, output, " API deleted successfully!. Status: 200", "Api delete is success with active subscriptions")
	assert.NotEqual(t, apisListBeforeDelete.Count, apisListAfterDelete.Count+1, "Expected number of APIs not deleted")

	t.Cleanup(func() {
		UnsubscribeAPI(args.SrcAPIM, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)
	})
}

func exportApiImportedFromProject(t *testing.T, APIName string, APIVersion string, EnvName string) (string, error) {
	return base.Execute(t, "export", "api", "-n", APIName, "-v", APIVersion, "-e", EnvName)
}

func exportAllApisOfATenant(t *testing.T, args *ApiImportExportTestArgs) (string, error) {
	//Setup environment
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	//Login to the environmeTestImportAndExportAPIWithJpegImagent
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, error := base.Execute(t, "export", "apis", "-e", args.SrcAPIM.GetEnvName(), "-k", "--force")
	return output, error
}

func validateAPIIsDeleted(t *testing.T, api *apim.API, apisListAfterDelete *apim.APIList) {
	for _, existingAPI := range apisListAfterDelete.List {
		assert.NotEqual(t, existingAPI.ID, api.ID, "API delete is not successful")
	}
}

func ValidateChangeLifeCycleStatusOfAPI(t *testing.T, args *ApiChangeLifeCycleStatusTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	//Execute apictl command to change life cycle of an Api
	output, _ := changeLifeCycleOfAPI(t, args)
	//Assert apictl output
	assert.Contains(t, output, "state changed successfully!", "Error while changing life cycle of API")

	base.WaitForIndexing()
	//Assert life cycle state after change
	api := GetAPI(t, args.APIM, args.Api.Name, args.CtlUser.Username, args.CtlUser.Password)
	assert.Equal(t, args.ExpectedState, api.LifeCycleStatus, "Expected Life cycle state change is not equals to actual status")
}

func ValidateChangeLifeCycleStatusOfAPIFailure(t *testing.T, args *ApiChangeLifeCycleStatusTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	//Execute apictl command to change life cycle of an Api
	output, _ := changeLifeCycleOfAPI(t, args)
	//Assert apictl output
	assert.NotContains(t, output, "state changed successfully!", "Error while changing life cycle of API")
	assert.NotEqual(t, args.Api.LifeCycleStatus, args.ExpectedState, "Life Cycle State changed successfully")
}

func ValidateApisListWithVersions(t *testing.T, args *InitTestArgs, newVersion string) {
	t.Helper()

	apis := getAPIs(args.SrcAPIM, args.CtlUser.Username, args.CtlUser.Password)

	isV2ApiExsits := false
	isV1ApiExsits := false

	// Validate required Apis in Apis List
	for _, api := range apis.List {
		if strings.EqualFold(api.Version, "1.0.0") && strings.EqualFold(args.APIName, api.Name) {
			isV1ApiExsits = true
		}
		if strings.EqualFold(api.Version, newVersion) && strings.EqualFold(args.APIName, api.Name) {
			isV2ApiExsits = true
		}
	}
	assert.Equal(t, true, isV1ApiExsits && isV2ApiExsits)
}

// Execute get apis command with query parameters
func searchAPIsWithQuery(t *testing.T, args *ApiImportExportTestArgs, query string) (string, error) {
	output, err := base.Execute(t, "get", "apis", "-e", args.SrcAPIM.EnvName, "--query", query, "-k", "--verbose")
	return output, err
}

// ValidateSearchApisList : Validate the received list of APIs and verify only the required ones are there and others
// are not in the command line output
func ValidateSearchApisList(t *testing.T, args *ApiImportExportTestArgs, searchQuery, matchQuery, unmatchedQuery string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := searchAPIsWithQuery(t, args, searchQuery)

	// Assert the match query is in the output
	assert.Truef(t, strings.Contains(output, matchQuery), "apisListFromCtl: "+output+
		" , does not contain the query: "+matchQuery)
	// Assert the unmatched query is not in the output
	assert.False(t, strings.Contains(output, unmatchedQuery), "apisListFromCtl: "+output+
		" , contains the query: "+unmatchedQuery)
}
