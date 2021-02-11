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
const deleteCmdShortDesc = "Delete an API in Microgateway"
const deleteCmdLongDesc = `Delete an API by specifying name, version, host, username 
and optionally vhost by specifying the flags (--name (-n), --version (-v), --host (-c), 
--username (-u), and optionally --vhost (-t). Note: The password can be included 
via the flag --password (-p) or entered at the prompt.`

const deleteCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteAPICmdLiteral +
	` -n petstore -v 0.0.1 --host https://localhost:9095 -u admin -t www.pets.com`

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
