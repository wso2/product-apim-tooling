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
const deleteCmdLiteral = "delete"
const deleteCmdShortDesc = "Delete an API/APIProduct/Application in an environment"
const deleteCmdLongDesc = `Delete an API available in the environment specified by flag (--environment, -e)
Delete an API Product available in the environment specified by flag (--environment, -e)
Delete an Application of a specific user in the environment specified by flag (--environment, -e)`

const deleteCmdExamples = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPIProductCmdLiteral + ` -n TwitterAPI -r admin -e dev 
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAppCmdLiteral + ` -n TestApplication -o admin -e dev`

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:                deleteCmdLiteral,
	Short:              deleteCmdShortDesc,
	Long:               deleteCmdLongDesc,
	Example:            deleteCmdExamples,
	DisableFlagParsing: isK8sEnabled(),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteCmdLiteral + " called")
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if configVars.Config.KubernetesMode {
			k8sArgs := []string{k8sUtils.K8sDelete}
			k8sArgs = append(k8sArgs, args...)
			executeKubernetes(k8sArgs...)
		}
	},
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
}
