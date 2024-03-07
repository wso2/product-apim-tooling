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
package managementserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddApplication(t *testing.T) {
	testApp := Application{
		UUID:         "uelu3y37a822828ye2hd2eo27yo8q822d2o7dwdhccbcwuw",
		Name:         "Test App",
		Owner:        "John Doe",
		Organization: "Org1",
		Attributes:   map[string]string{"k": "v"},
		TimeStamp:    123456789,
	}
	AddApplication(testApp)
	if _, ok := applicationMap[testApp.UUID]; !ok {
		t.Errorf("Application not added to the map")
	}
}

func TestAddSubscription(t *testing.T) {
	testSub := Subscription{
		UUID:         "c8d8fb750ece-0a5b1039-9836-4b05-baa8-e06c",
		SubStatus:    "active",
		Organization: "Org1",
		SubscribedAPI: &SubscribedAPI{
			Name:    "Test API",
			Version: "v1.0.0",
		},
		TimeStamp: 123456789,
	}
	AddSubscription(testSub)
	if _, ok := subscriptionMap[testSub.UUID]; !ok {
		t.Errorf("Subscription not added to the map")
	}
}

func TestAddApplicationMapping(t *testing.T) {
	applicationMapping := ApplicationMapping{
		UUID:            "mapping1",
		ApplicationRef:  "app1",
		SubscriptionRef: "sub1",
		Organization:    "Org1",
	}
	AddApplicationMapping(applicationMapping)
	if _, ok := applicationMappingMap[applicationMapping.UUID]; !ok {
		t.Errorf("Application mapping not added to the map")
	}
}

func TestAddApplicationKeyMapping(t *testing.T) {
	td := []struct {
		applicationKeyMapping ApplicationKeyMapping
		expectedUniqueID      string
	}{
		{ApplicationKeyMapping{"app1", "scheme1", "appID1", "keyType1", "Env1", 12334556, "Org1"}, "272fd0a54d4a515ddf5afd2d65720abdefb9f494"},
		{ApplicationKeyMapping{"app2", "scheme2", "appID2", "keyType2", "Env2", 12334556, "Org2"}, "45b387e3701e9a0d63ecb5fa6645f19b7c8a3795"},
		{ApplicationKeyMapping{"app3", "scheme3", "appID3", "keyType3", "Env3", 12334556, "Org3"}, "f2fbae37f4359f2bd05c759557e6cdead142bded"},
	}
	for _, test := range td {
		AddApplicationKeyMapping(test.applicationKeyMapping)
		if _, ok := applicationKeyMappingMap[test.expectedUniqueID]; !ok {
			t.Error("Application mapping not added to the map")
		}
	}
}
func TestGetAllApplications(t *testing.T) {
	application1 := Application{UUID: "app1", Name: "Test App 1", Owner: "John Doe", Organization: "Org1", Attributes: map[string]string{"key1": "value1"}, TimeStamp: 123456789}
	application2 := Application{UUID: "app2", Name: "Test App 2", Owner: "Jane Smith", Organization: "Org2", Attributes: map[string]string{"key2": "value2"}, TimeStamp: 987654321}

	applicationMap = map[string]Application{
		"app1": application1,
		"app2": application2,
	}

	// Create mappings using application UUIDs as keys
	applicationKeyMapping1 := ApplicationKeyMapping{ApplicationUUID: "app1", SecurityScheme: "scheme1", ApplicationIdentifier: "identifier1", KeyType: "type1", EnvID: "env1", Timestamp: 123456789, Organization: "Org1"}
	applicationKeyMapping2 := ApplicationKeyMapping{ApplicationUUID: "app2", SecurityScheme: "scheme2", ApplicationIdentifier: "identifier2", KeyType: "type2", EnvID: "env2", Timestamp: 987654321, Organization: "Org2"}
	applicationKeyMappingMap = map[string]ApplicationKeyMapping{
		"app1": applicationKeyMapping1,
		"app2": applicationKeyMapping2,
	}

	applications := GetAllApplications()
	assert.Len(t, applications, 2)
	for _, app := range applications {
		expApp := applicationMap[app.UUID]
		assert.Equal(t, app.UUID, expApp.UUID)
		assert.Equal(t, app.Name, expApp.Name)
		assert.Equal(t, app.Owner, expApp.Owner)
		assert.Equal(t, app.Organization, expApp.Organization)
		assert.Equal(t, app.TimeStamp, int64(expApp.TimeStamp))
		assert.Len(t, app.SecuritySchemes, 1) // Assuming each application has only one associated security scheme
		assert.Equal(t, app.SecuritySchemes[0].SecurityScheme, applicationKeyMappingMap[app.UUID].SecurityScheme)
	}
}

