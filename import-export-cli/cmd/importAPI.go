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
	"github.com/spf13/cobra"
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
	importAPIRotateRevision      bool
	importAPISkipDeployments     bool
)

const (
	// ImportAPI command related usage info
	ImportAPICmdLiteral   = "api"
	importAPICmdShortDesc = "Import API"
	importAPICmdLongDesc  = "Import an API to an environment"
)

const importAPICmdExamples = utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportAPICmdLiteral + ` -f qa/TwitterAPI.zip -e dev
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportAPICmdLiteral + ` -f staging/FacebookAPI.zip -e production
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportAPICmdLiteral + ` -f ~/myapi -e production --update --rotate-revision
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportAPICmdLiteral + ` -f ~/myapi -e production --update
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use: ImportAPICmdLiteral + " --file <path-to-api> --environment " +
		"<environment>",
	Short:   importAPICmdShortDesc,
	Long:    importAPICmdLongDesc,
	Example: importAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ImportAPICmdLiteral + " called")
		if importAPIRotateRevision == false {
			fmt.Println("The Rotate revision flag is set to false. Please note that you can only update and deploy the " +
				"API 5 times without deleting a revision")
		}
		cred, err := GetCredentials(importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(cred, importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing API", err)
		}
		err = impl.ImportAPIToEnv(accessOAuthToken, importEnvironment, importAPIFile, importAPIParamsFile, importAPIUpdate,
			importAPICmdPreserveProvider, importAPISkipCleanup, importAPIRotateRevision, importAPISkipDeployments)
		if err != nil {
			utils.HandleErrorAndExit("Error importing API", err)
			return
		}
	},
}

// init using Cobra
func init() {
	ImportCmd.AddCommand(ImportAPICmd)
	ImportAPICmd.Flags().StringVarP(&importAPIFile, "file", "f", "",
		"Name of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importEnvironment, "environment", "e",
		"", "Environment from the which the API should be imported")
	ImportAPICmd.Flags().BoolVar(&importAPICmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of API after importing")
	ImportAPICmd.Flags().BoolVar(&importAPIUpdate, "update", false, "Update an "+
		"existing API or create a new API")
	ImportAPICmd.Flags().BoolVar(&importAPIRotateRevision, "rotate-revision", false, "Rotate the "+
		"revisions with each update")
	ImportAPICmd.Flags().BoolVar(&importAPISkipDeployments, "skip-deployments", false, "Update only " +
		"the working copy and skip deployment steps in import")
	ImportAPICmd.Flags().StringVarP(&importAPIParamsFile, "params", "", utils.ParamFileAPI,
		"Provide a API Manager params file")
	ImportAPICmd.Flags().BoolVarP(&importAPISkipCleanup, "skip-cleanup", "", false, "Leave "+
		"all temporary files created during import process")
	// Mark required flags
	_ = ImportAPICmd.MarkFlagRequired("environment")
	_ = ImportAPICmd.MarkFlagRequired("file")
}
