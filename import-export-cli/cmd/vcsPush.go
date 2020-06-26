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

var flagVCSPushEnvName string           // name of the environment to be added

// push command related usage Info
const pushCmdLiteral = "push"
const pushCmdShortDesc = "Pushes project changes to the specified environment"
const pushCmdLongDesc = ``

const pushCmdExamples = utils.ProjectName + ` ` + pushCmdLiteral + ` `  + ` -e dev`

// pushCmd represents the push command
var PushCmd = &cobra.Command{
	Use:     pushCmdLiteral,
	Short:   pushCmdShortDesc,
	Long:    pushCmdLongDesc,
	Example: pushCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + pushCmdLiteral + " called")
		credential, err := getCredentials(flagVCSPushEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, flagVCSPushEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for pushing the project(s)", err)
		}
		git.PushChangedFiles(accessOAuthToken, flagVCSPushEnvName);
	},
}

func init() {
	VCSCmd.AddCommand(PushCmd)

	PushCmd.Flags().StringVarP(&flagVCSPushEnvName, "environment", "e", "", "Name of the " +
		"environment to push the project(s)")

	_ = PushCmd.MarkFlagRequired("environment")
}
