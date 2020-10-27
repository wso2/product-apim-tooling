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

package deprecated

import (
	"fmt"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"path/filepath"
)

var exportAPIName string
var exportAPIVersion string
var exportProvider string
var exportAPIPreserveStatus bool
var exportAPIFormat string
var runningExportApiCommand bool

// ExportAPI command related usage info
const exportAPICmdLiteral = "export-api"
const exportAPICmdShortDesc = "Export API"

const exportAPICmdLongDesc = "Export an API from an environment"

const exportAPICmdExamples = utils.ProjectName + ` ` + exportAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + exportAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -r admin -e production
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory`

// ExportAPICmd represents the exportAPI command
var ExportAPICmdDeprecated = &cobra.Command{
	Use: exportAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment " +
		"<environment-from-which-the-api-should-be-exported>)",
	Short:      exportAPICmdShortDesc,
	Long:       exportAPICmdLongDesc,
	Example:    exportAPICmdExamples,
	Deprecated: "instead use \"" + cmd.ExportCmdLiteral + " " + cmd.ExportAPICmdLiteral + "\".",
	Run: func(deprecatedCmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAPICmdLiteral + " called")
		var apisExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)

		cred, err := cmd.GetCredentials(cmd.CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportAPICmd(cred, apisExportDirectory)
	},
}

func executeExportAPICmd(credential credentials.Credential, exportDirectory string) {
	runningExportApiCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmd.CmdExportEnvironment)

	if preCommandErr == nil {
		resp, err := impl.ExportAPIFromEnv(accessToken, exportAPIName, exportAPIVersion, exportProvider, exportAPIFormat, cmd.CmdExportEnvironment, exportAPIPreserveStatus)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		apiZipLocationPath := filepath.Join(exportDirectory, cmd.CmdExportEnvironment)
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
	cmd.RootCmd.AddCommand(ExportAPICmdDeprecated)
	ExportAPICmdDeprecated.Flags().StringVarP(&exportAPIName, "name", "n", "",
		"Name of the API to be exported")
	ExportAPICmdDeprecated.Flags().StringVarP(&exportAPIVersion, "version", "v", "",
		"Version of the API to be exported")
	ExportAPICmdDeprecated.Flags().StringVarP(&exportProvider, "provider", "r", "",
		"Provider of the API")
	ExportAPICmdDeprecated.Flags().StringVarP(&cmd.CmdExportEnvironment, "environment", "e",
		"", "Environment to which the API should be exported")
	ExportAPICmdDeprecated.Flags().BoolVarP(&exportAPIPreserveStatus, "preserveStatus", "", true,
		"Preserve API status when exporting. Otherwise API will be exported in CREATED status")
	ExportAPICmdDeprecated.Flags().StringVarP(&exportAPIFormat, "format", "", utils.DefaultExportFormat, "File format of exported archive(json or yaml)")
	_ = ExportAPICmdDeprecated.MarkFlagRequired("name")
	_ = ExportAPICmdDeprecated.MarkFlagRequired("version")
	_ = ExportAPICmdDeprecated.MarkFlagRequired("environment")
}
