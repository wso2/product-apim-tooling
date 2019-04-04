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

package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/renstrom/dedent"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func TestImportAPI1(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected '%s', got '%s' instead\n", http.MethodPost, r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)
		w.Header().Set(utils.HeaderContentEncoding, utils.HeaderValueGZIP)
		w.Header().Set(utils.HeaderTransferEncoding, utils.HeaderValueChunked)

		body := dedent.Dedent(`
		`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	name := "sampleapi.zip"
	accessToken := "access-token"

	_, err := ImportAPI(name, server.URL, accessToken, "")
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
	utils.Insecure = true
	_, err = ImportAPI(name, server.URL, accessToken, "")
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}

func TestNewFileUploadRequest(t *testing.T) {
	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected '%s', got '%s' instead\n", http.MethodPut, r.Method)
		}

		if !strings.Contains(r.Header.Get(utils.HeaderAccept), utils.HeaderValueMultiPartFormData) {
			t.Errorf("Expected '%s', got '%s' instead\n", utils.HeaderValueApplicationZip,
				r.Header.Get(utils.HeaderContentType))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set(utils.HeaderContentType, utils.HeaderValueApplicationJSON)
		w.Header().Set(utils.HeaderContentEncoding, utils.HeaderValueGZIP)
		w.Header().Set(utils.HeaderTransferEncoding, utils.HeaderValueChunked)

		body := dedent.Dedent(`
		`)

		w.Write([]byte(body))
	}))
	defer server.Close()

	extraParams := map[string]string{}
	filePath := filepath.Join("sampleapi.zip")
	accessToken := "access-token"
	_, err := NewFileUploadRequest(server.URL, extraParams, "file", filePath, accessToken)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}

func TestExtractAPIInfo(t *testing.T) {
	// Correct json
	content := `{
	  "id": {
		"providerName": "admin",
		"apiName": "APIName",
		"version": "1.0.0"
	  },
	  "uuid": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
	  "description": "Some API Description",
	  "type": "HTTP",
	  "context": "/api/1.0.0",
	  "contextTemplate": "/api/{version}",
	  "tags": [
		"api"
	  ]
	}`

	api, err := extractAPIInfo([]byte(content))
	assert.Equal(t, api, &API{IdInfo{Provider: "admin", Version: "1.0.0", Name: "APIName"}},
		"Should parse correct json")
	assert.Equal(t, err, nil, "Should return nil error for correct json")

	// When ID tag missing
	content = `{
	  "uuid": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
	  "description": "Some API Description",
	  "type": "HTTP",
	  "context": "/api/1.0.0",
	  "contextTemplate": "/api/{version}",
	  "tags": [
		"api"
	  ]
	}`

	api, err = extractAPIInfo([]byte(content))
	assert.Equal(t, &API{}, api, "Should return empty IDInfo when ID tag missing")
	assert.Nil(t, err, "Should return nil error")

	// Malformed json
	content = `{
	  "uuid": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
	  "description": "Some API Description",
	  "type": "HTTP",
	  "context": "/api/1.0.0",
	  "contextTemplate": "/api/{version}",
	  "tags": [
		"api"
	  
	}`

	api, err = extractAPIInfo([]byte(content))
	assert.Nil(t, api, "Should return nil API struct")
	assert.Error(t, err, "Should return an error regarding malformed json")
}

func TestGetAPIInfo(t *testing.T) {
	api, err := getAPIInfo("testdata/PizzaShackAPI_1.0.0.zip")
	assert.Nil(t, err, "Should return nil error on reading correct zip files")
	assert.Equal(t, &API{IdInfo{Name: "PizzaShackAPI", Version: "1.0.0", Provider: "admin"}}, api,
		"Should return correct values for ID info")

	api, err = getAPIInfo("testdata/PizzaShackAPI-1.0.0")
	assert.Nil(t, err, "Should return nil error on reading correct directories")
	assert.Equal(t, &API{IdInfo{Name: "PizzaShackAPI", Version: "1.0.0", Provider: "admin"}}, api,
		"Should return correct values for ID info")

	api, err = getAPIInfo("testdata/PizzaShackAPI_1.0.0-malformed.zip")
	assert.Error(t, err, "Should return error on reading malformed zip files")
	assert.True(t, os.IsNotExist(err), "File not found error must be thrown")
	assert.Nil(t, api,
		"Should return nil for malformed directories")

	api, err = getAPIInfo("testdata/PizzaShackAPI_1.0.0-malformed")
	assert.Error(t, err, "Should return error on reading malformed directories")
	assert.True(t, os.IsNotExist(err), "File not found error must be thrown")
	assert.Nil(t, api,
		"Should return nil for malformed directories")
}
