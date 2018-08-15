package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"path/filepath"
	"fmt"
	"github.com/renstrom/dedent"
	"net/http"
	"errors"
)

const exportAllCmdShortDesc = "Export All artifacts"
const exportAllCmdLiteral = "export-all"

var exportAllCmdLongDesc = "Export all artifact from an environment to be migrated into newer version"

var exportAllCmdExamples = dedent.Dedent(`
		Examples:
		` + utils.ProjectName + ` ` + exportAllCmdLiteral + ` -e dev
		NOTE: flags --environment (-e) is mandatory
	`)

var ExportAllCmd = &cobra.Command{
	Use: exportAllCmdLiteral + " (--environment " +
		"<environment-from-which-artifacts-should-be-exported>)",
	Short: exportAllCmdShortDesc,
	Long:  exportAllCmdLongDesc + exportAllCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		//use log package for in-detail logging ; https://golangcode.com/add-line-numbers-to-log-output/
		utils.Logln(utils.LogPrefixInfo + exportAllCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)
		executeExportAllCmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, artifactExportDirectory)
	},
}

func executeExportAllCmd(mainConfigFilePath, envKeysAllFilePath, migrationsArtifactsPath string) {
	//migrationsArtifactsPath = /home/.../.wso2apimcli/exported/all   >> already created
	/* need to create
		<wso2apimcli_HOME>/exported/all/production-2.5/migration-1/apps
	 	<wso2apimcli_HOME>/exported/all/production-2.5/migration-1/apis
	*/
	var createDirError error
	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsPath)

	migrationsArtifactsEnvPath := filepath.Join(migrationsArtifactsPath, cmdExportEnvironment)
	migrationsArtifactsEnvDirPath := filepath.Join(migrationsArtifactsEnvPath, "migration-1")
	migrationsArtifactsEnvDirApisPath := filepath.Join(migrationsArtifactsEnvDirPath, utils.ExportedApisDirName)
	migrationsArtifactsEnvDirAppsPath := filepath.Join(migrationsArtifactsEnvDirPath, utils.ExportedAppsDirName)

	createDirError = utils.CreateDirIfNotExist(migrationsArtifactsEnvPath)
	createDirError = utils.CreateDir(migrationsArtifactsEnvDirPath)
	createDirError = utils.CreateDir(migrationsArtifactsEnvDirApisPath)
	createDirError = utils.CreateDir(migrationsArtifactsEnvDirAppsPath)

	if (createDirError != nil) {
		fmt.Print(utils.LogPrefixError, "Error in creating directory structure for the artifact export.", createDirError)
	}
	exportApiError:= executeExportAPIsCmd(mainConfigFilePath, envKeysAllFilePath, migrationsArtifactsEnvDirApisPath)

	if(exportApiError != nil) {
		fmt.Println(utils.LogPrefixError, "Error in exporting APIs during export-all operation", exportApiError.Error())
	} else {
		utils.Logln(utils.LogPrefixInfo, "Succeeded exporting all the APIs of the environment for the export-all operation")
	}

	exportAppsError := executeExportAPPsCmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, migrationsArtifactsEnvDirAppsPath)
	if(exportAppsError != nil) {
		fmt.Println(utils.LogPrefixError, "Error in exporting APIs during export-all operation", exportAppsError.Error())
	} else {
		utils.Logln(utils.LogPrefixInfo, "Succeeded exporting all the Applications of the environment for the export-all operation")
	}
}

