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
)

const deleteAPICmdLiteral = "api"
const deleteAPICmdShortDesc = "Delete an API in Microgateway"
const deleteAPICmdLongDesc = `Delete an API by specifying name, version, host, username and optionally vhost
 by specifying the flags (--name (-n), --version (-v), --host (-c), --username (-u), and optionally --vhost (-t)`

var deleteAPICmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteAPICmdLiteral + `--host https://localhost:9095 -u admin
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n petstore -v 0.0.1 -c https://localhost:9095 -u admin -t www.pets.com 
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n "petstore VIP" -v 0.0.1 --host https://localhost:9095 -u admin`

var mgDeleteAPIResourcePath = "/apis/delete"

// DeleteAPICmd represents the apis command
var DeleteAPICmd = &cobra.Command{
	Use:     deleteAPICmdLiteral,
	Short:   deleteAPICmdShortDesc,
	Long:    deleteAPICmdLongDesc,
	Example: deleteAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAPICmdLiteral + " called")

		// handle auth
		fmt.Print("Enter Password: ")
		password, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			utils.HandleErrorAndExit("Error reading password", err)
		}
		authToken := base64.StdEncoding.EncodeToString([]byte(deleteAPIUsername + ":" + string(password)))

		//handle parameters
		queryParams := make(map[string]string)
		queryParams["apiName"] = deleteAPICmdAPIName
		queryParams["version"] = deleteAPICmdAPIVersion
		queryParams["vhost"] = deleteAPICmdAPIVHost
		err = mgImpl.DeleteAPI(authToken,
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

	_ = DeleteAPICmd.MarkFlagRequired("host")
	_ = DeleteAPICmd.MarkFlagRequired("name")
	_ = DeleteAPICmd.MarkFlagRequired("version")
	_ = DeleteAPICmd.MarkFlagRequired("username")
}
