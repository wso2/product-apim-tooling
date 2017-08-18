// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/menuka94/wso2apim-cli/utils"
	"strings"
	"github.com/go-resty/resty"
	"encoding/json"
)

var listEnvironment string

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: utils.ListCmdShortDesc,
	Long:  utils.ListCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		if utils.EnvExistsInEndpointsFile(listEnvironment) {
			registrationEndpoint := utils.GetRegistrationEndpointOfEnv(listEnvironment)
			apiManagerEndpoint := utils.GetAPIMEndpointOfEnv(listEnvironment)
			tokenEndpoint := utils.GetTokenEndpointOfEnv(listEnvironment)
			var username string
			var password string
			var clientID string
			var clientSecret string

			if utils.EnvExistsInKeysFile(listEnvironment){
				// client_id, client_secret exists in file
				username = utils.GetUsernameOfEnv(listEnvironment)
				fmt.Println("Username:", username)
				password = utils.PromptForPassword()
				clientID = utils.GetClientIDOfEnv(listEnvironment)
				clientSecret = utils.GetClientSecretOfEnv(listEnvironment, password)

				fmt.Println("ClientID:", clientID)
				fmt.Println("ClientSecret:", clientSecret)
			}else{
				// env exists in endpoints file, but not in keys file
				// no client_id, client_secret in file
				// Get new values
				username = strings.TrimSpace(utils.PromptForUsername())
				password = utils.PromptForPassword()

				fmt.Println("\nUsername:", username)
				clientID, clientSecret = utils.GetClientIDSecret(username, password, registrationEndpoint)

				// Persist clientID, clientSecret, Username in file
				encryptedClientSecret := utils.Encrypt([]byte(utils.GetMD5Hash(password)), clientSecret)
				envKeys := utils.EnvKeys{clientID, encryptedClientSecret, username}
				utils.AddNewEnvToKeysFile(exportEnvironment, envKeys)
			}

			// Get OAuth Tokens
			m := utils.GetOAuthTokens(username, password, utils.GetBase64EncodedCredentials(clientID, clientSecret), tokenEndpoint)
			accessToken := m["access_token"]
			fmt.Println("AccessToken:", accessToken)

			resp, err := GetAPIList(accessToken, apiManagerEndpoint)
			fmt.Println("Status:", resp.Status())

			if err == nil{
				if resp.StatusCode() == 200 {
					apiListResponse := &utils.APIListResponse{}
					unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiListResponse)

					if unmarshalError != nil {
						fmt.Println("UnmarshalError")
						panic(unmarshalError)
					}

					fmt.Println("Count:", apiListResponse.Count)
					for _, api := range apiListResponse.List {
						fmt.Println(api.Name + " " + api.Version)
					}
				}else{
					fmt.Println("Error:", resp.Status())
				}
			}else{
				fmt.Println("Error:")
				panic(err)
			}
		}
	},
}

func GetAPIList(accessToken string, apiManagerEndpoint string) (*resty.Response, error){
	url := apiManagerEndpoint

	// append '/' to the end if there isn't one already
	if string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "apis"

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := resty.R().
			SetHeaders(headers).
			Get(url)
	return resp, err
}

func init() {
	RootCmd.AddCommand(ListCmd)
	ListCmd.Flags().StringVarP(&listEnvironment, "environment", "e", "", "Environment to be searched")


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
