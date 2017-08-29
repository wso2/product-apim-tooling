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
	constants "github.com/menuka94/wso2apim-cli/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"

	"crypto/tls"
)

var exportAPIName string
var exportAPIVersion string
var exportEnvironment string
var ExportAPICmdUsername string
var ExportAPICmdPassword string

// ExportAPICmd represents the exportAPI command
var ExportAPICmd = &cobra.Command{
	Use:   "export-api (--name <name-of-the-api> --version <version-of-the-api> --environment <environment-from-which-the-api-should-be-exported>)",
	Short: utils.ExportAPICmdShortDesc,
	Long:  utils.ExportAPICmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exportAPI called")

		accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(exportEnvironment, ExportAPICmdUsername, ExportAPICmdPassword)

		if preCommandErr == nil {
			resp := ExportAPI(exportAPIName, exportAPIVersion, apiManagerEndpoint, accessToken)

			// Print info on response
			fmt.Printf("ResponseStatus: %v\n", resp.Status())
			fmt.Printf("Error: %v\n", resp.Error())
			//fmt.Printf("Response Body: %v\n", resp.Body())

			if resp.StatusCode() == 200 {
				// Write to file
				directory := "./exported"

				// create directory if it doesn't exist
				if _, err := os.Stat(directory); os.IsNotExist(err) {
					os.Mkdir(directory, 0777)
					// permission 777 : Everyone can read, write, and execute
				}

				filename := exportAPIName + ".zip"
				err := ioutil.WriteFile(directory+"/"+filename, resp.Body(), 0644)
				// permissoin 644 : Only the owner can read and write.. Everyone else can only read.
				if err != nil {
					fmt.Println("Error creating zip archive")
					panic(err)
				}
				fmt.Println("Succesfully exported and wrote to file")

				numberOfAPIsExported, _, err := GetAPIList(exportAPIName, accessToken, apiManagerEndpoint)
				if err == nil {
					fmt.Println("Number of APIs exported: ", numberOfAPIsExported)
				} else {
					fmt.Println("Error:")
					panic(err)
				}

			} else if resp.StatusCode() == 500 {
				fmt.Println("Incorrect password")
			}

		} else {
			log.Fatal("Error: ", preCommandErr)
		}
	},
}

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
	headers[constants.HeaderAuthorization] = constants.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[constants.HeaderAccept] = constants.HeaderValueApplicationZip

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // To bypass errors in HTTPS certificates
	resp, err := resty.R().
		SetHeaders(headers).
		Get(url)

	if err != nil {
		fmt.Println("Error exporting API:", name)
		panic(err)
	}

	return resp
}

func init() {
	RootCmd.AddCommand(ExportAPICmd)
	ExportAPICmd.Flags().StringVarP(&exportAPIName, "name", "n", "", "Name of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportAPIVersion, "version", "v", "", "Version of the API to be exported")
	ExportAPICmd.Flags().StringVarP(&exportEnvironment, "environment", "e", "", "Environment to which the API should be exported")

	ExportAPICmd.Flags().StringVarP(&ExportAPICmdUsername, "username", "u", "", "Username")
	ExportAPICmd.Flags().StringVarP(&ExportAPICmdPassword, "password", "p", "", "Password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ExportAPICmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ExportAPICmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
