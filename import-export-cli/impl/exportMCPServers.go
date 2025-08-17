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

package impl

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/spf13/cast"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var mcpServerListOffset int //from which index of MCP Server, the MCP Servers will be fetched from APIM server
var mcpServerCount int32    // size of MCP Server list to be exported or number of MCP Servers left to be exported from last iteration
var mcpServers []utils.MCPServer

var startingMCPServerIndexFromList int

// Prepare resumption of previous-halted export-mcp-servers operation
func PrepareMCPServerResumption(credential credentials.Credential, exportRelatedFilesPath, cmdResourceTenantDomain, cmdUsername, cmdExportEnvironment string) {
	var lastSucceededMCPServer utils.MCPServer
	lastSucceededMCPServer = utils.ReadLastSucceededMCPServerFileData(exportRelatedFilesPath)
	var migrationMCPServersExportMetadata utils.MigrationMCPServersExportMetadata
	err := migrationMCPServersExportMetadata.ReadMigrationMCPServersExportMetadataFile(filepath.Join(exportRelatedFilesPath,
		utils.MigrationMCPServersExportMetadataFileName))
	if err != nil {
		utils.HandleErrorAndExit("Error loading metadata for resume from "+filepath.Join(exportRelatedFilesPath,
			utils.MigrationMCPServersExportMetadataFileName), err)
	}
	mcpServers = migrationMCPServersExportMetadata.MCPServerListToExport
	mcpServerListOffset = migrationMCPServersExportMetadata.MCPServerListOffset
	startingMCPServerIndexFromList = getLastSucceededMCPServerIndex(lastSucceededMCPServer) + 1

	//find count of MCP Servers left to be exported
	mcpServerCount = int32(len(mcpServers) - startingMCPServerIndexFromList)

	if mcpServerCount == 0 {
		//last iteration had been completed successfully but operation had halted at that point.
		//So get the next set of MCP Servers for next iteration
		startingMCPServerIndexFromList = 0
		mcpServerCount, mcpServers = getMCPServerList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
		if len(mcpServers)-startingMCPServerIndexFromList > 0 {
			utils.WriteMigrationMCPServersExportMetadataFile(mcpServers, cmdResourceTenantDomain, cmdUsername,
				exportRelatedFilesPath, mcpServerListOffset)
		} else {
			fmt.Println("Command: export mcp-servers execution completed !")
		}
	}
}

// Delete directories where the MCP Servers are exported, reset the indexes, get first MCP Server list and write the
// migration-mcp-servers-export-metadata.yaml file
func PrepareMCPServerStartFromBeginning(credential credentials.Credential, exportRelatedFilesPath, cmdResourceTenantDomain, cmdUsername, cmdExportEnvironment string) {
	fmt.Println("Cleaning all the previously exported MCP Servers of the given target tenant, in the given environment if " +
		"any, and prepare to export MCP Servers from beginning")
	//cleaning existing old files (if exists) related to exportation
	if err := utils.RemoveDirectoryIfExists(filepath.Join(exportRelatedFilesPath, utils.ExportedMCPServersDirName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}
	if err := utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.MigrationMCPServersExportMetadataFileName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}
	if err := utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.LastSucceededMCPServerFileName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}

	mcpServerListOffset = 0
	startingMCPServerIndexFromList = 0
	mcpServerCount, mcpServers = getMCPServerList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
	fmt.Printf("Found %d MCP Servers to be exported in the environment: %s, tenant: %s\n", mcpServerCount, cmdExportEnvironment, cmdResourceTenantDomain)
	//write migration-mcp-servers-export-metadata.yaml file
	utils.WriteMigrationMCPServersExportMetadataFile(mcpServers, cmdResourceTenantDomain, cmdUsername, exportRelatedFilesPath,
		mcpServerListOffset)
}

// get the index of the finally (successfully) exported MCP Server from the list of MCP Servers listed in migration-mcp-servers-export-metadata.yaml
func getLastSucceededMCPServerIndex(lastSucceededMCPServer utils.MCPServer) int {
	for i := 0; i < len(mcpServers); i++ {
		if (mcpServers[i].Name == lastSucceededMCPServer.Name) &&
			(mcpServers[i].Provider == lastSucceededMCPServer.Provider) &&
			(mcpServers[i].Version == lastSucceededMCPServer.Version) {
			return i
		}
	}
	return -1
}

