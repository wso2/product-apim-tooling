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
	"fmt"
	"os"
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
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
		OverrideProvider:  true,
	}

	testutils.ValidateMCPServerExport(t, args)

	// Since --preserve-provider=false both the mcpServerProvider and the ctlUser is tenant admin
	args.MCPServerProvider = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}
	args.CtlUser = testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword}

	// Import the MCP Server to env2 as tenant admin across domains
	testutils.ValidateMCPServerImport(t, args)
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
		OverrideProvider:  true,
	}

	testutils.ValidateMCPServerExport(t, args)

	// Since --preserve-provider=false both the mcpServerProvider and the ctlUser is tenant user with Internal/devops role
	args.MCPServerProvider = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}
	args.CtlUser = testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword}

	// Import the MCP Server to env2 as tenant user with Internal/devops role across domains
	testutils.ValidateMCPServerImport(t, args)
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
	testutils.ValidateMCPServerImportFailure(t, args)
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
	testutils.ValidateMCPServerImportFailure(t, args)
}

// Export a MCP Server with the life cycle status as Blocked and import to another environment
// and import update it
func TestExportImportMCPServerBlocked(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			mcpServer := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishMCPServer(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, mcpServer.ID)
			mcpServer = testutils.ChangeMCPServerLifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, mcpServer.ID, "Block")

			args := &testutils.MCPServerImportExportTestArgs{
				MCPServerProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:           testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				MCPServer:         mcpServer,
				SrcAPIM:           dev,
				DestAPIM:          prod,
			}

			importedMCPServer := testutils.ValidateMCPServerExportImport(t, args)

			// Change the lifecycle to Published in the prod environment
			testutils.ChangeMCPServerLifeCycle(prod, user.ApiPublisher.Username, user.ApiPublisher.Password, importedMCPServer.ID, "Re-Publish")

			args.Update = true
			testutils.ValidateMCPServerExportImport(t, args)
		})
	}
}

// Import an MCP Server with the default version. Change the version and import the same MCP Server again.

func TestMCPServerVersioning(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			// Add MCP Server with default version
			mcpServer1 := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)

			args := &testutils.MCPServerImportExportTestArgs{
				MCPServerProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:           testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				MCPServer:         mcpServer1,
				SrcAPIM:           dev,
				DestAPIM:          prod,
				OverrideProvider:  true,
			}

			// Export and import the MCP Server with default version
			testutils.ValidateMCPServerExport(t, args)
			importedMCPServer := testutils.ValidateMCPServerImportForMultipleVersions(t, args, "")

			// Change the version and update the MCP Server in dev
			mcpServer2 := testutils.AddCustomMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password,
				mcpServer1.Name, testutils.APIVersion2, mcpServer1.Context)

			args.MCPServer = mcpServer2

			// Export and import the MCP Server with new version
			testutils.ValidateMCPServerExport(t, args)
			testutils.ValidateMCPServerImportForMultipleVersions(t, args, importedMCPServer.ID)
		})
	}
}

// Export a MCP Server with the life cycle status as Deprecated and import to another environment
func TestExportImportMCPServerDeprecated(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			mcpServer := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishMCPServer(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, mcpServer.ID)
			mcpServer = testutils.ChangeMCPServerLifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, mcpServer.ID, "Deprecate")

			args := &testutils.MCPServerImportExportTestArgs{
				MCPServerProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:           testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				MCPServer:         mcpServer,
				SrcAPIM:           dev,
				DestAPIM:          prod,
			}

			testutils.ValidateMCPServerExportImport(t, args)
		})
	}
}

// Export a MCP Server with the life cycle status as Retired and import to another environment
func TestExportImportMCPServerRetired(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			dev := GetDevClient()
			prod := GetProdClient()

			mcpServer := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			testutils.PublishMCPServer(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, mcpServer.ID)
			testutils.ChangeMCPServerLifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, mcpServer.ID, "Deprecate")
			mcpServer = testutils.ChangeMCPServerLifeCycle(dev, user.ApiPublisher.Username, user.ApiPublisher.Password, mcpServer.ID, "Retire")

			args := &testutils.MCPServerImportExportTestArgs{
				MCPServerProvider: testutils.Credentials{Username: user.ApiCreator.Username, Password: user.ApiCreator.Password},
				CtlUser:           testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				MCPServer:         mcpServer,
				SrcAPIM:           dev,
				DestAPIM:          prod,
			}

			testutils.ValidateMCPServerExportImport(t, args)
		})
	}
}

func TestListMCPServersAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		// Add the MCP Server to env1
		testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)
	}

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: adminUsername, Password: adminPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateMCPServersList(t, args)
}

func TestListMCPServersDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		// Add the MCP Server to env1
		testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)
	}

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateMCPServersList(t, args)
}

func TestListMCPServersAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreator := creator.UserName + "@" + TENANT1
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		// Add the MCP Server to env1
		testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)
	}

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateMCPServersList(t, args)
}

func TestListMCPServersDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	mcpServerCreator := creator.UserName + "@" + TENANT1
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		// Add the MCP Server to env1
		testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)
	}

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser: testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		SrcAPIM: dev,
	}

	testutils.ValidateMCPServersList(t, args)
}

// MCP Servers listing with JsonArray format
func TestListMCPServersWithJsonArrayFormat(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
				// Add the MCP Server to env1
				testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			}

			args := &testutils.MCPServerImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}

			testutils.ValidateMCPServersListWithJsonArrayFormat(t, args)
		})
	}
}

func TestDeleteMCPServerAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	var mcpServer *apim.MCPServer
	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		mcpServer = testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)
	}

	// This will be the MCP Server that will be deleted by apictl, so no need to do cleaning
	mcpServer = testutils.AddMCPServerWithoutCleaning(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser:   testutils.Credentials{Username: adminUsername, Password: adminPassword},
		MCPServer: mcpServer,
		SrcAPIM:   dev,
	}

	testutils.ValidateMCPServerDelete(t, args)
}

func TestDeleteMCPServerDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	var mcpServer *apim.MCPServer
	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		mcpServer = testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)
	}

	// This will be the MCP Server that will be deleted by apictl, so no need to do cleaning
	mcpServer = testutils.AddMCPServerWithoutCleaning(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser:   testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		MCPServer: mcpServer,
		SrcAPIM:   dev,
	}

	testutils.ValidateMCPServerDelete(t, args)
}

func TestDeleteMCPServerAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantMCPServerCreator := creator.UserName + "@" + TENANT1
	tenantMCPServerCreatorPassword := creator.Password

	dev := GetDevClient()

	var mcpServer *apim.MCPServer
	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		mcpServer = testutils.AddMCPServer(t, dev, tenantMCPServerCreator, tenantMCPServerCreatorPassword)
	}

	// This will be the MCP Server that will be deleted by apictl, so no need to do cleaning
	mcpServer = testutils.AddMCPServerWithoutCleaning(t, dev, tenantMCPServerCreator, tenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser:   testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServer: mcpServer,
		SrcAPIM:   dev,
	}

	testutils.ValidateMCPServerDelete(t, args)
}

func TestDeleteMCPServerDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantMCPServerCreator := creator.UserName + "@" + TENANT1
	tenantMCPServerCreatorPassword := creator.Password

	dev := GetDevClient()

	var mcpServer *apim.MCPServer
	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		mcpServer = testutils.AddMCPServer(t, dev, tenantMCPServerCreator, tenantMCPServerCreatorPassword)
	}

	// This will be the MCP Server that will be deleted by apictl, so no need to do cleaning
	mcpServer = testutils.AddMCPServerWithoutCleaning(t, dev, tenantMCPServerCreator, tenantMCPServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser:   testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		MCPServer: mcpServer,
		SrcAPIM:   dev,
	}

	testutils.ValidateMCPServerDelete(t, args)
}

func TestDeleteMCPServerSuperTenantUser(t *testing.T) {
	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	var mcpServer *apim.MCPServer
	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		mcpServer = testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)
	}

	// This will be the MCP Server that will be deleted by apictl, so no need to do cleaning
	mcpServer = testutils.AddMCPServerWithoutCleaning(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser:   testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		MCPServer: mcpServer,
		SrcAPIM:   dev,
	}

	testutils.ValidateMCPServerDelete(t, args)
}

