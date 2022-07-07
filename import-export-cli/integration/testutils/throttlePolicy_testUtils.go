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
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	policyIDKey   = "policyId"
	policyNameKey = "policyName"
	policyIDReset = "-1"

	applicationPolicySubType  = "application policy"
	advancedPolicySubType     = "advanced policy"
	customPolicySubType       = "custom rule"
	subscriptionPolicySubType = "subscription policy"
)

// ValidateThrottlePolicyExportImport : Validates Exporting Throttling Policy from source env and Importing to destination env
func ValidateThrottlePolicyExportImport(t *testing.T, args *PolicyImportExportTestArgs, policyType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[policyNameKey])

	exportedOutput, err := exportThrottlePolicy(t, policyName, args.SrcAPIM.GetEnvName())
	assert.Nil(t, err, "Error while exporting the Throttling Policy")

	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))
	args.SrcAPIM.DeleteThrottlePolicy(fmt.Sprintf("%v", args.Policy[policyIDKey]), policyType)

	// Import Throttling Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	importedOutput, err := importThrottlePolicy(t, args)
	assert.Nil(t, err, "Error while importing the Throttling Policy")
	assert.Contains(t, importedOutput, "Successfully imported")
	// Give time for newly imported Throttling Policy to get indexed
	base.WaitForIndexing()

	// Get Throttle Policy from env 2
	importedPolicy, _ := getThrottlingPolicyByName(t, args, policyName, policyType)
	// Validate env 1 and env 2 policy is equal
	validatePoliciesEqual(t, args, importedPolicy)
	removeExportedThrottlingPolicyFile(t, args.ImportFilePath)
}

// ValidateThrottlePolicyImportUpdate : Validates importing existing Throttling Policy
func ValidateThrottlePolicyImportUpdate(t *testing.T, args *PolicyImportExportTestArgs, policyType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[policyNameKey])

	exportedOutput, _ := exportThrottlePolicy(t, policyName, args.DestAPIM.GetEnvName())
	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	// Import Throttling Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	output, err := importThrottlePolicy(t, args)
	assert.Contains(t, output, "Successfully updated")
	assert.Nil(t, err, "Error while importing the Throttling Policy")
	// Give time for newly imported Throttling Policy to get indexed
	base.WaitForIndexing()

	// Get Throttle Policy from env 2
	importedPolicy, _ := getThrottlingPolicyByName(t, args, policyName, policyType)
	// Validate env 1 and env 2 policy is equal
	validatePoliciesEqual(t, args, importedPolicy)
	removeExportedThrottlingPolicyFile(t, args.ImportFilePath)
}

// ValidateThrottlePolicyImportUpdateConflict : Validates importing existing Throttling Policy to create conflict
func ValidateThrottlePolicyImportUpdateConflict(t *testing.T, args *PolicyImportExportTestArgs, policyType string) {
	const conflictStatusCode = "409"
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[policyNameKey])

	exportedOutput, _ := exportThrottlePolicy(t, policyName, args.DestAPIM.GetEnvName())
	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	// Import Throttling Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	output, err := importThrottlePolicy(t, args)
	assert.Error(t, err, "Importation conflict expected")
	assert.Contains(t, output, conflictStatusCode, "Unexpected error code")
	// Give time for newly imported Throttling Policy to get indexed
	base.WaitForIndexing()

	// Get Throttle Policy from env 2
	importedPolicy, _ := getThrottlingPolicyByName(t, args, policyName, policyType)
	// Validate env 1 and env 2 policy is equal
	validatePoliciesEqual(t, args, importedPolicy)
	removeExportedThrottlingPolicyFile(t, args.ImportFilePath)
}

func ValidateThrottlePolicyImportFailureWithCorruptedFile(t *testing.T, args *PolicyImportExportTestArgs) {
	const internalServerError = "500"
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	path, _, _ := createExportedThrottlePolicyFile(t, args.DestAPIM, apim.AdvancedThrottlePolicyType, true)
	args.ImportFilePath = path
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	//// Import Throttling Policy to dest
	output, err := importThrottlePolicy(t, args)
	assert.Error(t, err, "Importation failure expected")
	assert.Contains(t, output, internalServerError, "Unexpected error code")
	removeExportedThrottlingPolicyFile(t, args.ImportFilePath)
}

func ValidateThrottlePolicyExportFailure(t *testing.T, args *PolicyImportExportTestArgs) {
	const policyString = "Policy"
	t.Helper()
	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := base.GenerateRandomString() + policyString

	output, _ := exportThrottlePolicy(t, policyName, args.SrcAPIM.GetEnvName())
	assert.Contains(t, output, "Error exporting Throttling Policies", "Exportation error expected")
}

