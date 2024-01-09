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
	logger "github.com/sirupsen/logrus"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/config"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
)

// SubscriptionList for struct list of applications
type SubscriptionList struct {
	List []Subscription `json:"list"`
}

// Application for struct application
type Application struct {
	UUID         string            `json:"uuid"`
	ID           int32             `json:"id" json:"applicationId"`
	Name         string            `json:"name" json:"applicationName"`
	SubName      string            `json:"subName" json:"subscriber"`
	Policy       string            `json:"policy" json:"applicationPolicy"`
	TokenType    string            `json:"tokenType"`
	Attributes   map[string]string `json:"attributes"`
	TenantID     int32             `json:"tenanId,omitempty"`
	TenantDomain string            `json:"tenanDomain,omitempty"`
	TimeStamp    int64             `json:"timeStamp,omitempty"`
}

// ApplicationList for struct list of application
type ApplicationList struct {
	List []Application `json:"list"`
}

// ApplicationKeyMapping for struct applicationKeyMapping
type ApplicationKeyMapping struct {
	ApplicationID   int32  `json:"applicationId"`
	ApplicationUUID string `json:"applicationUUID"`
	ConsumerKey     string `json:"consumerKey"`
	KeyType         string `json:"keyType"`
	KeyManager      string `json:"keyManager"`
	TenantID        int32  `json:"tenanId,omitempty"`
	TenantDomain    string `json:"tenanDomain,omitempty"`
	TimeStamp       int64  `json:"timeStamp,omitempty"`
}

// ApplicationKeyMappingList for struct list of applicationKeyMapping
type ApplicationKeyMappingList struct {
	List []ApplicationKeyMapping `json:"list"`
}

// APIs for struct Api
type APIs struct {
	APIID            int    `json:"apiId"`
	UUID             string `json:"uuid"`
	Provider         string `json:"provider" json:"apiProvider"`
	Name             string `json:"name" json:"apiName"`
	Version          string `json:"version" json:"apiVersion"`
	Context          string `json:"context" json:"apiContext"`
	Policy           string `json:"policy"`
	APIType          string `json:"apiType"`
	IsDefaultVersion bool   `json:"isDefaultVersion"`
	APIStatus        string `json:"status"`
	TenantID         int32  `json:"tenanId,omitempty"`
	TenantDomain     string `json:"tenanDomain,omitempty"`
	TimeStamp        int64  `json:"timeStamp,omitempty"`
}

// APIList for struct ApiList
type APIList struct {
	List []APIs `json:"list"`
}

// Subscription for struct subscription
type Subscription struct {
	SubscriptionID    int32  `json:"subscriptionId"`
	SubscriptionUUID  string `json:"subscriptionUUID"`
	PolicyID          string `json:"policyId"`
	APIID             int32  `json:"apiId"`
	APIUUID           string `json:"apiUUID"`
	AppID             int32  `json:"appId" json:"applicationId"`
	ApplicationUUID   string `json:"applicationUUID"`
	SubscriptionState string `json:"subscriptionState"`
	TenantID          int32  `json:"tenanId,omitempty"`
	TenantDomain      string `json:"tenanDomain,omitempty"`
	TimeStamp         int64  `json:"timeStamp,omitempty"`
}

var (
	// APIListMap has the following mapping label -> apiUUID -> API (Metadata)
	APIListMap map[string]map[string]APIs
	// SubscriptionMap contains the subscriptions recieved from API Manager Control Plane
	SubscriptionMap map[int32]Subscription
	// ApplicationMap contains the applications recieved from API Manager Control Plane
	ApplicationMap map[string]Application
	// ApplicationKeyMappingMap contains the application key mappings recieved from API Manager Control Plane
	ApplicationKeyMappingMap map[string]ApplicationKeyMapping
)

// MarshalMultipleApplications is used to update the applicationList during the startup where
// multiple applications are pulled at once. And then it returns the ApplicationList.
func MarshalMultipleApplications(appList *types.ApplicationList) map[string]Application {
	resourceMap := make(map[string]Application)
	for _, application := range appList.List {
		applicationSub := marshalApplication(&application)
		resourceMap[application.UUID] = applicationSub
	}
	ApplicationMap = resourceMap
	for appID, app := range ApplicationMap {
		logger.Info("Application: , Description:", appID, app)
	}
	return ApplicationMap
}

