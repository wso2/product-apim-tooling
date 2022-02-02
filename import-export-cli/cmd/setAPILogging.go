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
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var setApiLoggingEnvironment string
var setApiLoggingAPIId string
var setApiLoggingTenantDomain string
var setApiLoggingLogLevel string

const SetApiLoggingCmdLiteral = "api-logging"
const setApiLoggingCmdShortDesc = "Set the log level for an API in an environment"
const setApiLoggingCmdLongDesc = `Set the log level for an API in the environment specified`

var setApiLoggingCmdExamples = utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetApiLoggingCmdLiteral + ` --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level full -e dev --tenant-domain carbon.super
` + utils.ProjectName + ` ` + SetCmdLiteral + ` ` + SetApiLoggingCmdLiteral + ` --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level off -e dev --tenant-domain carbon.super`

var setApiLoggingCmd = &cobra.Command{
	Use:     SetApiLoggingCmdLiteral,
	Short:   setApiLoggingCmdShortDesc,
	Long:    setApiLoggingCmdLongDesc,
	Example: setApiLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + SetCmdLiteral + " " + SetApiLoggingCmdLiteral + " called")
		cred, err := GetCredentials(setApiLoggingEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeSetApiLoggingCmd(cred)
	},
}

func executeSetApiLoggingCmd(credential credentials.Credential) {
	resp, err := impl.SetAPILoggingLevel(credential, setApiLoggingEnvironment, setApiLoggingAPIId, setApiLoggingTenantDomain, setApiLoggingLogLevel)
	if err != nil {
		utils.HandleErrorAndExit("Error while setting the log level of the API", err)
	}
	// Print info on response
	utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
	if resp.StatusCode() == http.StatusOK {
		// 200 OK
		fmt.Println("Log level " + setApiLoggingLogLevel + " is successfully set to the API.")
	} else {
		fmt.Println("Setting the log level of the API: ", resp.Status(), "\n", string(resp.Body()))
	}
}

func init() {
	SetCmd.AddCommand(setApiLoggingCmd)

	setApiLoggingCmd.Flags().StringVarP(&setApiLoggingAPIId, "api-id", "i",
		"", "API ID")
	setApiLoggingCmd.Flags().StringVarP(&setApiLoggingTenantDomain, "tenant-domain", "",
		"", "Tenant Domain")
	setApiLoggingCmd.Flags().StringVarP(&setApiLoggingLogLevel, "log-level", "",
		"", "Log Level")
	setApiLoggingCmd.Flags().StringVarP(&setApiLoggingEnvironment, "environment", "e",
		"", "Environment of the API which the log level should be set")
	_ = setApiLoggingCmd.MarkFlagRequired("environment")
	_ = setApiLoggingCmd.MarkFlagRequired("api-id")
	_ = setApiLoggingCmd.MarkFlagRequired("log-level")
}
