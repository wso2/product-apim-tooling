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
	deleteAPICmdAPIName    string
	deleteAPICmdAPIVersion string
	deleteAPICmdAPIVHost   string
	deleteAPIUsername      string
	deleteAPIPassword      string
)

const deleteAPICmdLiteral = "api"
const deleteAPICmdShortDesc = "Delete an API in Microgateway"
const deleteAPICmdLongDesc = "Delete an API in Microgateway by specifying name, version, host, username " +
	"and optionally vhost"

var deleteAPICmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteAPICmdLiteral + ` --host https://localhost:9095 -u admin
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n petstore -v 0.0.1 -c https://localhost:9095 -u admin -t www.pets.com 
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n SwaggerPetstore -v 0.0.1 --host https://localhost:9095 -u admin -p admin` +

	"\n\nNote: The flags --name (-n), --version (-v), --host (-c), and --username (-u) are mandatory. " +
	"The password can be included via the flag --password (-p) or entered at the prompt."

var mgDeleteAPIResourcePath = "/api"

// DeleteAPICmd represents the apis command
var DeleteAPICmd = &cobra.Command{
	Use:     deleteAPICmdLiteral,
	Short:   deleteAPICmdShortDesc,
	Long:    deleteAPICmdLongDesc,
	Example: deleteAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAPICmdLiteral + " called")

		// handle auth
		if deleteAPIPassword == "" {
			fmt.Print("Enter Password: ")
			deleteAPIPasswordB, err := terminal.ReadPassword(0)
			deleteAPIPassword = string(deleteAPIPasswordB)
			fmt.Println()
			if err != nil {
				utils.HandleErrorAndExit("Error reading password", err)
			}
		}
		authToken := base64.StdEncoding.EncodeToString(
			[]byte(deleteAPIUsername + ":" + deleteAPIPassword))

		//handle parameters
		queryParams := make(map[string]string)
		queryParams["apiName"] = deleteAPICmdAPIName
		queryParams["version"] = deleteAPICmdAPIVersion
		queryParams["vhost"] = deleteAPICmdAPIVHost
		err := mgImpl.DeleteAPI(authToken,
			mgwAdapterHost+MgBasepath+mgDeleteAPIResourcePath,
			queryParams)
		if err != nil {
			utils.HandleErrorAndExit("Error deleting API", err)
		}
		fmt.Println("API deleted successfully!")
	},
}

func init() {
	DeleteCmd.AddCommand(DeleteAPICmd)

	DeleteAPICmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "The adapter host url with port")
	DeleteAPICmd.Flags().StringVarP(&deleteAPICmdAPIName, "name", "n", "", "API name")
	DeleteAPICmd.Flags().StringVarP(&deleteAPICmdAPIVersion, "version", "v", "", "API version")
	DeleteAPICmd.Flags().StringVarP(&deleteAPICmdAPIVHost, "vhost", "t", "", "Virtual host the API needs to be deleted from")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIUsername, "username", "u", "", "Username with delete permissions")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIPassword, "password", "p", "", "Password of the user (Can be provided at the prompt)")

	_ = DeleteAPICmd.MarkFlagRequired("host")
	_ = DeleteAPICmd.MarkFlagRequired("name")
	_ = DeleteAPICmd.MarkFlagRequired("version")
	_ = DeleteAPICmd.MarkFlagRequired("username")
}
