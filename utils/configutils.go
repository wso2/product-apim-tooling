package utils

import (
	"github.com/wso2/wum-client/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type EnvInfo struct {
	Endpoint     string `yaml:"endpoint"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RefreshToken string `yaml:"refresh_token"`
}

type EnvConfig struct {
	Environments map[string]EnvInfo `yaml:"environments"`
}

// variables
var envConfig EnvConfig

// Validates the configuration file
func (envConfig *EnvConfig) validate() {
	//
}

/**
Load the Environments Configuration file from the config.yaml file. If the file is not there
create a new config.yaml file and add default values
Validates the configuration, if it exists
*/
func LoadEnvConfig(envLocalConfig string) /* EnvConfig */ {
}

// Returns a pointer to env configuration
func GetEnvConfig() *EnvConfig {
	if &envConfig == nil {
		utils.HandleErrorAndExit("Env configuration is not available", nil)
	}
	return &envConfig
}

// Persists the given Env configuration
func WriteConfigFile(envConfig interface{}, envConfigFilePath string) {
	data, err := yaml.Marshal(&envConfig)
	if err != nil {
		utils.HandleErrorAndExit("Unable to create Env Configuration.", err)
	}

	err = ioutil.WriteFile(envConfigFilePath, data, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Unable to create Env Configuration.", err)
	}
}

/*
environments:
	dev:
		url: https://example.com
		client_id: eqwrewqr
		client_secret: 192430ijasj90

	staging:
		url: https://example.com/staging
		client_id: a930j
		client_secret: 24342jl
 */
