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
	defaultdataServiceListTableFormat = "table {{.ServiceName}}\t{{.Wsdl11}}\t{{.Wsdl20}}"
	defaultdataServiceDetailedFormat  = "detail Name - {{.ServiceName}}\n" +
		"Group Name - {{.ServiceGroupName}}\n" +
		"Description - {{.ServiceDescription}}\n" +
		"WSDL 1.1 - {{.Wsdl11}}\n" +
		"WSDL 2.0 - {{.Wsdl20}}\n" +
		"Queries :\n" +
		"ID\tNAMESPACE\n" +
		"{{range .Queries}}{{.Id}}\t{{.Namespace}}\n{{end}}"
)

// GetDataServiceList returns a list of data services deployed in the micro integrator in a given environment
func GetDataServiceList(env string) (*artifactutils.DataServicesList, error) {

	resp, err := getArtifactList(utils.MiManagementDataServiceResource, env, &artifactutils.DataServicesList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.DataServicesList), nil
}

// PrintDataServiceList print a list of data services according to the given format
func PrintDataServiceList(dataServiceList *artifactutils.DataServicesList, format string) {

	if dataServiceList.Count > 0 {

		dataServices := dataServiceList.List

		dataserviceListContext := getContextWithFormat(format, defaultdataServiceListTableFormat)

		// create a new renderer function which iterate collection
		renderer := func(w io.Writer, t *template.Template) error {
			for _, dataservice := range dataServices {
				if err := t.Execute(w, dataservice); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}

		// headers for table
		dataserviceListTableHeaders := map[string]string{
			"Name":   nameHeader,
			"Wsdl11": wsdl11Header,
			"Wsdl20": wsdl20Header,
		}

		// execute context
		if err := dataserviceListContext.Write(renderer, dataserviceListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Data Services found")
	}
}

// GetDataService returns information about a specific data service deployed in the micro integrator in a given environment
func GetDataService(env, dataserviceName string) (*artifactutils.DataServiceInfo, error) {

	resp, err := getArtifactInfo(utils.MiManagementDataServiceResource, "dataServiceName", dataserviceName, env, &artifactutils.DataServiceInfo{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.DataServiceInfo), nil
}

// PrintDataServiceDetails prints details about a data service according to the given format
func PrintDataServiceDetails(ds *artifactutils.DataServiceInfo, format string) {

	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultdataServiceDetailedFormat
	}

	dataserviceContext := formatter.NewContext(os.Stdout, format)

	renderer := getItemRenderer(ds)

	if err := dataserviceContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
