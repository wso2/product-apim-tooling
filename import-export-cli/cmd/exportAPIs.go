package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"path/filepath"
	"fmt"
	"net/http"
				"strconv"
	)

const exportAPIsCmdLiteral = "export-apis"
const exportAPIsCmdShortDesc = "Export APIs"

var exportAPIsCmdLongDesc = "Export all the APIs of the tenant from an APIM 2.x environment environment, to be imported " +
	"into 3.0.0 environment"
var exportAPIsCmdExamples = ""

var ExportAPIsCmd = &cobra.Command{
	Use: exportAPIsCmdLiteral + " (--environment " +
		"<environment-from-which-artifacts-should-be-exported>)",
	Short: exportAPIsCmdShortDesc,
	Long:  exportAPIsCmdLongDesc + exportAPIsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		//use log package for in-detail logging ; https://golangcode.com/add-line-numbers-to-log-output/
		utils.Logln(utils.LogPrefixInfo + exportAPIsCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)
		executeExportAPIsCmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, artifactExportDirectory)
	},
}

func executeExportAPIsCmd(mainConfigFilePath, envKeysAllFilePath, exportDirectory string) {
	//create dir structure
	var apiExportDir = createExportAPIsDirStructure(exportDirectory)
	var offset= 0 // this is the begining API index
	var count int32 = 0 // size of API list to be exported / number of  APIs left to be exported from last iteration
	var apis[] utils.API
	var exportRelatedFilesPath = filepath.Join(exportDirectory, cmdExportEnvironment,
		utils.GetMigrationExportTenantDirName(cmdResourceTenantDomain))
		//e.g. /home/samithac/.wso2apimcli/exported/migration/production-2.5/wso2-dot-org
	var startFromBeginning = false

	var lastSuceededAPI utils.API
	var isProcessCompleted = false
	var iterationNo int = 1
	var startingApiIndexFromList = 0

	fmt.Println("\nExporting APIs for the migration to APIM 3.0.0")

	/*

	with --force
	-------
	1- Delete last-succeeded-api.log, migration-apis-export-metadata.json files and tenant-x/apis dir if available.
		starting_from_beginning=true

			-
	Without --force set
	1 - Check if available last-succeeded-api.log file. If yes, then check for line 'COMPLETED'
		- If 'COMPLETED' is written >> stop process > Show message that already finished the process
		- else if 'COMPLETED' is not written >> the exportation has halted by an error. So read the first line content.
		- Get the last succeeded index and iteration number.  set
				offset +=1
				iteration = iteration number readed from file
	2- If Not available last-succeeded-api.log, set a variable
			starting_from_beginning= true
			if, tenant-x/apis/ dir exists
				show message > "Has to use -f option. Cannot resume the previous export process"
				exit program
	....
	...
	 */

	if (cmdForceStartFromBegin) {
		startFromBeginning = true
	}

	if ((utils.IsFileExist(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName))) && !startFromBeginning) {
		// set offset, iteration from files
		iterationNo, lastSuceededAPI, isProcessCompleted = utils.ReadLastSuceededAPIFileData(exportRelatedFilesPath);
		//read  migration-apis-export-metadata.json file and set
		// 	- apis[]
		// 	- startingApiIndexFromList
		//  - count = api_list_size_in_file - (indexoflast_suceeded+1)
	} else { // start from beginning

		//cleaning existing old files (if exists) related to exportation
		error := utils.RemoveDirectoryIfExists(filepath.Join(exportRelatedFilesPath, utils.ExportedApisDirName))
		error = utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.MigrationApisExportMetadataFileName))
		error = utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName))
		if (error != nil) {
			utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
				"exportation", error)
		}

		offset = 0; // this is the beginning API index
		iterationNo = 1;
		startingApiIndexFromList = 0;
		count, apis = getAPIList(mainConfigFilePath, envKeysAllFilePath, offset)
		//write  migration-apis-export-metadata.yaml file
	}

	if (count == 0) {
		fmt.Println("No APIs available to be exported..!")
		/*
		1- WRITE ON last-succeeded-api.log 'COMPLETED'
		2-

		*/
	} else {
		for (count > 0) {
			fmt.Println("Found ", count, "of APIs to be exported in the iteration #" + strconv.Itoa(iterationNo))
			//get basic Auth credentials
			b64encodedCredentials, preCommandErr :=
				utils.ExecutePreCommandWithBasicAuth(cmdExportEnvironment, cmdUsername, cmdPassword,
					mainConfigFilePath, envKeysAllFilePath)
			if preCommandErr != nil {
				utils.Logln("Error generating base64 encoded credentials for Basic Authentication:")
				//return preCommandErr
			}
			for i := startingApiIndexFromList; i < len(apis); i++ {
				exportAPIName := apis[i].Name
				exportAPIVersion := apis[i].Version
				exportApiProvider := apis[i].Provider
				apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(cmdExportEnvironment, mainConfigFilePath)
				resp := getExportApiResponse(exportAPIName, exportAPIVersion, exportApiProvider, apiImportExportEndpoint,
					b64encodedCredentials) // TODO Handle errors in getExportApiResponse()

				// Print info on response
				utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

				if resp.StatusCode() == http.StatusOK {
					//artifactExportDirectory >> use a suitable dynamic folder name instead of 'Migration-1' and also change
					// WriteToZip function
					WriteToZip(exportAPIName, exportAPIVersion, apiExportDir, resp)
				} else {
					fmt.Println("Error exporting API:", exportAPIName, "-", exportAPIVersion, " of Provider:", exportApiProvider)
					//return errors.New(resp.String())
				}
			}

			//iterate apis
			//Is it required to stop exporting process once an error in single API occurs?
			for _, api := range apis {
				exportAPIName := api.Name
				exportAPIVersion := api.Version
				exportApiProvider := api.Provider
				apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(cmdExportEnvironment, mainConfigFilePath)
				resp := getExportApiResponse(exportAPIName, exportAPIVersion, exportApiProvider, apiImportExportEndpoint,
					b64encodedCredentials) // TODO Handle errors in getExportApiResponse()

				// Print info on response
				utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

				if resp.StatusCode() == http.StatusOK {
					//artifactExportDirectory >> use a suitable dynamic folder name instead of 'Migration-1' and also change
					// WriteToZip function
					WriteToZip(exportAPIName, exportAPIVersion, apiExportDir, resp)
				} else {
					fmt.Println("Error exporting API:", exportAPIName, "-", exportAPIVersion, " of Provider:", exportApiProvider)
					//return errors.New(resp.String())
				}
			}
			offset = utils.MaxAPIsToExportOnce * iterationNo
			iterationNo++
			count, apis = getAPIList(mainConfigFilePath, envKeysAllFilePath, offset)
		}
	}
}


