/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
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
	"path/filepath"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var exportAppsWithKeys bool
var exportAppsFormat string
var startFromBeginningForApps bool
var isProcessCompletedForApps bool

// ExportApps command related usage info
const ExportAppsCmdLiteral = "apps"
const exportAppsCmdShortDesc = "Export Applications"

const exportAppsCmdLongDesc = "Export Applications of a given tenant from a specified environment"

const exportAppsCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAppsCmdLiteral + ` -e dev --force
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAppsCmdLiteral + ` -e prod
NOTE: The flag (--environment (-e)) is mandatory`

// ExportAppsCmd represents the exportApps command
var ExportAppsCmd = &cobra.Command{
	Use: ExportAppsCmdLiteral + " (--environment <environment-from-which-the-app-should-be-exported> " +
	"--format <export-format> --force)",
	Short:   exportAppsCmdShortDesc,
	Long:    exportAppsCmdLongDesc,
	Example: exportAppsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportAppsCmdLiteral + " called")
		var appsExportDirectoryPath = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeExportAppsCmd(cred, appsExportDirectoryPath)
	},
}

func executeExportAppsCmd(credential credentials.Credential, appsExportDirectoryPath string) {
    // Create dir structure
    appExportDir := impl.CreateExportAppsDirStructure(appsExportDirectoryPath, CmdResourceTenantDomain, CmdExportEnvironment,
        CmdForceStartFromBegin)
    exportRelatedFilesPath := filepath.Join(appsExportDirectoryPath, CmdExportEnvironment,
        utils.GetMigrationExportTenantDirName(CmdResourceTenantDomain))
    // e.g. /home/user/.wso2apictl/exported/migration/production/wso2-dot-org
    startFromBeginningForApps = false
    isProcessCompletedForApps = false

    fmt.Println("\nExporting Applications...")
    if CmdForceStartFromBegin {
        startFromBeginningForApps = true
    }

    if (utils.IsFileExist(filepath.Join(exportRelatedFilesPath, utils.LastSucceededAppFileName))) && !startFromBeginningForApps {
        impl.PrepareResumptionForApps(credential, exportRelatedFilesPath, CmdResourceTenantDomain, CmdUsername, CmdExportEnvironment)
    } else {
        impl.PrepareStartAppsFromBeginning(credential, exportRelatedFilesPath, CmdResourceTenantDomain, CmdUsername, CmdExportEnvironment)
    }

    impl.ExportApps(credential, exportRelatedFilesPath, CmdExportEnvironment, CmdResourceTenantDomain, exportAppsFormat,
     CmdUsername, appExportDir, exportAppsWithKeys)
}

// Init using Cobra
func init() {
	ExportCmd.AddCommand(ExportAppsCmd)
	ExportAppsCmd.PersistentFlags().BoolVarP(&CmdForceStartFromBegin, "force", "", false,
		"Clean all the previously exported Apps of the given target tenant, in the given environment if "+
            "any, and to export Apps from beginning")
	ExportAppsCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment from which the Applications should be exported")
	ExportAppsCmd.Flags().BoolVarP(&exportAppsWithKeys, "with-keys", "",
		false, "Export keys for the applications")
	ExportAppsCmd.Flags().StringVarP(&exportAppsFormat, "format", "", utils.DefaultExportFormat, "File format of exported archive (json or yaml)")
	_ = ExportAppCmd.MarkFlagRequired("environment")
}
