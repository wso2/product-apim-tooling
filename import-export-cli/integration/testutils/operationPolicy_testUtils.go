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
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"gopkg.in/yaml.v2"
)

const (
	PolicyVersionKey = "version"
	PolicyNameKey    = "name"
	PolicyIDReset    = "-1"
	PolicyIdKey      = "id"
)

const (
	DefaultPolicyListSize       = "-1"
	DefaultPolicyListOffsetSize = "-1"
	DefaultPolicyListLimit      = "25"
	DefaultPolicyListOffset     = "0"
)

// ValidateAPIPolicyExportImport : Validates Exporting API Policy from source env and Importing to destination env
func ValidateAPIPolicyExportImport(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[PolicyNameKey])
	policyVersion := fmt.Sprintf("%v", args.Policy[PolicyVersionKey])

	args.SrcAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	exportedOutput, _ := exportAPIPolicy(t, policyName, args)

	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)

	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importedOutput, err := importAPIPolicy(t, args.CtlUser.Username, args.CtlUser.Password, policyName, policyVersion, true, args)

	assert.Nil(t, err, "Error while importing the API Policy")
	assert.Contains(t, importedOutput, "Successfully Imported API Policy")
	// Give time for newly imported API Policy to get indexed
	base.WaitForIndexing()

	args.DestAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	policyId := args.DestAPIM.GetAPIPolicyID(t, policyName, policyVersion)

	// Get API Policy from env 2
	importedPolicy := args.DestAPIM.GetAPIPolicy(policyId)

	// Validate env 1 and env 2 policy is equal
	ValidatePoliciesEqual(t, args, importedPolicy)
}

// ValidateAPIPolicyImportWithDirectoryPath : Validates Importing API Policy with directory path given
func ValidateAPIPolicyImportWithDirectoryPath(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	pathToPolicySpecFile, err := filepath.Abs(CustomAddLogMessagePolicyDefinitionPathImport)

	assert.Nil(t, err, "Error in getting absolute path")

	policyContent := readAPIPolicyDefinition(t, pathToPolicySpecFile)

	assert.NotNil(t, policyContent, "Error in reading policy definition file")

	policyName := policyContent.Name
	policyVersion := policyContent.Version

	// Import API Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importedOutput, err := importAPIPolicy(t, args.CtlUser.Username, args.CtlUser.Password, policyName, policyVersion, true, args)

	assert.Nil(t, err, "Error while importing the API Policy")
	assert.Contains(t, importedOutput, "Successfully Imported API Policy")
	// Give time for newly imported API Policy to get indexed
	base.WaitForIndexing()

	args.DestAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	policyId := args.DestAPIM.GetAPIPolicyID(t, policyName, policyVersion)

	assert.NotNil(t, policyId, "Policy import was not successful!")
}

// ValidateAPIPolicyImportWithInconsistentFileNames : Validates Importing API Policy withinconsistent policy file names
func ValidateAPIPolicyImportWithInconsistentFileNames(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	policyContent := readAPIPolicyDefinition(t, DevSampleCaseInconsistentOperationPolicyDefinitionPath)

	assert.NotNil(t, policyContent, "Error in reading policy definition file")

	policyName := policyContent.Name
	policyVersion := policyContent.Version

	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	_, err := importAPIPolicy(t, args.CtlUser.Username, args.CtlUser.Password, policyName, policyVersion, false, args)

	assert.NotNil(t, err, "Error while importing the API Policy")

}

// ValidateAPIPolicyExportImport : Validates Exporting API Policy from source env and Importing to destination env
func ValidateAPIPolicyExportImportWithFormatFlag(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[PolicyNameKey])
	policyVersion := fmt.Sprintf("%v", args.Policy[PolicyVersionKey])

	args.SrcAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	exportedOutput, _ := exportAPIPolicyWithFormat(t, policyName, args)

	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)

	zipFile, err := zip.OpenReader(args.ImportFilePath)

	assert.Nil(t, err, "Error while reading ZIP file")

	defer zipFile.Close()

	policyDef := policyName + ".json"

	isPolicyFileExist := false

	for _, file := range zipFile.File {
		if policyDef == strings.Split(file.Name, "/")[1] {
			isPolicyFileExist = true
		}
	}

	assert.True(t, isPolicyFileExist, "Policy Definition File does not exist in JSON format")

	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importedOutput, err := importAPIPolicy(t, args.CtlUser.Username, args.CtlUser.Password, policyName, policyVersion, true, args)

	assert.Nil(t, err, "Error while importing the API Policy")
	assert.Contains(t, importedOutput, "Successfully Imported API Policy")
	// Give time for newly imported API Policy to get indexed
	base.WaitForIndexing()

	args.DestAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	policyId := args.DestAPIM.GetAPIPolicyID(t, policyName, policyVersion)

	// Get API Policy from env 2
	importedPolicy := args.DestAPIM.GetAPIPolicy(policyId)

	// Validate env 1 and env 2 policy is equal
	ValidatePoliciesEqual(t, args, importedPolicy)
}

