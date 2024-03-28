/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/tlsutils"
)

// Scope - token scope
type Scope string

const (
	// APIImportRelativePath is the relative path of API import in publisher rest API
	APIImportRelativePath = "api/am/publisher/v4/apis/import?preserveProvider=false&overwrite=true"
	// TokenRelativePath is the relative path for getting token in publisher rest API
	TokenRelativePath = "oauth2/token"
	// APIDeleteRelativePath is the relative path of delete api in publisher rest API
	APIDeleteRelativePath = "api/am/publisher/v4/apis/"
	payloadJSON           = `{
        "callbackUrl": "www.google.lk",
        "clientName": "rest_api_publisher",
        "owner": "admin",
        "grantType": "client_credentials password refresh_token",
        "saasApp": true
    }`
	// AdminScope admin scope
	AdminScope Scope = "apim:admin"
	// ImportExportScope import export api scope
	ImportExportScope Scope = "apim:api_import_export"
)

var (
	tokenURL             string
	apiImportURL         string
	apiDeleteURL         string
	username             string
	password             string
	skipSSL              bool
	clientID             string
	clientSecret         string
	basicAuthHeaderValue string
)

func init() {
	// Read configurations and derive the eventHub details
	conf, errReadConfig := config.ReadConfigs()
	if errReadConfig != nil {
		// This has to be error. For debugging purpose info
		logger.LoggerUtils.Errorf("Error reading configs: %v", errReadConfig)
	}
	// Populate data from the config
	cpConfigs := conf.ControlPlane
	cpURL := cpConfigs.ServiceURL
	// If the eventHub URL is configured with trailing slash
	if strings.HasSuffix(cpURL, "/") {
		apiImportURL = cpURL + APIImportRelativePath
		tokenURL = cpURL + TokenRelativePath
		apiDeleteURL = cpURL + APIDeleteRelativePath
	} else {
		apiImportURL = cpURL + "/" + APIImportRelativePath
		tokenURL = cpURL + "/" + TokenRelativePath
		apiDeleteURL = cpURL + "/" + APIDeleteRelativePath
	}
	username = cpConfigs.Username
	password = cpConfigs.Password
	clientID = cpConfigs.ClientID
	clientSecret = cpConfigs.ClientSecret
	skipSSL = cpConfigs.SkipSSLVerification

	// If clientId and clientSecret is not provided use username and password as basic auth to access rest apis.
	basicAuthHeaderValue = GetBasicAuthHeaderValue(username, password)
}

// Base64EncodeCredentials encodes the given username and password into a base64 string.
func Base64EncodeCredentials(username, password string) string {
	credentials := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(credentials))
}

// GetBasicAuthHeaderValue constructs and returns the Basic authentication header value using the provided username and password.
func GetBasicAuthHeaderValue(username, password string) string {
	return fmt.Sprintf("Basic %s", Base64EncodeCredentials(username, password))
}

// GetToken retrieves an OAuth token using the provided credentials and scopes.
func GetToken(scopes []string, clientID string, clientSecret string) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", strings.Join(scopes, " "))
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", GetBasicAuthHeaderValue(clientID, clientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check for non-200 response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	// Extract access_token
	accessToken, ok := response["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in response")
	}

	return accessToken, nil
}

// GetSuitableAuthHeadervalue returns an appropriate authentication header value based on whether client credentials are provided.
func GetSuitableAuthHeadervalue(scopes []string) (string, error) {
	if clientID != "" && clientSecret != "" {
		token, err := GetToken(scopes, clientID, clientSecret)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Bearer %s", token), nil
	}
	return basicAuthHeaderValue, nil
}

// ImportAPI imports an API from a zip file, returning the ID of the imported API.
func ImportAPI(apiZipName string, zipFileBytes *bytes.Buffer) (string, error) {
	authHeaderVal, err := GetSuitableAuthHeadervalue([]string{string(AdminScope), string(ImportExportScope)})
	if err != nil {
		return "", err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", apiZipName)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, zipFileBytes); err != nil {
		return "", err
	}
	writer.Close()
	req, err := http.NewRequest("POST", apiImportURL, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", authHeaderVal)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		logger.LoggerTLSUtils.Infof("API already exists in the CP hence ignoring the event. API zip name %s", apiZipName)
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	// try to parse the body as json and extract id from the response.
	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		// TODO after APIM is able to send json response, we should return error here, Until then return nil for error as its expected.
		return "", nil
	}

	// Assuming the response contains an ID field, you can extract it like this:
	id, ok := responseMap["id"].(string)
	if !ok {
		return "", nil
	}
	return id, nil
}

// DeleteAPI deletes an API given its UUID.
func DeleteAPI(apiUUID string) error {
	deleteURL := apiDeleteURL + apiUUID
	authheaderval, err := GetSuitableAuthHeadervalue([]string{string(AdminScope), string(ImportExportScope)})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", authheaderval)
	req.Header.Set("Content-Type", "application/json")
	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error occured while deleting the API. Status: %s", resp.Status)
	}

	return nil
}
