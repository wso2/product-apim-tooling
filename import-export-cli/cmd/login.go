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

var loginAPIMUsername string
var loginAPIMPassword string
var loginMIUsername string
var loginMIPassword string
var loginPasswordStdin bool

const loginCmdLiteral = "login [environment] [flags]"
const loginCmdShortDesc = "Login to an API Manager or Micro Integrator"
const loginCmdLongDesc = `Login to an API Manager or Micro Integrator using credentials`
const loginCmdExamples = utils.ProjectName + " login dev --apim-username admin --apim-password admin\n" +
	utils.ProjectName + " login dev --apim-username admin\n" +
	utils.ProjectName + " login test --mi-username admin --mi-password admin\n" +
	"cat ~/.mypassword | " + utils.ProjectName + " login dev --apim-username admin"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     loginCmdLiteral,
	Short:   loginCmdShortDesc,
	Long:    loginCmdLongDesc,
	Example: loginCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		environment := args[0]

		if !utils.EnvExistsInMainConfigFile(environment, utils.MainConfigFilePath) {
			fmt.Println(environment, "does not exists. Add it using add env")
			os.Exit(1)
		}

		miExists := utils.MIExistsInEnv(environment, utils.MainConfigFilePath)
		apimExists := utils.APIMExistsInEnv(environment, utils.MainConfigFilePath)

		onlyMi := miExists && !apimExists
		onlyAPIM := apimExists && !miExists

		if onlyMi {
			if loginAPIMPassword != "" || loginAPIMUsername != "" {
				fmt.Println("No APIM found to login using --apim-username and --apim-password")
				os.Exit(1)
			}
			if loginMIPassword != "" {
				fmt.Println("Warning: Using --mi-password in CLI is not secure. Use --password-stdin")
				if loginPasswordStdin {
					fmt.Println("--mi-password and --password-stdin are mutual exclusive")
					os.Exit(1)
				}
			}
		}

		if onlyAPIM {
			if loginMIPassword != "" || loginMIUsername != "" {
				fmt.Println("No MI found to login using --mi-username and --mi-password")
				os.Exit(1)
			}
			if loginAPIMPassword != "" {
				fmt.Println("Warning: Using --apim-password in CLI is not secure. Use --password-stdin")
				if loginPasswordStdin {
					fmt.Println("--apim-password and --password-stdin are mutual exclusive")
					os.Exit(1)
				}
			}
		}

		if loginAPIMPassword != "" && loginMIPassword != "" && loginPasswordStdin {
			fmt.Println("--apim-password, mi-password and --password-stdin are mutual exclusive")
			os.Exit(1)
		}

		if loginPasswordStdin {

			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if apimExists {
				if loginAPIMUsername == "" {
					fmt.Println("An username for APIM is required to use password-stdin")
					os.Exit(1)
				}
				if loginAPIMPassword == "" {
					loginAPIMPassword = strings.TrimRight(strings.TrimSuffix(string(data), "\n"), "\r")
				}
			}

			if miExists {
				if loginMIUsername == "" {
					fmt.Println("An username for MI is required to use password-stdin")
					os.Exit(1)
				}
				if loginMIPassword == "" {
					loginMIPassword = strings.TrimRight(strings.TrimSuffix(string(data), "\n"), "\r")
				}
			}
		}

		store, err := credentials.GetDefaultCredentialStore()
		if err != nil {
			fmt.Println("Error occurred while loading credential store : ", err)
			os.Exit(1)
		}
		err = runLogin(store, environment, loginAPIMUsername, loginAPIMPassword, loginMIUsername, loginMIPassword)
		if err != nil {
			fmt.Println("Error occurred while login : ", err)
			os.Exit(1)
		}
	},
}

func runLogin(store credentials.Store, environment, apimUsername, apimPassword, miUsername, miPassword string) error {

	if utils.APIMExistsInEnv(environment, utils.MainConfigFilePath) {
		err := runAPIMLogin(store, environment, apimUsername, apimPassword)
		if err != nil {
			return err
		}
	}

	if utils.MIExistsInEnv(environment, utils.MainConfigFilePath) {
		err := credentials.RunMILogin(store, environment, miUsername, miPassword)
		if err != nil {
			return err
		}
	}
	return nil
}

func runAPIMLogin(store credentials.Store, environment, username, password string) error {

	if username == "" {
		fmt.Print("APIM Username:")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			username = scanner.Text()
		}
	}

	if password == "" {
		fmt.Print("APIM Password:")
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		password = string(pass)
		fmt.Println()
	}

	registrationEndpoint := utils.GetRegistrationEndpointOfEnv(environment, utils.MainConfigFilePath)
	clientID, clientSecret, err := utils.GetClientIDSecret(username, password, registrationEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("Logged into APIM in", environment, "environment")
	err = store.Set(environment, username, password, clientID, clientSecret)
	if err != nil {
		return err
	}

	return nil
}

// GetCredentials functions get the credentials for the specified environment
func GetCredentials(env string) (credentials.Credential, error) {
	// get tokens or login
	store, err := credentials.GetDefaultCredentialStore()
	if err != nil {
		return credentials.Credential{}, err
	}

	if !utils.APIMExistsInEnv(env, utils.MainConfigFilePath) {
		fmt.Println("APIM instance does not exists in", env, "Add it using add env")
		os.Exit(1)
	}

	// check for creds
	if !store.Has(env) {

		fmt.Println("Login to APIM in", env)

		err = runAPIMLogin(store, env, "", "")
		if err != nil {
			return credentials.Credential{}, err
		}
		fmt.Println()
	}
	cred, err := store.Get(env)
	if err != nil {
		return credentials.Credential{}, err
	}
	return cred, nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginMIUsername, "mi-username", "", "", "Username for MI login")
	loginCmd.Flags().StringVarP(&loginMIPassword, "mi-password", "", "", "Password for MI login")
	loginCmd.Flags().StringVarP(&loginAPIMUsername, "apim-username", "", "", "Username for APIM login")
	loginCmd.Flags().StringVarP(&loginAPIMPassword, "apim-password", "", "", "Password for APIM login")
	loginCmd.Flags().BoolVarP(&loginPasswordStdin, "password-stdin", "", false, "Get password from stdin")
}
