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

package mi

import (
	"github.com/spf13/cobra"
	miActivateCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/activate"
	miAddCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/add"
	miDeactivateCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/deactivate"
	miDeleteCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/delete"
	miGetCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/get"
	miSecretCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/secret"
	miUpdateCmd "github.com/wso2/product-apim-tooling/import-export-cli/cmd/mi/update"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const miCmdShortDesc = "Micro Integrator related commands"

const miCmdLongDesc = `Micro Integrator related commands such as login, logout, get, add, update, delete, activate, deactivate.`

// MICmd represents the mi command
var MICmd = &cobra.Command{
	Use:   utils.MiCmdLiteral,
	Short: miCmdShortDesc,
	Long:  miCmdLongDesc,
	// Example: miCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + utils.MiCmdLiteral + " called")
		cmd.Help()
	},
}

func init() {
	MICmd.AddCommand(miGetCmd.GetCmd)
	MICmd.AddCommand(miAddCmd.AddCmd)
	MICmd.AddCommand(miDeleteCmd.DeleteCmd)
	MICmd.AddCommand(miUpdateCmd.UpdateCmd)
	MICmd.AddCommand(miActivateCmd.ActivateCmd)
	MICmd.AddCommand(miDeactivateCmd.DeactivateCmd)
	MICmd.AddCommand(miSecretCmd.SecretCmd)
}
