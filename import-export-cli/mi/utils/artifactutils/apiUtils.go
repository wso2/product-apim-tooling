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

package artifactutils

type IntegrationAPI struct {
	Name      string     `json:"name"`
	Url       string     `json:"url"`
	Version   string     `json:"version"`
	Stats     string     `json:"stats"`
	Tracing   string     `json:"tracing"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Methods []string `json:"methods"`
	Url     string   `json:"url"`
}

type IntegrationAPIList struct {
	Count int32                   `json:"count"`
	Apis  []IntegrationAPISummary `json:"list"`
}

type IntegrationAPISummary struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
