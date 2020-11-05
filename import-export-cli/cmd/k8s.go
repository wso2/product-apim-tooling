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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// K8s command related usage Info
const K8sCmdLiteral = "k8s"
const k8sCmdShortDesc = "Kubernetes mode based commands"

const k8sCmdLongDesc = `Kubernetes mode based commands such as install, uninstall, add/update api, change registry.`

const k8sCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallApiOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sUninstallCmdLiteral + ` ` + K8sUninstallApiOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sAddCmdLiteral + ` ` + AddApiCmdLiteral + ` ` + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sUpdateCmdLiteral + ` ` + AddApiCmdLiteral + ` ` + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sChangeCmdLiteral + ` ` + K8sChangeDockerRegistryCmdLiteral

// K8sCmd represents the import command
var K8sCmd = &cobra.Command{
	Use:     K8sCmdLiteral,
	Short:   k8sCmdShortDesc,
	Long:    k8sCmdLongDesc,
	Example: k8sCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ImportCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(K8sCmd)
}
