package utils

import (
	"github.com/go-resty/resty"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
)

// GetClientIDSecret implemented using go-resty
func GetClientIDSecret(username string, password string) (string, string) {
	url := "https://localhost:9443/identity/connect/register"
	body := `{"clientName": "Test", "redirect_uris": "www.google.lk", "grant_types":"password"}`
	headers := make(map[string]string)
	headers[HeaderContentType] = HeaderValueApplicationJSON
	headers[HeaderAuthorization] = "Basic " + GetBase64EncodedCredentials(username, password)

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTPS certificates

	resp, err := resty.R().
		SetHeaders(headers).
		SetBody(body).
		Post(url)

	if err != nil {
		panic(err)
	}
	// fmt.Printf("\nResponse Body: %v\n", resp)

	m := make(map[string]string)
	data := []byte(resp.Body())
	_ = json.Unmarshal(data, &m)

	clientID := m["client_id"]
	clientSecret := m["client_secret"]

	return clientID, clientSecret
}

func GetBase64EncodedCredentials(key string, secret string) string {
	line := key + ":" + secret
	encoded := base64.StdEncoding.EncodeToString([]byte(line))
	return encoded
}

// GetOAuthTokens implemented using go-resty/resty
func GetOAuthTokens(username string, password string) map[string]string {
	url := "https://localhost:9443/oauth2/token"
	body := "grant_type=password&username=" + username + "&password="+ password +"&validity_period=3600"

	headers := make(map[string]string)
	headers[HeaderContentType] = HeaderValueXWWWFormUrlEncoded
	headers[HeaderAuthorization] = "Bearer " + GetBase64EncodedCredentials(GetClientIDSecret(username, password))
	headers[HeaderAccept] = HeaderValueApplicationJSON
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTP certificates
	resp, err := resty.R().
		SetHeaders(headers).
		SetBody(body).
		Post(url)

	if err != nil {
		HandleErrorAndExit("Unable to Connect", err)
	}

	m := make(map[string]string)
	data := []byte(resp.Body())
	_ = json.Unmarshal(data, &m)

	return m // m contains 'access_token', 'refresh_token' etc
}

func isAccessTokenExpired() {
}

func refreshToken() {

}
