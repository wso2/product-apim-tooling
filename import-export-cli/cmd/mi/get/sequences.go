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

var getSequenceCmdEnvironment string
var getSequenceCmdFormat string

const artifactSequences = "sequences"
const getSequenceCmdLiteral = "sequences [sequence-name]"

var getSequenceCmd = &cobra.Command{
	Use:     getSequenceCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactSequences),
	Long:    generateGetCmdLongDescForArtifact(artifactSequences, "sequence-name"),
	Example: generateGetCmdExamplesForArtifact(artifactSequences, getTrimmedCmdLiteral(getSequenceCmdLiteral), "SampleSequence"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetSequenceCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getSequenceCmd)
	getSequenceCmd.Flags().StringVarP(&getSequenceCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getSequenceCmd.Flags().StringVarP(&getSequenceCmdFormat, "format", "", "", generateFormatFlagUsage(artifactSequences))
	getSequenceCmd.MarkFlagRequired("environment")
}

func handleGetSequenceCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getTrimmedCmdLiteral(getSequenceCmdLiteral))
	credentials.HandleMissingCredentials(getSequenceCmdEnvironment)
	if len(args) == 1 {
		var sequenceName = args[0]
		executeShowSequence(sequenceName)
	} else {
		executeListSequences()
	}
}

func executeListSequences() {

	sequenceList, err := impl.GetSequenceList(getSequenceCmdEnvironment)
	if err == nil {
		impl.PrintSequenceList(sequenceList, getSequenceCmdFormat)
	} else {
		printErrorForArtifactList(artifactSequences, err)
	}
}

func executeShowSequence(sequenceName string) {
	sequence, err := impl.GetSequence(getSequenceCmdEnvironment, sequenceName)
	if err == nil {
		impl.PrintSequenceDetails(sequence, getSequenceCmdFormat)
	} else {
		printErrorForArtifact(artifactSequences, sequenceName, err)
	}
}
