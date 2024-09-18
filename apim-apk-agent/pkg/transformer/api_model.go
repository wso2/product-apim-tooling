/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package transformer

// CustomParams holds the custom parameter values that has been enabled for the selected security mode
type CustomParams struct {
	CustomParamMapping map[string]string `json:"customParamMapping"`
}

// SecurityObj holds the idividual attribute values for each endpoint security config
type SecurityObj struct {
	Enabled          bool                   `json:"enabled"`
	Type             string                 `json:"type"`
	Username         string                 `json:"username"`
	Password         string                 `json:"password"`
	GrantType        string                 `json:"grantType"`
	TokenURL         string                 `json:"tokenUrl"`
	ClientID         string                 `json:"clientId"`
	ClientSecret     string                 `json:"clientSecret"`
	CustomParameters map[string]interface{} `json:"customParameters"`
}

// EndpointSecurityConfig holds security configs enabled for endpoints from the API level
type EndpointSecurityConfig struct {
	Production SecurityObj `json:"production"`
	Sandbox    SecurityObj `json:"sandbox"`
}

// EndpointDetails represents the details of an endpoint, containing its URL.
type EndpointDetails struct {
	URL string `json:"url"`
}

// EndpointConfig represents the configuration of an endpoint, including its type, sandbox, and production details.
type EndpointConfig struct {
	EndpointType        string                 `json:"endpoint_type"`
	SandboxEndpoints    EndpointDetails        `json:"sandbox_endpoints"`
	ProductionEndpoints EndpointDetails        `json:"production_endpoints"`
	EndpointSecurity    EndpointSecurityConfig `json:"endpoint_security"`
}

// CORSConfiguration represents the CORS (Cross-Origin Resource Sharing) configuration for an API.
type CORSConfiguration struct {
	CORSConfigurationEnabled      bool     `yaml:"corsConfigurationEnabled"`
	AccessControlAllowOrigins     []string `yaml:"accessControlAllowOrigins"`
	AccessControlAllowCredentials bool     `yaml:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     []string `yaml:"accessControlAllowHeaders"`
	AccessControlAllowMethods     []string `yaml:"accessControlAllowMethods"`
}

// AdditionalProperties represents additional properties for an API in the form of a map.
type AdditionalProperties struct {
	Name               string `yaml:"name"`
	Value              string `yaml:"value"`
	DisplayInDevPortal bool   `yaml:"display"`
}

// OperationPolicy defines policies, including interceptor parameters, for API operations.
type OperationPolicy struct {
	PolicyName    string    `yaml:"policyName,omitempty"`
	PolicyVersion string    `yaml:"policyVersion,omitempty"`
	PolicyID      string    `yaml:"policyId,omitempty"`
	Parameters    Parameter `yaml:"parameters,omitempty"`
}

// Parameter interface is used to define the type of parameters that can be used in an operation policy.
type Parameter interface {
	isParameter()
}

// RedirectPolicy contains the information for redirect request policies
type RedirectPolicy struct {
	URL        string `json:"url,omitempty" yaml:"url,omitempty"`
	StatusCode int    `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
}

func (u RedirectPolicy) isParameter() {}

// URLList contains the urls for mirror policies
type URLList struct {
	URLs []string `json:"urls,omitempty" yaml:"urls,omitempty"`
}

func (u URLList) isParameter() {}

// Header contains the information for header modification
type Header struct {
	HeaderName  string `yaml:"headerName"`
	HeaderValue string `yaml:"headerValue,omitempty"`
}

func (h Header) isParameter() {}

// InterceptorService holds configuration details for configuring interceptor
// for particular API requests or responses.
type InterceptorService struct {
	BackendURL      string `yaml:"backendUrl,omitempty"`
	HeadersEnabled  bool   `yaml:"headersEnabled,omitempty"`
	BodyEnabled     bool   `yaml:"bodyEnabled,omitempty"`
	TrailersEnabled bool   `yaml:"trailersEnabled,omitempty"`
	ContextEnabled  bool   `yaml:"contextEnabled,omitempty"`
	TLSSecretName   string `yaml:"tlsSecretName,omitempty"`
	TLSSecretKey    string `yaml:"tlsSecretKey,omitempty"`
}

func (s InterceptorService) isParameter() {}

// BackendJWT holds configuration details for configuring JWT for backend
type BackendJWT struct {
	Encoding         string `yaml:"encoding,omitempty"`
	Header           string `yaml:"header,omitempty"`
	SigningAlgorithm string `yaml:"signingAlgorithm,omitempty"`
	TokenTTL         int    `yaml:"tokenTTL,omitempty"`
}

