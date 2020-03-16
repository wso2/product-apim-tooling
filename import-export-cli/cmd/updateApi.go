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
	"errors"
	"fmt"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var updateflagApiName string
var updateflagSwaggerFilePath string
var updateflagReplicas int
var updateflagNamespace string

const updateCmdLiteral = "update"
const updateCmdShortDesc = "Update an API to the kubernetes cluster"
const updateCmdLongDesc = `Update an existing API with  Swagger file in the kubernetes cluster. JSON and YAML formats are accepted.`
const updateCmdExamples = utils.ProjectName + " " + updateCmdLiteral + " " + addApiCmdLiteral + " " + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2

` + utils.ProjectName + " " + updateCmdLiteral + " " + addApiCmdLiteral + " " + `-n petstore --from-file=./product-apim-tooling/import-export-cli/build/target/apictl/myapi --replicas=1 --namespace=wso2`

var updatedInterceptorConfName string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     updateCmdLiteral,
	Short:   updateCmdShortDesc,
	Long:    updateCmdLongDesc,
	Example: updateCmdExamples,
}

// updateApiCmd represents the updateApi command
var updateApiCmd = &cobra.Command{
	Use:     addApiCmdLiteral,
	Short:   addApiCmdShortDesc,
	Long:    addApiLongDesc,
	Example: addApiExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + updateCmdLiteral + " called")
		//check mode set to kubernetes
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if configVars.Config.KubernetesMode {
			if updateflagApiName == "" && updateflagSwaggerFilePath == "" {
				utils.HandleErrorAndExit("Required flags are missing. API name and swagger file paths are required",
					errors.New("required flags missing"))
			} else {
				//get current timestamp
				timestamp := time.Now().Format("20060102150405")
				//create new configmap with updated swagger file
				updateConfigMapName := updateflagApiName + "-swagger-up" + "-" + timestamp
				fi, err := os.Stat(updateflagSwaggerFilePath)
				if err != nil {
					fmt.Println(err)
					return
				}
				switch mode := fi.Mode(); {
				case mode.IsDir():
					//get swagger definition
					swaggerPath := filepath.Join(updateflagSwaggerFilePath, filepath.FromSlash("Meta-information/swagger.yaml"))
					//creating kubernetes configmap with swagger definition
					fmt.Println("creating configmap with swagger definition")
					errConf := createConfigMapWithNamespace(updateConfigMapName, swaggerPath, updateflagNamespace, k8sUtils.K8sCreate)
					if errConf != nil {
						utils.HandleErrorAndExit("Error creating configmap", err)
					}
					//handle interceptors
					updatedInterceptorConfName = updateflagApiName + "-interceptors-up-" + timestamp
					handleInterceptors(updatedInterceptorConfName, updateflagSwaggerFilePath, "create", updateflagNamespace)

				case mode.IsRegular():
					//creating kubernetes configmap with swagger definition
					fmt.Println("creating configmap with swagger definition")
					err := createConfigMapWithNamespace(updateConfigMapName, updateflagSwaggerFilePath, updateflagNamespace, k8sUtils.K8sCreate)
					if err != nil {
						utils.HandleErrorAndExit("Error creating configmap", err)
					}
				}
				//update the API
				fmt.Println("updating the API Kind")
				createAPI(updateflagApiName, updateflagNamespace, updateConfigMapName, updateflagReplicas, timestamp, updatedInterceptorConfName, false)
			}
		} else {
			utils.HandleErrorAndExit("set mode to kubernetes with command: apictl set --mode kubernetes",
				errors.New("mode should be set to kubernetes"))
		}

	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(updateApiCmd)
	updateApiCmd.Flags().StringVarP(&updateflagApiName, "name", "n", "", "Name of the API")
	updateApiCmd.Flags().StringVarP(&updateflagSwaggerFilePath, "from-file", "f", "", "Path to swagger file")
	updateApiCmd.Flags().IntVar(&updateflagReplicas, "replicas", 1, "replica set")
	updateApiCmd.Flags().StringVar(&updateflagNamespace, "namespace", "", "namespace of API")
}
