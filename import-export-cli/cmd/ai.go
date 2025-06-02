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

// Ai command related usage Info
const AiCmdLiteral = "ai"
const AiCmdShortDesc = "AI related commands."
const AiCmdLongDesc = `Perform AI related operations such as uploading APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.`
const AiCmdExamples = utils.ProjectName + ` ` + AiCmdLiteral + ` ` + UploadCmdLiteral + ` ` + UploadAPIsCmdLiteral + ` -e production
` + utils.ProjectName + ` ` + AiCmdLiteral + ` ` + UploadCmdLiteral + ` ` + UploadAPIProductsCmdLiteral + ` -e production
` + utils.ProjectName + ` ` + AiCmdLiteral + ` ` + UploadCmdLiteral + ` ` + UploadAPIsCmdLiteral + ` -e production --all
` + utils.ProjectName + ` ` + AiCmdLiteral + ` ` + PurgeCmdLiteral + ` ` + PurgeAPIsCmdLiteral + ` -e production --all
NOTE:The flag (--environment (-e)) is mandatory`

// AiCmd represents the Ai command
var AiCmd = &cobra.Command{
	Use:     AiCmdLiteral,
	Short:   AiCmdShortDesc,
	Long:    AiCmdLongDesc,
	Example: AiCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + AiCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(AiCmd)
}
