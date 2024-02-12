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
 * Package "synchronizer" contains artifacts relate to fetching APIs and
 * API related updates from the control plane event-hub.
 * This file contains functions to retrieve APIs and API updates.
 */

package synchronizer

import (
	"fmt"

	"archive/zip"
	"bytes"
	"strings"

	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/logging"
	sync "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/synchronizer"
	transformer "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/transformer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mapperUtil "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/mapper"
)

const (
	zipExt          string = ".zip"
	defaultCertPath string = "/home/wso2/security/controlplane.pem"
)

func init() {
	conf, _ := config.ReadConfigs()
	sync.InitializeWorkerPool(conf.ControlPlane.RequestWorkerPool.PoolSize, conf.ControlPlane.RequestWorkerPool.QueueSizePerPool,
		conf.ControlPlane.RequestWorkerPool.PauseTimeAfterFailure, conf.Agent.TrustStore.Location,
		conf.ControlPlane.SkipSSLVerification, conf.ControlPlane.HTTPClient.RequestTimeOut, conf.ControlPlane.RetryInterval,
		conf.ControlPlane.ServiceURL, conf.ControlPlane.Username, conf.ControlPlane.Password)
}

// FetchAPIsOnEvent  will fetch API from control plane during the API Notification Event
func FetchAPIsOnEvent(conf *config.Config, apiUUIDList []string, k8sClient client.Client) {
	// Populate data from config.
	envs := conf.ControlPlane.EnvironmentLabels

	// Create a channel for the byte slice (response from the APIs from control plane)
	c := make(chan sync.SyncAPIResponse)

	var queryParamMap map[string]string
	//Get API details.
	if apiUUIDList != nil {
		GetAPI(c, nil, envs, sync.APIArtifactEndpoint, true, apiUUIDList, queryParamMap)
	}
	for i := 0; i < 1; i++ {
		data := <-c
		logger.LoggerMsg.Info("Receiving data for an API")
		if data.Resp != nil {
			// Reading the root zip
			zipReader, err := zip.NewReader(bytes.NewReader(data.Resp), int64(len(data.Resp)))

			// apiFiles represents zipped API files fetched from API Manager
			apiFiles := make(map[string]*zip.File)
			// Read the .zip files within the root apis.zip and add apis to apiFiles array.
			for _, file := range zipReader.File {
				apiFiles[file.Name] = file
				logger.LoggerSync.Infof("API file found: " + file.Name)
				// Todo: Read the apis.zip and extract the api.zip,deployments.json
			}

			if err != nil {
				logger.LoggerSync.Errorf("Error while reading zip: %v", err)
				return
			}

			artifact, decodingError := transformer.DecodeAPIArtifact(data.Resp)

			if decodingError != nil {
				logger.LoggerSync.Errorf("Error while decoding the API Project Artifact: %v", decodingError)
				return
			}

			apkConf, apiUUID, revisionID, apkErr := transformer.GenerateAPKConf(artifact.APIJson, artifact.ClientCerts)

			if apkErr != nil {
				logger.LoggerSync.Errorf("Error while generating APK-Conf: %v", apkErr)
				return
			}

			logger.LoggerSync.Debugf("APK-Conf Content: %v", apkConf)

			k8ResourceEndpoint := conf.DataPlane.K8ResourceEndpoint

			deploymentDescriptor, descriptorErr := transformer.ProcessDeploymentDescriptor([]byte(artifact.DeploymentDescriptor))
			if descriptorErr != nil {
				logger.LoggerSync.Errorf("Error while decoding the Deployment Descriptor: %v", descriptorErr)
				return
			}

			crResponse, err := transformer.GenerateUpdatedCRs(apkConf, artifact.Swagger, k8ResourceEndpoint, deploymentDescriptor, artifact.APIFileName, apiUUID, fmt.Sprint(revisionID), artifact.CertMeta)
			if err != nil {
				logger.LoggerSync.Errorf("Error occured in receiving the updated CRDs: %v", err)
				return
			}

			mainZip, err := zip.NewReader(bytes.NewReader(crResponse.Bytes()), int64(crResponse.Len()))
			if err != nil {
				logger.LoggerSync.Errorf("Error creating zip reader for main zip buffer: %v", err)
				return
			}

			for _, file := range mainZip.File {
				if strings.HasSuffix(file.Name, ".zip") {
					subZipReader, err := file.Open()
					if err != nil {
						logger.LoggerSync.Errorf("Error opening sub zip file: %v", err)
						return
					}
					defer subZipReader.Close()

					var subZipBuffer bytes.Buffer
					_, err = subZipBuffer.ReadFrom(subZipReader)
					if err != nil {
						logger.LoggerSync.Errorf("Error reading sub zip file: %v", err)
						return
					}

					subZip, err := zip.NewReader(bytes.NewReader(subZipBuffer.Bytes()), int64(subZipBuffer.Len()))
					if err != nil {
						logger.LoggerSync.Errorf("Error creating zip reader for sub zip file: %v", err)
						return
					}

					for _, subFile := range subZip.File {
						mapperUtil.MapAndCreateCR(subFile, k8sClient, conf)
					}

				}
			}

			logger.LoggerMsg.Info("API applied successfully.\n")

		} else if data.ErrorCode == 204 {
			logger.LoggerMsg.Infof("No API Artifacts are available in the control plane for the envionments :%s",
				strings.Join(envs, ", "))
			//health.SetControlPlaneRestAPIStatus(true)
		} else if data.ErrorCode >= 400 && data.ErrorCode < 500 {
			logger.LoggerMsg.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error occurred when retrieving APIs from control plane(unrecoverable error): %v", data.Err.Error()),
				Severity:  logging.CRITICAL,
				ErrorCode: 1106,
			})
			//isNoAPIArtifacts := data.ErrorCode == 404 && strings.Contains(data.Err.Error(), "No Api artifacts found")
			//health.SetControlPlaneRestAPIStatus(isNoAPIArtifacts)
		} else {
			// Keep the iteration still until all the envrionment response properly.
			i--
			logger.LoggerMsg.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error occurred while fetching data from control plane: %v ..retrying..", data.Err),
				Severity:  logging.MINOR,
				ErrorCode: 1107,
			})
			//health.SetControlPlaneRestAPIStatus(false)
			sync.RetryFetchingAPIs(c, data, sync.RuntimeArtifactEndpoint, true, queryParamMap)
		}
	}
	logger.LoggerMsg.Info("Fetching API for an event is completed...")
}

// GetAPI function calls the FetchAPIs() with relevant environment labels defined in the config.
func GetAPI(c chan sync.SyncAPIResponse, id *string, envs []string, endpoint string, sendType bool, apiUUIDList []string,
	queryParamMap map[string]string) {
	if len(envs) > 0 {
		// If the envrionment labels are present, call the controle plane with labels.
		logger.LoggerAdapter.Debugf("Environment labels present: %v", envs)
		go sync.FetchAPIs(id, envs, c, endpoint, sendType, apiUUIDList, queryParamMap)
	} else {
		// If the environments are not give, fetch the APIs from default envrionment
		logger.LoggerAdapter.Debug("Environments label  NOT present. Hence adding \"default\"")
		envs = append(envs, "default")
		go sync.FetchAPIs(id, nil, c, endpoint, sendType, apiUUIDList, queryParamMap)
	}
}
