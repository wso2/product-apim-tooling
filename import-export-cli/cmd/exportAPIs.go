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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const ExportAPIsCmdLiteral = "apis"
const exportAPIsCmdShortDesc = "Export APIs for migration"

const exportAPIsCmdLongDesc = "Export all the APIs of a tenant from one environment, to be imported " +
	"into another environment"
const exportAPIsCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPIsCmdLiteral + ` -e production --force
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportAPIsCmdLiteral + ` -e production
NOTE: The flag (--environment (-e)) is mandatory`

var exportAPIsFormat string

//e.g. /home/samithac/.wso2apictl/exported/migration/production-2.5/wso2-dot-org
var startFromBeginning bool
var isProcessCompleted bool

var ExportAPIsCmd = &cobra.Command{
	Use: ExportAPIsCmdLiteral + " (--environment " +
		"<environment-from-which-artifacts-should-be-exported> --format <export-format> --preserveStatus --force)",
	Short:   exportAPIsCmdShortDesc,
	Long:    exportAPIsCmdLongDesc,
	Example: exportAPIsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportAPIsCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeExportAPIsCmd(cred, artifactExportDirectory)
	},
}

// Do operations to export APIs for the migration into the directory passed as exportDirectory
// <export_directory> is the patch defined in main_config.yaml
// exportDirectory = <export_directory>/migration/
func executeExportAPIsCmd(credential credentials.Credential, exportDirectory string) {
	//create dir structure
	apiExportDir := impl.CreateExportAPIsDirStructure(exportDirectory, CmdResourceTenantDomain, CmdExportEnvironment, CmdForceStartFromBegin)
	exportRelatedFilesPath := filepath.Join(exportDirectory, CmdExportEnvironment,
		utils.GetMigrationExportTenantDirName(CmdResourceTenantDomain))
	//e.g. /home/samithac/.wso2apictl/exported/migration/production-2.5/wso2-dot-org
	startFromBeginning = false
	isProcessCompleted = false

	fmt.Println("\nExporting APIs for the migration...")
	if CmdForceStartFromBegin {
		startFromBeginning = true
	}

	if (utils.IsFileExist(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName))) && !startFromBeginning {
		impl.PrepareResumption(credential, exportRelatedFilesPath, CmdResourceTenantDomain, CmdUsername, CmdExportEnvironment)
	} else {
		impl.PrepareStartFromBeginning(credential, exportRelatedFilesPath, CmdResourceTenantDomain, CmdUsername, CmdExportEnvironment)
	}

	impl.ExportAPIs(credential, exportRelatedFilesPath, CmdExportEnvironment, CmdResourceTenantDomain, exportAPIsFormat, CmdUsername,
		apiExportDir, exportAPIPreserveStatus, runningExportApiCommand)
}

func init() {
	ExportCmd.AddCommand(ExportAPIsCmd)
	ExportAPIsCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment from which the APIs should be exported")
	ExportAPIsCmd.PersistentFlags().BoolVarP(&CmdForceStartFromBegin, "force", "", false,
		"Clean all the previously exported APIs of the given target tenant, in the given environment if "+
			"any, and to export APIs from beginning")
	ExportAPIsCmd.Flags().BoolVarP(&exportAPIPreserveStatus, "preserveStatus", "", true,
		"Preserve API status when exporting. Otherwise API will be exported in CREATED status")
	ExportAPIsCmd.Flags().StringVarP(&exportAPIsFormat, "format", "", utils.DefaultExportFormat, "File format of exported archives(json or yaml)")
	_ = ExportAPIsCmd.MarkFlagRequired("environment")
}
