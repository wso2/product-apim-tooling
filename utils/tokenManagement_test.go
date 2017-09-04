package utils

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"net/http"
	"github.com/renstrom/dedent"
)

func TestGetClientIDSecret(t *testing.T) {
	one, two, three := GetClientIDSecret("admin", "admin", "")
	fmt.Println("one:", one)
	fmt.Println("two:", two)
	fmt.Println("three:", three)
}

func TestGetOAuthTokens(t *testing.T) {

}

func TestGetBase64EncodedCredentials(t *testing.T) {
	usernames := []string{"admin", "user", "admin"}
	passwords := []string{"admin", "password", "123456"}
	encodedPairs := []string{"YWRtaW46YWRtaW4=", "dXNlcjpwYXNzd29yZA==", "YWRtaW46MTIzNDU2"}

	for i, s := range encodedPairs {
		if s != GetBase64EncodedCredentials(usernames[i], passwords[i])	{
			t.Errorf("Error in Base64 Encoding. Base64(" + usernames[i] + ":" + passwords[i] + ") = " + encodedPairs[i])
		}
	}
}

func TestGetClientIDSecretUnreachable(t *testing.T) {
	username := "admin"
	password := "admin"
	url := "http://localhost:-41"
	_, _, err := GetClientIDSecret(username, password, url)
	if err == nil {
		t.Errorf("GetClientIDSecret() didn't return an error")
	}
}

func TestGetClientIDSecretOK(t *testing.T) {
	var registrationStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST', got '%s'\n", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(HeaderContentType, HeaderValueApplicationJSON)
		w.Header().Set("Sample-Header", "Sample-Value")
		body := dedent.Dedent(`{"clientName":"Test1",
									"client_id":"be88563b-21cb-417c-b574-bf1079959679",
									"client_secret":"ecb105a0-117c-463d-9376-442d24864f26"}`)


		w.Write([]byte(body))
		//io.WriteString(w, body)
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
	url := "http://localhost:-41"
	_, err := GetOAuthTokens("","", "", url)
	if err == nil {
		t.Errorf("GetOAuthTokens() didn't return an error")
	}
}

func TestGetOAuthTokensOK(t *testing.T) {

}

