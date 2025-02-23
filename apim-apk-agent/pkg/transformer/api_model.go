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

import "encoding/json"

// CustomParams holds the custom parameter values that has been enabled for the selected security mode
type CustomParams struct {
	CustomParamMapping map[string]string `json:"customParamMapping"`
}

// SecurityObj holds the idividual attribute values for each endpoint security config
type SecurityObj struct {
	Enabled          bool            `json:"enabled" yaml:"enabled"`
	EndpointUUID     string          `json:"endpointUUID" yaml:"endpointUUID"`
	Type             string          `json:"type" yaml:"type"`
	APIKeyValue      string          `json:"apiKeyValue" yaml:"apiKeyValue"`
	APIKeyIdentifier string          `json:"apiKeyIdentifier" yaml:"apiKeyIdentifier"`
	Username         string          `json:"username" yaml:"username"`
	Password         string          `json:"password" yaml:"password"`
	GrantType        string          `json:"grantType" yaml:"grantType"`
	TokenURL         string          `json:"tokenUrl" yaml:"tokenUrl"`
	ClientID         string          `json:"clientId" yaml:"clientId"`
	ClientSecret     string          `json:"clientSecret" yaml:"clientSecret"`
	CustomParameters json.RawMessage `json:"customParameters" yaml:"customParameters"`
}

// EndpointSecurityConfig holds security configs enabled for endpoints from the API level
type EndpointSecurityConfig struct {
	Production SecurityObj `json:"production" yaml:"production"`
	Sandbox    SecurityObj `json:"sandbox" yaml:"sandbox"`
}

// EndpointDetails represents the details of an endpoint, containing its URL.
type EndpointDetails struct {
	URL string `json:"url" yaml:"url"`
}

// EndpointConfig represents the configuration of an endpoint, including its type, sandbox, and production details.
type EndpointConfig struct {
	EndpointType        string                 `json:"endpoint_type" yaml:"endpoint_type"`
	SandboxEndpoints    EndpointDetails        `json:"sandbox_endpoints" yaml:"sandbox_endpoints"`
	ProductionEndpoints EndpointDetails        `json:"production_endpoints" yaml:"production_endpoints"`
	EndpointSecurity    EndpointSecurityConfig `json:"endpoint_security" yaml:"endpoint_security"`
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
	CORSConfiguration           CORSConfiguration     `yaml:"corsConfiguration"`
	EndpointConfig              EndpointConfig        `yaml:"endpointConfig"`
	PrimaryProductionEndpointID string                `yaml:"primaryProductionEndpointId"`
	PrimarySandboxEndpointID    string                `yaml:"primarySandboxEndpointId"`
	Operations                  []APIMOperation       `yaml:"operations"`
	OrganizationID              string                `yaml:"organizationId"`
	RevisionID                  uint32                `yaml:"revisionId"`
	RevisionedAPIID             string                `yaml:"revisionedApiId"`
	APIThrottlingPolicy         string                `yaml:"apiThrottlingPolicy"`
	APIPolicies                 APIMOperationPolicies `yaml:"apiPolicies"`
	SubtypeConfiguration        SubtypeConfiguration  `yaml:"subtypeConfiguration"`
	MaxTps                      *MaxTps               `yaml:"maxTps"`
}

// SubtypeConfiguration holds the details for Subtypes
type SubtypeConfiguration struct {
	Subtype       string `json:"subtype"`
	Configuration string `json:"_configuration"`
}

// Configuration holds the configuration details for the subtype
type Configuration struct {
	LLMProviderID string `json:"llmProviderId"`
}

// APIYaml is a wrapper struct for YAML representation of an API.
type APIYaml struct {
	Data APIMApi `json:"data"`
}

// EndpointsYaml is a wrapper struct for YAML representation of a list of endpoints.
type EndpointsYaml struct {
	Type    string     `json:"type"`
	Version string     `json:"version"`
	Data    []Endpoint `json:"data"`
}

