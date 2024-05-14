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
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const PurgeAPIsCmdLiteral = "artifacts"
const purgeAPIsCmdShortDesc = "Purge APIs and API Products of a tenant from one environment from a vector database."

const purgeAPIsCmdLongDesc = "Purge APIs and API Products of a tenant from one environment from a vector database."
const PurgeAPIsCmdLongDesc = `Purge APIs and API Products of a tenant from one environment specified by flag (--environment, -e)`
const purgeAPIsCmdExamples = utils.ProjectName + ` ` + AiCmdLiteral + ` ` + PurgeCmdLiteral + ` ` + PurgeAPIsCmdLiteral + ` --token 2fdca1b6-6a28-4aea-add6-77c97033bdb9 --endpoint https://dev-tools.wso2.com/apim-ai-service -e production
NOTE:The flag (--environment (-e)) is mandatory`

var PurgeAPIsCmd = &cobra.Command{
	Use: PurgeAPIsCmdLiteral + " (--endpoint <endpoint-url> --token <on-prem-key-of-the-organization> --environment " +
		"<environment-from-which-artifacts-should-be-purgeed>)",
	Short:   purgeAPIsCmdShortDesc,
	Long:    purgeAPIsCmdLongDesc,
	Example: purgeAPIsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + PurgeAPIsCmdLiteral + " called")

		cred, err := GetCredentials(CmdPurgeEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeAIDeleteAPIsCmd(cred, token, endpoint)
	},
}

// Do operations to Purge APIs to the vector database
func executeAIDeleteAPIsCmd(credential credentials.Credential, token, endpoint string) {
	var Tenant string
	if !strings.Contains(credential.Username, "@") {
		Tenant = "carbon.super"
	} else {
		Tenant = strings.Split(credential.Username, "@")[1]
	}
	impl.AIDeleteAPIs(credential, CmdPurgeEnvironment, token, endpoint, Tenant)
}

func init() {
	PurgeCmd.AddCommand(PurgeAPIsCmd)
	PurgeAPIsCmd.Flags().StringVarP(&CmdPurgeEnvironment, "environment", "e",
		"", "Environment from which the APIs should be Purgeed")
	PurgeAPIsCmd.Flags().StringVarP(&token, "token", "", "", "on-prem-key of the organization")
	PurgeAPIsCmd.Flags().StringVarP(&endpoint, "endpoint", "", "", "endpoint of the marketplace assistant service")
	_ = PurgeAPIsCmd.MarkFlagRequired("environment")
}
