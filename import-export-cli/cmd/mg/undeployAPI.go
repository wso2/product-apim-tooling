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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	mgImpl "github.com/wso2/product-apim-tooling/import-export-cli/impl/mg"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	undeployAPICmdAPIName       string
	undeployAPICmdAPIVersion    string
	undeployAPICmdAPIVHost      string
	undeployAPICmdAPIGatewayEnv string
	undeployAPIEnv              string
)

const gatewayNameSeparator = ":"

const (
	undeployAPICmdShortDesc = "Undeploy an API in Microgateway"
	undeployAPICmdLongDesc  = "Undeploy an API in Microgateway by specifying name, version, environment, username " +
		"and optionally vhost"
)

var undeployAPICmdExamples = `  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` --environment dev -n petstore -v 0.0.1
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n petstore -v 0.0.1 -e dev --vhost www.pets.com 
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n petstore -v 0.0.1 -e dev --gateway-env "Production and Sandbox" Label1 Label2
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n petstore -v 0.0.1 -e dev --vhost www.pets.com --gateway-env "Production and Sandbox" 
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n SwaggerPetstore -v 0.0.1 --environment dev

Note: The flags --name (-n), --version (-v), --environment (-e) are mandatory.
If the flag (--gateway-env (-g)) is not provided, API will be undeployed from all deployed gateway environments.
If the flag (--vhost (-t)) is not provided, API with all VHosts will be undeployed.
The user needs to be logged in to use this command.`

// UndeployAPICmd represents the undeploy api command
var UndeployAPICmd = &cobra.Command{
	Use:     apiCmdLiteral,
	Short:   undeployAPICmdShortDesc,
	Long:    undeployAPICmdLongDesc,
	Example: undeployAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + undeployCmdLiteral + " called")

		//handle parameters
		queryParams := make(map[string]string)
		queryParams["apiName"] = undeployAPICmdAPIName
		queryParams["version"] = undeployAPICmdAPIVersion
		queryParams["vhost"] = undeployAPICmdAPIVHost
		// Since other flags does not use args[], gateway-env flag will own all the args
		gatewayNames := append(args, undeployAPICmdAPIGatewayEnv)
		queryParams["environments"] = strings.Join(gatewayNames, gatewayNameSeparator)
		err := mgImpl.UndeployAPI(undeployAPIEnv, queryParams)
		if err != nil {
			utils.HandleErrorAndExit("Error undeploying API", err)
		}
		fmt.Println("API undeployed from microgateway successfully!")
	},
}

func init() {
	UndeployCmd.AddCommand(UndeployAPICmd)

	UndeployAPICmd.Flags().StringVarP(&undeployAPIEnv, "environment", "e", "", "Microgateway adapter environment to be undeployed from")
	UndeployAPICmd.Flags().StringVarP(&undeployAPICmdAPIName, "name", "n", "", "API name")
	UndeployAPICmd.Flags().StringVarP(&undeployAPICmdAPIVersion, "version", "v", "", "API version")
	UndeployAPICmd.Flags().StringVarP(&undeployAPICmdAPIVHost, "vhost", "t", "", "Virtual host the API needs to be undeployed from")
	UndeployAPICmd.Flags().StringVarP(&undeployAPICmdAPIGatewayEnv, "gateway-env", "g", "", "Gateway environment the API needs to be undeployed from")

	_ = UndeployAPICmd.MarkFlagRequired("environment")
	_ = UndeployAPICmd.MarkFlagRequired("name")
	_ = UndeployAPICmd.MarkFlagRequired("version")
}
