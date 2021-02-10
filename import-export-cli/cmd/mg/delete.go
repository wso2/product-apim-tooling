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

const deleteCmdLiteral = "delete"
const deleteCmdShortDesc = "Delete API"
const deleteCmdLongDesc = `Delete an API by specifying name, version and optionally vhost`

const deleteCmdExamples = utils.ProjectName + " " + mgCmdLiteral +
	" " + deleteCmdLiteral + " " + deleteApisCmdLiteral + " -h https://localhost:9095 -u admin " +
	"-n petstore -v version -vhost pets\n"

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:     deleteCmdLiteral,
	Short:   deleteCmdShortDesc,
	Long:    deleteCmdLongDesc,
	Example: deleteCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	MgCmd.AddCommand(DeleteCmd)
}
