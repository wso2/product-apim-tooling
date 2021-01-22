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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

// ListArtifacts return ctl out from the command get artifactType
func ListArtifacts(t *testing.T, artifactType string, config *MiConfig) (string, error) {
	t.Helper()
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	output, err := base.Execute(t, "mi", "get", artifactType, "-e", config.MIClient.GetEnvName())
	return output, err
}

// GetArtifact return ctl out from the command get artifactType artifactName
func GetArtifact(t *testing.T, artifactType, artifactName string, config *MiConfig) (string, error) {
	t.Helper()
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, config.MIClient.GetEnvName(), config.Username, config.Username)
	output, err := base.Execute(t, "mi", "get", artifactType, artifactName, "-e", config.MIClient.GetEnvName())
	return output, err
}

// GetArtifactListFromAPI : Get Artifact Lists from MI Management API
func (instance *MiRESTClient) GetArtifactListFromAPI(resource string, artifactListType interface{}) interface{} {
	apisURL := getResourceURL(instance.GetMiURL(), resource)

	request := base.CreateGet(apisURL)
	base.SetDefaultRestAPIHeaders(instance.accessToken, request)
	base.LogRequest("mi.GetArtifactList()", request)
	response := base.SendHTTPRequest(request)
	defer response.Body.Close()

	base.ValidateAndLogResponse("mi.GetArtifactList()", response, 200)

	artifactListResponse := artifactListType
	json.NewDecoder(response.Body).Decode(&artifactListResponse)
	return artifactListResponse
}

// GetArtifactFromAPI : Get Artifacts from MI Management API
func (instance *MiRESTClient) GetArtifactFromAPI(resource, artifactKey, artifactName string, artifactType interface{}) interface{} {
	apisURL := getResourceURLWithQueryParam(instance.GetMiURL(), resource, artifactKey, artifactName)

	request := base.CreateGet(apisURL)
	base.SetDefaultRestAPIHeaders(instance.accessToken, request)
	base.LogRequest("mi.GetArtifact()", request)
	response := base.SendHTTPRequest(request)
	defer response.Body.Close()

	base.ValidateAndLogResponse("mi.GetArtifact()", response, 200)

	artifactListResponse := artifactType
	json.NewDecoder(response.Body).Decode(&artifactListResponse)
	return artifactListResponse
}

// ExecGetCommandWithoutSettingEnv run get artifactType without setting up an environment
func ExecGetCommandWithoutSettingEnv(t *testing.T, artifactType string) {
	t.Helper()
	response, _ := base.Execute(t, "mi", "get", artifactType, "-e", "testing")
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, "MI does not exists in testing Add it using add env")
}

// ExecGetCommandWithoutLogin run get artifactType without login to MI
func ExecGetCommandWithoutLogin(t *testing.T, artifactType string, config *MiConfig) {
	t.Helper()
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	response, _ := base.Execute(t, "mi", "get", artifactType, "-e", "testing")
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, "Login to MI")
}

// ExecGetCommandWithoutEnvFlag run get artifactType without -e flag
func ExecGetCommandWithoutEnvFlag(t *testing.T, artifactType string, config *MiConfig) {
	t.Helper()
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, "testing", AdminUserName, AdminPassword)
	response, _ := base.Execute(t, "mi", "get", artifactType)
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	assert.Contains(t, response, `required flag(s) "environment" not set`)
}

// ExecGetCommandWithInvalidArgs run get artifactType with invalid number of args
func ExecGetCommandWithInvalidArgs(t *testing.T, config *MiConfig, required, passed int, args ...string) {
	t.Helper()
	base.SetupMIEnv(t, config.MIClient.GetEnvName(), config.MIClient.GetMiURL())
	base.MILogin(t, "testing", AdminUserName, AdminPassword)
	getCmdArgs := []string{"mi", "get"}
	getCmdArgs = append(getCmdArgs, args...)
	response, _ := base.Execute(t, getCmdArgs...)
	base.GetRowsFromTableResponse(response)
	base.Log(response)
	expected := fmt.Sprintf("accepts at most %v arg(s), received %v", required, passed)
	assert.Contains(t, response, expected)
}
