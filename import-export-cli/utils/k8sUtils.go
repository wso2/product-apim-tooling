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
	"regexp"
	"strings"
)

// GetValidK8sName returns a valid name from given name
func GetValidK8sName(name string) string {
	// a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.',
	// and must start and end with an alphanumeric character e.g. 'example.com',
	// regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*'

	// replace all special chars with "-"
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err) // no errors if regex in valid
	}
	// trim "-" if found
	replacedName := strings.Trim(reg.ReplaceAllString(name, "-"), "-")
	if replacedName == "" {
		return "default"
	}
	return strings.ToLower(replacedName)
}
