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
	"strings"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	defaultIntegrationAPIListTableFormat = "table {{.Name}}\t{{.Url}}"
	defaultIntegrationAPIDetailedFormat  = "detail Name - {{.Name}}\n" +
		"Version - {{.Version}}\n" +
		"Url - {{.Url}}\n" +
		"Stats - {{.Stats}}\n" +
		"Tracing - {{.Tracing}}\n" +
		"Resources :\n" +
		"URL\tMETHOD\n" +
		"{{range .Resources}}{{.Url}}\t{{.Methods}}\n{{end}}"
)

// GetIntegrationAPIList returns a list of apis deployed in the micro integrator in a given environment
func GetIntegrationAPIList(env string) (*artifactutils.IntegrationAPIList, error) {
	resp, err := getArtifactList(utils.MiManagementAPIResource, env, &artifactutils.IntegrationAPIList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.IntegrationAPIList), nil
}

// PrintIntegrationAPIList print a list of apis according to the given format
func PrintIntegrationAPIList(apiList *artifactutils.IntegrationAPIList, format string) {
	if apiList.Count > 0 {
		apis := apiList.Apis
		apiListContext := getContextWithFormat(format, defaultIntegrationAPIListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, api := range apis {
				if err := t.Execute(w, api); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}
		apiListTableHeaders := map[string]string{
			"Name": nameHeader,
			"Url":  urlHeader,
		}
		if err := apiListContext.Write(renderer, apiListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No APIs found")
	}
}

// GetIntegrationAPI returns a information about a specific api deployed in the micro integrator in a given environment
func GetIntegrationAPI(env, apiName string) (*artifactutils.IntegrationAPI, error) {
	resp, err := getArtifactInfo(utils.MiManagementAPIResource, "apiName", apiName, env, &artifactutils.IntegrationAPI{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.IntegrationAPI), nil
}

// PrintIntegrationAPIDetails prints details about an api according to the given format
func PrintIntegrationAPIDetails(api *artifactutils.IntegrationAPI, format string) {
	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultIntegrationAPIDetailedFormat
	}

	apiContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRenderer(api)

	if err := apiContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
