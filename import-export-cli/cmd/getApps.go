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

var getAppsCmdEnvironment string
var getAppsCmdAppOwner string
var getAppsCmdFormat string
var getAppsCmdLimit string
var defaultAppsOwner string

// GetAppsCmd related info
const GetAppsCmdLiteral = "apps"
const getAppsCmdShortDesc = "Display a list of Applications in an environment specific to an owner"

const getAppsCmdLongDesc = "Display a list of Applications of the user in the environment specified by the flag --environment, -e"

const getAppsCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAppsCmdLiteral + ` -e dev 
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAppsCmdLiteral + ` -e dev -o sampleUser
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAppsCmdLiteral + ` -e prod -o sampleUser
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAppsCmdLiteral + ` -e staging -o sampleUser
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAppsCmdLiteral + ` -e dev -l 40
NOTE: The flag (--environment (-e)) is mandatory`

// getAppsCmd represents the apps command
var getAppsCmd = &cobra.Command{
	Use:     GetAppsCmdLiteral,
	Short:   getAppsCmdShortDesc,
	Long:    getAppsCmdLongDesc,
	Example: getAppsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetAppsCmdLiteral + " called")
		cred, err := GetCredentials(getAppsCmdEnvironment)
		defaultAppsOwner = cred.Username
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetAppsCmd(cred, getAppsCmdAppOwner)
	},
}

func executeGetAppsCmd(credential credentials.Credential, appOwner string) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getAppsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetAppsCmdLiteral+"'", err)
	}

	_, apps, err := impl.GetApplicationListFromEnv(accessToken, getAppsCmdEnvironment, appOwner, getAppsCmdLimit)

	if err == nil {
		// Printing the list of available Applications
		impl.PrintApps(apps, getAppsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of Applications", err)
	}
}

func init() {
	GetCmd.AddCommand(getAppsCmd)

	getAppsCmd.Flags().StringVarP(&getAppsCmdAppOwner, "owner", "o", defaultAppsOwner,
		"Owner of the Application")
	getAppsCmd.Flags().StringVarP(&getAppsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getAppsCmd.Flags().StringVarP(&getAppsCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultAppsDisplayLimit), "Maximum number of applications to return")
	getAppsCmd.Flags().StringVarP(&getAppsCmdFormat, "format", "", "", "Pretty-print output"+
		"using Go templates. Use \"{{jsonPretty .}}\" to list all fields")
	_ = getAppsCmd.MarkFlagRequired("environment")
}
