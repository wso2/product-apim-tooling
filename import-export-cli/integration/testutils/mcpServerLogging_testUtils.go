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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

func ValidateGetMCPServerLogLevel(t *testing.T, args *MCPServerLoggingTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	base.WaitForIndexing()

	output, _ := getMCPServerLogLevel(t, args)
	response, _ := args.APIM.GetMCPServerLogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.MCPServerId)

	validateGetMCPServerLogLevel(t, output, response)
}

func ValidateGetMCPServerLogLevelError(t *testing.T, args *MCPServerLoggingTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	base.WaitForIndexing()

	output, _ := getMCPServerLogLevel(t, args)
	_, err := args.APIM.GetMCPServerLogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.MCPServerId)

	validateInvalidPermissionError(t, output, err)
}

func ValidateSetMCPServerLogLevel(t *testing.T, args *MCPServerLoggingTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	base.WaitForIndexing()

	output, _ := setMCPServerLogLevel(t, args)
	response, _ := args.APIM.SetMCPServerLogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.MCPServerId, args.LogLevel)

	validateSetMCPServerLogLevel(t, output, response)
}

func ValidateSetMCPServerLogLevelError(t *testing.T, args *MCPServerLoggingTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)
	base.WaitForIndexing()

	output, _ := setMCPServerLogLevel(t, args)
	_, err := args.APIM.SetMCPServerLogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.MCPServerId, args.LogLevel)

	validateInvalidPermissionError(t, output, err)
}

func getMCPServerLogLevel(t *testing.T, args *MCPServerLoggingTestArgs) (string, error) {
	tenant_domain := ""
	tenant_domain_flag := ""
	if args.TenantDomain != "" && args.TenantDomain != "carbon.super" {
		tenant_domain = args.TenantDomain
		tenant_domain_flag = "--tenant-domain"
	}

	mcp_server_id := ""
	mcp_server_id_flag := ""
	if args.MCPServerId != "" {
		mcp_server_id = args.MCPServerId
		mcp_server_id_flag = "--mcp-server-id"
	}

	output, err := base.Execute(t, "get", "mcp-server-logging", "-e", args.APIM.EnvName, tenant_domain_flag, tenant_domain, mcp_server_id_flag, mcp_server_id, "-k", "--verbose")
	return output, err
}

func setMCPServerLogLevel(t *testing.T, args *MCPServerLoggingTestArgs) (string, error) {
	output, err := base.Execute(t, "set", "mcp-server-logging", "-e", args.APIM.EnvName, "--tenant-domain", args.TenantDomain, "--mcp-server-id", args.MCPServerId, "--log-level", args.LogLevel, "-k", "--verbose")
	return output, err
}

func validateGetMCPServerLogLevel(t *testing.T, output string, response *apim.MCPServerLogLevelList) {
	for index, mcp := range response.McpServers {
		log_line := strings.Split(output, "\n")[index+1]
		assert.Contains(t, log_line, mcp.MCPServerId, "MCP Server '"+mcp.MCPServerId+"' is not listed in the apictl output.")
		assert.Contains(t, log_line, mcp.LogLevel, "Log level of MCP Server '"+mcp.MCPServerId+"' is not matching.")
	}
}

func validateSetMCPServerLogLevel(t *testing.T, output string, response *apim.MCPServerLogLevel) {
	assert.Equal(t, "Log level "+response.LogLevel+" is successfully set to the MCP Server.\n", output, "Invalid apictl output after setting MCP Server log level.")
}
