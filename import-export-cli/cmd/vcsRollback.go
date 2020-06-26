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
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/git"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagVCSRollbackEnvName string           // name of the environment to be added

// push command related usage Info
const vcsRollbackCmdLiteral = "rollback"
const vcsRollbackCmdShortDesc = "Rollback the environment to the last working state in case of an error"
const vcsRollbackCmdLongDesc = ``

const vcsRollbackCmdCmdExamples = utils.ProjectName + ` ` + vcsRollbackCmdLiteral + ` `  + ` -e dev`

// pushCmd represents the push command
var VCSRollbackCmd = &cobra.Command{
	Use:     vcsRollbackCmdLiteral,
	Short:   vcsRollbackCmdShortDesc,
	Long:    vcsRollbackCmdLongDesc,
	Example: vcsRollbackCmdCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + vcsStatusCmdLiteral + " called")
		credential, err := getCredentials(flagVCSRollbackEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, flagVCSRollbackEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for rolling back", err)
		}
		git.Rollback(accessOAuthToken, flagVCSRollbackEnvName)
	},
}

func init() {
	VCSCmd.AddCommand(VCSRollbackCmd)

	VCSRollbackCmd.Flags().StringVarP(&flagVCSRollbackEnvName, "environment", "e", "", "Name of the " +
		"environment to check the project(s) status")

	_ = VCSRollbackCmd.MarkFlagRequired("environment")
}
