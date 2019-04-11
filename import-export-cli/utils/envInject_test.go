package utils

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/Jeffail/gabs"

	"github.com/stretchr/testify/assert"
)

func TestInjectEnvShouldFailWhenEnvNotPresent(t *testing.T) {
	data := `$MYVAR`
	str, err := injectEnv(data)
	assert.Equal(t, "", str, "Should return empty string")
	assert.Error(t, err, "Should return an error")
	assert.EqualError(t, err, "$MYVAR is required, please set the environment variable")
}

func TestInjectEnvShouldPassWhenEnvPresents(t *testing.T) {
	data := `$MYVAR`
	_ = os.Setenv("MYVAR", "myval")
	str, err := injectEnv(data)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, "myval", str, "Should correctly replace environment variable")
}

func TestLoadConfigFromFileValidYAML(t *testing.T) {
	conf, err := LoadConfigFromFile("testdata/.apim-vars.yml")
	assert.Nil(t, err, "Should return nil for correctly parsed files")
	assert.Equal(t, 2, len(conf.Environments), "Should return two environments")
	assert.Equal(t, "dev", conf.Environments[0].Name, "Should have correct name for environment")
	assert.Equal(t, "test", conf.Environments[1].Name, "Should have correct name for environment")
	assert.Equal(t, 2, *conf.Environments[0].Endpoints.Production.Config.Factor, "Should return "+
		"correct values for factor")
	assert.Nil(t, conf.Environments[0].Endpoints.Sandbox, "Should return nil for ignored fields on yaml")
}

func TestLoadConfigFromFileInvalidYAML(t *testing.T) {
	conf, err := LoadConfigFromFile("testdata/.apim-vars-invalid.yml")
	assert.Error(t, err, "Should return an error for invalid yaml files")
	assert.Nil(t, conf, "Should return nil when errors are returned")
}

func TestLoadConfigFromFileWithoutEnv(t *testing.T) {
	conf, err := LoadConfigFromFile("testdata/.apim-vars-env.yml")
	assert.Error(t, err, "Should return error when environment variables not present")
	assert.Nil(t, conf, "Conf should be nil")
}

func TestLoadConfigWithEnv(t *testing.T) {
	_ = os.Setenv("FOO_DEV_RETRY", "10")
	_ = os.Setenv("FOO_SANDBOX", "http://127.0.0.1")
	conf, err := LoadConfigFromFile("testdata/.apim-vars-env.yml")
	assert.Nil(t, err, "Should return empty error on correct reading")
	assert.Equal(t, 10, *conf.Environments[0].Endpoints.Production.Config.RetryTimeOut)
	assert.Equal(t, "http://127.0.0.1", *conf.Environments[1].Endpoints.Sandbox.Url)
}

func TestLoadAPIFromFile(t *testing.T) {
	apiData, err := LoadAPIFromFile("testdata/api.json")
	assert.Nil(t, err, "Error should be nil when correct json loaded")
	assert.True(t, len(apiData) > 0, "API data should be greater than zero for correct file")
}

func TestExtractAPIEndpointConfig(t *testing.T) {
	apiData, err := LoadAPIFromFile("testdata/api.json")
	assert.Nil(t, err, "Error should be nil for correct json loading")
	endpointData, err := ExtractAPIEndpointConfig(apiData)
	assert.Nil(t, err, "Error should be nil for unmarshal json")
	assert.True(t, len(endpointData) > 0, "Correct endpoint data should be loaded")
	assert.True(t, strings.Contains(endpointData, "production_endpoint"), "Should contain correct data")
}

func TestMergeAPIConfig(t *testing.T) {
	apiData, err := LoadAPIFromFile("testdata/api.json")
	assert.Nil(t, err, "Error should be nil for correct json loading")
	endpointData, err := ExtractAPIEndpointConfig(apiData)
	assert.Nil(t, err, "Error should be nil for correct json extraction")
	configData, err := LoadConfigFromFile("testdata/.apim-vars.yml")
	assert.Nil(t, err, "Error should be nil for correct yaml loading")
	config, err := json.Marshal(configData.Environments[0].Endpoints)

	merged, err := MergeAPIConfig([]byte(endpointData), config)
	assert.Nil(t, err, "Error should be nil for successful merging")

	jsonObj, err := gabs.ParseJSON([]byte(merged))
	assert.Nil(t, err, "Merged should be valid json")

	suspendDuration, ok := jsonObj.Path("production_endpoints.config.suspendDuration").Data().(string)
	assert.True(t, ok, "Should return correct type for unchanged fields")
	assert.Equal(t, "40", suspendDuration, "Should return correct value for unchanged fields")

	retryTimeOut, ok := jsonObj.Path("production_endpoints.config.retryTimeOut").Data().(string)
	assert.True(t, ok, "Should return correct type for changed fields")
	assert.Equal(t, "60", retryTimeOut, "Should return correct value for changed fields")
}

func TestAPIConfig_ContainsEnv(t *testing.T) {
	configData, err := LoadConfigFromFile("testdata/.apim-vars.yml")
	assert.Nil(t, err, "Error should be nil for correct yaml loading")

	assert.True(t, configData.ContainsEnv("dev"), "Should contain correct environment")
	assert.False(t, configData.ContainsEnv("prod"), "Should not contain undefined environment")
}
