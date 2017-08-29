// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0 // // Unless required by applicable law or agreed to in writing, software // distributed under the License is distributed on an "AS IS" BASIS, // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"bytes"
	"crypto/tls"
	"github.com/menuka94/wso2apim-cli/utils"
	"github.com/spf13/cobra"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var importAPIName string
var importEnvironment string
var ImportAPICmdUsername string
var ImportAPICmdPassword string

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use:   "import-api (--name <name-of-the-api> --environment <environment-to-which-the-api-should-be-imported>)",
	Short: utils.ImportAPICmdShortDesc,
	Long:  utils.ImportAPICmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("importAPI called")
		for key, arg := range args {
			fmt.Println(key, ":", arg)
		}

		accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(importEnvironment, ImportAPICmdUsername, ImportAPICmdPassword)

		if preCommandErr == nil {

			resp := ImportAPI(importAPIName, apiManagerEndpoint, accessToken)
			fmt.Printf("Status: %v\n", resp.Status)
			if resp.StatusCode == 200 {
				fmt.Println("Header:", resp.Header)
				fmt.Printf("Body: %s\n", resp.Body)
			}
			//fmt.Printf("Errors: %v\n", resp.Error)
		} else {
			// env_endpoints_all.yaml file is not configured properly by the user
			log.Fatal("Error:", preCommandErr)
		}
	},
}

func ImportAPI(name string, url string, accessToken string) *http.Response {
	// append '/' to the end if there isn't one already
	if string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "import/apis"

	filePath, _ := os.Getwd()
	filePath += "/exported/" + name
	fmt.Println("filePath:", filePath)

	// check if '.zip' exists in the input 'name'
	hasZipExtension, _ := regexp.MatchString(`^\S+\.zip$`, name)

	if hasZipExtension {
		// import the zip file directly
		fmt.Println("hasZipExtension: ", true)

	}else{
		fmt.Println("hasZipExtension: ", false)
		// search for a directory with the given name
		destination, _ := os.Getwd()
		destination += "/exported/" + name + ".zip"
		err := utils.ZipDir(filePath,  destination)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		filePath += ".zip"
	}

	extraParams := map[string]string{}

	req, err := newFileUploadRequest(url, extraParams, "file", filePath, accessToken)
	if err != nil {
		fmt.Println("Error creating request")
		log.Fatal(err)
		panic(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	} else {
		var bodyContent []byte
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		fmt.Println(bodyContent)
	}

	return resp
}

func newFileUploadRequest(uri string, params map[string]string, paramName, path string, accessToken string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", uri, body)
	request.Header.Add(utils.HeaderAuthorization, utils.HeaderValueAuthBearerPrefix+" "+accessToken)
	request.Header.Add(utils.HeaderContentType, writer.FormDataContentType())
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Connection", "keep-alive")

	return request, err
}

func init() {
	RootCmd.AddCommand(ImportAPICmd)
	ImportAPICmd.Flags().StringVarP(&importAPIName, "name", "n", "", "Name of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importEnvironment, "environment", "e", "", "Environment from the which the API should be imported")
	ImportAPICmd.Flags().StringVarP(&ImportAPICmdUsername, "username", "u", "", "Username")
	ImportAPICmd.Flags().StringVarP(&ImportAPICmdPassword, "password", "p", "", "Password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ImportAPICmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ImportAPICmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
