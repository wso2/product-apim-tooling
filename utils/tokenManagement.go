package utils

import (
	"github.com/go-resty/resty"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

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

	return m // m contains 'access_token', 'refresh_token' etc
}
