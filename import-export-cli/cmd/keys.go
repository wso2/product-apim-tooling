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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	encodeURL "net/url"
	"strings"
)

// keys command related Info
const genKeyCmdLiteral = "get-keys"
const genKeyCmdShortDesc = "Generate access token to invoke the API"
const genKeyCmdLongDesc = `Generate JWT token to invoke the API by subscribing to a default application for testing purposes`
const genKeyCmdExamples = utils.ProjectName + " " + genKeyCmdLiteral + ` -n TwitterAPI -v 1.0.0 -e dev --provider admin`

var keyGenEnv string
var apiName string
var apiVersion string
var apiProvider string
var tokenType string
var throttlingTier string

var genKeyCmd = &cobra.Command{
	Use:     genKeyCmdLiteral,
	Short:   genKeyCmdShortDesc,
	Long:    genKeyCmdLongDesc,
	Example: genKeyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {

		utils.Logln(utils.LogPrefixInfo + genKeyCmdLiteral + " called")
		getKeys()
	},
}

//Subscribe the given API to the default application and generate an access token
func getKeys() {
	cred, err := getCredentials(keyGenEnv)
	if err != nil {
		utils.HandleErrorAndExit("Error getting credentials", err)
	}
	utils.Logln(utils.LogPrefixInfo + "Retrieved credentials of the environment successfully")
	//Calling the DCR endpoint to get the credentials of the env
	cred.ClientId, cred.ClientSecret, err = callDCREndpoint(cred)
	//If the DCR call fails exit with the error
	if err != nil {
		utils.HandleErrorAndExit("Internal error occurred", err)
	}
	utils.Logln(utils.LogPrefixInfo + "Called DCR endpoint successfully")
	//generating access token for the env based on the credentials
	accessToken, err := generateAccessToken(cred)
	if err != nil {
		utils.HandleErrorAndExit("Internal error occurred", err)
	}
	utils.Logln(utils.LogPrefixInfo + "Generated accesstoken to call the store rest APIs.")
	//retrieving subscription tiers
	tiers, err := getAvailableAPITiers(accessToken)

	if tiers != nil && err == nil {
		utils.Logln(utils.LogPrefixInfo + "Retrieved available subscription tiers of the API: ", tiers)
		//Needs an available subscription tier when creating application
		throttlingTier = tiers[0]
	} else {
		utils.HandleErrorAndExit("Please check the API details and try again.", err)
	}
	//search if the default cli application already exists
	appId, err := searchApplication(utils.DefaultCliApp, accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Internal error occurred", err)
	}
	utils.Logln(utils.LogPrefixInfo + "Searched if application exists.")
	//if the application exists
	if appId != "" {
		utils.Logln(utils.LogPrefixInfo + "Application already exists")
		//Search the if the given API is present
		subId, err := subscribe(appId, accessToken)
		//If subscrition fails
		if subId == "" && err != nil {
			utils.HandleErrorAndExit("Error occurred while subscribing.", err)
		}

		scopes, err := getScopes(appId, accessToken)
		//retrieve application specific details
		appDetails, err := getApplicationDetails(appId, accessToken)
		if appDetails != nil {
			//Reading configuration to check if the application needs to be updated
			configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
			tokenType = configVars.Config.TokenType
			//Check if the token type of the application has been updated
			if tokenType != appDetails.TokenType && tokenType != "" {
				//Request body for the store rest api
				appUpdateReq := utils.AppCreateRequest{
					Name:             utils.DefaultCliApp,
					ThrottlingPolicy: throttlingTier,
					Description:      "Default CLI Application",
					TokenType:        configVars.Config.TokenType,
				}
				body, err := json.Marshal(appUpdateReq);
				if body == nil && err != nil {
					utils.HandleErrorAndExit("Error occurred while creating CLI application update request.", err)
				}
				utils.Logln(utils.LogPrefixInfo + "Updating application as token type is changed to: " + tokenType)
				updatedApp, updateError := updateApplicationDetails(appId, string(body), accessToken)

				if updatedApp != nil && updateError == nil {
					utils.Logln(utils.LogPrefixInfo + "Updated application successfully")
				} else if updateError != nil {
					fmt.Println("Error while updating the application. : ", updateError)
				}
			}
			//retrieve keys of application to see if there are already generated keys
			appKeys, keysErr := getApplicationKeys(appId, accessToken)
			if keysErr != nil {
				utils.HandleErrorAndExit("Error occurred while getting application keys.", keysErr)
			}

			//if keys have been already generated before, then update the consumer key and secret
			if appKeys.Count != 0 {
				keygenResponse, keyGenErr := regenerateConsumerSecret(appId, "PRODUCTION", accessToken)
				if keyGenErr != nil {
					fmt.Println("Error occurred while regenerating the keys for the app ", appId)
				} else {
					appKeys.List[0].ConsumerSecret = keygenResponse.ConsumerSecret
					utils.Logln(utils.LogPrefixInfo + "Regenerated application keys successfully")
				}

				//If the keys have not been generated and the application is updated
				token, err := getNewToken(&appKeys.List[0], scopes)
				if accessToken != "" {
					fmt.Println("Access Token: ", token)
				} else {
					fmt.Println("Error while generating token: ", err)
				}
			} else {
				//If the application is already created but the keys have not generated in the first time
				keygenResponse, err := generateApplicationKeys(appId, accessToken)
				if keygenResponse == nil && err != nil {
					utils.HandleErrorAndExit("Error occurred while generating application keys.", err)
				}
				fmt.Println("Access Token: ", keygenResponse.Token.AccessToken)
			}
		} else {
			fmt.Println("Error while retrieving the application:", err)
		}
	} else {
		//If the default cli appId does not exist in the environment
		//Create the application
		createdAppId, appName, err := createApplication(accessToken)
		appId = createdAppId
		if createdAppId != "" || appName != "" {
			utils.Logln(utils.LogPrefixInfo + "Created application: ", appName)
		} else {
			//if error occurred while creating the application, then
			utils.HandleErrorAndExit("Error while creating the application:", err)
		}
		//Search the if the given API is present
		subId, err := subscribe(appId, accessToken)
		//If subscription failed
		if subId == "" && err != nil {
			utils.HandleErrorAndExit("Error occurred while subscribing.", err)
		}
		scopes, err := getScopes(appId, accessToken)
		//If errors occurred while retrieving scopes
		if scopes == nil && err != nil {
			utils.HandleErrorAndExit("Error while retrieving scopes ", err)
		}
		//Generate the tokens
		keygenResponse, err:= generateApplicationKeys(appId, accessToken)
		if err != nil {
			utils.HandleErrorAndExit("Error while generating application keys", err)
		}
		appKey := &utils.ApplicationKey{}
		appKey.ConsumerKey = keygenResponse.ConsumerKey;
		appKey.ConsumerSecret = keygenResponse.ConsumerSecret;
		token, err := getNewToken(appKey, scopes)
		if token != ""  {
			fmt.Println("Access Token: ", token)
		} else {
			fmt.Println("Error while generating token: ", err)
		}
	}
}

