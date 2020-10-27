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
const GetCmdLiteral = "get"
const getCmdShortDesc = "Get APIs/APIProducts/Applications in an environment or Get the environments"

const getCmdLongDesc = `Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
Display a list containing all the API Products available in the environment specified by flag (--environment, -e)/
Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)
OR
List all the environments`

const getCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetEnvsCmdLiteral + `
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApisCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiProductsCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAppsCmdLiteral + ` -e dev`

// ListCmd represents the list command
var GetCmd = &cobra.Command{
	Use:     GetCmdLiteral,
	Short:   getCmdShortDesc,
	Long:    getCmdLongDesc,
	Example: getCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(GetCmd)
}
