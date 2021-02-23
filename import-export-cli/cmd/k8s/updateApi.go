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

package k8s

import (
	"encoding/json"
	"fmt"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const K8sUpdateCmdLiteral = "update"
const k8sUpdateCmdShortDesc = "Update an API to the kubernetes cluster"
const k8sUpdateCmdLongDesc = `Update an existing API with Swagger file, project zip and API project in the kubernetes cluster. 
JSON, YAML, zip and API project folder formats are accepted.`
const k8sUpdateCmdExamples = utils.ProjectName + " " + K8sCmdLiteral + " " + K8sUpdateCmdLiteral + " " + AddApiCmdLiteral +
	" " + `-n petstore --file=./Swagger.json --namespace=wso2

` + utils.ProjectName + " " + K8sCmdLiteral + " " + K8sUpdateCmdLiteral + " " + AddApiCmdLiteral +
	" " + `-n petstore --file=./product-apim-tooling/import-export-cli/build/target/apictl/myapi --namespace=wso2`

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
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
		handleUpdateApi()
	},
}

func handleUpdateApi() {
	var errMsg string
	flagApiName = strings.ToLower(flagApiName)
	getApiCr, getApiErr:= k8sUtils.GetCommandOutput(
		k8sUtils.Kubectl, k8sUtils.K8sGet, k8sUtils.ApiOpCrdApi, flagApiName, "-n", flagNamespace, "-o", "json")
	if getApiErr != nil {

		if flagNamespace != "" {
			errMsg = fmt.Sprintf("Could not find the API \"%s\" in the namespace \"%s\"",
				flagApiName, flagNamespace)
		} else {
			errMsg = fmt.Sprintf("Could not find the API \"%s\"", flagApiName)
		}
		utils.HandleErrorAndExit(errMsg, nil)
	}
	var apiCr map[string]interface{}
	_ = json.Unmarshal([]byte(getApiCr), &apiCr)
	swaggerCmName := apiCr["spec"].(map[string]interface{})["swaggerConfigMapName"].(string)
	timestampSuffix := fmt.Sprint(time.Now().Unix())
	handleAddApi("-" + strings.ToLower(timestampSuffix))
	deleteApiErr := k8sUtils.ExecuteCommand(
		k8sUtils.Kubectl, k8sUtils.K8sDelete, "cm", swaggerCmName, "-n", flagNamespace)
	if deleteApiErr != nil {
		if flagNamespace != "" {
			errMsg = fmt.Sprintf("Could not find the config map \"%s\" in the namespace \"%s\"",
				flagApiName, flagNamespace)
		} else {
			errMsg = fmt.Sprintf("Could not find the config map \"%s\"", flagApiName)
		}
		utils.HandleErrorAndExit(errMsg, nil)
	}

}

func init() {
	UpdateCmd.AddCommand(updateApiCmd)
	updateApiCmd.Flags().StringVarP(&flagApiName, "name", "n", "", "Name of the API")
	updateApiCmd.Flags().StringVarP(&flagSwaggerFilePath, "file", "f", "",
		"Path to swagger, zip file or API project")
	updateApiCmd.Flags().StringVar(&flagNamespace, "namespace", "", "namespace of API")
	_ = updateApiCmd.MarkFlagRequired("name")
	_ = updateApiCmd.MarkFlagRequired("file")
}
