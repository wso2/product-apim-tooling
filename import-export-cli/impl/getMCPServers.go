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

package impl

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	mcpServerIdHeader       = "ID"
	mcpServerNameHeader     = "NAME"
	mcpContextHeader        = "CONTEXT"
	mcpServerVersionHeader  = "VERSION"
	mcpServerProviderHeader = "PROVIDER"
	mcpServerStatusHeader   = "STATUS"

	defaultMCPServerTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Version}}\t{{.Context}}\t{{.LifeCycleStatus}}\t{{.Provider}}"
)

// mcpServer holds information about an MCP Server for outputting
type mcpServer struct {
	id              string
	name            string
	context         string
	version         string
	provider        string
	lifeCycleStatus string
}

// creates a new mcpServer from utils.MCPServer
func newMCPServerDefinitionFromItem(s utils.MCPServer) *mcpServer {
	return &mcpServer{s.ID, s.Name, s.Context, s.Version, s.Provider, s.LifeCycleStatus}
}

// Id of mcpServer
func (s mcpServer) Id() string {
	return s.id
}

// Name of mcpServer
func (s mcpServer) Name() string {
	return s.name
}

// Context of mcpServer
func (s mcpServer) Context() string {
	return s.context
}

// Version of mcpServer
func (s mcpServer) Version() string {
	return s.version
}

// Lifecycle Status of mcpServer
func (s mcpServer) LifeCycleStatus() string {
	return s.lifeCycleStatus
}

// Provider of mcpServer
func (s mcpServer) Provider() string {
	return s.provider
}

// MarshalJSON marshals mcpServer using custom marshaller which uses methods instead of fields
func (s *mcpServer) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(s)
}

// GetMCPServerListFromEnv
// @param accessToken : Access Token for the environment
// @param environment : Environment name to use when getting the MCP Server List
// @param query : string to be matched against the MCP Server names
// @param limit : total # of results to return
// @return count (no. of MCP Servers)
// @return array of MCPServer objects
// @return error
func GetMCPServerListFromEnv(accessToken, environment, query, limit string) (count int32, servers []utils.MCPServer, err error) {
	mcpListEndpoint := utils.GetMcpServerListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return GetMCPServerList(accessToken, mcpListEndpoint, query, limit)
}

// PrintMCPServers prints the list of MCP servers in a specific format
func PrintMCPServers(servers []utils.MCPServer, format string) {
	if format == "" {
		format = defaultMCPServerTableFormat
	} else if format == utils.JsonArrayFormatType {
		utils.ListArtifactsInJsonArrayFormat(servers, utils.ProjectTypeMcpServer)
		return
	}

	// create mcp server context with standard output
	mcpServerContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, s := range servers {
			if err := t.Execute(w, newMCPServerDefinitionFromItem(s)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	mcpServerTableHeaders := map[string]string{
		"Id":              mcpServerIdHeader,
		"Name":            mcpServerNameHeader,
		"Context":         mcpContextHeader,
		"Version":         mcpServerVersionHeader,
		"LifeCycleStatus": mcpServerStatusHeader,
		"Provider":        mcpServerProviderHeader,
	}

	// execute context
	if err := mcpServerContext.Write(renderer, mcpServerTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
