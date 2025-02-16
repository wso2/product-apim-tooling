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

/*
 * Package "transformer" contains functions related to converting
 * API project to apk-conf and generating and modifying CRDs belonging to
 * a particular API.
 */

package transformer

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"io"
	"mime/multipart"
	"net/http"

	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/constants"
	eventHub "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/managementserver"
	k8Yaml "sigs.k8s.io/yaml"

	"gopkg.in/yaml.v2"
)

// GenerateAPKConf will Generate the mapped .apk-conf file for a given API Project zip
func GenerateAPKConf(APIJson string, certArtifact CertificateArtifact, endpoints string, organizationID string) (string, string, uint32, map[string]eventHub.RateLimitPolicy, []EndpointSecurityConfig, *API, *AIRatelimit, *AIRatelimit, error) {

	apk := &API{}

	var apiYaml APIYaml
	var endpointsYaml EndpointsYaml

	var configuredRateLimitPoliciesMap = make(map[string]eventHub.RateLimitPolicy)

	logger.LoggerTransformer.Infof("APIJson: %v", APIJson)
	logger.LoggerTransformer.Infof("Endpoints: %v", endpoints)

	apiYamlError := json.Unmarshal([]byte(APIJson), &apiYaml)
	if apiYamlError != nil {
		apiYamlError = yaml.Unmarshal([]byte(APIJson), &apiYaml)
	}

	endpointsYamlError := json.Unmarshal([]byte(endpoints), &endpointsYaml)
	if endpointsYamlError != nil {
		endpointsYamlError = yaml.Unmarshal([]byte(endpoints), &endpointsYaml)
	}

	if apiYamlError != nil {
		logger.LoggerTransformer.Error("Error while unmarshalling api.json/api.yaml content", apiYamlError)
		return "", "null", 0, nil, []EndpointSecurityConfig{}, nil, nil, nil, apiYamlError
	}

	if endpointsYamlError != nil {
		logger.LoggerTransformer.Error("Error while unmarshalling endpoints.json/endpoints.yaml content", endpointsYamlError)
		return "", "null", 0, nil, []EndpointSecurityConfig{}, nil, nil, nil, endpointsYamlError
	}

	apiYamlData := apiYaml.Data
	logger.LoggerTransformer.Debugf("apiYamlData: %v", apiYamlData)

	apk.Name = apiYamlData.Name
	apk.Context = apiYamlData.Context
	apk.Version = apiYamlData.Version
	apk.Type = getAPIType(apiYamlData.Type)
	apk.DefaultVersion = apiYamlData.DefaultVersion
	apk.DefinitionPath = "/definition"
	apk.SubscriptionValidation = true

	if apiYamlData.SubtypeConfiguration.Subtype == "AIAPI" && apiYamlData.SubtypeConfiguration.Configuration != "" {
		// Unmarshal the _configuration field into the Configuration struct
		var config Configuration
		err := json.Unmarshal([]byte(apiYamlData.SubtypeConfiguration.Configuration), &config)
		if err != nil {
			fmt.Println("Error unmarshalling _configuration:", err)
			return "", "null", 0, nil, []EndpointSecurityConfig{}, nil, nil, nil, err
		}
		sha1ValueforCRName := config.LLMProviderID
		apk.AIProvider = &AIProvider{
			Name:       sha1ValueforCRName,
			APIVersion: "1",
		}
	}

	if apiYamlData.APIThrottlingPolicy != "" {
		rateLimitPolicy := managementserver.GetRateLimitPolicy(apiYamlData.APIThrottlingPolicy, organizationID)
		logger.LoggerTransformer.Debugf("Rate Limit Policy: %v", rateLimitPolicy)
		if rateLimitPolicy.Name != "" && rateLimitPolicy.Name != "Unlimited" {
			var rateLimitPolicyConfigured = RateLimit{
				RequestsPerUnit: rateLimitPolicy.DefaultLimit.RequestCount.RequestCount,
				Unit:            rateLimitPolicy.DefaultLimit.RequestCount.TimeUnit,
			}
			apk.RateLimit = &rateLimitPolicyConfigured
			configuredRateLimitPoliciesMap["API"] = rateLimitPolicy
		}
	}
	apkOperations := make([]Operation, len(apiYamlData.Operations))

	for i, operation := range apiYamlData.Operations {

		reqPolicyCount := len(operation.OperationPolicies.Request)
		resPolicyCount := len(operation.OperationPolicies.Response)
		reqInterceptor, resInterceptor := getReqAndResInterceptors(reqPolicyCount, resPolicyCount,
			operation.OperationPolicies.Request, operation.OperationPolicies.Response)

		var opRateLimit *RateLimit
		if apiYamlData.APIThrottlingPolicy == "" && operation.ThrottlingPolicy != "" {
			rateLimitPolicy := managementserver.GetRateLimitPolicy(operation.ThrottlingPolicy, organizationID)
			logger.LoggerTransformer.Debugf("Op Rate Limit Policy Name: %v", rateLimitPolicy.Name)
			if rateLimitPolicy.Name != "" && rateLimitPolicy.Name != "Unlimited" {
				var rateLimitPolicyConfigured = RateLimit{
					RequestsPerUnit: rateLimitPolicy.DefaultLimit.RequestCount.RequestCount,
					Unit:            rateLimitPolicy.DefaultLimit.RequestCount.TimeUnit,
				}
				opRateLimit = &rateLimitPolicyConfigured
				configuredRateLimitPoliciesMap["Resource"] = rateLimitPolicy
			}
		}
		logger.LoggerTransformer.Debugf("Operation Auth Type: %v", operation.AuthType)
		AuthSecured := true
		if operation.AuthType == "None" {
			logger.LoggerTransformer.Debugf("Setting AuthSecured to false")
			AuthSecured = false
		}
		op := &Operation{
			Target:  operation.Target,
			Verb:    operation.Verb,
			Scopes:  operation.Scopes,
			Secured: AuthSecured,
			OperationPolicies: &OperationPolicies{
				Request:  *reqInterceptor,
				Response: *resInterceptor,
			},
			RateLimit: opRateLimit,
		}
		apkOperations[i] = *op
	}

	apk.Operations = &apkOperations

	//Adding API Level Operation Policies to the conf
	reqPolicyCount := len(apiYaml.Data.APIPolicies.Request)
	resPolicyCount := len(apiYaml.Data.APIPolicies.Response)
	reqInterceptor, resInterceptor := getReqAndResInterceptors(reqPolicyCount, resPolicyCount,
		apiYaml.Data.APIPolicies.Request, apiYaml.Data.APIPolicies.Response)

	apk.APIPolicies = &OperationPolicies{
		Request:  *reqInterceptor,
		Response: *resInterceptor,
	}

	//Adding Endpoint-certificate configurations to the conf
	var endpointCertList EndpointCertDescriptor
	endCertAvailable := false

	if certArtifact.EndpointCerts != "" {
		certErr := json.Unmarshal([]byte(certArtifact.EndpointCerts), &endpointCertList)
		if certErr != nil {
			logger.LoggerTransformer.Errorf("Error while unmarshalling endpoint_cert.json content: %v", apiYamlError)
			return "", "null", 0, nil, []EndpointSecurityConfig{}, nil, nil, nil, certErr
		}
		endCertAvailable = true
	}

	sandboxURL := apiYamlData.EndpointConfig.SandboxEndpoints.URL
	prodURL := apiYamlData.EndpointConfig.ProductionEndpoints.URL
	primaryProdEndpointID := apiYamlData.PrimaryProductionEndpointID
	primarySandboxEndpointID := apiYamlData.PrimarySandboxEndpointID
	endpointList := endpointsYaml.Data
	endpointSecurityDataList := []EndpointSecurityConfig{}

	apiUniqueID := GetUniqueIDForAPI(apiYamlData.Name, apiYamlData.Version, apiYamlData.OrganizationID)
	logger.LoggerTransformer.Infof("Maxtps: %+v", apiYamlData)
	prodAIRatelimit, sandAIRatelimit := prepareAIRatelimit(apiYamlData.MaxTps)
	endpointRes := EndpointConfigurations{}
	if len(endpointList) == 0 {
		endpointSecurityData := apiYamlData.EndpointConfig.EndpointSecurity
		endpointRes, endpointSecurityData = getEndpointConfigs(sandboxURL, prodURL, endCertAvailable, endpointCertList, endpointSecurityData, apiUniqueID, prodAIRatelimit, sandAIRatelimit)
		apk.EndpointConfigurations = &endpointRes
		endpointSecurityDataList = append(endpointSecurityDataList, endpointSecurityData)
	} else {
		endpointRes, endpointSecurityDataList = getMultiEndpointConfigs(endpointList, primaryProdEndpointID, primarySandboxEndpointID, endCertAvailable, endpointCertList, apiUniqueID, prodAIRatelimit, sandAIRatelimit)
		apk.EndpointConfigurations = &endpointRes
	}

	//Adding client-certificate configurations to the conf
	var certList CertDescriptor
	certAvailable := false

	if certArtifact.ClientCerts != "" {
		certErr := json.Unmarshal([]byte(certArtifact.ClientCerts), &certList)
		if certErr != nil {
			logger.LoggerTransformer.Errorf("Error while unmarshalling client_cert.json content: %v", apiYamlError)
			return "", "null", 0, nil, []EndpointSecurityConfig{}, nil, nil, nil, certErr
		}
		certAvailable = true
	}

	authConfigList := mapAuthConfigs(apiYamlData.ID, apiYamlData.AuthorizationHeader, apiYamlData.APIKeyHeader,
		apiYamlData.SecuritySchemes, certAvailable, certList, apiUniqueID)

	apk.Authentication = &authConfigList

	corsEnabled := apiYamlData.CORSConfiguration.CORSConfigurationEnabled

	if corsEnabled {
		apk.CorsConfig = &apiYamlData.CORSConfiguration
	}

	aditionalProperties := make([]AdditionalProperty, len(apiYamlData.AdditionalProperties))

	for i, property := range apiYamlData.AdditionalProperties {
		prop := &AdditionalProperty{
			Name:  property.Name,
			Value: property.Value,
		}
		aditionalProperties[i] = *prop
	}

	apk.AdditionalProperties = &aditionalProperties

	c, marshalError := yaml.Marshal(apk)

	if marshalError != nil {
		logger.LoggerTransformer.Error("Error while marshalling apk yaml", marshalError)
		return "", "null", 0, nil, []EndpointSecurityConfig{}, nil, prodAIRatelimit, sandAIRatelimit, marshalError
	}
	return string(c), apiYamlData.RevisionedAPIID, apiYamlData.RevisionID, configuredRateLimitPoliciesMap, endpointSecurityDataList, apk, prodAIRatelimit, sandAIRatelimit, nil
}

