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

package eventhub

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	loggers "github.com/sirupsen/logrus"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	common "github.com/wso2/product-apim-tooling/apim-apk-agent/internal/common"

	pkgAuth "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/auth"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/health"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/tlsutils"
)

const (
	authorizationBasic         string = "Basic "
	authorizationHeaderDefault string = "Authorization"
	internalWebAppEP           string = "internal/data/v1/"
	// ContextParam is required to call /apis endpoint
	ContextParam string = "context"
	// VersionParam is trequired to call /apis endpoint
	VersionParam string = "version"
	// GatewayLabelParam is trequired to call /apis endpoint
	GatewayLabelParam string = "gatewayLabel"
	// APIUUIDParam is required to call /apis endpoint
	APIUUIDParam string = "apiId"
	// ApisEndpoint is the resource path of /apis endpoint
	ApisEndpoint string = "apis"
)

var (
	// This set of variables are used just for Type resolution with reflect.
	// Hence no value needs to be assigned.
	subList           *types.SubscriptionList
	appList           *types.ApplicationList
	appKeyMappingList *types.ApplicationKeyMappingList
	apiList           *types.APIList

	resources = []resource{
		{
			endpoint:     "subscriptions",
			responseType: subList,
		},
		{
			endpoint:     "applications",
			responseType: appList,
		},
		{
			endpoint:     "application-key-mappings",
			responseType: appKeyMappingList,
		},
	}
	// APIListChannel is used to add apis
	APIListChannel chan response
	accessToken    string
	conf           *config.Config
)

type response struct {
	Error        error
	Payload      []byte
	ErrorCode    int
	Endpoint     string
	GatewayLabel string
	Type         interface{}
}

type resource struct {
	endpoint     string
	responseType interface{}
}

func init() {
	APIListChannel = make(chan response)
}

// LoadInitialData loads subscription/application and keymapping data from control-plane
func LoadInitialData(configFile *config.Config, initialAPIUUIDListMap map[string]int) {
	loggers.Info("Inside Load 1")
	conf = configFile
	accessToken = pkgAuth.GetBasicAuth(configFile.ControlPlane.Username, configFile.ControlPlane.Password)
	loggers.Info("accessToken: " + accessToken)
	var responseChannel = make(chan response)
	for _, url := range resources {
		// Create a local copy of the loop variable
		localURL := url
		loggers.Info("Inside loop" + localURL.endpoint)

		go InvokeService(localURL.endpoint, localURL.responseType, nil, responseChannel, 0)

		for {
			data := <-responseChannel
			logger.LoggerSync.Info("Receiving subscription data for an environment")
			if data.Payload != nil {
				logger.LoggerSync.Info("Payload data with subscription information recieved")
				loggers.Info("Payload data with subscription information recieved" + string(data.Payload))
				retrieveDataFromResponseChannel(data)
				break
			} else if data.ErrorCode >= 400 && data.ErrorCode < 500 {
				//Error handle
				health.SetControlPlaneRestAPIStatus(false)
			} else {
				// Keep the iteration going on until a response is received.
				// Error handle
				go func(d response, endpoint string, responseType interface{}) {
					// Retry fetching from control plane after a configured time interval
					if conf.ControlPlane.RetryInterval == 0 {
						// Assign default retry interval
						conf.ControlPlane.RetryInterval = 5
					}
					logger.LoggerSync.Debugf("Time Duration for retrying: %v", conf.ControlPlane.RetryInterval*time.Second)
					time.Sleep(conf.ControlPlane.RetryInterval * time.Second)
					logger.LoggerSync.Infof("Retrying to fetch APIs from control plane. Time Duration for the next retry: %v", conf.ControlPlane.RetryInterval*time.Second)
					go InvokeService(endpoint, responseType, nil, responseChannel, 0)
				}(data, localURL.endpoint, localURL.responseType)
			}
		}
	}
}

// InvokeService invokes the internal data resource
func InvokeService(endpoint string, responseType interface{}, queryParamMap map[string]string, c chan response,
	retryAttempt int) {

	serviceURL := conf.ControlPlane.ServiceURL + internalWebAppEP + endpoint
	loggers.Info("Payload service URL: " + serviceURL)
	// Create the request
	req, err := http.NewRequest("GET", serviceURL, nil)
	// gatewayLabel will only be required for apis endpoint
	gatewayLabel, ok := queryParamMap[GatewayLabelParam]
	if !ok {
		gatewayLabel = ""
	}
	if err != nil {
		c <- response{err, nil, 0, endpoint, gatewayLabel, responseType}
		logger.LoggerSubscription.Errorf("Error occurred while creating an HTTP request for serviceURL: "+serviceURL, err)
		return
	}
	queryParamMap = common.PopulateQueryParamForOrganizationID(queryParamMap)
	q := req.URL.Query()

	if queryParamMap != nil && len(queryParamMap) > 0 {
		// Making necessary query parameters for the request
		for queryParamKey, queryParamValue := range queryParamMap {
			q.Add(queryParamKey, queryParamValue)
		}
	}
	req.URL.RawQuery = q.Encode()

	// Check if TLS is enabled
	skipSSL := conf.ControlPlane.SkipSSLVerification

	// Setting authorization header
	req.Header.Set(authorizationHeaderDefault, authorizationBasic+accessToken)

	// Make the request
	//logger.LoggerSubscription.Debug("Sending the request to the control plane over the REST API: " + serviceURL)
	resp, err := tlsutils.InvokeControlPlane(req, skipSSL)

	if err != nil {
		if resp != nil {
			c <- response{err, nil, resp.StatusCode, endpoint, gatewayLabel, responseType}
		} else {
			c <- response{err, nil, 0, endpoint, gatewayLabel, responseType}
		}
		loggers.Info("Error occurred while calling the REST API: "+serviceURL, err)
		return
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		if err != nil {
			c <- response{err, nil, resp.StatusCode, endpoint, gatewayLabel, responseType}
			loggers.Info("Error occurred while reading the response received for: "+serviceURL, err)
			return
		}
		c <- response{nil, responseBytes, resp.StatusCode, endpoint, gatewayLabel, responseType}
	} else {
		c <- response{errors.New(string(responseBytes)), nil, resp.StatusCode, endpoint, gatewayLabel, responseType}
		loggers.Info("Failed to fetch data! "+serviceURL+" responded with "+strconv.Itoa(resp.StatusCode),
			err)
	}
}

func retrieveDataFromResponseChannel(response response) {
	responseType := reflect.TypeOf(response.Type).Elem()
	newResponse := reflect.New(responseType).Interface()
	err := json.Unmarshal(response.Payload, &newResponse)

	if err != nil {
		logger.LoggerSubscription.Errorf("Error occurred while unmarshalling the response received for: "+response.Endpoint, err)
	} else {
		switch t := newResponse.(type) {
		case *types.SubscriptionList:
			logger.LoggerSubscription.Debug("Received Subscription information.")
			subList = newResponse.(*types.SubscriptionList)
			MarshalMultipleSubscriptions(subList)
		case *types.ApplicationList:
			logger.LoggerSubscription.Debug("Received Application information.")
			appList = newResponse.(*types.ApplicationList)
			MarshalMultipleApplications(appList)
		case *types.ApplicationKeyMappingList:
			logger.LoggerSubscription.Debug("Received Application Key Mapping information.")
			appKeyMappingList = newResponse.(*types.ApplicationKeyMappingList)
			MarshalMultipleApplicationKeyMappings(appKeyMappingList)
		default:
			logger.LoggerSubscription.Debugf("Unknown type %T", t)
		}
	}
}
