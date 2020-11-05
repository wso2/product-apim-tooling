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

// GcrRegistry represents Google Container Registry
var GcrRegistry = &Registry{
	Name:       "GCR",
	Caption:    "GCR",
	Repository: Repository{},
	Option:     3,
	Read: func(reg *Registry, flagValues *map[string]FlagValue) {
		var svcAccKeyFile string

		// check input mode: interactive or batch
		if flagValues == nil {
			// get inputs in interactive mode
			svcAccKeyFile = readGcrInputs()
		} else {
			// get inputs in batch mode
			svcAccKeyFile = (*flagValues)[k8sUtils.FlagBmKeyFile].Value.(string)

			// validate required inputs
			if !utils.IsFileExist(svcAccKeyFile) {
				utils.HandleErrorAndExit("Invalid service account key file: "+svcAccKeyFile, nil)
			}
		}

		reg.Repository.Name = getGcrProjectName(svcAccKeyFile)
		reg.Repository.KeyFile = svcAccKeyFile
	},
	Run: func(reg *Registry) {
		data, err := ioutil.ReadFile(reg.Repository.KeyFile)
		if err != nil {
			utils.HandleErrorAndExit("Error reading GCR service account key json file", err)
		}

		k8sUtils.K8sCreateSecretFromFile(k8sUtils.GcrSvcAccKeySecret, k8sUtils.ApiOpWso2Namespace,
			reg.Repository.KeyFile, k8sUtils.GcrSvcAccKeyFile)
		k8sUtils.K8sCreateSecretFromInputs(k8sUtils.GcrPullSecret, k8sUtils.ApiOpWso2Namespace,
			"gcr.io", "_json_key", string(data))
	},
	Flags: Flags{
		RequiredFlags: &map[string]bool{k8sUtils.FlagBmKeyFile: true},
		OptionalFlags: &map[string]bool{},
	},
}

// readGcrInputs reads the GCR service account key json file from user
func readGcrInputs() string {
	isConfirm := false
	svcAccKeyFile := ""
	var err error

	for !isConfirm {
		svcAccKeyFile, err = utils.ReadInput("GCR service account key json file", utils.Default{IsDefault: false},
			utils.IsFileExist, "Invalid file", true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading GCR service account key json file from user", err)
		}

		fmt.Println("\nGCR service account key json file: " + svcAccKeyFile)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", utils.Default{Value: "Y", IsDefault: true},
			"", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirm = strings.EqualFold(isConfirmStr, "y") || strings.EqualFold(isConfirmStr, "yes")
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
