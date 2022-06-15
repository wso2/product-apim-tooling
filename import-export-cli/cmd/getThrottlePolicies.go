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
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getThrottlePoliciesCmdEnvironment string
var getThrottlePoliciesCmdFormat string
var getThrottlePoliciesCmdQuery []string

// GetThrottlePoliciesCmdLiteral related info
const GetThrottlePoliciesCmdLiteral = "rate-limiting"
const getThrottlePoliciesCmdShortDesc = "Display a list of APIs in an environment"

const getThrottlePoliciesCmdLongDesc = `Display a list of Throttling Policies in the environment specified by the flag --environment, -e`

var getThrottlePoliciesCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdExamples + ` ` + GetThrottlePoliciesCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetThrottlePoliciesCmdLiteral + ` -e prod -q type:api
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetThrottlePoliciesCmdLiteral + ` -e prod -q type:sub
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetThrottlePoliciesCmdLiteral + ` -e staging -q type:global
NOTE: The flag (--environment (-e)) is mandatory`

// getThrottlePoliciesCmd represents the get policies rate-limiting command
var getThrottlePoliciesCmd = &cobra.Command{
	Use:     GetThrottlePoliciesCmdLiteral,
	Short:   getThrottlePoliciesCmdShortDesc,
	Long:    getThrottlePoliciesCmdLongDesc,
	Example: getThrottlePoliciesCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetThrottlePoliciesCmdLiteral + " called")
		cred, err := GetCredentials(getThrottlePoliciesCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetThrottlePoliciesCmd(cred)
	},
}

func executeGetThrottlePoliciesCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, getThrottlePoliciesCmdEnvironment)
	if preCommandErr == nil {
		resp, err := impl.GetThrottlePolicyListFromEnv(accessToken, getThrottlePoliciesCmdEnvironment,
			strings.Join(getThrottlePoliciesCmdQuery, queryParamSeparator))
		if err != nil {
			utils.HandleErrorAndExit("Error while getting throttling policies", err)
		}
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

		if resp.StatusCode() == http.StatusOK {
			impl.PrintThrottlePolicies(resp, getThrottlePoliciesCmdFormat)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error getting Throttling Policies:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens while getting Throttling Policies:" + preCommandErr.Error())
	}
}

func init() {
	GetPoliciesCmd.AddCommand(getThrottlePoliciesCmd)
	getThrottlePoliciesCmd.Flags().StringVarP(&getThrottlePoliciesCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	getThrottlePoliciesCmd.Flags().StringSliceVarP(&getThrottlePoliciesCmdQuery, "query", "q",
		[]string{}, "Query pattern")
	getThrottlePoliciesCmd.Flags().StringVarP(&getThrottlePoliciesCmdFormat, "format", "", "", "Pretty-print throttle policies "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getThrottlePoliciesCmd.MarkFlagRequired("environment")
}
