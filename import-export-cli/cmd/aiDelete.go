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

// Purge command related usage Info
const PurgeCmdLiteral = "delete"
const PurgeCmdShortDesc = "Purge APIs and API Products of a tenant from one environment from the vector database."
const PurgeCmdLongDesc = `Purge APIs and API Products of a tenant from one environment specified by flag (--environment, -e)`
const PurgeCmdExamples = utils.ProjectName + ` ` + AiCmdLiteral + ` ` + PurgeCmdLiteral + ` ` + PurgeAPIsCmdLiteral + ` -e production
NOTE:The flag (--environment (-e)) is mandatory`

// PurgeCmd represents the Purge command
var PurgeCmd = &cobra.Command{
	Use:     PurgeCmdLiteral,
	Short:   PurgeCmdShortDesc,
	Long:    PurgeCmdLongDesc,
	Example: PurgeCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + PurgeCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	AiCmd.AddCommand(PurgeCmd)
}