// Get the list of MCP Servers from the defined offset index, upto the limit of constant value utils.MaxMCPServersToExportOnce
func getMCPServerList(credential credentials.Credential, cmdExportEnvironment, cmdResourceTenantDomain string) (count int32, mcpServers []utils.MCPServer) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)
	if preCommandErr == nil {
		mcpServerListEndpoint := utils.GetMcpServerListEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		mcpServerListEndpoint += "?limit=" + strconv.Itoa(utils.MaxMCPServersToExportOnce) + "&offset=" + strconv.Itoa(mcpServerListOffset)
		if cmdResourceTenantDomain != "" {
			mcpServerListEndpoint += "&tenantDomain=" + cmdResourceTenantDomain
		}
		count, mcpServers, err := GetMCPServerList(accessToken, mcpServerListEndpoint, "", "")
		if err == nil {
			return count, mcpServers
		} else {
			utils.HandleErrorAndExit(utils.LogPrefixError+"Getting List of MCP Servers.", utils.GetHttpErrorResponse(err))
		}
	} else {
		utils.HandleErrorAndExit(utils.LogPrefixError+"Error in getting access token for user while getting "+
			"the list of MCP Servers: ", preCommandErr)
	}
	return 0, nil
}

// Get the revisions associated with the MCP server
func getRevisionsListForMCPServer(accessToken, environment string, mcpServer utils.MCPServer,
	exportAllRevisions bool) (count int32, revisions []utils.Revisions, err error) {
	var query string
	if !exportAllRevisions {
		query = "deployed:true"
	}
	return GetMCPServerRevisionListFromEnv(accessToken, environment, mcpServer.Name, mcpServer.Version, mcpServer.Provider, query)
}

// Do the MCP Server exportation
func ExportMCPServers(credential credentials.Credential, exportRelatedFilesPath, environment, tenantDomain, format,
	username, mcpServerExportDir string, preserveStatus, runningExportMCPServerCommand, allRevisions,
	preserveCredentials bool) {

	if mcpServerCount == 0 {
		fmt.Println("No MCP Servers available to be exported..!")
	} else {
		var counterSucceededMCPServers = 0
		for mcpServerCount > 0 {
			utils.Logln(utils.LogPrefixInfo+"Found ", mcpServerCount, "of MCP Servers to be exported in the iteration beginning with the offset #"+
				strconv.Itoa(mcpServerListOffset)+". Maximum limit of MCP Servers exported in single iteration is "+
				strconv.Itoa(utils.MaxMCPServersToExportOnce))
			accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, environment)
			if preCommandErr == nil {
				for i := startingMCPServerIndexFromList; i < len(mcpServers); i++ {
					if allRevisions {
						// Export the working copy of the MCP server
						exportMCPServerAndWriteToZip(mcpServers[i], "", accessToken, environment, mcpServerExportDir,
							exportRelatedFilesPath, format, preserveStatus, runningExportMCPServerCommand,
							preserveCredentials)
						counterSucceededMCPServers++
					}
					revisionCount, revisions, err := getRevisionsListForMCPServer(accessToken, environment, mcpServers[i],
						allRevisions)
					if err != nil {
						fmt.Printf("An error occurred while getting the revisions list for MCP Server %s_%s: %v\n",
							mcpServers[i].Name, mcpServers[i].Version, err)
					} else if revisionCount > 0 {
						for j := 0; j < len(revisions); j++ {
							exportMCPServerRevision := utils.GetRevisionNumFromRevisionName(revisions[j].RevisionNumber)
							exportMCPServerAndWriteToZip(mcpServers[i], exportMCPServerRevision, accessToken, environment,
								mcpServerExportDir, exportRelatedFilesPath, format, preserveStatus, runningExportMCPServerCommand,
								preserveCredentials)
							counterSucceededMCPServers++
						}
					}
				}
			} else {
				// error getting OAuth tokens
				fmt.Println("Error getting OAuth Tokens : " + preCommandErr.Error())
			}
			fmt.Println("Batch of " + cast.ToString(mcpServerCount) + " MCP Servers exported successfully..!")

			mcpServerListOffset += utils.MaxMCPServersToExportOnce
			mcpServerCount, mcpServers = getMCPServerList(credential, environment, tenantDomain)
			startingMCPServerIndexFromList = 0
			if len(mcpServers) > 0 {
				utils.WriteMigrationMCPServersExportMetadataFile(mcpServers, tenantDomain, username,
					exportRelatedFilesPath, mcpServerListOffset)
			}
		}
		fmt.Println("\nTotal number of MCP Servers exported: " + cast.ToString(counterSucceededMCPServers))
		fmt.Println("MCP Server export path: " + mcpServerExportDir)
		fmt.Println("\nCommand: export-mcp-servers execution completed !")
	}
}

