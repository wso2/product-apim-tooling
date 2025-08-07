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
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getMCPServerRevisionsMCPServerName string
var getMCPServerRevisionsMCPServerVersion string
var getMCPServerRevisionsMCPServerProvider string
var getMCPServerRevisionsCmdEnvironment string
var getMCPServerRevisionsCmdFormat string
var getMCPServerRevisionsCmdQuery []string

// GetMCPServerRevisionsCmd related info
const GetMCPServerRevisionsCmdLiteral = "mcp-server-revisions"
const GetMCPServerRevisionsCmdShortDesc = "Display a list of Revisions for the MCP Server"

const GetMCPServerRevisionsCmdLongDesc = `Display a list of Revisions available for the MCP Server in the environment specified`

var getMCPServerRevisionsCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServerRevisionsCmdLiteral + ` -n ChoreoConnect -v 1.0.0 -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServerRevisionsCmdLiteral + ` -n ChoreoConnect -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetMCPServerRevisionsCmdLiteral + ` -n ChoreoConnect -v 1.0.0 -q deployed:true -e dev
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory.`

// getMCPServerRevisionsCmd represents the revisions command
var getMCPServerRevisionsCmd = &cobra.Command{
	Use:     GetMCPServerRevisionsCmdLiteral,
	Short:   GetMCPServerRevisionsCmdShortDesc,
	Long:    GetMCPServerRevisionsCmdLongDesc,
	Example: getMCPServerRevisionsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetMCPServerRevisionsCmdLiteral + " called")
		cred, err := GetCredentials(getMCPServerRevisionsCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetMCPServerRevisionsCmd(cred)
	},
}

func executeGetMCPServerRevisionsCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getMCPServerRevisionsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'get revisions' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetMCPServerRevisionsCmdLiteral+"'", err)
	}

	_, revisions, err := impl.GetMCPServerRevisionListFromEnv(accessToken, getMCPServerRevisionsCmdEnvironment, getMCPServerRevisionsMCPServerName,
		getMCPServerRevisionsMCPServerVersion, getMCPServerRevisionsMCPServerProvider, strings.Join(getMCPServerRevisionsCmdQuery, queryParamSeparator))
	if err == nil {
		impl.PrintMCPServerRevisions(revisions, getMCPServerRevisionsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of MCP Server Revisions", err)
		utils.HandleErrorAndExit("Error getting the list of MCP Server revisions.", err)
	}
}

func init() {
	GetCmd.AddCommand(getMCPServerRevisionsCmd)
	getMCPServerRevisionsCmd.Flags().StringVarP(&getMCPServerRevisionsMCPServerName, "name", "n", "",
		"Name of the MCP Server to get the revision")
	getMCPServerRevisionsCmd.Flags().StringVarP(&getMCPServerRevisionsMCPServerVersion, "version", "v", "",
		"Version of the MCP Server to get the revision")
	getMCPServerRevisionsCmd.Flags().StringVarP(&getMCPServerRevisionsMCPServerProvider, "provider", "r", "",
		"Provider of the MCP Server")
	getMCPServerRevisionsCmd.Flags().StringSliceVarP(&getMCPServerRevisionsCmdQuery, "query", "q",
		[]string{}, "Query pattern")
	getMCPServerRevisionsCmd.Flags().StringVarP(&getMCPServerRevisionsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getMCPServerRevisionsCmd.Flags().StringVarP(&getMCPServerRevisionsCmdFormat, "format", "", "", "Pretty-print revisions "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getMCPServerRevisionsCmd.MarkFlagRequired("name")
	_ = getMCPServerRevisionsCmd.MarkFlagRequired("version")
	_ = getMCPServerRevisionsCmd.MarkFlagRequired("environment")
}
