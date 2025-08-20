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

package apim

// MCPServer : MCPServer DTO
type MCPServer struct {
	ID                              string                `json:"id,omitempty" yaml:"id,omitempty"`
	Name                            string                `json:"name,omitempty" yaml:"name,omitempty"`
	DisplayName                     string                `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Description                     string                `json:"description,omitempty" yaml:"description,omitempty"`
	Context                         string                `json:"context,omitempty" yaml:"context,omitempty"`
	Version                         string                `json:"version,omitempty" yaml:"version,omitempty"`
	Provider                        string                `json:"provider,omitempty" yaml:"provider,omitempty"`
	LifeCycleStatus                 string                `json:"lifeCycleStatus,omitempty" yaml:"lifeCycleStatus,omitempty"`
	WsdlInfo                        interface{}           `json:"wsdlInfo,omitempty" yaml:"wsdlInfo,omitempty"`
	WsdlURL                         string                `json:"wsdlUrl,omitempty" yaml:"wsdlUrl,omitempty"`
	ResponseCachingEnabledKey       bool                  `json:"responseCachingEnabled,omitempty" yaml:"responseCachingEnabled,omitempty"`
	CacheTimeout                    int                   `json:"cacheTimeout,omitempty" yaml:"cacheTimeout,omitempty"`
	HasThumbnail                    bool                  `json:"hasThumbnail,omitempty" yaml:"hasThumbnail,omitempty"`
	IsDefaultVersion                bool                  `json:"isDefaultVersion,omitempty" yaml:"isDefaultVersion,omitempty"`
	IsRevision                      bool                  `json:"isRevision" yaml:"isRevision"`
	RevisionID                      int                   `json:"revisionId" yaml:"revisionId"`
	EnableSchemaValidation          bool                  `json:"enableSchemaValidation,omitempty" yaml:"enableSchemaValidation,omitempty"`
	Type                            string                `json:"type,omitempty" yaml:"type,omitempty"`
	Transport                       []string              `json:"transport,omitempty" yaml:"transport,omitempty"`
	Tags                            []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Policies                        []string              `json:"policies,omitempty" yaml:"policies,omitempty"`
	APIThrottlingPolicy             string                `json:"apiThrottlingPolicy,omitempty" yaml:"apiThrottlingPolicy,omitempty"`
	AuthorizationHeader             string                `json:"authorizationHeader,omitempty" yaml:"authorizationHeader,omitempty"`
	ApiKeyHeader                    string                `json:"apiKeyHeader,omitempty" yaml:"apiKeyHeader,omitempty"`
	SecurityScheme                  []string              `json:"securityScheme,omitempty" yaml:"securityScheme,omitempty"`
	MaxTPS                          interface{}           `json:"maxTps,omitempty" yaml:"maxTps,omitempty"`
	Visibility                      string                `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	VisibleRoles                    []string              `json:"visibleRoles,omitempty" yaml:"visibleRoles,omitempty"`
	VisibleTenants                  []string              `json:"visibleTenants,omitempty" yaml:"visibleTenants,omitempty"`
	MediationPolicies               []MediationPolicy     `json:"mediationPolicies,omitempty" yaml:"mediationPolicies,omitempty"`
	SubscriptionAvailability        string                `json:"subscriptionAvailability,omitempty" yaml:"subscriptionAvailability,omitempty"`
	SubscriptionAvailableTenants    []string              `json:"subscriptionAvailableTenants,omitempty" yaml:"subscriptionAvailableTenants,omitempty"`
	AdditionalProperties            []interface{}         `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Monetization                    interface{}           `json:"monetization,omitempty" yaml:"monetization,omitempty"`
	AccessControl                   string                `json:"accessControl,omitempty" yaml:"accessControl,omitempty"`
	AccessControlRoles              []string              `json:"accessControlRoles,omitempty" yaml:"accessControlRoles,omitempty"`
	BusinessInformation             interface{}           `json:"businessInformation,omitempty" yaml:"businessInformation,omitempty"`
	CorsConfiguration               interface{}           `json:"corsConfiguration,omitempty" yaml:"corsConfiguration,omitempty"`
	WorkflowStatus                  []string              `json:"workflowStatus,omitempty" yaml:"workflowStatus,omitempty"`
	CreatedTime                     string                `json:"createdTime,omitempty" yaml:"createdTime,omitempty"`
	LastUpdatedTimestamp            string                `json:"lastUpdatedTimestamp,omitempty" yaml:"lastUpdatedTimestamp,omitempty"`
	LastUpdatedTime                 string                `json:"lastUpdatedTime,omitempty" yaml:"lastUpdatedTime,omitempty"`
	EndpointConfig                  interface{}           `json:"endpointConfig,omitempty" yaml:"endpointConfig,omitempty"`
	EndpointImplementationType      string                `json:"endpointImplementationType,omitempty" yaml:"endpointImplementationType,omitempty"`
	Scopes                          []interface{}         `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	Operations                      []MCPServerOperations `json:"operations,omitempty" yaml:"operations,omitempty"`
	ThreatProtectionPolicies        interface{}           `json:"threatProtectionPolicies,omitempty" yaml:"threatProtectionPolicies,omitempty"`
	Categories                      []string              `json:"categories,omitempty" yaml:"categories,omitempty"`
	KeyManagers                     []string              `json:"keyManagers,omitempty" yaml:"keyManagers,omitempty"`
	AdvertiseInformation            AdvertiseInfo         `json:"advertiseInfo,omitempty" yaml:"advertiseInfo,omitempty"`
	WebsubSubscriptionConfiguration interface{}           `json:"websubSubscriptionConfiguration" yaml:"websubSubscriptionConfiguration"`
	GatewayVendor                   string                `json:"gatewayVendor,omitempty" yaml:"gatewayVendor,omitempty"`
	AsyncTransportProtocols         []string              `json:"asyncTransportProtocols,omitempty" yaml:"asyncTransportProtocols,omitempty"`
	GatewayType                     string                `json:"gatewayType,omitempty" yaml:"gatewayType,omitempty"`
	InitiatedFromGateway            bool                  `json:"initiatedFromGateway,omitempty" yaml:"initiatedFromGateway,omitempty"`
	EnableSubscriberVerification    bool                  `json:"enableSubscriberVerification,omitempty" yaml:"enableSubscriberVerification,omitempty"`
}

// MCPServerOperations represents operations available on a MCP Server
type MCPServerOperations struct {
	ID                string            `json:"id" yaml:"id"`
	Target            string            `json:"target" yaml:"target"`
	Verb              string            `json:"verb" yaml:"verb"`
	Feature           string            `json:"feature,omitempty" yaml:"feature,omitempty"`
	SchemaDefinition  interface{}       `json:"schemaDefinition,omitempty" yaml:"schemaDefinition,omitempty"`
	Description       string            `json:"description,omitempty" yaml:"description,omitempty"`
	AuthType          string            `json:"authType,omitempty" yaml:"authType,omitempty"`
	ThrottlingPolicy  string            `json:"throttlingPolicy,omitempty" yaml:"throttlingPolicy,omitempty"`
	Scopes            []string          `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	UsedProductIds    []string          `json:"usedProductIds,omitempty" yaml:"usedProductIds,omitempty"`
	OperationPolicies OperationPolicies `json:"operationPolicies,omitempty" yaml:"operationPolicies,omitempty"`
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

	return a[i].Verb < a[j].Verb
}

// Swap : Swap two elements in slice
func (a ByTargetVerbMCPServer) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
