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

	"github.com/wso2/product-apim-tooling/import-export-cli/integration/testutils"
)

func TestExportInvalidMCPServerRevision(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			mcpServer := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			args := &testutils.MCPServerImportExportTestArgs{
				MCPServerProvider: user.ApiCreator,
				CtlUser:           user.CtlUser,
				MCPServer:         mcpServer,
				SrcAPIM:           dev,
				Revision:          "999", // Invalid revision number
			}
			testutils.ValidateMCPServerRevisionExportFailure(t, args)
		})
	}
}

func TestExportMCPServerRevision(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			mcpServer := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			args := &testutils.MCPServerImportExportTestArgs{
				MCPServerProvider: user.ApiCreator,
				CtlUser:           user.CtlUser,
				MCPServer:         mcpServer,
				SrcAPIM:           dev,
			}
			testutils.ValidateMCPServerExport(t, args)
		})
	}
}

// Test export of MCP Server revision as devops super tenant user
func TestExportMCPServerRevisionDevopsSuperTenantUser(t *testing.T) {
	devopsUsername := devops.UserName
	devopsPassword := devops.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: devopsUsername, Password: devopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Test export of MCP Server revision as publisher super tenant user
func TestExportMCPServerRevisionPublisherSuperTenantUser(t *testing.T) {
	publisherUsername := publisher.UserName
	publisherPassword := publisher.Password

	mcpServerCreator := creator.UserName
	mcpServerCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, mcpServerCreator, mcpServerCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: mcpServerCreator, Password: mcpServerCreatorPassword},
		CtlUser:           testutils.Credentials{Username: publisherUsername, Password: publisherPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Test export of MCP Server revision as tenant admin user
func TestExportMCPServerRevisionAdminTenantUser(t *testing.T) {
	tenantAdminUsername := superAdminUser + "@" + TENANT1
	tenantAdminPassword := superAdminPassword

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: tenantCreator, Password: tenantCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantAdminUsername, Password: tenantAdminPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Test export of MCP Server revision as tenant devops user
func TestExportMCPServerRevisionDevopsTenantUser(t *testing.T) {
	tenantDevopsUsername := devops.UserName + "@" + TENANT1
	tenantDevopsPassword := devops.Password

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: tenantCreator, Password: tenantCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantDevopsUsername, Password: tenantDevopsPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

// Test export of MCP Server revision as tenant publisher user
func TestExportMCPServerRevisionTenantPublisherUser(t *testing.T) {
	tenantPublisherUsername := publisher.UserName + "@" + TENANT1
	tenantPublisherPassword := publisher.Password

	tenantCreator := creator.UserName + "@" + TENANT1
	tenantCreatorPassword := creator.Password

	dev := GetDevClient()

	mcpServer := testutils.AddMCPServer(t, dev, tenantCreator, tenantCreatorPassword)

	args := &testutils.MCPServerImportExportTestArgs{
		MCPServerProvider: testutils.Credentials{Username: tenantCreator, Password: tenantCreatorPassword},
		CtlUser:           testutils.Credentials{Username: tenantPublisherUsername, Password: tenantPublisherPassword},
		MCPServer:         mcpServer,
		SrcAPIM:           dev,
	}

	testutils.ValidateMCPServerExport(t, args)
}

func TestExportImportMCPServerRevision(t *testing.T) {
	for _, user := range testCaseUsers {
		t.Run(user.Description, func(t *testing.T) {
			dev := GetDevClient()
			prod := GetProdClient()
			mcpServer := testutils.AddMCPServer(t, dev, user.ApiCreator.Username, user.ApiCreator.Password)
			args := &testutils.MCPServerImportExportTestArgs{
				MCPServerProvider: user.ApiCreator,
				CtlUser:           user.CtlUser,
				MCPServer:         mcpServer,
				SrcAPIM:           dev,
				DestAPIM:          prod,
			}
			testutils.ValidateMCPServerExportImport(t, args)
		})
	}
}

// Test export import of MCP Server revision as devops super tenant user
func TestExportImportMCPServerRevisionDevopsSuperTenantUser(t *testing.T) {
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

// Test export import of MCP Server revision as tenant admin user
func TestExportImportMCPServerRevisionAdminTenantUser(t *testing.T) {
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

// Test export import of MCP Server revision as tenant devops user
func TestExportImportMCPServerRevisionDevopsTenantUser(t *testing.T) {
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
