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

// DockerHubRegistry represents Docker Hub registry
var DockerHubRegistry = &Registry{
	Name:       "DOCKER_HUB",
	Caption:    "Docker Hub",
	Repository: Repository{ServerUrl: "https://index.docker.io/v1/"},
	Option:     1,
	Read: func(reg *Registry, flagValues *map[string]FlagValue) {
		var repository, username, password string

		// check input mode: interactive or batch
		if flagValues == nil {
			// get inputs in interactive mode
			repository, username, password = readDockerHubInputs()
		} else {
			// get inputs in batch mode
			repository = (*flagValues)[k8sUtils.FlagBmRepository].Value.(string)
			username = (*flagValues)[k8sUtils.FlagBmUsername].Value.(string)
			password = (*flagValues)[k8sUtils.FlagBmPassword].Value.(string)

			// validate required inputs
			if !utils.ValidateValue(repository, utils.RepoValidRegex) {
				utils.HandleErrorAndExit("Invalid repository name: "+repository, nil)
			}

			// if "--password-stdin" is supplied get password from stdin
			if (*flagValues)[k8sUtils.FlagBmPasswordStdin].Value.(bool) {
				pwStdin, err := utils.ReadPassword("Enter password")
				if err != nil {
					utils.HandleErrorAndExit("Error reading password from user", err)
				}
				password = pwStdin
			}
		}

		// support prefix "docker.io/" to be compatible with older version of API-CTL
		// trim it if found
		reg.Repository.Name = strings.TrimPrefix(repository, "docker.io/")
		reg.Repository.Username = username
		reg.Repository.Password = password
	},
	Run: func(reg *Registry) {
		k8sUtils.K8sCreateSecretFromInputs(k8sUtils.DockerRegCredSecret, k8sUtils.ApiOpWso2Namespace,
			reg.Repository.ServerUrl, reg.Repository.Username, reg.Repository.Password)
		reg.Repository.Password = "" // clear password
	},
	Flags: Flags{
		RequiredFlags: &map[string]bool{k8sUtils.FlagBmRepository: true, k8sUtils.FlagBmUsername: true},
		OptionalFlags: &map[string]bool{k8sUtils.FlagBmPassword: true, k8sUtils.FlagBmPasswordStdin: true},
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
		repository, err = utils.ReadInputString("Enter repository name", utils.Default{Value: "", IsDefault: false},
			utils.RepoValidRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading repository name from user", err)
		}

		username, err = utils.ReadInputString("Enter username", utils.Default{Value: "", IsDefault: false},
			"", true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading username from user", err)
		}

		password, err = utils.ReadPassword("Enter password")
		if err != nil {
			utils.HandleErrorAndExit("Error reading password from user", err)
		}

		// only validate credentials if registry is DockerHub
		isCredentialsValid, err := validateDockerHubCredentials(repository, username, password)
		if err != nil {
			utils.HandleErrorAndExit("Error connecting to Docker Registry repository using credentials", err)
		}

		if !isCredentialsValid {
			utils.HandleErrorAndExit("Invalid credentials for Docker Hub", err)
		}

		fmt.Println("\nRepository: " + repository)
		fmt.Println("Username  : " + username)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations",
			utils.Default{Value: "Y", IsDefault: true}, "", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirm = strings.EqualFold(isConfirmStr, "y") || strings.EqualFold(isConfirmStr, "yes")
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
	return resp.StatusCode == 200, nil //TODO: use repository as well to validate
}

func init() {
	add(DockerHubRegistry)
}
