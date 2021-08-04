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

// ApplicationKey : Application Key Details
type ApplicationKey struct {
	KeyMappingId         string           `json:"keyMappingId" yaml:"keyMappingId"`
	KeyManager           string           `json:"keyManager" yaml:"keyManager"`
	ConsumerKey          string           `json:"consumerKey" yaml:"consumerKey"`
	ConsumerSecret       string           `json:"consumerSecret" yaml:"consumerSecret"`
	Mode                 string           `json:"mode" yaml:"mode"`
	SupportedGrantTypes  []string         `json:"supportedGrantTypes" yaml:"supportedGrantTypes"`
	CallbackURL          string           `json:"callbackUrl" yaml:"callbackUrl"`
	KeyState             string           `json:"keyState" yaml:"keyState"`
	KeyType              string           `json:"keyType" yaml:"keyType"`
	GroupID              string           `json:"groupId" yaml:"groupId"`
	Token                ApplicationToken `json:"token" yaml:"token"`
	AdditionalProperties interface{}      `json:"additionalProperties" yaml:"additionalProperties"`
}

// ApplicationToken : Application Token Details
type ApplicationToken struct {
	AccessToken  string   `json:"accessToken"`
	TokenScopes  []string `json:"tokenScopes"`
	ValidityTime int64    `json:"validityTime"`
}

// ApplicationKeysList : Applications list
type ApplicationKeysList struct {
	Count int              `json:"count"`
	List  []ApplicationKey `json:"list"`
}
