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

package mg

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const listCmdLiteral = "list"
const listCmdShortDesc = "List APIs"
const listCmdLongDesc = `Display a list containing all the APIs of type specified by flag (--type, -t)`

const listCmdExamples = utils.ProjectName + " " + mgCmdLiteral +
	" " + listCmdLiteral + " " + listApisCmdLiteral + " -h https://localhost:9095 -u admin \n"

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:     listCmdLiteral,
	Short:   listCmdShortDesc,
	Long:    listCmdLongDesc,
	Example: listCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	MgCmd.AddCommand(ListCmd)
}
