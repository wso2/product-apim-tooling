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

package utils

import (
	"testing"
)

func TestGetValidK8sName(t *testing.T) {
	tests := []struct {
		testName, name, validName string
	}{
		{
			testName:  "Name_with_spaces",
			name:      "# Swagger Pet#store 1..0.5",
			validName: "swagger-pet-store-1-0-5",
		},
		{
			testName:  "Name_with_special_chars",
			name:      "#Swagger##Pet store$",
			validName: "swagger-pet-store",
		},
		{
			testName:  "Name_with_digits",
			name:      "123Swagger.123Pet2store$#",
			validName: "123swagger-123pet2store",
		},
		{
			testName:  "Name_with_empty_result",
			name:      "#@$%^",
			validName: "default",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			validName := GetValidK8sResourceName(test.name)
			if validName != test.validName {
				t.Errorf("got %s, want %s", validName, test.validName)
			}
		})
	}
}