// prepareAIRatelimit Function that accepts apiYamlData and returns AIRatelimit
func prepareAIRatelimit(maxTps *MaxTps) (*AIRatelimit, *AIRatelimit) {
	if maxTps == nil {
		return nil, nil
	}
	prodAIRL := &AIRatelimit{}
	if maxTps.TokenBasedThrottlingConfiguration == nil ||
		maxTps.TokenBasedThrottlingConfiguration.IsTokenBasedThrottlingEnabled == nil ||
		maxTps.TokenBasedThrottlingConfiguration.ProductionMaxPromptTokenCount == nil ||
		maxTps.TokenBasedThrottlingConfiguration.ProductionMaxCompletionTokenCount == nil ||
		maxTps.TokenBasedThrottlingConfiguration.ProductionMaxTotalTokenCount == nil ||
		maxTps.ProductionTimeUnit == nil {
		prodAIRL = nil
	} else {
		prodAIRL = &AIRatelimit{
			Enabled: *maxTps.TokenBasedThrottlingConfiguration.IsTokenBasedThrottlingEnabled,
			Token: TokenAIRL{
				PromptLimit:     *maxTps.TokenBasedThrottlingConfiguration.ProductionMaxPromptTokenCount,
				CompletionLimit: *maxTps.TokenBasedThrottlingConfiguration.ProductionMaxCompletionTokenCount,
				TotalLimit:      *maxTps.TokenBasedThrottlingConfiguration.ProductionMaxTotalTokenCount,
				Unit:            CapitalizeFirstLetter(*maxTps.ProductionTimeUnit),
			},
			Request: RequestAIRL{
				RequestLimit: *maxTps.Production,
				Unit:         CapitalizeFirstLetter(*maxTps.ProductionTimeUnit),
			},
		}
	}
	sandAIRL := &AIRatelimit{}
	if maxTps.TokenBasedThrottlingConfiguration == nil ||
		maxTps.TokenBasedThrottlingConfiguration.IsTokenBasedThrottlingEnabled == nil ||
		maxTps.TokenBasedThrottlingConfiguration.SandboxMaxPromptTokenCount == nil ||
		maxTps.TokenBasedThrottlingConfiguration.SandboxMaxCompletionTokenCount == nil ||
		maxTps.TokenBasedThrottlingConfiguration.SandboxMaxTotalTokenCount == nil ||
		maxTps.SandboxTimeUnit == nil {
		sandAIRL = nil
	} else {
		sandAIRL = &AIRatelimit{
			Enabled: *maxTps.TokenBasedThrottlingConfiguration.IsTokenBasedThrottlingEnabled,
			Token: TokenAIRL{
				PromptLimit:     *maxTps.TokenBasedThrottlingConfiguration.SandboxMaxPromptTokenCount,
				CompletionLimit: *maxTps.TokenBasedThrottlingConfiguration.SandboxMaxCompletionTokenCount,
				TotalLimit:      *maxTps.TokenBasedThrottlingConfiguration.SandboxMaxTotalTokenCount,
				Unit:            CapitalizeFirstLetter(*maxTps.SandboxTimeUnit),
			},
			Request: RequestAIRL{
				RequestLimit: *maxTps.Sandbox,
				Unit:         CapitalizeFirstLetter(*maxTps.SandboxTimeUnit),
			},
		}
	}

	return prodAIRL, sandAIRL
}

