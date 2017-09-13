package utils

import (
	"testing"
	"os"
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
	writeCorrectEndpoints()

	returned := EnvExistsInEndpointsFile(devName, testEndpointsFilePath)

	if !returned {
		t.Errorf("Error in EnvExistsInEndpointsFile(). Returned: %t\n", returned)
	}

	returned = EnvExistsInEndpointsFile(qaName, testEndpointsFilePath)

	if !returned {
		t.Errorf("Error in EnvExistsInEndpointsFile(). Returned: %t\n", returned)
	}

	returned = EnvExistsInEndpointsFile("staging", testEndpointsFilePath) // not available
	if returned {
		t.Error("Error in EnvExistsInEndpointsFile(). False Positive")
	}
	defer os.Remove(testEndpointsFilePath)

}

func TestGetAPIMEndpointOfEnv(t *testing.T) {
	writeCorrectEndpoints()

	returnedEndpoint := GetAPIMEndpointOfEnv(devName, testEndpointsFilePath)
	expectedEndpoint := getSampleEndpoints().Environments[devName].APIManagerEndpoint
	if returnedEndpoint != expectedEndpoint {
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	defer os.Remove(testEndpointsFilePath)

}

func TestGetTokenEndpointOfEnv(t *testing.T) {
	writeCorrectEndpoints()

	returnedEndpoint := GetTokenEndpointOfEnv(devName, testEndpointsFilePath)
	expectedEndpoint := getSampleEndpoints().Environments[devName].TokenEndpoint
	if returnedEndpoint != expectedEndpoint {
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	defer os.Remove(testEndpointsFilePath)

}

func TestGetRegistrationEndpointOfEnv(t *testing.T) {
	writeCorrectEndpoints()

	returnedEndpoint := GetRegistrationEndpointOfEnv(devName, testEndpointsFilePath)
	expectedEndpoint := getSampleEndpoints().Environments[devName].RegistrationEndpoint
	if returnedEndpoint != expectedEndpoint {
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	defer os.Remove(testEndpointsFilePath)

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
	var envKeys EnvKeys = EnvKeys{"staging_username", "staging_client_id", "staging_client_secret"}

	AddNewEnvToKeysFile("staging", envKeys, testKeysFilePath)
	defer os.Remove(testKeysFilePath)
}

func TestAddNewEnvToKeysFile2(t *testing.T) {
	var envKeys EnvKeys = EnvKeys{"staging_username", "staging_client_id", "staging_client_secret"}

	AddNewEnvToKeysFile("staging", envKeys, testKeysFilePath)
	defer os.Remove(testKeysFilePath)
}

// Case 1: Correct Details
func TestRemoveEnvFromKeysFile1(t *testing.T) {
	writeCorrectEndpoints()
	writeCorrectKeys()
	err := RemoveEnvFromKeysFile(devName, testKeysFilePath, testEndpointsFilePath)
	if err != nil {
		t.Error("Error removing env from keys file: " + err.Error())
	}

	defer removeFiles()
}

// Case 2: Env not available in keys file
func TestRemoveEnvFromKeysFile2(t *testing.T) {
	writeCorrectEndpoints()

	// write incorrect keys
	envKeysAll.Environments = make(map[string]EnvKeys)
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa_client_secret")
	envKeysAll.Environments[qaName] = EnvKeys{"qa_client_id", qaEncryptedClientSecret, qaUsername}
	WriteConfigFile(envKeysAll, testKeysFilePath)
	// end of writing incorrect keys

	err := RemoveEnvFromKeysFile(devName, testKeysFilePath, testEndpointsFilePath)
	if err == nil {
		t.Error("No error returned. 'Env not found in keys file' error expected")
	}

	defer removeFiles()
}

// Case 4: Incorrect Endpoints
func TestRemoveEnvFromKeysFile4(t *testing.T) {
	// writing incorrect endpoints
	envEndpointsAll.Environments = make(map[string]EnvEndpoints)
	envEndpointsAll.Environments[devName] = EnvEndpoints{"dev_apim_endpoint",
		"dev_reg_endpoint", "dev_token_endpoint"}
	WriteConfigFile(envEndpointsAll, testEndpointsFilePath)
	// end of writing incorrect endpoints

	err := RemoveEnvFromKeysFile(qaName, testKeysFilePath, testEndpointsFilePath)
	if err == nil {
		t.Error("No error returned. 'Env not found in endpoints file' error expected")
	}

	defer os.Remove(testEndpointsFilePath)
}

// Case 3: Environment is blank
func TestRemoveEnvFromKeysFile3(t *testing.T) {
	err := RemoveEnvFromKeysFile("", testKeysFilePath, testEndpointsFilePath)
	if err == nil {
		t.Error("No error returned. 'Env cannot be blank' error expected")
	}
}

func removeFiles() {
	_ = os.Remove(testEndpointsFilePath)
	//fmt.Println("Error removing endpoints file:", err.Error())
	_ = os.Remove(testKeysFilePath)
	//fmt.Println("Error removing keys file:", err.Error())
}
