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
	apiProductIdHeader       = "ID"
	apiProductNameHeader     = "NAME"
	apiProductContextHeader  = "CONTEXT"
	apiProductProviderHeader = "PROVIDER"
	apiProductStatusHeader   = "STATUS"

	defaultApiProductTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Context}}\t{{.LifeCycleStatus}}\t{{.Provider}}"
)

// apiProduct holds information about an API Product for outputting
type apiProduct struct {
	id              string
	name            string
	context         string
	provider        string
	lifeCycleStatus string
}

// Id of API Product
func (a apiProduct) Id() string {
	return a.id
}

// Name of API Product
func (a apiProduct) Name() string {
	return a.name
}

// Context of API Product
func (a apiProduct) Context() string {
	return a.context
}

// Lifecycle Status of API Product
func (a apiProduct) LifeCycleStatus() string {
	return a.lifeCycleStatus
}

// Provider of API Product
func (a apiProduct) Provider() string {
	return a.provider
}

// MarshalJSON marshals apiProduct using custom marshaller which uses methods instead of fields
func (a *apiProduct) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

// GetAPIProductListFromEnv
// @param accessToken : Access Token for the environment
// @param environment : Environment where API Product should be imported to
// @param query : String to be matched against the API Product names
// @param limit : Total number of API Products to return
// @return count (no. of API Products)
// @return array of API Product objects
// @return error
func GetAPIProductListFromEnv(accessToken, environment, query, limit string) (count int32, apiProducts []utils.APIProduct, err error) {
	unifiedSearchEndpoint := utils.GetUnifiedSearchEndpointOfEnv(environment, utils.MainConfigFilePath)
	return GetAPIProductList(accessToken, unifiedSearchEndpoint, query, limit)
}

// GetAPIProductList
// @param accessToken : Access Token for the environment
// @param unifiedSearchEndpoint : Unified Search Endpoint for the environment to retreive API Product list
// @param query : String to be matched against the API Product names
// @return count (no. of API Products)
// @return array of API Product objects
// @return error
func GetAPIProductList(accessToken, unifiedSearchEndpoint, query, limit string) (count int32, apiProducts []utils.APIProduct, err error) {
	// Unified Search endpoint from the config file to search API Products
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken

	// To filter API Products from unified search
	unifiedSearchEndpoint += "?query=type:\"" + utils.DefaultApiProductType + "\""

	// Setting up the query parameter and limit parameter
	if query != "" {
		unifiedSearchEndpoint += " " + query
	}
	if limit != "" {
		unifiedSearchEndpoint += "&limit=" + limit
	}
	utils.Logln(utils.LogPrefixInfo+"URL:", unifiedSearchEndpoint)
	resp, err := utils.InvokeGETRequest(unifiedSearchEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+unifiedSearchEndpoint, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiProductListResponse := &utils.APIProductListResponse{}
		unmarshalError := json.Unmarshal([]byte(resp.Body()), &apiProductListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid JSON response", unmarshalError)
		}
		return apiProductListResponse.Count, apiProductListResponse.List, nil
	} else {
		return 0, nil, errors.New(string(resp.Body()))
	}
}

// PrintAPIProducts
func PrintAPIProducts(apiProducts []utils.APIProduct, format string) {
	if format == "" {
		format = defaultApiProductTableFormat
	}
	// create API Product context with standard output
	apiProductContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apiProducts {
			if err := t.Execute(w, newApiProductDefinitionFromAPI(a)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	apiProductTableHeaders := map[string]string{
		"Id":              apiProductIdHeader,
		"Name":            apiProductNameHeader,
		"Context":         apiProductContextHeader,
		"LifeCycleStatus": apiProductStatusHeader,
		"Provider":        apiProductProviderHeader,
	}

	// execute context
	if err := apiProductContext.Write(renderer, apiProductTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}

// creates a new API Product from utils.API
func newApiProductDefinitionFromAPI(a utils.APIProduct) *apiProduct {
	return &apiProduct{a.ID, a.Name, a.Context, a.Provider, a.LifeCycleStatus}
}
