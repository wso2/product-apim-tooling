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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ValidateLogger validate ctl output with the logger from the Management API
func ValidateLogger(t *testing.T, logLevelCmd string, config *MiConfig, loggerName string) {
	t.Helper()
	output, _ := GetArtifact(t, logLevelCmd, loggerName, config)
	artifactList := config.MIClient.GetArtifactFromAPI(utils.MiManagementLoggingResource, "loggerName", loggerName, &artifactutils.Logger{})
	validateLoggerEqual(t, output, (artifactList.(*artifactutils.Logger)))
}

func validateLoggerEqual(t *testing.T, loggerFromCtl string, logger *artifactutils.Logger) {
	assert.Contains(t, loggerFromCtl, logger.LoggerName)
	assert.Contains(t, loggerFromCtl, logger.ComponentName)
	assert.Contains(t, loggerFromCtl, logger.LogLevel)
}