func TestGetAllSubscriptions(t *testing.T) {
	subscription1 := Subscription{UUID: "sub1", SubStatus: "Active", Organization: "Org1"}
	subscription2 := Subscription{UUID: "sub2", SubStatus: "Inactive", Organization: "Org2"}
	subscriptionMap = map[string]Subscription{
		"sub1": subscription1,
		"sub2": subscription2,
	}
	subscriptions := GetAllSubscriptions()
	assert.Len(t, subscriptions, 2)
	for _, sub := range subscriptions {
		expSub := subscriptionMap[sub.UUID]
		assert.Equal(t, sub.UUID, expSub.UUID)
		assert.Equal(t, sub.SubStatus, expSub.SubStatus)
		assert.Equal(t, sub.Organization, expSub.Organization)
	}
}

func TestGetApplication(t *testing.T) {
	// Sample application
	application := Application{UUID: "app1", Name: "Test App", Owner: "John Doe", Organization: "Org1"}
	applicationMap = map[string]Application{
		"app1": application,
	}
	result := GetApplication("app1")
	assert.Equal(t, result, application)
}

func TestGetSubscription(t *testing.T) {
	subscription := Subscription{UUID: "sub1", SubStatus: "Active", Organization: "Org1"}
	subscriptionMap = map[string]Subscription{
		"sub1": subscription,
	}
	result := GetSubscription("sub1")
	assert.Equal(t, result, subscription)
}

func TestGetApplicationMapping(t *testing.T) {
	applicationMapping := ApplicationMapping{UUID: "map1", ApplicationRef: "app1", SubscriptionRef: "sub1", Organization: "Org1"}
	applicationMappingMap = map[string]ApplicationMapping{
		"map1": applicationMapping,
	}
	result := GetApplicationMapping("map1")
	assert.Equal(t, result, applicationMapping)
}

func TestGetApplicationKeyMapping(t *testing.T) {
	applicationKeyMapping := ApplicationKeyMapping{ApplicationUUID: "app1", KeyType: "OAuth", SecurityScheme: "Bearer", EnvID: "env1", ApplicationIdentifier: "app_identifier", Organization: "Org1"}
	applicationKeyMappingMap = map[string]ApplicationKeyMapping{
		"key1": applicationKeyMapping,
	}
	result := GetApplicationKeyMapping("key1")
	assert.Equal(t, result, applicationKeyMapping)
}

func TestDeleteApplication(t *testing.T) {
	applicationMap = map[string]Application{"app1": {UUID: "app1", Name: "Test App", Organization: "Org1"}}
	DeleteApplication("app1")
	assert.Empty(t, applicationMap)
}

func TestDeleteSubscription(t *testing.T) {
	subscriptionMap = map[string]Subscription{"sub1": {UUID: "sub1", Organization: "Org1"}}
	DeleteSubscription("sub1")
	assert.Empty(t, subscriptionMap)
}

func TestDeleteApplicationMapping(t *testing.T) {
	applicationMappingMap = map[string]ApplicationMapping{"map1": {UUID: "map1", Organization: "Org1"}}
	DeleteApplicationMapping("map1")
	assert.Empty(t, applicationMappingMap)
}

func TestDeleteApplicationKeyMapping(t *testing.T) {
	uuid := "mapping1"
	applicationKeyMappingMap = map[string]ApplicationKeyMapping{
		uuid: ApplicationKeyMapping{ApplicationUUID: "app1", SecurityScheme: "OAuth", KeyType: "APIKey", EnvID: "env1", ApplicationIdentifier: "app_identifier", Organization: "Org1"},
	}
	DeleteApplicationKeyMapping(uuid)
	_, exists := applicationKeyMappingMap[uuid]
	assert.False(t, exists)
}

func TestUpdateApplication(t *testing.T) {
	uuid := "app1"
	application := Application{UUID: "uuid", Name: "Test App", Owner: "John Doe", Organization: "Org1"}
	UpdateApplication(uuid, application)
	assert.Equal(t, applicationMap[uuid], application)
}

func TestUpdateSubscription(t *testing.T) {
	uuid := "sub1"
	subscription := Subscription{UUID: "uuid", SubStatus: "Active", Organization: "Org1", SubscribedAPI: &SubscribedAPI{Name: "Test API", Version: "v1"}}
	UpdateSubscription(uuid, subscription)
	assert.Equal(t, subscriptionMap[uuid], subscription)
}

func TestUpdateApplicationMapping(t *testing.T) {
	uuid := "mapping1"
	applicationMapping := ApplicationMapping{UUID: "uuid", ApplicationRef: "app1", SubscriptionRef: "sub1", Organization: "Org1"}
	UpdateApplicationMapping(uuid, applicationMapping)
	assert.Equal(t, applicationMappingMap[uuid], applicationMapping)
}

