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

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"path/filepath"
)

// var exportThrottlePolicyType string
// var exportThrottlePolicyName string
// var exportThrottlePolicyFormat string
// var runningExportThrottlePolicyCommand bool

// ExportOperationPolicy command related usage info
const ExportOperationPolicyCmdLiteral = "operation"
const exportOperationPolicyCmdShortDesc = "Export Operation Policies"
const exportOperationPolicyCmdLongDesc = "Export Operation Policies from an environment"

const exportOperationPolicyCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportPolicyCmdLiteral + ` ` + ExportOperationPolicyCmdLiteral + ` -n AddHeader -e dev --type sub 
 ` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportPolicyCmdLiteral + ` ` + ExportOperationPolicyCmdLiteral + ` -n AppPolicy -e prod --type app --format JSON
 ` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportPolicyCmdLiteral + ` ` + ExportOperationPolicyCmdLiteral + ` -n TestPolicy -e dev --type advanced 
 ` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportPolicyCmdLiteral + ` ` + ExportOperationPolicyCmdLiteral + ` -n CustomPolicy -e prod --type custom 
 NOTE: All the 2 flags (--name (-n) and --environment (-e)) are mandatory.`

// ExportOperationPolicyCmd represents the export policy operation command
var ExportOperationPolicyCmd = &cobra.Command{
	Use: ExportOperationPolicyCmdLiteral + " (--environment " +
		"<environment-from-which-the-operation-policies-should-be-exported>)",
	Short:   exportOperationPolicyCmdShortDesc,
	Long:    exportOperationPolicyCmdLongDesc,
	Example: exportOperationPolicyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportOperationPolicyCmdLiteral + " called")
		var throttlePoliciesExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedPoliciesDirName, utils.ExportedThrottlePoliciesDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportThrottlePolicyCmd(cred, throttlePoliciesExportDirectory)
	},
}

func executeExportOperationPolicyCmd(credential credentials.Credential, exportDirectory string) {
	runningExportThrottlePolicyCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, CmdExportEnvironment)
	if preCommandErr == nil {
		resp, err := impl.ExportThrottlingPolicyFromEnv(accessToken, CmdExportEnvironment, exportThrottlePolicyName, exportThrottlePolicyType, exportThrottlePolicyFormat)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			impl.WriteThrottlePolicyToFile(exportDirectory, resp, exportThrottlePolicyFormat, runningExportThrottlePolicyCommand)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting Throttling Policies:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// error exporting Throttling Policy
		fmt.Println("Error getting OAuth tokens while exporting Throttling Policies:" + preCommandErr.Error())
	}
}

// init using Cobra
func init() {
	ExportPolicyCmd.AddCommand(ExportOperationPolicyCmd)
	ExportOperationPolicyCmd.Flags().StringVarP(&exportThrottlePolicyName, "name", "n",
		"", "Name of the Operation Policy to be exported")
	ExportOperationPolicyCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment to which the Operation Policies should be exported")
	ExportOperationPolicyCmd.Flags().StringVarP(&exportThrottlePolicyFormat, "format", "", utils.DefaultExportFormat, "File format of exported archive(JSON or YAML)")
	_ = ExportOperationPolicyCmd.MarkFlagRequired("name")
	_ = ExportOperationPolicyCmd.MarkFlagRequired("environment")

}
