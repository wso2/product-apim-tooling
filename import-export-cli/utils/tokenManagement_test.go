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
	"github.com/renstrom/dedent"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

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
			t.Errorf("Expected '%s', got '%s'\n", http.MethodPost, r.Method)
		}
		if r.Header.Get(HeaderContentType) != HeaderValueApplicationJSON {
			t.Errorf("Expected '%s', got '%s'\n", HeaderValueApplicationJSON, r.Header.Get(HeaderContentType))
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
	SkipTLSVerification = true
	var oauthStub = getOAuthStubOK(t)
	defer oauthStub.Close()

	m, err := GetOAuthTokens("admin", "admin", "", oauthStub.URL)
	if err != nil {
		t.Error("Error in GetOAuthTokens()")
	}

	if m["refresh_token"] != "fe8f8400-05c9-430f-8e2f-4f3b2fbd01f8" {
		t.Error("Error in GetOAuthTokens(): Incorrect RefreshToken")
	}
	if m["access_token"] != "a2e5c3ac-68e6-4d78-a8a1-b2b0372cb575" {
		t.Error("Error in GetOAuthTokens(): Incorrect AccessToken")
	}
}

// test case 1 - MainConfig file available, flagUsername not blank, flagPassword not blank
func TestExecutePreCommand1(t *testing.T) {
	var registrationStub = getRegistrationStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var apimStub = getApimStubOK(t)

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

	ExecutePreCommand("dev", "dev_username", "dev_password", mainConfigFilePath, keysAllFilePath)

	defer func() {
		os.Remove(mainConfigFilePath)
		os.Remove(keysAllFilePath)
		apimStub.Close()
		oauthStub.Close()
		registrationStub.Close()
	}()
}

// test case 5 - MainConfig file available, flagUsername not blank, flagPassword blank
func TestExecutePreCommand5(t *testing.T) {
	var registrationStub = getRegistrationStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var apimStub = getApimStubOK(t)

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

	ExecutePreCommand("dev", "dev_username", "", mainConfigFilePath, keysAllFilePath)

	defer func() {
		os.Remove(mainConfigFilePath)
		os.Remove(keysAllFilePath)
		apimStub.Close()
		oauthStub.Close()
		registrationStub.Close()
	}()
}

// test case 6 - MainConfig file available, flagUsername blank, flagPassword not blank
func TestExecutePreCommand6(t *testing.T) {
	var registrationStub = getRegistrationStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var apimStub = getApimStubOK(t)

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

	ExecutePreCommand("dev", "", "dev-password", mainConfigFilePath, keysAllFilePath)

	defer func() {
		os.Remove(mainConfigFilePath)
		os.Remove(keysAllFilePath)
		apimStub.Close()
		oauthStub.Close()
		registrationStub.Close()
	}()
}

func TestExecutePreCommand2(t *testing.T) {
	var registrationStub = getRegistrationStubOK(t)
	var oauthStub = getOAuthStubOK(t)
	var apimStub = getApimStubOK(t)

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

	ExecutePreCommand("dev", "dev_username", "dev_password", mainConfigFilePath, keysAllFilePath)

	defer func() {
		os.Remove(mainConfigFilePath)
		os.Remove(keysAllFilePath)
		apimStub.Close()
		oauthStub.Close()
		registrationStub.Close()
	}()
}

// test case 3 - input environment not available in MainConfig (endpoints) file
func TestExecutePreCommand3(t *testing.T) {
	// endpoints
	mainConfig := new(MainConfig)
	mainConfigFileName := "test_main_config.yaml"
	mainConfigFilePath := filepath.Join(CurrentDir, mainConfigFileName)

	mainConfig.Config = Config{2500, "/home/exported"}
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[devName] = EnvEndpoints{
		"publisher-endpoint",
		"reg-endpoint",
		"token-endpoint",
	}
	WriteConfigFile(mainConfig, mainConfigFilePath)

	_, _, err := ExecutePreCommand("not-available", "", "", mainConfigFilePath, "")
	if err == nil {
		t.Errorf("Expected error, go nil instead\n")
	}

	defer os.Remove(mainConfigFilePath)
}

// test case 4 - input environment blank
func TestExecutePreCommand4(t *testing.T) {
	// endpoints
	mainConfig := new(MainConfig)
	mainConfigFileName := "test_main_config.yaml"
	mainConfigFilePath := filepath.Join(CurrentDir, mainConfigFileName)

	mainConfig.Config = Config{2500, "/home/exported"}
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[devName] = EnvEndpoints{
		"publisher-endpoint",
		"reg-endpoint",
		"token-endpoint",
	}
	WriteConfigFile(mainConfig, mainConfigFilePath)

	_, _, err := ExecutePreCommand("", "", "", mainConfigFilePath, "")
	if err == nil {
		t.Errorf("Expected error, go nil instead\n")
	}

	defer os.Remove(mainConfigFilePath)

}

// getOAuthStubOK - Helper for testing
// Token endpoint - OK
func getOAuthStubOK(t *testing.T) *httptest.Server {
	var oauthStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s'\n", http.MethodPost, r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueXWWWFormUrlEncoded {
			t.Errorf("Expected '%s', got '%s'\n", HeaderValueXWWWFormUrlEncoded, r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)

		body := dedent.Dedent(`
			{"access_token":"a2e5c3ac-68e6-4d78-a8a1-b2b0372cb575",
			 "refresh_token":"fe8f8400-05c9-430f-8e2f-4f3b2fbd01f8",
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

// getOAuthStubError - Helper for testing
// Token endpoint - Service Unavailable
func getOAuthStubError(t *testing.T) *httptest.Server {
	var oauthStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s'\n", http.MethodPost, r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueXWWWFormUrlEncoded {
			t.Errorf("Expected '%s', got '%s'\n", HeaderValueXWWWFormUrlEncoded, r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusServiceUnavailable)

		body := dedent.Dedent(``)

		w.Write([]byte(body))
	}))

	return oauthStub
}

// getRegistrationStubOK - Helper for testing
// Client registration - OK
func getRegistrationStubOK(t *testing.T) *httptest.Server {
	var registrationStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s'\n", http.MethodPost, r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueApplicationJSON {
			t.Errorf("Expected '%s', got '%s'\n", HeaderValueApplicationJSON, r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)
		body := dedent.Dedent(`{"client_name":"Test1",
									"client_id":"be88563b-21cb-417c-b574-bf1079959679",
									"client_secret":"ecb105a0-117c-463d-9376-442d24864f26"}`)

		w.Write([]byte(body))
	}))

	return registrationStub
}

// getRegistrationStubError - Helper for testing
// Client registration - Service Unavailable
func getRegistrationStubError(t *testing.T) *httptest.Server {
	var registrationStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s'\n", http.MethodPost, r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueApplicationJSON {
			t.Errorf("Expected '%s', got '%s'\n", HeaderValueApplicationJSON, r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		body := dedent.Dedent(``)

		w.Write([]byte(body))
	}))

	return registrationStub
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
		w.WriteHeader(http.StatusOK)
	}))

	return apimStub
}
