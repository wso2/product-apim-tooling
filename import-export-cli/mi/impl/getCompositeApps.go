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
	defaultCompositeAppListTableFormat = "table {{.Name}}\t{{.Version}}"
	defaultCompositeAppDetailedFormat  = "detail Name - {{.Name}}\n" +
		"Version - {{.Version}}\n" +
		"Artifacts :\n" +
		"NAME\tTYPE\n" +
		"{{range .Artifacts}}{{.Name}}\t{{.Type}}\n{{end}}"
)

// GetCompositeAppList returns a list of composite apps deployed in the micro integrator in a given environment
func GetCompositeAppList(env string) (*artifactutils.CompositeAppList, error) {

	resp, err := getArtifactList(utils.MiManagementCarbonAppResource, env, &artifactutils.CompositeAppList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.CompositeAppList), nil
}

// PrintCompositeAppList print a list of composite apps according to the given format
func PrintCompositeAppList(appList *artifactutils.CompositeAppList, format string) {

	if appList.Count > 0 {

		apps := appList.CompositeApps

		appListContext := getContextWithFormat(format, defaultCompositeAppListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, app := range apps {
				if err := t.Execute(w, app); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}

		appListTableHeaders := map[string]string{
			"Name":    nameHeader,
			"Version": versionHeader,
		}

		if err := appListContext.Write(renderer, appListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Composite Apps found")
	}
}

// GetCompositeApp returns a information about a specific composite app deployed in the micro integrator in a given environment
func GetCompositeApp(env, appname string) (*artifactutils.CompositeApp, error) {

	resp, err := getArtifactInfo(utils.MiManagementCarbonAppResource, "carbonAppName", appname, env, &artifactutils.CompositeApp{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.CompositeApp), nil
}

// PrintCompositeAppDetails prints details about a composite app according to the given format
func PrintCompositeAppDetails(app *artifactutils.CompositeApp, format string) {

	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultCompositeAppDetailedFormat
	}

	appContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRenderer(app)

	if err := appContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
