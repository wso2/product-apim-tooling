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
	mgwLabelAdding            string
	addedAPIName              string
	addedSwaggerPath          string
)

type MgwResponse struct {
	Message string
}

const (
	defaultMgwHostUrl = "http://localhost:9095"
	defaultMgwLabel   = "default"
	defaultAPIName    = "api_v1"
)

const addAPIMgwCmdExample = `apictl mgw-add-api --host http://localhost:9095 --labels mgw_lbl --api api_v1 --oas https://petstore.swagger.io/v2/swagger.json`

var addAPIMgwCmd = &cobra.Command{
	Use:     "mgw-add-api --host [control plane url] --labels [microgateway labels] --api [api name] --oas [swagger path]",
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
	if mgwControlPlaneHostAdding == "" {
		mgwControlPlaneHostAdding = defaultMgwHostUrl
	}
	if mgwLabelAdding == "" {
		mgwLabelAdding = defaultMgwLabel
	}
	if addedAPIName == "" {
		addedAPIName = defaultAPIName
	}
	if addedSwaggerPath != "" {
		addAPIMgw(mgwLabelAdding, addedAPIName, addedSwaggerPath)
	}
	return nil
}

func addAPIMgw(label string, apiName string, apiDefinition string) {
	var file *os.File
	//err := os.Mkdir(importedApiDefinitionsPath, os.ModePerm)
	//if err != nil {
	//	log.Fatal(err)
	//}
	if isUrl(apiDefinition) {
		// Get the data
		resp, err := http.Get(apiDefinition)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		// reading the file name from the link
		fileName := apiDefinition[strings.LastIndex(apiDefinition, "/")+1:]
		// Create the file
		out, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		file, err = os.Open(fileName)
	} else {
		// open the local file we want to upload
		var err error
		file, err = os.Open(apiDefinition)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create a buffer we can write the file to
	fileDataBuffer := bytes.Buffer{}
	multipartWritter := multipart.NewWriter(&fileDataBuffer)

	// create an http formfile. This wraps our local file in a format that can be sent to the server
	formFile, err := multipartWritter.CreateFormFile("swaggerFile", file.Name())
	if err != nil {
		log.Fatal(err)
	}
	// copy the file we want to upload into the form file wrapper
	_, err = io.Copy(formFile, file)
	if err != nil {
		log.Fatal(err)
	}

	// add label as a field to the body
	_ = multipartWritter.WriteField("label", label)
	// add api name as a field to the body
	_ = multipartWritter.WriteField("apiName", apiName)

	// close the file writer. This lets it know we're done copying in data
	multipartWritter.Close()
	// create the POST request to send the file data to the server
	url := mgwControlPlaneHostAdding + "/api/add"
	req, err := http.NewRequest("POST", url, &fileDataBuffer)
	if err != nil {
		log.Fatal(err)
	}
	// we set the header so the server knows about the files content
	req.Header.Set("Content-Type", multipartWritter.FormDataContentType())
	// send the POST request and receive the response data
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	log.Println(responseBody.Message)
}

func init() {
	RootCmd.AddCommand(addAPIMgwCmd)
	addAPIMgwCmd.Flags().StringVarP(&mgwControlPlaneHostAdding, "host", "", "", "Provide the host url "+
		"for the control plane with port")
	addAPIMgwCmd.Flags().StringVarP(&mgwLabelAdding, "labels", "", "", "Provide label for the "+
		"microgateway instances you want to add the API")
	addAPIMgwCmd.Flags().StringVarP(&addedAPIName, "api", "", "", "Provide the API name")
	addAPIMgwCmd.Flags().StringVarP(&addedSwaggerPath, "oas", "", "", "Provide an OpenAPI "+
		"specification file for the API")

	_ = addAPIMgwCmd.MarkFlagRequired("oas")
}
