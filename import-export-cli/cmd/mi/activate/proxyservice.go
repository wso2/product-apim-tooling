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

var activateProxyCmdEnvironment string

const artifactProxy = "proxy service"
const activateProxyCmdLiteral = "proxy-service [proxy-name]"

var activateProxyCmd = &cobra.Command{
	Use:     activateProxyCmdLiteral,
	Short:   generateActivateCmdShortDescForArtifact(artifactProxy),
	Long:    generateActivateCmdLongDescForArtifact(artifactProxy, "proxy-name"),
	Example: generateActivateCmdExamplesForArtifact(artifactProxy, miUtils.GetTrimmedCmdLiteral(activateProxyCmdLiteral), "SampleProxy"),
	Args:    cobra.ExactArgs(1),
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handleActivateProxyCmdArguments(args)
	},
}

func init() {
	ActivateCmd.AddCommand(activateProxyCmd)
	setEnvFlag(activateProxyCmd, &activateProxyCmdEnvironment, artifactProxy)
}

func handleActivateProxyCmdArguments(args []string) {
	printActivateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(activateProxyCmdLiteral))
	credentials.HandleMissingCredentials(activateProxyCmdEnvironment)
	executeActivateProxy(args[0])
}

func executeActivateProxy(proxyName string) {
	resp, err := impl.ActivateProxy(activateProxyCmdEnvironment, proxyName)
	if err != nil {
		printErrorForArtifact(artifactProxy, proxyName, err)
	} else {
		fmt.Println(resp)
	}
}
