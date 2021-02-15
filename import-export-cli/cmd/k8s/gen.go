package k8s

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// Get command related usage Info
const GenCmdLiteral = "gen"
const GenCmdShortDesc = "Generate deployment directory for VM and K8S operator"

const GenCmdLongDesc = `Generate sample directory with all the contents to use as the deployment directory` +
	`  when performing CI/CD pipeline tasks `

const GenCmdExamples = utils.ProjectName + ` ` + GenCmdLiteral + ` ` + GenDeploymentDirCmdLiteral

// ListCmd represents the list command
var GenCmd = &cobra.Command{
	Use:     GenCmdLiteral,
	Short:   GenCmdShortDesc,
	Long:    GenCmdLongDesc,
	Example: GenCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GenCmdLiteral + " called")

	},
}

