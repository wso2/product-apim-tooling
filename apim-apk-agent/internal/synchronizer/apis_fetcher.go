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
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/health"

	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	sync "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/synchronizer"
)

const (
	zipExt          string = ".zip"
	defaultCertPath string = "/home/wso2/security/controlplane.pem"
)

func init() {
	conf, _ := config.ReadConfigs()
	sync.InitializeWorkerPool(conf.ControlPlane.RequestWorkerPool.PoolSize, conf.ControlPlane.RequestWorkerPool.QueueSizePerPool,
		conf.ControlPlane.RequestWorkerPool.PauseTimeAfterFailure, conf.Adapter.Truststore.Location,
		conf.ControlPlane.SkipSSLVerification, conf.ControlPlane.HTTPClient.RequestTimeOut, conf.ControlPlane.RetryInterval,
		conf.ControlPlane.ServiceURL, conf.ControlPlane.Username, conf.ControlPlane.Password)
}

// FetchAPIsFromControlPlane method pulls API data for a given APIs according to a
// given API ID and a list of environments that API has been deployed to.
// updatedAPIID is the corresponding ID of the API in the form of an UUID
// updatedEnvs contains the list of environments the API deployed to.
func FetchAPIsFromControlPlane(updatedAPIID string, updatedEnvs []string) {
	// Read configurations and derive the eventHub details
	conf, errReadConfig := config.ReadConfigs()
	if errReadConfig != nil {
		// This has to be error. For debugging purpose info
		logger.LoggerSync.Errorf("Error reading configs: %v", errReadConfig)
	}
	// Populate data from config.
	configuredEnvs := conf.ControlPlane.EnvironmentLabels
	//finalEnvs contains the actual envrionments that the adapter should update
	var finalEnvs []string
	if len(configuredEnvs) > 0 {
		// If the configuration file contains environment list, then check if then check if the
		// affected environments are present in the provided configs. If so, add that environment
		// to the finalEnvs slice
		for _, updatedEnv := range updatedEnvs {
			for _, configuredEnv := range configuredEnvs {
				if updatedEnv == configuredEnv {
					finalEnvs = append(finalEnvs, updatedEnv)
				}
			}
		}
	} else {
		// If the labels are not configured, publish the APIS to the default environment
		finalEnvs = []string{config.DefaultGatewayName}
	}

	if len(finalEnvs) == 0 {
		// If the finalEnvs is empty -> it means, the configured envrionments  does not contains the affected/updated
		// environments. If that's the case, then APIs should not be fetched from the adapter.
		return
	}

	c := make(chan sync.SyncAPIResponse)
	logger.LoggerSync.Infof("API %s is added/updated to APIList for label %v", updatedAPIID, updatedEnvs)
	var queryParamMap map[string]string

	go sync.FetchAPIs(&updatedAPIID, finalEnvs, c, sync.RuntimeArtifactEndpoint, true, nil, queryParamMap)
	for {
		data := <-c
		logger.LoggerSync.Infof("Receiving data for the API: %q", updatedAPIID)
		if data.Resp != nil {
			// For successfull fetches, data.Resp would return a byte slice with API project(s)
			logger.LoggerSync.Infof("API Project %q", data.Resp)
			// err := PushAPIProjects(data.Resp, finalEnvs)
			// if err != nil {
			// 	logger.LoggerSync.Errorf("Error occurred while pushing API data for the API %q: %v ", updatedAPIID, err)
			// }
			break
		} else if data.ErrorCode >= 400 && data.ErrorCode < 500 {
			logger.LoggerSync.Errorf("Error occurred when retrieving API %q from control plane: %v", updatedAPIID, data.Err)
			health.SetControlPlaneRestAPIStatus(false)
		} else {
			// Keep the iteration still until all the envrionment response properly.
			logger.LoggerSync.Errorf("Error occurred while fetching data from control plane for the API %q: %v. Hence retrying..", updatedAPIID, data.Err)
			sync.RetryFetchingAPIs(c, data, sync.RuntimeArtifactEndpoint, true, queryParamMap)
		}
	}

}
