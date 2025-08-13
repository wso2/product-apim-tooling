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

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const ExportMCPServersCmdLiteral = "mcp-servers"
const exportMCPServersCmdShortDesc = "Export MCP Servers for migration"

const exportMCPServersCmdLongDesc = "Export all the MCP Servers of a tenant from one environment, to be imported " +
	"into another environment"
const exportMCPServersCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportMCPServersCmdLiteral + ` -e production --force
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportMCPServersCmdLiteral + ` -e production
NOTE: The flag (--environment (-e)) is mandatory`

var exportMCPServersFormat string
var exportMCPServersAllRevisions bool

var ExportMCPServersCmd = &cobra.Command{
	Use: ExportMCPServersCmdLiteral + " (--environment " +
		"<environment-from-which-artifacts-should-be-exported> --format <export-format> --preserve-status --force)",
	Short:   exportMCPServersCmdShortDesc,
	Long:    exportMCPServersCmdLongDesc,
	Example: exportMCPServersCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportMCPServersCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeExportMCPServersCmd(cred, artifactExportDirectory)
	},
}

// Do operations to export mcp servers for the migration into the directory passed as exportDirectory
// <export_directory> is the patch defined in main_config.yaml
// exportDirectory = <export_directory>/migration/
func executeExportMCPServersCmd(credential credentials.Credential, exportDirectory string) {
	//create dir structure
	mcpServerExportDir := impl.CreateExportAPIsDirStructure(exportDirectory, CmdResourceTenantDomain, CmdExportEnvironment,
		CmdForceStartFromBegin)
	exportRelatedFilesPath := filepath.Join(exportDirectory, CmdExportEnvironment,
		utils.GetMigrationExportTenantDirName(CmdResourceTenantDomain))
	//e.g. /home/samithac/.wso2apictl/exported/migration/production-2.5/wso2-dot-org
	startFromBeginning = false
	isProcessCompleted = false

	fmt.Println("\nExporting MCP Servers for the migration...")
	if CmdForceStartFromBegin {
		startFromBeginning = true
	}

	if (utils.IsFileExist(filepath.Join(exportRelatedFilesPath, utils.LastSucceededMCPServerFileName))) && !startFromBeginning {
		impl.PrepareMCPServerResumption(credential, exportRelatedFilesPath, CmdResourceTenantDomain, CmdUsername, CmdExportEnvironment)
	} else {
		impl.PrepareMCPServerStartFromBeginning(credential, exportRelatedFilesPath, CmdResourceTenantDomain, CmdUsername, CmdExportEnvironment)
	}

	impl.ExportMCPServers(credential, exportRelatedFilesPath, CmdExportEnvironment, CmdResourceTenantDomain, exportMCPServersFormat,
		CmdUsername, mcpServerExportDir, exportMCPServerPreserveStatus, runningExportMCPServerCommand, exportMCPServersAllRevisions, false)
}

func init() {
	ExportCmd.AddCommand(ExportMCPServersCmd)
	ExportMCPServersCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment from which the MCP Servers should be exported")
	ExportMCPServersCmd.PersistentFlags().BoolVarP(&CmdForceStartFromBegin, "force", "", false,
		"Clean all the previously exported MCP Servers of the given target tenant, in the given environment if "+
			"any, and to export MCP Servers from beginning")
	ExportMCPServersCmd.Flags().BoolVarP(&exportMCPServerPreserveStatus, "preserve-status", "", true,
		"Preserve MCP Server status when exporting. Otherwise MCP Server will be exported in CREATED status")
	ExportMCPServersCmd.Flags().BoolVarP(&exportMCPServerPreserveCredentials, "preserve-credentials", "", false,
		"Preserve endpoint credentials when exporting. Otherwise credentials will not be exported")
	ExportMCPServersCmd.Flags().BoolVarP(&exportMCPServersAllRevisions, "all", "", false,
		"Export working copy and all revisions for the MCP Servers in the environments ")
	ExportMCPServersCmd.Flags().StringVarP(&exportMCPServersFormat, "format", "", utils.DefaultExportFormat, "File format of exported archives(json or yaml)")
	_ = ExportMCPServersCmd.MarkFlagRequired("environment")
}
