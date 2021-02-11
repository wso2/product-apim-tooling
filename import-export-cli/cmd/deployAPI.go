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

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var deployAPIName string
var deployAPIVersion string
var deployRevisionNum string
var deployProvider string
var deployAPIEnvironment string
var deployAPIGateway string

// DeployAPI command related usage info
const DeployAPICmdLiteral = "api"
const deployAPICmdShortDesc = "Deploy API"

const deployAPICmdLongDesc = "Deploy an API to the given gateway environment"

const deployAPICmdExamples = utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -g Label1 -e dev
` + utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 6 -r admin -g Label1 -e production
` + utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 2 -r admin -g Label1 -e production
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e) --gateway (-g) ) are mandatory. 
If --rev is not provided, working copy of the API will create a new revision and deploy to the provided gateway.`

// DeployAPICmd represents the deploy API command
var DeployAPICmd = &cobra.Command{
	Use: DeployAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> " +
		"--gateway <gateway-environment> --environment <environment-from-which-the-api-should-be-deployed>)",
	Short:   deployAPICmdShortDesc,
	Long:    deployAPICmdLongDesc,
	Example: deployAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportAPICmdLiteral + " called")
		if exportRevisionNum == "" {
			fmt.Println("A Revision number is not provided. Working copy will be deployed to the gateway.")
		}

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeDeployAPICmd(cred)
	},
}

func executeDeployAPICmd(credential credentials.Credential) {

}

// init using Cobra
func init() {
	DeployCmd.AddCommand(DeployAPICmd)
	DeployAPICmd.Flags().StringVarP(&deployAPIName, "name", "n", "",
		"Name of the API to be exported")
	DeployAPICmd.Flags().StringVarP(&deployAPIVersion, "version", "v", "",
		"Version of the API to be exported")
	DeployAPICmd.Flags().StringVarP(&deployProvider, "provider", "r", "",
		"Provider of the API")
	DeployAPICmd.Flags().StringVarP(&deployAPIGateway, "gateway", "g", "",
		"Gateway which the API has to be deployed")
	DeployAPICmd.Flags().StringVarP(&deployRevisionNum, "rev", "", "",
		"Revision number of the API to be exported")
	DeployAPICmd.Flags().StringVarP(&deployAPIEnvironment, "environment", "e",
		"", "Environment to which the API should be deployed")
	_ = ExportAPICmd.MarkFlagRequired("name")
	_ = ExportAPICmd.MarkFlagRequired("version")
	_ = ExportAPICmd.MarkFlagRequired("gateway")
	_ = ExportAPICmd.MarkFlagRequired("environment")
}
