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

var (
	importAPIProductFile                string
	importAPIProductEnvironment         string
	importAPIProductCmdPreserveProvider bool
	importAPIs                          bool
	importAPIProductUpdate              bool
	importAPIsUpdate                    bool
	importAPIProductSkipCleanup         bool
	importAPIProductRotateRevision		bool
)

const (
	// ImportAPIProduct command related usage info
	importAPIProductCmdLiteral   = "api-product"
	importAPIProductCmdShortDesc = "Import API Product"
	importAPIProductCmdLongDesc  = "Import an API Product to an environment"
)

const importAPIProductCmdExamples = utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f qa/LeasingAPIProduct.zip -e dev
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f staging/CreditAPIProduct.zip -e production --update-api-product
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f ~/myapiproduct -e production
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + importAPIProductCmdLiteral + ` -f ~/myapiproduct -e production --update-api-product --update-apis
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// ImportAPIProductCmd represents the importAPIProduct command
var ImportAPIProductCmd = &cobra.Command{
	Use: importAPIProductCmdLiteral + " (--file <path-to-api-product> --environment " +
		"<environment-to-which-the-api-product-should-be-imported>)",
	Short:   importAPIProductCmdShortDesc,
	Long:    importAPIProductCmdLongDesc,
	Example: importAPIProductCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + importAPIProductCmdLiteral + " called")

		cred, err := GetCredentials(importAPIProductEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(cred, importAPIProductEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing API Product", err)
		}
		err = impl.ImportAPIProductToEnv(accessOAuthToken, importAPIProductEnvironment, importAPIProductFile, importAPIs, importAPIsUpdate,
			importAPIProductUpdate, importAPIProductCmdPreserveProvider, importAPIProductSkipCleanup, importAPIProductRotateRevision)
		if err != nil {
			utils.HandleErrorAndExit("Error importing API Product", err)
			return
		}
	},
}

// init using Cobra
func init() {
	ImportCmd.AddCommand(ImportAPIProductCmd)
	ImportAPIProductCmd.Flags().StringVarP(&importAPIProductFile, "file", "f", "",
		"Name of the API Product to be imported")
	ImportAPIProductCmd.Flags().StringVarP(&importAPIProductEnvironment, "environment", "e",
		"", "Environment from the which the API Product should be imported")
	ImportAPIProductCmd.Flags().BoolVar(&importAPIProductRotateRevision, "rotate-revision", false,
		"If the maximum revision limit is reached, undeploy and delete the earliest revision")
	ImportAPIProductCmd.Flags().BoolVar(&importAPIProductCmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of API Product after importing")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIs, "import-apis", "", false, "Import "+
		"dependent APIs associated with the API Product")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIProductUpdate, "update-api-product", "", false, "Update an "+
		"existing API Product or create a new API Product")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIsUpdate, "update-apis", "", false, "Update existing dependent APIs "+
		"associated with the API Product")
	ImportAPIProductCmd.Flags().BoolVarP(&importAPIProductSkipCleanup, "skipCleanup", "", false, "Leave "+
		"all temporary files created during import process")
	// Mark required flags
	_ = ImportAPIProductCmd.MarkFlagRequired("environment")
	_ = ImportAPIProductCmd.MarkFlagRequired("file")
}
