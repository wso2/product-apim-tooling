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

// ValidateEndpointList validate ctl output with list of endpoints from the Management API
func ValidateEndpointList(t *testing.T, endpointCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, endpointCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementEndpointResource, &artifactutils.EndpointList{})
	validateEndpointListEqual(t, output, (artifactList.(*artifactutils.EndpointList)))
}

func validateEndpointListEqual(t *testing.T, endpointsListFromCtl string, endpointList *artifactutils.EndpointList) {
	unmatchedCount := endpointList.Count
	for _, endpoint := range endpointList.Endpoints {
		assert.Truef(t, strings.Contains(endpointsListFromCtl, endpoint.Name), "endpointsListFromCtl: "+endpointsListFromCtl+
			" , does not contain endpoint.Name: "+endpoint.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Endpoint lists are not equal")
}

// ValidateEndpoint validate ctl output with the endpoint from the Management API
func ValidateEndpoint(t *testing.T, endpointCmd string, config *MiConfig, endpointName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, endpointCmd, endpointName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementEndpointResource, getParamMap("endpointName", endpointName), &artifactutils.Endpoint{})
	validateEndpointEqual(t, output, (artifact.(*artifactutils.Endpoint)))
}

func validateEndpointEqual(t *testing.T, endpointFromCtl string, endpoint *artifactutils.Endpoint) {
	assert.Contains(t, endpointFromCtl, endpoint.Name)
	assert.Contains(t, endpointFromCtl, endpoint.Type)
	assert.Contains(t, endpointFromCtl, endpoint.Method)
	assert.Contains(t, endpointFromCtl, endpoint.URITemplate)
}
