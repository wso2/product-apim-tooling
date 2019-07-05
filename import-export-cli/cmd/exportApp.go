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
	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var exportAppName string
var exportAppOwner string
var exportAppWithKeys bool

//var flagExportAPICmdToken string
// ExportApp command related usage info
const exportAppCmdLiteral = "export-app"
const exportAppCmdShortDesc = "Export App"

const exportAppCmdLongDesc = "Export an Application from a specified  environment"

const exportAppCmdExamples = utils.ProjectName + ` ` + exportAppCmdLiteral + ` -n SampleApp -o admin -e dev
` + utils.ProjectName + ` ` + exportAppCmdLiteral + ` -n SampleApp -o admin -e prod
NOTE: Flag --name (-n) and --owner (-o) are mandatory`

// exportAppCmd represents the exportApp command
var ExportAppCmd = &cobra.Command{
	Use: exportAppCmdLiteral + " (--name <name-of-the-application> --owner <owner-of-the-application> --environment " +
		"<environment-from-which-the-app-should-be-exported>)",
	Short:   exportAppCmdShortDesc,
	Long:    exportAppCmdLongDesc,
	Example: exportAppCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAppCmdLiteral + " called")
		var appsExportDirectoryPath = filepath.Join(utils.ExportDirectory, utils.ExportedAppsDirName, cmdExportEnvironment)

		cred, err := getCredentials(cmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeExportAppCmd(cred, appsExportDirectoryPath)
	},
}

func executeExportAppCmd(credential credentials.Credential, appsExportDirectoryPath string) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)

	if preCommandErr == nil {
		adminEndpiont := utils.GetAdminEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		resp := getExportAppResponse(exportAppName, exportAppOwner, adminEndpiont, accessToken)

		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			WriteApplicationToZip(exportAppName, exportAppOwner, appsExportDirectoryPath, resp)
		} else if resp.StatusCode() == http.StatusUnauthorized {
			// 401 Unauthenticated request
			fmt.Println("Incorrect Password!")
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting Application:", resp.Status())
		}
	} else {
		// error exporting Application
		fmt.Println("Error exporting Application:" + preCommandErr.Error())
	}
}

// WriteApplicationToZip
// @param exportAppName : Name of the Application to be exported
// @param exportAppOwner : Owner of the Application to be exported
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported Application will be written to a zip file
func WriteApplicationToZip(exportAppName, exportAppOwner, zipLocationPath string,
	resp *resty.Response) {
	// Write to file
	//directory := filepath.Join(exportDirectory, exportEnvironment)
	// create directory if it doesn't exist
	if _, err := os.Stat(zipLocationPath); os.IsNotExist(err) {
		os.Mkdir(zipLocationPath, 0777)
		// permission 777 : Everyone can read, write, and execute
	}
	zipFilename := exportAppOwner + "_" + exportAppName + ".zip" // admin_testApp.zip
	pFile := filepath.Join(zipLocationPath, zipFilename)
	err := ioutil.WriteFile(pFile, resp.Body(), 0644)
	// permission 644 : Only the owner can read and write.. Everyone else can only read.
	if err != nil {
		utils.HandleErrorAndExit("Error creating zip archive", err)
	}
	fmt.Println("Successfully exported Application!")
	fmt.Println("Find the exported Application at " + pFile)
}

// ExportApp
// @param name : Name of the Application to be exported
// @param apimEndpoint : API Manager Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func getExportAppResponse(name, owner, adminEndpoint, accessToken string) *resty.Response {
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	query := "export/applications?appName=" + name + utils.SearchAndTag + "appOwner=" + owner

	if exportAppWithKeys {
		query += "&withKeys=true"
	}

	url := adminEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportApp: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(url, headers)

	if err != nil {
		utils.HandleErrorAndExit("Error exporting Application: "+name, err)
	}

	return resp
}

//init using Cobra
func init() {
	RootCmd.AddCommand(ExportAppCmd)
	ExportAppCmd.Flags().StringVarP(&exportAppName, "name", "n", "",
		"Name of the Application to be exported")
	ExportAppCmd.Flags().StringVarP(&exportAppOwner, "owner", "o", "",
		"Owner of the Application to be exported")
	ExportAppCmd.Flags().StringVarP(&cmdExportEnvironment, "environment", "e",
		"", "Environment to which the Application should be exported")
	ExportAppCmd.Flags().BoolVarP(&exportAppWithKeys, "withKeys", "",
		false, "Export keys for the application")
	_ = ExportAppCmd.MarkFlagRequired("environment")
}
