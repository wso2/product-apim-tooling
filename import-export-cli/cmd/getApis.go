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

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getApisCmdEnvironment string
var getApisCmdFormat string
var getApisCmdQuery string
var getApisCmdLimit string

// GetApisCmd related info
const GetApisCmdLiteral = "apis"
const getApisCmdShortDesc = "Display a list of APIs in an environment"

const getApisCmdLongDesc = `Display a list of APIs in the environment specified by the flag --environment, -e`

var getApisCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApisCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApisCmdLiteral + ` -e dev -q version:1.0.0
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApisCmdLiteral + ` -e prod -q provider:admin
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApisCmdLiteral + ` -e prod -l 100
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApisCmdLiteral + ` -e staging
NOTE: The flag (--environment (-e)) is mandatory`

// getApisCmd represents the apis command
var getApisCmd = &cobra.Command{
	Use:     GetApisCmdLiteral,
	Short:   getApisCmdShortDesc,
	Long:    getApisCmdLongDesc,
	Example: getApisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetApisCmdLiteral + " called")
		cred, err := GetCredentials(getApisCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		//Since other flags does not use args[], query flag will own this
		if len(args) != 0 && getApisCmdQuery != "" {
			for _, argument := range args {
				getApisCmdQuery += " " + argument
			}
		}
		executeGetApisCmd(cred)
	},
}

func executeGetApisCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getApisCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetApisCmdLiteral+"'", err)
	}

	_, apis, err := impl.GetAPIListFromEnv(accessToken, getApisCmdEnvironment, getApisCmdQuery, getApisCmdLimit)
	if err == nil {
		impl.PrintAPIs(apis, getApisCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
	}
}

func init() {
	GetCmd.AddCommand(getApisCmd)

	getApisCmd.Flags().StringVarP(&getApisCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getApisCmd.Flags().StringVarP(&getApisCmdQuery, "query", "q",
		"", "Query pattern")
	getApisCmd.Flags().StringVarP(&getApisCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultApisDisplayLimit), "Maximum number of apis to return")
	getApisCmd.Flags().StringVarP(&getApisCmdFormat, "format", "", "", "Pretty-print apis "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getApisCmd.MarkFlagRequired("environment")
}
