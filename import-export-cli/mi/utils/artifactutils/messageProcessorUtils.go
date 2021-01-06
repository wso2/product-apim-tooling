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

type MessageProcessorList struct {
	Count             int32              `json:"count"`
	MessageProcessors []MessageProcessor `json:"list"`
}

type MessageProcessor struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type MessageProcessorData struct {
	Name       string            `json:"name"`
	FileName   string            `json:"fileName"`
	Type       string            `json:"type"`
	Store      string            `json:"messageStore"`
	Container  string            `json:"artifactContainer"`
	Parameters map[string]string `json:"parameters"`
	Status     string            `json:"status"`
}
