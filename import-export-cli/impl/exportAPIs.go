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

package impl

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/spf13/cast"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var apiExportDir string
var apiListOffset int //from which index of API, the APIs will be fetched from APIM server
var count int32       // size of API list to be exported or number of  APIs left to be exported from last iteration
var apis []utils.API
var exportRelatedFilesPath string
var exportAPIsFormat string

var startingApiIndexFromList int
var mainConfigFilePath string

//  Prepare resumption of previous-halted export-apis operation
func PrepareResumption(credential credentials.Credential, exportRelatedFilesPath, cmdResourceTenantDomain, cmdUsername, cmdExportEnvironment string) {
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
		count, apis = getAPIList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
		if len(apis) > 0 {
			utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
				exportRelatedFilesPath, apiListOffset)
		} else {
			fmt.Println("Command: export-apis execution completed !")
		}
	}
}

// Delete directories where the APIs are exported, reset the indexes, get first API list and write the
// migration-apis-export-metadata.yaml file
func PrepareStartFromBeginning(credential credentials.Credential, exportRelatedFilesPath, cmdResourceTenantDomain, cmdUsername, cmdExportEnvironment string) {
	fmt.Println("Cleaning all the previously exported APIs of the given target tenant, in the given environment if " +
		"any, and prepare to export APIs from beginning")
	//cleaning existing old files (if exists) related to exportation
	if err := utils.RemoveDirectoryIfExists(filepath.Join(exportRelatedFilesPath, utils.ExportedApisDirName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}
	if err := utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.MigrationAPIsExportMetadataFileName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}
	if err := utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}

	apiListOffset = 0
	startingApiIndexFromList = 0
	count, apis = getAPIList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
	//write  migration-apis-export-metadata.yaml file
	utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
		exportRelatedFilesPath, apiListOffset)
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

// Get the list of APIs from the defined offset index, upto the limit of constant value utils.MaxAPIsToExportOnce
func getAPIList(credential credentials.Credential, cmdExportEnvironment, cmdResourceTenantDomain string) (count int32, apis []utils.API) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)
	if preCommandErr == nil {
		apiListEndpoint := utils.GetApiListEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		apiListEndpoint += "?limit=" + strconv.Itoa(utils.MaxAPIsToExportOnce) + "&offset=" + strconv.Itoa(apiListOffset)
		if cmdResourceTenantDomain != "" {
			apiListEndpoint += "&tenantDomain=" + cmdResourceTenantDomain
		}
		count, apis, err := GetAPIList(accessToken, apiListEndpoint, "", "")
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

// Do the API exportation
func ExportAPIs(credential credentials.Credential, exportRelatedFilesPath, cmdExportEnvironment, cmdResourceTenantDomain, exportAPIsFormat, cmdUsername, apiExportDir string,
	exportAPIPreserveStatus, runningExportApiCommand bool) {
	if count == 0 {
		fmt.Println("No APIs available to be exported..!")
	} else {
		var counterSuceededAPIs = 0
		for count > 0 {
			utils.Logln(utils.LogPrefixInfo+"Found ", count, "of APIs to be exported in the iteration beginning with the offset #"+
				strconv.Itoa(apiListOffset)+". Maximum limit of APIs exported in single iteration is "+
				strconv.Itoa(utils.MaxAPIsToExportOnce))
			accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)
			if preCommandErr == nil {
				for i := startingApiIndexFromList; i < len(apis); i++ {
					exportAPIName := apis[i].Name
					exportAPIVersion := apis[i].Version
					exportApiProvider := apis[i].Provider
					resp, err := ExportAPIFromEnv(accessToken, exportAPIName, exportAPIVersion, exportApiProvider, exportAPIsFormat,
						cmdExportEnvironment, exportAPIPreserveStatus)
					if err != nil {
						utils.HandleErrorAndExit("Error exporting", err)
					}

					if resp.StatusCode() == http.StatusOK {
						utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
						WriteToZip(exportAPIName, exportAPIVersion, apiExportDir, runningExportApiCommand, resp)
						//write on last-succeeded-api.log
						counterSuceededAPIs++
						utils.WriteLastSuceededAPIFileData(exportRelatedFilesPath, apis[i])
					} else {
						fmt.Println("Error exporting API:", exportAPIName, "-", exportAPIVersion, " of Provider:", exportApiProvider)
						utils.PrintErrorResponseAndExit(resp)
					}
				}
			} else {
				// error getting OAuth tokens
				fmt.Println("Error getting OAuth Tokens : " + preCommandErr.Error())
			}
			fmt.Println("Batch of " + cast.ToString(count) + " APIs exported successfully..!")

			apiListOffset += utils.MaxAPIsToExportOnce
			count, apis = getAPIList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
			startingApiIndexFromList = 0
			if len(apis) > 0 {
				utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
					exportRelatedFilesPath, apiListOffset)
			}
		}
		fmt.Println("\nTotal number of APIs exported: " + cast.ToString(counterSuceededAPIs))
		fmt.Println("API export path: " + apiExportDir)
		fmt.Println("\nCommand: export-apis execution completed !")
	}
}

// Create the required directory structure to save the exported APIs
func CreateExportAPIsDirStructure(artifactExportDirectory, cmdResourceTenantDomain, cmdExportEnvironment string, cmdForceStartFromBegin bool) string {
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
