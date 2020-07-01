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

var flagVCSDeployEnvName string           // name of the environment to be added

// deploy command related usage Info
const deployCmdLiteral = "deploy"
const deployCmdShortDesc = "Deploys project changes to the specified environment"
const deployCmdLongDesc = ``

const deployCmdExamples = utils.ProjectName + ` ` + deployCmdLiteral + ` `  + ` -e dev`

// deployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:     deployCmdLiteral,
	Short:   deployCmdShortDesc,
	Long:    deployCmdLongDesc,
	Example: deployCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deployCmdLiteral + " called")
		credential, err := getCredentials(flagVCSDeployEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, flagVCSDeployEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for deploying the project(s)", err)
		}
		git.DeployChangedFiles(accessOAuthToken, flagVCSDeployEnvName);
	},
}

func init() {
	VCSCmd.AddCommand(DeployCmd)

	DeployCmd.Flags().StringVarP(&flagVCSDeployEnvName, "environment", "e", "", "Name of the " +
		"environment to deploy the project(s)")

	_ = DeployCmd.MarkFlagRequired("environment")
}
