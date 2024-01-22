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

// Package messaging holds the implementation for event listeners functions
package messaging

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/wso2/apk/common-go-libs/constants"
	event "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/internal/synchronizer"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	logger "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/logging"
	msg "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/messaging"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/utils"
)

// constant variables
const (
	apiEventType                = "API"
	applicationEventType        = "APPLICATION"
	subscriptionEventType       = "SUBSCRIPTION"
	scopeEvenType               = "SCOPE"
	policyEventType             = "POLICY"
	removeAPIFromGateway        = "REMOVE_API_FROM_GATEWAY"
	deployAPIToGateway          = "DEPLOY_API_IN_GATEWAY"
	applicationRegistration     = "APPLICATION_REGISTRATION_CREATE"
	removeApplicationKeyMapping = "REMOVE_APPLICATION_KEYMAPPING"
	apiLifeCycleChange          = "LIFECYCLE_CHANGE"
	applicationCreate           = "APPLICATION_CREATE"
	applicationUpdate           = "APPLICATION_UPDATE"
	applicationDelete           = "APPLICATION_DELETE"
	subscriptionCreate          = "SUBSCRIPTIONS_CREATE"
	subscriptionUpdate          = "SUBSCRIPTIONS_UPDATE"
	subscriptionDelete          = "SUBSCRIPTIONS_DELETE"
	policyCreate                = "POLICY_CREATE"
	policyUpdate                = "POLICY_UPDATE"
	policyDelete                = "POLICY_DELETE"
	blockedStatus               = "BLOCKED"
	apiUpdate                   = "API_UPDATE"
)

// var variables
var (
	ScopeList = make([]types.Scope, 0)
	// timestamps needs to be maintained as it is not guranteed to receive them in order,
	// hence older events should be discarded
	apiListTimeStampMap          = make(map[string]int64, 0)
	subsriptionsListTimeStampMap = make(map[string]int64, 0)
	applicationListTimeStampMap  = make(map[string]int64, 0)
)

// handleNotification to process
func handleNotification() {
	conf, _ := config.ReadConfigs()
	for d := range msg.NotificationChannel {
		var notification msg.EventNotification
		notificationErr := parseNotificationJSONEvent([]byte(string(d.Body)), &notification)
		if notificationErr != nil {
			continue
		}
		logger.LoggerMsg.Infof("Event %s is received", notification.Event.PayloadData.EventType)
		err := processNotificationEvent(conf, &notification)
		if err != nil {
			continue
		}
		d.Ack(false)
	}
	logger.LoggerMsg.Infof("handle: deliveries channel closed")
}

func processNotificationEvent(conf *config.Config, notification *msg.EventNotification) error {
	var eventType string
	var decodedByte, err = base64.StdEncoding.DecodeString(notification.Event.PayloadData.Event)
	if err != nil {
		if _, ok := err.(base64.CorruptInputError); ok {
			logger.LoggerMsg.Error("\nbase64 input is corrupt, check the provided key")
		}
		logger.LoggerMsg.Errorf("Error occurred while decoding the notification event %v. "+
			"Hence dropping the event", err)
		return err
	}
	logger.LoggerMsg.Debugf("\n\n[%s]", decodedByte)
	eventType = notification.Event.PayloadData.EventType
	if strings.Contains(eventType, apiLifeCycleChange) {
		handleLifeCycleEvents(decodedByte)
	} else if strings.Contains(eventType, apiEventType) {
		handleAPIEvents(decodedByte, eventType, conf)
	} else if strings.Contains(eventType, applicationEventType) {
		handleApplicationEvents(decodedByte, eventType)
	} else if strings.Contains(eventType, subscriptionEventType) {
		handleSubscriptionEvents(decodedByte, eventType)
	} else if strings.Contains(eventType, policyEventType) {
		handlePolicyEvents(decodedByte, eventType)
	}
	// other events will ignore including HEALTH_CHECK event
	return nil
}

