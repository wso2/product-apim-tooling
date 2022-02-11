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

// Get log levels of APIs of the carbon.super tenant in an environment as a super admin user
func TestGetAPILogLevelsSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreatorUsername := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api1 := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)
	api2 := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)
	api3 := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Apis:         []*apim.API{api1, api2, api3},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
	}

	testutils.ValidateGetAPILogLevel(t, args)
}

// Get log levels of APIs of the carbon.super tenant in an environment as a non super admin user
func TestGetAPILogLevelsNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreatorUsername := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
	}

	testutils.ValidateGetAPILogLevelError(t, args)
}

// Get log levels of APIs of another tenant in an environment as a super admin user
func TestGetAPILogLevelsAnotherTenantSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreatorUsername := testCaseUsers[1].ApiCreator.Username
	apiCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	api1 := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)
	api2 := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)
	api3 := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Apis:         []*apim.API{api1, api2, api3},
		APIM:         dev,
		TenantDomain: TENANT1,
	}

	testutils.ValidateGetAPILogLevel(t, args)
}

// Get log levels of APIs of the another tenant in an environment as a non super admin user
func TestGetAPILogLevelsAnotherTenantNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreatorUsername := testCaseUsers[1].ApiCreator.Username
	apiCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: TENANT1,
	}

	testutils.ValidateGetAPILogLevelError(t, args)
}

// Get log level of an API of the carbon.super tenant in an environment as a super admin user
func TestGetAPILogLevelSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreatorUsername := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		ApiId:        api.ID,
	}

	testutils.ValidateGetAPILogLevel(t, args)
}

// Get log level of an API of the carbon.super tenant in an environment as a non super admin user
func TestGetAPILogLevelNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreatorUsername := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		ApiId:        api.ID,
	}

	testutils.ValidateGetAPILogLevelError(t, args)
}

// Get log level of an API of another tenant in an environment as a super admin user
func TestGetAPILogLevelAnotherTenantSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreatorUsername := testCaseUsers[1].ApiCreator.Username
	apiCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: TENANT1,
		ApiId:        api.ID,
	}

	testutils.ValidateGetAPILogLevel(t, args)
}

// Get log level of an API of the another tenant in an environment as a non super admin user
func TestGetAPILogLevelAnotherTenantNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreatorUsername := testCaseUsers[1].ApiCreator.Username
	apiCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: TENANT1,
		ApiId:        api.ID,
	}

	testutils.ValidateGetAPILogLevelError(t, args)
}

// Set log level of an API of the carbon.super tenant in an environment as a super admin user
func TestSetAPILogLevelSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreatorUsername := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		ApiId:        api.ID,
		LogLevel:     "FULL",
	}

	testutils.ValidateSetAPILogLevel(t, args)
}

// Set log level of an API of the carbon.super tenant in an environment as a non super admin user
func TestSetAPILogLevelNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreatorUsername := creator.UserName
	apiCreatorPassword := creator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		ApiId:        api.ID,
		LogLevel:     "STANDARD",
	}

	testutils.ValidateSetAPILogLevelError(t, args)
}

// Set log level of an API of another tenant in an environment as a super admin user
func TestSetAPILogLevelAnotherTenantSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	apiCreatorUsername := testCaseUsers[1].ApiCreator.Username
	apiCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: TENANT1,
		ApiId:        api.ID,
		LogLevel:     "BASIC",
	}

	testutils.ValidateSetAPILogLevel(t, args)
}

// Set log level of an API of the another tenant in an environment as a non super admin user
func TestSetAPILogLevelAnotherTenantNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	apiCreatorUsername := testCaseUsers[1].ApiCreator.Username
	apiCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	api := testutils.AddAPI(t, dev, apiCreatorUsername, apiCreatorPassword)

	args := &testutils.ApiLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		Apis:         []*apim.API{api},
		APIM:         dev,
		TenantDomain: TENANT1,
		ApiId:        api.ID,
		LogLevel:     "OFF",
	}

	testutils.ValidateSetAPILogLevelError(t, args)
}