func getAvailableAPITiers(accessToken string) ([]string, error) {
	apiId, err := searchApi(accessToken)
	if apiId == "" && err != nil {
		return nil, err
	}
	api, err := getApi(apiId, accessToken)
	if err == nil && api != nil {
		return api.Policies, err
	} else {
		return nil, err
	}
}

// Calling DCR endpoint
// @param credential : Username and Password
// @return client_id, client_secret, error
func callDCREndpoint(credential credentials.Credential) (string, string, error) {
	//Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	//Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	//Request body for the store rest api
	body := dedent.Dedent(`{
								"clientName": "rest_api_store",
							   	"callbackUrl": "www.google.lk",
							   	"grantType":"password refresh_token",
							   	"saasApp": true,
							   	"owner": "` + credential.Username + `"
							}`)
	registrationEndpoint := utils.GetRegistrationEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	//Calling the DCR endpoint
	resp, err := utils.InvokePOSTRequest(registrationEndpoint, headers, body)
	if err != nil {
		utils.HandleErrorAndExit("DCR request failed. Reason: ", err)
	}

	utils.Logln(utils.LogPrefixInfo + "Getting ClientID, ClientSecret: Status - " + resp.Status())

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		registrationResponse := &utils.RegistrationResponse{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &registrationResponse)
		//Retrieving client credentials
		clientID := registrationResponse.ClientID
		clientSecret := registrationResponse.ClientSecret
		return clientID, clientSecret, err

	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", "", fmt.Errorf("invalid username/password combination")
		}
		return "", "", errors.New("Request didn't respond 200 OK for DCR request. Status: " + resp.Status())
	}
}

