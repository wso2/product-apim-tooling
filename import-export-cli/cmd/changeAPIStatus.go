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

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"

	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var apiStateChangeEnvironment string
var apiNameForStateChange string
var apiVersionForStateChange string
var apiProviderForStateChange string
var apiStateChangeAction string

// ChangeAPIStatus command related usage info
const changeAPIStatusCmdLiteral = "api"
const changeAPIStatusCmdShortDesc = "Change Status of an API"
const changeAPIStatusCmdLongDesc = "Change the lifecycle status of an API in an environment"

const changeAPIStatusCmdExamples = utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeAPIStatusCmdLiteral + ` -a Publish -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + changeStatusCmdLiteral + ` ` + changeAPIStatusCmdLiteral + ` -a Publish -n FacebookAPI -v 2.1.0 -e production
NOTE: The 4 flags (--action (-a), --name (-n), --version (-v), and --environment (-e)) are mandatory.`

// changeAPIStatusCmd represents change-status api command
var ChangeAPIStatusCmd = &cobra.Command{
	Use: changeAPIStatusCmdLiteral + " (--action <action-of-the-api-state-change> --name <name-of-the-api> --version <version-of-the-api> --provider " +
		"<provider-of-the-api> --environment <environment-from-which-the-api-state-should-be-changed>)",
	Short:   changeAPIStatusCmdShortDesc,
	Long:    changeAPIStatusCmdLongDesc,
	Example: changeAPIStatusCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + changeAPIStatusCmdLiteral + " called")
		cred, err := GetCredentials(apiStateChangeEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeChangeAPIStatusCmd(cred)
	},
}

// executeChangeAPIStatusCmd executes the change api status command
func executeChangeAPIStatusCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, apiStateChangeEnvironment)
	if preCommandErr == nil {
		changeAPIStatusEndpoint := utils.GetApiListEndpointOfEnv(apiStateChangeEnvironment, utils.MainConfigFilePath)
		resp, err := getChangeAPIStatusResponse(changeAPIStatusEndpoint, accessToken, credential)
		if err != nil {
			utils.HandleErrorAndExit("Error while changing the API status", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			// 200 OK
			fmt.Println(apiNameForStateChange + " API state changed successfully!")
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// Neither 200 nor 500
			fmt.Println("Error while changing API Status: ", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// Error changing the API status
		fmt.Println("Error getting OAuth tokens while changing status of the API:" + preCommandErr.Error())
	}
}

// getChangeAPIStatusResponse
// @param changeAPIStatusEndpoint : API Manager Publisher REST API Endpoint for the environment
// @param accessToken : Access Token for the resource
// @param credential : Credentials of the logged-in user
// @return response Response in the form of *resty.Response
func getChangeAPIStatusResponse(changeAPIStatusEndpoint, accessToken string, credential credentials.Credential) (*resty.Response, error) {
	changeAPIStatusEndpoint = utils.AppendSlashToString(changeAPIStatusEndpoint)
	apiId, err := getAPIIdForStateChange(accessToken, credential)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id for state change ", err)
	}
	url := changeAPIStatusEndpoint + "change-lifecycle?action=" + apiStateChangeAction + "&apiId=" + apiId
	utils.Logln(utils.LogPrefixInfo+"APIStateChange: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokePOSTRequest(url, headers, "")

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the ID of an API if available
// @param accessToken : Token to call the Publisher Rest API
// @return apiId, error
func getAPIIdForStateChange(accessToken string, credential credentials.Credential) (string, error) {
	// Unified Search endpoint from the config file to search APIs
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(apiStateChangeEnvironment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryVal string
	queryVal = "name:\"" + apiNameForStateChange + "\" version:\"" + apiVersionForStateChange + "\""
	if apiProviderForStateChange != "" {
		queryVal = queryVal + " provider:\"" + apiProviderForStateChange + "\""
	}
	resp, err := utils.InvokeGETRequestWithQueryParam("query", queryVal, unifiedSearchEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		apiData := &utils.ApiSearch{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &apiData)
		if apiData.Count != 0 {
			apiId := apiData.List[0].ID
			return apiId, err
		}
		if apiProviderForStateChange != "" {
			return "", errors.New("Requested API is not available in the Publisher. API: " + apiNameForStateChange +
				" Version: " + apiVersionForStateChange + " Provider: " + apiProviderForStateChange)
		}
		return "", errors.New("Requested API is not available in the Publisher. API: " + apiNameForStateChange +
			" Version: " + apiVersionForStateChange)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching API: " + apiNameForStateChange)
		}
		return "", errors.New("Request didn't respond 200 OK for searching APIs. Status: " + resp.Status())
	}
}

func init() {
	ChangeStatusCmd.AddCommand(ChangeAPIStatusCmd)
	ChangeAPIStatusCmd.Flags().StringVarP(&apiStateChangeAction, "action", "a", "",
		"Action to be taken to change the status of the API")
	ChangeAPIStatusCmd.Flags().StringVarP(&apiNameForStateChange, "name", "n", "",
		"Name of the API to be state changed")
	ChangeAPIStatusCmd.Flags().StringVarP(&apiVersionForStateChange, "version", "v", "",
		"Version of the API to be state changed")
	ChangeAPIStatusCmd.Flags().StringVarP(&apiProviderForStateChange, "provider", "r", "",
		"Provider of the API")
	ChangeAPIStatusCmd.Flags().StringVarP(&apiStateChangeEnvironment, "environment", "e",
		"", "Environment of which the API state should be changed")
	// Mark required flags
	_ = ChangeAPIStatusCmd.MarkFlagRequired("action")
	_ = ChangeAPIStatusCmd.MarkFlagRequired("name")
	_ = ChangeAPIStatusCmd.MarkFlagRequired("version")
	_ = ChangeAPIStatusCmd.MarkFlagRequired("environment")
}
