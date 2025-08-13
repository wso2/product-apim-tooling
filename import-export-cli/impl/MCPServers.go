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

package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// GetMCPServerId returns the id of an MCP server
func GetMCPServerId(accessToken, environment, mcpServerName, mcpServerVersion, mcpServerProvider string) (string, error) {
	// Unified Search endpoint from the config file to search MCP Servers
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(environment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryVal string

	queryVal = "type:\"" + utils.DefaultMcpServerType + "\" name:\"" + mcpServerName + "\" version:\"" + mcpServerVersion + "\""

	if mcpServerProvider != "" {
		queryVal = queryVal + " provider:\"" + mcpServerProvider + "\""
	}

	resp, err := utils.InvokeGETRequestWithQueryParam("query", queryVal, unifiedSearchEndpoint, headers)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		mcpServerData := &utils.MCPServerSearch{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &mcpServerData)
		if mcpServerData.Count != 0 {
			mcpServerId := mcpServerData.List[0].ID
			fmt.Printf("MCP Server ID: %s\n", mcpServerId)
			return mcpServerId, err
		}
		if mcpServerProvider != "" {
			return "", errors.New("Requested MCP Server is not available. Name: " + mcpServerName +
				" Version: " + mcpServerVersion + " Provider: " + mcpServerProvider)
		}
		return "", errors.New("Requested MCP Server is not available. Name: " + mcpServerName +
			" Version: " + mcpServerVersion)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			return "", fmt.Errorf("Authorization failed while searching MCP Server: " + mcpServerName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching MCP Servers. Status: " + resp.Status())
	}
}

// GetMCPServerList returns a list of MCP servers from the publisher endpoint
func GetMCPServerList(accessToken, publisherEndpoint, query, limit string) (count int32, servers []utils.MCPServer, err error) {
	queryParamAdded := false
	getQueryParamConnector := func() string {
		if queryParamAdded {
			return "&"
		}
		return ""
	}

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryParamString string
	if query != "" {
		queryParamString = "query=" + query
		queryParamAdded = true
	}
	if limit != "" {
		queryParamString += getQueryParamConnector() + "limit=" + limit
	}
	utils.Logln(utils.LogPrefixInfo+"URL:", publisherEndpoint+"?"+queryParamString)
	resp, err := utils.InvokeGETRequestWithQueryParamsString(publisherEndpoint, queryParamString, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+publisherEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		mcpServerListResponse := &utils.MCPServerListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &mcpServerListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return mcpServerListResponse.Count, mcpServerListResponse.List, nil
	} else {
		return 0, nil, errors.New(string(resp.Body()))
	}
}

// GetMCPServerRevisionsList Get the list of Revisions available for the given MCP Server
// @param accessToken 			: Access Token for the environment
// @param revisionListEndpoint 	: Revision List endpoint
// @return count (no. of revisions)
// @return array of revision objects
// @return error
func GetMCPServerRevisionsList(accessToken, revisionListEndpoint string) (count int32, revisions []utils.Revisions, err error) {

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	utils.Logln(utils.LogPrefixInfo+"URL:", revisionListEndpoint)
	resp, err := utils.InvokeGETRequest(revisionListEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+revisionListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		revisionListResponse := &utils.RevisionListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &revisionListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return revisionListResponse.Count, revisionListResponse.List, nil
	} else {
		return 0, nil, errors.New(string(resp.Body()))
	}
}
