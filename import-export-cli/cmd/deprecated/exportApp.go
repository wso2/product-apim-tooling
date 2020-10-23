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
	"path/filepath"

	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var exportAppName string
var exportAppOwner string
var exportAppWithKeys bool

//var flagExportAPICmdToken string
// ExportApp command related usage info
const exportAppCmdLiteral = "export-app"
const exportAppCmdShortDesc = "Export App"

const exportAppCmdLongDesc = "Export an Application from a specified  environment"

const exportAppCmdExamples = utils.ProjectName + ` ` + exportAppCmdLiteral + ` -n SampleApp -o admin -e dev
` + utils.ProjectName + ` ` + exportAppCmdLiteral + ` -n SampleApp -o admin -e prod
NOTE: All the 3 flags (--name (-n), --owner (-o) and --environment (-e)) are mandatory`

// exportAppCmd represents the exportApp command
var ExportAppCmdDeprecated = &cobra.Command{
	Use: exportAppCmdLiteral + " (--name <name-of-the-application> --owner <owner-of-the-application> --environment " +
		"<environment-from-which-the-app-should-be-exported>)",
	Short:      exportAppCmdShortDesc,
	Long:       exportAppCmdLongDesc,
	Example:    exportAppCmdExamples,
	Deprecated: "instead use \"" + cmd.ExportCmdLiteral + " " + cmd.ExportAppCmdLiteral + "\".",
	Run: func(deprecatedCmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAppCmdLiteral + " called")
		var appsExportDirectoryPath = filepath.Join(utils.ExportDirectory, utils.ExportedAppsDirName, cmd.CmdExportEnvironment)

		cred, err := cmd.GetCredentials(cmd.CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeExportAppCmd(cred, appsExportDirectoryPath)
	},
}

func executeExportAppCmd(credential credentials.Credential, appsExportDirectoryPath string) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmd.CmdExportEnvironment)

	if preCommandErr == nil {
		resp, err := impl.ExportAppFromEnv(accessToken, exportAppName, exportAppOwner, cmd.CmdExportEnvironment, exportAppWithKeys)
		if err != nil {
			utils.HandleErrorAndExit("Error exporting Application: "+exportAppName, err)
		}

		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			impl.WriteApplicationToZip(exportAppName, exportAppOwner, appsExportDirectoryPath, resp)
		} else {
			fmt.Println("Error " + string(resp.Body()))
		}
	} else {
		// error exporting Application
		fmt.Println("Error exporting Application:" + preCommandErr.Error())
	}
}

//init using Cobra
func init() {
	cmd.RootCmd.AddCommand(ExportAppCmdDeprecated)
	ExportAppCmdDeprecated.Flags().StringVarP(&exportAppName, "name", "n", "",
		"Name of the Application to be exported")
	ExportAppCmdDeprecated.Flags().StringVarP(&exportAppOwner, "owner", "o", "",
		"Owner of the Application to be exported")
	ExportAppCmdDeprecated.Flags().StringVarP(&cmd.CmdExportEnvironment, "environment", "e",
		"", "Environment to which the Application should be exported")
	ExportAppCmdDeprecated.Flags().BoolVarP(&exportAppWithKeys, "withKeys", "",
		false, "Export keys for the application ")
	_ = ExportAppCmdDeprecated.MarkFlagRequired("environment")
	_ = ExportAppCmdDeprecated.MarkFlagRequired("owner")
	_ = ExportAppCmdDeprecated.MarkFlagRequired("name")
}
