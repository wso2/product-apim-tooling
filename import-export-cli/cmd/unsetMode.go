/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const unsetModeCmdLiteral = "unset-mode"
const unsetModeCmdShortDesc = "unset mode in  " + utils.ProjectName
const unsetModeCmdLongDesc = `Unset mode allows to reset the mode of apimcli tool. use the -n flag to define the mode.
available modes are as follows
* kubernetes`
const unsetModeCmdExamples = utils.ProjectName + " " + unsetModeCmdLiteral + " " + kubernetesCmdLiteral

const kubernetesUnsetCmdShortDesc = "Reset mode to the default in " + utils.ProjectName
const kubernetesUnsetCmdLongDesc = `unset mode kubernetes allows to change the mode of apimcli tool to the default mode.`

// unsetModeCmd represents the unsetMode command
var unsetModeCmd = &cobra.Command{
	Use:     unsetModeCmdLiteral,
	Short:   unsetModeCmdShortDesc,
	Long:    unsetModeCmdLongDesc,
	Example: unsetModeCmdExamples,
}

var kubernetesUnsetCmd = &cobra.Command{
	Use:     kubernetesCmdLiteral,
	Short:   kubernetesUnsetCmdShortDesc,
	Long:    kubernetesUnsetCmdLongDesc,
	Example: unsetModeCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + kubernetesCmdLiteral + " called")
		//make the kubernetes_mode config false
		unsetMode(utils.MainConfigFilePath)

	},
}

func unsetMode(mainConfigFilePath string) {
	mainConfig := utils.GetMainConfigFromFile(mainConfigFilePath)
	mainConfig.Config.KubernetesMode = false
	utils.WriteConfigFile(mainConfig, mainConfigFilePath)
	fmt.Println("change mode to default")
}

func init() {
	RootCmd.AddCommand(unsetModeCmd)
	unsetModeCmd.AddCommand(kubernetesUnsetCmd)
}
