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
	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// List command related usage Info
const listCmdLiteral = "list"
const listCmdShortDesc = "List APIs/Applications in an environment or List the environments"

var listCmdLongDesc = dedent.Dedent(`
			Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
			Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)
			OR
			List all the environments
	`)

var listCmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + EnvsCmdLiteral + `
		` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e dev
	`)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   listCmdLiteral,
	Short: listCmdShortDesc,
	Long:  listCmdLongDesc + listCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ListCmd)
}
