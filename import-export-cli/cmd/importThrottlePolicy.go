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
	importThrottlingPolicyFile string
)

const (
	// ImportAPI command related usage info
	ImportThrottlingPolicyCmdLiteral   = "throttlepolicy"
	importThrottlingPolicyCmdShortDesc = "Import Throttling Policy"
	importThrottlingPolicyCmdLongDesc  = "Import a Throttling Policy to an environment"
)

const importThrottlingPolicyCmdExamples = utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportThrottlingPolicyCmdLiteral + ` -f qa/customadvanced.zip -e dev
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportThrottlingPolicyCmdLiteral + ` -f Env1/Exported/sub1.zip -e production
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportThrottlingPolicyCmdLiteral + ` -f ~/myapi -e production 
` + utils.ProjectName + ` ` + ImportCmdLiteral + ` ` + ImportThrottlingPolicyCmdLiteral + ` -f ~/myapi -e production
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// ImportAPICmd represents the importAPI command
var ImportThrottlingPolicyCmd = &cobra.Command{
	Use: ImportThrottlingPolicyCmdLiteral + " --file <path-to-api> --environment " +
		"<environment>",
	Short:   importThrottlingPolicyCmdShortDesc,
	Long:    importThrottlingPolicyCmdLongDesc,
	Example: importThrottlingPolicyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ImportThrottlingPolicyCmdLiteral + " called")
		cred, err := GetCredentials(importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(cred, importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing Throttling Policy", err)
		}
		err = impl.ImportThrottlingPolicyToEnv(accessOAuthToken, importEnvironment, importThrottlingPolicyFile)
		if err != nil {
			utils.HandleErrorAndExit("Error importing throttling Policy", err)
			return
		}
	},
}

// init using Cobra
func init() {
	ImportCmd.AddCommand(ImportThrottlingPolicyCmd)
	ImportThrottlingPolicyCmd.Flags().StringVarP(&importThrottlingPolicyFile, "file", "f", "",
		"File path of the Throttling Policy to be imported")
	ImportThrottlingPolicyCmd.Flags().StringVarP(&importEnvironment, "environment", "e",
		"", "Environment from the which the Throttling Policy should be imported")
	// Mark required flags
	_ = ImportThrottlingPolicyCmd.MarkFlagRequired("environment")
	_ = ImportThrottlingPolicyCmd.MarkFlagRequired("file")
}
