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

package update

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var updateHashiCorpSecretCmdEnvironment string

const updateHashiCorpSecretCmdLiteral = "hashicorp-secret [secret-id]"
const updateHashiCorpSecretCmdShortDesc = "Update the secret ID of HashiCorp configuration in a Micro Integrator"

const updateHashiCorpSecretCmdLongDesc = "Update the secret ID of the HashiCorp configuration in a Micro Integrator in the environment specified by the flag --environment, -e"

var updateHashiCorpSecretCmdExamples = "To update the secret ID\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + updateCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(updateHashiCorpSecretCmdLiteral) + " new_secret_id -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var updateHashiCorpSecretCmd = &cobra.Command{
	Use:     updateHashiCorpSecretCmdLiteral,
	Short:   updateHashiCorpSecretCmdShortDesc,
	Long:    updateHashiCorpSecretCmdLongDesc,
	Example: updateHashiCorpSecretCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleUpdateHashiCorpSecretCmdArguments(args)
	},
}

func init() {
	UpdateCmd.AddCommand(updateHashiCorpSecretCmd)
	updateHashiCorpSecretCmd.Flags().StringVarP(&updateHashiCorpSecretCmdEnvironment, "environment", "e", "", "Environment of the micro integrator of which the HashiCorp secret ID should be updated")
	updateHashiCorpSecretCmd.MarkFlagRequired("environment")
}

func handleUpdateHashiCorpSecretCmdArguments(args []string) {
	printUpdateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(updateHashiCorpSecretCmdLiteral))
	credentials.HandleMissingCredentials(updateHashiCorpSecretCmdEnvironment)
	executeUpdateHashiCorpSecretID(args[0])
}

func executeUpdateHashiCorpSecretID(hashiCorpSecretID string) {
	resp, err := impl.UpdateHashiCorpSecretID(updateHashiCorpSecretCmdEnvironment, hashiCorpSecretID)
	if err != nil {
		fmt.Println(utils.LogPrefixError+"updating secretID of HashiCorp configuration.", err)
	} else {
		fmt.Println(resp)
	}
}
