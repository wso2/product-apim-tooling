package k8s

import (
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const K8sInstallCmdLiteral = "install"
const k8sInstallCmdShortDesc = "Install an operator in the configured K8s cluster"
const k8sInstallCmdLongDesc = "Install an operator in the configured K8s cluster"
const k8sInstallCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallApiOperatorCmdLiteral

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     K8sInstallCmdLiteral,
	Short:   k8sInstallCmdShortDesc,
	Long:    k8sInstallCmdLongDesc,
	Example: k8sInstallCmdExamples,
}

func init() {
	K8sCmd.AddCommand(installCmd)
}

