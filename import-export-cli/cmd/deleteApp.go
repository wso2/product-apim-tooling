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
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"

	"github.com/spf13/cobra"
)

var deleteAppEnvironment string
var deleteAppName string
var deleteAppOwner string

// DeleteApp command related usage info
const deleteAppCmdLiteral = "app"
const deleteAppCmdShortDesc = "Delete App"
const deleteAppCmdLongDesc = "Delete an Application from an environment"

const deleteAppCmdExamples = utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAppCmdLiteral + ` -n TestApplication -o admin -e dev
` + utils.ProjectName + ` ` + deleteCmdLiteral + ` ` + deleteAppCmdLiteral + ` -n SampleApplication -e production
NOTE: Both the flags (--name (-n), and --environment (-e)) are mandatory and the flag --owner (-o) is optional.`

// DeleteAppCmd represents the delete app command
var DeleteAppCmd = &cobra.Command{
	Use:   deleteAppCmdLiteral + " (--name <name-of-the-application> --owner <owner-of-the-application> --environment " +
		"<environment-from-which-the-application-should-be-deleted>)",
	Short: deleteAppCmdShortDesc,
	Long: deleteAppCmdLongDesc,
	Example: deleteAppCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + deleteAppCmdLiteral + " called")
		cred, err := getCredentials(deleteAppEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials ", err)
		}
		executeDeleteAppCmd(cred)
	},
}

// executeDeleteAppCmd executes the delete app command
func executeDeleteAppCmd(credential credentials.Credential)  {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, deleteAppEnvironment)
	if preCommandErr == nil {
		deleteAppEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(deleteAppEnvironment, utils.MainConfigFilePath)
		resp, err := getDeleteAppResponse(deleteAppEndpoint, accessToken)
		if err != nil {
			utils.HandleErrorAndExit("Error while deleting Application ", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo + "ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			// 200 OK
			fmt.Println(deleteAppName + " Application deleted successfully!")
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// Neither 200 nor 500
			fmt.Println("Error deleting Application:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// Error deleting Application
		fmt.Println("Error getting OAuth tokens while deleting Application:" + preCommandErr.Error())
	}
}

// getDeleteAppResponse
// @param deleteEndpoint : API Manager Developer Portal REST API Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func getDeleteAppResponse(deleteEndpoint, accessToken string) (*resty.Response, error) {
	deleteEndpoint = utils.AppendSlashToString(deleteEndpoint)
	appId, err := getAppId(accessToken)
	if err != nil {
		utils.HandleErrorAndExit("Error while getting App Id for deletion ", err)
	}
	url := deleteEndpoint + appId
	utils.Logln(utils.LogPrefixInfo+"DeleteApplication: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeDELETERequest(url, headers)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get the ID of an Application if available
// @param accessToken : Token to call the Developer Portal Rest API
// @return appId, error
func getAppId(accessToken string) (string, error) {
	// Application REST API endpoint of the environment from the config file
	applicationEndpoint := utils.GetDevPortalApplicationListEndpointOfEnv(deleteAppEnvironment, utils.MainConfigFilePath)

	// Prepping headers
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequestWithQueryParam("query", deleteAppName, applicationEndpoint, headers)

	if resp.StatusCode() == http.StatusOK || resp.StatusCode() == http.StatusCreated {
		// 200 OK or 201 Created
		appData := &utils.AppList{}
		data := []byte(resp.Body())
		err = json.Unmarshal(data, &appData)
		if appData.Count != 0 {
			appId := appData.List[0].ApplicationID
			return appId, err
		}
		return "", nil

	} else {
		utils.Logf("Error: %s\n", resp.Error())
		utils.Logf("Body: %s\n", resp.Body())
		if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthorized
			return "", fmt.Errorf("Authorization failed while searching CLI application: " + deleteAppName)
		}
		return "", errors.New("Request didn't respond 200 OK for searching existing applications. " +
			"Status: " + resp.Status())
	}
}

// Init using Cobra
func init() {
	DeleteCmd.AddCommand(DeleteAppCmd)
	DeleteAppCmd.Flags().StringVarP(&deleteAppName, "name", "n", "",
		"Name of the Application to be deleted")
	DeleteAppCmd.Flags().StringVarP(&deleteAppOwner, "owner", "o", "",
		"Owner of the Application to be deleted")
	DeleteAppCmd.Flags().StringVarP(&deleteAppEnvironment, "environment", "e",
		"", "Environment from which the Application should be deleted")
	// Mark required flags
	_ = DeleteAppCmd.MarkFlagRequired("name")
	_ = DeleteAppCmd.MarkFlagRequired("environment")
}
