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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

// Invoke http-post request using go-resty
func InvokePOSTRequest(url string, headers map[string]string, body interface{}) (*resty.Response, error) {
	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetBody(body).Post(url)
}

// Invoke http-post request without body using go-resty
func InvokePOSTRequestWithoutBody(url string, headers map[string]string) (*resty.Response, error) {
	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).Post(url)
}

// Invoke http-post request with query parameters using go-resty
func InvokePOSTRequestWithQueryParam(queryParam map[string]string, url string, headers map[string]string,
	body string) (*resty.Response, error) {

	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetQueryParams(queryParam).SetBody(body).Post(url)
}

// Invoke http-post request with file & query parameters using go-resty
func InvokePOSTRequestWithFileAndQueryParams(queryParam map[string]string, url string, headers map[string]string,
	fileParamName, filePath string) (*resty.Response, error) {

	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetQueryParams(queryParam).
		SetFile(fileParamName, filePath).Post(url)
}

// Invoke http-post request with file using go-resty
func InvokePOSTRequestWithFile(url string, headers map[string]string,
	fileParamName, filePath string) (*resty.Response, error) {

	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).
		SetFile(fileParamName, filePath).Post(url)
}

// Invoke http-get request using go-resty
func InvokeGETRequest(url string, headers map[string]string) (*resty.Response, error) {
	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).Get(url)
}

// Invoke http-get request with query param
func InvokeGETRequestWithQueryParam(queryParam string, paramValue string, url string, headers map[string]string) (
	*resty.Response, error) {

	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetQueryParam(queryParam, paramValue).Get(url)
}

// Invoke http-get request with multiple query params
func InvokeGETRequestWithMultipleQueryParams(queryParam map[string]string, url string, headers map[string]string) (
	*resty.Response, error) {

	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetQueryParams(queryParam).Get(url)
}

// Invoke http-get request with query params as string
func InvokeGETRequestWithQueryParamsString(url, queryParams string, headers map[string]string) (
	*resty.Response, error) {

	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetQueryString(queryParams).Get(url)
}

// Invoke http-put request with multiple query params
func InvokePutRequest(queryParam map[string]string, url string, headers map[string]string, body string) (
	*resty.Response, error) {
	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetQueryParams(queryParam).SetBody(body).Put(url)
}

func InvokePUTRequestWithoutQueryParams(url string, headers map[string]string, body interface{}) (*resty.Response, error) {
	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetBody(body).Put(url)
}

// Invoke http-delete request using go-resty
func InvokeDELETERequest(url string, headers map[string]string) (*resty.Response, error) {
	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).Delete(url)
}

// Invoke http-delete request with multiple query params
func InvokeDELETERequestWithParams(url string, params map[string]string, headers map[string]string) (
	*resty.Response, error) {

	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetQueryParams(params).Delete(url)
}

// Invoke http-patch request using go-resty
func InvokePATCHRequest(url string, headers map[string]string, body map[string]string) (*resty.Response, error) {
	client := resty.New()

	if Insecure {
		client.SetTLSClientConfig(
			&tls.Config{InsecureSkipVerify: true, // To bypass errors in SSL certificates
				Renegotiation: TLSRenegotiationMode})
	} else {
		client.SetTLSClientConfig(GetTlsConfigWithCertificate())
	}

	client.SetTimeout(time.Duration(HttpRequestTimeout) * time.Millisecond)
	return client.R().SetHeaders(headers).SetBody(body).Patch(url)
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

// Remove revision from revision name
func GetRevisionNumFromRevisionName(input string) string {
	return strings.ReplaceAll(input, "Revision ", "")
}

// Add revision to revision number
func GetRevisionNamFromRevisionNum(input string) string {
	return "Revision-" + input
}

func GetPolicyNameByPolicyDefinitionFile(originalFilePath, ext string) (string, error) {
	policyDataImport := &PolicyDataImport{}
	policyDefFile, err := ioutil.ReadFile(originalFilePath)
	if err != nil {
		return "", err
	}

	if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(policyDefFile, &policyDataImport)
	} else if ext == ".json" {
		err = json.Unmarshal(policyDefFile, &policyDataImport)
	}

	if err != nil {
		return "", err
	}
	return policyDataImport.Data.Name, nil
}

// validate integer values are correctly provided
func ValidateFlagWithIntegerValues(value string) (int, error) {
	limit, err := strconv.Atoi(value)

	if err != nil {
		return -1, err
	}

	return limit, nil
}
