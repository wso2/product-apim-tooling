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

package mg

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
	apiNameHeader    = "NAME"
	apiVersionHeader = "VERSION"
	apiTypeHeader    = "TYPE"
	apiLabelsHeader  = "LABELS"

	defaultAPITableFormat = "table {{.Name}}\t{{.Version}}\t{{.Type}}\t{{.Labels}}"
)

type APIMeta struct {
	Total int               `json:"total"`
	Count int               `json:"count"`
	List  []APIMetaListItem `json:"list"`
}
type APIMetaListItem struct {
	APIName      string   `json:"apiName"`
	VersionParam string   `json:"version"`
	APIType      string   `json:"apiType"`
	LabelsParam  []string `json:"labels"`
}

// Name of the API
func (a APIMetaListItem) Name() string {
	return a.APIName
}

// Version of the API
func (a APIMetaListItem) Version() string {
	return a.VersionParam
}

// Type of the API
func (a APIMetaListItem) Type() string {
	return a.APIType
}

// Labels of the API
func (a APIMetaListItem) Labels() []string {
	return a.LabelsParam
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (a *APIMetaListItem) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

// GetAPIList sends GET request and returns the metadata of APIs
func GetAPIsList(accessToken, apiListEndpoint string, queryParam map[string]string) (
	total int, count int, apis []APIMetaListItem, err error) {
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequestWithMultipleQueryParams(queryParam, apiListEndpoint, headers)

	if err != nil {
		return 0, 0, nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		apiMetaResponse := APIMeta{}
		err := json.Unmarshal([]byte(string(resp.Body())), &apiMetaResponse)
		if err != nil {
			return 0, 0, nil, err
		}
		return apiMetaResponse.Total, apiMetaResponse.Count, apiMetaResponse.List, nil
	}
	return 0, 0, nil, errors.New(string(resp.Body()))
}

// PrintAPIs will print an array of APIs as a table
func PrintAPIs(apis []APIMetaListItem) {
	// create api context with standard output
	apiContext := formatter.NewContext(os.Stdout, defaultAPITableFormat)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apis {
			if err := t.Execute(w, a); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	apiTableHeaders := map[string]string{
		"Name":    apiNameHeader,
		"Version": apiVersionHeader,
		"Type":    apiTypeHeader,
		"Labels":  apiLabelsHeader,
	}

	// execute context
	if err := apiContext.Write(renderer, apiTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
