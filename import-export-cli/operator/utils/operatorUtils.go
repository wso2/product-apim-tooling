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
	"errors"
	"fmt"
	"net/http"
	"os"
)

// GetVersion returns version which is read from environment variable and verify it's existence
func GetVersion(name string, envVar string, defaultVersion string, versionValidationUrl string, findVersionUrl string) (string, error) {
	version := os.Getenv(envVar)
	if version == "" {
		version = defaultVersion

		// if error set it next time
		_ = os.Setenv(envVar, version)
	}

	resp, err := http.Head(fmt.Sprintf(versionValidationUrl, version))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf(
			"invalid %s version: %s\n"+
				"Set the environment variable \"%s\" with a valid %s version\n"+
				"Default %s version: %s\n"+
				"Find a version here: %s",
			name, version,
			envVar, name,
			name, defaultVersion,
			findVersionUrl,
		))
	}

	return version, nil
}
