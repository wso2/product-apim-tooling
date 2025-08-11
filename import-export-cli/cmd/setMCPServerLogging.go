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

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var setMCPServerLoggingEnvironment string
var setMCPServerLoggingMCPServerId string
var setMCPServerLoggingTenantDomain string
var setMCPServerLoggingLogLevel string

const SetMCPServerLoggingCmdLiteral = "mcp-server-logging"
const setMCPServerLoggingCmdShortDesc = "Set the log level for an MCP Server in an environment"
const setMCPServerLoggingCmdLongDesc = `Set the log level for an MCP Server in the environment specified`

var setMCPServerLoggingCmdExamples = utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetMCPServerLoggingCmdLiteral + ` --mcp-server-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level full -e dev --tenant-domain carbon.super
` + utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetMCPServerLoggingCmdLiteral + ` --mcp-server-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level off -e dev --tenant-domain carbon.super`

var setMCPServerLoggingCmd = &cobra.Command{
	Use:     SetMCPServerLoggingCmdLiteral,
	Short:   setMCPServerLoggingCmdShortDesc,
	Long:    setMCPServerLoggingCmdLongDesc,
	Example: setMCPServerLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + SetCmdLiteral + " " + SetMCPServerLoggingCmdLiteral + " called")
		cred, err := GetCredentials(setMCPServerLoggingEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeSetMCPServerLoggingCmd(cred)
	},
}

func executeSetMCPServerLoggingCmd(credential credentials.Credential) {
	resp, err := impl.SetMCPServerLoggingLevel(credential, setMCPServerLoggingEnvironment, setMCPServerLoggingMCPServerId, setMCPServerLoggingTenantDomain, setMCPServerLoggingLogLevel)
	if err != nil {
		utils.Logln(utils.LogPrefixError+"Setting the log level of the MCP Server", err)
		utils.HandleErrorAndExit("Error while setting the log level of the MCP Server", err)
	}
	// Print info on response
	utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
	if resp.StatusCode() == http.StatusOK {
		// 200 OK
		fmt.Println("Log level " + setMCPServerLoggingLogLevel + " is successfully set to the MCP Server.")
	} else {
		fmt.Println("Setting the log level of the MCP Server: ", resp.Status(), "\n", string(resp.Body()))
	}
}

func init() {
	SetCmd.AddCommand(setMCPServerLoggingCmd)

	setMCPServerLoggingCmd.Flags().StringVarP(&setMCPServerLoggingMCPServerId, "mcp-server-id", "i",
		"", "MCP Server ID")
	setMCPServerLoggingCmd.Flags().StringVarP(&setMCPServerLoggingTenantDomain, "tenant-domain", "",
		"", "Tenant Domain")
	setMCPServerLoggingCmd.Flags().StringVarP(&setMCPServerLoggingLogLevel, "log-level", "",
		"", "Log Level")
	setMCPServerLoggingCmd.Flags().StringVarP(&setMCPServerLoggingEnvironment, "environment", "e",
		"", "Environment of the MCP Server which the log level should be set")
	_ = setMCPServerLoggingCmd.MarkFlagRequired("environment")
	_ = setMCPServerLoggingCmd.MarkFlagRequired("mcp-server-id")
	_ = setMCPServerLoggingCmd.MarkFlagRequired("log-level")
}
