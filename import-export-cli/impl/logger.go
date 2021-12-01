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

package impl

import (
	"encoding/json"
	"errors"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
)

func GetPerAPILoggingList(accessToken, apiListEndpoint, query, limit string) (count int32, apis []utils.API, err error) {
	queryParamAdded := false
	getQueryParamConnector := func() (connector string) {
		if queryParamAdded {
			return "&"
		} else {
			return ""
		}
	}

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryParamSring string
	if query != "" {
		queryParamSring = "query=" + query
		queryParamAdded = true
	}
	if limit != "" {
		queryParamSring += getQueryParamConnector() + "limit=" + limit
	}
	utils.Logln(utils.LogPrefixInfo+"URL:", apiListEndpoint+"?"+queryParamSring)
	resp, err := utils.InvokeGETRequestWithQueryParamsString(apiListEndpoint, queryParamSring, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiListEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiListResponse := &utils.APIListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return apiListResponse.Count, apiListResponse.List, nil
	} else {
		return 0, nil, errors.New(string(resp.Body()))
	}
}
