package k8s

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// Get command related usage Info
const K8sGenCmdLiteral = "gen"
const k8sGenCmdShortDesc = "Generate deployment directory for K8S operator"

const k8sGenCmdLongDesc = `Generate sample directory with all the contents to use as the deployment directory` +
	`  when performing CI/CD pipeline tasks `

const k8sGenCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sGenCmdLiteral + ` ` + GenDeploymentDirCmdLiteral

// ListCmd represents the list command
var GenCmd = &cobra.Command{
	Use:     K8sGenCmdLiteral,
	Short:   k8sGenCmdShortDesc,
	Long:    k8sGenCmdLongDesc,
	Example: k8sGenCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + K8sGenCmdLiteral + " called")

	},
}
