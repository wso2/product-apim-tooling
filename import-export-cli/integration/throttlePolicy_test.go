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
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"testing"
)

const (
	ApplicationPolicy  = "application"
	CustomPolicy       = "custom"
	AdvancedPolicy     = "advanced"
	SubscriptionPolicy = "subscription"
)

// Export an Application Throttling Policy from one environment and import to another environment
func TestExportImportApplicationThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			adminUsername := user.CtlUser.Username
			adminPassword := user.CtlUser.Password

			apiCreator := creator.UserName
			apiCreatorPassword := creator.Password

			dev := GetDevClient()
			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, ApplicationPolicy)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
				CtlUser:     testutils.Credentials{Username: adminUsername, Password: adminPassword},
				Policy:      throttlePolicy,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, ApplicationPolicy)
		})
	}
}

// Export a Subscription Throttling Policy from one environment and import to another environment
func TestExportImportSubscriptionThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			adminUsername := user.CtlUser.Username
			adminPassword := user.CtlUser.Password

			apiCreator := creator.UserName
			apiCreatorPassword := creator.Password

			dev := GetDevClient()
			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, SubscriptionPolicy)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
				CtlUser:     testutils.Credentials{Username: adminUsername, Password: adminPassword},
				Policy:      throttlePolicy,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, SubscriptionPolicy)
		})
	}
}

// Export an Advanced Throttling Policy from one environment and import to another environment
func TestExportImportAdvancedThrottlePolicy(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			adminUsername := user.CtlUser.Username
			adminPassword := user.CtlUser.Password

			apiCreator := creator.UserName
			apiCreatorPassword := creator.Password

			dev := GetDevClient()
			prod := GetProdClient()

			newPolicy := testutils.AddNewThrottlePolicy(t, dev, user.Admin.Username, user.Admin.Password, AdvancedPolicy)
			throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
			args := &testutils.ThrottlePolicyImportExportTestArgs{
				ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
				CtlUser:     testutils.Credentials{Username: adminUsername, Password: adminPassword},
				Policy:      throttlePolicy,
				SrcAPIM:     dev,
				DestAPIM:    prod,
			}
			testutils.ValidateThrottlePolicyExportImport(t, args, AdvancedPolicy)
		})
	}
}

// Export a Custom Throttling Policy from one environment and import to another environment as super tenant admin
func TestExportImportCustomThrottlePolicyAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreator := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	newPolicy := testutils.AddNewThrottlePolicy(t, dev, adminUsername, adminPassword, CustomPolicy)
	throttlePolicy, _ := testutils.ThrottlePolicyStructToMap(newPolicy)
	args := &testutils.ThrottlePolicyImportExportTestArgs{
		ApiProvider: testutils.Credentials{Username: apiCreator, Password: apiCreatorPassword},
		CtlUser:     testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Policy:      throttlePolicy,
		SrcAPIM:     dev,
		DestAPIM:    prod,
	}
	testutils.ValidateThrottlePolicyExportImport(t, args, CustomPolicy)
}