// CapitalizeFirstLetter takes a string and returns it with the first letter capitalized and the rest in lowercase.
func CapitalizeFirstLetter(input string) string {
	if len(input) == 0 {
		return input // Return empty string if input is empty
	}
	// Capitalize the first letter and convert the rest to lowercase
	return strings.ToUpper(string(input[0])) + strings.ToLower(input[1:])
}

// getAPIType will be selecting the appropriate API type need to be added in the apk-conf
// based on the type mentioned in the api.json
func getAPIType(protocolType string) string {
	if protocolType == "" {
		logger.LoggerTransformer.Error("Protocol type found empty. Unable to map the API Type.")
	}
	var apiType string
	switch protocolType {
	case "HTTP", "HTTPS":
		apiType = "REST"
	case "GRAPHQL":
		apiType = "GRAPHQL"
	}
	return apiType
}

// Generate the interceptor policy if request or response policy exists
func getReqAndResInterceptors(reqPolicyCount, resPolicyCount int, reqPolicies []APIMOperationPolicy, resPolicies []APIMOperationPolicy) (*[]OperationPolicy, *[]OperationPolicy) {
	var requestPolicyList, responsePolicyList []OperationPolicy
	var interceptorParams *InterceptorService
	var requestInterceptorPolicy, responseInterceptorPolicy, requestBackendJWTPolicy OperationPolicy
	var mirrorRequestPolicy OperationPolicy
	var mirrorUrls []string

	if reqPolicyCount > 0 {
		for _, reqPolicy := range reqPolicies {
			logger.LoggerTransformer.Debugf("Request Policy: %v", reqPolicy)
			if strings.HasSuffix(reqPolicy.PolicyName, constants.InterceptorService) {
				logger.LoggerTransformer.Debugf("Interceptor Type Request Policy: %v", reqPolicy)
				logger.LoggerTransformer.Debugf("Interceptor Service URL: %v", reqPolicy.Parameters[interceptorServiceURL])
				logger.LoggerTransformer.Debugf("Interceptor Includes: %v", reqPolicy.Parameters[includes])
				interceptorServiceURL := reqPolicy.Parameters[interceptorServiceURL].(string)
				includes := reqPolicy.Parameters[includes].(string)
				substrings := strings.Split(includes, ",")
				bodyEnabled := false
				headerEnabled := false
				trailersEnabled := false
				contextEnabled := false
				sslEnabled := false
				tlsSecretName := ""
				tlsSecretKey := ""
				for _, substring := range substrings {
					if strings.Contains(substring, requestHeader) {
						headerEnabled = true
					} else if strings.Contains(substring, requestBody) {
						bodyEnabled = true
					} else if strings.Contains(substring, requestTrailers) {
						trailersEnabled = true
					} else if strings.Contains(substring, requestContext) {
						contextEnabled = true
					}
				}

				if strings.Contains(interceptorServiceURL, https) {
					sslEnabled = true
				}

				if sslEnabled {
					tlsSecretName = reqPolicy.PolicyID + requestInterceptorSecretName
					tlsSecretKey = tlsKey
				}

				interceptorParams = &InterceptorService{
					BackendURL:      interceptorServiceURL,
					HeadersEnabled:  headerEnabled,
					BodyEnabled:     bodyEnabled,
					TrailersEnabled: trailersEnabled,
					ContextEnabled:  contextEnabled,
					TLSSecretName:   tlsSecretName,
					TLSSecretKey:    tlsSecretKey,
				}

				// Create an instance of OperationPolicy
				requestInterceptorPolicy = OperationPolicy{
					PolicyName:    interceptorPolicy,
					PolicyVersion: v1,
					Parameters:    interceptorParams,
				}
			} else if reqPolicy.PolicyName == constants.BackendJWT {
				encoding := reqPolicy.Parameters[encoding].(string)
				header := reqPolicy.Parameters[header].(string)
				signingAlgorithm := reqPolicy.Parameters[signingAlgorithm].(string)
				tokenTTL := reqPolicy.Parameters[tokenTTL].(string)
				tokenTTLConverted, err := strconv.Atoi(tokenTTL)
				if err != nil {
					logger.LoggerTransformer.Errorf("Error while converting tokenTTL to integer: %v", err)
				}

				if encoding == base64Url {
					encoding = base64url
				}

				backendJWTParams := &BackendJWT{
					Encoding:         encoding,
					Header:           header,
					SigningAlgorithm: signingAlgorithm,
					TokenTTL:         tokenTTLConverted,
				}

				// Create an instance of OperationPolicy
				requestBackendJWTPolicy = OperationPolicy{
					PolicyName:    backendJWTPolicy,
					PolicyVersion: v1,
					Parameters:    backendJWTParams,
				}
			} else if reqPolicy.PolicyName == constants.AddHeader {
				logger.LoggerTransformer.Debugf("AddHeader Type Request Policy: %v", reqPolicy)
				requestAddHeader := OperationPolicy{
					PolicyName:    addHeaderPolicy,
					PolicyVersion: v1,
					Parameters: Header{
						HeaderName:  reqPolicy.Parameters[headerName].(string),
						HeaderValue: reqPolicy.Parameters[headerValue].(string),
					},
				}
				requestPolicyList = append(requestPolicyList, requestAddHeader)
			} else if reqPolicy.PolicyName == constants.RemoveHeader {
				logger.LoggerTransformer.Debugf("RemoveHeader Type Request Policy: %v", reqPolicy)
				requestRemoveHeader := OperationPolicy{
					PolicyName:    removeHeaderPolicy,
					PolicyVersion: v1,
					Parameters: Header{
						HeaderName: reqPolicy.Parameters[headerName].(string),
					},
				}
				requestPolicyList = append(requestPolicyList, requestRemoveHeader)
			} else if reqPolicy.PolicyName == constants.RedirectRequest {
				logger.LoggerTransformer.Debugf("RedirectRequest Type Request Policy: %v", reqPolicy)
				redirectRequestPolicy := OperationPolicy{
					PolicyName:    requestRedirectPolicy,
					PolicyVersion: v1,
				}
				parameters := RedirectPolicy{
					URL: reqPolicy.Parameters[url].(string),
				}
				switch v := reqPolicy.Parameters[statusCode].(type) {
				case int:
					parameters.StatusCode = v
				case string:
					if intValue, err := strconv.Atoi(v); err == nil {
						parameters.StatusCode = intValue
					} else {
						logger.LoggerTransformer.Error("Invalid status code provided.")
					}
				default:
					parameters.StatusCode = 302
				}
				redirectRequestPolicy.Parameters = parameters
				requestPolicyList = append(requestPolicyList, redirectRequestPolicy)
			} else if reqPolicy.PolicyName == constants.MirrorRequest {
				logger.LoggerTransformer.Debugf("MirrorRequest Type Request Policy: %v", reqPolicy)
				if mirrorRequestPolicy.PolicyName == "" {
					mirrorRequestPolicy = OperationPolicy{
						PolicyName:    requestMirrorPolicy,
						PolicyVersion: v1,
					}
					mirrorUrls = []string{}
				}
				if reqPolicyParameters, ok := reqPolicy.Parameters[url]; ok {
					url := reqPolicyParameters.(string)
					mirrorUrls = append(mirrorUrls, url)
				}
			}
		}
	}

	if resPolicyCount > 0 {
		for _, resPolicy := range resPolicies {
			if resPolicy.PolicyName == constants.InterceptorService {
				interceptorServiceURL := resPolicy.Parameters[interceptorServiceURL].(string)
				includes := resPolicy.Parameters[includes].(string)
				substrings := strings.Split(includes, ",")
				bodyEnabled := false
				headerEnabled := false
				trailersEnabled := false
				contextEnabled := false
				sslEnabled := false
				tlsSecretName := ""
				tlsSecretKey := ""
				for _, substring := range substrings {
					if strings.Contains(substring, requestHeader) {
						headerEnabled = true
					} else if strings.Contains(substring, requestBody) {
						bodyEnabled = true
					} else if strings.Contains(substring, requestTrailers) {
						trailersEnabled = true
					} else if strings.Contains(substring, requestContext) {
						contextEnabled = true
					}
				}

				if strings.Contains(interceptorServiceURL, https) {
					sslEnabled = true
				}

				if sslEnabled {
					tlsSecretName = resPolicies[0].PolicyID + responseInterceptorSecretName
					tlsSecretKey = tlsKey
				}

				interceptorParams = &InterceptorService{
					BackendURL:      interceptorServiceURL,
					HeadersEnabled:  headerEnabled,
					BodyEnabled:     bodyEnabled,
					TrailersEnabled: trailersEnabled,
					ContextEnabled:  contextEnabled,
					TLSSecretName:   tlsSecretName,
					TLSSecretKey:    tlsSecretKey,
				}

				// Create an instance of OperationPolicy
				responseInterceptorPolicy = OperationPolicy{
					PolicyName:    interceptorPolicy,
					PolicyVersion: v1,
					Parameters:    interceptorParams,
				}
			} else if resPolicy.PolicyName == constants.AddHeader {
				logger.LoggerTransformer.Debugf("AddHeader Type Response Policy: %v", resPolicy)

				responseAddHeader := OperationPolicy{
					PolicyName:    addHeaderPolicy,
					PolicyVersion: v2,
					Parameters: Header{
						HeaderName:  resPolicy.Parameters[headerName].(string),
						HeaderValue: resPolicy.Parameters[headerValue].(string),
					},
				}
				responsePolicyList = append(responsePolicyList, responseAddHeader)
			} else if resPolicy.PolicyName == constants.RemoveHeader {
				logger.LoggerTransformer.Debugf("RemoveHeader Type Response Policy: %v", resPolicy)
				responseRemoveHeader := OperationPolicy{
					PolicyName:    removeHeaderPolicy,
					PolicyVersion: v1,
					Parameters: Header{
						HeaderName: resPolicy.Parameters[headerName].(string),
					},
				}
				responsePolicyList = append(responsePolicyList, responseRemoveHeader)
			}
		}
	}

	if reqPolicyCount > 0 {
		if requestInterceptorPolicy.PolicyName != "" {
			requestPolicyList = append(requestPolicyList, requestInterceptorPolicy)
		}
		if requestBackendJWTPolicy.PolicyName != "" {
			requestPolicyList = append(requestPolicyList, requestBackendJWTPolicy)
		}
		if mirrorRequestPolicy.PolicyName != "" {
			mirrorRequestPolicy.Parameters = URLList{
				URLs: mirrorUrls,
			}
			requestPolicyList = append(requestPolicyList, mirrorRequestPolicy)
		}
	}

	if resPolicyCount > 0 {
		if responseInterceptorPolicy.PolicyName != "" {
			responsePolicyList = append(responsePolicyList, responseInterceptorPolicy)
		}
	}
	return &requestPolicyList, &responsePolicyList
}