func createExportedThrottlePolicyFile(t *testing.T, client *apim.Client, policyType string, corrupted bool) (string, interface{}, error) {

	t.Helper()
	var exportedPolicy utils.ExportThrottlePolicy
	exportedPolicy.Type = "throttling policy"

	if !corrupted {
		switch policyType {
		case apim.SubscriptionThrottlePolicyType:
			exportedPolicy.Subtype = subscriptionPolicySubType
		case apim.AdvancedThrottlePolicyType:
			exportedPolicy.Subtype = advancedPolicySubType
		case apim.CustomThrottlePolicyType:
			exportedPolicy.Subtype = customPolicySubType
		case apim.ApplicationThrottlePolicyType:
			exportedPolicy.Subtype = applicationPolicySubType
		}
	}

	exportedPolicy.Version = "v4.1.0"
	policyData := client.GenerateSampleThrottlePolicyData(policyType)
	policyMap, _ := PolicyStructToMap(policyData)
	var yamlMap yaml.MapSlice
	yamlBytes, err := yaml.Marshal(policyMap)
	err = yaml.Unmarshal(yamlBytes, &yamlMap)
	if err != nil {
		return "", policyData, err
	}
	exportedPolicy.Data = yamlMap
	policyMarshaledData, _ := yaml.Marshal(exportedPolicy)

	tmpDir, err := ioutil.TempDir("", "apim")
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", policyData, err
	}
	tempFilePath := tmpDir + "/temp.yaml"
	err = ioutil.WriteFile(tempFilePath, policyMarshaledData, os.ModePerm)
	if err != nil {
	}
	return tempFilePath, policyData, err
}

// Adds a new Throttling Policy to an env
func AddNewThrottlePolicy(t *testing.T, client *apim.Client, username, password, policyType string) interface{} {
	client.Login(username, password)
	generatedPolicy := client.GenerateSampleThrottlePolicyData(policyType)
	addedPolicy := client.AddThrottlePolicy(t, generatedPolicy, policyType)
	return addedPolicy
}

// Exports Throttling Policy from an env
func exportThrottlePolicy(t *testing.T, name, env string) (string, error) {
	output, err := base.Execute(t, "export", "policy", "rate-limiting", "-n", name, "-e", env, "-k", "--verbose")
	return output, err
}

// Imports Throttling policy to an env
func importThrottlePolicy(t *testing.T, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error
	if args.Update {
		output, err = base.Execute(t, "import", "policy", "rate-limiting", "-e", args.DestAPIM.GetEnvName(), "-f", args.ImportFilePath, "-u")
	} else {
		output, err = base.Execute(t, "import", "policy", "rate-limiting", "-e", args.DestAPIM.GetEnvName(), "-f", args.ImportFilePath)
	}
	return output, err
}

// Retrieve Throttling policy by name
func getThrottlingPolicyByName(t *testing.T, args *PolicyImportExportTestArgs, throttlePolicyName, policyType string) (map[string]interface{}, error) {
	client := args.DestAPIM
	uuid := client.GetThrottlePolicyID(t, args.Admin.Username, args.Admin.Password, throttlePolicyName, policyType)
	policy := client.GetThrottlePolicy(uuid, policyType)
	client.DeleteThrottlePolicy(uuid, policyType)
	return PolicyStructToMap(policy)
}

// Validates whether throttling policies are equal
func validatePoliciesEqual(t *testing.T, args *PolicyImportExportTestArgs, importedPolicy map[string]interface{}) {
	exportedPolicy := args.Policy
	exportedPolicy[policyIDKey] = policyIDReset
	importedPolicy[policyIDKey] = policyIDReset
	assert.Equal(t, exportedPolicy, importedPolicy)
}

// Converts policy struct to map
func PolicyStructToMap(policyStruct interface{}) (map[string]interface{}, error) {
	var policyMap map[string]interface{}
	marshalled, _ := json.Marshal(policyStruct)
	err := json.Unmarshal(marshalled, &policyMap)
	return policyMap, err
}

// Cleanup func for exported throttling policy file
func removeExportedThrottlingPolicyFile(t *testing.T, file string) {
	t.Log("base.RemoveExportedThrottlingPolicyFile() - file path:", file)
	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Validates whether the throttling policy list is complete
func ValidateThrottlePoliciesList(t *testing.T, doJsonPrettyFormatting bool, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listThrottlePolicies(t, doJsonPrettyFormatting, args)

	throttlePoliciesList := args.SrcAPIM.GetThrottlePolicies(t, args.CtlUser.Username, args.CtlUser.Password)

	validateListThrottlePoliciesEqual(t, output, throttlePoliciesList)
}

// get the Throttle policy list apictl output
func listThrottlePolicies(t *testing.T, doJsonPrettyFormatting bool, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error
	if doJsonPrettyFormatting {
		output, err = base.Execute(t, "get", "policies", "rate-limiting", "-e", args.SrcAPIM.EnvName, "--format", "\"{{ jsonPretty . }}\"", "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "get", "policies", "rate-limiting", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	}
	return output, err
}

// Checks whether apictl output contains all the available throttle policy UUIDs
func validateListThrottlePoliciesEqual(t *testing.T, throttlePoliciesListOutput string, throttlePoliciesList *utils.ThrottlingPoliciesDetailsList) {
	unmatchedCount := throttlePoliciesList.Count
	for _, policy := range throttlePoliciesList.List {
		// If the output string contains the same Policy ID, then decrement the count
		assert.Truef(t, strings.Contains(throttlePoliciesListOutput, policy.Uuid), "throttlePoliciesListFromCtl: "+throttlePoliciesListOutput+
			" , does not contain policy.uuid: "+policy.Uuid)
		unmatchedCount--
	}
	// Count == 0 means that all the policies from throttlePoliciesList were in throttlePoliciesListOutput
	assert.Equal(t, 0, unmatchedCount, "Throttle policies lists are not equal")
}
