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
	"github.com/wso2/product-apim-tooling/import-export-cli/operator/registry"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"time"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const installApiOperatorCmdLiteral = "api-operator"

var flagApiOperatorFile string

//var flagRegistryHost string
//var flagUsername string
//var flagPassword string
//var flagBatchMod bool

// installApiOperatorCmd represents the install api-operator command
var installApiOperatorCmd = &cobra.Command{
	Use:   installApiOperatorCmdLiteral,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + installApiOperatorCmdLiteral + " called")

		// read inputs for docker registry
		registry.ChooseRegistry()
		registry.ReadInputs()

		isLocalInstallation := flagApiOperatorFile != ""
		if !isLocalInstallation {
			installOLM(utils.OlmVersion)
			installApiOperatorOperatorHub()
		}

		createControllerConfigs(isLocalInstallation)
		registry.CreateSecret()
		setToK8sMode()
	},
}

// installOLM installs Operator Lifecycle Manager (OLM) with the given version
func installOLM(version string) {
	utils.Logln(utils.LogPrefixInfo + "Installing OLM")

	// this implements the logic in
	// https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.13.0/install.sh
	olmNamespace := "olm"
	csvPhaseSuccessed := "Succeeded"

	// apply OperatorHub CRDs
	if err := k8sUtils.K8sApplyFromFile(fmt.Sprintf(utils.OlmCrdUrlTemplate, version)); err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}

	// wait for OperatorHub CRDs
	if err := k8sUtils.K8sWaitForResourceType(10, "clusterserviceversions.operators.coreos.com", "catalogsources.operators.coreos.com", "operatorgroups.operators.coreos.com"); err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}

	// apply OperatorHub OLM
	if err := k8sUtils.K8sApplyFromFile(fmt.Sprintf(utils.OlmOlmUrlTemplate, version)); err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}

	// rolling out
	if err := k8sUtils.ExecuteCommand(utils.Kubectl, utils.K8sRollOut, "status", "-w", "deployment/olm-operator", "-n", olmNamespace); err != nil {
		utils.HandleErrorAndExit("Error installing OLM: Rolling out deployment OLM Operator", err)
	}
	if err := k8sUtils.ExecuteCommand(utils.Kubectl, utils.K8sRollOut, "status", "-w", "deployment/catalog-operator", "-n", olmNamespace); err != nil {
		utils.HandleErrorAndExit("Error installing OLM: Rolling out deployment Catalog Operator", err)
	}

	// wait max 50s to csv phase to be succeeded
	csvPhase := ""
	for i := 50; i > 0 && csvPhase != csvPhaseSuccessed; i-- {
		newCsvPhase, err := k8sUtils.GetCommandOutput(utils.Kubectl, utils.K8sGet, utils.OperatorCsv, "-n", olmNamespace, "packageserver", "-o", `jsonpath={.status.phase}`)
		if err != nil {
			utils.HandleErrorAndExit("Error installing OLM: Getting csv phase", err)
		}

		// only print new phase
		if csvPhase != newCsvPhase {
			fmt.Println("Package server phase: " + newCsvPhase)
			csvPhase = newCsvPhase
		}

		// sleep 1 second
		time.Sleep(1e9)
	}

	if csvPhase != csvPhaseSuccessed {
		utils.HandleErrorAndExit("Error installing OLM: CSV Package Server failed to reach phase succeeded", nil)
	}
}

// installApiOperatorOperatorHub installs WSO2 api-operator from Operator-Hub
func installApiOperatorOperatorHub() {
	utils.Logln(utils.LogPrefixInfo + "Installing API Operator from Operator-Hub")
	operatorFile := utils.OperatorYamlUrl

	err := k8sUtils.K8sApplyFromFile(operatorFile)
	if err != nil {
		utils.HandleErrorAndExit("Error installing API Operator from Operator-Hub", err)
	}
}

// createControllerConfigs creates configs
func createControllerConfigs(isLocalInstallation bool) {
	utils.Logln(utils.LogPrefixInfo + "Installing controller configs")
	fmt.Println("Installing controller configurations...")
	configFile := flagApiOperatorFile

	if !isLocalInstallation {
		// TODO: renuka replace this url (configuration?)
		configFile = `https://gist.githubusercontent.com/renuka-fernando/6d6c64c786e6d13742e802534de3da4e/raw/3a654b9f54d1115532ca757c108b98bff09bcd74/controller_conf.yaml`
	}

	// apply all files without printing errors
	if err := k8sUtils.ExecuteCommandWithoutPrintingErrors(utils.Kubectl, utils.K8sApply, "-f", configFile); err != nil {
		fmt.Println("Installing controller configurations...")

		// if error then wait for namespace and the resource type security
		_ = k8sUtils.K8sWaitForResourceType(20, utils.ApiOpCrdSecurity)

		// apply again with printing errors
		if err := k8sUtils.K8sApplyFromFile(configFile); err != nil {
			utils.HandleErrorAndExit("Error creating configurations", err)
		}
	}
}

// setToK8sMode sets the apictl mode to kubernetes
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
	//installApiOperatorCmd.Flags().StringVarP(&flagRegistryHost, "registry-host", "h", "", "URL of the registry host")
	//installApiOperatorCmd.Flags().StringVarP(&flagUsername, "username", "u", "", "Username for the registry repository")
	//installApiOperatorCmd.Flags().StringVarP(&flagPassword, "password", "p", "", "Password for the registry repository user")
	//installApiOperatorCmd.Flags().BoolVarP(&flagBatchMod, "batch-mod", "B", false, "Run in non-interactive (batch) mode")
}
