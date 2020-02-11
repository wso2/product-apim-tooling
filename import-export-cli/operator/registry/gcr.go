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

var gcrRepo = new(string)

var gcrValues = struct {
	repository    string
	svcAccKeyFile string
}{}

var GcrRegistry = &Registry{
	Name:       "GCR",
	Caption:    "GCR",
	Repository: gcrRepo,
	Option:     3,
	Read: func() {
		repository, svcAccKeyFile := readGcrInputs()
		*gcrRepo = repository
		gcrValues.repository = repository
		gcrValues.svcAccKeyFile = svcAccKeyFile
	},
	Run: func() {
		k8sUtils.K8sCreateSecretFromFile(k8sUtils.GcrSvcAccKeyVolume, gcrValues.svcAccKeyFile, k8sUtils.GcrSvcAccKeyFile)
	},
}

// readDockerHubInputs reads docker-registry URL, username and password from the user
func readGcrInputs() (string, string) {
	isConfirm := false
	repository := ""
	svcAccKeyFile := ""
	var err error

	for !isConfirm {
		repository, err = utils.ReadInputString("Enter project name", utils.Default{Value: "", IsDefault: true}, utils.UsernameValidRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading GCR project name from user", err)
		}

		svcAccKeyFile, err = utils.ReadInput("GCR service account key json file", utils.Default{IsDefault: false}, utils.IsFileExist, "Invalid file", true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading GCR service account key json file from user", err)
		}

		fmt.Println("")
		fmt.Println("Project                                  : " + repository)
		fmt.Println("UserGCR service account key json filename: " + svcAccKeyFile)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", utils.Default{Value: "Y", IsDefault: true}, "", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirmStr = strings.ToUpper(isConfirmStr)
		isConfirm = isConfirmStr == "Y" || isConfirmStr == "YES"
	}

	return repository, svcAccKeyFile
}

func init() {
	add(GcrRegistry)
}
