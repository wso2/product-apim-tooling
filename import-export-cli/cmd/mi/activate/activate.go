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

package activate

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const activateCmdLiteral = "activate"
const activateCmdShortDesc = "Activate artifacts deployed in a Micro Integrator instance"

const activateCmdLongDesc = "Activate artifacts deployed in a Micro Integrator instance in the environment specified by the flag (--environment, -e)"

const activateCmdExamples = utils.ProjectName + " " + utils.MiCmdLiteral + " " + activateCmdLiteral + " " + "endpoint" + " TestEP -e dev"

// ActivateCmd represents the activate command
var ActivateCmd = &cobra.Command{
	Use:     activateCmdLiteral,
	Short:   activateCmdShortDesc,
	Long:    activateCmdLongDesc,
	Example: activateCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + activateCmdLiteral + " called")
		cmd.Help()
	},
}
