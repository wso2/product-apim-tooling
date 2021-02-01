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

// ValidateMessageProcessorList validate ctl output with list of message processors from the Management API
func ValidateMessageProcessorList(t *testing.T, messageProcessorCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, messageProcessorCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementMessageProcessorResource, &artifactutils.MessageProcessorList{})
	validateMessageProcessorListEqual(t, output, (artifactList.(*artifactutils.MessageProcessorList)))
}

func validateMessageProcessorListEqual(t *testing.T, messageProcessorListFromCtl string, messageProcessorList *artifactutils.MessageProcessorList) {
	unmatchedCount := messageProcessorList.Count
	for _, messageProcessor := range messageProcessorList.MessageProcessors {
		assert.Truef(t, strings.Contains(messageProcessorListFromCtl, messageProcessor.Name), "messageProcessorListFromCtl: "+messageProcessorListFromCtl+
			" , does not contain messageProcessor.Name: "+messageProcessor.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Message Processor lists are not equal")
}

// ValidateMessageProcessor validate ctl output with the message processor from the Management API
func ValidateMessageProcessor(t *testing.T, messageProcessorCmd string, config *MiConfig, messageProcessorName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, messageProcessorCmd, messageProcessorName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementMessageProcessorResource, getParamMap("messageProcessorName", messageProcessorName), &artifactutils.MessageProcessorData{})
	validateMessageProcessorEqual(t, output, (artifact.(*artifactutils.MessageProcessorData)))
}

func validateMessageProcessorEqual(t *testing.T, messageProcessorFromCtl string, messageProcessor *artifactutils.MessageProcessorData) {
	assert.Contains(t, messageProcessorFromCtl, messageProcessor.Name)
	assert.Contains(t, messageProcessorFromCtl, messageProcessor.Type)
	assert.Contains(t, messageProcessorFromCtl, messageProcessor.Store)
	assert.Contains(t, messageProcessorFromCtl, messageProcessor.Container)
	for key, value := range messageProcessor.Parameters {
		assert.Contains(t, messageProcessorFromCtl, key)
		assert.Contains(t, messageProcessorFromCtl, value)
	}
}
