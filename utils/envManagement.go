package utils

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

// Get keys of environment 'env' from the file env_keys_all.yaml
// client_secret is not decrypted
func GetKeysOfEnvironment(env string) *EnvKeys {
	envKeysAll := GetEnvKeysFromFile()
	for _env, keys := range envKeysAll.Environments {
		if _env == env {
			return &keys
		}
	}

	// TODO: Throw error instead of returning an empty object
	return &EnvKeys{"", "", ""}
}


// Return EnvEndpoints for a given environment
func GetEndpointsOfEnvironment(env string) *EnvEndpoints {
	envEndpointsAll := GetEnvEndpointsFromFile()
	for _env, endpoints := range envEndpointsAll.Environments {
		if _env == env {
			return &endpoints
		}
	}

	// TODO: Throw error instead of returning an empty object
	return &EnvEndpoints{"", "", ""}
}

// Get APIMEndpoint of a given environment
func GetAPIMEndpointOfEnv(env string) string {
	envEndpoints := GetEndpointsOfEnvironment(env)
	return envEndpoints.APIManagerEndpoint
}

// Get TokenEndpoint of a given environment
func GetTokenEndpointOfEnv(env string) string {
	envEndpoints := GetEndpointsOfEnvironment(env)
	return envEndpoints.TokenEndpoint
}

// Get RegistrationEndpoint of a given environment
func GetRegistrationEndpointOfEnv(env string) string {
	envEndpoints := GetEndpointsOfEnvironment(env)
	return envEndpoints.RegistrationEndpoint
}

// Get username of an environment given the environment
func GetUsernameOfEnv(env string) string {
	envKeys := GetKeysOfEnvironment(env)
	return envKeys.Username
}

// Get client_id of an environment given the environment
func GetClientIDOfEnv(env string) string {
	envKeys := GetKeysOfEnvironment(env)
	return envKeys.ClientID
}

// Get decrypted client_secret of an environment given the environment and password
// password is needed to decrypt client_secret
// decryption_key = md5(password)
func GetClientSecretOfEnv(env string, password string) string {
	envKeys := GetKeysOfEnvironment(env)
	decryptedClientSecret := Decrypt([]byte(GetMD5Hash(password)), envKeys.ClientSecret)
	return decryptedClientSecret
}
