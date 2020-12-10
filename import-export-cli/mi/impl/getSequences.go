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
	defaultSequenceListTableFormat = "table {{.Name}}\t{{.Stats}}\t{{.Tracing}}"
	defaultSequenceDetailedFormat  = "detail Name - {{.Name}}\n" +
		"Container - {{.Container}}\n" +
		"Stats - {{.Stats}}\n" +
		"Tracing - {{.Tracing}}\n" +
		"Mediators - " +
		"{{range $index, $mediator := .Mediators}}" +
		"{{if $index}}, {{end}}" +
		"{{$mediator}}" +
		"{{end}}"
)

// GetSequenceList returns a list of sequences deployed in the micro integrator in a given environment
func GetSequenceList(env string) (*artifactutils.SequenceList, error) {

	resp, err := getArtifactList(utils.MiManagementSequenceResource, env, &artifactutils.SequenceList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.SequenceList), nil
}

// PrintSequenceList print a list of sequences according to the given format
func PrintSequenceList(sequenceList *artifactutils.SequenceList, format string) {

	if sequenceList.Count > 0 {

		sequences := sequenceList.Sequences

		sequenceListContext := getContextWithFormat(format, defaultSequenceListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, sequence := range sequences {
				if err := t.Execute(w, sequence); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}

		sequenceListTableHeaders := map[string]string{
			"Name":    nameHeader,
			"Stats":   statsHeader,
			"Tracing": tracingHeader,
		}

		if err := sequenceListContext.Write(renderer, sequenceListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Sequences found")
	}
}

// GetSequence returns a information about a specific sequence deployed in the micro integrator in a given environment
func GetSequence(env, sequenceName string) (*artifactutils.Sequence, error) {

	resp, err := getArtifactInfo(utils.MiManagementSequenceResource, "sequenceName", sequenceName, env, &artifactutils.Sequence{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.Sequence), nil
}

// PrintSequenceDetails prints details about a sequence according to the given format
func PrintSequenceDetails(sequence *artifactutils.Sequence, format string) {

	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultSequenceDetailedFormat
	}

	sequenceContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRendererEndsWithNewLine(sequence)

	if err := sequenceContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
