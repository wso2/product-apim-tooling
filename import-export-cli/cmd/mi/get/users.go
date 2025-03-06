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

package get

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getUserCmdEnvironment string
var getUserCmdFormat string
var getUserCmdPattern string
var getUserCmdRole string
var getUserCmdDomain string

const getUserCmdLiteral = "users [user-name]"

const getUserCmdShortDesc = "Get information about users"
const getUserCmdLongDesc = "Get information about the users filtered by username pattern and role.\n" +
	"If not provided list all users of the Micro Integrator in the environment specified by the flag --environment, -e"

var getUserCmdExamples = "Example:\n" +
	"To list all the users\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getUserCmdLiteral) + " -e dev\n" +
	"To get the list of users with specific role\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getUserCmdLiteral) + " -r [role-name] -e dev\n" +
	"To get the list of users with a username matching with the wild card Ex: \"*mi*\" matches with \"admin\"\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getUserCmdLiteral) + " -p [pattern] -e dev\n" +
	"To get details about a user by providing the user-id\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getUserCmdLiteral) + " [user-id] -e dev\n" +
	"To get details about a user in a secondary user store\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getUserCmdLiteral) + " [user-id] -d [domain] -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var getUserCmd = &cobra.Command{
	Use:     getUserCmdLiteral,
	Short:   getUserCmdShortDesc,
	Long:    getUserCmdLongDesc,
	Example: getUserCmdExamples,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			var errMessage = "accepts at most 1 arg(s), received " + fmt.Sprint(len(args))
			return errors.New(errMessage)
		} else if len(args) == 1 {
			if isInvalidFlagArgCombination(args[0], getUserCmdPattern, getUserCmdRole) {
				var errMessage = "[user-id] arg cannot be used with -p or -r flags"
				return errors.New(errMessage)
			}
		}
		return nil
	},
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handleGetUserCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getUserCmd)
	setEnvFlag(getUserCmd, &getUserCmdEnvironment)
	setFormatFlag(getUserCmd, &getUserCmdFormat)
	getUserCmd.Flags().StringVarP(&getUserCmdRole, "role", "r", "", "Filter users by role")
	getUserCmd.Flags().StringVarP(&getUserCmdPattern, "pattern", "p", "", "Filter users by regex")
	getUserCmd.Flags().StringVarP(&getUserCmdDomain, "domain", "d", "", "Filter users by domain")
}

func handleGetUserCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getUserCmdLiteral))
	credentials.HandleMissingCredentials(getUserCmdEnvironment)
	if len(args) == 1 {
		var userID = args[0]
		executeShowUser(userID)
	} else {
		executeListUsers()
	}
}

func executeShowUser(userID string) {
	userInfo, err := impl.GetUserInfo(getUserCmdEnvironment, userID, getUserCmdDomain)
	if err == nil {
		impl.PrintUserDetails(userInfo, getUserCmdFormat)
	} else {
		printErrorForArtifact("users", userID, err)
	}
}

func executeListUsers() {
	userList, err := impl.GetUserList(getUserCmdEnvironment, getUserCmdRole, getUserCmdPattern)
	if err == nil {
		impl.PrintUserList(userList, getUserCmdFormat)
	} else {
		printErrorForArtifactList("users", err)
	}
}

func isInvalidFlagArgCombination(userIDArg, patternFlag, roleFlag string) bool {
	return userIDArg != "" && (getUserCmdPattern != "" || getUserCmdRole != "")
}
