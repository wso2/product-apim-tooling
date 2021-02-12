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

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"net/http"
	"path/filepath"
)

var exportAPIName string
var exportAPIVersion string
var exportRevisionNum string
var exportProvider string
var exportAPIPreserveStatus bool
var exportAPIFormat string
var runningExportApiCommand bool
var exportAPILatestRevision bool

// ExportAPI command related usage info
const ExportAPICmdLiteral = "api"
const exportAPICmdShortDesc = "Export API"

const exportAPICmdLongDesc = "Export an API from an environment"

const exportAPICmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 6 -r admin -e production
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 2 -r admin -e production
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory. If --rev is not provided, working copy of the API
without deployment environments will be exported.`

// ExportAPICmd represents the exportAPI command
var ExportAPICmd = &cobra.Command{
	Use: ExportAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment " +
		"<environment-from-which-the-api-should-be-exported>)",
	Short:   exportAPICmdShortDesc,
	Long:    exportAPICmdLongDesc,
	Example: exportAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportAPICmdLiteral + " called")
		if exportRevisionNum == "" && !exportAPILatestRevision {
			fmt.Println("A Revision number is not provided. Only the working copy without deployment environments will be exported." +
				"To export the latest revision, please use --latest flag.")
		}
		var apisExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportAPICmd(cred, apisExportDirectory)
	},
}

func executeExportAPICmd(credential credentials.Credential, exportDirectory string) {
	runningExportApiCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, CmdExportEnvironment)

	if preCommandErr == nil {
		resp, err := impl.ExportAPIFromEnv(accessToken, exportAPIName, exportAPIVersion, exportRevisionNum, exportProvider,
			exportAPIFormat, CmdExportEnvironment, exportAPIPreserveStatus, exportAPILatestRevision)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		apiZipLocationPath := filepath.Join(exportDirectory, CmdExportEnvironment)
		if resp.StatusCode() == http.StatusOK {
			impl.WriteToZip(exportAPIName, exportAPIVersion, apiZipLocationPath, runningExportApiCommand, resp)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting API:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// error exporting Api
		fmt.Println("Error getting OAuth tokens while exporting API:" + preCommandErr.Error())
	}
}

// init using Cobra
func init() {
	ExportCmd.AddCommand(ExportAPICmd)
	ExportAPICmd.Flags().StringVarP(&exportAPIName, "name", "n", "",
		"Name of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportAPIVersion, "version", "v", "",
		"Version of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportProvider, "provider", "r", "",
		"Provider of the API")
	ExportAPICmd.Flags().StringVarP(&exportRevisionNum, "rev", "", "",
		"Revision number of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment to which the API should be exported")
	ExportAPICmd.Flags().BoolVarP(&exportAPIPreserveStatus, "preserveStatus", "", true,
		"Preserve API status when exporting. Otherwise API will be exported in CREATED status")
	ExportAPICmd.Flags().BoolVarP(&exportAPILatestRevision, "latest", "", false,
		"Export the latest revision of the API")
	ExportAPICmd.Flags().StringVarP(&exportAPIFormat, "format", "", utils.DefaultExportFormat, "File format of exported archive(json or yaml)")
	_ = ExportAPICmd.MarkFlagRequired("name")
	_ = ExportAPICmd.MarkFlagRequired("version")
	_ = ExportAPICmd.MarkFlagRequired("environment")
}
