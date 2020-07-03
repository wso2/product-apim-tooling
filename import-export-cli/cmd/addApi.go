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
	"bytes"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	wso2v1alpha1 "github.com/wso2/k8s-api-operator/api-operator/pkg/apis/wso2/v1alpha1"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var flagApiName string
var flagSwaggerFilePaths []string
var flagReplicas int
var flagNamespace string
var flagOverride bool
var flagApiVersion string
var flagApiMode string
var flagApiEndPoint string

const addApiCmdLiteral = "api"
const addApiCmdShortDesc = "handle APIs in kubernetes cluster "
const addApiLongDesc = `Add, Update and Delete APIs in kubernetes cluster. JSON and YAML formats are accepted.
available modes are as follows
* kubernetes`
const addApiExamples = utils.ProjectName + " add/update " + addApiCmdLiteral + ` -n petstore --from-file=./Swagger.json --replicas=3 --namespace=wso2`

// addApiCmd represents the api command
var addApiCmd = &cobra.Command{
	Use:     addApiCmdLiteral,
	Short:   addApiCmdShortDesc,
	Long:    addApiLongDesc,
	Example: addApiExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + addApiCmdLiteral + " called")
		handleAddApi("")
	},
}

func handleAddApi(nameSuffix string) {
	validateAddApiCommand()

	swaggerCmNames := make([]string, len(flagSwaggerFilePaths))
	balInterceptorsCmNames := make([]string, 0, len(flagSwaggerFilePaths))
	var javaInterceptorsCmNames []string

	for i, flagSwaggerFilePath := range flagSwaggerFilePaths {
		// log processing only if there are more projects
		utils.Logln(fmt.Sprintf("%sProcessing swagger %v: %v", utils.LogPrefixInfo, i+1, flagSwaggerFilePath))

		flagApiName = strings.ToLower(flagApiName)
		swaggerCmNames[i] = fmt.Sprintf("%v-%v-swagger%s", flagApiName, i+1, nameSuffix)

		fi, _ := os.Stat(flagSwaggerFilePath) // error already handled and ignore error
		switch mode := fi.Mode(); {
		//check if the swagger path is a Dir
		case mode.IsDir():
			//get swagger definition
			swaggerPath := filepath.Join(flagSwaggerFilePath, filepath.FromSlash("Meta-information/swagger.yaml"))
			//creating kubernetes configmap with swagger definition
			fmt.Println("creating configmap with swagger definition")
			errConf := createConfigMapWithNamespace(swaggerCmNames[i], swaggerPath, flagNamespace, k8sUtils.K8sCreate)
			if errConf != nil {
				utils.HandleErrorAndExit("Error creating configmap", errConf)
			}

			// copy all bal interceptors to the temp dir
			balInterceptorsCmName := fmt.Sprintf("%v-%v-bal-intcpt%s", flagApiName, i+1, nameSuffix)
			intceptFound := handleBalInterceptors(balInterceptorsCmName, flagSwaggerFilePath, "create", flagNamespace)
			if intceptFound {
				balInterceptorsCmNames = append(balInterceptorsCmNames, balInterceptorsCmName)
			}

			// handle java interceptors
			tempJavaIntCms := handleJavaInterceptors(nameSuffix, flagSwaggerFilePath, "create", flagNamespace, fmt.Sprintf("%v-%v", flagApiName, i+1))
			if tempJavaIntCms != nil {
				javaInterceptorsCmNames = append(javaInterceptorsCmNames, tempJavaIntCms...)
			}
		//check if the swagger path is a file
		case mode.IsRegular():
			//creating kubernetes configmap with swagger definition
			fmt.Println("creating configmap with swagger definition")
			errConf := createConfigMapWithNamespace(swaggerCmNames[i], flagSwaggerFilePath, flagNamespace, k8sUtils.K8sCreate)
			if errConf != nil {
				utils.HandleErrorAndExit("Error creating configmap", errConf)
			}
		}
	}

	//create API
	fmt.Println("creating API definition")
	createAPI(flagApiName, flagNamespace, swaggerCmNames, flagReplicas, nameSuffix, balInterceptorsCmNames, flagOverride, javaInterceptorsCmNames, flagApiMode, flagApiVersion, flagApiEndPoint)
}

