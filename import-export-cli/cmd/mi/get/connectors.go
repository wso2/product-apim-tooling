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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getConnectorCmdEnvironment string
var getConnectorCmdFormat string

const getConnectorCmdLiteral = "connectors"
const getConnectorCmdShortDesc = "Get information about connectors deployed in a Micro Integrator"

const getConnectorCmdLongDesc = "List all the connectors deployed in a Micro Integrator in the environment specified by the flag --environment, -e"

var getConnectorCmdExamples = "To list all the connectors\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + getConnectorCmdLiteral + " -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var getConnectorCmd = &cobra.Command{
	Use:     getConnectorCmdLiteral,
	Short:   getConnectorCmdShortDesc,
	Long:    getConnectorCmdLongDesc,
	Example: getConnectorCmdExamples,
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetConnectorCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getConnectorCmd)
	setEnvFlag(getConnectorCmd, &getConnectorCmdEnvironment)
	setFormatFlag(getConnectorCmd, &getConnectorCmdFormat)
}

func handleGetConnectorCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getConnectorCmdLiteral)
	credentials.HandleMissingCredentials(getConnectorCmdEnvironment)
	executeListConnectors()
}

func executeListConnectors() {
	connectorList, err := impl.GetConnectorList(getConnectorCmdEnvironment)
	if err == nil {
		impl.PrintConnectorList(connectorList, getConnectorCmdFormat)
	} else {
		printErrorForArtifactList(getConnectorCmdLiteral, err)
	}
}
