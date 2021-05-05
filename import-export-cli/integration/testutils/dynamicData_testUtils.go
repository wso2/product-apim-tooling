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

package testutils

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
)

func SetDynamicDataForAPI(t *testing.T, client *apim.Client) {

	os.Setenv("DEV_ENV_PROD_URL", "https://localhost:"+strconv.Itoa(9443+client.GetPortOffset())+
		"/am/sample/pizzashack/v1/api/")
	os.Setenv("DEV_ENV_SAND_URL", "https://localhost:"+strconv.Itoa(9443+client.GetPortOffset())+
		"/am/sample/pizzashack/v1/api/")
	os.Setenv("DEV_ENV_PROD_RE_DELAY", "10")
	os.Setenv("DEV_ENV_PROD_RE_TO", "5")

	t.Cleanup(func() {
		os.Unsetenv("DEV_ENV_PROD_URL")
		os.Unsetenv("DEV_ENV_SAND_URL")
		os.Unsetenv("DEV_ENV_PROD_RE_DELAY")
		os.Unsetenv("DEV_ENV_PROD_RE_TO")
	})
}

func ValidateDynamicData(t *testing.T, api *apim.API) {

	// Retrieve the endpointConfig of the imported API
	endpointConfig := api.GetEndPointConfig()

	// Check whether the production endpoint has the expected value set using the env variable
	productionEndpoints := endpointConfig.(map[string]interface{})["production_endpoints"].(map[string]interface{})
	assert.Equal(t, os.Getenv("DEV_ENV_PROD_URL"), productionEndpoints["url"], "Production endpoint value mismatched")

	// Check whether the sandbox endpoint has the expected value set using the env variable
	sandboxEndpoints := endpointConfig.(map[string]interface{})["sandbox_endpoints"].(map[string]interface{})
	assert.Equal(t, os.Getenv("DEV_ENV_SAND_URL"), sandboxEndpoints["url"], "Sandbox endpoint value mismatched")

	// Check whether the retryDelay and retryTimeOut roduction endpoint
	// config values has the expected values set using the env variables
	assert.Equal(t, os.Getenv("DEV_ENV_PROD_RE_DELAY"),
		productionEndpoints["config"].(map[string]interface{})["retryDelay"].(string),
		"Retry delay value of the production endpoint value mismatched")
	assert.Equal(t, os.Getenv("DEV_ENV_PROD_RE_TO"),
		productionEndpoints["config"].(map[string]interface{})["retryTimeOut"].(string),
		"Retry time out value of the production endpoint config mismatched")
}
