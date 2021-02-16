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
	"fmt"
	"net/http"
	"strconv"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

//DeployAPI creats or updates an API in the microgateway depending on the overwrite param
func DeployAPI(endpoint, filePath, accessToken string, extraParams map[string]string,
	importAPISkipCleanup bool, overwrite bool) {
	//TODO: (VirajSalaka) support substituting parameters with params file. At the moment it is in hold on state, as the decision to use environments is
	//not finalized yet.
	// if apiFilePath contains a directory, zip it. Otherwise, leave it as it is.
	filePath, err, cleanupFunc := utils.CreateZipFileFromProject(filePath, importAPISkipCleanup)
	if err != nil {
		utils.HandleErrorAndExit("Error adding API to microgateway", err)
	}
	//cleanup the temporary artifacts once consuming the zip file
	if cleanupFunc != nil {
		defer cleanupFunc()
	}

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + accessToken
	headers[utils.HeaderAccept] = "application/json"
	headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive

	if overwrite {
		UpdateAPI(endpoint, extraParams, headers, "file", filePath)
	} else {
		AddAPI(endpoint, extraParams, headers, "file", filePath)
	}
}

//AddAPI creats an API in the microgateway
func AddAPI(endpoint string, extraParams, headers map[string]string,
	fileParamName string, filePath string) {
	resp, err := utils.InvokePOSTRequestWithFileAndQueryParams(extraParams, endpoint, headers,
		"file", filePath)
	if err != nil {
		utils.HandleErrorAndExit("Error deploying API.", err)
	}
	if resp.StatusCode() == http.StatusOK {
		fmt.Println("Successfully deployed API.")
	} else if resp.StatusCode() == http.StatusConflict {
		fmt.Println("Unable to deploy API. API already exists. Status: " + resp.Status())
	} else {
		fmt.Println("Unable to deploy API. Error Status: " + resp.Status())
	}
}

//UpdateAPI updates an API in the microgateway
func UpdateAPI(endpoint string, extraParams, headers map[string]string,
	fileParamName string, filePath string) {

	endpoint += "?overwrite=" + strconv.FormatBool(true)
	resp, err := utils.InvokePOSTRequestWithFileAndQueryParams(extraParams, endpoint, headers,
		"file", filePath)
	if err != nil {
		utils.HandleErrorAndExit("Error updating API.", err)
	}
	if resp.StatusCode() == http.StatusOK {
		fmt.Println("Successfully updated the API.")
	} else if resp.StatusCode() == http.StatusNotFound {
		fmt.Println("Unable to update API. API does not exist. Status: " + resp.Status())
	} else {
		fmt.Println("Unable to update API. Error Status: " + resp.Status())
	}
}
