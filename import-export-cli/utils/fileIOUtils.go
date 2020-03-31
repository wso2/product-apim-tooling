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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// WriteConfigFile
// @param c : data
// @param envConfigFilePath : Path to file where env endpoints are stored
func WriteConfigFile(c interface{}, configFilePath string) {
	data, err := yaml.Marshal(&c)
	if err != nil {
		HandleErrorAndExit("Unable to write configuration to file.", err)
	}

	err = ioutil.WriteFile(configFilePath, data, 0644)
	if err != nil {
		HandleErrorAndExit("Unable to write configuration to file.", err)
	}
}

// Read and return EnvKeysAll
func GetEnvKeysAllFromFile(envKeysAllFilePath string) *EnvKeysAll {
	data, err := ioutil.ReadFile(envKeysAllFilePath)
	if err != nil {
		fmt.Println("Error reading " + envKeysAllFilePath)
		os.Create(envKeysAllFilePath)
		data, err = ioutil.ReadFile(envKeysAllFilePath)
	}

	var envKeysAll EnvKeysAll
	if err := envKeysAll.ParseEnvKeysFromFile(data); err != nil {
		fmt.Println(LogPrefixError + "parsing " + envKeysAllFilePath)
		return nil
	}

	return &envKeysAll
}

// Read and return MainConfig
func GetMainConfigFromFile(filePath string) *MainConfig {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		HandleErrorAndExit("MainConfig: File Not Found: "+filePath, err)
	}

	var mainConfig MainConfig
	if err := mainConfig.ParseMainConfigFromFile(data); err != nil {
		HandleErrorAndExit("MainConfig: Error parsing "+filePath, err)
	}

	return &mainConfig
}

// Read and return MainConfig. Silently catch the error  when config file is not found
func GetMainConfigFromFileSilently(filePath string) *MainConfig {
	var mainConfig MainConfig
	data, err := ioutil.ReadFile(filePath)
	if err == nil {
		if err := mainConfig.ParseMainConfigFromFile(data); err != nil {
			HandleErrorAndExit("MainConfig: Error parsing "+filePath, err)
		}
	}
	return &mainConfig
}

// Read and validate contents of main_config.yaml
// will throw errors if the any of the lines is blank
func (mainConfig *MainConfig) ParseMainConfigFromFile(data []byte) error {
	if err := yaml.Unmarshal(data, mainConfig); err != nil {
		return err
	}
	for name, endpoints := range mainConfig.Environments {
		if endpoints.ApiManagerEndpoint == "" {
			if endpoints.AdminEndpoint != "" && endpoints.DevPortalEndpoint != "" &&
				endpoints.PublisherEndpoint != "" && endpoints.RegistrationEndpoint != "" &&
				endpoints.TokenEndpoint != "" {
				return nil
			} else {
				return errors.New("Blank API Manager Endpoint for " + name)
			}
		}
		if endpoints.TokenEndpoint == "" {
			return errors.New("Blank Token Endpoint for " + name)
		}
		// ApiImportExportEndpoint is not mandatory
		// ApiListEndpoint is not mandatory
	}
	return nil
}

// Read and validate contents of env_keys_all.yaml
// will throw errors if the any of the lines is blank
func (envKeysAll *EnvKeysAll) ParseEnvKeysFromFile(data []byte) error {
	if err := yaml.Unmarshal(data, envKeysAll); err != nil {
		return err
	}
	for name, keys := range envKeysAll.Environments {
		if keys.ClientID == "" {
			return errors.New("Blank ClientID for " + name)
		}
		if keys.ClientSecret == "" {
			return errors.New("Blank ClientSecret for " + name)
		}
	}
	return nil
}

// Check whether the file exists.
func IsFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			HandleErrorAndExit(fmt.Sprintf(UnableToReadFileMsg, path), err)
		}
	}
	return true
}

// exists returns whether the given file or directory exists or not
func IsDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func CreateDirIfNotExist(path string) (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	return err
}

func CreateDir(path string) (err error) {
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		fmt.Println("Error in creating the directory:" + path + "\n" + err.Error())
	}
	return err
}

func RemoveDirectory(path string) (err error) {
	err = os.RemoveAll(path)
	if err != nil {
		fmt.Println("Error in deleting the directory:" + path + "\n" + err.Error())
	}
	return err
}

// Delete a directory if it exists in the given path
func RemoveDirectoryIfExists(path string) (err error) {
	if exists, err := IsDirExists(path); exists {
		err = os.RemoveAll(path)
		if err != nil {
			fmt.Println("Error in deleting the directory:" + path + "\n" + err.Error())
		}
	}
	return err
}

// Delete a file if it exists in the given path
func RemoveFileIfExists(path string) (err error) {
	if exists := IsFileExist(path); exists {
		err = os.Remove(path)
		if err != nil {
			fmt.Println("Error in deleting the directory:" + path + "\n" + err.Error())
		}
	}
	return err
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// CreateTempFile creates a temporary file in the OS' temp directory
// example pattern "docker-secret-*.yaml"
func CreateTempFile(pattern string, content []byte) (string, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), pattern)
	if err != nil {
		return "", err
	}
	if _, err = tmpFile.Write(content); err != nil {
		return "", err
	}
	// Close the file
	if err := tmpFile.Close(); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
