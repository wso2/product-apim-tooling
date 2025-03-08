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

const UploadAPIsCmdLiteral = "apis"
const uploadAPIsCmdShortDesc = "Upload APIs of a tenant from one environment to a vector database."
const uploadAPIsCmdLongDesc = "Upload APIs of a tenant from one environment to a vector database to provide context to the marketplace assistant."
const uploadAPIsCmdExamples = utils.ProjectName + ` ` + UploadCmdLiteral + ` ` + UploadAPIsCmdLiteral + ` -e production --all
` + utils.ProjectName + ` ` + UploadCmdLiteral + ` ` + UploadAPIsCmdLiteral + ` -e production
` + utils.ProjectName + ` ` + UploadCmdLiteral + ` ` + UploadAPIsCmdLiteral + ` -e production
NOTE:The flag (--environment (-e)) is mandatory`

var (
	oldToken     string
	oldEndpoint  string
	uploadAll bool
)

var UploadAPIsCmd = &cobra.Command{
	Use: UploadAPIsCmdLiteral + " (--endpoint <endpoint-url> --token <on-prem-key-of-the-organization> --environment " +
         "<environment-from-which-artifacts-should-be-uploaded> --all)",
	Short:   uploadAPIsCmdShortDesc,
	Long:    uploadAPIsCmdLongDesc,
	Example: uploadAPIsCmdExamples,
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
		executeAIUploadAPIsCmd(cred, token, oldEndpoint)
	},
}

// Do operations to upload APIs to the vector database
func executeAIUploadAPIsCmd(credential credentials.Credential, token, oldEndpoint string) {
	impl.AIUploadAPIs(credential, CmdUploadEnvironment, token, oldEndpoint, uploadAll, false)
}

func init() {
	UploadCmd.AddCommand(UploadAPIsCmd)
	UploadAPIsCmd.Flags().StringVarP(&CmdUploadEnvironment, "environment", "e",
		"", "Environment from which the APIs should be uploaded")
	UploadAPIsCmd.Flags().StringVarP(&oldToken, "token", "", "", "on-prem-key of the organization")
	UploadAPIsCmd.Flags().StringVarP(&oldEndpoint, "endpoint", "", "", "endpoint of the marketplace assistant service")
	UploadAPIsCmd.Flags().BoolVarP(&uploadAll, "all", "", false,
		"Upload both apis and api products")
	_ = UploadAPIsCmd.MarkFlagRequired("environment")
}
