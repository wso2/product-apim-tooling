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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
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
		data := map[interface{}]interface{}{
			"credential":              cred,
			"exportAPIsFormat":        exportAPIsFormat,
			"exportAPIPreserveStatus": exportAPIPreserveStatus,
		}
		cmd.ExecuteExportAPIsCmdByDeprecatedCommand(artifactExportDirectory, data)
	},
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
