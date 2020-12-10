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
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	defaultLocalEntryListTableFormat = "table {{.Name}}\t{{.Type}}"
	defaultLocalEntryDetailedFormat  = "detail Name - {{.Name}}\n" +
		"Type - {{.Type}}\n" +
		"Value - {{.Value}}"
)

// GetLocalEntryList returns a list of local entries deployed in the micro integrator in a given environment
func GetLocalEntryList(env string) (*artifactutils.LocalEntryList, error) {

	resp, err := getArtifactList(utils.MiManagementLocalEntrieResource, env, &artifactutils.LocalEntryList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.LocalEntryList), nil
}

// PrintLocalEntryList print a list of local entries according to the given format
func PrintLocalEntryList(localEntryList *artifactutils.LocalEntryList, format string) {

	if localEntryList.Count > 0 {

		localEntrys := localEntryList.LocalEntries

		localEntryListContext := getContextWithFormat(format, defaultLocalEntryListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, localEntry := range localEntrys {
				if err := t.Execute(w, localEntry); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}

		localEntryListTableHeaders := map[string]string{
			"Name": nameHeader,
			"Type": typeHeader,
		}

		if err := localEntryListContext.Write(renderer, localEntryListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Local Entries found")
	}
}

// GetLocalEntry returns a information about a specific local entry deployed in the micro integrator in a given environment
func GetLocalEntry(env, localEntryName string) (*artifactutils.LocalEntryData, error) {

	resp, err := getArtifactInfo(utils.MiManagementLocalEntrieResource, "name", localEntryName, env, &artifactutils.LocalEntryData{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.LocalEntryData), nil
}

// PrintLocalEntryDetails prints details about a local entry according to the given format
func PrintLocalEntryDetails(localEntry *artifactutils.LocalEntryData, format string) {

	localEntryContext := getContextWithFormat(format, defaultLocalEntryDetailedFormat)
	renderer := getItemRendererEndsWithNewLine(localEntry)

	if err := localEntryContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
