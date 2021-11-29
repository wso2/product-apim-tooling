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

package v2

const (
	EpHttp        = "http"
	EpLoadbalance = "load_balance"
	EpFailover    = "failover"
)

// APIDefinition represents an API artifact in APIM
type APIDefinitionFile struct {
	Type        string           `json:"type,omitempty" yaml:"type,omitempty"`
	ApimVersion string           `json:"version,omitempty" yaml:"version,omitempty"`
	Data        APIDTODefinition `json:"data,omitempty" yaml:"data,omitempty"`
}

// APIDTODefinition represents an APIDTO artifact in APIM
type APIDTODefinition struct {
	ID                              string        `json:"id,omitempty" yaml:"id,omitempty"`
	Name                            string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description                     string        `json:"description,omitempty" yaml:"description,omitempty"`
	Context                         string        `json:"context,omitempty" yaml:"context,omitempty"`
	Version                         string        `json:"version,omitempty" yaml:"version,omitempty"`
	Provider                        string        `json:"provider,omitempty" yaml:"provider,omitempty"`
	LifeCycleStatus                 string        `json:"lifeCycleStatus,omitempty" yaml:"lifeCycleStatus,omitempty"`
	WsdlInfo                        interface{}   `json:"wsdlInfo,omitempty" yaml:"wsdlInfo,omitempty"`
	WsdlURL                         string        `json:"wsdlUrl,omitempty" yaml:"wsdlUrl,omitempty"`
	ResponseCachingEnabledKey       bool          `json:"responseCachingEnabled,omitempty" yaml:"responseCachingEnabled,omitempty"`
	CacheTimeout                    int           `json:"cacheTimeout,omitempty" yaml:"cacheTimeout,omitempty"`
	HasThumbnail                    bool          `json:"hasThumbnail,omitempty" yaml:"hasThumbnail,omitempty"`
	IsDefaultVersion                bool          `json:"isDefaultVersion,omitempty" yaml:"isDefaultVersion,omitempty"`
	IsRevision                      bool          `json:"isRevision" yaml:"isRevision"`
	RevisionID                      int32         `json:"revisionId" yaml:"revisionId"`
	EnableSchemaValidation          bool          `json:"enableSchemaValidation,omitempty" yaml:"enableSchemaValidation,omitempty"`
	Type                            string        `json:"type,omitempty" yaml:"type,omitempty"`
	Transport                       []string      `json:"transport,omitempty" yaml:"transport,omitempty"`
	Tags                            []string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	Policies                        []string      `json:"policies,omitempty" yaml:"policies,omitempty"`
	APIThrottlingPolicy             string        `json:"apiThrottlingPolicy,omitempty" yaml:"apiThrottlingPolicy,omitempty"`
	AuthorizationHeader             string        `json:"authorizationHeader,omitempty" yaml:"authorizationHeader,omitempty"`
	SecurityScheme                  []string      `json:"securityScheme,omitempty" yaml:"securityScheme,omitempty"`
	MaxTPS                          interface{}   `json:"maxTps,omitempty" yaml:"maxTps,omitempty"`
	Visibility                      string        `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	VisibleRoles                    []string      `json:"visibleRoles,omitempty" yaml:"visibleRoles,omitempty"`
	VisibleTenants                  []string      `json:"visibleTenants,omitempty" yaml:"visibleTenants,omitempty"`
	MediationPolicies               []interface{} `json:"mediationPolicies,omitempty" yaml:"mediationPolicies,omitempty"`
	SubscriptionAvailability        string        `json:"subscriptionAvailability,omitempty" yaml:"subscriptionAvailability,omitempty"`
	SubscriptionAvailableTenants    []string      `json:"subscriptionAvailableTenants,omitempty" yaml:"subscriptionAvailableTenants,omitempty"`
	AdditionalProperties            []interface{} `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Monetization                    interface{}   `json:"monetization,omitempty" yaml:"monetization,omitempty"`
	AccessControl                   string        `json:"accessControl,omitempty" yaml:"accessControl,omitempty"`
	AccessControlRoles              []string      `json:"accessControlRoles,omitempty" yaml:"accessControlRoles,omitempty"`
	BusinessInformation             interface{}   `json:"businessInformation,omitempty" yaml:"businessInformation,omitempty"`
	CorsConfiguration               interface{}   `json:"corsConfiguration,omitempty" yaml:"corsConfiguration,omitempty"`
	WorkflowStatus                  []string      `json:"workflowStatus,omitempty" yaml:"workflowStatus,omitempty"`
	CreatedTime                     string        `json:"createdTime,omitempty" yaml:"createdTime,omitempty"`
	LastUpdatedTime                 string        `json:"lastUpdatedTime,omitempty" yaml:"lastUpdatedTime,omitempty"`
	EndpointConfig                  interface{}   `json:"endpointConfig,omitempty" yaml:"endpointConfig,omitempty"`
	EndpointImplementationType      string        `json:"endpointImplementationType,omitempty" yaml:"endpointImplementationType,omitempty"`
	Scopes                          []interface{} `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	Operations                      []interface{} `json:"operations,omitempty" yaml:"operations,omitempty"`
	ThreatProtectionPolicies        interface{}   `json:"threatProtectionPolicies,omitempty" yaml:"threatProtectionPolicies,omitempty"`
	Categories                      []string      `json:"categories,omitempty" yaml:"categories,omitempty"`
	KeyManagers                     []string      `json:"keyManagers,omitempty" yaml:"keyManagers,omitempty"`
	AdvertiseInformation            AdvertiseInfo `json:"advertiseInfo,omitempty" yaml:"advertiseInfo,omitempty"`
	WebsubSubscriptionConfiguration interface{}   `json:"websubSubscriptionConfiguration" yaml:"websubSubscriptionConfiguration"`
	GatewayVendor                   string        `json:"gatewayVendor,omitempty" yaml:"gatewayVendor,omitempty"`
	AsyncTransportProtocols         []string      `json:"asyncTransportProtocols,omitempty" yaml:"asyncTransportProtocols,omitempty"`
}

type CorsConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled,omitempty" yaml:"corsConfigurationEnabled,omitempty"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty" yaml:"accessControlAllowOrigins,omitempty"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty" yaml:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty" yaml:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty" yaml:"accessControlAllowMethods,omitempty"`
}
type Document struct {
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Data    Data   `json:"data,omitempty" yaml:"data,omitempty"`
}
type Data struct {
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	Type          string `json:"type,omitempty" yaml:"type,omitempty"`
	Summary       string `json:"summary,omitempty" yaml:"summary,omitempty"`
	SourceType    string `json:"sourceType,omitempty" yaml:"sourceType,omitempty"`
	OtherTypeName string `json:"otherTypeName,omitempty" yaml:"otherTypeName,omitempty"`
	Visibility    string `json:"visibility,omitempty" yaml:"visibility,omitempty"`
}

// AdvertiseInfo : Advertise only information
type AdvertiseInfo struct {
	Advertised           bool   `json:"advertised,omitempty" yaml:"advertised,omitempty"`
	OriginalDevPortalUrl string `json:"originalDevPortalUrl,omitempty" yaml:"originalDevPortalUrl,omitempty"`
	ApiOwner             string `json:"apiOwner,omitempty" yaml:"apiOwner,omitempty"`
	Vendor               string `json:"vendor,omitempty" yaml:"vendor,omitempty"`
}
