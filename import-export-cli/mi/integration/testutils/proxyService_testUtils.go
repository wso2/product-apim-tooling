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

// ValidateProxyServiceList validate ctl output with list of proxies from the Management API
func ValidateProxyServiceList(t *testing.T, proxyServiceCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, proxyServiceCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementProxyServiceResource, &artifactutils.ProxyServiceList{})
	validateProxyServiceListEqual(t, output, (artifactList.(*artifactutils.ProxyServiceList)))
}

func validateProxyServiceListEqual(t *testing.T, proxyServiceListFromCtl string, proxyServiceList *artifactutils.ProxyServiceList) {
	unmatchedCount := proxyServiceList.Count
	for _, proxyService := range proxyServiceList.Proxies {
		assert.Truef(t, strings.Contains(proxyServiceListFromCtl, proxyService.Name), "proxyServiceListFromCtl: "+proxyServiceListFromCtl+
			" , does not contain proxyService.Name: "+proxyService.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Proxy Service lists are not equal")
}

// ValidateProxyService validate ctl output with the proxy from the Management API
func ValidateProxyService(t *testing.T, proxyServiceCmd string, config *MiConfig, proxyServiceName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, proxyServiceCmd, proxyServiceName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementProxyServiceResource, getParamMap("proxyServiceName", proxyServiceName), &artifactutils.Proxy{})
	validateProxyServiceEqual(t, output, (artifact.(*artifactutils.Proxy)))
}

func validateProxyServiceEqual(t *testing.T, proxyServiceFromCtl string, proxyService *artifactutils.Proxy) {
	assert.Contains(t, proxyServiceFromCtl, proxyService.Name)
	assert.Contains(t, proxyServiceFromCtl, proxyService.Wsdl11)
	assert.Contains(t, proxyServiceFromCtl, proxyService.Wsdl20)
	assert.Contains(t, proxyServiceFromCtl, proxyService.Stats)
	assert.Contains(t, proxyServiceFromCtl, proxyService.Tracing)
}
