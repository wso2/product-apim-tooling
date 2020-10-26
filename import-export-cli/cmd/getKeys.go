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
	"net/http"
	"strings"

	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// keys command related Info
const GetKeysCmdLiteral = "keys"
const getKeysCmdShortDesc = "Generate access token to invoke the API or API Product"
const getKeysCmdLongDesc = `Generate JWT token to invoke the API or API Product by subscribing to a default application for testing purposes`
const getKeysCmdExamples = utils.ProjectName + " " + GetCmdLiteral + " " + GetKeysCmdLiteral + ` -n TwitterAPI -v 1.0.0 -e dev --provider admin
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory.
You can override the default token endpoint using --token (-t) optional flag providing a new token endpoint`

var keyGenEnv string
var apiName string
var apiVersion string
var apiProvider string
var tokenType string
var subscriptionThrottlingTier string
var applicationThrottlingPolicy string
var keyGenTokenEndpoint string

var getKeysCmd = &cobra.Command{
	Use:     GetKeysCmdLiteral,
	Short:   getKeysCmdShortDesc,
	Long:    getKeysCmdLongDesc,
	Example: getKeysCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {

		utils.Logln(utils.LogPrefixInfo + GetKeysCmdLiteral + " called")
		getKeys()
	},
}

//Subscribe the given API or API Product to the default application and generate an access token
func getKeys() {

	cred, err := GetCredentials(keyGenEnv)
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
	accessToken, err := credentials.GetOAuthAccessToken(cred, keyGenEnv)
	if err != nil {
		utils.HandleErrorAndExit("Internal error occurred", err)
	}
	utils.Logln(utils.LogPrefixInfo + "Generated a token to access the Publisher and DevPortal REST APIs.")
	//retrieving subscription tiers
	tiers, err := getAvailableAPITiers(accessToken)

	if tiers != nil && err == nil {
		utils.Logln(utils.LogPrefixInfo+"Retrieved available subscription tiers of the API or API Product: ", tiers)
		// Needs an available subscription tier when subscribing to the particular API or API Product using the application
		subscriptionThrottlingTier = tiers[0]
	} else {
		utils.HandleErrorAndExit("Internal error occurred", err)
	}
	// Retrieving application throttling policy
	applicationThrottlingPolicy, err := getApplicationThrottlingPolicy(accessToken)
	// If the application throttling policy call fails, exit with the error
	if err != nil {
		utils.HandleErrorAndExit("Internal error occurred", err)
	}
	utils.Logln(utils.LogPrefixInfo+"Retrieved application throttling policy successfully: ", applicationThrottlingPolicy)
	//search if the default cli application already exists
	appId, err := searchApplication(utils.DefaultCliApp, accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Internal error occurred", err)
	}
	utils.Logln(utils.LogPrefixInfo + "Searched if application exists.")
	//if the application exists
	if appId != "" {
		utils.Logln(utils.LogPrefixInfo + "CLI application already exists")
		// Subscribe API or API Product to a given application
		subId, err := subscribe(appId, accessToken)
		// If subscription fails
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

			//retrieve keys of application to see if there are already generated keys
			appKeys, keysErr := getApplicationKeys(appId, accessToken)
			if keysErr != nil {
				utils.HandleErrorAndExit("Error occurred while getting CLI application keys.", keysErr)
			}

			//if keys have been already generated before, then update the consumer key and secret
			if appKeys.Count != 0 {
				//If the keys have not been generated and the application is updated
				token, err := getNewToken(&appKeys.List[0], scopes)
				//Assert token endpoint related fails and errors
				if err != nil {
					utils.HandleErrorAndExit("Error while generating token. ", err)
				}

				if accessToken != "" {
					// Access Token generated successfully.
					fmt.Println(token)
				} else {
					utils.HandleErrorAndExit("Error while generating token: ", err)
				}
			} else {
				//If the application is already created but the keys have not generated in the first time
				keygenResponse, err := generateApplicationKeys(appId, accessToken)
				if keygenResponse == nil && err != nil {
					utils.HandleErrorAndExit("Error occurred while generating CLI application keys.", err)
				}
				// Access Token generated successfully.
				fmt.Println(keygenResponse.Token.AccessToken)
			}
		} else {
			utils.HandleErrorAndExit("Error while retrieving the CLI application:", err)
		}
	} else {
		//If the default cli appId does not exist in the environment
		//Create the application
		createdAppId, appName, err := createApplication(accessToken, applicationThrottlingPolicy)
		appId = createdAppId
		if createdAppId != "" || appName != "" {
			utils.Logln(utils.LogPrefixInfo+"Created CLI application: ", appName)
		} else {
			//if error occurred while creating the application, then
			utils.HandleErrorAndExit("Error while creating the CLI application:", err)
		}
		//Search the if the given API or API Product is present to subscribe
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
		keygenResponse, err := generateApplicationKeys(appId, accessToken)
		if err != nil {
			utils.HandleErrorAndExit("Error while generating CLI application keys", err)
		}
		appKey := &utils.ApplicationKey{}
		appKey.ConsumerKey = keygenResponse.ConsumerKey
		appKey.ConsumerSecret = keygenResponse.ConsumerSecret
		token, err := getNewToken(appKey, scopes)
		if token != "" {
			// Access Token generated successfully.
			fmt.Println(token)
		} else {
			utils.HandleErrorAndExit("Error while generating token: ", err)
		}
	}
}

