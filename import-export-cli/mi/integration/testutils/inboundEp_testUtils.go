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

// ValidateInboundEndpointList validate ctl output with list of data services from the Management API
func ValidateInboundEndpointList(t *testing.T, inboundEndpointCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, inboundEndpointCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementInboundEndpointResource, &artifactutils.InboundEndpointList{})
	validateinboundEndpointListEqual(t, output, (artifactList.(*artifactutils.InboundEndpointList)))
}

func validateinboundEndpointListEqual(t *testing.T, inboundEndpointsListFromCtl string, inboundEndpointsList *artifactutils.InboundEndpointList) {
	unmatchedCount := inboundEndpointsList.Count
	for _, inboundEndpoint := range inboundEndpointsList.InboundEndpoints {
		assert.Truef(t, strings.Contains(inboundEndpointsListFromCtl, inboundEndpoint.Name), "inboundEndpointsListFromCtl: "+inboundEndpointsListFromCtl+
			" , does not contain inboundEndpoint.Name: "+inboundEndpoint.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Inbound enpoint lists are not equal")
}

// ValidateInboundEndpoint validate ctl output with the data service from the Management API
func ValidateInboundEndpoint(t *testing.T, inboundEndpointCmd string, config *MiConfig, inboundEndpointName string) {
	t.Helper()
	output, _ := GetArtifact(t, inboundEndpointCmd, inboundEndpointName, config)
	artifactList := config.MIClient.GetArtifactFromAPI(utils.MiManagementInboundEndpointResource, getParamMap("inboundEndpointName", inboundEndpointName), &artifactutils.InboundEndpoint{})
	validateinboundEndpointEqual(t, output, (artifactList.(*artifactutils.InboundEndpoint)))
}

func validateinboundEndpointEqual(t *testing.T, inboundEndpointsListFromCtl string, inboundEndpoint *artifactutils.InboundEndpoint) {
	assert.Contains(t, inboundEndpointsListFromCtl, inboundEndpoint.Name)
	assert.Contains(t, inboundEndpointsListFromCtl, inboundEndpoint.Type)
	for _, param := range inboundEndpoint.Parameters {
		assert.Contains(t, inboundEndpointsListFromCtl, param.Name)
		assert.Contains(t, inboundEndpointsListFromCtl, param.Value)
	}
}
