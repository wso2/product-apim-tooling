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

package registry

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"strings"
)

const DockerRegistryUrl = "https://index.docker.io/v1/"

var dockerRepo = new(string)

var dockerHubValues = struct {
	repository string
	username   string
	password   string // TODO: renuka: password should be byte[], strings can be exploited from memory
}{}

var DockerHubRegistry = &Registry{
	Name:       "DOCKER_HUB",
	Caption:    "Docker Hub",
	Repository: dockerRepo,
	Option:     1,
	Read: func() {
		repository, username, password := readDockerHubInputs()
		*dockerRepo = repository
		dockerHubValues.repository = repository
		dockerHubValues.username = username
		dockerHubValues.password = password
	},
	Run: func() {
		createDockerSecret(dockerHubValues.username, dockerHubValues.password)
		dockerHubValues.password = "" // clear password
	},
}

// readDockerHubInputs reads docker-registry URL, username and password from the user
func readDockerHubInputs() (string, string, string) {
	isConfirm := false
	repository := ""
	username := ""
	password := ""
	var err error

	for !isConfirm {
		repository, err = utils.ReadInputString("Enter repository name", "", utils.UsernameValidRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading DockerHub repository name from user", err)
		}

		username, err = utils.ReadInputString("Enter username", "", utils.UsernameValidRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading username from user", err)
		}

		password, err = utils.ReadPassword("Enter password")
		if err != nil {
			utils.HandleErrorAndExit("Error reading password from user", err)
		}

		isCredentialsValid, err := validateDockerHubCredentials(repository, username, password)
		if err != nil {
			utils.HandleErrorAndExit("Error connecting to Docker Registry repository using credentials", err)
		}

		if !isCredentialsValid {
			utils.HandleErrorAndExit("Invalid credentials", err)
		}

		fmt.Println("")
		fmt.Println("Repository: " + repository)
		fmt.Println("Username  : " + username)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", "Y", "", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirmStr = strings.ToUpper(isConfirmStr)
		isConfirm = isConfirmStr == "Y" || isConfirmStr == "YES"
	}

	return repository, username, password
}

// validateDockerHubCredentials validates the credentials for the repository
func validateDockerHubCredentials(repository string, username string, password string) (bool, error) {
	cred, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	resp, err := http.Post("https://hub.docker.com/v2/users/login/", "application/json", bytes.NewBuffer(cred))
	if err != nil {
		return false, err
	}
	_ = resp.Body.Close()
	return resp.StatusCode == 200, nil //TODO: renuka: use reposity as well to validate
}

// createDockerSecret creates K8S secret with credentials for docker registry
func createDockerSecret(username string, password string) {
	encodedCredential := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	auth := map[string]map[string]map[string]string{
		"auths": {
			DockerRegistryUrl: {
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
	if err := k8sUtils.K8sApplyFromStdin(secretYaml); err != nil {
		utils.HandleErrorAndExit("Error creating docker secret credentials", err)
	}
}

func init() {
	add(DockerHubRegistry)
}
