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

// TestRemoveEnv1 - Correct Details
func TestRemoveEnv1(t *testing.T) {
	sampleMainConfigFileName := "sample_main_config.yaml"
	sampleMainConfigFilePath := filepath.Join(utils.ApplicationRoot, sampleMainConfigFileName)

	var sampleMainConfig = new(utils.MainConfig)
	sampleMainConfig.Config = utils.Config{10000, ""}
	sampleMainConfig.Environments = make(map[string]utils.EnvEndpoints)
	sampleMainConfig.Environments["dev"] = utils.EnvEndpoints{"sample-publisher-endpoint",
		"sample-reg-endpoint", "sample-token-endpoint"}
	utils.WriteConfigFile(sampleMainConfig, sampleMainConfigFilePath)

	sampleEnvKeysAllFileName := "sample_env_keys_all.yaml"
	sampleEnvKeysAllFilePath := filepath.Join(utils.ApplicationRoot, sampleEnvKeysAllFileName)
	var sampleEnvKeysAll = new(utils.EnvKeysAll)
	sampleEnvKeysAll.Environments = make(map[string]utils.EnvKeys)
	sampleEnvKeysAll.Environments["dev"] = utils.EnvKeys{"clien-id", "client-secret",
		"username"}
	utils.WriteConfigFile(sampleEnvKeysAll, sampleEnvKeysAllFilePath)

	err := removeEnv("dev", sampleMainConfigFilePath, sampleEnvKeysAllFilePath)

	if err != nil {
		t.Errorf("Expected nil, got '%s'\n", err.Error())
	}

	defer func() {
		os.Remove(sampleEnvKeysAllFilePath)
		os.Remove(sampleMainConfigFilePath)
		os.Remove(filepath.Join(utils.CurrentDir, utils.EnvKeysAllFileName))
	}()

}

// TestRemoveEnv2 - Remove blank env name
func TestRemoveEnv2(t *testing.T) {
	sampleMainConfigFileName := "sample_main_config.yaml"
	sampleMainConfigFilePath := filepath.Join(utils.ApplicationRoot, sampleMainConfigFileName)

	var sampleMainConfig = new(utils.MainConfig)
	sampleMainConfig.Config = utils.Config{10000, ""}
	sampleMainConfig.Environments = make(map[string]utils.EnvEndpoints)
	sampleMainConfig.Environments["dev"] = utils.EnvEndpoints{"sample-publisher-endpoint",
		"sample-reg-endpoint", "sample-token-endpoint"}
	utils.WriteConfigFile(sampleMainConfig, sampleMainConfigFilePath)

	err := removeEnv("", sampleMainConfigFilePath, "")
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}

	defer os.Remove(sampleMainConfigFilePath)
}

// TestRemoveEnv4 - Remove an environment that doesn't exist
func TestRemoveEnv3(t *testing.T) {
	sampleMainConfigFileName := "sample_main_config.yaml"
	sampleMainConfigFilePath := filepath.Join(utils.ApplicationRoot, sampleMainConfigFileName)

	var sampleMainConfig = new(utils.MainConfig)
	sampleMainConfig.Config = utils.Config{10000, ""}
	sampleMainConfig.Environments = make(map[string]utils.EnvEndpoints)
	sampleMainConfig.Environments["dev"] = utils.EnvEndpoints{"sample-publisher-endpoint",
		"sample-reg-endpoint", "sample-token-endpoint"}
	utils.WriteConfigFile(sampleMainConfig, sampleMainConfigFilePath)

	err := removeEnv("new-dev", sampleMainConfigFilePath, "")
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}

	defer os.Remove(sampleMainConfigFilePath)
}