// Retrieve an available throttling tiers of the API or API Product
// @param accessToken : Access token to authenticate the store REST API
// @return tiers, error
func getAvailableAPITiers(accessToken string) ([]string, error) {
	apiId, err := searchApiOrProduct(accessToken)
	if apiId == "" && err != nil {
		return nil, err
	}
	api, err := getApiOrProduct(apiId, accessToken)
	if err == nil && api != nil {
		return api.Policies, err
	} else {
		return nil, err
	}
}

// Retrieve an available application throttling policy
// @param accessToken : Access token to authenticate the store REST API
// @return throttlingPolicy, error
func getApplicationThrottlingPolicy(accessToken string) (string, error) {
	applicationThrottlingPoliciesEndpoint := utils.GetDevPortalThrottlingPoliciesEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) + "/application"
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(applicationThrottlingPoliciesEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		applicationThrottlingData := &utils.ThrottlingPoliciesList{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &applicationThrottlingData)
		if applicationThrottlingData.Count != 0 {
			// Needs an available throttling policy when creating/updating an application
			throttlingPolicy := applicationThrottlingData.List[0].Name
			return throttlingPolicy, err
		}
		return "", err
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("authorization failed while trying to retrieve the details of application throttling policies.")
		}
		return "", errors.New("Request didn't respond 200 OK for retrieving application throttling policies. Status: " + resp.Status())
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
	//Request body for the store REST API
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
			return "", "", fmt.Errorf("authorization failed during CLI client registration process")
		}
		return "", "", errors.New("Request didn't respond 200 OK for DCR request. Status: " + resp.Status())
	}
}

// Search if the application exists with the name
// @param appName : Name of the application
// @param accessToken : Access token to authenticate the store REST API
// @return appId, error
func searchApplication(appName string, accessToken string) (string, error) {
	//Application REST API endpoint of the environment from the config file
	applicationEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
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
			return "", fmt.Errorf("authorization failed while searching CLI application: " + appName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching existing applications. " +
			"Status: " + resp.Status())
	}
}

// Searching if the API or API Product is available
// @param accessToken : Access token to call the store REST API
// @return apiId, error
func searchApiOrProduct(accessToken string) (string, error) {
	// Unified Search endpoint from the config file to search APIs or API Products
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)

	//Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	//headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	var queryVal string
	if apiName != "" {
		queryVal = "name:\"" + apiName + "\""
		if apiVersion == "" {
			// If the user has not specified the version, use the version as 1.0.0
			apiVersion = utils.DefaultApiProductVersion
		}
		queryVal = queryVal + " version:\"" + apiVersion + "\""
		if apiProvider != "" {
			queryVal = queryVal + " provider:\"" + apiProvider + "\""
		}
	}
	resp, err := utils.InvokeGETRequestWithQueryParam("query", queryVal, unifiedSearchEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		apiData := &utils.ApiSearch{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &apiData)
		if apiData.Count != 0 {
			apiId := apiData.List[0].ID
			return apiId, err
		}
		if apiProvider != "" {
			return "", errors.New("Requested API is not available in the store. API: " + apiName +
				" Version: " + apiVersion + " Provider: " + apiProvider)
		}
		return "", errors.New("Requested API is not available in the store. API: " + apiName +
			" Version: " + apiVersion)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("authorization failed while searching API or API Product: " + apiName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching APIs and API Products. Status: " + resp.Status())
	}
}

