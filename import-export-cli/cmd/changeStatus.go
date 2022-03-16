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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

// ChangeStatus command related usage info
const changeStatusCmdLiteral = "change-status"
const changeStatusCmdShortDesc = "Change Status of an API or API Product"
const changeStatusCmdLongDesc = "Change the lifecycle status of an API or API Product in an environment"

const changeStatusCmdExamples = utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeAPIStatusCmdLiteral + ` -a Publish -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeAPIStatusCmdLiteral + ` -a Publish -n FacebookAPI -v 2.1.0 -e production
` + utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeAPIProductStatusCmdLiteral + ` -a Publish -n SocialMediaProduct -r admin -e dev`

// ChangeStatusCmd represents the change-status command
var ChangeStatusCmd = &cobra.Command{
	Use:     changeStatusCmdLiteral,
	Short:   changeStatusCmdShortDesc,
	Long:    changeStatusCmdLongDesc,
	Example: changeStatusCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + changeStatusCmdLiteral + " called")
	},
}

func init() {
	RootCmd.AddCommand(ChangeStatusCmd)
}
