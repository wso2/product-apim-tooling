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

package mg

import (
	"errors"
	"fmt"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// AddEnv adds a new environment and its endpoints and writes to config file
func AddEnv(envName string, mgwEndpoints *utils.MgwEndpoints) error {
	if envName == "" {
		// name of the environment is blank
		return errors.New("Name of the environment cannot be blank")
	}

	mainConfigFilePath := utils.MainConfigFilePath
	if utils.MgwAdapterEnvExistsInMainConfigFile(envName, mainConfigFilePath) {
		// environment already exists
		return errors.New("MgwAdapter Environment '" + envName + "' already exists in " + mainConfigFilePath)
	}

	mainConfig := utils.GetMainConfigFromFile(mainConfigFilePath)

	var validatedMgwEndpoints = utils.MgwEndpoints{}
	if mgwEndpoints.AdapterEndpoint == "" {
		return errors.New("Adapter url cannot be blank")
	} else {
		validatedMgwEndpoints.AdapterEndpoint = mgwEndpoints.AdapterEndpoint
	}

	mainConfig.MgwAdapterEnvs[envName] = validatedMgwEndpoints
	utils.WriteConfigFile(mainConfig, mainConfigFilePath)

	fmt.Printf("Successfully added environment '%s'\n", envName)

	return nil
}
