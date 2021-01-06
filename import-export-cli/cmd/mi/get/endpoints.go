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

var getEndpointCmdEnvironment string
var getEndpointCmdFormat string

const artifactEndpoints = "endpoints"
const getEndpointCmdLiteral = "endpoints [endpoint-name]"

var getEndpointCmd = &cobra.Command{
	Use:     getEndpointCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactEndpoints),
	Long:    generateGetCmdLongDescForArtifact(artifactEndpoints, "endpoint-name"),
	Example: generateGetCmdExamplesForArtifact(artifactEndpoints, getTrimmedCmdLiteral(getEndpointCmdLiteral), "SampleEndpoint"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetEndpointCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getEndpointCmd)
	setEnvFlag(getEndpointCmd, &getEndpointCmdEnvironment)
	setFormatFlag(getEndpointCmd, &getEndpointCmdFormat)
}

func handleGetEndpointCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getTrimmedCmdLiteral(getEndpointCmdLiteral))
	credentials.HandleMissingCredentials(getEndpointCmdEnvironment)
	if len(args) == 1 {
		var EndpointName = args[0]
		executeShowEndpoint(EndpointName)
	} else {
		executeListEndpoints()
	}
}

func executeListEndpoints() {
	epList, err := impl.GetEndpointList(getEndpointCmdEnvironment)
	if err == nil {
		impl.PrintEndpointList(epList, getEndpointCmdFormat)
	} else {
		printErrorForArtifactList(artifactEndpoints, err)
	}
}

func executeShowEndpoint(epName string) {
	endpoint, err := impl.GetEndpoint(getEndpointCmdEnvironment, epName)
	if err == nil {
		impl.PrintEndpointDetails(endpoint, getEndpointCmdFormat)
	} else {
		printErrorForArtifact(artifactEndpoints, epName, err)
	}
}
