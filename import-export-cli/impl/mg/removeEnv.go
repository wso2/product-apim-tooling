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

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// removeEnv
// to be used with 'remove env' command
// @param envName : Name of the environment to be removed from the
// @param mainConfigFilePath : Path to file where env endpoints are stored
// @return error
func RemoveEnv(envName, mainConfigFilePath string) error {
	if envName == "" {
		return errors.New("Name of the environment cannot be blank")
	}
	if utils.MgwAdapterEnvExistsInMainConfigFile(envName, mainConfigFilePath) {
		var err error

		// remove access tokens, if user has already logged into this environment
		store, err := credentials.GetDefaultCredentialStore()
		if store.HasMG(envName) {
			err = RunLogout(envName)
			if err != nil {
				utils.Logln("Unable to log out from Microgateway Adapter in environment: "+
					envName, err)
			}
		}

		// remove env from mainConfig file (endpoints file)
		err = utils.RemoveMgwAdapterEnvFromMainConfigFile(envName, mainConfigFilePath)
		if err != nil {
			return err
		}
	} else {
		// environment does not exist in mainConfig file (endpoints file). Nothing to remove
		return errors.New("Microgateway Adapter '" + envName + "' not found in " + mainConfigFilePath)
	}
	return nil
}
