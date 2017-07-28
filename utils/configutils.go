package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"errors"
	"debug/dwarf"
)

// ------------------- Structs for YAML Config Files ----------------------------------

// For env_keys_all.yaml
// Not to be manually edited
type EnvKeysAll struct {
	Environments map[string]EnvKeys `yaml:"environments"`
}

// For env_endpoints_all.yaml
// To be manually edited by the user
type EnvEndpointsAll struct {
	Environments map[string]EnvEndpoints `yaml:"environments"`
}

type EnvKeys struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"` // to be encrypted (with the user's password) and stored
}

type EnvEndpoints struct {
	APIManagerEndpoint   string `yaml:"api_manager_endpoint"`
	RegistrationEndpoint string `yaml:"registration_endpoint"`
	TokenEndpoint        string `yaml:"token_endpoint"`
}

// ---------------- End of Structs for YAML Config Files ---------------------------------

// variables
var envEndpointsAll EnvEndpointsAll
var envKeysAll EnvKeysAll

// Validates the configuration file
func (envEndpointsAll *EnvEndpointsAll) validate() {
	//
}

// Read contents of env_endpoints_all.yaml
func (envEndpointsAll *EnvEndpointsAll) ReadFromFile(data []byte) error {
	if err := yaml.Unmarshal(data, envEndpointsAll); err != nil {
		return err
	}
	for name, endpoints := range envEndpointsAll.Environments {
		if endpoints.APIManagerEndpoint == "" {
			return errors.New("Invalid API Manager Endpoint for " + name)
		}
		if endpoints.RegistrationEndpoint == "" {
			return errors.New("Invalid Registration Endpoint for " + name)
		}
		if endpoints.TokenEndpoint == "" {
			return errors.New("Invalid Token Endpoint for " + name)
		}
	}
	return nil
}

// Read contents of env_keys_all.yaml
func (envKeysAll *EnvKeysAll) ReadFromFile(data []byte) error {
	if err := yaml.Unmarshal(data, envKeysAll); err != nil {
		return err
	}
	for name, keys := range envKeysAll.Environments {
		if keys.ClientID == "" {
			return errors.New("Invalid ClientID for " + name)
		}
		if keys.ClientSecret == "" {
			return errors.New("Invalid ClientSecret for " + name)
		}
		return nil
	}
}

/**
Load the Environments Configuration file from the config.yaml file. If the file is not there
create a new config.yaml file and add default values
Validates the configuration, if it exists
*/
func LoadEnvConfig(envLocalConfig string) /* EnvEndpointsAll */ {
}

// Returns a pointer to EnvEndpointsAll
func GetEnvEndpointsAll() *EnvEndpointsAll {
	if &envEndpointsAll == nil {
		HandleErrorAndExit("Env configuration is not available", nil)
	}
	return &envEndpointsAll
}

// Returns a pointer to EnvKeysAll
func GetEnvKeysAll() *EnvKeysAll {
	if &envKeysAll == nil {
		HandleErrorAndExit("EnvKeys configuration is not available", nil)
	}
	return &envKeysAll
}

// Persists the given Env configuration
func WriteConfigFile(envConfig interface{}, envConfigFilePath string) {
	data, err := yaml.Marshal(&envConfig)
	if err != nil {
		HandleErrorAndExit("Unable to create Env Configuration.", err)
	}

	err = ioutil.WriteFile(envConfigFilePath, data, 0644)
	if err != nil {
		HandleErrorAndExit("Unable to create Env Configuration.", err)
	}
}

/*
env_keys_config.yaml (Programmatically edited)
===============
environments:
	dev:
		client_id: xxxxxxxxxx
		client_secret: xxxxxxxxxx
		refresh_token: xxxxxxxxxx

	staging:
		client_id: xxxxxxxxxx
		client_secret: xxxxxxxxxx
		refresh_token: xxxxxxxxxx
 */

/*
env_config.yaml (Manually edited)
===============
environments:
	dev:
		apim_endpoint: xxxxxxxxx
		registration_endpoint: xxxxxxxxxx
		token_endpoint: xxxxxxxxx

	staging:
		apim_endpoint: xxxxxxxxx
		registration_endpoint: xxxxxxxxxx
		token_endpoint: xxxxxxxxx
*/