// ValidateMalformedAPIPolicyExportImport : Validates Exporting API Policy from source env and Importing to destination env
func ValidateMalformedAPIPolicyExportImport(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	pathToPolicyDirectory, _ := filepath.Abs(DevFirstSampleCaseMalformedOperationPolicyArtifactPath)
	pathToPolicyhSpecFile, _ := filepath.Abs(DevSampleCaseMalformedOperationPolicyDefinitionPath)

	args.ImportFilePath = pathToPolicyDirectory + "/"

	policyName := TestSamplePolicyName
	policyVersion := DefaultAPIPolicyVersion

	assert.True(t, base.IsFileAvailable(t, pathToPolicyhSpecFile))

	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importOutput, err := importAPIPolicy(t, args.CtlUser.Username, args.CtlUser.Password, policyName, policyVersion, false, args)

	assert.NotNil(t, err, "Error while importing the API Policy")
	assert.Contains(t, importOutput, "500", "Error importing API Policy")
}

// ValidateMalformedAPIPolicyExportImport : Validates Exporting API Policy from source env and Importing to destination env
func ValidateAPIPolicyImportFailureWhenPolicyExisted(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[PolicyNameKey])
	policyVersion := fmt.Sprintf("%v", args.Policy[PolicyVersionKey])

	args.SrcAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	exportedOutput, _ := exportAPIPolicy(t, policyName, args)

	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)

	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Import API Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	args.DestAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	importedOutput, err := importAPIPolicy(t, args.CtlUser.Username, args.CtlUser.Password, policyName, policyVersion, false, args)

	// policyId := args.DestAPIM.GetAPIPolicyID(t, policyName, policyVersion)

	assert.Nil(t, err, "Error while importing the API Policy")
	assert.Contains(t, importedOutput, "Successfully Imported API Policy")

	// Give time for newly imported API Policy to get indexed
	base.WaitForIndexing()

	importedOutput, err = importAPIPolicy(t, args.CtlUser.Username, args.CtlUser.Password, policyName, policyVersion, true, args)

	assert.NotNil(t, err, "Error importing API Policy")
	assert.Contains(t, importedOutput, "Error importing API Policy")

}

// Adds a new API Policy to an env
func AddNewAPIPolicy(t *testing.T, client *apim.Client, username, password, pathToSpec, pathToSynapse string) interface{} {
	client.Login(username, password)
	pathToPolicySpecFile, _ := filepath.Abs(pathToSpec)
	pathToSynapseDefFile, _ := filepath.Abs(pathToSynapse)
	policyContent := readAPIPolicyDefinition(t, pathToPolicySpecFile)
	policySpecFileData, err := json.Marshal(policyContent)
	doClean := true

	if err != nil {
		t.Fatal(err)
	}

	createdPolicy := client.AddAPIPolicy(t, policySpecFileData, pathToSynapseDefFile, username, password, CleanUpFunction, doClean)
	return createdPolicy
}

// Exports API Policy from an env
func exportAPIPolicy(t *testing.T, name string, args *PolicyImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "export", "policy", "api", "-n", name, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	t.Cleanup(func() {
		removeExportedAPIPolicyFile(t, args.ImportFilePath)
	})
	return output, err
}

