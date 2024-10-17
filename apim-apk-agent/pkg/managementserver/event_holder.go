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
	eventHub "github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/eventhub/types"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/loggers"
	"github.com/wso2/product-apim-tooling/apim-apk-agent/pkg/utils"
)

var (
	applicationMap           map[string]Application
	subscriptionMap          map[string]Subscription
	applicationMappingMap    map[string]ApplicationMapping
	applicationKeyMappingMap map[string]ApplicationKeyMapping
	rateLimitPolicyMap       map[string]eventHub.RateLimitPolicy
	aiProviderMap            map[string]eventHub.AIProvider
	subscriptionPolicyMap    map[string]eventHub.SubscriptionPolicy
)

func init() {
	applicationMap = make(map[string]Application)
	subscriptionMap = make(map[string]Subscription)
	applicationMappingMap = make(map[string]ApplicationMapping)
	applicationKeyMappingMap = make(map[string]ApplicationKeyMapping)
	rateLimitPolicyMap = make(map[string]eventHub.RateLimitPolicy)
	aiProviderMap = make(map[string]eventHub.AIProvider)
	subscriptionPolicyMap = make(map[string]eventHub.SubscriptionPolicy)
}

// AddAIProvider adds an AI provider to the aiProviderMap
func AddAIProvider(aiProvider eventHub.AIProvider) {
	aiProviderMap[aiProvider.ID] = aiProvider
}

// GetAIProvider returns an AI provider from the aiProviderMap
func GetAIProvider(id string) eventHub.AIProvider {
	return aiProviderMap[id]
}

// DeleteAIProvider deletes an AI provider from the aiProviderMap
func DeleteAIProvider(id string) {
	delete(aiProviderMap, id)
}

// GetAllAIProviders returns all the AI providers in the aiProviderMap
func GetAllAIProviders() []eventHub.AIProvider {
	var aiProviders []eventHub.AIProvider
	for _, aiProvider := range aiProviderMap {
		aiProviders = append(aiProviders, aiProvider)
	}
	return aiProviders
}

// AddRateLimitPolicy adds a rate limit policy to the rateLimitPolicyMap
func AddRateLimitPolicy(rateLimitPolicy eventHub.RateLimitPolicy) {
	rateLimitPolicyMap[rateLimitPolicy.Name+rateLimitPolicy.TenantDomain] = rateLimitPolicy
}

// AddSubscriptionPolicy adds a rate limit policy to the subscriptionPolicyMap
func AddSubscriptionPolicy(rateLimitPolicy eventHub.SubscriptionPolicy) {
	subscriptionPolicyMap[rateLimitPolicy.Name+rateLimitPolicy.TenantDomain] = rateLimitPolicy
}

// GetSubscriptionPolicies return the subscription policy map
func GetSubscriptionPolicies() map[string]eventHub.SubscriptionPolicy {
	return subscriptionPolicyMap
}

// GetRateLimitPolicy returns a rate limit policy from the rateLimitPolicyMap
func GetRateLimitPolicy(name string, tenantDomain string) eventHub.RateLimitPolicy {
	return rateLimitPolicyMap[name+tenantDomain]
}

// GetAllRateLimitPolicies returns all the rate limit policies in the rateLimitPolicyMap
func GetAllRateLimitPolicies() []eventHub.RateLimitPolicy {
	var rateLimitPolicies []eventHub.RateLimitPolicy
	for _, rateLimitPolicy := range rateLimitPolicyMap {
		rateLimitPolicies = append(rateLimitPolicies, rateLimitPolicy)
	}
	return rateLimitPolicies
}

// DeleteRateLimitPolicy deletes a rate limit policy from the rateLimitPolicyMap
func DeleteRateLimitPolicy(name string, tenantDomain string) {
	delete(rateLimitPolicyMap, name+tenantDomain)
}

// DeleteSubscriptionPolicy deletes a subscription policy from the subscriptionPolicyMap
func DeleteSubscriptionPolicy(name string, tenantDomain string) {
	delete(subscriptionPolicyMap, name+tenantDomain)
}

// UpdateRateLimitPolicy updates a rate limit policy in the rateLimitPolicyMap
func UpdateRateLimitPolicy(name string, tenantDomain string, rateLimitPolicy eventHub.RateLimitPolicy) {
	rateLimitPolicyMap[name+tenantDomain] = rateLimitPolicy
}

// AddApplication adds an application to the applicationMap
func AddApplication(application Application) {
	applicationMap[application.UUID] = application
}

// AddSubscription adds a subscription to the subscriptionMap
func AddSubscription(subscription Subscription) {
	subscriptionMap[subscription.UUID] = subscription
}

// AddApplicationMapping adds an application mapping to the applicationMappingMap
func AddApplicationMapping(applicationMapping ApplicationMapping) {
	applicationMappingMap[applicationMapping.UUID] = applicationMapping
}