// handleDefaultVersionUpdate will redeploy default versioned API.
// API runtime artifact doesn't get updated in CP side when default version is updated
// (isDefaultVersion prop in apiYaml is not updated). API deployment or should happen
// for it to get updated. However we need to redeploy the API when there is a default
// version change. For that we call `/apis` endpoint to get updated API metadata (this
// contains the updated `isDefaultVersion` field). Now we proceed with fetching runtime
// artifact from the CP. When creating CC deployment objects we refer to updated `APIList`
// map and update runtime artifact's `isDefaultVersion` field to correctly deploy default
// versioned API.
func handleDefaultVersionUpdate(event msg.APIEvent) {
	// deployedEnvs := xds.GetDeployedEnvironments(event.UUID)
	// for _, env := range deployedEnvs {
	// 	query := make(map[string]string, 3)
	// 	query[eh.GatewayLabelParam] = env
	// 	query[eh.ContextParam] = event.APIContext
	// 	query[eh.VersionParam] = event.APIVersion
	// 	eh.UpdateAPIMetadataFromCP(query)
	// }

	// synchronizer.FetchAPIsFromControlPlane(event.UUID, deployedEnvs)
}

// handleAPIEvents to process api related data
func handleAPIEvents(data []byte, eventType string, conf *config.Config) {
	var (
		apiEvent              msg.APIEvent
		currentTimeStamp      int64 = apiEvent.Event.TimeStamp
		isDefaultVersionEvent bool
		apiUUIDList           []string
	)

	apiEventErr := json.Unmarshal([]byte(string(data)), &apiEvent)
	if apiEventErr != nil {
		logger.LoggerMsg.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error occurred while unmarshalling API event data %v", apiEventErr),
			Severity:  logging.MAJOR,
			ErrorCode: 2004,
		})
		return
	}

	if !belongsToTenant(apiEvent.TenantDomain) {
		apiName := apiEvent.APIName
		if apiEvent.APIName == "" {
			apiName = apiEvent.Name
		}
		apiVersion := apiEvent.Version
		if apiEvent.Version == "" {
			apiVersion = apiEvent.Version
		}
		logger.LoggerMsg.Debugf("API event for the API %s:%s is dropped due to having non related tenantDomain : %s",
			apiName, apiVersion, apiEvent.TenantDomain)
		return
	}

	apiEventObj := types.API{UUID: apiEvent.UUID, APIID: apiEvent.APIID, Name: apiEvent.APIName,
		Context: apiEvent.APIContext, Version: apiEvent.APIVersion, Provider: apiEvent.APIProvider}

	logger.LoggerMsg.Infof("API event data %v", apiEventObj)

	isDefaultVersionEvent = isDefaultVersionUpdate(apiEvent)

	if isDefaultVersionEvent {
		handleDefaultVersionUpdate(apiEvent)
		return
	}

	//Per each revision, synchronization should happen.
	if strings.EqualFold(deployAPIToGateway, apiEvent.Event.Type) {
		//go synchronizer.FetchAPIsFromControlPlane(apiEvent.UUID, apiEvent.GatewayLabels)
		apiUUIDList = append(apiUUIDList, apiEvent.UUID)
		go synchronizer.FetchAPIsOnEvent(conf, apiUUIDList)
	}

	for _, env := range apiEvent.GatewayLabels {
		if isLaterEvent(apiListTimeStampMap, apiEvent.UUID+":"+env, currentTimeStamp) {
			break
		}
		// removeFromGateway event with multiple labels could only appear when the API is subjected
		// to delete. Hence we could simply delete after checking against just one iteration.
		if strings.EqualFold(removeAPIFromGateway, apiEvent.Event.Type) {
			// xds.DeleteAPIWithAPIMEvent(apiEvent.UUID, apiEvent.TenantDomain, apiEvent.GatewayLabels, "")
			// for _, env := range apiEvent.GatewayLabels {
			// 	xdsAPIList := xds.DeleteAPIAndReturnList(apiEvent.UUID, apiEvent.TenantDomain, env)
			// 	if xdsAPIList != nil {
			// 		xds.UpdateEnforcerAPIList(env, xdsAPIList)
			// 	}
			// }
			break
		}
		if strings.EqualFold(deployAPIToGateway, apiEvent.Event.Type) {
			conf, _ := config.ReadConfigs()
			configuredEnvs := conf.ControlPlane.EnvironmentLabels
			if len(configuredEnvs) == 0 {
				configuredEnvs = append(configuredEnvs, config.DefaultGatewayName)
			}
			// for _, configuredEnv := range configuredEnvs {
			// 	if configuredEnv == env {
			// 			if xds.CheckIfAPIMetadataIsAlreadyAvailable(apiEvent.UUID, env) {
			// 				logger.LoggerInternalMsg.Debugf("API Metadata for api Id: %s is not updated as it already exists", apiEvent.UUID)
			// 				continue
			// 			}
			// 			logger.LoggerInternalMsg.Debugf("Fetching Metadata for api Id: %s ", apiEvent.UUID)
			// 			queryParamMap := make(map[string]string, 3)
			// 			queryParamMap[eh.GatewayLabelParam] = configuredEnv
			// 			queryParamMap[eh.ContextParam] = apiEvent.Context
			// 			queryParamMap[eh.VersionParam] = apiEvent.Version
			// 			var apiList *types.APIList
			// 			go eh.InvokeService(eh.ApisEndpoint, apiList, queryParamMap, eh.APIListChannel, 0)
			// 		}
			// 	}
		}
	}
}

