package utils

import (
	"net/url"
	"fmt"
	"os"
	"net/http"
	"bytes"
	"time"
	"strings"
	"errors"
	"io/ioutil"
	"encoding/json"
)

// Token management according to wso2/wum-client
type TokenResponse struct {
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int `json:expires_in`
	RefreshToken string `json:"refres_token"`
	AccessToken  string `json:"access_token"`
}

type TokenErrResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type APIErrResp struct {
	Fault       string `xml:"fault"`
	Code        string `xml:"code"`
	Message     string `xml:"message"`
	Description string `xml:"description"`
}

// Get an access token from WSO2 Update with the given username and the password using the
// 'password' grant type of Oauth2.
// This method returns an error only if the username or password is incorrect.
func GetAccessToken(username string, password []byte, envInfo EnvInfo, envKeys EnvKeys, scope string) (TokenResponse, error) {
	payload := url.Values{}
	payload.Add("grant_type", "password")
	payload.Add("username", username)
	payload.Add("password", string(password))

	if len(scope) > 0 {
		payload.Add("scope", scope)
	}
	// Get an access token and a refresh token
	fmt.Fprintln(os.Stderr, "Authenticating...")
	return InvokeTokenAPI(payload, envInfo, envKeys)
}

func ProcessResponseFromServer(resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logln(ErrorReadingResponseMsg)
		HandleUnableToConnectErrorAndExit(err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		Logln(ErrorProcessingResponseMsg)
		Logln(err.Error())
		HandleUnableToConnectErrorAndExit(nil)
	}
}

func ProcessServerResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logln(ErrorReadingResponseMsg)
		HandleUnableToConnectErrorAndExit(err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return errors.New(ErrorProcessingResponseMsg)
	}

	return nil
}

// Invokes the configured token API of the API gateway. This method can be used to get access tokens
// as well as renew access tokens using the refresh token.
func InvokeTokenAPI(payload url.Values, envInfo EnvInfo, envKeys EnvKeys) (TokenResponse, error) {
	Logln("Invoking WSO2 Token API")
	req, err := http.NewRequest("POST", envInfo.TokenEndpoint, bytes.NewBufferString(payload.Encode()))
	if err != nil {
		HandleUnableToConnectErrorAndExit(err)
	}

	appKey := GetBase64EncodedCredentials(envKeys.ClientID, envKeys.ClientSecret)
	req.Header.Add(HeaderAuthorization, "Basic " + appKey)
	req.Header.Add(HeaderContentType, HeaderValueXWWWFormUrlEncoded)
	resp := SendRequest(req, time.Duration(WSO2APIMUpdateTokenTimeout*time.Minute))

	Logf("Response status code is %d.\n", resp.StatusCode)
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		ter := TokenErrResp{}
		ProcessResponseFromServer(resp, &ter)

		if ter.Error == "invalid_grant" && strings.HasPrefix(ter.ErrorDescription,
			"Authentication failed for") {
			return TokenResponse{}, errors.New("Authentication failed")

		} else if ter.Error == "invalid_grant" && strings.HasPrefix(ter.ErrorDescription,
			"Provided Authorization Grant is invalid") {
			PrintErrorMessageAndExit(InvalidExpiredRefreshTokenMsg, nil)
		} else {
			HandleUnableToConnectErrorAndExit(nil)
		}
	}

	tr := TokenResponse{}
	ProcessResponseFromServer(resp, &tr)
	return tr, nil
}
