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
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"sort"
)

// Registry represents Docker Registry
type Registry struct {
	Name       string                                                // Unique Name
	Caption    string                                                // Text to display in the CLI about registry details
	Repository Repository                                            // Repository name
	Option     int                                                   // Option to be choose the CLI registry list
	Read       func(reg *Registry, flagValues *map[string]FlagValue) // Function to be called when getting inputs, if flagValues is nil get inputs interactively
	Run        func(reg *Registry)                                   // Function to be called when updating k8s secrets
	Flags      Flags                                                 // Required and Optional flags
}

// Repository represents a Docker repository of a Docker registry
type Repository struct {
	Name      string
	ServerUrl string
	Username  string
	Password  string
	KeyFile   string
}

// Flags represents Required and Optional flags that supports the specified registry type
type Flags struct {
	RequiredFlags *map[string]bool // Map of flag name and bool value of the flag is required (true) or not (false)
	OptionalFlags *map[string]bool // Map of flag name and bool value of the flag is optional (true) or not (false)
}

// FlagValue represents a value of a flag and its value supplied by user
type FlagValue struct {
	Value      interface{} // Value of the flag
	IsProvided bool        // Is the value provided by the user
}

// registries represents a map of registries
var registries = make(map[int]*Registry)

// optionToExec represents the choice use selected
var optionToExec int

// ReadInputsInteractive reads inputs with respect to the selected registry type interactively
func ReadInputsInteractive() {
	reg := registries[optionToExec]
	reg.Read(reg, nil)
}

// ReadInputsFromFlags reads inputs from flags with respect to the selected registry type
func ReadInputsFromFlags(flagValues *map[string]FlagValue) {
	reg := registries[optionToExec]
	reg.Read(reg, flagValues)
}

// UpdateConfigsSecrets updates controller config with registry type and creates secrets with credentials
func UpdateConfigsSecrets() {
	// set registry first since this can throw error if api operator not installed. If error occur no need to rollback secret.
	updateDockerRegistryConfig(registries[optionToExec].Name, registries[optionToExec].Repository.Name)
	// create secret
	registries[optionToExec].Run(registries[optionToExec])
}

// ChooseRegistryInteractive lists registries in the CLI and reads a choice from user
func ChooseRegistryInteractive() {
	keys := make([]int, 0, len(registries))
	for key := range registries {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	// print all repository types
	fmt.Println("Choose registry type:")
	for _, key := range keys {
		fmt.Printf("%d: %s\n", key, registries[key].Caption)
	}

	option, err := utils.ReadOption("Choose a number", 1, len(keys), true)
	if err != nil {
		utils.HandleErrorAndExit("Error reading registry type", err)
	}

	optionToExec = option
}

// SetRegistry set the private value 'optionToExec' that match with 'registryType' un-interactively
func SetRegistry(registryType string) {
	for opt, reg := range registries {
		if reg.Name == registryType {
			optionToExec = opt
			return
		}
	}

	// if not found throw error: invalid registry type
	utils.HandleErrorAndExit("Invalid registry type: "+registryType, nil)
}

// ValidateFlags validates if any additional flag is given or any required flag is missing
// throw error if invalid
func ValidateFlags(flagsValues *map[string]FlagValue) {
	// check for required flags
	for flg, flgRequired := range *registries[optionToExec].Flags.RequiredFlags {
		if flgRequired && !(*flagsValues)[flg].IsProvided {
			// required flag is missing
			utils.HandleErrorAndExit("Required flag is missing in batch mode. Flag: "+flg, nil)
		}
	}

	// check for additional flags
	for flg, flgVal := range *flagsValues {
		if flgVal.IsProvided && !(*registries[optionToExec].Flags.RequiredFlags)[flg] && !(*registries[optionToExec].Flags.OptionalFlags)[flg] {
			// additional, not supported flag
			utils.HandleErrorAndExit("Invalid, not supported flag found in batch mode. Flag: "+flg, nil)
		}
	}

	// flag validation success and continue the flow
}

// updateDockerRegistryConfig sets the repository type value and the repository in the config: `controller-config`
func updateDockerRegistryConfig(registryType string, repository string) {
	// get controller config config map
	registryConfigMapYaml, _ := box.Get("/kubernetes_resources/docker_registry_conf.yaml")

	registryConfigMap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(registryConfigMapYaml), &registryConfigMap); err != nil {
		utils.HandleErrorAndExit("Error reading controller-config", err)
	}

	// set configurations
	registryConfigMap["data"].(map[interface{}]interface{})[k8sUtils.CtrlConfigRegType] = registryType
	registryConfigMap["data"].(map[interface{}]interface{})[k8sUtils.CtrlConfigReg] = repository

	configuredRegConfigMap, err := yaml.Marshal(registryConfigMap)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering controller-config", err)
	}

	// apply controller config config map back
	if err := k8sUtils.K8sApplyFromStdin(string(configuredRegConfigMap)); err != nil {
		utils.HandleErrorAndExit("Error creating controller-configs", err)
	}
}

// add adds a registry to the registries maps
// using pointers for memory optimization
func add(registry *Registry) {
	if registry.Option < 1 {
		utils.HandleErrorAndExit("Error adding registry: "+registry.Name, errors.New("'option' should be positive"))
	}
	if registries[registry.Option] != nil {
		utils.HandleErrorAndExit("Error adding registry"+registry.Name, errors.New("duplicate 'options' values"))
	}

	registries[registry.Option] = registry
}
