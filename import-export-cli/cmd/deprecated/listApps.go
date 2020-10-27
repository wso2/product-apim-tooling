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

var listAppsCmdEnvironment string
var listAppsCmdAppOwner string
var listAppsCmdFormat string
var listAppsCmdLimit string
var defaultAppsOwner string

// appsCmd related info
const appsCmdLiteral = "apps"
const appsCmdShortDesc = "Display a list of Applications in an environment specific to an owner"

const appsCmdLongDesc = "Display a list of Applications of the user in the environment specified by the flag --environment, -e"

const appsCmdExamples = utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e dev 
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e dev -o sampleUser
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e prod -o sampleUser
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e staging -o sampleUser
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e dev -l 40
NOTE: The flag (--environment (-e)) is mandatory`

// appsCmd represents the apps command
var appsCmdDeprected = &cobra.Command{
	Use:        appsCmdLiteral,
	Short:      appsCmdShortDesc,
	Long:       appsCmdLongDesc,
	Example:    appsCmdExamples,
	Deprecated: "use \"" + cmd.GetCmdLiteral + " " + cmd.GetAppsCmdLiteral + "\" " + "instead of \"" + listCmdLiteral + " " + appsCmdLiteral + "\".",
	Run: func(depcrecatedCmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + appsCmdLiteral + " called")
		cred, err := cmd.GetCredentials(listAppsCmdEnvironment)
		defaultAppsOwner = cred.Username
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeAppsCmd(cred, listAppsCmdAppOwner)
	},
}

func executeAppsCmd(credential credentials.Credential, appOwner string) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, listAppsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+appsCmdLiteral+"'", err)
	}

	_, apps, err := impl.GetApplicationListFromEnv(accessToken, listAppsCmdEnvironment, appOwner, listAppsCmdLimit)

	if err == nil {
		// Printing the list of available Applications
		impl.PrintApps(apps, listAppsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of Applications", err)
	}
}

func init() {
	ListCmdDeprecated.AddCommand(appsCmdDeprected)

	appsCmdDeprected.Flags().StringVarP(&listAppsCmdAppOwner, "owner", "o", defaultAppsOwner,
		"Owner of the Application")
	appsCmdDeprected.Flags().StringVarP(&listAppsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	appsCmdDeprected.Flags().StringVarP(&listAppsCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultAppsDisplayLimit), "Maximum number of applications to return")
	appsCmdDeprected.Flags().StringVarP(&listAppsCmdFormat, "format", "", "", "Pretty-print output"+
		"using Go templates. Use \"{{jsonPretty .}}\" to list all fields")
	_ = appsCmdDeprected.MarkFlagRequired("environment")
}
