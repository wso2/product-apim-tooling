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
	"fmt"
	"io"
	"os"
	"strings"

	loggers "github.com/sirupsen/logrus"
)

const clientUtils = "client_utils"

// DecodeAPIArtifact decodes a zip-encoded API payload, extracting API details like JSON, Swagger, and deployment configuration.
// Returns the APIArtifact or an error if decoding or extraction fails.
func DecodeAPIArtifact(payload []byte) (*APIArtifact, error) {
	zipReader, zipReaderError := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	var apiArtifact = &APIArtifact{}
	if zipReaderError != nil {
		loggers.Error("Error reading zip file", zipReaderError.Error())
		return nil, zipReaderError
	}

	// Read the zip file and get the
	for _, file := range zipReader.File {
		loggers.Info("Reading " + file.Name)
		err := readZipFile(file, apiArtifact)

		if err != nil {
			loggers.Error("Error reading zip file", err.Error())
			return nil, err
		}
	}
	return apiArtifact, nil
}

func readZipFile(file *zip.File, apiArtifact *APIArtifact) error {
	f, fileOpenErr := file.Open()

	if fileOpenErr != nil {
		//logger.GetLogger(ctx, clientUtils).Error(fmt.Sprintf("error opening file %s in zip archieve", file.Name), fileOpenErr.Error())
		return fmt.Errorf("error opening file %s in zip archieve", file.Name)
	}
	defer f.Close()

	content, contentError := io.ReadAll(f)

	if contentError != nil {
		loggers.Error(fmt.Sprintf("error reading file %s in zip archieve", file.Name), contentError.Error())
		return fmt.Errorf("error reading file %s in zip archieve", file.Name)
	}

	if strings.Contains(file.Name, ".zip") {
		// revisionID := strings.Split(strings.Split(file.Name, ".zip")[0], "-")[1]
		// apiArtifact.RevisionID = revId
		zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))

		if err != nil {
			loggers.Error("Error reading zip file", err.Error())
			return fmt.Errorf("error reading zip file %s", file.Name)
		}

		for _, file := range zipReader.File {
			er := readZipFile(file, apiArtifact)

			if er != nil {
				loggers.Error("Error while reading the embedded zip file", er.Error())
				return fmt.Errorf("error reading embedded zip file %s", file.Name)
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
	return nil
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

// Prints the content inside the files in a given zipfile
func printZipContents(zipBytes []byte) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return err
	}

	for _, zipFile := range zipReader.File {
		fmt.Printf("\nFile Name: %s\n", zipFile.Name)

		fileContent, err := getZipFileBytes(zipFile)
		if err != nil {
			return err
		}

		fmt.Printf("File Content:\n%s\n", fileContent)
	}

	return nil
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
