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

var deactivateProxyCmdEnvironment string

const artifactProxy = "proxy service"
const deactivateProxyCmdLiteral = "proxy-service [proxy-name]"

var deactivateProxyCmd = &cobra.Command{
	Use:     deactivateProxyCmdLiteral,
	Short:   generateDeactivateCmdShortDescForArtifact(artifactProxy),
	Long:    generateDeactivateCmdLongDescForArtifact(artifactProxy, "proxy-name"),
	Example: generateDeactivateCmdExamplesForArtifact(artifactProxy, miUtils.GetTrimmedCmdLiteral(deactivateProxyCmdLiteral), "SampleProxy"),
	Args:    cobra.ExactArgs(1),
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handleDeactivateProxyCmdArguments(args)
	},
}

func init() {
	DeactivateCmd.AddCommand(deactivateProxyCmd)
	setEnvFlag(deactivateProxyCmd, &deactivateProxyCmdEnvironment, artifactProxy)
}

func handleDeactivateProxyCmdArguments(args []string) {
	printDeactivateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(deactivateProxyCmdLiteral))
	credentials.HandleMissingCredentials(deactivateProxyCmdEnvironment)
	executeDeactivateProxy(args[0])
}

func executeDeactivateProxy(proxyName string) {
	resp, err := impl.DeactivateProxy(deactivateProxyCmdEnvironment, proxyName)
	if err != nil {
		printErrorForArtifact(artifactProxy, proxyName, err)
	} else {
		fmt.Println(resp)
	}
}
