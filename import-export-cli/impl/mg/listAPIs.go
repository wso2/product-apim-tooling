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
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
)

const (
	apiNameHeader    = "NAME"
	apiVersionHeader = "VERSION"
	apiTypeHeader    = "TYPE"
	apiLabelsHeader  = "LABELS"

	defaultAPITableFormat = "table {{.Name}}\t{{.Version}}\t{{.Type}}\t{{.Labels}}"
)

var queryParamAdded bool = false

// api holds information about an API for outputting
type api struct {
	name    string
	version string
	apiType string
	labels  []string
}

// Name of api
func (a api) Name() string {
	return a.name
}

// Version of api
func (a api) Version() string {
	return a.version
}

// Lifecycle Status of api
func (a api) Type() string {
	return a.apiType
}

// Provider of api
func (a api) Labels() []string {
	return a.labels
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (a *api) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

// PrintAPIs will print an array of APIs as a table
func PrintAPIs(apis []api) {
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
