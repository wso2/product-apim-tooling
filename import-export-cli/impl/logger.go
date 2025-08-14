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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	loggingApiIdHeader       = "API_ID"
	loggingApiContextHeader  = "API_CONTEXT"
	loggingApiLogLevelHeader = "LOG_LEVEL"

	loggingMCPServerIdHeader       = "MCP_SERVER_ID"
	loggingMCPServerContextHeader  = "MCP_SERVER_CONTEXT"
	loggingMCPServerLogLevelHeader = "LOG_LEVEL"

	defaultLoggingApiTableFormat         = "table {{.Id}}\t{{.Context}}\t{{.LogLevel}}"
	defaultLoggingMCPServerTableFormat   = "table {{.Id}}\t{{.Context}}\t{{.LogLevel}}"
	defaultLoggingCorrelationTableFormat = "table {{.Name}}\t{{.Enabled}}\t{{.Properties}}"
)

// LoggingApi holds information about an API for outputting
type loggingApi struct {
	id       string
	context  string
	logLevel string
}

// LoggingMCPServer holds information about an MCP Server for outputting
type loggingMCPServer struct {
	id       string
	context  string
	logLevel string
}

type loggingCorrelationComponent struct {
	name       string
	enabled    string
	properties string
}

func newLoggingCorrelationComponentFromComponent(component utils.CorrelationComponent) *loggingCorrelationComponent {
	if len(component.Properties) > 0 {
		return &loggingCorrelationComponent{component.Name, component.Enabled, component.Properties[len(component.Properties)-1].Name +
			" : " + strings.Join(component.Properties[len(component.Properties)-1].Value, ", ")}
	}
	return &loggingCorrelationComponent{component.Name, component.Enabled, "-"}
}

func (component loggingCorrelationComponent) Name() string {
	return component.name
}

func (component loggingCorrelationComponent) Enabled() string {
	return component.enabled
}

func (component loggingCorrelationComponent) Properties() string {
	return component.properties
}

// Creates a new api from utils.APILogger
func newLoggingApiDefinitionFromAPI(a utils.APILogger) *loggingApi {
	return &loggingApi{a.ID, a.Context, a.LogLevel}
}

// Creates a new mcp server from utils.MCPServerLogger
func newLoggingMCPServerDefinitionFromMCPServer(m utils.MCPServerLogger) *loggingMCPServer {
	return &loggingMCPServer{m.ID, m.Context, m.LogLevel}
}

// Id of loggingApi
func (a loggingApi) Id() string {
	return a.id
}

// Context of loggingApi
func (a loggingApi) Context() string {
	return a.context
}

// LogLevel of loggingApi
func (a loggingApi) LogLevel() string {
	return a.logLevel
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (a *loggingApi) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

// Id of loggingMCPServer
func (m loggingMCPServer) Id() string {
	return m.id
}

// Context of loggingMCPServer
func (m loggingMCPServer) Context() string {
	return m.context
}

// LogLevel of loggingMCPServer
func (m loggingMCPServer) LogLevel() string {
	return m.logLevel
}

// MarshalJSON marshals mcp server using custom marshaller which uses methods instead of fields
func (m *loggingMCPServer) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(m)
}

