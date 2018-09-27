package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"path/filepath"
	"fmt"
	"net/http"
	"strconv"
	"os"
	)

const exportAPIsCmdLiteral = "export-apis"
const exportAPIsCmdShortDesc = "Export APIs"

var exportAPIsCmdLongDesc = "Export all the APIs of the tenant from an APIM 2.x environment environment, to be imported " +
	"into 3.0.0 environment"
var exportAPIsCmdExamples = ""
var apiExportDir string
//var offset int  //from which # of API, the APIs will be fetched from APIM server
var apiListOffset int
var count int32 // size of API list to be exported / number of  APIs left to be exported from last iteration
var apis [] utils.API
var exportRelatedFilesPath string
//e.g. /home/samithac/.wso2apimcli/exported/migration/production-2.5/wso2-dot-org
var startFromBeginning bool

var isProcessCompleted bool
//var iterationNo int
var startingApiIndexFromList int
var mainConfigFilePath string

var ExportAPIsCmd = &cobra.Command{
	Use: exportAPIsCmdLiteral + " (--environment " +
		"<environment-from-which-artifacts-should-be-exported>)",
	Short: exportAPIsCmdShortDesc,
	Long:  exportAPIsCmdLongDesc + exportAPIsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		//use log package for in-detail logging ; https://golangcode.com/add-line-numbers-to-log-output/
		utils.Logln(utils.LogPrefixInfo + exportAPIsCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)
		executeExportAPIsCmd(artifactExportDirectory)
	},
}

func executeExportAPIsCmd(exportDirectory string) {
	//create dir structure
	apiExportDir = createExportAPIsDirStructure(exportDirectory)
	exportRelatedFilesPath = filepath.Join(exportDirectory, cmdExportEnvironment,
		utils.GetMigrationExportTenantDirName(cmdResourceTenantDomain))
	//e.g. /home/samithac/.wso2apimcli/exported/migration/production-2.5/wso2-dot-org
	startFromBeginning = false
	isProcessCompleted = false

	fmt.Println("\nExporting APIs for the migration to APIM 3.0.0")
	if (cmdForceStartFromBegin) {
		startFromBeginning = true
	}

	if ((utils.IsFileExist(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName))) && !startFromBeginning) {
		var lastSuceededAPI utils.API
		lastSuceededAPI = utils.ReadLastSucceededAPIFileData(exportRelatedFilesPath);
		var migrationApisExportMetadata utils.MigrationApisExportMetadata
		migrationApisExportMetadata.ReadMigrationApisExportMetadataFile(filepath.Join(exportRelatedFilesPath,
			utils.MigrationAPIsExportMetadataFileName))
		apis = migrationApisExportMetadata.ApiListToExport
		apiListOffset = migrationApisExportMetadata.ApiListOffset
		startingApiIndexFromList = getLastSuceededApiIndex(lastSuceededAPI) + 1

		//find count of APIs left to be exported
		var lastSucceededAPInumber = getLastSuceededApiIndex(lastSuceededAPI) + 1
		count = int32(len(apis) - lastSucceededAPInumber)

		if (count == 0) {
			//last iteration had been completed successfully but operation had halted at that point.
			//So get the next set of APIs for next iteration
			apiListOffset += utils.MaxAPIsToExportOnce
			startingApiIndexFromList = 0
			count, apis = getAPIList()
			if (len(apis) > 0) {
				utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
					exportRelatedFilesPath, apiListOffset)
			} else {
				println("All the APIs has been exported and so the execution of export-apis command is completed successfully")
				os.Exit(1)
			}
		}
	} else { // start from beginning
		prepareStartFromBeginning()
	}

	if (count == 0) {
		fmt.Println("No APIs available to be exported..!")
	} else {
		for (count > 0) {
			fmt.Println("Found ", count, "of APIs to be exported in the iteration beginning with the offset #"+
				strconv.Itoa(apiListOffset) + " with a maximum limit of " + strconv.Itoa(utils.MaxAPIsToExportOnce))
			//get basic Auth credentials
			b64encodedCredentials, preCommandErr :=
				utils.ExecutePreCommandWithBasicAuth(cmdExportEnvironment, cmdUsername, cmdPassword,
					utils.MainConfigFilePath, utils.EnvKeysAllFilePath)
			if preCommandErr != nil {
				utils.Logln("Error generating base64 encoded credentials for Basic Authentication:")
				//return preCommandErr
			}
			for i := startingApiIndexFromList; i < len(apis); i++ {
				/*if( (i==2) && (apiListOffset == 8)) {
					os.Exit(1)
				}*/

				exportAPIName := apis[i].Name
				exportAPIVersion := apis[i].Version
				exportApiProvider := apis[i].Provider
				apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
				resp := getExportApiResponse(exportAPIName, exportAPIVersion, exportApiProvider, apiImportExportEndpoint,
					b64encodedCredentials)

				// Print info on response
				utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

				if resp.StatusCode() == http.StatusOK {
					WriteToZip(exportAPIName, exportAPIVersion, apiExportDir, resp)
					//write on last-succeeded-api.log
					utils.WriteLastSuceededAPIFileData(exportRelatedFilesPath, apis[i])
				} else {
					fmt.Println("Error exporting API:", exportAPIName, "-", exportAPIVersion, " of Provider:", exportApiProvider)
				}
			}

			//offset = utils.MaxAPIsToExportOnce * iterationNo
			//iterationNo++
			apiListOffset += utils.MaxAPIsToExportOnce
			count, apis = getAPIList()
			startingApiIndexFromList = 0
			if (len(apis) > 0) {
				utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
					exportRelatedFilesPath, apiListOffset)
			}
		}
		fmt.Println("All the APIs has been exported and so the execution of export-apis command is completed successfully")
	}
}
func getLastSuceededApiIndex(lastSuceededApi utils.API) (int) {
	for i := 0; i < len(apis); i++ {
		if ((apis[i].Name == lastSuceededApi.Name) &&
			(apis[i].Provider == lastSuceededApi.Provider) &&
			(apis[i].Version == lastSuceededApi.Version)) {
			return i
		}
	}
	return -1
}

