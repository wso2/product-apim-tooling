/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package utils

import (
	"archive/zip"
	"io"
)

// ZipFile holds the content and the path of the file inside the zip folder
type ZipFile struct {
	Path    string
	Content string
}

// CreateZipFile creates a zip file using the provided io.Writer.
// It takes a slice of ZipFile structs containing information about the files to be added to the zip.
// Each ZipFile struct specifies the file path within the zip and its content.
// It returns an error if any operation fails.
func CreateZipFile(writer io.Writer, zipFiles []ZipFile) error {
	zipWriter := zip.NewWriter(writer)
	for _, zipFile := range zipFiles {
		fileWriter, err := zipWriter.Create(zipFile.Path)
		if err != nil {
			return err
		}
		_, err = fileWriter.Write([]byte(zipFile.Content))
		if err != nil {
			return err
		}
	}
	return zipWriter.Close()
}
