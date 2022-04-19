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
	"github.com/json-iterator/go"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io/ioutil"
	"path/filepath"
)

// ExportThrottlePoliciesFromEnv function is used with export throttlepolicies command
func ExportThrottlingPoliciesFromEnv(accessToken, exportEnvironment string, exportThrottlePoliciesType string) (*resty.Response, error) {
	adminEndpoint := utils.GetAdminEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportThrottlePolicies(adminEndpoint, accessToken, exportThrottlePoliciesType)
}

func exportThrottlePolicies(adminEndpoint, accessToken string, ThrottlePoliciesType string) (*resty.Response, error) {
	var policytype string
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	ThrottlePolicyresource := "throttling/policies/export?name=TestPolicy&type=custom&format=JSON"
	//ThrottlePolicyresource := "throttling/policies/"
	switch ThrottlePoliciesType {
	case "sub":
		policytype = "subscription"
	case "app":
		policytype = "application"
	case "deny":
		policytype = "deny-policies"
		ThrottlePolicyresource = "throttling/"
	case "advanced":
		policytype = "advanced"
	case "custom":
		policytype = "custom"
	case "test":
		policytype = "export"
	}

	url := adminEndpoint + ThrottlePolicyresource

	policytype += "hi"

	utils.Logln(utils.LogPrefixInfo+"ExportThrottlingPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)
	fmt.Println(ThrottlePolicyresource)

	fmt.Println(resp.String())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ThrottlePoliciesWriteToZip(PolicyType, zipLocationPath string, runningExportThrottlePoliciesCommand bool, resp *resty.Response) {
	var path = ""
	switch PolicyType {
	case "sub":
		path = WriteSubscriptionThrottlingPolicies(zipLocationPath, resp)
	case "app":
		path = WriteApplicationThrottlingPolicies(zipLocationPath, resp)
	case "deny":
		path = WriteDenyThrottlingPolicies(zipLocationPath, resp)
	case "advanced":
		path = WriteAdvancedThrottlingPolicies(zipLocationPath, resp)
	case "custom":
		path = WriteCustomThrottlingPolicies(zipLocationPath, resp)
	}
	// Output the final zip file location.
	if runningExportThrottlePoliciesCommand {
		fmt.Println("Successfully exported Throttling Policies!")
		fmt.Println("Find the exported Throttling Policies at " + path)
	}

}

func WriteSubscriptionThrottlingPolicies(zipLocationPath string, resp *resty.Response) string {
	zipLocationPath = filepath.Join(zipLocationPath, "Subscription-Policies")
	var zipFilename string

	err := utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+zipLocationPath, err)
	}

	var ThrottlingPolicyListResponse utils.SubscriptionExportThrottlePolicyList
	err = json.Unmarshal(resp.Body(), &ThrottlingPolicyListResponse)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshelling data", err)
	}

	var ThrottlePolicyList = ThrottlingPolicyListResponse.List
	var Count = ThrottlingPolicyListResponse.Count

	for i := 0; i < Count; i++ {
		policyContent := ThrottlePolicyList[i]
		zipFilename = policyContent.PolicyName
		marshaledData, err := jsoniter.MarshalIndent(policyContent, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshelling policy content", err)
		}
		_, _ = WriteThrottlingPolicy(zipLocationPath, zipFilename, marshaledData)
	}
	return zipLocationPath
}

func WriteApplicationThrottlingPolicies(zipLocationPath string, resp *resty.Response) string {
	zipLocationPath = filepath.Join(zipLocationPath, "Application-Policies")
	var zipFilename string

	err := utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+zipLocationPath, err)
	}

	var ThrottlingPolicyListResponse utils.ApplicationExportThrottlePolicyList

	err = json.Unmarshal(resp.Body(), &ThrottlingPolicyListResponse)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshelling data", err)
	}

	var ThrottlePolicyList = ThrottlingPolicyListResponse.List
	var Count = ThrottlingPolicyListResponse.Count

	for i := 0; i < Count; i++ {
		policyContent := ThrottlePolicyList[i]
		zipFilename = policyContent.PolicyName
		marshaledData, err := jsoniter.MarshalIndent(policyContent, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshelling policy content", err)
		}
		_, _ = WriteThrottlingPolicy(zipLocationPath, zipFilename, marshaledData)
	}
	return zipLocationPath
}

