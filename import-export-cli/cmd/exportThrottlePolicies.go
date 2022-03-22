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
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"path/filepath"
)

var exportThrottlePoliciesType string

var runningExportThrottlePoliciesCommand bool

// ExportThrottlePolicies command related usage info
const ExportThrottlePoliciesCmdLiteral = "throttlepolicies"
const exportThrottlePoliciesCmdShortDesc = "Export Throttling Policies"

const exportThrottlePoliciesCmdLongDesc = "Export ThrottlingPolicies from an environment"

const exportThrottlePoliciesCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportThrottlePoliciesCmdLiteral + `-type custom -e dev
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportThrottlePoliciesCmdLiteral + `-type app -e production
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportThrottlePoliciesCmdLiteral + `-type sub -e dev
NOTE: All the 2 flags (--type (-t) and --environment (-e)) are mandatory.`

// ExportThrottlePoliciesCmd represents the export throttlepolicies command
var ExportThrottlePoliciesCmd = &cobra.Command{
	Use: ExportThrottlePoliciesCmdLiteral + " (--type <type-of-the-throttling-policy> --environment " +
		"<environment-from-which-the-throttling-policies-should-be-exported>)",
	Short:   exportThrottlePoliciesCmdShortDesc,
	Long:    exportThrottlePoliciesCmdLongDesc,
	Example: exportThrottlePoliciesCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportThrottlePoliciesCmdLiteral + " called")
		var throttlepoliciesExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedThrottlePoliciesDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportThrottlePoliciesCmd(cred, throttlepoliciesExportDirectory)
	},
}

func executeExportThrottlePoliciesCmd(credential credentials.Credential, exportDirectory string) {
	runningExportThrottlePoliciesCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, CmdExportEnvironment)

	if preCommandErr == nil {
		resp, err := impl.ExportThrottlingPoliciesFromEnv(accessToken, CmdExportEnvironment, exportThrottlePoliciesType)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		ThrottlePolicyZipLocationPath := filepath.Join(exportDirectory, CmdExportEnvironment)
		if resp.StatusCode() == http.StatusOK {
			//fmt.Println(string(resp.Body()))
			impl.ThrottlePoliciesWriteToZip(exportThrottlePoliciesType, ThrottlePolicyZipLocationPath, runningExportThrottlePoliciesCommand, resp)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting Throttling Policies:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// error exporting Api
		fmt.Println("Error getting OAuth tokens while exporting Throttling Policies:" + preCommandErr.Error())
	}
}

// init using Cobra
func init() {
	ExportCmd.AddCommand(ExportThrottlePoliciesCmd)
	ExportThrottlePoliciesCmd.Flags().StringVarP(&exportThrottlePoliciesType, "type", "t",
		"", "Type of the Throttling Policies to be exported (sub,app,custom,advanced,deny)")
	ExportThrottlePoliciesCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment to which the Throttling Policies should be exported")
	_ = ExportThrottlePoliciesCmd.MarkFlagRequired("environment")
	_ = ExportThrottlePoliciesCmd.MarkFlagRequired("type")

}
