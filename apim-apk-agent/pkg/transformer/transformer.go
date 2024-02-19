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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"io"
	"mime/multipart"
	"net/http"

	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	k8Yaml "sigs.k8s.io/yaml"

	"gopkg.in/yaml.v2"
)

// GenerateAPKConf will Generate the mapped .apk-conf file for a given API Project zip
func GenerateAPKConf(APIJson string, clientCerts string) (string, string, uint32, error) {

	apk := &API{}

	var apiYaml APIYaml

	apiYamlError := json.Unmarshal([]byte(APIJson), &apiYaml)

	if apiYamlError != nil {
		logger.LoggerTransformer.Error("Error while unmarshalling api.json content", apiYamlError)
		return "", "null", 0, apiYamlError
	}

	apiYamlData := apiYaml.Data

	apk.Name = apiYamlData.Name
	apk.Context = apiYamlData.Context
	apk.Version = apiYamlData.Version
	apk.Type = getAPIType(apiYamlData.Type)
	apk.DefaultVersion = apiYamlData.DefaultVersion
	apk.DefinitionPath = "/definition"
	apk.SubscriptionValidation = true
	apkOperations := make([]Operation, len(apiYamlData.Operations))

	for i, operation := range apiYamlData.Operations {

		reqPolicyCount := len(operation.OperationPolicies.Request)
		resPolicyCount := len(operation.OperationPolicies.Response)
		reqInterceptor, resInterceptor := getReqAndResInterceptors(reqPolicyCount, resPolicyCount)

		op := &Operation{
			Target:  operation.Target,
			Verb:    operation.Verb,
			Scopes:  operation.Scopes,
			Secured: true,
			OperationPolicies: &OperationPolicies{
				Request:  *reqInterceptor,
				Response: *resInterceptor,
			},
		}
		apkOperations[i] = *op
	}

	apk.Operations = &apkOperations

	apk.EndpointConfigurations = &EndpointConfiguration{
		// For private PPDPs, we need to treat the token type to be SANDBOX as it is tested by developers.
		Sandbox: &Endpoint{
			Endpoint: apiYamlData.EndpointConfig.SandboxEndpoints.URL},
		Production: &Endpoint{
			Endpoint: apiYamlData.EndpointConfig.ProductionEndpoints.URL},
	}

	var certList CertDescriptor
	certAvailable := false

	if clientCerts != "" {
		certErr := json.Unmarshal([]byte(clientCerts), &certList)
		if certErr != nil {
			logger.LoggerTransformer.Errorf("Error while unmarshalling client_cert.json content: ", apiYamlError)
			return "", "null", 0, certErr
		}
		certAvailable = true
	}

	authConfigList := mapAuthConfigs(apiYamlData.AuthorizationHeader, apiYamlData.SecuritySchemes, certAvailable, certList)
	apk.Authentication = &authConfigList

	apk.CorsConfig = &apiYamlData.CORSConfiguration

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
		return "", "null", 0, marshalError
	}
	return string(c), apiYamlData.RevisionedAPIID, apiYamlData.RevisionID, nil
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
func getReqAndResInterceptors(reqPolicyCount int, resPolicyCount int) (*[]OperationPolicy, *[]OperationPolicy) {
	var reqInterceptor, resInterceptor []OperationPolicy
	var interceptorParams *InterceptorService
	var interceptorPolicy OperationPolicy

	if reqPolicyCount > 0 || resPolicyCount > 0 {
		interceptorParams = &InterceptorService{
			BackendURL:      "https://interceptor-svc.ns:9081",
			HeadersEnabled:  true,
			BodyEnabled:     true,
			TrailersEnabled: true,
			ContextEnabled:  true,
			TLSSecretName:   "interceptor-cert",
			TLSSecretKey:    "ca.crt",
		}

		// Create an instance of OperationPolicy
		interceptorPolicy = OperationPolicy{
			PolicyName:    "Interceptor",
			PolicyVersion: "v1",
			Parameters:    interceptorParams,
		}
	}

	if reqPolicyCount > 0 {
		reqInterceptor = append(reqInterceptor, interceptorPolicy)
	}

	if resPolicyCount > 0 {
		resInterceptor = append(resInterceptor, interceptorPolicy)

	}

	return &reqInterceptor, &resInterceptor
}

// mapAuthConfigs will take the security schemes as the parameter and will return the mapped auth configs to be
// added into the apk-conf
func mapAuthConfigs(authHeader string, secSchemes []string, certAvailable bool, certList CertDescriptor) []AuthConfiguration {
	var authConfigs []AuthConfiguration
	if StringExists("oauth2", secSchemes) {
		var newConfig AuthConfiguration
		newConfig.AuthType = "OAuth2"
		newConfig.Enabled = true
		newConfig.HeaderName = authHeader
		if StringExists("oauth_basic_auth_api_key_mandatory", secSchemes) {
			newConfig.Required = "mandatory"
		} else {
			newConfig.Required = "optional"
		}

		authConfigs = append(authConfigs, newConfig)
	}
	if StringExists("mutualssl", secSchemes) && certAvailable {
		var newConfig AuthConfiguration
		newConfig.AuthType = "mTLS"
		newConfig.Enabled = true
		if StringExists("mutualssl_mandatory", secSchemes) {
			newConfig.Required = "mandatory"
		} else {
			newConfig.Required = "optional"
		}

		clientCerts := make([]Certificate, len(certList.CertData))

		for i, cert := range certList.CertData {
			prop := &Certificate{
				Name: cert.Alias,
				Key:  cert.Certificate,
			}
			clientCerts[i] = *prop
		}
		newConfig.Certificates = clientCerts
		authConfigs = append(authConfigs, newConfig)
	}
	return authConfigs
}

