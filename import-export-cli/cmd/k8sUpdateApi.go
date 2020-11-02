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
	"strings"
	"time"

	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

const K8sUpdateCmdLiteral = "update"
const k8sUpdateCmdShortDesc = "Update an API to the kubernetes cluster"
const k8sUpdateCmdLongDesc = `Update an existing API with Swagger file in the kubernetes cluster. JSON and YAML formats are accepted.`
const k8sUpdateCmdExamples = utils.ProjectName + " " + K8sCmdLiteral + " " + K8sUpdateCmdLiteral + " " + AddApiCmdLiteral + " " + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2

` + utils.ProjectName + " " + K8sCmdLiteral + " " + K8sUpdateCmdLiteral + " " + AddApiCmdLiteral + " " + `-n petstore --from-file=./product-apim-tooling/import-export-cli/build/target/apictl/myapi --replicas=1 --namespace=wso2`

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     K8sUpdateCmdLiteral,
	Short:   k8sUpdateCmdShortDesc,
	Long:    k8sUpdateCmdLongDesc,
	Example: k8sUpdateCmdExamples,
}

// updateApiCmd represents the updateApi command
var updateApiCmd = &cobra.Command{
	Use:     AddApiCmdLiteral,
	Short:   addApiCmdShortDesc,
	Long:    addApiLongDesc,
	Example: addApiExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + K8sUpdateCmdLiteral + " called")
		validateAddApiCommand()

		// check the existence of the API
		getApiErr := k8sUtils.ExecuteCommandWithoutOutputs(
			k8sUtils.Kubectl, k8sUtils.K8sGet, k8sUtils.ApiOpCrdApi, flagApiName, "-n", flagNamespace)
		if getApiErr != nil {
			var errMsg string
			if flagNamespace != "" {
				errMsg = fmt.Sprintf("Could not find the API \"%s\" in the namespace \"%s\"",
					flagApiName, flagNamespace)
			} else {
				errMsg = fmt.Sprintf("Could not find the API \"%s\"", flagApiName)
			}
			utils.HandleErrorAndExit(errMsg, nil)
		}

		//get current timestamp
		timestampSuffix := time.Now().Format("2Jan2006150405")
		handleAddApi("-" + strings.ToLower(timestampSuffix))
	},
}

func init() {
	K8sCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(updateApiCmd)
	updateApiCmd.Flags().StringVarP(&flagApiName, "name", "n", "", "Name of the API")
	updateApiCmd.Flags().StringArrayVarP(&flagSwaggerFilePaths, "from-file", "f", []string{}, "Path to swagger file")
	updateApiCmd.Flags().IntVar(&flagReplicas, "replicas", 1, "replica set")
	updateApiCmd.Flags().StringVar(&flagNamespace, "namespace", "", "namespace of API")
	updateApiCmd.Flags().StringVarP(&flagApiVersion, "version", "v", "", "Property to override the existing docker image with same name and version")
	updateApiCmd.Flags().StringVarP(&flagApiMode, "mode", "m", "",
		fmt.Sprintf("Property to override the deploying mode. Available modes: %v, %v", utils.PrivateJetModeConst, utils.SidecarModeConst))
}
