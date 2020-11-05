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
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// DeleteAPI
// @param accessToken : Access Token for the resource
// @param environment : Environment where API should be deleted
// @param deleteAPIName : Name of the API to delete
// @param deleteAPIVersion : Version of the API to delete
// @param deleteAPIProvider : Provider of API
// @return response Response in the form of *resty.Response
func DeleteAPI(accessToken, environment, deleteAPIName, deleteAPIVersion, deleteAPIProvider string) (*resty.Response, error) {
	deleteEndpoint := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)
	apiId, err := GetAPIId(accessToken, environment, deleteAPIName, deleteAPIVersion, deleteAPIProvider)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id for deletion ", err)
	}
	url := deleteEndpoint + apiId
	utils.Logln(utils.LogPrefixInfo+"DeleteAPI: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeDELETERequest(url, headers)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return nil, errors.New(strconv.Itoa(resp.StatusCode()) + ":<" + string(resp.Body()) + ">")
	}
	return resp, nil
}

func PrintDeleteAPIResponse(resp *resty.Response, err error) {
	if err != nil {
		fmt.Println("Error deleting API:", err)
	} else {
		fmt.Println("API deleted successfully!. Status: " + strconv.Itoa(resp.StatusCode()))
	}
}
