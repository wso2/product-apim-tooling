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

var apiProductStateChangeEnvironment string
var apiProductNameForStateChange string
var apiProductVersionForStateChange string
var apiProductProviderForStateChange string
var apiProductStateChangeAction string

// ChangeAPIProductStatus command related usage info
const changeAPIProductStatusCmdLiteral = "api-product"
const changeAPIProductStatusCmdShortDesc = "Change Status of an API Product"
const changeAPIProductStatusCmdLongDesc = "Change the lifecycle status of an API Product in an environment"

const changeAPIProductStatusCmdExamples = utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeAPIProductStatusCmdLiteral + ` -a Publish -n TwitterAPI -r admin -e dev
` + utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeAPIProductStatusCmdLiteral + ` -a Publish -n FacebookAPI -e production
NOTE: The 3 flags (--action (-a), --name (-n), --version (-v) and --environment (-e)) are mandatory.`

// changeAPIProductStatusCmd represents change-status api command
var ChangeAPIProductStatusCmd = &cobra.Command{
	Use: changeAPIProductStatusCmdLiteral + " (--action <action-of-the-api-product-state-change> --name <name-of-the-api-product> --version <version-of-the-api-product> --provider " +
		"<provider-of-the-api-product> --environment <environment-from-which-the-api-product-state-should-be-changed>)",
	Short:   changeAPIProductStatusCmdShortDesc,
	Long:    changeAPIProductStatusCmdLongDesc,
	Example: changeAPIProductStatusCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + changeAPIProductStatusCmdLiteral + " called")
		cred, err := GetCredentials(apiProductStateChangeEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeChangeAPIProductStatusCmd(cred)
	},
}

// executeChangeAPIProductStatusCmd executes the change api product status command
func executeChangeAPIProductStatusCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, apiProductStateChangeEnvironment)
	if preCommandErr == nil {
		resp, err := impl.ChangeAPIProductStatusInEnv(accessToken, apiProductStateChangeEnvironment, apiProductStateChangeAction,
			apiProductNameForStateChange, apiProductVersionForStateChange, apiProductProviderForStateChange)
		if err != nil {
			utils.HandleErrorAndExit("Error while changing the API Product status", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			// 200 OK
			fmt.Println(apiNameForStateChange + " API Product state changed successfully!")
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// Neither 200 nor 500
			fmt.Println("Error while changing API Product Status: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// Error changing the API Product status
		fmt.Println("Error getting OAuth tokens while changing status of the API Product : " + preCommandErr.Error())
	}
}

func init() {
	ChangeStatusCmd.AddCommand(ChangeAPIProductStatusCmd)
	ChangeAPIProductStatusCmd.Flags().StringVarP(&apiProductStateChangeAction, "action", "a", "",
		"Action to be taken to change the status of the API Product")
	ChangeAPIProductStatusCmd.Flags().StringVarP(&apiProductNameForStateChange, "name", "n", "",
		"Name of the API Product to be state changed")
	ChangeAPIProductStatusCmd.Flags().StringVarP(&apiProductVersionForStateChange, "version", "v", "",
		"Version of the API Product to be state changed")
	ChangeAPIProductStatusCmd.Flags().StringVarP(&apiProductProviderForStateChange, "provider", "r", "",
		"Provider of the API Product")
	ChangeAPIProductStatusCmd.Flags().StringVarP(&apiProductStateChangeEnvironment, "environment", "e",
		"", "Environment of which the API Product state should be changed")
	// Mark required flags
	_ = ChangeAPIProductStatusCmd.MarkFlagRequired("action")
	_ = ChangeAPIProductStatusCmd.MarkFlagRequired("name")
	_ = ChangeAPIProductStatusCmd.MarkFlagRequired("version")
	_ = ChangeAPIProductStatusCmd.MarkFlagRequired("environment")
}
