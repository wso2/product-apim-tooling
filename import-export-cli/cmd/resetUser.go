/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"github.com/renstrom/dedent"
)

var resetUserEnvironment string

// ResetUser command related usage Info

const ResetUserCmdLiteral string = "reset-user"
const ResetUserCmdShortDesc string = "Reset user of an environment"

var ResetUserCmdLongDesc = dedent.Dedent(`
		Reset user data of a particular environment (Clear the entry in env_keys_all.yaml file)
	`)

var ResetUserCmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + ResetUserCmdLiteral + ` -e dev
		` + utils.ProjectName + ` ` + ResetUserCmdLiteral + `reset-user -e staging
	`)



// ResetUserCmd represents the resetUser command
var ResetUserCmd = &cobra.Command{
	Use:   ResetUserCmdLiteral,
	Short: ResetUserCmdShortDesc,
	Long:  ResetUserCmdLongDesc + ResetUserCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "reset-user called")

		err := utils.RemoveEnvFromKeysFile(resetUserEnvironment, utils.EnvKeysAllFilePath, utils.MainConfigFilePath)
		if err != nil {
			utils.HandleErrorAndExit("Error clearing user data for environment "+resetUserEnvironment, err)
		} else {
			fmt.Println("Successfully cleared user data for environment: " + resetUserEnvironment)
		}
	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ResetUserCmd)
	ResetUserCmd.Flags().StringVarP(&resetUserEnvironment, "environment", "e",
		utils.GetDefaultEnvironment(utils.MainConfigFilePath), "Clear user details of an environment")
}
