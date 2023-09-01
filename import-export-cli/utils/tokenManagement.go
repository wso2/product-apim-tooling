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

package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	encodeURL "net/url"
	"strings"

	"github.com/renstrom/dedent"
)

// ExecutePreCommandWithBasicAuth deals with generating tokens needed for executing a particular command
// @param environment : Environment on which the particular command is run
// @param flagUsername : Username entered using the flag --username (-u). Could be blank
// @param flagPassword : Password entered using the flag --password (-p). Could be blank
// @return b64encodedCredentials, ApiManagerEndpoint, Errors
// including (export-api, import-api)
func ExecutePreCommandWithBasicAuth(environment, flagUsername, flagPassword, mainConfigFilePath,
	envKeysAllFilePath string) (b64encodedCredentials string, err error) {
	if EnvExistsInMainConfigFile(environment, mainConfigFilePath) {
		Logln(LogPrefixInfo + "Environment: '" + environment + "'")

		var username string
		var password string
		var err error

		if EnvExistsInKeysFile(environment, envKeysAllFilePath) {
			// client_id, client_secret, and username exist in file
			username = GetUsernameOfEnv(environment, envKeysAllFilePath)

			if flagUsername != "" {
				// flagUsername is not blank
				if flagUsername != username {
					// username entered with flag -u is not the same as username found
					// in env_keys_all.yaml file
					fmt.Println("Username entered with flag -u for the environment '" +
						environment + "' is not the same as username found in file '" + EnvKeysAllFilePath + "'")
					fmt.Println("Execute '" + ProjectName + " reset-user -e " + environment + "' to clear user data")
					HandleErrorAndExit("Username mismatch", nil)
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

			if err != nil {
				fmt.Println("Error:", err)
			}
		}

		// Get Base64 Encoded Username:Password
		b64encodedCredentials := GetBase64EncodedCredentials(username, password)

		return b64encodedCredentials, nil
	} else {
		// env does not exist in main config file
		if environment == "" {
			return "",
				errors.New("no environment specified. Either specify it using the -e flag or name one of " +
					"the environments in '" + MainConfigFileName + "' to 'default'")
		}

		return "",
			errors.New("Details incorrect/unavailable for environment '" + environment + "' in " + mainConfigFilePath)
	}
}

// ExecutePreCommandWithOAuth deals with generating tokens needed for executing a particular command
// @param environment : Environment on which the particular command is run
// @param flagUsername : Username entered using the flag --username (-u). Could be blank
// @param flagPassword : Password entered using the flag --password (-p). Could be blank
// @return AccessToken, ApiManagerEndpoint, Errors
// including (export-api, import-api, list)
func ExecutePreCommandWithOAuth(environment, flagUsername, flagPassword, mainConfigFilePath,
	envKeysAllFilePath string) (accessToken string, err error) {
	if EnvExistsInMainConfigFile(environment, mainConfigFilePath) {
		registrationEndpoint := GetRegistrationEndpointOfEnv(environment, mainConfigFilePath)
		tokenEndpoint := GetInternalTokenEndpointOfEnv(environment, mainConfigFilePath)

		Logln(LogPrefixInfo + "Environment: '" + environment + "'")
		Logln(LogPrefixInfo+"Reg Endpoint read:", registrationEndpoint)

		var username string
		var password string
		var clientID string
		var clientSecret string
		var err error

		if EnvExistsInKeysFile(environment, envKeysAllFilePath) {
			// client_id, client_secret, and username exist in file
			username = GetUsernameOfEnv(environment, envKeysAllFilePath)

			if flagUsername != "" {
				// flagUsername is not blank
				if flagUsername != username {
					// username entered with flag -u is not the same as username found
					// in env_keys_all.yaml file
					fmt.Println("Username entered with flag -u for the environment '" +
						environment + "' is not the same as username found in file '" + EnvKeysAllFilePath + "'")
					fmt.Println("Execute '" + ProjectName + " reset-user -e " + environment + "' to clear user data")
					HandleErrorAndExit("Username mismatch", nil)
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

			clientID = GetClientIDOfEnv(environment, envKeysAllFilePath)
			clientSecret = GetClientSecretOfEnv(environment, password, envKeysAllFilePath)

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
				// doesn't matter if flagPassword is blank or not
				username = strings.TrimSpace(PromptForUsername())
				password = PromptForPassword()
			}

			clientID, clientSecret, err = GetClientIDSecret(username, password, registrationEndpoint)

			if err != nil {
				HandleErrorAndExit("Error:", err)
			}

			// Persist clientID, clientSecret, Username in file
			encryptedClientSecret := Encrypt([]byte(GetMD5Hash(password)), clientSecret)
			envKeys := EnvKeys{clientID, encryptedClientSecret, username}
			AddNewEnvToKeysFile(environment, envKeys, envKeysAllFilePath)
		}

		// Get OAuth Tokens
		responseDataMap, _ := GetOAuthTokens(username, password,
			GetBase64EncodedCredentials(clientID, clientSecret), tokenEndpoint)
		accessToken := responseDataMap["access_token"]

		return accessToken, nil
	} else {
		// env does not exist in main config file
		if environment == "" {
			return "",
				errors.New("no environment specified. Either specify it using the -e flag or rename one of " +
					"the environments in '" + MainConfigFileName + "' to 'default'")
		}

		return "", errors.New("No environment specified or details incorrect/unavailable for environment '" +
			environment + "' in " + mainConfigFilePath + "")
	}
}

// GetClientIDSecret implemented using go-resty
// @param username : Username for application server account
// @param password : Password for application server account
// @param url : Registration Endpoint for the environment
// @return client_id, client_secret, error
func GetClientIDSecret(username, password, url string) (clientID string, clientSecret string, err error) {
	updatedUsername := strings.ReplaceAll(username, "@", "_")
	body := dedent.Dedent(`{"clientName": "rest_api_import_export_` + updatedUsername +  `",
								  "callbackUrl": "www.google.lk",
								  "grantType":"password refresh_token",
								  "saasApp": true,
								  "owner": "` + username + `",
								  "tokenScope": "Production"
							     }`)
	headers := make(map[string]string)

	headers[HeaderContentType] = HeaderValueApplicationJSON

	headers[HeaderAuthorization] = HeaderValueAuthBasicPrefix + " " + GetBase64EncodedCredentials(username, password)

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
		Logf("Error: %s\n", resp.Error())
		Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", "",
				fmt.Errorf("authorization failed during CLI client registration process")
		}
		return "", "", errors.New("Request didn't respond 200 OK: " + resp.Status())
	}
}

// Encode the concatenation of two strings (using ":")
// provide two strings
// returns base64Encode(key:secret)
func GetBase64EncodedCredentials(key, secret string) (encodedValue string) {
	line := key + ":" + secret
	encoded := base64.StdEncoding.EncodeToString([]byte(line))
	return encoded
}

// GetOAuthTokens implemented using go-resty/resty
// @param username
// @param password
// @param b64EncodedClientIDClientSecret
// @param url : OAuth token endpoint
// @return response as a map
// @return error
func GetOAuthTokens(username, password, b64EncodedClientIDClientSecret, url string) (map[string]string, error) {
	body := "grant_type=password&username=" + username + "&password=" + encodeURL.QueryEscape(password) +
		"&scope=apim:app_import_export+apim:api_import_export+apim:api_product_import_export+apim:app_manage+" +
		"apim:sub_manage+apim:api_view+apim:api_delete+apim:app_owner_change+apim:subscribe+apim:api_publish+" +
		"apim:admin+apim:policies_import_export"

	// set headers
	headers := make(map[string]string)
	headers[HeaderContentType] = HeaderValueXWWWFormUrlEncoded
	headers[HeaderAuthorization] = HeaderValueAuthBasicPrefix + " " + b64EncodedClientIDClientSecret
	headers[HeaderAccept] = HeaderValueApplicationJSON

	Logln(LogPrefixInfo + "connecting to " + url)
	resp, err := InvokePOSTRequest(url, headers, body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("Unable to connect. " +
			"Status: " + resp.Status())
	}

	responseDataMap := make(map[string]string) // a map to hold response data
	data := []byte(resp.Body())
	json.Unmarshal(data, &responseDataMap) // add response data to the map

	return responseDataMap, nil // contains 'access_token', 'refresh_token' etc
}
