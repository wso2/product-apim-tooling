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

package deprecated

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// List command related usage Info
const listCmdLiteral = "list"
const listCmdShortDesc = "List APIs/APIProducts/Applications in an environment or List the environments"

const listCmdLongDesc = `Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
Display a list containing all the API Products available in the environment specified by flag (--environment, -e)/
Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)
OR
List all the environments`

const listCmdExamples = utils.ProjectName + ` ` + listCmdLiteral + ` ` + EnvsCmdLiteral + `
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e dev`

// ListCmd represents the list command
var ListCmdDeprecated = &cobra.Command{
	Use:        listCmdLiteral,
	Short:      listCmdShortDesc,
	Long:       listCmdLongDesc,
	Example:    listCmdExamples,
	Deprecated: "instead use \"" + cmd.GetCmdLiteral + "\".",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	cmd.RootCmd.AddCommand(ListCmdDeprecated)
}
