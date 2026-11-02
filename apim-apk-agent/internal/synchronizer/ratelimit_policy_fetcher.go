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
	subscriptionsPoliciesEndpoint       string = "internal/data/v1/subscription-policies"
	subscriptionsPoliciesByNameEndpoint string = "internal/data/v1/subscription-policies?policyName="
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
			// API level policies do not have burst control, so only the request count time unit is resolved.
			// API level rate limits in the data plane do not support a time unit of seconds.
			requestCountUnit, ok := resolveCRTimeUnitWithoutSeconds(policy.DefaultLimit.RequestCount.TimeUnit)
			if !ok {
				if policy.DefaultLimit.RequestCount.TimeUnit == "" {
					logger.LoggerSynchronizer.Debugf("Empty timeunit for policy %s, skipping", policy.Name)
				} else {
					logger.LoggerSynchronizer.Errorf("Unsupported timeunit '%s' for policy %s, skipping",
						policy.DefaultLimit.RequestCount.TimeUnit, policy.Name)
				}
				continue
			}
			policy.DefaultLimit.RequestCount.TimeUnit = requestCountUnit
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
func FetchSubscriptionRateLimitPoliciesOnEvent(ratelimitName string, organization string, c client.Client, cleanupDeletedPolicies bool) {
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
		if cleanupDeletedPolicies {
			// This logic is executed once at the startup time so no need to worry about the nested for loops for performance.
			// Fetch all AiRatelimitPolicies
			airls, _, retrieveAllAIRLErr := k8sclient.RetrieveAllAIRatelimitPoliciesSFromK8s(c, "")
			rls, _, retrieveAllRLErr := k8sclient.RetrieveAllRatelimitPoliciesSFromK8s(c, "")
			if retrieveAllAIRLErr == nil {
				for _, airl := range airls {
					if cpName, exists := airl.ObjectMeta.Labels["CPName"]; exists {
						found := false
						for _, policy := range rateLimitPolicies {
							if policy.Name == cpName {
								found = true
								break
							}
						}
						if !found {
							// Delete the airatelimitpolicy
							k8sclient.UndeploySubscriptionAIRateLimitPolicyCR(airl.Name, c)
						}
					}
				}
			} else {
				logger.LoggerSynchronizer.Errorf("Error while fetching subscription airatelimitpolicies for cleaning up outdated crs. Error: %+v", retrieveAllAIRLErr)
			}
			if retrieveAllRLErr == nil {
				for _, rl := range rls {
					if cpName, exists := rl.ObjectMeta.Labels["CPName"]; exists {
						found := false
						for _, policy := range rateLimitPolicies {
							if policy.Name == cpName {
								found = true
								break
							}
						}
						if !found {
							// Delete the airatelimitpolicy
							k8sclient.UnDeploySubscriptionRateLimitPolicyCR(rl.Name, c)
						}
					}
				}
			} else {
				logger.LoggerSynchronizer.Errorf("Error while fetching subscription ratelimitpolicies for cleaning up outdated crs. Error: %+v", retrieveAllRLErr)
			}
		}

		for _, rateLimitPolicy := range rateLimitPolicies {
			if rateLimitPolicy.QuotaType == "aiApiQuota" {
				if rateLimitPolicy.DefaultLimit.AiAPIQuota != nil {
					// AI rate limits in the data plane do not support a time unit of seconds.
					aiQuotaUnit, ok := resolveCRTimeUnitWithoutSeconds(rateLimitPolicy.DefaultLimit.AiAPIQuota.TimeUnit)
					if !ok {
						if rateLimitPolicy.DefaultLimit.AiAPIQuota.TimeUnit == "" {
							logger.LoggerSynchronizer.Debugf("Empty AiAPIQuota timeunit for policy %s, skipping", rateLimitPolicy.Name)
						} else {
							logger.LoggerSynchronizer.Errorf("Unsupported AiAPIQuota timeunit '%s' for policy %s, skipping",
								rateLimitPolicy.DefaultLimit.AiAPIQuota.TimeUnit, rateLimitPolicy.Name)
						}
						continue
					}
					rateLimitPolicy.DefaultLimit.AiAPIQuota.TimeUnit = aiQuotaUnit
					if rateLimitPolicy.DefaultLimit.AiAPIQuota.PromptTokenCount == nil && rateLimitPolicy.DefaultLimit.AiAPIQuota.TotalTokenCount != nil {
						rateLimitPolicy.DefaultLimit.AiAPIQuota.PromptTokenCount = rateLimitPolicy.DefaultLimit.AiAPIQuota.TotalTokenCount
					}
					if rateLimitPolicy.DefaultLimit.AiAPIQuota.CompletionTokenCount == nil && rateLimitPolicy.DefaultLimit.AiAPIQuota.TotalTokenCount != nil {
						rateLimitPolicy.DefaultLimit.AiAPIQuota.CompletionTokenCount = rateLimitPolicy.DefaultLimit.AiAPIQuota.TotalTokenCount
					}
					if rateLimitPolicy.DefaultLimit.AiAPIQuota.TotalTokenCount == nil && rateLimitPolicy.DefaultLimit.AiAPIQuota.PromptTokenCount != nil && rateLimitPolicy.DefaultLimit.AiAPIQuota.CompletionTokenCount != nil {
						total := *rateLimitPolicy.DefaultLimit.AiAPIQuota.PromptTokenCount + *rateLimitPolicy.DefaultLimit.AiAPIQuota.CompletionTokenCount
						rateLimitPolicy.DefaultLimit.AiAPIQuota.TotalTokenCount = &total
					}
					managementserver.AddSubscriptionPolicy(rateLimitPolicy)
					k8sclient.DeployAIRateLimitPolicyFromCPPolicy(rateLimitPolicy, c)
				} else {
					logger.LoggerSynchronizer.Errorf("AIQuota type response recieved but no data found. %+v", rateLimitPolicy.DefaultLimit)
				}
			} else {
				if !normalizeSubscriptionRateLimitPolicy(&rateLimitPolicy) {
					continue
				}
				managementserver.AddSubscriptionPolicy(rateLimitPolicy)
				logger.LoggerSynchronizer.Infof("Subscription RateLimit Policy added to internal map: %v", rateLimitPolicy)
				// Update the exisitng rate limit policies with current policy
				k8sclient.DeploySubscriptionRateLimitPolicyCR(rateLimitPolicy, c)
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
	FetchSubscriptionRateLimitPoliciesOnEvent("", "", c, false)
	retryAttempt++
	if retryAttempt >= retryCount {
		logger.LoggerSynchronizer.Errorf(errorMessage, err)
		return
	}
}
