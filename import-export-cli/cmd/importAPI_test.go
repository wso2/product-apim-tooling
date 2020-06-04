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

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

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

	name := "PizzaShackAPI-1.0.0"

	err := ImportAPI(credentials.Credential{}, name, server.URL, "testdata", "")
	assert.Nil(t, err, "Error should be nil")

	utils.Insecure = true
	err = ImportAPI(credentials.Credential{}, name, server.URL, "testdata", "")
	assert.Nil(t, err, "Error should be nil")
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
	filePath := filepath.FromSlash("testdata/sampleapi.zip")
	accessToken := "access-token"
	_, err := NewFileUploadRequest(server.URL, http.MethodPost, extraParams, "file", filePath, accessToken)
	if err != nil {
		t.Errorf("Error: %s\n", err.Error())
	}
}

func TestExtractAPIInfoWithCorrectJSON(t *testing.T) {
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

	api, err := extractAPIDefinition([]byte(content))
	assert.Equal(t, err, nil, "Should return nil error for correct json")
	assert.Equal(t, api.ID, v2.ID{ProviderName: "admin", Version: "1.0.0", APIName: "APIName"},
		"Should parse correct json")
}

func TestExtractAPIInfoWhenIDTagMissing(t *testing.T) {
	// When ID tag missing
	content := `{
	  "description": "Some API Description",
	  "type": "HTTP",
	  "context": "/api/1.0.0",
	  "contextTemplate": "/api/{version}",
	  "tags": [
		"api"
	  ]
	}`

	api, err := extractAPIDefinition([]byte(content))
	assert.Nil(t, err, "Should return nil error")
	assert.Equal(t, v2.ID{}, api.ID, "Should return empty IDInfo when ID tag missing")
}

func TestExtractAPIInfoWithMalformedJSON(t *testing.T) {
	// Malformed json
	content := `{
	  "uuid": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
	  "description": "Some API Description",
	  "type": "HTTP",
	  "context": "/api/1.0.0",
	  "contextTemplate": "/api/{version}",
	  "tags": [
		"api"
	  
	}`

	api, err := extractAPIDefinition([]byte(content))
	assert.Nil(t, api, "Should return nil API struct")
	assert.Error(t, err, "Should return an error regarding malformed json")
}

func TestGetAPIInfoCorrectDirectoryStructure(t *testing.T) {
	api, _, err := getAPIDefinition("testdata/PizzaShackAPI-1.0.0")
	assert.Nil(t, err, "Should return nil error on reading correct directories")
	assert.Equal(t, v2.ID{APIName: "PizzaShackAPI", Version: "1.0.0", ProviderName: "admin"}, api.ID,
		"Should return correct values for ID info")
}

func TestGetAPIInfoMalformedDirectory(t *testing.T) {
	api, _, err := getAPIDefinition("testdata/PizzaShackAPI_1.0.0-malformed")
	assert.Error(t, err, "Should return error on reading malformed directories")
	assert.True(t, os.IsNotExist(err), "File not found error must be thrown")
	assert.Nil(t, api,
		"Should return nil for malformed directories")
}
