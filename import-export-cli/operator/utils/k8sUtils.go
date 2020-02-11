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
	"errors"
	"fmt"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io"
	"os"
	"os/exec"
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
			if err := ExecuteCommandWithoutPrintingErrors(utils.Kubectl, utils.K8sGet, resourceType); err != nil {
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

func K8sCreateSecretFromFile(secretName string, filePath string, renamedFile string) {
	var fromFile string
	if renamedFile == "" {
		fromFile = fmt.Sprintf("--from-file=%s", filePath)
	} else {
		fromFile = fmt.Sprintf("--from-file=%s=%s", renamedFile, filePath)
	}

	// render secret
	secret, err := GetCommandOutput(
		utils.Kubectl, utils.Create, utils.K8sSecret, "generic",
		secretName, fromFile,
		"--dry-run", "-o", "yaml",
	)
	if err != nil {
		utils.HandleErrorAndExit("Error creating secret from file", err)
	}

	// apply secret
	if err = K8sApplyFromStdin(secret); err != nil {
		utils.HandleErrorAndExit("Error creating secret from file", err)
	}
}

// K8sApplyFromFile applies resources from list of files, urls or directories
func K8sApplyFromFile(fileList ...string) error {
	kubectlArgs := []string{utils.K8sApply}
	for _, file := range fileList {
		kubectlArgs = append(kubectlArgs, "-f", file)
	}

	return ExecuteCommand(utils.Kubectl, kubectlArgs...)
}

// K8sApplyFromStdin applies resources from standard input
func K8sApplyFromStdin(stdInput string) error {
	return ExecuteCommandFromStdin(stdInput, utils.Kubectl, utils.K8sApply, "-f", "-")
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
