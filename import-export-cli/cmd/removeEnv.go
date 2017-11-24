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
	"errors"
	"fmt"
	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagNameOfEnvToBeRemoved string

// RemoveEnv command related Info
const removeEnvCmdLiteral = "remove-env"
const removeEnvCmdShortDesc = "Remove Environment from Config file"

var removeEnvCmdLongDesc = dedent.Dedent(`
		Remove Environment and its related endpoints from the config file
    `)

var removeEnvCmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + removeEnvCmdLiteral + ` -n production
	`)


// removeEnvCmd represents the removeEnv command
var removeEnvCmd = &cobra.Command{
	Use:   removeEnvCmdLiteral,
	Short: removeEnvCmdShortDesc,
	Long:  removeEnvCmdLongDesc + removeEnvCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + removeEnvCmdLiteral + " called")
		err := removeEnv(flagNameOfEnvToBeRemoved, utils.MainConfigFilePath, utils.EnvKeysAllFilePath)
		if err != nil {
			utils.HandleErrorAndExit("Error removing environment", err)
		}
	},
}

// removeEnv
// to be used with 'remove-env' command
// @param envName : Name of the environment to be removed from the
// @param mainConfigFilePath : Path to file where env endpoints are stored
// @param envKeysFilePath : Path to file where env keys are stored
// @return error
func removeEnv(envName string, mainConfigFilePath string, envKeysFilePath string) error {
	if envName == "" {
		return errors.New("name of the environment cannot be blank")
	}
	if utils.EnvExistsInMainConfigFile(envName, mainConfigFilePath) {
		var err error
		if utils.EnvExistsInKeysFile(envName, envKeysFilePath) {
			// environment exists in keys file, it has to be cleared first
			err = utils.RemoveEnvFromKeysFile(envName, envKeysFilePath, mainConfigFilePath)
			if err != nil {
				return err
			}
		}
		// remove env from mainConfig file (endpoints file)
		err = utils.RemoveEnvFromMainConfigFile(envName, mainConfigFilePath)

		if err != nil {
			return err
		}
	} else {
		// environment does not exist in mainConfig file (endpoints file). Nothing to remove
		return errors.New("environment '" + envName + "' not found in " + mainConfigFilePath)
	}

	fmt.Println("Successfully removed environment '" + envName + "'")
	fmt.Println("Execute '" + utils.ProjectName + " add-env" + " --help' to see how to add a new environment")

	return nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(removeEnvCmd)
	removeEnvCmd.Flags().StringVarP(&flagNameOfEnvToBeRemoved, "name", "n",
		utils.GetDefaultEnvironment(utils.MainConfigFilePath), "Name of the environment to be removed")
}
