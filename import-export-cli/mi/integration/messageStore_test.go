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

const validMessageStore = "in-memory-message-store"
const invalidMessageStore = "abc-msg-store"
const messageStoreCmd = "message-stores"

func TestGetMessageStores(t *testing.T) {
	testutils.ValidateMessageStoreList(t, messageStoreCmd, config)
}

func TestGetMessageStoreByName(t *testing.T) {
	testutils.ValidateMessageStore(t, messageStoreCmd, config, validMessageStore)
}

func TestGetNonExistingMessageStoreByName(t *testing.T) {
	response, _ := testutils.GetArtifact(t, messageStoreCmd, invalidMessageStore, config)
	base.Log(response)
	assert.Contains(t, response, "[ERROR]: Getting Information of message stores [ "+invalidMessageStore+" ]  Specified message store ('"+invalidMessageStore+"') not found")
}

func TestGetMessageStoresWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, messageStoreCmd)
}

func TestGetMessageStoresWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, messageStoreCmd, config)
}

func TestGetMessageStoresWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, messageStoreCmd, config)
}

func TestGetMessageStoresWithInvalidArgs(t *testing.T) {
	testutils.ExecGetCommandWithInvalidArgCount(t, config, 1, 2, false, messageStoreCmd, validMessageStore, invalidMessageStore)
}
