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
	"io"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	apiIdHeader       = "ID"
	apiNameHeader     = "NAME"
	apiContextHeader  = "CONTEXT"
	apiVersionHeader  = "VERSION"
	apiProviderHeader = "PROVIDER"
	apiStatusHeader   = "STATUS"

	defaultApiTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Version}}\t{{.Context}}\t{{.LifeCycleStatus}}\t{{.Provider}}"
)

var listApisCmdEnvironment string
var listApisCmdFormat string
var listApisCmdQuery string
var listApisCmdLimit string
var queryParamAdded bool = false

// apisCmd related info
const apisCmdLiteral = "apis"
const apisCmdShortDesc = "Display a list of APIs in an environment"

const apisCmdLongDesc = `Display a list of APIs in the environment specified by the flag --environment, -e`

var apisCmdExamples = utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e dev -q version:1.0.0
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e prod -q provider:admin
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e prod -l 100
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e staging
NOTE: The flag (--environment (-e)) is mandatory`

// apisCmd represents the apis command
var apisCmd = &cobra.Command{
	Use:     apisCmdLiteral,
	Short:   apisCmdShortDesc,
	Long:    apisCmdLongDesc,
	Example: apisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apisCmdLiteral + " called")
		cred, err := GetCredentials(listApisCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		//Since other flags does not use args[], query flag will own this
		if len(args) != 0 && listApisCmdQuery != "" {
			for _, argument := range args {
				listApisCmdQuery += " " + argument
			}
		}
		executeApisCmd(cred)
	},
}

// api holds information about an API for outputting
type api struct {
	id              string
	name            string
	context         string
	version         string
	provider        string
	lifeCycleStatus string
}

// creates a new api from utils.API
func newApiDefinitionFromAPI(a utils.API) *api {
	return &api{a.ID, a.Name, a.Context, a.Version, a.Provider,
		a.LifeCycleStatus}
}

// Id of api
func (a api) Id() string {
	return a.id
}

// Name of api
func (a api) Name() string {
	return a.name
}

// Context of api
func (a api) Context() string {
	return a.context
}

// Version of api
func (a api) Version() string {
	return a.version
}

// Lifecycle Status of api
func (a api) LifeCycleStatus() string {
	return a.lifeCycleStatus
}

// Provider of api
func (a api) Provider() string {
	return a.provider
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (a *api) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

func executeApisCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, listApisCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+apisCmdLiteral+"'", err)
	}

	apiListEndpoint := utils.GetApiListEndpointOfEnv(listApisCmdEnvironment, utils.MainConfigFilePath)
	_, apis, err := GetAPIList(listApisCmdQuery, listApisCmdLimit, accessToken, apiListEndpoint)
	if err == nil {
		printAPIs(apis, listApisCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
	}
}

// GetAPIList
// @param query : string to be matched against the API names
// @param accessToken : Access Token for the environment
// @param apiManagerEndpoint : API Manager Endpoint for the environment
// @return count (no. of APIs)
// @return array of API objects
// @return error
func GetAPIList(query, limit, accessToken, apiListEndpoint string) (count int32, apis []utils.API, err error) {
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	if query != "" {
		apiListEndpoint += getQueryParamConnector() + "query=" + query
	}
	if limit != "" {
		apiListEndpoint += getQueryParamConnector() + "limit=" + limit
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
		return 0, nil, errors.New(string(resp.Body()))
	}

}

func getQueryParamConnector() (connector string) {
	if queryParamAdded {
		return "&"
	} else {
		queryParamAdded = true
		return "?"
	}
}

// printAPIs
func printAPIs(apis []utils.API, format string) {
	if format == "" {
		format = defaultApiTableFormat
	}
	// create api context with standard output
	apiContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apis {
			if err := t.Execute(w, newApiDefinitionFromAPI(a)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	apiTableHeaders := map[string]string{
		"Id":              apiIdHeader,
		"Name":            apiNameHeader,
		"Context":         apiContextHeader,
		"Version":         apiVersionHeader,
		"LifeCycleStatus": apiStatusHeader,
		"Provider":        apiProviderHeader,
	}

	// execute context
	if err := apiContext.Write(renderer, apiTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}

func init() {
	ListCmd.AddCommand(apisCmd)

	apisCmd.Flags().StringVarP(&listApisCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	apisCmd.Flags().StringVarP(&listApisCmdQuery, "query", "q",
		"", "Query pattern")
	apisCmd.Flags().StringVarP(&listApisCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultApisDisplayLimit), "Maximum number of apis to return")
	apisCmd.Flags().StringVarP(&listApisCmdFormat, "format", "", "", "Pretty-print apis "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = apisCmd.MarkFlagRequired("environment")
}
