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

package mg

import (
	"errors"
	"net/http"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var mgUndeployAPIResourcePath = "/apis"

// UndeployAPI sends a DELETE request to delete an API
func UndeployAPI(env string, queryParam map[string]string) (err error) {
	mgwAdapterInfo, err := GetStoredTokenAndHost(env)
	if err != nil {
		return err
	}
	apiDeleteEndpoint := mgwAdapterInfo.Host + utils.DefaultMgwAdapterEndpointSuffix + mgUndeployAPIResourcePath

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + mgwAdapterInfo.AccessToken
	resp, err := utils.InvokeDELETERequestWithParams(apiDeleteEndpoint, queryParam, headers)

	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	} else if resp.StatusCode() == http.StatusNotFound {
		return errors.New("the API does not exist")
	}
	return errors.New(string(resp.Body()))
}
