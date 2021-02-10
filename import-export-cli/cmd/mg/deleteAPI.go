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
	deleteApisCmdAPIName    string
	deleteApisCmdAPIVersion string
	deleteApisCmdAPIVHost   string
	deleteApisUsername      string
)

const deleteApisCmdLiteral = "apis"
const deleteApisCmdShortDesc = "Delete API"
const deleteApisCmdLongDesc = `Delete an API by specifying name, version and optionally vhost
 by specifying the flags (--name (-n), --version (-v), and optionally --vhost (-vh)`

var deleteApisCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteApisCmdLiteral + `-h https://localhost:9095 -u admin
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteApisCmdLiteral + ` -n petstore -v 0.0.1 -vh www.pets.com http -h https://localhost:9095 -u admin
  ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + deleteApisCmdLiteral + ` -n "petstore VIP" -v 0.0.1 -h https://localhost:9095 -u admin`

var mgDeleteAPIsResourcePath = "/apis/delete"

// DeleteApisCmd represents the apis command
var DeleteApisCmd = &cobra.Command{
	Use:     deleteApisCmdLiteral,
	Short:   deleteApisCmdShortDesc,
	Long:    deleteApisCmdLongDesc,
	Example: deleteApisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteApisCmdLiteral + " called")

		// handle auth
		fmt.Print("Enter Password: ")
		password, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			utils.HandleErrorAndExit("Error reading password", err)
		}
		authToken := base64.StdEncoding.EncodeToString([]byte(deleteApisUsername + ":" + string(password)))

		//handle parameters
		queryParams := make(map[string]string)
		queryParams["apiName"] = deleteApisCmdAPIName
		queryParams["version"] = deleteApisCmdAPIVersion
		queryParams["vhost"] = deleteApisCmdAPIVHost
		err = mgImpl.DeleteAPI(authToken,
			mgwAdapterHost+MgBasepath+mgDeleteAPIsResourcePath,
			queryParams)
		if err != nil {
			utils.HandleErrorAndExit("Error deleting API", err)
		}
		fmt.Println("API deleted successfully!.")
	},
}

func init() {
	DeleteCmd.AddCommand(DeleteApisCmd)

	DeleteApisCmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "The adapter host url with port")
	DeleteApisCmd.Flags().StringVarP(&deleteApisCmdAPIName, "name", "n", "", "API name")
	DeleteApisCmd.Flags().StringVarP(&deleteApisCmdAPIVersion, "version", "v", "", "API version")
	DeleteApisCmd.Flags().StringVarP(&deleteApisCmdAPIVHost, "vhost", "t", "", "Virtual host the API needs to be deleted from")
	DeleteApisCmd.Flags().StringVarP(&deleteApisUsername, "username", "u", "", "Username with delete permissions")

	_ = DeleteApisCmd.MarkFlagRequired("host")
	_ = DeleteApisCmd.MarkFlagRequired("name")
	_ = DeleteApisCmd.MarkFlagRequired("version")
	_ = DeleteApisCmd.MarkFlagRequired("username")
}
