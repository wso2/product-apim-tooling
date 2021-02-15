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

var undeployAPIProductName string
var undeployAPIProductRevisionNum string
var undeployAPIProductProvider string
var undeployAPIProductEnvironment string
var undeployAPIProductGateway string
var undeployAPIProductAllGateways bool

// Undeploy API Product command related usage info
const UndeployAPIProductCmdLiteral = "api-product"
const undeployAPIProductCmdShortDesc = "Undeploy API Product"

const undeployAPIProductmdLongDesc = "Undeploy an API Product revision from the given gateway environment"

const undeployAPIProductCmdExamples = utils.
	ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPIProductCmdLiteral + ` -n TwitterAPIProduct -v 1.0.0 -r admin -g Label1 -e dev
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPIProductCmdLiteral + ` -n StoreProduct -v 2.1.0 --rev 6 --all-gateways -e production
` + utils.ProjectName + ` ` + UndeployCmdLiteral + ` ` + UndeployAPIProductCmdLiteral + ` -n FacebookProduct -v 2.1.0 --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e)) and one from (--gateway (-g) or --all-gateways) are mandatory.`

// UndeployAPIProductCmd represents the deploy API command
var UndeployAPIProductCmd = &cobra.Command{
	Use: UndeployAPIProductCmdLiteral + " (--name <name-of-the-api-product> --version <version-of-the-api-product> " +
		"--rev<revision-number-of-the-api-product> --gateway <gateway-environment> " +
		"--environment <environment-from-which-the-api-product-should-be-undeployed>)",
	Short:   undeployAPIProductCmdShortDesc,
	Long:    undeployAPIProductmdLongDesc,
	Example: undeployAPIProductCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + UndeployAPIProductCmdLiteral + " called")
		if undeployAPIProductGateway != "" || undeployAPIProductAllGateways {
			gateways := generateGatewayArray(args, undeployAPIProductGateway)

			cred, err := GetCredentials(undeployAPIProductEnvironment)
			if err != nil {
				utils.HandleErrorAndExit("Error getting credentials", err)
			}
			executeUndeployAPIProductCmd(cred, gateways)
		} else {
			fmt.Println("Invalid Arguments. Atleast one gateway environment or --all-gateways " +
				"flag should be provided to undeploy the revision")
		}
	},
}

func executeUndeployAPIProductCmd(credential credentials.Credential, deployments []utils.Deployment) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, undeployAPIProductEnvironment)
	if preCommandErr == nil {
		resp, err := impl.UndeployAPIProductRevisionFromGateways(accessToken,
			undeployAPIProductEnvironment, undeployAPIProductName, undeployAPIProductProvider,
			undeployAPIProductRevisionNum, deployments, undeployAPIProductAllGateways)
		if err != nil {
			utils.HandleErrorAndExit("Error while undeploying the API Product", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusCreated {
			fmt.Println(undeployAPIProductName + " API Product revision successfully undeployed from the gateways!")
		} else {
			fmt.Println("Error while undeploying the  APIProduct: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		fmt.Println("Error getting OAuth tokens to undeploy the API Product:" + preCommandErr.Error())
	}
}

// init using Cobra
func init() {
	UndeployCmd.AddCommand(UndeployAPIProductCmd)
	UndeployAPIProductCmd.Flags().StringVarP(&undeployAPIProductName, "name", "n", "",
		"Name of the API Product to be exported")
	UndeployAPIProductCmd.Flags().StringVarP(&undeployAPIProductProvider, "provider", "r", "",
		"Provider of the API")
	UndeployAPIProductCmd.Flags().StringVarP(&undeployAPIProductGateway, "gateway", "g", "",
		"Gateway which the revision has to be undeployed")
	UndeployAPIProductCmd.Flags().StringVarP(&undeployAPIProductRevisionNum, "rev", "", "",
		"Revision number of the API Product to undeploy")
	UndeployAPIProductCmd.Flags().BoolVar(&undeployAPIProductAllGateways, "all-gateways", false,
		"Undeploy the revision from all the deployed gateways at once")
	UndeployAPIProductCmd.Flags().StringVarP(&undeployAPIProductEnvironment, "environment", "e",
		"", "Environment of which the API Product should be undeployed")
	_ = UndeployAPIProductCmd.MarkFlagRequired("name")
	_ = UndeployAPIProductCmd.MarkFlagRequired("version")
	_ = UndeployAPIProductCmd.MarkFlagRequired("rev")
	_ = UndeployAPIProductCmd.MarkFlagRequired("environment")
}
