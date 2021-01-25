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

// ValidateCAppsList validate ctl output with list of capps from the Management API
func ValidateCAppsList(t *testing.T, cAppCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, cAppCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementCarbonAppResource, &artifactutils.CompositeAppList{})
	validateCAppListEqual(t, output, (artifactList.(*artifactutils.CompositeAppList)))
}

func validateCAppListEqual(t *testing.T, cAppsListFromCtl string, cAppsList *artifactutils.CompositeAppList) {
	unmatchedCount := cAppsList.Count
	for _, cApp := range cAppsList.CompositeApps {
		assert.Truef(t, strings.Contains(cAppsListFromCtl, cApp.Name), "cAppsListFromCtl: "+cAppsListFromCtl+
			" , does not contain cApp.Name: "+cApp.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "CApp lists are not equal")
}

// ValidateCApp validate ctl output with the capp from the Management API
func ValidateCApp(t *testing.T, cAppCmd string, config *MiConfig, cAppName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, cAppCmd, cAppName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementCarbonAppResource, getParamMap("carbonAppName", cAppName), &artifactutils.CompositeApp{})
	calidateCAppEqual(t, output, (artifact.(*artifactutils.CompositeApp)))
}

func calidateCAppEqual(t *testing.T, CAppsListFromCtl string, cApp *artifactutils.CompositeApp) {
	assert.Contains(t, CAppsListFromCtl, cApp.Name)
	assert.Contains(t, CAppsListFromCtl, cApp.Version)
	for _, artifact := range cApp.Artifacts {
		assert.Contains(t, CAppsListFromCtl, artifact.Name)
		assert.Contains(t, CAppsListFromCtl, artifact.Type)
	}
}
