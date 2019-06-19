package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const logoutCmdLiteral = "logout [environment]"
const logoutCmdShortDesc = "logout to from an API Manager"
const logoutCmdLongDesc = `logout from an API Manager environment`
const logoutCmdExamples = utils.ProjectName + " logout dev"

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:     logoutCmdLiteral,
	Short:   logoutCmdShortDesc,
	Long:    logoutCmdLongDesc,
	Example: logoutCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runLogout(args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func runLogout(environment string) error {
	store, err := credentials.GetDefaultCredentialStore()
	if err != nil {
		return err
	}
	return store.Erase(environment)
}

// init using Cobra
func init() {
	RootCmd.AddCommand(logoutCmd)
}
