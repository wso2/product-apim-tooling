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

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"

	"io"
	"strings"

	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
)

// DecodeAPIArtifact decodes a zip-encoded API payload, extracting API details like JSON, Swagger, and deployment configuration.
// Returns the APIArtifact or an error if decoding or extraction fails.
func DecodeAPIArtifact(apiZip *zip.File) (*APIArtifact, error) {
	logger.LoggerTransformer.Info("Reading " + apiZip.Name)
	apiArtifact, err := readZipFile(apiZip)
	if err != nil {
		logger.LoggerTransformer.Errorf("Error reading zip file %v", err)
		return nil, err
	}
	return apiArtifact, nil
}

// ReadContent read the content of file
func ReadContent(file *zip.File) ([]byte, error) {
	f, fileOpenErr := file.Open()
	if fileOpenErr != nil {
		logger.LoggerTransformer.Errorf("error opening file %s in zip archieve", file.Name)
		return nil, fileOpenErr
	}
	defer f.Close()
	content, contentError := io.ReadAll(f)
	if contentError != nil {
		logger.LoggerTransformer.Errorf("Error reading file %s in zip archieve", file.Name)
		return nil, contentError
	}
	return content, nil
}

// readZipfile will recursively go through the zip file, read and maps the content inside
// to its appropriate artifact attribute
func readZipFile(file *zip.File) (*APIArtifact, error) {
	var apiArtifact = &APIArtifact{}

	content, err := ReadContent(file)
	if err != nil {
		return nil, err
	}
	apiArtifact.APIFileName = file.Name
	zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		logger.LoggerTransformer.Errorf("Error reading zip file: ", err)
		return nil, err
	}
	for _, file := range zipReader.File {
		if strings.Contains(file.Name, "swagger.json") {
			openAPIContent, err := ReadContent(file)
			if err != nil {
				return nil, err
			}
			apiArtifact.Schema = string(openAPIContent)
		}

		if strings.Contains(file.Name, "schema.graphql") {
			graphqlContent, err := ReadContent(file)
			if err != nil {
				return nil, err
			}
			apiArtifact.Schema = string(graphqlContent)
		}

		if strings.Contains(file.Name, "api.json") {
			apiJSON, err := ReadContent(file)
			if err != nil {
				return nil, err
			}
			apiArtifact.APIJson = string(apiJSON)
		}
		if strings.Contains(file.Name, "client_certificates.json") {
			certificateJSON, err := ReadContent(file)
			if err != nil {
				return nil, err
			} else if string(certificateJSON) == "" {
				return nil, errors.New("empty Client_Certificate content detected")
			}
			apiArtifact.CertArtifact.ClientCerts = string(certificateJSON)
			apiArtifact.CertMeta.CertAvailable = true
		}

		if strings.Contains(file.Name, "endpoint_certificates.json") {
			endpointCertificateJSON, err := ReadContent(file)
			if err != nil {
				return nil, err
			} else if string(endpointCertificateJSON) == "" {
				return nil, errors.New("empty Endpoint_Certificate content detected")
			}
			apiArtifact.CertArtifact.EndpointCerts = string(endpointCertificateJSON)
			apiArtifact.EndpointCertMeta.CertAvailable = true
		}

		if strings.Contains(file.Name, ".crt") {
			certificateData, err := ReadContent(file)
			if err != nil {
				return nil, err
			}
			//NOTE:There is an issue in reading the certificate content. Even the same logic is been used, the
			// .crt files in the Client-certificates gets parsed in base64 encoded version while the cert files in
			// Endpoint-certtificate folder gets parsed as original value
			pathSegments := strings.Split(file.Name, "/")
			if strings.Contains(file.Name, "Client-certificates") {
				if apiArtifact.CertMeta.ClientCertFiles == nil {
					apiArtifact.CertMeta.ClientCertFiles = make(map[string]string)
				}
				apiArtifact.CertMeta.ClientCertFiles[pathSegments[len(pathSegments)-1]] = string(certificateData)
			}
			if strings.Contains(file.Name, "Endpoint-certificates") {
				if apiArtifact.EndpointCertMeta.EndpointCertFiles == nil {
					apiArtifact.EndpointCertMeta.EndpointCertFiles = make(map[string]string)
				}
				apiArtifact.EndpointCertMeta.EndpointCertFiles[pathSegments[len(pathSegments)-1]] = string(certificateData)
			}

		}

	}
	return apiArtifact, nil
}

// ProcessDeploymentDescriptor processes the artifact and returns the deployment descriptor struct
func ProcessDeploymentDescriptor(deploymentDescriptor []byte) (*DeploymentDescriptor, error) {
	var deployment DeploymentDescriptor

	umErr := json.Unmarshal([]byte(deploymentDescriptor), &deployment)
	if umErr != nil {
		return nil, umErr
	}
	return &deployment, nil
}

// StringExists checks for the existance of a particular string a string slice
func StringExists(target string, slice []string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, exists := set[target]
	return exists
}
