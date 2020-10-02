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

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func K8sWaitForResourceType(maxTimeSec int, resourceTypes ...string) error {
	if maxTimeSec < 0 {
		return errors.New("'maxTimeSec' should be non negative")
	}

	noErrors := false
	for i := maxTimeSec; i > 0 && !noErrors; i-- {
		for _, resourceType := range resourceTypes {
			noErrors = true
			if err := ExecuteCommandWithoutPrintingErrors(Kubectl, K8sGet, resourceType); err != nil {
				noErrors = false
				continue
			}
		}

		time.Sleep(1e9) // sleep 1 second
	}

	if !noErrors {
		return errors.New("kubernetes resources not installed")
	}

	return nil
}

// K8sCreateSecretFromInputs creates K8S a docker-registry secret with given inputs
func K8sCreateSecretFromInputs(secretName string, namespace string, server string, username string, password string) {
	if username == "" {
		username = "N/A"
		password = "N/A"
	}

	// render secret
	dockerCredentialsMap := RenderSecretTemplate(secretName, namespace, K8sSecret)
	dockerAuth := credentials.Base64Encode(username + ":" + password)

	type JSON map[string]interface{}
	authJson := JSON{"auths":JSON{server:JSON{"username":username,"password":password,"auth":dockerAuth}}}
	out, err := json.Marshal(authJson)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering JSON configurations for .dockerconfigjson", err)
	}
	dockerConfEncoded := credentials.Base64Encode(string(out))

	dockerCredentialsMap["data"] = make(map[interface{}]interface{})
	dockerCredentialsMap["data"].(map[interface{}]interface{})[DockerConfigJson] = dockerConfEncoded
	dockerCredentialsMap["type"] = "kubernetes.io/dockerconfigjson"

	configuredCredentials, err := yaml.Marshal(dockerCredentialsMap)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering docker credentials file", err)
	}

	if err := K8sApplyFromStdin(string(configuredCredentials)); err != nil {
		utils.HandleErrorAndExit("Error creating docker credentials", err)
	}
}

// K8sCreateSecretFromFile creates K8S a generic secret with give file
func K8sCreateSecretFromFile(secretName string, namespace string, filePath string, renamedFile string) {

	// render secret
	secretsMap := RenderSecretTemplate(secretName, namespace, K8sSecret)

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.HandleErrorAndExit("Error reading the specified credentials file", err)
	}
	encodedData := credentials.Base64Encode(string(fileData))

	secretsMap["data"] = make(map[interface{}]interface{})
	secretsMap["data"].(map[interface{}]interface{})[renamedFile] = encodedData

	configuredSecret, err := yaml.Marshal(secretsMap)
	if err != nil {
		utils.HandleErrorAndExit("Error rendering registry credentials from file", err)
	}

	if err := K8sApplyFromStdin(string(configuredSecret)); err != nil {
		utils.HandleErrorAndExit("Error creating registry credentials from file", err)
	}
}

// K8sApplyFromFile applies resources from list of files, urls or directories
func K8sApplyFromFile(fileList ...string) error {
	kubectlArgs := []string{K8sApply}
	for _, file := range fileList {
		kubectlArgs = append(kubectlArgs, "-f", file)
	}

	return ExecuteCommand(Kubectl, kubectlArgs...)
}

// K8sApplyFromBytes applies resources by content
func K8sApplyFromBytes(data [][]byte) error {
	dir, _ := ioutil.TempDir("", "example")
	defer os.RemoveAll(dir) // clean up

	for i, d := range data {
		tmpFile := filepath.Join(dir, fmt.Sprintf("config-file-%v.yml", i))
		err := ioutil.WriteFile(tmpFile, d, 0666) // permission -rw-rw-rw
		if err != nil {
			return err
		}
	}

	return ExecuteCommand(Kubectl, K8sApply, "-f", dir)
}

// K8sApplyFromStdin applies resources from standard input
func K8sApplyFromStdin(stdInputs string) error {
	return ExecuteCommandFromStdin(stdInputs, Kubectl, K8sApply, "-f", "-")
}

// ExecuteCommand executes the command with args and prints output, errors in standard output, error
func ExecuteCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	setCommandOutAndError(cmd)
	return cmd.Run()
}

// ExecuteCommandWithoutPrintingErrors executes the command with args and prints output, standard output
func ExecuteCommandWithoutPrintingErrors(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	setCommandOutOnly(cmd)
	return cmd.Run()
}

// ExecuteCommandWithoutOutputs executes the command with args without printing outputs and errors
func ExecuteCommandWithoutOutputs(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

// ExecuteCommandFromStdin executes the command with args and prints output the standard output
func ExecuteCommandFromStdin(stdInput string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	setCommandOutAndError(cmd)

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if _, err = pipe.Write([]byte(stdInput)); err != nil {
		return err
	}
	if err := pipe.Close(); err != nil {
		return err
	}

	return cmd.Run()
}

// GetCommandOutput executes a command and returns the output
func GetCommandOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var errBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)

	output, err := cmd.Output()
	return string(output), err
}

// setCommandOutAndError sets the output and error of the command cmd to the standard output and error
func setCommandOutAndError(cmd *exec.Cmd) {
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
}

// setCommandOutOnly sets the output the command cmd to the standard output and not sets the error
func setCommandOutOnly(cmd *exec.Cmd) {
	var outBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
}

// RenderSecretTemplate renders the registry related configmaps/secrets using docker_credentials.yaml
func RenderSecretTemplate(secretName string, namespace string, kind string) (map[interface{}]interface{}){
	secretConfigsYaml, _ := box.Get("/kubernetes_resources/docker_credentials.yaml")
	secretConfigsMap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(secretConfigsYaml), &secretConfigsMap); err != nil {
		utils.HandleErrorAndExit("Error reading registry config template", err)
	}

	secretConfigsMap[kindKey] = kind
	secretConfigsMap["metadata"].(map[interface{}]interface{})["name"] = secretName
	secretConfigsMap["metadata"].(map[interface{}]interface{})["namespace"] = namespace

	return secretConfigsMap
}
