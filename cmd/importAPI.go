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
	"crypto/tls"
	"net/http"
	"io/ioutil"
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

			var username string
			var password string
			var clientID string
			var clientSecret string

			if utils.EnvExistsInKeysFile(importEnvironment) {
				// client_id, client_secret,username exists in file
				// get username from file
				username = utils.GetUsernameOfEnv(importEnvironment)
				fmt.Println("Username:", username)

				// get client_id from file
				clientID = utils.GetClientIDOfEnv(importEnvironment)

				// get client_secret from file, password needed to decrypt client_secret
				password = utils.PromptForPassword()
				clientSecret = utils.GetClientSecretOfEnv(importEnvironment, password)

				fmt.Println("ClientID:", clientID)
				fmt.Println("ClientSecret:", clientSecret)
			} else {
				// env exists in endpoints file, but not in keys file (first use of the tool)
				// no client_id, client_secret in file
				// Get new values
				username = strings.TrimSpace(utils.PromptForUsername())
				password = utils.PromptForPassword()

				fmt.Println("\nUsername: " + username + "\n")
				clientID, clientSecret = utils.GetClientIDSecret(username, password, registrationEndpoint)

				// Persist clientID, clientSecret, Username in file
				encryptedClientSecret := utils.Encrypt([]byte(utils.GetMD5Hash(password)), clientSecret)
				envKeys := utils.EnvKeys{clientID, encryptedClientSecret, username}
				utils.AddNewEnvToKeysFile(importEnvironment, envKeys)
			}

			// Get OAuth Tokens
			m := utils.GetOAuthTokens(username, password, utils.GetBase64EncodedCredentials(clientID, clientSecret), tokenEndpoint)
			accessToken := m["access_token"]
			fmt.Println("AccessToken:", accessToken)

			resp := ImportAPI(importAPIName, importAPIVersion, apiManagerEndpoint, accessToken)
			fmt.Printf("Status: %v\n", resp.Status)
			//fmt.Printf("Errors: %v\n", resp.Error)
			fmt.Println("Header:", resp.Header)
			fmt.Printf("Body: %s\n", resp.Body)
		} else {
			// env_endpoints_all.yaml file is not configured properly by the user
			log.Fatal("Error: env_endpoints_all.yaml does not contain necessary information for the environment " + importEnvironment)
		}
	},
}


func ImportAPI(name string, version string, url string, accessToken string) *http.Response{
	// append '/' to the end if there isn't one already
	if string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "imports/api"

	file := "./exported/" + name + ".zip"
	fmt.Println("File:", file)
	fmt.Println("ImportAPI: URL:", url)
	//headers[constants.HeaderConsumes]  = constants.HeaderValueMultiPartFormData

	//openFile, _ := ioutil.ReadFile(file)
	//osOpen, err := os.Open(file)
	//if err != nil {
	//	fmt.Println("Error opening file:")
	//	panic(err)
	//}

	payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"file\"; filename=\"" + file + "\r\nContent-Type: application/zip\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")

	req, _ := http.NewRequest("PUT", url, payload)

	req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	req.Header.Add("Authorization", "Bearer " + accessToken)
	 //req.Header.Add("cache-control", "no-cache")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, _ := client.Do(req)

	if resp != nil {

		body, _ := ioutil.ReadAll(resp.Body)

		fmt.Println(resp)
		fmt.Println(string(body))
	} else {
		fmt.Println("Null Response")
	}

	return resp
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
