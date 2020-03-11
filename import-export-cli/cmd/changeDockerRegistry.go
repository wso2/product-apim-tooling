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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/operator/registry"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const changeCmdLiteral = "change"
const changeCmdShortDesc = "Change a configuration"
const changeCmdLongDesc = "Change a configuration in k8s cluster resource"
const changeCmdExamples = utils.ProjectName + ` ` + changeCmdLiteral + ` ` + changeDockerRegistryCmdLiteral

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:     changeCmdLiteral,
	Short:   changeCmdShortDesc,
	Long:    changeCmdLongDesc,
	Example: changeCmdExamples,
}

const changeDockerRegistryCmdLiteral = "registry"
const changeDockerRegistryCmdShortDesc = "Change the registry"
const changeDockerRegistryCmdLongDesc = "Change the registry to be pushed the built micro-gateway image"
const changeDockerRegistryCmdExamples = utils.ProjectName + ` ` + changeCmdLiteral + ` ` + changeDockerRegistryCmdLiteral

// changeDockerRegistryCmd represents the change registry command
var changeDockerRegistryCmd = &cobra.Command{
	Use:     changeDockerRegistryCmdLiteral,
	Short:   changeDockerRegistryCmdShortDesc,
	Long:    changeDockerRegistryCmdLongDesc,
	Example: changeDockerRegistryCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(fmt.Sprintf("%s%s %s called", utils.LogPrefixInfo, changeCmdLiteral, changeDockerRegistryCmdLiteral))
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if !configVars.Config.KubernetesMode {
			utils.HandleErrorAndExit("set mode to kubernetes with command: apictl set --mode kubernetes",
				errors.New("mode should be set to kubernetes"))
		}

		// read inputs for docker registry
		registry.ChooseRegistryInteractive()
		registry.ReadInputs()
		registry.UpdateConfigsSecrets()
	},
}

func init() {
	RootCmd.AddCommand(changeCmd)
	changeCmd.AddCommand(changeDockerRegistryCmd)
}
