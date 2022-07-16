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

const ImportPolicyCmdLiteral = "policy"
const ImportPolicyCmdShortDesc = "Import a Policy"
const ImportPolicyCmdLongDesc = "Import a Policy in an environment or Import a Policy to an environment"
const ImportPolicyCmdExamples = utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportPolicyCmdLiteral + ` ` + ImportThrottlingPolicyCmdLiteral + ` -f ~/CustomPolicy -e production --u
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportPolicyCmdLiteral + ` ` + ImportAPIPolicyCmdLiteral + ` ` + ` -f ~/AddHeader -e production`

// ImportPolicyCmd represents the import command for a policy
var ImportPolicyCmd = &cobra.Command{
	Use:     ImportPolicyCmdLiteral,
	Short:   ImportPolicyCmdShortDesc,
	Long:    ImportPolicyCmdLongDesc,
	Example: ImportPolicyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	ImportCmd.AddCommand(ImportPolicyCmd)
}