// MarshalMultipleApplicationKeyMappings is used to update the application key mappings during the startup where
// multiple key mappings are pulled at once. And then it returns the ApplicationKeyMappingList.
func MarshalMultipleApplicationKeyMappings(keymappingList *types.ApplicationKeyMappingList) map[string]ApplicationKeyMapping {
	resourceMap := make(map[string]ApplicationKeyMapping)
	for _, keyMapping := range keymappingList.List {
		applicationKeyMappingReference := GetApplicationKeyMappingReference(&keyMapping)
		keyMappingSub := marshalKeyMapping(&keyMapping)
		resourceMap[applicationKeyMappingReference] = keyMappingSub
	}
	ApplicationKeyMappingMap = resourceMap
	return ApplicationKeyMappingMap
}

// MarshalMultipleSubscriptions is used to update the subscriptions during the startup where
// multiple subscriptions are pulled at once. And then it returns the SubscriptionList.
func MarshalMultipleSubscriptions(subscriptionsList *types.SubscriptionList) map[int32]Subscription {
	resourceMap := make(map[int32]Subscription)
	for _, sb := range subscriptionsList.List {
		resourceMap[sb.SubscriptionID] = marshalSubscription(&sb)
	}
	SubscriptionMap = resourceMap
	return SubscriptionMap
}

func marshalSubscription(subscriptionInternal *types.Subscription) Subscription {
	sub := Subscription{
		SubscriptionID:    subscriptionInternal.SubscriptionID,
		PolicyID:          subscriptionInternal.PolicyID,
		APIID:             subscriptionInternal.APIID,
		AppID:             subscriptionInternal.AppID,
		SubscriptionState: subscriptionInternal.SubscriptionState,
		TimeStamp:         subscriptionInternal.TimeStamp,
		TenantID:          subscriptionInternal.TenantID,
		TenantDomain:      subscriptionInternal.TenantDomain,
		SubscriptionUUID:  subscriptionInternal.SubscriptionUUID,
		APIUUID:           subscriptionInternal.APIUUID,
		ApplicationUUID:   subscriptionInternal.ApplicationUUID,
	}
	if sub.TenantDomain == "" {
		sub.TenantDomain = config.GetControlPlaneConnectedTenantDomain()
	}
	return sub
}

func marshalApplication(appInternal *types.Application) Application {
	app := Application{
		UUID:         appInternal.UUID,
		ID:           appInternal.ID,
		Name:         appInternal.Name,
		SubName:      appInternal.SubName,
		Policy:       appInternal.Policy,
		TokenType:    appInternal.TokenType,
		Attributes:   appInternal.Attributes,
		TenantID:     appInternal.TenantID,
		TenantDomain: appInternal.TenantDomain,
		TimeStamp:    appInternal.TimeStamp,
	}
	if app.TenantDomain == "" {
		app.TenantDomain = config.GetControlPlaneConnectedTenantDomain()
	}
	return app
}

func marshalKeyMapping(keyMappingInternal *types.ApplicationKeyMapping) ApplicationKeyMapping {
	return ApplicationKeyMapping{
		ConsumerKey:     keyMappingInternal.ConsumerKey,
		KeyType:         keyMappingInternal.KeyType,
		KeyManager:      keyMappingInternal.KeyManager,
		ApplicationID:   keyMappingInternal.ApplicationID,
		ApplicationUUID: keyMappingInternal.ApplicationUUID,
		TenantID:        keyMappingInternal.TenantID,
		TenantDomain:    keyMappingInternal.TenantDomain,
		TimeStamp:       keyMappingInternal.TimeStamp,
	}
}

// GetApplicationKeyMappingReference returns unique reference for each key Mapping event.
// It is the combination of consumerKey:keyManager
func GetApplicationKeyMappingReference(keyMapping *types.ApplicationKeyMapping) string {
	return keyMapping.ConsumerKey + ":" + keyMapping.KeyManager
}

// CheckIfAPIMetadataIsAlreadyAvailable returns true only if the API Metadata for the given API UUID
// is already available
func CheckIfAPIMetadataIsAlreadyAvailable(apiUUID, label string) bool {
	if _, labelAvailable := APIListMap[label]; labelAvailable {
		if _, apiAvailale := APIListMap[label][apiUUID]; apiAvailale {
			return true
		}
	}
	return false
}
