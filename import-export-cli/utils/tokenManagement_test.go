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
	"os"
	"path/filepath"
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

func TestExecutePreCommandWithBasicAuth(t *testing.T) {
	var apimStub = getApimStubOK(t)

	// endpoints
	mainConfig := new(MainConfig)
	mainConfigFileName := "test_main_config.yaml"
	mainConfigFilePath := filepath.Join(CurrentDir, mainConfigFileName)
	WriteConfigFile(mainConfig, mainConfigFilePath)

	defer func() {
		os.Remove(mainConfigFilePath)
		apimStub.Close()
	}()
}

// test case 1 - env exists in both endpoints (mainConfig) file and keys file
func TestExecutePreCommandWithOAuth1(t *testing.T) {
	var apimStub = getApimStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var registrationStub = getRegistrationStubOK(t)

	// endpoints
	mainConfig := new(MainConfig)
	mainConfigFileName := "test_main_config.yaml"
	mainConfigFilePath := filepath.Join(CurrentDir, mainConfigFileName)

	mainConfig.Config = Config{2500, "/home/exported"}
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[devName] = EnvEndpoints{apimStub.URL,
		registrationStub.URL, oauthStub.URL}
	WriteConfigFile(mainConfig, mainConfigFilePath)

	// keys
	envKeysAll := new(EnvKeysAll)
	keysAllFileName := "test_keys_all.yaml"
	keysAllFilePath := filepath.Join(CurrentDir, keysAllFileName)
	envKeysAll.Environments = make(map[string]EnvKeys)
	devEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(devPassword)), "dev_client_secret")
	envKeysAll.Environments[devName] = EnvKeys{"dev_client_id", devEncryptedClientSecret, devUsername}
	WriteConfigFile(envKeysAll, keysAllFilePath)

	accessToken, apimEndpoint, err := ExecutePreCommandWithOAuth(devName, devUsername, "admin", mainConfigFilePath, keysAllFilePath)
	if accessToken != sampleAccessToken {
		t.Errorf("Expected accessToken: '%s', got '%s' instead\n", sampleAccessToken, accessToken)
	}

	if apimEndpoint != apimStub.URL {
		t.Errorf("Expected apimEndpoint: '%s', got '%s' instead\n", apimStub.URL, apimEndpoint)
	}

	if err != nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "nil", err)
	}

	defer func() {
		os.Remove(mainConfigFilePath)
		os.Remove(keysAllFilePath)
		apimStub.Close()
		registrationStub.Close()
		oauthStub.Close()
	}()
}

// test case 2 - env exists only in endpoints (mainConfg) file
func TestExecutePreCommandWithOAuth2(t *testing.T) {
	var apimStub = getApimStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var registrationStub = getRegistrationStubOK(t)

	// endpoints
	mainConfig := new(MainConfig)
	mainConfigFileName := "test_main_config.yaml"
	mainConfigFilePath := filepath.Join(CurrentDir, mainConfigFileName)

	mainConfig.Config = Config{2500, "/home/exported"}
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[devName] = EnvEndpoints{apimStub.URL,
		registrationStub.URL, oauthStub.URL}
	WriteConfigFile(mainConfig, mainConfigFilePath)

	// keys
	envKeysAll := new(EnvKeysAll)
	keysAllFileName := "test_keys_all.yaml"
	keysAllFilePath := filepath.Join(CurrentDir, keysAllFileName)
	envKeysAll.Environments = make(map[string]EnvKeys)
	WriteConfigFile(envKeysAll, keysAllFilePath)

	accessToken, apimEndpoint, err := ExecutePreCommandWithOAuth(devName, devUsername, "admin", mainConfigFilePath, keysAllFilePath)
	if accessToken != sampleAccessToken {
		t.Errorf("Expected accessToken: '%s', got '%s' instead\n", sampleAccessToken, accessToken)
	}

	if apimEndpoint != apimStub.URL {
		t.Errorf("Expected apimEndpoint: '%s', got '%s' instead\n", apimStub.URL, apimEndpoint)
	}

	if err != nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "nil", err)
	}

	defer func() {
		os.Remove(mainConfigFilePath)
		os.Remove(keysAllFilePath)
		apimStub.Close()
		registrationStub.Close()
		oauthStub.Close()
	}()
}

// test case 3 - env does not exist in either file
func TestExecutePreCommandWithOAuth3(t *testing.T) {
	var apimStub = getApimStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var registrationStub = getRegistrationStubOK(t)

	// endpoints
	mainConfig := new(MainConfig)
	mainConfigFileName := "test_main_config.yaml"
	mainConfigFilePath := filepath.Join(CurrentDir, mainConfigFileName)

	mainConfig.Config = Config{2500, "/home/exported"}
	mainConfig.Environments = make(map[string]EnvEndpoints)
	WriteConfigFile(mainConfig, mainConfigFilePath)

	// keys
	envKeysAll := new(EnvKeysAll)
	keysAllFileName := "test_keys_all.yaml"
	keysAllFilePath := filepath.Join(CurrentDir, keysAllFileName)
	envKeysAll.Environments = make(map[string]EnvKeys)
	WriteConfigFile(envKeysAll, keysAllFilePath)

	accessToken, apimEndpoint, err := ExecutePreCommandWithOAuth(devName, devUsername, "admin", mainConfigFilePath, keysAllFilePath)
	if accessToken != "" {
		t.Errorf("Expected accessToken: '%s', got '%s' instead\n", "", accessToken)
	}

	if apimEndpoint != "" {
		t.Errorf("Expected apimEndpoint: '%s', got '%s' instead\n", "", apimEndpoint)
	}

	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", err, "nil")
	}

	defer func() {
		os.Remove(mainConfigFilePath)
		os.Remove(keysAllFilePath)
		apimStub.Close()
		registrationStub.Close()
		oauthStub.Close()
	}()
}

// test case 4 - blank env name
func TestExecutePreCommandWithOAuth4(t *testing.T) {
	var apimStub = getApimStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var registrationStub = getRegistrationStubOK(t)

	// endpoints
	mainConfig := new(MainConfig)
	mainConfigFileName := "test_main_config.yaml"
	mainConfigFilePath := filepath.Join(CurrentDir, mainConfigFileName)

	mainConfig.Config = Config{2500, "/home/exported"}
	mainConfig.Environments = make(map[string]EnvEndpoints)
	WriteConfigFile(mainConfig, mainConfigFilePath)

	accessToken, apimEndpoint, err := ExecutePreCommandWithOAuth("", devUsername, "admin", mainConfigFilePath, "")
	if accessToken != "" {
		t.Errorf("Expected accessToken: '%s', got '%s' instead\n", "", accessToken)
	}

	if apimEndpoint != "" {
		t.Errorf("Expected apimEndpoint: '%s', got '%s' instead\n", "", apimEndpoint)
	}

	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", err, "nil")
	}

	defer func() {
		os.Remove(mainConfigFilePath)
		apimStub.Close()
		registrationStub.Close()
		oauthStub.Close()
	}()
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
					 "apim:subscription_block","apim:subscribe"
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
