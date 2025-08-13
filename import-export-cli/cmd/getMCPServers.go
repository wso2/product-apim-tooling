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
	"strconv"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getMCPServersCmdEnvironment string
var getMCPServersCmdFormat string
var getMCPServersCmdQuery []string
var getMCPServersCmdLimit string

// GetMCPServersCmd related info
const GetMCPServersCmdLiteral = "mcp-servers"
const getMCPServersCmdShortDesc = "Display a list of MCP Servers in an environment"

const getMCPServersCmdLongDesc = `Display a list of MCP Servers in the environment specified by the flag --environment, -e`

var getMCPServersCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServersCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServersCmdLiteral + ` -e dev -q version:1.0.0
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServersCmdLiteral + ` -e prod -q provider:admin -q version:1.0.0
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServersCmdLiteral + ` -e prod -l 100
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServersCmdLiteral + ` -e staging
NOTE: The flag (--environment (-e)) is mandatory`

// getMCPServersCmd represents the mcp-servers command
var getMCPServersCmd = &cobra.Command{
	Use:     GetMCPServersCmdLiteral,
	Short:   getMCPServersCmdShortDesc,
	Long:    getMCPServersCmdLongDesc,
	Example: getMCPServersCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetMCPServersCmdLiteral + " called")
		cred, err := GetCredentials(getMCPServersCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetMCPServersCmd(cred)
	},
}

func executeGetMCPServersCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getMCPServersCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetMCPServersCmdLiteral+"'", err)
	}

	_, mcpServers, err := impl.GetMCPServerListFromEnv(accessToken, getMCPServersCmdEnvironment, strings.Join(getMCPServersCmdQuery, queryParamSeparator), getMCPServersCmdLimit)
	if err == nil {
		impl.PrintMCPServers(mcpServers, getMCPServersCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of MCP Servers", err)
		utils.HandleErrorAndExit("Error getting the list of MCP Servers.", err)
	}
}

func init() {
	GetCmd.AddCommand(getMCPServersCmd)

	getMCPServersCmd.Flags().StringVarP(&getMCPServersCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getMCPServersCmd.Flags().StringSliceVarP(&getMCPServersCmdQuery, "query", "q",
		[]string{}, "Query pattern")
	getMCPServersCmd.Flags().StringVarP(&getMCPServersCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultApisDisplayLimit), "Maximum number of MCP servers to return")
	getMCPServersCmd.Flags().StringVarP(&getMCPServersCmdFormat, "format", "", "", "Pretty-print mcp-servers "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getMCPServersCmd.MarkFlagRequired("environment")
}
