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
	"encoding/json"
	"fmt"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"strings"
)

const DockerRegistryUrl = "https://index.docker.io/v1/"

var dockerHubRepo = new(string)

var dockerHubValues = struct {
	repository string
	username   string
	password   string // TODO: renuka: password should be byte[], strings can be exploited from memory
}{}

var DockerHubRegistry = &Registry{
	Name:       "DOCKER_HUB",
	Caption:    "Docker Hub",
	Repository: dockerHubRepo,
	Option:     1,
	Read: func() {
		repository, username, password := readDockerHubInputs()
		*dockerHubRepo = repository
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
		repository, err = utils.ReadInputString("Enter repository name", utils.Default{Value: "", IsDefault: true}, utils.UsernameValidRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading DockerHub repository name from user", err)
		}

		username, err = utils.ReadInputString("Enter username", utils.Default{Value: "", IsDefault: false}, utils.UsernameValidRegex, true)
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

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", utils.Default{Value: "Y", IsDefault: true}, "", false)
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
	return resp.StatusCode == 200, nil //TODO: renuka: use repository as well to validate
}

// createDockerSecret creates K8S secret with credentials for Docker Hub
func createDockerSecret(username string, password string) {
	dockerSecret, err := k8sUtils.GetCommandOutput(
		utils.Kubectl, utils.Create, utils.K8sSecret, utils.K8sSecretDockerRegType, utils.ConfigJsonVolume,
		"--docker-server", DockerRegistryUrl,
		"--docker-username", username,
		"--docker-password", password,
		"--dry-run", "-o", "yaml",
	)

	if err != nil {
		utils.HandleErrorAndExit("Error rendering kubernetes secret for Docker Hub", err)
	}

	// apply created secret yaml file
	if err := k8sUtils.K8sApplyFromStdin(dockerSecret); err != nil {
		utils.HandleErrorAndExit("Error creating docker secret credentials", err)
	}
}

func init() {
	add(DockerHubRegistry)
}