func handleLifeCycleEvents(data []byte) {
	var apiEvent msg.APIEvent
	apiLCEventErr := json.Unmarshal([]byte(string(data)), &apiEvent)
	if apiLCEventErr != nil {
		logger.LoggerMsg.Errorf("Error occurred while unmarshalling Lifecycle event data %v", apiLCEventErr)
		return
	}
	if !belongsToTenant(apiEvent.TenantDomain) {
		logger.LoggerMsg.Debugf("API Lifecycle event for the API %s:%s is dropped due to having non related tenantDomain : %s",
			apiEvent.APIName, apiEvent.APIVersion, apiEvent.TenantDomain)
		return
	}

	apiEventObj := types.API{UUID: apiEvent.UUID, APIID: apiEvent.APIID, Name: apiEvent.APIName,
		Context: apiEvent.APIContext, Version: apiEvent.APIVersion, Provider: apiEvent.APIProvider}

	logger.LoggerMsg.Infof("API event data %v", apiEventObj)

	conf, _ := config.ReadConfigs()
	configuredEnvs := conf.ControlPlane.EnvironmentLabels
	logger.LoggerMsg.Debugf("%s : %s API life cycle state change event triggered", apiEvent.APIName, apiEvent.APIVersion)
	if len(configuredEnvs) == 0 {
		configuredEnvs = append(configuredEnvs, config.DefaultGatewayName)
	}
	// for _, configuredEnv := range configuredEnvs {
	// 	xdsAPIList := xds.MarshalAPIForLifeCycleChangeEventAndReturnList(apiEvent.UUID, apiEvent.APIStatus, configuredEnv)
	// 	if xdsAPIList != nil {
	// 		xds.UpdateEnforcerAPIList(configuredEnv, xdsAPIList)
	// 	}
	// }
}

