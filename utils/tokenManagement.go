package utils

import (
	"github.com/go-resty/resty"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// GetClientIDSecret implemented using go-resty
func GetClientIDSecret(username string, password string) (string, string) {
	url := "https://localhost:9443/identity/connect/register"
	body := `{"clientName": "Test", "redirect_uris": "www.google.lk", "grant_types":"password"}`
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Basic " + GetBase64EncodedCredentials(username, password)

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

	client_id := m["client_id"]
	client_secret := m["client_secret"]

	return client_id, client_secret
}

func GetBase64EncodedCredentials(key string, secret string) string {
	line := key + ":" + secret
	encoded := base64.StdEncoding.EncodeToString([]byte(line))
	return encoded
}

// GetOAuthAccessToken implemented using go-resty/resty
func GetOAuthAccessToken() {
	url := "https://localhost:9443/oauth2/token"
	body := `{"grant_type": "password",
		"username": "admin",
		"password": "admin",
		"validity_period": "3600"}`

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = "Bearer " + GetBase64EncodedCredentials(GetClientIDSecret())
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := resty.R().
		SetHeaders(headers).
		SetBody(body).
		Post(url)

	if err != nil {
		panic(err)
	}
	fmt.Printf("\nResponse Body: %v\n", resp)

	fmt.Println(body)
}

func isExpired() {

}

func refreshToken() {

}
