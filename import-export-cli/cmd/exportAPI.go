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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/renstrom/dedent"
	"net/http"
	"path/filepath"
)

var exportAPIName string
var exportAPIVersion string
var exportProvider string
var runnigExportApiCommand bool

// ExportAPI command related usage info
const exportAPICmdLiteral = "export-api"
const exportAPICmdShortDesc = "Export API"

var exportAPICmdLongDesc = "Export APIs from an environment"

var exportAPICmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + exportAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -e dev --provider admin
		` + utils.ProjectName + ` ` + exportAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -e production --provider admin
		NOTE: all three flags (--name (-n), --version (-v), --provider (-r)) are mandatory
	`)

// ExportAPICmd represents the exportAPI command
var ExportAPICmd = &cobra.Command{
	Use: exportAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --environment " +
		"<environment-from-which-the-api-should-be-exported>)",
	Short: exportAPICmdShortDesc,
	Long:  exportAPICmdLongDesc + exportAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAPICmdLiteral + " called")
		var apisExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)
		executeExportAPICmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, apisExportDirectory)
	},
}

func executeExportAPICmd(mainConfigFilePath, envKeysAllFilePath, exportDirectory string) {
	runnigExportApiCommand = true
	b64encodedCredentials, preCommandErr :=
		utils.ExecutePreCommandWithBasicAuth(cmdExportEnvironment, cmdUsername, cmdPassword,
			mainConfigFilePath, envKeysAllFilePath)

	if preCommandErr == nil {
		apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(cmdExportEnvironment, mainConfigFilePath)
		resp := getExportApiResponse(exportAPIName, exportAPIVersion, exportProvider, apiImportExportEndpoint,
			b64encodedCredentials)

		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		apiZipLocationPath := filepath.Join(exportDirectory, cmdExportEnvironment)
		if resp.StatusCode() == http.StatusOK {
			WriteToZip(exportAPIName, exportAPIVersion, apiZipLocationPath, resp)
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
}

// WriteToZip
// @param exportAPIName : Name of the API to be exported
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported API will be written to a zip file
func WriteToZip(exportAPIName, exportAPIVersion, zipLocationPath string, resp *resty.Response) {
	// Write to file
	//directory := filepath.Join(exportDirectory, cmdExportEnvironment)
	// create directory if it doesn't exist
	if _, err := os.Stat(zipLocationPath); os.IsNotExist(err) {
		os.Mkdir(zipLocationPath, 0777)
		// permission 777 : Everyone can read, write, and execute
	}
	zipFilename := exportAPIName + "_" + exportAPIVersion + ".zip" // MyAPI_1.0.0.zip
	pFile := filepath.Join(zipLocationPath, zipFilename)
	err := ioutil.WriteFile(pFile, resp.Body(), 0644)
	// permission 644 : Only the owner can read and write.. Everyone else can only read.
	if err != nil {
		utils.HandleErrorAndExit("Error creating zip archive", err)
	}
	if(runnigExportApiCommand) {
		fmt.Println("Successfully exported API!")
		fmt.Println("Find the exported API at " + pFile)
	}
}

// ExportAPI
// @param name : Name of the API to be exported
// @param version : Version of the API to be exported
// @param apiImportExportEndpoint : API Import Export Endpoint for the environment
// @param  b64encodedCredentials: Base64 Encoded 'username:password'
// @return response Response in the form of *resty.Response
func getExportApiResponse(name, version, provider, apiImportExportEndpoint, b64encodedCredentials string) *resty.Response {
	apiImportExportEndpoint = utils.AppendSlashToString(apiImportExportEndpoint)
	query := "export-api?name=" + name + "&version=" + version + "&provider=" + provider

	url := apiImportExportEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportAPI: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBasicPrefix + " " + b64encodedCredentials
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(url, headers)

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
	ExportAPICmd.Flags().StringVarP(&exportProvider, "provider", "r", "",
		"Provider of the API")
	ExportAPICmd.Flags().StringVarP(&cmdExportEnvironment, "environment", "e",
		utils.DefaultEnvironmentName, "Environment to which the API should be exported")

	ExportAPICmd.Flags().StringVarP(&cmdUsername, "username", "u", "", "Username")
	ExportAPICmd.Flags().StringVarP(&cmdPassword, "password", "p", "", "Password")
}
