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

package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var deleteRoleCmdEnvironment string
var deleteRoleCmdDomain string

const deleteRoleCmdLiteral = "role [role-name]"
const deleteRoleCmdShortDesc = "Delete a role from the Micro Integrator"

const deleteRoleCmdLongDesc = "Delete a role with the name specified by the command line argument [role-name] from a Micro Integrator in the environment specified by the flag --environment, -e"

var deleteRoleCmdExamples = "To delete a role\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + deleteCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(deleteRoleCmdLiteral) + " [role-name] -e dev\n" +
	"To delete a role in a secondary user store\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + deleteCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(deleteRoleCmdLiteral) + " [role-name] -d [domain] -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var deleteRoleCmd = &cobra.Command{
	Use:     deleteRoleCmdLiteral,
	Short:   deleteRoleCmdShortDesc,
	Long:    deleteRoleCmdLongDesc,
	Example: deleteRoleCmdExamples,
	Args:    cobra.ExactArgs(1),
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handledeleteRoleCmdArguments(args)
	},
}

func init() {
	DeleteCmd.AddCommand(deleteRoleCmd)
	deleteRoleCmd.Flags().StringVarP(&deleteRoleCmdDomain, "domain", "d", "", "Select the domain of the role")
	deleteRoleCmd.Flags().StringVarP(&deleteRoleCmdEnvironment, "environment", "e", "", "Environment of the Micro Integrator from which a role should be deleted")
	deleteRoleCmd.MarkFlagRequired("environment")
}

func handledeleteRoleCmdArguments(args []string) {
	printDeleteCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(deleteRoleCmdLiteral))
	credentials.HandleMissingCredentials(deleteRoleCmdEnvironment)
	executeDeleteRole(args[0])
}

func executeDeleteRole(roleName string) {
	resp, err := impl.DeleteMIRole(deleteRoleCmdEnvironment, roleName, deleteRoleCmdDomain)
	if err != nil {
		fmt.Println(utils.LogPrefixError + "deleting role [ "+roleName+" ]", err)
	} else {
		fmt.Println("Deleting role [ "+roleName+" ] status:", resp)
	}
}
