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
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// k8sDeleteAPI command related usage info
const k8sDeleteAPICmdLiteral = "api"
const k8sDeleteAPICmdShortDesc = "Delete API resources"
const k8sDeleteAPICmdLongDesc = "Delete API resources by API name or label selector in kubernetes mode"

const k8sDeleteAPICmdExamples = utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` petstore
` + "  " + utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` -l name=myLabel`

// k8sDeleteAPICmd represents the delete api command in kubernetes mode
var k8sDeleteAPICmd = &cobra.Command{
	Use:                utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + " (<name-of-the-api> or -l name=<name-of-the-label>)",
	Short:              k8sDeleteAPICmdShortDesc,
	Long:               k8sDeleteAPICmdLongDesc,
	Example:            k8sDeleteAPICmdExamples,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAPICmdLiteral + " called")
		k8sArgs := []string{k8sUtils.K8sDelete, k8sUtils.K8sApi}
		k8sArgs = append(k8sArgs, args...)
		executeKubernetes(k8sArgs...)
	},
}

// Init using Cobra
func init() {
	k8sDeleteCmd.AddCommand(k8sDeleteAPICmd)
}
