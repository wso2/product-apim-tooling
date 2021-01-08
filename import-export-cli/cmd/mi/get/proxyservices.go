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

var getProxyServiceCmdEnvironment string
var getProxyServiceCmdFormat string

const artifactProxyServices = "proxy services"
const getProxyServiceCmdLiteral = "proxy-services [proxy-name]"

var getProxyServiceCmd = &cobra.Command{
	Use:     getProxyServiceCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactProxyServices),
	Long:    generateGetCmdLongDescForArtifact(artifactProxyServices, "proxy-name"),
	Example: generateGetCmdExamplesForArtifact(artifactProxyServices, miUtils.GetTrimmedCmdLiteral(getProxyServiceCmdLiteral), "SampleProxy"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetProxyServiceCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getProxyServiceCmd)
	setEnvFlag(getProxyServiceCmd, &getProxyServiceCmdEnvironment)
	setFormatFlag(getProxyServiceCmd, &getProxyServiceCmdFormat)
}

func handleGetProxyServiceCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getProxyServiceCmdLiteral))
	credentials.HandleMissingCredentials(getProxyServiceCmdEnvironment)
	if len(args) == 1 {
		var proxyServiceName = args[0]
		executeShowProxyService(proxyServiceName)
	} else {
		executeListProxyServices()
	}
}

func executeListProxyServices() {
	proxyList, err := impl.GetProxyServiceList(getProxyServiceCmdEnvironment)
	if err == nil {
		impl.PrintProxyServiceList(proxyList, getProxyServiceCmdFormat)
	} else {
		printErrorForArtifactList(artifactProxyServices, err)
	}
}

func executeShowProxyService(proxyName string) {
	proxyService, err := impl.GetProxyService(getProxyServiceCmdEnvironment, proxyName)
	if err == nil {
		impl.PrintProxyServiceDetails(proxyService, getProxyServiceCmdFormat)
	} else {
		printErrorForArtifact(artifactProxyServices, proxyName, err)
	}
}
