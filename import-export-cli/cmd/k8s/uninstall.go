package k8s

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const K8sUninstallCmdLiteral = "uninstall"
const k8sUninstallCmdShortDesc = "Uninstall an operator in the configured K8s cluster"
const k8sUninstallCmdLongDesc = "Uninstall an operator in the configured K8s cluster"
const k8sUninstallCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sUninstallCmdLiteral + ` ` + K8sUninstallApiOperatorCmdLiteral

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:     K8sUninstallCmdLiteral,
	Short:   k8sUninstallCmdShortDesc,
	Long:    k8sUninstallCmdLongDesc,
	Example: k8sUninstallCmdExamples,
}

func init() {
	K8sCmd.AddCommand(uninstallCmd)
}

