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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/stretchr/testify/assert"
)

func TestExtractAPIInfoWithCorrectJSON(t *testing.T) {
	// Correct json
	content := `{
		"type": "api",
		"version": "v4.4.0",
		"data": {
		  "id": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
		  "name": "APIName",
		  "provider": "devops",
		  "version": "1.0.0"
		}
	  }`

	api, err := extractAPIDefinition([]byte(content))
	assert.Equal(t, err, nil, "Should return nil error for correct json")
	assert.Equal(t, api.Data.Name, "APIName", "Should parse correct json")
	assert.Equal(t, api.Data.Provider, "devops", "Should parse correct json")
	assert.Equal(t, api.Data.Version, "1.0.0", "Should parse correct json")
}

func TestExtractAPIInfoWhenDataTagMissing(t *testing.T) {
	// When ID tag missing
	content := `{
		"type": "api",
		"version": "v4.4.0"
	  }`

	api, err := extractAPIDefinition([]byte(content))
	assert.Nil(t, err, "Should return nil error")
	assert.Equal(t, v2.APIDTODefinition{}, api.Data, "Should return empty Data when ID tag missing")
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
	api, _, err := GetAPIDefinition(utils.GetRelativeTestDataPathFromImpl() + "PizzaShackAPI-1.0.0")
	assert.Nil(t, err, "Should return nil error on reading correct directories")
	assert.Equal(t, api.Data.Name, "PizzaShackAPI", "Should return correct values for API name")
	assert.Equal(t, api.Data.Provider, "admin", "Should return correct values for API provider")
	assert.Equal(t, api.Data.Version, "1.0.0", "Should return correct values for API version")
}

func TestGetAPIInfoMalformedDirectory(t *testing.T) {
	api, _, err := GetAPIDefinition("testdata/PizzaShackAPI_1.0.0-malformed")
	assert.Error(t, err, "Should return error on reading malformed directories")
	assert.True(t, os.IsNotExist(err), "File not found error must be thrown")
	assert.Nil(t, api,
		"Should return nil for malformed directories")
}
