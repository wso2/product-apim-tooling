package cmd

import (
	"fmt"
	//"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	//"net/http"
	"path/filepath"
)

var exportThrottlePolicyName string
var exportThrottlePolicyType string

//var runningExportThrottlePolicyCommand bool

// ExportThrottlePolicy command related usage info
const ExportThrottlePolicyCmdLiteral = "throttlepolicy"
const exportThrottlePolicyCmdShortDesc = "Export Throttling Policy"

const exportThrottlePolicyCmdLongDesc = "Export a ThrottlingPolicy from an environment"

const exportThrottlePolicyCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportThrottlePolicyCmdLiteral + ` -n Policy2 -type custom -e dev
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportThrottlePolicyCmdLiteral + ` -n AppPolicy -type app -e production
` + utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + ExportThrottlePolicyCmdLiteral + ` -n SubPolicy -type sub -e dev
NOTE: All the 3 flags (--name (-n), --type and --environment (-e)) are mandatory.`

// ExportAPICmd represents the exportAPI command
var ExportThrottlePolicyCmd = &cobra.Command{
	Use: ExportThrottlePolicyCmdLiteral + " (--name <name-of-the-throttling-policy> --type <type-of-the-throttling-policy> --environment " +
		"<environment-from-which-the-throttling-policy-should-be-exported>)",
	Short:   exportThrottlePolicyCmdShortDesc,
	Long:    exportThrottlePolicyCmdLongDesc,
	Example: exportThrottlePolicyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportThrottlePolicyCmdLiteral + " called")
		var apisExportDirectory = filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)

		cred, err := GetCredentials(CmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}

		executeExportThrottlePolicyCmd(cred, apisExportDirectory)
	},
}

func executeExportThrottlePolicyCmd(credential credentials.Credential, exportDirectory string) {
	fmt.Println("Code Still in Progress")
	//runningExportThrottlePolicyCommand = true
	//accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, CmdExportEnvironment)
	//
	//if preCommandErr == nil {
	//	resp, err := impl.ExportThrottlingFromEnv(accessToken, exportAPIName, exportAPIVersion, exportRevisionNum, exportProvider,
	//		exportAPIFormat, CmdExportEnvironment, exportAPIPreserveStatus, exportAPILatestRevision)
	//	if err != nil {
	//		utils.HandleErrorAndExit("Error while exporting", err)
	//	}
	//	// Print info on response
	//	utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
	//	apiZipLocationPath := filepath.Join(exportDirectory, CmdExportEnvironment)
	//	if resp.StatusCode() == http.StatusOK {
	//		impl.WriteToZip(exportAPIName, exportAPIVersion, "", apiZipLocationPath, runningExportApiCommand, resp)
	//	} else if resp.StatusCode() == http.StatusInternalServerError {
	//		// 500 Internal Server Error
	//		fmt.Println(string(resp.Body()))
	//	} else {
	//		// neither 200 nor 500
	//		fmt.Println("Error exporting API:", resp.Status(), "\n", string(resp.Body()))
	//	}
	//} else {
	//	// error exporting Api
	//	fmt.Println("Error getting OAuth tokens while exporting API:" + preCommandErr.Error())
	//}
}

// init using Cobra
func init() {
	ExportCmd.AddCommand(ExportThrottlePolicyCmd)
	ExportThrottlePolicyCmd.Flags().StringVarP(&exportThrottlePolicyName, "name", "n", "",
		"Name of the ThrottlingPolicy to be exported")
	ExportThrottlePolicyCmd.Flags().StringVarP(&exportThrottlePolicyType, "type", "t",
		"", "Type of the ThrottlingPolicy to be exported")
	ExportThrottlePolicyCmd.Flags().StringVarP(&CmdExportEnvironment, "environment", "e",
		"", "Environment to which the Throttling Policy should be exported")
	_ = ExportThrottlePolicyCmd.MarkFlagRequired("name")
	_ = ExportThrottlePolicyCmd.MarkFlagRequired("environment")
	_ = ExportThrottlePolicyCmd.MarkFlagRequired("type")

}
