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

// ValidateTaskList validate ctl output with list of tasks from the Management API
func ValidateTaskList(t *testing.T, taskCmd string, config *MiConfig) {
	t.Helper()
	output, _ := ListArtifacts(t, taskCmd, config)
	artifactList := config.MIClient.GetArtifactListFromAPI(utils.MiManagementTaskResource, &artifactutils.TaskList{})
	validateTaskListEqual(t, output, (artifactList.(*artifactutils.TaskList)))
}

func validateTaskListEqual(t *testing.T, taskListFromCtl string, taskList *artifactutils.TaskList) {
	unmatchedCount := taskList.Count
	for _, task := range taskList.Tasks {
		assert.Truef(t, strings.Contains(taskListFromCtl, task.Name), "taskListFromCtl: "+taskListFromCtl+
			" , does not contain task.Name: "+task.Name)
		unmatchedCount--
	}
	assert.Equal(t, 0, int(unmatchedCount), "task lists are not equal")
}

// ValidateTask validate ctl output with the task from the Management API
func ValidateTask(t *testing.T, taskCmd string, config *MiConfig, taskName string) {
	t.Helper()
	output, _ := GetArtifact(t, taskCmd, taskName, config)
	artifactList := config.MIClient.GetArtifactFromAPI(utils.MiManagementTaskResource, "taskName", taskName, &artifactutils.Task{})
	validateTaskEqual(t, output, (artifactList.(*artifactutils.Task)))
}

func validateTaskEqual(t *testing.T, taskFromCtl string, task *artifactutils.Task) {
	assert.Contains(t, taskFromCtl, task.Name)
	assert.Contains(t, taskFromCtl, task.Type)
	assert.Contains(t, taskFromCtl, task.TriggerCron)
	assert.Contains(t, taskFromCtl, task.TriggerCount)
	assert.Contains(t, taskFromCtl, task.TriggerInterval)
}
