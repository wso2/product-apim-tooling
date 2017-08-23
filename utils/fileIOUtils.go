package utils

import (
	"gopkg.in/yaml.v2"
	"errors"
	"io/ioutil"
	"fmt"
	"os"
)

func WriteConfigFile(c interface{}, envConfigFilePath string) {
	data, err := yaml.Marshal(&c)
	if err != nil {
		HandleErrorAndExit("Unable to create Env Configuration.", err)
	}

	err = ioutil.WriteFile(envConfigFilePath, data, 0644)
	if err != nil {
		HandleErrorAndExit("Unable to create Env Configuration.", err)
	}
}

// Read and return EnvKeysAll
func GetEnvKeysAllFromFile() *EnvKeysAll {
	data, err := ioutil.ReadFile("./env_keys_all.yaml")
	if err != nil {
		fmt.Println("Error reading env_keys_all.yaml")
		_, _ = os.Create("./env_keys_all.yaml")
		data, err = ioutil.ReadFile("./env_keys_all.yaml")
	}

	var envKeysAll EnvKeysAll
	if err := envKeysAll.ReadEnvKeysFromFile(data); err != nil {
		fmt.Println("Error parsing env_keys_all.yaml")
		panic(err)
	}
	//fmt.Printf("%+v\n", envKeysAll)

	return &envKeysAll
}

// Read and return EnvEndpointsAll
func GetEnvEndpointsAllFromFile() *EnvEndpointsAll {
	data, err := ioutil.ReadFile("./env_endpoints_all.yaml")
	if err != nil {
		fmt.Println("Error reading env_endpoints_all.yaml")
		panic(err)
	}

	var envEndpointsAll EnvEndpointsAll
	if err := envEndpointsAll.ReadEnvEndpointsFromFile(data); err != nil {
		fmt.Println("Error parsing env_endpoints_all.yaml")
		panic(err)
	}

	return &envEndpointsAll
}

// Read and validate contents of env_endpoints_all.yaml
// will throw errors if the any of the lines is blank
func (envEndpointsAll *EnvEndpointsAll) ReadEnvEndpointsFromFile(data []byte) error {
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

// Read and validate contents of env_keys_all.yaml
// will throw errors if the any of the lines is blank
func (envKeysAll *EnvKeysAll) ReadEnvKeysFromFile(data []byte) error {
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
	}
	return nil
}
