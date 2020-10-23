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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var importAppFile string
var importAppEnvironment string
var importAppOwner string
var preserveOwner bool
var skipSubscriptions bool
var importAppSkipKeys bool
var importAppUpdateApplication bool
var importAppSkipCleanup bool

// ImportApp command related usage info
const importAppCmdLiteral = "import-app"
const importAppCmdShortDesc = "Import App"

const importAppCmdLongDesc = "Import an Application to an environment"

const importAppCmdExamples = utils.ProjectName + ` ` + importAppCmdLiteral + ` -f qa/apps/sampleApp.zip -e dev
` + utils.ProjectName + ` ` + importAppCmdLiteral + ` -f staging/apps/sampleApp.zip -e prod -o testUser
` + utils.ProjectName + ` ` + importAppCmdLiteral + ` -f qa/apps/sampleApp.zip --preserveOwner --skipSubscriptions -e prod
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// importAppCmd represents the importApp command
var ImportAppCmd = &cobra.Command{
	Use: importAppCmdLiteral + " (--file <app-zip-file> --environment " +
		"<environment-to-which-the-app-should-be-imported>)",
	Short:   importAppCmdShortDesc,
	Long:    importAppCmdLongDesc,
	Example: importAppCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + importAppCmdLiteral + " called")
		cred, err := GetCredentials(importAppEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeImportAppCmd(cred)
	},
}

func executeImportAppCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, importAppEnvironment)
	if err != nil {
		utils.HandleErrorAndExit("Error getting OAuth Tokens", err)
	}
	_, err = impl.ImportApplicationToEnv(accessToken, importAppEnvironment, importAppFile, importAppOwner,
		importAppUpdateApplication, preserveOwner, skipSubscriptions, importAppSkipKeys, importAppSkipCleanup)
	if err != nil {
		utils.HandleErrorAndExit("Error importing Application", err)
	}
}

func init() {
	RootCmd.AddCommand(ImportAppCmd)
	ImportAppCmd.Flags().StringVarP(&importAppFile, "file", "f", "",
		"Name of the ZIP file of the Application to be imported")
	ImportAppCmd.Flags().StringVarP(&importAppOwner, "owner", "o", "",
		"Name of the target owner of the Application as desired by the Importer")
	ImportAppCmd.Flags().StringVarP(&importAppEnvironment, "environment", "e",
		"", "Environment from the which the Application should be imported")
	ImportAppCmd.Flags().BoolVarP(&preserveOwner, "preserveOwner", "", false,
		"Preserves app owner")
	ImportAppCmd.Flags().BoolVarP(&skipSubscriptions, "skipSubscriptions", "s", false,
		"Skip subscriptions of the Application")
	ImportAppCmd.Flags().BoolVarP(&importAppSkipKeys, "skipKeys", "", false,
		"Skip importing keys of the Application")
	ImportAppCmd.Flags().BoolVarP(&importAppUpdateApplication, "update", "", false,
		"Update the Application if it is already imported")
	ImportAppCmd.Flags().BoolVarP(&importAppSkipCleanup, "skipCleanup", "", false, "Leave "+
		"all temporary files created during import process")
	_ = ImportAppCmd.MarkFlagRequired("file")
	_ = ImportAppCmd.MarkFlagRequired("environment")
}
