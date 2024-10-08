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
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/eventhub"
	k8sclient "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/k8sClient"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/logging"
	pkgAuth "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/auth"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	sync "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/synchronizer"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/tlsutils"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	keyManagersEndpoint string = "internal/data/v1/keymanagers"
	retryCount          int    = 5
)

var retryAttempt int

// FetchKeyManagersOnStartUp pulls the Key managers calling to the API manager
// API Manager returns a .zip file as a response and this function
// returns a byte slice of that ZIP file.
func FetchKeyManagersOnStartUp(c client.Client) {
	logger.LoggerSynchronizer.Info("Fetching KeyManagers from Control Plane.")

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
		ehURL += keyManagersEndpoint
	} else {
		ehURL += "/" + keyManagersEndpoint
	}
	logger.LoggerSynchronizer.Debugf("Fetching KeyManagers from the URL %v: ", ehURL)

	ehUname := ehConfigs.Username
	ehPass := ehConfigs.Password
	basicAuth := "Basic " + pkgAuth.GetBasicAuth(ehUname, ehPass)

	// Check if TLS is enabled
	skipSSL := ehConfigs.SkipSSLVerification

	// Create a HTTP request
	req, err := http.NewRequest("GET", ehURL, nil)
	if err != nil {
		logger.LoggerSynchronizer.Errorf("Error while creating http request for Key Manager Endpoint : %v", err)
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

	//Todo: Need to set ALL when APIM Fix is available
	req.Header.Set("xWSO2Tenant", "ALL")

	// Make the request
	logger.LoggerSynchronizer.Debug("Sending the control plane request")
	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)
	var errorMsg string
	if err != nil {
		errorMsg = "Error occurred while calling the REST API: " + keyManagersEndpoint
		go retryFetchData(conf, errorMsg, err, c)
		return
	}
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMsg = "Error occurred while reading the response received for: " + keyManagersEndpoint
		go retryFetchData(conf, errorMsg, err, c)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var keyManagers []eventhubTypes.KeyManager
		err := json.Unmarshal(responseBytes, &keyManagers)
		if err != nil {
			logger.LoggerSynchronizer.Errorf("Error occurred while unmarshelling Key Managers event data %v", err)
			return
		}
		logger.LoggerSynchronizer.Infof("Key Managers received: %v", keyManagers)
		resolvedKeyManagers := eventhub.MarshalKeyManagers(&keyManagers)
		applyAllKeymanagerConfifuration(c, resolvedKeyManagers)
	} else {
		errorMsg = "Failed to fetch data! " + keyManagersEndpoint + " responded with " +
			strconv.Itoa(resp.StatusCode)
		go retryFetchData(conf, errorMsg, err, c)
	}
}

func retryFetchData(conf *config.Config, errorMessage string, err error, c client.Client) {
	logger.LoggerSynchronizer.Debugf("Time Duration for retrying: %v",
		conf.ControlPlane.RetryInterval*time.Second)
	time.Sleep(conf.ControlPlane.RetryInterval * time.Second)
	FetchKeyManagersOnStartUp(c)
	retryAttempt++
	if retryAttempt >= retryCount {
		logger.LoggerSynchronizer.Errorf(errorMessage, err)
		return
	}
}
func applyAllKeymanagerConfifuration(c client.Client, resolvedKeyManagers []eventhubTypes.ResolvedKeyManager) error {
	tokenIssuersFromK8s, _, err := retrieveAllTokenIssuers(c, "")
	if err != nil {
		return err
	}
	clonedTokenIssuerListFromK8s := make([]dpv1alpha2.TokenIssuer, len(tokenIssuersFromK8s))
	copy(clonedTokenIssuerListFromK8s, tokenIssuersFromK8s)
	clonedTokenIssuers := make([]eventhubTypes.ResolvedKeyManager, len(resolvedKeyManagers))
	copy(clonedTokenIssuers, resolvedKeyManagers)
	newTokenissuers := make([]eventhubTypes.ResolvedKeyManager, 0)
	sameTokenissuers := make([]eventhubTypes.ResolvedKeyManager, 0)
	for _, tokenIssuer := range clonedTokenIssuers {
		found := false
		unFilteredTokenIssuersFRomK8s := make([]dpv1alpha2.TokenIssuer, 0)
		for _, tokenIssuersFromK8s := range clonedTokenIssuerListFromK8s {
			if tokenIssuer.UUID == tokenIssuersFromK8s.Name {
				sameTokenissuers = append(sameTokenissuers, tokenIssuer)
				found = true
				break
			}
			unFilteredTokenIssuersFRomK8s = append(unFilteredTokenIssuersFRomK8s, tokenIssuersFromK8s)
		}
		clonedTokenIssuerListFromK8s = unFilteredTokenIssuersFRomK8s
		if !found {
			newTokenissuers = append(newTokenissuers, tokenIssuer)
		}
	}
	for _, tokenIssuer := range newTokenissuers {
		err := k8sclient.CreateAndUpdateTokenIssuersCR(tokenIssuer, c)
		if err != nil {
			return err
		}
		logger.LoggerSynchronizer.Debugf("Token Issuer created: %v", tokenIssuer)

	}
	for _, tokenIssuer := range sameTokenissuers {
		err := k8sclient.UpdateTokenIssuersCR(tokenIssuer, c)
		if err != nil {
			return err
		}
		logger.LoggerSynchronizer.Debugf("Token Issuer updated: %v", tokenIssuer)
	}
	logger.LoggerSynchronizer.Debugf("Deleted Token Issuers from K8s: %v", clonedTokenIssuerListFromK8s)
	for _, tokenissuer := range clonedTokenIssuerListFromK8s {
		err := k8sclient.DeleteTokenIssuerCR(c, tokenissuer)
		if err != nil {
			return err
		}
		logger.LoggerSynchronizer.Debugf("Token Issuer deleted: %v", tokenissuer)
	}
	return nil
}
func retrieveAllTokenIssuers(c client.Client, nextToken string) ([]dpv1alpha2.TokenIssuer, string, error) {
	conf, _ := config.ReadConfigs()
	tokenIssuerList := dpv1alpha2.TokenIssuerList{}
	resolvedTokenIssuerList := make([]dpv1alpha2.TokenIssuer, 0)
	var err error
	// Convert the label selector string into a labels.Selector
	labelSelector := labels.SelectorFromSet(labels.Set{"InitiateFrom": "CP"})

	opts := &client.ListOptions{
		Namespace:     conf.DataPlane.Namespace,
		LabelSelector: labelSelector,
	}
	if nextToken == "" {
		err = c.List(context.Background(), &tokenIssuerList, opts)
	} else {
		opts.Continue = nextToken
		err = c.List(context.Background(), &tokenIssuerList, opts)
	}
	if err != nil {
		logger.LoggerSynchronizer.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Failed to get application from k8s %v", err.Error()))
		return nil, "", err
	}
	resolvedTokenIssuerList = append(resolvedTokenIssuerList, tokenIssuerList.Items...)
	if tokenIssuerList.Continue != "" {
		tempTokenIssuerList, _, err := retrieveAllTokenIssuers(c, tokenIssuerList.Continue)
		if err != nil {
			return nil, "", err
		}
		resolvedTokenIssuerList = append(resolvedTokenIssuerList, tempTokenIssuerList...)
	}
	return resolvedTokenIssuerList, tokenIssuerList.Continue, nil
}
