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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	devEnvProdUrl          = "DEV_ENV_PROD_URL"
	devEnvSandUrl          = "DEV_ENV_SAND_URL"
	devEnvProdRetryDelay   = "DEV_ENV_PROD_RE_DELAY"
	devEnvProdRetryTimeOut = "DEV_ENV_PROD_RE_TO"
	envKey                 = "ENV_KEY"
)

func SetEnvVariablesForAPI(t *testing.T, client *apim.Client) {

	t.Log("Setting up the environment variable values")
	os.Setenv(devEnvProdUrl, "https://localhost:"+strconv.Itoa(9443+client.GetPortOffset())+
		"/am/sample/pizzashack/v1/api/")
	os.Setenv(devEnvSandUrl, "https://localhost:"+strconv.Itoa(9443+client.GetPortOffset())+
		"/am/sample/pizzashack/v1/api/")
	os.Setenv(devEnvProdRetryDelay, "10")
	os.Setenv(devEnvProdRetryTimeOut, "5")
	os.Setenv(envKey, "dev_101")

	t.Cleanup(func() {
		t.Log("Unsetting the environment variable values")
		os.Unsetenv(devEnvProdUrl)
		os.Unsetenv(devEnvSandUrl)
		os.Unsetenv(devEnvProdRetryDelay)
		os.Unsetenv(devEnvProdRetryTimeOut)
		os.Unsetenv(envKey)
	})
}

func ValidateDynamicData(t *testing.T, api *apim.API) {

	// Retrieve the endpointConfig of the imported API
	endpointConfig := api.GetEndPointConfig()

	// Check whether the production endpoint has the expected value set using the env variable
	productionEndpoints := endpointConfig.(map[string]interface{})["production_endpoints"].(map[string]interface{})
	assert.Equal(t, os.Getenv(devEnvProdUrl), productionEndpoints["url"], "Production endpoint value mismatched")

	// Check whether the sandbox endpoint has the expected value set using the env variable
	sandboxEndpoints := endpointConfig.(map[string]interface{})["sandbox_endpoints"].(map[string]interface{})
	assert.Equal(t, os.Getenv(devEnvSandUrl), sandboxEndpoints["url"], "Sandbox endpoint value mismatched")

	// Check whether the retryDelay and retryTimeOut roduction endpoint
	// config values has the expected values set using the env variables
	assert.Equal(t, os.Getenv(devEnvProdRetryDelay),
		productionEndpoints["config"].(map[string]interface{})["retryDelay"].(string),
		"Retry delay value of the production endpoint value mismatched")
	assert.Equal(t, os.Getenv(devEnvProdRetryTimeOut),
		productionEndpoints["config"].(map[string]interface{})["retryTimeOut"].(string),
		"Retry time out value of the production endpoint config mismatched")
}

func AddSequenceWithDynamicDataToAPIProject(t *testing.T, args *InitTestArgs) apim.APIFile {
	inSequencePath := args.InitFlag + string(os.PathSeparator) + utils.InitProjectSequencesInCustom

	base.CreateTempDir(t, inSequencePath)
	// Move sequence file to created project
	srcPathForSequence, _ := filepath.Abs(DynamicDataSampleCaseArtifactPath + string(os.PathSeparator) + CustomDynamicInSequenceName)
	destPathForSequence := inSequencePath + string(os.PathSeparator) + CustomDynamicInSequenceName
	base.Copy(srcPathForSequence, destPathForSequence)

	// Read the API definition file in the project
	apiDefinitionFilePath := args.InitFlag + string(os.PathSeparator) + utils.APIDefinitionFileYaml
	apiDefinitionFileContent := ReadAPIDefinition(t, apiDefinitionFilePath)

	mediationPolicy := apim.MediationPolicy{
		Name:   "custom-header-in",
		Type:   "IN",
		Shared: false,
	}
	mediationPolicies := []apim.MediationPolicy{mediationPolicy}
	apiDefinitionFileContent.Data.MediationPolicies = mediationPolicies

	// Write the modified API definition to the directory
	WriteToAPIDefinition(t, apiDefinitionFileContent, apiDefinitionFilePath)

	return apiDefinitionFileContent
}

func ValidateExportedSequenceWithDynamicData(t *testing.T, args *InitTestArgs, api apim.API) {
	exportedOutput := ValidateExportImportedAPI(t, args, api.Name, api.Version)

	// Unzip exported API and check whether the imported sequence is in there
	exportedPath := base.GetExportedPathFromOutput(exportedOutput)
	relativePath := strings.ReplaceAll(exportedPath, ".zip", "")
	base.Unzip(relativePath, exportedPath)

	exportedAPIRelativePath := relativePath + string(os.PathSeparator) + api.Name + "-" + api.Version
	sequencePathOfExportedAPI := exportedAPIRelativePath + string(os.PathSeparator) +
		utils.InitProjectSequencesInCustom + string(os.PathSeparator) + CustomDynamicInSequenceName

	// Check whether the sequence file is available
	isSequenceExported := base.IsFileAvailable(t, sequencePathOfExportedAPI)
	base.Log("Sequence is Exported ", isSequenceExported)
	assert.Equal(t, true, isSequenceExported, "Error while exporting API with document")

	// Read the file content of the exported sequence
	sequenceData, err := ioutil.ReadFile(sequencePathOfExportedAPI)
	if err != nil {
		t.Error(err)
	}

	// The environment variable must have been substituted twice in the sequence
	assert.Equal(t, strings.Count(string(sequenceData), os.Getenv(envKey)), 2,
		"Env variable is not substituted correctly in the sequence")

	t.Cleanup(func() {
		base.RemoveDir(relativePath)
		base.RemoveDir(exportedPath)
	})
}

func ValidateImportProjectFailedWithoutSettingEnvVariables(t *testing.T, args *InitTestArgs, paramsPath string, preserveProvider bool) {
	t.Helper()

	result, _ := importApiFromProject(t, args.InitFlag, args.APIName, paramsPath, args.SrcAPIM, &args.CtlUser, false,
		preserveProvider)

	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")

	base.WaitForIndexing()

	//Remove Created project and logout
	t.Cleanup(func() {
		base.RemoveDir(args.InitFlag)
	})
}
