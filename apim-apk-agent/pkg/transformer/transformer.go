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
	"strings"

	"io"
	"mime/multipart"
	"net/http"

	"github.com/wso2/apk/common-go-libs/utils"

	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"

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
	} else {
		logger.LoggerTransformer.Warn("Warn:client_cert.json empty or not exist for the given zip.")
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
	case "GraphQL":
		apiType = "GraphQL"
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

// GenerateUpdatedCRs takes the .apk-conf, api definition, vHost and the organization for a particular API and then generate and returns
// the relavant CRD set as a zip
func GenerateUpdatedCRs(apkConf string, apiDefinition string, k8ResourceGenEndpoint string, deploymentDescriptor *DeploymentDescriptor, apiFileName string, apiID string, revisionID string, certMeta CertMetadata) (*bytes.Buffer, error) {
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

	var allZipsBuffer bytes.Buffer
	combinedZip := zip.NewWriter(&allZipsBuffer)

	for _, deployment := range *deploymentDescriptor.Data.Deployments {
		if deployment.APIFile == apiFileName {
			for _, environment := range *deployment.Environments {

				modifiedZip, err := transformCRD(body, environment.Vhost, deployment.OrganizationID, apiID, revisionID, certMeta)

				if err != nil {
					logger.LoggerTransformer.Error("Unable to transform the initial CRDs:", err)
					return nil, err
				}

				// Write the modifiedZipData to the combined zip
				fileName := fmt.Sprintf("%s_%s.zip", deployment.OrganizationID, environment.Vhost)
				writer, err := combinedZip.Create(fileName)
				if err != nil {
					logger.LoggerTransformer.Error("Error creating zip file entry:", err)
					return nil, err
				}

				_, err = writer.Write(modifiedZip)
				if err != nil {
					logger.LoggerTransformer.Error("Error writing  to the zip file:", err)
					return nil, err
				}

			}
		}
	}

	// Close the combined zip
	err = combinedZip.Close()
	if err != nil {
		logger.LoggerTransformer.Error("Error closing zip file:", err)
		return nil, err
	}

	return &allZipsBuffer, nil
}

// transformCRD converts the APK CRDs and returns the modified CRDs with modified
func transformCRD(crdZip []byte, vHost string, organization string, apiID string, revisionID string, certMeta CertMetadata) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(crdZip), int64(len(crdZip)))
	if err != nil {
		logger.LoggerTransformer.Fatal(err)
	}

	//create a new zip writer
	var modifiedZipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&modifiedZipBuffer)

	defer zipWriter.Close()

	namespace := utils.GetOperatorPodNamespace()

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		logger.LoggerTransformer.Debugf("Reading file: %s", zipFile.Name)
		apkCRDFileBytes, err := getZipFileBytes(zipFile)
		if err != nil {
			logger.LoggerTransformer.Error(err)
			continue
		}

		_ = apkCRDFileBytes // this is unzipped file bytes
		yamlCrd, err := generateAPKCrdsFromYaml(apkCRDFileBytes, organization, vHost, namespace, apiID, revisionID)
		if err != nil {
			logger.LoggerTransformer.Error("Error occured while retrieving the modified CRDs", err)
			return nil, err
		}

		// Create a new file in the modified zip with the same name
		modifiedFile, err := zipWriter.Create(zipFile.Name)
		if err != nil {
			logger.LoggerTransformer.Error("Error in creating new file in the modified zip", err)
			return nil, err
		}

		// Write the modified content to the new file in the modified zip
		_, err = modifiedFile.Write(yamlCrd)
		if err != nil {
			logger.LoggerTransformer.Error("Error in writing modified content to the new file", err)
			return nil, err
		}
	}

	// Create ConfigMap to store the cert data if mTLS has enabled
	if certMeta.CertAvailable {
		for confKey, confValue := range certMeta.ClientCertFiles {
			i := 0
			i++
			pathSegments := strings.Split(confKey, ".")
			configName := pathSegments[0]
			cm := createCongigMap(configName, confKey, confValue)
			logger.LoggerTransformer.Infof("New ConfigMap Data:", string(cm))

			// Create a new yaml file for the ConfigMap yaml
			cmFileName := fmt.Sprintf(apiID, "-configmap", "-", i, ".yaml")
			// Write the new yaml file to the modified zip
			modifiedFile, err := zipWriter.Create(cmFileName)
			if err != nil {
				logger.LoggerTransformer.Error("Error in creating new file in the modified zip", err)
				return nil, err
			}

			_, err = modifiedFile.Write(cm)
			if err != nil {
				logger.LoggerTransformer.Error("Error in writing modified content to the new file", err)
				return nil, err
			}
		}

	}

	// Finish writing the modified zip file
	err = zipWriter.Close()
	if err != nil {
		logger.LoggerTransformer.Error("Error occured in closing the zip with modified files", err)
		return nil, err
	}

	return modifiedZipBuffer.Bytes(), nil

}

