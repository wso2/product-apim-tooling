/*
Copyright © 2020 Renuka Fernando <renuka@wso2.com>

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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// changeDockerRegistryCmd represents the change registry command
var changeDockerRegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if !configVars.Config.KubernetesMode {
			utils.HandleErrorAndExit("set mode to kubernetes with command: apictl set --mode kubernetes",
				errors.New("mode should be set to kubernetes"))
		}

		registryUrl, username, password := readDockerRegistryInputs()

		// set registry first since this can throw error if api operator not installed. If error occur no need to rollback secret.
		setRegistryRepositoryOnControllerConfig(registryUrl)
		createDockerSecret(registryUrl, username, password)
	},
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

// setRegistryRepositoryOnControllerConfig sets the repository value in the config: `controller-config`
func setRegistryRepositoryOnControllerConfig(dockerRegistryUrl string) {
	// get controller config config map
	controllerConfigMapYaml, err := utils.GetCommandOutput(
		utils.Kubectl, utils.K8sGet, "cm", utils.ApiOpControllerConfigMap,
		"-n", utils.ApiOpWso2Namespace,
		"-o", "yaml",
	)
	if err != nil {
		utils.HandleErrorAndExit("Error reading controller-config.\nInstall api operator using the command: apictl install api-operator",
			errors.New("error reading controller-config"))
	}

	// replace registry
	controllerConfigMap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(controllerConfigMapYaml), &controllerConfigMap); err != nil {
		utils.HandleErrorAndExit("Error reading controller-config", err)
	}

	// remove protocol from docker-registry url
	regExpString := `^(http|https)://`
	regExp := regexp.MustCompile(regExpString)
	repository := regExp.ReplaceAllString(dockerRegistryUrl, "")

	controllerConfigMap["data"].(map[interface{}]interface{})["dockerRegistry"] = repository
	configuredConfigMap, err := yaml.Marshal(controllerConfigMap)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering controller-config", err)
	}

	// apply controller config config map back
	if err := utils.K8sApplyFromStdin(string(configuredConfigMap)); err != nil {
		utils.HandleErrorAndExit("Error creating controller-configs", err)
	}
}

// readDockerRegistryInputs reads docker-registry URL, username and password from the user
func readDockerRegistryInputs() (string, string, string) {
	isConfirm := false
	registryUrl := ""
	username := ""
	password := ""
	var err error

	for !isConfirm {
		registryUrl, err = utils.ReadInputString("Enter Docker-Registry URL", utils.DockerRegistryUrl, "", true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading Docker-Registry URL", err)
		}

		username, err = utils.ReadInputString("Enter Username", "", utils.UsernameValidationRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading Username", err)
		}

		password, err = utils.ReadPassword("Enter Password")
		if err != nil {
			utils.HandleErrorAndExit("Error reading Password", err)
		}

		isCredentialsValid, err := validateDockerRegistry(registryUrl, username, password)
		if !isCredentialsValid {
			utils.HandleErrorAndExit("Error connecting to Docker Registry repository using credentials", err)
		}

		fmt.Println("")
		fmt.Println("Docker-Registry URL: " + registryUrl)
		fmt.Println("Username           : " + username)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", "Y", "", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirmStr = strings.ToUpper(isConfirmStr)
		isConfirm = isConfirmStr == "Y" || isConfirmStr == "YES"
	}

	return registryUrl, username, password
}

// validateDockerRegistry validates the credentials against registry url and repository
func validateDockerRegistry(registryUrl string, username string, password string) (bool, error) {
	// remove version tag if exists
	//regExpString := `\/v[\d.-]*\/?$`
	//regExp := regexp.MustCompile(regExpString)
	//registryUrl = regExp.ReplaceAllString(registryUrl, "")
	//
	//hub, err := registry.New(registryUrl, username, password)
	//if err != nil {
	//	return false, err
	//}
	//
	//if _, err := hub.Repositories(); err != nil {
	//	return false, err
	//}

	return true, nil
}

func init() {
	RootCmd.AddCommand(changeCmd)
	changeCmd.AddCommand(changeDockerRegistryCmd)
}
