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

const numberOfMCPServers = 5 // Number of MCP Servers to be added in a loop

// Export a MCP Server from one environment and check the structure of the DTO whether it is similar to what is being
// maintained by APICTL
func TestExportMCPServerCompareStruct(t *testing.T) {
	mcpServerPublisher := publisher.UserName
	mcpServerPublisherPassword := publisher.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: mcpServerPublisher, Password: mcpServerPublisherPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateExportedMCPServerStructure(t, args)
}

// Export a MCP Server from one environment as non-admin super tenant user
func TestExportMCPServerNonAdminSuperTenantUser(t *testing.T) {
	mcpServerPublisher := publisher.UserName
	mcpServerPublisherPassword := publisher.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: mcpServerPublisher, Password: mcpServerPublisherPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Export a MCP Server from one environment and import to another environment as super tenant admin user
func TestExportImportMCPServerAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportImport(t, args)
}

// Export a MCP Server from one environment and import to another environment as super tenant user with
// Internal/devops role
func TestExportImportMCPServerDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportImport(t, args)
}

// Export a MCP Server from one environment as super tenant publisher user
func TestExportMCPServerSuperTenantPublisherUser(t *testing.T) {
	mcpServerPublisher := publisher.UserName
	mcpServerPublisherPassword := publisher.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: mcpServerPublisher, Password: mcpServerPublisherPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Export a MCP Server from one environment as tenant publisher user
func TestExportMCPServerTenantPublisherUser(t *testing.T) {
	tenantPublisher := publisher.UserName + "@" + TENANT1
	tenantPublisherPassword := publisher.Password

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: tenantCreator, Password: tenantCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantPublisher, Password: tenantPublisherPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Export a MCP Server from one environment as super tenant subscriber user
func TestExportMCPServerSuperTenantSubscriberUser(t *testing.T) {
	mcpServerSubscriber := subscriber.UserName
	mcpServerSubscriberPassword := subscriber.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: mcpServerSubscriber, Password: mcpServerSubscriberPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExportFailure(t, args)
}

// Export a MCP Server from one environment as tenant subscriber user
func TestExportMCPServerTenantSubscriberUser(t *testing.T) {
	tenantSubscriber := subscriber.UserName + "@" + TENANT1
	tenantSubscriberPassword := subscriber.Password

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: tenantCreator, Password: tenantCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantSubscriber, Password: tenantSubscriberPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExportFailure(t, args)
}

// Export a MCP Server from one environment and import to another environment as tenant admin user
func TestExportImportMCPServerAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: tenantCreator, Password: tenantCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportImport(t, args)
}

// Export a MCP Server from one environment and import to another environment as tenant user with
// Internal/devops role
func TestExportImportMCPServerDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: tenantCreator, Password: tenantCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportImport(t, args)
}

// Export a MCP Server from one environment as super tenant admin user without specifying provider
func TestExportMCPServerAdminSuperTenantUserWithoutProvider(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: "", Password: ""},
		CtlUser:           testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Export a MCP Server from one environment as super tenant user with Internal/devops role without specifying provider
func TestExportMCPServerDevopsSuperTenantUserWithoutProvider(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: "", Password: ""},
		CtlUser:           testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Export a MCP Server from one environment as tenant admin user without specifying provider
func TestExportMCPServerAdminTenantUserWithoutProvider(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: "", Password: ""},
		CtlUser:           testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Export a MCP Server from one environment as tenant user with Internal/devops role without specifying provider
func TestExportMCPServerDevopsTenantUserWithoutProvider(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: "", Password: ""},
		CtlUser:           testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Export a MCP Server using a tenant user by specifying the provider name - MCP Server is in a different tenant
func TestExportMCPServerAdminTenantUserFromAnotherTenant(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: superTenantMCPServerCreator, Password: superTenantMCPServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportFailure(t, args)
}

// Export a MCP Server using a tenant user with Internal/devops role by specifying the provider name - MCP Server is in a different tenant
func TestExportMCPServerDevopsTenantUserFromAnotherTenant(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: superTenantMCPServerCreator, Password: superTenantMCPServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportFailure(t, args)
}

// Export a MCP Server using a tenant user without specifying the provider name - MCP Server is in a different tenant
func TestExportMCPServerAdminTenantUserFromAnotherTenantWithoutProvider(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: "", Password: ""},
		CtlUser:           testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportFailure(t, args)
}

// Export a MCP Server using a tenant user with Internal/devops role without specifying the provider name - MCP Server is in a different tenant
func TestExportMCPServerDevopsTenantUserFromAnotherTenantWithoutProvider(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: "", Password: ""},
		CtlUser:           testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExportFailure(t, args)
}

// Export a MCP Server from one environment as super tenant admin and import to another environment as cross tenant admin
// (with preserve-provider=false)
func TestExportImportMCPServerCrossTenantUserWithoutPreserveProvider(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: superTenantMCPServerCreator, Password: superTenantMCPServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExport(t, args)

	// Since --preserve-provider=false both the mcpServerProvider and the ctlUser is tenant admin
	args.MCPServerProvider = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the MCP Server to env2 as tenant admin across domains
	testutils.ValidateMCPServerExportImport(t, args)
}

// Export a MCP Server from one environment as super tenant user with Internal/devops role
// and import to another environment as cross tenant user with Internal/devops role (with preserve-provider=false)
func TestExportImportMCPServerCrossTenantDevopsUserWithoutPreserveProvider(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: superTenantMCPServerCreator, Password: superTenantMCPServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExport(t, args)

	// Since --preserve-provider=false both the mcpServerProvider and the ctlUser is tenant user with Internal/devops role
	args.MCPServerProvider = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}
	args.CtlUser = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}

	// Import the MCP Server to env2 as tenant user with Internal/devops role across domains
	testutils.ValidateMCPServerExportImport(t, args)
}

// Export a MCP Server from one environment as super tenant admin and import to another environment as cross tenant admin
// (without preserve-provider=false)
func TestExportImportMCPServerCrossTenantUser(t *testing.T) {
	superTenantAdminUsername := superAdminUser
	superTenantAdminPassword := superAdminPassword

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: superTenantMCPServerCreator, Password: superTenantMCPServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: superTenantAdminUsername, Password: superTenantAdminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExport(t, args)

	// Since --preserve-provider=false is not specified, the mcpServerProvider remain as it is and the ctlUser is tenant admin
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the MCP Server to env2 as tenant admin across domains - this should fail
	testutils.ValidateMCPServerExportFailure(t, args)
}

// Export a MCP Server from one environment as super tenant user with Internal/devops role
// and import to another environment as cross tenant user with Internal/devops role (without preserve-provider=false)
func TestExportImportMCPServerCrossTenantDevopsUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	superTenantMCPServerCreator := creator.UserName
	superTenantMCPServerCreatorPassword := creator.Password

	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	dev := GetDevClient()
	prod := GetProdClient()

	mcpServer := testutils.AddMCPServer(t, dev, superTenantMCPServerCreator, superTenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: superTenantMCPServerCreator, Password: superTenantMCPServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
		DestAPIM:          prod,
	}

	testutils.ValidateMCPServerExport(t, args)

	// Since --preserve-provider=false is not specified, the mcpServerProvider remain as it is and the ctlUser is tenant user
	// with Internal/devops role
	args.CtlUser = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}

	// Import the MCP Server to env2 as tenant admin across domains - this should fail
	testutils.ValidateMCPServerExportFailure(t, args)
}
