/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package utils

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty"
	"net/http"
	"os"
	"strings"
)

// ExecutePreCommand deals with generating tokens needed for executing a particular command
// @param environment : Environment on which the particular command is run
// @param flagUsername : Username entered using the flag --username (-u). Could be blank
// @param flagPassword : Password entered using the flag --password (-p). Could be blank
// @return AccessToken, PublisherEndpoint, Errors
// including (export-api, import-api, list)
func ExecutePreCommand(environment string, flagUsername string, flagPassword string) (string, string, error) {
	if EnvExistsInMainConfigFile(environment, MainConfigFilePath) {
		registrationEndpoint := GetRegistrationEndpointOfEnv(environment, MainConfigFilePath)
		apiManagerEndpoint := GetAPIMEndpointOfEnv(environment, MainConfigFilePath)
		tokenEndpoint := GetTokenEndpointOfEnv(environment, MainConfigFilePath)

		Logln(LogPrefixInfo + "Environment: '" + environment + "'")
		Logln(LogPrefixInfo + "Reg Endpoint read:", registrationEndpoint)

		var username string
		var password string
		var clientID string
		var clientSecret string
		var err error

		if EnvExistsInKeysFile(environment, EnvKeysAllFilePath) {
			// client_id, client_secret, and username exist in file
			username = GetUsernameOfEnv(environment, EnvKeysAllFilePath)

			if flagUsername != "" {
				// flagUsername is not blank
				if flagUsername != username {
					// username entered with flag -u is not the same as username found
					// in env_keys_all.yaml file
					fmt.Println(LogPrefixWarning + "Username entered with flag -u for the environment '" + environment +
						"' is not the same as username found in file '" + EnvKeysAllFilePath + "'")
					fmt.Println("Execute '" + ProjectName + " reset-user -e " + environment + "' to clear user data")
					os.Exit(1)
				} else {
					// username entered with flag -u is the same as username found in env_keys_all.yaml file
					if flagPassword == "" {
						fmt.Println("For Username: " + username)
						password = PromptForPassword()
					} else {
						// flagPassword is not blank
						// no need of prompting for password now
						password = flagPassword
					}
				}
			} else {
				// flagUsername is blank
				if flagPassword != "" {
					// flagPassword is not blank
					password = flagPassword
				} else {
					// flagPassword is blank
					fmt.Println("For username: " + username)
					password = PromptForPassword()
				}
			}

			clientID = GetClientIDOfEnv(environment, EnvKeysAllFilePath)
			clientSecret = GetClientSecretOfEnv(environment, password, EnvKeysAllFilePath)

			Logln(LogPrefixInfo+"Username:", username)
			Logln(LogPrefixInfo+"ClientID:", clientID)
		} else {
			// env exists in endpoints file, but not in keys file
			// no client_id, client_secret in file
			// first use of the environment
			// Get new values

			if flagUsername != "" {
				// flagUsername is not blank
				username = flagUsername
				if flagPassword == "" {
					// flagPassword is blank
					fmt.Println("For Username: " + username)
					password = PromptForPassword()
				} else {
					// flagPassword is not blank
					password = flagPassword
				}
			} else {
				// flagUsername is blank
				// doesn't matter is flagPassword is blank or not
				username = strings.TrimSpace(PromptForUsername())
				password = PromptForPassword()
			}

			fmt.Println("\nUsername: " + username + "\n")
			clientID, clientSecret, err = GetClientIDSecret(username, password, registrationEndpoint)

			if err != nil {
				fmt.Println("Error:", err)
			}

			// Persist clientID, clientSecret, Username in file
			encryptedClientSecret := Encrypt([]byte(GetMD5Hash(password)), clientSecret)
			envKeys := EnvKeys{clientID, encryptedClientSecret, username}
			AddNewEnvToKeysFile(environment, envKeys, EnvKeysAllFilePath)
		}

		// Get OAuth Tokens
		responseDataMap, _ := GetOAuthTokens(username, password,
			GetBase64EncodedCredentials(clientID, clientSecret), tokenEndpoint)
		accessToken := responseDataMap["access_token"]

		Logln(LogPrefixInfo+"[Remove in Production] AccessToken:", accessToken) // TODO:: Remove in production

		return accessToken, apiManagerEndpoint, nil
	} else {
		// env does not exist in main config file
		if environment == "" {
			return "", "", errors.New("no environment specified. Either specify it using the -e flag or name one of " +
				"the environments in '" + MainConfigFileName + "' to 'default'")
		}

		return "", "", errors.New("Details incorrect/unavailable for environment '" + environment + "' in " +
			MainConfigFilePath)
	}
}

