package params

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// loadAPIFromFile loads API file from the path and returns a slice of bytes or an error
func loadAPIFromFile(path string) ([]byte, error) {
	r, err := os.Open(filepath.FromSlash(path))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	_ = r.Close()

	return data, err
}

func TestLoadApiParamsFromFileInvalidYAML(t *testing.T) {
	conf, err := LoadApiParamsFromFile("testdata/api_params-invalid.yml")
	assert.Error(t, err, "Should return an error for invalid yaml files")
	assert.Nil(t, conf, "Should return nil when errors are returned")
}

func TestLoadApiParamsFromFileWithoutEnv(t *testing.T) {
	conf, err := LoadApiParamsFromFile("testdata/api_params-env.yml")
	assert.Error(t, err, "Should return error when environment variables not present")
	assert.Nil(t, conf, "Conf should be nil")
}

func TestLoadAPIFromFile(t *testing.T) {
	apiData, err := loadAPIFromFile("testdata/api.json")
	assert.Nil(t, err, "Error should be nil when correct json loaded")
	assert.True(t, len(apiData) > 0, "API data should be greater than zero for correct file")
}

func TestExtractAPIEndpointConfig(t *testing.T) {
	apiData, err := loadAPIFromFile("testdata/api.json")
	assert.Nil(t, err, "Error should be nil for correct json loading")
	endpointData, err := ExtractAPIEndpointConfig(apiData)
	assert.Nil(t, err, "Error should be nil for unmarshal json")
	assert.True(t, len(endpointData) > 0, "Correct endpoint data should be loaded")
	assert.True(t, strings.Contains(endpointData, "production_endpoint"), "Should contain correct data")
}

func TestAPIConfig_ContainsEnv(t *testing.T) {
	configData, err := LoadApiParamsFromFile("testdata/api_params.yml")
	assert.Nil(t, err, "Error should be nil for correct yaml loading")

	assert.NotNil(t, configData.GetEnv("dev"), "Should contain correct environment")
	assert.Nil(t, configData.GetEnv("prod"), "Should not contain undefined environment")
}
