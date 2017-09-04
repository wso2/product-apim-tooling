package utils

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"net/http"
	"github.com/renstrom/dedent"
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
			t.Errorf("Expected 'POST', got '%s'\n", r.Method)
		}
		if r.Header.Get(HeaderContentType) != HeaderValueApplicationJSON {
			t.Errorf("Expected '"+HeaderValueApplicationJSON+"', got '%s'\n", r.Header.Get(HeaderContentType))
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
	var registrationStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST', got '%s'\n", r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueApplicationJSON {
			t.Errorf("Expected '"+HeaderValueApplicationJSON+"', got '%s'\n", r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)
		body := dedent.Dedent(`{"client_name":"Test1",
									"client_id":"be88563b-21cb-417c-b574-bf1079959679",
									"client_secret":"ecb105a0-117c-463d-9376-442d24864f26"}`)

		w.Write([]byte(body))
	}))
	defer registrationStub.Close()

	fmt.Println("URL:", registrationStub.URL)

	clientID, clientSecret, err := GetClientIDSecret("admin", "admin", registrationStub.URL)
	fmt.Println("ClientID:", clientID)
	fmt.Println("ClientSecret:", clientSecret)
	fmt.Println("Error:", err)
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

func TestGetOAuthTokensUnreachable(t *testing.T) {
	var oauthStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST', got '%s'\n", r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueXWWWFormUrlEncoded {
			t.Errorf("Exptected '"+HeaderValueXWWWFormUrlEncoded+"', got '%s'\n", r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)
	}))
	defer oauthStub.Close()

	_, err := GetOAuthTokens("", "", "", oauthStub.URL)
	if err == nil {
		t.Errorf("GetOAuthTokens() didn't return an error")
	}
}

func TestGetOAuthTokensOK(t *testing.T) {
	var oauthStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST', got '%s'\n", r.Method)
		}

		if r.Header.Get(HeaderContentType) != HeaderValueXWWWFormUrlEncoded {
			t.Errorf("Expected '"+HeaderValueXWWWFormUrlEncoded+"', got '%s'\n", r.Header.Get(HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)

		body := dedent.Dedent(`
			{"access_token":"a2e5c3ac-68e6-4d78-a8a1-b2b0372cb575",
			 "refresh_token":"fe8f8400-05c9-430f-8e2f-4f3b2fbd01f8",
			 "expires_in":1487166427829,
			 "scopes":["apim:api_view","apim:api_create","apim:api_publish",
				"apim:tier_view","apim:tier_manage","apim:subscription_view",
				"apim:subscription_block","apim:subscribe"
				]
			}
		`)

		w.Write([]byte(body))
	}))
	defer oauthStub.Close()

	m, err := GetOAuthTokens("admin", "admin", "", oauthStub.URL)
	if err != nil {
		t.Error("Error in GetOAuthTokens()")
	}

	if m["refresh_token"] != "fe8f8400-05c9-430f-8e2f-4f3b2fbd01f8"{
		t.Error("Error in GetOAuthTokens(): Incorrect RefreshToken")
	}
	if m["access_token"] != "a2e5c3ac-68e6-4d78-a8a1-b2b0372cb575"{
		t.Error("Error in GetOAuthTokens(): Incorrect AccessToken")
	}
}
