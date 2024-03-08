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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// Define testResourcesDir
var testResourcesDir = "../../resources/test-resources/"

// Read HTTPk8Json file
var httpFilePath = filepath.Join(testResourcesDir, "httpk8Json.json")
var httpBytes, _ = os.ReadFile(httpFilePath)
var HTTPk8Json = string(httpBytes)

// Read GQLk8Json file
var gqlFilePath = filepath.Join(testResourcesDir, "gqlk8Json.json")
var gqlBytes, _ = os.ReadFile(gqlFilePath)
var GQLk8Json = string(gqlBytes)

var sampleK8Artifacts = []string{HTTPk8Json, GQLk8Json}

func TestAPIArtifactDecoding(t *testing.T) {
	apiFiles := make(map[string]*zip.File)
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

			for _, file := range zipReader.File {
				apiFiles[file.Name] = file
			}
			if err != nil {
				t.Errorf("Error while reading zip: %v", err)

			}
			deploymentJSON, exists := apiFiles["deployments.json"]
			if !exists {
				t.Errorf("deployments.json not found")

			}
			deploymentJSONBytes, err := ReadContent(deploymentJSON)
			assert.NotNil(t, deploymentJSONBytes)
			assert.NoError(t, err)
			assert.IsType(t, []byte{}, deploymentJSONBytes)

			deploymentDescriptor, err := ProcessDeploymentDescriptor(deploymentJSONBytes)

			assert.NotNil(t, deploymentDescriptor)
			assert.NoError(t, err)
			assert.IsType(t, &DeploymentDescriptor{}, deploymentDescriptor)
			apiDeployments := deploymentDescriptor.Data.Deployments
			if apiDeployments != nil {
				for _, apiDeployment := range *apiDeployments {
					apiZip, exists := apiFiles[apiDeployment.APIFile]
					if exists {
						artifact, decodingError := DecodeAPIArtifact(apiZip)
						if decodingError != nil {
							t.Errorf("Error while decoding the API Project Artifact: %v", decodingError)

						}
						assert.NotNil(t, artifact)
						assert.NoError(t, err)
						assert.IsType(t, &APIArtifact{}, artifact)
					}
				}
			}
		}
	}
}

func TestAPKConfGeneration(t *testing.T) {
	testResourcesDir := "../../resources/test-resources/Base/"
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

			for _, zipFile := range zipReader.File {
				apiArtifact, err := DecodeAPIArtifact(zipFile)
				if err != nil {
					t.Logf("Error decoding API artifact from %s: %v", zipFile.Name, err)
					continue
				}
				assert.NotNil(t, apiArtifact)
				assert.NoError(t, err)
				assert.IsType(t, &APIArtifact{}, apiArtifact)

				apkConf, apiUUID, revisionID, apkErr := GenerateAPKConf(apiArtifact.APIJson, apiArtifact.CertArtifact)

				assert.NoError(t, apkErr)
				assert.NotEmpty(t, apkConf)
				assert.NotEqual(t, "null", apiUUID)
				assert.NotEqual(t, uint32(0), revisionID)
			}
		}
	}
}

func TestAddRevisionAndAPIUUID(t *testing.T) {
	for _, k8Json := range sampleK8Artifacts {
		var k8sArtifact K8sArtifacts
		err := json.Unmarshal([]byte(k8Json), &k8sArtifact)
		if err != nil {
			t.Error("Unable to unmarshal the dummy k8artifact")
		}

		// Call the function multiple times with different apiID and revisionID combinations
		addRevisionAndAPIUUID(&k8sArtifact, "c1f22832-2859-4849-b6c1-9a48644dedc9", "revisionID1")
		addRevisionAndAPIUUID(&k8sArtifact, "0d822a52-dde1-45e9-8896-203dfe1f4d22", "revisionID2")
		addRevisionAndAPIUUID(&k8sArtifact, "dcef31c3-aa9b-447b-8907-90bb05000801", "revisionID3")

		// Check whether the labels are properly populated
		assert.Equal(t, "dcef31c3-aa9b-447b-8907-90bb05000801", k8sArtifact.API.ObjectMeta.Labels[k8APIUuidField])
		assert.Equal(t, "revisionID3", k8sArtifact.API.ObjectMeta.Labels[k8RevisionField])
	}
}

