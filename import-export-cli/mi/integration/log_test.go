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

const validLogFileName = "wso2carbon.log"
const invalidLogFileName = "abc.log"
const logsCmd = "logs"

func TestGetLogs(t *testing.T) {
	testutils.ValidateLogFileList(t, config, logsCmd)
}

func TestGetLogByName(t *testing.T) {
	testutils.ValidateLogFile(t, config, logsCmd, validLogFileName)
}

func TestGetNonExistingLogFileByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, config, logsCmd, invalidLogFileName)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of log file [ "+invalidLogFileName+" ]  500 Internal Server Error")
}

func TestGetLogsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, logsCmd)
}

func TestGetLogsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, logsCmd, config)
}

func TestGetLogsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, logsCmd, config)
}

func TestGetLogsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, logsCmd, validLogFileName, invalidLogFileName)
}
