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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aybabtme/orderedjson"
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
)

// const (
// 	CmdPolicyTypeSubscription = "sub"
// 	CmdPolicyTypeApplication  = "app"
// 	CmdPolicyTypeAdvanced     = "advanced"
// 	CmdPolicyTypeCustom       = "custom"

// 	QueryPolicyTypeSubscription = "sub"
// 	QueryPolicyTypeApplication  = "app"
// 	QueryPolicyTypeAdvanced     = "api"
// 	QueryCmdPolicyTypeCustom    = "global"

// 	ExportPolicyTypeSubscription = "subscription policy"
// 	ExportPolicyTypeApplication  = "application policy"
// 	ExportPolicyTypeAdvanced     = "advanced policy"
// 	ExportPolicyTypeCustom       = "custom rule"

// 	ExportPolicyFileNamePrefixSubscription = "Subscription"
// 	ExportPolicyFileNamePrefixApplication  = "Application"
// 	ExportPolicyFileNamePrefixAdvanced     = "Advanced"
// 	ExportPolicyFileNamePrefixCustom       = "Custom"
// )

// ExportOperationPolicyFromEnv function is used with export policy rate-limiting command
func ExportOperationPolicyFromEnv(accessToken string, exportEnvironment string, operationPolicyName string, operationPolicyVersion string) (*resty.Response, error) {
	operationPolicyEndpoint := utils.GetPublisherEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	// var query string
	operationPolicyEndpoint = utils.AppendSlashToString(operationPolicyEndpoint)
	// operationPolicyResource := "operation-policies/export?"
	operationPolicyResource := "operation-policies/3931fe31-69fc-4a43-a4fe-6054a7d96a5d" + "/content"

	// query = `name=` + operationPolicyName + `&version=` + operationPolicyVersion

	// operationPolicyResource += query
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

	// zipLocationPath := "/Users/benura/Desktop/"

	// err = utils.CreateDirIfNotExist(zipLocationPath)
	// if err != nil {
	// 	utils.HandleErrorAndExit("Error creating dir to store zip archive: "+zipLocationPath, err)
	// }
	// exportedFinalZip := filepath.Join(zipLocationPath, zipFilename)

	// err = utils.Zip(tmpClonedLoc, exportedFinalZip)
	// if err != nil {
	// 	utils.HandleErrorAndExit("Error creating the final zip archive", err)
	// }

	// fileName, marshaledData := resolveOperationPolicy(ExportFormat, resp)
	// _, _ = operationPolicyWrite(ExportLocationPath, fileName, marshaledData)
	if runningExportThrottlePolicyCommand {
		fmt.Println("Successfully exported Operation Policy!")
		fmt.Println("Find the exported Operation Policies at " +
			utils.AppendSlashToString(exportLocationPath) + zipFileName)
	}
}

// // resolves the policy file name with the policy type
// func resolveExportFileName(policyType, policyName string) string {
// 	var fileName string
// 	switch policyType {
// 	case ExportPolicyTypeSubscription:
// 		fileName = ExportPolicyFileNamePrefixSubscription
// 	case ExportPolicyTypeApplication:
// 		fileName = ExportPolicyFileNamePrefixApplication
// 	case ExportPolicyTypeAdvanced:
// 		fileName = ExportPolicyFileNamePrefixAdvanced
// 	case ExportPolicyTypeCustom:
// 		fileName = ExportPolicyFileNamePrefixCustom
// 	}
// 	fileName = fileName + `-` + policyName
// 	return fileName
// }

// resolveThrottlePolicy resolves the policy file name and the marshalled data
func resolveOperationPolicy(exportThrottlePolicyFormat string, resp *resty.Response) (string, []byte) {
	var marshaledData []byte
	var ExportOperationPolicy utils.ExportThrottlePolicy
	fmt.Println(resp)
	err := yaml.Unmarshal(resp.Body(), &ExportOperationPolicy)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshalling response data", err)
	}
	policyType := ExportOperationPolicy.Subtype
	policyName := ExportOperationPolicy.Data[1].Value
	throttlingPolicyType := fmt.Sprintf("%v", policyType)
	throttlePolicyName := fmt.Sprintf("%v", policyName)
	fileName := resolveExportFileName(throttlingPolicyType, throttlePolicyName)
	if exportThrottlePolicyFormat == utils.DefaultExportFormat {
		fileName += ".yaml"
		if err != nil {
			utils.HandleErrorAndExit("Error marshaling policy content", err)
		}
		marshaledData, _ = yaml.Marshal(ExportOperationPolicy)
	} else {
		var s orderedjson.Map
		err = json.Unmarshal(resp.Body(), &s)
		marshaledData, _ = json.MarshalIndent(s, "", " ")
		fileName += ".json"
	}
	return fileName, marshaledData
}

func operationPolicyWrite(filePath string, fileName string, marshaledData []byte) (string, error) {
	fileName = filepath.Join(filePath, fileName)
	err := ioutil.WriteFile(fileName, marshaledData, os.ModePerm)
	if err != nil {
		utils.HandleErrorAndExit("Error writing file"+fileName, err)
	}
	return filePath, err
}