func TestAddOrganization(t *testing.T) {
	for _, k8Json := range sampleK8Artifacts {
		var k8sArtifact K8sArtifacts
		err := json.Unmarshal([]byte(k8Json), &k8sArtifact)
		if err != nil {
			t.Error("Unable to unmarshal the dummy k8artifact")
		}
		// Call the function with organization
		organization := "test.org.wso2.com"
		addOrganization(&k8sArtifact, organization)

		// Generate SHA1 hash of the organization
		organizationHash := generateSHA1Hash(organization)

		// Check whether the labels are properly populated for API and other artifacts
		assert.Equal(t, organizationHash, k8sArtifact.API.ObjectMeta.Labels[k8sOrganizationField])
		for _, apiPolicy := range k8sArtifact.APIPolicies {
			assert.Equal(t, organizationHash, apiPolicy.ObjectMeta.Labels[k8sOrganizationField])
		}
		for _, httpRoute := range k8sArtifact.HTTPRoutes {
			assert.Equal(t, organizationHash, httpRoute.ObjectMeta.Labels[k8sOrganizationField])
		}
		for _, gqlRoute := range k8sArtifact.GQLRoutes {
			assert.Equal(t, organizationHash, gqlRoute.ObjectMeta.Labels[k8sOrganizationField])
		}
		for _, authentication := range k8sArtifact.Authentication {
			assert.Equal(t, organizationHash, authentication.ObjectMeta.Labels[k8sOrganizationField])
		}
		for _, backend := range k8sArtifact.Backends {
			assert.Equal(t, organizationHash, backend.ObjectMeta.Labels[k8sOrganizationField])
		}
		for _, configMap := range k8sArtifact.ConfigMaps {
			assert.Equal(t, organizationHash, configMap.ObjectMeta.Labels[k8sOrganizationField])
		}
		for _, secret := range k8sArtifact.Secrets {
			assert.Equal(t, organizationHash, secret.ObjectMeta.Labels[k8sOrganizationField])
		}
		for _, scope := range k8sArtifact.Scopes {
			assert.Equal(t, organizationHash, scope.ObjectMeta.Labels[k8sOrganizationField])
		}
	}
}

