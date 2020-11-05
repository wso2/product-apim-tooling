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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// keys command related Info
const GetKeysCmdLiteral = "keys"
const getKeysCmdShortDesc = "Generate access token to invoke the API or API Product"
const getKeysCmdLongDesc = `Generate JWT token to invoke the API or API Product by subscribing to a default application for testing purposes`
const getKeysCmdExamples = utils.ProjectName + " " + GetCmdLiteral + " " + GetKeysCmdLiteral + ` -n TwitterAPI -v 1.0.0 -e dev --provider admin
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory.
You can override the default token endpoint using --token (-t) optional flag providing a new token endpoint`

var keyGenEnv string
var apiName string
var apiVersion string
var apiProvider string
var keyGenTokenEndpoint string

var getKeysCmd = &cobra.Command{
	Use:     GetKeysCmdLiteral,
	Short:   getKeysCmdShortDesc,
	Long:    getKeysCmdLongDesc,
	Example: getKeysCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetKeysCmdLiteral + " called")
		cred, err := GetCredentials(keyGenEnv)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		utils.Logln(utils.LogPrefixInfo + "Retrieved credentials of the environment successfully")
		//Calling the DCR endpoint to get the credentials of the env
		cred.ClientId, cred.ClientSecret, err = impl.CallDCREndpoint(cred, keyGenEnv)
		//If the DCR call fails exit with the error
		if err != nil {
			utils.HandleErrorAndExit("Internal error occurred", err)
		}
		utils.Logln(utils.LogPrefixInfo + "Called DCR endpoint successfully")
		impl.GetKeys(cred, keyGenEnv, apiName, apiVersion, apiProvider, keyGenTokenEndpoint)
	},
}

//init function to add the cli command to the root command
func init() {
	GetCmd.AddCommand(getKeysCmd)
	getKeysCmd.Flags().StringVarP(&keyGenEnv, "environment", "e", "", "Key generation environment")
	getKeysCmd.Flags().StringVarP(&apiName, "name", "n", "", "API or API Product to generate keys")
	getKeysCmd.Flags().StringVarP(&apiVersion, "version", "v", "", "Version of the API")
	getKeysCmd.Flags().StringVarP(&apiProvider, "provider", "r", "", "Provider of the API or API Product")
	getKeysCmd.Flags().StringVarP(&keyGenTokenEndpoint, "token", "t", "", "Token endpoint URL of Environment")
	_ = getKeysCmd.MarkFlagRequired("name")
	_ = getKeysCmd.MarkFlagRequired("environment")
}
