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

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func AIDeleteAPIs(OnPremKey, Endpoint, Tenant string) error {

	fmt.Println("Removing existing APIs and API Products from vector DB for tenant:", Tenant)

	headers := make(map[string]string)
	headers["API-KEY"] = OnPremKey
	headers["TENANT-DOMAIN"] = Tenant

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
			continue
		}

		jsonResp := map[string]map[string]int32{}

		err := json.Unmarshal(resp.Body(), &jsonResp)

		if err != nil {
			utils.HandleErrorAndContinue("Error in unmarshalling response:", err)
			continue
		}

		fmt.Printf("Removed %d APIs and API Products successfully from vector database (attempt %d)\n", jsonResp["message"]["delete_count"], attempt)
		return nil
	}

	if deleteErr != nil {
		return fmt.Errorf("Error removing existing APIs and API Products after retry: %v", deleteErr)
	}
	return fmt.Errorf("Removing existing APIs and API Products failed after retry")
}
