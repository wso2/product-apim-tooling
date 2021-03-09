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

package k8s

import (
	"fmt"
	"github.com/spf13/cobra"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"strings"
)

// k8sDeleteAPI command related usage info
const k8sDeleteAPICmdLiteral = "api"
const k8sDeleteAPICmdShortDesc = "Delete API resources"
const k8sDeleteAPICmdLongDesc = "Delete API resources by API name or label selector in kubernetes mode"

const k8sDeleteAPICmdExamples = utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` petstore`

// k8sDeleteAPICmd represents the delete api command in kubernetes mode
var deleteAPICmd = &cobra.Command{
	Use:     k8sDeleteAPICmdLiteral,
	Short:   k8sDeleteAPICmdShortDesc,
	Long:    k8sDeleteAPICmdLongDesc,
	Example: k8sDeleteAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + k8sDeleteAPICmdLiteral + " called")
		handleDeleteApi()
	},
}

func handleDeleteApi() {
	flagApiName = strings.ToLower(flagApiName)
	var errMsg string
	deleteApiErr := k8sUtils.ExecuteCommand(
		k8sUtils.Kubectl, k8sUtils.K8sDelete, k8sUtils.ApiOpCrdApi, flagApiName, "-n", flagNamespace)
	if deleteApiErr != nil {
		if flagNamespace != "" {
			errMsg = fmt.Sprintf("Could not find the API \"%s\" in the namespace \"%s\"",
				flagApiName, flagNamespace)
		} else {
			errMsg = fmt.Sprintf("Could not find the API \"%s\"", flagApiName)
		}
		utils.HandleErrorAndExit(errMsg, nil)
	}
}

// Init using Cobra
func init() {
	DeleteCmd.AddCommand(deleteAPICmd)
	deleteAPICmd.Flags().StringVarP(&flagApiName, "name", "n", "", "API name")
	_ = deleteAPICmd.MarkFlagRequired("name")
}