// AddApplicationKeyMapping adds an application key mapping to the applicationKeyMappingMap
func AddApplicationKeyMapping(applicationKeyMapping ApplicationKeyMapping) {
	uuid := utils.GetUniqueIDOfApplicationKeyMapping(applicationKeyMapping.ApplicationUUID, applicationKeyMapping.KeyType, applicationKeyMapping.SecurityScheme, applicationKeyMapping.EnvID, applicationKeyMapping.Organization)
	loggers.LoggerMgtServer.Infof("Adding application key mapping with uuid: %v", uuid)
	applicationKeyMappingMap[uuid] = applicationKeyMapping
}

// GetAllApplications returns all the applications in the applicationMap
func GetAllApplications() []ResolvedApplication {
	var applications []ResolvedApplication
	for _, application := range applicationMap {
		resolvedApplication := marshalApplication(application)
		applications = append(applications, resolvedApplication)
	}
	return applications
}
func marshalApplication(application Application) ResolvedApplication {
	resolvedApplication := ResolvedApplication{UUID: application.UUID, Name: application.Name, Owner: application.Owner, Organization: application.Organization, Attributes: application.Attributes, TimeStamp: application.TimeStamp, SecuritySchemes: make([]SecurityScheme, 0)}
	for _, applicationKeyMapping := range applicationKeyMappingMap {
		if applicationKeyMapping.ApplicationUUID == application.UUID {
			securityScheme := SecurityScheme{SecurityScheme: applicationKeyMapping.SecurityScheme, KeyType: applicationKeyMapping.KeyType, EnvID: applicationKeyMapping.EnvID, ApplicationIdentifier: applicationKeyMapping.ApplicationIdentifier}
			resolvedApplication.SecuritySchemes = append(resolvedApplication.SecuritySchemes, securityScheme)
		}
	}
	return resolvedApplication
}

