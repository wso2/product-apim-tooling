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
	defaultMessageProcessorListTableFormat = "table {{.Name}}\t{{.Type}}\t{{.Status}}"
	defaultMessageProcessorDetailedFormat  = "detail Name - {{.Name}}\n" +
		"Type - {{.Type}}\n" +
		"File Name - {{.FileName}}\n" +
		"Message Store - {{.Store}}\n" +
		"Artifact Container - {{.Container}}\n" +
		"Status - {{.Status}}\n" +
		"Parameters :\n" +
		"{{ if eq (len .Parameters) 0 }}" +
		"No Parameters found\n" +
		"{{else}}" +
		"{{ range $key, $value := .Parameters }}" +
		" {{ $key }} = {{ $value }}\n" +
		"{{ end }}" +
		"{{ end }}"
)

// GetMessageProcessorList returns a list of message processors deployed in the micro integrator in a given environment
func GetMessageProcessorList(env string) (*artifactutils.MessageProcessorList, error) {
	resp, err := getArtifactList(utils.MiManagementMessageProcessorResource, env, &artifactutils.MessageProcessorList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.MessageProcessorList), nil
}

// PrintMessageProcessorList print a list of message processors according to the given format
func PrintMessageProcessorList(messageProcessorList *artifactutils.MessageProcessorList, format string) {
	if messageProcessorList.Count > 0 {
		messageProcessors := messageProcessorList.MessageProcessors
		messageProcessorListContext := getContextWithFormat(format, defaultMessageProcessorListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, messageProcessor := range messageProcessors {
				if err := t.Execute(w, messageProcessor); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}
		messageProcessorListTableHeaders := map[string]string{
			"Name":   nameHeader,
			"Type":   typeHeader,
			"Status": statusHeader,
		}
		if err := messageProcessorListContext.Write(renderer, messageProcessorListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Message Processors found")
	}
}

// GetMessageProcessor returns a information about a specific message processor deployed in the micro integrator in a given environment
func GetMessageProcessor(env, messageProcessorName string) (*artifactutils.MessageProcessorData, error) {
	resp, err := getArtifactInfo(utils.MiManagementMessageProcessorResource, "name", messageProcessorName, env, &artifactutils.MessageProcessorData{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.MessageProcessorData), nil
}

// PrintMessageProcessorDetails prints details about a message processor according to the given format
func PrintMessageProcessorDetails(messageProcessor *artifactutils.MessageProcessorData, format string) {
	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultMessageProcessorDetailedFormat
	}

	messageProcessorContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRenderer(messageProcessor)

	if err := messageProcessorContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
