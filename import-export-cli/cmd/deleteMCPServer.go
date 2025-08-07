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

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var deleteMCPServerEnvironment string
var deleteMCPServerName string
var deleteMCPServerVersion string
var deleteMCPServerProvider string

// DeleteMCPServer command related usage info
const deleteMCPServerCmdLiteral = "mcp-server"
const deleteMCPServerCmdShortDesc = "Delete MCP Server"
const deleteMCPServerCmdLongDesc = "Delete an MCP Server from an environment"

const deleteMCPServerCmdExamplesDefault = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteMCPServerCmdLiteral + ` -n ChoreoConnect -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteMCPServerCmdLiteral + ` -n ChoreoConnect -v 2.1.0 -e production
NOTE: The 3 flags (--name (-n), --version (-v), and --environment (-e)) are mandatory.`

// DeleteMCPServerCmd represents the delete mcp-server command
var DeleteMCPServerCmd = &cobra.Command{
	Use: deleteMCPServerCmdLiteral + " (--name <name-of-the-mcp-server> --version <version-of-the-mcp-server> --provider <provider-of-the-mcp-server> --environment " +
		"<environment-from-which-the-mcp-server-should-be-deleted>)",
	Short:              deleteMCPServerCmdShortDesc,
	Long:               deleteMCPServerCmdLongDesc,
	Example:            deleteMCPServerCmdExamplesDefault,
	DisableFlagParsing: isK8sEnabled(),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteMCPServerCmdLiteral + " called")
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if configVars.Config.KubernetesMode {
			k8sArgs := []string{k8sUtils.K8sDelete, k8sUtils.ApiOpCrdApi}
			k8sArgs = append(k8sArgs, args...)
			ExecuteKubernetes(k8sArgs...)
		} else {
			cred, err := GetCredentials(deleteMCPServerEnvironment)
			if err != nil {
				utils.HandleErrorAndExit("Error getting credentials ", err)
			}
			executeDeleteMCPServerCmd(cred)
		}
	},
}

// executeDeleteMCPServerCmd executes the delete mcp-server command
func executeDeleteMCPServerCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteMCPServerEnvironment)
	if preCommandErr == nil {
		resp, err := impl.DeleteMCPServer(accessToken, deleteMCPServerEnvironment, deleteMCPServerName, deleteMCPServerVersion, deleteMCPServerProvider)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting MCP Server ", err)
		}
		impl.PrintDeleteMCPServerResponse(resp, err)
	} else {
		// Error deleting MCP Server
		fmt.Println("Error getting OAuth tokens while deleting MCP Server:" + preCommandErr.Error())
	}
}

// Init using Cobra
func init() {
	DeleteCmd.AddCommand(DeleteMCPServerCmd)
	DeleteMCPServerCmd.Flags().StringVarP(&deleteMCPServerName, "name", "n", "",
		"Name of the MCP Server to be deleted")
	DeleteMCPServerCmd.Flags().StringVarP(&deleteMCPServerVersion, "version", "v", "",
		"Version of the MCP Server to be deleted")
	DeleteMCPServerCmd.Flags().StringVarP(&deleteMCPServerProvider, "provider", "r", "",
		"Provider of the MCP Server to be deleted")
	DeleteMCPServerCmd.Flags().StringVarP(&deleteMCPServerEnvironment, "environment", "e",
		"", "Environment from which the MCP Server should be deleted")

	// fetches the main-config.yaml file silently; i.e. if it's not created, ignore the error and assume that
	//	this is the default mode.
	configVars := utils.GetMainConfigFromFileSilently(utils.MainConfigFilePath)
	if configVars == nil || !configVars.Config.KubernetesMode {
		// Mark required flags
		_ = DeleteMCPServerCmd.MarkFlagRequired("name")
		_ = DeleteMCPServerCmd.MarkFlagRequired("version")
		_ = DeleteMCPServerCmd.MarkFlagRequired("environment")
	}
}
