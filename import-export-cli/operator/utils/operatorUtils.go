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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
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

// CreateControllerConfigs creates configs
func CreateControllerConfigs(configFile string, maxTimeSec int, resourceTypes ...string) {
	utils.Logln(utils.LogPrefixInfo + "Installing controller configs")

	// apply all files without printing errors
	if err := ExecuteCommandWithoutPrintingErrors(Kubectl, K8sApply, "-f", configFile); err != nil {
		fmt.Println("Waiting for resource creation...")

		// if error then wait for namespace and the resource type security
		if len(resourceTypes) > 0 {
			_ = K8sWaitForResourceType(maxTimeSec, resourceTypes...)
		}

		// apply again with printing errors
		if err := K8sApplyFromFile(configFile); err != nil {
			utils.HandleErrorAndExit("Error creating configurations", err)
		}
	}
}
