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
	"os"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/git"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagVCSDeployEnvName string    // name of the environment the project changes need to be deployed
var flagVCSDeploySkipRollback bool // specifies whether rolling back on error needs to be avoided

// deploy command related usage Info
const deployCmdLiteral = "deploy"
const deployCmdShortDesc = "Deploys projects to the specified environment"
const deployCmdLongDesc = `Deploys projects to the specified environment specified by --environment(-e). 
Only the changed projects compared to the revision at the last successful deployment will be deployed. 
If any project(s) got failed during the deployment, by default, the operation will rollback the environment to the last successful state. If this needs to be avoided, use --skipRollback=true
NOTE: --environment (-e) flag is mandatory`

const deployCmdExamples = utils.ProjectName + ` ` + vcsCmdLiteral + ` ` + deployCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + vcsCmdLiteral + ` ` + deployCmdLiteral + ` -e dev --skipRollback=true`

// deployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:     deployCmdLiteral,
	Short:   deployCmdShortDesc,
	Long:    deployCmdLongDesc,
	Example: deployCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deployCmdLiteral + " called")
		if !utils.EnvExistsInMainConfigFile(flagVCSDeployEnvName, utils.MainConfigFilePath) {
			fmt.Println(flagVCSDeployEnvName, "does not exists. Add it using add-env")
			os.Exit(1)
		}

		credential, err := GetCredentials(flagVCSDeployEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, flagVCSDeployEnvName)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for deploying the project(s)", err)
		}
		failedProjects := git.DeployChangedFiles(accessOAuthToken, flagVCSDeployEnvName)
		if failedProjects != nil && len(failedProjects) > 0 && flagVCSDeploySkipRollback == false {
			fmt.Println("\nRolling back to the last successful revision as there are failures..")
			err = git.Rollback(accessOAuthToken, flagVCSDeployEnvName)
			if err != nil {
				utils.HandleErrorAndExit("There are project deployment failures. Failed to rollback.", err)
			} else {
				utils.HandleErrorAndExit("There are project deployment failures. Rolled back to the last successful revision.", err)
			}
		}
	},
}

func init() {
	VCSCmd.AddCommand(DeployCmd)

	DeployCmd.Flags().StringVarP(&flagVCSDeployEnvName, "environment", "e", "", "Name of the "+
		"environment to deploy the project(s)")
	DeployCmd.Flags().BoolVarP(&flagVCSDeploySkipRollback, "skipRollback", "", false,
		"Specifies whether rolling back to the last successful revision during an error situation should be skipped")

	_ = DeployCmd.MarkFlagRequired("environment")
}
