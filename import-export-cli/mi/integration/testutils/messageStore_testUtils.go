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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ValidateMessageStoreList validate ctl output with list of message stores from the Management API
func ValidateMessageStoreList(t *testing.T, messageStoreCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, messageStoreCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementMessageStoreResource, &artifactutils.MessageStoreList{})
	validateMessageStoreListEqual(t, output, (artifactList.(*artifactutils.MessageStoreList)))
}

func validateMessageStoreListEqual(t *testing.T, messageStoreListFromCtl string, messageStoreList *artifactutils.MessageStoreList) {
	unmatchedCount := messageStoreList.Count
	for _, messageStore := range messageStoreList.MessageStores {
		assert.Truef(t, strings.Contains(messageStoreListFromCtl, messageStore.Name), "messageStoreListFromCtl: "+messageStoreListFromCtl+
			" , does not contain messageStore.Name: "+messageStore.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Message Store lists are not equal")
}

// ValidateMessageStore validate ctl output with the message store from the Management API
func ValidateMessageStore(t *testing.T, messageStoreCmd string, config *MiConfig, messageStoreName string) {
	t.Helper()
	output, _ := GetArtifact(t, messageStoreCmd, messageStoreName, config)
	artifactList := config.MIClient.GetArtifactFromAPI(utils.MiManagementMessageStoreResource, getParamMap("messageStoreName", messageStoreName), &artifactutils.MessageStoreData{})
	validateMessageStoreEqual(t, output, (artifactList.(*artifactutils.MessageStoreData)))
}

func validateMessageStoreEqual(t *testing.T, messageStoreFromCtl string, messageStore *artifactutils.MessageStoreData) {
	assert.Contains(t, messageStoreFromCtl, messageStore.Name)
	assert.Contains(t, messageStoreFromCtl, messageStore.Container)
	assert.Contains(t, messageStoreFromCtl, messageStore.FileName)
	assert.Contains(t, messageStoreFromCtl, messageStore.Consumer)
	assert.Contains(t, messageStoreFromCtl, messageStore.Producer)
	assert.Contains(t, messageStoreFromCtl, fmt.Sprint(messageStore.Size))
	for key, value := range messageStore.Properties {
		assert.Contains(t, messageStoreFromCtl, key)
		assert.Contains(t, messageStoreFromCtl, value)
	}
}
