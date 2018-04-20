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

var flagAddEnvName string              // name of the environment to be added
var flagTokenEndpoint string           // token endpoint of the environment to be added
var flagApiImportExportEndpoint string // ApiImportExportEndpoint of the environment to be added
var flagApiListEndpoint string         // ApiListEnvironment of the environment to be added
var flagAppListEndpoint string         // ApplicationListEndpoint of the environment to be added
var flagRegistrationEndpoint string    // registration endpoint of the environment to be added
var flagApiManagerEndpoint string      // api manager endpoint of the environment to be added
var flagAdminEndpoint string           // admin endpoint of the environment to be added

// AddEnv command related Info
const addEnvCmdLiteral = "add-env"
const addEnvCmdShortDesc = "Add Environment to Config file"

var addEnvCmdLongDesc = dedent.Dedent(`
		Add new environment and its related endpoints to the config file
	`)

var addEnvCmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + addEnvCmdLiteral + ` -n production \
						--registration https://localhost:9763/client-registration/v0.12/register \
						--apim  https://localhost:9443 \
						--token https://localhost:8243/token

		` + utils.ProjectName + ` ` + addEnvCmdLiteral + ` -n test \
						--registration https://localhost:9763/client-registration/v0.12/register \
					    --import-export https://localhost:9443/api-import-export-2.2.0-v2 \
						--list https://localhsot:9443/api/am/publisher/v0.12/apis \
						--apim  https://localhost:9443 \
						--token https://localhost:8243/token

		` + utils.ProjectName + ` ` + addEnvCmdLiteral + ` -n dev --apim https://localhost:9443 \
						--token	https://localhost:8243/token \
						--registration http://localhost:9763/client-registration/v0.12/register

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
	envEndpoints := new(utils.EnvEndpoints)
	envEndpoints.ApiManagerEndpoint = flagApiManagerEndpoint
	envEndpoints.RegistrationEndpoint = flagRegistrationEndpoint

	envEndpoints.ApiImportExportEndpoint = flagApiImportExportEndpoint
	envEndpoints.ApiListEndpoint = flagApiListEndpoint
	envEndpoints.AppListEndpoint = flagAppListEndpoint
	envEndpoints.AdminEndpoint = flagAdminEndpoint
	envEndpoints.TokenEndpoint = flagTokenEndpoint
	err := addEnv(flagAddEnvName, envEndpoints, mainConfigFilePath)
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
func addEnv(envName string, envEndpoints *utils.EnvEndpoints, mainConfigFilePath string) error {
	if envName == "" {
		// name of the environment is blank
		return errors.New("name of the environment cannot be blank")
	}
	if envEndpoints.ApiManagerEndpoint == "" || envEndpoints.RegistrationEndpoint == "" || envEndpoints.TokenEndpoint == "" {
		// at least one of the 3 mandatory endpoints is blank
		utils.ShowHelpCommandTip(addEnvCmdLiteral)
		return errors.New("endpoints cannot be blank")
	}
	if utils.EnvExistsInMainConfigFile(envName, mainConfigFilePath) {
		// environment already exists
		return errors.New("environment '" + envName + "' already exists in " + mainConfigFilePath)
	}

	mainConfig := utils.GetMainConfigFromFile(mainConfigFilePath)

	var validatedEnvEndpoints = utils.EnvEndpoints{
		ApiManagerEndpoint:   envEndpoints.ApiManagerEndpoint,
		TokenEndpoint:        envEndpoints.TokenEndpoint,
		RegistrationEndpoint: envEndpoints.RegistrationEndpoint,
	}

	if envEndpoints.ApiImportExportEndpoint != "" {
		validatedEnvEndpoints.ApiImportExportEndpoint = envEndpoints.ApiImportExportEndpoint
	}

	if envEndpoints.ApiListEndpoint != "" {
		validatedEnvEndpoints.ApiListEndpoint = envEndpoints.ApiListEndpoint
	}

	mainConfig.Environments[envName] = validatedEnvEndpoints
	utils.WriteConfigFile(mainConfig, mainConfigFilePath)

	fmt.Printf("Successfully added environment '%s'\n", envName)

	return nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(addEnvCmd)

	addEnvCmd.Flags().StringVarP(&flagAddEnvName, "name", "n", "", "Name of the environment to be added")
	addEnvCmd.Flags().StringVar(&flagApiManagerEndpoint, "apim", "", "API Manager endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagApiImportExportEndpoint, "import-export", "",
		"API Import Export endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagApiListEndpoint, "api_list", "", "API List endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagAppListEndpoint, "app_list", "", "Application List endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagTokenEndpoint, "token", "", "Token endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagRegistrationEndpoint, "registration", "",
		"Registration endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagAdminEndpoint, "admin", "", "Admin endpoint for the environment")
}