func createExportAPIsDirStructure(artifactExportDirectory string) string {
	//create required directory structure
	var resourceTenantDirName = utils.GetMigrationExportTenantDirName(cmdResourceTenantDomain)

	var createDirError error
	createDirError = utils.CreateDirIfNotExist(artifactExportDirectory)

	migrationsArtifactsEnvPath := filepath.Join(artifactExportDirectory, cmdExportEnvironment)
	migrationsArtifactsEnvTenantPath := filepath.Join(migrationsArtifactsEnvPath, resourceTenantDirName)
	migrationsArtifactsEnvTenantApisPath := filepath.Join(migrationsArtifactsEnvTenantPath, utils.ExportedApisDirName)

	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvPath)
	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvTenantPath)

	if dirExists,_ := utils.IsDirExists(migrationsArtifactsEnvTenantApisPath); dirExists {
		if (cmdForceStartFromBegin) {
			utils.RemoveDirectory(migrationsArtifactsEnvTenantApisPath)
			createDirError = utils.CreateDir(migrationsArtifactsEnvTenantApisPath)
		}
	} else {
		createDirError = utils.CreateDir(migrationsArtifactsEnvTenantApisPath)
	}

	if (createDirError != nil) {
		utils.HandleErrorAndExit("Error in creating directory structure for the API export for migration .",
			createDirError)
	}
	return migrationsArtifactsEnvTenantApisPath
}

func getAPIList(mainConfigFilePath string, envKeysAllFilePath string, offset int) (count int32, apis []utils.API) {
	accessToken, preCommandErr :=
		utils.ExecutePreCommandWithOAuth(cmdExportEnvironment, cmdUsername, cmdPassword,
			mainConfigFilePath, envKeysAllFilePath)
	if preCommandErr == nil {
		apiListEndpoint := utils.GetApiListEndpointOfEnv(cmdExportEnvironment, mainConfigFilePath)
		apiListEndpoint += "?limit=" + strconv.Itoa(utils.MaxAPIsToExportOnce) + "&offset=" + strconv.Itoa(offset)
		if (cmdResourceTenantDomain != "") {
			apiListEndpoint += "&tenantDomain=" + cmdResourceTenantDomain
		}
		count, apis, err := GetAPIList("", accessToken, apiListEndpoint)
		if err == nil {
			return count, apis
		} else {
			utils.HandleErrorAndExit(utils.LogPrefixError+"Getting List of APIs", err)
		}
	} else {
		utils.HandleErrorAndExit(utils.LogPrefixError + "Error in getting access token for user while getting " +
			"the list of APIs: ", preCommandErr)
	}
	return 0, nil
}

func init() {
	RootCmd.AddCommand(ExportAPIsCmd)
	ExportAPIsCmd.Flags().StringVarP(&cmdExportEnvironment, "environment", "e",
		utils.DefaultEnvironmentName, "Environment to which the API should be exported")
	ExportAPIsCmd.Flags().StringVarP(&cmdResourceTenantDomain, "tenant", "t", "",
		"Tenant domain of the resources to be exported")
	ExportAPIsCmd.Flags().StringVarP(&cmdUsername, "username", "u", "", "Username")
	ExportAPIsCmd.Flags().StringVarP(&cmdPassword, "password", "p", "", "Password")
	ExportAPIsCmd.PersistentFlags().BoolVarP(&cmdForceStartFromBegin, "force", "", false,
		"Allow connections to SSL endpoints without certs")
}