// Get tokens for
// @param credential : ClientID and ClientSecret
// @return accessToken, error
func generateAccessToken(credential credentials.Credential) (string, error) {
	//Base64 encoding the credentials
	b64encodedCredentials := credentials.Base64Encode(fmt.Sprintf("%s:%s", credential.ClientId, credential.ClientSecret))
	//Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderContentType] = utils.HeaderValueXWWWFormUrlEncoded
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON

	//Retrieving the token endpoint of the relevant environment
	tokenEndpoint := utils.GetTokenEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	//Prepping query params
	body := "grant_type=password&username=" + credential.Username + "&password=" +
		encodeURL.QueryEscape(credential.Password) + "&validity_period=" + string(utils.DefaultTokenValidityPeriod) +
		"&scope=apim:api_view+apim:subscribe+apim:api_publish"

	//Call to the token endpoint with the necessary payload
	resp, err := utils.InvokePOSTRequest(tokenEndpoint, headers, body)
	//If the response is erroneous
	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+tokenEndpoint, err)
	}
	//Logging the response
	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())
	//If the token generation response is success
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		keygenResponse := &utils.TokenResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &keygenResponse)
		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return keygenResponse.AccessToken, err
	} else {
		return "", errors.New("Request didn't respond 200 OK for generating an access token. Status: " + resp.Status())
	}
}

// Regenerate consumer secret of the application
// @param appId : ID of the application
// @param keyType : key Type of the application. Allowed values: PRODUCTION, SANDBOX
// @return KeygenResponse, error
func regenerateConsumerSecret(appId string, keyType string, accessToken string) (*utils.ConsumerSecretRegenResponse,
		error) {
	applicationEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	url := applicationEndpoint + "/" + appId + "/keys/" + keyType + "/regenerate-secret"
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON

	resp, err := utils.InvokePOSTRequestWithoutBody(url, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		keygenResp := &utils.ConsumerSecretRegenResponse{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &keygenResp)

		return keygenResp, err

	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return nil, fmt.Errorf("invalid username/password combination")
		}
		return nil, errors.New("Request didn't respond 200 OK for regenerating the consumer secret. " +
			"Status: " + resp.Status())
	}
}

// Search if the application exists with the name
// @param appName : Name of the application
// @param accessToken : Access token to authenticate the store REST API
// @return appId, error
func searchApplication(appName string, accessToken string) (string, error) {
	//Application rest API endpoint of the environment from the config file
	applicationEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	//Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequestWithQueryParam("query", appName, applicationEndpoint, headers)

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		appData := &utils.AppList{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &appData)
		if appData.Count != 0 {
			appId := appData.List[0].ApplicationID
			return appId, err
		}
		return "", nil

	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("invalid username/password combination")
		}
		return "", errors.New("Request didn't respond 200 OK for searching existing applications. " +
			"Status: " + resp.Status())
	}
}

