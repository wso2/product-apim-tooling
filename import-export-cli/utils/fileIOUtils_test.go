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

// +integration
package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const testKeysFileName = "test_keys_config.yaml"
const testMainConfigFileName = "test_main_config.yaml"
var testKeysFilePath = filepath.Join(ApplicationRoot, testKeysFileName)
var testMainConfigFilePath = filepath.Join(ApplicationRoot, testMainConfigFileName)

var envKeysAll = new(EnvKeysAll)
var mainConfig = new(MainConfig)

const devName = "dev"
const qaName = "qa"
const devUsername = "dev_username"
const qaUsername = "qa_username"
const devPassword = "dev_password"
const qaPassword = "qa_password"

// helper function for testing
func writeCorrectKeys() {
	initSampleKeys()
	WriteConfigFile(envKeysAll, testKeysFilePath)
}

func getSampleKeys() *EnvKeysAll {
	initSampleKeys()
	return envKeysAll
}

func getSampleMainConfig() *MainConfig {
	return mainConfig
}

func initSampleKeys() {
	envKeysAll.Environments = make(map[string]EnvKeys)
	devEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(devPassword)), "dev_client_secret")
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa_client_secret")
	envKeysAll.Environments[devName] = EnvKeys{"dev_client_id", devEncryptedClientSecret, devUsername}
	envKeysAll.Environments[qaName] = EnvKeys{"qa_client_id", qaEncryptedClientSecret, qaUsername}
}

// helper function for unit-testing
func WriteCorrectMainConfig() {
	initSampleMainConfig()
	WriteConfigFile(mainConfig, testMainConfigFilePath)
}

func initSampleMainConfig() {
	mainConfig.Config = Config{2500,"/home/exported"}
	mainConfig.Environments = make(map[string]EnvEndpoints)
	mainConfig.Environments[devName] = EnvEndpoints{"dev_apim_endpoint",
		"dev_reg_endpoint", "dev_token_endpoint"}
	mainConfig.Environments[qaName] = EnvEndpoints{"qa_apim_endpoint",
		"qa_reg_endpoint", "dev_token_endpoint"}
}

func TestWriteConfigFile(t *testing.T) {
	writeCorrectKeys()
	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

func TestGetEnvKeysAllFromFile1(t *testing.T) {
	writeCorrectKeys()
	envKeysAllReturned := GetEnvKeysAllFromFile(testKeysFilePath)

	if envKeysAllReturned.Environments[devName] != envKeysAll.Environments[devName] ||
		envKeysAllReturned.Environments[qaName] != envKeysAll.Environments[qaName] {
		t.Errorf("Error in GetEnvKeysAllFromFile()")
	}

	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

/*
func TestGetEnvKeysAllFromFile3(t *testing.T) {
	// testing for incorrect data
	var envIncorrectKeysAll *EnvKeysAll = new(EnvKeysAll)
	envIncorrectKeysAll.Environments  = make(map[string]EnvKeys)
	envIncorrectKeysAll.Environments[devName] = EnvKeys{"dev_client_id", "", devUsername}

	WriteConfigFile(envIncorrectKeysAll, testKeysFilePath)
	envKeysAllReturned := GetEnvKeysAllFromFile(testKeysFilePath)

	if envKeysAllReturned.Environments[devName] != envIncorrectKeysAll.Environments[devName] ||
		envKeysAllReturned.Environments[qaName] != envIncorrectKeysAll.Environments[qaName] {
		t.Errorf("Error in GetEnvKeysAllFromFile()")
	}

	var err = os.Remove(testKeysFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}
*/

func TestGetMainConfigFromFile(t *testing.T) {
	// testing for correct data
	WriteCorrectMainConfig()
	mainConfigReturned := GetMainConfigFromFile(testMainConfigFilePath)

	if mainConfigReturned.Environments[devName] != mainConfig.Environments[devName] ||
		mainConfigReturned.Environments[qaName] != mainConfig.Environments[qaName] {
		t.Errorf("Error in GetMainConfigFromFile()")
	}

	var err = os.Remove(testMainConfigFilePath)
	if err != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

func TestMainConfig_ParseMainConfigFromFile(t *testing.T) {
	var envEndpointsAllLocal MainConfig
	WriteCorrectMainConfig()
	data, err := ioutil.ReadFile(testMainConfigFilePath)
	if err != nil {
		t.Error("Error")
	}
	envEndpointsAllLocal.ParseMainConfigFromFile(data)

	var err1 = os.Remove(testMainConfigFilePath)
	if err1 != nil {
		t.Errorf("Error deleting file " + testMainConfigFilePath)
	}
}

func TestEnvKeysAll_ParseEnvKeysFromFile(t *testing.T) {
	var envKeysAllLocal EnvKeysAll
	writeCorrectKeys()
	data, err := ioutil.ReadFile(testKeysFilePath)
	if err != nil {
		t.Error("Error")
	}
	envKeysAllLocal.ParseEnvKeysFromFile(data)

	var err1 = os.Remove(testKeysFilePath)
	if err1 != nil {
		t.Errorf("Error deleting file " + testKeysFilePath)
	}
}

// test for a file that does not exist
func TestIsFileExist1(t *testing.T) {
	isFileExist := IsFileExist("random-string")
	if isFileExist {
		t.Errorf("Expected '%t' for a file that does not exist, got '%t' instead\n", false, true)
	}
}

// test for a file that does exist
func TestIsFileExist2(t *testing.T) {
	isFileExist := IsFileExist("./fileIOUtils.go")
	if !isFileExist {
		t.Errorf("Expected '%t' for a file that does exist,  got '%t' instead\n", true, false)
	}
}

// test for directory that does not exist
func TestIsDirExist1(t *testing.T) {
	isDirExist, _ := IsDirExist("random-string")
	if isDirExist {
		t.Errorf("Expected '%t' for a directoroy that does not exist, got '%t' instead\n", false, true)
	}
}

// test for a directory that does exist
func TestIsDirExist2(t *testing.T) {
	tempDirName := "tempDir"
	os.Mkdir(tempDirName, os.ModePerm)
	isDirExist, _ := IsDirExist(tempDirName)
	if !isDirExist {
		t.Error("Expected '%t' for a directory that does exist, got '%t' instead\n", true, false)
	}

	defer os.Remove(tempDirName)
}
