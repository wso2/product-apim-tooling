// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0 // // Unless required by applicable law or agreed to in writing, software // distributed under the License is distributed on an "AS IS" BASIS, // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/menuka94/wso2apim-cli/utils"
	"log"
	"strings"
)

var importAPIName string
var importAPIVersion string
var importEnvironment string

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use:   "importAPI (--name <name-of-the-api> --version <version-of-the-api> --environment <environment-to-which-the-api-should-be-imported>)",
	Short: utils.ImportAPICmdShortDesc,
	Long:  utils.ImportAPICmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("importAPI called")
		for key, arg := range args {
			fmt.Println(key, ":", arg)
		}

		if utils.EnvExistsInEndpointsFile(importEnvironment) {

			registrationEndpoint := utils.GetRegistrationEndpointOfEnv(importEnvironment)
			apiManagerEndpoint := utils.GetAPIMEndpointOfEnv(importEnvironment)
			tokenEndpoint := utils.GetTokenEndpointOfEnv(importEnvironment)

			if utils.EnvExistsInKeysFile(importEnvironment) {
				// client_id, client_secret,username exists in file
				// get username from file
				username := utils.GetUsernameOfEnv(importEnvironment)
				fmt.Println("Username:", username)

				// get client_id from file
				clientID := utils.GetClientIDOfEnv(importEnvironment)

				// password is needed to decrypt client_secret
				password := utils.PromptForPassword()
				clientSecret := utils.GetClientSecretOfEnv(importEnvironment, password)

				fmt.Println("ClientID:", clientID)
				fmt.Println("ClientSecret:", clientSecret)
			} else {
				// env exists in endpoints file, but not in keys file
				// no client_id, client_secret in file
				// Get new values
				username := strings.TrimSpace(utils.PromptForUsername())
				password := utils.PromptForPassword()

				fmt.Println("\nUsername: " + username + "\n")
				clientID, clientSecret := utils.GetClientIDSecret(username, password, registrationEndpoint)

				// Persist clientID, clientSecret, Username in file
				encryptedClientSecret := utils.Encrypt([]byte(utils.GetMD5Hash(password)), clientSecret)
				envKeys := utils.EnvKeys{clientID, encryptedClientSecret, username}
				utils.AddNewEnvToKeysFile(importEnvironment, envKeys)

				// Get OAuth Tokens
				m := utils.GetOAuthTokens(username, password, utils.GetBase64EncodedCredentials(clientID, clientSecret), tokenEndpoint)
				accessToken := m["access_token"]

				resp := utils.ImportAPI(importAPIName, importAPIVersion, apiManagerEndpoint, accessToken)
				fmt.Printf("Response: %+v\n", resp)
			}
		} else {
			// env_endpoints_all.yaml file is not configured properly by the user
			log.Fatal("Error: env_endpoints_all.yaml does not contain necessary information for environment " + importEnvironment)
		}
	},
}

func init() {
	RootCmd.AddCommand(ImportAPICmd)
	ImportAPICmd.Flags().StringVarP(&importAPIName, "name", "n", "", "Name of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importAPIVersion, "version", "v", "", "Version of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importEnvironment, "environment", "e", "", "Environment from the which the API should be imported")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ImportAPICmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ImportAPICmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
