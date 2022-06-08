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
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"os"
	"strconv"
)

func ImportThrottlingPolicyToEnv(accessOAuthToken, importEnvironment, importThrottlingPolicyFile string, importThrottlePolicyUpdate bool) error {
	adminEndpoint := utils.GetAdminEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	return ImportThrottlingPolicy(accessOAuthToken, adminEndpoint, importThrottlingPolicyFile, importThrottlePolicyUpdate)
}

func ImportThrottlingPolicy(accessOAuthToken, adminEndpoint, importPath string, importThrottlePolicyUpdate bool) error {
	if _, err := os.Stat(importPath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	uri := adminEndpoint + "/throttling/policies/import"
	err := importThrottlingPolicy(uri, importPath, accessOAuthToken, true, importThrottlePolicyUpdate)
	return err
}

func importThrottlingPolicy(endpoint string, importPath string, accessToken string, isOauth bool, ThrottlePolicyUpdate bool) error {
	resp, err := executeThrottlingPolicyUploadRequest(endpoint, importPath, ThrottlePolicyUpdate, accessToken, isOauth)
	utils.Logf("Response : %v", resp)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}
	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		// 201 Created or 200 OK
		fmt.Println(resp.String())
		return nil
	} else {
		// We have an HTTP error
		if resp.StatusCode() == http.StatusConflict && ThrottlePolicyUpdate {
			fmt.Println("Cannot Update")
		}
		fmt.Println("Error importing Throttling Policy.")
		fmt.Println("Status: " + resp.Status())
		fmt.Println("Response:", resp.IsSuccess())

		return errors.New(resp.Status())
	}
}

func executeThrottlingPolicyUploadRequest(uri string, importPath string, update bool, accessToken string, isOAuthToken bool) (*resty.Response, error) {

	headers := make(map[string]string)
	if isOAuthToken {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	} else {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + accessToken
	}
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive
	params := make(map[string]string)
	params["overwrite"] = strconv.FormatBool(update)
	return utils.InvokePOSTRequestWithFileAndQueryParams(params, uri, headers, "file", importPath)
}
