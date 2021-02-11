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

// Deploy command related usage Info
const DeployCmdLiteral = "deploy"
const deployCmdShortDesc = "Deploy an API/API Product in a gateway environment"

const deployCmdLongDesc = `Deploy an API/API Product available in the environment specified by flag (--environment, -e)
to the gateway specified by flag (--gateway, -g)`

const deployCmdExamples = utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin --rev 1 -g Label1 -e dev
` + utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 6 -g Label1 Label2 Label3 -e production
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPIProductCmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 2 -r admin -g Label1 -e production --hide-on-devportal`

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:     DeployCmdLiteral,
	Short:   deployCmdShortDesc,
	Long:    deployCmdLongDesc,
	Example: deployCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + DeployCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	RootCmd.AddCommand(DeployCmd)
}