// Searching if the API is available
// @param accessToken : token to call the store rest API
// @return apiId, error
func searchApi(accessToken string) (string, error) {
	//API endpoint of the environment from the config file
	apiEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	apiEndpoint = strings.Replace(apiEndpoint, "applications", "apis", -1)

	//Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	//headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	var queryVal string
	if apiName != "" {
		queryVal = "name:\"" + apiName + "\""
		if apiVersion != "" {
			queryVal = queryVal + " version:\"" + apiVersion + "\""
		}
		if apiProvider != "" {
			queryVal = queryVal + " provider:\"" + apiProvider + "\""
		}
	}
	resp, err := utils.InvokeGETRequestWithQueryParam("query", queryVal, apiEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		apiData := &utils.ApiSearch{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &apiData)
		if apiData.Count != 0 {
			apiId := apiData.List[0].ID
			return apiId, err
		}
		return "", errors.New("Requested API is not available in the store. API: " + apiName +
										" Version: " + apiVersion + " Provider: " + apiProvider)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("invalid username/password combination")
		}
		return "", errors.New("Request didn't respond 200 OK for searching APIs. Status: " + resp.Status())
	}
}

// Subscribe API to a given application
// @param appId : Appplication ID to subscribe the API
// @param accessToken : Token to call rest API
// @return subscriptionId, error
func subscribe(appId string, accessToken string) (string, error) {
	apiId, err := searchApi(accessToken)
	if apiId != "" && err == nil {
		//If the API is present, subscribe that API to the application
		fmt.Println("API name: ", apiName, "& version: ", apiVersion, "exists")
		subId, err := subscribeApi(apiId, appId, accessToken)
		if subId != "" {
			fmt.Println("API ", apiName, ":", apiVersion, "subscribed successfully.")
		} else {
			fmt.Println("Error while subscribing to the application:", err)
		}
		return subId, err
	} else {
		return "", errors.New("API is not found. Name: " + apiName + " version: " + apiVersion)
	}
}

// Get API specific details of a given API
// @param apiId : API ID to retrieve the information
// @param accessToken : token to call the rest API
// @return API, error
func getApi(apiId string, accessToken string) (*utils.APIData, error) {
	apiEndpoint := utils.GetApiListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) + "/" + apiId
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(apiEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		apiData := &utils.APIData{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &apiData)
		return apiData, err
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return nil, fmt.Errorf("invalid username/password combination")
		}
		return nil, errors.New("Request didn't respond 200 OK for retrieving API details. Status: " + resp.Status())
	}
}

// Subscribe the API to a given Application
// @param apiId : apiId to be subscribed
// @param appId : appId to be subscribed
// @param accessToken : token to call the rest API
// @return subscriptionId, error
func subscribeApi(apiId string, appId string, accessToken string) (string, error) {
	//todo: subscription endpoint to be included in conf
	subEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	subEndpoint = strings.Replace(subEndpoint, "applications", "subscriptions", -1)
	//prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	//Prepping query parameters
	queryParams := map[string]string{
		utils.ApiId: apiId}
	//Checking if there is a subscription of given API to the give application
	subResp, subErr := utils.InvokeGETRequestWithMultipleQueryParams(queryParams, subEndpoint, headers)

	if subResp.StatusCode() == http.StatusOK || subResp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		subscription := &utils.SubscriptionList{}
		data := []byte(subResp.Body())
		subErr = json.Unmarshal(data, &subscription)
		if subscription.Count != 0 {

			for _, sub := range subscription.List {
				//If an subscription exists, then return the subscription ID
				if sub.ApplicationID == appId {
					return sub.SubscriptionID, subErr
				}
			}
		}
		subscriptionReq := &utils.SubscriptionCreateRequest{
			APIID:            apiId,
			ApplicationID:    appId,
			ThrottlingPolicy: throttlingTier,
		}
		//If there is no subscription, make a subscription
		body, err := json.Marshal(subscriptionReq);
		if body == nil && err != nil {
			utils.HandleErrorAndExit("Error occurred while creating CLI application subscription request.", err)
		}
		resp, err := utils.InvokePOSTRequest(subEndpoint, headers, string(body))
		if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
			// 200 OK or 201 Created
			subscription := &utils.Subscription{}
			data := []byte(resp.Body())
			err = json.Unmarshal(data, &subscription)
			return subscription.SubscriptionID, err
		} else {
			utils.Logf("Error: %s\n", resp.Error())
			utils.Logf("Body: %s\n", resp.Body())
			if resp.StatusCode() == http.StatusUnauthorized {
				// 401 Unauthorized
				return "", fmt.Errorf("invalid username/password combination")
			}
			return "", errors.New("Request didn't respond 200 OK for subscribing to the API. Status: " + resp.Status())
		}
	} else {
		utils.Logf("Error: %s\n", subResp.Error())
		utils.Logf("Body: %s\n", subResp.Body())
		if subResp.StatusCode() == http.StatusUnauthorized {
			return "", fmt.Errorf("invalid username/password combination")
		}
		return "", errors.New("Request didn't respond 200 OK: " + subResp.Status())
	}
}

