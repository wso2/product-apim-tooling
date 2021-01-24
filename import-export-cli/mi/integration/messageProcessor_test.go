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
const messageProcessorsCmd = "message-processors"
const messageProcessorCmd = "message-processor"

func TestGetMessageProcessors(t *testing.T) {
	testutils.ValidateMessageProcessorList(t, messageProcessorsCmd, config)
}

func TestGetMessageProcessorByName(t *testing.T) {
	testutils.ValidateMessageProcessor(t, messageProcessorsCmd, config, validMessageProcessor)
}

func TestGetNonExistingMessageProcessorByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, messageProcessorsCmd, invalidMessageProcessor, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of message processors [ "+invalidMessageProcessor+" ]  Specified message processor ('"+invalidMessageProcessor+"') not found")
}

func TestGetMessageProcessorsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, messageProcessorsCmd)
}

func TestGetMessageProcessorsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, messageProcessorsCmd, config)
}

func TestGetMessageProcessorsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, messageProcessorsCmd, config)
}

func TestGetMessageProcessorsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, messageProcessorsCmd, validMessageProcessor, invalidMessageProcessor)
}

func TestActivateMessageProcessor(t *testing.T) {
	expected := validMessageProcessor + " : is activated"
	testutils.ExecActivateCommand(t, config, messageProcessorCmd, validMessageProcessor, expected)
}

func TestActivateNonExistingMessageProcessor(t *testing.T) {
	expected := "[ERROR]: Activating message processor [ " + invalidMessageProcessor + " ] Message processor does not exist"
	testutils.ExecActivateCommand(t, config, messageProcessorCmd, invalidMessageProcessor, expected)
}

func TestActivateMessageProcessorWithoutEnvFlag(t *testing.T) {
	testutils.ExecActivateCommandWithoutEnvFlag(t, config, messageProcessorCmd, validMessageProcessor)
}

func TestActivateMessageProcessorWithInvalidArgs(t *testing.T) {
	testutils.ExecActivateCommandWithInvalidArgCount(t, config, 1, 0, messageProcessorCmd)
}

func TestActivateMessageProcessorWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecActivateCommandWithoutSettingEnv(t, messageProcessorCmd, validMessageProcessor)
}

func TestActivateMessageProcessorWithoutLogin(t *testing.T) {
	testutils.ExecActivateCommandWithoutLogin(t, config, messageProcessorCmd, validMessageProcessor)
}
