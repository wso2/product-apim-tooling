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
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var undeployAPIName string
var undeployAPIVersion string
var undeployRevisionNum string
var undeployProvider string
var undeployAPIEnvironment string
var undeployAPICmdAPIGatewayEnvs []string
var undeployAllGatewayEnvs = true

// UndeployAPICmd command related usage info
const UndeployAPICmdLiteral = "api"
const undeployAPICmdShortDesc = "Undeploy API"

const undeployAPICmdLongDesc = "Undeploy an API revision from gateway environments"

const undeployAPICmdExamples = utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 --rev 2 -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 6 -g Label1 -g Label2 -g Label3 -e production
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -r alice --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e)) are mandatory. 
If the flag (--gateway-env (-g)) is not provided, revision will be undeployed from all deployed gateway environments.`

// UndeployAPICmd represents the deploy API command
var UndeployAPICmd = &cobra.Command{
	Use: UndeployAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> " +
		"--rev <revision-number-of-the-api> --gateway-env <gateway-environment> " +
		"--environment <environment-from-which-the-api-should-be-undeployed>)",
	Short:   undeployAPICmdShortDesc,
	Long:    undeployAPICmdLongDesc,
	Example: undeployAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + UndeployAPICmdLiteral + " called")
		if len(undeployAPICmdAPIGatewayEnvs) > 0 {
			undeployAllGatewayEnvs = false
		}
		gateways := generateGatewayEnvsArray(undeployAPICmdAPIGatewayEnvs)

		cred, err := GetCredentials(undeployAPIEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeUndeployAPICmd(cred, gateways)
	},
}

func executeUndeployAPICmd(credential credentials.Credential, deployments []utils.Deployment) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, undeployAPIEnvironment)
	if preCommandErr == nil {
		resp, err := impl.UndeployRevisionFromGateways(accessToken,
			undeployAPIEnvironment, undeployAPIName, undeployAPIVersion, undeployProvider, undeployRevisionNum,
			deployments, undeployAllGatewayEnvs)
		if err != nil {
			utils.HandleErrorAndExit("Error while undeploying the API", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusCreated {
			fmt.Println("Revision " + undeployRevisionNum + " of API " + undeployAPIName + "_" + undeployAPIVersion +
				" successfully undeployed from the specified gateways environments")
		} else {
			fmt.Println("Error while undeploying the  API: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens to undeploy the API:" + preCommandErr.Error())
	}
}

// Process the gatewayEnvs array and create deployments array
func generateGatewayEnvsArray(gatewayEnvs []string) []utils.Deployment {
	var deployments []utils.Deployment
	if len(gatewayEnvs) != 0 {
		for _, gatewayEnv := range gatewayEnvs {
			var deployment utils.Deployment
			deployment.Name = gatewayEnv
			deployment.DisplayOnDevportal = true
			deployments = append(deployments, deployment)
		}
	}
	return deployments
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
	UndeployAPICmd.Flags().StringSliceVarP(&undeployAPICmdAPIGatewayEnvs, "gateway-env", "g", []string{},
		"Gateway environment which the revision has to be undeployed")
	UndeployAPICmd.Flags().StringVarP(&undeployRevisionNum, "rev", "", "",
		"Revision number of the API to undeploy")
	UndeployAPICmd.Flags().StringVarP(&undeployAPIEnvironment, "environment", "e",
		"", "Environment of which the API should be undeployed")
	_ = UndeployAPICmd.MarkFlagRequired("name")
	_ = UndeployAPICmd.MarkFlagRequired("version")
	_ = UndeployAPICmd.MarkFlagRequired("rev")
	_ = UndeployAPICmd.MarkFlagRequired("environment")
}
