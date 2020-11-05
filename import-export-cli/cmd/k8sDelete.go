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
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

// Delete command related usage Info
const k8sDeleteCmdLiteral = "delete"
const k8sDeleteCmdShortDesc = "Delete resources related to kubernetes"
const k8sDeleteCmdLongDesc = `Delete resources by filenames, stdin, resources and names, or by resources and label selector in kubernetes mode`

const k8sDeleteCmdExamples = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPICmdLiteral + ` petstore
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPICmdLiteral + ` -l name=myLabel`

// k8sDeleteCmd represents the delete command
var k8sDeleteCmd = &cobra.Command{
	Use:                k8sDeleteCmdLiteral,
	Short:              k8sDeleteCmdShortDesc,
	Long:               k8sDeleteCmdLongDesc,
	Example:            k8sDeleteCmdExamples,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + k8sDeleteCmdLiteral + " called")
		k8sArgs := []string{k8sUtils.K8sDelete}
		k8sArgs = append(k8sArgs, args...)
		executeKubernetes(k8sArgs...)
	},
}

func init() {
	K8sCmd.AddCommand(k8sDeleteCmd)
}
