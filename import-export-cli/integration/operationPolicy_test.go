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

package integration

import (
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

// Export an API Policy from one environment and import to another environment
func TestExportImportAPIPolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			// Add two versions of the same policy to a map
			operationPoliciesMap := map[int]testutils.PolicySpecFile{
				1: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition1Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy1Path,
				},
				2: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition2Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy2Path,
				},
			}

			// Export and import two versions of the same policy
			for _, policy := range operationPoliciesMap {
				newPolicy := testutils.AddNewAPIPolicy(t, dev, user.ApiCreator.Username, user.ApiCreator.Password,
					policy, true)
				operationPolicy, _ := testutils.APIPolicyStructToMap(newPolicy)

				args := &testutils.PolicyImportExportTestArgs{
					CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
					Policy:   operationPolicy,
					SrcAPIM:  dev,
					DestAPIM: prod,
				}
				testutils.ValidateAPIPolicyExportImport(t, args)
			}
		})
	}
}

// Export an API Policy that supports both Synapse and Choreo Connect gateways from one environment
// and import to another environment
func TestExportImportAPIPolicyWithSynapseChoreoConnectTypes(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			// Add a policy that supports both Synapse and Choreo Connect gateways
			operationPolicyFiles := testutils.PolicySpecFile{
				Definition:        testutils.TestSynapseChoreoConnectPolicyDefinitionPath,
				SynapsePolicyFile: testutils.TestSynapseChoreoConnectPolicyPathForSynapseType,
				CcPolicyFile:      testutils.TestSynapseChoreoConnectPolicyPathForChoreoConnectType,
			}

			// Export and import the same policy
			newPolicy := testutils.AddNewAPIPolicy(t, dev, user.ApiCreator.Username, user.ApiCreator.Password,
				operationPolicyFiles, true)
			operationPolicy, _ := testutils.APIPolicyStructToMap(newPolicy)

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   operationPolicy,
				SrcAPIM:  dev,
				DestAPIM: prod,
			}
			testutils.ValidateAPIPolicyExportImport(t, args)

		})
	}
}

// Import an API Policy with the directory path
func TestImportAPIPolicyWithDirectoryPath(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			pathToPolicyDir := testutils.CustomAddLogMessage

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:        testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:        dev,
				DestAPIM:       prod,
				ImportFilePath: pathToPolicyDir,
			}
			testutils.ValidateAPIPolicyImportWithDirectoryPath(t, args)
		})
	}
}

// Import an API Policy with inconsistent policy file names
func TestImportAPIPolicyWithInconsistentFileNames(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			pathToPolicyDir := testutils.DevSampleCaseOperationPolicyArtifactsWithInconsistentFileNames

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:        testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:        dev,
				DestAPIM:       prod,
				ImportFilePath: pathToPolicyDir,
			}
			testutils.ValidateAPIPolicyImportWithInconsistentFileNames(t, args)
		})
	}
}

// Export an API Policy from one environment and import to another environment with JSON format
func TestExportImportAPIPolicyWithFormatFlag(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			// Add two versions of the same policy to a map
			operationPoliciesMap := map[int]testutils.PolicySpecFile{
				1: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition1Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy1Path,
				},
				2: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition2Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy2Path,
				},
			}

			// Export and import two versions of the same policy
			for _, policy := range operationPoliciesMap {
				newPolicy := testutils.AddNewAPIPolicy(t, dev, user.ApiCreator.Username, user.ApiCreator.Password,
					policy, true)
				operationPolicy, _ := testutils.APIPolicyStructToMap(newPolicy)

				args := &testutils.PolicyImportExportTestArgs{
					CtlUser:      testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
					Policy:       operationPolicy,
					SrcAPIM:      dev,
					DestAPIM:     prod,
					ExportFormat: "JSON",
				}
				testutils.ValidateAPIPolicyExportImportWithFormatFlag(t, args)
			}
		})
	}
}

// Get API Policy List APICTL output and check whether all policies are included
func TestGetAPIPoliciesList(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateAPIPoliciesList(t, false, args)
		})
	}
}

// Get API Policy List APICTL output in JsonArray format and check whether all policies are included
func TestGetAPIPoliciesListWithJsonArrayFormat(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateAPIPoliciesList(t, true, args)
		})
	}
}

// Get API Policy List APICTL output in JsonArray format and check whether all policies are included under limit
func TestGetAPIPoliciesListWithLimit(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateAPIPoliciesListWithLimit(t, args)
		})
	}
}

// Get API Policy List APICTL output in JsonArray format and check whether all policies are included
func TestGetAPIPoliciesListWithAllFlag(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateAPIPoliciesListWithAllFlag(t, args)
		})
	}
}

// Get API Policy List APICTL output in JsonArray format and check whether all policies are included with default limit
func TestGetAPIPoliciesListWithDefaultLimit(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateAPIPoliciesListWithDefaultLimit(t, args)
		})
	}
}

// Get API Policy List with both the limit and all flags
func TestGetAPIPoliciesListWithAllAndLimitFlags(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateAPIPoliciesListWithLimitAndAllFlagsReturnError(t, args)
		})
	}
}

// Delete API Policy by comparing the status code.
func TestAPIPolicyDelete(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			// Add two versions of the same policy to a map
			operationPoliciesMap := map[int]testutils.PolicySpecFile{
				1: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition1Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy1Path,
				},
				2: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition2Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy2Path,
				},
			}

			// Delete two versions of the same policy
			for _, policy := range operationPoliciesMap {
				newPolicy := testutils.AddNewAPIPolicy(t, dev, user.ApiCreator.Username, user.ApiCreator.Password,
					policy, false)
				operationPolicy, _ := testutils.APIPolicyStructToMap(newPolicy)

				args := &testutils.PolicyImportExportTestArgs{
					CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
					SrcAPIM: dev,
					Policy:  operationPolicy,
				}
				testutils.ValidateAPIPolicyDelete(t, args)
			}
		})
	}
}

// Import a malformed API Policy to an environment by comparing the status code.
func TestMalformedAPIPolicyExportImport(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:  dev,
				DestAPIM: prod,
			}
			testutils.ValidateMalformedAPIPolicyExportImport(t, args)
		})
	}
}

// Exprt API Policy from one environment and import twice to another environment by comparing the status code.
func TestAPIPolicyImportFailureWhenPolicyExisted(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			// Add two versions of the same policy to a map
			operationPoliciesMap := map[int]testutils.PolicySpecFile{
				1: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition1Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy1Path,
				},
				2: {
					Definition:        testutils.DevSampleCaseOperationPolicyDefinition2Path,
					SynapsePolicyFile: testutils.DevSampleCaseOperationPolicy2Path,
				},
			}

			// Export and import two versions of the same policy
			for _, policy := range operationPoliciesMap {
				newPolicy := testutils.AddNewAPIPolicy(t, dev, user.ApiCreator.Username, user.ApiCreator.Password,
					policy, true)
				operationPolicy, _ := testutils.APIPolicyStructToMap(newPolicy)

				args := &testutils.PolicyImportExportTestArgs{
					CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
					SrcAPIM:  dev,
					DestAPIM: prod,
					Policy:   operationPolicy,
				}
				testutils.ValidateAPIPolicyExportImportFailureWhenPolicyExisted(t, args)
			}
		})
	}
}