// mapAuthConfigs will take the security schemes as the parameter and will return the mapped auth configs to be
// added into the apk-conf
func mapAuthConfigs(apiUUID string, authHeader string, configuredAPIKeyHeader string, securitySchemes []string, certAvailable bool, certList CertDescriptor, apiUniqueID string) []AuthConfiguration {
	var authConfigs []AuthConfiguration
	if StringExists(oAuth2SecScheme, securitySchemes) {
		var oauth2Config AuthConfiguration
		oauth2Config.AuthType = oAuth2
		oauth2Config.Enabled = true
		oauth2Config.HeaderName = authHeader
		if StringExists(applicationSecurityMandatory, securitySchemes) {
			oauth2Config.Required = mandatory
		} else {
			oauth2Config.Required = optional
		}

		authConfigs = append(authConfigs, oauth2Config)
	}
	if !StringExists("oauth2", securitySchemes) {
		oAuth2DisabledConfig := AuthConfiguration{
			AuthType: oAuth2,
			Enabled:  false,
		}
		authConfigs = append(authConfigs, oAuth2DisabledConfig)
	}
	if StringExists(mutualSSL, securitySchemes) && certAvailable {
		var mtlsConfig AuthConfiguration
		mtlsConfig.AuthType = mTLS
		mtlsConfig.Enabled = true
		if StringExists(mutualSSLMandatory, securitySchemes) {
			mtlsConfig.Required = mandatory
		} else {
			mtlsConfig.Required = optional
		}

		clientCerts := make([]Certificate, len(certList.CertData))

		for i, cert := range certList.CertData {
			prop := &Certificate{
				Name: apiUniqueID + "-" + cert.Alias,
				Key:  cert.Certificate,
			}
			clientCerts[i] = *prop
		}
		mtlsConfig.Certificates = clientCerts
		authConfigs = append(authConfigs, mtlsConfig)
	}

	internalKeyAuthConfig := AuthConfiguration{
		AuthType:   jwt,
		Enabled:    true,
		Audience:   []string{apiUUID},
		HeaderName: internalKeyHeader,
	}
	authConfigs = append(authConfigs, internalKeyAuthConfig)

	if StringExists(apiKeySecScheme, securitySchemes) {
		apiKeyAuthConfig := AuthConfiguration{
			AuthType:       apiKey,
			Enabled:        true,
			HeaderName:     configuredAPIKeyHeader,
			HeaderEnabled:  true,
			QueryParamName: apiKeyHeader,
		}
		if StringExists(applicationSecurityMandatory, securitySchemes) {
			apiKeyAuthConfig.Required = mandatory
		} else if StringExists(applicationSecurityOptional, securitySchemes) {
			apiKeyAuthConfig.Required = optional
		}
		authConfigs = append(authConfigs, apiKeyAuthConfig)
	}
	return authConfigs
}

