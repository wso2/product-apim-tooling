/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"path/filepath"
	"net/http"
)

var exportAPIName string
var exportAPIVersion string
var exportEnvironment string
var exportAPICmdUsername string
var exportAPICmdPassword string

// ExportAPICmd represents the exportAPI command
var ExportAPICmd = &cobra.Command{
	Use: "export-api (--name <name-of-the-api> --version <version-of-the-api> --environment " +
		"<environment-from-which-the-api-should-be-exported>)",
	Short: utils.ExportAPICmdShortDesc,
	Long:  utils.ExportAPICmdLongDesc + utils.ExportAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln("export-api called")

		accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(exportEnvironment, exportAPICmdUsername,
			exportAPICmdPassword)

		if preCommandErr == nil {
			resp := ExportAPI(exportAPIName, exportAPIVersion, apiManagerEndpoint, accessToken)

			// Print info on response
			utils.Logln("ResponseStatus: %v\n", resp.Status())
			utils.Logln("Error: %v\n", resp.Error())
			//fmt.Printf("Response Body: %v\n", resp.Body())

			if resp.StatusCode() == http.StatusOK {
				WriteToZip(exportAPIName, resp)

				// only to get the number of APIs exported
				numberOfAPIsExported, _, err := GetAPIList(exportAPIName, accessToken, apiManagerEndpoint)
				if err == nil {
					fmt.Println("Number of APIs exported:", numberOfAPIsExported)
				} else {
					utils.HandleErrorAndExit("Error getting list of APIs", err)
				}
			} else if resp.StatusCode() == http.StatusInternalServerError {
				// 500 Internal Server Error
				fmt.Println("Incorrect password")
			} else {
				// neither 200 nor 500
				fmt.Println("Error exporting API:", resp.Status())
			}
		} else {
			// error exporting API
			fmt.Println("Error exporting API:" + preCommandErr.Error())
		}
	},
}

// WriteToZip
// @param exportAPIName : Name of the API to be exported
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported API will be written to a zip file
func WriteToZip(exportAPIName string, resp *resty.Response) {
	// Write to file
	directory := utils.ExportDirectory
	// create directory if it doesn't exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0777)
		// permission 777 : Everyone can read, write, and execute
	}
	filename := exportAPIName + ".zip"
	pFile := filepath.Join(directory, filename)
	err := ioutil.WriteFile(pFile, resp.Body(), 0644)
	// permission 644 : Only the owner can read and write.. Everyone else can only read.
	if err != nil {
		utils.HandleErrorAndExit("Error creating zip archive", err)
	}
	fmt.Println("Succesfully exported and wrote to file")
}

// ExportAPI
// @param name : Name of the API to be exported
// @param version : Version of the API to be exported
// @param apimEndpoint : API Manager Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func ExportAPI(name string, version string, apimEndpoint string, accessToken string) *resty.Response {
	// append '/' to the end if there isn't one already
	if string(apimEndpoint[len(apimEndpoint)-1]) != "/" {
		apimEndpoint += "/"
	}

	apimEndpoint += "export/apis"

	query := "?query=" + name

	// TODO:: Add 'version' to the query (make sure the backend supports attribute searching)

	apimEndpoint += query
	fmt.Println("ExportAPI: URL:", apimEndpoint)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(apimEndpoint, headers)

	if err != nil {
		utils.HandleErrorAndExit("Error exporting API: "+name, err)
	}

	return resp
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ExportAPICmd)
	ExportAPICmd.Flags().StringVarP(&exportAPIName, "name", "n", "",
		"Name of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportAPIVersion, "version", "v", "",
		"Version of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportEnvironment, "environment", "e",
		utils.GetDefaultEnvironment(utils.MainConfigFilePath), "Environment to which the API should be exported")

	ExportAPICmd.Flags().StringVarP(&exportAPICmdUsername, "username", "u", "", "Username")
	ExportAPICmd.Flags().StringVarP(&exportAPICmdPassword, "password", "p", "", "Password")
}
