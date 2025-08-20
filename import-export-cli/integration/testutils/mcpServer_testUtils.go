/*
*  Copyright (c) 2025 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 LLC. licenses this file to you under the Apache License,
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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/adminservices"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
)

// // MCPServerImportExportTestArgs holds the arguments needed for MCP Server import/export tests
// type MCPServerImportExportTestArgs struct {
// 	MCPServerProvider Credentials
// 	CtlUser           Credentials
// 	MCPServer         *apim.MCPServer
// 	SrcAPIM           *apim.Client
// 	DestAPIM          *apim.Client
// 	IsDeployed        bool
// 	ImportFilePath    string
// 	UpdateFlag        bool
// 	PreserveProvider  bool
// }

// // MCPServerChangeLifeCycleStatusTestArgs holds the arguments needed for MCP Server lifecycle change tests
// type MCPServerChangeLifeCycleStatusTestArgs struct {
// 	CtlUser       Credentials
// 	MCPServer     *apim.MCPServer
// 	APIM          *apim.Client
// 	Action        string
// 	ExpectedState string
// }

// const (
// 	DevFirstDefaultMCPServerVersion = "1.0.0"
// 	SampleMCPServerYamlFilePath     = "testdata/TestArtifactDirectory/MCPServerSampleArtifacts/mcp_server.yaml"
// 	// Additional constants for MCP Server testing
// 	MCPServerYamlFilePath = "mcp_server.yaml"
// )

// GetMCPServerById returns MCP Server by ID
func GetMCPServerById(t *testing.T, client *apim.Client, username, password, mcpServerId string) *apim.MCPServer {
	client.Login(username, password)
	return client.GetMCPServer(mcpServerId)
}

// AddMCPServer adds a new MCP Server to APIM
func AddMCPServer(t *testing.T, client *apim.Client, username, password string) *apim.MCPServer {
	client.Login(username, password)
	mcpServer := client.GenerateSampleMCPServerData(username, "", DevFirstDefaultMCPServerVersion, "")
	id := client.AddMCPServer(t, mcpServer, username, password, true)
	mcpServer = client.GetMCPServer(id)
	return mcpServer
}

// AddCustomMCPServer adds a custom MCP Server to APIM
func AddCustomMCPServer(t *testing.T, client *apim.Client, username, password, name, version, context string) *apim.MCPServer {
	client.Login(username, password)
	mcpServer := client.GenerateSampleMCPServerData(username, name, version, context)
	id := client.AddMCPServer(t, mcpServer, username, password, true)
	mcpServer = client.GetMCPServer(id)
	return mcpServer
}

// UpdateMCPServer updates an existing MCP Server
func UpdateMCPServer(t *testing.T, client *apim.Client, mcpServer *apim.MCPServer, username, password string) *apim.MCPServer {
	client.Login(username, password)
	id := client.UpdateMCPServer(t, mcpServer, username, password)
	mcpServer = client.GetMCPServer(id)
	return mcpServer
}

func CreateAndDeployMCPServerRevision(t *testing.T, client *apim.Client, username, password, mcpServerID string) string {
	client.Login(username, password)
	revision := client.CreateMCPServerRevision(mcpServerID)
	client.DeployMCPServerRevision(t, mcpServerID, "", "", revision.ID)
	base.WaitForIndexing()
	return revision.ID
}

func DeployAndPublishMCPServer(t *testing.T, client *apim.Client, username, password, mcpServerID string) {
	CreateAndDeployMCPServerRevision(t, client, username, password, mcpServerID)
	PublishMCPServer(client, username, password, mcpServerID)
	base.WaitForIndexing()
}

func GetDeployedMCPServerRevisions(t *testing.T, client *apim.Client, username, password, mcpServerID string) *apim.MCPServerRevisionList {
	client.Login(username, password)
	revisionsList := client.GetMCPServerRevisions(mcpServerID, "deployed:true")
	return revisionsList
}

// AddMCPServerWithoutCleaning adds MCP Server without cleanup
func AddMCPServerWithoutCleaning(t *testing.T, client *apim.Client, username string, password string) *apim.MCPServer {
	client.Login(username, password)
	mcpServer := client.GenerateSampleMCPServerData(username, "", DevFirstDefaultMCPServerVersion, "")
	doClean := false
	id := client.AddMCPServer(t, mcpServer, username, password, doClean)
	mcpServer = client.GetMCPServer(id)
	return mcpServer
}

// AddMCPServerToTwoEnvs adds the same MCP Server to two environments
func AddMCPServerToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.MCPServer, *apim.MCPServer) {

	client1.Login(username, password)
	mcpServer1 := client1.GenerateSampleMCPServerData(username, "", DevFirstDefaultMCPServerVersion, "")
	id1 := client1.AddMCPServer(t, mcpServer1, username, password, true)
	mcpServer1 = client1.GetMCPServer(id1)

	client2.Login(username, password)
	mcpServer2 := client2.GenerateSampleMCPServerData(username, mcpServer1.Name, mcpServer1.Version, mcpServer1.Context)
	id2 := client2.AddMCPServer(t, mcpServer2, username, password, true)
	mcpServer2 = client2.GetMCPServer(id2)

	return mcpServer1, mcpServer2
}

func CreateAndDeploySeriesOfMCPServerRevisions(t *testing.T, client *apim.Client, mcpServer *apim.MCPServer,
	apiCreator *Credentials, apiPublisher *Credentials) map[int]*apim.MCPServer {

	mcpServerRevisions := make(map[int]*apim.MCPServer)

	revisionIds := make([]string, 0, 3)

	originalMcpServerId := mcpServer.ID

	// Create and Deploy Revision 1 of the above MCP Server
	revisionIds = append(revisionIds, CreateAndDeployMCPServerRevision(t, client, apiPublisher.Username, apiPublisher.Password, originalMcpServerId))

	mcpServer.Transport = []string{"https"}

	mcpServer = UpdateMCPServer(t, client, mcpServer, apiCreator.Username, apiCreator.Password)

	// Create and Deploy Revision 2 of the above MCP Server
	revisionIds = append(revisionIds, CreateAndDeployMCPServerRevision(t, client, apiPublisher.Username, apiPublisher.Password, originalMcpServerId))

	mcpServer.AuthorizationHeader = "AuthorizationNew"

	UpdateMCPServer(t, client, mcpServer, apiCreator.Username, apiCreator.Password)

	// Create and Deploy Revision 3 of the above MCP Server
	revisionIds = append(revisionIds, CreateAndDeployMCPServerRevision(t, client, apiPublisher.Username, apiPublisher.Password, originalMcpServerId))

	for _, rev := range revisionIds {
		mcpServer := client.GetMCPServer(rev)
		mcpServerRevisions[mcpServer.RevisionID] = mcpServer
		t.Log("CreateAndDeploySeriesOfMCPServerRevisions mcpServer.RevisionID: ", mcpServer.RevisionID, ", mcpServer ID: ", mcpServer.ID)
	}

	return mcpServerRevisions
}

func GetMCPServer(t *testing.T, client *apim.Client, name string, username string, password string) *apim.MCPServer {
	if username == adminservices.DevopsUsername {
		client.Login(adminservices.AdminUsername, adminservices.AdminPassword)
	} else if username == adminservices.DevopsUsername+"@"+adminservices.Tenant1 {
		client.Login(adminservices.AdminUsername+"@"+adminservices.Tenant1, adminservices.AdminPassword)
	} else {
		client.Login(username, password)
	}
	mcpServerInfo, err := client.GetMCPServerByName(name)
	if err != nil {
		t.Fatal(err)
	}
	return client.GetMCPServer(mcpServerInfo.ID)
}

func getMCPServers(client *apim.Client, username string, password string) *apim.MCPServerList {
	client.Login(username, password)
	return client.GetMCPServers()
}

func deleteMCPServerByCtl(t *testing.T, args *MCPServerImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "mcp-server", "-n", args.MCPServer.Name, "-v", args.MCPServer.Version, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func PublishMCPServer(client *apim.Client, username string, password string, mcpServerID string) {
	base.WaitForIndexing()
	client.Login(username, password)
	client.PublishMCPServer(mcpServerID)
}

func ChangeMCPServerLifeCycle(client *apim.Client, username, password, mcpServerID, action string) *apim.MCPServer {
	base.WaitForIndexing()
	client.Login(username, password)
	client.ChangeMCPServerLifeCycle(mcpServerID, action)
	mcpServer := client.GetMCPServer(mcpServerID)
	return mcpServer
}

func UnsubscribeMCPServer(client *apim.Client, username string, password string, mcpServerID string) {
	client.Login(username, password)
	client.DeleteSubscriptions(mcpServerID)
}

func GetEnvMCPServerExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedMCPServersDirName, envName)
}

// exportMCPServer exports a MCP Server from the environment
func exportMCPServer(t *testing.T, name, version, provider, env string) (string, error) {
	var output string
	var err error

	if provider == "" {
		output, err = base.Execute(t, "export", "mcp-server", "-n", name, "-v", version, "-e", env, "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "export", "mcp-server", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")
	}

	t.Cleanup(func() {
		base.RemoveMCPServerArchive(t, GetEnvMCPServerExportPath(env), name, version)
	})

	return output, err
}

func exportMCPServerRevision(t *testing.T, args *MCPServerImportExportTestArgs) (string, error) {
	var output string
	var err error

	flags := []string{"export", "mcp-server", "-n", args.MCPServer.Name, "-v", args.MCPServer.Version, "-e", args.SrcAPIM.GetEnvName(), "-k", "--verbose"}

	if args.MCPServerProvider.Username != "" {
		flags = append(flags, "-r", args.MCPServerProvider.Username)
	}

	if args.IsLatest {
		flags = append(flags, "--latest")
	} else {
		flags = append(flags, "--rev", args.Revision)
	}

	output, err = base.Execute(t, flags...)

	t.Cleanup(func() {
		base.RemoveMCPServerArchive(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()), args.MCPServer.Name, args.MCPServer.Version)
	})

	return output, err
}

func ValidateAllMCPServersOfATenantIsExported(t *testing.T, args *MCPServerImportExportTestArgs, mcpServersAdded int) {
	output, error := exportAllMCPServersOfATenant(t, args)
	assert.Nil(t, error, "Error while exporting MCP Servers")
	assert.Contains(t, output, "export-mcp-servers execution completed", "Error while exporting MCP Servers")

	//Derive exported path from output
	exportedPath := base.GetExportedPathFromOutput(strings.ReplaceAll(output, "Command: export-mcp-servers execution completed !", ""))
	count, _ := base.CountFiles(t, exportedPath)
	assert.GreaterOrEqual(t, count, mcpServersAdded, "Error while exporting MCP Servers")

	t.Cleanup(func() {
		//Remove Exported MCP Servers and logout
		pathToCleanUp := utils.DefaultExportDirPath + TestMigrationDirectorySuffix
		base.RemoveDir(pathToCleanUp)
	})
}

func importMCPServer(t *testing.T, args *MCPServerImportExportTestArgs, doClean bool) (string, error) {
	var fileName string
	if args.ImportFilePath == "" {
		fileName = base.GetAPIArchiveFilePath(t, args.SrcAPIM.GetEnvName(), args.MCPServer.Name, args.MCPServer.Version)
	} else {
		fileName = args.ImportFilePath
	}

	params := []string{"import", "mcp-server", "-f", fileName, "-e", args.DestAPIM.EnvName, "-k", "--verbose"}

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
			if strings.EqualFold("DEPRECATED", args.MCPServer.LifeCycleStatus) {
				args.CtlUser.Username, args.CtlUser.Password =
					apim.RetrieveAdminCredentialsInsteadCreator(args.CtlUser.Username, args.CtlUser.Password)
				args.DestAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)
			}
			err := args.DestAPIM.DeleteMCPServerByName(args.MCPServer.Name)

			if err != nil {
				t.Fatal(err)
			}
			base.WaitForIndexing()
		})

	}
	return output, err
}

func importMCPServerPreserveProviderFailure(t *testing.T, sourceEnv string, mcpServer *apim.MCPServer, client *apim.Client) (string, error) {
	fileName := base.GetMCPServerArchiveFilePath(t, sourceEnv, mcpServer.Name, mcpServer.Version)
	output, err := base.Execute(t, "import", "mcp-server", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")
	return output, err
}

func listMCPServers(t *testing.T, args *MCPServerImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "get", "mcp-servers", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func listMCPServersWithJsonArrayFormat(t *testing.T, args *MCPServerImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "get", "mcp-servers", "-e", args.SrcAPIM.EnvName, "--format", "jsonArray", "-k", "--verbose")
	return output, err
}

func changeMCPServerLifeCycle(t *testing.T, args *MCPServerChangeLifeCycleStatusTestArgs) (string, error) {
	output, err := base.Execute(t, "change-status", "mcp-server", "-a", args.Action, "-n", args.MCPServer.Name,
		"-v", args.MCPServer.Version, "-e", args.APIM.EnvName, "-k", "--verbose")
	return output, err
}

func GetMCPServerResourceURL(apim *apim.Client, mcpServer *apim.MCPServer) string {
	return apim.GetApimURL() + "/" + mcpServer.Context + "/" + mcpServer.Version
}

func ValidateMCPServerExportFailure(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting MCP Server from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportMCPServer(t, args.MCPServer.Name, args.MCPServer.Version, args.MCPServerProvider.Username, args.SrcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version), "Test failed because the MCP Server was exported successfully")
}

// ValidateMCPServerRevisionExportFailure validates that MCP Server revision export fails
func ValidateMCPServerRevisionExportFailure(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting MCP Server revision from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportMCPServerRevision(t, args)

	// Validate that export failed
	assert.False(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version), "Test failed because the MCP Server revision was exported successfully")
}

// ValidateMCPServerExportFailureUnauthenticated validates that MCP Server export fails for unauthenticated user
func ValidateMCPServerExportFailureUnauthenticated(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting MCP Server from env without login
	result, _ := exportMCPServer(t, args.MCPServer.Name, args.MCPServer.Version, args.MCPServerProvider.Username, args.SrcAPIM.GetEnvName())
	assert.Contains(t, result, "401", "Test failed because the response does not contain Unauthenticated request")

	// Validate that export failed
	assert.False(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version), "Test failed because the MCP Server was exported successfully")
}

// ValidateMCPServerExportImport validates MCP Server export and import
func ValidateMCPServerExportImport(t *testing.T, args *MCPServerImportExportTestArgs) *apim.MCPServer {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export MCP Server from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportMCPServer(t, args.MCPServer.Name, args.MCPServer.Version, args.MCPServerProvider.Username, args.SrcAPIM.GetEnvName())

	validateMCPServerProject(t, args)

	assert.True(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version), "Exported MCP Server archive does not exist")

	// Import MCP Server to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	result, err := importMCPServer(t, args, true)
	assert.Nil(t, err, "Error while importing the MCP Server")
	assert.Contains(t, result, "Successfully imported MCP Server", "Error while importing the MCP Server")

	base.WaitForIndexing()

	// Get MCP Server from env 2
	importedMCPServer := GetMCPServer(t, args.DestAPIM, args.MCPServer.Name, args.MCPServerProvider.Username, args.MCPServerProvider.Password)

	// Validate env 1 and env 2 MCP Server is equal
	ValidateMCPServersEqual(t, args.MCPServer, importedMCPServer)

	return importedMCPServer
}

func ValidateMCPServerRevisionExportImport(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export MCP Server revision from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportMCPServerRevision(t, args)

	// Import MCP Server revision to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importMCPServer(t, args, true)

	// Give time for newly imported MCP Server to get indexed
	base.WaitForIndexing()

	// Get MCP Server from env 2
	importedMCPServer := GetMCPServer(t, args.DestAPIM, args.MCPServer.Name, args.MCPServerProvider.Username, args.MCPServerProvider.Password)

	// Validate env 1 and env 2 MCP Server are equal
	ValidateMCPServersEqual(t, args.MCPServer, importedMCPServer)
}

// ValidateMCPServerExport validates MCP Server export
func ValidateMCPServerExport(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export mcp server from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportMCPServer(t, args.MCPServer.Name, args.MCPServer.Version, args.MCPServerProvider.Username, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version))
}

func ValidateExportedMCPServerStructure(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export MCP Server from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportMCPServer(t, args.MCPServer.Name, args.MCPServer.Version, args.MCPServerProvider.Username, args.SrcAPIM.GetEnvName())

	validateMCPServerProject(t, args)

	assert.True(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version))
}

func ValidateExportedMCPServerRevisionStructure(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export MCP Server revision from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportMCPServerRevision(t, args)

	validateMCPServerProject(t, args)

	assert.True(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version))
}

func ValidateExportedMCPServerRevisionFailure(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export MCP Server revision from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, _ := exportMCPServerRevision(t, args)

	assert.Contains(t, output, "404", "Test failed because the response does not contain a not found request")

	// Validate that export failed
	assert.False(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version), "Test failed because the MCP Server Revision was exported successfully")
}

func validateMCPServerProject(t *testing.T, args *MCPServerImportExportTestArgs) {
	assert.True(t, base.IsMCPServerArchiveExists(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		args.MCPServer.Name, args.MCPServer.Version))

	// Check for the existence of the main MCP Server definition file in the archive
	assert.True(t, base.IsFileExistsInMCPServerArchive(t, GetEnvMCPServerExportPath(args.SrcAPIM.GetEnvName()),
		utils.InitProjectDefinitionsSwagger, args.MCPServer.Name, args.MCPServer.Version))
}

func validateMCPServerStructure(t *testing.T, fileData *[]byte, sampleFile string) {
	// Extract the "data" field to an interface
	fileContent := make(map[string]interface{})
	err := yaml.Unmarshal(*fileData, &fileContent)
	if err != nil {
		t.Error(err)
	}
	mcpServerData := fileContent["data"].(map[interface{}]interface{})

	// Read the sample-mcp-server.yaml file in the testdata directory
	sampleData, err := os.ReadFile(sampleFile)
	if err != nil {
		t.Error(err)
	}

	// Extract the "data" field to an interface
	sampleDataContent := make(map[string]interface{})
	err = yaml.Unmarshal(sampleData, &sampleDataContent)
	if err != nil {
		t.Error(err)
	}
	sampleMCPServerData := sampleDataContent["data"].(map[interface{}]interface{})

	// Compare the artifact versions
	exportedMCPServerArtifactVersion := fileContent["version"].(string)
	sampleMCPServerArtifactVersion := sampleDataContent["version"].(string)
	assert.Equal(t, exportedMCPServerArtifactVersion, sampleMCPServerArtifactVersion,
		"Exported artifact version: "+exportedMCPServerArtifactVersion+
			" does not matches with the sample artifact version: "+sampleMCPServerArtifactVersion)

	// Check whether the fields of the MCP Server DTO structure in APICTL has all the fields in MCP Server DTO structure from APIM
	base.Log("\n-----------------------------------------------------------------------------------------")
	base.Log("Checking whether the fields of APICTL MCP Server DTO struct has all the fields from APIM MCP Server DTO struct")
	for key := range mcpServerData {
		keyValue := key.(string)
		_, ok := sampleMCPServerData[key]
		base.Log("\"" + keyValue + "\" is in both the structures")
		if !ok {
			t.Error("Missing \"" + keyValue + "\" in the MCP Server DTO structure from APICTL")
		}
	}

	// Check whether the fields of the MCP Server DTO structure in APIM has all the fields in MCP Server DTO structure from APICTL
	base.Log("\n-----------------------------------------------------------------------------------------")
	base.Log("Checking whether the fields of APIM MCP Server DTO struct has all the fields from APICTL MCP Server DTO struct")
	for key := range sampleMCPServerData {
		keyValue := key.(string)
		_, ok := mcpServerData[key]
		base.Log("\"" + keyValue + "\" is in both the structures")
		if !ok {
			t.Error("Missing \"" + keyValue + "\" in the MCP Server DTO structure from APIM")
		}
	}
}
func GetImportedMCPServer(t *testing.T, args *MCPServerImportExportTestArgs) *apim.MCPServer {
	t.Helper()

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())
	// Import MCP Server to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	_, err := importMCPServer(t, args, true)

	if err != nil {
		t.Fatal(err)
	}

	// Give time for newly imported MCP Server to get indexed, or else GetMCPServer by name will fail
	base.WaitForIndexing()

	// Get MCP Server from env 2
	importedMCPServer := GetMCPServer(t, args.DestAPIM, args.MCPServer.Name, args.MCPServerProvider.Username, args.MCPServerProvider.Password)

	return importedMCPServer
}

func ValidateMCPServerImport(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import MCP Server to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importMCPServer(t, args, true)

	// Give time for newly imported MCP Server to get indexed, or else GetMCPServer by name will fail
	base.WaitForIndexing()

	// Get MCP Server from env 2
	importedMCPServer := GetMCPServer(t, args.DestAPIM, args.MCPServer.Name, args.MCPServerProvider.Username, args.MCPServerProvider.Password)

	// Validate env 1 and env 2 MCP Server is equal
	validateMCPServersEqualCrossTenant(t, args.MCPServer, importedMCPServer)
}

func ValidateMCPServerImportForMultipleVersions(t *testing.T, args *MCPServerImportExportTestArgs, firstImportedMCPServerID string) *apim.MCPServer {
	t.Helper()

	isFirstImport := false
	if strings.EqualFold(firstImportedMCPServerID, "") {
		isFirstImport = true
	}

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import MCP Server to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importMCPServer(t, args, isFirstImport)

	// Give time for newly imported MCP Server to get indexed, or else getMCPServer by name will fail
	base.WaitForIndexing()

	if !isFirstImport {
		args.DestAPIM.DeleteMCPServer(firstImportedMCPServerID)
		base.WaitForIndexing()
	}

	// Get MCP Server from env 2
	importedMCPServer := GetMCPServer(t, args.DestAPIM, args.MCPServer.Name, args.MCPServerProvider.Username, args.MCPServerProvider.Password)

	// Validate env 1 and env 2 MCP Server is equal
	validateMCPServersEqualCrossTenant(t, args.MCPServer, importedMCPServer)

	return importedMCPServer
}

func ValidateMCPServerImportFailure(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Add env2
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import MCP Server to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// importMCPServerPreserveProviderFailure is used to eliminate cleaning the MCP Server after importing
	result, err := importMCPServerPreserveProviderFailure(t, args.SrcAPIM.GetEnvName(), args.MCPServer, args.DestAPIM)

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
}

// ValidateMCPServersEqual : Validate if two MCP Servers are equal while ignoring unique fields
func ValidateMCPServersEqual(t *testing.T, mcpServer1 *apim.MCPServer, mcpServer2 *apim.MCPServer) {
	t.Helper()

	mcpServer1Copy := apim.CopyMCPServer(mcpServer1)
	mcpServer2Copy := apim.CopyMCPServer(mcpServer2)

	same := "override_with_same_value"
	// Since the MCP Servers are from different envs, their respective ID will differ.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	mcpServer1Copy.ID = same
	mcpServer2Copy.ID = same

	mcpServer1Copy.CreatedTime = same
	mcpServer2Copy.CreatedTime = same

	mcpServer1Copy.LastUpdatedTimestamp = same
	mcpServer2Copy.LastUpdatedTimestamp = same

	mcpServer1Copy.LastUpdatedTime = same
	mcpServer2Copy.LastUpdatedTime = same

	mcpServer1Copy.Provider = same
	mcpServer2Copy.Provider = same

	// If an MCP Server is not advertise only, the API owner will be changed during export and import to the current provider
	if (mcpServer1Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		mcpServer1Copy.AdvertiseInformation.ApiOwner = same
	}
	if (mcpServer2Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		mcpServer2Copy.AdvertiseInformation.ApiOwner = same
	}

	if len(mcpServer1Copy.Operations) > 0 {
		overrideOperationPolicyIdsInMcpDefinition(&mcpServer1Copy)
		overrideOperationPolicyIdsInMcpDefinition(&mcpServer2Copy)
	}

	// Sort member collections to make equality check possible
	apim.SortMCPServerMembers(&mcpServer1Copy)
	apim.SortMCPServerMembers(&mcpServer2Copy)

	assert.Equal(t, mcpServer1Copy, mcpServer2Copy, "MCP Server objects are not equal")
}

func overrideOperationPolicyIdsInMcpDefinition(apiCopy *apim.MCPServer) {
	same := "override_with_same_value"
	for _, operation := range apiCopy.Operations {
		requestPolicies := operation.OperationPolicies.Request.([]interface{})
		if len(requestPolicies) > 0 {
			for _, requestPolicy := range requestPolicies {
				requestPolicy.(map[string]interface{})["policyId"] = same
			}
		}
		responsePolicies := operation.OperationPolicies.Response.([]interface{})
		if len(responsePolicies) > 0 {
			for _, responsePolicy := range responsePolicies {
				responsePolicy.(map[string]interface{})["policyId"] = same
			}
		}
		faultPolicies := operation.OperationPolicies.Fault.([]interface{})
		if len(faultPolicies) > 0 {
			for _, faultPolicy := range faultPolicies {
				faultPolicy.(map[string]interface{})["policyId"] = same
			}
		}
	}
}

// ValidateImportedMCPServersEqualToRevision : Validate if the imported MCP Server and exported revision are the same by ignoring
// the unique details and revision specific details.
func ValidateImportedMCPServersEqualToRevision(t *testing.T, mcpServer1 *apim.MCPServer, mcpServer2 *apim.MCPServer) {
	t.Helper()

	mcpServer1Copy := apim.CopyMCPServer(mcpServer1)
	mcpServer2Copy := apim.CopyMCPServer(mcpServer2)

	same := "override_with_same_value"
	// Since the MCP Servers are from different envs, their respective ID will differ.
	mcpServer1Copy.ID = same
	mcpServer2Copy.ID = same

	mcpServer1Copy.CreatedTime = same
	mcpServer2Copy.CreatedTime = same

	mcpServer1Copy.LastUpdatedTimestamp = same
	mcpServer2Copy.LastUpdatedTimestamp = same

	mcpServer1Copy.LastUpdatedTime = same
	mcpServer2Copy.LastUpdatedTime = same

	// When imported the revision as MCP Server, the "IsRevision" property will be false for the imported MCP Server. Hence,
	// the property of the imported MCP Server should be changed
	mcpServer2Copy.IsRevision = true

	// When imported revision as MCP Server, the "RevisionID" property will be 0 for the imported MCP Server. Hence, the property of
	// the imported MCP Server should be changed
	mcpServer2Copy.RevisionID = 1

	// If an API is not advertise only, the API owner will be changed during export and import to the current provider
	if (mcpServer1Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		mcpServer1Copy.AdvertiseInformation.ApiOwner = same
	}
	if (mcpServer2Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		mcpServer2Copy.AdvertiseInformation.ApiOwner = same
	}

	// Sort member collections to make equality check possible
	apim.SortMCPServerMembers(&mcpServer1Copy)
	apim.SortMCPServerMembers(&mcpServer2Copy)

	assert.Equal(t, mcpServer1Copy, mcpServer2Copy, "MCP Server objects are not equal")

}

// ValidateMCPServersList validates the list of MCP Servers
func ValidateMCPServersList(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List MCP Servers of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listMCPServers(t, args)

	mcpServersList := args.SrcAPIM.GetMCPServers()

	ValidateListMCPServersEqual(t, output, mcpServersList)
}

func ValidateListMCPServersEqual(t *testing.T, mcpServersListFromCtl string, mcpServersList *apim.MCPServerList) {
	unmatchedCount := mcpServersList.Count
	for _, mcpServer := range mcpServersList.List {
		// If the output string contains the same MCP Server ID, then decrement the count
		assert.Truef(t, strings.Contains(mcpServersListFromCtl, mcpServer.ID), "mcpServersListFromCtl: "+mcpServersListFromCtl+
			" , does not contain mcpServer.ID: "+mcpServer.ID)
		unmatchedCount--
	}

	// Count == 0 means that all the MCP Servers from mcpServersList were in mcpServersListFromCtl
	assert.Equal(t, 0, unmatchedCount, "MCP Server lists are not equal")
}

// ValidateMCPServersListWithJsonArrayFormat validates the received list of MCP Servers are in JsonArray format and verifies only
// the required ones are there and others are not in the command line output
func ValidateMCPServersListWithJsonArrayFormat(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List MCP Servers of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listMCPServersWithJsonArrayFormat(t, args)

	mcpServersList := args.SrcAPIM.GetMCPServers()

	// Validate MCP Servers list with added MCP Servers
	ValidateListMCPServersEqual(t, output, mcpServersList)

	// Validate JsonArray format
	assert.Contains(t, output, "[\n {\n", "Error while listing MCP Servers in JsonArray format")
}

func validateMCPServersEqualCrossTenant(t *testing.T, mcpServer1 *apim.MCPServer, mcpServer2 *apim.MCPServer) {
	t.Helper()

	mcpServer1Copy := apim.CopyMCPServer(mcpServer1)
	mcpServer2Copy := apim.CopyMCPServer(mcpServer2)

	same := "override_with_same_value"
	// Since the MCP Servers are from different envs, their respective ID will differ.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	mcpServer1Copy.ID = same
	mcpServer2Copy.ID = same

	mcpServer1Copy.CreatedTime = same
	mcpServer2Copy.CreatedTime = same

	mcpServer1Copy.LastUpdatedTimestamp = same
	mcpServer2Copy.LastUpdatedTimestamp = same

	mcpServer1Copy.LastUpdatedTime = same
	mcpServer2Copy.LastUpdatedTime = same

	// The contexts and providers will differ since this is a cross tenant import
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	mcpServer1Copy.Context = same
	mcpServer2Copy.Context = same

	mcpServer1Copy.Provider = same
	mcpServer2Copy.Provider = same

	// If an MCP Server is not advertise only, the API owner will be changed during export and import to the current provider
	if (mcpServer1Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		mcpServer1Copy.AdvertiseInformation.ApiOwner = same
	}
	if (mcpServer2Copy.AdvertiseInformation != apim.AdvertiseInfo{}) {
		mcpServer2Copy.AdvertiseInformation.ApiOwner = same
	}

	// Sort member collections to make equality check possible
	apim.SortMCPServerMembers(&mcpServer1Copy)
	apim.SortMCPServerMembers(&mcpServer2Copy)

	assert.Equal(t, mcpServer1Copy, mcpServer2Copy, "MCP Server objects are not equal")
}

func ValidateMCPServerDelete(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Delete an MCP Server of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()
	mcpServersListBeforeDelete := args.SrcAPIM.GetMCPServers()

	deleteMCPServerByCtl(t, args)

	mcpServersListAfterDelete := args.SrcAPIM.GetMCPServers()
	base.WaitForIndexing()

	// Validate whether the expected number of MCP Server count is there
	assert.Equal(t, mcpServersListBeforeDelete.Count, mcpServersListAfterDelete.Count+1, "Expected number of MCP Servers not deleted")

	// Validate that the delete is a success
	validateMCPServerIsDeleted(t, args.MCPServer, mcpServersListAfterDelete)
}

func ValidateMCPServerDeleteFailure(t *testing.T, args *MCPServerImportExportTestArgs) {
	t.Helper()

	mcpServersListBeforeDelete := args.SrcAPIM.GetMCPServers()

	output, _ := deleteMCPServerByCtl(t, args)

	mcpServersListAfterDelete := args.SrcAPIM.GetMCPServers()
	base.WaitForIndexing()

	// Validate whether the expected number of MCP Server count is there
	assert.NotContains(t, output, "MCP Server deleted successfully!. Status: 200", "MCP Server delete is success with active subscriptions")
	assert.NotEqual(t, mcpServersListBeforeDelete.Count, mcpServersListAfterDelete.Count+1, "Expected number of MCP Servers not deleted")

	t.Cleanup(func() {
		UnsubscribeMCPServer(args.SrcAPIM, args.CtlUser.Username, args.CtlUser.Password, args.MCPServer.ID)
	})
}

func exportMCPServerImportedFromProject(t *testing.T, MCPServerName string, MCPServerVersion string, EnvName string) (string, error) {
	return base.Execute(t, "export", "mcp-server", "-n", MCPServerName, "-v", MCPServerVersion, "-e", EnvName)
}

func exportAllMCPServersOfATenant(t *testing.T, args *MCPServerImportExportTestArgs) (string, error) {
	// Setup environment
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	// Login to the environment
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, err := base.Execute(t, "export", "mcp-servers", "-e", args.SrcAPIM.GetEnvName(), "-k", "--force")
	return output, err
}

func validateMCPServerIsDeleted(t *testing.T, mcpServer *apim.MCPServer, mcpServersListAfterDelete *apim.MCPServerList) {
	for _, existingMCPServer := range mcpServersListAfterDelete.List {
		assert.NotEqual(t, existingMCPServer.ID, mcpServer.ID, "MCP Server delete is not successful")
	}
}

func ValidateChangeLifeCycleStatusOfMCPServer(t *testing.T, args *MCPServerChangeLifeCycleStatusTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	// Execute apictl command to change life cycle of an MCP Server
	output, _ := changeMCPServerLifeCycle(t, args)
	// Assert apictl output
	assert.Contains(t, output, "state changed successfully!", "Error while changing life cycle of MCP Server")

	base.WaitForIndexing()
	// Assert life cycle state after change
	mcpServer := GetMCPServer(t, args.APIM, args.MCPServer.Name, args.CtlUser.Username, args.CtlUser.Password)
	assert.Equal(t, args.ExpectedState, mcpServer.LifeCycleStatus, "Expected Life cycle state change is not equals to actual status")
}

func ValidateChangeLifeCycleStatusOfMCPServerFailure(t *testing.T, args *MCPServerChangeLifeCycleStatusTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	// Execute apictl command to change life cycle of an MCP Server
	output, _ := changeMCPServerLifeCycle(t, args)
	// Assert apictl output
	assert.NotContains(t, output, "state changed successfully!", "Error while changing life cycle of MCP Server")
	assert.NotEqual(t, args.ExpectedState, args.MCPServer.LifeCycleStatus, "Life Cycle State changed successfully")
}

func ValidateMCPServersListWithVersions(t *testing.T, args *MCPServerImportExportTestArgs, newVersion string) {
	t.Helper()

	mcpServers := getMCPServers(args.SrcAPIM, args.CtlUser.Username, args.CtlUser.Password)

	isV1MCPServerExists := false
	isV2MCPServerExists := false

	for _, mcpServer := range mcpServers.List {
		if strings.EqualFold(mcpServer.Version, "1.0.0") && strings.EqualFold(args.MCPServer.Name, mcpServer.Name) {
			isV1MCPServerExists = true
		}
		if strings.EqualFold(mcpServer.Version, newVersion) && strings.EqualFold(args.MCPServer.Name, mcpServer.Name) {
			isV2MCPServerExists = true
		}
	}
	assert.Equal(t, true, isV1MCPServerExists && isV2MCPServerExists)
}

func ValidateMCPServersListWithVersionsFromInitArgs(t *testing.T, args *InitTestArgs, newVersion string) {
	t.Helper()

	mcpServers := getMCPServers(args.SrcAPIM, args.CtlUser.Username, args.CtlUser.Password)

	isV1MCPServerExists := false
	isV2MCPServerExists := false

	for _, mcpServer := range mcpServers.List {
		if strings.EqualFold(mcpServer.Version, "1.0.0") && strings.EqualFold(args.APIName, mcpServer.Name) {
			isV1MCPServerExists = true
		}
		if strings.EqualFold(mcpServer.Version, newVersion) && strings.EqualFold(args.APIName, mcpServer.Name) {
			isV2MCPServerExists = true
		}
	}
	assert.Equal(t, true, isV1MCPServerExists && isV2MCPServerExists)
}

func searchMCPServersWithQuery(t *testing.T, args *MCPServerImportExportTestArgs, query string) (string, error) {
	output, err := base.Execute(t, "get", "mcp-servers", "-e", args.SrcAPIM.EnvName, query, "-k", "--verbose")
	return output, err
}

func ValidateSearchMCPServersList(t *testing.T, args *MCPServerImportExportTestArgs, searchQuery, matchQuery, unmatchedQuery string) {
	t.Helper()

	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	base.WaitForIndexing()

	output, _ := searchMCPServersWithQuery(t, args, searchQuery)

	assert.Truef(t, strings.Contains(output, matchQuery), "mcpServersListFromCtl: "+output+
		" , does not contain the query: "+matchQuery)
	assert.False(t, strings.Contains(output, unmatchedQuery), "mcpServersListFromCtl: "+output+
		" , contains the query: "+unmatchedQuery)
}