// validateAddApiCommand validates for required flags and if invalid print error and exit
// mode should be k8s
func validateAddApiCommand() {
	// validate mode
	configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	if !configVars.Config.KubernetesMode {
		utils.HandleErrorAndExit("set mode to kubernetes with command: apictl set --mode kubernetes",
			errors.New("mode should be set to kubernetes"))
	}

	// validate required flags
	if flagApiName == "" || len(flagSwaggerFilePaths) == 0 {
		utils.HandleErrorAndExit("required flags are missing. API name and swagger file paths are requiredn"+
			"required flags: --name <name>, --from-file <swagger-file>",
			errors.New("required flags missing"))
	}

	// validate --from-file flag values
	for _, swaggerFilePath := range flagSwaggerFilePaths {
		if _, err := os.Stat(swaggerFilePath); err != nil {
			utils.HandleErrorAndExit("swagger file path or project not found", err)
		}
	}

	// validate --mode flag
	if flagApiMode != "" && flagApiMode != utils.PrivateJetModeConst && flagApiMode != utils.SidecarModeConst {
		utils.HandleErrorAndExit(fmt.Sprintf("invalid api mode. available modes: %v, %v", utils.PrivateJetModeConst, utils.SidecarModeConst), nil)
	}
}

//create configmap with swagger definition
func createConfigMapWithNamespace(configMapName string, filePath string, namespace string, operation string) error {
	cmd := exec.Command(
		k8sUtils.Kubectl,
		operation,
		"configmap",
		configMapName,
		"--from-file",
		filePath,
		"-n", namespace,
	)
	//print kubernetes error commands
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func createAPI(name string, namespace string, configMapNames []string, replicas int, timestamp string, balInterceptors []string, override bool, javaInterceptors []string, apiMode string, apiVersion string,
	apiEndPoint string) {
	//get API definition from file
	apiConfigMapData, _ := box.Get("/kubernetes_resources/api_cr.yaml")
	apiConfigMap := &wso2v1alpha1.API{}
	errUnmarshal := yaml.Unmarshal(apiConfigMapData, apiConfigMap)
	if errUnmarshal != nil {
		utils.HandleErrorAndExit("Error unmarshal api configmap into struct ", errUnmarshal)
	}
	//assigning values to API cr
	apiConfigMap.Name = name
	apiConfigMap.Namespace = namespace
	apiConfigMap.Spec.Definition.SwaggerConfigmapNames = configMapNames
	apiConfigMap.Spec.Replicas = replicas
	apiConfigMap.Spec.Override = override
	apiConfigMap.Spec.ApiEndPoint = apiEndPoint

	k8sOperation := k8sUtils.K8sCreate
	k8sSaveConfig := true
	if timestamp != "" {
		//set update timestamp
		apiConfigMap.Spec.UpdateTimeStamp = timestamp
		k8sOperation = k8sUtils.K8sApply
		k8sSaveConfig = false
	}
	if len(balInterceptors) > 0 {
		// set bal interceptors configmap name in API cr
		apiConfigMap.Spec.Definition.Interceptors.Ballerina = balInterceptors
	}
	if len(javaInterceptors) > 0 {
		//set java interceptors configmaps names in API cr
		apiConfigMap.Spec.Definition.Interceptors.Java = javaInterceptors
	} else {
		apiConfigMap.Spec.Definition.Interceptors.Java = []string{}
	}
	if apiMode != "" {
		apiConfigMap.Spec.Mode = wso2v1alpha1.Mode(apiMode)
	}
	if apiVersion != "" {
		apiConfigMap.Spec.Version = apiVersion
	}
	if apiEndPoint != "" {
		apiConfigMap.Spec.ApiEndPoint = apiEndPoint
	}
	if replicas != 0 {
		apiConfigMap.Status.Replicas = replicas
	}

	byteVal, errMarshal := yaml.Marshal(apiConfigMap)
	if errMarshal != nil {
		utils.HandleErrorAndExit("Error marshal api configmap ", errMarshal)
	}
	//write configmap to a temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "apicr-*.yaml")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err = tmpFile.Write(byteVal); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}
	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	k8sArgs := []string{k8sOperation, "-f", tmpFile.Name(), "-n", namespace}
	if k8sSaveConfig {
		k8sArgs = append(k8sArgs, "--save-config")
	}

	//execute kubernetes command to create or update api from file
	errAddApi := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sArgs...)
	if errAddApi != nil {
		fmt.Println("error configuring API")
		// delete all configs if any error
		rollbackConfigs(apiConfigMap)
	}
}

