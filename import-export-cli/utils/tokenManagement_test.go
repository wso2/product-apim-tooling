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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renstrom/dedent"
)

const sampleAccessToken = "a2e5c3ac-68e6-4d78-a8a1-b2b0372cb575"
const sampleRefreshToken = "fe8f8400-05c9-430f-8e2f-4f3b2fbd01f8"

func TestGetBase64EncodedCredentials(t *testing.T) {
	usernames := []string{"admin", "user", "admin"}
	passwords := []string{"admin", "password", "123456"}
	encodedPairs := []string{"YWRtaW46YWRtaW4=", "dXNlcjpwYXNzd29yZA==", "YWRtaW46MTIzNDU2"}

	for i, s := range encodedPairs {
		if s != GetBase64EncodedCredentials(usernames[i], passwords[i]) {
			t.Errorf("Error in Base64 Encoding. Base64(" + usernames[i] + ":" + passwords[i] + ") = " + encodedPairs[i])
		}
	}
}

func TestGetClientIDSecretUnreachable(t *testing.T) {
	var registrationStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST', got '%s' instead\n", r.Method)
		}
		if r.Header.Get(HeaderContentType) != HeaderValueApplicationJSON {
			t.Errorf("Expected '%s', got '%s' instead\n", HeaderValueApplicationJSON, r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)
	}))
	defer registrationStub.Close()

	username := "admin"
	password := "admin"
	_, _, err := GetClientIDSecret(username, password, registrationStub.URL)
	if err == nil {
		t.Errorf("GetClientIDSecret() didn't return an error")
	}
}

func TestGetClientIDSecretOK(t *testing.T) {
	var registrationStub = getRegistrationStubOK(t)
	defer registrationStub.Close()

	clientID, clientSecret, err := GetClientIDSecret("admin", "admin", registrationStub.URL)
	if err != nil {
		t.Error("Error receving response")
	}

	if clientID != "be88563b-21cb-417c-b574-bf1079959679" {
		t.Error("Invalid ClientID")
	}

	if clientSecret != "ecb105a0-117c-463d-9376-442d24864f26" {
		t.Error("Invalid ClientSecret")
	}
}

func TestGetOAuthTokensOK(t *testing.T) {
	var oauthStub = getOAuthStubOK(t)
	defer oauthStub.Close()

	m, err := GetOAuthTokens("admin", "admin", "", oauthStub.URL)
	if err != nil {
		t.Error("Error in GetOAuthTokens()")
	}

	if m["refresh_token"] != sampleRefreshToken {
		t.Error("Error in GetOAuthTokens(): Incorrect RefreshToken")
	}
	if m["access_token"] != sampleAccessToken {
		t.Error("Error in GetOAuthTokens(): Incorrect AccessToken")
	}
}

// Registration Server - OK
func getRegistrationStubOK(t *testing.T) *httptest.Server {
	var registrationStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s' instead\n", http.MethodPost, r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueApplicationJSON {
			t.Errorf("Expected '%s', got '%s' instead\n", HeaderValueApplicationJSON, r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)
		body := dedent.Dedent(`{"client_name":"Test1",
									"clientId":"be88563b-21cb-417c-b574-bf1079959679",
									"clientSecret":"ecb105a0-117c-463d-9376-442d24864f26"}`)

		w.Write([]byte(body))
	}))
	return registrationStub
}

// OAuth Server - OK
func getOAuthStubOK(t *testing.T) *httptest.Server {
	var oauthStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s' instead\n", http.MethodPost, r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueXWWWFormUrlEncoded {
			t.Errorf("Expected '%s', got '%s' instead\n", HeaderValueXWWWFormUrlEncoded, r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)

		body := dedent.Dedent(`
			{"access_token": "` + sampleAccessToken + `",
			 "refresh_token":"` + sampleRefreshToken + `",
			 "expires_in":1487166427829,
			 "scopes":
					["apim:api_view","apim:api_create","apim:api_publish",
					 "apim:tier_view","apim:tier_manage","apim:subscription_view",
					 "apim:subscription_block","apim:subscribe","apim:policy_import_export"
					]
			}
		`)

		w.Write([]byte(body))
	}))
	return oauthStub
}

// API Manager - OK
func getApimStubOK(t *testing.T) *httptest.Server {
	var apimStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	return apimStub
}

// API Manager - Service Unavailable
func getApimStubError(t *testing.T) *httptest.Server {
	var apimStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	return apimStub
}
