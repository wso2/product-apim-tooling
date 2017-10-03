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

var flagAddEnvName string           // name of the environment to be added
var flagTokenEndpoint string        // token endpoint of the environment to be added
var flagRegistrationEndpoint string // registration endpoint of the environment to be added
var flagAPIManagerEndpoint string   // api manager endpoint of the environment to be added

// addEnvCmd represents the addEnv command
var addEnvCmd = &cobra.Command{
	Use:   "add-env",
	Short: utils.AddEnvCmdShortDesc,
	Long:  utils.AddEnvCmdLongDesc + utils.AddEnvCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "add-env called")
		err := addEnv(flagAddEnvName, flagAPIManagerEndpoint, flagRegistrationEndpoint, flagTokenEndpoint)
		if err != nil {
			utils.HandleErrorAndExit("Error adding environment", err)
		}
	},
}

func addEnv(envName string, apimEndpoint string, regEndpoint string, tokenEndpoint string) error {
	mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)

	if envName == "" {
		// name of the environment is blank
		return errors.New("name of the environment cannot be blank")
	}
	if apimEndpoint == "" || regEndpoint == "" || tokenEndpoint == "" {
		// at least one of the 3 endpoints is blank
		return errors.New("none of the 3 endpoints can be blank")
	}
	if utils.EnvExistsInMainConfigFile(envName, utils.MainConfigFilePath) {
		// environment already exists
		return errors.New("environment '" + envName + "' already exists in " + utils.MainConfigFilePath)
	}

	var envEndpoints utils.EnvEndpoints = utils.EnvEndpoints{APIManagerEndpoint: apimEndpoint,
		TokenEndpoint: tokenEndpoint, RegistrationEndpoint: regEndpoint}

	mainConfig.Environments[envName] = envEndpoints
	utils.WriteConfigFile(mainConfig, utils.MainConfigFilePath)

	fmt.Println("Successfully added environment '" + envName + "'")

	return nil
}

func init() {
	RootCmd.AddCommand(addEnvCmd)

	addEnvCmd.Flags().StringVarP(&flagAddEnvName, "name", "n", "",
		"Name of the environment to be added")
	addEnvCmd.Flags().StringVar(&flagAPIManagerEndpoint, "apim", "",
		"API Manager endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagTokenEndpoint, "token", "",
		"Token endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagRegistrationEndpoint, "registration", "",
		"Registration endpoint for the environment")
}
