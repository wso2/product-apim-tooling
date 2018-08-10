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

var exportAllCmdUsername string
var exportAllCmdPassword string

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
		utils.Logln(utils.LogPrefixInfo + exportAllCmdLiteral + " called")
		var artifactExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedMigrationArtifactsDirName)
		executeExportAllCmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, artifactExportDirectory)
	},
}

func executeExportAllCmd(mainConfigFilePath, envKeysAllFilePath, artifactExportDirectory string) {
	error:= executeExportAPIsCmd(mainConfigFilePath, envKeysAllFilePath, artifactExportDirectory)
	utils.Logln("this is calling Logln")

	//executeExportAPPsCmd(utils.MainConfigFilePath, utils.EnvKeysAllFilePath, artifactExportDirectory)
	if(error != nil) {
		utils.Logln(utils.LogPrefixInfo, error.Error())
	} else {
		utils.Logln(utils.LogPrefixInfo, "Succeeded exporting all the artifacts of the environment")
	}
}

func executeExportAPIsCmd(mainConfigFilePath, envKeysAllFilePath, artifactExportDirectory string) (error) {
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
		utils.ExecutePreCommandWithBasicAuth(exportEnvironment, exportAllCmdUsername, exportAllCmdPassword,
			mainConfigFilePath, envKeysAllFilePath)
	if preCommandErr != nil {
		utils.Logln("Error generating base64 encoded credentials for Basic Authentication:")
		return preCommandErr
	}

	utils.CreateDirIfNotExist(filepath.Join(artifactExportDirectory, exportEnvironment))
	utils.CreateDirIfNotExist(filepath.Join(artifactExportDirectory, exportEnvironment, "migration-1"))

	//iterate apis
	//Is it required to stop exporting process once an error in single API occurs?
	for _, api := range apis {
		exportAPIName := api.Name
		exportAPIVersion := api.Version
		exportProvider := api.Provider
		apiImportExportEndpoint := utils.GetApiImportExportEndpointOfEnv(exportEnvironment, mainConfigFilePath)
		resp := getExportApiResponse(exportAPIName, exportAPIVersion, exportProvider, apiImportExportEndpoint,
			b64encodedCredentials)// TODO Handle errors in getExportApiResponse()

		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())

		if resp.StatusCode() == http.StatusOK {
			//artifactExportDirectory >> use a suitable dynamic folder name instead of 'Migration-1' and also change
			// WriteToZip function
			//TODO: Change WriteToZip() function and Handle errors
			WriteToZip(exportAPIName, exportAPIVersion, filepath.Join(exportEnvironment, "migration-1"),
				artifactExportDirectory, resp)
		} else {
			fmt.Println("Error exporting API:", exportAPIName, "-", exportAPIVersion," of Provider:", exportProvider)
			return errors.New(resp.String())
		}
	}
	return nil
}
func getAPIList(mainConfigFilePath string, envKeysAllFilePath string) (count int32, apis []utils.API, err error) {
	accessToken, preCommandErr :=
		utils.ExecutePreCommandWithOAuth(listApisCmdEnvironment, exportAllCmdUsername, exportAllCmdPassword,
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
	ExportAllCmd.Flags().StringVarP(&exportEnvironment, "environment", "e",
		"", "Environment from which the artifacts should be exported") //creates default dir
	ExportAllCmd.Flags().StringVarP(&exportAllCmdUsername, "username", "u", "", "Username")
	ExportAllCmd.Flags().StringVarP(&exportAllCmdPassword, "password", "p", "", "Password")
}