// getEndpointConfigs will map the endpoints and there security configurations and returns them
// TODO: Currently the APK-Conf does not support giving multiple certs for a particular endpoint.
// After fixing this, the following logic should be changed to map multiple cert configs
func getEndpointConfigs(sandboxURL string, prodURL string, endCertAvailable bool, endpointCertList EndpointCertDescriptor, endpointSecurityData EndpointSecurityConfig, apiUniqueID string, prodAIRatelimit *AIRatelimit, sandAIRatelimit *AIRatelimit) (EndpointConfigurations, EndpointSecurityConfig) {
	var sandboxEndpointConf, prodEndpointConf EndpointConfiguration
	var sandBoxEndpointEnabled = false
	var prodEndpointEnabled = false
	if sandboxURL != "" {
		sandBoxEndpointEnabled = true
	}
	if prodURL != "" {
		prodEndpointEnabled = true
	}
	if prodAIRatelimit != nil {
		prodEndpointConf.AIRatelimit = *prodAIRatelimit
	}
	if sandAIRatelimit != nil {
		sandboxEndpointConf.AIRatelimit = *sandAIRatelimit
	}
	sandboxEndpointConf.Endpoint = sandboxURL
	prodEndpointConf.Endpoint = prodURL
	if endCertAvailable {
		for _, endCert := range endpointCertList.EndpointCertData {
			if endCert.Endpoint == sandboxURL {
				sandboxEndpointConf.EndCertificate = EndpointCertificate{
					Name: endCert.Alias,
					Key:  endCert.Certificate,
				}
			}
			if endCert.Endpoint == prodURL {
				prodEndpointConf.EndCertificate = EndpointCertificate{
					Name: endCert.Alias,
					Key:  endCert.Certificate,
				}
			}
		}
	}

	if endpointSecurityData.Sandbox.Enabled {
		endpointSecurityData.Sandbox.EndpointUUID = "primary"
		sandboxEndpointConf.EndSecurity.Enabled = true
		if endpointSecurityData.Sandbox.Type == "apikey" {
			sandboxEndpointConf.EndSecurity.SecurityType = SecretInfo{
				SecretName:     strings.Join([]string{apiUniqueID, generateSHA1Hash(endpointSecurityData.Sandbox.EndpointUUID), "sandbox", "secret"}, "-"),
				In:             "Header",
				APIKeyNameKey:  endpointSecurityData.Sandbox.APIKeyIdentifier,
				APIKeyValueKey: "apiKey",
			}
		} else {
			sandboxEndpointConf.EndSecurity.SecurityType = SecretInfo{
				SecretName:  strings.Join([]string{apiUniqueID, generateSHA1Hash(endpointSecurityData.Sandbox.EndpointUUID), "sandbox", "secret"}, "-"),
				UsernameKey: "username",
				PasswordKey: "password",
			}
		}
	}

	if endpointSecurityData.Production.Enabled {
		endpointSecurityData.Production.EndpointUUID = "primary"
		prodEndpointConf.EndSecurity.Enabled = true
		if endpointSecurityData.Production.Type == "apikey" {
			prodEndpointConf.EndSecurity.SecurityType = SecretInfo{
				SecretName:     strings.Join([]string{apiUniqueID, generateSHA1Hash(endpointSecurityData.Production.EndpointUUID), "production", "secret"}, "-"),
				In:             "Header",
				APIKeyNameKey:  endpointSecurityData.Production.APIKeyIdentifier,
				APIKeyValueKey: "apiKey",
			}
		} else {
			prodEndpointConf.EndSecurity.SecurityType = SecretInfo{
				SecretName:  strings.Join([]string{apiUniqueID, generateSHA1Hash(endpointSecurityData.Production.EndpointUUID), "production", "secret"}, "-"),
				UsernameKey: "username",
				PasswordKey: "password",
			}
		}
	}

	epconfigs := EndpointConfigurations{}
	sandboxEndpoints := []EndpointConfiguration{}
	productionEndpoints := []EndpointConfiguration{}
	if sandBoxEndpointEnabled && prodEndpointEnabled {
		sandboxEndpoints = append(sandboxEndpoints, sandboxEndpointConf)
		productionEndpoints = append(productionEndpoints, prodEndpointConf)
		epconfigs = EndpointConfigurations{
			Sandbox:    &sandboxEndpoints,
			Production: &productionEndpoints,
		}
	} else if sandBoxEndpointEnabled {
		sandboxEndpoints = append(sandboxEndpoints, sandboxEndpointConf)
		epconfigs = EndpointConfigurations{
			Sandbox: &sandboxEndpoints,
		}
	} else if prodEndpointEnabled {
		productionEndpoints = append(productionEndpoints, prodEndpointConf)
		epconfigs = EndpointConfigurations{
			Production: &productionEndpoints,
		}
	}
	return epconfigs, endpointSecurityData
}

// getMultiEndpointConfigs will map the endpoints and there security configurations and returns them
func getMultiEndpointConfigs(endpointList []Endpoint, primaryProdEndpointID string, primarySandboxEndpointID string, endCertAvailable bool, endpointCertList EndpointCertDescriptor, apiUniqueID string, prodAIRatelimit *AIRatelimit, sandAIRatelimit *AIRatelimit) (EndpointConfigurations, []EndpointSecurityConfig) {
	sandboxEndpoints := []EndpointConfiguration{}
	productionEndpoints := []EndpointConfiguration{}
	endpointSecurityConfigs := []EndpointSecurityConfig{}
	if primaryProdEndpointID == "" && primarySandboxEndpointID == "" {
		logger.LoggerTransformer.Error("Primary Production and Sandbox Endpoint ID's are empty. Unable to map the endpoints.")
	}
	for _, endpoint := range endpointList {
		prodURL := endpoint.EndpointConfig.ProductionEndpoints.URL
		sandboxURL := endpoint.EndpointConfig.SandboxEndpoints.URL
		endpointSecurityData := endpoint.EndpointConfig.EndpointSecurity
		endpointSecurityData.Production.EndpointUUID = endpoint.EndpointUUID
		endpointSecurityData.Sandbox.EndpointUUID = endpoint.EndpointUUID
		endpointSecurityConfigs = append(endpointSecurityConfigs, endpointSecurityData)
		var sandboxEndpointConf, prodEndpointConf EndpointConfiguration
		var sandBoxEndpointEnabled = false
		var prodEndpointEnabled = false
		if sandboxURL != "" {
			sandBoxEndpointEnabled = true
		}
		if prodURL != "" {
			prodEndpointEnabled = true
		}
		if prodAIRatelimit != nil {
			prodEndpointConf.AIRatelimit = *prodAIRatelimit
		}
		if sandAIRatelimit != nil {
			sandboxEndpointConf.AIRatelimit = *sandAIRatelimit
		}
		sandboxEndpointConf.Endpoint = sandboxURL
		prodEndpointConf.Endpoint = prodURL
		if endCertAvailable {
			for _, endCert := range endpointCertList.EndpointCertData {
				if endCert.Endpoint == sandboxURL {
					sandboxEndpointConf.EndCertificate = EndpointCertificate{
						Name: endCert.Alias,
						Key:  endCert.Certificate,
					}
				}
				if endCert.Endpoint == prodURL {
					prodEndpointConf.EndCertificate = EndpointCertificate{
						Name: endCert.Alias,
						Key:  endCert.Certificate,
					}
				}
			}
		}

		if endpointSecurityData.Sandbox.Enabled {
			sandboxEndpointConf.EndSecurity.Enabled = true

			if endpointSecurityData.Sandbox.Type == "apikey" {
				sandboxEndpointConf.EndSecurity.SecurityType = SecretInfo{
					SecretName:     strings.Join([]string{apiUniqueID, generateSHA1Hash(endpoint.EndpointUUID), "sandbox", "secret"}, "-"),
					In:             "Header",
					APIKeyNameKey:  endpointSecurityData.Sandbox.APIKeyIdentifier,
					APIKeyValueKey: "apiKey",
				}
			} else {
				sandboxEndpointConf.EndSecurity.SecurityType = SecretInfo{
					SecretName:  strings.Join([]string{apiUniqueID, generateSHA1Hash(endpoint.EndpointUUID), "sandbox", "secret"}, "-"),
					UsernameKey: "username",
					PasswordKey: "password",
				}
			}
		}

		if endpointSecurityData.Production.Enabled {
			prodEndpointConf.EndSecurity.Enabled = true
			if endpointSecurityData.Production.Type == "apikey" {
				prodEndpointConf.EndSecurity.SecurityType = SecretInfo{
					SecretName:     strings.Join([]string{apiUniqueID, generateSHA1Hash(endpoint.EndpointUUID), "production", "secret"}, "-"),
					In:             "Header",
					APIKeyNameKey:  endpointSecurityData.Production.APIKeyIdentifier,
					APIKeyValueKey: "apiKey",
				}
			} else {
				prodEndpointConf.EndSecurity.SecurityType = SecretInfo{
					SecretName:  strings.Join([]string{apiUniqueID, generateSHA1Hash(endpoint.EndpointUUID), "production", "secret"}, "-"),
					UsernameKey: "username",
					PasswordKey: "password",
				}
			}
		}

		if sandBoxEndpointEnabled && prodEndpointEnabled {
			sandboxEndpoints = append(sandboxEndpoints, sandboxEndpointConf)
			productionEndpoints = append(productionEndpoints, prodEndpointConf)
		} else if sandBoxEndpointEnabled {
			sandboxEndpoints = append(sandboxEndpoints, sandboxEndpointConf)
		} else if prodEndpointEnabled {
			productionEndpoints = append(productionEndpoints, prodEndpointConf)
		}
	}
	epconfigs := EndpointConfigurations{
		Sandbox:    &sandboxEndpoints,
		Production: &productionEndpoints,
	}
	logger.LoggerTransformer.Debugf("Endpoint Configurations: %v", epconfigs)
	return epconfigs, endpointSecurityConfigs
}