func TestCreateConfigMaps(t *testing.T) {
	for _, k8Json := range sampleK8Artifacts {
		// Define input parameters
		certFiles := map[string]string{
			"cert1": "-----BEGIN CERTIFICATE----- MIIDWTCCAkGgAwIBAgIUbiBM1STcH3a8LjqLjelY1+jD2KwwDQYJKoZIhvcNAQEL BQAwgYgxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRMwEQYDVQQH DApTYW4gRnJhbmNpc2NvMRcwFQYDVQQKDA5PcGVuU3BhY2UgTGltaXRlZDEXMBUG A1UEAwwOcGxpdHRlc3QuY29tMB4XDTIxMDgyNTE0MDUyN1oXDTIyMDgyNTE0MDUy N1owgYgxCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRMwEQYDVQQH DApTYW4gRnJhbmNpc2NvMRcwFQYDVQQKDA5PcGVuU3BhY2UgTGltaXRlZDEXMBUG A1UEAwwOcGxpdHRlc3QuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC AQEAyH+VgZ5mcIiDw5YNvVD9sPsGQ5zUzOqK4JbQ2QipVt2mIXRMOHZjBsoLOIHw +U0mqyKTsDZN2zSq9N8Nc58VLyG2DLvOQqzqSC9P6hfrCed09pb3xRP2EnB16rli iC/DzN4Ou4gQ0JHh8THHIKd+OydQJpj1qoE/cpOpqkTx61Gd8RaN9YOm87dvyoYx kYzK9jsm24eX7l7pYzrQ/8oG++J4Cqof1f+bBjx8ZYxx92EhwGqRuBUVnROAv9WS vhJt7zk4H3ugVTJ9CBTmkdz+j5QZw4b36vJpySfu+DlDC6ZzuoXKZcc9k5l9MPnQ eG+MlH2sHwvtSfhiFpFbFQIDAQABo1MwUTAdBgNVHQ4EFgQUwCljqo6ES4rT3o+X ofc/m1j8i58wHwYDVR0jBBgwFoAUwCljqo6ES4rT3o+Xofc/m1j8i58wDwYDVR0T AQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAMnR2MzIRoB25/jl5g69GQ4ju eynDDr7GwRqfV8bJN05zmgnlGxBXkT/3jpiwuK+PBdzG6Dw1qRxbN52Z1QUzYpFq eN0B4K9Zmc4d82z9/4M+7tNLx09JKe7ky+f1QGkSZBxIjAKxPUyT8GCOVvQj0x9C 8q6ht3R4miq/rGpUXjJiWYTBZ2V/X33RlDfH38QrhqRYPltp++UDs+8LwTp3Dx4N 8cjplhh9lyM4lH33D20CNUw2T+3JOGtzgTn1ffwsxgDbW5Vf2RU8Qs5iTYoi8epF OnMzCqBt/t9gKGJ1oXc6T/URQKKKfGZL+RWbqFb1wUOuYfzL9nxI63tJxvK7yg== -----END CERTIFICATE-----",
			"cert2": "-----BEGIN CERTIFICATE----- MIIEqDCCA5CgAwIBAgIJAIql0LphfqC5MA0GCSqGSIb3DQEBCwUAMIGMMQswCQYD VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEUMBIGA1UEBwwLU2FuIEZyYW5j aXNjbzEXMBUGA1UECgwOT3BlblNwYWNlIExpbWl0ZWQxFzAVBgNVBAMMDnBsaXR0 ZXMudGVzdDAeFw0yMDEyMDExNjUzMTFaFw0yMTAzMzAxNjUzMTFaMIGMMQswCQYD VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEUMBIGA1UEBwwLU2FuIEZyYW5j aXNjbzEXMBUGA1UECgwOT3BlblNwYWNlIExpbWl0ZWQxFzAVBgNVBAMMDnBsaXR0 ZXMudGVzdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAIs8piOIsHH1 iTelPml0g9TLjnUWtAtIJpLPz5gI5cP16zBZMTnR3Qw4M2Kjwos9iRlfGWSttg0t nGo8DqP7d9HoP7chqey7hx7YmMvIgfrjT+wtbB69GBCW0vqn/3rPfb/IB4fVZ+RQ +D9r3k1Y3LR0ewIplbTfDit8xjqkFmGvvho8GpwP8P/yOIEteXJL3GceH1ap3Sre GwPGkNLwBe6AY8Hh1PcX4QXgUbA9tIYpqYwwVweRvdGRTkZrO5YiAbLOZJzVf6zM W2+3Xl86HkM0a/DJcGx1N7hZwWxyb+XX5OLbTtCLb6KAl2q28zxdBC6ftKysZVCH GOUCAwEAAaMnMCUwCwYDVR0PBAQDAgXgMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZI hvcNAQELBQADggEBAFUwdqeE6McZmK2Iq61xZlL2iZ3pE6/DqDP0BSjWvWSrvnfb RAR4+5WuS/q80MJd2iXXfC2AfIC0EGzzCkMP05gblhzLRp/J/VRW5uPrqzBp1h6F e/LI3bJjCpUMZ0WJi7HXVjRA7n/N1gGsZB5vTI5nmKTeCSMPA8V7R/q+QzhM5NL7 XwlVyxuyYlWuXUNp5GqHRQNNQem6v44tRM56NmN4nIylWvFjBwnunlvdqCr83HcE /bjXKSoBUERm9k7dRbUcm1JrQiJ6LCZbxdPhzXnccHkE4r8qq0WReJb6l5nHeEQE R1r2HqNwnHMtEvfUHJQpHYRrU06VQVvdQrTDhmQ= -----END CERTIFICATE-----",
		}

		var k8sArtifact K8sArtifacts
		err := json.Unmarshal([]byte(k8Json), &k8sArtifact)
		if err != nil {
			t.Error("Unable to unmarshal the dummy k8artifact")
		}

		// Call the function
		createConfigMaps(certFiles, &k8sArtifact)

		// Check whether the ConfigMaps are properly created and populated
		for confKey, confValue := range certFiles {
			pathSegments := strings.Split(confKey, ".")
			configName := k8sArtifact.API.Name + "-" + pathSegments[0]

			// Check whether the ConfigMap is created and stored in K8sArtifacts
			_, exists := k8sArtifact.ConfigMaps[configName]
			assert.True(t, exists, "ConfigMap should exist")

			// Check whether the data is properly set in the ConfigMap
			cm := k8sArtifact.ConfigMaps[configName]
			assert.NotNil(t, cm, "ConfigMap should not be nil")
			assert.Equal(t, "v1", cm.APIVersion, "APIVersion should be 'v1'")
			assert.Equal(t, "ConfigMap", cm.Kind, "Kind should be 'ConfigMap'")
			assert.Equal(t, confValue, cm.Data[confKey], "Data should match the provided certificate content")
		}
	}
}

