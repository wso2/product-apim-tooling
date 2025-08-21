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

package apim

// MCPServer : MCPServer DTO
type MCPServer struct {
	ID                           string                `json:"id,omitempty" yaml:"id,omitempty"`
	Name                         string                `json:"name,omitempty" yaml:"name,omitempty"`
	DisplayName                  string                `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Description                  string                `json:"description,omitempty" yaml:"description,omitempty"`
	Context                      string                `json:"context,omitempty" yaml:"context,omitempty"`
	EndpointConfig               interface{}           `json:"endpointConfig,omitempty" yaml:"endpointConfig,omitempty"`
	Version                      string                `json:"version,omitempty" yaml:"version,omitempty"`
	Provider                     string                `json:"provider,omitempty" yaml:"provider,omitempty"`
	LifeCycleStatus              string                `json:"lifeCycleStatus,omitempty" yaml:"lifeCycleStatus,omitempty"`
	HasThumbnail                 bool                  `json:"hasThumbnail,omitempty" yaml:"hasThumbnail,omitempty"`
	IsDefaultVersion             bool                  `json:"isDefaultVersion,omitempty" yaml:"isDefaultVersion,omitempty"`
	IsRevision                   bool                  `json:"isRevision,omitempty" yaml:"isRevision,omitempty"`
	RevisionedMCPServerId        string                `json:"revisionedMCPServerId,omitempty" yaml:"revisionedMCPServerId,omitempty"`
	RevisionID                   int                   `json:"revisionId,omitempty" yaml:"revisionId,omitempty"`
	EnableSchemaValidation       bool                  `json:"enableSchemaValidation,omitempty" yaml:"enableSchemaValidation,omitempty"`
	Audiences                    []string              `json:"audiences,omitempty" yaml:"audiences,omitempty"`
	Transport                    []string              `json:"transport,omitempty" yaml:"transport,omitempty"`
	Tags                         []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Policies                     []string              `json:"policies,omitempty" yaml:"policies,omitempty"`
	OrganizationPolicies         interface{}           `json:"organizationPolicies,omitempty" yaml:"organizationPolicies,omitempty"`
	ThrottlingPolicy             string                `json:"throttlingPolicy,omitempty" yaml:"throttlingPolicy,omitempty"`
	AuthorizationHeader          string                `json:"authorizationHeader,omitempty" yaml:"authorizationHeader,omitempty"`
	ApiKeyHeader                 string                `json:"apiKeyHeader,omitempty" yaml:"apiKeyHeader,omitempty"`
	SecurityScheme               []string              `json:"securityScheme,omitempty" yaml:"securityScheme,omitempty"`
	MaxTps                       interface{}           `json:"maxTps,omitempty" yaml:"maxTps,omitempty"`
	Visibility                   string                `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	VisibleRoles                 []string              `json:"visibleRoles,omitempty" yaml:"visibleRoles,omitempty"`
	VisibleTenants               []string              `json:"visibleTenants,omitempty" yaml:"visibleTenants,omitempty"`
	VisibleOrganizations         []string              `json:"visibleOrganizations,omitempty" yaml:"visibleOrganizations,omitempty"`
	MCPServerPolicies            interface{}           `json:"mcpServerPolicies,omitempty" yaml:"mcpServerPolicies,omitempty"`
	SubscriptionAvailability     string                `json:"subscriptionAvailability,omitempty" yaml:"subscriptionAvailability,omitempty"`
	SubscriptionAvailableTenants []string              `json:"subscriptionAvailableTenants,omitempty" yaml:"subscriptionAvailableTenants,omitempty"`
	AdditionalPropertiesMap      interface{}           `json:"additionalPropertiesMap,omitempty" yaml:"additionalPropertiesMap,omitempty"`
	Monetization                 interface{}           `json:"monetization,omitempty" yaml:"monetization,omitempty"`
	AccessControl                string                `json:"accessControl,omitempty" yaml:"accessControl,omitempty"`
	AccessControlRoles           []string              `json:"accessControlRoles,omitempty" yaml:"accessControlRoles,omitempty"`
	BusinessInformation          interface{}           `json:"businessInformation,omitempty" yaml:"businessInformation,omitempty"`
	CorsConfiguration            interface{}           `json:"corsConfiguration,omitempty" yaml:"corsConfiguration,omitempty"`
	WorkflowStatus               string                `json:"workflowStatus,omitempty" yaml:"workflowStatus,omitempty"`
	ProtocolVersion              string                `json:"protocolVersion,omitempty" yaml:"protocolVersion,omitempty"`
	CreatedTime                  string                `json:"createdTime,omitempty" yaml:"createdTime,omitempty"`
	LastUpdatedTimestamp         string                `json:"lastUpdatedTimestamp,omitempty" yaml:"lastUpdatedTimestamp,omitempty"`
	LastUpdatedTime              string                `json:"lastUpdatedTime,omitempty" yaml:"lastUpdatedTime,omitempty"`
	SubtypeConfiguration         interface{}           `json:"subtypeConfiguration,omitempty" yaml:"subtypeConfiguration,omitempty"`
	Scopes                       []interface{}         `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	Operations                   []MCPServerOperations `json:"operations,omitempty" yaml:"operations,omitempty"`
	Categories                   []string              `json:"categories,omitempty" yaml:"categories,omitempty"`
	KeyManagers                  interface{}           `json:"keyManagers,omitempty" yaml:"keyManagers,omitempty"`
	GatewayVendor                string                `json:"gatewayVendor,omitempty" yaml:"gatewayVendor,omitempty"`
	GatewayType                  string                `json:"gatewayType,omitempty" yaml:"gatewayType,omitempty"`
	InitiatedFromGateway         bool                  `json:"initiatedFromGateway,omitempty" yaml:"initiatedFromGateway,omitempty"`
}

// MCPServerOperations represents operations available on a MCP Server
type MCPServerOperations struct {
	ID                  string            `json:"id,omitempty" yaml:"id,omitempty"`
	Target              string            `json:"target,omitempty" yaml:"target,omitempty"`
	Feature             string            `json:"feature,omitempty" yaml:"feature,omitempty"`
	AuthType            string            `json:"authType,omitempty" yaml:"authType,omitempty"`
	ThrottlingPolicy    string            `json:"throttlingPolicy,omitempty" yaml:"throttlingPolicy,omitempty"`
	Scopes              []string          `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	PayloadSchema       string            `json:"payloadSchema,omitempty" yaml:"payloadSchema,omitempty"`
	URIMapping          string            `json:"uriMapping,omitempty" yaml:"uriMapping,omitempty"`
	SchemaDefinition    string            `json:"schemaDefinition,omitempty" yaml:"schemaDefinition,omitempty"`
	Description         string            `json:"description,omitempty" yaml:"description,omitempty"`
	OperationPolicies   OperationPolicies `json:"operationPolicies,omitempty" yaml:"operationPolicies,omitempty"`
	APIOperationMapping interface{}       `json:"apiOperationMapping,omitempty" yaml:"apiOperationMapping,omitempty"`
}

// ByTargetVerbMCPServer implements sort.Interface based on the Target and Verb fields for MCPServerOperations.
type ByTargetVerbMCPServer []MCPServerOperations

// Len : Returns length of slice
func (a ByTargetVerbMCPServer) Len() int { return len(a) }

// Less : Compare two elements in slice for sorting
func (a ByTargetVerbMCPServer) Less(i, j int) bool {
	if a[i].Target < a[j].Target {
		return true
	}

	if a[i].Target > a[j].Target {
		return false
	}

	return a[i].ID < a[j].ID
}

// Swap : Swap two elements in slice
func (a ByTargetVerbMCPServer) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
