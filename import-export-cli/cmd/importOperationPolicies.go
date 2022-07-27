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
	importAPIPolicyFile   string
	importAPIPolicyUpdate bool
)

const (
	// ImportAPIPolicyCmdLiteral command related usage info
	ImportAPIPolicyCmdLiteral   = "api"
	importAPIPolicyCmdShortDesc = "Import API Policy"
	importAPIPolicyCmdLongDesc  = "Import a API Policy to an environment"
)

const importAPIPolicyCmdExamples = utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportPolicyCmdLiteral + ` ` + ImportAPIPolicyCmdLiteral + ` -f addHeader_v1.zip -e dev
 ` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportPolicyCmdLiteral + ImportAPIPolicyCmdLiteral + ` -f AddHeader -e production
 NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

var ImportAPIPolicyCmd = &cobra.Command{
	Use: ImportAPIPolicyCmdLiteral + " --file <path-to-api> --environment " +
		"<environment>",
	Short:   importAPIPolicyCmdShortDesc,
	Long:    importAPIPolicyCmdLongDesc,
	Example: importAPIPolicyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ImportAPIPolicyCmdLiteral + " called")
		cred, err := GetCredentials(importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(cred, importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing API Policy", err)
		}
		err = impl.ImportAPIPolicyToEnv(accessOAuthToken, importEnvironment, importAPIPolicyFile)
		if err != nil {
			utils.HandleErrorAndExit("Error importing api Policy", err)
		}
	},
}

// init using Cobra
func init() {
	ImportPolicyCmd.AddCommand(ImportAPIPolicyCmd)
	ImportAPIPolicyCmd.Flags().StringVarP(&importAPIPolicyFile, "file", "f", "",
		"File path of the API Policy to be imported")
	ImportAPIPolicyCmd.Flags().StringVarP(&importEnvironment, "environment", "e",
		"", "Environment from the which the API Policy should be imported")
	// Mark required flags
	_ = ImportAPIPolicyCmd.MarkFlagRequired("environment")
	_ = ImportAPIPolicyCmd.MarkFlagRequired("file")
}
