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
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	mgImpl "github.com/wso2/product-apim-tooling/import-export-cli/impl/mg"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	undeployAPICmdAPIName    string
	undeployAPICmdAPIVersion string
	undeployAPICmdAPIVHost   string
	undeployAPIUsername      string
	undeployAPIPassword      string
)

const (
	undeployAPICmdShortDesc = "Undeploy an API in Microgateway"
	undeployAPICmdLongDesc  = "Undeploy an API in Microgateway by specifying name, version, host, username " +
		"and optionally vhost"
)

var undeployAPICmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` --host https://localhost:9095 -n petstore -v 0.0.1 -u admin
   ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n petstore -v 0.0.1 -c https://localhost:9095 -u admin --vhost www.pets.com 
   ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + undeployCmdLiteral + ` ` + apiCmdLiteral + ` -n SwaggerPetstore -v 0.0.1 --host https://localhost:9095 -u admin -p admin` +

	"\n\nNote: The flags --name (-n), --version (-v), --host (-c), and --username (-u) are mandatory. " +
	"The password can be included via the flag --password (-p) or entered at the prompt."

var mgUndeployAPIResourcePath = "/apis"

// UndeployAPICmd represents the undeploy api command
var UndeployAPICmd = &cobra.Command{
	Use:     apiCmdLiteral,
	Short:   undeployAPICmdShortDesc,
	Long:    undeployAPICmdLongDesc,
	Example: undeployAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + undeployCmdLiteral + " called")

		// handle auth
		if undeployAPIPassword == "" {
			fmt.Print("Enter Password: ")
			undeployAPIPasswordB, err := terminal.ReadPassword(0)
			undeployAPIPassword = string(undeployAPIPasswordB)
			fmt.Println()
			if err != nil {
				utils.HandleErrorAndExit("Error reading password", err)
			}
		}
		authToken := base64.StdEncoding.EncodeToString(
			[]byte(undeployAPIUsername + ":" + undeployAPIPassword))

		//handle parameters
		queryParams := make(map[string]string)
		queryParams["apiName"] = undeployAPICmdAPIName
		queryParams["version"] = undeployAPICmdAPIVersion
		queryParams["vhost"] = undeployAPICmdAPIVHost
		err := mgImpl.UndeployAPI(authToken,
			mgwAdapterHost+MgBasepath+mgUndeployAPIResourcePath,
			queryParams)
		if err != nil {
			utils.HandleErrorAndExit("Error deleting API", err)
		}
		fmt.Println("API undeployed from microgateway successfully!")
	},
}

func init() {
	UndeployCmd.AddCommand(UndeployAPICmd)

	UndeployAPICmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "The adapter host url with port")
	UndeployAPICmd.Flags().StringVarP(&undeployAPICmdAPIName, "name", "n", "", "API name")
	UndeployAPICmd.Flags().StringVarP(&undeployAPICmdAPIVersion, "version", "v", "", "API version")
	UndeployAPICmd.Flags().StringVarP(&undeployAPICmdAPIVHost, "vhost", "t", "", "Virtual host the API needs to be undeployed from")
	UndeployAPICmd.Flags().StringVarP(&undeployAPIUsername, "username", "u", "", "Username with undeploy permissions")
	UndeployAPICmd.Flags().StringVarP(&undeployAPIPassword, "password", "p", "", "Password of the user (Can be provided at the prompt)")

	_ = UndeployAPICmd.MarkFlagRequired("host")
	_ = UndeployAPICmd.MarkFlagRequired("name")
	_ = UndeployAPICmd.MarkFlagRequired("version")
	_ = UndeployAPICmd.MarkFlagRequired("username")
}