// handleApplicationEvents to process application related events
func handleApplicationEvents(data []byte, eventType string) {
	if strings.EqualFold(applicationRegistration, eventType) ||
		strings.EqualFold(removeApplicationKeyMapping, eventType) {
		var applicationRegistrationEvent msg.ApplicationRegistrationEvent
		appRegEventErr := json.Unmarshal([]byte(string(data)), &applicationRegistrationEvent)
		if appRegEventErr != nil {
			logger.LoggerMsg.Errorf("Error occurred while unmarshalling Application Registration event data %v", appRegEventErr)
			return
		}

		if !belongsToTenant(applicationRegistrationEvent.TenantDomain) {
			logger.LoggerMsg.Debugf("Application Registration event for the Consumer Key : %s is dropped due to having non related tenantDomain : %s",
				applicationRegistrationEvent.ConsumerKey, applicationRegistrationEvent.TenantDomain)
			return
		}
		applicationKeyMappingEvent := event.ApplicationKeyMapping{ApplicationUUID: applicationRegistrationEvent.ApplicationUUID,
			SecurityScheme:        "OAuth2",
			ApplicationIdentifier: applicationRegistrationEvent.ConsumerKey,
			KeyType:               applicationRegistrationEvent.KeyType,
			Organization:          applicationRegistrationEvent.TenantDomain,
			EnvID:                 "Default",
		}
		if strings.EqualFold(applicationRegistration, eventType) {
			event := event.Event{Type: constants.ApplicationKeyMappingCreated,
				Uuid:                  uuid.New().String(),
				TimeStamp:             applicationRegistrationEvent.TimeStamp,
				ApplicationKeyMapping: &applicationKeyMappingEvent,
			}
			utils.SendEvent(&event)
		} else if strings.EqualFold(removeApplicationKeyMapping, eventType) {
			event := event.Event{Type: constants.ApplicationKeyMappingDeleted,
				Uuid:                  uuid.New().String(),
				TimeStamp:             applicationRegistrationEvent.TimeStamp,
				ApplicationKeyMapping: &applicationKeyMappingEvent,
			}
			utils.SendEvent(&event)
		}
	} else {
		var applicationEvent msg.ApplicationEvent
		appEventErr := json.Unmarshal([]byte(string(data)), &applicationEvent)
		if appEventErr != nil {
			logger.LoggerMsg.Errorf("Error occurred while unmarshalling Application event data %v", appEventErr)
			return
		}

		if !belongsToTenant(applicationEvent.TenantDomain) {
			logger.LoggerMsg.Debugf("Application event for the Application : %s (with uuid %s) is dropped due to having non related tenantDomain : %s",
				applicationEvent.ApplicationName, applicationEvent.UUID, applicationEvent.TenantDomain)
			return
		}

		logger.LoggerMsg.Infof("Application event data %v", applicationEvent)

		if isLaterEvent(applicationListTimeStampMap, fmt.Sprint(applicationEvent.ApplicationID), applicationEvent.TimeStamp) {
			return
		}

		applicationGrpcEvent := event.Application{Uuid: applicationEvent.UUID,
			Name:         applicationEvent.ApplicationName,
			Owner:        applicationEvent.Subscriber,
			Organization: applicationEvent.TenantDomain,
			Attributes:   marshalAppAttributes(applicationEvent.Attributes),
		}
		if applicationEvent.Event.Type == applicationCreate {
			event := event.Event{Type: constants.ApplicationCreated, Uuid: uuid.New().String(), TimeStamp: applicationEvent.TimeStamp, Application: &applicationGrpcEvent}
			utils.SendEvent(&event)
		} else if applicationEvent.Event.Type == applicationUpdate {
			event := event.Event{Type: constants.ApplicationUpdated, Uuid: uuid.New().String(), TimeStamp: applicationEvent.TimeStamp, Application: &applicationGrpcEvent}
			utils.SendEvent(&event)
		} else if applicationEvent.Event.Type == applicationDelete {
			event := event.Event{Type: constants.ApplicationDeleted, Uuid: uuid.New().String(), TimeStamp: applicationEvent.TimeStamp, Application: &applicationGrpcEvent}
			utils.SendEvent(&event)
		} else {
			logger.LoggerAdapter.Warnf("Application Event Type is not recognized for the Event under "+
				"Application UUID %s", applicationEvent.UUID)
			return
		}
	}
}
func marshalAppAttributes(attributes interface{}) map[string]string {
	attributesMap := make(map[string]string)
	if attributes != nil {
		for key, value := range attributes.(map[string]interface{}) {
			attributesMap[key] = value.(string)
		}
	}
	return attributesMap
}

