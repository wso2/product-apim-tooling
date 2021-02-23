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

package mg

import (
	"github.com/spf13/cobra"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/impl/mg"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var loginUsername string
var loginPassword string
var loginPasswordStdin bool

const loginCmdLiteral = "login [environment]"
const loginCmdShortDesc = "Login to an Microgateway Adapter environment"
const loginCmdLongDesc = `Login to an Microgateway Adapter environment using username and password`
const loginCmdExamples = utils.ProjectName + " " + mgCmdLiteral + " login dev -u admin -p admin\n" +
	utils.ProjectName + " " + mgCmdLiteral + " login dev -u admin\n" +
	"cat ~/.mypassword | " + utils.ProjectName + " " + mgCmdLiteral + " login dev -u admin --password-stdin"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     loginCmdLiteral,
	Short:   loginCmdShortDesc,
	Long:    loginCmdLongDesc,
	Example: loginCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environment := args[0]

		err := impl.RunLogin(environment, loginUsername, loginPassword,
			loginPasswordStdin)
		if err != nil {
			utils.HandleErrorAndExit("Error occurred while login : ", err)
		}
	},
}

// init using Cobra
func init() {
	MgCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for login")
	loginCmd.Flags().BoolVarP(&loginPasswordStdin, "password-stdin", "", false, "Get password from stdin")
}
