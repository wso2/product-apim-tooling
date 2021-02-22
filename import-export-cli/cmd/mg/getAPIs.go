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
	getAPIsQuery    string
	getAPIsLimit    string
	getAPIsUsername string
	getAPIsPassword string
)

const getAPIsCmdShortDesc = "List APIs in Microgateway"
const getAPIsCmdLongDesc = "Display a list of all the APIs in Microgateway or a set of APIs " +
	"with a limit set or filtered by apiType"

var getAPIsCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + apisCmdLiteral + ` --host https://localhost:9095 -u admin
` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + apisCmdLiteral + ` -q type:http --host https://localhost:9095 -u admin -l 100
` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + apisCmdLiteral + ` -q type:ws --host https://localhost:9095 -u admin

Note: The flags --host (-c), --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.`

var mgGetAPIsResourcePath = "/apis"

// GetAPIsCmd represents the apis command
var GetAPIsCmd = &cobra.Command{
	Use:     apisCmdLiteral,
	Short:   getAPIsCmdShortDesc,
	Long:    getAPIsCmdLongDesc,
	Example: getAPIsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apisCmdLiteral + " called")

		// handle auth
		if getAPIsPassword == "" {
			fmt.Print("Enter Password: ")
			getAPIsPasswordB, err := terminal.ReadPassword(0)
			getAPIsPassword = string(getAPIsPasswordB)
			fmt.Println()
			if err != nil {
				utils.HandleErrorAndExit("Error reading password", err)
			}
		}

		authToken := base64.StdEncoding.EncodeToString([]byte(
			getAPIsUsername + ":" + getAPIsPassword))

		//handle parameters
		if getAPIsLimit == "" {
			getAPIsLimit = strconv.Itoa(utils.DefaultApisDisplayLimit)
			fmt.Print("Limit flag not set. Set to default: " + getAPIsLimit + "\n")
		}
		queryParams := make(map[string]string)
		queryParams["limit"] = getAPIsLimit
		queryParams["query"] = getAPIsQuery
		total, count, apis, err := mgImpl.GetAPIsList(authToken,
			mgwAdapterHost+MgBasepath+mgGetAPIsResourcePath,
			queryParams)
		if err != nil {
			utils.HandleErrorAndExit("Error while retrieving or processing received APIs", err)
		}
		fmt.Fprintf(os.Stderr, "APIs total: %v received: %v\n", total, count)
		mgImpl.PrintAPIs(apis)
	},
}

func init() {
	GetCmd.AddCommand(GetAPIsCmd)

	GetAPIsCmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "The adapter host url with port")
	GetAPIsCmd.Flags().StringVarP(&getAPIsQuery, "query", "q", "", "Query to filter the APIs")
	GetAPIsCmd.Flags().StringVarP(&getAPIsLimit, "limit", "l", "", "Maximum number of APIs to return")
	GetAPIsCmd.Flags().StringVarP(&getAPIsUsername, "username", "u", "", "The username")
	GetAPIsCmd.Flags().StringVarP(&getAPIsPassword, "password", "p", "", "Password of the user (Can be provided at the prompt)")

	_ = GetAPIsCmd.MarkFlagRequired("host")
	_ = GetAPIsCmd.MarkFlagRequired("username")
}
