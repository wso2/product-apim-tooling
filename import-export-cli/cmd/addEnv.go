/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var flagTokenEndpoint string
var flagRegistrationEndpoint string
var flagAPIManagerEndpoint string

// addEnvCmd represents the addEnv command
var addEnvCmd = &cobra.Command{
	Use:   "add-env",
	Short: utils.AddEnvCmdShortDesc,
	Long: utils.AddEnvCmdLongDesc + utils.AddEnvCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "add-env called")
	},
}

func addEnv(apimEndpoint string, regEndpoint string, tokenEndpoint string){

}

func init() {
	RootCmd.AddCommand(addEnvCmd)

	addEnvCmd.Flags().StringVar(&flagAPIManagerEndpoint, "apim", "",
		"API Manager endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagTokenEndpoint, "token", "",
		"Token endpoint for the environment")
	addEnvCmd.Flags().StringVar(&flagRegistrationEndpoint, "registration", "",
		"Registration endpoint for the environment")
}
