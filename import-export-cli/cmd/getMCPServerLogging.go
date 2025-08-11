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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getMCPServerLoggingEnvironment string
var getMCPServerLoggingMCPServerId string
var getMCPServerLoggingTenantDomain string
var getMCPServerLoggingCmdFormat string

const GetMCPServerLoggingCmdLiteral = "mcp-server-logging"
const getMCPServerLoggingCmdShortDesc = "Display a list of MCP Server loggers in an environment"
const getMCPServerLoggingCmdLongDesc = `Display a list of MCP Server loggers available for the MCP Servers in the environment specified`

var getMCPServerLoggingCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServerLoggingCmdLiteral + ` -e dev --tenant-domain carbon.super
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServerLoggingCmdLiteral + ` --mcp-server-id bf36ca3a-0332-49ba-abce-e9992228ae06 -e dev --tenant-domain carbon.super`

var getMCPServerLoggingCmd = &cobra.Command{
	Use:     GetMCPServerLoggingCmdLiteral,
	Short:   getMCPServerLoggingCmdShortDesc,
	Long:    getMCPServerLoggingCmdLongDesc,
	Example: getMCPServerLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetCmdLiteral + " " + GetMCPServerLoggingCmdLiteral + " called")
		cred, err := GetCredentials(getMCPServerLoggingEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetMCPServerLoggingCmd(cred)
	},
}

func executeGetMCPServerLoggingCmd(credential credentials.Credential) {
	if getMCPServerLoggingMCPServerId != "" {
		mcpServer, err := impl.GetPerMCPServerLoggingDetailsFromEnv(credential, getMCPServerLoggingEnvironment, getMCPServerLoggingMCPServerId, getMCPServerLoggingTenantDomain)
		if err == nil {
			impl.PrintMCPServerLoggers(mcpServer, getMCPServerLoggingCmdFormat)
		} else {
			utils.Logln(utils.LogPrefixError+"Getting the log level of the MCP Server", err)
			utils.HandleErrorAndExit("Error while getting the log level of the MCP Server", err)
		}
	} else {
		mcpServers, err := impl.GetPerMCPServerLoggingListFromEnv(credential, getMCPServerLoggingEnvironment, getMCPServerLoggingTenantDomain)
		if err == nil {
			impl.PrintMCPServerLoggers(mcpServers, getMCPServerLoggingCmdFormat)
		} else {
			utils.Logln(utils.LogPrefixError+"Getting list of MCP Server log levels for the MCP Servers", err)
			utils.HandleErrorAndExit("Error while getting list of MCP Server log levels for the MCP Servers", err)
		}
	}
}

func init() {
	GetCmd.AddCommand(getMCPServerLoggingCmd)

	getMCPServerLoggingCmd.Flags().StringVarP(&getMCPServerLoggingMCPServerId, "mcp-server-id", "i",
		"", "MCP Server ID")
	getMCPServerLoggingCmd.Flags().StringVarP(&getMCPServerLoggingTenantDomain, "tenant-domain", "",
		"", "Tenant Domain")
	getMCPServerLoggingCmd.Flags().StringVarP(&getMCPServerLoggingEnvironment, "environment", "e",
		"", "Environment of the MCP Servers which the MCP Server loggers should be displayed")
	getMCPServerLoggingCmd.Flags().StringVarP(&getMCPServerLoggingCmdFormat, "format", "", "", "Pretty-print MCP Server loggers "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getMCPServerLoggingCmd.MarkFlagRequired("environment")
}