func TestUpdateApplicationKeyMapping(t *testing.T) {
	uuid := "key_mapping1"
	applicationKeyMapping := ApplicationKeyMapping{ApplicationUUID: "app1", SecurityScheme: "OAuth", KeyType: "APIKey", EnvID: "env1", ApplicationIdentifier: "app_identifier", Organization: "Org1"}
	UpdateApplicationKeyMapping(uuid, applicationKeyMapping)
	assert.Equal(t, applicationKeyMappingMap[uuid], applicationKeyMapping)
}

func TestGetApplicationKeyMappingByApplicationUUID(t *testing.T) {
	applicationKeyMapping1 := ApplicationKeyMapping{ApplicationUUID: "app1", KeyType: "OAuth", SecurityScheme: "Bearer", EnvID: "env1", ApplicationIdentifier: "app_identifier1", Organization: "Org1"}
	applicationKeyMapping2 := ApplicationKeyMapping{ApplicationUUID: "app2", KeyType: "APIKey", SecurityScheme: "Basic", EnvID: "env2", ApplicationIdentifier: "app_identifier2", Organization: "Org1"}
	applicationKeyMappingMap = map[string]ApplicationKeyMapping{"mapping1": applicationKeyMapping1, "mapping2": applicationKeyMapping2}
	result := GetApplicationKeyMappingByApplicationUUID("app1")
	assert.Equal(t, result, applicationKeyMapping1)
}

func TestGetApplicationKeyMappingByApplicationUUIDAndEnvID(t *testing.T) {
	applicationKeyMapping1 := ApplicationKeyMapping{ApplicationUUID: "app1", KeyType: "OAuth", SecurityScheme: "Bearer", EnvID: "env1", ApplicationIdentifier: "app_identifier1", Organization: "Org1"}
	applicationKeyMapping2 := ApplicationKeyMapping{ApplicationUUID: "app2", KeyType: "APIKey", SecurityScheme: "Basic", EnvID: "env2", ApplicationIdentifier: "app_identifier2", Organization: "Org1"}
	applicationKeyMappingMap = map[string]ApplicationKeyMapping{"mapping1": applicationKeyMapping1, "mapping2": applicationKeyMapping2}
	result := GetApplicationKeyMappingByApplicationUUIDAndEnvID("app2", "env2")
	assert.Equal(t, result, applicationKeyMapping2)
}

func TestGetApplicationKeyMappingByApplicationUUIDAndSecurityScheme(t *testing.T) {
	applicationKeyMapping1 := ApplicationKeyMapping{ApplicationUUID: "app1", KeyType: "OAuth", SecurityScheme: "Bearer", EnvID: "env1", ApplicationIdentifier: "app_identifier1", Organization: "Org1"}
	applicationKeyMapping2 := ApplicationKeyMapping{ApplicationUUID: "app2", KeyType: "APIKey", SecurityScheme: "Basic", EnvID: "env2", ApplicationIdentifier: "app_identifier2", Organization: "Org1"}
	applicationKeyMappingMap = map[string]ApplicationKeyMapping{"mapping1": applicationKeyMapping1, "mapping2": applicationKeyMapping2}
	result := GetApplicationKeyMappingByApplicationUUIDAndSecurityScheme("app2", "Basic")
	assert.Equal(t, result, applicationKeyMapping2)
}

func TestGetApplicationKeyMappingByApplicationUUIDAndSecuritySchemeAndEnvID(t *testing.T) {
	applicationKeyMapping1 := ApplicationKeyMapping{ApplicationUUID: "app1", KeyType: "OAuth", SecurityScheme: "Bearer", EnvID: "env1", ApplicationIdentifier: "app_identifier1", Organization: "Org1"}
	applicationKeyMapping2 := ApplicationKeyMapping{ApplicationUUID: "app2", KeyType: "APIKey", SecurityScheme: "Basic", EnvID: "env2", ApplicationIdentifier: "app_identifier2", Organization: "Org1"}
	applicationKeyMappingMap = map[string]ApplicationKeyMapping{"mapping1": applicationKeyMapping1, "mapping2": applicationKeyMapping2}
	result := GetApplicationKeyMappingByApplicationUUIDAndSecuritySchemeAndEnvID("app2", "Basic", "env2")
	assert.Equal(t, result, applicationKeyMapping2)
}

func TestGetApplicationMappingByApplicationUUID(t *testing.T) {
	applicationMapping1 := ApplicationMapping{UUID: "mapping1", ApplicationRef: "app1", SubscriptionRef: "sub1", Organization: "Org1"}
	applicationMapping2 := ApplicationMapping{UUID: "mapping2", ApplicationRef: "app2", SubscriptionRef: "sub2", Organization: "Org2"}
	applicationMappingMap = map[string]ApplicationMapping{"mapping1": applicationMapping1, "mapping2": applicationMapping2}
	result := GetApplicationMappingByApplicationUUID("app1")
	assert.Equal(t, result, applicationMapping1)
}