// handleSubscriptionRelatedEvents to process subscription related events
func handleSubscriptionEvents(data []byte, eventType string) {
	var subscriptionEvent msg.SubscriptionEvent
	subEventErr := json.Unmarshal([]byte(string(data)), &subscriptionEvent)
	if subEventErr != nil {
		logger.LoggerMsg.Errorf("Error occurred while unmarshalling Subscription event data %v", subEventErr)
		return
	}
	if !belongsToTenant(subscriptionEvent.TenantDomain) {
		logger.LoggerMsg.Debugf("Subscription event for the Application : %s and API %s is dropped due to having non related tenantDomain : %s",
			subscriptionEvent.ApplicationUUID, subscriptionEvent.APIUUID, subscriptionEvent.TenantDomain)
		return
	}

	if isLaterEvent(subsriptionsListTimeStampMap, fmt.Sprint(subscriptionEvent.SubscriptionID), subscriptionEvent.TimeStamp) {
		return
	}

	subscription := event.Subscription{Uuid: subscriptionEvent.SubscriptionUUID,
		SubStatus:     subscriptionEvent.SubscriptionState,
		Organization:  subscriptionEvent.TenantDomain,
		SubscribedApi: &event.SubscribedAPI{Name: subscriptionEvent.APIUUID, Version: "1.0.0"},
	}

	hasher := sha1.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s", subscriptionEvent.ApplicationUUID, subscriptionEvent.SubscriptionUUID)))

	// Get the final hash sum
	hashSum := hasher.Sum(nil)

	// Convert the hash sum to a hexadecimal string
	hashString := hex.EncodeToString(hashSum)
	applicationMapping := event.ApplicationMapping{Uuid: hashString, ApplicationRef: subscriptionEvent.ApplicationUUID, SubscriptionRef: subscriptionEvent.SubscriptionUUID, Organization: subscriptionEvent.TenantDomain}
	if subscriptionEvent.Event.Type == subscriptionCreate {
		subsEvent := event.Event{Uuid: uuid.New().String(), Type: constants.SubscriptionCreated, TimeStamp: subscriptionEvent.TimeStamp, Subscription: &subscription}
		utils.SendEvent(&subsEvent)
		applicationMappingEvent := event.Event{Uuid: uuid.New().String(), Type: constants.ApplicationMappingCreated, TimeStamp: subscriptionEvent.TimeStamp, ApplicationMapping: &applicationMapping}
		utils.SendEvent(&applicationMappingEvent)
	} else if subscriptionEvent.Event.Type == subscriptionUpdate {
		subsEvent := event.Event{Uuid: uuid.New().String(), Type: constants.SubscriptionUpdated, TimeStamp: subscriptionEvent.TimeStamp, Subscription: &subscription}
		utils.SendEvent(&subsEvent)
		applicationMappingEvent := event.Event{Uuid: uuid.New().String(), Type: constants.ApplicationMappingUpdated, TimeStamp: subscriptionEvent.TimeStamp, ApplicationMapping: &applicationMapping}
		utils.SendEvent(&applicationMappingEvent)

	} else if subscriptionEvent.Event.Type == subscriptionDelete {
		subsEvent := event.Event{Uuid: uuid.New().String(), Type: constants.SubscriptionDeleted, TimeStamp: subscriptionEvent.TimeStamp, Subscription: &subscription}
		utils.SendEvent(&subsEvent)
		applicationMappingEvent := event.Event{Uuid: uuid.New().String(), Type: constants.ApplicationMappingDeleted, TimeStamp: subscriptionEvent.TimeStamp, ApplicationMapping: &applicationMapping}
		utils.SendEvent(&applicationMappingEvent)
	}
}