// Endpoint represents an endpoint with its UUID, name, configuration, and deployment stage.
type Endpoint struct {
	ID              string         `json:"id" yaml:"id"`
	Name            string         `json:"name" yaml:"name"`
	EndpointConfig  EndpointConfig `json:"endpointConfig" yaml:"endpointConfig"`
	DeploymentStage string         `json:"deploymentStage" yaml:"deploymentStage"`
}

// MaxTps represents the maximum transactions per second (TPS) settings for both
// production and sandbox environments. It also includes an optional configuration
// for token-based throttling.
//
// Fields:
// - Production: Maximum TPS for the production environment.
// - ProductionTimeUnit: The time unit for the production TPS limit (e.g., seconds, minutes).
// - Sandbox: Maximum TPS for the sandbox environment.
// - SandboxTimeUnit: The time unit for the sandbox TPS limit.
// - TokenBasedThrottlingConfiguration: Configuration for token-based throttling.
type MaxTps struct {
	Production                        *int                        `yaml:"production"`
	ProductionTimeUnit                *string                     `yaml:"productionTimeUnit"`
	Sandbox                           *int                        `yaml:"sandbox"`
	SandboxTimeUnit                   *string                     `yaml:"sandboxTimeUnit"`
	TokenBasedThrottlingConfiguration *TokenBasedThrottlingConfig `yaml:"tokenBasedThrottlingConfiguration"`
}

// TokenBasedThrottlingConfig defines the token-based throttling limits for
// both production and sandbox environments. Token-based throttling places
// a limit on the number of prompt and completion tokens that can be used.
//
// Fields:
// - ProductionMaxPromptTokenCount: Maximum number of prompt tokens for production.
// - ProductionMaxCompletionTokenCount: Maximum number of completion tokens for production.
// - ProductionMaxTotalTokenCount: Maximum total token count (prompt + completion) for production.
// - SandboxMaxPromptTokenCount: Maximum number of prompt tokens for sandbox.
// - SandboxMaxCompletionTokenCount: Maximum number of completion tokens for sandbox.
// - SandboxMaxTotalTokenCount: Maximum total token count (prompt + completion) for sandbox.
// - IsTokenBasedThrottlingEnabled: Flag to enable or disable token-based throttling.
type TokenBasedThrottlingConfig struct {
	ProductionMaxPromptTokenCount     *int  `yaml:"productionMaxPromptTokenCount"`
	ProductionMaxCompletionTokenCount *int  `yaml:"productionMaxCompletionTokenCount"`
	ProductionMaxTotalTokenCount      *int  `yaml:"productionMaxTotalTokenCount"`
	SandboxMaxPromptTokenCount        *int  `yaml:"sandboxMaxPromptTokenCount"`
	SandboxMaxCompletionTokenCount    *int  `yaml:"sandboxMaxCompletionTokenCount"`
	SandboxMaxTotalTokenCount         *int  `yaml:"sandboxMaxTotalTokenCount"`
	IsTokenBasedThrottlingEnabled     *bool `yaml:"isTokenBasedThrottlingEnabled"`
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
	Endpoints            string               `json:"endpoints"`
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
	SecretData      []EndpointSecurityConfig
}

// ModelConfig holds the configuration details of a model
type ModelConfig struct {
	Model      string `json:"model"`
	EndpointID string `json:"endpointId"`
	Weight     int    `json:"weight"`
}

// Config holds the configuration details of the transformer
type Config struct {
	Production      []ModelConfig `json:"production"`
	Sandbox         []ModelConfig `json:"sandbox"`
	SuspendDuration string        `json:"suspendDuration"`
}

// ModelBasedRoundRobin holds the configuration details of the model based round robin
type ModelBasedRoundRobin struct {
	OnQuotaExceedSuspendDuration int              `yaml:"onQuotaExceedSuspendDuration"`
	ProductionModels             []ModelEndpoints `yaml:"productionModels"`
	SandboxModels                []ModelEndpoints `yaml:"sandboxModels"`
}

// ModelEndpoints holds the model and endpoint details
type ModelEndpoints struct {
	Model    string `yaml:"model"`
	Endpoint string `yaml:"endpoint"`
	Weight   int    `yaml:"weight"`
}

func (u ModelBasedRoundRobin) isParameter() {}
