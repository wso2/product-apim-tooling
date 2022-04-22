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
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"os"
	"strconv"
)

func ImportThrottlingPolicyToEnv(accessOAuthToken string, importEnvironment string, importThrottlingPolicyFile string, importThrottlePolicyUpdate bool) error {
	adminEndpoint := utils.GetAdminEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	return ImportThrottlingPolicy(accessOAuthToken, adminEndpoint, importThrottlingPolicyFile, importThrottlePolicyUpdate)
}

func ImportThrottlingPolicy(accessOAuthToken string, adminEndpoint string, importPath string, importThrottlePolicyUpdate bool) error {

	if _, err := os.Stat(importPath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	_, data, err := resolveYamlOrJSON(importPath)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}

	var policy utils.ExportThrottlePolicy

	err = json.Unmarshal(data, &policy)

	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}

	uri := adminEndpoint + "/throttling/policies/import"
	query := `?overwrite=` + strconv.FormatBool(importThrottlePolicyUpdate)
	uri += query
	err = importThrottlingPolicy(uri, policy, accessOAuthToken, true, importThrottlePolicyUpdate)
	return err
}

func importThrottlingPolicy(endpoint string, PolicyDetails interface{}, accessToken string, isOauth bool, ThrottlePolicyUpdate bool) error {
	resp, err := ExecuteThrottlingPolicyUploadRequest(endpoint, PolicyDetails, accessToken, isOauth)
	utils.Logf("Response : %v", resp)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}
	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		// 201 Created or 200 OK
		fmt.Println("Successfully imported Throttling Policy.")
		return nil
	} else {
		// We have an HTTP error

		if resp.StatusCode() == http.StatusConflict && ThrottlePolicyUpdate {
			fmt.Println("Cannot Update")
			//Execute Throttle Policy update
		}
		fmt.Println("Error importing Throttling Policy.")
		fmt.Println("Status: " + resp.Status())
		fmt.Println("Response:", resp.IsSuccess())

		return errors.New(resp.Status())
	}
}

func ExecuteThrottlingPolicyUploadRequest(uri string, PolicyDetails interface{}, accessToken string, isOAuthToken bool) (*resty.Response, error) {

	headers := make(map[string]string)
	if isOAuthToken {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	} else {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + accessToken
	}
	headers[utils.HeaderContentType] = "application/json"
	headers[utils.HeaderAccept] = "application/json"
	headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive

	return utils.InvokePOSTRequest(uri, headers, PolicyDetails)
}