// generateAPKCrdsFromYaml processes the returned APK CRD yaml, replaces the vhost, adds the organization
// and namespace and returns the json
func generateAPKCrdsFromYaml(crdYaml []byte, orgUUID, vhost string, namespace string, apiID string, revisionID string) ([]byte, error) {
	var crdYml map[interface{}]interface{}
	unMarshalErr := yaml.Unmarshal(crdYaml, &crdYml)

	if unMarshalErr != nil {
		return nil, unMarshalErr
	}
	replaceVhost(crdYml, vhost)
	addOrganization(crdYml, orgUUID)
	addNamespace(crdYml, namespace)
	addRevisionAndAPIUUID(crdYml, apiID, revisionID)

	processdCrdYml := convertMap(crdYml)

	yamlCrd, err := yaml.Marshal(processdCrdYml)
	if err != nil {
		return nil, err
	}
	return yamlCrd, nil
}

// ConvertMap recursively converts a map[interface{}]interface{} to map[string]interface{}
func convertMap(inputMap map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range inputMap {
		keyStr, ok := k.(string)
		if !ok {
			// If the key is not a string, try to convert it
			keyStr = fmt.Sprintf("%v", k)
		}

		switch value := v.(type) {
		case map[interface{}]interface{}:
			// If the value is a map, recursively convert it
			result[keyStr] = convertMap(value)
		case []interface{}:
			// If the value is an array, convert each element recursively
			var arr []interface{}
			for _, elem := range value {
				if childMap, ok := elem.(map[interface{}]interface{}); ok {
					arr = append(arr, convertMap(childMap))
				} else {
					arr = append(arr, elem)
				}
			}
			result[keyStr] = arr
		case string:
			result[keyStr] = value
		default:
			// Otherwise, keep the value as is
			result[keyStr] = v
		}
	}

	return result
}

// replaceVhost will take the httpRoute CR and replace the default vHost with the one passed inside
// the deploymemt descriptor
func replaceVhost(inputMap map[interface{}]interface{}, vhost string) {
	if kind, ok := inputMap[k8sKindField].(string); ok && kind == k8sKindHTTPRoute {
		if spec, ok := inputMap[k8sSpecField].(map[interface{}]interface{}); ok {
			if hostnames, ok := spec[k8sHostnamesField].([]interface{}); ok {
				hostnames[0] = vhost
			}
		}
	}
}

// addOrganization will take the API CR and change the organization to the one passed inside
// the deploymemt descriptor
func addOrganization(inputMap map[interface{}]interface{}, organization string) {
	if kind, ok := inputMap[k8sKindField].(string); ok && kind == k8sKindAPI {
		if spec, ok := inputMap[k8sSpecField].(map[interface{}]interface{}); ok {
			if _, ok := spec[k8sOrganizationField]; ok {
				spec[k8sOrganizationField] = organization
			}
		}
	}
	organizationHash := generateSHA1Hash(organization)
	if metadata, ok := inputMap[k8sMetadataField].(map[interface{}]interface{}); ok {
		if labels, ok := metadata[k8sLabelsField].(map[interface{}]interface{}); ok {
			if _, ok := labels[k8sOrganizationField]; ok {
				labels[k8sOrganizationField] = organizationHash
			}
		}
	}
}

// addNamespace will set the namespace attribute in the CR to the pods currently operating namespace
func addNamespace(inputMap map[interface{}]interface{}, namespace string) {
	if metadata, ok := inputMap[k8sMetadataField].(map[interface{}]interface{}); ok {
		metadata[k8sNamespaceField] = namespace
	}
}

// addRevisionAndAPIUUID will add the API ID and the revision field attributes to the API CR
func addRevisionAndAPIUUID(inputMap map[interface{}]interface{}, apiID string, revisionID string) {
	if kind, ok := inputMap[k8sKindField].(string); ok && kind == k8sKindAPI {
		if metadata, ok := inputMap[k8sMetadataField].(map[interface{}]interface{}); ok {
			if labels, ok := metadata[k8sLabelsField].(map[interface{}]interface{}); ok {
				labels[k8APIUuidField] = apiID
				labels[k8RevisionField] = revisionID
			}
		}
	}
}

// generateSHA1Hash returns the SHA1 hash for the given string
func generateSHA1Hash(input string) string {
	h := sha1.New() /* #nosec */
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// createConfigMap returns a marshalled yaml of ConfigMap kind after adding the given values
func createCongigMap(configName, dataKey, dataValue string) []byte {
	cm := CertConfigYaml{}
	cm.APIVersion = "v1"
	cm.Kind = "ConfigMap"
	cm.Metadata = MetadataBlock{
		Name:      configName,
		Namespace: "apk-integration-test",
	}
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
	cm.Data[dataKey] = dataValue
	configYaml, marshalErr := yaml.Marshal(cm)
	if marshalErr != nil {
		logger.LoggerTransformer.Errorf("Error occured while marshalling the config yaml")
	}
	return configYaml
}
