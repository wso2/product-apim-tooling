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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	k8sclient "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/k8sClient"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	pkgAuth "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/auth"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/managementserver"
	sync "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/synchronizer"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/tlsutils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	policiesEndpoint                    string = "internal/data/v1/api-policies"
	policiesByNameEndpoint              string = "internal/data/v1/api-policies?policyName="
	subscriptionsPoliciesEndpoint       string = "internal/data/v1//subscription-policies"
	subscriptionsPoliciesByNameEndpoint string = "internal/data/v1//subscription-policies?policyName="
)

// FetchRateLimitPoliciesOnEvent fetches the policies from the control plane on the start up and notification event updates
func FetchRateLimitPoliciesOnEvent(ratelimitName string, organization string, c client.Client) {
	logger.LoggerSynchronizer.Info("Fetching RateLimit Policies from Control Plane.")

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
		if ratelimitName != "" {
			ehURL += policiesByNameEndpoint + ratelimitName
		} else {
			ehURL += policiesEndpoint
		}
	} else {
		if ratelimitName != "" {
			ehURL += "/" + policiesByNameEndpoint + ratelimitName
		} else {
			ehURL += "/" + policiesEndpoint
		}
	}

	logger.LoggerSynchronizer.Debugf("Fetching RateLimit Policies from the URL %v: ", ehURL)

	ehUname := ehConfigs.Username
	ehPass := ehConfigs.Password
	basicAuth := "Basic " + pkgAuth.GetBasicAuth(ehUname, ehPass)

	// Check if TLS is enabled
	skipSSL := ehConfigs.SkipSSLVerification

	// Create a HTTP request
	req, err := http.NewRequest("GET", ehURL, nil)
	if err != nil {
		logger.LoggerSynchronizer.Errorf("Error while creating http request for RateLimit Policies Endpoint : %v", err)
	}

	var queryParamMap map[string]string

	if queryParamMap != nil && len(queryParamMap) > 0 {
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
	logger.LoggerSynchronizer.Debug("Sending the control plane request")
	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)
	var errorMsg string
	if err != nil {
		errorMsg = "Error occurred while calling the REST API: " + policiesEndpoint
		go retryRLPFetchData(conf, errorMsg, err, c)
		return
	}
	responseBytes, err := ioutil.ReadAll(resp.Body)
	logger.LoggerSynchronizer.Debugf("Response String received for Policies: %v", string(responseBytes))

	if err != nil {
		errorMsg = "Error occurred while reading the response received for: " + policiesEndpoint
		go retryRLPFetchData(conf, errorMsg, err, c)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var rateLimitPolicyList eventhubTypes.RateLimitPolicyList
		err := json.Unmarshal(responseBytes, &rateLimitPolicyList)
		if err != nil {
			logger.LoggerSynchronizer.Errorf("Error occurred while unmarshelling RateLimit Policies event data %v", err)
			return
		}
		logger.LoggerSynchronizer.Debugf("Policies received: %v", rateLimitPolicyList.List)
		var rateLimitPolicies []eventhubTypes.RateLimitPolicy = rateLimitPolicyList.List
		for _, policy := range rateLimitPolicies {
			if policy.DefaultLimit.RequestCount.TimeUnit == "min" {
				policy.DefaultLimit.RequestCount.TimeUnit = "Minute"
			} else if policy.DefaultLimit.RequestCount.TimeUnit == "hour" {
				policy.DefaultLimit.RequestCount.TimeUnit = "Hour"
			} else if policy.DefaultLimit.RequestCount.TimeUnit == "day" {
				policy.DefaultLimit.RequestCount.TimeUnit = "Day"
			}
			managementserver.AddRateLimitPolicy(policy)
			logger.LoggerSynchronizer.Infof("RateLimit Policy added to internal map: %v", policy)
			// Update the exisitng rate limit policies with current policy
			k8sclient.UpdateRateLimitPolicyCR(policy, c)

		}
	} else {
		errorMsg = "Failed to fetch data! " + policiesEndpoint + " responded with " +
			strconv.Itoa(resp.StatusCode)
		go retryRLPFetchData(conf, errorMsg, err, c)
	}
}