func handleBalInterceptors(configMapName string, path string, operation string, namespace string) bool {
	//get interceptors if available
	interceptorsPath := filepath.Join(path, "Interceptors")
	//check interceptors dir is not empty
	file, err := os.Open(interceptorsPath)
	if err != nil {
		utils.HandleErrorAndContinue("cannot open Interceptors directory", err)
	} else {
		defer file.Close()
		if _, err = file.Readdir(1); err != nil {
			return false
		}

		//creating kubernetes configmap with interceptors
		fmt.Println("creating configmap with ballerina interceptors")
		if err := createConfigMapWithNamespace(configMapName, interceptorsPath, namespace, operation); err != nil {
			utils.HandleErrorAndExit("Error creating configmap for interceptors", err)
		}

		return true
	}

	return false
}

func handleJavaInterceptors(nameSuffix string, path string, operation string, namespace string, cmPrefixName string) []string {
	var interceptors []string
	var javaInterceptorsConfNames []string
	//get interceptors if available
	interceptorsPath := filepath.Join(path, "libs")
	//check interceptors dir is not empty
	exists, err := utils.IsDirExists(interceptorsPath)
	if !exists {
		utils.HandleErrorAndContinue("cannot open java interceptors directory", err)
	} else {
		//get all jars in libs dir
		errReadInterceptors := filepath.Walk(interceptorsPath, func(path string, info os.FileInfo, err error) error {
			interceptors = append(interceptors, path)
			return nil
		})
		if errReadInterceptors != nil {
			utils.HandleErrorAndExit("cannot read interceptors in the libs", errReadInterceptors)
		}
	}

	const jarExt = ".jar"
	for _, filePath := range interceptors {
		if filepath.Ext(filePath) == jarExt {
			// creating kubernetes configmap for each java interceptor
			// added the random number instead of file name to omit lengthy names and k8s resource name constraints
			cmName := fmt.Sprintf("%s-%v-jar-intcpt%s", cmPrefixName, rand.Intn(10000), nameSuffix)
			javaInterceptorsConfNames = append(javaInterceptorsConfNames, cmName)

			fmt.Println("creating configmap with java interceptor " + cmName)
			errConfInt := createConfigMapWithNamespace(cmName, filePath, namespace, operation)
			if errConfInt != nil {
				utils.HandleErrorAndExit("Error creating configmap for java-interceptor "+cmName, errConfInt)
			}
		}
	}
	return javaInterceptorsConfNames
}

// rollbackConfigs deletes configs defined in the API CR given
func rollbackConfigs(apiCr *wso2v1alpha1.API) {
	var rollbackConfMaps []string // configmap names to be deleted

	// swagger configmaps
	rollbackConfMaps = append(rollbackConfMaps, apiCr.Spec.Definition.SwaggerConfigmapNames...)
	// ballerina interceptor configmaps
	rollbackConfMaps = append(rollbackConfMaps, apiCr.Spec.Definition.Interceptors.Ballerina...)
	// java interceptor configmaps
	rollbackConfMaps = append(rollbackConfMaps, apiCr.Spec.Definition.Interceptors.Java...)

	if len(rollbackConfMaps) == 0 {
		return
	}

	// execute kubernetes command to delete
	fmt.Println("Deleting created configs")
	k8sArgs := []string{k8sUtils.K8sDelete, "cm"}
	k8sArgs = append(k8sArgs, rollbackConfMaps...)

	delConfErr := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sArgs...)
	if delConfErr != nil {
		utils.HandleErrorAndExit("error deleting configmaps of the API: "+apiCr.Name, delConfErr)
	}
}

func init() {
	addCmd.AddCommand(addApiCmd)
	addApiCmd.Flags().StringVarP(&flagApiEndPoint, "apiEndPoint", "a", "","")
	addApiCmd.Flags().StringVarP(&flagApiName, "name", "n", "", "Name of the API")
	addApiCmd.Flags().StringArrayVarP(&flagSwaggerFilePaths, "from-file", "f", []string{}, "Path to swagger file")
	addApiCmd.Flags().IntVar(&flagReplicas, "replicas", 1, "replica set")
	addApiCmd.Flags().StringVar(&flagNamespace, "namespace", "", "namespace of API")
	addApiCmd.Flags().BoolVarP(&flagOverride, "override", "", false, "Property to override the existing docker image with same name and version")
	addApiCmd.Flags().StringVarP(&flagApiVersion, "version", "v", "", "Version of the API")
	addApiCmd.Flags().StringVarP(&flagApiMode, "mode", "m", "",
		fmt.Sprintf("Property to override the deploying mode. Available modes: %v, %v", utils.PrivateJetModeConst, utils.SidecarModeConst))
}