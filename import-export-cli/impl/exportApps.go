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

var appsExportDir string
var appListOffset int // from which index of App, the Apps will be fetched from APIM server
var appCount int32    // size of App list to be exported or number of  Apps left to be exported from last iteration
var apps []utils.Application
var exportAppsRelatedFilesPath string
var exportAppsFormat string

var startingAppIndexFromList int

// Prepare resumption of previous-halted export-apps operation
func PrepareResumptionForApps(credential credentials.Credential, exportAppsRelatedFilesPath, cmdResourceTenantDomain, cmdUsername, cmdExportEnvironment string) {
	var lastSuceededApp utils.Application
	lastSuceededApp = utils.ReadLastSucceededAppFileData(exportAppsRelatedFilesPath)
	var migrationAppsExportMetadata utils.MigrationAppsExportMetadata
	err := migrationAppsExportMetadata.ReadMigrationAppsExportMetadataFile(filepath.Join(exportAppsRelatedFilesPath,
		utils.MigrationAppsExportMetadataFileName))
	if err != nil {
		utils.HandleErrorAndExit("Error loading metadata for resume from"+filepath.Join(exportAppsRelatedFilesPath,
			utils.MigrationAppsExportMetadataFileName), err)
	}
	apps = migrationAppsExportMetadata.AppListToExport
	appListOffset = migrationAppsExportMetadata.AppListOffset
	startingAppIndexFromList = getLastSuceededAppIndex(lastSuceededApp) + 1

	// Find count of Apps left to be exported
	appCount = int32(len(apps) - startingAppIndexFromList)

	if appCount == 0 {
		// Last iteration had been completed successfully but operation had halted at that point.
		// So get the next set of Apps for next iteration
		startingAppIndexFromList = 0
		appCount, apps = getAppList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
		if len(apps)-startingAppIndexFromList > 0 {
			utils.WriteMigrationAppsExportMetadataFile(apps, cmdResourceTenantDomain, cmdUsername,
				exportAppsRelatedFilesPath, appListOffset)
		} else {
			fmt.Println("Command: export apps execution completed !")
		}
	}
}

// Delete directories where the Apps are exported, reset the indexes, get first App list and write the
// migration-apps-export-metadata.yaml file
func PrepareStartAppsFromBeginning(credential credentials.Credential, exportAppsRelatedFilesPath, cmdResourceTenantDomain, cmdUsername, cmdExportEnvironment string) {
	fmt.Println("Cleaning all the previously exported Apps of the given target tenant, in the given environment if " +
		"any, and prepare to export Apps from beginning")
	// Cleaning existing old files (if exists) related to exportation
	if err := utils.RemoveDirectoryIfExists(filepath.Join(exportAppsRelatedFilesPath, utils.ExportedAppsDirName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}
	if err := utils.RemoveFileIfExists(filepath.Join(exportAppsRelatedFilesPath, utils.MigrationAppsExportMetadataFileName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}
	if err := utils.RemoveFileIfExists(filepath.Join(exportAppsRelatedFilesPath, utils.LastSucceededAppFileName)); err != nil {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", err)
	}

	appListOffset = 0
	startingAppIndexFromList = 0
	appCount, apps = getAppList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
	// Write  migration-apps-export-metadata.yaml file
	utils.WriteMigrationAppsExportMetadataFile(apps, cmdResourceTenantDomain, cmdUsername, exportAppsRelatedFilesPath,
		appListOffset)
}

// Get the index of the finally (successfully) exported App from the list of Apps listed in migration-apps-export-metadata.yaml
func getLastSuceededAppIndex(lastSuceededApp utils.Application) int {
	for i := 0; i < len(apps); i++ {
		if (apps[i].Name == lastSuceededApp.Name) &&
			(apps[i].Owner == lastSuceededApp.Owner) {
			return i
		}
	}
	return -1
}

// Get the list of Apps from the defined offset index, upto the limit of constant value utils.MaxAppsToExportOnce
func getAppList(credential credentials.Credential, cmdExportEnvironment, cmdResourceTenantDomain string) (appCount int32, apps []utils.Application) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)
	if preCommandErr == nil {
		appListEndpoint := utils.GetAdminApplicationListEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		appListEndpoint += "?limit=" + strconv.Itoa(utils.MaxAppsToExportOnce) + "&offset=" + strconv.Itoa(appListOffset)
		if cmdResourceTenantDomain != "" {
			appListEndpoint += "&tenantDomain=" + cmdResourceTenantDomain
		}
		appCount, apps, err := GetApplicationList(accessToken, appListEndpoint, "", "")
		if err == nil {
			return appCount, apps
		} else {
			utils.HandleErrorAndExit(utils.LogPrefixError+"Getting List of Apps.", utils.GetHttpErrorResponse(err))
		}
	} else {
		utils.HandleErrorAndExit(utils.LogPrefixError+"Error in getting access token for user while getting "+
			"the list of Apps: ", preCommandErr)
	}
	return 0, nil
}

