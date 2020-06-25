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
	"regexp"
	"strings"
)

var (
	mgwControlPlaneHostUpdating string
	mgwLabelsUpdating           string
	updatedAPIName              string
	updatedSwaggerPath          string
)

const (
	defaultUpdatedMgwHostUrl = "http://localhost:9095"
	defaultUpdatedMgwLabel   = "default"
	defaultUpdatedAPIName    = "api_v1"
	defaultUpdatedAPIDest    = "./mgw-api-definitions/"
	updateAPIMgwCmdExample   = `apictl mgw-update-api --host http://localhost:9095 --labels label1,label2 
									--api api_v1 --oas https://petstore.swagger.io/v2/swagger.json`
)

var updateAPIMgwCmd = &cobra.Command{
	Use: "mgw-update-api --host [control plane url] --labels [microgateway labels] --api [api name] " +
		"--oas [swagger path]",
	Short:   "Update an API swagger definition in Microgateway.",
	Long:    "Update an API swagger definition in Microgateway. You can provide either a file on the disk or a link.",
	Example: updateAPIMgwCmdExample,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "mgw-update-api called")
		err := executeUpdateAPIMgw()
		if err != nil {
			utils.HandleErrorAndExit("Error updating swagger definition in microgateway", err)
		}
	},
}

func executeUpdateAPIMgw() error {
	// TODO: add control plane url to env
	if mgwControlPlaneHostUpdating == "" {
		mgwControlPlaneHostUpdating = defaultUpdatedMgwHostUrl
	}
	if mgwLabelsUpdating == "" {
		mgwLabelsUpdating = defaultUpdatedMgwLabel
	}
	if updatedAPIName == "" {
		updatedAPIName = defaultUpdatedAPIName
	}
	if updatedSwaggerPath != "" {
		updateAPIMgw(mgwLabelsUpdating, updatedAPIName, updatedSwaggerPath)
	}
	return nil
}

func isUrl(apiDefinition string) bool {
	match, _ := regexp.MatchString("^(http|https)://(.)+", apiDefinition)
	return match
}

func updateAPIMgw(labels string, apiName string, apiDefinition string) {
	var file *os.File
	//err := os.Mkdir(importedApiDefinitionsPath, os.ModePerm)
	//if err != nil {
	//	log.Fatal(err)
	//}
	if isUrl(apiDefinition) {
		// Get the data
		resp, err := http.Get(apiDefinition)
		if err != nil {
			utils.HandleErrorAndExit("Error downloading the file from the link", err)
		}
		defer resp.Body.Close()
		if _, err := os.Stat(defaultUpdatedAPIDest); os.IsNotExist(err) {
			err = os.Mkdir(defaultUpdatedAPIDest, os.ModePerm)
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
	url := mgwControlPlaneHostUpdating + "/api/update"
	req, err := http.NewRequest("PUT", url, &fileDataBuffer)
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

	var responseBody MgwResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		utils.HandleErrorAndExit("Error reading the response", err)
	}
	log.Println(responseBody.Message)
}

func init() {
	RootCmd.AddCommand(updateAPIMgwCmd)
	updateAPIMgwCmd.Flags().StringVarP(&mgwControlPlaneHostUpdating, "host", "", "",
		"Provide the host url for the control plane with port")
	updateAPIMgwCmd.Flags().StringVarP(&mgwLabelsUpdating, "labels", "", "",
		"Provide label for the microgateway instances you want to add the API")
	updateAPIMgwCmd.Flags().StringVarP(&updatedAPIName, "api", "", "", "Provide the API name")
	updateAPIMgwCmd.Flags().StringVarP(&updatedSwaggerPath, "oas", "", "",
		"Provide an OpenAPI specification file for the API")

	_ = updateAPIMgwCmd.MarkFlagRequired("oas")
	_ = updateAPIMgwCmd.MarkFlagRequired("api")
}
