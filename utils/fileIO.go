package utils

import (
	"gopkg.in/yaml.v2"
	"errors"
	"io/ioutil"
	"fmt"
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

func WriteEnvKeysToFile() {
	fmt.Println("WriteEnvKeysToFlie() called")
	var envKeysAll EnvKeysAll
	envKeysAll.Environments = make(map[string]EnvKeys)

	username := "admin"
	password := "admin"
	hashedPassword := GetMD5Hash(password)

	// Generate (client_id, client_secret) pairs based on registration endpoints in env_endpoints_all.yaml
	envEndpointsAll := GetEnvEndpointsFromFile()
	for env, endpoints := range envEndpointsAll.Environments {
		clientID, clientSecret := GetClientIDSecret(username, password, endpoints.RegistrationEndpoint)
		clientSecretEncrypted := Encrypt([]byte(hashedPassword), clientSecret)
		envKeysAll.Environments[env] = EnvKeys{clientID,clientSecretEncrypted , username}
	}

	WriteConfigFile(envKeysAll, "env_keys_all.yaml")

}

// Get keys of environment 'env' from the file env_keys_all.yaml
func GetKeysOfEnvironment(env string, password string) *EnvKeys {
	envKeysAll := GetEnvKeysFromFile()
	for _env, keys := range envKeysAll.Environments {
		if _env == env {
			keys.ClientSecret = Decrypt([]byte(password), keys.ClientSecret)
			return &keys
		}
	}

	return &EnvKeys{"", "", ""}
}

// Return true if 'env' exists in the env_keys_all.yaml
// and false otherwise
func EnvExistsInKeysFile(env string) bool {
	envKeysAll := GetEnvKeysFromFile()
	for _env, _ := range envKeysAll.Environments {
		if _env == env {
			return true
		}
	}
	return false
}

// Returns true if 'env' exists in env_endpoints_all.yaml
// and false otherwise
func EnvExistsInEndpointsFile(env string) bool {
	envEndpointsAll := GetEnvEndpointsFromFile()
	for _env, _ := range envEndpointsAll.Environments {
		if _env == env {
			return true
		}
	}

	return false
}

// Read and return EnvKeysAll
func GetEnvKeysFromFile() EnvKeysAll{
	data, err := ioutil.ReadFile("./env_keys_all.yaml")
	if err != nil {
		fmt.Println("Error reading env_keys_all.yaml")
		panic(err)
	}

	var envKeysAll EnvKeysAll
	if err := envKeysAll.ReadEnvKeysFromFile(data); err != nil {
		fmt.Println("Error parsing env_keys_all.yaml")
		panic(err)
	}
	//fmt.Printf("%+v\n", envKeysAll)

	return envKeysAll
}

// Read and return EnvEndpointsAll
func GetEnvEndpointsFromFile() *EnvEndpointsAll {
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