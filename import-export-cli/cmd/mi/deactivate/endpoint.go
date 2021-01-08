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

var deactivateEndpointCmdEnvironment string

const artifactEndpoint = "endpoint"
const deactivateEndpointCmdLiteral = "endpoint [endpoint-name]"

var deactivateEndpointCmd = &cobra.Command{
	Use:     deactivateEndpointCmdLiteral,
	Short:   generateDeactivateCmdShortDescForArtifact(artifactEndpoint),
	Long:    generateDeactivateCmdLongDescForArtifact(artifactEndpoint, "endpoint-name"),
	Example: generateDeactivateCmdExamplesForArtifact(artifactEndpoint, miUtils.GetTrimmedCmdLiteral(deactivateEndpointCmdLiteral), "TestEP"),
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleDeactivateEndpointCmdArguments(args)
	},
}

func init() {
	DeactivateCmd.AddCommand(deactivateEndpointCmd)
	setEnvFlag(deactivateEndpointCmd, &deactivateEndpointCmdEnvironment, artifactEndpoint)
}

func handleDeactivateEndpointCmdArguments(args []string) {
	printDeactivateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(deactivateEndpointCmdLiteral))
	credentials.HandleMissingCredentials(deactivateEndpointCmdEnvironment)
	executeDeactivateEndpoint(args[0])
}

func executeDeactivateEndpoint(endpointName string) {
	resp, err := impl.DeactivateEndpoint(deactivateEndpointCmdEnvironment, endpointName)
	if err != nil {
		printErrorForArtifact(artifactEndpoint, endpointName, err)
	} else {
		fmt.Println(resp)
	}
}