func executeExportAPPsCmd(mainConfigFilePath, envKeysAllFilePath, directoryMigrationArtifactsApps string) error{
	fmt.Println("\nExporting Applications...")
	//get list of Apps
	accessToken, preCommandErr :=
		utils.ExecutePreCommandWithOAuth(cmdExportEnvironment, cmdUsername, cmdPassword,
			mainConfigFilePath, envKeysAllFilePath)

	if preCommandErr != nil {
		utils.Logln(utils.LogPrefixError + "Error in generating token for Getting list of Applications : " + preCommandErr.Error())
		return preCommandErr
	}

	applicationListEndpoint := utils.GetApplicationListEndpointOfEnv(cmdExportEnvironment, mainConfigFilePath)
	count, apps, err := GetApplicationList("", accessToken, applicationListEndpoint)

	if err == nil {
		// Printing the list of available Applications
		utils.Logln("Environment:", cmdExportEnvironment)
		utils.Logln("No. of Applications:", count)
	} else {
		fmt.Println(utils.LogPrefixError+"Error in getting List of Applications", err)
		return err
	}

	//iterate and export apps
	for _, app := range apps {
		exportAppName = app.Name
		exportAppOwner = "admin" // app.Subscriber
		executeExportAppCmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, directoryMigrationArtifactsApps)
	}
	return nil
}

func executeExportAPIsCmd(mainConfigFilePath, envKeysAllFilePath, artifactExportDirectory string) (error) {
	fmt.Println("\nExporting APIs...")
	//get list of APIs
	count, apis, error := getAPIList(mainConfigFilePath, envKeysAllFilePath)

	if(error != nil) {
		utils.Logln("Error in calling getAPIList with mainConfigFilePath:" + mainConfigFilePath + ", envKeysAllFilePath:" +
			envKeysAllFilePath)
		return error
	}
	fmt.Println("Found ", count, "of APIs to be exported")
	//get basic Auth credentials
	b64encodedCredentials, preCommandErr :=
		utils.ExecutePreCommandWithBasicAuth(cmdExportEnvironment, cmdUsername, cmdPassword,
			mainConfigFilePath, envKeysAllFilePath)
	if preCommandErr != nil {
		utils.Logln("Error generating base64 encoded credentials for Basic Authentication:")
		return preCommandErr
	}

	//iterate apis
	//Is it required to stop exporting process once an error in single API occurs?
	for _, api := range apis {
		exportAPIName := api.Name
		exportAPIVersion := api.Version
		exportApiProvider := api.Provider
		apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(cmdExportEnvironment, mainConfigFilePath)
		resp := getExportApiResponse(exportAPIName, exportAPIVersion, exportApiProvider, apiImportExportEndpoint,
			b64encodedCredentials)// TODO Handle errors in getExportApiResponse()

		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

		if resp.StatusCode() == http.StatusOK {
			//artifactExportDirectory >> use a suitable dynamic folder name instead of 'Migration-1' and also change
			// WriteToZip function
			//TODO: Change WriteToZip() function and Handle errors
			WriteToZip(exportAPIName, exportAPIVersion, artifactExportDirectory, resp)
		} else {
			fmt.Println("Error exporting API:", exportAPIName, "-", exportAPIVersion," of Provider:", exportApiProvider)
			return errors.New(resp.String())
		}
	}
	return nil
}
func getAPIList(mainConfigFilePath string, envKeysAllFilePath string) (count int32, apis []utils.API, err error) {
	accessToken, preCommandErr :=
		utils.ExecutePreCommandWithOAuth(listApisCmdEnvironment, cmdUsername, cmdPassword,
			mainConfigFilePath, envKeysAllFilePath)
	if preCommandErr == nil {
		apiListEndpoint := utils.GetApiListEndpointOfEnv(listApisCmdEnvironment, mainConfigFilePath)
		count, apis, err = GetAPIList("", accessToken, apiListEndpoint)
		if err == nil {
			return count, apis, err
		} else {
			utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
			return 0,nil, err
		}
	} else {
		utils.Logln(utils.LogPrefixError + "Error in getting access token for user : " + preCommandErr.Error())
		//utils.HandleErrorAndExit("Error calling '"+apisCmdLiteral+"'", preCommandErr)
		return 0, nil, preCommandErr
	}
	return 0, nil, err
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ExportAllCmd)
	ExportAllCmd.Flags().StringVarP(&cmdExportEnvironment, "environment", "e",
		"", "Environment from which the artifacts should be exported") //creates default dir
	ExportAllCmd.Flags().StringVarP(&cmdUsername, "username", "u", "", "Username")
	ExportAllCmd.Flags().StringVarP(&cmdPassword, "password", "p", "", "Password")
}
