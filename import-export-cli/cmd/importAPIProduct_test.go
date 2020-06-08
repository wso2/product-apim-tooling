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
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/renstrom/dedent"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func TestImportAPIProduct1(t *testing.T) {
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

	name := "MyProduct-1.0.0"

	err := ImportAPIProduct(credentials.Credential{}, name, server.URL, "testdata")
	assert.Nil(t, err, "Error should be nil")

	utils.Insecure = true
	err = ImportAPIProduct(credentials.Credential{}, name, server.URL, "testdata")
	assert.Nil(t, err, "Error should be nil")
}

func TestExtractAPIProductInfoWithCorrectJSON(t *testing.T) {
	// Correct json
	content := `{
	  "id": {
		"providerName": "admin",
		"apiProductName": "APIProductName",
		"version": "1.0.0"
	  },
	  "uuid": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
	  "description": "Some API Product Description",
	  "type": "HTTP",
	  "context": "/api-product/1.0.0",
	  "contextTemplate": "/api-product/{version}",
	  "tags": [
		"api-product"
	  ]
	}`

	apiProduct, err := extractAPIProductDefinition([]byte(content))
	assert.Equal(t, err, nil, "Should return nil error for correct json")
	assert.Equal(t, apiProduct.ID, v2.ProductID{ProviderName: "admin", Version: "1.0.0", APIProductName: "APIProductName"},
		"Should parse correct json")
}

func TestExtractAPIProductInfoWhenIDTagMissing(t *testing.T) {
	// When ID tag missing
	content := `{
	  "description": "Some API Product Description",
	  "type": "HTTP",
	  "context": "/api-product/1.0.0",
	  "contextTemplate": "/api-product/{version}",
	  "tags": [
		"api-product"
	  ]
	}`

	apiProduct, err := extractAPIProductDefinition([]byte(content))
	assert.Nil(t, err, "Should return nil error")
	assert.Equal(t, v2.ProductID{}, apiProduct.ID, "Should return empty IDInfo when ID tag missing")
}

func TestExtractAPIProductInfoWithMalformedJSON(t *testing.T) {
	// Malformed json
	content := `{
	  "uuid": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
	  "description": "Some API Product Description",
	  "type": "HTTP",
	  "context": "/api-product/1.0.0",
	  "contextTemplate": "/api-product/{version}",
	  "tags": [
		"api-product"
	  
	}`

	apiProduct, err := extractAPIProductDefinition([]byte(content))
	assert.Nil(t, apiProduct, "Should return nil API Product struct")
	assert.Error(t, err, "Should return an error regarding malformed json")
}

func TestGetAPIProductInfoCorrectDirectoryStructure(t *testing.T) {
	apiProduct, _, err := getAPIProductDefinition("testdata/MyProduct-1.0.0")
	assert.Nil(t, err, "Should return nil error on reading correct directories")
	assert.Equal(t, v2.ProductID{APIProductName: "MyProduct", Version: "1.0.0", ProviderName: "admin"}, apiProduct.ID,
		"Should return correct values for ID info")
}

func TestGetAPIProductInfoMalformedDirectory(t *testing.T) {
	apiProduct, _, err := getAPIProductDefinition("testdata/MyProduct-1.0.0-malformed")
	fmt.Println(reflect.TypeOf(err))
	assert.Error(t, err, "Should return error on reading malformed directories")
	assert.Contains(t, err.Error(), "was not found as a YAML or JSON", "Should contain this message")
	assert.Nil(t, apiProduct,
		"Should return nil for malformed directories")
}
