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

package deactivate

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const deactivateCmdLiteral = "deactivate"
const deactivateCmdShortDesc = "Deactivate artifacts deployed in a Micro Integrator instance"

const deactivateCmdLongDesc = "Deactivate artifacts deployed in a Micro Integrator instance in the environment specified by the flag (--environment, -e)"

const deactivateCmdExamples = utils.ProjectName + " " + utils.MiCmdLiteral + " " + deactivateCmdLiteral + " " + "endpoint" + " TestEP -e dev"

// DeactivateCmd represents the deactivate command
var DeactivateCmd = &cobra.Command{
	Use:     deactivateCmdLiteral,
	Short:   deactivateCmdShortDesc,
	Long:    deactivateCmdLongDesc,
	Example: deactivateCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deactivateCmdLiteral + " called")
		cmd.Help()
	},
}
