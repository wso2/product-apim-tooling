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

package impl

import (
	"errors"
	"fmt"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// AddEnv adds a new environment and its endpoints and writes to config file
// @param envName : Name of the Environment
// @param publisherEndpoint : API Manager Endpoint for the environment
// @param regEndpoint : Registration Endpoint for the environment
// @param tokenEndpoint : Token Endpoint for the environment
// @param mainConfigFilePath : Path to file where env endpoints are stored
// @return error
func AddEnv(envName string, envEndpoints *utils.EnvEndpoints, mainConfigFilePath, addEnvCmdLiteral string) error {
	var isDefaultTokenEndpointSet bool = false

	if envName == "" {
		// name of the environment is blank
		return errors.New("Name of the environment cannot be blank")
	}

	if envEndpoints.TokenEndpoint == "" {
		// If token endpoint string is empty,then assign the default value
		if envEndpoints.ApiManagerEndpoint != "" && !isDefaultTokenEndpointSet {
			isDefaultTokenEndpointSet = true
			envEndpoints.TokenEndpoint = utils.GetTokenEndPointFromAPIMEndpoint(envEndpoints.ApiManagerEndpoint)
		}
		if envEndpoints.PublisherEndpoint != "" && !isDefaultTokenEndpointSet {
			envEndpoints.TokenEndpoint = utils.GetTokenEndPointFromPublisherEndpoint(envEndpoints.PublisherEndpoint)
		}
		fmt.Printf("Default token endpoint '%s' is added as the token endpoint \n", envEndpoints.TokenEndpoint)
	}

	if envEndpoints.ApiManagerEndpoint == "" {
		if envEndpoints.AdminEndpoint == "" || envEndpoints.DevPortalEndpoint == "" ||
			envEndpoints.PublisherEndpoint == "" || envEndpoints.RegistrationEndpoint == "" ||
			envEndpoints.TokenEndpoint == "" {
			utils.ShowHelpCommandTip(addEnvCmdLiteral)
			return errors.New("Endpoint(s) cannot be blank")
		}
	}

	if utils.EnvExistsInMainConfigFile(envName, mainConfigFilePath) {
		// environment already exists
		return errors.New("Environment '" + envName + "' already exists in " + mainConfigFilePath)
	}

	mainConfig := utils.GetMainConfigFromFile(mainConfigFilePath)

	var validatedEnvEndpoints = utils.EnvEndpoints{
		TokenEndpoint: envEndpoints.TokenEndpoint,
	}

	if envEndpoints.ApiManagerEndpoint != "" {
		validatedEnvEndpoints.ApiManagerEndpoint = envEndpoints.ApiManagerEndpoint
	}

	if envEndpoints.RegistrationEndpoint != "" {
		validatedEnvEndpoints.RegistrationEndpoint = envEndpoints.RegistrationEndpoint
	}

	if envEndpoints.PublisherEndpoint != "" {
		validatedEnvEndpoints.PublisherEndpoint = envEndpoints.PublisherEndpoint
	}

	if envEndpoints.DevPortalEndpoint != "" {
		validatedEnvEndpoints.DevPortalEndpoint = envEndpoints.DevPortalEndpoint
	}

	if envEndpoints.AdminEndpoint != "" {
		validatedEnvEndpoints.AdminEndpoint = envEndpoints.AdminEndpoint
	}

	mainConfig.Environments[envName] = validatedEnvEndpoints
	utils.WriteConfigFile(mainConfig, mainConfigFilePath)

	fmt.Printf("Successfully added environment '%s'\n", envName)

	return nil
}
