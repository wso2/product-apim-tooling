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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	importMCPServerFile                string
	importMCPServerEnvironment         string
	importMCPServerCmdPreserveProvider bool
	importMCPServerUpdate              bool
	importMCPServerParamsFile          string
	importMCPServerSkipCleanup         bool
	importMCPServerRotateRevision      bool
	importMCPServerSkipDeployments     bool
	mcpServerDryRun                    bool
	mcpServerLoggingCmdFormat          string
)

const (
	// ImportMCPServer command related usage info
	ImportMCPServerCmdLiteral   = "mcp-server"
	importMCPServerCmdShortDesc = "Import MCP Server"
	importMCPServerCmdLongDesc  = "Import an MCP Server to an environment"
)

const importMCPServerCmdExamples = utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportMCPServerCmdLiteral + ` -f qa/ChoreoConnect.zip -e dev
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportMCPServerCmdLiteral + ` -f staging/ChoreoConnect.zip -e production
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportMCPServerCmdLiteral + ` -f ~/my-mcp-server -e production --update --rotate-revision
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportMCPServerCmdLiteral + ` -f ~/my-mcp-server -e production --update
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// ImportMCPServerCmd represents the importMCPServer command
var ImportMCPServerCmd = &cobra.Command{
	Use: ImportMCPServerCmdLiteral + " --file <path-to-mcp-server> --environment " +
		"<environment>",
	Short:   importMCPServerCmdShortDesc,
	Long:    importMCPServerCmdLongDesc,
	Example: importMCPServerCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ImportMCPServerCmdLiteral + " called")
		cred, err := GetCredentials(importMCPServerEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(cred, importMCPServerEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing MCP Server", err)
		}
		err = impl.ImportMCPServerToEnv(accessOAuthToken, importMCPServerEnvironment, importMCPServerFile, importMCPServerParamsFile, importMCPServerUpdate,
			importMCPServerCmdPreserveProvider, importMCPServerSkipCleanup, importMCPServerRotateRevision, importMCPServerSkipDeployments, mcpServerDryRun,
			mcpServerLoggingCmdFormat)
		if err != nil {
			utils.HandleErrorAndExit("Error importing MCP Server", err)
			return
		}
	},
}

// init using Cobra
func init() {
	ImportCmd.AddCommand(ImportMCPServerCmd)
	ImportMCPServerCmd.Flags().StringVarP(&importMCPServerFile, "file", "f", "",
		"Name of the MCP Server to be imported")
	ImportMCPServerCmd.Flags().StringVarP(&importMCPServerEnvironment, "environment", "e",
		"", "Environment from the which the MCP Server should be imported")
	ImportMCPServerCmd.Flags().BoolVar(&importMCPServerCmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of MCP Server after importing")
	ImportMCPServerCmd.Flags().BoolVar(&importMCPServerUpdate, "update", false, "Update an "+
		"existing MCP Server or create a new MCP Server")
	ImportMCPServerCmd.Flags().BoolVar(&importMCPServerRotateRevision, "rotate-revision", false, "Rotate the "+
		"revisions with each update")
	ImportMCPServerCmd.Flags().BoolVar(&importMCPServerSkipDeployments, "skip-deployments", false, "Update only "+
		"the working copy and skip deployment steps in import")
	ImportMCPServerCmd.Flags().StringVarP(&importMCPServerParamsFile, "params", "", "", "Provide an API Manager params file "+
		"or a directory generated using \"gen deployment-dir\" command")
	ImportMCPServerCmd.Flags().BoolVarP(&importMCPServerSkipCleanup, "skip-cleanup", "", false, "Leave "+
		"all temporary files created during import process")
	ImportMCPServerCmd.Flags().BoolVarP(&mcpServerDryRun, "dry-run", "", false, "Get "+
		"verification of the governance compliance of the MCP Server without importing it")
	ImportMCPServerCmd.Flags().StringVarP(&mcpServerLoggingCmdFormat, "format", "", "", "Output format of violation results in "+
		"dry-run mode. Supported formats: [table, json, list]. If not provided, the default format is table.")
	// Mark required flags
	_ = ImportMCPServerCmd.MarkFlagRequired("environment")
	_ = ImportMCPServerCmd.MarkFlagRequired("file")
}
