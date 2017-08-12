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
	constants "github.com/menuka94/wso2apim-cli/utils"
	"strings"
	"log"
	"io/ioutil"
	"os"
	"github.com/go-resty/resty"

)

var exportAPIName string
var exportAPIVersion string
var exportEnvironment string

// ExportAPICmd represents the exportAPI command
var ExportAPICmd = &cobra.Command{
	Use:   "exportAPI (--name <name-of-the-api> --version <version-of-the-api> --environment <environment-from-which-the-api-should-be-exported>)",
	Short: utils.ExportAPICmdLongDesc,
	Long:  utils.ExportAPICmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exportAPI called")

		if utils.EnvExistsInEndpointsFile(exportEnvironment) {

			registrationEndpoint := utils.GetRegistrationEndpointOfEnv(exportEnvironment)
			apiManagerEndpoint := utils.GetAPIMEndpointOfEnv(exportEnvironment)
			tokenEndpoint := utils.GetTokenEndpointOfEnv(exportEnvironment)
			var username string
			var password string
			var clientID string
			var clientSecret string

			if utils.EnvExistsInKeysFile(exportEnvironment) {
				// client_id, client_secret exists in file
				username = utils.GetUsernameOfEnv(exportEnvironment)
				fmt.Println("Username:", username)
				password = utils.PromptForPassword()
				clientID = utils.GetClientIDOfEnv(exportEnvironment)
				clientSecret = utils.GetClientSecretOfEnv(exportEnvironment, password)

				fmt.Println("ClientID:", clientID)
				fmt.Println("ClientSecret:", clientSecret)
			} else {
				// env exists in endpoints file, but not in keys file
				// no client_id, client_secret in file
				// Get new values
				username = strings.TrimSpace(utils.PromptForUsername())
				password = utils.PromptForPassword()

				fmt.Println("\nUsername: " + username + "\n")
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

			resp := ExportAPI(exportAPIName, exportAPIVersion, apiManagerEndpoint, accessToken)

			// Print info on response
			fmt.Printf("ResponseStatus: %v\n", resp.Status())
			fmt.Printf("Error: %v\n", resp.Error())
			//fmt.Printf("Response Body: %v\n", resp.Body())

			if resp.StatusCode() == 200 {
				// Write to file
				directory := "./exported"

				// create directory if it doesn't exist
				if _, err := os.Stat(directory); os.IsNotExist(err) {
					os.Mkdir(directory, 0777)
				}

				filename := exportAPIName + ".zip"
				err := ioutil.WriteFile(directory+"/"+filename, resp.Body(), 0644)
				if err != nil {
					fmt.Println("Error creating zip archive")
					panic(err)
				}
				fmt.Println("Succesfully wrote to file")
			} else if resp.StatusCode() == 500 {
				fmt.Println("Incorrect password")
			}

		} else {
			// env_endpoints_all.yaml file is not configured properly by the user
			log.Fatal("Error: env_endpoints_all.yaml does not contain necessary information for the environment " + exportEnvironment)
		}
	},
}


func ExportAPI(name string, version string, url string, accessToken string) *resty.Response {
	// append '/' to the end if there isn't one already
	if string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "export/apis"

	query := "?query=" + name
	url += query
	fmt.Println("ExportAPI: URL:", url)
	headers := make(map[string]string)
	headers[constants.HeaderAuthorization] = constants.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[constants.HeaderAccept] = constants.HeaderValueApplicationZip

	resp, err := resty.R().
		SetHeaders(headers).
		Get(url)

	if err != nil {
		fmt.Println("Error exporting API:", name)
		panic(err)
	}

	return resp
}


func init() {
	RootCmd.AddCommand(ExportAPICmd)
	ExportAPICmd.Flags().StringVarP(&exportAPIName, "name", "n", "", "Name of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportAPIVersion, "version", "v", "", "Version of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportEnvironment, "environment", "e", "", "Environment to which the API should be exported")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ExportAPICmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ExportAPICmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
