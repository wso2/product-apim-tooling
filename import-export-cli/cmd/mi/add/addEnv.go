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

package add

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var envToBeAdded string         // Name of the environment to be added
var miManagementEndpoint string // mi management endpoint of the environment to be added

// AddEnv command related Info
const AddEnvCmdLiteral = "env [environment]"
const AddEnvCmdLiteralTrimmed = "env"
const addEnvCmdShortDesc = "Add Environment to Config file"
const addEnvCmdLongDesc = "Add new environment and its related endpoints to the config file"

var addEnvCmdExamples = utils.MICmd + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` production https://localhost:9164

` + utils.MICmd + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` dev https://localhost:9164

` + utils.MICmd + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` prod https://localhost:9164

` + utils.MICmd + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` test https://localhost:9164

` + utils.MICmd + ` ` + AddCmdLiteral + ` ` + AddEnvCmdLiteralTrimmed + ` dev https://localhost:9164`

// addEnvCmd represents the addEnv command
var addEnvCmd = &cobra.Command{
	Use:     AddEnvCmdLiteral,
	Short:   addEnvCmdShortDesc,
	Long:    addEnvCmdLongDesc,
	Example: addEnvCmdExamples,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envToBeAdded = args[0]
		miManagementEndpoint = args[1]
		utils.Logln(utils.LogPrefixInfo + AddCmdLiteral + " " + AddEnvCmdLiteralTrimmed + " called")
		executeAddEnvCmd(utils.MainConfigFilePath)
	},
}

func executeAddEnvCmd(mainConfigFilePath string) {
	envEndpoints := new(utils.EnvEndpoints)
	envEndpoints.MiManagementEndpoint = miManagementEndpoint
	err := impl.AddMIEnv(envToBeAdded, envEndpoints, mainConfigFilePath, AddEnvCmdLiteral)
	if err != nil {
		utils.HandleErrorAndExit("Error adding environment", err)
	}
}

// init using Cobra
func init() {
	if utils.GetMICmdName() == "" {
		AddCmd.AddCommand(addEnvCmd)
	}
}
