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
	"os"
	"strconv"

	"github.com/spf13/cobra"
	mgImpl "github.com/wso2/product-apim-tooling/import-export-cli/impl/mg"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	listApisCmdAPIType string
	listApisCmdLimit   string
	mgwAdapterHost     string
	listApisUsername   string
)

const listApisCmdLiteral = "apis"
const listApisCmdShortDesc = "Display a list of all APIs or a filtered set of APIs"
const listApisCmdLongDesc = `Display a list of all APIs or filtered by apiType using the flag --type, -t`

var listApisCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + listApisCmdLiteral + `-h https://localhost:9095 -u admin
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + listApisCmdLiteral + ` -t http -h https://localhost:9095 -u admin -l 100
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + listApisCmdLiteral + ` -t ws -h https://localhost:9095 -u admin`

var mgListAPIsResourcePath = "/apis"

// ListApisCmd represents the apis command
var ListApisCmd = &cobra.Command{
	Use:     listApisCmdLiteral,
	Short:   listApisCmdShortDesc,
	Long:    listApisCmdLongDesc,
	Example: listApisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listApisCmdLiteral + " called")

		// handle auth
		fmt.Print("Enter Password: ")
		password, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			utils.HandleErrorAndExit("Error reading password", err)
		}
		authToken := base64.StdEncoding.EncodeToString([]byte(listApisUsername + ":" + string(password)))

		//handle parameters
		if listApisCmdLimit == "" {
			listApisCmdLimit = strconv.Itoa(utils.DefaultApisDisplayLimit)
			fmt.Print("Limit flag not set. Set to default: " + listApisCmdLimit + "\n")
		}
		queryParams := make(map[string]string)
		queryParams["limit"] = listApisCmdLimit
		queryParams["apiType"] = listApisCmdAPIType
		total, count, apis, err := mgImpl.GetAPIList(authToken,
			mgwAdapterHost+MgBasepath+mgListAPIsResourcePath,
			queryParams)
		if err != nil {
			utils.HandleErrorAndExit("Error retrieving APIs", err)
		}
		fmt.Fprintf(os.Stderr, "APIs total: %v received: %v\n", total, count)
		mgImpl.PrintAPIs(apis)
	},
}

func init() {
	ListCmd.AddCommand(ListApisCmd)

	ListApisCmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "The adapter host url with port")
	ListApisCmd.Flags().StringVarP(&listApisCmdAPIType, "type", "t", "", "API type to filter the APIs")
	ListApisCmd.Flags().StringVarP(&listApisCmdLimit, "limit", "l", "", "Maximum number of apis to return")
	ListApisCmd.Flags().StringVarP(&listApisUsername, "username", "u", "", "The username")

	_ = ListApisCmd.MarkFlagRequired("host")
	_ = ListApisCmd.MarkFlagRequired("username")
}
