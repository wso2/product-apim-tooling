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

package deprecated

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/operator/registry"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const changeCmdLiteral = "change"
const changeCmdShortDesc = "Change a configuration in K8s cluster resource"
const changeCmdLongDesc = "Change a configuration in K8s cluster resource"
const changeCmdExamples = utils.ProjectName + ` ` + changeCmdLiteral + ` ` + changeDockerRegistryCmdLiteral

// changeCmdDeprecated represents the change command
var changeCmdDeprecated = &cobra.Command{
	Use:        changeCmdLiteral,
	Short:      changeCmdShortDesc,
	Long:       changeCmdLongDesc,
	Example:    changeCmdExamples,
	Deprecated: "instead use \"" + cmd.K8sCmdLiteral + " " + cmd.K8sChangeCmdLiteral + "\".",
}

const changeDockerRegistryCmdLiteral = "registry"
const changeDockerRegistryCmdShortDesc = "Change the registry"
const changeDockerRegistryCmdLongDesc = "Change the registry to be pushed the built micro-gateway image"
const changeDockerRegistryCmdExamples = utils.ProjectName + ` ` + changeCmdLiteral + ` ` + changeDockerRegistryCmdLiteral

// changeDockerRegistryCmdDeprecated represents the change registry command
var changeDockerRegistryCmdDeprecated = &cobra.Command{
	Use:        changeDockerRegistryCmdLiteral,
	Short:      changeDockerRegistryCmdShortDesc,
	Long:       changeDockerRegistryCmdLongDesc,
	Example:    changeDockerRegistryCmdExamples,
	Deprecated: "instead use \"" + cmd.K8sCmdLiteral + " " + cmd.K8sChangeCmdLiteral + " " + cmd.K8sChangeDockerRegistryCmdLiteral + "\".",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(fmt.Sprintf("%s%s %s called", utils.LogPrefixInfo, changeCmdLiteral, changeDockerRegistryCmdLiteral))
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if !configVars.Config.KubernetesMode {
			utils.HandleErrorAndExit("set mode to kubernetes with command: apictl set --mode kubernetes",
				errors.New("mode should be set to kubernetes"))
		}

		// check for installation mode: interactive or batch mode
		if flagBmRegistryType == "" {
			// run in interactive mode
			// read inputs for docker registry
			registry.ChooseRegistryInteractive()
			registry.ReadInputsInteractive()
		} else {
			// run in batch mode
			// set registry type
			registry.SetRegistry(flagBmRegistryType)

			flagsValues := getGivenFlagsValues()
			registry.ValidateFlags(flagsValues)       // validate flags with respect to registry type
			registry.ReadInputsFromFlags(flagsValues) // read values from flags with respect to registry type
		}

		registry.UpdateConfigsSecrets()
	},
}

func init() {
	cmd.RootCmd.AddCommand(changeCmdDeprecated)
	changeCmdDeprecated.AddCommand(changeDockerRegistryCmdDeprecated)

	// flags for installing api-operator in batch mode
	// only the flag "registry-type" is required and others are registry specific flags
	// same flags defined in 'installApiOperator'
	changeDockerRegistryCmdDeprecated.Flags().StringVarP(&flagBmRegistryType, "registry-type", "R", "", "Registry type: DOCKER_HUB | AMAZON_ECR |GCR | HTTP")
	changeDockerRegistryCmdDeprecated.Flags().StringVarP(&flagBmRepository, k8sUtils.FlagBmRepository, "r", "", "Repository name or URI")
	changeDockerRegistryCmdDeprecated.Flags().StringVarP(&flagBmUsername, k8sUtils.FlagBmUsername, "u", "", "Username of the repository")
	changeDockerRegistryCmdDeprecated.Flags().StringVarP(&flagBmPassword, k8sUtils.FlagBmPassword, "p", "", "Password of the given user")
	changeDockerRegistryCmdDeprecated.Flags().BoolVar(&flagBmPasswordStdin, k8sUtils.FlagBmPasswordStdin, false, "Prompt for password of the given user in the stdin")
	changeDockerRegistryCmdDeprecated.Flags().StringVarP(&flagBmKeyFile, k8sUtils.FlagBmKeyFile, "c", "", "Credentials file")
}
