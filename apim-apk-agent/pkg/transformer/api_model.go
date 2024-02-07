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

// EndpointDetails represents the details of an endpoint, containing its URL.
type EndpointDetails struct {
	URL string `json:"url"`
}

// EndpointConfig represents the configuration of an endpoint, including its type, sandbox, and production details.
type EndpointConfig struct {
	EndpointType        string          `json:"endpoint_type"`
	SandboxEndpoints    EndpointDetails `json:"sandbox_endpoints"`
	ProductionEndpoints EndpointDetails `json:"production_endpoints"`
}

// CORSConfiguration represents the CORS (Cross-Origin Resource Sharing) configuration for an API.
type CORSConfiguration struct {
	CORSConfigurationEnabled      bool     `yaml:"corsConfigurationEnabled"`
	AccessControlAllowOrigins     []string `yaml:"accessControlAllowOrigins"`
	AccessControlAllowCredentials bool     `yaml:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     []string `yaml:"accessControlAllowHeaders"`
	AccessControlAllowMethods     []string `yaml:"accessControlAllowMethods"`
}

// AdditionalPropertiesMap represents additional properties for an API in the form of a map.
type AdditionalPropertiesMap struct{}

// InterceptorService holds configuration details for configuring interceptor
// for a aperticular API requests or responses.
type InterceptorService struct {
	BackendURL      string `yaml:"backendUrl,omitempty"`
	HeadersEnabled  bool   `yaml:"headersEnabled,omitempty"`
	BodyEnabled     bool   `yaml:"bodyEnabled,omitempty"`
	TrailersEnabled bool   `yaml:"trailersEnabled,omitempty"`
	ContextEnabled  bool   `yaml:"contextEnabled,omitempty"`
	TLSSecretName   string `yaml:"tlsSecretName,omitempty"`
	TLSSecretKey    string `yaml:"tlsSecretKey,omitempty"`
}

// OperationPolicy defines policies, including interceptor parameters, for API operations.
type OperationPolicy struct {
	PolicyName    string              `yaml:"policyName,omitempty"`
	PolicyVersion string              `yaml:"policyVersion,omitempty"`
	PolicyID      string              `yaml:"policyId,omitempty"`
	Parameters    *InterceptorService `yaml:"parameters,omitempty"`
}

// APIMOperationPolicies organizes request, response, and fault policies for an API operation.
type APIMOperationPolicies struct {
	Request  []OperationPolicy `yaml:"request"`
	Response []OperationPolicy `yaml:"response"`
	Fault    []OperationPolicy `yaml:"fault"`
}

// APIMOperation represents an API operation with its target, verb, scopes, and associated policies.
type APIMOperation struct {
	Target            string                 `yaml:"target"`
	Verb              string                 `yaml:"verb"`
	Scopes            []string               `yaml:"scopes"`
	OperationPolicies *APIMOperationPolicies `yaml:"operationPolicies"`
}

// APIMApi represents an API along with it's all basic information and the operations.
type APIMApi struct {
	ID                      string                  `yaml:"id"`
	Name                    string                  `yaml:"name"`
	Version                 string                  `yaml:"version"`
	Context                 string                  `yaml:"context"`
	DefaultVersion          bool                    `yaml:"isDefaultVersion"`
	Type                    string                  `yaml:"type"`
	AuthorizationHeader     string                  `yaml:"authorizationHeader"`
	SecuritySchemes         []string                `json:"securityScheme"`
	AdditionalProperties    []string                `yaml:"additionalProperties"`
	AdditionalPropertiesMap AdditionalPropertiesMap `yaml:"additionalPropertiesMap"`
	CORSConfiguration       CORSConfiguration       `yaml:"corsConfiguration"`
	EndpointConfig          EndpointConfig          `yaml:"endpointConfig"`
	Operations              []APIMOperation         `yaml:"operations"`
	OrganizationID          string                  `yaml:"organizationId"`
	RevisionID              uint32                  `yaml:"revisionId"`
}

// APIYaml is a wrapper struct for YAML representation of an API.
type APIYaml struct {
	Data APIMApi `json:"data"`
}

// APIArtifact represents the artifact details of an API, including api details, environment configuration,
// Swagger definition, deployment descriptor, and revision ID extracted from the API Project Zip.
type APIArtifact struct {
	APIJson              string `json:"apiJson"`
	APIFileName          string `json:"apiFileName"`
	EnvConfig            string `json:"envConfig"`
	Swagger              string `json:"swagger"`
	DeploymentDescriptor string `json:"deploymentDescriptor"`
	RevisionID           uint32 `json:"revisionId"`
}
