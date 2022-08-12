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
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
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
func ValidateThrottlePolicyExportImport(t *testing.T, args *PolicyImportExportTestArgs, adminUsername, adminPassword, policyType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[policyNameKey])

	exportedOutput, err := exportThrottlePolicy(t, policyName, args)
	assert.Nil(t, err, "Error while exporting the Throttling Policy")

	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	// Import Throttling Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	importedOutput, err := importThrottlePolicy(t, adminUsername, adminPassword, policyName, policyType, true, args)
	assert.Nil(t, err, "Error while importing the Throttling Policy")
	assert.Contains(t, importedOutput, "Successfully imported")
	// Give time for newly imported Throttling Policy to get indexed
	base.WaitForIndexing()

	// Get Throttle Policy from env 2
	importedPolicy, _ := getThrottlingPolicyByName(t, args, adminUsername, adminPassword, policyName, policyType)
	// Validate env 1 and env 2 policy is equal
	validatePoliciesEquality(t, true, args, importedPolicy)
}

// Validates whether the throttling policy deletion is complete
func ValidateThrottlingPolicyDelete(t *testing.T, args *PolicyImportExportTestArgs, username, password, policyType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	policyName := fmt.Sprintf("%v", args.Policy[policyNameKey])

	_, err := deleteThrottlingPolicy(t, policyName, policyType, args)

	assert.Nil(t, err, "Error while deleting the API Policy")

}

func deleteThrottlingPolicy(t *testing.T, name string, policyType string, args *PolicyImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "policy", "rate-limiting", "-e", args.SrcAPIM.EnvName, "-n", name, "-t", policyType, "-k", "--verbose")
	return output, err
}

// ValidateThrottlePolicyImportUpdate : Validates importing existing Throttling Policy
func ValidateThrottlePolicyImportUpdate(t *testing.T, args *PolicyImportExportTestArgs, adminUsername, adminPassword, policyType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[policyNameKey])

	exportedOutput, _ := exportThrottlePolicy(t, policyName, args)
	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	// Import Throttling Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	importedOutput, err := importThrottlePolicy(t, adminUsername, adminPassword, policyName, policyType, false, args)
	assert.Nil(t, err, "Error while importing the Throttling Policy")
	assert.Contains(t, importedOutput, "Successfully imported")
	// Give time for newly imported Throttling Policy to get indexed
	base.WaitForIndexing()

	changeExportedThrottlePolicyFile(t, args)

	// Import Throttling Policy to env 2 to update
	args.Update = true // to update existing throttling policy
	output, err := importThrottlePolicy(t, adminUsername, adminPassword, policyName, policyType, true, args)
	assert.Contains(t, output, "Successfully updated")
	assert.Nil(t, err, "Error while importing the Throttling Policy")
	// Give time for newly imported Throttling Policy to get indexed
	base.WaitForIndexing()
	// Get Throttle Policy from env 2
	importedPolicy, _ := getThrottlingPolicyByName(t, args, adminUsername, adminPassword, policyName, policyType)
	// Validate  policies are not equal
	validatePoliciesEquality(t, false, args, importedPolicy)
}

// ValidateThrottlePolicyImportFailureWhenPolicyExisted : Validates importing existing Throttling Policy to create conflict
func ValidateThrottlePolicyImportFailureWhenPolicyExisted(t *testing.T, args *PolicyImportExportTestArgs, adminUsername, adminPassword, policyType string) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := fmt.Sprintf("%v", args.Policy[policyNameKey])

	exportedOutput, _ := exportThrottlePolicy(t, policyName, args)
	args.ImportFilePath = base.GetExportedPathFromOutput(exportedOutput)
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	// Import Throttling Policy to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importedOutput, err := importThrottlePolicy(t, adminUsername, adminPassword, policyName, policyType, false, args)
	assert.Nil(t, err, "Error while importing the Throttling Policy")
	assert.Contains(t, importedOutput, "Successfully imported")
	// Give time for newly imported Throttling Policy to get indexed
	base.WaitForIndexing()

	changeExportedThrottlePolicyFile(t, args)
	// Import Throttling Policy to env 2 for update conflict

	output, err := importThrottlePolicy(t, adminUsername, adminPassword, policyName, policyType, true, args)
	assert.Error(t, err, "Importation conflict expected")
	assert.Contains(t, output, "409", "Unexpected error code")
}

func ValidateThrottlePolicyImportFailureWithCorruptedFile(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	path, _, _ := createExportedThrottlePolicyFile(t, args.DestAPIM, apim.AdvancedThrottlePolicyType, true)
	args.ImportFilePath = path
	assert.True(t, base.IsFileAvailable(t, args.ImportFilePath))

	// Import Throttling Policy to dest
	output, err := importThrottlePolicy(t, "", "", "", "", false, args)
	assert.Error(t, err, "Importation failure expected")
	assert.Contains(t, output, "500", "Unexpected error code")
}