// GenerateCRs takes the .apk-conf, api definition, vHost and the organization for a particular API and then generate and returns
// the relavant CRD set as a zip
func GenerateCRs(apkConf string, apiDefinition string, k8ResourceGenEndpoint string) (*K8sArtifacts, error) {
	k8sArtifact := K8sArtifacts{HTTPRoutes: make(map[string]*gwapiv1b1.HTTPRoute), GQLRoutes: make(map[string]*dpv1alpha2.GQLRoute), Backends: make(map[string]*dpv1alpha1.Backend), Scopes: make(map[string]*dpv1alpha1.Scope), Authentication: make(map[string]*dpv1alpha2.Authentication), APIPolicies: make(map[string]*dpv1alpha2.APIPolicy), InterceptorServices: make(map[string]*dpv1alpha1.InterceptorService), ConfigMaps: make(map[string]*corev1.ConfigMap), Secrets: make(map[string]*corev1.Secret), RateLimitPolicies: make(map[string]*dpv1alpha1.RateLimitPolicy)}
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

	// Create the HTTP request
	request, err := http.NewRequest(postHTTPMethod, k8ResourceGenEndpoint, &requestBody)
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
			var apiPolicy dpv1alpha2.APIPolicy
			err = k8Yaml.Unmarshal(yamlData, &apiPolicy)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling APIPolicy YAML: %v", err)
				continue
			}
			k8sArtifact.APIPolicies[apiPolicy.ObjectMeta.Name] = &apiPolicy
		case "HTTPRoute":
			var httpRoute gwapiv1b1.HTTPRoute
			err = k8Yaml.Unmarshal(yamlData, &httpRoute)
			if err != nil {
				logger.LoggerSync.Errorf("Error unmarshaling HTTPRoute YAML: %v", err)
				continue
			}
			k8sArtifact.HTTPRoutes[httpRoute.ObjectMeta.Name] = &httpRoute

		case "Backend":
			var backend dpv1alpha1.Backend
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
			var api dpv1alpha2.API
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
	return &k8sArtifact, nil
}

// UpdateCRS cr update
func UpdateCRS(k8sArtifact *K8sArtifacts, environments *[]Environment, organizationID string, apiUUID string, revisionID string, namespace string) {
	addOrganization(k8sArtifact, organizationID)
	addRevisionAndAPIUUID(k8sArtifact, apiUUID, revisionID)
	for _, environment := range *environments {
		replaceVhost(k8sArtifact, environment.Vhost, environment.Type)
	}
}

// replaceVhost will take the httpRoute CR and replace the default vHost with the one passed inside
// the deploymemt descriptor
func replaceVhost(k8sArtifact *K8sArtifacts, vhost string, deploymentType string) {
	if deploymentType == "hybrid" {
		// append sandbox. part to available vhost to generate sandbox vhost
		if k8sArtifact.API.Spec.Production != nil {
			for _, routeName := range k8sArtifact.API.Spec.Production {
				for _, routes := range routeName.HTTPRouteRefs {
					httprouteRef, ok := k8sArtifact.HTTPRoutes[routes]
					if ok {
						httprouteRef.Spec.Hostnames = []gwapiv1b1.Hostname{gwapiv1b1.Hostname(vhost)}
					}
				}
			}
		}
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.HTTPRouteRefs {
					httprouteRef, ok := k8sArtifact.HTTPRoutes[routes]
					if ok {
						httprouteRef.Spec.Hostnames = []gwapiv1b1.Hostname{gwapiv1b1.Hostname("sandbox." + vhost)}
					}
				}
			}
		}
	} else if deploymentType == "sandbox" {
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.HTTPRouteRefs {
					httprouteRef, ok := k8sArtifact.HTTPRoutes[routes]
					if ok {
						httprouteRef.Spec.Hostnames = []gwapiv1b1.Hostname{gwapiv1b1.Hostname(vhost)}
					}
				}
			}
		}
		if k8sArtifact.API.Spec.Production != nil {
			for _, routeName := range k8sArtifact.API.Spec.Production {
				for _, routes := range routeName.HTTPRouteRefs {
					delete(k8sArtifact.HTTPRoutes, routes)
				}
			}
			k8sArtifact.API.Spec.Production = []dpv1alpha2.EnvConfig{}
		}
	} else {
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.HTTPRouteRefs {
					httprouteRef, ok := k8sArtifact.HTTPRoutes[routes]
					if ok {
						httprouteRef.Spec.Hostnames = []gwapiv1b1.Hostname{gwapiv1b1.Hostname(vhost)}
					}
				}
			}
		}
		if k8sArtifact.API.Spec.Sandbox != nil {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.HTTPRouteRefs {
					delete(k8sArtifact.HTTPRoutes, routes)
				}
			}
			k8sArtifact.API.Spec.Sandbox = []dpv1alpha2.EnvConfig{}
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

// generateSHA1Hash returns the SHA1 hash for the given string
func generateSHA1Hash(input string) string {
	h := sha1.New() /* #nosec */
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