// GetAllSubscriptions returns all the subscriptions in the subscriptionMap
func GetAllSubscriptions() []Subscription {
	var subscriptions []Subscription
	for _, subscription := range subscriptionMap {
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions
}

// GetAllApplicationMappings returns all the application mappings in the applicationMappingMap
func GetAllApplicationMappings() []ApplicationMapping {
	var applicationMappings []ApplicationMapping
	for _, applicationMapping := range applicationMappingMap {
		applicationMappings = append(applicationMappings, applicationMapping)
	}
	return applicationMappings
}

// GetApplication returns an application from the applicationMap
func GetApplication(uuid string) Application {
	return applicationMap[uuid]
}

// GetSubscription returns a subscription from the subscriptionMap
func GetSubscription(uuid string) Subscription {
	return subscriptionMap[uuid]
}

// GetApplicationMapping returns an application mapping from the applicationMappingMap
func GetApplicationMapping(uuid string) ApplicationMapping {
	return applicationMappingMap[uuid]
}

// GetApplicationKeyMapping returns an application key mapping from the applicationKeyMappingMap
func GetApplicationKeyMapping(uuid string) ApplicationKeyMapping {
	return applicationKeyMappingMap[uuid]
}

// DeleteApplication deletes an application from the applicationMap
func DeleteApplication(uuid string) {
	delete(applicationMap, uuid)
}

// DeleteSubscription deletes a subscription from the subscriptionMap
func DeleteSubscription(uuid string) {
	delete(subscriptionMap, uuid)
}

// DeleteApplicationMapping deletes an application mapping from the applicationMappingMap
func DeleteApplicationMapping(uuid string) {
	delete(applicationMappingMap, uuid)
}

// DeleteApplicationKeyMapping deletes an application key mapping from the applicationKeyMappingMap
func DeleteApplicationKeyMapping(uuid string) {
	loggers.LoggerMgtServer.Infof("Deleting application key mapping with uuid: %v", uuid)
	delete(applicationKeyMappingMap, uuid)
}

// UpdateApplication updates an application in the applicationMap
func UpdateApplication(uuid string, application Application) {
	applicationMap[uuid] = application
}

// UpdateSubscription updates a subscription in the subscriptionMap
func UpdateSubscription(uuid string, subscription Subscription) {
	subscriptionMap[uuid] = subscription
}

// UpdateApplicationMapping updates an application mapping in the applicationMappingMap
func UpdateApplicationMapping(uuid string, applicationMapping ApplicationMapping) {
	applicationMappingMap[uuid] = applicationMapping
}

// UpdateApplicationKeyMapping updates an application key mapping in the applicationKeyMappingMap
func UpdateApplicationKeyMapping(uuid string, applicationKeyMapping ApplicationKeyMapping) {
	applicationKeyMappingMap[uuid] = applicationKeyMapping
}

// GetApplicationKeyMappingByApplicationUUID returns an application key mapping from the applicationKeyMappingMap
func GetApplicationKeyMappingByApplicationUUID(uuid string) ApplicationKeyMapping {
	for _, applicationKeyMapping := range applicationKeyMappingMap {
		if applicationKeyMapping.ApplicationUUID == uuid {
			return applicationKeyMapping
		}
	}
	return ApplicationKeyMapping{}
}

// GetApplicationKeyMappingByApplicationUUIDAndEnvID returns an application key mapping from the applicationKeyMappingMap
func GetApplicationKeyMappingByApplicationUUIDAndEnvID(uuid string, envID string) ApplicationKeyMapping {
	for _, applicationKeyMapping := range applicationKeyMappingMap {
		if applicationKeyMapping.ApplicationUUID == uuid && applicationKeyMapping.EnvID == envID {
			return applicationKeyMapping
		}
	}
	return ApplicationKeyMapping{}
}

// GetApplicationKeyMappingByApplicationUUIDAndSecurityScheme returns an application key mapping from the applicationKeyMappingMap
func GetApplicationKeyMappingByApplicationUUIDAndSecurityScheme(uuid string, securityScheme string) ApplicationKeyMapping {
	for _, applicationKeyMapping := range applicationKeyMappingMap {
		if applicationKeyMapping.ApplicationUUID == uuid && applicationKeyMapping.SecurityScheme == securityScheme {
			return applicationKeyMapping
		}
	}
	return ApplicationKeyMapping{}
}

// GetApplicationKeyMappingByApplicationUUIDAndSecuritySchemeAndEnvID returns an application key mapping from the applicationKeyMappingMap
func GetApplicationKeyMappingByApplicationUUIDAndSecuritySchemeAndEnvID(uuid string, securityScheme string, envID string) ApplicationKeyMapping {
	for _, applicationKeyMapping := range applicationKeyMappingMap {
		if applicationKeyMapping.ApplicationUUID == uuid && applicationKeyMapping.SecurityScheme == securityScheme && applicationKeyMapping.EnvID == envID {
			return applicationKeyMapping
		}
	}
	return ApplicationKeyMapping{}
}

// GetApplicationMappingByApplicationUUID returns an application mapping from the applicationMappingMap
func GetApplicationMappingByApplicationUUID(uuid string) ApplicationMapping {
	for _, applicationMapping := range applicationMappingMap {
		if applicationMapping.ApplicationRef == uuid {
			return applicationMapping
		}
	}
	return ApplicationMapping{}
}

// GetApplicationMappingByApplicationUUIDAndSubscriptionUUID returns an application mapping from the applicationMappingMap
func GetApplicationMappingByApplicationUUIDAndSubscriptionUUID(uuid string, subscriptionUUID string) ApplicationMapping {
	for _, applicationMapping := range applicationMappingMap {
		if applicationMapping.ApplicationRef == uuid && applicationMapping.SubscriptionRef == subscriptionUUID {
			return applicationMapping
		}
	}
	return ApplicationMapping{}
}

// DeleteAllApplications deletes all the applications in the applicationMap
func DeleteAllApplications() {
	applicationMap = make(map[string]Application)
}

// DeleteAllSubscriptions deletes all the subscriptions in the subscriptionMap
func DeleteAllSubscriptions() {
	subscriptionMap = make(map[string]Subscription)
}

// DeleteAllApplicationMappings deletes all the application mappings in the applicationMappingMap
func DeleteAllApplicationMappings() {
	applicationMappingMap = make(map[string]ApplicationMapping)
}

// DeleteAllApplicationKeyMappings deletes all the application key mappings in the applicationKeyMappingMap
func DeleteAllApplicationKeyMappings() {
	applicationKeyMappingMap = make(map[string]ApplicationKeyMapping)
}

// AddAllSubscriptions adds all the subscriptions in the subscriptionMap
func AddAllSubscriptions(subscriptionMapTemp map[string]Subscription) {
	subscriptionMap = subscriptionMapTemp
}

// AddAllApplications adds all the applications in the applicationMap
func AddAllApplications(applicationMapTemp map[string]Application) {
	applicationMap = applicationMapTemp
}

// AddAllApplicationMappings adds all the application mappings in the applicationMappingMap
func AddAllApplicationMappings(applicationMappingMapTemp map[string]ApplicationMapping) {
	applicationMappingMap = applicationMappingMapTemp
}

// AddAllApplicationKeyMappings adds all the application key mappings in the applicationKeyMappingMap
func AddAllApplicationKeyMappings(applicationKeyMappingMapTemp map[string]ApplicationKeyMapping) {
	applicationKeyMappingMap = applicationKeyMappingMapTemp
}

// DeleteAllSubscriptionsByApplicationsUUID deletes all the subscriptions in the subscriptionMap
func DeleteAllSubscriptionsByApplicationsUUID(uuid string) {
	for _, subscription := range subscriptionMap {
		if subscription.Organization == uuid {
			delete(subscriptionMap, subscription.UUID)
		}
	}
}

// DeleteAllApplicationMappingsByApplicationsUUID deletes all the application mappings in the applicationMappingMap
func DeleteAllApplicationMappingsByApplicationsUUID(uuid string) {
	for _, applicationMapping := range applicationMappingMap {
		if applicationMapping.UUID == uuid {
			delete(applicationMappingMap, applicationMapping.UUID)
		}
	}
}
