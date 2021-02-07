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
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/mg/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	listApisCmdQuery string
	listApisCmdLimit string
	mgwAdapterHost   string
	listApisUsername string
)

const listApisCmdLiteral = "apis"
const listApisCmdShortDesc = "Display a list of all APIs or a filtered set of APIs"
const listApisCmdLongDesc = `Display a list of all APIs or filtered by a query specified by the flag --query, -q`

var listApisCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + listApisCmdLiteral + `-h https://localhost:9095 -u admin
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + listApisCmdLiteral + ` -q type:http -h https://localhost:9095 -u admin -l 100
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + listApisCmdLiteral + ` -q type:ws -h https://localhost:9095 -u admin`

// ListApisCmd represents the apis command
var ListApisCmd = &cobra.Command{
	Use:     listApisCmdLiteral,
	Short:   listApisCmdShortDesc,
	Long:    listApisCmdLongDesc,
	Example: listApisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listApisCmdLiteral + " called")
		fmt.Print("Enter Password: ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			utils.HandleErrorAndExit("Error reading password", err)
		}
		authToken := base64.StdEncoding.EncodeToString([]byte(listApisUsername + ":" + string(password)))
		if listApisCmdLimit == "" {
			listApisCmdLimit = strconv.Itoa(utils.DefaultApisDisplayLimit)
			fmt.Print("Limit flag not set. Maximum APIs to retrieve set to  :" + listApisCmdLimit)

		}
		count, apis, err := impl.GetAPIList(authToken, mgwAdapterHost, listApisCmdQuery, listApisCmdLimit)
		if err != nil {
			utils.HandleErrorAndExit("Error retrieving APIs", err)
		}
		fmt.Print("APIs received: " + count)
		impl.PrintAPIs(apis)
	},
}

func init() {
	ListCmd.AddCommand(ListApisCmd)

	ListApisCmd.Flags().StringVarP(&mgwAdapterHost, "host", "h", "", "The adapter host url with port")
	ListApisCmd.Flags().StringVarP(&listApisCmdQuery, "query", "q", "", "Query to filter the APIs")
	ListApisCmd.Flags().StringVarP(&listApisCmdLimit, "limit", "l", "", "Maximum number of apis to return")
	ListApisCmd.Flags().StringVarP(&listApisUsername, "username", "u", "", "The username")

	_ = MgDeployCmd.MarkFlagRequired("host")
	_ = MgDeployCmd.MarkFlagRequired("username")
}
