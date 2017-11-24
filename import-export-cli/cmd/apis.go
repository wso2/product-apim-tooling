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
	"github.com/olekukonko/tablewriter"
	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"net/url"
	"os"
)

var listApisCmdEnvironment string
var listApisCmdUsername string
var listApisCmdPassword string
var listApisCmdQuery string
var listApisCmdToken string

// apisCmd related info
const apisCmdLiteral = "apis"
const apisCmdShortDesc = "Display a list of APIs in an environment"

var apisCmdLongDesc = dedent.Dedent(`
		Display a list of APIs in the environment specified by the flag --environment, -e
	`)
var apisCmdExamples = dedent.Dedent(`
	` + utils.ProjectName + ` ` + apisCmdLiteral + ` ` + listCmdLiteral + ` -e dev
	` + utils.ProjectName + ` ` + apisCmdLiteral + ` ` + listCmdLiteral + ` -e staging
	`)

// apisCmd represents the apis command
var apisCmd = &cobra.Command{
	Use:   apisCmdLiteral,
	Short: apisCmdShortDesc,
	Long:  apisCmdLongDesc + apisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + listCmdLiteral + " " + apisCmdLiteral + " called")

		if exportAPICmdToken != "" {
			// token provided with --token (-t) flag
			if exportAPICmdUsername != "" || exportAPICmdPassword != "" {
				// username and/or password provided with -u and/or -p flags
				// Error
				utils.HandleErrorAndExit("username/password provided with OAuth token.", nil)
			} else {
				// token only, proceed with token
			}
		} else {
			// no token provided with --token (-t) flag
			// proceed with username and password
			accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(listApisCmdEnvironment, listApisCmdUsername,
				listApisCmdPassword, utils.MainConfigFilePath, utils.EnvKeysAllFilePath)

			if preCommandErr == nil {
				if listApisCmdQuery != "" {
					fmt.Println("Search query:", listApisCmdQuery)
				}
				count, apis, err := GetAPIList(listApisCmdQuery, accessToken, apiManagerEndpoint)

				if err == nil {
					// Printing the list of available APIs
					fmt.Println("Environment:", listApisCmdEnvironment)
					fmt.Println("No. of APIs:", count)
					if count > 0 {
						printAPIs(apis)
					}
				} else {
					utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
				}
			} else {
				utils.HandleErrorAndExit("Error calling '"+listCmdLiteral+" "+apisCmdLiteral+"'", preCommandErr)
			}
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
	// append '/' to the end if there isn't one already
	if apiManagerEndpoint != "" && string(apiManagerEndpoint[len(apiManagerEndpoint)-1]) != "/" {
		apiManagerEndpoint += "/"
	}

	// url path encode query
	encodedQuery := (&url.URL{Path: query}).String()

	finalUrl := apiManagerEndpoint + "apis?query=" + encodedQuery

	utils.Logln(utils.LogPrefixInfo+"URL:", finalUrl)

	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	resp, err := utils.InvokeGETRequest(finalUrl, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+finalUrl, err)
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

	var data [][]string

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

	apisCmd.Flags().StringVarP(&listApisCmdEnvironment, "environment", "e",
		utils.GetDefaultEnvironment(utils.MainConfigFilePath), "Environment to be searched")
	apisCmd.Flags().StringVarP(&listApisCmdQuery, "query", "q", "",
		"(Optional) query to search for APIs")
	apisCmd.Flags().StringVarP(&listApisCmdToken, "token", "t", "",
		"OAuth token to be used instead of username and password")
	apisCmd.Flags().StringVarP(&listApisCmdUsername, "username", "u", "", "Username")
	apisCmd.Flags().StringVarP(&listApisCmdPassword, "password", "p", "", "Password")
}