// Cleanup func for exported api policy file
func removeExportedAPIPolicyFile(t *testing.T, file string) {
	t.Log("base.RemoveExportedAPIPolicyFile() - file path:", file)
	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Exports API Policy with JSON/YAML Format Policy Definition from an env
func exportAPIPolicyWithFormat(t *testing.T, name string, args *PolicyImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "export", "policy", "api", "-n", name, "--format", args.ExportFormat, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")

	t.Cleanup(func() {
		removeExportedAPIPolicyFile(t, args.ImportFilePath)
	})

	return output, err
}

// Imports API policy to an env
func importAPIPolicy(t *testing.T, username, password, policyName, policyVersion string, doClean bool, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error

	output, err = base.Execute(t, "import", "policy", "api", "-e", args.DestAPIM.GetEnvName(), "-f", args.ImportFilePath)

	if doClean {
		t.Cleanup(func() {
			args.DestAPIM.Login(username, password)
			args.DestAPIM.DeleteAPIPolicyByNameAndVersion(t, policyName, policyVersion, CleanUpFunction)
		})
	}

	return output, err
}

// Validates whether throttling policies are equal
func ValidatePoliciesEqual(t *testing.T, args *PolicyImportExportTestArgs, importedPolicy map[string]interface{}) {
	exportedPolicy := args.Policy
	exportedPolicy[PolicyIdKey] = PolicyIDReset
	importedPolicy[PolicyIdKey] = PolicyIDReset
	assert.Equal(t, exportedPolicy, importedPolicy)
}

// Converts API policy struct to map
func APIPolicyStructToMap(policy interface{}) (map[string]interface{}, error) {
	var apiPolicy map[string]interface{}
	marshalled, _ := json.Marshal(policy)
	err := json.Unmarshal(marshalled, &apiPolicy)
	return apiPolicy, err
}

// Validates whether the api policy list is complete
func ValidateAPIPoliciesList(t *testing.T, jsonArray bool, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIPolicies(t, jsonArray, args)

	args.SrcAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	apiPolicyList := args.SrcAPIM.GetAPIPolicies(t, DefaultPolicyListOffsetSize, DefaultPolicyListSize)

	validateListAPIPoliciesEqual(t, output, apiPolicyList)
}

// Validates whether the api policy list is complete
func ValidateAPIPoliciesListWithLimit(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIPoliciesWithLimit(t, args)

	var policySpecDataList []apim.PolicySpecData

	err := json.Unmarshal([]byte(output), &policySpecDataList)

	if err != nil {
		t.Fatal(err)
	}

	args.SrcAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	apiPolicyList := args.SrcAPIM.GetAPIPolicies(t, TestAPIPolicyOffset, TestAPIPolicyLimit)

	limit, err := strconv.Atoi(TestAPIPolicyLimit)
	if err != nil {
		t.Fatal(err)
	}

	validateListAPIPoliciesEqualWithLimit(t, policySpecDataList, apiPolicyList, limit)
}

// Validates whether the api policy list is complete with no limit flag
func ValidateAPIPoliciesListWithDefaultLimit(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIPoliciesWithNoLimit(t, args)

	var policySpecDataList []apim.PolicySpecData

	err := json.Unmarshal([]byte(output), &policySpecDataList)

	if err != nil {
		t.Fatal(err)
	}

	args.SrcAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	apiPolicyList := args.SrcAPIM.GetAPIPolicies(t, DefaultPolicyListOffset, DefaultPolicyListLimit)

	limit, err := strconv.Atoi(DefaultPolicyListLimit)
	if err != nil {
		t.Fatal(err)
	}

	validateListAPIPoliciesEqualWithLimit(t, policySpecDataList, apiPolicyList, limit)
}

// Validates whether the api policy list is complete with no limit flag
func ValidateAPIPoliciesListWithLimitAndAllFlagsReturnError(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, err := listAPIPoliciesWithLimitAndAllFlags(t, args)

	assert.Contains(t, output, "if any flags in the group [limit all] are set none of the others can be")

	assert.NotNil(t, err, "--all and --limit flags cannot be used at the same time")

}

// Validates whether the api policy list is complete
func ValidateAPIPoliciesListWithAllFlag(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listAPIPoliciesWithAllFlag(t, args)

	var policySpecDataList []apim.PolicySpecData

	err := json.Unmarshal([]byte(output), &policySpecDataList)

	if err != nil {
		t.Fatal(err)
	}

	args.SrcAPIM.Login(args.CtlUser.Username, args.CtlUser.Password)

	apiPolicyResponse := args.SrcAPIM.GetAPIPolicies(t, DefaultPolicyListOffsetSize, DefaultPolicyListSize)

	validateListAPIPoliciesEqualWithAllFlag(t, policySpecDataList, apiPolicyResponse)
}

// Validates whether the api policy deletion is complete
func ValidateAPIPolicyDelete(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	policyName := fmt.Sprintf("%v", args.Policy[PolicyNameKey])

	_, err := deleteAPIPolicy(t, policyName, args)

	assert.Nil(t, err, "Error while deleting the API Policy")

}

func deleteAPIPolicy(t *testing.T, name string, args *PolicyImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "policy", "api", "-e", args.SrcAPIM.EnvName, "-n", name, "-k", "--verbose")

	return output, err
}

