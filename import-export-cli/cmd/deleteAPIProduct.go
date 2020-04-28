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

var deleteAPIProductEnvironment string
var deleteAPIProductName string
var deleteAPIProductProvider string

// DeleteAPIProduct command related usage info
const deleteAPIProductCmdLiteral = "delete-api-product"
const deleteAPIProductCmdShortDesc = "Delete API Product"
const deleteAPIProductCmdLongDesc = "Delete an API Product from an environment"

const deleteAPIProductCmdExamples = utils.ProjectName + ` ` + deleteAPIProductCmdLiteral + ` -n TwitterAPI -r admin -e dev
` + utils.ProjectName + ` ` + deleteAPIProductCmdLiteral + ` -n FacebookAPI -v 2.1.0 -e production
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory.
If the --provider (-r) is not specified, the logged-in user will be considered as the provider.`

// TODO Introduce a version flag and mandate it when the versioning support has been implemented for API Products

// DeleteAPIProductCmd represents the delete-api-product command
var DeleteAPIProductCmd = &cobra.Command{
	Use: deleteAPIProductCmdLiteral + " (--name <name-of-the-api-product> --provider <provider-of-the-api-product> --environment " +
		"<environment-from-which-the-api-product-should-be-deleted>)",
	Short:   deleteAPIProductCmdShortDesc,
	Long:    deleteAPIProductCmdLongDesc,
	Example: deleteAPIProductCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAPIProductCmdLiteral + " called")
		cred, err := getCredentials(deleteAPIProductEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeDeleteAPIProductCmd(cred)
	},
}

// executeDeleteAPIProductCmd executes the delete api command
func executeDeleteAPIProductCmd(credential credentials.Credential) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteAPIProductEnvironment)
	if preCommandErr == nil {
		deleteAPIProductEndpoint := utils.GetApiProductListEndpointOfEnv(deleteAPIProductEnvironment, utils.MainConfigFilePath)
		resp, err := getDeleteAPIProductResponse(deleteAPIProductEndpoint, accessToken, credential)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting API Product ", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			// 200 OK
			fmt.Println(deleteAPIProductName + " API Product deleted successfully!")
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// Neither 200 nor 500
			fmt.Println("Error deleting API Product:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// Error deleting API Product
		fmt.Println("Error getting OAuth tokens while deleting API Product:" + preCommandErr.Error())
	}
}

// getDeleteAPIProductResponse
// @param deleteEndpoint : API Manager Publisher REST API Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func getDeleteAPIProductResponse(deleteEndpoint, accessToken string, credential credentials.Credential) (*resty.Response, error) {
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)
	apiProductId, err := getAPIProductId(accessToken, credential)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting API Product Id for deletion ", err)
	}
	url := deleteEndpoint + apiProductId
	utils.Logln(utils.LogPrefixInfo+"DeleteAPIProduct: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeDELETERequest(url, headers)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the ID of an API Product if available
// @param accessToken : Access token to call the Publisher Rest API
// @return apiId, error
func getAPIProductId(accessToken string, credential credentials.Credential) (string, error) {
	// Unified Search endpoint from the config file to search API Products
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(deleteAPIProductEnvironment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	var queryVal string
	if deleteAPIProductProvider == "" {
		deleteAPIProductProvider = credential.Username
	}
	// TODO Search by version as well when the versioning support has been implemented for API Products
	queryVal = "type:\"APIProduct\" name:\"" + deleteAPIProductName + "\" provider:\"" + deleteAPIProductProvider + "\""
	resp, err := utils.InvokeGETRequestWithQueryParam("query", queryVal, unifiedSearchEndpoint, headers)
	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		apiData := &utils.ApiSearch{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &apiData)
		if apiData.Count != 0 {
			apiProductId := apiData.List[0].ID
			return apiProductId, err
		}
		// TODO Print the version as well when the versioning support has been implemented for API Products
		return "", errors.New("Requested API Product is not available in the Publisher. API Product: " + deleteAPIProductName +
			" Provider: " + deleteAPIProductProvider)
	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching API Product: " + deleteAPIProductName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching API Products. Status: " + resp.Status())
	}
}

// Init using Cobra
func init() {
	RootCmd.AddCommand(DeleteAPIProductCmd)
	DeleteAPIProductCmd.Flags().StringVarP(&deleteAPIProductName, "name", "n", "",
		"Name of the API to be deleted")
	DeleteAPIProductCmd.Flags().StringVarP(&deleteAPIProductProvider, "provider", "r", "",
		"Provider of the API")
	DeleteAPIProductCmd.Flags().StringVarP(&deleteAPIProductEnvironment, "environment", "e",
		"", "Environment from which the API should be deleted")
	// Mark required flags
	_ = DeleteAPIProductCmd.MarkFlagRequired("name")
	_ = DeleteAPIProductCmd.MarkFlagRequired("environment")
}
