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

var deleteUserCmdEnvironment string
var deleteUserCmdDomain string

const deleteUserCmdLiteral = "user [user-name]"
const deleteUserCmdShortDesc = "Delete a user from the Micro Integrator"

const deleteUserCmdLongDesc = "Delete a user with the name specified by the command line argument [user-name] from a Micro Integrator in the environment specified by the flag --environment, -e"

var deleteUserCmdExamples = "To delete a user\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + deleteCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(deleteUserCmdLiteral) + " [user-id] -e dev\n" +
	"To delete a user in a secondary user store\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + deleteCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(deleteUserCmdLiteral) + " [user-id] -d [domain] -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var deleteUserCmd = &cobra.Command{
	Use:     deleteUserCmdLiteral,
	Short:   deleteUserCmdShortDesc,
	Long:    deleteUserCmdLongDesc,
	Example: deleteUserCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handledeleteUserCmdArguments(args)
	},
}

func init() {
	DeleteCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().StringVarP(&deleteUserCmdDomain, "domain", "d", "", "select user's domain")
	deleteUserCmd.Flags().StringVarP(&deleteUserCmdEnvironment, "environment", "e", "", "Environment of the micro integrator from which a user should be deleted")
	deleteUserCmd.MarkFlagRequired("environment")
}

func handledeleteUserCmdArguments(args []string) {
	printDeleteCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(deleteUserCmdLiteral))
	credentials.HandleMissingCredentials(deleteUserCmdEnvironment)
	executeDeleteUser(args[0])
}

func executeDeleteUser(userName string) {
	resp, err := impl.DeleteMIUser(deleteUserCmdEnvironment, userName, deleteUserCmdDomain)
	if err != nil {
		fmt.Println(utils.LogPrefixError+"deleting user [ "+userName+" ]", err)
	} else {
		fmt.Println("Deleting user [ "+userName+" ] status:", resp)
	}
}

func printDeleteCmdVerboseLog(cmd string) {
	utils.Logln(utils.LogPrefixInfo + deleteCmdLiteral + " " + cmd + " called")
}