// Do the App exportation
func ExportApps(credential credentials.Credential, exportAppsRelatedFilesPath, cmdExportEnvironment, cmdResourceTenantDomain,
	exportAppsFormat, cmdUsername, appExportDir string, exportAppsWithKeys bool) {
	if appCount == 0 {
		fmt.Println("No Apps available to be exported..!")
	} else {
		var counterSuceededApps = 0
		for appCount > 0 {
			utils.Logln(utils.LogPrefixInfo+"Found ", appCount, "of Apps to be exported in the iteration beginning with the offset #"+
				strconv.Itoa(appListOffset)+". Maximum limit of Apps exported in single iteration is "+
				strconv.Itoa(utils.MaxAppsToExportOnce))
			accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)
			if preCommandErr == nil {
				for i := startingAppIndexFromList; i < len(apps); i++ {
                    exportAppAndWriteToZip(apps[i], accessToken, cmdExportEnvironment, appExportDir,
                        exportAppsRelatedFilesPath, exportAppsFormat, exportAppsWithKeys)
                    counterSuceededApps++
				}
			} else {
				// Error getting OAuth tokens
				fmt.Println("Error getting OAuth Tokens : " + preCommandErr.Error())
			}
			fmt.Println("Batch of " + cast.ToString(appCount) + " Apps exported successfully..!")

			appListOffset += utils.MaxAppsToExportOnce
			appCount, apps = getAppList(credential, cmdExportEnvironment, cmdResourceTenantDomain)
			startingAppIndexFromList = 0
			if len(apps) > 0 {
				utils.WriteMigrationAppsExportMetadataFile(apps, cmdResourceTenantDomain, cmdUsername,
					exportAppsRelatedFilesPath, appListOffset)
			}
		}
		fmt.Println("\nTotal number of Apps exported: " + cast.ToString(counterSuceededApps))
		fmt.Println("App export path: " + appExportDir)
		fmt.Println("\nCommand: export-apps execution completed !")
	}
}

// Export the App and archive to zip format
func exportAppAndWriteToZip(app utils.Application, accessToken, cmdExportEnvironment, appExportDir,
	exportAppsRelatedFilesPath, exportAppsFormat string, exportAppsWithKeys bool) {

	exportAppName := app.Name
	exportAppOwner := app.Owner
	resp, err := ExportAppFromEnv(accessToken, exportAppName, exportAppOwner, exportAppsFormat, cmdExportEnvironment,
	 exportAppsWithKeys)
	if err != nil {
		utils.HandleErrorAndExit("Error exporting", err)
	}

	if resp.StatusCode() == http.StatusOK {
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		WriteApplicationToZip(exportAppName, exportAppOwner, appExportDir, resp)
		// write on last-succeeded-app.log
		utils.WriteLastSuceededAppFileData(exportAppsRelatedFilesPath, app)
	} else {
		fmt.Println("Error exporting App:", exportAppName, " of Owner:", exportAppOwner)
		utils.PrintErrorResponseAndExit(resp)
	}
}

// Create the required directory structure to save the exported Apps
func CreateExportAppsDirStructure(artifactExportDirectory, cmdResourceTenantDomain, cmdExportEnvironment string, cmdForceStartFromBegin bool) string {
	var resourceTenantDirName = utils.GetMigrationExportTenantDirName(cmdResourceTenantDomain)

	var createDirError error
	createDirError = utils.CreateDirIfNotExist(artifactExportDirectory)

	migrationsArtifactsEnvPath := filepath.Join(artifactExportDirectory, cmdExportEnvironment)
	migrationsArtifactsEnvTenantPath := filepath.Join(migrationsArtifactsEnvPath, resourceTenantDirName)
	migrationsArtifactsEnvTenantAppsPath := filepath.Join(migrationsArtifactsEnvTenantPath, utils.ExportedAppsDirName)

	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvPath)
	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvTenantPath)

	if dirExists, _ := utils.IsDirExists(migrationsArtifactsEnvTenantAppsPath); dirExists {
		if cmdForceStartFromBegin {
			utils.RemoveDirectory(migrationsArtifactsEnvTenantAppsPath)
			createDirError = utils.CreateDir(migrationsArtifactsEnvTenantAppsPath)
		}
	} else {
		createDirError = utils.CreateDir(migrationsArtifactsEnvTenantAppsPath)
	}

	if createDirError != nil {
		utils.HandleErrorAndExit("Error in creating directory structure for the Apps export for migration .",
			createDirError)
	}
	return migrationsArtifactsEnvTenantAppsPath
}
