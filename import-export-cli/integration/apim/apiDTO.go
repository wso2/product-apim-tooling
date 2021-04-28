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

type APIFile struct {
	Type    string `json:"type" yaml:"type"`
	Version string `json:"version" yaml:"version"`
	Data    API    `json:"data" yaml:"data"`
}

// API : API DTO
type API struct {
	ID                              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name                            string            `json:"name,omitempty" yaml:"name,omitempty"`
	Description                     string            `json:"description,omitempty" yaml:"description,omitempty"`
	Context                         string            `json:"context,omitempty" yaml:"context,omitempty"`
	Version                         string            `json:"version,omitempty" yaml:"version,omitempty"`
	Provider                        string            `json:"provider,omitempty" yaml:"provider,omitempty"`
	LifeCycleStatus                 string            `json:"lifeCycleStatus,omitempty" yaml:"lifeCycleStatus,omitempty"`
	WsdlInfo                        interface{}       `json:"wsdlInfo,omitempty" yaml:"wsdlInfo,omitempty"`
	WsdlURL                         string            `json:"wsdlUrl,omitempty" yaml:"wsdlUrl,omitempty"`
	ResponseCachingEnabledKey       bool              `json:"responseCachingEnabled,omitempty" yaml:"responseCachingEnabled,omitempty"`
	CacheTimeout                    int               `json:"cacheTimeout,omitempty" yaml:"cacheTimeout,omitempty"`
	HasThumbnail                    bool              `json:"hasThumbnail,omitempty" yaml:"hasThumbnail,omitempty"`
	IsDefaultVersion                bool              `json:"isDefaultVersion,omitempty" yaml:"isDefaultVersion,omitempty"`
	IsRevision                      bool              `json:"isRevision" yaml:"isRevision"`
	RevisionID                      int32             `json:"revisionId" yaml:"revisionId"`
	EnableSchemaValidation          bool              `json:"enableSchemaValidation,omitempty" yaml:"enableSchemaValidation,omitempty"`
	Type                            string            `json:"type,omitempty" yaml:"type,omitempty"`
	Transport                       []string          `json:"transport,omitempty" yaml:"transport,omitempty"`
	Tags                            []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Policies                        []string          `json:"policies,omitempty" yaml:"policies,omitempty"`
	APIThrottlingPolicy             string            `json:"apiThrottlingPolicy,omitempty" yaml:"apiThrottlingPolicy,omitempty"`
	AuthorizationHeader             string            `json:"authorizationHeader,omitempty" yaml:"authorizationHeader,omitempty"`
	SecurityScheme                  []string          `json:"securityScheme,omitempty" yaml:"securityScheme,omitempty"`
	MaxTPS                          interface{}       `json:"maxTps,omitempty" yaml:"maxTps,omitempty"`
	Visibility                      string            `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	VisibleRoles                    []string          `json:"visibleRoles,omitempty" yaml:"visibleRoles,omitempty"`
	VisibleTenants                  []string          `json:"visibleTenants,omitempty" yaml:"visibleTenants,omitempty"`
	MediationPolicies               []MediationPolicy `json:"mediationPolicies,omitempty" yaml:"mediationPolicies,omitempty"`
	SubscriptionAvailability        string            `json:"subscriptionAvailability,omitempty" yaml:"subscriptionAvailability,omitempty"`
	SubscriptionAvailableTenants    []string          `json:"subscriptionAvailableTenants,omitempty" yaml:"subscriptionAvailableTenants,omitempty"`
	AdditionalProperties            []interface{}     `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Monetization                    interface{}       `json:"monetization,omitempty" yaml:"monetization,omitempty"`
	AccessControl                   string            `json:"accessControl,omitempty" yaml:"accessControl,omitempty"`
	AccessControlRoles              []string          `json:"accessControlRoles,omitempty" yaml:"accessControlRoles,omitempty"`
	BusinessInformation             interface{}       `json:"businessInformation,omitempty" yaml:"businessInformation,omitempty"`
	CorsConfiguration               interface{}       `json:"corsConfiguration,omitempty" yaml:"corsConfiguration,omitempty"`
	WorkflowStatus                  []string          `json:"workflowStatus,omitempty" yaml:"workflowStatus,omitempty"`
	CreatedTime                     string            `json:"createdTime,omitempty" yaml:"createdTime,omitempty"`
	LastUpdatedTime                 string            `json:"lastUpdatedTime,omitempty" yaml:"lastUpdatedTime,omitempty"`
	EndpointConfig                  interface{}       `json:"endpointConfig,omitempty" yaml:"endpointConfig,omitempty"`
	EndpointImplementationType      string            `json:"endpointImplementationType,omitempty" yaml:"endpointImplementationType,omitempty"`
	Scopes                          []interface{}     `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	Operations                      []APIOperations   `json:"operations,omitempty" yaml:"operations,omitempty"`
	ThreatProtectionPolicies        interface{}       `json:"threatProtectionPolicies,omitempty" yaml:"threatProtectionPolicies,omitempty"`
	Categories                      []string          `json:"categories,omitempty" yaml:"categories,omitempty"`
	KeyManagers                     []string          `json:"keyManagers,omitempty" yaml:"keyManagers,omitempty"`
	AdvertiseInformation            AdvertiseInfo     `json:"advertiseInfo,omitempty" yaml:"advertiseInfo,omitempty"`
	WebsubSubscriptionConfiguration interface{}       `json:"websubSubscriptionConfiguration" yaml:"websubSubscriptionConfiguration"`
}

// GetProductionURL : Get APIs production URL
func (instance *API) GetProductionURL() string {
	endpoint := instance.EndpointConfig.(map[string]interface{})["production_endpoints"]
	return endpoint.(map[string]interface{})["url"].(string)
}

// GetEndpointType : Get the Endpoint Type (http/address/aws)
func (instance *API) GetEndpointType() string {
	endpoint := instance.EndpointConfig.(map[string]interface{})["endpoint_type"].(string)
	return endpoint
}

// SetEndpointType : Get the Endpoint Type (http/address/aws)
func (instance *API) SetEndpointType(endPointType string) {
	instance.EndpointConfig.(map[string]interface{})["endpoint_type"] = endPointType
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

// GetEndPointConfig : Get APIs endpoint config
func (instance *API) GetEndPointConfig() interface{} {
	return instance.EndpointConfig
}

// SetEndPointConfig : Set APIs endpoint config
func (instance *API) SetEndPointConfig(config interface{}) {
	instance.EndpointConfig = config
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

// AdvertiseInfo : Advertise only information
type AdvertiseInfo struct {
	Advertised           bool   `json:"advertised" yaml:"advertised"`
	OriginalDevPortalUrl string `json:"originalDevPortalUrl,omitempty" yaml:"originalDevPortalUrl,omitempty"`
	ApiOwner             string `json:"apiOwner,omitempty" yaml:"apiOwner,omitempty"`
	Vendor               string `json:"vendor,omitempty" yaml:"vendor,omitempty"`
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
