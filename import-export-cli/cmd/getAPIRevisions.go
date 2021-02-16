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

var getAPIRevisionsAPIName string
var getAPIRevisionsAPIVersion string
var getAPIRevisionsAPIProvider string
var getAPIRevisionsCmdEnvironment string
var getAPIRevisionsCmdFormat string
var getAPIRevisionsCmdQuery string

// GetRevisionsCmd related info
const GetAPIRevisionsCmdLiteral = "api-revisions"
const GetAPIRevisionsCmdShortDesc = "Display a list of Revisions for the API"

const GetAPIRevisionsCmdLongDesc = `Display a list of Revisions available for the API in the environment specified`

var getAPIRevisionsCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPIRevisionsCmdLiteral + ` -n PizzaAPI -v 1.0.0 -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPIRevisionsCmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory.`

// getRevisionsCmd represents the revisions command
var getAPIRevisionsCmd = &cobra.Command{
	Use:     GetAPIRevisionsCmdLiteral,
	Short:   GetAPIRevisionsCmdShortDesc,
	Long:    GetAPIRevisionsCmdLongDesc,
	Example: getAPIRevisionsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetAPIRevisionsCmdLiteral + " called")
		cred, err := GetCredentials(getAPIRevisionsCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetAPIRevisionsCmd(cred)
	},
}

func executeGetAPIRevisionsCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getAPIRevisionsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'get revisions' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetAPIRevisionsCmdLiteral+"'", err)
	}

	_, revisions, err := impl.GetRevisionListFromEnv(accessToken, getAPIRevisionsCmdEnvironment, getAPIRevisionsAPIName,
		getAPIRevisionsAPIVersion, getAPIRevisionsAPIProvider, getAPIRevisionsCmdQuery)
	if err == nil {
		impl.PrintRevisions(revisions, getAPIRevisionsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of API Revisions", err)
	}
}

func init() {
	GetCmd.AddCommand(getAPIRevisionsCmd)
	getAPIRevisionsCmd.Flags().StringVarP(&getAPIRevisionsAPIName, "name", "n", "",
		"Name of the API to get the revision")
	getAPIRevisionsCmd.Flags().StringVarP(&getAPIRevisionsAPIVersion, "version", "v", "",
		"Version of the API to get the revision")
	getAPIRevisionsCmd.Flags().StringVarP(&getAPIRevisionsAPIProvider, "provider", "r", "",
		"Provider of the API")
	getAPIRevisionsCmd.Flags().StringVarP(&getAPIRevisionsCmdQuery, "query", "q",
		"", "Query pattern")
	getAPIRevisionsCmd.Flags().StringVarP(&getAPIRevisionsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getAPIRevisionsCmd.Flags().StringVarP(&getAPIRevisionsCmdFormat, "format", "", "", "Pretty-print revisions "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getAPIRevisionsCmd.MarkFlagRequired("name")
	_ = getAPIRevisionsCmd.MarkFlagRequired("version")
	_ = getAPIRevisionsCmd.MarkFlagRequired("environment")
}
