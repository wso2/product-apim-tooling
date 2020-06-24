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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

var (
	mgwControlPlaneHostDeleting string
	mgwLabelDeleting            string
	deletedAPIName              string
)

const deleteSwaggerMgwCmdExample = `apictl mgw-delete-api --host http://localhost:9095 --labels mgw_lbl --api swagger.json`

var deleteAPIMgwCmd = &cobra.Command{
	Use:     "mgw-delete-api --host [control plane url] --labels [microgateway labels] --api [api name]",
	Short:   "Delete a API definition from Microgateway.",
	Long:    "Delete a API definition from Microgateway.",
	Example: deleteSwaggerMgwCmdExample,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "mgw-delete-api called")

		err := executeDeleteSwaggerMgw()
		if err != nil {
			utils.HandleErrorAndExit("Error deleting API from microgateway", err)
		}
	},
}

func executeDeleteSwaggerMgw() error {
	if mgwControlPlaneHostDeleting == "" {
		mgwControlPlaneHostDeleting = "http://localhost:9095"
	}
	if mgwLabelDeleting == "" {
		mgwLabelDeleting = "default"
	}
	if deletedAPIName != "" {
		deleteAPIMgw(mgwLabelDeleting, deletedAPIName)
	}
	return nil
}

func deleteAPIMgw(label string, apiName string) string {
	// create a buffer we can write the file to
	fileDataBuffer := bytes.Buffer{}
	multipartWriter := multipart.NewWriter(&fileDataBuffer)
	// add label as a field to the body
	_ = multipartWriter.WriteField("label", label)
	// add api name as a field to the body
	_ = multipartWriter.WriteField("apiName", apiName)
	// close the file writer. This lets it know we're done copying in data
	multipartWriter.Close()
	// create the POST request to send the file data to the server
	url := mgwControlPlaneHostUpdating + "/api/delete"
	req, err := http.NewRequest("DELETE", url, &fileDataBuffer)
	if err != nil {
		log.Fatal(err)
	}
	// we set the header so the server knows about the files content
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	// send the DELETE request and receive the response data
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// get data from the response body
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// return the response data
	return string(data)
}

func init() {
	RootCmd.AddCommand(deleteAPIMgwCmd)
	deleteAPIMgwCmd.Flags().StringVarP(&mgwControlPlaneHostDeleting, "host", "", "", "Provide the host url "+
		"for the control plane with port")
	deleteAPIMgwCmd.Flags().StringVarP(&mgwLabelDeleting, "labels", "", "", "Provide label for the "+
		"microgateway instances you want to add the API")
	deleteAPIMgwCmd.Flags().StringVarP(&deletedAPIName, "api", "", "", "Provide the API name")
}
