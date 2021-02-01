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

// ValidateLocalEntryList validate ctl output with list of local entries from the Management API
func ValidateLocalEntryList(t *testing.T, localEntryCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, localEntryCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementLocalEntrieResource, &artifactutils.LocalEntryList{})
	validateLocalEntryListEqual(t, output, (artifactList.(*artifactutils.LocalEntryList)))
}

func validateLocalEntryListEqual(t *testing.T, localEntryListFromCtl string, localEntryList *artifactutils.LocalEntryList) {
	unmatchedCount := localEntryList.Count
	for _, localEntry := range localEntryList.LocalEntries {
		assert.Truef(t, strings.Contains(localEntryListFromCtl, localEntry.Name), "localEntryListFromCtl: "+localEntryListFromCtl+
			" , does not contain localEntry.Name: "+localEntry.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Local entry lists are not equal")
}

// ValidateLocalEntry validate ctl output with the local entry from the Management API
func ValidateLocalEntry(t *testing.T, localEntryCmd string, config *MiConfig, localEntryName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, localEntryCmd, localEntryName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementLocalEntrieResource, getParamMap("localEntryName", localEntryName), &artifactutils.LocalEntryData{})
	validateLocalEntryEqual(t, output, (artifact.(*artifactutils.LocalEntryData)))
}

func validateLocalEntryEqual(t *testing.T, localEntryListFromCtl string, localEntry *artifactutils.LocalEntryData) {
	assert.Contains(t, localEntryListFromCtl, localEntry.Name)
	assert.Contains(t, localEntryListFromCtl, localEntry.Type)
	assert.Contains(t, localEntryListFromCtl, localEntry.Value)
}
