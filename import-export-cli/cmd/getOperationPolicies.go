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

package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getOperationPoliciesCmdEnvironment string
var getOperationPoliciesCmdFormat string
var getOperationPoliciesCmdQuery []string

// GetOperationPoliciesCmdLiteral related info
const GetOperationPoliciesCmdLiteral = "operation"
const getOperationPoliciesCmdShortDesc = "Display a list of Operation Policies in an environment"

const getOperationPoliciesCmdLongDesc = `Display a list of Operation Policies in the environment specified by the flag --environment, -e`

var getOperationPoliciesCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdExamples + ` ` + GetOperationPoliciesCmdExample + ` -e dev
 ` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetOperationPoliciesCmdExample + ` -e prod -q gateway:choreo
 NOTE: The flag (--environment (-e)) is mandatory`

// getOperationPoliciesCmd represents the get policies rate-limiting command
var getOperationPoliciesCmd = &cobra.Command{
	Use:     GetOperationPoliciesCmdLiteral,
	Short:   getOperationPoliciesCmdShortDesc,
	Long:    getOperationPoliciesCmdLongDesc,
	Example: getOperationPoliciesCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetOperationPoliciesCmdLiteral + " called")
		cred, err := GetCredentials(getOperationPoliciesCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetOperationPoliciesCmd(cred)
	},
}

func executeGetOperationPoliciesCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, getOperationPoliciesCmdEnvironment)
	fmt.Println("Access Token: ", accessToken)
	if preCommandErr == nil {
		resp, err := impl.GetOperationPolicyListFromEnv(accessToken, getOperationPoliciesCmdEnvironment,
			strings.Join(getOperationPoliciesCmdQuery, queryParamSeparator))
		if err != nil {
			utils.HandleErrorAndExit("Error while getting operation policies", err)
		}
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

		if resp.StatusCode() == http.StatusOK {
			impl.PrintOperationPolicies(resp, getOperationPoliciesCmdFormat)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error getting Operation Policies:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens while getting Operation Policies:" + preCommandErr.Error())
	}
}

func init() {
	GetPoliciesCmd.AddCommand(getOperationPoliciesCmd)
	getOperationPoliciesCmd.Flags().StringVarP(&getOperationPoliciesCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getOperationPoliciesCmd.Flags().StringSliceVarP(&getOperationPoliciesCmdQuery, "query", "q",
		[]string{}, "Query pattern")
	getOperationPoliciesCmd.Flags().StringVarP(&getOperationPoliciesCmdFormat, "format", "", "", "Pretty-print throttle policies "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getOperationPoliciesCmd.MarkFlagRequired("environment")
}
