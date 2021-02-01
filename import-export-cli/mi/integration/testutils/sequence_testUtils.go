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

// ValidateSequenceList validate ctl output with list of sequences from the Management API
func ValidateSequenceList(t *testing.T, sequenceCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, sequenceCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementSequenceResource, &artifactutils.SequenceList{})
	validateSequenceListEqual(t, output, (artifactList.(*artifactutils.SequenceList)))
}

func validateSequenceListEqual(t *testing.T, sequenceListFromCtl string, sequenceList *artifactutils.SequenceList) {
	unmatchedCount := sequenceList.Count
	for _, sequence := range sequenceList.Sequences {
		assert.Truef(t, strings.Contains(sequenceListFromCtl, sequence.Name), "sequenceListFromCtl: "+sequenceListFromCtl+
			" , does not contain sequence.Name: "+sequence.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "Sequence lists are not equal")
}

// ValidateSequence validate ctl output with the sequence from the Management API
func ValidateSequence(t *testing.T, sequenceCmd string, config *MiConfig, sequenceName string) {
	t.Helper()
	output, _ := GetArtifact(t, config, sequenceCmd, sequenceName)
	artifact := config.MIClient.GetArtifactFromAPI(utils.MiManagementSequenceResource, getParamMap("sequenceName", sequenceName), &artifactutils.Sequence{})
	validateSequenceEqual(t, output, (artifact.(*artifactutils.Sequence)))
}

func validateSequenceEqual(t *testing.T, sequenceFromCtl string, sequence *artifactutils.Sequence) {
	assert.Contains(t, sequenceFromCtl, sequence.Name)
	assert.Contains(t, sequenceFromCtl, sequence.Container)
	assert.Contains(t, sequenceFromCtl, sequence.Stats)
	assert.Contains(t, sequenceFromCtl, sequence.Tracing)
	for _, mediator := range sequence.Mediators {
		assert.Contains(t, sequenceFromCtl, mediator)
	}
}
