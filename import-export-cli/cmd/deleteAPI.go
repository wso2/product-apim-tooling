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

	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"net/http"
)

var deleteAPIEnvironment string
var deleteAPIName string
var deleteAPIVersion string
var deleteAPIProvider string

// DeleteAPI command related usage info
const deleteAPICmdLiteral = "api"
const deleteAPICmdShortDesc = "Delete API"
const deleteAPICmdLongDesc = "Delete an API from an environment"

const deleteAPICmdExamples = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -e production
NOTE: The 3 flags (--name (-n), --version (-v), and --environment (-e)) are mandatory.`

// DeleteAPICmd represents the delete api command
var DeleteAPICmd = &cobra.Command{
	Use:   deleteAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment " +
		"<environment-from-which-the-api-should-be-deleted>)",
	Short: deleteAPICmdShortDesc,
	Long: deleteAPICmdLongDesc,
	Example: deleteAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAPICmdLiteral + " called")
		cred, err := getCredentials(deleteAPIEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeDeleteAPICmd(cred)
	},
}

// executeDeleteAPICmd executes the delete api command
func executeDeleteAPICmd(credential credentials.Credential)  {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteAPIEnvironment)
	if preCommandErr == nil {
		deleteAPIEndpoint := utils.GetApiListEndpointOfEnv(deleteAPIEnvironment, utils.MainConfigFilePath)
		resp, err := getDeleteAPIResponse(deleteAPIEndpoint, accessToken)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting API ", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo + "ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			// 200 OK
			fmt.Println(deleteAPIName + " API deleted successfully!")
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// Neither 200 nor 500
			fmt.Println("Error deleting API:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// Error deleting API
		fmt.Println("Error getting OAuth tokens while deleting API:" + preCommandErr.Error())
	}
}

// getDeleteAPIResponse
// @param deleteEndpoint : API Manager Publisher REST API Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func getDeleteAPIResponse(deleteEndpoint, accessToken string) (*resty.Response, error) {
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)
	apiId, err := getAPIId(accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Id for deletion ", err)
	}
	url := deleteEndpoint + apiId
	utils.Logln(utils.LogPrefixInfo+"DeleteAPI: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeDELETERequest(url, headers)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the ID of an API if available
// @param accessToken : Token to call the Publisher Rest API
// @return apiId, error
func getAPIId(accessToken string) (string, error) {
	// Unified Search endpoint from the config file to search APIs
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(deleteAPIEnvironment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryVal string
	queryVal = "name:\"" + deleteAPIName + "\" version:\"" + deleteAPIVersion + "\""
	if deleteAPIProvider != "" {
		queryVal = queryVal + " provider:\"" + deleteAPIProvider + "\""
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
		if deleteAPIProvider != "" {
			return "", errors.New("Requested API is not available in the Publisher. API: " + deleteAPIName +
				" Version: " + deleteAPIVersion + " Provider: " + deleteAPIProvider)
		}
		return "", errors.New("Requested API is not available in the Publisher. API: " + deleteAPIName +
			" Version: " + deleteAPIVersion)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching API: " + deleteAPIName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching APIs. Status: " + resp.Status())
	}
}

// Init using Cobra
func init() {
	DeleteCmd.AddCommand(DeleteAPICmd)
	DeleteAPICmd.Flags().StringVarP(&deleteAPIName, "name", "n", "",
		"Name of the API to be deleted")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIVersion, "version", "v", "",
		"Version of the API to be deleted")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIProvider, "provider", "r", "",
		"Provider of the API to be deleted")
	DeleteAPICmd.Flags().StringVarP(&deleteAPIEnvironment, "environment", "e",
		"", "Environment from which the API should be deleted")
	// Mark required flags
	_ = DeleteAPICmd.MarkFlagRequired("name")
	_ = DeleteAPICmd.MarkFlagRequired("version")
	_ = DeleteAPICmd.MarkFlagRequired("environment")
}
