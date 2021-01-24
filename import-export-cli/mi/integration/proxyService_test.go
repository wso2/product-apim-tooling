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

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/integration/testutils"
)

const validProxyServiceName = "StockQuoteProxy"
const invalidProxyServiceName = "abc-proxy"
const proxyServicesCmd = "proxy-services"
const proxyServiceCmd = "proxy-service"

func TestGetProxyServices(t *testing.T) {
	testutils.ValidateProxyServiceList(t, proxyServicesCmd, config)
}

func TestGetProxyServiceByName(t *testing.T) {
	testutils.ValidateProxyService(t, proxyServicesCmd, config, validProxyServiceName)
}

func TestGetNonExistingProxyServiceByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, proxyServicesCmd, invalidProxyServiceName, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of proxy services [ "+invalidProxyServiceName+" ]  404 Not Found")
}

func TestGetProxyServicesWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, proxyServicesCmd)
}

func TestGetProxyServicesWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, proxyServicesCmd, config)
}

func TestGetProxyServicesWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, proxyServicesCmd, config)
}

func TestGetProxyServicesWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, proxyServicesCmd, validProxyServiceName, invalidProxyServiceName)
}

func TestActivateProxyService(t *testing.T) {
	expected := validProxyServiceName + " started successfully"
	testutils.ExecActivateCommand(t, config, proxyServiceCmd, validProxyServiceName, expected)
}

func TestActivateNonExistingProxyService(t *testing.T) {
	expected := "[ERROR]: Activating proxy service [ " + invalidProxyServiceName + " ] Proxy service could not be found"
	testutils.ExecActivateCommand(t, config, proxyServiceCmd, invalidProxyServiceName, expected)
}

func TestActivateProxyServiceWithoutEnvFlag(t *testing.T) {
	testutils.ExecActivateCommandWithoutEnvFlag(t, config, proxyServiceCmd, validProxyServiceName)
}

func TestActivateProxyServiceWithInvalidArgs(t *testing.T) {
	testutils.ExecActivateCommandWithInvalidArgCount(t, config, 1, 0, proxyServiceCmd)
}

func TestActivateProxyServiceWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecActivateCommandWithoutSettingEnv(t, proxyServiceCmd, validProxyServiceName)
}

func TestActivateProxyServiceWithoutLogin(t *testing.T) {
	testutils.ExecActivateCommandWithoutLogin(t, config, proxyServiceCmd, validProxyServiceName)
}
