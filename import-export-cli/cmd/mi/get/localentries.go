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

var getLocalEntryCmdEnvironment string
var getLocalEntryCmdFormat string

const artifactLocalEntries = "local entries"
const getLocalEntryCmdLiteral = "local-entries [localentry-name]"

var getLocalEntryCmd = &cobra.Command{
	Use:     getLocalEntryCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactLocalEntries),
	Long:    generateGetCmdLongDescForArtifact(artifactLocalEntries, "localentry-name"),
	Example: generateGetCmdExamplesForArtifact(artifactLocalEntries, getTrimmedCmdLiteral(getLocalEntryCmdLiteral), "SampleLocalEntry"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetLocalEntryCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getLocalEntryCmd)
	setEnvFlag(getLocalEntryCmd, &getLocalEntryCmdEnvironment)
	setFormatFlag(getLocalEntryCmd, &getLocalEntryCmdFormat)
}

func handleGetLocalEntryCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getTrimmedCmdLiteral(getLocalEntryCmdLiteral))
	credentials.HandleMissingCredentials(getLocalEntryCmdEnvironment)
	if len(args) == 1 {
		var LocalEntryName = args[0]
		executeShowLocalEntry(LocalEntryName)
	} else {
		executeListLocalEntrys()
	}
}

func executeListLocalEntrys() {
	localEntryList, err := impl.GetLocalEntryList(getLocalEntryCmdEnvironment)
	if err == nil {
		impl.PrintLocalEntryList(localEntryList, getLocalEntryCmdFormat)
	} else {
		printErrorForArtifactList(artifactLocalEntries, err)
	}
}

func executeShowLocalEntry(localEntryName string) {
	localEntry, err := impl.GetLocalEntry(getLocalEntryCmdEnvironment, localEntryName)
	if err == nil {
		impl.PrintLocalEntryDetails(localEntry, getLocalEntryCmdFormat)
	} else {
		printErrorForArtifact(artifactLocalEntries, localEntryName, err)
	}
}
