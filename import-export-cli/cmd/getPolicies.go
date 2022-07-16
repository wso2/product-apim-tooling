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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const GetPoliciesCmdLiteral = "policies"
const GetPoliciesCmdShortDesc = "Get Policy list"
const GetPoliciesCmdLongDesc = "Get a list of Policies in an environment"
const GetPoliciesCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetThrottlePoliciesCmdLiteral + ` -e production -q type:sub`
const GetAPIPoliciesCmdExample = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetPoliciesCmdLiteral + ` ` + GetAPIPoliciesCmdLiteral + ` -e production -q gateway:choreo`

// GetPoliciesCmd  represents the get command for policies
var GetPoliciesCmd = &cobra.Command{
	Use:     GetPoliciesCmdLiteral,
	Short:   GetPoliciesCmdShortDesc,
	Long:    GetPoliciesCmdLongDesc,
	Example: GetPoliciesCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetCmdLiteral + " called")
	},
}

// init using Cobra
func init() {
	GetCmd.AddCommand(GetPoliciesCmd)
}
