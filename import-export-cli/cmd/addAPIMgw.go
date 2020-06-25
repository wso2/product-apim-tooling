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

package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var (
	mgwControlPlaneHostAdding string
	mgwLabelsAdding           string
	addedAPIName              string
	addedSwaggerPath          string
)

type MgwResponse struct {
	Message string
}

const (
	defaultAddedMgwHostUrl = "http://localhost:9095"
	defaultAddedMgwLabel   = "default"
	defaultAddedAPIName    = "api_v1"
	defaultAddedAPIDest    = "./mgw-api-definitions/"
	addAPIMgwCmdExample    = `apictl mgw-add-api --host http://localhost:9095 --labels label1,label2 --api api_v1 
								--oas https://petstore.swagger.io/v2/swagger.json`
)

var addAPIMgwCmd = &cobra.Command{
	Use: "mgw-add-api --host [control plane url] --labels [microgateway labels] --api [api name] " +
		"--oas [swagger path]",
	Short:   "Add a swagger file to Microgateway.",
	Long:    "Add a swagger file to Microgateway. You can provide either a file on the disk or a link.",
	Example: addAPIMgwCmdExample,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "mgw-add-api called")
		err := executeAddAPIMgw()
		if err != nil {
			utils.HandleErrorAndExit("Error adding swagger to microgateway", err)
		}
	},
}

func executeAddAPIMgw() error {
	// TODO: add control plane url to env
	if mgwControlPlaneHostAdding == "" {
		mgwControlPlaneHostAdding = defaultAddedMgwHostUrl
	}
	if mgwLabelsAdding == "" {
		mgwLabelsAdding = defaultAddedMgwLabel
	}
	if addedAPIName == "" {
		addedAPIName = defaultAddedAPIName
	}
	if addedSwaggerPath != "" {
		addAPIMgw(mgwLabelsAdding, addedAPIName, addedSwaggerPath)
	}
	return nil
}

func addAPIMgw(labels string, apiName string, apiDefinition string) {
	var file *os.File
	if isUrl(apiDefinition) {
		// Get the data
		resp, err := http.Get(apiDefinition)
		if err != nil {
			utils.HandleErrorAndExit("Error downloading the file from the link", err)
		}
		defer resp.Body.Close()
		if _, err := os.Stat(defaultAddedAPIDest); os.IsNotExist(err) {
			err = os.Mkdir(defaultAddedAPIDest, os.ModePerm)
			if err != nil {
				utils.HandleErrorAndExit("Error creating the destination directory", err)
			}
		}
		// reading the file name from the link
		fileName := apiDefinition[strings.LastIndex(apiDefinition, "/")+1:]
		//TODO: add dest folder location to env
		filePath := defaultAddedAPIDest + fileName
		// Create the file
		out, err := os.Create(filePath)
		if err != nil {
			utils.HandleErrorAndExit("Error creating the API definition file", err)
		}
		defer out.Close()
		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		file, err = os.Open(filePath)
		if err != nil {
			utils.HandleErrorAndExit("Error opening the file", err)
		}
	} else {
		// open the local file we want to upload
		var err error
		file, err = os.Open(apiDefinition)
		if err != nil {
			utils.HandleErrorAndExit("Error opening the file", err)
		}
	}

	// create a buffer we can write the file to
	fileDataBuffer := bytes.Buffer{}
	multipartWriter := multipart.NewWriter(&fileDataBuffer)

	// create an http formfile. This wraps our local file in a format that can be sent to the server
	formFile, err := multipartWriter.CreateFormFile("swaggerFile", file.Name())
	if err != nil {
		utils.HandleErrorAndExit("Error adding file to the request", err)
	}
	// copy the file we want to upload into the form file wrapper
	_, err = io.Copy(formFile, file)
	if err != nil {
		utils.HandleErrorAndExit("Error adding file to the request", err)
	}

	// add label as a field to the body
	_ = multipartWriter.WriteField("labels", labels)
	// add api name as a field to the body
	_ = multipartWriter.WriteField("apiName", apiName)

	// close the file writer. This lets it know we're done copying in data
	multipartWriter.Close()
	// create the POST request to send the file data to the server
	url := mgwControlPlaneHostAdding + "/api/add"
	req, err := http.NewRequest("POST", url, &fileDataBuffer)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the request", err)
	}
	// we set the header so the server knows about the files content
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	// send the POST request and receive the response data
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.HandleErrorAndExit("Error sending the request to the control plane", err)
	}
	// get data from the response body
	defer response.Body.Close()
	//data, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}

	var responseBody MgwResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		utils.HandleErrorAndExit("Error reading the response", err)
	}
	log.Println(responseBody.Message)
}

func init() {
	RootCmd.AddCommand(addAPIMgwCmd)
	addAPIMgwCmd.Flags().StringVarP(&mgwControlPlaneHostAdding, "host", "", "",
		"Provide the host url for the control plane with port")
	addAPIMgwCmd.Flags().StringVarP(&mgwLabelsAdding, "labels", "", "",
		"Provide label for the microgateway instances you want to add the API")
	addAPIMgwCmd.Flags().StringVarP(&addedAPIName, "api", "", "", "Provide the API name")
	addAPIMgwCmd.Flags().StringVarP(&addedSwaggerPath, "oas", "", "",
		"Provide an OpenAPI specification file for the API")

	_ = addAPIMgwCmd.MarkFlagRequired("oas")
	_ = addAPIMgwCmd.MarkFlagRequired("api")
}
