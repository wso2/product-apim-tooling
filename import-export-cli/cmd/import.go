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

// Import command related usage Info
const ImportCmdLiteral = "import"
const importCmdShortDesc = "Import an API/API Product/Application to an environment"

const importCmdLongDesc = `Import an API to the environment specified by flag (--environment, -e)
Import an API Product to the environment specified by flag (--environment, -e)
Import an Application to the environment specified by flag (--environment, -e)`

const importCmdExamples = utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportAPICmdLiteral + ` -f qa/TwitterAPI.zip -e dev
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f qa/LeasingAPIProduct.zip -e dev
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportAppCmdLiteral + ` -f qa/apps/sampleApp.zip -e dev`

// ImportCmd represents the import command
var ImportCmd = &cobra.Command{
	Use:     ImportCmdLiteral,
	Short:   importCmdShortDesc,
	Long:    importCmdLongDesc,
	Example: importCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ImportCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ImportCmd)
}
