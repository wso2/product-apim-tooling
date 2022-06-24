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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// ExportOperationPolicyFromEnv function is used with export policy rate-limiting command
func ExportOperationPolicyFromEnv(accessToken string, exportEnvironment string, operationPolicyName string, operationPolicyVersion string) (*resty.Response, error) {
	operationPolicyEndpoint := utils.GetPublisherEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	// var query string
	operationPolicyEndpoint = utils.AppendSlashToString(operationPolicyEndpoint)
	// operationPolicyResource := "operation-policies/export?"
	operationPolicyResource := "operation-policies/export?"

	query := `name=` + operationPolicyName + `&version=` + operationPolicyVersion

	operationPolicyResource += query
	url := operationPolicyEndpoint + operationPolicyResource
	utils.Logln(utils.LogPrefixInfo+"ExportOperationPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WriteOperationPolicyToFile writes the policy to a specified location
func WriteOperationPolicyToFile(exportLocationPath string, resp *resty.Response, exportOperationPolicyVersion string, exportOperationPolicyName string,
	runningExportThrottlePolicyCommand bool) {
	err := utils.CreateDirIfNotExist(exportLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+exportLocationPath, err)
	}
	zipFileName := exportOperationPolicyName + "_" + exportOperationPolicyVersion + ".zip"
	zipFile := filepath.Join(exportLocationPath, zipFileName)

	fmt.Println(zipFile)

	err = ioutil.WriteFile(zipFile, resp.Body(), 0644)
	if err != nil {
		return
	}

	if err != nil {
		utils.HandleErrorAndExit("Error creating the temporary zip file to store the exported API", err)
	}

	if runningExportThrottlePolicyCommand {
		fmt.Println("Successfully exported Operation Policy!")
		fmt.Println("Find the exported Operation Policies at " +
			utils.AppendSlashToString(exportLocationPath) + zipFileName)
	}
}