func prepareStartFromBeginning() {
	//cleaning existing old files (if exists) related to exportation
	error := utils.RemoveDirectoryIfExists(filepath.Join(exportRelatedFilesPath, utils.ExportedApisDirName))
	error = utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.MigrationAPIsExportMetadataFileName))
	error = utils.RemoveFileIfExists(filepath.Join(exportRelatedFilesPath, utils.LastSucceededApiFileName))
	if (error != nil) {
		utils.HandleErrorAndExit("Error occurred while cleaning existing old files (if exists) related to "+
			"exportation", error)
	}

	//offset = 0;
	apiListOffset = 0
	//iterationNo = 1;
	startingApiIndexFromList = 0;
	count, apis = getAPIList()
	//write  migration-apis-export-metadata.yaml file
	utils.WriteMigrationApisExportMetadataFile(apis, cmdResourceTenantDomain, cmdUsername,
		exportRelatedFilesPath, apiListOffset)

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

	if dirExists, _ := utils.IsDirExists(migrationsArtifactsEnvTenantApisPath); dirExists {
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

func getAPIList() (count int32, apis []utils.API) {
	accessToken, preCommandErr :=
		utils.ExecutePreCommandWithOAuth(cmdExportEnvironment, cmdUsername, cmdPassword,
			utils.MainConfigFilePath, utils.EnvKeysAllFilePath)
	if preCommandErr == nil {
		apiListEndpoint := utils.GetApiListEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		apiListEndpoint += "?limit=" + strconv.Itoa(utils.MaxAPIsToExportOnce) + "&offset=" + strconv.Itoa(apiListOffset)
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
		utils.HandleErrorAndExit(utils.LogPrefixError+"Error in getting access token for user while getting "+
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
