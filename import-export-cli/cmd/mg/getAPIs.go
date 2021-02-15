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
	getAPIsCmdAPIType string
	getAPIsCmdLimit   string
	getAPIsUsername   string
	getAPIsPassword   string
)

const getAPIsCmdLiteral = "apis"
const getAPIsCmdShortDesc = "List APIs in Microgateway"
const getAPIsCmdLongDesc = `Display a list of all the APIs or 
a set of APIs with a limit or filtered by apiType using the flags --limit (-l), --type (-t). 
Note: The flags --host (-c), --username (-u) are mandatory. The password can be included 
via the flag --password (-p) or entered at the prompt.`

var getAPIsCmdExamples = utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + getAPIsCmdLiteral + ` --host https://localhost:9095 -u admin
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + getAPIsCmdLiteral + ` -t http --host https://localhost:9095 -u admin -l 100
 ` + utils.ProjectName + ` ` + mgCmdLiteral + ` ` + getCmdLiteral + ` ` + getAPIsCmdLiteral + ` -t ws --host https://localhost:9095 -u admin`

var mgGetAPIsResourcePath = "/apis"

// GetAPIsCmd represents the apis command
var GetAPIsCmd = &cobra.Command{
	Use:     getAPIsCmdLiteral,
	Short:   getAPIsCmdShortDesc,
	Long:    getAPIsCmdLongDesc,
	Example: getAPIsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + getAPIsCmdLiteral + " called")

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
		if getAPIsCmdLimit == "" {
			getAPIsCmdLimit = strconv.Itoa(utils.DefaultApisDisplayLimit)
			fmt.Print("Limit flag not set. Set to default: " + getAPIsCmdLimit + "\n")
		}
		queryParams := make(map[string]string)
		queryParams["limit"] = getAPIsCmdLimit
		queryParams["apiType"] = getAPIsCmdAPIType
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
	GetAPIsCmd.Flags().StringVarP(&getAPIsCmdAPIType, "type", "t", "", "API type to filter the APIs")
	GetAPIsCmd.Flags().StringVarP(&getAPIsCmdLimit, "limit", "l", "", "Maximum number of APIs to return")
	GetAPIsCmd.Flags().StringVarP(&getAPIsUsername, "username", "u", "", "The username")
	GetAPIsCmd.Flags().StringVarP(&getAPIsPassword, "password", "p", "", "Password of the user (Can be provided at the prompt)")

	_ = GetAPIsCmd.MarkFlagRequired("host")
	_ = GetAPIsCmd.MarkFlagRequired("username")
}