// Export the MCP Server and archive to zip format
func exportMCPServerAndWriteToZip(mcpServer utils.MCPServer, revisionNumber, accessToken, environment, mcpServerExportDir,
	exportRelatedFilesPath, format string, preserveStatus, runningExportMCPServerCommand,
	preserveCredentials bool) {

	exportMCPServerName := mcpServer.Name
	exportMCPServerVersion := mcpServer.Version
	exportMCPServerProvider := mcpServer.Provider
	var exportMCPServerRevision string
	if revisionNumber != "" {
		exportMCPServerRevision = utils.GetRevisionNumFromRevisionName(revisionNumber)
	}

	resp, err := ExportMCPServerFromEnv(accessToken, exportMCPServerName, exportMCPServerVersion, exportMCPServerRevision,
		exportMCPServerProvider, format, environment, preserveStatus, false, preserveCredentials)
	if err != nil {
		utils.HandleErrorAndExit("Error exporting MCP Server", err)
	}

	if resp.StatusCode() == http.StatusOK {
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		WriteMCPServerToZip(exportMCPServerName, exportMCPServerVersion, exportMCPServerRevision, mcpServerExportDir, runningExportMCPServerCommand, resp)
		//write on last-succeeded-mcp-server.log
		utils.WriteLastSucceededMCPServerFileData(exportRelatedFilesPath, mcpServer)
	} else {
		fmt.Printf("Error exporting MCP Server: %s - %s of Provider: %s\n", exportMCPServerName, exportMCPServerVersion, exportMCPServerProvider)
		utils.PrintErrorResponseAndExit(resp)
	}
}

// Create the required directory structure to save the exported APIs
func CreateExportMCPServersDirStructure(artifactExportDirectory, cmdResourceTenantDomain, cmdExportEnvironment string, cmdForceStartFromBegin bool) string {
	var resourceTenantDirName = utils.GetMigrationExportTenantDirName(cmdResourceTenantDomain)

	var createDirError error
	createDirError = utils.CreateDirIfNotExist(artifactExportDirectory)

	migrationsArtifactsEnvPath := filepath.Join(artifactExportDirectory, cmdExportEnvironment)
	migrationsArtifactsEnvTenantPath := filepath.Join(migrationsArtifactsEnvPath, resourceTenantDirName)
	migrationsArtifactsEnvTenantApisPath := filepath.Join(migrationsArtifactsEnvTenantPath, utils.ExportedMCPServersDirName)

	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvPath)
	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvTenantPath)

	if dirExists, _ := utils.IsDirExists(migrationsArtifactsEnvTenantApisPath); dirExists {
		if cmdForceStartFromBegin {
			utils.RemoveDirectory(migrationsArtifactsEnvTenantApisPath)
			createDirError = utils.CreateDir(migrationsArtifactsEnvTenantApisPath)
		}
	} else {
		createDirError = utils.CreateDir(migrationsArtifactsEnvTenantApisPath)
	}

	if createDirError != nil {
		utils.HandleErrorAndExit("Error in creating directory structure for the API export for migration .",
			createDirError)
	}
	return migrationsArtifactsEnvTenantApisPath
}
