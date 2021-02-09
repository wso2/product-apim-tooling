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
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getRevisionsAPIName string
var getRevisionsAPIVersion string
var getRevisionsAPIProvider string
var getRevisionsCmdEnvironment string
var getRevisionsCmdFormat string
var getRevisionsCmdQuery string

// GetRevisionsCmd related info
const GetRevisionsCmdLiteral = "revisions"
const GetRevisionsCmdShortDesc = "Display a list of Revisions for the API"

const GetRevisionsCmdLongDesc = `Display a list of Revisions available for the API in the environment specified`

var getRevisionsCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetRevisionsCmdLiteral + ` -n PizzaAPI -v 1.0.0 -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetRevisionsCmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory.`

// getRevisionsCmd represents the revisions command
var getRevisionsCmd = &cobra.Command{
	Use:     GetRevisionsCmdLiteral,
	Short:   GetRevisionsCmdShortDesc,
	Long:    GetRevisionsCmdLongDesc,
	Example: getRevisionsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetRevisionsCmdLiteral + " called")
		cred, err := GetCredentials(getRevisionsCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetRevisionsCmd(cred)
	},
}

func executeGetRevisionsCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getRevisionsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'get revisions' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetRevisionsCmdLiteral+"'", err)
	}

	_, revisions, err := impl.GetRevisionListFromEnv(accessToken, getRevisionsCmdEnvironment, getRevisionsAPIName,
		getRevisionsAPIVersion, getRevisionsAPIProvider, getRevisionsCmdQuery)
	if err == nil {
		impl.PrintRevisions(revisions, getRevisionsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of Revisions", err)
	}
}

func init() {
	GetCmd.AddCommand(getRevisionsCmd)
	getRevisionsCmd.Flags().StringVarP(&getRevisionsAPIName, "name", "n", "",
		"Name of the API to get the revision")
	getRevisionsCmd.Flags().StringVarP(&getRevisionsAPIVersion, "version", "v", "",
		"Version of the API to get the revision")
	getRevisionsCmd.Flags().StringVarP(&getRevisionsAPIProvider, "provider", "r", "",
		"Provider of the API")
	getRevisionsCmd.Flags().StringVarP(&getRevisionsCmdQuery, "query", "q",
		"", "Query pattern")
	getRevisionsCmd.Flags().StringVarP(&getRevisionsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getRevisionsCmd.Flags().StringVarP(&getRevisionsCmdFormat, "format", "", "", "Pretty-print revisions "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getRevisionsCmd.MarkFlagRequired("name")
	_ = getRevisionsCmd.MarkFlagRequired("version")
	_ = getRevisionsCmd.MarkFlagRequired("environment")
}