func TestDeleteMCPServerWithActiveSubscriptionsSuperTenantUser(t *testing.T) {
	adminUser := superAdminUser
	adminPassword := superAdminPassword

	mcpServerPublisher := publisher.UserName
	mcpServerPublisherPassword := publisher.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	// Create and Deploy Revision of the above MCP Server
	testutils.CreateAndDeployMCPServerRevision(t, dev, mcpServerPublisher, mcpServerPublisherPassword, mcpServer.ID)

	//Publish created MCP Server
	testutils.PublishMCPServer(dev, mcpServerPublisher, mcpServerPublisherPassword, mcpServer.ID)

	//args to delete MCP Server
	argsToDelete := &testutils.MCPServerImportExportTestArgs{
		CtlUser:   testutils.Credentials{Username: adminUser, Password: adminPassword},
		MCPServer: mcpServer,
		SrcAPIM:   dev,
	}

	//validate MCP Server with active subscriptions delete failure
	testutils.ValidateMCPServerDeleteFailure(t, argsToDelete)
}

func TestExportMCPServersWithExportMCPServersCommand(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	dev := GetDevClient()

	var mcpServer *apim.MCPServer
	var mcpServersAdded = 0
	for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
		mcpServer = testutils.AddMCPServer(t, dev, tenantAdminUsername, tenantAdminPassword)
		testutils.CreateAndDeployMCPServerRevision(t, dev, tenantAdminUsername, tenantAdminPassword, mcpServer.ID)
		mcpServersAdded++
	}

	// This will be the MCP Server that will be deleted by apictl, so no need to do cleaning
	mcpServer = testutils.AddMCPServerWithoutCleaning(t, dev, tenantAdminUsername, tenantAdminPassword)
	testutils.CreateAndDeployMCPServerRevision(t, dev, tenantAdminUsername, tenantAdminPassword, mcpServer.ID)

	args := &testutils.MCPServerImportExportTestArgs{
		CtlUser:   testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServer: mcpServer,
		SrcAPIM:   dev,
	}

	testutils.ValidateAllMCPServersOfATenantIsExported(t, args, mcpServersAdded)
}

// Export MCP Servers bunch at once with export mcp-servers command and then add new MCP Servers and export MCP Servers once again to check whether
// the new MCP Servers exported
func TestExportMCPServersTwiceWithAfterAddingMCPServers(t *testing.T) {

	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			var mcpServer *apim.MCPServer
			var mcpServersAdded = 0
			for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
				mcpServer := testutils.AddMCPServer(t, dev, user.Admin.Username, user.Admin.Password)
				testutils.CreateAndDeployMCPServerRevision(t, dev, user.Admin.Username, user.Admin.Password, mcpServer.ID)
				mcpServersAdded++
			}

			// This will be the MCP Server that will be deleted by apictl, so no need to do cleaning
			mcpServer = testutils.AddMCPServerWithoutCleaning(t, dev, user.Admin.Username, user.Admin.Password)
			testutils.CreateAndDeployMCPServerRevision(t, dev, user.Admin.Username, user.Admin.Password, mcpServer.ID)

			args := &testutils.MCPServerImportExportTestArgs{
				CtlUser:   testutils.Credentials{Username: user.Admin.Username, Password: user.Admin.Password},
				MCPServer: mcpServer,
				SrcAPIM:   dev,
			}

			testutils.ValidateAllMCPServersOfATenantIsExported(t, args, mcpServersAdded)

			// Add new MCP Server and deploy
			mcpServer = testutils.AddMCPServer(t, dev, user.Admin.Username, user.Admin.Password)
			testutils.CreateAndDeployMCPServerRevision(t, dev, user.Admin.Username, user.Admin.Password, mcpServer.ID)
			newMCPServerCount := mcpServersAdded + 1

			// Validate again to check whether the newly added MCP Server exported properly.
			testutils.ValidateAllMCPServersOfATenantIsExported(t, args, newMCPServerCount)
		})
	}
}

func TestChangeLifeCycleStatusOfMCPServerAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	// Add the MCP Server to env
	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	//Change life cycle state of MCP Server from CREATED to PUBLISHED
	args := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: adminUsername, Password: adminPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, args)

	//Change life cycle state of MCP Server from PUBLISHED to CREATED
	argsToNextChange := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: adminUsername, Password: adminPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfMCPServerDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	// Add the MCP Server to env
	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	//Change life cycle state of MCP Server from CREATED to PUBLISHED
	args := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, args)

	//Change life cycle state of MCP Server from PUBLISHED to CREATED
	argsToNextChange := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfMCPServerAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	mcpServerCreator := creator.UserName + "@" + TENANT1
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	// Add the MCP Server to env
	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	//Change life cycle state of MCP Server from CREATED to PUBLISHED
	args := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, args)

	//Change life cycle state of MCP Server from PUBLISHED to CREATED
	argsToNextChange := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfMCPServerDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	mcpServerCreator := creator.UserName + "@" + TENANT1
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()
	// Add the MCP Server to env
	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	//Change life cycle state of MCP Server from CREATED to PUBLISHED
	args := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, args)

	//Change life cycle state of MCP Server from PUBLISHED to CREATED
	argsToNextChange := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, argsToNextChange)
}

