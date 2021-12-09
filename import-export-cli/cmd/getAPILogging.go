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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
)

var getApiLoggingEnvironment string
var getApiLoggingAPIId string
var getAPILoggingCmdFormat string

const GetApiLoggingCmdLiteral = "api-logging"
const getApiLoggingCmdShortDesc = "Display a list of API loggers in an environment"
const getApiLoggingCmdLongDesc = `Display a list of API loggers available for the APIs in the environment specified`

var getApiLoggingCmdExamples =
	utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiLoggingCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiLoggingCmdLiteral + ` --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 -e dev`

var getApiLoggingCmd = &cobra.Command{
	Use:     GetApiLoggingCmdLiteral,
	Short:   getApiLoggingCmdShortDesc,
	Long:    getApiLoggingCmdLongDesc,
	Example: getApiLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetCmdLiteral + " " + GetApiLoggingCmdLiteral + " called")
		cred, err := GetCredentials(getApiLoggingEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetApiLoggingCmd(cred)
	},
}

func executeGetApiLoggingCmd(credential credentials.Credential) {
	if getApiLoggingAPIId != "" {
		api, err := impl.GetPerAPILoggingDetailsFromEnv(credential, getApiLoggingEnvironment, getApiLoggingAPIId)
		if err == nil {
			impl.PrintAPILoggers(api, getAPILoggingCmdFormat)
		} else {
			utils.Logln(utils.LogPrefixError+"Getting API logger details of the API", err)
		}
	} else {
		apis, err := impl.GetPerAPILoggingListFromEnv(credential, getApiLoggingEnvironment)
		if err == nil {
			impl.PrintAPILoggers(apis, getAPILoggingCmdFormat)
		} else {
			utils.Logln(utils.LogPrefixError+"Getting List of API loggers for the APIs", err)
		}
	}
}

func init() {
	GetCmd.AddCommand(getApiLoggingCmd)

	getApiLoggingCmd.Flags().StringVarP(&getApiLoggingAPIId, "api-id", "i",
		"", "API ID")
	getApiLoggingCmd.Flags().StringVarP(&getApiLoggingEnvironment, "environment", "e",
		"", "Environment of the APIs which the API loggers should be displayed")
	getApiLoggingCmd.Flags().StringVarP(&getAPILoggingCmdFormat, "format", "", "", "Pretty-print API loggers "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getApiLoggingCmd.MarkFlagRequired("environment")
}
