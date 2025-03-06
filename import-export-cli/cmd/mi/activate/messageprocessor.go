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

package activate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
)

var activateMessageProcessorCmdEnvironment string

const artifactMessageProcessor = "message processor"
const activateMessageProcessorCmdLiteral = "message-processor [messageprocessor-name]"

var activateMessageProcessorCmd = &cobra.Command{
	Use:     activateMessageProcessorCmdLiteral,
	Short:   generateActivateCmdShortDescForArtifact(artifactMessageProcessor),
	Long:    generateActivateCmdLongDescForArtifact(artifactMessageProcessor, "messageprocessor-name"),
	Example: generateActivateCmdExamplesForArtifact(artifactMessageProcessor, miUtils.GetTrimmedCmdLiteral(activateMessageProcessorCmdLiteral), "TestMessageProcessor"),
	Args:    cobra.ExactArgs(1),
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handleActivateMessageProcessorCmdArguments(args)
	},
}

func init() {
	ActivateCmd.AddCommand(activateMessageProcessorCmd)
	setEnvFlag(activateMessageProcessorCmd, &activateMessageProcessorCmdEnvironment, artifactMessageProcessor)
}

func handleActivateMessageProcessorCmdArguments(args []string) {
	printActivateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(activateMessageProcessorCmdLiteral))
	credentials.HandleMissingCredentials(activateMessageProcessorCmdEnvironment)
	executeActivateMessageProcessor(args[0])
}

func executeActivateMessageProcessor(messageProcessorName string) {
	resp, err := impl.ActivateMessageProcessor(activateMessageProcessorCmdEnvironment, messageProcessorName)
	if err != nil {
		printErrorForArtifact(artifactMessageProcessor, messageProcessorName, err)
	} else {
		fmt.Println(resp)
	}
}
