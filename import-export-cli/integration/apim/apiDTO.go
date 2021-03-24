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

// API : API DTO
type API struct {
	ID                              string               `json:"id"`
	Name                            string               `json:"name"`
	Description                     string               `json:"description"`
	Context                         string               `json:"context"`
	Version                         string               `json:"version"`
	Provider                        string               `json:"provider"`
	LifeCycleStatus                 string               `json:"lifeCycleStatus"`
	ResponseCachingEnabled          bool                 `json:"responseCachingEnabled"`
	CacheTimeout                    int32                `json:"cacheTimeout"`
	DestinationStatsEnabled         string               `json:"destinationStatsEnabled"`
	HasThumbnail                    bool                 `json:"hasThumbnail"`
	IsDefaultVersion                bool                 `json:"isDefaultVersion"`
	IsRevision                      bool                 `json:"isRevision"`
	RevisionID                      int32                `json:"revisionId"`
	EnableSchemaValidation          bool                 `json:"enableSchemaValidation"`
	Type                            string               `json:"type"`
	Transport                       []string             `json:"transport"`
	Tags                            []string             `json:"tags"`
	Policies                        []string             `json:"policies"`
	APIThrottlingPolicy             string               `json:"apiThrottlingPolicy"`
	AuthorizationHeader             string               `json:"authorizationHeader"`
	SecurityScheme                  []string             `json:"securityScheme"`
	MaxTps                          APIMaxTps            `json:"maxTps,omitempty"`
	Visibility                      string               `json:"visibility"`
	VisibleRoles                    []string             `json:"visibleRoles"`
	VisibleTenants                  []string             `json:"visibleTenants"`
	GatewayEnvironments             []string             `json:"gatewayEnvironments"`
	MediationPolicies               []MediationPolicy    `json:"mediationPolicies,omitempty"`
	SubscriptionAvailability        string               `json:"subscriptionAvailability"`
	SubscriptionAvailableTenants    []string             `json:"subscriptionAvailableTenants"`
	AdditionalProperties            map[string]string    `json:"additionalProperties"`
	Monetization                    APIMonetization      `json:"monetization,omitempty"`
	AccessControl                   string               `json:"accessControl"`
	AccessControlRoles              []string             `json:"accessControlRoles"`
	BusinessInformation             BusinessInfo         `json:"businessInformation,omitempty"`
	CorsConfiguration               APICorsConfiguration `json:"corsConfiguration,omitempty"`
	WorkflowStatus                  string               `json:"workflowStatus"`
	CreatedTime                     string               `json:"createdTime"`
	LastUpdatedTime                 string               `json:"lastUpdatedTime"`
	EndpointConfig                  interface{}          `json:"endpointConfig"`
	EndpointImplementationType      string               `json:"endpointImplementationType"`
	Scopes                          []OAuthScopes        `json:"scopes,omitempty"`
	Operations                      []APIOperations      `json:"operations"`
	ThreatProtectionPolicies        interface{}          `json:"threatProtectionPolicies"`
	WebsubSubscriptionConfiguration interface{}          `json:"websubSubscriptionConfiguration"`
}

// GetProductionURL : Get APIs production URL
func (instance *API) GetProductionURL() string {
	endpoint := instance.EndpointConfig.(map[string]interface{})["production_endpoints"]
	return endpoint.(map[string]interface{})["url"].(string)
}

// GetSandboxURL : Get APIs sandbox URL
func (instance *API) GetSandboxURL() string {
	endpoint := instance.EndpointConfig.(map[string]interface{})["sandbox_endpoints"]
	return endpoint.(map[string]interface{})["url"].(string)
}

// SetProductionURL : Set APIs production URL
func (instance *API) SetProductionURL(url string) {
	endpoint := instance.EndpointConfig.(map[string]interface{})["production_endpoints"]
	endpoint.(map[string]interface{})["url"] = url
}

// SetSandboxURL : Set APIs sandbox URL
func (instance *API) SetSandboxURL(url string) {
	endpoint := instance.EndpointConfig.(map[string]interface{})["sandbox_endpoints"]
	endpoint.(map[string]interface{})["url"] = url
}

