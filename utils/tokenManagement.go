package utils

import (
	"github.com/go-resty/resty"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"errors"
)

func ExecutePreCommand(environment string) (string, string, error){
	if EnvExistsInEndpointsFile(environment) {
		registrationEndpoint := GetRegistrationEndpointOfEnv(environment)
		apiManagerEndpoint := GetAPIMEndpointOfEnv(environment)
		tokenEndpoint := GetTokenEndpointOfEnv(environment)
		var username string
		var password string
		var clientID string
		var clientSecret string

		if EnvExistsInKeysFile(environment) {
			// client_id, client_secret exists in file
			username = GetUsernameOfEnv(environment)
			fmt.Println("Username:", username)
			password = PromptForPassword()
			clientID = GetClientIDOfEnv(environment)
			clientSecret = GetClientSecretOfEnv(environment, password)

			fmt.Println("ClientID:", clientID)
			fmt.Println("ClientSecret:", clientSecret)
		} else {
			// env exists in endpoints file, but not in keys file
			// no client_id, client_secret in file
			// Get new values
			username = strings.TrimSpace(PromptForUsername())
			password = PromptForPassword()

			fmt.Println("\nUsername: " + username + "\n")
			clientID, clientSecret = GetClientIDSecret(username, password, registrationEndpoint)

			// Persist clientID, clientSecret, Username in file
			encryptedClientSecret := Encrypt([]byte(GetMD5Hash(password)), clientSecret)
			envKeys := EnvKeys{clientID, encryptedClientSecret, username}
			AddNewEnvToKeysFile(environment, envKeys)
		}

		// Get OAuth Tokens
		m := GetOAuthTokens(username, password, GetBase64EncodedCredentials(clientID, clientSecret), tokenEndpoint)
		accessToken := m["access_token"]
		fmt.Println("AccessToken:", accessToken)

		return accessToken, apiManagerEndpoint, nil
	}else{
		return "", "", errors.New("Details incorrect/unavailable for environment "+ environment)
	}
}


// GetClientIDSecret implemented using go-resty
// provide username, password
// returns client_id, client_secret
func GetClientIDSecret(username string, password string, url string) (string, string) {
	body := `{"clientName": "Test", "redirect_uris": "www.google.lk", "grant_types":"password"}`
	headers := make(map[string]string)

	headers[HeaderContentType] = HeaderValueApplicationJSON
	// headers["Content-Type"] = "application/json"

	headers[HeaderAuthorization] = HeaderValueAuthBasicPrefix + " " + GetBase64EncodedCredentials(username, password)
	// headers["Authorization"] = "Basic " + GetBase64EncodedCredentials(username, password)

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTPS certificates

	// POST request using resty
	resp, err := InvokePOSTRequest(url, headers, body)

	if err != nil {
		fmt.Println("Error in connecting:")
		panic(err)
	}

	m := make(map[string]string) // a map to hold response data
	data := []byte(resp.Body())
	_ = json.Unmarshal(data, &m) // add response data to m

	clientID := m["client_id"]
	clientSecret := m["client_secret"]

	return clientID, clientSecret
}

// Encode the concatenation of two strings (using ":")
// provide two strings
// returns base64Encode(strOne:strTwo)
func GetBase64EncodedCredentials(key string, secret string) string {
	line := key + ":" + secret
	encoded := base64.StdEncoding.EncodeToString([]byte(line))
	return encoded
}

// GetOAuthTokens implemented using go-resty/resty
// provide username, password, and validity period for the access token
// returns the response as a map
func GetOAuthTokens(username string, password string, b64EncodedClientIDClientSecret string, url string) map[string]string {
	validityPeriod := DefaultTokenValidityPeriod
	body := "grant_type=password&username=" + username + "&password=" + password + "&validity_period=" + validityPeriod

	// set headers
	headers := make(map[string]string)
	headers[HeaderContentType] = HeaderValueXWWWFormUrlEncoded
	headers[HeaderAuthorization] = HeaderValueAuthBearerPrefix + " " + b64EncodedClientIDClientSecret
	headers[HeaderAccept] = HeaderValueApplicationJSON

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTP certificates
	resp, err := InvokePOSTRequest(url, headers, body)

	if err != nil {
		HandleErrorAndExit("Unable to Connect", nil)
	}

	m := make(map[string]string) // a map to hold response data
	data := []byte(resp.Body())
	_ = json.Unmarshal(data, &m) // add response data to m

	return m // m contains 'access_token', 'refresh_token' etc
}

// GetAccessTokenUsingRefreshToken implemented using resty
// provide refreshToken (decrypted), and base64(clientID:clientSecret)
// returns the response as a map
func GetAccessTokenUsingRefreshToken(refreshToken string, b64encodedKeySecret string) map[string]string {
	url := "https://localhost:9443/oauth2/token"
	body := "grant_type=refresh_token&refresh_token=" + refreshToken + "&validity_period=3600&scopes="

	// set headers
	headers := make(map[string]string)
	headers[HeaderAuthorization] = HeaderValueAuthBearerPrefix + " " + b64encodedKeySecret
	headers[HeaderContentType] = HeaderValueXWWWFormUrlEncoded
	headers[HeaderAccept] = HeaderValueApplicationJSON

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTP certificates

	resp, err := resty.R().
		SetHeaders(headers).
		SetBody(body).
		Post(url)

	if err != nil {
		HandleErrorAndExit("Unable to Connect", err)
	}

	m := make(map[string]string) // a map to hold response data
	data := []byte(resp.Body())
	_ = json.Unmarshal(data, &m) // add response data to m

	if resp.StatusCode() == 200 {
		fmt.Println("OAuth Tokens Received")
	}else if resp.StatusCode() == 400{
		fmt.Println("Error:", m["error_description"])
		os.Exit(1)
	}
	return m
}
