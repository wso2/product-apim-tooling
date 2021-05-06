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

// ApplicationFile : Application File content
type ApplicationFile struct {
	Type    string          `json:"type" yaml:"type"`
	Version string          `json:"version" yaml:"version"`
	Data    ApplicationData `json:"data" yaml:"data"`
}

// ApplicationData : Application Data
type ApplicationData struct {
	ApplicationInfo Application `json:"applicationInfo" yaml:"applicationInfo"`
	SubscribedAPIs  interface{} `json:"subscribedAPIs" yaml:"subscribedAPIs"`
}

// Application : Application DTO
type Application struct {
	ApplicationID      string           `json:"applicationId" yaml:"applicationId"`
	Name               string           `json:"name" yaml:"name"`
	ThrottlingPolicy   string           `json:"throttlingPolicy" yaml:"throttlingPolicy"`
	Description        string           `json:"description" yaml:"description"`
	TokenType          string           `json:"tokenType" yaml:"tokenType"`
	Status             string           `json:"status" yaml:"status"`
	Groups             []string         `json:"groups" yaml:"groups"`
	SubscriptionCount  int              `json:"subscriptionCount" yaml:"subscriptionCount"`
	Keys               []ApplicationKey `json:"keys" yaml:"keys"`
	SubscriptionScopes []string         `json:"subscriptionScopes" yaml:"subscriptionScopes"`
	Owner              string           `json:"owner" yaml:"owner"`
	HashEnabled        bool             `json:"hashEnabled" yaml:"hashEnabled"`
}
