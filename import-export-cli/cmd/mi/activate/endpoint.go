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

var activateEndpointCmdEnvironment string

const artifactEndpoint = "endpoint"
const activateEndpointCmdLiteral = "endpoint [endpoint-name]"

var activateEndpointCmd = &cobra.Command{
	Use:     activateEndpointCmdLiteral,
	Short:   generateActivateCmdShortDescForArtifact(artifactEndpoint),
	Long:    generateActivateCmdLongDescForArtifact(artifactEndpoint, "endpoint-name"),
	Example: generateActivateCmdExamplesForArtifact(artifactEndpoint, miUtils.GetTrimmedCmdLiteral(activateEndpointCmdLiteral), "TestEP"),
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleActivateEndpointCmdArguments(args)
	},
}

func init() {
	ActivateCmd.AddCommand(activateEndpointCmd)
	setEnvFlag(activateEndpointCmd, &activateEndpointCmdEnvironment, artifactEndpoint)
}

func handleActivateEndpointCmdArguments(args []string) {
	printActivateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(activateEndpointCmdLiteral))
	credentials.HandleMissingCredentials(activateEndpointCmdEnvironment)
	executeActivateEndpoint(args[0])
}

func executeActivateEndpoint(endpointName string) {
	resp, err := impl.ActivateEndpoint(activateEndpointCmdEnvironment, endpointName)
	if err != nil {
		printErrorForArtifact(artifactEndpoint, endpointName, err)
	} else {
		fmt.Println(resp)
	}
}
