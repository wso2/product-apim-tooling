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

package impl

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	defaultTaskListTableFormat = "table {{.Name}}"
	defaultTaskDetailedFormat  = "detail Name - {{.Name}}\n" +
		"Trigger Type - {{.Type}}\n" +
		"{{if .TriggerCron}}Cron Expression - {{.TriggerCron}}" +
		"{{else}}" +
		"Trigger Count - {{.TriggerCount}}\n" +
		"Trigger Interval - {{.TriggerInterval}}" +
		"{{end}}"
)

// GetTaskList returns a list of Tasks deployed in the micro integrator in a given environment
func GetTaskList(env string) (*artifactutils.TaskList, error) {
	resp, err := getArtifactList(utils.MiManagementTaskResource, env, &artifactutils.TaskList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.TaskList), nil
}

// PrintTaskList print a list of Tasks according to the given format
func PrintTaskList(taskList *artifactutils.TaskList, format string) {
	if taskList.Count > 0 {
		tasks := taskList.Tasks
		taskListContext := getContextWithFormat(format, defaultTaskListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, task := range tasks {
				if err := t.Execute(w, task); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}
		taskListTableHeaders := map[string]string{
			"Name": nameHeader,
		}
		if err := taskListContext.Write(renderer, taskListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Tasks found")
	}
}

// GetTask returns a information about a specific Task deployed in the micro integrator in a given environment
func GetTask(env, taskName string) (*artifactutils.Task, error) {
	resp, err := getArtifactInfo(utils.MiManagementTaskResource, "taskName", taskName, env, &artifactutils.Task{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.Task), nil
}

// PrintTaskDetails prints details about a Task according to the given format
func PrintTaskDetails(task *artifactutils.Task, format string) {
	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultTaskDetailedFormat
	}

	taskContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRendererEndsWithNewLine(task)

	if err := taskContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
