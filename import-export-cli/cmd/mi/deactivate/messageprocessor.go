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

package deactivate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
)

var deactivateMessageProcessorCmdEnvironment string

const artifactMessageProcessor = "message processor"
const deactivateMessageProcessorCmdLiteral = "message-processor [messageprocessor-name]"

var deactivateMessageProcessorCmd = &cobra.Command{
	Use:     deactivateMessageProcessorCmdLiteral,
	Short:   generateDeactivateCmdShortDescForArtifact(artifactMessageProcessor),
	Long:    generateDeactivateCmdLongDescForArtifact(artifactMessageProcessor, "messageprocessor-name"),
	Example: generateDeactivateCmdExamplesForArtifact(artifactMessageProcessor, miUtils.GetTrimmedCmdLiteral(deactivateMessageProcessorCmdLiteral), "TestMessageProcessor"),
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleDeactivateMessageProcessorCmdArguments(args)
	},
}

func init() {
	DeactivateCmd.AddCommand(deactivateMessageProcessorCmd)
	setEnvFlag(deactivateMessageProcessorCmd, &deactivateMessageProcessorCmdEnvironment, artifactMessageProcessor)
}

func handleDeactivateMessageProcessorCmdArguments(args []string) {
	printDeactivateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(deactivateMessageProcessorCmdLiteral))
	credentials.HandleMissingCredentials(deactivateMessageProcessorCmdEnvironment)
	executeDeactivateMessageProcessor(args[0])
}

func executeDeactivateMessageProcessor(messageProcessorName string) {
	resp, err := impl.DeactivateMessageProcessor(deactivateMessageProcessorCmdEnvironment, messageProcessorName)
	if err != nil {
		printErrorForArtifact(artifactMessageProcessor, messageProcessorName, err)
	} else {
		fmt.Println(resp)
	}
}