func TestGetApplicationMappingByApplicationUUIDAndSubscriptionUUID(t *testing.T) {
	applicationMapping1 := ApplicationMapping{UUID: "mapping1", ApplicationRef: "app1", SubscriptionRef: "sub1", Organization: "Org1"}
	applicationMapping2 := ApplicationMapping{UUID: "mapping2", ApplicationRef: "app2", SubscriptionRef: "sub2", Organization: "Org2"}
	applicationMappingMap = map[string]ApplicationMapping{"mapping1": applicationMapping1, "mapping2": applicationMapping2}
	result := GetApplicationMappingByApplicationUUIDAndSubscriptionUUID("app1", "sub1")
	assert.Equal(t, result, applicationMapping1)
}

func TestDeleteAllApplications(t *testing.T) {
	DeleteAllApplications()
	assert.Empty(t, applicationMap)
}

func TestDeleteAllSubscriptions(t *testing.T) {
	DeleteAllSubscriptions()
	assert.Empty(t, subscriptionMap)
}

func TestDeleteAllApplicationMappings(t *testing.T) {
	DeleteAllApplicationMappings()
	assert.Empty(t, applicationMappingMap)
}

func TestDeleteAllApplicationKeyMappings(t *testing.T) {
	DeleteAllApplicationKeyMappings()
	assert.Empty(t, applicationKeyMappingMap)
}

func TestAddAllSubscriptions(t *testing.T) {
	subscriptionMapTemp := map[string]Subscription{
		"sub1": {UUID: "sub1", SubStatus: "Active", Organization: "Org1"},
		"sub2": {UUID: "sub2", SubStatus: "Inactive", Organization: "Org2"},
	}
	AddAllSubscriptions(subscriptionMapTemp)
	assert.Equal(t, subscriptionMapTemp, subscriptionMap)
}

func TestAddAllApplications(t *testing.T) {
	applicationMapTemp := map[string]Application{
		"app1": {UUID: "app1", Name: "Test App 1", Owner: "John Doe", Organization: "Org1", Attributes: map[string]string{"key1": "value1"}, TimeStamp: 123456789},
		"app2": {UUID: "app2", Name: "Test App 2", Owner: "Jane Smith", Organization: "Org2", Attributes: map[string]string{"key2": "value2"}, TimeStamp: 987654321},
	}
	AddAllApplications(applicationMapTemp)
	assert.Equal(t, applicationMapTemp, applicationMap)
}

func TestAddAllApplicationMappings(t *testing.T) {
	applicationMappingMapTemp := map[string]ApplicationMapping{
		"mapping1": {UUID: "mapping1", ApplicationRef: "app1", SubscriptionRef: "sub1", Organization: "Org1"},
		"mapping2": {UUID: "mapping2", ApplicationRef: "app2", SubscriptionRef: "sub2", Organization: "Org2"},
	}
	AddAllApplicationMappings(applicationMappingMapTemp)
	assert.Equal(t, applicationMappingMapTemp, applicationMappingMap)
}

func TestAddAllApplicationKeyMappings(t *testing.T) {
	applicationKeyMappingMapTemp := map[string]ApplicationKeyMapping{
		"keyMapping1": {ApplicationUUID: "app1", KeyType: "OAuth", SecurityScheme: "Bearer", EnvID: "env1", ApplicationIdentifier: "app_identifier", Organization: "Org1"},
		"keyMapping2": {ApplicationUUID: "app2", KeyType: "APIKey", SecurityScheme: "APIKey", EnvID: "env2", ApplicationIdentifier: "app_identifier", Organization: "Org2"},
	}
	AddAllApplicationKeyMappings(applicationKeyMappingMapTemp)
	assert.Equal(t, applicationKeyMappingMapTemp, applicationKeyMappingMap)
}

func TestDeleteAllSubscriptionsByApplicationsUUID(t *testing.T) {
	uuid := "Org1"
	DeleteAllSubscriptionsByApplicationsUUID(uuid)
	for _, sub := range subscriptionMap {
		assert.NotEqual(t, uuid, sub.Organization)
	}
}

func TestDeleteAllApplicationMappingsByApplicationsUUID(t *testing.T) {
	uuid := "mapping1"
	DeleteAllApplicationMappingsByApplicationsUUID(uuid)
	for _, appMapping := range applicationMappingMap {
		assert.NotEqual(t, uuid, appMapping.ApplicationRef)
	}
}