func ValidateThrottlePolicyExportFailure(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()
	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Export policy from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	policyName := base.GenerateRandomString() + "Policy"

	output, _ := exportThrottlePolicy(t, policyName, args)
	assert.Contains(t, output, "Error exporting Throttling Policies", "Exportation error expected")
}

func changeExportedThrottlePolicyFile(t *testing.T, args *PolicyImportExportTestArgs) {
	t.Helper()
	var exportedPolicy utils.ExportThrottlePolicy
	exportedFile, err := ioutil.ReadFile(args.ImportFilePath)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(exportedFile, &exportedPolicy)
	if err != nil {
		t.Fatal(err)
	}
	policyData := exportedPolicy.Data
	policyData[2].Value = base.GenerateRandomString()
	policyData[3].Value = base.GenerateRandomString()
	marshaledData, _ := yaml.Marshal(exportedPolicy)

	err = ioutil.WriteFile(args.ImportFilePath, marshaledData, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
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

	if err != nil {
		return "", "", err
	}

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
		t.Fatal(err)
	}

	t.Cleanup(func() {
		removeExportedThrottlingPolicyFile(t, tempFilePath)
	})

	return tempFilePath, policyData, err
}

// Adds a new Throttling Policy to an env
func AddNewThrottlePolicy(t *testing.T, client *apim.Client, adminUsername, adminPassword, policyType string, doClean bool) map[string]interface{} {
	client.Login(adminUsername, adminPassword)
	generatedPolicy := client.GenerateSampleThrottlePolicyData(policyType)
	addedPolicy := client.AddThrottlePolicy(t, generatedPolicy, adminUsername, adminPassword, policyType, doClean)
	return addedPolicy
}

// Exports Throttling Policy from an env
func exportThrottlePolicy(t *testing.T, name string, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error
	if args.Type != "" {
		output, err = base.Execute(t, "export", "policy", "rate-limiting", "-n", name, "-t", args.Type, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "export", "policy", "rate-limiting", "-n", name, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	}

	t.Cleanup(func() {
		removeExportedThrottlingPolicyFile(t, args.ImportFilePath)
	})

	return output, err
}

// Imports Throttling policy to an env
func importThrottlePolicy(t *testing.T, username, password, policyName, policyType string, doClean bool, args *PolicyImportExportTestArgs) (string, error) {
	var output string
	var err error
	if args.Update {
		output, err = base.Execute(t, "import", "policy", "rate-limiting", "-e", args.DestAPIM.GetEnvName(), "-f", args.ImportFilePath, "-u")
	} else {
		output, err = base.Execute(t, "import", "policy", "rate-limiting", "-e", args.DestAPIM.GetEnvName(), "-f", args.ImportFilePath)
	}

	if doClean {
		t.Cleanup(func() {
			args.DestAPIM.Login(username, password)
			args.DestAPIM.DeleteThrottlePolicyByName(t, policyName, policyType, doClean)
		})
	}

	return output, err
}

// Retrieve Throttling policy by name
func getThrottlingPolicyByName(t *testing.T, args *PolicyImportExportTestArgs, adminUsername, adminPassword, throttlePolicyName, policyType string) (map[string]interface{}, error) {
	client := args.DestAPIM
	client.Login(adminUsername, adminPassword)
	uuid := client.GetThrottlePolicyID(t, throttlePolicyName, policyType)
	policy := client.GetThrottlePolicy(uuid, policyType)
	return PolicyStructToMap(policy)
}

// Validates whether throttling policies are equal
func validatePoliciesEquality(t *testing.T, checkEquality bool, args *PolicyImportExportTestArgs, importedPolicy map[string]interface{}) {
	exportedPolicy := args.Policy
	exportedPolicy[policyIDKey] = policyIDReset
	importedPolicy[policyIDKey] = policyIDReset
	if checkEquality {
		assert.Equal(t, exportedPolicy, importedPolicy)
	} else {
		assert.NotEqual(t, exportedPolicy, importedPolicy)
	}
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
func ValidateThrottlePoliciesList(t *testing.T, doJsonPrettyFormatting bool, adminUsername, adminPassword string, args *PolicyImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	base.WaitForIndexing()

	output, _ := listThrottlePolicies(t, doJsonPrettyFormatting, args)

	args.SrcAPIM.Login(adminUsername, adminPassword)
	throttlePoliciesList := args.SrcAPIM.GetThrottlePolicies(t)

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
