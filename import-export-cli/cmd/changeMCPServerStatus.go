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

	"net/http"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var mcpServerStateChangeEnvironment string
var mcpServerNameForStateChange string
var mcpServerVersionForStateChange string
var mcpServerProviderForStateChange string
var mcpServerStateChangeAction string

// ChangeMCPServerStatus command related usage info
const changeMCPServerStatusCmdLiteral = "mcp-server"
const changeMCPServerStatusCmdShortDesc = "Change Status of an MCP Server"
const changeMCPServerStatusCmdLongDesc = "Change the lifecycle status of an MCP Server in an environment"

const changeMCPServerStatusCmdExamples = utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeMCPServerStatusCmdLiteral + ` -a Publish -n MyMCPServer -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeMCPServerStatusCmdLiteral + ` -a Publish -n MyMCPServer -v 2.1.0 -e production
NOTE: The 4 flags (--action (-a), --name (-n), --version (-v), and --environment (-e)) are mandatory.`

// changeMCPServerStatusCmd represents change-status mcp-server command
var ChangeMCPServerStatusCmd = &cobra.Command{
	Use: changeMCPServerStatusCmdLiteral + " (--action <action-of-the-mcpserver-state-change> --name <name-of-the-mcpserver> --version <version-of-the-mcpserver> --provider " +
		"<provider-of-the-mcpserver> --environment <environment-from-which-the-mcpserver-state-should-be-changed>)",
	Short:   changeMCPServerStatusCmdShortDesc,
	Long:    changeMCPServerStatusCmdLongDesc,
	Example: changeMCPServerStatusCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + changeMCPServerStatusCmdLiteral + " called")
		cred, err := GetCredentials(mcpServerStateChangeEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeChangeMCPServerStatusCmd(cred)
	},
}

// executeChangeMCPServerStatusCmd executes the change mcp server status command
func executeChangeMCPServerStatusCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, mcpServerStateChangeEnvironment)
	if preCommandErr == nil {
		resp, err := impl.ChangeMCPServerStatusInEnv(accessToken, mcpServerStateChangeEnvironment, mcpServerStateChangeAction,
			mcpServerNameForStateChange, mcpServerVersionForStateChange, mcpServerProviderForStateChange)
		if err != nil {
			utils.HandleErrorAndExit("Error while changing the MCP Server status", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			// 200 OK
			fmt.Println(mcpServerNameForStateChange + " MCP Server state changed successfully!")
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// Neither 200 nor 500
			fmt.Println("Error while changing MCP Server Status: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// Error changing the MCP Server status
		fmt.Println("Error getting OAuth tokens while changing status of the MCP Server:" + preCommandErr.Error())
	}
}

func init() {
	ChangeStatusCmd.AddCommand(ChangeMCPServerStatusCmd)
	ChangeMCPServerStatusCmd.Flags().StringVarP(&mcpServerStateChangeAction, "action", "a", "",
		"Action to be taken to change the status of the MCP Server")
	ChangeMCPServerStatusCmd.Flags().StringVarP(&mcpServerNameForStateChange, "name", "n", "",
		"Name of the MCP Server to be state changed")
	ChangeMCPServerStatusCmd.Flags().StringVarP(&mcpServerVersionForStateChange, "version", "v", "",
		"Version of the MCP Server to be state changed")
	ChangeMCPServerStatusCmd.Flags().StringVarP(&mcpServerProviderForStateChange, "provider", "r", "",
		"Provider of the MCP Server")
	ChangeMCPServerStatusCmd.Flags().StringVarP(&mcpServerStateChangeEnvironment, "environment", "e",
		"", "Environment of which the MCP Server state should be changed")
	// Mark required flags
	_ = ChangeMCPServerStatusCmd.MarkFlagRequired("action")
	_ = ChangeMCPServerStatusCmd.MarkFlagRequired("name")
	_ = ChangeMCPServerStatusCmd.MarkFlagRequired("version")
	_ = ChangeMCPServerStatusCmd.MarkFlagRequired("environment")
}
