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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const queryParamSeparator = " "

// Get command related usage Info
const GetCmdLiteral = "get"
const getCmdShortDesc = "Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments"

const getCmdLongDesc = `Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
Display a list containing all the MCP Servers available in the environment specified by flag (--environment, -e)/
Display a list containing all the API Products available in the environment specified by flag (--environment, -e)/
Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)/
Display a list of API revisions of a specific API in the environment specified by flag (--environment, -e)/
Display a list of MCP Server revisions of a specific MCP Server in the environment specified by flag (--environment, -e)/
Display a list of API Product revisions of a specific API Product in the environment specified by flag (--environment, -e)/
Get a generated JWT token to invoke an API or API Product by subscribing to a default application for testing purposes in the environment specified by flag (--environment, -e)/
Get the log level of each API in the environment specified by flag (--environment, -e)/
Get the log level of each MCP Server in the environment specified by flag (--environment, -e)/
Get the correlation log configurations in the environment specified by flag (--environment, -e)
OR
List all the environments`

const getCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetEnvsCmdLiteral + `
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApisCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServersCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiProductsCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAppsCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPIRevisionsCmdLiteral + ` -n PizzaAPI -v 1.0.0 -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServerRevisionsCmdLiteral + ` -n WeatherMCPServer -v 1.0.0 -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPIProductRevisionsCmdLiteral + ` -n PizzaProduct -v 1.0.0 -e dev
` + utils.ProjectName + " " + GetCmdLiteral + " " + GetKeysCmdLiteral + ` -n TwitterAPI -v 1.0.0 -e dev
` + utils.ProjectName + " " + GetCmdLiteral + " " + GetApiLoggingCmdLiteral + ` -e dev --tenant-domain carbon.super
` + utils.ProjectName + " " + GetCmdLiteral + " " + GetApiLoggingCmdLiteral + ` --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 -e dev --tenant-domain carbon.super
` + utils.ProjectName + " " + GetCmdLiteral + " " + GetMCPServerLoggingCmdLiteral + ` --mcp-server-id bf36ca3a-0332-49ba-abce-e9992228ae06 -e dev --tenant-domain carbon.super
` + utils.ProjectName + " " + GetCmdLiteral + " " + GetCorrelationLoggingCmdLiteral + ` -e dev`

// ListCmd represents the list command
var GetCmd = &cobra.Command{
	Use:     GetCmdLiteral,
	Short:   getCmdShortDesc,
	Long:    getCmdLongDesc,
	Example: getCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(GetCmd)
}
