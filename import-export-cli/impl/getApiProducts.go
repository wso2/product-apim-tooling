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
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	apiProductIdHeader       = "ID"
	apiProductNameHeader     = "NAME"
	apiProductContextHeader  = "CONTEXT"
	apiProductVersionHeader  = "VERSION"
	apiProductProviderHeader = "PROVIDER"
	apiProductStatusHeader   = "STATUS"

	defaultApiProductTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Context}}\t{{.Version}}\t{{.LifeCycleStatus}}\t{{.Provider}}"
)

// apiProduct holds information about an API Product for outputting
type apiProduct struct {
	id              string
	name            string
	context         string
	version         string
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

// Version of API Product
func (a apiProduct) Version() string {
	return a.version
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

// PrintAPIProducts
func PrintAPIProducts(apiProducts []utils.APIProduct, format string) {
	if format == "" {
		format = defaultApiProductTableFormat
	} else if format == utils.JsonArrayFormatType {
		utils.ListArtifactsInJsonArrayFormat(apiProducts, utils.ProjectTypeApiProduct)
		return
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
		"Version":         apiProductVersionHeader,
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
	return &apiProduct{a.ID, a.Name, a.Context, a.Version, a.Provider, a.LifeCycleStatus}
}
