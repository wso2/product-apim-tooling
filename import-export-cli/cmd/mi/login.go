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

package mi

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var loginUsername string
var loginPassword string
var loginPasswordStdin bool

const loginCmdLiteral = "login [environment] [flags]"
const loginCmdShortDesc = "Login to a Micro Integrator"
const loginCmdLongDesc = `Login to a Micro Integrator using credentials`
var loginCmdExamples = utils.GetMICmdName() + " " + utils.MiCmdLiteral + " login dev -u admin -p admin\n" +
	utils.GetMICmdName() + " " + utils.MiCmdLiteral + " login dev -u admin\n" +
	"cat ~/.mypassword | " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " login dev -u admin"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     loginCmdLiteral,
	Short:   loginCmdShortDesc,
	Long:    loginCmdLongDesc,
	Example: loginCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environment := args[0]

		if loginPassword != "" {
			fmt.Println("Warning: Using --password in CLI is not secure. Use --password-stdin")
			if loginPasswordStdin {
				fmt.Println("--password and --password-stdin are mutual exclusive")
				os.Exit(1)
			}
		}

		if loginPasswordStdin {
			if loginUsername == "" {
				fmt.Println("An username is required to use password-stdin")
				os.Exit(1)
			}

			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			loginPassword = strings.TrimRight(strings.TrimSuffix(string(data), "\n"), "\r")
		}

		store, err := credentials.GetDefaultCredentialStore()
		if err != nil {
			fmt.Println("Error occurred while loading credential store : ", err)
			os.Exit(1)
		}
		err = credentials.RunMILogin(store, environment, loginUsername, loginPassword)
		if err != nil {
			fmt.Println("Error occurred while login : ", err)
			os.Exit(1)
		}
	},
}

func init() {
	MICmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for login")
	loginCmd.Flags().BoolVarP(&loginPasswordStdin, "password-stdin", "", false, "Get password from stdin")
}
