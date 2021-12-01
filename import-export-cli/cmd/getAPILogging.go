package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"strings"
	"fmt"
)

var getApiLoggingCmdEnvironment string
var getApiLoggingCmdFormat string
var getApiLoggingCmdQuery []string
var getApiLoggingCmdLimit string

const GetApiLoggingCmdLiteral = "api-logging"
const getApiLoggingCmdShortDesc = "Display a list of API logger in an environment"
const getApiLoggingCmdLongDesc = `Display a list of API logger in the environment`

var getApiLoggingCmdExamples =
	utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiLoggingCmdLiteral + `
` + utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetApiLoggingCmdLiteral + ` --api-id`

var getApiLoggingCmd = &cobra.Command{
	Use:     GetApiLoggingCmdLiteral,
	Short:   getApiLoggingCmdShortDesc,
	Long:    getApiLoggingCmdLongDesc,
	Example: getApiLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("hello world")
		utils.Logln(utils.LogPrefixInfo + GetApiLoggingCmdLiteral + " called")
		cred, err := GetCredentials(getApiLoggingCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetApiLoggingCmd(cred)
	},
}

func executeGetApiLoggingCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, getApiLoggingCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+GetApiLoggingCmdLiteral+"'", err)
	}
	fmt.Print("env is", getApiLoggingCmdEnvironment)
	_, apis, err := impl.GetAPIListFromEnv(accessToken, getApiLoggingCmdEnvironment,
		strings.Join(getApiLoggingCmdQuery, queryParamSeparator), getApiLoggingCmdLimit)
	if err == nil {
		impl.PrintAPIs(apis, getApiLoggingCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
	}
}

func init() {
	GetCmd.AddCommand(getApiLoggingCmd)

	getApiLoggingCmd.Flags().StringVarP(&getApiLoggingCmdEnvironment, "api-id", "i",
		"", "Api ID")
	getApiLoggingCmd.Flags().StringVarP(&apiStateChangeEnvironment, "environment", "e",
		"", "Environment of which the API state should be changed")
}