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
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"strings"

	"github.com/spf13/cobra"
)

const uninstallWso2amOperatorCmdLiteral = "wso2am-operator"
const uninstallWso2amOperatorCmdShortDesc = "Uninstall API Operator"
const uninstallWso2amOperatorCmdLongDesc = "Uninstall API Operator in the configured K8s cluster"
const uninstallWso2amOperatorCmdExamples = utils.ProjectName + ` ` + uninstallCmdLiteral + ` ` + uninstallWso2amOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + uninstallCmdLiteral + ` ` + uninstallWso2amOperatorCmdLiteral + ` --force`

var flagForceUninstallWso2amOperator bool

// uninstallWso2amOperatorCmd represents the uninstallWso2amOperator command
var uninstallWso2amOperatorCmd = &cobra.Command{
	Use:     uninstallWso2amOperatorCmdLiteral,
	Short:   uninstallWso2amOperatorCmdShortDesc,
	Long:    uninstallWso2amOperatorCmdLongDesc,
	Example: uninstallWso2amOperatorCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		isConfirm := flagForceUninstallWso2amOperator

		// getting API Operator version
		operatorVersion, err := k8sUtils.GetVersion(
			"WSO2AM Operator",
			k8sUtils.Wso2AmOperatorVersionEnvVariable,
			k8sUtils.DefaultWso2AmOperatorVersion,
			k8sUtils.Wso2AmOperatorVersionValidationUrlTemplate,
			k8sUtils.Wso2AmOperatorFindVersionUrl,
		)
		if err != nil {
			utils.HandleErrorAndExit("Error in WSO2AM Operator version", err)
		}

		if !flagForceUninstallWso2amOperator {
			isConfirmStr, err := utils.ReadInputString(
				fmt.Sprintf("\nUninstall \"%s-%s\" and all related resources: APIs, Securities, Rate Limitings and Target Endpoints\n"+
					"[WARNING] Remove the namespace: %s\n"+
					"Are you sure",
					k8sUtils.ApiOperator, operatorVersion, k8sUtils.ApiOpWso2Namespace),
				utils.Default{Value: "N", IsDefault: true},
				"",
				false,
			)
			if err != nil {
				utils.HandleErrorAndExit("Error reading user input Confirmation", err)
			}

			isConfirmStr = strings.ToUpper(isConfirmStr)
			isConfirm = isConfirmStr == "Y" || isConfirmStr == "YES"
		}

		if isConfirm {
			fmt.Println("Deleting kubernetes resources for API Operator")

			// delete the OLM subscription
			fmt.Println("Removing OLM subscription")
			_ = k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "subscriptions.operators.coreos.com", "my-api-operator", "-n", "operators")

			// delete the namespace "wso2-system"
			// namespace, "wso2-system" contains all the artifacts and configs
			// deleting the namespace: "wso2-system", will remove all the artifacts and configs
			fmt.Printf("Removing namespace: %s\nThis operation will take some minutes...\n", k8sUtils.ApiOpWso2Namespace)

			deleteErrors := []error{
				k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "namespace", k8sUtils.ApiOpWso2Namespace),
				k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "crd", k8sUtils.ApiOpCrdApi),
				k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "crd", k8sUtils.ApiOpCrdSecurity),
				k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "crd", k8sUtils.ApiOpCrdRateLimiting),
				k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "crd", k8sUtils.ApiOpCrdTargetEndpoint),
			}
			// ignore csv deletion
			_ = k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, "csv", fmt.Sprintf("%s.%s", k8sUtils.ApiOperator, operatorVersion), "-n", "operators")

			for _, err := range deleteErrors {
				if err != nil {
					utils.HandleErrorAndExit("Error uninstalling API Operator", err)
				}
			}
		} else {
			fmt.Println("Cancelled")
		}
	},
}

func init() {
	uninstallCmd.AddCommand(uninstallWso2amOperatorCmd)
	uninstallWso2amOperatorCmd.Flags().BoolVar(&flagForceUninstallWso2amOperator, "force", false, "Force uninstall WSO2AM Operator")
}
