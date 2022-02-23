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

package add

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var addRoleCmdEnvironment string

const addRoleCmdLiteral = "role [role-name]"
const addRoleCmdShortDesc = "Add new role to a Micro Integrator"

const addRoleCmdLongDesc = "Add a new role with the name specified by the command line argument [role-name] to a Micro Integrator in the environment specified by the flag --environment, -e"

var addRoleCmdExamples = "To add a new role\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + addCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(addRoleCmdLiteral) + " [role-name] -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var addRoleCmd = &cobra.Command{
	Use:     addRoleCmdLiteral,
	Short:   addRoleCmdShortDesc,
	Long:    addRoleCmdLongDesc,
	Example: addRoleCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleAddRoleCmdArguments(args)
	},
}

func init() {
	AddCmd.AddCommand(addRoleCmd)
	addRoleCmd.Flags().StringVarP(&addRoleCmdEnvironment, "environment", "e", "", "Environment of the micro integrator to which a new user should be added")
	addRoleCmd.MarkFlagRequired("environment")
}

func handleAddRoleCmdArguments(args []string) {
	printAddCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(addRoleCmdLiteral))
	credentials.HandleMissingCredentials(addRoleCmdEnvironment)
	startConsoleToAddRole(args[0])
}

func startConsoleToAddRole(roleName string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter user store (domain) for " + roleName + " default (primary): ")
	domain, _ := reader.ReadString('\n')
	domain = strings.TrimSuffix(domain, "\n")

	executeAddNewRole(roleName, domain)
}

func executeAddNewRole(roleName, domain string) {
	resp, err := impl.AddMIRole(addRoleCmdEnvironment, roleName, domain)
	if err != nil {
		fmt.Println(utils.LogPrefixError+"Adding new role [ "+roleName+" ]", err)
	} else {
		fmt.Println("Adding new role [ "+roleName+" ] status:", resp)
	}
}
