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
	"fmt"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var undeployMCPServerName string
var undeployMCPServerVersion string
var undeployMCPServerRevisionNum string
var undeployMCPServerProvider string
var undeployMCPServerEnvironment string
var undeployMCPServerCmdGatewayEnvs []string

const UndeployMCPServerCmdLiteral = "mcp-server"
const undeployMCPServerCmdShortDesc = "Undeploy MCP Server"
const undeployMCPServerCmdLongDesc = "Undeploy an MCP Server revision from gateway environments"
const undeployMCPServerCmdExamples = utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployMCPServerCmdLiteral + ` -n MyMCPServer -v 1.0.0 --rev 2 -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployMCPServerCmdLiteral + ` -n MyMCPServer -v 2.1.0 --rev 6 -g Label1 -g Label2 -g Label3 -e production
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployMCPServerCmdLiteral + ` -n MyMCPServer -v 2.1.0 -r alice --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e)) are mandatory.
If the flag (--gateway-env (-g)) is not provided, revision will be undeployed from all deployed gateway environments.`

var UndeployMCPServerCmd = &cobra.Command{
	Use: UndeployMCPServerCmdLiteral + " (--name <name-of-the-mcpserver> --version <version-of-the-mcpserver> --provider <provider-of-the-mcpserver> " +
		"--rev <revision-number-of-the-mcpserver> --gateway-env <gateway-environment> " +
		"--environment <environment-from-which-the-mcpserver-should-be-undeployed>)",
	Short:   undeployMCPServerCmdShortDesc,
	Long:    undeployMCPServerCmdLongDesc,
	Example: undeployMCPServerCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + UndeployMCPServerCmdLiteral + " called")
		if len(undeployMCPServerCmdGatewayEnvs) > 0 {
			undeployAllGatewayEnvs = false
		}
		gateways := generateGatewayEnvsArray(undeployMCPServerCmdGatewayEnvs)

		cred, err := GetCredentials(undeployMCPServerEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeUndeployMCPServerCmd(cred, gateways)
	},
}

func executeUndeployMCPServerCmd(credential credentials.Credential, deployments []utils.Deployment) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, undeployMCPServerEnvironment)
	if preCommandErr == nil {
		resp, err := impl.UndeployRevisionFromGatewaysMCPServer(accessToken,
			undeployMCPServerEnvironment, undeployMCPServerName, undeployMCPServerVersion, undeployMCPServerProvider, undeployMCPServerRevisionNum,
			deployments, undeployAllGatewayEnvs)
		if err != nil {
			utils.HandleErrorAndExit("Error while undeploying the MCP Server", err)
		}
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusCreated {
			fmt.Println("Revision " + undeployMCPServerRevisionNum + " of MCP Server " + undeployMCPServerName + "_" + undeployMCPServerVersion +
				" successfully undeployed from the specified gateways environments")
		} else {
			fmt.Println("Error while undeploying the MCP Server: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens to undeploy the MCP Server:" + preCommandErr.Error())
	}
}

func init() {
	UndeployCmd.AddCommand(UndeployMCPServerCmd)
	UndeployMCPServerCmd.Flags().StringVarP(&undeployMCPServerName, "name", "n", "",
		"Name of the MCP Server to be undeployed")
	UndeployMCPServerCmd.Flags().StringVarP(&undeployMCPServerVersion, "version", "v", "",
		"Version of the MCP Server to be undeployed")
	UndeployMCPServerCmd.Flags().StringVarP(&undeployMCPServerProvider, "provider", "r", "",
		"Provider of the MCP Server")
	UndeployMCPServerCmd.Flags().StringSliceVarP(&undeployMCPServerCmdGatewayEnvs, "gateway-env", "g", []string{},
		"Gateway environment which the revision has to be undeployed")
	UndeployMCPServerCmd.Flags().StringVarP(&undeployMCPServerRevisionNum, "rev", "", "",
		"Revision number of the MCP Server to undeploy")
	UndeployMCPServerCmd.Flags().StringVarP(&undeployMCPServerEnvironment, "environment", "e",
		"", "Environment of which the MCP Server should be undeployed")
	_ = UndeployMCPServerCmd.MarkFlagRequired("name")
	_ = UndeployMCPServerCmd.MarkFlagRequired("version")
	_ = UndeployMCPServerCmd.MarkFlagRequired("rev")
	_ = UndeployMCPServerCmd.MarkFlagRequired("environment")
}