// GenerateCRs takes the .apk-conf, api definition, vHost and the organization for a particular API and then generate and returns
// the relavant CRD set as a zip
func GenerateCRs(apkConf string, apiDefinition string, certContainer CertContainer, k8ResourceGenEndpoint string, organizationID string) (*K8sArtifacts, error) {
	k8sArtifact := K8sArtifacts{HTTPRoutes: make(map[string]*gwapiv1.HTTPRoute), GQLRoutes: make(map[string]*dpv1alpha2.GQLRoute), Backends: make(map[string]*dpv1alpha2.Backend), Scopes: make(map[string]*dpv1alpha1.Scope), Authentication: make(map[string]*dpv1alpha2.Authentication), APIPolicies: make(map[string]*dpv1alpha4.APIPolicy), InterceptorServices: make(map[string]*dpv1alpha1.InterceptorService), ConfigMaps: make(map[string]*corev1.ConfigMap), Secrets: make(map[string]*corev1.Secret), RateLimitPolicies: make(map[string]*dpv1alpha1.RateLimitPolicy), AIRateLimitPolicies: make(map[string]*dpv1alpha3.AIRateLimitPolicy)}
	if apkConf == "" {
		logger.LoggerTransformer.Error("Empty apk-conf parameter provided. Unable to generate CRDs.")
		return nil, errors.New("Error: APK-Conf can't be empty")
	}

	if apiDefinition == "" {
		logger.LoggerTransformer.Error("Empty api definition provided. Unable to generate CRDs.")
		return nil, errors.New("Error: API Definition can't be empty")
	}

	// Create a buffer to store the request body
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add apkConfiguration field and store the passed APK Conf file
	if err := writer.WriteField("apkConfiguration", apkConf); err != nil {
		logger.LoggerTransformer.Error("Error writing apkConfiguration field:", err)
		return nil, err
	}

	// Add apiDefinition field and store the passed API Definition file
	if err := writer.WriteField("definitionFile", apiDefinition); err != nil {
		logger.LoggerTransformer.Error("Error writing definitionFile field:", err)
		return nil, err
	}

	// Close the multipart writer
	writer.Close()

	k8sResourceEndpointWithOrg := k8ResourceGenEndpoint + "?organization=" + organizationID

	// Create the HTTP request
	request, err := http.NewRequest(postHTTPMethod, k8sResourceEndpointWithOrg, &requestBody)
	if err != nil {
		logger.LoggerTransformer.Error("Error creating HTTP request:", err)
		return nil, err
	}

	// Set the Content-Type header
	request.Header.Set(contentTypeHeader, writer.FormDataContentType())
	// Certificate validation is turned off as linkerd would be used for mTLS between the two components
	tr := &http.Transport{
		/* #nosec */
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Make the request
	client := &http.Client{Transport: tr}

	response, err := client.Do(request)
	if err != nil {
		logger.LoggerTransformer.Error("Error making HTTP request:", err)
		return nil, err
	}

	defer response.Body.Close()

	// Check the HTTP status code
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		logger.LoggerTransformer.Errorf("HTTP request failed with status code: %d", response.StatusCode)
		return nil, fmt.Errorf("HTTP request failed with status code: %v", response.Body)
	}

	//Extracting response body to get the CRD zipfile
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.LoggerTransformer.Error("Error reading response body:", err)
		panic(err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		logger.LoggerTransformer.Error("Unable to transform the initial CRDs:", err)
		return nil, err
	}
	for _, zipFile := range zipReader.File {
		fileReader, err := zipFile.Open()
		if err != nil {
			logger.LoggerTransformer.Errorf("Failed to open YAML file inside zip: %v", err)
			return nil, err
		}
		defer fileReader.Close()

		yamlData, err := io.ReadAll(fileReader)
		if err != nil {
			logger.LoggerTransformer.Errorf("Failed to read YAML file inside zip: %v", err)
			return nil, err
		}

		var crdData map[string]interface{}
		if err := yaml.Unmarshal(yamlData, &crdData); err != nil {
			logger.LoggerTransformer.Errorf("Failed to unmarshal YAML data to parse the Kind: %v", err)
			return nil, err
		}

		kind, ok := crdData["kind"].(string)
		if !ok {
			logger.LoggerTransformer.Errorf("Kind attribute not found in the given yaml file.")
			return nil, err
		}

		switch kind {
		case "APIPolicy":
			var apiPolicy dpv1alpha4.APIPolicy
			err = k8Yaml.Unmarshal(yamlData, &apiPolicy)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling APIPolicy YAML: %v", err)
				continue
			}
			k8sArtifact.APIPolicies[apiPolicy.ObjectMeta.Name] = &apiPolicy
		case "HTTPRoute":
			var httpRoute gwapiv1.HTTPRoute
			err = k8Yaml.Unmarshal(yamlData, &httpRoute)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling HTTPRoute YAML: %v", err)
				continue
			}
			k8sArtifact.HTTPRoutes[httpRoute.ObjectMeta.Name] = &httpRoute

		case "Backend":
			var backend dpv1alpha2.Backend
			err = k8Yaml.Unmarshal(yamlData, &backend)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling Backend YAML: %v", err)
				continue
			}
			k8sArtifact.Backends[backend.ObjectMeta.Name] = &backend

		case "ConfigMap":
			var configMap corev1.ConfigMap
			err = k8Yaml.Unmarshal(yamlData, &configMap)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling ConfigMap YAML: %v", err)
				continue
			}
			k8sArtifact.ConfigMaps[configMap.ObjectMeta.Name] = &configMap
		case "Authentication":
			var authPolicy dpv1alpha2.Authentication
			err = k8Yaml.Unmarshal(yamlData, &authPolicy)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling Authentication YAML: %v", err)
				continue
			}
			k8sArtifact.Authentication[authPolicy.ObjectMeta.Name] = &authPolicy

		case "API":
			var api dpv1alpha3.API
			err = k8Yaml.Unmarshal(yamlData, &api)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling API YAML: %v", err)
				continue
			}
			k8sArtifact.API = api
		case "InterceptorService":
			var interceptorService dpv1alpha1.InterceptorService
			err = k8Yaml.Unmarshal(yamlData, &interceptorService)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling InterceptorService YAML: %v", err)
				continue
			}
			k8sArtifact.InterceptorServices[interceptorService.Name] = &interceptorService
		case "BackendJWT":
			var backendJWT *dpv1alpha1.BackendJWT
			err = k8Yaml.Unmarshal(yamlData, &backendJWT)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling BackendJWT YAML: %v", err)
				continue
			}
			k8sArtifact.BackendJWT = backendJWT
		case "Scope":
			var scope dpv1alpha1.Scope
			err = k8Yaml.Unmarshal(yamlData, &scope)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling Scope YAML: %v", err)
				continue
			}
			k8sArtifact.Scopes[scope.Name] = &scope
		case "RateLimitPolicy":
			var rateLimitPolicy dpv1alpha1.RateLimitPolicy
			err = k8Yaml.Unmarshal(yamlData, &rateLimitPolicy)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling RateLimitPolicy YAML: %v", err)
				continue
			}
			k8sArtifact.RateLimitPolicies[rateLimitPolicy.Name] = &rateLimitPolicy
		case "AIRateLimitPolicy":
			var aiRateLimitPolicy dpv1alpha3.AIRateLimitPolicy
			err = k8Yaml.Unmarshal(yamlData, &aiRateLimitPolicy)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling AIRateLimitPolicy YAML: %v", err)
				continue
			}
			k8sArtifact.AIRateLimitPolicies[aiRateLimitPolicy.Name] = &aiRateLimitPolicy
		case "Secret":
			var secret corev1.Secret
			err = k8Yaml.Unmarshal(yamlData, &secret)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling Secret YAML: %v", err)
				continue
			}
			k8sArtifact.Secrets[secret.Name] = &secret
		case "GQLRoute":
			var gqlRoute dpv1alpha2.GQLRoute
			err = k8Yaml.Unmarshal(yamlData, &gqlRoute)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling GQLRoute YAML: %v", err)
				continue
			}
			k8sArtifact.GQLRoutes[gqlRoute.Name] = &gqlRoute
		default:
			logger.LoggerSync.Errorf("[!]Unknown Kind parsed from the YAML File: %v", kind)
		}
	}
	// Create ConfigMap to store the cert data if mTLS has enabled
	if certContainer.ClientCertObj.CertAvailable {
		createConfigMaps(certContainer.ClientCertObj.ClientCertFiles, &k8sArtifact)
	}

	// Create ConfigMap to store the cert data if endpoint security has enabled
	if certContainer.EndpointCertObj.CertAvailable {
		createConfigMaps(certContainer.EndpointCertObj.EndpointCertFiles, &k8sArtifact)
	}

	createEndpointSecrets(certContainer.SecretData, &k8sArtifact)

	return &k8sArtifact, nil
}

