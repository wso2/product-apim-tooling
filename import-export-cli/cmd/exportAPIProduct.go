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

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"net/http"
	"path/filepath"
)

var exportAPIProductName string
var exportAPIProductVersion string
var exportAPIProductRevisionNum string
var exportAPIProductProvider string
var exportAPIProductFormat string
var runningExportAPIProductCommand bool
var exportAPIProductLatestRevision bool


// ExportAPIProduct command related usage info
const ExportAPIProductCmdLiteral = "api-product"
const exportAPIProductCmdShortDesc = "Export API Product"

const exportAPIProductCmdLongDesc = "Export an API Product in an environment"

const exportAPIProductCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPIProductCmdLiteral + ` -n LeasingAPIProduct -e dev
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPIProductCmdLiteral + ` -n CreditAPIProduct -v 1.0.0 -r admin -e production
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory`

// ExportAPIProductCmd represents the exportAPIProduct command
var ExportAPIProductCmd = &cobra.Command{
	Use: ExportAPIProductCmdLiteral + " (--name <name-of-the-api-product> --provider <provider-of-the-api-product> --environment " +
		"<environment-from-which-the-api-product-should-be-exported>)",
	Short:   exportAPIProductCmdShortDesc,
	Long:    exportAPIProductCmdLongDesc,
	Example: exportAPIProductCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportAPIProductCmdLiteral + " called")
		if exportAPIProductRevisionNum == "" && !exportAPIProductLatestRevision {
			fmt.Println("A Revision number is not provided. Only the working copy without deployment environments will be exported." +
				"To export the latest revision, please use --latest flag.")
		}
		var apiProductsExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApiProductsDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportAPIProductCmd(cred, apiProductsExportDirectory)
	},
}

func executeExportAPIProductCmd(credential credentials.Credential, exportDirectory string) {
	runningExportAPIProductCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, CmdExportEnvironment)

	if preCommandErr == nil {
		if exportAPIProductVersion == "" {
			// If the user has not specified the version, use the version as 1.0.0
			exportAPIProductVersion = utils.DefaultApiProductVersion
		}
		resp, err := impl.ExportAPIProductFromEnv(accessToken, exportAPIProductName, exportAPIProductVersion,
			exportAPIProductRevisionNum, exportAPIProductProvider, exportAPIProductFormat, CmdExportEnvironment,
			exportAPIProductLatestRevision)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		apiProductZipLocationPath := filepath.Join(exportDirectory, CmdExportEnvironment)
		if resp.StatusCode() == http.StatusOK {
			impl.WriteAPIProductToZip(exportAPIProductName, exportAPIProductVersion, apiProductZipLocationPath, runningExportAPIProductCommand, resp)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting API Product:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// error exporting API Product
		fmt.Println("Error getting OAuth tokens while exporting API Product:" + preCommandErr.Error())
	}
}

// init using Cobra
func init() {
	ExportCmd.AddCommand(ExportAPIProductCmd)
	ExportAPIProductCmd.Flags().StringVarP(&exportAPIProductName, "name", "n", "",
		"Name of the API Product to be exported")
	ExportAPIProductCmd.Flags().StringVarP(&exportAPIProductRevisionNum, "rev", "", "",
		"Revision number of the API product to be exported")
	ExportAPIProductCmd.Flags().StringVarP(&exportAPIProductProvider, "provider", "r", "",
		"Provider of the API Product")
	ExportAPIProductCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment to which the API Product should be exported")
	ExportAPIProductCmd.Flags().BoolVarP(&exportAPIProductLatestRevision, "latest", "", false,
		"Export the latest revision of the API")
	ExportAPIProductCmd.Flags().StringVarP(&exportAPIProductFormat, "format", "", utils.DefaultExportFormat, "File format of exported archive (json or yaml)")
	_ = ExportAPIProductCmd.MarkFlagRequired("name")
	_ = ExportAPIProductCmd.MarkFlagRequired("environment")
}
