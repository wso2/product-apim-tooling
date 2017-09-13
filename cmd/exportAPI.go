// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/menuka94/wso2apim-cli/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"

	"crypto/tls"
)

var exportAPIName string
var exportAPIVersion string
var exportEnvironment string
var exportAPICmdUsername string
var exportAPICmdPassword string

// ExportAPICmd represents the exportAPI command
var ExportAPICmd = &cobra.Command{
	Use:   "export-api (--name <name-of-the-api> --version <version-of-the-api> --environment <environment-from-which-the-api-should-be-exported>)",
	Short: utils.ExportAPICmdShortDesc,
	Long:  utils.ExportAPICmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln("export-api called")

		accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(exportEnvironment, exportAPICmdUsername, exportAPICmdPassword)

		if preCommandErr == nil {
			resp := ExportAPI(exportAPIName, exportAPIVersion, apiManagerEndpoint, accessToken)

			// Print info on response
			utils.Logln("ResponseStatus: %v\n", resp.Status())
			utils.Logln("Error: %v\n", resp.Error())
			//fmt.Printf("Response Body: %v\n", resp.Body())

			if resp.StatusCode() == 200 {
				WriteToZip(exportAPIName, resp)

				numberOfAPIsExported, _, err := GetAPIList(exportAPIName, accessToken, apiManagerEndpoint)
				if err == nil {
					fmt.Println("Number of APIs exported: ", numberOfAPIsExported)
				} else {
					utils.HandleErrorAndExit("Error getting list of APIs", err)
				}

			} else if resp.StatusCode() == 500 {
				fmt.Println("Incorrect password")
			} else {
				utils.Logln(utils.LogPrefixWarning + resp.Status())
			}

		} else {
			utils.Logln(utils.LogPrefixError + preCommandErr.Error())
		}
	},
}

// Input name of the API, and the resty response (only 200 OK) returned from APIM Endpoint
// Exported API will be written to a zip file
func WriteToZip(exportAPIName string, resp *resty.Response) {
	// Write to file
	directory := utils.ExportedAPIDirectoryPath
	// create directory if it doesn't exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0777)
		// permission 777 : Everyone can read, write, and execute
	}
	filename := exportAPIName + ".zip"
	err := ioutil.WriteFile(directory+"/"+filename, resp.Body(), 0644)
	// permission 644 : Only the owner can read and write.. Everyone else can only read.
	if err != nil {
		utils.HandleErrorAndExit("Error creating zip archive", err)
	}
	fmt.Println("Succesfully exported and wrote to file")
}


// Input name of the API, version of the API, APIM Endpoint, and AccessToken
// Response will be returned after processing the request
func ExportAPI(name string, version string, url string, accessToken string) *resty.Response {
	// append '/' to the end if there isn't one already
	if string(url[len(url)-1]) != "/" {
		url += "/"
	}
	url += "export/apis"

	query := "?query=" + name
	url += query
	fmt.Println("ExportAPI: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTPS certificates
	utils.Logln("")
	resp, err := resty.R().
		SetHeaders(headers).
		Get(url)

	if err != nil {
		utils.HandleErrorAndExit("Error exporting API:" + name, err)
	}

	return resp
}

// Generated with Cobra
func init() {
	RootCmd.AddCommand(ExportAPICmd)
	ExportAPICmd.Flags().StringVarP(&exportAPIName, "name", "n", "", "Name of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportAPIVersion, "version", "v", "", "Version of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportEnvironment, "environment", "e", "", "Environment to which the API should be exported")

	ExportAPICmd.Flags().StringVarP(&exportAPICmdUsername, "username", "u", "", "Username")
	ExportAPICmd.Flags().StringVarP(&exportAPICmdPassword, "password", "p", "", "Password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ExportAPICmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ExportAPICmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
