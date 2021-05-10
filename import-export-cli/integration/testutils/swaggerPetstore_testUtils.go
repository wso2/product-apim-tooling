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

package testutils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	yaml2 "gopkg.in/yaml.v2"
)

func GetSwaggerPetstoreDefinition(t *testing.T, username string) string {
	path := "testdata/petstore.yaml"
	tempSwaggerPetstoreFileNameforUser := "testdata" + string(os.PathSeparator) + base.GenerateRandomString() + "-" + username + ".yaml"

	sampleData, _ := ioutil.ReadFile(path)

	// Extract the content to a structure
	petstoreDefinitionMap := make(map[interface{}]interface{})
	err := yaml2.Unmarshal(sampleData, &petstoreDefinitionMap)
	if err != nil {
		t.Error(err)
	}
	//fmt.Printf("--- m:\n%v\n\n", petstoreDefinitionMap)

	scope1 := base.GenerateRandomString() + username + "Scope1"
	scope2 := base.GenerateRandomString() + username + "Scope2"
	petstoreDefinitionMapKeys := keys(petstoreDefinitionMap)

	if contains(petstoreDefinitionMapKeys, "paths") {
		paths := petstoreDefinitionMap["paths"].(map[interface{}]interface{})
		for _, pathKey := range keys(paths) {
			resourcePath := paths[pathKey].(map[interface{}]interface{})
			for _, resourcePathKey := range keys(resourcePath) {
				resourcePathHtppVerb := resourcePath[resourcePathKey].(map[interface{}]interface{})
				resourcePathHtppVerbKeys := keys(resourcePathHtppVerb)
				if contains(resourcePathHtppVerbKeys, "security") {
					petstoreAuth := map[string]interface{}{
						"petstore_auth": []interface{}{scope1, scope2},
					}
					updatedSecurityField := []interface{}{petstoreAuth}
					resourcePathHtppVerb["security"] = updatedSecurityField
				}
			}
		}
	}

	if contains(petstoreDefinitionMapKeys, "securityDefinitions") {
		securityDefinitions := petstoreDefinitionMap["securityDefinitions"].(map[interface{}]interface{})
		for _, securityDefinitionKey := range keys(securityDefinitions) {
			securityDefinitionKeyProperties := securityDefinitions[securityDefinitionKey].(map[interface{}]interface{})
			securityDefinitionKeyPropertiesKeys := keys(securityDefinitionKeyProperties)
			if contains(securityDefinitionKeyPropertiesKeys, "scopes") {
				updatedScopesField := map[string]interface{}{
					scope1: "Description of " + scope1,
					scope2: "Description of " + scope2,
				}
				securityDefinitionKeyProperties["scopes"] = updatedScopesField
			}
		}
	}

	swaggerData, err := yaml2.Marshal(petstoreDefinitionMap)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile(tempSwaggerPetstoreFileNameforUser, swaggerData, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		// Remove extracted archive
		base.RemoveDir(tempSwaggerPetstoreFileNameforUser)
	})

	return tempSwaggerPetstoreFileNameforUser
}

// keys returns a string array of the keys of a map
func keys(m map[interface{}]interface{}) []interface{} {
	keys := make([]interface{}, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// contains checks if a string is present in a slice
func contains(s []interface{}, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
