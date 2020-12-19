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

// Get command related usage Info
const GenCmdLiteral = "gen"
const GenCmdShortDesc = "Generate deployment directory for VM and K8S operator"

const GenCmdLongDesc = `Generate sample directory with all the contents to use as the deployment directory` +
	`  when performing CI/CD pipeline tasks `

const GenCmdExamples = utils.ProjectName + ` ` + GenCmdLiteral + ` ` + GenDeploymentDirCmdLiteral

// ListCmd represents the list command
var GenCmd = &cobra.Command{
	Use:     GenCmdLiteral,
	Short:   GenCmdShortDesc,
	Long:    GenCmdLongDesc,
	Example: GenCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GenCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(GenCmd)
}
