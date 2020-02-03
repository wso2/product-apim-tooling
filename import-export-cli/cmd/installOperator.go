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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const installOperatorCmdLiteral = "operator"

var flagApiOperatorFile string

//var flagRegistryHost string
//var flagUsername string
//var flagPassword string
//var flagBatchMod bool

// installOperatorCmd represents the installOperator command
var installOperatorCmd = &cobra.Command{
	Use:   installOperatorCmdLiteral,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + installOperatorCmdLiteral + " called")

		// OLM installation requires time to install before installing the WSO2 API Operator
		// Hence getting user inputs
		registryUrl, repository, username, password := readInputs()
		isLocalInstallation := flagApiOperatorFile != ""
		if !isLocalInstallation {
			installOLM("0.13.0")
			installApiOperatorOperatorHub()
		}

		createDockerSecret(registryUrl, username, password)
		createControllerConfigs(repository, isLocalInstallation)
	},
}

// installOLM installs Operator Lifecycle Manager (OLM) with the given version
func installOLM(version string) {
	utils.Logln(utils.LogPrefixInfo + "Installing OLM")

	// this implements the logic in
	// https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.13.0/install.sh
	olmNamespace := "olm"
	csvPhaseSuccessed := "Succeeded"

	// apply OperatorHub CRDs and OLM
	err := utils.K8sApplyFromFile(fmt.Sprintf(utils.OlmCrdUrlTemplate, version), fmt.Sprintf(utils.OlmOlmUrlTemplate, version))
	if err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}

	// rolling out
	if err := utils.ExecuteCommand(utils.Kubectl, utils.K8sRollOut, "status", "-w", "deployment/olm-operator", "-n", olmNamespace); err != nil {
		utils.HandleErrorAndExit("Error installing OLM: Rolling out deployment OLM Operator", err)
	}
	if err := utils.ExecuteCommand(utils.Kubectl, utils.K8sRollOut, "status", "-w", "deployment/catalog-operator", "-n", olmNamespace); err != nil {
		utils.HandleErrorAndExit("Error installing OLM: Rolling out deployment Catalog Operator", err)
	}

	// wait max 50s to csv phase to be succeeded
	csvPhase := ""
	for i := 50; i > 0 && csvPhase != csvPhaseSuccessed; i-- {
		newCsvPhase, err := utils.GetCommandOutput(utils.Kubectl, utils.K8sGet, utils.K8sCsv, "-n", olmNamespace, "packageserver", "-o", `jsonpath={.status.phase}`)
		if err != nil {
			utils.HandleErrorAndExit("Error installing OLM: Getting csv phase", err)
		}

		// only print new phase
		if csvPhase != newCsvPhase {
			fmt.Println("Package server phase: " + newCsvPhase)
			csvPhase = newCsvPhase
		}

		// sleep 1 second
		time.Sleep(1e9)
	}

	if csvPhase != csvPhaseSuccessed {
		utils.HandleErrorAndExit("Error installing OLM: CSV Package Server failed to reach phase succeeded", nil)
	}
}

// installApiOperatorOperatorHub installs WSO2 api-operator from Operator-Hub
func installApiOperatorOperatorHub() {
	utils.Logln(utils.LogPrefixInfo + "Installing API Operator from Operator-Hub")
	operatorFile := utils.OperatorYamlUrl

	err := utils.K8sApplyFromFile(operatorFile)
	if err != nil {
		utils.HandleErrorAndExit("Error installing API Operator from Operator-Hub", err)
	}
}

// createDockerSecret creates K8S secret with credentials for docker registry
func createDockerSecret(registryUrl string, username string, password string) {
	encodedCredential := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	auth := map[string]map[string]map[string]string{
		"auths": {
			registryUrl: {
				"auth":     encodedCredential,
				"username": username,
				"password": password,
			},
		},
	}
	authJsonByte, err := json.Marshal(auth)
	if err != nil {
		utils.HandleErrorAndExit("Error marshalling docker secret credentials ", err)
	}

	encodedAuthJson := base64.StdEncoding.EncodeToString(authJsonByte)
	secretTemplate, _ := box.Get("/kubernetes_resources/registry_secret_mustache.yaml")
	secretYaml, err := mustache.Render(string(secretTemplate), map[string]string{
		"encodedJson": encodedAuthJson,
	})
	if err != nil {
		utils.HandleErrorAndExit("Error rendering docker secret credentials", err)
	}

	// apply created secret yaml file
	if err := utils.K8sApplyFromStdin(secretYaml); err != nil {
		utils.HandleErrorAndExit("Error creating docker secret credentials", err)
	}
}