// handlePolicyRelatedEvents to process policy related events
func handlePolicyEvents(data []byte, eventType string) {
	var policyEvent msg.PolicyInfo
	policyEventErr := json.Unmarshal([]byte(string(data)), &policyEvent)
	if policyEventErr != nil {
		logger.LoggerMsg.Errorf("Error occurred while unmarshalling Throttling Policy event data %v", policyEventErr)
		return
	}
	// TODO: Handle policy events
	if strings.EqualFold(eventType, policyCreate) {
		logger.LoggerMsg.Infof("Policy: %s for policy type: %s", policyEvent.PolicyName, policyEvent.PolicyType)
	} else if strings.EqualFold(eventType, policyUpdate) {
		logger.LoggerMsg.Infof("Policy: %s for policy type: %s", policyEvent.PolicyName, policyEvent.PolicyType)
	} else if strings.EqualFold(eventType, policyDelete) {
		logger.LoggerMsg.Infof("Policy: %s for policy type: %s", policyEvent.PolicyName, policyEvent.PolicyType)
	}

	if strings.EqualFold(applicationEventType, policyEvent.PolicyType) {
		applicationPolicy := types.ApplicationPolicy{ID: policyEvent.PolicyID, TenantID: policyEvent.Event.TenantID,
			Name: policyEvent.PolicyName, QuotaType: policyEvent.QuotaType}

		logger.LoggerMsg.Infof("ApplicationPolicy event data %v", applicationPolicy)
		// var applicationPolicyList *subscription.ApplicationPolicyList
		// if policyEvent.Event.Type == policyCreate {
		// 	applicationPolicyList = xds.MarshalApplicationPolicyEventAndReturnList(&applicationPolicy, xds.CreateEvent)
		// } else if policyEvent.Event.Type == policyUpdate {
		// 	applicationPolicyList = xds.MarshalApplicationPolicyEventAndReturnList(&applicationPolicy, xds.UpdateEvent)
		// } else if policyEvent.Event.Type == policyDelete {
		// 	applicationPolicyList = xds.MarshalApplicationPolicyEventAndReturnList(&applicationPolicy, xds.DeleteEvent)
		// } else {
		// 	logger.LoggerInternalMsg.Warnf("ApplicationPolicy Event Type is not recognized for the Event under "+
		// 		" policy name %s", policyEvent.PolicyName)
		// 	return
		// }
		// xds.UpdateEnforcerApplicationPolicies(applicationPolicyList)

	} else if strings.EqualFold(subscriptionEventType, policyEvent.PolicyType) {
		var subscriptionPolicyEvent msg.SubscriptionPolicyEvent
		subPolicyErr := json.Unmarshal([]byte(string(data)), &subscriptionPolicyEvent)
		if subPolicyErr != nil {
			logger.LoggerMsg.Errorf("Error occurred while unmarshalling Subscription Policy event data %v", subPolicyErr)
			return
		}

		subscriptionPolicy := types.SubscriptionPolicy{ID: subscriptionPolicyEvent.PolicyID, TenantID: -1,
			Name: subscriptionPolicyEvent.PolicyName, QuotaType: subscriptionPolicyEvent.QuotaType,
			GraphQLMaxComplexity: subscriptionPolicyEvent.GraphQLMaxComplexity,
			GraphQLMaxDepth:      subscriptionPolicyEvent.GraphQLMaxDepth, RateLimitCount: subscriptionPolicyEvent.RateLimitCount,
			RateLimitTimeUnit: subscriptionPolicyEvent.RateLimitTimeUnit, StopOnQuotaReach: subscriptionPolicyEvent.StopOnQuotaReach,
			TenantDomain: subscriptionPolicyEvent.TenantDomain, TimeStamp: subscriptionPolicyEvent.TimeStamp}

		logger.LoggerMsg.Infof("SubscriptionPolicy event data %v", subscriptionPolicy)

		// var subscriptionPolicyList *subscription.SubscriptionPolicyList
		// if subscriptionPolicyEvent.Event.Type == policyCreate {
		// 	subscriptionPolicyList = xds.MarshalSubscriptionPolicyEventAndReturnList(&subscriptionPolicy, xds.CreateEvent)
		// } else if subscriptionPolicyEvent.Event.Type == policyUpdate {
		// 	subscriptionPolicyList = xds.MarshalSubscriptionPolicyEventAndReturnList(&subscriptionPolicy, xds.UpdateEvent)
		// } else if subscriptionPolicyEvent.Event.Type == policyDelete {
		// 	subscriptionPolicyList = xds.MarshalSubscriptionPolicyEventAndReturnList(&subscriptionPolicy, xds.DeleteEvent)
		// } else {
		// 	logger.LoggerInternalMsg.Warnf("SubscriptionPolicy Event Type is not recognized for the Event under "+
		// 		" policy name %s", policyEvent.PolicyName)
		// 	return
		// }
		// xds.UpdateEnforcerSubscriptionPolicies(subscriptionPolicyList)
	}
}

func isLaterEvent(timeStampMap map[string]int64, mapKey string, currentTimeStamp int64) bool {
	if timeStamp, ok := timeStampMap[mapKey]; ok {
		if timeStamp > currentTimeStamp {
			return true
		}
	}
	timeStampMap[mapKey] = currentTimeStamp
	return false
}

func isDefaultVersionUpdate(event msg.APIEvent) bool {
	return strings.EqualFold(apiUpdate, event.Event.Type) && strings.EqualFold("DEFAULT_VERSION", event.Action)
}

func belongsToTenant(tenantDomain string) bool {
	// TODO : enable this once the events are fixed in apim
	// return config.GetControlPlaneConnectedTenantDomain() == tenantDomain
	return true
}

func parseNotificationJSONEvent(data []byte, notification *msg.EventNotification) error {
	unmarshalErr := json.Unmarshal(data, &notification)
	if unmarshalErr != nil {
		logger.LoggerMsg.Errorf("Error occurred while unmarshalling "+
			"notification event data %v. Hence dropping the event", unmarshalErr)
	}
	return unmarshalErr
}
