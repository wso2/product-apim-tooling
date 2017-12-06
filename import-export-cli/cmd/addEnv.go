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

var flagAddEnvName string           // name of the environment to be added
var flagTokenEndpoint string        // token endpoint of the environment to be added
var flagRegistrationEndpoint string // registration endpoint of the environment to be added
var flagPublisherEndpoint string    // api manager endpoint of the environment to be added

// AddEnv command related Info
const addEnvCmdLiteral = "add-env"
const addEnvCmdShortDesc = "Add Environment to Config file"

var addEnvCmdLongDesc = dedent.Dedent(`
		Add new environment and its related endpoints to the config file
	`)

var addEnvCmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + addEnvCmdLiteral + ` -n production \
						--registration http://localhost:9763/client-registration/v0.11/register \
						--apim  https://localhost:9443 \
						--token https://localhost:8243/token
	`)

// addEnvCmd represents the addEnv command
var addEnvCmd = &cobra.Command{
	Use:   addEnvCmdLiteral,
	Short: addEnvCmdShortDesc,
	Long:  addEnvCmdLongDesc + addEnvCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + addEnvCmdLiteral + " called")
		executeAddEnvCmd(utils.MainConfigFilePath)
	},
}

func executeAddEnvCmd(mainConfigFilePath string) {
	err := addEnv(flagAddEnvName, flagPublisherEndpoint, flagRegistrationEndpoint, flagTokenEndpoint,
		mainConfigFilePath)
	if err != nil {
		utils.HandleErrorAndExit("Error adding environment", err)
	}

}

// addEnv adds a new environment and its endpoints and writes to config file
// @param envName : Name of the Environment
// @param publisherEndpoint : API Manager Endpoint for the environment
// @param regEndpoint : Registratiopin Endpoint for the environment
// @param tokenEndpoint : Token Endpoint for the environment
// @param mainConfigFilePath : Path to file where env endpoints are stored
// @return error
func addEnv(envName, apiManagerEndpoint, regEndpoint, tokenEndpoint, mainConfigFilePath string) error {
	if envName == "" {
		// name of the environment is blank
		return errors.New("name of the environment cannot be blank")
	}
	if apiManagerEndpoint == "" || regEndpoint == "" || tokenEndpoint == "" {
		// at least one of the 3 endpoints is blank
		utils.ShowHelpCommandTip(addEnvCmdLiteral)
		return errors.New("endpoints cannot be blank")
	}
	if utils.EnvExistsInMainConfigFile(envName, mainConfigFilePath) {
		// environment already exists
		return errors.New("environment '" + envName + "' already exists in " + mainConfigFilePath)
	}

	mainConfig := utils.GetMainConfigFromFile(mainConfigFilePath)

	var envEndpoints = utils.EnvEndpoints{
		APIManagerEndpoint:   apiManagerEndpoint,
		TokenEndpoint:        tokenEndpoint,
		RegistrationEndpoint: regEndpoint,
	}

	mainConfig.Environments[envName] = envEndpoints
	utils.WriteConfigFile(mainConfig, mainConfigFilePath)

	fmt.Printf("Successfully added environment '%s'\n", envName)

	return nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(addEnvCmd)

	addEnvCmd.Flags().StringVarP(&flagAddEnvName, "name", "n", "",
		"Name of the environment to be added")
	addEnvCmd.Flags().StringVarP(&flagPublisherEndpoint, "apim", "a", "",
		"API Manager endpoint for the environment")
	addEnvCmd.Flags().StringVarP(&flagTokenEndpoint, "token", "t", "",
		"Token endpoint for the environment")
	addEnvCmd.Flags().StringVarP(&flagRegistrationEndpoint, "registration", "r", "",
		"Registration endpoint for the environment")
}