// Subscribe API or API Product to a given application
// @param appId : Application ID to subscribe the API or API Product
// @param accessToken : Token to call REST API
// @return subscriptionId, error
func subscribe(appId string, accessToken string) (string, error) {
	apiId, err := searchApiOrProduct(accessToken)
	if apiId != "" && err == nil {
		//If the API or API Product is present, subscribe that API or API Product to the application
		utils.Logln(utils.LogPrefixInfo+"API or API Product name: ", apiName, "& version: ", apiVersion, "exists")
		subId, err := subscribeApiOrProduct(apiId, appId, accessToken)
		if subId != "" {
			utils.Logln(utils.LogPrefixInfo+"API or API Product", apiName, ":", apiVersion, "subscribed successfully.")
		} else {
			utils.HandleErrorAndExit("Error while subscribing the CLI application to the API: "+appId, err)
		}
		return subId, err
	} else {
		return "", errors.New("API or API Product is not found. Name: " + apiName + " version: " + apiVersion)
	}
}

// Get API or API Product specific details of a given API or API Product
// @param apiId : API ID to retrieve the information
// @param accessToken : Access token to call the REST API
// @return API, error
func getApiOrProduct(apiId string, accessToken string) (*utils.APIData, error) {
	// Since apis/{api-id} supports retrieving details of both APIs and API Products, we can use it here.
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
			return nil, fmt.Errorf("authorization failed while trying to retrieve the details of API or API Prodcut: " + apiId)
		}
		return nil, errors.New("Request didn't respond 200 OK for retrieving API or API Product details. Status: " + resp.Status())
	}
}

// Subscribe the API or API Product to a given Application
// @param apiId : API or API Product ID to be subscribed
// @param appId : Application ID to be subscribed
// @param accessToken : Access token to call the REST API
// @return subscriptionId, error
func subscribeApiOrProduct(apiId string, appId string, accessToken string) (string, error) {
	//todo: subscription endpoint to be included in conf
	subEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
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
			ThrottlingPolicy: subscriptionThrottlingTier,
		}
		//If there is no subscription, make a subscription
		body, err := json.Marshal(subscriptionReq)
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
				return "", fmt.Errorf("authorization failed while trying to subscribe to the API or API Product: " + apiId)
			}
			return "", errors.New("Request didn't respond 200 OK for subscribing to the API or API Product. Status: " + resp.Status())
		}
	} else {
		utils.Logf("Error: %s\n", subResp.Error())
		utils.Logf("Body: %s\n", subResp.Body())
		if subResp.StatusCode() == http.StatusUnauthorized {
			return "", fmt.Errorf("authorization failed while trying to check existing subscriptions of API or API Product: " +
				apiId)
		}
		return "", errors.New("Request didn't respond 200 OK: " + subResp.Status())
	}
}

// Get application details
// @param appId : Application ID
// @param accessToken : Access token to call the store REST API
// @return AppDetails, error
func getApplicationDetails(appId string, accessToken string) (*utils.AppDetails, error) {

	applicationEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) + "/" + appId
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
			return nil, fmt.Errorf("authorization failed while trying to retrieve the details of application: " +
				appId)
		}
		return nil, errors.New("Request didn't respond 200 OK for retrieving application details. " +
			"Status: " + resp.Status())
	}
}

// Get application keys
// @param appId : Application ID
// @param accessToken : Access token to call the store REST API
// @return AppDetails, error
func getApplicationKeys(appId string, accessToken string) (*utils.AppKeyList, error) {

	applicationEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) +
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
			return nil, fmt.Errorf("authorization failed while trying to retrieve the existing keys of application: " +
				appId)
		}
		return nil, errors.New("Request didn't respond 200 OK for retrieving App key information. " +
			"Status: " + resp.Status())
	}
}

// Update application details
// @param appId : Application ID
// @param accessToken : Access token to call the store REST API
// @return AppDetails, error
func updateApplicationDetails(appId string, body string, accessToken string) (*utils.AppDetails, error) {

	applicationEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) + "/" + appId
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
			return nil, fmt.Errorf("authorization failed while trying to update the CLI application: " + appId)
		}
		return nil, errors.New("Request didn't respond 200 OK for updating CLI application. Status: " + resp.Status())
	}
}