// createControllerConfigs downloads the mustache, replaces repository value and creates the config: `controller-config`
func createControllerConfigs(repository string, isLocalInstallation bool) {
	utils.Logln(utils.LogPrefixInfo + "Installing controller configs")
	configFile := flagApiOperatorFile

	if !isLocalInstallation {
		// TODO: renuka replace this url (configuration?)
		configFile = `https://gist.githubusercontent.com/renuka-fernando/6d6c64c786e6d13742e802534de3da4e/raw/3a654b9f54d1115532ca757c108b98bff09bcd74/controller_conf.yaml`
	}

	if err := utils.K8sApplyFromFile(configFile); err != nil {
		utils.HandleErrorAndExit("Error creating configurations", err)
	}

	controllerConfigMapYaml, err := utils.GetCommandOutput(utils.Kubectl, utils.K8sGet, "cm", "controller-config", "-n", "wso2-system", "-o", "yaml")
	if err != nil {
		utils.HandleErrorAndExit("Error reading controller-config", err)
	}

	controllerConfigMap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(controllerConfigMapYaml), &controllerConfigMap); err != nil {
		utils.HandleErrorAndExit("Error reading controller-config", err)
	}
	controllerConfigMap["data"].(map[interface{}]interface{})["dockerRegistry"] = repository

	configuredConfigMap, err := yaml.Marshal(controllerConfigMap)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering controller-config", err)
	}

	if err := utils.K8sApplyFromStdin(string(configuredConfigMap)); err != nil {
		utils.HandleErrorAndExit("Error creating controller-configs", err)
	}
}

// readInputs reads docker-registry URL, repository, username and password from the user
func readInputs() (string, string, string, string) {
	isConfirm := false
	registryUrl := ""
	repository := ""
	username := ""
	password := ""
	var err error

	for !isConfirm {
		registryUrl, err = utils.ReadInputString("Enter Docker-Registry URL", utils.DockerRegistryUrl, "", true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading Docker-Registry URL", err)
		}

		repository, err = utils.ReadInputString("Enter Repository Name", "docker.io/library", utils.UsernameValidationRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading Repository Name", err)
		}

		username, err = utils.ReadInputString("Enter Username", "", utils.UsernameValidationRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading Username", err)
		}

		password, err = utils.ReadPassword("Enter Password")
		if err != nil {
			utils.HandleErrorAndExit("Error reading Password", err)
		}

		isCredentialsValid, err := validateDockerRegistry(registryUrl, repository, username, password)
		if !isCredentialsValid {
			utils.HandleErrorAndExit("Error connecting to Docker Registry", err)
		}

		fmt.Println("")
		fmt.Println("Docker-Registry URL: " + registryUrl)
		fmt.Println("Repository         : " + repository)
		fmt.Println("Username           : " + username)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", "Y", "", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirmStr = strings.ToUpper(isConfirmStr)
		isConfirm = isConfirmStr == "Y" || isConfirmStr == "YES"
	}

	return registryUrl, repository, username, password
}

func validateDockerRegistry(registryUrl string, repository string, username string, password string) (bool, error) {
	// remove version tag if exists
	regExpString := `\/v[\d.-]*\/?$`
	regExp := regexp.MustCompile(regExpString)
	registryUrl = regExp.ReplaceAllString(registryUrl, "")

	_, err := registry.New(registryUrl+"/"+repository, username, password)
	if err != nil {
		return false, err
	}
	return true, nil
}

// init using Cobra
func init() {
	installCmd.AddCommand(installOperatorCmd)
	installOperatorCmd.Flags().StringVarP(&flagApiOperatorFile, "from-file", "f", "", "Path to API Operator directory")
	//installOperatorCmd.Flags().StringVarP(&flagRegistryHost, "registry-host", "h", "", "URL of the registry host")
	//installOperatorCmd.Flags().StringVarP(&flagUsername, "username", "u", "", "Username for the registry repository")
	//installOperatorCmd.Flags().StringVarP(&flagPassword, "password", "p", "", "Password for the registry repository user")
	//installOperatorCmd.Flags().BoolVarP(&flagBatchMod, "batch-mod", "B", false, "Run in non-interactive (batch) mode")
}
