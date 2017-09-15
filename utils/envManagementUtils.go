/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/

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

func RemoveEnvFromKeysFile(env string, keysFilePath string, endpointsFilePath string) (error) {
	fmt.Println("RemoveEnvFromKeysFile(): KeysFilePath:", keysFilePath)
	fmt.Println("RemoveEnvFromKeysFile(): EndpointsFilePath:", endpointsFilePath)
	if env == "" {
		return errors.New("environment cannot be blank")
	}
	envKeysAll := GetEnvKeysAllFromFile(keysFilePath)
	if EnvExistsInEndpointsFile(env, endpointsFilePath) {
		if EnvExistsInKeysFile(env, keysFilePath) {
			delete(envKeysAll.Environments, env)
			WriteConfigFile(envKeysAll, keysFilePath)
			return nil
		} else {
			// env doesn't exist in keys file
			return errors.New("environment is not initialized yet. No user data to reset")
		}
	} else {
		// env doesn't exist in endpoints file
		return errors.New("environment not found in " + endpointsFilePath)
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
