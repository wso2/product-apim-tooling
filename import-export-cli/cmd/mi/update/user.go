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
	"bufio"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var updateUserCmdEnvironment string

const updateUserCmdLiteral = "user [user-name]"
const updateUserCmdShortDesc = "Update roles of a user in a Micro Integrator"

const updateUserCmdLongDesc = "Update the roles of a user named [user-name] specified by the command line arguments in a Micro Integrator in the environment specified by the flag --environment, -e"

var updateUserCmdExamples = "To update the roles\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + updateCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(updateUserCmdLiteral) + " [user-name] -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var updateUserCmd = &cobra.Command{
	Use:     updateUserCmdLiteral,
	Short:   updateUserCmdShortDesc,
	Long:    updateUserCmdLongDesc,
	Example: updateUserCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleupdateUserCmdArguments(args)
	},
}

func init() {
	UpdateCmd.AddCommand(updateUserCmd)
	updateUserCmd.Flags().StringVarP(&updateUserCmdEnvironment, "environment", "e", "", "Environment of the Micro Integrator of which the user's roles should be updated")
	updateUserCmd.MarkFlagRequired("environment")
}

func handleupdateUserCmdArguments(args []string) {
	printUpdateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(updateUserCmdLiteral))
	credentials.HandleMissingCredentials(updateUserCmdEnvironment)
	startConsoleToUpdateUser(args[0])
}

func executeUpdateUser(userName, domain string, addedRoles, removedRoles []string) {
	resp, err := impl.UpdateMIUser(updateUserCmdEnvironment, userName, domain, addedRoles, removedRoles)
	if err != nil {
		fmt.Println(utils.LogPrefixError+"updating roles of user [ "+userName+" ] ", err)
	} else {
		fmt.Println(resp)
	}
}

func startConsoleToUpdateUser(userName string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter user store(domain) for " + userName + " default (primary): ")
	domain, _ := reader.ReadString('\n')
	domain = strings.TrimSuffix(domain, "\n")

	fmt.Printf("Enter list of new roles to be assigned to " + userName + ": ")
	addedRolesStr, _ := reader.ReadString('\n')
	addedRolesStr = strings.TrimSuffix(addedRolesStr, "\n")
	addedRoles := strings.Fields(addedRolesStr)

	fmt.Printf("Enter list of existing roles to be revoked from " + userName + ": ")
	removedRolesStr, _ := reader.ReadString('\n')
	removedRolesStr = strings.TrimSuffix(removedRolesStr, "\n")
	removedRoles := strings.Fields(removedRolesStr)

	executeUpdateUser(userName,domain,addedRoles,removedRoles)
}
