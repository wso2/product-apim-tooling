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

package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip will create an archive from source and store it in target
func Zip(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close() // close the archive when exit

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	fileInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Get base directory if this is a directory
	var baseDir string
	if fileInfo.IsDir() {
		baseDir = filepath.Base(source)
	}

	// Walk through the source path to generate an archive
	// Walk accepts a WalkFn which has signature of func(path string, info os.FileInfo, err error) error
	// Walk will return any error to err while walking
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a partial zip header from current file or directory
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// If baseDir is not empty it means we need to strip source from path, so we can get a relative filename from
		// base.
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}
		if info.IsDir() {
			// add directory to zip archive
			header.Name += "/"
		} else {
			// add a file to zip archive using deflate algorithm
			header.Method = zip.Deflate
		}
		Logln("Creating:", header.Name)

		// Create an archive writer
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		// if this is a directory we don't copy, we only add header
		if info.IsDir() {
			return nil
		}

		// open the file for reading
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// copy contents of the file to the archive
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
