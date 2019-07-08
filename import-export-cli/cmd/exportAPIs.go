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
	"github.com/hashicorp/go-multierror"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const exportAPIsCmdLiteral = "export-apis"
const exportAPIsCmdShortDesc = "Export APIs for migration"

const exportAPIsCmdLongDesc = "Export all the APIs of a tenant from an APIM 2.6.0 environment environment, to be imported " +
	"into 3.0.0 environment"
const exportAPIsCmdExamples = utils.ProjectName + ` ` + exportAPIsCmdLiteral + ` -e production-2.6.0 -u wso2admin@wso2.org -p 12345 -t wso2.org -k --force
		` + utils.ProjectName + ` ` + exportAPIsCmdLiteral + ` -e production-2.6.0 -u admin -p admin -k`

var apiExportDir string
var apiListOffset int //from which index of API, the APIs will be fetched from APIM server
var count int32       // size of API list to be exported or number of  APIs left to be exported from last iteration
var apis []utils.API
var exportRelatedFilesPath string

//e.g. /home/samithac/.wso2apimcli/exported/migration/production-2.5/wso2-dot-org
var startFromBeginning bool
var isProcessCompleted bool
var startingApiIndexFromList int
var mainConfigFilePath string
var credential credentials.Credential

var ExportAPIsCmd = &cobra.Command{
	Use: exportAPIsCmdLiteral + " [--environment " +
		"<environment-from-which-artifacts-should-be-exported>] -u <user_name> -p <password> [-t <Tenant-domain-of-the-resources-to-be-exported>] [--force]",
	Short:   exportAPIsCmdShortDesc,
	Long:    exportAPIsCmdLongDesc,
	Example: exportAPIsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAPIsCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)

		cred, err := getCredentials(cmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		credential = cred
		executeExportAPIsCmd(artifactExportDirectory)
	},
}

// Do operations to export APIs for the migration into the directory passed as exportDirectory
// <export_directory> is the patch defined in main_config.yaml
// exportDirectory = <export_directory>/migration/
func executeExportAPIsCmd(exportDirectory string) {
	//create dir structure
	apiExportDir = createExportAPIsDirStructure(exportDirectory)
	exportRelatedFilesPath = filepath.Join(exportDirectory, cmdExportEnvironment,
		utils.GetMigrationExportTenantDirName(cmdResourceTenantDomain))
	//e.g. /home/samithac/.wso2apimcli/exported/migration/production-2.5/wso2-dot-org
	startFromBeginning = false
	isProcessCompleted = false

	fmt.Println("\nExporting APIs for the migration...")
	if cmdForceStartFromBegin {
		startFromBeginning = true
	}

	if (utils.IsFileExist(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName))) && !startFromBeginning {
		prepareResumption()
	} else {
		prepareStartFromBeginning()
	}

	exportAPIs()
}

// Do the API exportation
func exportAPIs() {
	if count == 0 {
		fmt.Println("No APIs available to be exported..!")
	} else {
		var counterSuceededAPIs = 0
		for count > 0 {
			utils.Logln("Found ", count, "of APIs to be exported in the iteration beginning with the offset #"+
				strconv.Itoa(apiListOffset)+". Maximum limit of APIs exported in single iteration is "+
				strconv.Itoa(utils.MaxAPIsToExportOnce))
			//get basic Auth credentials
			b64encodedCredentials := credentials.GetBasicAuth(credential)
			for i := startingApiIndexFromList; i < len(apis); i++ {
				exportAPIName := apis[i].Name
				exportAPIVersion := apis[i].Version
				exportApiProvider := apis[i].Provider
				apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
				resp, err := getExportApiResponse(exportAPIName, exportAPIVersion, exportApiProvider, exportAPIFormat, apiImportExportEndpoint,
					b64encodedCredentials, exportAPIPreserveStatus)
				if err != nil {
					utils.HandleErrorAndExit("Error exporting", err)
				}

				if resp.StatusCode() == http.StatusOK {
					utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
					WriteToZip(exportAPIName, exportAPIVersion, apiExportDir, resp)
					//write on last-succeeded-api.log
					counterSuceededAPIs++
					utils.WriteLastSuceededAPIFileData(exportRelatedFilesPath, apis[i])
				} else {
					fmt.Println("Error exporting API:", exportAPIName, "-", exportAPIVersion, " of Provider:", exportApiProvider)
					utils.PrintErrorResponseAndExit(resp)
				}
			}
			fmt.Println("Batch of " + cast.ToString(count) + " APIs exported successfully..!")

			apiListOffset += utils.MaxAPIsToExportOnce
			count, apis = getAPIList()
			startingApiIndexFromList = 0
			if len(apis) > 0 {
				utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
					exportRelatedFilesPath, apiListOffset)
			}
		}
		fmt.Println("\nTotal number of APIs exported : " + cast.ToString(counterSuceededAPIs))
		fmt.Println("API export path: " + apiExportDir)
		fmt.Println("\nCommand: export-apis execution completed !")
	}
}