// GetProductionConfig : Get APIs production Config
func (instance *API) GetProductionConfig() map[string]interface{} {
	endpoint := instance.EndpointConfig.(map[string]interface{})["production_endpoints"]
	return endpoint.(map[string]interface{})["config"].(map[string]interface{})
}

// GetSandboxConfig : Get APIs sandbox Config
func (instance *API) GetSandboxConfig() map[string]interface{} {
	endpoint := instance.EndpointConfig.(map[string]interface{})["sandbox_endpoints"]
	return endpoint.(map[string]interface{})["config"].(map[string]interface{})
}

// SetProductionConfig : Set APIs production Config
func (instance *API) SetProductionConfig(config map[interface{}]interface{}) {
	endpoint := instance.EndpointConfig.(map[string]interface{})["production_endpoints"]
	endpoint.(map[string]interface{})["config"] = config
}

// SetSandboxConfig : Set APIs sandbox Config
func (instance *API) SetSandboxConfig(config map[interface{}]interface{}) {
	endpoint := instance.EndpointConfig.(map[string]interface{})["sandbox_endpoints"]
	endpoint.(map[string]interface{})["config"] = config
}

// GetProductionSecurityConfig : Get APIs production security config
func (instance *API) GetProductionSecurityConfig() map[string]interface{} {
	endpoint := instance.EndpointConfig.(map[string]interface{})["endpoint_security"]
	return endpoint.(map[string]interface{})["production"].(map[string]interface{})
}

// GetSandboxSecurityConfig : Get APIs sandbox security config
func (instance *API) GetSandboxSecurityConfig() map[string]interface{} {
	endpoint := instance.EndpointConfig.(map[string]interface{})["endpoint_security"]
	return endpoint.(map[string]interface{})["sandbox"].(map[string]interface{})
}

// APIMaxTps : Defines Max TPS of backends
type APIMaxTps struct {
	Production int64 `json:"production,omitempty"`
	Sandbox    int64 `json:"sandbox,omitempty"`
}

// APIEndpointSecurity : API Endpoint Security config
type APIEndpointSecurity struct {
	Type     string `json:"type,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// MediationPolicy : Mediation Policy config
type MediationPolicy struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Shared bool   `json:"shared,omitempty"`
}

// ByID implements sort.Interface based on the ID field.
type ByID []MediationPolicy

// Len : Returns lendth of slice
func (a ByID) Len() int { return len(a) }

// Less : Compare two elements in slice for sorting
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Swap : Swap two elements in slice
func (a ByID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// APIMonetization : API Monetization config
type APIMonetization struct {
	Enabled    bool              `json:"enabled,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// BusinessInfo : Business Information
type BusinessInfo struct {
	BusinessOwner       string `json:"businessOwner,omitempty"`
	BusinessOwnerEmail  string `json:"businessOwnerEmail,omitempty"`
	TechnicalOwner      string `json:"technicalOwner,omitempty"`
	TechnicalOwnerEmail string `json:"technicalOwnerEmail,omitempty"`
}

// APICorsConfiguration : API CORS Configuration
type APICorsConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled,omitempty"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty"`
}

// OAuthScopes : OAuth Scope definition
type OAuthScopes struct {
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Bindings    RoleBindings `json:"bindings,omitempty"`
}

// RoleBindings : Role bindings
type RoleBindings struct {
	Type   string   `json:"type,omitempty"`
	Values []string `json:"values,omitempty"`
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

// ByTargetVerb implements sort.Interface based on the Target and Verb fields.
type ByTargetVerb []APIOperations

// Len : Returns lendth of slice
func (a ByTargetVerb) Len() int { return len(a) }

// Less : Compare two elements in slice for sorting
func (a ByTargetVerb) Less(i, j int) bool {
	if a[i].Target < a[j].Target {
		return true
	}

	if a[i].Target > a[j].Target {
		return false
	}

	return a[i].Verb < a[j].Verb
}

// Swap : Swap two elements in slice
func (a ByTargetVerb) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// HTTPEndpoint : HTTP Endpoint definition
type HTTPEndpoint struct {
	EndpointType        string     `json:"endpoint_type"`
	SandboxEndpoints    *URLConfig `json:"sandbox_endpoints"`
	ProductionEndpoints *URLConfig `json:"production_endpoints"`
}

// URLConfig : URL Configuration
type URLConfig struct {
	URL string `json:"url"`
}
