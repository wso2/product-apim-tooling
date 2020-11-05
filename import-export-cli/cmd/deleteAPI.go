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

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var deleteAPIEnvironment string
var deleteAPIName string
var deleteAPIVersion string
var deleteAPIProvider string

// DeleteAPI command related usage info
const deleteAPICmdLiteral = "api"
const deleteAPICmdShortDesc = "Delete API"
const deleteAPICmdLongDesc = "Delete an API from an environment"

const deleteAPICmdExamplesDefault = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -e production
NOTE: The 3 flags (--name (-n), --version (-v), and --environment (-e)) are mandatory.`

// DeleteAPICmd represents the delete api command
var DeleteAPICmd = &cobra.Command{
	Use: deleteAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment " +
		"<environment-from-which-the-api-should-be-deleted>)",
	Short:              deleteAPICmdShortDesc,
	Long:               deleteAPICmdLongDesc,
	Example:            deleteAPICmdExamplesDefault,
	DisableFlagParsing: isK8sEnabled(),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAPICmdLiteral + " called")
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if configVars.Config.KubernetesMode {
			k8sArgs := []string{k8sUtils.K8sDelete, k8sUtils.ApiOpCrdApi}
			k8sArgs = append(k8sArgs, args...)
			executeKubernetes(k8sArgs...)
		} else {
			cred, err := GetCredentials(deleteAPIEnvironment)
			if err != nil {
				utils.HandleErrorAndExit("Error getting credentials ", err)
			}
			executeDeleteAPICmd(cred)
		}
	},
}

// executeDeleteAPICmd executes the delete api command
func executeDeleteAPICmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteAPIEnvironment)
	if preCommandErr == nil {
		resp, err := impl.DeleteAPI(accessToken, deleteAPIEnvironment, deleteAPIName, deleteAPIVersion, deleteAPIProvider)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting API ", err)
		}
		impl.PrintDeleteAPIResponse(resp, err)
	} else {
		// Error deleting API
		fmt.Println("Error getting OAuth tokens while deleting API:" + preCommandErr.Error())
	}
}

// Init using Cobra
func init() {
	DeleteCmd.AddCommand(DeleteAPICmd)
	DeleteAPICmd.Flags().StringVarP(&deleteAPIName, "name", "n", "",
		"Name of the API to be deleted")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIVersion, "version", "v", "",
		"Version of the API to be deleted")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIProvider, "provider", "r", "",
		"Provider of the API to be deleted")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIEnvironment, "environment", "e",
		"", "Environment from which the API should be deleted")

	// fetches the main-config.yaml file silently; i.e. if it's not created, ignore the error and assume that
	//	this is the default mode.
	configVars := utils.GetMainConfigFromFileSilently(utils.MainConfigFilePath)
	if configVars == nil || !configVars.Config.KubernetesMode {
		// Mark required flags
		_ = DeleteAPICmd.MarkFlagRequired("name")
		_ = DeleteAPICmd.MarkFlagRequired("version")
		_ = DeleteAPICmd.MarkFlagRequired("environment")
	}
}
