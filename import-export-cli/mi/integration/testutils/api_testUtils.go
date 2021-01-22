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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ValidateAPIsList validate ctl output with list of apis from the Management API
func ValidateAPIsList(t *testing.T, apisCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, apisCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementAPIResource, &artifactutils.IntegrationAPIList{})
	validateAPIListEqual(t, output, (artifactList.(*artifactutils.IntegrationAPIList)))
}

func validateAPIListEqual(t *testing.T, apisListFromCtl string, apisList *artifactutils.IntegrationAPIList) {
	unmatchedCount := apisList.Count
	for _, api := range apisList.Apis {
		assert.Truef(t, strings.Contains(apisListFromCtl, api.Name), "apisListFromCtl: "+apisListFromCtl+
			" , does not contain api.Name: "+api.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "API lists are not equal")
}

// ValidateAPI validate ctl output with the api from the Management API
func ValidateAPI(t *testing.T, apisCmd string, config *MiConfig, apiName string) {
	t.Helper()
	output, _ := GetArtifact(t, apisCmd, apiName, config)
	artifactList := config.MIClient.GetArtifactFromAPI(utils.MiManagementAPIResource, "apiName", apiName, &artifactutils.IntegrationAPI{})
	validateAPIEqual(t, output, (artifactList.(*artifactutils.IntegrationAPI)))
}

func validateAPIEqual(t *testing.T, apisListFromCtl string, api *artifactutils.IntegrationAPI) {
	assert.Contains(t, apisListFromCtl, api.Name)
	assert.Contains(t, apisListFromCtl, api.Stats)
	assert.Contains(t, apisListFromCtl, api.Url)
	assert.Contains(t, apisListFromCtl, api.Version)
}
