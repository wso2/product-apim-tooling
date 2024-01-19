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
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	loggers "github.com/sirupsen/logrus"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"

	pkgAuth "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/auth"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/logging"
	sync "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/synchronizer"
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

const (
	// OrganizationID query parameter key.
	organizationID string = "organization"
	// OrganizationID query parameter value used when the global adapter is enabled and it is a shared gateway.
	commonOrganizationIDValue string = "ALL"
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
	apiUUIDList    []string
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
func LoadInitialData(configFile *config.Config) {
	conf = configFile
	accessToken = pkgAuth.GetBasicAuth(configFile.ControlPlane.Username, configFile.ControlPlane.Password)
	var responseChannel = make(chan response)
	for _, url := range resources {
		// Create a local copy of the loop variable
		localURL := url

		go InvokeService(localURL.endpoint, localURL.responseType, nil, responseChannel, 0)

		for {
			data := <-responseChannel
			logger.LoggerSync.Info("Receiving subscription data for an environment")
			if data.Payload != nil {
				logger.LoggerSync.Info("Payload data information received" + string(data.Payload))
				retrieveDataFromResponseChannel(data)
				break
			} else if data.ErrorCode >= 400 && data.ErrorCode < 500 {
				//Error handle
				logger.LoggerSync.Info("Error data information received")
				//health.SetControlPlaneRestAPIStatus(false)
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
	// Take the configured labels from the adapter
	configuredEnvs := conf.ControlPlane.EnvironmentLabels

	// If no environments are configured, default gateway label value is assigned.
	if len(configuredEnvs) == 0 {
		configuredEnvs = append(configuredEnvs, config.DefaultGatewayName)
	}
	for _, configuredEnv := range configuredEnvs {
		queryParamMap := make(map[string]string, 1)
		queryParamMap[GatewayLabelParam] = configuredEnv
		queryParamMap[organizationID] = commonOrganizationIDValue
		go InvokeService(ApisEndpoint, apiList, queryParamMap, APIListChannel, 0)
		for {
			data := <-APIListChannel
			logger.LoggerSync.Debug("Receiving API information for an environment")
			if data.Payload != nil {
				loggers.Info("Payload data with API information recieved" + string(data.Payload))
				retrieveAPIList(data)
				break
			} else if data.ErrorCode >= 400 && data.ErrorCode < 500 {
				logger.LoggerSync.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprintf("Error occurred when retrieving Subscription information from the control plane: %v", data.Error),
					Severity:  logging.CRITICAL,
					ErrorCode: 1600,
				})
				//health.SetControlPlaneRestAPIStatus(false)
			} else {
				// Keep the iteration going on until a response is recieved.
				logger.LoggerSync.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprintf("Error occurred while fetching data from control plane: %v", data.Error),
					Severity:  logging.MAJOR,
					ErrorCode: 1601,
				})
				go func(d response) {
					// Retry fetching from control plane after a configured time interval
					if conf.ControlPlane.RetryInterval == 0 {
						// Assign default retry interval
						conf.ControlPlane.RetryInterval = 5
					}
					logger.LoggerSync.Debugf("Time Duration for retrying: %v", conf.ControlPlane.RetryInterval*time.Second)
					time.Sleep(conf.ControlPlane.RetryInterval * time.Second)
					logger.LoggerSync.Infof("Retrying to fetch APIs from control plane. Time Duration for the next retry: %v", conf.ControlPlane.RetryInterval*time.Second)
					go InvokeService(ApisEndpoint, apiList, queryParamMap, APIListChannel, 0)
				}(data)
			}
		}
	}
	if apiUUIDList == nil || len(apiUUIDList) == 0 {
		loggers.Info("Empty API List Recieved in fetching")
	} else {
		FetchAPIsOnStartUp(conf, apiUUIDList)
	}
}