func TestReplaceVhost(t *testing.T) {
	var k8sArtifact K8sArtifacts
	err := json.Unmarshal([]byte(HTTPk8Json), &k8sArtifact)
	if err != nil {
		t.Error("Unable to unmarshal the dummy k8artifact")
	}

	var envList = []Environment{}
	env1 := Environment{
		Name:              "Default",
		Vhost:             "env1.gw.wso2.com",
		DeployedTimeStamp: 1707977355909,
		Type:              "sandbox",
	}
	env2 := Environment{
		Name:              "Default",
		Vhost:             "env2.gw.wso2.com",
		DeployedTimeStamp: 1707977355909,
		Type:              "hybrid",
	}
	env3 := Environment{
		Name:              "Default",
		Vhost:             "env3.gw.wso2.com",
		DeployedTimeStamp: 1707977355909,
		Type:              "production",
	}
	envList = append(envList, env1, env2, env3)

	for _, env := range envList {
		replaceVhost(&k8sArtifact, env.Vhost, env.Type)
		if env.Type == "hybrid" {
			for _, routeName := range k8sArtifact.API.Spec.Production {
				for _, routes := range routeName.RouteRefs {
					assert.Equal(t, gwapiv1b1.Hostname(env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
				}
			}

			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					assert.Equal(t, gwapiv1b1.Hostname("sandbox."+env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
				}
			}

		} else if env.Type == "sandbox" {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					assert.Equal(t, gwapiv1b1.Hostname(env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
				}
			}
			assert.IsType(t, []dpv1alpha2.EnvConfig{}, k8sArtifact.API.Spec.Production)
			assert.Empty(t, k8sArtifact.API.Spec.Production, "Production should be empty")

		} else {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					assert.Equal(t, gwapiv1b1.Hostname(env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
				}
			}
			assert.IsType(t, []dpv1alpha2.EnvConfig{}, k8sArtifact.API.Spec.Sandbox)
			assert.Empty(t, k8sArtifact.API.Spec.Sandbox, "Sandbox should be empty")
		}
	}
}

func TestBrokenZipHandlingFlow(t *testing.T) {
	testResourcesDir := "../../resources/test-resources/Broken"
	k8ResourceGenEndpoint := "https://api.am.wso2.com:9095/api/configurator/1.0.0/apis/generate-k8s-resources"

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

			for _, zipFile := range zipReader.File {
				apiArtifact, err := DecodeAPIArtifact(zipFile)
				if err != nil {
					t.Logf("Error decoding API artifact from %s: %v", zipFile.Name, err)
					continue
				}
				// When these files are empty if they are present, then Decoding error should occur
				if strings.Contains(zipFile.Name, "Client_Cert_Empty") || strings.Contains(zipFile.Name, "Endpoint_Cert_Empty") {
					assert.Nil(t, apiArtifact)
					assert.Error(t, err)
				}

				apkConf, apiUUID, revisionID, apkErr := GenerateAPKConf(apiArtifact.APIJson, apiArtifact.CertArtifact)

				//When all the contents are empty or some properties are missing, an unmarshalling error should occur when creating the apiArtifact
				if strings.Contains(zipFile.Name, "All_Empty") {
					assert.Equal(t, "", apkConf)
					assert.Equal(t, "null", apiUUID)
					assert.Equal(t, uint32(0), revisionID)
					assert.Error(t, apkErr)
				}
				// If API_Json is broken then the generate conf is invalid hence it will be failed in CR generation
				if strings.Contains(zipFile.Name, "Empty_Definition") {
					certContainer := CertContainer{
						ClientCertObj:   apiArtifact.CertMeta,
						EndpointCertObj: apiArtifact.EndpointCertMeta,
					}
					crResponse, err := GenerateCRs(apkConf, apiArtifact.Schema, certContainer, k8ResourceGenEndpoint)
					assert.Error(t, err)
					assert.Nil(t, crResponse)
				}

			}
		}
	}
}
