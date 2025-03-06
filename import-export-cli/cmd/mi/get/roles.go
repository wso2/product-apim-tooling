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

var getRoleCmdEnvironment string
var getRoleCmdFormat string
var getRoleCmdDomain string

const getRoleCmdLiteral = "roles [role-name]"

const getRoleCmdShortDesc = "Get information about roles"
const getRoleCmdLongDesc = "Get information about the roles in primary and secondary user stores.\n" +
	"List all roles of the Micro Integrator in the environment specified by the flag --environment, -e"

var getRoleCmdExamples = "To list all the roles\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getRoleCmdLiteral) + " -e dev\n" +
	"To get details about a role by providing the role name\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getRoleCmdLiteral) + " [role-name] -e dev\n" +
	"To get details about a role in a secondary user store\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getRoleCmdLiteral) + " [role-name] -d [domain] -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var getRoleCmd = &cobra.Command{
	Use:     getRoleCmdLiteral,
	Short:   getRoleCmdShortDesc,
	Long:    getRoleCmdLongDesc,
	Example: getRoleCmdExamples,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			var errMessage = "accepts at most 1 arg(s), received " + fmt.Sprint(len(args))
			return errors.New(errMessage)
		}
		return nil
	},
	Deprecated: "instead refer to https://mi.docs.wso2.com/en/latest/observe-and-manage/managing-integrations-with-micli/ for updated usage.",
	Run: func(cmd *cobra.Command, args []string) {
		handleGetRoleCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getRoleCmd)
	setEnvFlag(getRoleCmd, &getRoleCmdEnvironment)
	setFormatFlag(getRoleCmd, &getRoleCmdFormat)
	getRoleCmd.Flags().StringVarP(&getRoleCmdDomain, "domain", "d", "", "Filter roles by domain")
}

func handleGetRoleCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getRoleCmdLiteral))
	credentials.HandleMissingCredentials(getRoleCmdEnvironment)
	if len(args) == 1 {
		var role = args[0]
		executeShowRole(role)
	} else {
		executeListRoles()
	}
}

func executeShowRole(role string) {
	roleInfo, err := impl.GetRoleInfo(getRoleCmdEnvironment, role, getRoleCmdDomain)
	if err == nil {
		impl.PrintRoleDetails(roleInfo, getRoleCmdFormat)
	} else {
		printErrorForArtifact("roles", role, err)
	}
}

func executeListRoles() {
	roleList, err := impl.GetRoleList(getRoleCmdEnvironment)
	if err == nil {
		impl.PrintRoleList(roleList, getRoleCmdFormat)
	} else {
		printErrorForArtifactList("roles", err)
	}
}


