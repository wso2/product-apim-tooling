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

package base

// API : API Model
type API struct {
	ID                           string                `json:"id"`
	Name                         string                `json:"name"`
	Description                  string                `json:"description"`
	Context                      string                `json:"context"`
	Version                      string                `json:"version"`
	Provider                     string                `json:"provider"`
	LifeCycleStatus              string                `json:"lifeCycleStatus"`
	WsdlInfo                     string                `json:"wsdlInfo"`
	ResponseCachingEnabled       bool                  `json:"responseCachingEnabled"`
	CacheTimeout                 int32                 `json:"cacheTimeout"`
	DestinationStatsEnabled      string                `json:"destinationStatsEnabled"`
	HasThumbnail                 bool                  `json:"hasThumbnail"`
	IsDefaultVersion             bool                  `json:"isDefaultVersion"`
	EnableSchemaValidation       bool                  `json:"enableSchemaValidation"`
	Type                         string                `json:"type"`
	Transport                    []string              `json:"transport"`
	Tags                         []string              `json:"tags"`
	Policies                     []string              `json:"policies"`
	APIThrottlingPolicy          string                `json:"apiThrottlingPolicy"`
	AuthorizationHeader          string                `json:"authorizationHeader"`
	SecurityScheme               []string              `json:"securityScheme"`
	MaxTps                       *APIMaxTps            `json:"maxTps"`
	Visibility                   string                `json:"visibility"`
	VisibleRoles                 []string              `json:"visibleRoles"`
	VisibleTenants               []string              `json:"visibleTenants"`
	EndpointSecurity             *APIEndpointSecurity  `json:"endpointSecurity"`
	GatewayEnvironments          []string              `json:"gatewayEnvironments"`
	Labels                       []string              `json:"labels"`
	MediationPolicies            []*MediationPolicy    `json:"mediationPolicies"`
	SubscriptionAvailability     string                `json:"subscriptionAvailability"`
	SubscriptionAvailableTenants []string              `json:"subscriptionAvailableTenants"`
	AdditionalProperties         map[string]string     `json:"additionalProperties"`
	Monetization                 *APIMonetization      `json:"monetization"`
	AccessControl                string                `json:"accessControl"`
	AccessControlRoles           []string              `json:"accessControlRoles"`
	BusinessInformation          *BusinessInfo         `json:"businessInformation"`
	CorsConfiguration            *APICorsConfiguration `json:"corsConfiguration"`
	WorkflowStatus               string                `json:"workflowStatus"`
	CreatedTime                  string                `json:"createdTime"`
	LastUpdatedTime              string                `json:"lastUpdatedTime"`
	EndpointConfig               *interface{}          `json:"endpointConfig"`
	EndpointImplementationType   string                `json:"endpointImplementationType"`
	Scopes                       *OAuthScopes          `json:"scopes"`
	Operations                   *APIOperations        `json:"operations"`
	ThreatProtectionPolicies     *interface{}          `json:"threatProtectionPolicies"`
}

// APIMaxTps : Defines Max TPS of backends
type APIMaxTps struct {
	Production int64 `json:"production"`
	Sandbox    int64 `json:"sandbox"`
}

// APIEndpointSecurity : API Endpoint Security config
type APIEndpointSecurity struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// MediationPolicy : Mediation Policy config
type MediationPolicy struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Shared bool   `json:"shared"`
}

// APIMonetization : API Monetization config
type APIMonetization struct {
	Enabled    bool              `json:"enabled"`
	Properties map[string]string `json:"properties"`
}

// BusinessInfo : Business Information
type BusinessInfo struct {
	BusinessOwner       string `json:"businessOwner"`
	BusinessOwnerEmail  string `json:"businessOwnerEmail"`
	TechnicalOwner      string `json:"technicalOwner"`
	TechnicalOwnerEmail string `json:"technicalOwnerEmail"`
}

// APICorsConfiguration : API CORS Configuration
type APICorsConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods"`
}

// OAuthScopes : OAuth Scope definition
type OAuthScopes struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Bindings    *RoleBindings `json:"bindings"`
}

// RoleBindings : Role bindings
type RoleBindings struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

// APIOperations : API Operations definition
type APIOperations struct {
	ID               string   `json:"id"`
	Target           string   `json:"target"`
	Verb             string   `json:"verb"`
	AuthType         string   `json:"authType"`
	ThrottlingPolicy string   `json:"throttlingPolicy"`
	Scopes           []string `json:"scopes"`
	UsedProductIds   []string `json:"usedProductIds"`
}
