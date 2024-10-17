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
 * Package "synchronizer" contains artifacts relate to fetching AI Provider
 * related updates from the control plane event-hub.
 * This file contains functions to retrieve AI Providers and AI Provider updates.
 */

package synchronizer

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	k8sclient "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/k8sClient"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	pkgAuth "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/auth"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/managementserver"
	sync "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/synchronizer"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/tlsutils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	aiProviderEndpoint string = "internal/data/v1/llm-providers"
)

// FetchAIProvidersOnEvent fetches the AI Providers from the control plane on the start up and notification event updates
func FetchAIProvidersOnEvent(aiProviderName string, aiProviderVersion string, organization string, c client.Client, cleanupDeletedProviders bool) {
	logger.LoggerSynchronizer.Info("Fetching AI Providers from Control Plane.")

	// Read configurations and derive the eventHub details
	conf, errReadConfig := config.ReadConfigs()
	if errReadConfig != nil {
		// This has to be error. For debugging purpose info
		logger.LoggerSynchronizer.Errorf("Error reading configs: %v", errReadConfig)
	}
	// Populate data from the config
	ehConfigs := conf.ControlPlane
	ehURL := ehConfigs.ServiceURL

	// If the eventHub URL is configured with trailing slash
	if strings.HasSuffix(ehURL, "/") {
		ehURL += aiProviderEndpoint
	} else {
		ehURL += "/" + aiProviderEndpoint
	}
	logger.LoggerSynchronizer.Debugf("Fetching AI Providers from the URL %v: ", ehURL)

	ehUname := ehConfigs.Username
	ehPass := ehConfigs.Password
	basicAuth := "Basic " + pkgAuth.GetBasicAuth(ehUname, ehPass)

	// Check if TLS is enabled
	skipSSL := ehConfigs.SkipSSLVerification

	// Create a HTTP request
	req, err := http.NewRequest("GET", ehURL, nil)
	if err != nil {
		logger.LoggerSynchronizer.Errorf("Error while creating http request for AI Providers Endpoint : %v", err)
	}

	queryParamMap := make(map[string]string)
	if aiProviderName != "" {
		queryParamMap["name"] = aiProviderName
	}
	if aiProviderVersion != "" {
		queryParamMap["apiVersion"] = aiProviderVersion
	}
	if organization != "" {
		queryParamMap["organization"] = organization
	}

	if len(queryParamMap) > 0 {
		q := req.URL.Query()
		// Making necessary query parameters for the request
		for queryParamKey, queryParamValue := range queryParamMap {
			q.Add(queryParamKey, queryParamValue)
		}
		req.URL.RawQuery = q.Encode()
	}
	// Setting authorization header
	req.Header.Set(sync.Authorization, basicAuth)

	if organization != "" {
		logger.LoggerSynchronizer.Debugf("Setting the organization header for the request: %v", organization)
		req.Header.Set("xWSO2Tenant", organization)
	} else {
		logger.LoggerSynchronizer.Debugf("Setting the organization header for the request: %v", "ALL")
		req.Header.Set("xWSO2Tenant", "ALL")
	}

	// Make the request
	logger.LoggerSynchronizer.Debugf("Sending the control plane request" + req.RequestURI)
	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)
	var errorMsg string
	if err != nil {
		errorMsg = "Error occurred while calling the REST API: " + aiProviderEndpoint
		go retryRLPFetchData(conf, errorMsg, err, c)
		return
	}
	responseBytes, err := io.ReadAll(resp.Body)
	logger.LoggerSynchronizer.Infof("Response String received for AI Providers: %v", string(responseBytes))

	if err != nil {
		errorMsg = "Error occurred while reading the response received for: " + aiProviderEndpoint
		go retryRLPFetchData(conf, errorMsg, err, c)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var aiProviderList eventhubTypes.AIProviderList
		err := json.Unmarshal(responseBytes, &aiProviderList)
		if err != nil {
			logger.LoggerSynchronizer.Errorf("Error occurred while unmarshelling AI Provider event data %v", err)
			return
		}
		logger.LoggerSynchronizer.Debugf("AI Providers received: %v", aiProviderList.APIs)
		var aiProviders []eventhubTypes.AIProvider = aiProviderList.APIs

		if cleanupDeletedProviders {
			aiProvidersFromK8, _, errK8 := k8sclient.RetrieveAllAIProvidersFromK8s(c, "")
			if errK8 == nil {
				for _, aiP := range aiProvidersFromK8 {
					if cpName, exists := aiP.ObjectMeta.Labels["CPName"]; exists {
						found := false
						for _, aiProviderFromCP := range aiProviders {
							if aiProviderFromCP.Name == cpName {
								found = true
								break
							}
						}
						if !found {
							// Delete the airatelimitpolicy
							k8sclient.DeleteAIProviderCR(aiP.Name, c)
						}
					}
				}
			} else {
				logger.LoggerSynchronizer.Errorf("Error while fetching aiproviders for cleaning up outdataed crs. Error: %+v", errK8)
			}
		}
		for _, aiProvider := range aiProviders {
			managementserver.AddAIProvider(aiProvider)
			logger.LoggerSynchronizer.Debugf("AI Provider added to internal map: %v", aiProvider)
			// Generate the AI Provider CR
			crAIProvider := createAIProvider(&aiProvider)
			// Deploy the AI Provider CR
			k8sclient.DeployAIProviderCR(&crAIProvider, c)
			logger.LoggerSynchronizer.Infof("AI Provider CR Deployed Successfully: %v", crAIProvider)
		}
	} else {
		errorMsg = "Failed to fetch data! " + aiProviderEndpoint + " responded with " +
			strconv.Itoa(resp.StatusCode)
		go retryRLPFetchData(conf, errorMsg, err, c)
	}

}

