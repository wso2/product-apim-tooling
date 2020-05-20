/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package v2

// APIProductDefinition represents an API Product artifact in APIM
type APIProductDefinition struct {
	ID                                 ProductID            `json:"id,omitempty" yaml:"id,omitempty"`
	UUID                               string               `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Description                        string               `json:"description,omitempty" yaml:"description,omitempty"`
	Type                               string               `json:"type,omitempty" yaml:"type,omitempty"`
	Context                            string               `json:"context" yaml:"context"`
	ContextTemplate                    string               `json:"contextTemplate,omitempty" yaml:"contextTemplate,omitempty"`
	Tags                               []string             `json:"tags" yaml:"tags,omitempty"`
	Documents                          []interface{}        `json:"documents,omitempty" yaml:"documents,omitempty"`
	LastUpdated                        string               `json:"lastUpdated,omitempty" yaml:"lastUpdated,omitempty"`
	AvailableTiers                     []AvailableTiers     `json:"availableTiers,omitempty" yaml:"availableTiers,omitempty"`
	AvailableSubscriptionLevelPolicies []interface{}        `json:"availableSubscriptionLevelPolicies,omitempty" yaml:"availableSubscriptionLevelPolicies,omitempty"`
	ProductResources                   []APIProductResource `json:"productResources" yaml:"productResources,omitempty"`
	State                              string               `json:"state,omitempty" yaml:"state,omitempty"`
	TechnicalOwner                     string               `json:"technicalOwner,omitempty" yaml:"technicalOwner,omitempty"`
	TechnicalOwnerEmail                string               `json:"technicalOwnerEmail,omitempty" yaml:"technicalOwnerEmail,omitempty"`
	BusinessOwner                      string               `json:"businessOwner,omitempty" yaml:"businessOwner,omitempty"`
	BusinessOwnerEmail                 string               `json:"businessOwnerEmail,omitempty" yaml:"businessOwnerEmail,omitempty"`
	Visibility                         string               `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	Transports                         string               `json:"transports,omitempty" yaml:"transports,omitempty"`
	CorsConfiguration                  *CorsConfiguration   `json:"corsConfiguration,omitempty" yaml:"corsConfiguration,omitempty"`
	ResponseCache                      string               `json:"responseCache,omitempty" yaml:"responseCache,omitempty"`
	CacheTimeout                       int                  `json:"cacheTimeout,omitempty" yaml:"cacheTimeout,omitempty"`
	AuthorizationHeader                string               `json:"authorizationHeader,omitempty" yaml:"authorizationHeader,omitempty"`
	Scopes                             []interface{}        `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	Environments                       []string             `json:"environments,omitempty" yaml:"environments,omitempty"`
	CreatedTime                        string               `json:"createdTime,omitempty" yaml:"createdTime,omitempty"`
	AdditionalProperties               map[string]string    `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	APISecurity                        string               `json:"apiSecurity,omitempty" yaml:"apiSecurity,omitempty"`
	AccessControl                      string               `json:"accessControl,omitempty" yaml:"accessControl,omitempty"`
	Rating                             string               `json:"rating,omitempty" yaml:"rating,omitempty"`
}
type ProductID struct {
	ProviderName   string `json:"providerName" yaml:"providerName"`
	APIProductName string `json:"apiProductName" yaml:"apiProductName"`
	Version        string `json:"version" yaml:"version"`
}
type APIProductResource struct {
	APIProductName       string                 `json:"apiName,omitempty" yaml:"apiName,omitempty"`
	APIProductId         string                 `json:"apiId,omitempty" yaml:"apiId,omitempty"`
	APIIdentifier        ID                     `json:"apiIdentifier,omitempty" yaml:"apiIdentifier,omitempty"`
	APIProductIdentifier ID                     `json:"productIdentifier,omitempty" yaml:"productIdentifier,omitempty"`
	URITemplate          URITemplates           `json:"uriTemplate,omitempty" yaml:"uriTemplate,omitempty"`
	EndpointConfig       string                 `json:"endpointConfig,omitempty" yaml:"endpointConfig,omitempty"`
	EndpointSecurityMap  map[string]interface{} `json:"endpointSecurityMap,omitempty" yaml:"endpointSecurityMap,omitempty"`
	InSequenceName       string                 `json:"inSequenceName,omitempty" yaml:"inSequenceName,omitempty"`
	OutSequenceName      string                 `json:"outSequenceName,omitempty" yaml:"outSequenceName,omitempty"`
	FaultSequenceName    string                 `json:"faultSequenceName,omitempty" yaml:"faultSequenceName,omitempty"`
}
