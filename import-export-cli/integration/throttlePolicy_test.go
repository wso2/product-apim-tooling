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
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"testing"
)

// Export an Application Throttling Policy from one environment and import to another environment
func TestExportImportApplicationThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.ApplicationThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, apim.ApplicationThrottlePolicyType)
		})
	}
}

// Export a Subscription Throttling Policy from one environment and import to another environment
func TestExportImportSubscriptionThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.SubscriptionThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, apim.SubscriptionThrottlePolicyType)
		})
	}
}

// Export an Advanced Throttling Policy from one environment and import to another environment
func TestExportImportAdvancedThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				SrcAPIM:  dev,
				DestAPIM: prod,
				Update:   false,
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Export a Custom Throttling Policy from one environment and import to another environment as super tenant admin
func TestExportImportCustomThrottlePolicyAdminSuperTenantUser(t *testing.T) {
	//Custom throttling polices is accessible only by AdminSuperTenant
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	dev := GetDevClient()
	prod := GetProdClient()

	newPolicy := testutils.AddNewThrottlePolicy(t, dev, adminUsername, adminPassword, apim.CustomThrottlePolicyType)
	throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
	args := &testutils.ThrottlePolicyImportExportTestArgs{
		Admin:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:  testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Policy:   throttlePolicy,
		SrcAPIM:  dev,
		DestAPIM: prod,
		Update:   false,
	}
	testutils.ValidateThrottlePolicyExportImport(t, args, apim.CustomThrottlePolicyType)
}

// Import an already existing Application Throttling Policy to the destination env with update
func TestImportUpdateApplicationThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, prod, user.Admin.Username, user.Admin.Password, apim.ApplicationThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				DestAPIM: prod,
				Update:   true,
			}
			testutils.ValidateThrottlePolicyImportUpdate(t, args, apim.ApplicationThrottlePolicyType)
		})
	}
}

// Import an already existing Subscription Throttling Policy to the destination env with update
func TestImportUpdateSubscriptionThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, prod, user.Admin.Username, user.Admin.Password, apim.SubscriptionThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				DestAPIM: prod,
				Update:   true,
			}
			testutils.ValidateThrottlePolicyImportUpdate(t, args, apim.SubscriptionThrottlePolicyType)
		})
	}
}

// Import an already existing Advanced Throttling Policy to the destination env with update
func TestImportUpdateAdvancedThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, prod, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				DestAPIM: prod,
				Update:   true,
			}
			testutils.ValidateThrottlePolicyImportUpdate(t, args, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Import an already existing Custom Throttling Policy to the destination env as super tenant admin with update
func TestImportUpdateCustomThrottlePolicyAdminSuperTenantUser(t *testing.T) {
	//Custom throttling polices is accessible only by AdminSuperTenant
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	prod := GetProdClient()

	newPolicy := testutils.AddNewThrottlePolicy(t, prod, adminUsername, adminPassword, apim.CustomThrottlePolicyType)
	throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
	args := &testutils.ThrottlePolicyImportExportTestArgs{
		Admin:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:  testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Policy:   throttlePolicy,
		DestAPIM: prod,
		Update:   true,
	}
	testutils.ValidateThrottlePolicyImportUpdate(t, args, apim.CustomThrottlePolicyType)
}

// Import an already existing Subscription Throttling Policy to the destination env without update
func TestImportUpdateConflictSubscriptionThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, prod, user.Admin.Username, user.Admin.Password, apim.SubscriptionThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				DestAPIM: prod,
				Update:   false,
			}
			testutils.ValidateThrottlePolicyImportUpdateConflict(t, args, apim.SubscriptionThrottlePolicyType)
		})
	}
}

// Import an already existing Application Throttling Policy to the destination env without update
func TestImportUpdateConflictApplicationThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, prod, user.Admin.Username, user.Admin.Password, apim.ApplicationThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				DestAPIM: prod,
				Update:   false,
			}
			testutils.ValidateThrottlePolicyImportUpdateConflict(t, args, apim.ApplicationThrottlePolicyType)
		})
	}
}

// Import an already existing Advanced Throttling Policy to the destination env without update
func TestImportUpdateConflictAdvancedThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, prod, user.Admin.Username, user.Admin.Password, apim.AdvancedThrottlePolicyType)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				Policy:   throttlePolicy,
				DestAPIM: prod,
				Update:   false,
			}
			testutils.ValidateThrottlePolicyImportUpdateConflict(t, args, apim.AdvancedThrottlePolicyType)
		})
	}
}

// Import an already existing Custom Throttling Policy to the destination env as super tenant admin without update
func TestImportUpdateConflictCustomThrottlePolicyAdminSuperTenantUser(t *testing.T) {
	//Custom throttling polices is accessible only by AdminSuperTenant
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	prod := GetProdClient()

	newPolicy := testutils.AddNewThrottlePolicy(t, prod, adminUsername, adminPassword, apim.CustomThrottlePolicyType)
	throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
	args := &testutils.ThrottlePolicyImportExportTestArgs{
		Admin:    testutils.Credentials{Username: adminUsername, Password: adminPassword},
		CtlUser:  testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Policy:   throttlePolicy,
		DestAPIM: prod,
		Update:   false,
	}
	testutils.ValidateThrottlePolicyImportUpdateConflict(t, args, apim.CustomThrottlePolicyType)
}

// Get  Throttle Policy List APICTL output and check whether all policies are included
func TestGetThrottlePoliciesList(t *testing.T) {
	//devops users don't have access to view throttling policies
	for _, user := range testCaseUsers[:2] {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.ThrottlePolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateThrottlePoliciesList(t, false, args)
		})
	}
}

// Get  Throttle Policy List APICTL output in JsonPretty format and check whether all policies are included
func TestGetThrottlePoliciesListWithJsonPrettyFormat(t *testing.T) {
	//devops users don't have access to view throttling policies
	for _, user := range testCaseUsers[:2] {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()

			args := &testutils.ThrottlePolicyImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}
			testutils.ValidateThrottlePoliciesList(t, true, args)
		})
	}
}

// Import failure with a corrupted file
func TestImportInvalidThrottlePolicyFile(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			prod := GetProdClient()

			args := &testutils.ThrottlePolicyImportExportTestArgs{
				Admin:    testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				CtlUser:  testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				DestAPIM: prod,
				Update:   true,
			}
			testutils.ValidateThrottlePolicyImportFailureWithCorruptedFile(t, args)
		})
	}
}
