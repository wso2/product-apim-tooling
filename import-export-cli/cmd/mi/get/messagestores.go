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
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
)

var getMessageStoreCmdEnvironment string
var getMessageStoreCmdFormat string

const artifactMessageStores = "message stores"
const getMessageStoreCmdLiteral = "message-stores [messagestore-name]"

var getMessageStoreCmd = &cobra.Command{
	Use:     getMessageStoreCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactMessageStores),
	Long:    generateGetCmdLongDescForArtifact(artifactMessageStores, "messagestore-name"),
	Example: generateGetCmdExamplesForArtifact(artifactMessageStores, miUtils.GetTrimmedCmdLiteral(getMessageStoreCmdLiteral), "TestMessageStore"),
	Args:    cobra.MaximumNArgs(1),
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handleGetMessageStoreCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getMessageStoreCmd)
	setEnvFlag(getMessageStoreCmd, &getMessageStoreCmdEnvironment)
	setFormatFlag(getMessageStoreCmd, &getMessageStoreCmdFormat)
}

func handleGetMessageStoreCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getMessageStoreCmdLiteral))
	credentials.HandleMissingCredentials(getMessageStoreCmdEnvironment)
	if len(args) == 1 {
		var messageStoreName = args[0]
		executeShowMessageStore(messageStoreName)
	} else {
		executeListMessageStores()
	}
}

func executeListMessageStores() {
	messageStoreList, err := impl.GetMessageStoreList(getMessageStoreCmdEnvironment)
	if err == nil {
		impl.PrintMessageStoreList(messageStoreList, getMessageStoreCmdFormat)
	} else {
		printErrorForArtifactList(artifactMessageStores, err)
	}
}

func executeShowMessageStore(messageStoreName string) {
	messageStore, err := impl.GetMessageStore(getMessageStoreCmdEnvironment, messageStoreName)
	if err == nil {
		impl.PrintMessageStoreDetails(messageStore, getMessageStoreCmdFormat)
	} else {
		printErrorForArtifact(artifactMessageStores, messageStoreName, err)
	}
}
