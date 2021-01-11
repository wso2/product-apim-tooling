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

package base

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"log"
)

// logTransport : Flag which determines if http transport level requests and responses are logged
var logTransport = false

// indexingDelay : Time in milliseconds that tests need to wait for to allow APIM solr indexing to take place
var indexingDelay = 1000

func init() {
	flag.BoolVar(&logTransport, "logtransport", false, "Log http transport level requests and responses")
}

// Execute : Run apictl command
//
func Execute(t *testing.T, args ...string) (string, error) {
	cmd := exec.Command(RelativeBinaryPath+BinaryName, args...)

	t.Log("base.Execute() - apictl command:", cmd.String())
	// run command
	output, err := cmd.Output()

	t.Log("base.Execute() - apictl command output:", string(output))
	return string(output), err
}

// GetRowsFromTableResponse : Parse tabular apictl output to retrieve an array of rows.
// This friendly format aids in asserting results during testing via simple string comparison.
//
func GetRowsFromTableResponse(response string) []string {
	// Replace Windows carriage return if exists and split by new line to get rows
	result := strings.Split(strings.ReplaceAll(response, "\r\n", "\n"), "\n")

	// Remove Column header row and trailing new line character row
	return result[1 : len(result)-1]
}

// GetValueOfUniformResponse : Trim uniformally formatted output from apictl.
//
func GetValueOfUniformResponse(response string) string {
	return strings.TrimSpace(strings.Split(response, "output: ")[0])
}

// SetupEnv : Adds a new environment and automatically removes it when the calling test function execution ends
//
func SetupEnv(t *testing.T, env string, apim string, tokenEp string) {
	Execute(t, "add", "env", env, "--apim", apim, "--token", tokenEp)

	t.Cleanup(func() {
		Execute(t, "remove", "env", env)
	})
}

// SetupEnv : Adds a new environment just with apim endpoint and automatically removes it when the
// calling test function execution ends
func SetupEnvWithoutTokenFlag(t *testing.T, env string, apim string) {
	Execute(t, "add", "env", env, "--apim", apim)

	t.Cleanup(func() {
		Execute(t, "remove", "env", env)
	})
}

// SetupMIEnv : Adds a new mi environment and automatically removes it when the calling test function execution ends
//
func SetupMIEnv(t *testing.T, env, mi string) {
	Execute(t, "add", "env", env, "--mi", mi)

	t.Cleanup(func() {
		Execute(t, "remove", "env", env)
	})
}

// Login : Logs into an environment and automatically logs out when the calling test function execution ends
//
func Login(t *testing.T, env string, username string, password string) {
	Execute(t, "login", env, "-u", username, "-p", password, "-k", "--verbose")

	t.Cleanup(func() {
		Execute(t, "logout", env)
	})
}

// MILogin : Logs into an mi environment and automatically logs out when the calling test function execution ends
//
func MILogin(t *testing.T, env string, username string, password string) {
	Execute(t, "mi", "login", env, "-u", username, "-p", password, "-k", "--verbose")

	t.Cleanup(func() {
		Execute(t, "mi", "logout", env)
	})
}

// IsAPIArchiveExists : Returns true if exported application archive exists on file system, else returns false
//
func IsAPIArchiveExists(t *testing.T, path string, name string, version string) bool {
	file := constructAPIFilePath(path, name, version)

	t.Log("base.IsAPIArchiveExists() - archive file path:", file)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}

// RemoveAPIArchive : Remove exported api archive from file system
//
func RemoveAPIArchive(t *testing.T, path string, name string, version string) {
	file := constructAPIFilePath(path, name, version)

	t.Log("base.RemoveAPIArchive() - archive file path:", file)

	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)

		if err != nil {
			t.Fatal(err)
		}
	}
}

// GetAPIArchiveFilePath : Get API archive file path
func GetAPIArchiveFilePath(t *testing.T, path string, name string, version string) string {
	file := constructAPIFilePath(path, name, version)

	t.Log("base.GetAPIArchiveFilePath() - archive file path:", file)

	return file
}

func constructAPIFilePath(path string, name string, version string) string {
	return filepath.Join(path, name+"_"+version+".zip")
}

// IsApplicationArchiveExists : Returns true if exported application archive exists on file system, else returns false
//
func IsApplicationArchiveExists(t *testing.T, path string, name string, owner string) bool {
	file := constructAppFilePath(path, name, owner)

	t.Log("base.IsApplicationArchiveExists() - archive file path:", file)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}

// RemoveApplicationArchive : Remove exported application archive from file system
//
func RemoveApplicationArchive(t *testing.T, path string, name string, owner string) {
	file := constructAppFilePath(path, name, owner)

	t.Log("base.RemoveApplicationArchive() - archive file path:", file)

	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)

		if err != nil {
			t.Fatal(err)
		}
	}
}

