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
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"sort"
)

// Registry represents Docker Registry
type Registry struct {
	Name       string  // Unique Name
	Caption    string  // Text to display in the CLI about registry details
	Repository *string // Repository name
	Option     int     // Option to be choose the CLI registry list
	Read       func()  // Function to be called when getting inputs
	Run        func()  // Function to be called when updating k8s secrets
	Flags      Flags
}

type Flags struct {
	RequiredFlags *[]string
	OptionalFlags *[]string
}

// registries represents a map of registries
var registries = make(map[int]*Registry)

// optionToExec represents the choice use selected
var optionToExec int

func ReadInputs() {
	registries[optionToExec].Read()
}

// UpdateConfigsSecrets updates controller config with registry type and creates secrets with credentials
func UpdateConfigsSecrets() {
	// set registry first since this can throw error if api operator not installed. If error occur no need to rollback secret.
	updateCtrlConfig(registries[optionToExec].Name, *registries[optionToExec].Repository)
	// create secret
	registries[optionToExec].Run()
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
func ValidateFlags(flags *[]string) {
	// check for required flags
	for _, flag := range *registries[optionToExec].Flags.RequiredFlags {
		if !k8sUtils.StringArrayContains(*flags, flag) {
			// required flag is missing
			utils.HandleErrorAndExit("Required flag is missing in un-interactive mode. Flag: "+flag, nil)
		}
	}

	// check for additional flags
	for _, flag := range *flags {
		if !k8sUtils.StringArrayContains(append(
			*registries[optionToExec].Flags.RequiredFlags,
			*registries[optionToExec].Flags.OptionalFlags...,
		), flag) {
			// additional, not supported flag
			utils.HandleErrorAndExit("Additional not supported flag found in un-interactive mode. Flag: "+flag, nil)
		}
	}

	// flag validation success and continue the flow
}

// updateCtrlConfig sets the repository type value and the repository in the config: `controller-config`
func updateCtrlConfig(registryType string, repository string) {
	// get controller config config map
	controllerConfigMapYaml, err := k8sUtils.GetCommandOutput(
		k8sUtils.Kubectl, k8sUtils.K8sGet, k8sUtils.K8sConfigMap, k8sUtils.ApiOpControllerConfigMap,
		"-n", k8sUtils.ApiOpWso2Namespace,
		"-o", "yaml",
	)
	if err != nil {
		utils.HandleErrorAndExit("Error reading controller-config.\nInstall api operator using the command: apictl install api-operator",
			errors.New("error reading controller-config"))
	}

	controllerConfigMap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(controllerConfigMapYaml), &controllerConfigMap); err != nil {
		utils.HandleErrorAndExit("Error reading controller-config", err)
	}

	// set configurations
	controllerConfigMap["data"].(map[interface{}]interface{})[k8sUtils.CtrlConfigRegType] = registryType
	controllerConfigMap["data"].(map[interface{}]interface{})[k8sUtils.CtrlConfigReg] = repository

	configuredConfigMap, err := yaml.Marshal(controllerConfigMap)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering controller-config", err)
	}

	// apply controller config config map back
	if err := k8sUtils.K8sApplyFromStdin(string(configuredConfigMap)); err != nil {
		utils.HandleErrorAndExit("Error creating controller-configs", err)
	}
}

// add adds a registry to the registries maps
func add(registry *Registry) {
	if registry.Option < 1 {
		utils.HandleErrorAndExit("Error adding registry: "+registry.Name, errors.New("'option' should be positive"))
	}
	if registries[registry.Option] != nil {
		utils.HandleErrorAndExit("Error adding registry"+registry.Name, errors.New("duplicate 'options' values"))
	}

	registries[registry.Option] = registry
}