func (j BackendJWT) isParameter() {}

// APIMOperationPolicy defines policies, including interceptor parameters, for API operations.
type APIMOperationPolicy struct {
	PolicyName    string                 `yaml:"policyName,omitempty"`
	PolicyVersion string                 `yaml:"policyVersion,omitempty"`
	PolicyID      string                 `yaml:"policyId,omitempty"`
	Parameters    map[string]interface{} `yaml:"parameters,omitempty"`
}

// APIMOperationPolicies organizes request, response, and fault policies for an API operation.
type APIMOperationPolicies struct {
	Request  []APIMOperationPolicy `yaml:"request"`
	Response []APIMOperationPolicy `yaml:"response"`
	Fault    []APIMOperationPolicy `yaml:"fault"`
}

// APIMOperation represents an API operation with its target, verb, scopes, and associated policies.
type APIMOperation struct {
	Target            string                 `yaml:"target"`
	Verb              string                 `yaml:"verb"`
	Scopes            []string               `yaml:"scopes"`
	OperationPolicies *APIMOperationPolicies `yaml:"operationPolicies"`
	ThrottlingPolicy  string                 `yaml:"throttlingPolicy"`
	AuthType          string                 `yaml:"authType"`
}

// APIMApi represents an API along with it's all basic information and the operations.
type APIMApi struct {
	ID                   string                 `yaml:"id"`
	Name                 string                 `yaml:"name"`
	Version              string                 `yaml:"version"`
	Context              string                 `yaml:"context"`
	DefaultVersion       bool                   `json:"isDefaultVersion"`
	Type                 string                 `yaml:"type"`
	AuthorizationHeader  string                 `yaml:"authorizationHeader"`
	APIKeyHeader         string                 `yaml:"apiKeyHeader"`
	SecuritySchemes      []string               `json:"securityScheme"`
	AdditionalProperties []AdditionalProperties `yaml:"additionalProperties"`
	// AdditionalPropertiesMap []AdditionalPropertiesMap `yaml:"additionalPropertiesMap"`
	CORSConfiguration   CORSConfiguration     `yaml:"corsConfiguration"`
	EndpointConfig      EndpointConfig        `yaml:"endpointConfig"`
	Operations          []APIMOperation       `yaml:"operations"`
	OrganizationID      string                `yaml:"organizationId"`
	RevisionID          uint32                `yaml:"revisionId"`
	RevisionedAPIID     string                `yaml:"revisionedApiId"`
	APIThrottlingPolicy string                `yaml:"apiThrottlingPolicy"`
	APIPolicies         APIMOperationPolicies `yaml:"apiPolicies"`
}

// APIYaml is a wrapper struct for YAML representation of an API.
type APIYaml struct {
	Data APIMApi `json:"data"`
}

// APIArtifact represents the artifact details of an API, including api details, environment configuration,
// Swagger definition, deployment descriptor, and revision ID extracted from the API Project Zip.
type APIArtifact struct {
	APIJson              string               `json:"apiJson"`
	APIFileName          string               `json:"apiFileName"`
	EnvConfig            string               `json:"envConfig"`
	Schema               string               `json:"schema"`
	DeploymentDescriptor string               `json:"deploymentDescriptor"`
	CertArtifact         CertificateArtifact  `json:"certArtifact"`
	RevisionID           uint32               `json:"revisionId"`
	CertMeta             CertMetadata         `json:"certMeta"`
	EndpointCertMeta     EndpointCertMetadata `json:"endpintCertMeta"`
}

// CertificateArtifact stores the parsed file content created inside the API project zip upon enabling certificate aided security options
type CertificateArtifact struct {
	ClientCerts   string `json:"clientCert"`
	EndpointCerts string `json:"endpointCert"`
}

// CertMetadata marks the availability of the cert files provided by the client and their contents
type CertMetadata struct {
	CertAvailable   bool              `json:"certAvailable"`
	ClientCertFiles map[string]string `json:"clientCertFiles"`
}

// EndpointCertMetadata marks the availability of the endpoint certificates and stores the cert contents
type EndpointCertMetadata struct {
	CertAvailable     bool              `json:"certAvailable"`
	EndpointCertFiles map[string]string `json:"endpointCertFiles"`
}

// CertContainer acts as a wrapper to hold onto all the certificate details for both endpoint and client-side security configs
// belong to a particular API Project
type CertContainer struct {
	ClientCertObj   CertMetadata
	EndpointCertObj EndpointCertMetadata
	SecretData      EndpointSecurityConfig
}
