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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagNameOfEnvToBeRemoved string // name of the environment to be removed

// removeEnvCmd represents the removeEnv command
var removeEnvCmd = &cobra.Command{
	Use:   "remove-env",
	Short: utils.RemoveEnvCmdShortDesc,
	Long:  utils.RemoveEnvCmdLongDesc + utils.RemoveEnvCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "remove-env called")
		err := removeEnv(flagNameOfEnvToBeRemoved)
		if err != nil {
			utils.HandleErrorAndExit("Error removing environment", err)
		}
	},
}

// to be used only with 'remove-env' command
func removeEnv(envName string) error {
	if envName == "" {
		return errors.New("name of the environment cannot be blank")
	}
	if utils.EnvExistsInMainConfigFile(envName, utils.MainConfigFilePath) {
		var err error
		if utils.EnvExistsInKeysFile(envName, utils.EnvKeysAllFilePath) {
			// environment exists in keys file, it has to be cleared first
			err = utils.RemoveEnvFromKeysFile(envName, utils.EnvKeysAllFilePath, utils.MainConfigFilePath)
			if err != nil {
				return err
			}
		}
		// remove env from mainConfig file (endpoints file)
		err = utils.RemoveEnvFromMainConfigFile(envName, utils.MainConfigFilePath)

		if err != nil {
			return err
		}
	} else {
		// environment does not exist in mainConfig file (endpoints file). Nothing to remove
		return errors.New("environment '" + envName + "' not found in " + utils.MainConfigFilePath)
	}

	fmt.Println("Successfully removed environment '" + envName + "'")
	fmt.Println("Execute '" + utils.ProjectName + " add-env" + " --help' to see how to add a new environment")

	return nil
}

func init() {
	RootCmd.AddCommand(removeEnvCmd)
	removeEnvCmd.Flags().StringVarP(&flagNameOfEnvToBeRemoved, "name", "n",
		utils.GetDefaultEnvironment(utils.MainConfigFilePath), "Name of the environment to be removed")
}
