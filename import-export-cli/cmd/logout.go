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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const logoutCmdLiteral = "logout [environment]"
const logoutCmdShortDesc = "Logout to from an API Manager"
const logoutCmdLongDesc = `Logout from an API Manager environment`
const logoutCmdExamples = utils.ProjectName + " logout dev"

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:     logoutCmdLiteral,
	Short:   logoutCmdShortDesc,
	Long:    logoutCmdLongDesc,
	Example: logoutCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runLogout(args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func runLogout(environment string) error {
	cred, err := GetCredentials(environment)
	//Get current access token for
	accessToken, err := credentials.GetOAuthAccessToken(cred, environment)
	error := credentials.RevokeAccessToken(cred, environment, accessToken)
	if error != nil {
		return err
	}
	store, err := credentials.GetDefaultCredentialStore()
	if err != nil {
		return err
	}
	fmt.Println("Logged out from", environment, "environment")
	return store.Erase(environment)
}

// init using Cobra
func init() {
	RootCmd.AddCommand(logoutCmd)
}
