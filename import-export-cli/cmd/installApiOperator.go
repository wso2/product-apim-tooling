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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/operator/olm"
	"github.com/wso2/product-apim-tooling/import-export-cli/operator/registry"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const installApiOperatorCmdLiteral = "api-operator"
const installApiOperatorCmdShortDesc = "Install API Operator"
const installApiOperatorCmdLongDesc = "Install API Operator in the configured K8s cluster"
const installApiOperatorCmdExamples = utils.ProjectName + ` ` + installCmdLiteral + ` ` + installApiOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + installApiOperatorCmdLiteral + ` -f path/to/operator/configs
` + utils.ProjectName + ` ` + installApiOperatorCmdLiteral + ` -f path/to/operator/config/file.yaml`

// flags
var flagApiOperatorFile string

// flags for installing api-operator in batch mode
var flagBmRegistryType string
var flagBmRepository string
var flagBmUsername string
var flagBmPassword string
var flagBmPasswordStdin bool
var flagBmKeyFile string

// installApiOperatorCmd represents the install api-operator command
var installApiOperatorCmd = &cobra.Command{
	Use:     installApiOperatorCmdLiteral,
	Short:   installApiOperatorCmdShortDesc,
	Long:    installApiOperatorCmdLongDesc,
	Example: installApiOperatorCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + installApiOperatorCmdLiteral + " called")

		// is -f or --from-file flag specified
		isLocalInstallation := flagApiOperatorFile != ""
		configFile := flagApiOperatorFile
		var olmVersion string

		if !isLocalInstallation {
			// getting API Operator version
			operatorVersion, err := k8sUtils.GetVersion(
				"API Operator",
				k8sUtils.ApiOperatorVersionEnvVariable,
				k8sUtils.DefaultApiOperatorVersion,
				k8sUtils.ApiOperatorVersionValidationUrlTemplate,
				k8sUtils.ApiOperatorFindVersionUrl,
			)
			if err != nil {
				utils.HandleErrorAndExit("Error in API Operator version", err)
			}
			configFile = fmt.Sprintf(k8sUtils.ApiOperatorConfigsUrlTemplate, operatorVersion)

			// getting OLM version
			olmVersion, err = k8sUtils.GetVersion(
				"OLM",
				olm.VersionEnvVariable,
				olm.DefaultVersion,
				olm.OlmVersionValidationUrlTemplate,
				olm.OlmVersionFindVersionUrl,
			)
			if err != nil {
				utils.HandleErrorAndExit("Error in OLM version", err)
			}
		}

		// check for installation mode: interactive or batch mode
		if flagBmRegistryType == "" {
			// run api-operator installation in interactive mode
			// read inputs for docker registry
			registry.ChooseRegistryInteractive()
			registry.ReadInputsInteractive()
		} else {
			// run api-operator installation in batch mode
			// set registry type
			registry.SetRegistry(flagBmRegistryType)

			flagsValues := getGivenFlagsValues()
			registry.ValidateFlags(flagsValues)       // validate flags with respect to registry type
			registry.ReadInputsFromFlags(flagsValues) // read values from flags with respect to registry type
		}

		if !isLocalInstallation {
			fmt.Println("[Installing OLM]")
			olm.InstallOLM(olmVersion)

			fmt.Println("[Installing API Operator]")
			olm.InstallApiOperator()
		}

		// installing operator and configs if -f flag given
		// otherwise settings configs only
		createControllerConfigs(configFile)
		registry.UpdateConfigsSecrets()

		fmt.Println("[Setting to K8s Mode]")
		setToK8sMode()
	},
}

// getGivenFlagsValues returns flags that user given in the batch mode except the "registry type"
func getGivenFlagsValues() *map[string]registry.FlagValue {
	flags := make(map[string]registry.FlagValue)
	flags[k8sUtils.FlagBmRepository] = registry.FlagValue{Value: flagBmRepository, IsProvided: flagBmRepository != ""}
	flags[k8sUtils.FlagBmUsername] = registry.FlagValue{Value: flagBmUsername, IsProvided: flagBmUsername != ""}
	flags[k8sUtils.FlagBmPassword] = registry.FlagValue{Value: flagBmPassword, IsProvided: flagBmPassword != ""}
	flags[k8sUtils.FlagBmPasswordStdin] = registry.FlagValue{Value: flagBmPasswordStdin, IsProvided: flagBmPasswordStdin}
	flags[k8sUtils.FlagBmKeyFile] = registry.FlagValue{Value: flagBmKeyFile, IsProvided: flagBmKeyFile != ""}

	return &flags
}

// createControllerConfigs creates configs
func createControllerConfigs(configFile string) {
	utils.Logln(utils.LogPrefixInfo + "Installing controller configs")

	// apply all files without printing errors
	if err := k8sUtils.ExecuteCommandWithoutPrintingErrors(k8sUtils.Kubectl, k8sUtils.K8sApply, "-f", configFile); err != nil {
		fmt.Println("Waiting for resource creation...")

		// if error then wait for namespace and the resource type security
		_ = k8sUtils.K8sWaitForResourceType(20, k8sUtils.ApiOpCrdSecurity)

		// apply again with printing errors
		if err := k8sUtils.K8sApplyFromFile(configFile); err != nil {
			utils.HandleErrorAndExit("Error creating configurations", err)
		}
	}
}

// setToK8sMode sets the "api-ctl" mode to kubernetes
func setToK8sMode() {
	// read the existing config vars
	configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	configVars.Config.KubernetesMode = true
	utils.WriteConfigFile(configVars, utils.MainConfigFilePath)
}

// init using Cobra
func init() {
	installCmd.AddCommand(installApiOperatorCmd)
	installApiOperatorCmd.Flags().StringVarP(&flagApiOperatorFile, "from-file", "f", "", "Path to API Operator directory")

	// flags for installing api-operator in batch mode
	// only the flag "registry-type" is required and others are registry specific flags
	installApiOperatorCmd.Flags().StringVarP(&flagBmRegistryType, "registry-type", "R", "", "Registry type: DOCKER_HUB | AMAZON_ECR |GCR | HTTP")
	installApiOperatorCmd.Flags().StringVarP(&flagBmRepository, k8sUtils.FlagBmRepository, "r", "", "Repository name or URI")
	installApiOperatorCmd.Flags().StringVarP(&flagBmUsername, k8sUtils.FlagBmUsername, "u", "", "Username of the repository")
	installApiOperatorCmd.Flags().StringVarP(&flagBmPassword, k8sUtils.FlagBmPassword, "p", "", "Password of the given user")
	installApiOperatorCmd.Flags().BoolVar(&flagBmPasswordStdin, k8sUtils.FlagBmPasswordStdin, false, "Prompt for password of the given user in the stdin")
	installApiOperatorCmd.Flags().StringVarP(&flagBmKeyFile, k8sUtils.FlagBmKeyFile, "c", "", "Credentials file")
}
