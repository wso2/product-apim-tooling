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

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var deleteAPIPolicyEnvironment string
var deleteAPIPolicyName string
var deleteAPIPolicyVersion string

// DeleteAPIPolicy command related usage info
const DeleteAPIPolicyCmdLiteral = "api"
const DeleteAPIPolicyCmdShortDesc = "Delete an API Policy"
const DeleteAPIPolicyCmdLongDesc = "Delete an API Policy from an environment"

const DeleteAPIPolicyCmdExamplesDefault = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + DeletePolicyCmdLiteral + ` ` + DeleteAPIPolicyCmdLiteral + ` -n addHeader -e dev
 NOTE: The 2 flags (--name (-n) and --environment (-e)) are mandatory.`

// DeleteAPIPolicyCmd represents the delete api policy command
var DeleteAPIPolicyCmd = &cobra.Command{
	Use: DeleteAPIPolicyCmdLiteral + " (--name <name-of-the-api-policy> --environment " +
		"<environment-from-which-the-policy-should-be-deleted>)",
	Short:   DeleteAPIPolicyCmdShortDesc,
	Long:    DeleteAPIPolicyCmdLongDesc,
	Example: DeleteAPIPolicyCmdExamplesDefault,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + DeleteAPIPolicyCmdLiteral + " called")

		cred, err := GetCredentials(deleteAPIPolicyEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		deleteAPIPolicyVersion = utils.DefaultAPIPolicyVersion
		executeDeleteAPIPolicyCmd(cred)

	},
}

// executeDeleteAPIPolicyCmd executes the delete api policy command
func executeDeleteAPIPolicyCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteAPIPolicyEnvironment)
	if preCommandErr == nil {
		_, err := impl.DeleteAPIPolicy(accessToken, deleteAPIPolicyName, deleteAPIPolicyVersion, deleteAPIPolicyEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting API Policy ", err)
		}
		impl.PrintDeleteAPIPolicyResponse(deleteAPIPolicyName, err)
	} else {
		// Error deleting API Policy
		fmt.Println("Error getting OAuth tokens while deleting API Policy:" + preCommandErr.Error())
	}
}

// Init using Cobra
func init() {
	DeletePolicyCmd.AddCommand(DeleteAPIPolicyCmd)
	DeleteAPIPolicyCmd.Flags().StringVarP(&deleteAPIPolicyName, "name", "n", "",
		"Name of the API Policy to be deleted")
	DeleteAPIPolicyCmd.Flags().StringVarP(&deleteAPIPolicyEnvironment, "environment", "e",
		"", "Environment from which the API Policy should be deleted")

	_ = DeleteAPIPolicyCmd.MarkFlagRequired("name")
	_ = DeleteAPIPolicyCmd.MarkFlagRequired("environment")
}
