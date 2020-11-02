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

const K8sInstallCmdLiteral = "install"
const k8sInstallCmdShortDesc = "Install an operator in the configured K8s cluster"
const k8sInstallCmdLongDesc = "Install an operator in the configured K8s cluster"
const k8sInstallCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallApiOperatorCmdLiteral

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     K8sInstallCmdLiteral,
	Short:   k8sInstallCmdShortDesc,
	Long:    k8sInstallCmdLongDesc,
	Example: k8sInstallCmdExamples,
}

func init() {
	K8sCmd.AddCommand(installCmd)
}
