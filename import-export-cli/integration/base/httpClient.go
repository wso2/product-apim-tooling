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

package base

import (
	"crypto/tls"
	"encoding/base64"
	"io"
	"net/http"
)

// CreateGet : Construct GET http request
func CreateGet(url string) *http.Request {
	Log("base.CreateGet() - url:", url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		Fatal(err)
	}

	return req
}

// CreatePost : Construct POST http request
func CreatePost(url string, body io.Reader) *http.Request {
	Log("base.CreatePost() - url:", url)

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		Fatal(err)
	}

	return req
}

// CreatePostEmptyBody : Construct POST http request with empty body
func CreatePostEmptyBody(url string) *http.Request {
	Log("base.CreatePostEmptyBody() - url:", url)

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		Fatal(err)
	}

	return req
}

// CreateDelete : Construct DELETE http request
func CreateDelete(url string) *http.Request {
	Log("base.CreateDelete() - url:", url)

	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		Fatal(err)
	}

	return req
}

// SetDefaultRestAPIHeaders : Set HTTP headers required for APIM REST calls
func SetDefaultRestAPIHeaders(token string, request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
}

// SetDefaultRestAPIHeadersToConsumeFormData : Set HTTP headers required for APIM REST calls to consume form data
func SetDefaultRestAPIHeadersToConsumeFormData(token string, request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "multipart/form-data")
}

// SetTokenAPIHeaders : Set HTTP headers for token API invocation
func SetTokenAPIHeaders(clientID string, clientSecret string, request *http.Request) {
	authHeader := clientID + ":" + clientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(authHeader))

	request.Header.Set("Authorization", "Basic "+encoded)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

// SetSOAPHeaders : Set HTTP headers for SOAP invocation
func SetSOAPHeaders(userName string, password string, action string, request *http.Request) {
	authHeader := userName + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(authHeader))

	request.Header.Set("Authorization", "Basic "+encoded)
	request.Header.Set("Content-Type", "text/xml")
	request.Header.Set("SOAPAction", action)
}

// SendHTTPRequest : Send HTTP request
func SendHTTPRequest(request *http.Request) *http.Response {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(request)

	if err != nil {
		Fatal(err)
	}

	return resp
}
