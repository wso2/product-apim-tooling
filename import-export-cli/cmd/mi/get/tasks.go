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

package get

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
)

var getTaskCmdEnvironment string
var getTaskCmdFormat string

const artifactTasks = "tasks"
const getTaskCmdLiteral = "tasks [task-name]"

var getTasksCmd = &cobra.Command{
	Use:     getTaskCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactTasks),
	Long:    generateGetCmdLongDescForArtifact(artifactTasks, "task-name"),
	Example: generateGetCmdExamplesForArtifact(artifactTasks, getTrimmedCmdLiteral(getTaskCmdLiteral), "SampleTask"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetTaskCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getTasksCmd)
	getTasksCmd.Flags().StringVarP(&getTaskCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getTasksCmd.Flags().StringVarP(&getTaskCmdFormat, "format", "", "", generateFormatFlagUsage(artifactTasks))
	getTasksCmd.MarkFlagRequired("environment")
}

func handleGetTaskCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getTrimmedCmdLiteral(getTaskCmdLiteral))
	credentials.HandleMissingCredentials(getTaskCmdEnvironment)
	if len(args) == 1 {
		var taskName = args[0]
		executeShowTask(taskName)
	} else {
		executeListTasks()
	}
}

func executeListTasks() {

	taskList, err := impl.GetTaskList(getTaskCmdEnvironment)
	if err == nil {
		impl.PrintTaskList(taskList, getTaskCmdFormat)
	} else {
		printErrorForArtifactList(artifactTasks, err)
	}
}

func executeShowTask(taskName string) {
	task, err := impl.GetTask(getTaskCmdEnvironment, taskName)
	if err == nil {
		impl.PrintTaskDetails(task, getTaskCmdFormat)
	} else {
		printErrorForArtifact(artifactTasks, taskName, err)
	}
}
