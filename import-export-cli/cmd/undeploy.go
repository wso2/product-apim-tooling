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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// Undeploy command related usage Info
const UndeployCmdLiteral = "undeploy"
const undeployCmdShortDesc = "Undeploy an API/MCP Server/API Product revision from a gateway environment"

const undeployCmdLongDesc = `Undeploy an API/MCP Server/API Product revision available in the environment specified by flag (--environment, -e) from the gateway specified by flag (--gateway, -g)`

const undeployCmdExamples = utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin --rev 1 -g Label1 -g Label2 -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n PizzaAPI -v 1.0.0 --rev 2 --all-gateways -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployMCPServerCmdLiteral + ` -n WeatherMCPServer -v 1.0.0 --rev 1 -g Label1 -g Label2 -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPIProductCmdLiteral + ` -n LeasingAPIProduct -v 1.0.0 --rev 3 -g Label1 -e dev`

// UndeployCmd represents the undeploy command
var UndeployCmd = &cobra.Command{
	Use:     UndeployCmdLiteral,
	Short:   undeployCmdShortDesc,
	Long:    undeployCmdLongDesc,
	Example: undeployCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + UndeployCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(UndeployCmd)
}
