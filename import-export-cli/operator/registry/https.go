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
	"fmt"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"strings"
)

// HttpsRegistry represents private HTTPS registry
var HttpsRegistry = &Registry{
	Name:       "HTTPS",
	Caption:    "HTTPS Private Registry",
	Repository: Repository{},
	Option:     5,
	Read: func(reg *Registry, flagValues *map[string]FlagValue) {
		var repository, username, password string

		// check input mode: interactive or batch
		if flagValues == nil {
			// get inputs in interactive mode
			repository, username, password = readHttpsRepInputs()
		} else {
			// get inputs in batch mode
			repository = (*flagValues)[k8sUtils.FlagBmRepository].Value.(string)
			username = (*flagValues)[k8sUtils.FlagBmUsername].Value.(string)
			password = (*flagValues)[k8sUtils.FlagBmPassword].Value.(string)

			// validate required inputs
			if !utils.ValidateValue(repository, utils.RepoValidRegex) {
				utils.HandleErrorAndExit("Invalid repository name: "+repository, nil)
			}

			// validate optional inputs
			if username != "" && !utils.ValidateValue(username, utils.UsernameValidRegex) {
				utils.HandleErrorAndExit("Invalid username : "+username, nil)
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

		reg.Repository.Name = repository
		reg.Repository.Username = username
		reg.Repository.Password = password
	},
	Run: func(reg *Registry) {
		if reg.Repository.ServerUrl == "" {
			reg.Repository.ServerUrl = getRegistryUrl(reg.Repository.Name)
		}

		k8sUtils.K8sCreateSecretFromInputs(
			k8sUtils.DockerRegCredSecret, k8sUtils.ApiOpWso2Namespace,
			reg.Repository.ServerUrl, reg.Repository.Username, reg.Repository.Password,
		)
		reg.Repository.Password = "" // clear password
	},
	Flags: Flags{
		RequiredFlags: &map[string]bool{k8sUtils.FlagBmRepository: true},
		OptionalFlags: &map[string]bool{k8sUtils.FlagBmUsername: true, k8sUtils.FlagBmPassword: true,
			k8sUtils.FlagBmPasswordStdin: true},
	},
}

// readHttpsRepInputs reads https private registry URL, username and password from the user
func readHttpsRepInputs() (string, string, string) {
	isConfirm := false
	repository := ""
	username := ""
	password := ""
	var err error

	for !isConfirm {
		repository, err = utils.ReadInputString("Enter repository",
			utils.Default{Value: "", IsDefault: false}, utils.RepoValidRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading registry repository name from user", err)
		}

		username, err = utils.ReadInputString("Enter username", utils.Default{Value: "", IsDefault: false},
			utils.UsernameValidRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading username from user", err)
		}

		password, err = utils.ReadPassword("Enter password")
		if err != nil {
			utils.HandleErrorAndExit("Error reading password from user", err)
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

// getRegistryUrl returns the registry URL for given repository
func getRegistryUrl(repository string) string {
	names := strings.SplitN(repository, "/", 2)
	// if "myDomain.com:5000/foo" return "myDomain.com:5000"
	if len(names) == 2 {
		return names[0]
	}

	// if "myDomain.com:5000" return "myDomain.com:5000"
	return repository
}

func init() {
	add(HttpsRegistry)
}
