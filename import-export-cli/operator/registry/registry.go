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
	"errors"
	"fmt"
	utils2 "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"sort"
)

type Registry struct {
	Name    string
	Caption string
	Option  int
	Read    func()
	Run     func()
}

var registries = make(map[int]*Registry)
var optionToExec int

func ReadInputs() {
	registries[optionToExec].Read()
}

func CreateSecret() {
	registries[optionToExec].Run()
}

func ChooseRegistry() {
	keys := make([]int, 0, len(registries))
	for key := range registries {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	// print all repository types
	fmt.Println("Choose repository type:")
	for _, key := range keys {
		fmt.Printf("%d: %s\n", key, registries[key].Caption)
	}

	option, err := utils.ReadOption("Choose a number", 1, len(keys), true)
	if err != nil {
		utils.HandleErrorAndExit("Error reading registry type", err)
	}

	optionToExec = option
}

// setRegistryRepositoryOnControllerConfig sets the repository value in the config: `controller-config`
func setRegistryRepositoryOnControllerConfig(repository string) {
	// get controller config config map
	controllerConfigMapYaml, err := utils2.GetCommandOutput(
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

	controllerConfigMap["data"].(map[interface{}]interface{})["dockerRegistry"] = repository
	configuredConfigMap, err := yaml.Marshal(controllerConfigMap)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering controller-config", err)
	}

	// apply controller config config map back
	if err := utils2.K8sApplyFromStdin(string(configuredConfigMap)); err != nil {
		utils.HandleErrorAndExit("Error creating controller-configs", err)
	}
}

func add(registry *Registry) {
	if registry.Option < 1 {
		utils.HandleErrorAndExit("Error adding registry", errors.New("option should be positive"))
	}

	registries[registry.Option] = registry
}