// GetClientIDSecret implemented using go-resty
// @param username : Username for application server account
// @param password : Password for application server account
// @param url : Registration Endpoint for the environment
// @return client_id, client_secret, error
func GetClientIDSecret(username string, password string, url string) (string, string, error) {
	body := `{"clientName": "Test", "redirect_uris": "www.google.lk", "grant_types":"password"}`
	headers := make(map[string]string)

	headers[HeaderContentType] = HeaderValueApplicationJSON
	// headers["Content-Type"] = "application/json"

	headers[HeaderAuthorization] = HeaderValueAuthBasicPrefix + " " + GetBase64EncodedCredentials(username, password)
	// headers["Authorization"] = "Basic " + GetBase64EncodedCredentials(username, password)


	// POST request using resty
	resp, err := InvokePOSTRequest(url, headers, body)

	if err != nil {
		HandleErrorAndExit("Error in connecting.", err)
	}

	Logln("Getting ClientID, ClientSecret: Status - " + resp.Status())

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		registrationResponse := RegistrationResponse{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &registrationResponse)

		clientID := registrationResponse.ClientID
		clientSecret := registrationResponse.ClientSecret

		return clientID, clientSecret, err

	} else {
		//fmt.Println("Error:", resp.Error())
		//fmt.Printf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			HandleErrorAndExit("Incorrect Username/Password combination.", errors.New("401 Unauthorized"))
		}
		return "", "", errors.New("Request didn't respond 200 OK: " + resp.Status())
	}
}

// Encode the concatenation of two strings (using ":")
// provide two strings
// returns base64Encode(key:secret)
func GetBase64EncodedCredentials(key string, secret string) string {
	line := key + ":" + secret
	encoded := base64.StdEncoding.EncodeToString([]byte(line))
	return encoded
}

// GetOAuthTokens implemented using go-resty/resty
// provide username, password, and validity period for the access token
// returns the response as a map
func GetOAuthTokens(username string, password string,
	b64EncodedClientIDClientSecret string, url string) (map[string]string, error) {
	validityPeriod := DefaultTokenValidityPeriod
	body := "grant_type=password&username=" + username + "&password=" + password + "&validity_period=" + validityPeriod

	// set headers
	headers := make(map[string]string)
	headers[HeaderContentType] = HeaderValueXWWWFormUrlEncoded
	headers[HeaderAuthorization] = HeaderValueAuthBearerPrefix + " " + b64EncodedClientIDClientSecret
	headers[HeaderAccept] = HeaderValueApplicationJSON

	if SkipTLSVerification {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTP certificates
	}

	resp, err := InvokePOSTRequest(url, headers, body)

	if err != nil {
		Logln(LogPrefixError + "connecting to " + url)
		HandleErrorAndExit("Unable to Connect.", err)
	}

	if resp.StatusCode() != http.StatusOK {
		HandleErrorAndExit("Unable to connect.", errors.New("Status: "+resp.Status()))
		return nil, nil
	}

	responseDataMap := make(map[string]string) // a map to hold response data
	data := []byte(resp.Body())
	_ = json.Unmarshal(data, &responseDataMap) // add response data to the map

	return responseDataMap, nil // contains 'access_token', 'refresh_token' etc
}