// get the API policy list apictl output
func listAPIPolicies(t *testing.T, jsonArray bool, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error
	if jsonArray {
		output, err = base.Execute(t, "get", "policies", "api", "-e", args.SrcAPIM.EnvName, "--format", "jsonArray", "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "get", "policies", "api", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	}
	return output, err
}

// get the API policy list apictl output
func listAPIPoliciesWithLimit(t *testing.T, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error
	limit := TestAPIPolicyLimit

	output, err = base.Execute(t, "get", "policies", "api", "-e", args.SrcAPIM.EnvName, "-l", limit, "--format", "jsonArray", "-k", "--verbose")

	return output, err
}

// get the API policy list apictl output
func listAPIPoliciesWithNoLimit(t *testing.T, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error

	output, err = base.Execute(t, "get", "policies", "api", "-e", args.SrcAPIM.EnvName, "--format", "jsonArray", "-k", "--verbose")

	return output, err
}

func listAPIPoliciesWithLimitAndAllFlags(t *testing.T, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error
	limit := TestAPIPolicyLimit

	output, err = base.Execute(t, "get", "policies", "api", "-e", args.SrcAPIM.EnvName, "--all", "--limit", limit, "--format", "jsonArray", "-k", "--verbose")

	return output, err
}

// get the API policy list apictl output
func listAPIPoliciesWithAllFlag(t *testing.T, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error

	output, err = base.Execute(t, "get", "policies", "api", "-e", args.SrcAPIM.EnvName, "--all", "--format", "jsonArray", "-k", "--verbose")

	return output, err
}

// Checks whether apictl output contains all the available api policy UUIDs
func validateListAPIPoliciesEqual(t *testing.T, apiPoliciesListOutput string, apiPoliciesList *apim.APIPoliciesList) {
	unmatchedCount := apiPoliciesList.Count

	for _, policy := range apiPoliciesList.List {
		// If the output string contains the same Policy ID, then decrement the count
		assert.Truef(t, strings.Contains(apiPoliciesListOutput, policy.Id), "APIPoliciesListFromCtl: "+apiPoliciesListOutput+
			" , does not contain policy.id: "+policy.Id)
		unmatchedCount--
	}
	// Count == 0 means that all the policies from apiPoliciesList were in apiPoliciesListOutput
	assert.Equal(t, 0, unmatchedCount, "API policies lists are not equal")
}

// Checks whether apictl output contains all the available api policy UUIDs
func validateListAPIPoliciesEqualWithLimit(t *testing.T, apiPoliciesListOutput []apim.PolicySpecData, apiPoliciesList *apim.APIPoliciesList, limit int) {
	assert.LessOrEqual(t, len(apiPoliciesListOutput), limit, "API Policy list output size is not less than or equivalent with Limit")
	assert.Equal(t, apiPoliciesList.Count, len(apiPoliciesListOutput), "API policies list sizes are not equal")

	for i, policy := range apiPoliciesList.List {
		assert.Equal(t, apiPoliciesListOutput[i].Id, policy.Id, "API Policies are not equal")
	}
}

// Checks whether apictl output has all policies
func validateListAPIPoliciesEqualWithAllFlag(t *testing.T, apiPoliciesListOutput []apim.PolicySpecData, apiPoliciesList *apim.APIPoliciesList) {
	assert.Equal(t, apiPoliciesList.Count, len(apiPoliciesListOutput), "API policies list sizes are not equal")

	for i, policy := range apiPoliciesList.List {
		assert.Equal(t, apiPoliciesListOutput[i].Id, policy.Id, "API Policies are not equal")
	}
}

func readAPIPolicyDefinition(t *testing.T, path string) apim.PolicySpecData {

	// Read the file in the path
	sampleData, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// Extract the content to a structure
	sampleContent := apim.APIPolicyFile{}
	err = yaml.Unmarshal(sampleData, &sampleContent)
	if err != nil {
		t.Fatal(err)
	}

	policyContent, err := apiPolicyDataStructToMap(sampleContent.Data)
	policyContent.Type = sampleContent.Type
	policyContent.Version = DefaultAPIPolicyVersion

	if err != nil {
		t.Fatal(err)
	}

	return policyContent
}

// Converts API policy struct to map
func apiPolicyDataStructToMap(policy interface{}) (apim.PolicySpecData, error) {
	var policyData apim.PolicySpecData
	marshalled, _ := json.Marshal(policy)
	err := json.Unmarshal(marshalled, &policyData)
	return policyData, err
}
