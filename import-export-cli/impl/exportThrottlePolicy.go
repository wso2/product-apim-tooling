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
	"github.com/aybabtme/orderedjson"
	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

const (
	CmdPolicyTypeSubscription = "sub"
	CmdPolicyTypeApplication  = "app"
	CmdPolicyTypeAdvanced     = "advanced"
	CmdPolicyTypeCustom       = "custom"

	QueryPolicyTypeSubscription = "sub"
	QueryPolicyTypeApplication  = "app"
	QueryPolicyTypeAdvanced     = "api"
	QueryCmdPolicyTypeCustom    = "global"

	ExportPolicyTypeSubscription = "subscription policy"
	ExportPolicyTypeApplication  = "application policy"
	ExportPolicyTypeAdvanced     = "advanced policy"
	ExportPolicyTypeCustom       = "custom rule"

	ExportPolicyFileNamePrefixSubscription = "Subscription"
	ExportPolicyFileNamePrefixApplication  = "Application"
	ExportPolicyFileNamePrefixAdvanced     = "Advanced"
	ExportPolicyFileNamePrefixCustom       = "Custom"
)

// ExportThrottlingPolicyFromEnv function is used with export policy rate-limiting command
func ExportThrottlingPolicyFromEnv(accessToken string, exportEnvironment string, exportThrottlePolicyName string,
	exportThrottlePolicyType string, exportFormat string) (*resty.Response, error) {
	adminEndpoint := utils.GetAdminEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportThrottlePolicy(adminEndpoint, accessToken, exportThrottlePolicyName, exportThrottlePolicyType, exportFormat)
}

func exportThrottlePolicy(adminEndpoint, accessToken string, ThrottlePolicyName string, ThrottlePolicyType string,
	exportFormat string) (*resty.Response, error) {
	var PolicyType string
	var query string
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	ThrottlePolicyResource := "throttling/policies/export?"
	if ThrottlePolicyType != "" {
		switch ThrottlePolicyType {
		case CmdPolicyTypeSubscription:
			PolicyType = QueryPolicyTypeSubscription
		case CmdPolicyTypeApplication:
			PolicyType = QueryPolicyTypeApplication
		case CmdPolicyTypeAdvanced:
			PolicyType = QueryPolicyTypeAdvanced
		case CmdPolicyTypeCustom:
			PolicyType = QueryCmdPolicyTypeCustom
		}
		query = `name=` + ThrottlePolicyName + `&type=` + PolicyType + `&format=` + exportFormat
	} else {
		query = `name=` + ThrottlePolicyName + `&format=` + exportFormat
	}

	ThrottlePolicyResource += query
	url := adminEndpoint + ThrottlePolicyResource
	utils.Logln(utils.LogPrefixInfo+"ExportThrottlingPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func WriteThrottlePolicyToFile(ExportLocationPath string, resp *resty.Response, ExportFormat string,
	runningExportThrottlePolicyCommand bool) {
	err := utils.CreateDirIfNotExist(ExportLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+ExportLocationPath, err)
	}
	fileName, marshaledData := ResolveThrottlePolicy(ExportFormat, resp)
	_, _ = ThrottlingPolicyWrite(ExportLocationPath, fileName, marshaledData)
	if runningExportThrottlePolicyCommand {
		fmt.Println("Successfully exported Throttling Policy!")
		fmt.Println("Find the exported Throttling Policies at " +
			utils.AppendSlashToString(ExportLocationPath) + fileName)
	}
}

func ResolveExportFileName(policyType string, policyName string) string {
	var fileName string
	switch policyType {
	case ExportPolicyTypeSubscription:
		fileName = ExportPolicyFileNamePrefixSubscription
	case ExportPolicyTypeApplication:
		fileName = ExportPolicyFileNamePrefixApplication
	case ExportPolicyTypeAdvanced:
		fileName = ExportPolicyFileNamePrefixAdvanced
	case ExportPolicyTypeCustom:
		fileName = ExportPolicyFileNamePrefixCustom
	}
	fileName = fileName + `-` + policyName
	return fileName
}

// ResolveThrottlePolicy resolves the policy file name and the marshalled data
func ResolveThrottlePolicy(exportThrottlePolicyFormat string, resp *resty.Response) (string, []byte) {
	var marshaledData []byte
	var ExportThrottlingPolicy utils.ExportThrottlePolicy
	err := yaml.Unmarshal(resp.Body(), &ExportThrottlingPolicy)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshalling response data", err)
	}
	policyType := ExportThrottlingPolicy.Subtype
	policyName := ExportThrottlingPolicy.Data[1].Value
	throttlingPolicyType := fmt.Sprintf("%v", policyType)
	throttlePolicyName := fmt.Sprintf("%v", policyName)
	fileName := ResolveExportFileName(throttlingPolicyType, throttlePolicyName)
	if exportThrottlePolicyFormat == utils.DefaultExportFormat {
		fileName += ".yaml"
		if err != nil {
			utils.HandleErrorAndExit("Error marshaling policy content", err)
		}
		marshaledData, _ = yaml.Marshal(ExportThrottlingPolicy)
	} else {
		var s orderedjson.Map
		err = json.Unmarshal(resp.Body(), &s)
		marshaledData, _ = json.MarshalIndent(s, "", " ")
		fileName += ".json"
	}
	return fileName, marshaledData
}

func ThrottlingPolicyWrite(FilePath string, Filename string, marshaledData []byte) (string, error) {
	Filename = filepath.Join(FilePath, Filename)
	err := ioutil.WriteFile(Filename, marshaledData, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Error writing file", err)
	}
	return FilePath, err
}