//  Prepare resumption of previous-halted export-apis operation
func prepareResumption() {
	var lastSuceededAPI utils.API
	lastSuceededAPI = utils.ReadLastSucceededAPIFileData(exportRelatedFilesPath)
	var migrationApisExportMetadata utils.MigrationApisExportMetadata
	err := migrationApisExportMetadata.ReadMigrationApisExportMetadataFile(filepath.Join(exportRelatedFilesPath,
		utils.MigrationAPIsExportMetadataFileName))
	if err != nil {
		utils.HandleErrorAndExit("Error loading metadata for resume from"+filepath.Join(exportRelatedFilesPath,
			utils.MigrationAPIsExportMetadataFileName), err)
	}
	apis = migrationApisExportMetadata.ApiListToExport
	apiListOffset = migrationApisExportMetadata.ApiListOffset
	startingApiIndexFromList = getLastSuceededApiIndex(lastSuceededAPI) + 1

	//find count of APIs left to be exported
	var lastSucceededAPInumber = getLastSuceededApiIndex(lastSuceededAPI) + 1
	count = int32(len(apis) - lastSucceededAPInumber)

	if count == 0 {
		//last iteration had been completed successfully but operation had halted at that point.
		//So get the next set of APIs for next iteration
		apiListOffset += utils.MaxAPIsToExportOnce
		startingApiIndexFromList = 0
		count, apis = getAPIList()
		if len(apis) > 0 {
			utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
				exportRelatedFilesPath, apiListOffset)
		} else {
			fmt.Println("Command: export-apis execution completed !")
			os.Exit(1)
		}
	}
}

// get the index of the finally (successfully) exported API from the list of APIs listed in migration-apis-export-metadata.yaml
func getLastSuceededApiIndex(lastSuceededApi utils.API) int {
	for i := 0; i < len(apis); i++ {
		if (apis[i].Name == lastSuceededApi.Name) &&
			(apis[i].Provider == lastSuceededApi.Provider) &&
			(apis[i].Version == lastSuceededApi.Version) {
			return i
		}
	}
	return -1
}

// Delete directories where the APIs are exported, reset the indexes, get first API list and write the
// migration-apis-export-metadata.yaml file
func prepareStartFromBeginning() {
	fmt.Println("Cleaning all the previously exported APIs of the given target tenant, in the given environment if " +
		"any, and prepare to export APIs from beginning")
	//cleaning existing old files (if exists) related to exportation
	var err error
	err = multierror.Append(err,
		utils.RemoveDirectoryIfExists(filepath.Join(exportRelatedFilesPath, utils.ExportedApisDirName)),
		utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.MigrationAPIsExportMetadataFileName)),
		utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName)),
	)
	if err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}

	apiListOffset = 0
	startingApiIndexFromList = 0
	count, apis = getAPIList()
	//write  migration-apis-export-metadata.yaml file
	utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
		exportRelatedFilesPath, apiListOffset)
}

// Create the required directory structure to save the exported APIs
func createExportAPIsDirStructure(artifactExportDirectory string) string {
	var resourceTenantDirName = utils.GetMigrationExportTenantDirName(cmdResourceTenantDomain)

	var createDirError error
	createDirError = utils.CreateDirIfNotExist(artifactExportDirectory)

	migrationsArtifactsEnvPath := filepath.Join(artifactExportDirectory, cmdExportEnvironment)
	migrationsArtifactsEnvTenantPath := filepath.Join(migrationsArtifactsEnvPath, resourceTenantDirName)
	migrationsArtifactsEnvTenantApisPath := filepath.Join(migrationsArtifactsEnvTenantPath, utils.ExportedApisDirName)

	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvPath)
	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvTenantPath)

	if dirExists, _ := utils.IsDirExists(migrationsArtifactsEnvTenantApisPath); dirExists {
		if cmdForceStartFromBegin {
			utils.RemoveDirectory(migrationsArtifactsEnvTenantApisPath)
			createDirError = utils.CreateDir(migrationsArtifactsEnvTenantApisPath)
		}
	} else {
		createDirError = utils.CreateDir(migrationsArtifactsEnvTenantApisPath)
	}

	if createDirError != nil {
		utils.HandleErrorAndExit("Error in creating directory structure for the API export for migration .",
			createDirError)
	}
	return migrationsArtifactsEnvTenantApisPath
}

// Get the list of APIs from the defined offset index, upto the limit of constant value utils.MaxAPIsToExportOnce
func getAPIList() (count int32, apis []utils.API) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)
	if preCommandErr == nil {
		apiListEndpoint := utils.GetApiListEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		apiListEndpoint += "?limit=" + strconv.Itoa(utils.MaxAPIsToExportOnce) + "&offset=" + strconv.Itoa(apiListOffset)
		if cmdResourceTenantDomain != "" {
			apiListEndpoint += "&tenantDomain=" + cmdResourceTenantDomain
		}
		count, apis, err := GetAPIList("", accessToken, apiListEndpoint)
		if err == nil {
			return count, apis
		} else {
			utils.HandleErrorAndExit(utils.LogPrefixError+"Getting List of APIs.", utils.GetHttpErrorResponse(err))
		}
	} else {
		utils.HandleErrorAndExit(utils.LogPrefixError+"Error in getting access token for user while getting "+
			"the list of APIs: ", preCommandErr)
	}
	return 0, nil
}

func init() {
	RootCmd.AddCommand(ExportAPIsCmd)
	ExportAPIsCmd.Flags().StringVarP(&cmdExportEnvironment, "environment", "e",
		"", "Environment to which the API should be exported")
	ExportAPIsCmd.Flags().StringVarP(&cmdResourceTenantDomain, "tenant", "t", "",
		"Tenant domain of the resources to be exported")
	ExportAPIsCmd.PersistentFlags().BoolVarP(&cmdForceStartFromBegin, "force", "", false,
		"Clean all the previously exported APIs of the given target tenant, in the given environment if "+
			"any, and to export APIs from beginning")
	ExportAPIsCmd.Flags().BoolVarP(&exportAPIPreserveStatus, "preserveStatus", "", true,
		"Preserve API status when exporting. Otherwise API will be exported in CREATED status")
	ExportAPIsCmd.Flags().StringVarP(&exportAPIFormat, "format", "", "json", "File format of exported archive")
	_ = ExportAPIsCmd.MarkFlagRequired("environment")
}
