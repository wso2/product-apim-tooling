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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var loginUsername string
var loginPassword string
var loginPasswordStdin bool
var clientId string
var clientSecret string
var personalAccessToken string

const loginCmdLiteral = "login [environment] [flags]"
const loginCmdShortDesc = "Login to an API Manager"
const loginCmdLongDesc = `Login to an API Manager using credentials or set token for authentication`
const loginCmdExamples = utils.ProjectName + " login dev -u admin -p admin\n" +
	utils.ProjectName + " login dev -u admin\n" +
	"cat ~/.mypassword | " + utils.ProjectName + " login dev -u admin\n" +
	utils.ProjectName + " login dev --token e79bda48-3406-3178-acce-f6e4dbdcbb12"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     loginCmdLiteral,
	Short:   loginCmdShortDesc,
	Long:    loginCmdLongDesc,
	Example: loginCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environment := args[0]
		store, err := credentials.GetDefaultCredentialStore()
		if err != nil {
			fmt.Println("Error occurred while loading credential store : ", err)
			os.Exit(1)
		}
		if personalAccessToken != "" {
			err = runLogin(store, environment, loginUsername, loginPassword, personalAccessToken)
			if err != nil {
				fmt.Println("Error occurred while login using the token : ", err)
				os.Exit(1)
			}
		} else {
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
			err = runLogin(store, environment, loginUsername, loginPassword, personalAccessToken)
			if err != nil {
				fmt.Println("Error occurred while login : ", err)
				os.Exit(1)
			}
		}
	},
}

func runLogin(store credentials.Store, environment, username, password, personalAccessToken string) error {
	if !utils.APIMExistsInEnv(environment, utils.MainConfigFilePath) {
		fmt.Println("APIM does not exists in", environment, "Add it using add env")
		os.Exit(1)
	}

	if username == "" && personalAccessToken == "" {
		fmt.Print("Username:")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			username = scanner.Text()
		}
	}

	if password == "" && personalAccessToken == "" {
		fmt.Print("Password:")
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		password = string(pass)
		fmt.Println()
	}

	if username != "" && password != "" {
		registrationEndpoint := utils.GetRegistrationEndpointOfEnv(environment, utils.MainConfigFilePath)
		id, secret, err := utils.GetClientIDSecret(username, password, registrationEndpoint)
		if err != nil {
			return err
		}
		clientId = id
		clientSecret = secret
	}

	fmt.Println("Logged into APIM in ", environment, "environment")
	err := store.SetAPIMCredentials(environment, username, password, clientId, clientSecret, personalAccessToken)
	if err != nil {
		return err
	}

	return nil
}

// GetCredentials function gets the credentials for the specified environment
func GetCredentials(env string) (credentials.Credential, error) {
	// get tokens or login
	store, err := credentials.GetDefaultCredentialStore()
	if err != nil {
		return credentials.Credential{}, err
	}

	if !utils.APIMExistsInEnv(env, utils.MainConfigFilePath) {
		fmt.Println("APIM does not exists in", env, "Add it using add env")
		os.Exit(1)
	}

	// check for creds
	if !store.HasAPIM(env) {
		fmt.Println("Login to APIM in ", env)
		err = runLogin(store, env, "", "", "")
		if err != nil {
			return credentials.Credential{}, err
		}
		fmt.Println()
	}
	cred, err := store.GetAPIMCredentials(env)
	if err != nil {
		return credentials.Credential{}, err
	}
	return cred, nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for login")
	loginCmd.Flags().BoolVarP(&loginPasswordStdin, "password-stdin", "", false, "Get password from stdin")
	loginCmd.Flags().StringVarP(&personalAccessToken, "token", "", "", "Personal access token")
}