// createAIProvider creates the AI provider CR
func createAIProvider(aiProvider *eventhubTypes.AIProvider) dpv1alpha3.AIProvider {
	conf, _ := config.ReadConfigs()
	sha1ValueofAIProviderName := GetSha1Value(aiProvider.Name)
	sha1ValueOfOrganization := GetSha1Value(aiProvider.Organization)
	labelMap := map[string]string{"name": sha1ValueofAIProviderName,
		"organization": sha1ValueOfOrganization,
		"InitiateFrom": "CP",
		"CPName":       aiProvider.Name,
	}
	var modelInputSource string
	var modelAttributeIdentifier string
	var promptTokenCountInputSource string
	var promptTokenCountAttributeIdentifier string
	var completionTokenCountInputSource string
	var completionTokenCountAttributeIdentifier string
	var totalTokenCountInputSource string
	var totalTokenCountAttributeIdentifier string

	var config eventhubTypes.Config
	err := json.Unmarshal([]byte(aiProvider.Configurations), &config)
	if err != nil {
		logger.LoggerSynchronizer.Errorf("Error unmarshalling configurations metadata in AI Provider: %v", err)
	}

	for _, field := range config.Metadata {
		if field.AttributeName == "model" {
			modelInputSource = field.InputSource
			modelAttributeIdentifier = field.AttributeIdentifier
		} else if field.AttributeName == "promptTokenCount" {
			promptTokenCountInputSource = field.InputSource
			promptTokenCountAttributeIdentifier = field.AttributeIdentifier
		} else if field.AttributeName == "completionTokenCount" {
			completionTokenCountInputSource = field.InputSource
			completionTokenCountAttributeIdentifier = field.AttributeIdentifier
		} else if field.AttributeName == "totalTokenCount" {
			totalTokenCountInputSource = field.InputSource
			totalTokenCountAttributeIdentifier = field.AttributeIdentifier
		}
	}

	crAIProvider := dpv1alpha3.AIProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aiProvider.ID,
			Namespace: conf.DataPlane.Namespace,
			Labels:    labelMap,
		},
		Spec: dpv1alpha3.AIProviderSpec{
			ProviderName:       aiProvider.Name,
			ProviderAPIVersion: aiProvider.APIVersion,
			Organization:       aiProvider.Organization,
			Model: dpv1alpha3.ValueDetails{
				In:    modelInputSource,
				Value: modelAttributeIdentifier,
			},
			RateLimitFields: dpv1alpha3.RateLimitFields{
				PromptTokens: dpv1alpha3.ValueDetails{
					In:    promptTokenCountInputSource,
					Value: promptTokenCountAttributeIdentifier,
				},
				CompletionToken: dpv1alpha3.ValueDetails{
					In:    completionTokenCountInputSource,
					Value: completionTokenCountAttributeIdentifier,
				},
				TotalToken: dpv1alpha3.ValueDetails{
					In:    totalTokenCountInputSource,
					Value: totalTokenCountAttributeIdentifier,
				},
			},
		},
	}
	return crAIProvider
}

// GetSha1Value returns the SHA1 value of the input string
func GetSha1Value(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
