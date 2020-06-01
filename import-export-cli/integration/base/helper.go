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
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"log"
)

// logTransport : Flag which determinse if http transport level requests and responses are logged
var logTransport = false

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

// GetRowsFromTableResponse : Parse tabular apictl output to retreive an array of rows.
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

// SetupEnv : Adds a new environment and automitcally removes it when the calling test function execution ends
//
func SetupEnv(t *testing.T, env string, apim string, tokenEp string) {
	Execute(t, "add-env", "-e", env, "--apim", apim, "--token", tokenEp)

	t.Cleanup(func() {
		Execute(t, "remove", "env", env)
	})
}

// Login : Logs into an environment and automitcally logs out when the calling test function execution ends
//
func Login(t *testing.T, env string, username string, password string) {
	Execute(t, "login", env, "-u", username, "-p", password, "-k", "--verbose")

	t.Cleanup(func() {
		Execute(t, "logout", env)
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

// Fatal : Log and Fail execution now. This is not equivalant to testing.T.Fatal(),
// which will exit the calling go routine. This function will result in the process exiting.
// It should be used in scenarios where the testing context testing.T is not available
// in order to call Fatal(), such as code that is executed from TestMain(m *testing.M)
func Fatal(v ...interface{}) {
	log.Fatalln(v...)
}

// Log : This is equivalant to testing.T.Log() and is dependent on test -v(verbose) flag.
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

// ValidateAndLogResponse : Validate reposne against expected status code and optionally log the response
func ValidateAndLogResponse(logString string, response *http.Response, expectedStatusCode int) {
	if response.StatusCode != expectedStatusCode {
		FatalStatusCodeResponse(logString, response)
	}

	LogResponse(logString, response)
}

// FatalStatusCodeResponse : Handle response with Status Code that is considerd fatal.
// Log response and exit the process
func FatalStatusCodeResponse(logString string, response *http.Response) {
	logString += " - Unexpected Status Code in response:"
	logResponse(logString, response)
	os.Exit(1)
}

// FatalContentTypeResponse : Handle response with Content-Type that is considerd fatal.
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
