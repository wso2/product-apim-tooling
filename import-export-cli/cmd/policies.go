package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// Export command related usage Info
const policyCmdLiteral = "policy"
const policyCmdShortDesc = "Export/Import a Policy"

const policyCmdLongDesc = "Export a Policy in an environment or Import a Policy to an environment"

const policyCmdExamples = utils.ProjectName + ` ` + ExportCmdLiteral + ` ` + policyCmdLiteral + ` ` + ExportThrottlePolicyCmdLiteral + ` -n Silver -e prod --type subscription`

// ExportCmd represents the export command
var policyCmd = &cobra.Command{
	Use:     policyCmdLiteral,
	Short:   policyCmdShortDesc,
	Long:    policyCmdLongDesc,
	Example: policyCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + ExportCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	ExportCmd.AddCommand(policyCmd)
}
