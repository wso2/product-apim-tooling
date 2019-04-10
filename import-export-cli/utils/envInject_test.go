package utils

import (
	"os"
	"testing"

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
