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
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	uploadedAPIs         int32
	totalAPIs            int32
	Credential           credentials.Credential
	CmdUploadEnvironment string
	UploadProducts       bool
	UploadAll            bool
	AIToken              string
	Tenant               string
	Endpoint             = utils.DefaultAIEndpoint
)

func AIDeleteAPIs(credential credentials.Credential, CmdPurgeEnvironment, aiToken, oldEndpoint, tenant string) {

	headers := make(map[string]string)
	if aiToken != "" {
		headers["Authorization"] = "Bearer " + aiToken
	} else {
		AIToken = utils.AIToken
		headers["API-KEY"] = AIToken
	}

	if (oldEndpoint != "") {
		Endpoint = oldEndpoint
	} else {
		Endpoint = utils.GetAIServiceEndpointOfEnv(CmdPurgeEnvironment, utils.MainConfigFilePath)
	}

	fmt.Println("Removing existing APIs and API Products from vector database for tenant:", tenant)

	headers["TENANT-DOMAIN"] = tenant
	headers["User-Agent"] = "WSO2-API-Controller"
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON

	var resp *resty.Response
	var deleteErr error

	for attempt := 1; attempt <= 2; attempt++ {
		resp, deleteErr = utils.InvokeDELETERequest(Endpoint+"/ai/spec-populator/bulk-remove", headers)
		if deleteErr != nil {
			fmt.Printf("Error removing existing APIs and API Products (attempt %d): %v\n", attempt, deleteErr)
			continue
		}

		if resp.StatusCode() != 200 {
			fmt.Printf("Removing existing APIs and API Products failed with status %d %s (attempt %d)\n", resp.StatusCode(), resp.Body(), attempt)
			if attempt == 2 {
				fmt.Println("Removing existing APIs and API Products failed.")
				os.Exit(1)
			}
			continue
		}

		jsonResp := map[string]map[string]int32{}

		err := json.Unmarshal(resp.Body(), &jsonResp)

		if err != nil {
			utils.HandleErrorAndContinue("Error in unmarshalling response:", err)
			continue
		}

		fmt.Printf("Removed %d APIs and API Products successfully from vector database for tenant: %s (attempt %d)\n", jsonResp["message"]["delete_count"], tenant, attempt)
		return
	}

	if deleteErr != nil {
		utils.HandleErrorAndExit("Error removing existing APIs and API Products after retry:", deleteErr)
	}
}
