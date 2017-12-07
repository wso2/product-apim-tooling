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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"os"
	"path/filepath"
	"testing"
)

// TestAddEnv1 - Blank Env Name
func TestAddEnv1(t *testing.T) {
	envName := ""
	publisherEndpoint := "test-publisher-endpoint"
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"
	mainConfigFilePath := ""

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, mainConfigFilePath)
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}

}

// TestAddEnv2 - Blank Publisher Endpoint
func TestAddEnv2(t *testing.T) {
	envName := "test-env"
	publisherEndpoint := ""
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"
	mainConfigFilePath := ""

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, mainConfigFilePath)
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}

}

// TestAddEnv3 - Already existing environment
func TestAddEnv3(t *testing.T) {

	sampleMainConfigFileName := "sample_main_config.yaml"
	sampleMainConfigFilePath := filepath.Join(utils.ConfigDirPath, sampleMainConfigFileName)

	var sampleMainConnfig = new(utils.MainConfig)
	sampleMainConnfig.Config = utils.Config{10000, ""}
	sampleMainConnfig.Environments = make(map[string]utils.EnvEndpoints)
	sampleMainConnfig.Environments["dev"] = utils.EnvEndpoints{"sample-publisher-endpoint",
				"sample-reg-endpoint", "sample-token-endpoint"}
	utils.WriteConfigFile(sampleMainConnfig, sampleMainConfigFilePath)

	envName := "dev"
	publisherEndpoint := "sample-publisher-endpoint"
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, sampleMainConfigFilePath)
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}


	defer os.Remove(sampleMainConfigFilePath)
}

// TetsAddEnv4 - Correct Details - Successfully add new environment
func TestAddEnv4(t *testing.T) {
	sampleMainConfigFileName := "sample_main_config.yaml"
	sampleMainConfigFilePath := filepath.Join(utils.ConfigDirPath, sampleMainConfigFileName)

	var sampleMainConnfig = new(utils.MainConfig)
	sampleMainConnfig.Config = utils.Config{10000, ""}
	sampleMainConnfig.Environments = make(map[string]utils.EnvEndpoints)
	sampleMainConnfig.Environments["dev"] = utils.EnvEndpoints{"sample-publisher-endpoint",
		"sample-reg-endpoint", "sample-token-endpoint"}
	utils.WriteConfigFile(sampleMainConnfig, sampleMainConfigFilePath)

	envName := "new-env"
	publisherEndpoint := "sample-publisher-endpoint"
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, sampleMainConfigFilePath)
	if err != nil {
		t.Errorf("Expected nil, got '%s' instead\n", err.Error())
	}


	defer os.Remove(sampleMainConfigFilePath)
}