// GetApplicationArchiveFilePath : Get Application archive file path
func GetApplicationArchiveFilePath(t *testing.T, path string, name string, owner string) string {
	file := constructAppFilePath(path, name, owner)

	t.Log("base.GetApplicationArchiveFilePath() - archive file path:", file)

	return file
}

func constructAppFilePath(path string, name string, owner string) string {
	return filepath.Join(path, strings.ReplaceAll(owner, "/", "#")+"_"+name+".zip")
}

// Fatal : Log and Fail execution now. This is not equivalent to testing.T.Fatal(),
// which will exit the calling go routine. This function will result in the process exiting.
// It should be used in scenarios where the testing context testing.T is not available
// in order to call Fatal(), such as code that is executed from TestMain(m *testing.M)
func Fatal(v ...interface{}) {
	log.Fatalln(v...)
}

// Log : This is equivalent to testing.T.Log() and is dependent on test -v(verbose) flag.
// It should be used in scenarios where the testing context testing.T is not available
// in order to call Log(), such as code that is executed from TestMain(m *testing.M)
func Log(v ...interface{}) {
	if testing.Verbose() {
		log.Println(v...)
	}
}

// LogResponse : Log http response received
func LogResponse(logString string, response *http.Response) {
	if logTransport {
		logString += " - response:"
		logResponse(logString, response)
	}
}

// LogRequest : Log http request sent
func LogRequest(logString string, resquest *http.Request) {
	if logTransport {
		logString += " - request:"
		logRequest(logString, resquest)
	}
}

// ValidateAndLogResponse : Validate response against expected status code and optionally log the response
func ValidateAndLogResponse(logString string, response *http.Response, expectedStatusCode int) {
	if response.StatusCode != expectedStatusCode {
		FatalStatusCodeResponse(logString, response)
	}

	LogResponse(logString, response)
}

// FatalStatusCodeResponse : Handle response with Status Code that is considered fatal.
// Log response and exit the process
func FatalStatusCodeResponse(logString string, response *http.Response) {
	logString += " - Unexpected Status Code in response:"
	logResponse(logString, response)
	os.Exit(1)
}

// FatalContentTypeResponse : Handle response with Content-Type that is considered fatal.
// Log response and exit the process
func FatalContentTypeResponse(logString string, response *http.Response) {
	logString += " - Unexpected Content-Type in response:"
	logResponse(logString, response)
	os.Exit(1)
}

func logRequest(logString string, resquest *http.Request) {
	dump, err := httputil.DumpRequest(resquest, true)

	if err != nil {
		Fatal(err)
	}

	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", logString, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	log.Println(string(dump))
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}

func logResponse(logString string, response *http.Response) {
	dump, err := httputil.DumpResponse(response, true)

	if err != nil {
		Fatal(err)
	}

	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", logString, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	log.Println(string(dump))
	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}

// SetIndexingDelay : Set time in milliseconds that tests need to wait for to allow APIM solr indexing to take place
func SetIndexingDelay(delay int) {
	indexingDelay = delay
}

// WaitForIndexing : Wait for specified interval to allow APIM solr indexes to be updated
func WaitForIndexing() {
	time.Sleep(time.Duration(indexingDelay) * time.Millisecond)
}

func RemoveDir(projectName string) {
	error := os.RemoveAll(projectName)
	if error != nil {
		log.Fatal(error)
	}
}

// CreateTempDir : Create temp directory at the specified root path.
// The directory will be removed when the calling test exits.
func CreateTempDir(t *testing.T, path string) {
	t.Log("base.CreateTempDir() - path:", path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() {
			os.RemoveAll(path)
		})
	}
}

func GetExportedPathFromOutput(output string) string {
	//Check directory path to omit changes due to OS differences
	if strings.Contains(output, ":\\") {
		arrayOutput := []rune(output)
		extractedPath := string(arrayOutput[strings.Index(output, ":\\")-1:])
		return strings.ReplaceAll(strings.ReplaceAll(extractedPath, "\n", ""), " ", "")
	} else {
		return strings.ReplaceAll(strings.ReplaceAll(output[strings.Index(output, string(os.PathSeparator)):], "\n", ""), " ", "")
	}
}

//Count number of files in a directory
func CountFiles(path string) (int, error) {
	i := 0
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}
	for _, file := range files {
		if !file.IsDir() {
			i++
		}
	}
	return i, nil
}

// IsFileAvailable checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func IsFileAvailable(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Copy the src file to dst. Any existing file will be overwritten
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

//Generate random strings with given length
func GenerateRandomName(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
