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
	"os"
	"path/filepath"
	"testing"
)

func TestEnvExistsInKeysFile(t *testing.T) {
	writeCorrectKeys()

	returned := EnvExistsInKeysFile(devName, testKeysFilePath)

	if !returned {
		t.Errorf("Error in EnvExistsInKeysFile(). Returned: %t\n", returned)
	}

	returned = EnvExistsInKeysFile(qaName, testKeysFilePath)

	if !returned {
		t.Errorf("Error in EnvExistsInKeysFile(). Returned: %t\n", returned)
	}

	returned = EnvExistsInKeysFile("staging", testKeysFilePath) // not available
	if returned {
		t.Error("Error in EnvExistsInKeysFile(). False Positive")
	}
	defer os.Remove(testKeysFilePath)

}

func TestGetClientIDOfEnv(t *testing.T) {
	writeCorrectKeys()

	returnedKey := GetClientIDOfEnv(devName, testKeysFilePath)
	expectedKey := getSampleKeys().Environments[devName].ClientID
	if returnedKey != expectedKey {
		t.Errorf("Expected '%s', got '%s'\n", expectedKey, returnedKey)
	}
	defer os.Remove(testKeysFilePath)

}

func TestGetClientSecretOfEnv(t *testing.T) {
	writeCorrectKeys()

	returnedKey := GetClientSecretOfEnv(devName, devPassword, testKeysFilePath)
	expectedKey := Decrypt([]byte(GetMD5Hash(devPassword)), getSampleKeys().Environments[devName].ClientSecret)

	if returnedKey != expectedKey {
		t.Errorf("Expected '%s', got '%s'\n", expectedKey, returnedKey)
	}
	defer os.Remove(testKeysFilePath)
}

func TestGetUsernameOfEnv(t *testing.T) {
	writeCorrectKeys()

	returnedKey := GetUsernameOfEnv(devName, testKeysFilePath)
	expectedKey := getSampleKeys().Environments[devName].Username

	if returnedKey != expectedKey {
		t.Errorf("Expected '%s', got '%s'\n", expectedKey, returnedKey)
	}
	defer os.Remove(testKeysFilePath)
}

// file exists
func TestAddNewEnvToKeysFile1(t *testing.T) {
	writeCorrectKeys()

	if !IsFileExist(testKeysFilePath) {
		t.Error("test keys file does not exist. test is not going to function properly")
	}

	var envKeys = EnvKeys{"staging_username", "staging_client_id", "staging_client_secret"}

	AddNewEnvToKeysFile("staging", envKeys, testKeysFilePath)
	defer os.Remove(testKeysFilePath)
}

func TestAddNewEnvToKeysFile2(t *testing.T) {
	if IsFileExist(testKeysFilePath) {
		t.Error("test keys file exists. test is not going to function properly")
	}
	var envKeys = EnvKeys{"staging_username", "staging_client_id", "staging_client_secret"}

	AddNewEnvToKeysFile("staging", envKeys, testKeysFilePath)
	defer os.Remove(testKeysFilePath)
}

// Case 3: Environment is blank
func TestRemoveEnvFromKeysFile3(t *testing.T) {
	err := RemoveEnvFromKeysFile("", testKeysFilePath, testMainConfigFilePath)
	if err == nil {
		t.Error("No error returned. 'Env cannot be blank' error expected")
	}
}

func TestIsDefaultEnvPresent2(t *testing.T) {
	testMainConfigFileName := "test_main_config.yaml"
	testMainConfigFilePath := filepath.Join(CurrentDir, testMainConfigFileName)
	mainConfig := new(MainConfig)
	mainConfig.Environments = make(map[string]EnvEndpoints)

	WriteConfigFile(mainConfig, testMainConfigFilePath)

	isDefaultEnvPresent := IsDefaultEnvPresent(testMainConfigFilePath)
	if isDefaultEnvPresent {
		t.Errorf("Expected '%t', got '%t'\n", false, true)
	}

	defer os.Remove(testMainConfigFilePath)
}

// test case 2 - default env absent
func TestGetDefaultEnvironment2(t *testing.T) {
	testMainConfigFileName := "test_main_config.yaml"
	testMainConfigFilePath := filepath.Join(CurrentDir, testMainConfigFileName)
	mainConfig := new(MainConfig)
	mainConfig.Environments = make(map[string]EnvEndpoints)

	WriteConfigFile(mainConfig, testMainConfigFilePath)

	defaultEnv := GetDefaultEnvironment(testMainConfigFilePath)
	if defaultEnv != "" {
		t.Errorf("Expected '%s', got '%s'\n", " defaultEnv", "empty-string")
	}

	defer os.Remove(testMainConfigFilePath)
}

// test case 1 - input env blank
func TestRemoveEnvFromMainConfigFile1(t *testing.T) {
	err := RemoveEnvFromMainConfigFile("", "")
	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", err, "nil")
	}
}

func TestGetKeysOfEnvironment(t *testing.T) {
	var envKeysAll = new(EnvKeysAll)

	// write incorrect keys
	envKeysAll.Environments = make(map[string]EnvKeys)
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa_client_secret")
	envKeysAll.Environments[qaName] = EnvKeys{"qa_client_id", qaEncryptedClientSecret, qaUsername}
	WriteConfigFile(envKeysAll, testKeysFilePath)

	envKeys, err := GetKeysOfEnvironment("incorrect-env", testKeysFilePath)

	if envKeys != nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "nil", envKeys)
	}

	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "error", err)
	}

	defer os.Remove(testKeysFilePath)

}
