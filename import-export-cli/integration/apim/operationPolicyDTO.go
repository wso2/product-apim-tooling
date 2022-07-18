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

type APIPolicyFileData struct {
	Id                string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name              string            `json:"name,omitempty" yaml:"name,omitempty"`
	DisplayName       string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Description       string            `json:"description"`
	Category          string            `json:"category,omitempty" yaml:"category,omitempty"`
	ApplicableFlows   []string          `json:"applicableFlows" yaml:"applicableFlows"`
	SupportedGateways []string          `json:"supportedGateways" yaml:"supportedGateways"`
	SupportedApiTypes []string          `json:"supportedApiTypes" yaml:"supportedApiTypes"`
	PolicyAttributes  []PolicyAttribute `json:"policyAttributes,omitempty" yaml:"policyAttributes,omitempty"`
}

type PolicyAttribute struct {
	Name          string   `json:"name" yaml:"name"`
	DisplayName   string   `json:"displayName" yaml:"displayName"`
	Description   string   `json:"description,omitempty" yaml:"description,omitempty"`
	Type          string   `json:"type,omitempty" yaml:"type,omitempty"`
	AllowedValues []string `json:"allowedValues,omitempty" yaml:"allowedValues,omitempty"`
	Required      bool     `json:"required,omitempty" yaml:"required,omitempty"`
}

type APIPolicyRequest struct {
	PolicySpecFile              []byte `json:"policySpecFile"`
	SynapsePolicyDefinitionFile []byte `json:"synapsePolicyDefinitionFile,omitempty"`
	CCPolicyDefinitionFile      []byte `json:"ccPolicyDefinitionFile,omitempty"`
}

type APIPolicyFile struct {
	Type    string            `json:"type" yaml:"type"`
	Version string            `json:"version" yaml:"version"`
	Data    APIPolicyFileData `json:"data" yaml:"data"`
}

type PolicySpecData struct {
	Type              string            `json:"type,omitempty" yaml:"type,omitempty"`
	Version           string            `json:"version,omitempty" yaml:"version,omitempty"`
	Id                string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name              string            `json:"name,omitempty" yaml:"name,omitempty"`
	DisplayName       string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Category          string            `json:"category,omitempty" yaml:"category,omitempty"`
	Description       string            `json:"description,omitempty" yaml:"description,omitempty"`
	ApplicableFlows   []string          `json:"applicableFlows" yaml:"applicableFlows"`
	SupportedGateways []string          `json:"supportedGateways" yaml:"supportedGateways"`
	SupportedApiTypes []string          `json:"supportedApiTypes" yaml:"supportedApiTypes"`
	PolicyAttributes  []PolicyAttribute `json:"policyAttributes,omitempty" yaml:"policyAttributes,omitempty"`
}

type APIPoliciesList struct {
	Count int                 `json:"count"`
	List  []APIPolicyFileData `json:"list"`
}
