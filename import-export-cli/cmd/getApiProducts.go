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
	"strconv"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getApiProductsCmdEnvironment string
var getApiProductsCmdFormat string
var getApiProductsCmdQuery string
var getApiProductsCmdLimit string

// GetApiProductsCmd related info
const GetApiProductsCmdLiteral = "api-products"
const getApiProductsCmdShortDesc = "Display a list of API Products in an environment"

const getApiProductsCmdLongDesc = `Display a list of API Products in the environment specified by the flag --environment, -e`

var getApiProductsCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiProductsCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiProductsCmdLiteral + ` -e dev -q provider:devops
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiProductsCmdLiteral + ` -e prod -q provider:admin context:/myproduct
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiProductsCmdLiteral + ` -e prod -l 25
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiProductsCmdLiteral + ` -e staging
NOTE: The flag (--environment (-e)) is mandatory`

// getApiProductsCmd represents the api-products command
var getApiProductsCmd = &cobra.Command{
	Use:     GetApiProductsCmdLiteral,
	Short:   getApiProductsCmdShortDesc,
	Long:    getApiProductsCmdLongDesc,
	Example: getApiProductsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetApiProductsCmdLiteral + " called")
		cred, err := GetCredentials(getApiProductsCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		//Since other flags does not use args[], query flag will own this
		if len(args) != 0 && getApiProductsCmdQuery != "" {
			for _, argument := range args {
				getApiProductsCmdQuery += " " + argument
			}
		}
		executeGetApiProductsCmd(cred)
	},
}

func executeGetApiProductsCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getApiProductsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetApiProductsCmdLiteral+"'", err)
	}

	// Unified Search endpoint from the config file to search API Products
	_, apiProducts, err := impl.GetAPIProductListFromEnv(accessToken, getApiProductsCmdEnvironment, getApiProductsCmdQuery,
		getApiProductsCmdLimit)
	if err == nil {
		impl.PrintAPIProducts(apiProducts, getApiProductsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of API Products", err)
	}
}

func init() {
	GetCmd.AddCommand(getApiProductsCmd)

	getApiProductsCmd.Flags().StringVarP(&getApiProductsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getApiProductsCmd.Flags().StringVarP(&getApiProductsCmdQuery, "query", "q",
		"", "Query pattern")
	getApiProductsCmd.Flags().StringVarP(&getApiProductsCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultApiProductsDisplayLimit), "Maximum number of API Products to return")
	getApiProductsCmd.Flags().StringVarP(&getApiProductsCmdFormat, "format", "", "", "Pretty-print API Products "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getApiProductsCmd.MarkFlagRequired("environment")
}
