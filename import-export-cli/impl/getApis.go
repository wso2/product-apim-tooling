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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

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

var queryParamAdded bool = false

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

// GetAPIListFromEnv
// @param accessToken : Access Token for the environment
// @param environment : Environment name to use when getting the API List
// @param query : string to be matched against the API names
// @param limit : total # of results to return
// @return count (no. of APIs)
// @return array of API objects
// @return error
func GetAPIListFromEnv(accessToken, environment, query, limit string) (count int32, apis []utils.API, err error) {
	apiListEndpoint := utils.GetApiListEndpointOfEnv(environment, utils.MainConfigFilePath)
	return GetAPIList(accessToken, apiListEndpoint, query, limit)
}

// GetAPIList
// @param accessToken : Access Token for the environment
// @param apiListEndpoint : API List endpoint
// @param query : string to be matched against the API names
// @param limit : total # of results to return
// @return count (no. of APIs)
// @return array of API objects
// @return error
func GetAPIList(accessToken, apiListEndpoint, query, limit string) (count int32, apis []utils.API, err error) {
	queryParamAdded := false
	getQueryParamConnector := func() (connector string) {
		if queryParamAdded {
			return "&"
		} else {
			queryParamAdded = true
			return "?"
		}
	}

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

// PrintAPIs
func PrintAPIs(apis []utils.API, format string) {
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
