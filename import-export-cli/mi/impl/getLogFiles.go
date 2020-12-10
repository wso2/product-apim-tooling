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
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	defaultLogFileListTableFormat = "table {{.FileName}}\t{{.Size}}"
)

// GetLogFileList returns a list of log files created by the micro integrator in a given environment
func GetLogFileList(env string) (*artifactutils.LogFileList, error) {
	resp, err := getArtifactList(utils.MiManagementLogResource, env, &artifactutils.LogFileList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.LogFileList), nil
}

// PrintLogFileList print a list of log file names and sizes according to the given format
func PrintLogFileList(logFileList *artifactutils.LogFileList, format string) {

	if logFileList.Count > 0 {

		logFiles := logFileList.LogFiles

		logFileListContext := getContextWithFormat(format, defaultLogFileListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, logFile := range logFiles {
				if err := t.Execute(w, logFile); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}

		logFileListTableHeaders := map[string]string{
			"FileName": nameHeader,
			"Size":     sizeHeader,
		}

		if err := logFileListContext.Write(renderer, logFileListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Log Files found")
	}
}

// FilterOnlyLogFiles filter the files and return only a list of log files that has .log suffix
func FilterOnlyLogFiles(logFileList *artifactutils.LogFileList) *artifactutils.LogFileList {
	filteredList := new(artifactutils.LogFileList)
	for _, logFile := range logFileList.LogFiles {
		if strings.HasSuffix(logFile.FileName, ".log") {
			filteredList.LogFiles = append(filteredList.LogFiles, logFile)
		}
	}
	filteredList.Count = int32(len(filteredList.LogFiles))
	return filteredList
}

// GetLogFile downloads the specified log file created by the micro integrator in a given environment as a byte stream
func GetLogFile(env, logFileName string) ([]byte, error) {

	params := make(map[string]string)
	params["file"] = logFileName

	url := utils.GetMIManagementEndpointOfResource(utils.MiManagementLogResource, env, utils.MainConfigFilePath)

	resp, err := downloadLogFileData(url, params, env)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WriteLogFile writes the log file in the specified target directory
func WriteLogFile(logFileData []byte, filePath string) {
	err := ioutil.WriteFile(filePath, logFileData, 0644)
	if err != nil {
		fmt.Println("Error writing the log file", err.Error())
	} else {
		fmt.Println("Log file downloaded to", filePath)
	}
}
