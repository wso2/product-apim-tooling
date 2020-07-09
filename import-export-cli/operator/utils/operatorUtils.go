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
	"bytes"
	"errors"
	"fmt"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// CreateControllerConfigs apply (kubectl apply) configs to the k8s cluster
func CreateControllerConfigs(configFile string, maxTimeSec int, resourceTypes ...string) {
	configData := *readConfigData(configFile)

	// filter CRDs and other configs
	type YAML map[string]interface{}
	var crds []YAML
	var nonCrds []YAML
	for _, data := range configData {
		dec := yaml.NewDecoder(bytes.NewReader(data))
		for yml := make(YAML); dec.Decode(&yml) == nil; yml = make(YAML) {
			if strings.EqualFold(fmt.Sprint(yml[kindKey]), CrdKind) ||
				strings.EqualFold(fmt.Sprint(yml[kindKey]), namespaceKey) {
				crds = append(crds, yml)
			} else {
				nonCrds = append(nonCrds, yml)
			}
		}
	}

	// applying CRDs and namespaces
	crdsData := make([][]byte, 0, 5) // make capacity as CRD count for high performance
	for _, crd := range crds {
		data, err := yaml.Marshal(crd)
		if err != nil {
			utils.HandleErrorAndExit("Error parsing yaml content", err)
		}
		crdsData = append(crdsData, data)
	}
	if len(crdsData) > 0 {
		// apply all crds once to lower request count to k8s cluster
		err := K8sApplyFromBytes(crdsData)
		if err != nil {
			utils.HandleErrorAndExit("Error applying CRDs to K8s cluster", err)
		}
	}

	// applying non CRD configs
	nonCrdsData := make([][]byte, 0, 16) // make capacity for high performance
	for _, nonCrd := range nonCrds {
		data, err := yaml.Marshal(nonCrd)
		if err != nil {
			utils.HandleErrorAndExit("Error parsing yaml content", err)
		}
		nonCrdsData = append(nonCrdsData, data)
	}
	if len(nonCrdsData) > 0 {
		// waiting for resource creation if CRDs are applied
		if len(crdsData) > 0 {
			fmt.Println("Waiting for resource creation...")
			// if error then wait for namespace and the resource type security
			if len(resourceTypes) > 0 {
				_ = K8sWaitForResourceType(maxTimeSec, resourceTypes...)
			}
		}

		// apply all configs once to lower request count to k8s cluster
		err := K8sApplyFromBytes(nonCrdsData)
		if err != nil {
			utils.HandleErrorAndExit("Error applying configs to k8s cluster", err)
		}
	}
}

// readConfigData reads content of configFile from configFile of type: URL, local file or dir
func readConfigData(configFile string) *[][]byte {
	// read from URL
	if utils.IsValidUrl(configFile) {
		utils.Logln(utils.LogPrefixInfo + "Installing controller configs using URL")

		data, err := utils.ReadFromUrl(configFile)
		if err != nil {
			utils.HandleErrorAndExit("Error reading configs from URL: "+configFile, err)
		}
		return &[][]byte{data}
	}

	// read from local file or dir
	if stat, err := os.Stat(configFile); !os.IsNotExist(err) {
		// local file
		if !stat.IsDir() {
			utils.Logln(utils.LogPrefixInfo + "Installing controller configs using local file")

			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				utils.HandleErrorAndExit("Error reading configs from local file: "+configFile, err)
			}
			return &[][]byte{data}
		}

		// local dir
		utils.Logln(utils.LogPrefixInfo + "Installing controller configs using local dir")

		configDir, err := ioutil.ReadDir(configFile)
		if err != nil {
			utils.HandleErrorAndExit("Error reading configs from local dir: "+configFile, err)
		}
		var configData [][]byte
		for _, file := range configDir {
			if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
				f := filepath.Join(configFile, file.Name())
				data, err := ioutil.ReadFile(f)
				if err != nil {
					utils.HandleErrorAndExit("Error reading configs from local file: "+f, err)
				}
				configData = append(configData, data)
			}
		}
		return &configData
	} else {
		utils.HandleErrorAndExit("Error reading configs", errors.New("config file does not exists"))
		return nil
	}
}
