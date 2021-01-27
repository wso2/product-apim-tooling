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

package testutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

// ExecActivateCommand run activate artifactType artifactName
func ExecActivateCommand(t *testing.T, config *MiConfig, artifactType, artifactName, expected string) {
	t.Helper()
	execActivateDeactivateCommand(t, config, "activate", artifactType, artifactName, expected)
}

// ExecDeactivateCommand run deactivate artifactType artifactName
func ExecDeactivateCommand(t *testing.T, config *MiConfig, artifactType, artifactName, expected string) {
	t.Helper()
	execActivateDeactivateCommand(t, config, "deactivate", artifactType, artifactName, expected)
}

func execActivateDeactivateCommand(t *testing.T, config *MiConfig, mode, artifactType, artifactName, expected string) {
	SetupAndLoginToMI(t, config)
	response, _ := base.Execute(t, "mi", mode, artifactType, artifactName, "-e", config.MIClient.GetEnvName(), "-k")
	base.Log(response)
	assert.Contains(t, response, expected)
}

// ExecActivateCommandWithoutSettingEnv run activate without setting up an environment
func ExecActivateCommandWithoutSettingEnv(t *testing.T, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithoutSettingEnv(t, "activate", args)
}

// ExecDeactivateCommandWithoutSettingEnv run deactivate without setting up an environment
func ExecDeactivateCommandWithoutSettingEnv(t *testing.T, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithoutSettingEnv(t, "deactivate", args)
}

func execActivateDeactivateCommandWithoutSettingEnv(t *testing.T, mode string, args []string) {
	getCmdArgs := []string{"mi", mode, "-e", "testing", "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	assert.Contains(t, response, "MI does not exists in testing Add it using add env")
}

// ExecActivateCommandWithoutLogin run activate artifactType artifactName without login to MI
func ExecActivateCommandWithoutLogin(t *testing.T, config *MiConfig, artifactType, artifactName string, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithoutLogin(t, config, "activate", artifactType, artifactName, args)
}

// ExecDeactivateCommandWithoutLogin run deactivate artifactType artifactName without login to MI
func ExecDeactivateCommandWithoutLogin(t *testing.T, config *MiConfig, artifactType, artifactName string, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithoutLogin(t, config, "deactivate", artifactType, artifactName, args)
}

func execActivateDeactivateCommandWithoutLogin(t *testing.T, config *MiConfig, mode, artifactType, artifactName string, args []string) {
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	getCmdArgs := []string{"mi", mode, artifactType, artifactName, "-e", config.MIClient.GetEnvName(), "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	assert.Contains(t, response, "Login to MI")
}

// ExecActivateCommandWithoutEnvFlag run activate artifactType artifactName without -e flag
func ExecActivateCommandWithoutEnvFlag(t *testing.T, config *MiConfig, artifactType, artifactName string, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithoutEnvFlag(t, config, "activate", artifactType, artifactName, args)
}

// ExecDeactivateCommandWithoutEnvFlag run deactivate artifactType artifactName without -e flag
func ExecDeactivateCommandWithoutEnvFlag(t *testing.T, config *MiConfig, artifactType, artifactName string, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithoutEnvFlag(t, config, "deactivate", artifactType, artifactName, args)
}

func execActivateDeactivateCommandWithoutEnvFlag(t *testing.T, config *MiConfig, mode, artifactType, artifactName string, args []string) {
	SetupAndLoginToMI(t, config)
	getCmdArgs := []string{"mi", mode, artifactType, artifactName, "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	assert.Contains(t, response, `required flag(s) "environment" not set`)
}

// ExecActivateCommandWithInvalidArgCount run activate artifactType artifactName with invalid number of args
func ExecActivateCommandWithInvalidArgCount(t *testing.T, config *MiConfig, required, passed int, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithInvalidArgs(t, config, "activate", required, passed, args)
}

// ExecDeactivateCommandWithInvalidArgCount run deactivate artifactType artifactName with invalid number of args
func ExecDeactivateCommandWithInvalidArgCount(t *testing.T, config *MiConfig, required, passed int, args ...string) {
	t.Helper()
	execActivateDeactivateCommandWithInvalidArgs(t, config, "deactivate", required, passed, args)
}

func execActivateDeactivateCommandWithInvalidArgs(t *testing.T, config *MiConfig, mode string, required, passed int, args []string) {
	SetupAndLoginToMI(t, config)
	getCmdArgs := []string{"mi", mode, "-k"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.Log(response)
	expected := fmt.Sprintf("accepts %v arg(s), received %v", required, passed)
	assert.Contains(t, response, expected)
}
