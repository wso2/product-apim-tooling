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

var deployCmdAPIName string
var deployCmdAPIVersion string
var deployCmdRevisionNum string
var deployCmdProvider string
var deployCmdEnvironment string
var deployCmdGateway string
var deployCmdHideOnDevportal bool

// DeployAPI command related usage info
const DeployAPICmdLiteral = "api"
const deployAPICmdShortDesc = "Deploy API"

const deployAPICmdLongDesc = "Deploy an API to the given gateway environment"

const deployAPICmdExamples = utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin --rev 1 -g Label1 -e dev
` + utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 6 -g Label1 Label2 Label3 -e production
` + utils.ProjectName + ` ` + DeployCmdLiteral + ` ` + DeployAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 --rev 2 -r admin -g Label1 -e production --hide-on-devportal
NOTE: All the 5 flags (--name (-n), --version (-v) , --rev and --gateway (-g) --environment (-e)) are mandatory.`

// DeployAPICmd represents the deploy API command
var DeployAPICmd = &cobra.Command{
	Use: DeployAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> " +
		"--rev <revision_number> --gateway <gateway-environment> " +
		"--environment <environment-from-which-the-api-should-be-deployed>)",
	Short:   deployAPICmdShortDesc,
	Long:    deployAPICmdLongDesc,
	Example: deployAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + DeployAPICmdLiteral + " called")
		cred, err := GetCredentials(deployCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		gateways := impl.GenerateGatewayArray(args, deployCmdGateway, deployCmdHideOnDevportal)
		executeDeployAPICmd(cred, gateways)
	},
}

func executeDeployAPICmd(credential credentials.Credential, deployments []utils.Deployment) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deployCmdEnvironment)
	if preCommandErr == nil {
		resp, err := impl.ReflectDeploymentChangeInGatewayEnv("deploy-revision", accessToken, deployCmdEnvironment,
			deployCmdAPIName, deployCmdAPIVersion, deployCmdProvider, deployCmdRevisionNum, deployments, false)
		if err != nil {
			utils.HandleErrorAndExit("Error while deploying the API", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusCreated {
			fmt.Println(apiNameForStateChange + " API state changed successfully!")
		} else {
			fmt.Println("Error while deploying the API: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens to deploy the API:" + preCommandErr.Error())
	}

}

// init using Cobra
func init() {
	DeployCmd.AddCommand(DeployAPICmd)
	DeployAPICmd.Flags().StringVarP(&deployCmdAPIName, "name", "n", "",
		"Name of the API to be deployed")
	DeployAPICmd.Flags().StringVarP(&deployCmdAPIVersion, "version", "v", "",
		"Version of the API to be deployed")
	DeployAPICmd.Flags().StringVarP(&deployCmdProvider, "provider", "r", "",
		"Provider of the API")
	DeployAPICmd.Flags().StringVarP(&deployCmdGateway, "gateway", "g", "",
		"Gateways which the revision has to be deployed")
	DeployAPICmd.Flags().StringVarP(&deployCmdRevisionNum, "rev", "", "",
		"Revision number of the API to be deployed")
	DeployAPICmd.Flags().BoolVar(&deployCmdHideOnDevportal, "hide-on-devportal", false,
		"Hide the gateway environment on devportal")
	DeployAPICmd.Flags().StringVarP(&deployCmdEnvironment, "environment", "e",
		"", "Environment to which the API should be deployed")
	_ = DeployAPICmd.MarkFlagRequired("name")
	_ = DeployAPICmd.MarkFlagRequired("version")
	_ = DeployAPICmd.MarkFlagRequired("gateway")
	_ = DeployAPICmd.MarkFlagRequired("rev")
	_ = DeployAPICmd.MarkFlagRequired("environment")
}
