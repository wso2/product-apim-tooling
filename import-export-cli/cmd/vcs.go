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

// vcs command related usage Info
const vcsCmdLiteral = "vcs"
const vcsCmdShortDesc = "Checks status and deploys projects"
const vcsCmdLongDesc = `Checks status and deploys projects to the specified environment. In order to 
use this command, 'git' must be installed in the system.'`
const vcsCmdExamples = utils.ProjectName + ` ` + vcsInitCmdLiteral + `
` + utils.ProjectName + ` ` + vcsStatusCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + deployCmdLiteral + ` -e dev`

// vcsCmd represents the vcs command
var VCSCmd = &cobra.Command{
	Use:     vcsCmdLiteral,
	Short:   vcsCmdShortDesc,
	Long:    vcsCmdLongDesc,
	Example: vcsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + vcsCmdLiteral + " called")
	},
}

func init() {
	RootCmd.AddCommand(VCSCmd)
}
