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

var getDataServiceCmdEnvironment string
var getDataServiceCmdFormat string

const artifactDataServices = "data services"
const getDataServiceCmdLiteral = "data-services [dataservice-name]"

var getDataServiceCmd = &cobra.Command{
	Use:     getDataServiceCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactDataServices),
	Long:    generateGetCmdLongDescForArtifact(artifactDataServices, "dataservice-name"),
	Example: generateGetCmdExamplesForArtifact(artifactDataServices, miUtils.GetTrimmedCmdLiteral(getDataServiceCmdLiteral), "SampleDataService"),
	Args:    cobra.MaximumNArgs(1),
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handleGetDataServiceCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getDataServiceCmd)
	setEnvFlag(getDataServiceCmd, &getDataServiceCmdEnvironment)
	setFormatFlag(getDataServiceCmd, &getDataServiceCmdFormat)
}

func handleGetDataServiceCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getDataServiceCmdLiteral))
	credentials.HandleMissingCredentials(getDataServiceCmdEnvironment)
	if len(args) == 1 {
		var dataServiceName = args[0]
		executeShowDataService(dataServiceName)
	} else {
		executeListDataServices()
	}
}

func executeListDataServices() {
	dataServiceList, err := impl.GetDataServiceList(getDataServiceCmdEnvironment)
	if err == nil {
		impl.PrintDataServiceList(dataServiceList, getDataServiceCmdFormat)
	} else {
		printErrorForArtifactList(artifactDataServices, err)
	}
}

func executeShowDataService(dataserviceName string) {
	dataservice, err := impl.GetDataService(getDataServiceCmdEnvironment, dataserviceName)
	if err == nil {
		impl.PrintDataServiceDetails(dataservice, getDataServiceCmdFormat)
	} else {
		printErrorForArtifact(artifactDataServices, dataserviceName, err)
	}
}
