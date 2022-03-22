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
	"os"
	"path/filepath"
)

// ExportThrottlePoliciesFromEnv function is used with export throttlepolicies command
func ExportThrottlingPoliciesFromEnv(accessToken, exportEnvironment string, exportThrottlePoliciesType string) (*resty.Response, error) {
	adminEndpoint := utils.GetAdminEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportThrottlePolicies(adminEndpoint, accessToken, exportThrottlePoliciesType)
}

// exportAPI function is used with export api command
// @param name : Name of the API to be exported
// @param version : Version of the API to be exported
// @param provider : Provider of the API
// @param publisherEndpoint : API Manager Publisher Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func exportThrottlePolicies(adminEndpoint, accessToken string, ThrottlePoliciesType string) (*resty.Response, error) {
	var policytype string
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	ThrottlePolicyresource := "throttling/policies/"
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
	}

	url := adminEndpoint + ThrottlePolicyresource + policytype

	utils.Logln(utils.LogPrefixInfo+"ExportThrottlingPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// WriteToZip
// @param exportAPIName : Name of the API to be exported
// @param exportAPIVersion: Version of the API to be exported
// @param exportAPIRevisionNumber: Revision number of the api
// @param zipLocationPath: Path to the export directory
// @param runningExportApiCommand: Whether the export API command is running
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported API will be written to a zip file
func ThrottlePoliciesWriteToZip(PolicyType, zipLocationPath string, runningExportThrottlePoliciesCommand bool, resp *resty.Response) {

	switch PolicyType {
	case "sub":
		WriteSubscriptionThrottlingPolicies(zipLocationPath, resp)
	case "app":
		WriteApplicationThrottlingPolicies(zipLocationPath, resp)
	case "deny":
		WriteDenyThrottlingPolicies(zipLocationPath, resp)
	case "advanced":
		WriteAdvancedThrottlingPolicies(zipLocationPath, resp)
	case "custom":
		WriteCustomThrottlingPolicies(zipLocationPath, resp)
	}
	// Output the final zip file location.
	if runningExportThrottlePoliciesCommand {
		fmt.Println("Successfully exported Throttling Policies!")
		fmt.Println("Find the exported Throttling Policies at " + zipLocationPath)
	}

}

func WriteSubscriptionThrottlingPolicies(zipLocationPath string, resp *resty.Response) {

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
}

func WriteApplicationThrottlingPolicies(zipLocationPath string, resp *resty.Response) {

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
}

func WriteAdvancedThrottlingPolicies(zipLocationPath string, resp *resty.Response) {

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
}

func WriteCustomThrottlingPolicies(zipLocationPath string, resp *resty.Response) {

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
}

func WriteDenyThrottlingPolicies(zipLocationPath string, resp *resty.Response) {

	zipLocationPath = filepath.Join(zipLocationPath, "Deny-Policies")
	var zipFilename string

	err := utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+zipLocationPath, err)
	}

	var ThrottlingPolicyListResponse utils.DenyExportThrottlePolicyList
	fmt.Println(resp)
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
}

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

//////////////////////////////////////////////////
func WriteThrottlingPolicy(zipLocationPath string, zipFilename string, marshaledData []byte) (string, error) {
	TempJsonFile := zipFilename
	TempJsonFile += ".json"
	tmpDir, err := ioutil.TempDir("", "apim")
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}

	tempFile := filepath.Join(tmpDir, TempJsonFile)

	err = ioutil.WriteFile(tempFile, marshaledData, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Error writing temp json", err)
	}

	targetZipFile := filepath.Join(zipLocationPath, zipFilename)
	targetZipFile += ".zip"
	err = utils.Zip(tempFile, targetZipFile)
	return targetZipFile, err
}

/////////////////////////////////////////////////////////
