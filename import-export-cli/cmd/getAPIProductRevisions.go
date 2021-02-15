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

var getRevisionsAPIProductName string
var getRevisionsAPIProductProvider string
var getAPIProductRevisionsCmdEnvironment string
var getAPIProductRevisionsCmdFormat string
var getAPIProductRevisionsCmdQuery string

// GetAPIProductRevisionsCmd related info
const GetAPIProductRevisionsCmdLiteral = "api-product-revisions"
const GetAPIProductRevisionsCmdShortDesc = "Display a list of Revisions for the API"

const GetAPIProductRevisionsCmdLongDesc = `Display a list of Revisions available for the API in the environment specified`

var getRevisionsCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPIProductRevisionsCmdLiteral + ` -n PizzaAPI -v 1.0.0 -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPIProductRevisionsCmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory.`

// getRevisionsCmd represents the revisions command
var getAPIProductRevisionsCmd = &cobra.Command{
	Use:     GetAPIProductRevisionsCmdLiteral,
	Short:   GetAPIProductRevisionsCmdShortDesc,
	Long:    GetAPIProductRevisionsCmdLongDesc,
	Example: getRevisionsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetAPIProductRevisionsCmdLiteral + " called")
		cred, err := GetCredentials(getAPIProductRevisionsCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetAPIProductRevisionsCmd(cred)
	},
}

func executeGetAPIProductRevisionsCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getAPIProductRevisionsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'get revisions' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetAPIProductRevisionsCmdLiteral+"'", err)
	}

	_, revisions, err := impl.GetAPIProductRevisionListFromEnv(accessToken, getAPIProductRevisionsCmdEnvironment,
		getRevisionsAPIProductName, getRevisionsAPIProductProvider, getAPIProductRevisionsCmdQuery)
	if err == nil {
		impl.PrintRevisions(revisions, getAPIProductRevisionsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of Revisions", err)
	}
}

func init() {
	GetCmd.AddCommand(getAPIProductRevisionsCmd)
	getAPIProductRevisionsCmd.Flags().StringVarP(&getRevisionsAPIProductName, "name", "n", "",
		"Name of the API Product to get the revision")
	getAPIProductRevisionsCmd.Flags().StringVarP(&getRevisionsAPIProductProvider, "provider", "r", "",
		"Provider of the API Product")
	getAPIProductRevisionsCmd.Flags().StringVarP(&getAPIProductRevisionsCmdQuery, "query", "q",
		"", "Query pattern")
	getAPIProductRevisionsCmd.Flags().StringVarP(&getAPIProductRevisionsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getAPIProductRevisionsCmd.Flags().StringVarP(&getAPIProductRevisionsCmdFormat, "format", "", "", "Pretty-print revisions "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getAPIProductRevisionsCmd.MarkFlagRequired("name")
	_ = getAPIProductRevisionsCmd.MarkFlagRequired("environment")
}
