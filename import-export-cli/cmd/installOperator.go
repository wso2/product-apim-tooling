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
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const installOperatorCmdLiteral = "operator"

var flagApiOperatorFile string

//var flagRegistryHost string
//var flagUsername string
//var flagPassword string
//var flagBatchMod bool

// These types define authorization credentials for docker-config
// Credential represents a credential for a docker registry
type Credential struct {
	Auth     string `json:"auth"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Auth represents list of docker registries with credentials
type Auth struct {
	Auths map[string]Credential `json:"auths"`
}

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

		isLocalInstallation := flagApiOperatorFile != ""
		if !isLocalInstallation {
			installOLM("0.13.0")
		}

		// OLM installation requires time to install before installing the WSO2 API Operator
		// Hence getting user inputs
		registryUrl, repository, username, password := readInputs()

		createDockerSecret(registryUrl, username, password)
		installApiOperator(isLocalInstallation)
		createControllerConfigs(repository, isLocalInstallation) //TODO: renuka have to configure repository

	},
}

// installOLM installs Operator Lifecycle Manager (OLM) with the given version
func installOLM(version string) {
	utils.Logln(utils.LogPrefixInfo + "Installing OLM")

	// apply OperatorHub CRDs and OLM
	err := utils.K8sApplyFromFile(fmt.Sprintf(utils.OlmCrdUrlTemplate, version), fmt.Sprintf(utils.OlmOlmUrlTemplate, version))
	if err != nil {
		utils.HandleErrorAndExit("Error installing OLM", err)
	}
}

// installApiOperator installs WSO2 api-operator
func installApiOperator(isLocalInstallation bool) {
	utils.Logln(utils.LogPrefixInfo + "Installing API Operator")
	operatorFile := flagApiOperatorFile
	if !isLocalInstallation {
		operatorFile = utils.OperatorYamlUrl
	}

	// apply WSO2 api-operator via OperatorHub
	err := utils.K8sApplyFromFile(operatorFile)
	if err != nil {
		utils.HandleErrorAndExit("Error API Operator through Operator-Hub", err)
	}
}

// createDockerSecret creates K8S secret with credentials for docker registry
func createDockerSecret(registryUrl string, username string, password string) {
	encodedCredential := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	auth := Auth{Auths: map[string]Credential{registryUrl: {
		Auth:     encodedCredential,
		Username: username,
		Password: password,
	}}}

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
	var mustacheTemplate string

	if !isLocalInstallation {
		// read from GitHub
		mustacheGistUrl := `https://gist.githubusercontent.com/renuka-fernando/6d6c64c786e6d13742e802534de3da4e/raw/d6191bc60f3bae659749e9db5f882bef6d1d062a/controller_conf.yaml`

		templateBytes, err := utils.ReadFromUrl(mustacheGistUrl)
		if err != nil {
			utils.HandleErrorAndExit("Error reading controller-configs from server", err)
		}
		mustacheTemplate = string(templateBytes)
	} else {
		// read from local file
		// TODO: renuka read from file
		mustacheTemplate = ""
	}

	k8sConfigs, err := mustache.Render(mustacheTemplate, map[string]string{
		"usernameDockerRegistry": repository,
	})
	if err != nil {
		utils.HandleErrorAndExit("Error rendering controller-configs", err)
	}

	// apply created secret yaml file
	if err := utils.K8sApplyFromStdin(k8sConfigs); err != nil {
		utils.HandleErrorAndExit("Error creating controller-configs", err)
	}
}

// readInputs reads docker-registry URL, repository, username and password from the user
func readInputs() (string, string, string, string) {
	isConfirm := false
	registryUrl := ""
	repository := "renukafernando-test"
	username := ""
	password := ""
	var err error

	for !isConfirm {
		registryUrl, err = utils.ReadInputString("Enter Docker-Registry URL", utils.DockerRegistryUrl, utils.UrlValidationRegex, true)
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

// init using Cobra
func init() {
	installCmd.AddCommand(installOperatorCmd)
	installOperatorCmd.Flags().StringVarP(&flagApiOperatorFile, "from-file", "f", "", "Path to API Operator directory")
	//installOperatorCmd.Flags().StringVarP(&flagRegistryHost, "registry-host", "h", "", "URL of the registry host")
	//installOperatorCmd.Flags().StringVarP(&flagUsername, "username", "u", "", "Username for the registry repository")
	//installOperatorCmd.Flags().StringVarP(&flagPassword, "password", "p", "", "Password for the registry repository user")
	//installOperatorCmd.Flags().BoolVarP(&flagBatchMod, "batch-mod", "B", false, "Run in non-interactive (batch) mode")
}
