package utils

import "errors"

// Return true if 'env' exists in the env_keys_all.yaml
// and false otherwise
func EnvExistsInKeysFile(env string) bool {
	envKeysAll := GetEnvKeysAllFromFile()
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
	envEndpointsAll := GetEnvEndpointsAllFromFile()
	for _env, _ := range envEndpointsAll.Environments {
		if _env == env {
			return true
		}
	}

	return false
}

// Insert new env entry to env_keys_all.yaml
func AddNewEnvToKeysFile(name string, envKeys EnvKeys) {
	envKeysAll := GetEnvKeysAllFromFile()
	envKeysAll.Environments[name] = envKeys

	WriteConfigFile(envKeysAll, "./env_keys_all.yaml")
}

func RemoveEnvFromKeysFile(env string) (error) {
	if env == "" {
		return errors.New("Environment cannot be blank")
	}
	envKeysAll := GetEnvKeysAllFromFile()
	if EnvExistsInEndpointsFile(env) {
		if EnvExistsInKeysFile(env) {
			delete(envKeysAll.Environments, env)
			WriteConfigFile(envKeysAll, "./env_keys_all.yaml")
			return nil
		} else {
			return errors.New("Environment is not initialized yet. No user data to reset")
		}
	} else {
		return errors.New("Environment not found in env_endpoints_all.yaml")
	}

}

// Get keys of environment 'env' from the file env_keys_all.yaml
// client_secret is not decrypted
func GetKeysOfEnvironment(env string) (*EnvKeys, error) {
	envKeysAll := GetEnvKeysAllFromFile()
	for _env, keys := range envKeysAll.Environments {
		if _env == env {
			return &keys, nil
		}
	}

	return nil, errors.New("Error getting keys of environment '" + env + "'")
}

// Return EnvEndpoints for a given environment
func GetEndpointsOfEnvironment(env string) (*EnvEndpoints, error) {
	envEndpointsAll := GetEnvEndpointsAllFromFile()
	for _env, endpoints := range envEndpointsAll.Environments {
		if _env == env {
			return &endpoints, nil
		}
	}

	return nil, errors.New("Error getting endpoints of environment '" + env + "'")
}

// Get APIMEndpoint of a given environment
func GetAPIMEndpointOfEnv(env string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env)
	return envEndpoints.APIManagerEndpoint
}

// Get TokenEndpoint of a given environment
func GetTokenEndpointOfEnv(env string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env)
	return envEndpoints.TokenEndpoint
}

// Get RegistrationEndpoint of a given environment
func GetRegistrationEndpointOfEnv(env string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env)
	return envEndpoints.RegistrationEndpoint
}

// Get username of an environment given the environment
func GetUsernameOfEnv(env string) string {
	envKeys, _ := GetKeysOfEnvironment(env)
	return envKeys.Username
}

// Get client_id of an environment given the environment
func GetClientIDOfEnv(env string) string {
	envKeys, _ := GetKeysOfEnvironment(env)
	return envKeys.ClientID
}

// Get decrypted client_secret of an environment given the environment and password
// password is needed to decrypt client_secret
// decryption_key = md5(password)
func GetClientSecretOfEnv(env string, password string) string {
	envKeys, _ := GetKeysOfEnvironment(env)
	decryptedClientSecret := Decrypt([]byte(GetMD5Hash(password)), envKeys.ClientSecret)
	return decryptedClientSecret
}
