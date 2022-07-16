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
)

// Delete policy command related usage Info
const DeletePolicyCmdLiteral = "policy"
const DeletePolicyCmdShortDesc = "Delete a Policy"
const DeletePolicyCmdLongDesc = "Delete a Policy in an environment"
const DeletePolicyCmdExamples = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + DeletePolicyCmdLiteral + ` ` + DeleteAPIPolicyCmdLiteral + ` -n addHeader -e prod`

// DeletePolicyCmd represents the export policy command
var DeletePolicyCmd = &cobra.Command{
	Use:     DeletePolicyCmdLiteral,
	Short:   DeletePolicyCmdShortDesc,
	Long:    DeletePolicyCmdLongDesc,
	Example: DeletePolicyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + DeletePolicyCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	DeleteCmd.AddCommand(DeletePolicyCmd)

}
