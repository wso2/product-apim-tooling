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

package apim

// APIProduct : API Product DTO
type APIProduct struct {
	ID                           string               `json:"id"`
	Name                         string               `json:"name"`
	Description                  string               `json:"description"`
	Context                      string               `json:"context"`
	Provider                     string               `json:"provider"`
	HasThumbnail                 bool                 `json:"hasThumbnail"`
	State                        string               `json:"state"`
	EnableSchemaValidation       bool                 `json:"enableSchemaValidation"`
	ResponseCacheEnabled         bool                 `json:"responseCachingEnabled"`
	CacheTimeout                 int32                `json:"cacheTimeout"`
	Visibility                   string               `json:"visibility"`
	VisibleRoles                 []string             `json:"visibleRoles"`
	VisibleTenants               []string             `json:"visibleTenants"`
	AccessControl                string               `json:"accessControl"`
	AccessControlRoles           []string             `json:"accessControlRoles"`
	GatewayEnvironments          []string             `json:"gatewayEnvironments"`
	APIType                      string               `json:"apiType"`
	Transport                    []string             `json:"transport"`
	Tags                         []string             `json:"tags"`
	Policies                     []string             `json:"policies"`
	APIThrottlingPolicy          string               `json:"apiThrottlingPolicy"`
	AuthorizationHeader          string               `json:"authorizationHeader"`
	SecurityScheme               []string             `json:"securityScheme"`
	SubscriptionAvailability     string               `json:"subscriptionAvailability"`
	SubscriptionAvailableTenants []string             `json:"subscriptionAvailableTenants"`
	AdditionalProperties         map[string]string    `json:"additionalProperties"`
	Monetization                 APIMonetization      `json:"monetization,omitempty"`
	BusinessInformation          BusinessInfo         `json:"businessInformation,omitempty"`
	CorsConfiguration            APICorsConfiguration `json:"corsConfiguration,omitempty"`
	CreatedTime                  string               `json:"createdTime"`
	LastUpdatedTime              string               `json:"lastUpdatedTime"`
	APIs                         []interface{}        `json:"apis"`
	Scopes                       []interface{}        `json:"scopes"`
	Categories                   []interface{}        `json:"categories"`
}
