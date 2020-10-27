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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const exportAPIsCmdLiteral = "export-apis"
const exportAPIsCmdShortDesc = "Export APIs for migration"

const exportAPIsCmdLongDesc = "Export all the APIs of a tenant from one environment, to be imported " +
	"into another environment"
const exportAPIsCmdExamples = utils.ProjectName + ` ` + exportAPIsCmdLiteral + ` -e production --force
` + utils.ProjectName + ` ` + exportAPIsCmdLiteral + ` -e production
NOTE: The flag (--environment (-e)) is mandatory`

var exportAPIsFormat string

//e.g. /home/samithac/.wso2apictl/exported/migration/production-2.5/wso2-dot-org
var startFromBeginning bool
var isProcessCompleted bool

var ExportAPIsCmdDeprecated = &cobra.Command{
	Use: exportAPIsCmdLiteral + " (--environment " +
		"<environment-from-which-artifacts-should-be-exported> --format <export-format> --preserveStatus --force)",
	Short:      exportAPIsCmdShortDesc,
	Long:       exportAPIsCmdLongDesc,
	Example:    exportAPIsCmdExamples,
	Deprecated: "instead use \"" + cmd.ExportCmdLiteral + " " + cmd.ExportAPIsCmdLiteral + "\".",
	Run: func(deprecatedCmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAPIsCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)

		cred, err := cmd.GetCredentials(cmd.CmdExportEnvironment)
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
	apiExportDir := impl.CreateExportAPIsDirStructure(exportDirectory, cmd.CmdResourceTenantDomain, cmd.CmdExportEnvironment, cmd.CmdForceStartFromBegin)
	exportRelatedFilesPath := filepath.Join(exportDirectory, cmd.CmdExportEnvironment,
		utils.GetMigrationExportTenantDirName(cmd.CmdResourceTenantDomain))
	//e.g. /home/samithac/.wso2apictl/exported/migration/production-2.5/wso2-dot-org
	startFromBeginning = false
	isProcessCompleted = false

	fmt.Println("\nExporting APIs for the migration...")
	if cmd.CmdForceStartFromBegin {
		startFromBeginning = true
	}

	if (utils.IsFileExist(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName))) && !startFromBeginning {
		impl.PrepareResumption(credential, exportRelatedFilesPath, cmd.CmdResourceTenantDomain, cmd.CmdUsername, cmd.CmdExportEnvironment)
	} else {
		impl.PrepareStartFromBeginning(credential, exportRelatedFilesPath, cmd.CmdResourceTenantDomain, cmd.CmdUsername, cmd.CmdExportEnvironment)
	}

	impl.ExportAPIs(credential, exportRelatedFilesPath, cmd.CmdExportEnvironment, cmd.CmdResourceTenantDomain, exportAPIsFormat, cmd.CmdUsername,
		apiExportDir, exportAPIPreserveStatus, runningExportApiCommand)
}

func init() {
	cmd.RootCmd.AddCommand(ExportAPIsCmdDeprecated)
	ExportAPIsCmdDeprecated.Flags().StringVarP(&cmd.CmdExportEnvironment, "environment", "e",
		"", "Environment from which the APIs should be exported")
	ExportAPIsCmdDeprecated.PersistentFlags().BoolVarP(&cmd.CmdForceStartFromBegin, "force", "", false,
		"Clean all the previously exported APIs of the given target tenant, in the given environment if "+
			"any, and to export APIs from beginning")
	ExportAPIsCmdDeprecated.Flags().BoolVarP(&exportAPIPreserveStatus, "preserveStatus", "", true,
		"Preserve API status when exporting. Otherwise API will be exported in CREATED status")
	ExportAPIsCmdDeprecated.Flags().StringVarP(&exportAPIsFormat, "format", "", utils.DefaultExportFormat, "File format of exported archives(json or yaml)")
	_ = ExportAPIsCmdDeprecated.MarkFlagRequired("environment")
}
