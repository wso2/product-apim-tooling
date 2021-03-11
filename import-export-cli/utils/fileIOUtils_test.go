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
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testKeysFileName = "test_keys_config.yaml"
const testMainConfigFileName = "test_main_config.yaml"

var testKeysFilePath = filepath.Join(ConfigDirPath, testKeysFileName)
var testMainConfigFilePath = filepath.Join(ConfigDirPath, testMainConfigFileName)

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

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func initSampleKeys() {
	envKeysAll.Environments = make(map[string]EnvKeys)
	devEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(devPassword)), "dev_client_secret")
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa_client_secret")
	envKeysAll.Environments[devName] = EnvKeys{"dev_client_id", devEncryptedClientSecret, devUsername}
	envKeysAll.Environments[qaName] = EnvKeys{"qa_client_id", qaEncryptedClientSecret, qaUsername}
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

	defer os.Remove(testKeysFilePath)
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

	defer os.Remove(testKeysFilePath)
}
*/

// test case1 - correct keys
func TestEnvKeysAll_ParseEnvKeysFromFile1(t *testing.T) {
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

// test case 2 - incorrect keys (blank clientID)
func TestEnvKeysAll_ParseEnvKeysFromFile2(t *testing.T) {
	envKeysAll := new(EnvKeysAll)
	testKeysFileName := "test_env_keys_all.yaml"
	testKeysFilePath := filepath.Join(CurrentDir, testKeysFileName)

	// write incorrect keys
	envKeysAll.Environments = make(map[string]EnvKeys)
	qaEncryptedClientSecret := Encrypt([]byte(GetMD5Hash(qaPassword)), "qa_client_secret")
	envKeysAll.Environments[qaName] = EnvKeys{"", qaEncryptedClientSecret, qaUsername}
	WriteConfigFile(envKeysAll, testKeysFilePath)

	data, _ := ioutil.ReadFile(testKeysFilePath)
	err := envKeysAll.ParseEnvKeysFromFile(data)
	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "error", err)
	}

	defer os.Remove(testKeysFilePath)
}

// test case 3 - incorrect keys (blank clientSecret)
func TestEnvKeysAll_ParseEnvKeysFromFile3(t *testing.T) {
	envKeysAll := new(EnvKeysAll)
	testKeysFileName := "test_env_keys_all.yaml"
	testKeysFilePath := filepath.Join(CurrentDir, testKeysFileName)

	// write incorrect keys
	envKeysAll.Environments = make(map[string]EnvKeys)
	envKeysAll.Environments[qaName] = EnvKeys{"client-id", "", qaUsername}
	WriteConfigFile(envKeysAll, testKeysFilePath)

	data, _ := ioutil.ReadFile(testKeysFilePath)
	err := envKeysAll.ParseEnvKeysFromFile(data)
	if err == nil {
		t.Errorf("Expected '%s', got '%s' instead\n", "error", err)
	}

	defer os.Remove(testKeysFilePath)
}

// test case 1 - for a file that does not exist
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

// test case 1 - for directory that does not exist
func TestIsDirExist1(t *testing.T) {
	isDirExist, _ := IsDirExists("random-string")
	if isDirExist {
		t.Errorf("Expected '%t' for a directoroy that does not exist, got '%t' instead\n", false, true)
	}
}

// test case 2 - for a directory that does exist
func TestIsDirExist2(t *testing.T) {
	tempDirName := "tempDir"
	os.Mkdir(tempDirName, os.ModePerm)
	isDirExist, _ := IsDirExists(tempDirName)
	if !isDirExist {
		t.Errorf("Expected '%t' for a directory that does exist, got '%t' instead\n", true, false)
	}

	defer os.Remove(tempDirName)
}

func TestCopyFileNotExists(t *testing.T) {
	tmpFile, err := ioutil.TempFile("testdata", "")
	assert.Nil(t, err, "Should be able to create a temp file")

	err = CopyFile("notexists", tmpFile.Name())
	assert.Error(t, err, "Should return error when copying non existing file")

	// delete temp file
	_ = os.Remove(tmpFile.Name())
}

func TestCopyFileExists(t *testing.T) {
	tmpFile, err := ioutil.TempFile("testdata", "")
	assert.Nil(t, err, "Should be able to create a temp file")

	err = CopyFile("testdata/api.json", tmpFile.Name())
	assert.Nil(t, err, "Should return no error when copying existing file")

	original, err := ioutil.ReadFile("testdata/api.json")
	assert.Nil(t, err, "Should be able to read original file")

	copied, err := ioutil.ReadFile(tmpFile.Name())
	assert.Nil(t, err, "Should be able to read copied file")
	assert.ElementsMatch(t, original, copied, "should be same content")

	// delete temp file
	_ = os.Remove(tmpFile.Name())
}

func TestCopyDirNotExists(t *testing.T) {
	tmpDir, err := ioutil.TempDir("testdata", "")
	assert.Nil(t, err, "Should be able to create a temp directory")

	err = CopyDir("notexists", tmpDir)
	assert.Error(t, err, "Should return error when copying non existing directory")

	// delete temp file
	_ = os.Remove(tmpDir)
}

func TestCopyDirTargetExists(t *testing.T) {
	tmpDir, err := ioutil.TempDir("testdata", "")
	assert.Nil(t, err, "Should be able to create a temp directory")

	err = CopyDir("testdata/fakedir", tmpDir)
	assert.Error(t, err, "Should return error when target exists")

	// delete temp file
	_ = os.Remove(tmpDir)
}

func TestCopyDirTargetNotExists(t *testing.T) {
	tmpDir := "tmp"

	err := CopyDir("testdata/fakedir", tmpDir)
	assert.Nil(t, err, "Should return no error when target not exists")

	assert.True(t, fileExists(path.Join(tmpDir, "BAR")), "BAR should exist in tmpdir")
	assert.True(t, fileExists(path.Join(tmpDir, "subdir", "FOO")), "FOO should exist in tmpdir")

	// delete temp file
	_ = os.RemoveAll(tmpDir)
}
