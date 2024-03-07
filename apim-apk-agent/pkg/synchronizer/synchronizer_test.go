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
package synchronizer

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// NOTE: This test case may need some refinements as there are no errors returned in the original function
func TestConstructControlPlaneRequest(t *testing.T) {
	// Dummy values for testing
	var (
		dummyID                 = "6ge3e2777gwsgwsjt2gj7wi7h72yw72yw27"
		dummyGWLabel            = []string{"test1.gw.wso2.com", "test2.gw.wso2.com"}
		dummyServiceURL         = "http://example.com"
		dummyUsername           = "testusername"
		dummyPassword           = "testpassword"
		dummyResourceEndpoint   = "test-endpoint"
		dummySendType           = true
		dummyControlPlaneParams = controlPlaneParameters{
			serviceURL:    dummyServiceURL,
			username:      dummyUsername,
			password:      dummyPassword,
			retryInterval: time.Second * 5,
		}
	)

	td := []struct {
		name             string
		id               *string
		gwLabel          []string
		controlParams    controlPlaneParameters
		resourceEndpoint string
		sendType         bool
		expectedURL      string
	}{
		{name: "ValidCase1", id: &dummyID, gwLabel: dummyGWLabel, controlParams: dummyControlPlaneParams, resourceEndpoint: dummyResourceEndpoint, sendType: dummySendType, expectedURL: "http://example.com/test-endpoint?apiId=6ge3e2777gwsgwsjt2gj7wi7h72yw72yw27&gatewayLabel=dGVzdDEuZ3cud3NvMi5jb218dGVzdDIuZ3cud3NvMi5jb20%3D&type=Envoy"},
		{name: "EmptyID", id: nil, gwLabel: dummyGWLabel, controlParams: dummyControlPlaneParams, resourceEndpoint: dummyResourceEndpoint, sendType: dummySendType, expectedURL: "http://example.com/test-endpoint?gatewayLabel=dGVzdDEuZ3cud3NvMi5jb218dGVzdDIuZ3cud3NvMi5jb20%3D&type=Envoy"},
		{name: "EmptyServiceURL", id: &dummyID, gwLabel: dummyGWLabel, controlParams: controlPlaneParameters{}, resourceEndpoint: dummyResourceEndpoint, sendType: dummySendType, expectedURL: "/test-endpoint?apiId=6ge3e2777gwsgwsjt2gj7wi7h72yw72yw27&gatewayLabel=dGVzdDEuZ3cud3NvMi5jb218dGVzdDIuZ3cud3NvMi5jb20%3D&type=Envoy"},
		{name: "EmptyGatewayLabel", id: &dummyID, gwLabel: []string{}, controlParams: dummyControlPlaneParams, resourceEndpoint: dummyResourceEndpoint, sendType: dummySendType, expectedURL: "http://example.com/test-endpoint?apiId=6ge3e2777gwsgwsjt2gj7wi7h72yw72yw27&type=Envoy"},
	}

	for _, tc := range td {
		t.Run(tc.name, func(t *testing.T) {
			req := ConstructControlPlaneRequest(tc.id, tc.gwLabel, tc.controlParams, tc.resourceEndpoint, tc.sendType)
			assert.NotNil(t, req, "Expected a non-nil request for test case: %s", tc.name)
			assert.Equal(t, tc.expectedURL, req.URL.String())
		})
	}
}

func TestReadRootFiles(t *testing.T) {
	testResourcesDir := "../../resources/test-resources/"
	files, err := os.ReadDir(testResourcesDir)
	if err != nil {
		t.Fatal("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".zip" {
			zipPath := filepath.Join(testResourcesDir, file.Name())
			zipFileBytes, err := os.ReadFile(zipPath)
			if err != nil {
				t.Logf("Error reading zip file %s: %v", file.Name(), err)
				continue
			}

			zipReader, err := zip.NewReader(bytes.NewReader(zipFileBytes), int64(len(zipFileBytes)))
			if err != nil {
				t.Logf("Error creating zip reader for file %s: %v", file.Name(), err)
				continue
			}
			deploymentDesc, apiEnvProps, err := ReadRootFiles(zipReader)
			assert.Nil(t, err)
			assert.IsType(t, &DeploymentDescriptor{}, deploymentDesc)
			assert.IsType(t, map[string]map[string]APIEnvProps{}, apiEnvProps)
			assert.Empty(t, apiEnvProps, "apiEnvProps should be empty")
			assert.NotEmpty(t, deploymentDesc.Data, "Empty data field found. Descriptor have't populated properly.")
			for _, data := range deploymentDesc.Data.Deployments {
				assert.NotEqual(t, "", data.APIFile)
			}
		}
	}
}

func TestNewWorkerPool(t *testing.T) {
	maxWorkers := 5
	jobQueueCapacity := 10
	delayForFaultRequests := time.Second * 2

	workerPool := newWorkerPool(maxWorkers, jobQueueCapacity, delayForFaultRequests)
	if cap(workerPool.internalQueue) != jobQueueCapacity {
		t.Errorf("Unexpected internal queue capacity. Got: %d, Expected: %d", cap(workerPool.internalQueue), jobQueueCapacity)
	}
	if len(workerPool.workers) != maxWorkers {
		t.Errorf("Unexpected number of workers. Got: %d, Expected: %d", len(workerPool.workers), maxWorkers)
	}
}
