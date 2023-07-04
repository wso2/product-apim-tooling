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

package delete

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const deleteCmdLiteral = "delete"
const deleteCmdShortDesc = "Delete users from a Micro Integrator instance"

const deleteCmdLongDesc = "Delete users from a Micro Integrator instance in the environment specified by the flag (--environment, -e)"

var deleteCmdExamples = utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + deleteCmdLiteral + " " + "user" + " capp-tester -e dev"

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:     deleteCmdLiteral,
	Short:   deleteCmdShortDesc,
	Long:    deleteCmdLongDesc,
	Example: deleteCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteCmdLiteral + " called")
		cmd.Help()
	},
}
