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

const K8sUninstallCmdLiteral = "uninstall"
const k8sUninstallCmdShortDesc = "Uninstall an operator in the configured K8s cluster"
const k8sUninstallCmdLongDesc = "Uninstall an operator in the configured K8s cluster"
const k8sUninstallCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sUninstallCmdLiteral + ` ` + K8sUninstallApiOperatorCmdLiteral

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:     K8sUninstallCmdLiteral,
	Short:   k8sUninstallCmdShortDesc,
	Long:    k8sUninstallCmdLongDesc,
	Example: k8sUninstallCmdExamples,
}

func init() {
	K8sCmd.AddCommand(uninstallCmd)
}
