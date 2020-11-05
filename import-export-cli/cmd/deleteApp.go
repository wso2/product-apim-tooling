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

package cmd

import (
	"fmt"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

var deleteAppEnvironment string
var deleteAppName string
var deleteAppOwner string

// DeleteApp command related usage info
const deleteAppCmdLiteral = "app"
const deleteAppCmdShortDesc = "Delete App"
const deleteAppCmdLongDesc = "Delete an Application from an environment"

const deleteAppCmdExamples = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAppCmdLiteral + ` -n TestApplication -o admin -e dev
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAppCmdLiteral + ` -n SampleApplication -e production
NOTE: Both the flags (--name (-n), and --environment (-e)) are mandatory and the flag --owner (-o) is optional.`

// DeleteAppCmd represents the delete app command
var DeleteAppCmd = &cobra.Command{
	Use: deleteAppCmdLiteral + " (--name <name-of-the-application> --owner <owner-of-the-application> --environment " +
		"<environment-from-which-the-application-should-be-deleted>)",
	Short:   deleteAppCmdShortDesc,
	Long:    deleteAppCmdLongDesc,
	Example: deleteAppCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAppCmdLiteral + " called")
		cred, err := GetCredentials(deleteAppEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeDeleteAppCmd(cred)
	},
}

// executeDeleteAppCmd executes the delete app command
func executeDeleteAppCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteAppEnvironment)
	if preCommandErr == nil {
		if deleteAppOwner == "" {
			deleteAppOwner = credential.Username
		}
		resp, err := impl.DeleteApplication(accessToken, deleteAppEnvironment, deleteAppName, deleteAppOwner)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting Application ", err)
		}
		impl.PrintDeleteAppResponse(resp, err)
	} else {
		// Error deleting Application
		fmt.Println("Error getting OAuth tokens while deleting Application:" + preCommandErr.Error())
	}
}

// Init using Cobra
func init() {
	DeleteCmd.AddCommand(DeleteAppCmd)
	DeleteAppCmd.Flags().StringVarP(&deleteAppName, "name", "n", "",
		"Name of the Application to be deleted")
	DeleteAppCmd.Flags().StringVarP(&deleteAppOwner, "owner", "o", "",
		"Owner of the Application to be deleted")
	DeleteAppCmd.Flags().StringVarP(&deleteAppEnvironment, "environment", "e",
		"", "Environment from which the Application should be deleted")
	// Mark required flags
	_ = DeleteAppCmd.MarkFlagRequired("name")
	_ = DeleteAppCmd.MarkFlagRequired("environment")
}
