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
	"fmt"

	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"github.com/olekukonko/tablewriter"
	"os"
	"net/http"
	"github.com/renstrom/dedent"
)

var listEnvironment string
var listCmdUsername string
var listCmdPassword string

// List command related usage Info

const listCmdLiteral string = "list"
const listCmdShortDesc string = "List APIs in an environment"

var listCmdLongDesc string = dedent.Dedent(`
			Display a list containing all the APIs available in the environment specified by flag (--environment, -e)
	`)

var listCmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + listCmdLiteral + ` -e dev
		` + utils.ProjectName + ` ` + listCmdLiteral + ` -e staging
	`)



// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   listCmdLiteral,
	Short: listCmdShortDesc,
	Long:  listCmdLongDesc + listCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listCmdLiteral + " called")

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
			utils.Logln(utils.LogPrefixError + "calling 'list' " + preCommandErr.Error())
			fmt.Println("Error calling 'list'", preCommandErr.Error())
		}
	},
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

// init using Cobra
func init() {
	RootCmd.AddCommand(ListCmd)
	ListCmd.Flags().StringVarP(&listEnvironment, "environment", "e",
		utils.GetDefaultEnvironment(utils.MainConfigFilePath), "Environment to be searched")
	ListCmd.Flags().StringVarP(&listCmdUsername, "username", "u", "", "Username")
	ListCmd.Flags().StringVarP(&listCmdPassword, "password", "p", "", "Password")
}
