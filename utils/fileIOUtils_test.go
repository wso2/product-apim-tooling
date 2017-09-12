// +integration
package utils

import (
	"testing"
	"os"
	"io/ioutil"
)

const testKeysFilePath string = ApplicationRoot + "/test_keys_config.yaml"
const testEndpointsFilePath string = ApplicationRoot + "/test_endpoints_config.yaml"
var envKeysAll *EnvKeysAll = new(EnvKeysAll)
var envEndpointsAll *EnvEndpointsAll = new(EnvEndpointsAll)
const devName string = "dev"
const qaName string = "qa"
const devUsername string = "dev-username"
const qaUsername string = "qa-username"
const devPassword string = "dev-password"
const qaPassword string = "qa-password"


// helper function for testing
func writeCorrectKeys(){
	initSampleKeys()
	WriteConfigFile(envKeysAll, testKeysFilePath)
}

func getSampleKeys() *EnvKeysAll {
	initSampleKeys()
	return envKeysAll
}

func getSampleEndpoints() *EnvEndpointsAll {
	return envEndpointsAll
}

func initSampleKeys() {
	envKeysAll.Environments = make(map[string]EnvKeys)
	devEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(devPassword)), "dev-client-secret")
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa-client-secret")
	envKeysAll.Environments[devName] = EnvKeys{"dev-client-id", devEncryptedClientSecret, devUsername}
	envKeysAll.Environments[qaName] = EnvKeys{"qa-client-id", qaEncryptedClientSecret, qaUsername}
}

// helper function for testing
func writeCorrectEndpoints() {
	initSampleEndpoints()
	WriteConfigFile(envEndpointsAll, testEndpointsFilePath)
}
func initSampleEndpoints() {
	envEndpointsAll.Environments = make(map[string]EnvEndpoints)
	envEndpointsAll.Environments[devName] = EnvEndpoints{"dev-apim-endpoint",
		"dev-reg-endpoint", "dev-token-endpoint"}
	envEndpointsAll.Environments[qaName] = EnvEndpoints{"qa-apim-endpoint",
		"qa-reg-endpoint", "dev-token-endpoint"}
}

func TestWriteConfigFile(t *testing.T) {
	writeCorrectKeys()
	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

func TestGetEnvKeysAllFromFile1(t *testing.T) {
	writeCorrectKeys()
	envKeysAllReturned := GetEnvKeysAllFromFile(testKeysFilePath)

	if envKeysAllReturned.Environments[devName] != envKeysAll.Environments[devName] ||
		 envKeysAllReturned.Environments[qaName] != envKeysAll.Environments[qaName] {
		t.Errorf("Error in GetEnvKeysAllFromFile()")
	}

	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

/*
func TestGetEnvKeysAllFromFile3(t *testing.T) {
	// testing for incorrect data
	var envIncorrectKeysAll *EnvKeysAll = new(EnvKeysAll)
	envIncorrectKeysAll.Environments  = make(map[string]EnvKeys)
	envIncorrectKeysAll.Environments[devName] = EnvKeys{"dev-client-id", "", devUsername}

	WriteConfigFile(envIncorrectKeysAll, testKeysFilePath)
	envKeysAllReturned := GetEnvKeysAllFromFile(testKeysFilePath)

	if envKeysAllReturned.Environments[devName] != envIncorrectKeysAll.Environments[devName] ||
		envKeysAllReturned.Environments[qaName] != envIncorrectKeysAll.Environments[qaName] {
		t.Errorf("Error in GetEnvKeysAllFromFile()")
	}

	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}
*/

func TestGetEnvEndpointsAllFromFile(t *testing.T) {
	// testing for correct data
	writeCorrectEndpoints()
	envEndpointsAllReturned := GetEnvEndpointsAllFromFile(testEndpointsFilePath)

	if envEndpointsAllReturned.Environments[devName] != envEndpointsAllReturned.Environments[devName] ||
		envEndpointsAllReturned.Environments[qaName] != envEndpointsAllReturned.Environments[qaName] {
		t.Errorf("Error in GetEnvEndpointsAllFromFile()")
	}

	var err = os.Remove(testEndpointsFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

func TestEnvEndpointsAll_ParseEnvEndpointsFromFile(t *testing.T) {
	var envEndpointsAllLocal EnvEndpointsAll
	writeCorrectEndpoints()
	data, err := ioutil.ReadFile(testEndpointsFilePath)
	if err != nil {
		t.Error("Error")
	}
	envEndpointsAllLocal.ParseEnvEndpointsFromFile(data)

	var err1 = os.Remove(testEndpointsFilePath)
	if err1 != nil {
		t.Errorf("Error deleting file " + testEndpointsFilePath)
	}
}

func TestEnvKeysAll_ParseEnvKeysFromFile(t *testing.T) {
	var envKeysAllLocal EnvKeysAll
	writeCorrectKeys()
	data, err := ioutil.ReadFile(testKeysFilePath)
	if err != nil {
		t.Error("Error")
	}
	envKeysAllLocal.ParseEnvKeysFromFile(data)

	var err1 = os.Remove(testKeysFilePath)
	if err1 != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

