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

var deleteOperationPolicyEnvironment string
var deleteOperationPolicyName string
var deleteOperationPolicyVersion string

// DeleteOperationPolicy command related usage info
const DeleteOperationPolicyCmdLiteral = "operation"
const DeleteOperationPolicyCmdShortDesc = "Delete Operation Policy"
const DeleteOperationPolicyCmdLongDesc = "Delete an operation policy from an environment"

const DeleteOperationPolicyCmdExamplesDefault = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + DeletePolicyCmdLiteral + ` ` + DeleteOperationPolicyCmdLiteral + ` -n addHeader -e dev
 NOTE: The 2 flags (--name (-n) and --environment (-e)) are mandatory.`

// DeleteOperationPolicyCmd represents the delete api command
var DeleteOperationPolicyCmd = &cobra.Command{
	Use: DeleteOperationPolicyCmdLiteral + " (--name <name-of-the-operation-policy> --environment " +
		"<environment-from-which-the-policy-should-be-deleted>)",
	Short:   DeleteOperationPolicyCmdShortDesc,
	Long:    DeleteOperationPolicyCmdLongDesc,
	Example: DeleteOperationPolicyCmdExamplesDefault,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + DeleteOperationPolicyCmdLiteral + " called")

		cred, err := GetCredentials(deleteOperationPolicyEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeDeleteOperationPolicyCmd(cred)

	},
}

// executeDeleteOperationPolicyCmd executes the delete operation policy command
func executeDeleteOperationPolicyCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteOperationPolicyEnvironment)
	if preCommandErr == nil {
		deleteOperationPolicyVersion = utils.OperationPolicyVersion
		resp, err := impl.DeleteOperationPolicy(accessToken, deleteOperationPolicyName, deleteOperationPolicyVersion, deleteOperationPolicyEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting Operation Policy ", err)
		}
		impl.PrintDeleteOperationPolicyResponse(resp, err)
	} else {
		// Error deleting Operation Policy
		fmt.Println("Error getting OAuth tokens while deleting Operation Policy:" + preCommandErr.Error())
	}
}

// Init using Cobra
func init() {
	DeletePolicyCmd.AddCommand(DeleteOperationPolicyCmd)
	DeleteOperationPolicyCmd.Flags().StringVarP(&deleteOperationPolicyName, "name", "n", "",
		"Name of the Operation Policy to be deleted")
	DeleteOperationPolicyCmd.Flags().StringVarP(&deleteOperationPolicyEnvironment, "environment", "e",
		"", "Environment from which the Operation Policy should be deleted")

	_ = DeleteOperationPolicyCmd.MarkFlagRequired("name")
	_ = DeleteOperationPolicyCmd.MarkFlagRequired("environment")
}
