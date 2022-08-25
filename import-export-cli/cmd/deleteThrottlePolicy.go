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

var deleteThrottlingPolicyEnvironment string
var deleteThrottlingPolicyName string
var deleteThrottlingPolicyType string

// DeleteThrottlingPolicy command related usage info
const DeleteThrottlingPolicyCmdLiteral = "rate-limiting"
const DeleteThrottlingPolicyCmdShortDesc = "Delete Throttling Policy"
const DeleteThrottlingPolicyCmdLongDesc = "Delete a throttling policy from an environment"

const DeleteThrottlingPolicyCmdExamplesDefault = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + DeletePolicyCmdLiteral + ` ` + DeleteThrottlingPolicyCmdLiteral + ` -n Gold -e dev --type sub 
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + DeletePolicyCmdLiteral + ` ` + DeleteThrottlingPolicyCmdLiteral + ` -n AppPolicy -e prod --type app
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + DeletePolicyCmdLiteral + ` ` + DeleteThrottlingPolicyCmdLiteral + ` -n TestPolicy -e dev --type advanced 
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + DeletePolicyCmdLiteral + ` ` + DeleteThrottlingPolicyCmdLiteral + ` -n CustomPolicy -e prod --type custom 
NOTE: All the 2 flags (--name (-n) and --environment (-e)) are mandatory.`

// DeleteThrottlingPolicyCmd represents the delete Throttling policy command
var DeleteThrottlingPolicyCmd = &cobra.Command{
	Use: DeleteThrottlingPolicyCmdLiteral + " (--name <name-of-the-throttling-policy> --environment " +
		"<environment-from-which-the-policy-should-be-deleted>)" +
		"--type <type-of-the-throttling-policy>",
	Short:   DeleteThrottlingPolicyCmdShortDesc,
	Long:    DeleteThrottlingPolicyCmdLongDesc,
	Example: DeleteThrottlingPolicyCmdExamplesDefault,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + DeleteThrottlingPolicyCmdLiteral + " called")

		cred, err := GetCredentials(deleteThrottlingPolicyEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeDeleteThrottlingPolicyCmd(cred)
	},
}

// executeDeleteThrottlingPolicyCmd executes the delete Throttling policy command
func executeDeleteThrottlingPolicyCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteThrottlingPolicyEnvironment)
	if preCommandErr == nil {
		_, err := impl.DeleteThrottlingPolicy(accessToken, deleteThrottlingPolicyName, deleteThrottlingPolicyType, deleteThrottlingPolicyEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting Throttling Policy ", err)
		}
		impl.PrintDeleteThrottlingPolicyResponse(deleteThrottlingPolicyName, deleteThrottlingPolicyType, err)
	} else {
		// Error deleting Throttling Policy
		fmt.Println("Error getting OAuth tokens while deleting Throttling Policy:" + preCommandErr.Error())
	}
}

// Init using Cobra
func init() {
	DeletePolicyCmd.AddCommand(DeleteThrottlingPolicyCmd)
	DeleteThrottlingPolicyCmd.Flags().StringVarP(&deleteThrottlingPolicyName, "name", "n", "",
		"Name of the Throttling Policy to be deleted")
	DeleteThrottlingPolicyCmd.Flags().StringVarP(&deleteThrottlingPolicyEnvironment, "environment", "e",
		"", "Environment from which the Throttling Policy should be deleted")
	DeleteThrottlingPolicyCmd.Flags().StringVarP(&deleteThrottlingPolicyType, "type", "t",
		"", "Type of the Throttling Policies to be exported (sub,app,custom,advanced)")
	_ = DeleteThrottlingPolicyCmd.MarkFlagRequired("name")
	_ = DeleteThrottlingPolicyCmd.MarkFlagRequired("environment")
	_ = DeleteThrottlingPolicyCmd.MarkFlagRequired("type")
}
