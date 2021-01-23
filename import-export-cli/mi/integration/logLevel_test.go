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

const validLoggerName = "org-apache-coyote"
const invalidLoggerName = "abc-logger"
const logLevelCmd = "log-levels"

func TestGetLoggerByName(t *testing.T) {
	testutils.ValidateLogger(t, logLevelCmd, config, validLoggerName)
}

func TestGetNonExistingLoggerByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, logLevelCmd, invalidLoggerName, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of logger [ "+invalidLoggerName+" ]  Logger name ('"+invalidLoggerName+"') not found")
}

func TestGetLoggersWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, logLevelCmd, validLoggerName)
}

func TestGetLoggersWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, logLevelCmd, config, validLoggerName)
}

func TestGetLoggersWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, logLevelCmd, config, validLoggerName)
}

func TestGetLoggersWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, true, logLevelCmd, validLoggerName, invalidLoggerName)
}
