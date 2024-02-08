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

	"io"
	"os"
	"strings"

	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
)

// DecodeAPIArtifact decodes a zip-encoded API payload, extracting API details like JSON, Swagger, and deployment configuration.
// Returns the APIArtifact or an error if decoding or extraction fails.
func DecodeAPIArtifact(payload []byte) (*APIArtifact, error) {
	zipReader, zipReaderError := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	var apiArtifact = &APIArtifact{}
	if zipReaderError != nil {
		logger.LoggerTransformer.Errorf("Error reading zip file: %v", zipReaderError)
		return nil, zipReaderError
	}

	// Read the zip file and get the
	for _, file := range zipReader.File {
		logger.LoggerTransformer.Info("Reading " + file.Name)
		err := readZipFile(file, apiArtifact)

		if err != nil {
			logger.LoggerTransformer.Errorf("Error reading zip file %v", err)
			return nil, err
		}
	}
	return apiArtifact, nil
}

// readZipfile will recursively go through the zip file, read and maps the content inside
// to its appropriate artifact attribute
func readZipFile(file *zip.File, apiArtifact *APIArtifact) error {
	f, fileOpenErr := file.Open()

	if fileOpenErr != nil {
		logger.LoggerTransformer.Errorf("error opening file %s in zip archieve", file.Name)
		return fileOpenErr
	}
	defer f.Close()

	content, contentError := io.ReadAll(f)

	if contentError != nil {
		logger.LoggerTransformer.Errorf("Error reading file %s in zip archieve", file.Name)
		return contentError
	}

	if strings.Contains(file.Name, ".zip") {
		apiArtifact.APIFileName = file.Name
		zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))

		if err != nil {
			logger.LoggerTransformer.Errorf("Error reading zip file: ", err)
			return err
		}

		for _, file := range zipReader.File {
			err := readZipFile(file, apiArtifact)

			if err != nil {
				logger.LoggerTransformer.Errorf("Error while reading the embedded zip file: ", err)
				return err
			}
		}
	}

	if strings.Contains(file.Name, "swagger.json") {
		apiArtifact.Swagger = string(content)

	}

	if strings.Contains(file.Name, "api.json") {
		apiArtifact.APIJson = string(content)
	}

	if strings.Contains(file.Name, "env_properties.json") {
		apiArtifact.EnvConfig = string(content)
	}

	if strings.Contains(file.Name, "deployments.json") {
		apiArtifact.DeploymentDescriptor = string(content)
	}

	if strings.Contains(file.Name, "client_certificates.json") {
		apiArtifact.ClientCerts = string(content)
	}

	return nil
}

// readZipfile will recursively go through the zip file, read and maps the content inside
// to its appropriate artifact attribute
func readAPIZipFile(file *zip.File, apiArtifact *APIArtifact) error {
	f, fileOpenErr := file.Open()

	if fileOpenErr != nil {
		logger.LoggerTransformer.Errorf("error opening file %s in zip archieve", file.Name)
		return fileOpenErr
	}
	defer f.Close()

	content, contentError := io.ReadAll(f)

	if contentError != nil {
		logger.LoggerTransformer.Errorf("Error reading file %s in zip archieve", file.Name)
		return contentError
	}

	if strings.Contains(file.Name, "swagger.json") {
		apiArtifact.Swagger = string(content)

	}

	if strings.Contains(file.Name, "api.json") {
		apiArtifact.APIJson = string(content)
	}

	if strings.Contains(file.Name, "client_certificates.json") {
		apiArtifact.ClientCerts = string(content)
	}

	return nil
}

// DecodeAPIArtifacts decodes a zip-encoded API payload, extracting API details like JSON, Swagger, and deployment configuration.
// Returns an array of APIArtifacts or an error if decoding or extraction fails.
func DecodeAPIArtifacts(payload []byte) ([]APIArtifact, error) {
	zipReader, zipReaderError := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if zipReaderError != nil {
		logger.LoggerTransformer.Errorf("Error reading zip file: %v", zipReaderError)
		return nil, zipReaderError
	}

	var apiArtifacts []APIArtifact

	var apiArtifact APIArtifact
	// Read the zip file and get the
	for _, file := range zipReader.File {
		logger.LoggerTransformer.Info("Reading " + file.Name)

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

		if strings.Contains(file.Name, "api.json") {
			apiArtifact.APIJson = string(content)
		}

		if strings.Contains(file.Name, "env_properties.json") {
			apiArtifact.EnvConfig = string(content)
		}

		if strings.Contains(file.Name, "deployments.json") {
			apiArtifact.DeploymentDescriptor = string(content)
		}
	}

	for _, file := range zipReader.File {

		if strings.Contains(file.Name, ".zip") {
			logger.LoggerTransformer.Info("Reading " + file.Name)

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
			apiArtifact.APIFileName = file.Name
			zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))

			if err != nil {
				logger.LoggerTransformer.Errorf("Error reading zip file: ", err)
				return nil, err
			}

			for _, file := range zipReader.File {
				err := readAPIZipFile(file, &apiArtifact)

				if err != nil {
					logger.LoggerTransformer.Errorf("Error while reading the embedded zip file: ", err)
					return nil, err
				}
			}
			apiArtifacts = append(apiArtifacts, apiArtifact)
		}
		apiArtifact.APIJson = ""
		apiArtifact.Swagger = ""
		apiArtifact.RevisionID = 0
		apiArtifact.APIFileName = ""
	}

	return apiArtifacts, nil
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

// readZipFile reads the zip file and returns the file bytes
func getZipFileBytes(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// SaveZipToFile saves the given zipfile to the defined location
func SaveZipToFile(zipContent *bytes.Buffer, filePath string) error {
	zipFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	parentWriter, err := zipWriter.Create("parent.zip")
	if err != nil {
		return err
	}

	_, err = parentWriter.Write(zipContent.Bytes())
	if err != nil {
		return err
	}

	return nil
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
