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
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/git"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

// push command related usage Info
const pushCmdLiteral = "push"
const pushCmdShortDesc = "push an API/APIProduct/Application in an environment"
const pushCmdLongDesc = `push an API available in the environment specified by flag (--environment, -e) in default mode
push an API Product available in the environment specified by flag (--environment, -e) in default mode
push an Application of a specific user in the environment specified by flag (--environment, -e) in default mode
push resources by filenames, stdin, resources and names, or by resources and label selector in kubernetes mode`

const pushCmdExamples = utils.ProjectName + ` ` + pushCmdLiteral + ` `  + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + pushCmdLiteral + ` ` + ` -n TwitterAPI -r admin -e dev 
` + utils.ProjectName + ` ` + pushCmdLiteral + ` ` + ` -n TestApplication -o admin -e dev
` + utils.ProjectName + ` ` + pushCmdLiteral + ` ` + ` petstore
` + utils.ProjectName + ` ` + pushCmdLiteral + ` ` + ` -l name=myLabel`

// pushCmd represents the push command
var PushCmd = &cobra.Command{
	Use:     pushCmdLiteral,
	Short:   pushCmdShortDesc,
	Long:    pushCmdLongDesc,
	Example: pushCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + pushCmdLiteral + " called")
		environment := "dev"
		credential, err := getCredentials(environment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, environment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing API", err)
		}
		git.GetChangedFiles(accessOAuthToken, environment);
	},
}

func init() {
	RootCmd.AddCommand(PushCmd)
}
