/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	wso2v1alpha1 "github.com/wso2/k8s-apim-operator/apim-operator/pkg/apis/wso2/v1alpha1"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var flagApiName string
var flagSwaggerFilePath string
var flagReplicas int
var flagNamespace string

const apiCmdLiteral = "api"
const apiCmdShortDesc = "handle APIs in kubernetes cluster "
const apiLongDesc = `Add, Update and Delete APIs in kubernetes cluster. JSON and YAML formats are accepted.
available modes are as follows
* kubernetes`
const apiExamples = utils.ProjectName + " add/update " + apiCmdLiteral + ` -n petstore --from-file=./Swagger.json --replicas=3 --namespace=wso2`

const addCmdLiteral = "add"
const addCmdShortDesc = "Add an API to the kubernetes cluster"
const addCmdLongDesc = `Add an API from a Swagger file to the kubernetes cluster. JSON and YAML formats are accepted.
available modes are as follows
* kubernetes`
const addCmdExamples = utils.ProjectName + " " + apiCmdLiteral + " " + `-n petstore --from-file=./Swagger.json --replicas=3 --namespace=wso2`

// addApiCmd represents the addApi command
var addCmd = &cobra.Command{
	Use:     addCmdLiteral,
	Short:   addCmdShortDesc,
	Long:    addCmdLongDesc,
	Example: addCmdExamples,
}

var addApiCmd = &cobra.Command{
	Use:     apiCmdLiteral,
	Short:   apiCmdShortDesc,
	Long:    apiLongDesc,
	Example: apiExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apiCmdLiteral + " called")
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if configVars.Config.KubernetesMode {
			if flagApiName == "" && flagSwaggerFilePath == "" {
				utils.HandleErrorAndExit("Required flags are missing. API name and swagger file paths are required",
					errors.New("required flags missing"))
			} else {
				err := createConfigMapWithNamespace()
				if err != nil {
					utils.HandleErrorAndExit("Error creating configmap", err)
				}
				createAPI()
			}
		} else {
			utils.HandleErrorAndExit("set mode to kubernetes with command - apimcli set-mode kubernetes ",
				errors.New("mode should be set to kubernetes"))
		}

	},
}

func createConfigMapWithNamespace() error {

	configMapName := flagApiName + "-swagger"
	fmt.Println("creating configmap with swagger definition")
	cmd := exec.Command(
		utils.Kubectl,
		"create",
		"configmap",
		configMapName,
		"--from-file",
		flagSwaggerFilePath,
		"-n", flagNamespace,
	)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func createAPI() {

	apiConfigMapData, _ := box.Get("/kubernetes_resources/api_cr.yaml")

	apiConfigMap := &wso2v1alpha1.API{}
	errUnmarshal := yaml.Unmarshal(apiConfigMapData, apiConfigMap)
	if errUnmarshal != nil {
		utils.HandleErrorAndExit("Error unmarshal api configmap into struct ", errUnmarshal)
	}

	//assigning values to API cr
	apiConfigMap.Name = flagApiName
	apiConfigMap.Namespace = flagNamespace
	apiConfigMap.Spec.Definition.ConfigmapName = flagApiName + "-swagger"
	apiConfigMap.Spec.Replicas = flagReplicas

	fmt.Println(apiConfigMap)

	byteVal, errMarshal := yaml.Marshal(apiConfigMap)
	if errMarshal != nil {
		utils.HandleErrorAndExit("Error marshal api configmap ", errMarshal)
	}

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
	cmd := exec.Command(
		utils.Kubectl,
		"apply",
		"-f",
		tmpFile.Name(),
		"-n", flagNamespace,
	)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))

}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addApiCmd)

	addApiCmd.Flags().StringVarP(&flagApiName, "name", "n", "", "Name of the API")
	addApiCmd.Flags().StringVar(&flagSwaggerFilePath, "from-file", "", "Path to swagger file")
	addApiCmd.Flags().IntVar(&flagReplicas, "replicas", 1, "replica set")
	addApiCmd.Flags().StringVar(&flagNamespace, "namespace", "", "namespace of API")

}
