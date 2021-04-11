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
type APIProductDefinitionFile struct {
	Type        string                  `json:"type,omitempty" yaml:"type,omitempty"`
	ApimVersion string                  `json:"version,omitempty" yaml:"version,omitempty"`
	Data        APIProductDTODefinition `json:"data,omitempty" yaml:"data,omitempty"`
}

// APIProductDTODefinition represents an API Product artifact in APIM
type APIProductDTODefinition struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`
}
