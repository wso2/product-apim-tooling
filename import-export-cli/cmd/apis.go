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
	"github.com/olekukonko/tablewriter"
	"github.com/renstrom/dedent"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"net/http"
	"os"
)

var listApisCmdEnvironment string
var listApisCmdUsername string
var listApisCmdPassword string

// apisCmd related info
const apisCmdLiteral = "apis"
const apisCmdShortDesc = "Display a list of APIs in an environment"

var apisCmdLongDesc = dedent.Dedent(`
		Display a list of APIs in the environment specified by the flag --environment, -e
	`)
var apisCmdExamples = dedent.Dedent(`
	` + utils.ProjectName + ` ` + apisCmdLiteral + ` ` + listCmdLiteral + ` -e dev
	` + utils.ProjectName + ` ` + apisCmdLiteral + ` ` + listCmdLiteral + ` -e dev -q version:1.0.0
	` + utils.ProjectName + ` ` + apisCmdLiteral + ` ` + listCmdLiteral + ` -e prod -q provider:admin
	` + utils.ProjectName + ` ` + apisCmdLiteral + ` ` + listCmdLiteral + ` -e staging -u admin -p admin
	`)

// apisCmd represents the apis command
var apisCmd = &cobra.Command{
	Use:   apisCmdLiteral,
	Short: apisCmdShortDesc,
	Long:  apisCmdLongDesc + apisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apisCmdLiteral + " called")
		executeApisCmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath)
	},
}

func executeApisCmd(mainConfigFilePath, envKeysAllFilePath string) {
	accessToken, preCommandErr :=
		utils.ExecutePreCommandWithOAuth(listApisCmdEnvironment, listApisCmdUsername, listApisCmdPassword,
			mainConfigFilePath, envKeysAllFilePath)

	if preCommandErr == nil {
		apiListEndpoint := utils.GetApiListEndpointOfEnv(listApisCmdEnvironment, mainConfigFilePath)
		count, apis, err := GetAPIList("", accessToken, apiListEndpoint)

		if err == nil {
			printAPIs(apis)
			// Printing the list of available APIs
			if count == 0 {
				utils.Logln(utils.LogPrefixWarning + "No APIs found in environment - " + listApisCmdEnvironment)
			}
		} else {
			utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
		}
	} else {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + preCommandErr.Error())
		utils.HandleErrorAndExit("Error calling '"+apisCmdLiteral+"'", preCommandErr)
	}
}

// GetAPIList
// @param query : string to be matched against the API names
// @param accessToken : Access Token for the environment
// @param apiManagerEndpoint : API Manager Endpoint for the environment
// @return count (no. of APIs)
// @return array of API objects
// @return error
func GetAPIList(query, accessToken, apiListEndpoint string) (count int32, apis []utils.API, err error) {
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	if query != "" {
		apiListEndpoint += "?query=" + query
	}
	utils.Logln(utils.LogPrefixInfo+"URL:", apiListEndpoint)
	resp, err := utils.InvokeGETRequest(apiListEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+apiListEndpoint, err)
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
		return 0, nil, errors.New(cast.ToString(resp.Body()))
	}

}

// printAPIs
// @param apis : array of API objects
func printAPIs(apis []utils.API) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Version", "Context", "Status", "Provider", "ID"})

	var data [][]string

	for _, api := range apis {
		data = append(data, []string{api.Name, api.Version, api.Context, api.Status, api.Provider, api.ID})
	}

	for _, v := range data {
		table.Append(v)
	}

	// Change table lines
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	table.SetBorder(false)
	table.Render() // Send output
}

func init() {
	ListCmd.AddCommand(apisCmd)

	apisCmd.Flags().StringVarP(&listApisCmdEnvironment, "environment", "e",
		utils.DefaultEnvironmentName, "Environment to be searched")
	apisCmd.Flags().StringVarP(&listApisCmdUsername, "username", "u", "", "Username")
	apisCmd.Flags().StringVarP(&listApisCmdPassword, "password", "p", "", "Password")
}
