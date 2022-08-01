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

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

const (
	applicationThrottlePolicyFlag  = "app"
	subscriptionThrottlePolicyFlag = "sub"
	advancedThrottlePolicyFlag     = "advanced"
	customThrottlePolicyFlag       = "custom"
)

// Export an Application Throttling Policy from one environment and import to another environment
func TestExportImportApplicationThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.ApplicationThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     applicationThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, adminUsername, adminPassword, apim.ApplicationThrottlePolicyType)
		})
	}
}

// Export a Subscription Throttling Policy from one environment and import to another environment
func TestExportImportSubscriptionThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.SubscriptionThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     subscriptionThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, adminUsername, adminPassword, apim.SubscriptionThrottlePolicyType)
		})
	}
}

// Export an Advanced Throttling Policy from one environment and import to another environment
func TestExportImportAdvancedThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     advancedThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, adminUsername, adminPassword, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Export a Custom Throttling Policy from one environment and import to another environment as super tenant admin and devops admin
func TestExportImportCustomThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		if isTenantUser(user.CtlUser.Username, TENANT1) {
			continue
		}
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.CustomThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     customThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			testutils.ValidateThrottlePolicyExportImport(t, args, adminUsername, adminPassword, apim.CustomThrottlePolicyType)
		})
	}
}

// Import an already existing Application Throttling Policy to the destination env with update
func TestImportUpdateApplicationThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.ApplicationThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     applicationThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyImportUpdate(t, args, adminUsername, adminPassword, apim.ApplicationThrottlePolicyType)
		})
	}
}

// Import an already existing Subscription Throttling Policy to the destination env with update
func TestImportUpdateSubscriptionThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.SubscriptionThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     subscriptionThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyImportUpdate(t, args, adminUsername, adminPassword, apim.SubscriptionThrottlePolicyType)
		})
	}
}

// Import an already existing Advanced Throttling Policy to the destination env with update
func TestImportUpdateAdvancedThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     advancedThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyImportUpdate(t, args, adminUsername, adminPassword, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Import an already existing Custom Throttling Policy to the destination env with update
func TestImportUpdateCustomThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		if isTenantUser(user.CtlUser.Username, TENANT1) {
			continue
		}
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.CustomThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     customThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			testutils.ValidateThrottlePolicyImportUpdate(t, args, adminUsername, adminPassword, apim.CustomThrottlePolicyType)
		})
	}
}

// Import an already existing Subscription Throttling Policy to the destination env without update
func TestSubscriptionThrottlePolicyImportFailureWhenPolicyExisted(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.SubscriptionThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     subscriptionThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyImportFailureWhenPolicyExisted(t, args, adminUsername, adminPassword, apim.SubscriptionThrottlePolicyType)
		})
	}
}

// Import an already existing Application Throttling Policy to the destination env without update
func TestApplicationThrottlePolicyImportFailureWhenPolicyExisted(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.ApplicationThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     applicationThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyImportFailureWhenPolicyExisted(t, args, adminUsername, adminPassword, apim.ApplicationThrottlePolicyType)
		})
	}
}

// Import an already existing Advanced Throttling Policy to the destination env without update
func TestAdvancedThrottlePolicyImportFailureWhenPolicyExisted(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     advancedThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyImportFailureWhenPolicyExisted(t, args, adminUsername, adminPassword, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Import an already existing Custom Throttling Policy to the destination env without update
func TestCustomThrottlePolicyImportFailureWhenPolicyExisted(t *testing.T) {
	for _, user := range testCaseUsers {
		// Custom throttling policy importation
		if isTenantUser(user.CtlUser.Username, TENANT1) {
			continue
		}
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.CustomThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				Type:     customThrottlePolicyFlag,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			testutils.ValidateThrottlePolicyImportFailureWhenPolicyExisted(t, args, adminUsername, adminPassword, apim.CustomThrottlePolicyType)
		})
	}
}

// Get  Throttle Policy List APICTL output and check whether all policies are included
func TestGetThrottlePoliciesList(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePoliciesList(t, true, adminUsername, adminPassword, args)
		})
	}
}

// Get  Throttle Policy List APICTL output in JsonPretty format and check whether all policies are included
func TestGetThrottlePoliciesListWithJsonPrettyFormat(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePoliciesList(t, true, adminUsername, adminPassword, args)
		})
	}
}

// Import failure with a corrupted file
func TestImportInvalidThrottlePolicyFile(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			prod := GetProdClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				DestAPIM: prod,
				Update:   false,
			}
			testutils.ValidateThrottlePolicyImportFailureWithCorruptedFile(t, args)
		})
	}
}

// Export an invalid throttling policy
func TestExportInvalidThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			args := &testutils.PolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateThrottlePolicyExportFailure(t, args)
		})
	}
}

// Export a Throttling Policy from one environment and import to another environment without type flag
func TestExportImportThrottlePolicyWithoutTypeFlag(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			adminUsername := superAdminUser
			adminPassword := superAdminPassword
			if isTenantUser(args.CtlUser.Username, TENANT1) {
				adminUsername = adminUsername + "@" + TENANT1
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, adminUsername, adminPassword, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Delete Throttling Policy by comparing the status code.
func TestThrottlingPoliciesDelete(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}

			testutils.ValidateThrottlingPolicyDelete(t, args, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Delete Non Existing Throttling Policy by comparing the status code.
func TestNonExistingThrottlingPolicyDelete(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			throttlePolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			args := &testutils.PolicyImportExportTestArgs{
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}

			testutils.ValidateThrottlingPolicyDelete(t, args, apim.AdvancedThrottlePolicyType)
		})
	}
}
