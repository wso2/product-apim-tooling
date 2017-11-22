/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

	returnedEndpoint := GetAPIMEndpointOfEnv(devName, testMainConfigFilePath)
	expectedEndpoint := getSampleMainConfig().Environments[devName].APIManagerEndpoint
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

func TestAddNewEnvToKeysFile1(t *testing.T) {
	writeCorrectKeys()
	var envKeys = EnvKeys{"staging_username", "staging_client_id", "staging_client_secret"}

	AddNewEnvToKeysFile("staging", envKeys, testKeysFilePath)
	defer os.Remove(testKeysFilePath)
}

func TestAddNewEnvToKeysFile2(t *testing.T) {
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

	defer removeFiles()
}

// Case 2: Env not available in keys file
func TestRemoveEnvFromKeysFile2(t *testing.T) {
	WriteCorrectMainConfig()

	// write incorrect keys
	envKeysAll.Environments = make(map[string]EnvKeys)
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa_client_secret")
	envKeysAll.Environments[qaName] = EnvKeys{"qa_client_id", qaEncryptedClientSecret, qaUsername}
	WriteConfigFile(envKeysAll, testKeysFilePath)
	// end of writing incorrect keys

	err := RemoveEnvFromKeysFile(devName, testKeysFilePath, testMainConfigFilePath)
	if err == nil {
		t.Error("No error returned. 'Env not found in keys file' error expected")
	}

	defer removeFiles()
}

// Case 4: Incorrect Endpoints
func TestRemoveEnvFromKeysFile4(t *testing.T) {
	// writing incorrect endpoints
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[devName] = EnvEndpoints{"dev_apim_endpoint",
		"dev_reg_endpoint", "dev_token_endpoint"}
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

func removeFiles() {
	_ = os.Remove(testMainConfigFilePath)
	//fmt.Println("Error removing endpoints file:", err.Error())
	_ = os.Remove(testKeysFilePath)
	//fmt.Println("Error removing keys file:", err.Error())
}
