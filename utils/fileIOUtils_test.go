package utils

import (
	"testing"
	"os"
	"io/ioutil"
)

var testKeysFilePath string = ApplicationRoot + "/test_keys_config.yaml"
var testEndpointsFilePath string = ApplicationRoot + "/test_endpoints_config.yaml"
var envKeysAll *EnvKeysAll = new(EnvKeysAll)
var envEndpointsAll *EnvEndpointsAll = new(EnvEndpointsAll)

// helper function for testing
func writeCorrectKeys(){
	envKeysAll.Environments  = make(map[string]EnvKeys)
	envKeysAll.Environments["dev"] = EnvKeys{"dev-client-id", "dev-client-secret", "dev-username"}
	envKeysAll.Environments["qa"] = EnvKeys{"qa-client-id", "qa-client-secret", "qa-username"}
	WriteConfigFile(envKeysAll, testKeysFilePath)
}

// helper function for testing
func writeCorrectEndpoints() {
	envEndpointsAll.Environments = make(map[string]EnvEndpoints)
	envEndpointsAll.Environments["dev"] = EnvEndpoints{"dev-apim-endpoint", "dev-reg-endpoint", "dev-token-endpoint"}
	envEndpointsAll.Environments["qa"] = EnvEndpoints{"qa-apim-endpoint", "qa-reg-endpoint", "dev-token-endpoint"}
	WriteConfigFile(envEndpointsAll, testEndpointsFilePath)
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

	if envKeysAllReturned.Environments["dev"] != envKeysAll.Environments["dev"] ||
		 envKeysAllReturned.Environments["qa"] != envKeysAll.Environments["qa"] {
		t.Errorf("Error in GetEnvKeysAllFromFile()")
	}

	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

func TestGetEnvKeysAllFromFile2(t *testing.T) {
	// testing for incorrect data
	var envIncorrectKeysAll *EnvKeysAll = new(EnvKeysAll)
	envIncorrectKeysAll.Environments  = make(map[string]EnvKeys)
	envIncorrectKeysAll.Environments["dev"] = EnvKeys{"", "", "dev-username"}

	WriteConfigFile(envIncorrectKeysAll, testKeysFilePath)
	envKeysAllReturned := GetEnvKeysAllFromFile(testKeysFilePath)

	if envKeysAllReturned.Environments["dev"] != envIncorrectKeysAll.Environments["dev"] ||
		envKeysAllReturned.Environments["qa"] != envIncorrectKeysAll.Environments["qa"] {
		t.Errorf("Error in GetEnvKeysAllFromFile()")
	}

	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

func TestGetEnvKeysAllFromFile3(t *testing.T) {
	// testing for incorrect data
	var envIncorrectKeysAll *EnvKeysAll = new(EnvKeysAll)
	envIncorrectKeysAll.Environments  = make(map[string]EnvKeys)
	envIncorrectKeysAll.Environments["dev"] = EnvKeys{"dev-client-id", "", "dev-username"}

	WriteConfigFile(envIncorrectKeysAll, testKeysFilePath)
	envKeysAllReturned := GetEnvKeysAllFromFile(testKeysFilePath)

	if envKeysAllReturned.Environments["dev"] != envIncorrectKeysAll.Environments["dev"] ||
		envKeysAllReturned.Environments["qa"] != envIncorrectKeysAll.Environments["qa"] {
		t.Errorf("Error in GetEnvKeysAllFromFile()")
	}

	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

func TestGetEnvEndpointsAllFromFile(t *testing.T) {
	// testing for correct data
	writeCorrectEndpoints()
	envEndpointsAllReturned := GetEnvEndpointsAllFromFile(testEndpointsFilePath)

	if envEndpointsAllReturned.Environments["dev"] != envEndpointsAllReturned.Environments["dev"] ||
		envEndpointsAllReturned.Environments["qa"] != envEndpointsAllReturned.Environments["qa"] {
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

