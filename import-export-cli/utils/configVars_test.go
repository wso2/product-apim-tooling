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

func TestSetConfigVarsCorrect(t *testing.T) {
	WriteTestMainConfig()
	err := SetConfigVars(testMainConfigFilePath)

	if err != nil {
		t.Errorf("Error in setting Config Vars: %s\n", err.Error())
	}

	defer os.Remove(testMainConfigFilePath)
}

func TestSetConfigVarsIncorrect1(t *testing.T) {
	testIncorrectMainConfig := new(MainConfig)
	testIncorrectMainConfigFileName := "test_incorrect_main_config.yaml"
	testIncorrectMainConfigFilePath := filepath.Join(ConfigDirPath, testIncorrectMainConfigFileName)
	testIncorrectMainConfig.Config = Config{0, ""}
	WriteConfigFile(testIncorrectMainConfig, testIncorrectMainConfigFilePath)

	err := SetConfigVars(testIncorrectMainConfigFilePath)

	if err == nil {
		t.Errorf("Expected error, got nil\n")
	}

	defer os.Remove(testIncorrectMainConfigFilePath)
}

// test case 2 - negative value for httpRequestTimeout
func TestSetConfigVarsIncorrect2(t *testing.T) {
	testIncorrectMainConfig := new(MainConfig)
	testIncorrectMainConfigFileName := "test_incorrect_main_config.yaml"
	testIncorrectMainConfigFilePath := filepath.Join(ConfigDirPath, testIncorrectMainConfigFileName)
	testIncorrectMainConfig.Config = Config{-10, ""}
	WriteConfigFile(testIncorrectMainConfig, testIncorrectMainConfigFilePath)

	err := SetConfigVars(testIncorrectMainConfigFilePath)

	if err == nil {
		t.Errorf("Expected error, got nil\n")
	}

	defer os.Remove(testIncorrectMainConfigFilePath)
}

// TestIsValid1 - Create new file
func TestIsValid1(t *testing.T) {
	filePath := filepath.Join(CurrentDir, "test.txt")
	IsValid(filePath)
}

// TestIsValid2 - Create new file
func TestIsValid2(t *testing.T) {
	fileName := "test.txt"
	os.Create(fileName)
	filePath := filepath.Join(CurrentDir, fileName)
	IsValid(filePath)
	os.Remove(filePath)
}
