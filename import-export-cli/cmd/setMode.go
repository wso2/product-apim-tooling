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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const setModeCmdLiteral = "set-mode"
const setModeCmdShortDesc = "Set mode in " + utils.ProjectName
const setModeCmdLongDesc = `Set mode allows to change the mode of apimcli tool.
available modes are as follows
* kubernetes`
const kubernetesCmdLiteral = "kubernetes"
const kubernetesCmdShortDesc = "Set mode to kubernetes in " + utils.ProjectName
const kubernetesCmdLongDesc = `Set mode kubernetes allows to change the mode of apimcli tool to kubernetes and execute 
kubectl commads via apimcli tool.`
const setModeCmdExamples = utils.ProjectName + " " + setModeCmdLiteral + " " + kubernetesCmdLiteral

// setModeCmd represents the setMode command
var setModeCmd = &cobra.Command{
	Use:     setModeCmdLiteral,
	Short:   setModeCmdShortDesc,
	Long:    setModeCmdLongDesc,
	Example: setModeCmdExamples,
}

var kubernetesCmd = &cobra.Command{
	Use:     kubernetesCmdLiteral,
	Short:   kubernetesCmdShortDesc,
	Long:    kubernetesCmdLongDesc,
	Example: setModeCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + kubernetesCmdLiteral + " called")
		setMode(utils.MainConfigFilePath)

	},
}

func setMode(mainConfigFilePath string) {
	mainConfig := utils.GetMainConfigFromFile(mainConfigFilePath)
	mainConfig.Config.KubernetesMode = true
	utils.WriteConfigFile(mainConfig, mainConfigFilePath)
	fmt.Println("change mode to kubernetes")

}

func init() {
	RootCmd.AddCommand(setModeCmd)
	setModeCmd.AddCommand(kubernetesCmd)
}