func TestChangeLifeCycleStatusOfMCPServerFailWithAUserWithoutPermissions(t *testing.T) {
	subscriberUsername := subscriber.UserName
	subscriberDevopsPassword := subscriber.Password

	mcpServerCreator := creator.UserName + "@" + TENANT1
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()
	// Add the MCP Server to env
	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	//Change life cycle state of MCP Server from CREATED to PUBLISHED
	args := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: subscriberUsername, Password: subscriberDevopsPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Publish",
		ExpectedState: "PUBLISHED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServerFailure(t, args)
}

func TestChangeLifeCycleStatusOfMCPServerWithActiveSubscriptionWithAdminSuperTenantUser(t *testing.T) {
	adminUsername := superAdminUser
	adminPassword := superAdminPassword

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	mcpServerPublisher := publisher.UserName
	mcpServerPublisherPassword := publisher.Password

	mcpServerSubscriber := subscriber.UserName
	mcpServerSubscriberPassword := subscriber.Password

	dev := GetDevClient()

	// Add the MCP Server to env
	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	// Create and Deploy Revision of the above MCP Server
	testutils.CreateAndDeployMCPServerRevision(t, dev, mcpServerPublisher, mcpServerPublisherPassword, mcpServer.ID)

	testutils.PublishMCPServer(dev, adminUsername, adminPassword, mcpServer.ID)

	// Create an App
	app := testutils.AddApp(t, dev, mcpServerSubscriber, mcpServerSubscriberPassword)

	//Create an active subscription for MCP Server
	testutils.AddSubscription(t, dev, mcpServer.ID, app.ApplicationID, testutils.UnlimitedPolicy,
		mcpServerSubscriber, mcpServerSubscriberPassword)

	base.WaitForIndexing()

	//Change life cycle state of MCP Server from PUBLISHED to CREATED
	argsToLifeCycleStateChange := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: adminUsername, Password: adminPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, argsToLifeCycleStateChange)
}

func TestChangeLifeCycleStatusOfMCPServerWithActiveSubscriptionDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	mcpServerPublisher := publisher.UserName
	mcpServerPublisherPassword := publisher.Password

	mcpServerSubscriber := subscriber.UserName
	mcpServerSubscriberPassword := subscriber.Password

	dev := GetDevClient()

	// Add the MCP Server to env
	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	// Create and Deploy Revision of the above MCP Server
	testutils.CreateAndDeployMCPServerRevision(t, dev, mcpServerPublisher, mcpServerPublisherPassword, mcpServer.ID)

	testutils.PublishMCPServer(dev, devopsUsername, devopsPassword, mcpServer.ID)

	// Create an App
	app := testutils.AddApp(t, dev, mcpServerSubscriber, mcpServerSubscriberPassword)

	//Create an active subscription for MCP Server
	testutils.AddSubscription(t, dev, mcpServer.ID, app.ApplicationID, testutils.UnlimitedPolicy,
		mcpServerSubscriber, mcpServerSubscriberPassword)

	base.WaitForIndexing()

	//Change life cycle state of MCP Server from PUBLISHED to CREATED
	argsToLifeCycleStateChange := &testutils.MCPServerChangeLifeCycleStatusTestArgs{
		CtlUser:       testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		APIM:          dev,
		MCPServer:     mcpServer,
		Action:        "Demote to Created",
		ExpectedState: "CREATED",
	}

	testutils.ValidateChangeLifeCycleStatusOfMCPServer(t, argsToLifeCycleStateChange)
}