// UpdateCRS cr update
func UpdateCRS(k8sArtifact *K8sArtifacts, environments *[]Environment, organizationID string, apiUUID string, revisionID string, namespace string, configuredRateLimitPoliciesMap map[string]eventHub.RateLimitPolicy) {
	addOrganization(k8sArtifact, organizationID)
	addRevisionAndAPIUUID(k8sArtifact, apiUUID, revisionID)
	for _, environment := range *environments {
		replaceVhost(k8sArtifact, environment.Vhost, environment.Type)
	}
	addRateLimitPolicyNames(k8sArtifact, configuredRateLimitPoliciesMap)
}

// replaceVhost will take the httpRoute CR and replace the default vHost with the one passed inside
// the deploymemt descriptor
func replaceVhost(k8sArtifact *K8sArtifacts, vhost string, deploymentType string) {
	if deploymentType == "hybrid" {
		// append sandbox. part to available vhost to generate sandbox vhost
		if k8sArtifact.API.Spec.Production != nil {
			for _, routeName := range k8sArtifact.API.Spec.Production {
				for _, routes := range routeName.RouteRefs {
					httprouteRef, foundHTTPRoute := k8sArtifact.HTTPRoutes[routes]
					if foundHTTPRoute {
						httprouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname(vhost)}
					}
					gqlRouteRef, foundGQLRoute := k8sArtifact.GQLRoutes[routes]
					if foundGQLRoute {
						gqlRouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname(vhost)}
					}
				}
			}
		}
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					httprouteRef, foundHTTPRoute := k8sArtifact.HTTPRoutes[routes]
					if foundHTTPRoute {
						httprouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname("sandbox." + vhost)}
					}
					gqlRouteRef, foundGQLRoute := k8sArtifact.GQLRoutes[routes]
					if foundGQLRoute {
						gqlRouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname("sandbox." + vhost)}
					}
				}
			}
		}
	} else if deploymentType == "sandbox" {
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					httprouteRef, foundHTTPRoute := k8sArtifact.HTTPRoutes[routes]
					if foundHTTPRoute {
						httprouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname(vhost)}
					}
					gqlRouteRef, foundGQLRoute := k8sArtifact.GQLRoutes[routes]
					if foundGQLRoute {
						gqlRouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname(vhost)}
					}
				}
			}
		}
		if k8sArtifact.API.Spec.Production != nil {
			for _, routeName := range k8sArtifact.API.Spec.Production {
				for _, routes := range routeName.RouteRefs {
					delete(k8sArtifact.HTTPRoutes, routes)
					delete(k8sArtifact.GQLRoutes, routes)
				}
			}
			k8sArtifact.API.Spec.Production = []dpv1alpha3.EnvConfig{}
		}
	} else {
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					httprouteRef, foundHTTPRoute := k8sArtifact.HTTPRoutes[routes]
					if foundHTTPRoute {
						httprouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname(vhost)}
					}
					gqlRouteRef, foundGQLRoute := k8sArtifact.GQLRoutes[routes]
					if foundGQLRoute {
						gqlRouteRef.Spec.Hostnames = []gwapiv1.Hostname{gwapiv1.Hostname(vhost)}
					}
				}
			}
		}
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					delete(k8sArtifact.HTTPRoutes, routes)
					delete(k8sArtifact.GQLRoutes, routes)
				}
			}
			k8sArtifact.API.Spec.Sandbox = []dpv1alpha3.EnvConfig{}
		}
	}
}

