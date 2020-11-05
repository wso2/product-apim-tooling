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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	importAPIFile                string
	importEnvironment            string
	importAPICmdPreserveProvider bool
	importAPIUpdate              bool
	importAPIParamsFile          string
	importAPISkipCleanup         bool
)

const (
	// ImportAPI command related usage info
	importAPICmdLiteral   = "import-api"
	importAPICmdShortDesc = "Import API"
	importAPICmdLongDesc  = "Import an API to an environment"
)

const importAPICmdExamples = utils.ProjectName + ` ` + importAPICmdLiteral + ` -f qa/TwitterAPI.zip -e dev
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f staging/FacebookAPI.zip -e production
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f ~/myapi -e production --update
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f ~/myapi -e production --update
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// ImportAPICmd represents the importAPI command
var ImportAPICmdDeprecated = &cobra.Command{
	Use: importAPICmdLiteral + " --file <path-to-api> --environment " +
		"<environment>",
	Short:      importAPICmdShortDesc,
	Long:       importAPICmdLongDesc,
	Example:    importAPICmdExamples,
	Deprecated: "instead use \"" + cmd.ImportCmdLiteral + " " + cmd.ImportAPICmdLiteral + "\".",
	Run: func(deprecatedCmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + importAPICmdLiteral + " called")
		cred, err := cmd.GetCredentials(importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(cred, importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing API", err)
		}
		err = impl.ImportAPIToEnv(accessOAuthToken, importEnvironment, importAPIFile, importAPIParamsFile, importAPIUpdate,
			importAPICmdPreserveProvider, importAPISkipCleanup)
		if err != nil {
			utils.HandleErrorAndExit("Error importing API", err)
			return
		}
	},
}

// init using Cobra
func init() {
	cmd.RootCmd.AddCommand(ImportAPICmdDeprecated)
	ImportAPICmdDeprecated.Flags().StringVarP(&importAPIFile, "file", "f", "",
		"Name of the API to be imported")
	ImportAPICmdDeprecated.Flags().StringVarP(&importEnvironment, "environment", "e",
		"", "Environment from the which the API should be imported")
	ImportAPICmdDeprecated.Flags().BoolVar(&importAPICmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of API after importing")
	ImportAPICmdDeprecated.Flags().BoolVar(&importAPIUpdate, "update", false, "Update an "+
		"existing API or create a new API")
	ImportAPICmdDeprecated.Flags().StringVarP(&importAPIParamsFile, "params", "", utils.ParamFileAPI,
		"Provide a API Manager params file")
	ImportAPICmdDeprecated.Flags().BoolVarP(&importAPISkipCleanup, "skipCleanup", "", false, "Leave "+
		"all temporary files created during import process")
	// Mark required flags
	_ = ImportAPICmdDeprecated.MarkFlagRequired("environment")
	_ = ImportAPICmdDeprecated.MarkFlagRequired("file")
}
