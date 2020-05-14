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

package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

// APIMClient : Enables interacting with an instance of APIM
type APIMClient struct {
	accessToke string
}

// APIMConfigs : Configurations of an instance of APIM
type APIMConfigs struct {
	DcrURL   string
	TokenURL string
	Username string
	Password string
}

type dcrRequest struct {
	CallbackURL         string `json:"callbackUrl"`
	ClientName          string `json:"clientName"`
	Owner               string `json:"owner"`
	SupportedGrantTypes string `json:"grantType"`
	IsSaaSApp           bool   `json:"saasApp"`
}

type jsonString struct {
	UserName     string `json:"username"`
	ClientName   string `json:"client_name"`
	RedirectURIs string `json:"redirect_uris"`
	GrantTypes   string `json:"grant_types"`
}

type dcrResponse struct {
	CallbackURL  string      `json:"callBackURL"`
	ClientName   string      `json:"clientName"`
	JSONString   *jsonString `json:"jsonString"`
	ClientID     string      `json:"clientId"`
	ClientSecret string      `json:"clientSecret"`
	IsSaaSApp    bool        `json:"isSaasApplication"`
	Owner        string      `json:"appOwner"`
}

type tokenResponse struct {
	Scope        string `json:"scope"`
	TokeType     string `json:"token_type"`
	ValidTime    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

// Init : Initialise APIM Client instance
func (instance *APIMClient) Init(apimConfigs *APIMConfigs) {
	instance.accessToke = getToken(apimConfigs)
}

func getToken(apimConfigs *APIMConfigs) string {
	registrationResponse := registerClient(apimConfigs.DcrURL, apimConfigs.Username, apimConfigs.Password)

	request := CreatePostEmptyBody(apimConfigs.TokenURL)
	request.SetBasicAuth(registrationResponse.ClientID, registrationResponse.ClientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	values := url.Values{}
	values.Add("grant_type", "password")
	values.Add("username", apimConfigs.Username)
	values.Add("password", apimConfigs.Password)
	values.Add("scope", "apim:api_create apim:api_publish apim:subscribe")

	request.URL.RawQuery = values.Encode()

	response := SendHTTPRequest(request)

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Error when invoking token API")
		os.Exit(1)
	}

	var jsonResp tokenResponse
	json.NewDecoder(response.Body).Decode(&jsonResp)
	return jsonResp.AccessToken
}

func registerClient(dcrURL string, username string, password string) dcrResponse {
	dcrPayload := dcrRequest{}

	dcrPayload.CallbackURL = "http://localhost"
	dcrPayload.ClientName = "apictl-integration-tests"
	dcrPayload.IsSaaSApp = true
	dcrPayload.Owner = "admin"
	dcrPayload.SupportedGrantTypes = "password refresh_token"

	data, err := json.Marshal(dcrPayload)

	if err != nil {
		panic(err)
	}

	request := CreatePost(dcrURL, bytes.NewBuffer(data))

	request.SetBasicAuth(username, password)
	request.Header.Set("Content-Type", "application/json")

	response := SendHTTPRequest(request)

	defer response.Body.Close()

	contentType := response.Header["Content-Type"][0]

	// If DCR endpoint ebing invoked is invalid, an HTML error page will be returned.
	// We cannot rely on checking the response code since it will always be 200.
	// Therefore need to validate the Content Type of the response to detect this condition.
	if contentType != "application/json" {
		fmt.Println("\nInvalid response received for DCR request. Please check if configured DCR endpoint is correct.")
		os.Exit(1)
	}

	var jsonResp dcrResponse
	json.NewDecoder(response.Body).Decode(&jsonResp)

	return jsonResp
}
