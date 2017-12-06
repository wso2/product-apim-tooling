/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
)

// EnvExistsInKeysFile
// @param env : Name of the Environment
// @param filePath : Path to file where env keys are stored
// @return bool : true if 'env' exists in the env_keys_all.yaml
// and false otherwise
func EnvExistsInKeysFile(env, filePath string) bool {
	envKeysAll := GetEnvKeysAllFromFile(filePath)
	for _env := range envKeysAll.Environments {
		if _env == env {
			return true
		}
	}
	return false
}

// EnvExistsMainConfigFile
// @param env : Name of the Environment
// @param filePath : Path to file where env endpoints are stored
// @return bool : true if 'env' exists in the main_config.yaml
// and false otherwise
func EnvExistsInMainConfigFile(env, filePath string) bool {
	envEndpointsAll := GetMainConfigFromFile(filePath)
	for _env := range envEndpointsAll.Environments {
		if _env == env {
			return true
		}
	}

	return false
}

// AndNewEnvToKeysFile
// Insert new env entry to keys file (env_keys_all.yaml)
// @param name : Name of the environment
// @param envKeys : EnvKeys object for the environment
// @param filePath : Path to file where env keys are stored
func AddNewEnvToKeysFile(name string, envKeys EnvKeys, filePath string) {
	envKeysAll := GetEnvKeysAllFromFile(filePath)
	Logln(LogPrefixInfo+"EnvKeysAll:", envKeysAll)
	if envKeysAll == nil {
		envKeysAll = new(EnvKeysAll)
	}

	if envKeysAll.Environments == nil {
		envKeysAll.Environments = make(map[string]EnvKeys)
	}
	envKeysAll.Environments[name] = envKeys

	WriteConfigFile(envKeysAll, filePath)
}

// RemoveEnvFromKeysFiles
// used with 'reset-user' command
// does not remove env from endpoints file
// @param env
func RemoveEnvFromKeysFile(env, keysFilePath, mainConfigFilePath string) error {
	/*
	 mainConfigFilePath is passed to check if it exists in endpoints
	 env CANNOT exist only in keys file
	 env CAN exist only in endpoints file (env not initialized i.e. not used with a command)
	*/
	if env == "" {
		return errors.New("environment cannot be blank")
	}
	envKeysAll := GetEnvKeysAllFromFile(keysFilePath)
	if EnvExistsInMainConfigFile(env, mainConfigFilePath) {
		Logln(LogPrefixInfo + "Environment '" + env + "' exists in file " + mainConfigFilePath)
		if EnvExistsInKeysFile(env, keysFilePath) {
			Logln(LogPrefixInfo + "Environment '" + env + "' exists in file " + keysFilePath)
			delete(envKeysAll.Environments, env)
			Logln(LogPrefixInfo + "removing environment '" + env + "' from '" + keysFilePath + "'")
			WriteConfigFile(envKeysAll, keysFilePath)
			return nil
		} else {
			// env doesn't exist in keys file
			return errors.New("environment is not initialized yet. No user data to reset")
		}
	} else {
		// env doesn't exist in endpoints file
		// nothing to remove
		return errors.New("environment not found in " + mainConfigFilePath)
	}
}

// @param env : Environment to be removed from file
// @param endpointsFilePath : Path to file where env endpoints are stored
func RemoveEnvFromMainConfigFile(env, endpointsFilePath string) error {
	if env == "" {
		return errors.New("environment cannot be blank")
	}
	mainConfig := GetMainConfigFromFile(endpointsFilePath)
	if EnvExistsInMainConfigFile(env, endpointsFilePath) {
		Logln(LogPrefixInfo + "Environment '" + env + "' exists in file " + endpointsFilePath)
		delete(mainConfig.Environments, env)
		WriteConfigFile(mainConfig, endpointsFilePath)
		return nil
	} else {
		// env doesn't exist in endpoints file
		return errors.New("environment not found in " + endpointsFilePath)
	}
}

// Get keys of environment 'env' from the file env_keys_all.yaml
// client_secret is not decrypted
// @param env : name of the environment
// @param filePath : Path to file where env keys are stored
// @return *EnvKeys
// @return error
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
	mainConfig := GetMainConfigFromFile(filePath)
	for _env, endpoints := range mainConfig.Environments {
		if _env == env {
			return &endpoints, nil
		}
	}

	return nil, errors.New("error getting endpoints of environment '" + env + "'")
}

// Get APIMEndpoint of a given environment
func GetAPIMEndpointOfEnv(env, filePath string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env, filePath)
	return envEndpoints.ApiManagerEndpoint
}

// Get TokenEndpoint of a given environment
func GetTokenEndpointOfEnv(env, filePath string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env, filePath)
	return envEndpoints.TokenEndpoint
}

// Get RegistrationEndpoint of a given environment
func GetRegistrationEndpointOfEnv(env, filePath string) string {
	envEndpoints, _ := GetEndpointsOfEnvironment(env, filePath)
	return envEndpoints.RegistrationEndpoint
}

// Get username of an environment given the environment
func GetUsernameOfEnv(env, filePath string) string {
	envKeys, _ := GetKeysOfEnvironment(env, filePath)
	return envKeys.Username
}

// Get client_id of an environment given the environment
func GetClientIDOfEnv(env, filePath string) string {
	envKeys, _ := GetKeysOfEnvironment(env, filePath)
	return envKeys.ClientID
}

// Get decrypted client_secret of an environment given the environment and password
// password is needed to decrypt client_secret
// decryption_key = md5(password)
func GetClientSecretOfEnv(env, password, filePath string) string {
	envKeys, _ := GetKeysOfEnvironment(env, filePath)
	decryptedClientSecret := Decrypt([]byte(GetMD5Hash(password)), envKeys.ClientSecret)
	return decryptedClientSecret
}

// check if an environment by the name 'default' exists in the mainConfig file
// input the path to main_config file
func IsDefaultEnvPresent(mainConfigFilePath string) bool {
	mainConfig := GetMainConfigFromFile(mainConfigFilePath)
	for envName := range mainConfig.Environments {
		if envName == DefaultEnvironmentName {
			return true
		}
	}
	return false
}

// return the name of default environment, if it exists
// Currently, the name should be literally 'default'
func GetDefaultEnvironment(mainConfigFilePath string) string {
	if IsDefaultEnvPresent(mainConfigFilePath) {
		return DefaultEnvironmentName
	}
	return ""
}
