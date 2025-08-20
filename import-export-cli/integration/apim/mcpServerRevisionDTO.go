/*
*  Copyright (c) 2025 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 LLC. licenses this file to you under the Apache License,
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

// MCPServerRevision : MCP Server Revision DTO
type MCPServerRevision struct {
	DisplayName    string                        `json:"displayName"`
	ID             string                        `json:"id"`
	Description    string                        `json:"description"`
	DeploymentInfo []MCPServerRevisionDeployment `json:"deploymentInfo"`
}

// MCPServerRevisionDeployment : MCP Server Revision Deployment DTO
type MCPServerRevisionDeployment struct {
	RevisionUUID       string `json:"revisionUuid"`
	Name               string `json:"name"`
	VHost              string `json:"vhost"`
	DisplayOnDevportal bool   `json:"displayOnDevportal"`
}

// MCPServerRevisionList : MCP Server Revisions List DTO
type MCPServerRevisionList struct {
	Count string              `json:"count"`
	List  []MCPServerRevision `json:"list"`
}