// Get application details
// @param appId : Application ID
// @param accessToken : token to call the store rest API
// @return AppDetails, error
func getApplicationDetails(appId string, accessToken string) (*utils.AppDetails, error) {

	applicationEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) + "/" + appId
	//Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	//Retrieving the details of the particular application
	resp, err := utils.InvokeGETRequest(applicationEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		appData := &utils.AppDetails{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &appData)
		return appData, err
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return nil, fmt.Errorf("invalid username/password combination")
		}
		return nil, errors.New("Request didn't respond 200 OK for retrieving application details. " +
			"Status: " + resp.Status())
	}
}

// Get application keys
// @param appId : Application ID
// @param accessToken : token to call the store rest API
// @return AppDetails, error
func getApplicationKeys(appId string, accessToken string) (*utils.AppKeyList, error) {

	applicationEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) +
		"/" + appId + "/keys"
	//Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	//Retrieving the details of the particular application
	resp, err := utils.InvokeGETRequest(applicationEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		keyData := &utils.AppKeyList{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &keyData)
		return keyData, err
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return nil, fmt.Errorf("invalid username/password combination")
		}
		return nil, errors.New("Request didn't respond 200 OK for retrieving App key information. " +
			"Status: " + resp.Status())
	}
}

// Update application details
// @param appId : Application ID
// @param accessToken : token to call the store rest API
// @return AppDetails, error
func updateApplicationDetails(appId string, body string, accessToken string) (*utils.AppDetails, error) {

	applicationEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) + "/" + appId
	//Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON

	//Retrieving the details of the particular application
	resp, err := utils.InvokePutRequest(nil, applicationEndpoint, headers, body)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		appData := &utils.AppDetails{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &appData)
		return appData, err
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return nil, fmt.Errorf("invalid username/password combination")
		}
		return nil, errors.New("Request didn't respond 200 OK for updating application. Status: " + resp.Status())
	}
}

// Create application with a default name in a given environment
// @param accessToken : access token to call the store rest API
// @return client_id, client_secret, error
func createApplication(accessToken string) (string, string, error) {

	applicationEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	conf := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	appUpdateReq := utils.AppCreateRequest{
		Name:             utils.DefaultCliApp,
		ThrottlingPolicy: throttlingTier,
		Description:      "Default application for apictl testing purposes",
		TokenType:        conf.Config.TokenType,
	}
	body, err := json.Marshal(appUpdateReq);
	if body == nil && err != nil {
		utils.HandleErrorAndExit("Error occurred while creating CLI application update request.", err)
	}
	resp, err := utils.InvokePOSTRequest(applicationEndpoint, headers, string(body))
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		applicationResponse := &utils.Application{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &applicationResponse)
		appId := applicationResponse.ID
		appName := applicationResponse.Name
		return appId, appName, err

	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", "", fmt.Errorf("invalid username/password combination")
		}
		return "", "", errors.New("Request didn't respond 200 OK for application creation. Status: " + resp.Status())
	}
}

