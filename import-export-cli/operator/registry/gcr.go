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
	"encoding/json"
	"fmt"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io/ioutil"
	"strings"
)

var gcrRepo = new(string)

var gcrValues = struct {
	project       string
	svcAccKeyFile string
}{}

var GcrRegistry = &Registry{
	Name:       "GCR",
	Caption:    "GCR",
	Repository: gcrRepo,
	Option:     3,
	Read: func() {
		svcAccKeyFile := readGcrInputs()
		gcrValues.project = getGcrProjectName(svcAccKeyFile)
		gcrValues.svcAccKeyFile = svcAccKeyFile
		*gcrRepo = gcrValues.project
	},
	Run: func() {
		data, err := ioutil.ReadFile(gcrValues.svcAccKeyFile)
		if err != nil {
			utils.HandleErrorAndExit("Error reading GCR service account key json file", err)
		}

		k8sUtils.K8sCreateSecretFromFile(k8sUtils.GcrSvcAccKeyVolume, gcrValues.svcAccKeyFile, k8sUtils.GcrSvcAccKeyFile)
		k8sUtils.K8sCreateSecretFromInputs(k8sUtils.ConfigJsonVolume, "gcr.io", "_json_key", string(data))
	},
}

// readGcrInputs reads the GCR service account key json file from user
func readGcrInputs() string {
	isConfirm := false
	svcAccKeyFile := ""
	var err error

	for !isConfirm {
		svcAccKeyFile, err = utils.ReadInput("GCR service account key json file", utils.Default{IsDefault: false}, utils.IsFileExist, "Invalid file", true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading GCR service account key json file from user", err)
		}

		fmt.Println("")
		fmt.Println("UserGCR service account key json filename: " + svcAccKeyFile)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", utils.Default{Value: "Y", IsDefault: true}, "", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirmStr = strings.ToUpper(isConfirmStr)
		isConfirm = isConfirmStr == "Y" || isConfirmStr == "YES"
	}

	return svcAccKeyFile
}

// getGcrProjectName returns project name from the given GCR service account key json file
func getGcrProjectName(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.HandleErrorAndExit("Error reading GCR service account key json file", err)
	}
	svcAccKey := make(map[string]string)
	if err = json.Unmarshal(data, &svcAccKey); err != nil {
		utils.HandleErrorAndExit("Error unmarshal GCR service account key json file", err)
	}

	return svcAccKey["project_id"]
}

func init() {
	add(GcrRegistry)
}
