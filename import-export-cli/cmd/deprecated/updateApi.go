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
package deprecated

import (
	"fmt"
	"strings"
	"time"

	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

const updateCmdLiteral = "update"
const updateCmdShortDesc = "Update an API to the kubernetes cluster"
const updateCmdLongDesc = `Update an existing API with Swagger file in the kubernetes cluster. JSON and YAML formats are accepted.`
const updateCmdExamples = utils.ProjectName + " " + updateCmdLiteral + " " + addApiCmdLiteral + " " + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2

` + utils.ProjectName + " " + updateCmdLiteral + " " + addApiCmdLiteral + " " + `-n petstore --from-file=./product-apim-tooling/import-export-cli/build/target/apictl/myapi --replicas=1 --namespace=wso2`

// updateCmd represents the update command
var updateCmdDeprecated = &cobra.Command{
	Use:        updateCmdLiteral,
	Short:      updateCmdShortDesc,
	Long:       updateCmdLongDesc,
	Example:    updateCmdExamples,
	Deprecated: "instead use \"" + cmd.K8sCmdLiteral + " " + cmd.K8sUpdateCmdLiteral + "\".",
}

// updateApiCmd represents the updateApi command
var updateApiCmdDeprecated = &cobra.Command{
	Use:        addApiCmdLiteral,
	Short:      addApiCmdShortDesc,
	Long:       addApiLongDesc,
	Example:    addApiExamples,
	Deprecated: "instead use \"" + cmd.K8sCmdLiteral + " " + cmd.K8sUpdateCmdLiteral + " " + cmd.AddApiCmdLiteral + "\".",

	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + updateCmdLiteral + " called")
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
	cmd.RootCmd.AddCommand(updateCmdDeprecated)
	updateCmdDeprecated.AddCommand(updateApiCmdDeprecated)
	updateApiCmdDeprecated.Flags().StringVarP(&flagApiName, "name", "n", "", "Name of the API")
	updateApiCmdDeprecated.Flags().StringArrayVarP(&flagSwaggerFilePaths, "from-file", "f", []string{}, "Path to swagger file")
	updateApiCmdDeprecated.Flags().IntVar(&flagReplicas, "replicas", 1, "replica set")
	updateApiCmdDeprecated.Flags().StringVar(&flagNamespace, "namespace", "", "namespace of API")
	updateApiCmdDeprecated.Flags().StringVarP(&flagApiVersion, "version", "v", "", "Property to override the existing docker image with same name and version")
	updateApiCmdDeprecated.Flags().StringVarP(&flagApiMode, "mode", "m", "",
		fmt.Sprintf("Property to override the deploying mode. Available modes: %v, %v", utils.PrivateJetModeConst, utils.SidecarModeConst))
}