// Create application with a default name in a given environment
// @param accessToken : Access token to call the store REST API
// @param throttlingPolicy : Throttling policy to create the application
// @return client_id, client_secret, error
func createApplication(accessToken string, throttlingPolicy string) (string, string, error) {

	applicationEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	conf := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	appUpdateReq := utils.AppCreateRequest{
		Name:             utils.DefaultCliApp,
		ThrottlingPolicy: throttlingPolicy,
		Description:      "Default application for apictl testing purposes",
		TokenType:        conf.Config.TokenType,
	}
	body, err := json.Marshal(appUpdateReq)
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
			return "", "", fmt.Errorf("authorization failed while trying to create the CLI application")
		}
		return "", "", errors.New("Request didn't respond 200 OK for application creation. Status: " + resp.Status())
	}
}

// Calling token endpoint to get access token for the already created application
// @param key : Details of the particular key
// @param scopes[] : Scopes to generate the token
// @return accessToken, error
func getNewToken(key *utils.ApplicationKey, scopes []string) (string, error) {
	var tokenEndpoint string
	if keyGenTokenEndpoint == "" {
		tokenEndpoint = utils.GetTokenEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath)
	} else {
		tokenEndpoint = keyGenTokenEndpoint
	}
	body := "grant_type=client_credentials&scope=" + strings.Join(scopes, " ")

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " +
		utils.GetBase64EncodedCredentials(key.ConsumerKey, key.ConsumerSecret)

	headers[utils.HeaderContentType] = utils.HeaderValueXWWWFormUrlEncoded
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON

	resp, err := utils.InvokePOSTRequest(tokenEndpoint, headers, body)

	if err != nil {
		return "", errors.New("Token Endpoint is not valid. " + err.Error())
	}

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
			return "", fmt.Errorf("authorization failed while generating a token for the CLI application")
		}
		return "", errors.New("Request didn't respond 200 OK for generating a new token. Status: " + resp.Status())
	}

}

// Get all the scopes of the APIs and API Products subscribed to a particular application
// @param appId : Application ID to get the scopes of subscribed APIs and API Products
// @param accessToken : Access token to call the store REST API
// @return scope[], error
func getScopes(appId string, accessToken string) ([]string, error) {
	appDetails, err := getApplicationDetails(appId, accessToken)
	if err != nil || appDetails == nil {
		utils.HandleErrorAndContinue("Error occurred while retrieving subscribed scopes. "+
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
// @param appId : Application ID of the app to be generated keys
// @param token : Token to invoke the store REST API
// @return client_id, client_secret, error
func generateApplicationKeys(appId string, token string) (*utils.KeygenResponse, error) {

	applicationEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(keyGenEnv, utils.MainConfigFilePath) +
		"/" + appId + "/generate-keys"
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + token
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	generateKeyReq := utils.KeygenRequest{
		KeyType:                 "PRODUCTION",
		GrantTypesToBeSupported: []string{"refresh_token", "password", "client_credentials"},
		ValidityTime:            utils.DefaultTokenValidityPeriod,
	}
	body, err := json.Marshal(generateKeyReq)
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
			return nil, fmt.Errorf("authorization failed while generating keys of the CLI application: " + appId)
		}
		return nil, errors.New("Request didn't respond 200 OK for application key generation. Status: " + resp.Status())
	}
}

// Preparing scope values to compatible with request payload
// @param scopes []string : Scopes of the APIs and API Products subscribed to an application
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
	GetCmd.AddCommand(getKeysCmd)
	getKeysCmd.Flags().StringVarP(&keyGenEnv, "environment", "e", "", "Key generation environment")
	getKeysCmd.Flags().StringVarP(&apiName, "name", "n", "", "API or API Product to generate keys")
	getKeysCmd.Flags().StringVarP(&apiVersion, "version", "v", "", "Version of the API")
	getKeysCmd.Flags().StringVarP(&apiProvider, "provider", "r", "", "Provider of the API or API Product")
	getKeysCmd.Flags().StringVarP(&keyGenTokenEndpoint, "token", "t", "", "Token endpoint URL of Environment")
	_ = getKeysCmd.MarkFlagRequired("name")
	_ = getKeysCmd.MarkFlagRequired("environment")
}
