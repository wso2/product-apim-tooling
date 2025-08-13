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

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"net/http"
	"path/filepath"
)

var exportMCPServerName string
var exportMCPServerVersion string
var exportMCPServerRevisionNum string
var exportMCPServerProvider string
var exportMCPServerPreserveStatus bool
var exportMCPServerFormat string
var runningExportMCPServerCommand bool
var exportMCPServerLatestRevision bool
var exportMCPServerPreserveCredentials bool

// ExportMCPServer command related usage info
const ExportMCPServerCmdLiteral = "mcp-server"
const exportMCPServerCmdShortDesc = "Export MCP Server"

const exportMCPServerCmdLongDesc = "Export an MCP Server from an environment"

const exportMCPServerCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportMCPServerCmdLiteral + ` -n ChoreoConnect -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportMCPServerCmdLiteral + ` -n ChoreoConnect -v 2.1.0 --rev 6 -r admin -e production
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportMCPServerCmdLiteral + ` -n ChoreoConnect -v 2.1.0 --rev 2 -r admin -e production
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory. If --rev is not provided, working copy of the MCP Server
without deployment environments will be exported.`

// ExportMCPServerCmd represents the exportMCPServer command
var ExportMCPServerCmd = &cobra.Command{
	Use: ExportMCPServerCmdLiteral + " (--name <name-of-the-mcp-server> --version <version-of-the-mcp-server> --provider <provider-of-the-mcp-server> --environment " +
		"<environment-from-which-the-mcp-server-should-be-exported>)",
	Short:   exportMCPServerCmdShortDesc,
	Long:    exportMCPServerCmdLongDesc,
	Example: exportMCPServerCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportMCPServerCmdLiteral + " called")
		var mcpServersExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMCPServersDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportMCPServerCmd(cred, mcpServersExportDirectory)
	},
}

func executeExportMCPServerCmd(credential credentials.Credential, exportDirectory string) {
	runningExportMCPServerCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, CmdExportEnvironment)

	if preCommandErr == nil {
		resp, err := impl.ExportMCPServerFromEnv(accessToken, exportMCPServerName, exportMCPServerVersion, exportMCPServerRevisionNum, exportMCPServerProvider,
			exportMCPServerFormat, CmdExportEnvironment, exportMCPServerPreserveStatus, exportMCPServerLatestRevision,
			exportMCPServerPreserveCredentials)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		mcpServerZipLocationPath := filepath.Join(exportDirectory, CmdExportEnvironment)
		if resp.StatusCode() == http.StatusOK {
			impl.WriteToZip(exportMCPServerName, exportMCPServerVersion, "", mcpServerZipLocationPath, runningExportMCPServerCommand, resp)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting MCP Server:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// error exporting MCP Server
		fmt.Println("Error getting OAuth tokens while exporting MCP Server:" + preCommandErr.Error())
	}
}

// init using Cobra
func init() {
	ExportCmd.AddCommand(ExportMCPServerCmd)
	ExportMCPServerCmd.Flags().StringVarP(&exportMCPServerName, "name", "n", "",
		"Name of the MCP Server to be exported")
	ExportMCPServerCmd.Flags().StringVarP(&exportMCPServerVersion, "version", "v", "",
		"Version of the MCP Server to be exported")
	ExportMCPServerCmd.Flags().StringVarP(&exportMCPServerProvider, "provider", "r", "",
		"Provider of the MCP Server")
	ExportMCPServerCmd.Flags().StringVarP(&exportMCPServerRevisionNum, "rev", "", "",
		"Revision number of the MCP Server to be exported")
	ExportMCPServerCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment to which the MCP Server should be exported")
	ExportMCPServerCmd.Flags().BoolVarP(&exportMCPServerPreserveStatus, "preserve-status", "", true,
		"Preserve MCP Server status when exporting. Otherwise MCP Server will be exported in CREATED status")
	ExportMCPServerCmd.Flags().BoolVarP(&exportMCPServerPreserveCredentials, "preserve-credentials", "", false,
		"Preserve endpoint credentials when exporting. Otherwise credentials will not be exported")
	ExportMCPServerCmd.Flags().BoolVarP(&exportMCPServerLatestRevision, "latest", "", false,
		"Export the latest revision of the MCP Server")
	ExportMCPServerCmd.Flags().StringVarP(&exportMCPServerFormat, "format", "", utils.DefaultExportFormat, "File format of exported archive(json or yaml)")
	_ = ExportMCPServerCmd.MarkFlagRequired("name")
	_ = ExportMCPServerCmd.MarkFlagRequired("version")
	_ = ExportMCPServerCmd.MarkFlagRequired("environment")
}
