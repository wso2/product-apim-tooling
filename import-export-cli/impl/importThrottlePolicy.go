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
)

func ImportThrottlingPolicyToEnv(accessOAuthToken string, importEnvironment string, importThrottlingPolicyFile string) error {
	adminEndpoint := utils.GetAdminEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	return ImportThrottlingPolicy(accessOAuthToken, adminEndpoint, importThrottlingPolicyFile)
}

func ImportThrottlingPolicy(accessOAuthToken string, adminEndpoint string, importPath string) error {
	//exportDirectory := filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)
	//resolvedThrottlingPolicyFilePath, err := resolveImportFilePath(importPath, exportDirectory)
	//if err != nil {
	//	return err
	//}
	//utils.Logln(utils.LogPrefixInfo+"throttling Policy Location:", resolvedThrottlingPolicyFilePath)

	utils.Logln(utils.LogPrefixInfo + "Creating workspace")
	//tmpPath, err := utils.GetTempCloneFromDirOrZip(resolvedThrottlingPolicyFilePath)
	//if err != nil {
	//	return err
	//}
	uri := adminEndpoint + "/throttling/deny-policies"
	/////////////////////////////////
	if _, err := os.Stat(importPath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	//path, data, err := resolveYamlOrJSON(importPath)
	//fmt.Println(path)
	//fmt.Println(data)
	//if err != nil {
	//	utils.Logln(utils.LogPrefixError, err)
	//	return err
	//}
	//fmt.Println(uri)
	//fmt.Println(importPath)
	tmpPath := importPath
	extraParams := map[string]string{}
	err := importThrottlingPolicy(uri, tmpPath, accessOAuthToken, extraParams, true)
	return err
}

// importAPI imports an API to the API manager
func importThrottlingPolicy(endpoint, filePath, accessToken string, extraParams map[string]string, isOauth bool) error {
	resp, err := ExecuteThrottlingPolicyUploadRequest(endpoint, extraParams, "file",
		filePath, accessToken, isOauth)
	//fmt.Println(resp.RawResponse)
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
		fmt.Println("Error importing Throttling Policy.")
		fmt.Println("Status: " + resp.Status())
		fmt.Println("Response:", resp.IsSuccess())
		return errors.New(resp.Status())
	}
}

func ExecuteThrottlingPolicyUploadRequest(uri string, params map[string]string, paramName, path,
	accessToken string, isOAuthToken bool) (*resty.Response, error) {

	headers := make(map[string]string)
	if isOAuthToken {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	} else {
		headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + accessToken
	}
	headers[utils.HeaderContentType] = "application/json"
	headers[utils.HeaderAccept] = "application/json"
	headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive
	fmt.Println(headers)
	fmt.Println(uri)
	type Data struct {
		PolicyName   string `json:"policyName"`
		DisplayName  string `json:"displayName"`
		Description  string `json:"description"`
		IsDeployed   bool   `json:"isDeployed"`
		Type         string `json:"type"`
		DefaultLimit struct {
			Type         string `json:"type"`
			RequestCount struct {
				TimeUnit     string `json:"timeUnit"`
				UnitTime     int    `json:"unitTime"`
				RequestCount int    `json:"requestCount"`
			} `json:"requestCount"`
			Bandwidth struct {
				TimeUnit   string `json:"timeUnit"`
				UnitTime   int    `json:"unitTime"`
				DataAmount int    `json:"dataAmount"`
				DataUnit   string `json:"dataUnit"`
			} `json:"bandwidth"`
			EventCount struct {
				TimeUnit   string `json:"timeUnit"`
				UnitTime   int    `json:"unitTime"`
				EventCount int    `json:"eventCount"`
			} `json:"eventCount"`
		} `json:"defaultLimit"`
	}

	json_string := `
{
"conditionId": "b513eb68-69e8-4c32-92cf-852c101363cf",
"conditionType": "IP",
"conditionValue": {
"fixedIp": "192.168.1.1",
"invert": false
},
"conditionStatus": true
}
`

	var data utils.DenyThrottlingPolicy
	err := json.Unmarshal([]byte(json_string), &data)

	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		fmt.Println("Cannot unmarshal")
	}
	return utils.InvokePOSTRequest(uri, headers, data)
	//return utils.InvokePUTRequestWithoutQueryParams(uri, headers, data)
}
