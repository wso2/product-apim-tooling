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

var getApplicationCmdEnvironment string
var getApplicationCmdFormat string

const artifactCompositeApps = "composite apps"
const getApplicationCmdLiteral = "composite-apps [app-name]"

var getApplicationCmd = &cobra.Command{
	Use:     getApplicationCmdLiteral,
	Short:   generateGetCmdShortDescForArtifact(artifactCompositeApps),
	Long:    generateGetCmdLongDescForArtifact(artifactCompositeApps, "app-name"),
	Example: generateGetCmdExamplesForArtifact(artifactCompositeApps, getTrimmedCmdLiteral(getApplicationCmdLiteral), "SampleApp"),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetApplicationCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getApplicationCmd)
	setEnvFlag(getApplicationCmd, &getApplicationCmdEnvironment)
	setFormatFlag(getApplicationCmd, &getApplicationCmdFormat)
}

func handleGetApplicationCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(getTrimmedCmdLiteral(getApplicationCmdLiteral))
	credentials.HandleMissingCredentials(getApplicationCmdEnvironment)
	if len(args) == 1 {
		var appName = args[0]
		executeShowCarbonApp(appName)
	} else {
		executeListCarbonApps()
	}
}

func executeListCarbonApps() {
	appList, err := impl.GetCompositeAppList(getApplicationCmdEnvironment)
	if err == nil {
		impl.PrintCompositeAppList(appList, getApplicationCmdFormat)
	} else {
		printErrorForArtifactList(artifactCompositeApps, err)
	}
}

func executeShowCarbonApp(appname string) {
	app, err := impl.GetCompositeApp(getApplicationCmdEnvironment, appname)
	if err == nil {
		impl.PrintCompositeAppDetails(app, getApplicationCmdFormat)
	} else {
		printErrorForArtifact(artifactCompositeApps, appname, err)
	}
}