func WriteAdvancedThrottlingPolicies(zipLocationPath string, resp *resty.Response) string {
	zipLocationPath = filepath.Join(zipLocationPath, "Advanced-Policies")
	var zipFilename string

	err := utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+zipLocationPath, err)
	}

	var ThrottlingPolicyListResponse utils.AdvancedExportThrottlePolicyList

	err = json.Unmarshal(resp.Body(), &ThrottlingPolicyListResponse)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshelling data", err)
	}

	var ThrottlePolicyList = ThrottlingPolicyListResponse.List
	var Count = ThrottlingPolicyListResponse.Count

	for i := 0; i < Count; i++ {
		policyContent := ThrottlePolicyList[i]
		zipFilename = policyContent.PolicyName
		marshaledData, err := jsoniter.MarshalIndent(policyContent, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshelling policy content", err)
		}
		_, _ = WriteThrottlingPolicy(zipLocationPath, zipFilename, marshaledData)
	}
	return zipLocationPath
}

func WriteCustomThrottlingPolicies(zipLocationPath string, resp *resty.Response) string {
	zipLocationPath = filepath.Join(zipLocationPath, "Custom-Policies")
	var zipFilename string

	err := utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+zipLocationPath, err)
	}

	var ThrottlingPolicyListResponse utils.CustomExportThrottlePolicyList

	err = json.Unmarshal(resp.Body(), &ThrottlingPolicyListResponse)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshelling data", err)
	}

	var ThrottlePolicyList = ThrottlingPolicyListResponse.List
	var Count = ThrottlingPolicyListResponse.Count

	for i := 0; i < Count; i++ {
		policyContent := ThrottlePolicyList[i]
		zipFilename = policyContent.PolicyName
		marshaledData, err := jsoniter.MarshalIndent(policyContent, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshelling policy content", err)
		}
		_, _ = WriteThrottlingPolicy(zipLocationPath, zipFilename, marshaledData)
	}
	return zipLocationPath
}

func WriteDenyThrottlingPolicies(zipLocationPath string, resp *resty.Response) string {
	zipLocationPath = filepath.Join(zipLocationPath, "Deny-Policies")
	var zipFilename string

	err := utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+zipLocationPath, err)
	}

	var ThrottlingPolicyListResponse utils.DenyExportThrottlePolicyList
	err = json.Unmarshal(resp.Body(), &ThrottlingPolicyListResponse)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshelling data", err)
	}

	var ThrottlePolicyList = ThrottlingPolicyListResponse.List
	var Count = ThrottlingPolicyListResponse.Count

	for i := 0; i < Count; i++ {
		policyContent := ThrottlePolicyList[i]
		zipFilename = policyContent.ConditionId
		marshaledData, err := jsoniter.MarshalIndent(policyContent, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshelling policy content", err)
		}
		_, _ = WriteThrottlingPolicy(zipLocationPath, zipFilename, marshaledData)
	}
	return zipLocationPath
}

//Writing as .yaml
//func WriteThrottlingPolicy(zipLocationPath string, zipFilename string, marshaledData []byte) (string, error) {
//	TempJsonFile := zipFilename
//	TempJsonFile += ".yaml"
//	tmpDir, err := ioutil.TempDir("", "apim")
//	if err != nil {
//		_ = os.RemoveAll(tmpDir)
//		return "", err
//	}
//
//	tempFile := filepath.Join(tmpDir, TempJsonFile)
//
//	jsonMetaData, err := gabs.ParseJSON(marshaledData)
//	metaContent, err := utils.JsonToYaml(jsonMetaData.Bytes())
//
//	//write the content to temp file
//	err = ioutil.WriteFile(tempFile, metaContent, 0644)
//	if err != nil {
//		utils.HandleErrorAndExit("Error creating temp file", err)
//	}
//
//	targetZipFile := filepath.Join(zipLocationPath, zipFilename)
//	targetZipFile += ".zip"
//	err = utils.Zip(tempFile, targetZipFile)
//	return targetZipFile, err
//}

////////////////////////////////////////////////////
//func WriteThrottlingPolicy(zipLocationPath string, zipFilename string, marshaledData []byte) (string, error) {
//	TempJsonFile := zipFilename
//	TempJsonFile += ".json"
//	tmpDir, err := ioutil.TempDir("", "apim")
//	if err != nil {
//		_ = os.RemoveAll(tmpDir)
//		return "", err
//	}
//
//	tempFile := filepath.Join(tmpDir, TempJsonFile)
//
//	err = ioutil.WriteFile(tempFile, marshaledData, 0644)
//	if err != nil {
//		utils.HandleErrorAndExit("Error writing temp json", err)
//	}
//
//	targetZipFile := filepath.Join(zipLocationPath, zipFilename)
//	targetZipFile += ".zip"
//	err = utils.Zip(tempFile, targetZipFile)
//	return targetZipFile, err
//}

func WriteThrottlingPolicy(FilePath string, JsonFilename string, marshaledData []byte) (string, error) {

	JsonFilename += ".json"
	Filename := filepath.Join(FilePath, JsonFilename)
	err := ioutil.WriteFile(Filename, marshaledData, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Error writing json", err)
	}
	return FilePath, err
}

/////////////////////////////////////////////////////////
