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

const getCmdLiteral = "get"
const getCmdShortDesc = "List APIs in Microgateway"
const getCmdLongDesc = `Display a get of all the APIs or 
a set of APIs with a limit or filtered by apiType using the flags --limit (-l), --type (-t). 
Note: The flags --host (-c), --username (-u) are mandatory. The password can be included 
via the flag --password (-p) or entered at the prompt.`

const getCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` +
	getCmdLiteral + ` ` + getAPIsCmdLiteral + ` -t http --host https://localhost:9095 -u admin -l 100`

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:     getCmdLiteral,
	Short:   getCmdShortDesc,
	Long:    getCmdLongDesc,
	Example: getCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + getCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	MgCmd.AddCommand(GetCmd)
}
