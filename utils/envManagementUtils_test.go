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

	returned = EnvExistsInKeysFile("staging", testKeysFilePath)		// not available
	if returned {
		t.Error("Error in EnvExistsInKeysFile(). False Positive")
	}
	os.Remove(testKeysFilePath)

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

	returned = EnvExistsInEndpointsFile("staging", testEndpointsFilePath)		// not available
	if returned {
		t.Error("Error in EnvExistsInEndpointsFile(). False Positive")
	}
	os.Remove(testEndpointsFilePath)

}

func TestGetAPIMEndpointOfEnv(t *testing.T) {
	writeCorrectEndpoints()

	returnedEndpoint := GetAPIMEndpointOfEnv(devName, testEndpointsFilePath)
	expectedEndpoint := getSampleEndpoints().Environments[devName].APIManagerEndpoint
	if returnedEndpoint != expectedEndpoint{
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	os.Remove(testEndpointsFilePath)

}

func TestGetTokenEndpointOfEnv(t *testing.T) {
	writeCorrectEndpoints()

	returnedEndpoint := GetTokenEndpointOfEnv(devName, testEndpointsFilePath)
	expectedEndpoint := getSampleEndpoints().Environments[devName].TokenEndpoint
	if returnedEndpoint != expectedEndpoint{
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	os.Remove(testEndpointsFilePath)

}

func TestGetRegistrationEndpointOfEnv(t *testing.T) {
	writeCorrectEndpoints()

	returnedEndpoint := GetRegistrationEndpointOfEnv(devName, testEndpointsFilePath)
	expectedEndpoint := getSampleEndpoints().Environments[devName].RegistrationEndpoint
	if returnedEndpoint != expectedEndpoint{
		t.Errorf("Expected '%s', got '%s'\n", expectedEndpoint, returnedEndpoint)
	}
	os.Remove(testEndpointsFilePath)

}

func TestGetClientIDOfEnv(t *testing.T) {
	writeCorrectKeys()

	returnedKey := GetClientIDOfEnv(devName, testKeysFilePath)
	expectedKey := getSampleKeys().Environments[devName].ClientID
	if returnedKey != expectedKey {
		t.Errorf("Expected '%s', got '%s'\n", expectedKey, returnedKey)
	}
	os.Remove(testKeysFilePath)

}

func TestGetClientSecretOfEnv(t *testing.T) {
	writeCorrectKeys()

	returnedKey := GetClientSecretOfEnv(devName, devPassword, testKeysFilePath)
	expectedKey := Decrypt([]byte(GetMD5Hash(devPassword)), getSampleKeys().Environments[devName].ClientSecret)

	if returnedKey != expectedKey {
		t.Errorf("Expected '%s', got '%s'\n", expectedKey, returnedKey)
	}
	os.Remove(testKeysFilePath)
}

func TestGetUsernameOfEnv(t *testing.T) {
	writeCorrectKeys()

	returnedKey := GetUsernameOfEnv(devName, testKeysFilePath)
	expectedKey := getSampleKeys().Environments[devName].Username

	if returnedKey != expectedKey {
		t.Errorf("Expected '%s', got '%s'\n", expectedKey, returnedKey)
	}
	os.Remove(testKeysFilePath)
}