// addOrganization will take the API CR and change the organization to the one passed inside
// the deploymemt descriptor
func addOrganization(k8sArtifact *K8sArtifacts, organization string) {
	k8sArtifact.API.Spec.Organization = organization
	organizationHash := generateSHA1Hash(organization)
	k8sArtifact.API.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	for _, apiPolicies := range k8sArtifact.APIPolicies {
		apiPolicies.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
	for _, httproutes := range k8sArtifact.HTTPRoutes {
		httproutes.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
	for _, gqlroutes := range k8sArtifact.GQLRoutes {
		gqlroutes.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
	for _, authentication := range k8sArtifact.Authentication {
		authentication.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
	for _, backend := range k8sArtifact.Backends {
		backend.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
	for _, configMap := range k8sArtifact.ConfigMaps {
		configMap.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
	for _, secret := range k8sArtifact.Secrets {
		secret.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
	for _, scope := range k8sArtifact.Scopes {
		scope.ObjectMeta.Labels[k8sOrganizationField] = organizationHash
	}
}

// addRevisionAndAPIUUID will add the API ID and the revision field attributes to the API CR
func addRevisionAndAPIUUID(k8sArtifact *K8sArtifacts, apiID string, revisionID string) {
	k8sArtifact.API.ObjectMeta.Labels[k8APIUuidField] = apiID
	k8sArtifact.API.ObjectMeta.Labels[k8RevisionField] = revisionID
}

// addRateLimitPolicyNames will add the rate limit policy names to the respective CRs
func addRateLimitPolicyNames(k8sArtifact *K8sArtifacts, configuredRateLimitPoliciesMap map[string]eventHub.RateLimitPolicy) {
	logger.LoggerTransformer.Infof("Rate Limit Policies: %v", configuredRateLimitPoliciesMap)
	for _, rateLimitPolicy := range k8sArtifact.RateLimitPolicies {
		if strings.Contains(rateLimitPolicy.Name, "api-") {
			rateLimitPolicy.ObjectMeta.Labels[k8sRateLimitPolicyNameField] = generateSHA1Hash(configuredRateLimitPoliciesMap["API"].Name)
			logger.LoggerTransformer.Infof("Rate Limit Policy Name: %v", rateLimitPolicy.ObjectMeta.Labels[k8sRateLimitPolicyNameField])
		} else if strings.Contains(rateLimitPolicy.Name, "resource-") {
			rateLimitPolicy.ObjectMeta.Labels[k8sRateLimitPolicyNameField] = generateSHA1Hash(configuredRateLimitPoliciesMap["Resource"].Name)
			logger.LoggerTransformer.Infof("Rate Limit Policy Name: %v", rateLimitPolicy.ObjectMeta.Labels[k8sRateLimitPolicyNameField])
		}
	}
}

// generateSHA1Hash returns the SHA1 hash for the given string
func generateSHA1Hash(input string) string {
	h := sha1.New() /* #nosec */
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// createConfigMaps returns a marshalled yaml of ConfigMap kind after adding the given values
func createConfigMaps(certFiles map[string]string, k8sArtifact *K8sArtifacts) {
	for confKey, confValue := range certFiles {
		pathSegments := strings.Split(confKey, ".")
		configName := pathSegments[0]

		//TODO: Have to take the version, namespace as parameters instead of hardcoding
		cm := corev1.ConfigMap{}
		cm.APIVersion = "v1"
		cm.Kind = "ConfigMap"
		cm.ObjectMeta.Name = k8sArtifact.API.Name + "-" + configName

		if cm.ObjectMeta.Labels == nil {
			cm.ObjectMeta.Labels = make(map[string]string)
		}

		if cm.Data == nil {
			cm.Data = make(map[string]string)
		}
		apimCert := confValue
		// Remove "-----BEGIN CERTIFICATE-----" and "-----END CERTIFICATE-----" strings
		pemCert := strings.ReplaceAll(apimCert, "-----BEGIN CERTIFICATE-----", "")
		pemCert = strings.ReplaceAll(pemCert, "-----END CERTIFICATE-----", "")
		pemCert = strings.TrimSpace(pemCert)
		// Decode the Base64 encoded certificate content
		decodedCert, err := base64.StdEncoding.DecodeString(pemCert)
		logger.LoggerTransformer.Debugf("Decoded Certificate: %v", decodedCert)
		if err != nil {
			logger.LoggerTransformer.Errorf("Error decoding the certificate: %v", err)
		}
		cm.Data[confKey] = string(decodedCert)
		certConfigMap := &cm

		logger.LoggerTransformer.Debugf("New ConfigMap Data: %v", *certConfigMap)
		k8sArtifact.ConfigMaps[certConfigMap.ObjectMeta.Name] = certConfigMap
	}
}

// createEndpointSecrets creates and links the secret CRs need to be created for handling the endpoint security
func createEndpointSecrets(secretDataList []EndpointSecurityConfig, k8sArtifact *K8sArtifacts) {
	createSecret := func(environment string, username, password string, apiKeyValue string, securityType string, endpointUUID string) {
		var secret corev1.Secret
		if securityType == "apikey" {
			secret = corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.Join([]string{k8sArtifact.API.Name, generateSHA1Hash(endpointUUID), environment, "secret"}, "-"),
					Namespace: k8sArtifact.API.Namespace,
					Labels:    make(map[string]string),
				},
				Data: map[string][]byte{
					"apiKey": []byte(apiKeyValue),
				},
			}
		} else {
			secret = corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.Join([]string{k8sArtifact.API.Name, generateSHA1Hash(endpointUUID), environment, "secret"}, "-"),
					Namespace: k8sArtifact.API.Namespace,
					Labels:    make(map[string]string),
				},
				Data: map[string][]byte{
					"username": []byte(username),
					"password": []byte(password),
				},
			}
		}
		logger.LoggerTransformer.Debugf("New Secret Data for %s: %v", environment, secret)
		k8sArtifact.Secrets[secret.ObjectMeta.Name] = &secret
	}

	for _, secretData := range secretDataList {
		if secretData.Production.Enabled {
			createSecret("production", secretData.Production.Username, secretData.Production.Password, secretData.Production.APIKeyValue, secretData.Production.Type, secretData.Production.EndpointUUID)
		}
		if secretData.Sandbox.Enabled {
			createSecret("sandbox", secretData.Sandbox.Username, secretData.Sandbox.Password, secretData.Sandbox.APIKeyValue, secretData.Sandbox.Type, secretData.Sandbox.EndpointUUID)
		}
	}
}
