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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ValidateLogFileList validate ctl output with list of log file names from the Management API
func ValidateLogFileList(t *testing.T, config *MiConfig, logCmd string) {
	t.Helper()
	output, _ := ListArtifacts(t, logCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementLogResource, &artifactutils.LogFileList{})
	validateLogFileListEqual(t, output, (artifactList.(*artifactutils.LogFileList)))
}

func validateLogFileListEqual(t *testing.T, logFileListFromCtl string, logFileList *artifactutils.LogFileList) {
	filteredLogFileList := filterOnlyLogFiles(logFileList)
	unmatchedCount := filteredLogFileList.Count
	for _, logFile := range filteredLogFileList.LogFiles {
		assert.Truef(t, strings.Contains(logFileListFromCtl, logFile.FileName), "logFileListFromCtl: "+logFileListFromCtl+
			" , does not contain logFile.FileName: "+logFile.FileName)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "log file lists are not equal")
}

// ValidateLogFile validate wether the log file is downloaded
func ValidateLogFile(t *testing.T, config *MiConfig, logCmd, loggerName string) {
	GetArtifact(t, config, logCmd, loggerName)
	assert.True(t, base.IsFileAvailable(loggerName))
	t.Cleanup(func() {
		os.RemoveAll(loggerName)
	})
}

func filterOnlyLogFiles(logFileList *artifactutils.LogFileList) *artifactutils.LogFileList {
	filteredList := new(artifactutils.LogFileList)
	for _, logFile := range logFileList.LogFiles {
		if strings.HasSuffix(logFile.FileName, ".log") {
			filteredList.LogFiles = append(filteredList.LogFiles, logFile)
		}
	}
	filteredList.Count = int32(len(filteredList.LogFiles))
	return filteredList
}
