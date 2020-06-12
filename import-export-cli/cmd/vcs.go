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

// vcs command related usage Info
const vcsCmdLiteral = "vcs"
const vcsCmdShortDesc = "vcs an API/APIProduct/Application in an environment"
const vcsCmdLongDesc = `vcs an API available in the environment specified by flag (--environment, -e) in default mode
vcs an API Product available in the environment specified by flag (--environment, -e) in default mode
vcs an Application of a specific user in the environment specified by flag (--environment, -e) in default mode
vcs resources by filenames, stdin, resources and names, or by resources and label selector in kubernetes mode`

const vcsCmdExamples = utils.ProjectName + ` ` + vcsCmdLiteral + ` `  + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + vcsCmdLiteral + ` ` + ` -n TwitterAPI -r admin -e dev 
` + utils.ProjectName + ` ` + vcsCmdLiteral + ` ` + ` -n TestApplication -o admin -e dev
` + utils.ProjectName + ` ` + vcsCmdLiteral + ` ` + ` petstore
` + utils.ProjectName + ` ` + vcsCmdLiteral + ` ` + ` -l name=myLabel`

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
