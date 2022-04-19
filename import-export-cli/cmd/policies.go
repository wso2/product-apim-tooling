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

// Export command related usage Info
const policyCmdLiteral = "policy"
const policyCmdShortDesc = "Export/Import a Policy"

const policyCmdLongDesc = "Export a Policy in an environment or Import a Policy to an environment"

const policyCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + policyCmdLiteral + ` ` + ExportThrottlePolicyCmdLiteral + ` -n Silver -e prod --type subscription`

// ExportCmd represents the export command
var policyCmd = &cobra.Command{
	Use:     policyCmdLiteral,
	Short:   policyCmdShortDesc,
	Long:    policyCmdLongDesc,
	Example: policyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	ExportCmd.AddCommand(policyCmd)
}
