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
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"strconv"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/go-resty/resty"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"net/http"
	"path/filepath"
)

var exportAPIName string
var exportAPIVersion string
var exportProvider string
var exportAPIPreserveStatus bool
var exportAPIFormat string
var runnigExportApiCommand bool

// ExportAPI command related usage info
const exportAPICmdLiteral = "export-api"
const exportAPICmdShortDesc = "Export API"

const exportAPICmdLongDesc = "Export APIs from an environment"

const exportAPICmdExamples = utils.ProjectName + ` ` + exportAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -r admin -e dev
` + utils.ProjectName + ` ` + exportAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -r admin -e production
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory`

// ExportAPICmd represents the exportAPI command
var ExportAPICmd = &cobra.Command{
	Use: exportAPICmdLiteral + " (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment " +
		"<environment-from-which-the-api-should-be-exported>)",
	Short:   exportAPICmdShortDesc,
	Long:    exportAPICmdLongDesc,
	Example: exportAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAPICmdLiteral + " called")
		var apisExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)

		cred, err := getCredentials(cmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportAPICmd(cred, apisExportDirectory)
	},
}

func executeExportAPICmd(credential credentials.Credential, exportDirectory string) {
	runnigExportApiCommand = true
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)

	if preCommandErr == nil {
		adminEndpoint := utils.GetAdminEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		resp, err := getExportApiResponse(exportAPIName, exportAPIVersion, exportProvider, exportAPIFormat, adminEndpoint,
			accessToken, exportAPIPreserveStatus)
		if err != nil {
			utils.HandleErrorAndExit("Error while exporting", err)
		}
		// Print info on response
		utils.Logf(utils.LogPrefixInfo + "ResponseStatus: %v\n", resp.Status())
		apiZipLocationPath := filepath.Join(exportDirectory, cmdExportEnvironment)
		if resp.StatusCode() == http.StatusOK {
			WriteToZip(exportAPIName, exportAPIVersion, apiZipLocationPath, resp)
		} else if resp.StatusCode() == http.StatusInternalServerError {
			// 500 Internal Server Error
			fmt.Println(string(resp.Body()))
		} else {
			// neither 200 nor 500
			fmt.Println("Error exporting API:", resp.Status(), "\n", string(resp.Body()))
		}
	} else {
		// error exporting Api
		fmt.Println("Error getting OAuth tokens while exporting API:" + preCommandErr.Error())
	}
}

// WriteToZip
// @param exportAPIName : Name of the API to be exported
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported API will be written to a zip file
func WriteToZip(exportAPIName, exportAPIVersion, zipLocationPath string, resp *resty.Response) {
	zipFilename := exportAPIName + "_" + exportAPIVersion + ".zip" // MyAPI_1.0.0.zip
	// Writes the REST API response to a temporary zip file
	tempZipFile, err := utils.WriteResponseToTempZip(zipFilename, resp)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the temporary zip file to store the exported API" , err)
	}

	err = utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archive: " + zipLocationPath, err)
	}
	exportedFinalZip := filepath.Join(zipLocationPath, zipFilename)
	// Add api_params.yaml file inside the zip and create a new zip file in exportedFinalZip location
	err = impl.IncludeParamsFileToZip(tempZipFile, exportedFinalZip, utils.ParamFileAPI)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive", err)
	}

	// Output the final zip file location.
	if runnigExportApiCommand {
		fmt.Println("Successfully exported API!")
		fmt.Println("Find the exported API at " + exportedFinalZip)
	}
}

// ExportAPI
// @param name : Name of the API to be exported
// @param version : Version of the API to be exported
// @param provider : Provider of the API 
// @param adminEndpoint : API Manager Admin Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func getExportApiResponse(name, version, provider, format, adminEndpoint, accessToken string, preserveStatus bool) (*resty.Response, error) {
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	query := "export/api?name=" + name + "&version=" + version + "&providerName=" + provider +
		"&preserveStatus=" + strconv.FormatBool(preserveStatus)
	if format != "" {
		query += "&format=" + format
	}

	url := adminEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportAPI: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(url, headers)

	if err != nil {
		return nil, err
	}
	
	return resp, nil
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
		"", "Environment to which the API should be exported")
	ExportAPICmd.Flags().BoolVarP(&exportAPIPreserveStatus, "preserveStatus", "", true,
		"Preserve API status when exporting. Otherwise API will be exported in CREATED status")
	ExportAPICmd.Flags().StringVarP(&exportAPIFormat, "format", "", utils.DefaultExportFormat, "File format of exported archive(json or yaml)")
	_ = ExportAPICmd.MarkFlagRequired("name")
	_ = ExportAPICmd.MarkFlagRequired("version")
	_ = ExportAPICmd.MarkFlagRequired("environment")
}
