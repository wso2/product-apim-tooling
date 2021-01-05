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

var getInboundEndpointCmdEnvironment string
var getInboundEndpointCmdFormat string

const artifactInboundEndpoints = "inbound endpoints"
const getInboundEndpointCmdLiteral = "inbound-endpoints [inbound-name]"

var getInboundEndpointCmd = &cobra.Command{
	Use:     getInboundEndpointCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactInboundEndpoints),
	Long:    generateGetCmdLongDescForArtifact(artifactInboundEndpoints, "inbound-name"),
	Example: generateGetCmdExamplesForArtifact(artifactInboundEndpoints, getTrimmedCmdLiteral(getInboundEndpointCmdLiteral), "SampleInboundEndpoint"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetInboundEndpointCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getInboundEndpointCmd)
	setEnvFlag(getInboundEndpointCmd, &getInboundEndpointCmdEnvironment)
	setFormatFlag(getInboundEndpointCmd, &getInboundEndpointCmdFormat)
}

func handleGetInboundEndpointCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getTrimmedCmdLiteral(getInboundEndpointCmdLiteral))
	credentials.HandleMissingCredentials(getInboundEndpointCmdEnvironment)
	if len(args) == 1 {
		var inboundEndpointName = args[0]
		executeShowInboundEndpoint(inboundEndpointName)
	} else {
		executeListInboundEndpoints()
	}
}

func executeListInboundEndpoints() {
	inboundEpList, err := impl.GetInboundEndpointList(getInboundEndpointCmdEnvironment)
	if err == nil {
		impl.PrintInboundEndpointList(inboundEpList, getInboundEndpointCmdFormat)
	} else {
		printErrorForArtifactList(artifactInboundEndpoints, err)
	}
}

func executeShowInboundEndpoint(inboundEpName string) {
	inboundEndpoint, err := impl.GetInboundEndpoint(getInboundEndpointCmdEnvironment, inboundEpName)
	if err == nil {
		impl.PrintInboundEndpointDetails(inboundEndpoint, getInboundEndpointCmdFormat)
	} else {
		printErrorForArtifact(artifactInboundEndpoints, inboundEpName, err)
	}
}
