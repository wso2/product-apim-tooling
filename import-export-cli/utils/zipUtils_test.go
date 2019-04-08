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
	"os"
	"path/filepath"
	"testing"
)

func TestZipDirError(t *testing.T) {
	err := Zip("", "")
	if err == nil {
		t.Errorf("Zip() didn't return an error for invalid source and destination")
	}
}

func TestZipDirOK(t *testing.T) {
	directoryName := "wso2apimZipTest"

	directoryPath := filepath.Join(ConfigDirPath, directoryName)
	fileName := "test.txt"
	filePath := filepath.Join(directoryPath, fileName)

	os.Mkdir(directoryPath, os.ModePerm)

	// check if directory exists
	var _, err = os.Stat(directoryPath)
	if err != nil {
		t.Errorf("Error opening directory")
	}

	// create directory if it doesn't already exist
	if os.IsNotExist(err) {
		var file, err = os.Create(directoryPath)
		if err != nil {
			t.Errorf("Error creating sample directory for compressing: %s\n", err)
		}

		defer file.Close()
	}

	// check if file exists
	_, err = os.Stat(filePath)

	// create file if it doesn't already exist
	if os.IsNotExist(err) {
		var file, err = os.Create(filePath)
		if err != nil {
			t.Errorf("Error creating sample file for compressing: %s\n", err)
		}
		defer file.Close()
	}

	// Open file using READ & WRITE permissions
	var file, err1 = os.OpenFile(filePath, os.O_RDWR, 0644)
	if err1 != nil {
		t.Errorf("Error opening sample file: %s\n", err1)
	}
	defer file.Close()

	// Write content to file
	_, err = file.WriteString("abcdefgh\n")
	if err != nil {
		t.Errorf("Error writing content to file: %s\n", err)
	}

	// Save changes
	err = file.Sync()
	if err != nil {
		t.Errorf("Error saving file: %s\n", err)
	}

	zipFile := filepath.Join(directoryPath, "testZip.zip")

	// now try compressing
	err = Zip(directoryPath, zipFile)

	if err != nil {
		t.Errorf("Error compressing directory: %s\n", err)
	}

	// delete file
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error deleting file: %s\n", err)
	}

	// delete zip file
	err = os.Remove(zipFile)
	if err != nil {
		t.Errorf("Error deleting file: %s\n", err)
	}

	// delete directory
	err = os.Remove(directoryPath)
	if err != nil {
		t.Errorf("Error deleting directory: %s\n", err)
	}
}
