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

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var deleteAPIProductEnvironment string
var deleteAPIProductName string
var deleteAPIProductProvider string

// DeleteAPIProduct command related usage info
const deleteAPIProductCmdLiteral = "api-product"
const deleteAPIProductCmdShortDesc = "Delete API Product"
const deleteAPIProductCmdLongDesc = "Delete an API Product from an environment"

const deleteAPIProductCmdExamples = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPIProductCmdLiteral + ` -n LeasingAPIProduct -r admin -e dev
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPIProductCmdLiteral + ` -n CreditAPIProduct -e production
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory.`

// TODO Introduce a version flag and mandate it when the versioning support has been implemented for API Products

// DeleteAPIProductCmd represents the delete api-product command
var DeleteAPIProductCmd = &cobra.Command{
	Use: deleteAPIProductCmdLiteral + " (--name <name-of-the-api-product> --provider <provider-of-the-api-product> --environment " +
		"<environment-from-which-the-api-product-should-be-deleted>)",
	Short:   deleteAPIProductCmdShortDesc,
	Long:    deleteAPIProductCmdLongDesc,
	Example: deleteAPIProductCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAPIProductCmdLiteral + " called")
		cred, err := GetCredentials(deleteAPIProductEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeDeleteAPIProductCmd(cred)
	},
}

// executeDeleteAPIProductCmd executes the delete api command
func executeDeleteAPIProductCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteAPIProductEnvironment)
	if preCommandErr == nil {
		resp, err := impl.DeleteAPIProduct(accessToken, deleteAPIProductEnvironment, deleteAPIProductName, deleteAPIProductProvider)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting API Product", err)
		}
		impl.PrintDeleteAPIProductResponse(resp, err)
	} else {
		// Error deleting API Product
		fmt.Println("Error getting OAuth tokens while deleting API Product:" + preCommandErr.Error())
	}
}

// Init using Cobra
func init() {
	DeleteCmd.AddCommand(DeleteAPIProductCmd)
	DeleteAPIProductCmd.Flags().StringVarP(&deleteAPIProductName, "name", "n", "",
		"Name of the API Product to be deleted")
	DeleteAPIProductCmd.Flags().StringVarP(&deleteAPIProductProvider, "provider", "r", "",
		"Provider of the API Product to be deleted")
	DeleteAPIProductCmd.Flags().StringVarP(&deleteAPIProductEnvironment, "environment", "e",
		"", "Environment from which the API Product should be deleted")
	// Mark required flags
	_ = DeleteAPIProductCmd.MarkFlagRequired("name")
	_ = DeleteAPIProductCmd.MarkFlagRequired("environment")
}
