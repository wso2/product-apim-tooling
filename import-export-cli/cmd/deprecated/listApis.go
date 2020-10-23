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
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	apiIdHeader       = "ID"
	apiNameHeader     = "NAME"
	apiContextHeader  = "CONTEXT"
	apiVersionHeader  = "VERSION"
	apiProviderHeader = "PROVIDER"
	apiStatusHeader   = "STATUS"

	defaultApiTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Version}}\t{{.Context}}\t{{.LifeCycleStatus}}\t{{.Provider}}"
)

var listApisCmdEnvironment string
var listApisCmdFormat string
var listApisCmdQuery string
var listApisCmdLimit string
var queryParamAdded bool = false

// apisCmd related info
const apisCmdLiteral = "apis"
const apisCmdShortDesc = "Display a list of APIs in an environment"

const apisCmdLongDesc = `Display a list of APIs in the environment specified by the flag --environment, -e`

var apisCmdExamples = utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e dev -q version:1.0.0
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e prod -q provider:admin
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e prod -l 100
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e staging
NOTE: The flag (--environment (-e)) is mandatory`

// apisCmd represents the apis command
var apisCmdDeprecated = &cobra.Command{
	Use:        apisCmdLiteral,
	Short:      apisCmdShortDesc,
	Long:       apisCmdLongDesc,
	Example:    apisCmdExamples,
	Deprecated: "instead use \"" + cmd.GetCmdLiteral + " " + cmd.GetApisCmdLiteral + "\".",
	Run: func(deprecatedCmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apisCmdLiteral + " called")
		cred, err := cmd.GetCredentials(listApisCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		//Since other flags does not use args[], query flag will own this
		if len(args) != 0 && listApisCmdQuery != "" {
			for _, argument := range args {
				listApisCmdQuery += " " + argument
			}
		}
		executeApisCmd(cred)
	},
}

func executeApisCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, listApisCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+apisCmdLiteral+"'", err)
	}

	_, apis, err := impl.GetAPIListFromEnv(accessToken, listApisCmdEnvironment, listApisCmdQuery, listApisCmdLimit)
	if err == nil {
		impl.PrintAPIs(apis, listApisCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
	}
}

func init() {
	ListCmdDeprecated.AddCommand(apisCmdDeprecated)

	apisCmdDeprecated.Flags().StringVarP(&listApisCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	apisCmdDeprecated.Flags().StringVarP(&listApisCmdQuery, "query", "q",
		"", "Query pattern")
	apisCmdDeprecated.Flags().StringVarP(&listApisCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultApisDisplayLimit), "Maximum number of apis to return")
	apisCmdDeprecated.Flags().StringVarP(&listApisCmdFormat, "format", "", "", "Pretty-print apis "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = apisCmdDeprecated.MarkFlagRequired("environment")
}