// Calling token endpoint to get access token for the already created application
// @param key : Details of the particular key
// @param scopes[] : scopes to generate the token
// @return accessToken, error
func getNewToken(key *utils.ApplicationKey, scopes []string) (string, error) {
	tokenEndpoint := utils.GetTokenEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)

	body := "grant_type=client_credentials&scope=" + strings.Join(scopes, " ")

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " +
		utils.GetBase64EncodedCredentials(key.ConsumerKey, key.ConsumerSecret)

	headers[utils.HeaderContentType] = utils.HeaderValueXWWWFormUrlEncoded
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON

	resp, err := utils.InvokePOSTRequest(tokenEndpoint, headers, body)

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		keygenResponse := &utils.TokenResponse{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &keygenResponse)

		accessToken := keygenResponse.AccessToken
		return accessToken, err

	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("invalid username/password combination")
		}
		return "", errors.New("Request didn't respond 200 OK for generating a new token. Status: " + resp.Status())
	}

}

// Get all the scopes of the APIs subscribed to a particular application
// @param appId : application ID to get the scopes of subscribed APIs
// @param accessToken : accesstoken to call the store rest API
// @return scope[], error
func getScopes(appId string, accessToken string) ([]string, error) {
	appDetails, err := getApplicationDetails(appId, accessToken)
	if err != nil || appDetails == nil {
		utils.HandleErrorAndContinue("Error occurred while retrieving subscribed scopes. " +
			"Scopes may not be included in the access token", err)
	}
	if len(appDetails.SubscriptionScopes) > 0 {
		scopesCount := len(appDetails.SubscriptionScopes)
		var scopes = make([]string, scopesCount)
		for i := 0; i < scopesCount; i++ {
			scopes[i] = appDetails.SubscriptionScopes[i].Key
		}
		return scopes, nil
	} else {
		return nil, nil
	}
}

// Generate client credentials for the application first time and generate access token
// @param appId : application ID of the app to be generated keys
// @param token : token to invoke the store rest API
// @return client_id, client_secret, error
func generateApplicationKeys(appId string, token string) (*utils.KeygenResponse, error) {

	applicationEndpoint := utils.GetApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) +
		"/" + appId + "/generate-keys"
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + token
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	generateKeyReq := utils.KeygenRequest{
		KeyType:                 "PRODUCTION",
		GrantTypesToBeSupported: []string{"refresh_token", "password", "client_credentials"},
		ValidityTime:            utils.DefaultTokenValidityPeriod,
	}
	body, err := json.Marshal(generateKeyReq);
	if body == nil && err != nil {
		utils.HandleErrorAndExit("Error occurred while creating CLI application key generation request.", err)
	}

	resp, err := utils.InvokePOSTRequest(applicationEndpoint, headers, string(body))
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		keygenResponse := &utils.KeygenResponse{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &keygenResponse)
		return keygenResponse, err
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return nil, fmt.Errorf("invalid username/password combination")
		}
		return nil, errors.New("Request didn't respond 200 OK for application key generation. Status: " + resp.Status())
	}
}

// Preparing scope values to compatible with request payload
// @param scopes []string : Scopes of the APIs subscribed to an application
// @param password : Password for application server account
// @param url : Registration Endpoint for the environment
// @return string with formatted scope
func prepScopeValues(scope []string) string {
	scopeParam := ""
	for i := 0; i < len(scope); i++ {
		if i == len(scope)-1 {
			scopeParam += "\"" + scope[i] + "\""
		} else {
			scopeParam += "\"" + scope[i] + "\", "
		}
	}
	return scopeParam
}

//init function to add the cli command to the root command
func init() {
	RootCmd.AddCommand(genKeyCmd)
	genKeyCmd.Flags().StringVarP(&keyGenEnv, "environment", "e", "", "Key generation environment")
	genKeyCmd.Flags().StringVarP(&apiName, "name", "n", "", "API to be generated keys")
	genKeyCmd.Flags().StringVarP(&apiVersion, "version", "v", "", "Version of the API")
	genKeyCmd.Flags().StringVarP(&apiProvider, "provider", "r", "", "Provider of the API")
	_ = genKeyCmd.MarkFlagRequired("environment")
}
