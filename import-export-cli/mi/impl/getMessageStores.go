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
	defaultMessageStoreListTableFormat = "table {{.Name}}\t{{.Type}}\t{{.Size}}"
	defaultMessageStoreDetailedFormat  = "detail Name - {{.Name}}\n" +
		"File Name - {{.FileName}}\n" +
		"Container - {{.Container}}\n" +
		"Producer - {{.Producer}}\n" +
		"Consumer - {{.Consumer}}\n" +
		"Size - {{.Size}}\n" +
		"Properties :\n" +
		"{{ if eq (len .Properties) 0 }}" +
		"No Properties found\n" +
		"{{else}}" +
		"{{ range $key, $value := .Properties }}" +
		" {{ $key }} = {{ $value }}\n" +
		"{{ end }}" +
		"{{ end }}"
)

// GetMessageStoreList returns a list of message stores deployed in the micro integrator in a given environment
func GetMessageStoreList(env string) (*artifactutils.MessageStoreList, error) {
	resp, err := getArtifactList(utils.MiManagementMessageStoreResource, env, &artifactutils.MessageStoreList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.MessageStoreList), nil
}

// PrintMessageStoreList print a list of message stores according to the given format
func PrintMessageStoreList(messageStoreList *artifactutils.MessageStoreList, format string) {
	if messageStoreList.Count > 0 {
		messageStores := messageStoreList.MessageStores
		messageStoreListContext := getContextWithFormat(format, defaultMessageStoreListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, messageStore := range messageStores {
				if err := t.Execute(w, messageStore); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}
		messageStoreListTableHeaders := map[string]string{
			"Name": nameHeader,
			"Type": typeHeader,
			"Size": sizeHeader,
		}
		if err := messageStoreListContext.Write(renderer, messageStoreListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Message Stores found")
	}
}

// GetMessageStore returns a information about a specific message store deployed in the micro integrator in a given environment
func GetMessageStore(env, messageStoreName string) (*artifactutils.MessageStoreData, error) {
	resp, err := getArtifactInfo(utils.MiManagementMessageStoreResource, "name", messageStoreName, env, &artifactutils.MessageStoreData{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.MessageStoreData), nil
}

// PrintMessageStoreDetails prints details about a message store according to the given format
func PrintMessageStoreDetails(messageStore *artifactutils.MessageStoreData, format string) {
	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultMessageStoreDetailedFormat
	}

	messageStoreContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRenderer(messageStore)

	if err := messageStoreContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
