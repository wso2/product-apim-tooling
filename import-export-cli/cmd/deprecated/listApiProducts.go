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

package deprecated

import (
	"strconv"

	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var listApiProductsCmdEnvironment string
var listApiProductsCmdFormat string
var listApiProductsCmdQuery string
var listApiProductsCmdLimit string

// apiProductsCmd related info
const apiProductsCmdLiteral = "api-products"
const apiProductsCmdShortDesc = "Display a list of API Products in an environment"

const apiProductsCmdLongDesc = `Display a list of API Products in the environment specified by the flag --environment, -e`

var apiProductsCmdExamples = utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e dev -q provider:devops
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e prod -q provider:admin context:/myproduct
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e prod -l 25
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e staging
NOTE: The flag (--environment (-e)) is mandatory`

// apiProductsCmd represents the api-products command
var apiProductsCmdDeprecated = &cobra.Command{
	Use:        apiProductsCmdLiteral,
	Short:      apiProductsCmdShortDesc,
	Long:       apiProductsCmdLongDesc,
	Example:    apiProductsCmdExamples,
	Deprecated: "instead use \"" + cmd.GetCmdLiteral + " " + cmd.GetApiProductsCmdLiteral + "\".",
	Run: func(deprecatedCmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apiProductsCmdLiteral + " called")
		cred, err := cmd.GetCredentials(listApiProductsCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		//Since other flags does not use args[], query flag will own this
		if len(args) != 0 && listApiProductsCmdQuery != "" {
			for _, argument := range args {
				listApiProductsCmdQuery += " " + argument
			}
		}
		executeApiProductsCmd(cred)
	},
}

func executeApiProductsCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, listApiProductsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+apiProductsCmdLiteral+"'", err)
	}

	// Unified Search endpoint from the config file to search API Products
	_, apiProducts, err := impl.GetAPIProductListFromEnv(accessToken, listApiProductsCmdEnvironment, listApiProductsCmdQuery,
		listApiProductsCmdLimit)
	if err == nil {
		impl.PrintAPIProducts(apiProducts, listApiProductsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of API Products", err)
	}
}

func init() {
	ListCmdDeprecated.AddCommand(apiProductsCmdDeprecated)

	apiProductsCmdDeprecated.Flags().StringVarP(&listApiProductsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	apiProductsCmdDeprecated.Flags().StringVarP(&listApiProductsCmdQuery, "query", "q",
		"", "Query pattern")
	apiProductsCmdDeprecated.Flags().StringVarP(&listApiProductsCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultApiProductsDisplayLimit), "Maximum number of API Products to return")
	apiProductsCmdDeprecated.Flags().StringVarP(&listApiProductsCmdFormat, "format", "", "", "Pretty-print API Products "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = apiProductsCmdDeprecated.MarkFlagRequired("environment")
}
