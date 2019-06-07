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
	"fmt"
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

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
// returns a slice of extracted files with relative paths(dest is not appended)
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip.
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return filenames, err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