// FetchSubscriptionRateLimitPoliciesOnEvent fetches the policies from the control plane on the start up and notification event updates
func FetchSubscriptionRateLimitPoliciesOnEvent(ratelimitName string, organization string, c client.Client) {
	logger.LoggerSynchronizer.Info("Fetching Subscription RateLimit Policies from Control Plane.")

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
		if ratelimitName != "" {
			ehURL += subscriptionsPoliciesByNameEndpoint + ratelimitName
		} else {
			ehURL += subscriptionsPoliciesEndpoint
		}
	} else {
		if ratelimitName != "" {
			ehURL += "/" + subscriptionsPoliciesByNameEndpoint + ratelimitName
		} else {
			ehURL += "/" + subscriptionsPoliciesEndpoint
		}
	}

	logger.LoggerSynchronizer.Infof("Fetching Subscription RateLimit Policies from the URL %v: ", ehURL)

	ehUname := ehConfigs.Username
	ehPass := ehConfigs.Password
	basicAuth := "Basic " + pkgAuth.GetBasicAuth(ehUname, ehPass)

	// Check if TLS is enabled
	skipSSL := ehConfigs.SkipSSLVerification

	// Create a HTTP request
	req, err := http.NewRequest("GET", ehURL, nil)
	if err != nil {
		logger.LoggerSynchronizer.Errorf("Error while creating http request for Subscription RateLimit Policies Endpoint : %v", err)
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
	logger.LoggerSynchronizer.Debug("Sending the control plane request")
	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)
	var errorMsg string
	if err != nil {
		errorMsg = "Error occurred while calling the REST API: " + policiesEndpoint
		go retrySubscriptionRLPFetchData(conf, errorMsg, err, c)
		return
	}
	responseBytes, err := ioutil.ReadAll(resp.Body)
	var jsonResponse interface{}
	err = json.Unmarshal(responseBytes, &jsonResponse)
	if err != nil {
		logger.LoggerSynchronizer.Infof("Error unmarshalling JSON:", err)
		return
	}

	// Print the JSON response in a pretty format
	jsonOutput, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		logger.LoggerSynchronizer.Infof("Error formatting JSON:", err)
		return
	}
	logger.LoggerSynchronizer.Infof("Json string: %s",string(jsonOutput))
	logger.LoggerSynchronizer.Debugf("Response String received for Policies: %v", string(responseBytes))

	if err != nil {
		errorMsg = "Error occurred while reading the response received for: " + policiesEndpoint
		go retrySubscriptionRLPFetchData(conf, errorMsg, err, c)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var rateLimitPolicyList eventhubTypes.SubscriptionPolicyList
		err := json.Unmarshal(responseBytes, &rateLimitPolicyList)
		if err != nil {
			logger.LoggerSynchronizer.Errorf("Error occurred while unmarshelling Subscription RateLimit Policies event data %v", err)
			return
		}
		logger.LoggerSynchronizer.Debugf("Policies received: %v", rateLimitPolicyList.List)
		var rateLimitPolicies []eventhubTypes.SubscriptionPolicy = rateLimitPolicyList.List
		for _, policy := range rateLimitPolicies {
			if policy.QuotaType == "aiApiQuota" {
				if policy.DefaultLimit.AiApiQuota != nil {
					switch policy.DefaultLimit.AiApiQuota.TimeUnit {
					case "min":
						policy.DefaultLimit.AiApiQuota.TimeUnit = "Minute"
					case "hours":
						policy.DefaultLimit.AiApiQuota.TimeUnit = "Hour"
					case "days":
						policy.DefaultLimit.AiApiQuota.TimeUnit = "Day"
					default:
						logger.LoggerSynchronizer.Errorf("Unsupported timeunit %s", policy.DefaultLimit.AiApiQuota.TimeUnit)
						continue
					}
					if policy.DefaultLimit.AiApiQuota.PromptTokenCount == nil && policy.DefaultLimit.AiApiQuota.TotalTokenCount != nil {
						policy.DefaultLimit.AiApiQuota.PromptTokenCount = policy.DefaultLimit.AiApiQuota.TotalTokenCount
					}
					if policy.DefaultLimit.AiApiQuota.CompletionTokenCount == nil && policy.DefaultLimit.AiApiQuota.TotalTokenCount != nil {
						policy.DefaultLimit.AiApiQuota.CompletionTokenCount = policy.DefaultLimit.AiApiQuota.TotalTokenCount
					}
					if policy.DefaultLimit.AiApiQuota.TotalTokenCount == nil && policy.DefaultLimit.AiApiQuota.PromptTokenCount != nil && policy.DefaultLimit.AiApiQuota.CompletionTokenCount != nil {
						total := *policy.DefaultLimit.AiApiQuota.PromptTokenCount + *policy.DefaultLimit.AiApiQuota.CompletionTokenCount
						policy.DefaultLimit.AiApiQuota.TotalTokenCount = &total
					}
					managementserver.AddSubscriptionPolicy(policy)
					k8sclient.DeployAIRateLimitPolicyCR(policy, c)
				} else {
					logger.LoggerSynchronizer.Errorf("AIQuota type response recieved but no data found. %+v", policy.DefaultLimit)
				}
			} else {
				if policy.DefaultLimit.RequestCount.TimeUnit == "min" {
					policy.DefaultLimit.RequestCount.TimeUnit = "Minute"
				} else if policy.DefaultLimit.RequestCount.TimeUnit == "hours" {
					policy.DefaultLimit.RequestCount.TimeUnit = "Hour"
				} else if policy.DefaultLimit.RequestCount.TimeUnit == "days" {
					policy.DefaultLimit.RequestCount.TimeUnit = "Day"
				}
				managementserver.AddSubscriptionPolicy(policy)
				logger.LoggerSynchronizer.Infof("RateLimit Policy added to internal map: %v", policy)
				// Update the exisitng rate limit policies with current policy
				k8sclient.DeploySubscriptionRateLimitPolicyCR(policy, c)
			}
		}
	} else {
		errorMsg = "Failed to fetch data! " + policiesEndpoint + " responded with " +
			strconv.Itoa(resp.StatusCode)
		go retrySubscriptionRLPFetchData(conf, errorMsg, err, c)
	}
}

func retryRLPFetchData(conf *config.Config, errorMessage string, err error, c client.Client) {
	logger.LoggerSynchronizer.Debugf("Time Duration for retrying: %v",
		conf.ControlPlane.RetryInterval*time.Second)
	time.Sleep(conf.ControlPlane.RetryInterval * time.Second)
	FetchRateLimitPoliciesOnEvent("", "", c)
	retryAttempt++
	if retryAttempt >= retryCount {
		logger.LoggerSynchronizer.Errorf(errorMessage, err)
		return
	}
}

func retrySubscriptionRLPFetchData(conf *config.Config, errorMessage string, err error, c client.Client) {
	logger.LoggerSynchronizer.Debugf("Time Duration for retrying: %v",
		conf.ControlPlane.RetryInterval*time.Second)
	time.Sleep(conf.ControlPlane.RetryInterval * time.Second)
	FetchSubscriptionRateLimitPoliciesOnEvent("", "", c)
	retryAttempt++
	if retryAttempt >= retryCount {
		logger.LoggerSynchronizer.Errorf(errorMessage, err)
		return
	}
}