// Import a MCP Server and then create a new version of that MCP Server by updating the context and version only and import again
func TestCreateNewVersionOfMCPServerByUpdatingVersion(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {

			apim := GetDevClient()
			projectName := base.GenerateRandomName(16)

			args := &testutils.InitTestArgs{
				CtlUser:   testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM:   apim,
				InitFlag:  projectName,
				OasFlag:   testutils.TestSwagger2DefinitionPath,
				APIName:   testutils.DevFirstSwagger2APIName,
				ForceFlag: false,
			}

			//Initialize a project with MCP Server definition
			testutils.ValidateInitializeProjectWithOASFlag(t, args)

			//Assert that project import to publisher portal is successful
			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))

			// Read the MCP Server definition file in the project
			mcpServerDefinitionFilePath := args.InitFlag + string(os.PathSeparator) + utils.APIDefinitionFileYaml
			mcpServerDefinitionFileContent := testutils.ReadAPIDefinition(t, mcpServerDefinitionFilePath)

			//Change the version
			newVersion := base.GenerateRandomString()
			mcpServerDefinitionFileContent.Data.Version = newVersion

			// Write the modified MCP Server definition to the directory
			testutils.WriteToAPIDefinition(t, mcpServerDefinitionFileContent, mcpServerDefinitionFilePath)

			// Import and validate new MCP Server with version change
			testutils.ValidateImportProject(t, args, "", !isTenantUser(user.CtlUser.Username, TENANT1))

			testutils.ValidateMCPServersListWithVersionsFromInitArgs(t, args, newVersion)
		})
	}
}

// MCP Server search using query parameters
func TestMCPServerSearchWithQueryParams(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()

			var searchQuery string

			// Add set of MCP Servers to env and store mcp server details
			var addedMCPServersList [numberOfMCPServers + 1]*apim.MCPServer
			for mcpServerCount := 0; mcpServerCount <= numberOfMCPServers; mcpServerCount++ {
				// Add the MCP Server to env1
				mcpServer := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
				addedMCPServersList[mcpServerCount] = mcpServer
			}

			// Add custom MCP Server
			customMCPServer := addedMCPServersList[3]
			customMCPServer.Name = testutils.CustomAPIName
			customMCPServer.Version = testutils.CustomAPIVersion
			customMCPServer.Context = testutils.CustomAPIContext
			dev.AddMCPServer(t, customMCPServer, user.ApiCreator.Username, user.ApiCreator.Password, true)

			args := &testutils.MCPServerImportExportTestArgs{
				CtlUser: testutils.Credentials{Username: user.CtlUser.Username, Password: user.CtlUser.Password},
				SrcAPIM: dev,
			}

			for i := 0; i < len(addedMCPServersList); i++ {
				mcpServerNameToSearch := addedMCPServersList[i].Name
				mcpServerNameNotToSearch := addedMCPServersList[len(addedMCPServersList)-(i+1)].Name
				searchQuery = fmt.Sprintf("--query name:%v", mcpServerNameToSearch)

				//Search MCP Servers using query
				testutils.ValidateSearchMCPServersList(t, args, searchQuery, mcpServerNameToSearch, mcpServerNameNotToSearch)

				//Select random context from the added MCP Servers
				mcpServerContextToSearch := addedMCPServersList[i].Context
				mcpServerContextNotToSearch := addedMCPServersList[len(addedMCPServersList)-(i+1)].Context
				searchQuery = fmt.Sprintf("--query context:%v", mcpServerContextToSearch)

				//Search MCP Servers using query
				testutils.ValidateSearchMCPServersList(t, args, searchQuery, mcpServerContextToSearch, mcpServerContextNotToSearch)
			}

			// Search custom MCP Server with name
			searchQuery = fmt.Sprintf("--query name:%v", testutils.CustomAPIName)
			testutils.ValidateSearchMCPServersList(t, args, searchQuery, testutils.CustomAPIName,
				addedMCPServersList[1].Name)

			// Search custom MCP Server with context
			searchQuery = fmt.Sprintf("--query context:%v", testutils.CustomAPIContext)
			testutils.ValidateSearchMCPServersList(t, args, searchQuery, testutils.CustomAPIContext,
				addedMCPServersList[1].Context)

			// Search custom MCP Server with version
			searchQuery = fmt.Sprintf("--query version:%v", testutils.CustomAPIVersion)
			testutils.ValidateSearchMCPServersList(t, args, searchQuery, testutils.CustomAPIVersion,
				addedMCPServersList[1].Version)

			// Search custom MCP Server with version and name
			searchQuery = fmt.Sprintf("--query version:%v --query name:%v", testutils.CustomAPIVersion, testutils.CustomAPIName)
			testutils.ValidateSearchMCPServersList(t, args, searchQuery, testutils.CustomAPIVersion,
				addedMCPServersList[1].Version)
		})
	}
}
