/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"github.com/renstrom/dedent"
	"github.com/olekukonko/tablewriter"
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
)

var listEnvironment string
var listCmdUsername string
var listCmdPassword string

// apisCmd related info
const apisCmdLiteral string = "apis"
const apisCmdShortDesc string = "Display a list of APIs in an environment"

var apisCmdLongDesc = dedent.Dedent(`
		Display a list of APIs in the environment specified by the flag --environment, -e
	`)
var apisCmdExamples = dedent.Dedent(`
	`+utils.ProjectName + ` `+ apisCmdLiteral +` `+ listCmdLiteral +` -e dev
	`+utils.ProjectName + ` `+ apisCmdLiteral +` `+ listCmdLiteral +` -e staging
	`)

// apisCmd represents the apis command
var apisCmd = &cobra.Command{
	Use:   apisCmdLiteral,
	Short: apisCmdShortDesc,
	Long:  apisCmdLongDesc + apisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listCmdLiteral + " " + apisCmdLiteral + " called")

		accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(listEnvironment, listCmdUsername,
			listCmdPassword)

		if preCommandErr == nil {
			count, apis, err := GetAPIList("", accessToken, apiManagerEndpoint)

			if err == nil {
				// Printing the list of available APIs
				fmt.Println("Environment:", listEnvironment)
				fmt.Println("No. of APIs:", count)
				if count > 0 {
					printAPIs(apis)
				}
			} else {
				utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
			}
		} else {
			utils.HandleErrorAndExit("Error calling '"+listCmdLiteral + " " + apisCmdLiteral+"'", preCommandErr)
		}
	},
}


// GetAPIList
// @param query : string to be matched against the API names
// @param accessToken : Access Token for the environment
// @param apiManagerEndpoint : API Manager Endpoint for the environment
// @return count (no. of APIs)
// @return array of API objects
// @return error
func GetAPIList(query string, accessToken string, apiManagerEndpoint string) (int32, []utils.API, error) {
	url := apiManagerEndpoint

	// append '/' to the end if there isn't one already
	if url != "" && string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "apis?query=" + query
	utils.Logln(utils.LogPrefixInfo+"URL:", url)

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeGETRequest(url, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+url, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiListResponse := &utils.APIListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}

		return apiListResponse.Count, apiListResponse.List, nil
	} else {
		return 0, nil, errors.New(resp.Status())
	}

}


// printAPIs
// @param apis : array of API objects
func printAPIs(apis []utils.API) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Version", "Context", "LifeCycleStatus", "ID"})

	data := [][]string{}

	for _, api := range apis {
		data = append(data, []string{api.Name, api.Version, api.Context, api.LifeCycleStatus, api.ID})
	}

	for _, v := range data {
		table.Append(v)
	}

	table.Render() // Send output

}

func init() {
	ListCmd.AddCommand(apisCmd)

	apisCmd.Flags().StringVarP(&listEnvironment, "environment", "e",
		utils.GetDefaultEnvironment(utils.MainConfigFilePath), "Environment to be searched")
	apisCmd.Flags().StringVarP(&listCmdUsername, "username", "u", "", "Username")
	apisCmd.Flags().StringVarP(&listCmdPassword, "password", "p", "", "Password")
}
