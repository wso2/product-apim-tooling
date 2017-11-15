package cmd

import (
	"testing"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"path/filepath"
	"os"
)

// TestAddEnv1 - Blank Env Name
func TestAddEnv1(t *testing.T) {
	envName := ""
	publisherEndpoint := "test-publisher-endpoint"
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"
	mainConfigFilePath := ""

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, mainConfigFilePath)
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}

}

// TestAddEnv2 - Blank Publisher Endpoint
func TestAddEnv2(t *testing.T) {
	envName := "test-env"
	publisherEndpoint := ""
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"
	mainConfigFilePath := ""

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, mainConfigFilePath)
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}

}

// TestAddEnv3 - Already existing environment
func TestAddEnv3(t *testing.T) {

	sampleMainConfigFileName := "sample_main_config.yaml"
	sampleMainConfigFilePath := filepath.Join(utils.ApplicationRoot, sampleMainConfigFileName)

	var sampleMainConnfig = new(utils.MainConfig)
	sampleMainConnfig.Config = utils.Config{10000, ""}
	sampleMainConnfig.Environments = make(map[string]utils.EnvEndpoints)
	sampleMainConnfig.Environments["dev"] = utils.EnvEndpoints{"sample-publisher-endpoint",
				"sample-reg-endpoint", "sample-token-endpoint"}
	utils.WriteConfigFile(sampleMainConnfig, sampleMainConfigFilePath)

	envName := "dev"
	publisherEndpoint := "sample-publisher-endpoint"
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, sampleMainConfigFilePath)
	if err == nil {
		t.Errorf("Expected error, got nil instead\n")
	}


	defer os.Remove(sampleMainConfigFilePath)
}

// TetsAddEnv4 - Correct Details - Successfully add new environment
func TestAddEnv4(t *testing.T) {
	sampleMainConfigFileName := "sample_main_config.yaml"
	sampleMainConfigFilePath := filepath.Join(utils.ApplicationRoot, sampleMainConfigFileName)

	var sampleMainConnfig = new(utils.MainConfig)
	sampleMainConnfig.Config = utils.Config{10000, ""}
	sampleMainConnfig.Environments = make(map[string]utils.EnvEndpoints)
	sampleMainConnfig.Environments["dev"] = utils.EnvEndpoints{"sample-publisher-endpoint",
		"sample-reg-endpoint", "sample-token-endpoint"}
	utils.WriteConfigFile(sampleMainConnfig, sampleMainConfigFilePath)

	envName := "new-env"
	publisherEndpoint := "sample-publisher-endpoint"
	tokenEndpoint := "test-token-endpoint"
	regEndpint := "test-reg-endpoint"

	err := addEnv(envName, publisherEndpoint, regEndpint, tokenEndpoint, sampleMainConfigFilePath)
	if err != nil {
		t.Errorf("Expected nil, got '%s' instead\n", err.Error())
	}


	defer os.Remove(sampleMainConfigFilePath)
}

