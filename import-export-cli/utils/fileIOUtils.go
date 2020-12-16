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
	"strings"

	"github.com/go-resty/resty"

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
		if !HasOnlyMIEndpoint(&endpoints) {
			if endpoints.ApiManagerEndpoint == "" {
				if RequiredAPIMEndpointsExists(&endpoints) {
					return nil
				}
				return errors.New("Blank API Manager Endpoint for " + name)
			}
			if endpoints.TokenEndpoint == "" {
				return errors.New("Blank Token Endpoint for " + name)
			}
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
		os.MkdirAll(path, os.ModePerm)
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

// CopyDirectoryContents recursively copies all the content of a directory, attempting to preserve permissions.
// Source directory must exist,and the destination directory exist.
func CopyDirectoryContents (src string, dst string) (err error)  {
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
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}
	return
}

// MoveDirectoryContentsToNewDirectory recursively moves all the content of a directory to a given directory
// attempting to preserve permissions.
// Source directory must exist,and the destination directory exist.
// @param src source directory path
// @param dst destiny directory path
// @return error
func MoveDirectoryContentsToNewDirectory(src string, dst string) (err error)  {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	err = CreateDirIfNotExist(filepath.Join(dst))
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			//Copy directory from source to destination
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
			//remove directory from source after copying
			err = RemoveDirectoryIfExists(srcPath)
			if err != nil {
				return
			}
		} else {
			//Copy file from source to destination
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
			//remove file from source after copying
			err = RemoveFileIfExists(srcPath)
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

// CreateZipFileFromProject if the given projectPath contains a directory, zip it and return the zip file path.
//	Otherwise, leave it as it is.
// @param projectPath Project path
// @param skipCleanup Whether to clean the temporary files after the program exists
// @return string Path to the zip file
// @return error
// @return func() can be called to cleanup the temporary items created during this function execution. Needs to call
//	this once the zip file is consumed
func CreateZipFileFromProject(projectPath string, skipCleanup bool) (string, error, func()) {
	// If the projectPath contains a directory, zip it
	if info, err := os.Stat(projectPath); err == nil && info.IsDir() {
		tmp, err := ioutil.TempFile("", "project-artifact*.zip")
		if err != nil {
			return "", err, nil
		}
		Logln(LogPrefixInfo+"Creating the project artifact", tmp.Name())
		err = Zip(projectPath, tmp.Name())
		if err != nil {
			return "", err, nil
		}
		//creates a function to cleanup the temporary folders
		cleanup := func() {
			if skipCleanup {
				Logln(LogPrefixInfo+"Leaving", tmp.Name())
				return
			}
			Logln(LogPrefixInfo+"Deleting", tmp.Name())
			err := os.Remove(tmp.Name())
			if err != nil {
				Logln(LogPrefixError + err.Error())
			}
		}
		projectPath = tmp.Name()
		return projectPath, nil, cleanup
	}
	return projectPath, nil, nil
}

// Get a cloned copy of a given folder or a ZIP file path. If a zip file path is given, it will be extracted to the
//	tmp folder. The returned string will be the path to the extracted temp folder.
func GetTempCloneFromDirOrZip(path string) (string, error) {
	fileIsDir := false
	// create a temp directory
	tmpDir, err := ioutil.TempDir("", "apim")
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}

	if info, err := os.Stat(path); err == nil {
		fileIsDir = info.IsDir()
	} else {
		return "", err
	}
	if fileIsDir {
		// copy dir to a temp location
		Logln(LogPrefixInfo+"Copying from", path, "to", tmpDir)
		dest := filepath.Join(tmpDir, filepath.Base(path))
		err = CopyDir(path, dest)
		if err != nil {
			return "", err
		}
		return dest, nil
	} else {
		// try to extract archive
		Logln(LogPrefixInfo+"Extracting", path, "to", tmpDir)
		finalPath, err := extractArchive(path, tmpDir)
		if err != nil {
			return "", err
		}
		return finalPath, nil
	}
}

// extractArchive extracts the API and give the path.
// In API Manager archive there is a directory in the root which contains the API
// this function returns it appended to the destination path
func extractArchive(src, dest string) (string, error) {
	files, err := Unzip(src, dest)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("invalid API archive")
	}
	r := strings.TrimPrefix(files[0], src)
	return filepath.Join(dest, strings.Split(filepath.Clean(r), string(os.PathSeparator))[0]), nil
}

// Creates a temporary folder and creates a zip file with a given name (zipFileName) from the given REST API response.
//	Returns the location of the created zip file.
func WriteResponseToTempZip(zipFileName string, resp *resty.Response) (string, error) {
	// Create a temp directory to save the original zip from the REST API
	tmpDir, err := ioutil.TempDir("", "apim")
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}

	tempZipFile := filepath.Join(tmpDir, zipFileName)

	// Save the zip file in the temp directory.
	// permission 644 : Only the owner can read and write.. Everyone else can only read.
	err = ioutil.WriteFile(tempZipFile, resp.Body(), 0644)
	if err != nil {
		return "", err
	}
	return tempZipFile, err
}

// CreateZipFile if the given filePath contains a directory, zip it
// @param filePath Project path
// @param skipCleanup Whether to clean the temporary files after the program exists
// @return error
// @return func() can be called to cleanup the temporary items created during this function execution. Needs to call
//	this once the zip file is consumed
func CreateZipFile(filePath string, skipCleanup bool) (error, func()) {
	// If the filePath contains a directory, zip it
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		if err != nil {
			return err, nil
		}
		sourceZipPath := filePath + ZipFileSuffix
		Logln(LogPrefixInfo+"Creating the zipDirectory artifact", sourceZipPath)
		err = Zip(filePath, sourceZipPath)
		//creates a function to cleanup the temporary folders
		cleanup := func() {
			if skipCleanup {
				Logln(LogPrefixInfo+"Leaving", filePath)
				return
			}
			Logln(LogPrefixInfo+"Deleting", filePath)
			err := RemoveDirectoryIfExists(filePath)
			if err != nil {
				Logln(LogPrefixError + err.Error())
			}
		}
		return nil, cleanup
	}
	return nil, nil
}
