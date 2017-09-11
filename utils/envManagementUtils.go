package utils

import (
	"errors"
	"fmt"
)

// Return true if 'env' exists in the env_keys_all.yaml
// and false otherwise
func EnvExistsInKeysFile(env string, filePath string) bool {
	envKeysAll := GetEnvKeysAllFromFile(filePath)
	for _env, _ := range envKeysAll.Environments {
		if _env == env {
			return true
		}
	}
	return false
}

// Returns true if 'env' exists in env_endpoints_all.yaml
// and false otherwise
func EnvExistsInEndpointsFile(env string, filePath string) bool {
	envEndpointsAll := GetEnvEndpointsAllFromFile(filePath)
	for _env, _ := range envEndpointsAll.Environments {
		if _env == env {
			return true
		}
	}

	return false
}

// Insert new env entry to env_keys_all.yaml
func AddNewEnvToKeysFile(name string, envKeys EnvKeys, filePath string) {
	envKeysAll := GetEnvKeysAllFromFile(filePath)
	fmt.Println("EnvKeysAll:", envKeysAll)
	if envKeysAll == nil {
		fmt.Println("envKeysAll is nil")
		envKeysAll = new(EnvKeysAll)
	}

	if envKeysAll.Environments == nil {
		fmt.Println("envKeysAll.Environments is nil")
		envKeysAll.Environments = make(map[string]EnvKeys)
	}
	envKeysAll.Environments[name] = envKeys

	WriteConfigFile(envKeysAll, filePath)
}

func RemoveEnvFromKeysFile(env string, filePath string) (error) {
	if env == "" {
		return errors.New("environment cannot be blank")
	}
	envKeysAll := GetEnvKeysAllFromFile(filePath)
	if EnvExistsInEndpointsFile(env, filePath) {
		if EnvExistsInKeysFile(env, filePath) {
			delete(envKeysAll.Environments, env)
			WriteConfigFile(envKeysAll, filePath)
			return nil
		} else {
			return errors.New("environment is not initialized yet. No user data to reset")
		}
	} else {
		return errors.New("environment not found in " + filePath)
	}

}

// Get keys of environment 'env' from the file env_keys_all.yaml
// client_secret is not decrypted
func GetKeysOfEnvironment(env string, filePath string) (*EnvKeys, error) {
	envKeysAll := GetEnvKeysAllFromFile(filePath)
	for _env, keys := range envKeysAll.Environments {
		if _env == env {
			return &keys, nil
		}
	}

	return nil, errors.New("error getting keys of environment '" + env + "'")
}

// Return EnvEndpoints for a given environment
func GetEndpointsOfEnvironment(env string, filePath string) (*EnvEndpoints, error) {
	envEndpointsAll := GetEnvEndpointsAllFromFile(filePath)
	for _env, endpoints := range envEndpointsAll.Environments {
		if _env == env {
			return &endpoints, nil
		}
	}

	return nil, errors.New("error getting endpoints of environment '" + env + "'")
}

// Get APIMEndpoint of a given environment
func GetAPIMEndpointOfEnv(env string, filePath string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env, filePath)
	return envEndpoints.APIManagerEndpoint
}

// Get TokenEndpoint of a given environment
func GetTokenEndpointOfEnv(env string, filePath string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env, filePath)
	return envEndpoints.TokenEndpoint
}

// Get RegistrationEndpoint of a given environment
func GetRegistrationEndpointOfEnv(env string, filePath string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env, filePath)
	return envEndpoints.RegistrationEndpoint
}

// Get username of an environment given the environment
func GetUsernameOfEnv(env string, filePath string) string {
	envKeys, _ := GetKeysOfEnvironment(env, filePath)
	return envKeys.Username
}

// Get client_id of an environment given the environment
func GetClientIDOfEnv(env string, filePath string) string {
	envKeys, _ := GetKeysOfEnvironment(env, filePath)
	return envKeys.ClientID
}

// Get decrypted client_secret of an environment given the environment and password
// password is needed to decrypt client_secret
// decryption_key = md5(password)
func GetClientSecretOfEnv(env string, password string, filePath string) string {
	envKeys, _ := GetKeysOfEnvironment(env, filePath)
	decryptedClientSecret := Decrypt([]byte(GetMD5Hash(password)), envKeys.ClientSecret)
	return decryptedClientSecret
}