// Add new APILogger object and use it in the array
func GetPerAPILoggingListFromEnv(credential credentials.Credential, environment, tenantDomain string) (apis []utils.APILogger, err error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	if tenantDomain == "" {
		tenantDomain = utils.DefaultTenantDomain
	}
	apiListEndpoint := utils.GetAPILoggingListEndpointOfEnv(environment, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", apiListEndpoint)
	resp, err := utils.InvokeGETRequest(apiListEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiLoggerListResponse := &utils.APILoggerListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiLoggerListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return apiLoggerListResponse.Apis, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}

// PrintAPILoggers
func PrintAPILoggers(apis []utils.APILogger, format string) {
	if format == "" {
		format = defaultLoggingApiTableFormat
	}
	// Create API context with standard output
	apiContext := formatter.NewContext(os.Stdout, format)

	// Create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apis {
			if err := t.Execute(w, newLoggingApiDefinitionFromAPI(a)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// Headers for table
	apiLoggerTableHeaders := map[string]string{
		"Id":       loggingApiIdHeader,
		"Context":  loggingApiContextHeader,
		"LogLevel": loggingApiLogLevelHeader,
	}

	// Execute context
	if err := apiContext.Write(renderer, apiLoggerTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}

// Check how to send the API details in a single object instead of an array
func GetPerAPILoggingDetailsFromEnv(credential credentials.Credential, environment, apiId, tenantDomain string) (apis []utils.APILogger, err error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	if tenantDomain == "" {
		tenantDomain = utils.DefaultTenantDomain
	}
	apiDetailsEndpoint := utils.GetAPILoggingDetailsEndpointOfEnv(environment, apiId, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", apiDetailsEndpoint)
	resp, err := utils.InvokeGETRequest(apiDetailsEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiDetailsEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiLoggerListResponse := &utils.APILoggerListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiLoggerListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return apiLoggerListResponse.Apis, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}

// SetAPILoggingLevel
func SetAPILoggingLevel(credential credentials.Credential, environment, apiId, tenantDomain, logLevel string) (*resty.Response, error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	if tenantDomain == "" {
		tenantDomain = utils.DefaultTenantDomain
	}
	apiSetEndpoint := utils.GetAPILoggingSetEndpointOfEnv(environment, apiId, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", apiSetEndpoint)
	body := `{"logLevel":"` + logLevel + `"}`
	resp, err := utils.InvokePutRequest(nil, apiSetEndpoint, headers, body)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiSetEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		return resp, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}

// GET Correlation Components List
func GetCorrelationLogComponentListFromEnv(credential credentials.Credential, environment string) (components []utils.CorrelationComponent, err error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials

	correlationDevOpsEP := utils.GetCorrelationLoggingEndPointOfEnv(environment, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo + "URL : " + correlationDevOpsEP)
	resp, err := utils.InvokeGETRequest(correlationDevOpsEP, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+correlationDevOpsEP, err)
	}

	utils.Logln(utils.LogPrefixInfo + "Response : " + resp.Status())

	if resp.StatusCode() == http.StatusOK {
		correlationComponentsResponse := &utils.CorrelationComponentList{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &correlationComponentsResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return correlationComponentsResponse.Components, nil
	}
	return nil, errors.New(string(resp.Body()))
}

func PrintCorrelationLoggers(components []utils.CorrelationComponent, format string) {
	if format == "" {
		format = defaultLoggingCorrelationTableFormat
	}

	formatContext := formatter.NewContext(os.Stdout, format)

	renderer := func(w io.Writer, t *template.Template) error {
		for _, component := range components {
			if err := t.Execute(w, newLoggingCorrelationComponentFromComponent(component)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	correlationTableHeaders := map[string]string{
		"Name":       "COMPONENT_NAME",
		"Enabled":    "ENABLED",
		"Properties": "PROPERTIES",
	}

	if err := formatContext.Write(renderer, correlationTableHeaders); err != nil {
		fmt.Println("Error executing template", err.Error())
	}
}
func SetCorrelationLoggingComponent(credential credentials.Credential, environment, componentName, enabled, deniedThreads string) (*resty.Response, error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)

	correlationDevOpsEP := utils.GetCorrelationLoggingEndPointOfEnv(environment, utils.MainConfigFilePath)

	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON

	resp, err := utils.InvokeGETRequest(correlationDevOpsEP, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+correlationDevOpsEP, err)
	}

	utils.Logln(utils.LogPrefixInfo + "GET Response : " + resp.Status())

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	correlationComponentsResponse := &utils.CorrelationComponentList{}
	unmarshalError := json.Unmarshal([]byte(resp.Body()), &correlationComponentsResponse)

	if unmarshalError != nil {
		utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
	}

	responseComponents := correlationComponentsResponse.Components
	requestComponents := make([]utils.CorrelationComponent, 0)

	for _, cc := range responseComponents {
		if cc.Name == componentName {
			cc.Enabled = enabled
			if len(cc.Properties) > 0 {
				if cc.Properties[len(cc.Properties)-1].Name == "deniedThreads" {
					cc.Properties[len(cc.Properties)-1].Value = strings.Split(deniedThreads, ",")
				}
			}
		}
		requestComponents = append(requestComponents, cc)
	}

	b, err := json.Marshal(requestComponents)
	if err != nil {
		utils.Logln("Error when creating a json ")
		return nil, nil
	}

	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	body := string(b)
	body = "{\"components\":" + body + "}"
	putResp, err := utils.InvokePutRequest(nil, correlationDevOpsEP, headers, body)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+correlationDevOpsEP, err)
	}

	utils.Logln(utils.LogPrefixInfo+"PUT Response:", putResp.Status())
	if putResp.StatusCode() == http.StatusOK {
		return putResp, nil
	}
	return nil, errors.New(string(putResp.Body()))
}

// GetPerMCPServerLoggingListFromEnv gets MCP Server logging list from environment
func GetPerMCPServerLoggingListFromEnv(credential credentials.Credential, environment, tenantDomain string) (mcpServers []utils.MCPServerLogger, err error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	if tenantDomain == "" {
		tenantDomain = utils.DefaultTenantDomain
	}
	mcpServerListEndpoint := utils.GetMCPServerLoggingListEndpointOfEnv(environment, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", mcpServerListEndpoint)
	resp, err := utils.InvokeGETRequest(mcpServerListEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+mcpServerListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+" Response: ", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		mcpServerLoggerListResponse := &utils.MCPServerLoggerListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &mcpServerLoggerListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response ", unmarshalError)
		}

		return mcpServerLoggerListResponse.MCPServers, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}

// GetPerMCPServerLoggingDetailsFromEnv gets specific MCP Server logging details from environment
func GetPerMCPServerLoggingDetailsFromEnv(credential credentials.Credential, environment, mcpServerId, tenantDomain string) (mcpServers []utils.MCPServerLogger, err error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	if tenantDomain == "" {
		tenantDomain = utils.DefaultTenantDomain
	}
	mcpServerDetailsEndpoint := utils.GetMCPServerLoggingDetailsEndpointOfEnv(environment, mcpServerId, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", mcpServerDetailsEndpoint)
	resp, err := utils.InvokeGETRequest(mcpServerDetailsEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+mcpServerDetailsEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response: ", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		mcpServerLoggerListResponse := &utils.MCPServerLoggerListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &mcpServerLoggerListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response ", unmarshalError)
		}

		return mcpServerLoggerListResponse.MCPServers, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}

// SetMCPServerLoggingLevel sets MCP Server logging level
func SetMCPServerLoggingLevel(credential credentials.Credential, environment, mcpServerId, tenantDomain, logLevel string) (*resty.Response, error) {
	// Base64 encoding the credentials
	b64encodedCredentials := credentials.GetBasicAuth(credential)
	// Prepping the headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	if tenantDomain == "" {
		tenantDomain = utils.DefaultTenantDomain
	}
	mcpServerSetEndpoint := utils.GetMCPServerLoggingSetEndpointOfEnv(environment, mcpServerId, tenantDomain, utils.MainConfigFilePath)
	utils.Logln(utils.LogPrefixInfo+"URL:", mcpServerSetEndpoint)
	body := `{"logLevel":"` + logLevel + `"}`
	resp, err := utils.InvokePutRequest(nil, mcpServerSetEndpoint, headers, body)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+mcpServerSetEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response: ", resp.Status())

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		return resp, nil
	} else {
		return nil, errors.New(string(resp.Body()))
	}
}

// PrintMCPServerLoggers prints MCP Server loggers
func PrintMCPServerLoggers(mcpServers []utils.MCPServerLogger, format string) {
	if format == "" {
		format = defaultLoggingMCPServerTableFormat
	}
	// Create MCP Server context with standard output
	mcpServerContext := formatter.NewContext(os.Stdout, format)

	// Create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, m := range mcpServers {
			if err := t.Execute(w, newLoggingMCPServerDefinitionFromMCPServer(m)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// Headers for table
	mcpServerLoggerTableHeaders := map[string]string{
		"Id":       loggingMCPServerIdHeader,
		"Context":  loggingMCPServerContextHeader,
		"LogLevel": loggingMCPServerLogLevelHeader,
	}

	// Execute context
	if err := mcpServerContext.Write(renderer, mcpServerLoggerTableHeaders); err != nil {
		fmt.Println("Error executing template: ", err.Error())
	}
}
