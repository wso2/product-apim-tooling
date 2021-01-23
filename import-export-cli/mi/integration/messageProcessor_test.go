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

const validMessageProcessor = "scheduled-msg-processor"
const invalidMessageProcessor = "abc-msg-processor"
const messageProcessorCmd = "message-processors"

func TestGetMessageProcessors(t *testing.T) {
	testutils.ValidateMessageProcessorList(t, messageProcessorCmd, config)
}

func TestGetMessageProcessorByName(t *testing.T) {
	testutils.ValidateMessageProcessor(t, messageProcessorCmd, config, validMessageProcessor)
}

func TestGetNonExistingMessageProcessorByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, messageProcessorCmd, invalidMessageProcessor, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of message processors [ "+invalidMessageProcessor+" ]  Specified message processor ('"+invalidMessageProcessor+"') not found")
}

func TestGetMessageProcessorsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, messageProcessorCmd)
}

func TestGetMessageProcessorsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, messageProcessorCmd, config)
}

func TestGetMessageProcessorsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, messageProcessorCmd, config)
}

func TestGetMessageProcessorsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, messageProcessorCmd, validMessageProcessor, invalidMessageProcessor)
}
