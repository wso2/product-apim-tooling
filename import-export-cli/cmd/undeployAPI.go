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
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var undeployAPIName string
var undeployAPIVersion string
var undeployRevisionNum string
var undeployProvider string
var undeployAPIEnvironment string
var undeployAPIGateway string
var undeployAllGateways bool

// UndeployAPICmd command related usage info
const UndeployAPICmdLiteral = "api"
const undeployAPICmdShortDesc = "Undeploy API"

const undeployAPICmdLongDesc = "Undeploy an API to the given gateway environment"

const undeployAPICmdExamples = utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -g Label1 -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 6 --all-gateways -e production
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e) and one from --gateway (-g) or --all-gateways) are mandatory.`

// UndeployAPICmd represents the deploy API command
var UndeployAPICmd = &cobra.Command{
	Use: UndeployAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> " +
		"--rev <revision-number-of-the-api> --gateway <gateway-environment> " +
		"--environment <environment-from-which-the-api-should-be-undeployed>)",
	Short:   undeployAPICmdShortDesc,
	Long:    undeployAPICmdLongDesc,
	Example: undeployAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + UndeployAPICmdLiteral + " called")
		if undeployAPIGateway != "" || undeployAllGateways {
			gateways := impl.GenerateGatewayArray(args, undeployAPIGateway, true)

			cred, err := GetCredentials(undeployAPIEnvironment)
			if err != nil {
				utils.HandleErrorAndExit("Error getting credentials", err)
			}
			executeUndeployAPICmd(cred, gateways)
		} else {
			fmt.Println("Invalid Arguments. Atleast one gateway environment or --all-gateways " +
				"flag should be provided to undeploy the revision")
		}
	},
}

func executeUndeployAPICmd(credential credentials.Credential, deployments []utils.Deployment) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, undeployAPIEnvironment)
	if preCommandErr == nil {
		resp, err := impl.ReflectDeploymentChangeInGatewayEnv("undeploy-revision", accessToken,
			undeployAPIEnvironment, undeployAPIName, undeployAPIVersion, undeployProvider, undeployRevisionNum, deployments,
			undeployAllGateways)
		if err != nil {
			utils.HandleErrorAndExit("Error while undeploying the API", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusCreated {
			fmt.Println(apiNameForStateChange + " API artifact successfully undeployed from the gateway!")
		} else {
			fmt.Println("Error while undeploying the  API: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens to undeploy the API:" + preCommandErr.Error())
	}

}

// init using Cobra
func init() {
	UndeployCmd.AddCommand(UndeployAPICmd)
	UndeployAPICmd.Flags().StringVarP(&undeployAPIName, "name", "n", "",
		"Name of the API to be exported")
	UndeployAPICmd.Flags().StringVarP(&undeployAPIVersion, "version", "v", "",
		"Version of the API to be exported")
	UndeployAPICmd.Flags().StringVarP(&undeployProvider, "provider", "r", "",
		"Provider of the API")
	UndeployAPICmd.Flags().StringVarP(&undeployAPIGateway, "gateway", "g", "",
		"Gateway which the revision has to be undeployed")
	UndeployAPICmd.Flags().StringVarP(&undeployRevisionNum, "rev", "", "",
		"Revision number of the API to undeploy")
	UndeployAPICmd.Flags().BoolVar(&undeployAllGateways, "all-gateways", false,
		"Undeploy the revision from all the deployed gateways at once")
	UndeployAPICmd.Flags().StringVarP(&undeployAPIEnvironment, "environment", "e",
		"", "Environment of which the API should be undeployed")
	_ = UndeployAPICmd.MarkFlagRequired("name")
	_ = UndeployAPICmd.MarkFlagRequired("version")
	_ = UndeployAPICmd.MarkFlagRequired("rev")
	_ = UndeployAPICmd.MarkFlagRequired("environment")
}
