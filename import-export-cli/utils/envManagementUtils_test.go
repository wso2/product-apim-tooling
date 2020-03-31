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

func TestEnvExistsInEndpointsFile(t *testing.T) {
	WriteCorrectMainConfig()

	returned := EnvExistsInMainConfigFile(devName, testMainConfigFilePath)

	if !returned {
		t.Errorf("Error in EnvExistsInMainConfigFile(). Returned: %t\n", returned)
	}

	returned = EnvExistsInMainConfigFile(qaName, testMainConfigFilePath)

	if !returned {
		t.Errorf("Error in EnvExistsInMainConfigFile(). Returned: %t\n", returned)
	}

	returned = EnvExistsInMainConfigFile("staging", testMainConfigFilePath) // not available
	if returned {
		t.Error("Error in EnvExistsInMainConfigFile(). False Positive")
	}
	defer os.Remove(testMainConfigFilePath)

}

func TestGetAPIMEndpointOfEnv(t *testing.T) {
	WriteCorrectMainConfig()

	returnedEndpoint := GetApiManagerEndpointOfEnv(devName, testMainConfigFilePath)
	expectedEndpoint := getSampleMainConfig().Environments[devName].ApiManagerEndpoint
	if returnedEndpoint != expectedEndpoint {
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	defer os.Remove(testMainConfigFilePath)

}

func TestGetTokenEndpointOfEnv(t *testing.T) {
	WriteCorrectMainConfig()

	returnedEndpoint := GetTokenEndpointOfEnv(devName, testMainConfigFilePath)
	expectedEndpoint := getSampleMainConfig().Environments[devName].TokenEndpoint
	if returnedEndpoint != expectedEndpoint {
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	defer os.Remove(testMainConfigFilePath)

}

func TestGetRegistrationEndpointOfEnv(t *testing.T) {
	WriteCorrectMainConfig()

	returnedEndpoint := GetRegistrationEndpointOfEnv(devName, testMainConfigFilePath)
	expectedEndpoint := getSampleMainConfig().Environments[devName].RegistrationEndpoint
	if returnedEndpoint != expectedEndpoint {
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	defer os.Remove(testMainConfigFilePath)

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

// Case 1: Correct Details
func TestRemoveEnvFromKeysFile1(t *testing.T) {
	WriteCorrectMainConfig()
	writeCorrectKeys()
	err := RemoveEnvFromKeysFile(devName, testKeysFilePath, testMainConfigFilePath)
	if err != nil {
		t.Error("Error removing env from keys file: " + err.Error())
	}

	defer func() {
		os.Remove(testMainConfigFilePath)
		os.Remove(testKeysFilePath)
	}()
}

// Case 2: Env not available in keys file
func TestRemoveEnvFromKeysFile2(t *testing.T) {
	WriteCorrectMainConfig()

	// write incorrect keys
	envKeysAll.Environments = make(map[string]EnvKeys)
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa_client_secret")
	envKeysAll.Environments[qaName] = EnvKeys{"qa_client_id", qaEncryptedClientSecret, qaUsername}
	WriteConfigFile(envKeysAll, testKeysFilePath)

	err := RemoveEnvFromKeysFile(devName, testKeysFilePath, testMainConfigFilePath)
	if err == nil {
		t.Error("No error returned. 'Env not found in keys file' error expected")
	}

	defer func() {
		os.Remove(testMainConfigFilePath)
		os.Remove(testKeysFilePath)
	}()
}

// Case 4: Incorrect Endpoints
func TestRemoveEnvFromKeysFile4(t *testing.T) {
	// writing incorrect endpoints
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[devName] = EnvEndpoints{
		"dev_apim_endpoint",
		"dev_publisher_endpoint",
		"dev_devportal_endpoint",
		"dev_reg_endpoint",
		"dev_admin_endpoint",
		"dev_token_endpoint",
	}
	WriteConfigFile(mainConfig, testMainConfigFilePath)
	// end of writing incorrect endpoints

	err := RemoveEnvFromKeysFile(qaName, testKeysFilePath, testMainConfigFilePath)
	if err == nil {
		t.Error("No error returned. 'Env not found in endpoints file' error expected")
	}

	defer os.Remove(testMainConfigFilePath)
}

// Case 3: Environment is blank
func TestRemoveEnvFromKeysFile3(t *testing.T) {
	err := RemoveEnvFromKeysFile("", testKeysFilePath, testMainConfigFilePath)
	if err == nil {
		t.Error("No error returned. 'Env cannot be blank' error expected")
	}
}

func TestIsDefaultEnvPresent1(t *testing.T) {
	testMainConfigFileName := "test_main_config.yaml"
	testMainConfigFilePath := filepath.Join(CurrentDir, testMainConfigFileName)
	mainConfig := new(MainConfig)
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[DefaultEnvironmentName] = EnvEndpoints{
		"default-publisher",
		"default-api-list",
		"default-application-list",
		"default-reg",
		"default-admin",
		"default-token",
	}

	WriteConfigFile(mainConfig, testMainConfigFilePath)

	isDefaultEnvPresent := IsDefaultEnvPresent(testMainConfigFilePath)
	if !isDefaultEnvPresent {
		t.Errorf("Expected '%t', got '%t'\n", true, false)
	}

	defer os.Remove(testMainConfigFilePath)
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

// test case 1 - default env present
func TestGetDefaultEnvironment1(t *testing.T) {
	testMainConfigFileName := "test_main_config.yaml"
	testMainConfigFilePath := filepath.Join(CurrentDir, testMainConfigFileName)
	mainConfig := new(MainConfig)
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[DefaultEnvironmentName] = EnvEndpoints{
		"default-publisher",
		"default-api-list",
		"default-application-list",
		"default-reg",
		"default-admin",
		"default-token",
	}

	WriteConfigFile(mainConfig, testMainConfigFilePath)

	defaultEnv := GetDefaultEnvironment(testMainConfigFilePath)
	if defaultEnv == "" {
		t.Errorf("Expected '%s', got '%s'\n", "emtpy-string", defaultEnv)
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

// test case 2 - input env valid and available in file
func TestRemoveEnvFromMainConfigFile2(t *testing.T) {
	testMainConfigFileName := "test_main_config.yaml"
	testMainConfigFilePath := filepath.Join(CurrentDir, testMainConfigFileName)
	mainConfig := new(MainConfig)
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments["dev"] = EnvEndpoints{
		"default-publisher",
		"default-api-list",
		"default-application-list",
		"default-reg",
		"default-admin",
		"default-token",
	}

	WriteConfigFile(mainConfig, testMainConfigFilePath)
	err := RemoveEnvFromMainConfigFile("dev", testMainConfigFilePath)
	if err != nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "nil", err)
	}

	defer os.Remove(testMainConfigFilePath)
}

// test case 3 - input env valid but not available in file
func TestRemoveEnvFromMainConfigFile3(t *testing.T) {
	testMainConfigFileName := "test_main_config.yaml"
	testMainConfigFilePath := filepath.Join(CurrentDir, testMainConfigFileName)
	mainConfig := new(MainConfig)
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments["dev"] = EnvEndpoints{
		"default-publisher",
		"default-api-list",
		"default-application-list",
		"default-reg",
		"default-admin",
		"default-token",
	}

	WriteConfigFile(mainConfig, testMainConfigFilePath)
	err := RemoveEnvFromMainConfigFile("not-available", testMainConfigFilePath)
	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", err, "nil")
	}

	defer os.Remove(testMainConfigFilePath)
}

func TestGetEndpointsOfEnvironment(t *testing.T) {
	mainConfig := new(MainConfig)
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments["dev"] = EnvEndpoints{
		"default-publisher",
		"default-api-list",
		"default-application-list",
		"default-reg",
		"default-admin",
		"default-token",
	}

	WriteConfigFile(mainConfig, testMainConfigFilePath)
	endpoints, err := GetEndpointsOfEnvironment("not-available", testMainConfigFilePath)

	if endpoints != nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "nil", endpoints)
	}

	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", err, "nil")
	}

	defer os.Remove(testMainConfigFilePath)

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
