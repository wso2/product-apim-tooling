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
// password is needed to decrypt client_secret (decryption_key = md5(password))
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

// Return EnvEndpoints for a given environment
func GetEndpointsOfEnvironment(env string) *EnvEndpoints {
	envEndpointsAll := GetEnvEndpointsFromFile()
	for _env, endpoints := range envEndpointsAll.Environments {
		if _env == env {
			return &endpoints
		}
	}

	return &EnvEndpoints{"", "", ""}
}

func GetAPIMEndpointOfEnv(env string) string {
	envEndpoints := GetEndpointsOfEnvironment(env)
	return envEndpoints.APIManagerEndpoint
}

func GetTokenEndpointOfEnv(env string) string {
	envEndpoints := GetEndpointsOfEnvironment(env)
	return envEndpoints.TokenEndpoint
}

func GetRegistrationEndpointOfEnv(env string) string {
	envEndpoints := GetEndpointsOfEnvironment(env)
	return envEndpoints.RegistrationEndpoint
}
