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

package impl

import (
	"os"
	"testing"

	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/stretchr/testify/assert"
)

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
	api, _, err := GetAPIDefinition("testdata/PizzaShackAPI-1.0.0")
	assert.Nil(t, err, "Should return nil error on reading correct directories")
	assert.Equal(t, v2.ID{APIName: "PizzaShackAPI", Version: "1.0.0", ProviderName: "admin"}, api.ID,
		"Should return correct values for ID info")
}

func TestGetAPIInfoMalformedDirectory(t *testing.T) {
	api, _, err := GetAPIDefinition("testdata/PizzaShackAPI_1.0.0-malformed")
	assert.Error(t, err, "Should return error on reading malformed directories")
	assert.True(t, os.IsNotExist(err), "File not found error must be thrown")
	assert.Nil(t, api,
		"Should return nil for malformed directories")
}
