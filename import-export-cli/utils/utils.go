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
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty"
	"golang.org/x/crypto/ssh/terminal"
)

func InvokePOSTRequest(url string, headers map[string]string, body interface{}) (*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).SetBody(body).Post(url)

	return resp, err
}

// Invoke http-post request without body using go-resty
func InvokePOSTRequestWithoutBody(url string, headers map[string]string) (*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).Post(url)

	return resp, err
}

// Invoke http-get request using go-resty
func InvokeGETRequest(url string, headers map[string]string) (*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).Get(url)

	return resp, err
}

// Invoke http-get request with query param
func InvokeGETRequestWithQueryParam(queryParam string, paramValue string, url string, headers map[string]string) (
	*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).SetQueryParam(queryParam, paramValue).Get(url)

	return resp, err
}

// Invoke http-get request with multiple query params
func InvokeGETRequestWithMultipleQueryParams(queryParam map[string]string, url string, headers map[string]string) (
	*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).SetQueryParams(queryParam).Get(url)

	return resp, err
}

// Invoke http-put request
func InvokePutRequest(queryParam map[string]string, url string, headers map[string]string, body string) (
	*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).SetQueryParams(queryParam).SetBody(body).Put(url)

	return resp, err
}

//Invoke POST request with query parameters
func InvokePostRequestWithQueryParam(queryParam map[string]string, url string, headers map[string]string, body string) (
	*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).SetQueryParams(queryParam).SetBody(body).Post(url)

	return resp, err
}

// Invoke http-delete request using go-resty
func InvokeDELETERequest(url string, headers map[string]string) (*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).Delete(url)

	return resp, err
}

// Invoke http-patch request using go-resty
func InvokePATCHRequest(url string, headers map[string]string, body map[string]string) (*resty.Response, error) {
	if Insecure {
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in SSL certificates
	} else {
		resty.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}
	if os.Getenv("HTTP_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTP_PROXY"))
	} else if os.Getenv("HTTPS_PROXY") != "" {
		resty.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("http_proxy") != "" {
		resty.SetProxy(os.Getenv("http_proxy"))
	} else if os.Getenv("https_proxy") != "" {
		resty.SetProxy(os.Getenv("https_proxy"))
	}
	resty.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	resp, err := resty.R().SetHeaders(headers).SetBody(body).Patch(url)

	return resp, err
}

func PromptForUsername() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	return username
}

func PromptForPassword() string {
	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(0)
	password := string(bytePassword)
	fmt.Println()
	return password
}

// ShowHelpCommandTip function will print the instructions for displaying help info on a specific command
// @params cmdLiteral Command on which help command is to be displayed
func ShowHelpCommandTip(cmdLiteral string) {
	fmt.Printf("Execute '%s %s --help' for more info.\n", ProjectName, cmdLiteral)
}

// return a string containing the file name, function name
// and the line number of a specified entry on the call stack
func WhereAmI(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("File: %s  Function: %s Line: %d", chopPath(file), runtime.FuncForPC(function).Name(), line)
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}

// Append a slash to a string if there isn't one already
func AppendSlashToString(input string) string {
	if string(input[len(input)-1]) != "/" {
		input += "/"
	}
	return input
}

func WriteToFileSystem(exportAPIName, exportAPIVersion, exportEnvironment, exportDirectory string, resp *resty.Response) {
	// Write to file
	directory := filepath.Join(exportDirectory, exportEnvironment)
	// create directory if it doesn't exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0777)
		// permission 777 : Everyone can read, write, and execute
	}
	zipFilename := exportAPIName + "_" + exportAPIVersion + ".zip" // MyAPI_1.0.0.zip
	pFile := filepath.Join(directory, zipFilename)
	err := ioutil.WriteFile(pFile, resp.Body(), 0644)
	// permission 644 : Only the owner can read and write.. Everyone else can only read.
	if err != nil {
		HandleErrorAndExit("Error creating zip archive", err)
	}
	fmt.Println("Successfully exported API!")
	fmt.Println("Find the exported API at " + pFile)

}

// SetToK8sMode sets the "api-ctl" mode to kubernetes
func SetToK8sMode() {
	// read the existing config vars
	configVars := GetMainConfigFromFile(MainConfigFilePath)
	configVars.Config.KubernetesMode = true
	WriteConfigFile(configVars, MainConfigFilePath)
}

// returns min of two ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
