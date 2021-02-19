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

package mg

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	undeployCmdLiteral   = "undeploy"
	undeployCmdShortDesc = "Undeploy an API in Microgateway"
	undeployCmdLongDesc  = "Undeploy an API in Microgateway by specifying name, version, host, username " +
		"and optionally vhost"
)

const undeployCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` --host https://localhost:9095 -n petstore -v 0.0.1 -u admin
` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n petstore -v 0.0.1 -c https://localhost:9095 -u admin --vhost www.pets.com 
` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n SwaggerPetstore -v 0.0.1 --host https://localhost:9095 -u admin -p admin` +

	"\n\nNote: The flags --name (-n), --version (-v), --host (-c), and --username (-u) are mandatory. " +
	"The password can be included via the flag --password (-p) or entered at the prompt."

// UndeployCmd represents the undeploy command
var UndeployCmd = &cobra.Command{
	Use:     undeployCmdLiteral,
	Short:   undeployCmdShortDesc,
	Long:    undeployCmdLongDesc,
	Example: undeployCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + undeployCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	MgCmd.AddCommand(UndeployCmd)
}
