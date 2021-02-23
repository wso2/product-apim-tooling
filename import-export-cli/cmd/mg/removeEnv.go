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

package mg

import (
	"fmt"

	"github.com/spf13/cobra"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/impl/mg"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var envToBeRemoved string // name of the environment to be removed

const (
	removeEnvCmdShortDesc = "Remove an environment for the Microgateway Adapter(s)"
	removeEnvCmdLongDesc  = "Remove Environment and its configurations for the Microgateway " +
		"Adapter(s) from the config file."
)

const removeEnvCmdExamples = utils.ProjectName + " " + removeCmdLiteral + " " +
	envCmdLiteral + " prod"

// RemoveEnvCmd represents the removeEnv command
var RemoveEnvCmd = &cobra.Command{
	Use:     envCmdLiteral,
	Short:   removeEnvCmdShortDesc,
	Long:    removeEnvCmdLongDesc,
	Example: removeEnvCmdExamples,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envToBeRemoved := args[0]

		utils.Logln(utils.LogPrefixInfo + envCmdLiteral + " called")
		executeRemoveEnvCmd(envToBeRemoved, utils.MainConfigFilePath)
	},
}

func executeRemoveEnvCmd(env, mainConfigFilePath string) {
	err := impl.RemoveEnv(env, mainConfigFilePath)
	if err != nil {
		utils.HandleErrorAndExit("Error occurred when removing environment", err)
	}
	fmt.Println("Successfully removed environment '" + env + "'")
	fmt.Println("Execute '" + utils.ProjectName + " " + mgCmdLiteral + " " +
		addCmdLiteral + " " + envCmdLiteral + " --help' to see how to add a new environment")

}

// init using Cobra
func init() {
	RemoveCmd.AddCommand(RemoveEnvCmd)
}
