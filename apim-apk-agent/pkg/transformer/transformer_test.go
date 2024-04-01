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
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
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

				apkConf, apiUUID, revisionID, configuredRateLimitPoliciesMap, endpointSecurityData, apkErr := GenerateAPKConf(apiArtifact.APIJson, apiArtifact.CertArtifact, "default")

				assert.NoError(t, apkErr)
				assert.NotEmpty(t, apkConf)
				assert.NotEqual(t, "null", apiUUID)
				assert.NotEqual(t, uint32(0), revisionID)
				assert.NotNil(t, configuredRateLimitPoliciesMap)
				assert.IsType(t, EndpointSecurityConfig{}, endpointSecurityData) // Need to be refined maybe
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
			"cert1": "-----BEGIN CERTIFICATE----- LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHVENDQWdFQ0ZBTklrTFFCa2Q3NnFpVFh6U1hqQlMyc2NQSnNNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1FMHgKQ3pBSkJnTlZCQVlUQWt4TE1STXdFUVlEVlFRSURBcFRiMjFsTFZOMFlYUmxNUTB3Q3dZRFZRUUtEQVIzYzI4eQpNUXd3Q2dZRFZRUUxEQU5oY0dzeEREQUtCZ05WQkFNTUEyRndhekFlRncweU16RXlNRFl4TURFeU5EaGFGdzB5Ck5UQTBNVGt4TURFeU5EaGFNRVV4Q3pBSkJnTlZCQVlUQWt4TE1STXdFUVlEVlFRSURBcFRiMjFsTFZOMFlYUmwKTVNFd0h3WURWUVFLREJoSmJuUmxjbTVsZENCWGFXUm5hWFJ6SUZCMGVTQk1kR1F3Z2dFaU1BMEdDU3FHU0liMwpEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUNkRzkwVy9UbGs0dTlhd0hQdGVENXpwVmNUaFVLd01MdkFLdzlpCnZWUUJDMEFHNkd6UGJha29sNWdLVm0ra0JVREZ6enpGNmVheUVYS1dieWFaRHR5NjZBMis3SExMY0tCb3A1TS8KYTU3UTlYdFUzbFJZdm90Z3V0TFd1SGNJN21MQ1NjWkRyakEzcm5iL0tqamJoWjYwMlpTMXBwNWp0eVV6NkR3TAptN3c0d1EvUlByb3FDZEJqOFFxb0F2bkRETFNQZURmc3gxNEo1VmVOSlZHSlYyd2F4NjVqV1JqUmtqNndFN3oyCnF6V0FsUDV2RGVFRDZib2dZWVZEcEM4RHRnYXlRK3ZLQVFMaTF1aitJOVlxYi9uUFVyZFVoOUlseHVkbHFpRlEKUXh5dnNYTUpFemJXV21sYkQwa1hZa0htSHpldEpOUEs5YXlPUy9mSmNBY2ZBYjAxQWdNQkFBRXdEUVlKS29aSQpodmNOQVFFTEJRQURnZ0VCQUZtVWM3K2NJOGQwRGw0d1RkcStnZnlXZHFqUWI3QVlWTzlEdkppM1hHeGRjNUtwCjFuQ1NzS3pLVXo5Z3Z4WEhlYVlLckJOWWY0U1NVK1BrZGYvQldlUHFpN1VYL1NJeE5YYnkyZGE4eldnK1c2VWgKeFpmS2xMWUdNcDNtQ2p1ZVpwWlRKN1NLT09HRkE4SUlnRXpqSkQ5TG4xZ2wzeXdNYUN3bE5yRzlScGlEMU1jVApDT0t2eVdOS25TUlZyL1J2Q2tsTFZyQU1USnI1MGtjZTJjemNkRmwveEY0SG02NnZwN2NQL2JZSktXQUw4aEJHCnpVYTlhUUJLbmNPb0FPK3pRL1NHeTd1SnhURFVGOFN2ZXJEc21qT2M2QVU2SWhCR1ZVeVgvSlFiWXlKZlppbkIKWWx2aVl4VnpJbTZJYU5KSHg0c2lodzRVMS9qTUZXUlhUNDcwemNRPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t -----END CERTIFICATE-----",
			"cert2": "-----BEGIN CERTIFICATE----- LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHVENDQWdFQ0ZBTklrTFFCa2Q3NnFpVFh6U1hqQlMyc2NQSnNNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1FMHgKQ3pBSkJnTlZCQVlUQWt4TE1STXdFUVlEVlFRSURBcFRiMjFsTFZOMFlYUmxNUTB3Q3dZRFZRUUtEQVIzYzI4eQpNUXd3Q2dZRFZRUUxEQU5oY0dzeEREQUtCZ05WQkFNTUEyRndhekFlRncweU16RXlNRFl4TURFeU5EaGFGdzB5Ck5UQTBNVGt4TURFeU5EaGFNRVV4Q3pBSkJnTlZCQVlUQWt4TE1STXdFUVlEVlFRSURBcFRiMjFsTFZOMFlYUmwKTVNFd0h3WURWUVFLREJoSmJuUmxjbTVsZENCWGFXUm5hWFJ6SUZCMGVTQk1kR1F3Z2dFaU1BMEdDU3FHU0liMwpEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUNkRzkwVy9UbGs0dTlhd0hQdGVENXpwVmNUaFVLd01MdkFLdzlpCnZWUUJDMEFHNkd6UGJha29sNWdLVm0ra0JVREZ6enpGNmVheUVYS1dieWFaRHR5NjZBMis3SExMY0tCb3A1TS8KYTU3UTlYdFUzbFJZdm90Z3V0TFd1SGNJN21MQ1NjWkRyakEzcm5iL0tqamJoWjYwMlpTMXBwNWp0eVV6NkR3TAptN3c0d1EvUlByb3FDZEJqOFFxb0F2bkRETFNQZURmc3gxNEo1VmVOSlZHSlYyd2F4NjVqV1JqUmtqNndFN3oyCnF6V0FsUDV2RGVFRDZib2dZWVZEcEM4RHRnYXlRK3ZLQVFMaTF1aitJOVlxYi9uUFVyZFVoOUlseHVkbHFpRlEKUXh5dnNYTUpFemJXV21sYkQwa1hZa0htSHpldEpOUEs5YXlPUy9mSmNBY2ZBYjAxQWdNQkFBRXdEUVlKS29aSQpodmNOQVFFTEJRQURnZ0VCQUZtVWM3K2NJOGQwRGw0d1RkcStnZnlXZHFqUWI3QVlWTzlEdkppM1hHeGRjNUtwCjFuQ1NzS3pLVXo5Z3Z4WEhlYVlLckJOWWY0U1NVK1BrZGYvQldlUHFpN1VYL1NJeE5YYnkyZGE4eldnK1c2VWgKeFpmS2xMWUdNcDNtQ2p1ZVpwWlRKN1NLT09HRkE4SUlnRXpqSkQ5TG4xZ2wzeXdNYUN3bE5yRzlScGlEMU1jVApDT0t2eVdOS25TUlZyL1J2Q2tsTFZyQU1USnI1MGtjZTJjemNkRmwveEY0SG02NnZwN2NQL2JZSktXQUw4aEJHCnpVYTlhUUJLbmNPb0FPK3pRL1NHeTd1SnhURFVGOFN2ZXJEc21qT2M2QVU2SWhCR1ZVeVgvSlFiWXlKZlppbkIKWWx2aVl4VnpJbTZJYU5KSHg0c2lodzRVMS9qTUZXUlhUNDcwemNRPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t -----END CERTIFICATE-----",
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
			assert.NotNil(t, confValue, "Data should not be nil")
			decodedCert := "-----BEGIN CERTIFICATE-----\nMIIDGTCCAgECFANIkLQBkd76qiTXzSXjBS2scPJsMA0GCSqGSIb3DQEBCwUAME0x\nCzAJBgNVBAYTAkxLMRMwEQYDVQQIDApTb21lLVN0YXRlMQ0wCwYDVQQKDAR3c28y\nMQwwCgYDVQQLDANhcGsxDDAKBgNVBAMMA2FwazAeFw0yMzEyMDYxMDEyNDhaFw0y\nNTA0MTkxMDEyNDhaMEUxCzAJBgNVBAYTAkxLMRMwEQYDVQQIDApTb21lLVN0YXRl\nMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3\nDQEBAQUAA4IBDwAwggEKAoIBAQCdG90W/Tlk4u9awHPteD5zpVcThUKwMLvAKw9i\nvVQBC0AG6GzPbakol5gKVm+kBUDFzzzF6eayEXKWbyaZDty66A2+7HLLcKBop5M/\na57Q9XtU3lRYvotgutLWuHcI7mLCScZDrjA3rnb/KjjbhZ602ZS1pp5jtyUz6DwL\nm7w4wQ/RProqCdBj8QqoAvnDDLSPeDfsx14J5VeNJVGJV2wax65jWRjRkj6wE7z2\nqzWAlP5vDeED6bogYYVDpC8DtgayQ+vKAQLi1uj+I9Yqb/nPUrdUh9IlxudlqiFQ\nQxyvsXMJEzbWWmlbD0kXYkHmHzetJNPK9ayOS/fJcAcfAb01AgMBAAEwDQYJKoZI\nhvcNAQELBQADggEBAFmUc7+cI8d0Dl4wTdq+gfyWdqjQb7AYVO9DvJi3XGxdc5Kp\n1nCSsKzKUz9gvxXHeaYKrBNYf4SSU+Pkdf/BWePqi7UX/SIxNXby2da8zWg+W6Uh\nxZfKlLYGMp3mCjueZpZTJ7SKOOGFA8IIgEzjJD9Ln1gl3ywMaCwlNrG9RpiD1McT\nCOKvyWNKnSRVr/RvCklLVrAMTJr50kce2czcdFl/xF4Hm66vp7cP/bYJKWAL8hBG\nzUa9aQBKncOoAO+zQ/SGy7uJxTDUF8SverDsmjOc6AU6IhBGVUyX/JQbYyJfZinB\nYlviYxVzIm6IaNJHx4sihw4U1/jMFWRXT470zcQ=\n-----END CERTIFICATE-----"
			assert.Equal(t, decodedCert, cm.Data[confKey], "Data should match the provided certificate content")
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
					assert.Equal(t, gwapiv1.Hostname(env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
				}
			}

			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					assert.Equal(t, gwapiv1.Hostname("sandbox."+env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
				}
			}

		} else if env.Type == "sandbox" {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					assert.Equal(t, gwapiv1.Hostname(env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
				}
			}
			assert.IsType(t, []dpv1alpha2.EnvConfig{}, k8sArtifact.API.Spec.Production)
			assert.Empty(t, k8sArtifact.API.Spec.Production, "Production should be empty")

		} else {
			for _, routeName := range k8sArtifact.API.Spec.Sandbox {
				for _, routes := range routeName.RouteRefs {
					assert.Equal(t, gwapiv1.Hostname(env.Vhost), k8sArtifact.HTTPRoutes[routes].Spec.Hostnames[0])
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

				apkConf, apiUUID, revisionID, configuredRateLimitPoliciesMap, endpointSecurityData, apkErr := GenerateAPKConf(apiArtifact.APIJson, apiArtifact.CertArtifact, "orgID")

				//When all the contents are empty or some properties are missing, an unmarshalling error should occur when creating the apiArtifact
				if strings.Contains(zipFile.Name, "All_Empty") {
					assert.Equal(t, "", apkConf)
					assert.Equal(t, "null", apiUUID)
					assert.Equal(t, uint32(0), revisionID)
					assert.Error(t, apkErr)
					assert.NotNil(t, configuredRateLimitPoliciesMap)
					assert.IsType(t, EndpointSecurityConfig{}, endpointSecurityData) // Need to be refined maybe
				}
				// If API_Json is broken then the generate conf is invalid hence it will be failed in CR generation
				if strings.Contains(zipFile.Name, "Empty_Definition") {
					certContainer := CertContainer{
						ClientCertObj:   apiArtifact.CertMeta,
						EndpointCertObj: apiArtifact.EndpointCertMeta,
					}
					crResponse, err := GenerateCRs(apkConf, apiArtifact.Schema, certContainer, k8ResourceGenEndpoint, "orgID")
					assert.Error(t, err)
					assert.Nil(t, crResponse)
				}

			}
		}
	}
}
