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
	"syscall"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var addUserCmdEnvironment string

const addUserCmdLiteral = "user [user-name]"
const addUserCmdShortDesc = "Add new user to a Micro Integrator"

const addUserCmdLongDesc = "Add a new user with the name specified by the command line argument [user-name] to a Micro Integrator in the environment specified by the flag --environment, -e"

var addUserCmdExamples = "To add a new user\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + addCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(addUserCmdLiteral) + " capp-tester -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var addUserCmd = &cobra.Command{
	Use:     addUserCmdLiteral,
	Short:   addUserCmdShortDesc,
	Long:    addUserCmdLongDesc,
	Example: addUserCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleAddUserCmdArguments(args)
	},
}

func init() {
	AddCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().StringVarP(&addUserCmdEnvironment, "environment", "e", "", "Environment of the micro integrator to which a new user should be added")
	addUserCmd.MarkFlagRequired("environment")
}

func handleAddUserCmdArguments(args []string) {
	printAddCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(addUserCmdLiteral))
	credentials.HandleMissingCredentials(addUserCmdEnvironment)
	startConsoleToAddUser(args[0])
}

func startConsoleToAddUser(userName string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Is " + userName + " an admin [y/N]: ")
	isAdmin, _ := reader.ReadString('\n')

	fmt.Printf("Enter password for " + userName + ": ")
	byteUserPassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	userPassword := string(byteUserPassword)
	fmt.Println()

	fmt.Printf("Re-Enter password for " + userName + ": ")
	byteUserConfirmationPassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	userConfirmPassword := string(byteUserConfirmationPassword)
	fmt.Println()

	fmt.Printf("Enter user store for " + userName + " default (primary): ")
	domain, _ := reader.ReadString('\n')
	domain = strings.TrimSuffix(domain, "\n")

	if userConfirmPassword == userPassword {
		executeAddNewUser(userName, userPassword, isAdmin, domain)
	} else {
		fmt.Println("Passwords are not matching.")
	}
}

func executeAddNewUser(userName, userPassword, isAdmin, domain string) {
	resp, err := impl.AddMIUser(addUserCmdEnvironment, userName, userPassword, isAdmin, domain)
	if err != nil {
		fmt.Println(utils.LogPrefixError+"Adding new user [ "+userName+" ]", err)
	} else {
		fmt.Println("Adding new user [ "+userName+" ] status:", resp)
	}
}
