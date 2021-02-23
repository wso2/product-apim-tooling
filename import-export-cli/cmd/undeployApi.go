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
var undeployAllGateways = true

// UndeployAPICmd command related usage info
const UndeployAPICmdLiteral = "api"
const undeployAPICmdShortDesc = "Undeploy API"

const undeployAPICmdLongDesc = "Undeploy an API revision from gateway environments"

const undeployAPICmdExamples = utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -rev 2 -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 6 -g Label1 Label2 Label3 -e production
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -r alice --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e)) are mandatory. 
If the flag (--gateway (-g)) is not provided, revision will be undeployed from all deployed gateway environments.`

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
		if undeployAPIGateway != "" {
			undeployAllGateways = false
		}
		gateways := generateGatewayArray(args, undeployAPIGateway)

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
			deployments, undeployAllGateways)
		if err != nil {
			utils.HandleErrorAndExit("Error while undeploying the API", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusCreated {
			fmt.Println(apiNameForStateChange + " API revision successfully undeployed from the gateways!")
		} else {
			fmt.Println("Error while undeploying the  API: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens to undeploy the API:" + preCommandErr.Error())
	}
}

//process the args array and create deployments array
func generateGatewayArray(args []string, initialGateway string) []utils.Deployment {
	//Since other flags does not use args[], gateways flag will own all the args
	var deployments = []utils.Deployment{{initialGateway, true}}
	if len(args) != 0 {
		for _, argument := range args {
			var deployment utils.Deployment
			deployment.Name = argument
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
	UndeployAPICmd.Flags().StringVarP(&undeployAPIGateway, "gateway", "g", "",
		"Gateway which the revision has to be undeployed")
	UndeployAPICmd.Flags().StringVarP(&undeployRevisionNum, "rev", "", "",
		"Revision number of the API to undeploy")
	UndeployAPICmd.Flags().StringVarP(&undeployAPIEnvironment, "environment", "e",
		"", "Environment of which the API should be undeployed")
	_ = UndeployAPICmd.MarkFlagRequired("name")
	_ = UndeployAPICmd.MarkFlagRequired("version")
	_ = UndeployAPICmd.MarkFlagRequired("rev")
	_ = UndeployAPICmd.MarkFlagRequired("environment")
}