// InvokeService invokes the internal data resource
func InvokeService(endpoint string, responseType interface{}, queryParamMap map[string]string, c chan response,
	retryAttempt int) {

	serviceURL := conf.ControlPlane.ServiceURL + internalWebAppEP + endpoint
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
	if reflect.TypeOf(responseType) == reflect.TypeOf(types.APIList{}) {
		req.Header.Set("x-wso2-tenant", "ALL")
	}

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
		loggers.Info("Error occurred while unmarshalling the response received for: "+response.Endpoint, err)
	} else {
		switch t := newResponse.(type) {
		case *types.SubscriptionList:
			loggers.Info("Received Subscription information.")
			subList = newResponse.(*types.SubscriptionList)
			MarshalMultipleSubscriptions(subList)
		case *types.ApplicationList:
			loggers.Info("Received Application information.")
			appList = newResponse.(*types.ApplicationList)
			MarshalMultipleApplications(appList)
		case *types.ApplicationKeyMappingList:
			loggers.Info("Received Application Key Mapping information.")
			appKeyMappingList = newResponse.(*types.ApplicationKeyMappingList)
			MarshalMultipleApplicationKeyMappings(appKeyMappingList)
		default:
			logger.LoggerSubscription.Debugf("Unknown type %T", t)
		}
	}
}

// FetchAPIsOnStartUp APIs from control plane during the server start up and push them
// to the router and enforcer components.
func FetchAPIsOnStartUp(conf *config.Config, apiUUIDList []string) {
	// Populate data from config.
	envs := conf.ControlPlane.EnvironmentLabels

	// Create a channel for the byte slice (response from the APIs from control plane)
	c := make(chan sync.SyncAPIResponse)

	var queryParamMap map[string]string
	//Get API details.
	if apiUUIDList != nil {
		GetAPIs(c, nil, envs, sync.APIArtifactEndpoint, true, apiUUIDList, queryParamMap)
	}
	for i := 0; i < 1; i++ {
		data := <-c
		logger.LoggerMsg.Info("Receiving data for an environment")
		if data.Resp != nil {
			// Reading the root zip
			zipReader, err := zip.NewReader(bytes.NewReader(data.Resp), int64(len(data.Resp)))

			// apiFiles represents zipped API files fetched from API Manager
			apiFiles := make(map[string]*zip.File)
			// Read the .zip files within the root apis.zip and add apis to apiFiles array.
			for _, file := range zipReader.File {
				apiFiles[file.Name] = file
				loggers.Info("API file found: " + file.Name)
				// Todo:
			}
			logger.LoggerMsg.Info("Err", err)
			//health.SetControlPlaneRestAPIStatus(err == nil)

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
	logger.LoggerMsg.Info("Fetching APIs at startup is completed...")
}

// GetAPIs function calls the FetchAPIs() with relevant environment labels defined in the config.
func GetAPIs(c chan sync.SyncAPIResponse, id *string, envs []string, endpoint string, sendType bool, apiUUIDList []string,
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

func retrieveAPIList(response response) []string {

	responseType := reflect.TypeOf(response.Type).Elem()
	newResponse := reflect.New(responseType).Interface()
	if response.Error == nil && response.Payload != nil {
		err := json.Unmarshal(response.Payload, &newResponse)
		if err != nil {
			logger.LoggerSubscription.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error occurred while unmarshalling the APIList response received for: %v %v", response.Endpoint, err.Error()),
				Severity:  logging.MAJOR,
				ErrorCode: 1602,
			})
		} else {
			switch t := newResponse.(type) {
			case *types.APIList:
				apiListResponse := newResponse.(*types.APIList)
				if apiListResponse.List != nil {
					for _, api := range apiListResponse.List {
						apiUUIDList = append(apiUUIDList, api.UUID)
					}
				}
				loggers.Info("Received API information.", apiUUIDList)
				return apiUUIDList
			default:
				logger.LoggerSubscription.Warnf("APIList Type DTO is not recieved. Unknown type %T", t)
			}
		}
	}
	return nil
}
