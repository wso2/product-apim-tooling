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

	"github.com/wso2/product-apim-tooling/import-export-cli/mi/integration/testutils"
)

const connectorCmd = "connectors"

func TestGetConnectors(t *testing.T) {
	testutils.ValidateConnectorList(t, connectorCmd, config)
}

func TestGetConnectorsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, cAppCmd)
}

func TestGetConnectorsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, connectorCmd, config)
}

func TestGetConnectorsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, connectorCmd, config)
}

func TestGetConnectorsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 0, 2, true, connectorCmd, "abc", "123")
}
