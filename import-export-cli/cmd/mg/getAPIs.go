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
	getApisCmdAPIType string
	getApisCmdLimit   string
	getApisUsername   string
)

const getApisCmdLiteral = "apis"
const getApisCmdShortDesc = "List APIs in Microgateway"
const getApisCmdLongDesc = `Display a list of all the APIs or 
a set of APIs with a limit or filtered by apiType using the flags --limit (-l), --type (-t). 
Note: The flags --host (-c), --username (-u) are mandatory`

var getApisCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + getApisCmdLiteral + `--host https://localhost:9095 -u admin
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + getApisCmdLiteral + ` -t http --host https://localhost:9095 -u admin -l 100
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + getApisCmdLiteral + ` -t ws --host https://localhost:9095 -u admin`

var mgGetAPIsResourcePath = "/apis"

// GetApisCmd represents the apis command
var GetApisCmd = &cobra.Command{
	Use:     getApisCmdLiteral,
	Short:   getApisCmdShortDesc,
	Long:    getApisCmdLongDesc,
	Example: getApisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + getApisCmdLiteral + " called")

		// handle auth
		fmt.Print("Enter Password: ")
		password, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			utils.HandleErrorAndExit("Error reading password", err)
		}
		authToken := base64.StdEncoding.EncodeToString([]byte(getApisUsername + ":" + string(password)))

		//handle parameters
		if getApisCmdLimit == "" {
			getApisCmdLimit = strconv.Itoa(utils.DefaultApisDisplayLimit)
			fmt.Print("Limit flag not set. Set to default: " + getApisCmdLimit + "\n")
		}
		queryParams := make(map[string]string)
		queryParams["limit"] = getApisCmdLimit
		queryParams["apiType"] = getApisCmdAPIType
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
	GetCmd.AddCommand(GetApisCmd)

	GetApisCmd.Flags().StringVarP(&mgwAdapterHost, "host", "c", "", "The adapter host url with port")
	GetApisCmd.Flags().StringVarP(&getApisCmdAPIType, "type", "t", "", "API type to filter the APIs")
	GetApisCmd.Flags().StringVarP(&getApisCmdLimit, "limit", "l", "", "Maximum number of APIs to return")
	GetApisCmd.Flags().StringVarP(&getApisUsername, "username", "u", "", "The username")

	_ = GetApisCmd.MarkFlagRequired("host")
	_ = GetApisCmd.MarkFlagRequired("username")
}
