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
	wso2v1alpha1 "github.com/wso2/k8s-apim-operator/apim-operator/pkg/apis/wso2/v1alpha1"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var flagApiName string
var flagSwaggerFilePaths []string
var flagReplicas int
var flagNamespace string
var flagOverride bool
var flagApiVersion string
var flagApiMode string

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
		validateAddApiCommand()

		swaggerCmNames := make([]string, len(flagSwaggerFilePaths))
		var balInterceptorsCmNames []string
		var javaInterceptorsCmNames []string

		for i, flagSwaggerFilePath := range flagSwaggerFilePaths {
			fmt.Println(fmt.Sprintf("Processing swagger %v: %v", i+1, flagSwaggerFilePath))
			swaggerCmNames[i] = fmt.Sprintf("%v-%v-swagger", flagApiName, i+1)

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
				balInterceptorsCmName := fmt.Sprintf("%v-%v-bal-interceptors", flagApiName, i+1)
				balInterceptorsCmNames = append(balInterceptorsCmNames, balInterceptorsCmName)
				handleBalInterceptors(balInterceptorsCmName, flagSwaggerFilePath, "create", flagNamespace)

				// handle java interceptors
				tempJavaIntCms := handleJavaInterceptors(flagSwaggerFilePath, "create", flagNamespace, fmt.Sprintf("%v-%v", flagApiName, i+1))
				javaInterceptorsCmNames = append(javaInterceptorsCmNames, tempJavaIntCms...)
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
		createAPI(flagApiName, flagNamespace, swaggerCmNames, flagReplicas, "", balInterceptorsCmNames, flagOverride, javaInterceptorsCmNames, flagApiMode, flagApiVersion)
	},
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
		utils.HandleErrorAndExit("required flags are missing. API name and swagger file paths are required",
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

func createAPI(name string, namespace string, configMapNames []string, replicas int, timestamp string, balInterceptors []string, override bool, javaInterceptors []string, apiMode string, apiVersion string) {
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
	if timestamp != "" {
		//set update timestamp
		apiConfigMap.Spec.UpdateTimeStamp = timestamp
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
	//execute kubernetes command to create or update api from file
	cmd := exec.Command(
		k8sUtils.Kubectl,
		"apply",
		"-f",
		tmpFile.Name(),
		"-n", namespace,
	)
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	errAddApi := cmd.Run()
	if errAddApi != nil {
		fmt.Println(errAddApi)
	}
}

func handleBalInterceptors(configMapName string, path string, operation string, namespace string) {
	//get interceptors if available
	interceptorsPath := filepath.Join(path, "Interceptors")
	//check interceptors dir is not empty
	file, err := os.Open(interceptorsPath)
	if err != nil {
		utils.HandleErrorAndExit("cannot open interceptors Dir", err)
	}
	defer file.Close()
	if _, err = file.Readdir(1); err != nil {
		return
	}

	//creating kubernetes configmap with interceptors
	fmt.Println("creating configmap with ballerina interceptors")
	if err := createConfigMapWithNamespace(configMapName, interceptorsPath, namespace, operation); err != nil {
		utils.HandleErrorAndExit("Error creating configmap for interceptors", err)
	}
}

func handleJavaInterceptors(path string, operation string, namespace string, cmPrefixName string) []string {
	var interceptors []string
	var javaInterceptorsConfNames []string
	//get interceptors if available
	interceptorsPath := filepath.Join(path, "libs")
	//check interceptors dir is not empty
	if exists, err := utils.IsDirExists(interceptorsPath); !exists {
		utils.HandleErrorAndExit("cannot open java interceptors Dir", err)
	}

	//get all jars in libs dir
	errReadInterceptors := filepath.Walk(interceptorsPath, func(path string, info os.FileInfo, err error) error {
		interceptors = append(interceptors, path)
		return nil
	})
	if errReadInterceptors != nil {
		utils.HandleErrorAndExit("cannot read interceptors in the libs", errReadInterceptors)
	}
	for _, filePath := range interceptors {
		if filepath.Ext(filePath) == ".jar" {
			//creating kubernetes configmap for each java interceptor
			cmName := cmPrefixName + "-" + filepath.Base(filePath)
			javaInterceptorsConfNames = append(javaInterceptorsConfNames, cmName)
			fmt.Println("creating configmap with java interceptor " + cmName)
			errConfInt := createConfigMapWithNamespace(cmPrefixName+"-"+filepath.Base(filePath), filePath, namespace, operation)
			if errConfInt != nil {
				utils.HandleErrorAndExit("Error creating configmap for java-interceptor "+cmName, errConfInt)
			}
		}
	}
	return javaInterceptorsConfNames
}

func init() {
	addCmd.AddCommand(addApiCmd)
	addApiCmd.Flags().StringVarP(&flagApiName, "name", "n", "", "Name of the API")
	addApiCmd.Flags().StringArrayVarP(&flagSwaggerFilePaths, "from-file", "f", []string{}, "Path to swagger file")
	addApiCmd.Flags().IntVar(&flagReplicas, "replicas", 1, "replica set")
	addApiCmd.Flags().StringVar(&flagNamespace, "namespace", "", "namespace of API")
	addApiCmd.Flags().BoolVarP(&flagOverride, "override", "", false, "Property to override the existing docker image with same name and version")
	addApiCmd.Flags().StringVarP(&flagApiVersion, "version", "v", "", "Property to override the existing docker image with same name and version")
	addApiCmd.Flags().StringVarP(&flagApiMode, "mode", "m", "",
		fmt.Sprintf("Property to override the deploying mode. Available modes: %v, %v", utils.PrivateJetModeConst, utils.SidecarModeConst))
}
