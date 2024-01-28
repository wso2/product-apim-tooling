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
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	eventhubTypes "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/managementserver"
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

// KeyManager for struct keyManager
type KeyManager struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Issuer      string `json:"issuer"`
	Certificate string `json:"certificate"`
}

var (
	// KeyManagerMap is a map of key managers
	KeyManagerMap map[string]KeyManager
)

// MarshalKeyManagers is used to update the key managers during the startup where
// multiple key managers are pulled at once. And then it returns the KeyManagerMap.
func MarshalKeyManagers(keyManagersList *[]eventhubTypes.KeyManager) map[string]KeyManager {
	resourceMap := make(map[string]KeyManager)
	for _, keyManager := range *keyManagersList {
		resourceMap[keyManager.Name] = MarshalKeyManager(&keyManager)
	}
	KeyManagerMap = resourceMap
	return KeyManagerMap
}

// MarshalMultipleApplications is used to update the applicationList during the startup where
func MarshalMultipleApplications(appList *types.ApplicationList) {
	applicationMap := make(map[string]managementserver.Application)
	for _, application := range appList.List {
		applicationSub := MarshalApplication(&application)
		applicationMap[applicationSub.UUID] = applicationSub
	}
	managementserver.AddAllApplications(applicationMap)
}

// MarshalMultipleApplicationKeyMappings is used to update the application key mappings during the startup where
// multiple key mappings are pulled at once. And then it returns the ApplicationKeyMappingList.
func MarshalMultipleApplicationKeyMappings(keymappingList *types.ApplicationKeyMappingList) {
	resourceMap := make(map[string]managementserver.ApplicationKeyMapping)
	for _, keyMapping := range keymappingList.List {
		applicationKeyMappingReference := GetApplicationKeyMappingReference(&keyMapping)
		keyMappingSub := marshalKeyMapping(&keyMapping)
		resourceMap[applicationKeyMappingReference] = keyMappingSub
	}
	managementserver.AddAllApplicationKeyMappings(resourceMap)
}

// MarshalMultipleSubscriptions is used to update the subscriptions during the startup where
// multiple subscriptions are pulled at once. And then it returns the SubscriptionList.
func MarshalMultipleSubscriptions(subscriptionsList *types.SubscriptionList) {
	subscriptionMap := make(map[string]managementserver.Subscription)
	applicationMappingMap := make(map[string]managementserver.ApplicationMapping)
	for _, subscription := range subscriptionsList.List {
		subscriptionSub := MarshalSubscription(&subscription)
		subscriptionMap[subscriptionSub.UUID] = subscriptionSub
		applicationMappingMap[subscriptionSub.UUID] = managementserver.ApplicationMapping{
			UUID:            subscriptionSub.UUID,
			ApplicationRef:  subscriptionSub.SubscribedAPI.Name,
			SubscriptionRef: subscriptionSub.SubscribedAPI.Version,
			Organization:    subscriptionSub.Organization,
		}
	}
	managementserver.AddAllApplicationMappings(applicationMappingMap)
	managementserver.AddAllSubscriptions(subscriptionMap)

}

// MarshalSubscription is used to map to internal Subscription struct
func MarshalSubscription(subscriptionInternal *types.Subscription) managementserver.Subscription {
	sub := managementserver.Subscription{
		SubStatus:     subscriptionInternal.SubscriptionState,
		UUID:          subscriptionInternal.SubscriptionUUID,
		Organization:  subscriptionInternal.TenantDomain,
		SubscribedAPI: &managementserver.SubscribedAPI{Name: subscriptionInternal.APIUUID, Version: subscriptionInternal.APIUUID},
		TimeStamp:     subscriptionInternal.TimeStamp,
	}
	return sub
}

// MarshalApplication is used to map to internal Application struct
func MarshalApplication(appInternal *types.Application) managementserver.Application {
	app := managementserver.Application{
		UUID:         appInternal.UUID,
		Name:         appInternal.Name,
		Owner:        appInternal.SubName,
		Organization: appInternal.TenantDomain,
		Attributes:   appInternal.Attributes,
		TimeStamp:    appInternal.TimeStamp,
	}
	return app
}

func marshalKeyMapping(keyMappingInternal *types.ApplicationKeyMapping) managementserver.ApplicationKeyMapping {
	return managementserver.ApplicationKeyMapping{
		ApplicationUUID:       keyMappingInternal.ApplicationUUID,
		ApplicationIdentifier: keyMappingInternal.ConsumerKey,
		KeyType:               keyMappingInternal.KeyType,
		SecurityScheme:        "OAuth2",
		EnvID:                 "Default",
		Timestamp:             keyMappingInternal.TimeStamp,
	}
}

// MarshalKeyManager is used to map Internal key manager
func MarshalKeyManager(keyManagerInternal *types.KeyManager) KeyManager {
	return KeyManager{
		Name:    keyManagerInternal.Name,
		Enabled: keyManagerInternal.Enabled,
		// Issuer:      keyManagerInternal.Configuration.Issuer,
		// Certificate: keyManagerInternal.Configuration.Certificate,
	}
}

// GetApplicationKeyMappingReference returns unique reference for each key Mapping event.
// It is the combination of consumerKey:keyManager
func GetApplicationKeyMappingReference(keyMapping *types.ApplicationKeyMapping) string {
	return keyMapping.ConsumerKey + ":" + keyMapping.KeyManager
}

// CheckIfAPIMetadataIsAlreadyAvailable returns true only if the API Metadata for the given API UUID
// is already available
// func CheckIfAPIMetadataIsAlreadyAvailable(apiUUID, label string) bool {
// 	if _, labelAvailable := APIListMap[label]; labelAvailable {
// 		if _, apiAvailale := APIListMap[label][apiUUID]; apiAvailale {
// 			return true
// 		}
// 	}
// 	return false
// }
