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

const UploadAPIProductsCmdLiteral = "api-products"
const uploadAPIProductsCmdShortDesc = "Upload API Products of a tenant from one environment to a vector database."
const uploadAPIProductsCmdLongDesc = `Upload public API Products of a tenant from one environment specified by flag (--environment, -e)`
const uploadAPIProductsCmdExamples = utils.ProjectName + ` ` + UploadCmdLiteral + ` ` + UploadAPIProductsCmdLiteral + ` -e production --all
` + utils.ProjectName + ` ` + UploadCmdLiteral + ` ` + UploadAPIProductsCmdLiteral + ` -e production
` + utils.ProjectName + ` ` + UploadCmdLiteral + ` ` + UploadAPIProductsCmdLiteral + ` -e production
NOTE:The flag (--environment (-e)) is mandatory`

var UploadAPIProductsCmd = &cobra.Command{
	Use: UploadAPIProductsCmdLiteral + " (--endpoint <endpoint-url> --token <on-prem-key-of-the-organization> --environment " +
         "<environment-from-which-artifacts-should-be-uploaded>)",
	Short:   uploadAPIProductsCmdShortDesc,
	Long:    uploadAPIProductsCmdLongDesc,
	Example: uploadAPIProductsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + UploadAPIsCmdLiteral + " called")

		cred, err := GetCredentials(CmdUploadEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		if (oldToken != "") {
			token = oldToken
		} else {
			key, err := utils.GetAIKeyOfEnv(CmdUploadEnvironment, utils.MainConfigFilePath)
			if err != nil {
			    utils.HandleErrorAndExit("Error getting AI key", err)
			}
			token, err = impl.GetAIToken(key, CmdUploadEnvironment)
			if err != nil {
			    utils.HandleErrorAndExit("Error getting AI token", err)
			}
		}
		executeAIUploadAPIProductsCmd(cred, token, oldEndpoint)
	},
}

// Do operations to upload APIs to the vector database
func executeAIUploadAPIProductsCmd(credential credentials.Credential, token, oldEndpoint string) {
	impl.AIUploadAPIs(credential, CmdUploadEnvironment, token, oldEndpoint, uploadAll, true)
}

func init() {
	UploadCmd.AddCommand(UploadAPIProductsCmd)
	UploadAPIProductsCmd.Flags().StringVarP(&CmdUploadEnvironment, "environment", "e",
		"", "Environment from which the APIs should be uploaded")
	UploadAPIProductsCmd.Flags().StringVarP(&oldToken, "token", "", "", "on-prem-key of the organization")
	UploadAPIProductsCmd.Flags().StringVarP(&oldEndpoint, "endpoint", "", "", "endpoint of the marketplace assistant service")
	UploadAPIProductsCmd.Flags().BoolVarP(&uploadAll, "all", "", false,
		"Upload both apis and api products")
	_ = UploadAPIProductsCmd.MarkFlagRequired("environment")

}
