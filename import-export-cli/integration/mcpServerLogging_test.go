/*
*  Copyright (c) 2025 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 LLC. licenses this file to you under the Apache License,
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

// Get log levels of MCP Servers of the carbon.super tenant in an environment as a super admin user
func TestGetMCPServerLogLevelsSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreatorUsername := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer1 := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)
	mcpServer2 := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)
	mcpServer3 := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer1, mcpServer2, mcpServer3},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
	}

	testutils.ValidateGetMCPServerLogLevel(t, args)
}

// Get log levels of MCP Servers of the carbon.super tenant in an environment as a non super admin user
func TestGetMCPServerLogLevelsNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreatorUsername := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
	}

	testutils.ValidateGetMCPServerLogLevelError(t, args)
}

// Get log levels of MCP Servers of another tenant in an environment as a super admin user
func TestGetMCPServerLogLevelsAnotherTenantSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreatorUsername := testCaseUsers[1].ApiCreator.Username
	mcpServerCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	mcpServer1 := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)
	mcpServer2 := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)
	mcpServer3 := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer1, mcpServer2, mcpServer3},
		APIM:         dev,
		TenantDomain: TENANT1,
	}

	testutils.ValidateGetMCPServerLogLevel(t, args)
}

// Get log levels of MCP Servers of another tenant in an environment as a non super admin user
func TestGetMCPServerLogLevelsAnotherTenantNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreatorUsername := testCaseUsers[1].ApiCreator.Username
	mcpServerCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: TENANT1,
	}

	testutils.ValidateGetMCPServerLogLevelError(t, args)
}

// Get log level of an MCP Server of the carbon.super tenant in an environment as a super admin user
func TestGetMCPServerLogLevelSuperAdminUserSingle(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreatorUsername := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		MCPServerId:  mcpServer.ID,
	}

	testutils.ValidateGetMCPServerLogLevel(t, args)
}

// Get log level of an MCP Server of the carbon.super tenant in an environment as a non super admin user
func TestGetMCPServerLogLevelNonSuperAdminUserSingle(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreatorUsername := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		MCPServerId:  mcpServer.ID,
	}

	testutils.ValidateGetMCPServerLogLevelError(t, args)
}

// Get log level of an MCP Server of another tenant in an environment as a super admin user
func TestGetMCPServerLogLevelAnotherTenantSuperAdminUserSingle(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreatorUsername := testCaseUsers[1].ApiCreator.Username
	mcpServerCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: TENANT1,
		MCPServerId:  mcpServer.ID,
	}

	testutils.ValidateGetMCPServerLogLevel(t, args)
}

// Get log level of an MCP Server of another tenant in an environment as a non super admin user
func TestGetMCPServerLogLevelAnotherTenantNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreatorUsername := testCaseUsers[1].ApiCreator.Username
	mcpServerCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: TENANT1,
		MCPServerId:  mcpServer.ID,
	}

	testutils.ValidateGetMCPServerLogLevelError(t, args)
}

// Set log level of an MCP Server of the carbon.super tenant in an environment as a super admin user
func TestSetMCPServerLogLevelSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreatorUsername := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		MCPServerId:  mcpServer.ID,
		LogLevel:     "FULL",
	}

	testutils.ValidateSetMCPServerLogLevel(t, args)
}

// Set log level of an MCP Server of the carbon.super tenant in an environment as a non super admin user
func TestSetMCPServerLogLevelNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreatorUsername := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: DEFAULT_TENANT_DOMAIN,
		MCPServerId:  mcpServer.ID,
		LogLevel:     "STANDARD",
	}

	testutils.ValidateSetMCPServerLogLevelError(t, args)
}

// Set log level of an MCP Server of another tenant in an environment as a super admin user
func TestSetMCPServerLogLevelAnotherTenantSuperAdminUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreatorUsername := testCaseUsers[1].ApiCreator.Username
	mcpServerCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: TENANT1,
		MCPServerId:  mcpServer.ID,
		LogLevel:     "BASIC",
	}

	testutils.ValidateSetMCPServerLogLevel(t, args)
}

// Set log level of an MCP Server of another tenant in an environment as a non super admin user
func TestSetMCPServerLogLevelAnotherTenantNonSuperAdminUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreatorUsername := testCaseUsers[1].ApiCreator.Username
	mcpServerCreatorPassword := testCaseUsers[1].ApiCreator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreatorUsername, mcpServerCreatorPassword)

	args := &testutils.MCPServerLoggingTestArgs{
		CtlUser:      testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServers:   []*apim.MCPServer{mcpServer},
		APIM:         dev,
		TenantDomain: TENANT1,
		MCPServerId:  mcpServer.ID,
		LogLevel:     "OFF",
	}

	testutils.ValidateSetMCPServerLogLevelError(t, args)
}
