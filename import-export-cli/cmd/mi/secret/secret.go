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

package secret

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const secretCmdLiteral = "secret"
const secretCmdShortDesc = "Manage sensitive information"

const secretCmdLongDesc = "Encrypt secrets to be used in the Micro Integrator"

// SecretCmd represents the secret command
var SecretCmd = &cobra.Command{
	Use:   secretCmdLiteral,
	Short: secretCmdShortDesc,
	Long:  secretCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + secretCmdLiteral + " called")
		cmd.Help()
	},
}
