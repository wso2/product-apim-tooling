package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"strings"
)

var setApiLoggingCmdEnvironment string
var setApiLoggingCmdFormat string
var setApiLoggingCmdQuery []string
var setApiLoggingCmdLimit string

const SetApiLoggingCmdLiteral = "api-logging"
const setApiLoggingCmdShortDesc = "Display a list of API logger in an environment"
const setApiLoggingCmdLongDesc = `Display a list of API logger in the environment`

var setApiLoggingCmdExamples =
	utils.ProjectName + ` ` + setCmdLiteral + ` ` + SetApiLoggingCmdLiteral + `
` + utils.ProjectName + ` ` + setCmdLiteral + ` ` + SetApiLoggingCmdLiteral + ` --api-id`

var setApiLoggingCmd = &cobra.Command{
	Use:     SetApiLoggingCmdLiteral,
	Short:   setApiLoggingCmdShortDesc,
	Long:    setApiLoggingCmdLongDesc,
	Example: setApiLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + SetApiLoggingCmdLiteral + " called")
		cred, err := GetCredentials(setApiLoggingCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeSetApiLoggingCmd(cred)
	},
}

func executeSetApiLoggingCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, setApiLoggingCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+SetApiLoggingCmdLiteral+"'", err)
	}

	_, apis, err := impl.GetAPIListFromEnv(accessToken, setApiLoggingCmdEnvironment,
		strings.Join(setApiLoggingCmdQuery, queryParamSeparator), setApiLoggingCmdLimit)
	if err == nil {
		impl.PrintAPIs(apis, setApiLoggingCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
	}
}

func init() {
	SetCmd.AddCommand(setApiLoggingCmd)

	setApiLoggingCmd.Flags().StringVarP(&setApiLoggingCmdEnvironment, "api-id", "i",
		"", "Api ID")
}