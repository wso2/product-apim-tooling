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

// validation regex
const httpRepoValidRegex = `^[\w\d\-\.\:]*\/?[\w\d\-]+$`
const httpRepoUsernameRegex = utils.UsernameValidRegex

var httpRepo = new(string)

var httpValues = struct {
	repository string
	username   string
	password   string
}{}

// HttpRegistry represents private HTTP registry
var HttpRegistry = &Registry{
	Name:       "HTTP",
	Caption:    "HTTP Private Registry",
	Repository: httpRepo,
	Option:     4,
	Read: func(flagValues *map[string]FlagValue) {
		var repository, username, password string

		// check input mode: interactive or batch
		if flagValues == nil {
			// get inputs in interactive mode
			repository, username, password = readHttpRepInputs()
		} else {
			// get inputs in batch mode
			repository = (*flagValues)[k8sUtils.FlagBmRepository].Value.(string)
			username = (*flagValues)[k8sUtils.FlagBmUsername].Value.(string)
			password = (*flagValues)[k8sUtils.FlagBmPassword].Value.(string)

			// validate required inputs
			if !utils.ValidateValue(repository, httpRepoValidRegex) {
				utils.HandleErrorAndExit("Invalid repository name: "+repository, nil)
			}

			// validate optional inputs
			if username != "" && !utils.ValidateValue(username, httpRepoUsernameRegex) {
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

		*httpRepo = repository
		httpValues.repository = repository
		httpValues.username = username
		httpValues.password = password
	},
	Run: func() {
		k8sUtils.K8sCreateSecretFromInputs(k8sUtils.ConfigJsonVolume, getRegistryUrl(httpValues.repository), httpValues.username, httpValues.password)
		httpValues.password = "" // clear password
	},
	Flags: Flags{
		RequiredFlags: &map[string]bool{k8sUtils.FlagBmRepository: true},
		OptionalFlags: &map[string]bool{k8sUtils.FlagBmUsername: true, k8sUtils.FlagBmPassword: true, k8sUtils.FlagBmPasswordStdin: true},
	},
}

// readHttpRepInputs reads http private registry URL, username and password from the user
func readHttpRepInputs() (string, string, string) {
	isConfirm := false
	repository := ""
	username := ""
	password := ""
	var err error

	for !isConfirm {
		repository, err = utils.ReadInputString("Enter private registry (10.100.5.225:5000/jennifer)", utils.Default{Value: "", IsDefault: false}, httpRepoValidRegex, true)
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

		fmt.Println("\nRepository: " + repository)
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

func init() {
	add(HttpRegistry)
}
