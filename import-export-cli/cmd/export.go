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
const exportCmdLiteral = "export"
const exportCmdShortDesc = "Export an API Product in an environment"

const exportCmdLongDesc = `Export an API Product available in the environment specified by flag (--environment, -e)`

const exportCmdExamples = utils.ProjectName + ` ` + exportCmdLiteral + ` ` + exportAPIProductCmdLiteral + ` -n LeasingAPIProduct -e dev`

// ExportCmd represents the export command
var ExportCmd = &cobra.Command{
	Use:     exportCmdLiteral,
	Short:   exportCmdShortDesc,
	Long:    exportCmdLongDesc,
	Example: exportCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ExportCmd)
}
