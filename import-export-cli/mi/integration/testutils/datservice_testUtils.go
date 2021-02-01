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

// ValidateDataServicesList validate ctl output with list of data services from the Management API
func ValidateDataServicesList(t *testing.T, dataServiceCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, dataServiceCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementDataServiceResource, &artifactutils.DataServicesList{})
	validateDataServiceListEqual(t, output, (artifactList.(*artifactutils.DataServicesList)))
}

func validateDataServiceListEqual(t *testing.T, dataServicesListFromCtl string, dataServicesList *artifactutils.DataServicesList) {
	unmatchedCount := dataServicesList.Count
	for _, dataService := range dataServicesList.List {
		assert.Truef(t, strings.Contains(dataServicesListFromCtl, dataService.ServiceName), "dataServicesListFromCtl: "+dataServicesListFromCtl+
			" , does not contain dataService.ServiceName: "+dataService.ServiceName)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Data Service lists are not equal")
}

// ValidateDataService validate ctl output with the data service from the Management API
func ValidateDataService(t *testing.T, dataServiceCmd string, config *MiConfig, dataServiceName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, dataServiceCmd, dataServiceName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementDataServiceResource, getParamMap("dataServiceName", dataServiceName), &artifactutils.DataServiceInfo{})
	validateDataServiceEqual(t, output, (artifact.(*artifactutils.DataServiceInfo)))
}

func validateDataServiceEqual(t *testing.T, dataServicesListFromCtl string, dataService *artifactutils.DataServiceInfo) {
	assert.Contains(t, dataServicesListFromCtl, dataService.ServiceName)
	assert.Contains(t, dataServicesListFromCtl, dataService.ServiceGroupName)
	assert.Contains(t, dataServicesListFromCtl, dataService.Wsdl11)
	assert.Contains(t, dataServicesListFromCtl, dataService.Wsdl20)
	for _, query := range dataService.Queries {
		assert.Contains(t, dataServicesListFromCtl, query.Id)
	}
}
