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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/git"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagVCSInitForce bool    // Forcefully create and replace the existing vcs.yaml file if exists

// "vcs init" command related usage Info
const vcsInitCmdLiteral = "init"
const vcsInitCmdShortDesc = "Initializes a GIT repository with API Controller"
const vcsInitCmdLongDesc = `Initializes a GIT repository with API Controller (apictl). Before start using a GIT repository 
for 'vcs' commands, the GIT repository should be initialized once via 'vcs init'. This will create a file 'vcs.yaml'
in the root location of the GIT repository, which is used by API Controller  to uniquely identify the GIT repository. 
'vcs.yaml' should be committed to the GIT repository.`

const vcsInitCmdExamples = utils.ProjectName + ` ` + vcsCmdLiteral + ` ` + vcsInitCmdLiteral

// deployCmd represents the deploy command
var VcsInitCmd = &cobra.Command{
	Use:     vcsInitCmdLiteral,
	Short:   vcsInitCmdShortDesc,
	Long:    vcsInitCmdLongDesc,
	Example: vcsInitCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + vcsInitCmdLiteral + " called")
		err := git.InitializeRepo(flagVCSInitForce)
		if err != nil {
			utils.HandleErrorAndExit("Error initializing repository", err)
		}
		fmt.Println("Successfully initialized GIT repository")
	},
}

func init() {
	VCSCmd.AddCommand(VcsInitCmd)

	VcsInitCmd.Flags().BoolVarP(&flagVCSInitForce, "force", "f", false,
		"Forcefully reinitialize and replace vcs.yaml if already exists in the repository root.")
}
