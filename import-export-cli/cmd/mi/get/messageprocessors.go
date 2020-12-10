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

var getMessageProcessorCmdEnvironment string
var getMessageProcessorCmdFormat string

const artifactMessageProcessors = "message processors"
const getMessageProcessorCmdLiteral = "message-processors [messageprocessor-name]"

var getMessageProcessorCmd = &cobra.Command{
	Use:     getMessageProcessorCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactMessageProcessors),
	Long:    generateGetCmdLongDescForArtifact(artifactMessageProcessors, "messageprocessor-name"),
	Example: generateGetCmdExamplesForArtifact(artifactMessageProcessors, getTrimmedCmdLiteral(getMessageProcessorCmdLiteral), "TestMessageProcessor"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetMessageProcessorCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getMessageProcessorCmd)
	getMessageProcessorCmd.Flags().StringVarP(&getMessageProcessorCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getMessageProcessorCmd.Flags().StringVarP(&getMessageProcessorCmdFormat, "format", "", "", generateFormatFlagUsage(artifactMessageProcessors))
	getMessageProcessorCmd.MarkFlagRequired("environment")
}

func handleGetMessageProcessorCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getTrimmedCmdLiteral(getMessageProcessorCmdLiteral))
	credentials.HandleMissingCredentials(getMessageProcessorCmdEnvironment)
	if len(args) == 1 {
		var messageProcessorName = args[0]
		executeShowMessageProcessor(messageProcessorName)
	} else {
		executeListMessageProcessors()
	}
}

func executeListMessageProcessors() {

	msgProcessorList, err := impl.GetMessageProcessorList(getMessageProcessorCmdEnvironment)
	if err == nil {
		impl.PrintMessageProcessorList(msgProcessorList, getMessageProcessorCmdFormat)
	} else {
		printErrorForArtifactList(artifactMessageProcessors, err)
	}
}

func executeShowMessageProcessor(msgProcessorName string) {
	msgProcessor, err := impl.GetMessageProcessor(getMessageProcessorCmdEnvironment, msgProcessorName)
	if err == nil {
		impl.PrintMessageProcessorDetails(msgProcessor, getMessageProcessorCmdFormat)
	} else {
		printErrorForArtifact(artifactMessageProcessors, msgProcessorName, err)
	}
}
