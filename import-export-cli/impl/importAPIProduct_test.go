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
	"fmt"
	"reflect"
	"testing"

	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func TestExtractAPIProductInfoWithCorrectJSON(t *testing.T) {
	// Correct json
	content := `{
	  "type": "api_product",
      "version": "v4.3.0",
	  "data": {
		"id": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
		"name": "APIProductName",
		"provider" : "devops"
	  }
	}`

	apiProduct, err := extractAPIProductDefinition([]byte(content))
	assert.Equal(t, err, nil, "Should return nil error for correct json")
	assert.Equal(t, apiProduct.Data, v2.APIProductDTODefinition{Provider: "devops", Name: "APIProductName"},
		"Should parse correct json")
}

func TestExtractAPIProductInfoWhenDataTagMissing(t *testing.T) {
	// When ID tag missing
	content := `{
		"type": "api_product",
		"version": "v4.3.0"
	  }`
	apiProduct, err := extractAPIProductDefinition([]byte(content))
	assert.Nil(t, err, "Should return nil error")
	assert.Equal(t, v2.APIProductDTODefinition{}, apiProduct.Data, "Should return empty Data when ID tag missing")
}

func TestExtractAPIProductInfoWithMalformedJSON(t *testing.T) {
	// Malformed json
	content := `{
		"type": "api_product",
		"version": "v4.3.0",
		"data": {
		  "id": "e4d0c1be-44e9-43ad-b434-f8e2f02dad11",
		  "name": "APIProductName",
		  "provider" : "devops"
	  }`

	apiProduct, err := extractAPIProductDefinition([]byte(content))
	assert.Nil(t, apiProduct, "Should return nil API Product struct")
	assert.Error(t, err, "Should return an error regarding malformed json")
}

func TestGetAPIProductInfoCorrectDirectoryStructure(t *testing.T) {
	apiProduct, _, err := GetAPIProductDefinition(utils.GetRelativeTestDataPathFromImpl() + "MyProduct-1.0.0")
	assert.Nil(t, err, "Should return nil error on reading correct directories")
	assert.Equal(t, v2.APIProductDTODefinition{Provider: "admin", Name: "MyProduct"}, apiProduct.Data,
		"Should return correct values for ID info")
}

func TestGetAPIProductInfoMalformedDirectory(t *testing.T) {
	apiProduct, _, err := GetAPIProductDefinition(utils.GetRelativeTestDataPathFromImpl() + "MyProduct-1.0.0-malformed")
	fmt.Println(reflect.TypeOf(err))
	assert.Error(t, err, "Should return error on reading malformed directories")
	assert.Contains(t, err.Error(), "was not found as a YAML or JSON", "Should contain this message")
	assert.Nil(t, apiProduct,
		"Should return nil for malformed directories")
}
