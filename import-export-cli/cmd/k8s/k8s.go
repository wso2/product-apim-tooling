package k8s

import (
	"github.com/spf13/cobra"
	cmd2 "github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// K8s command related usage Info
const K8sCmdLiteral = "k8s"
const k8sCmdShortDesc = "Kubernetes mode based commands"

const k8sCmdLongDesc = `Kubernetes mode based commands such as install, uninstall, add/update api, change registry.`

const k8sCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallApiOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sUninstallCmdLiteral + ` ` + K8sUninstallApiOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sAddCmdLiteral + ` ` + AddApiCmdLiteral + ` ` + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sUpdateCmdLiteral + ` ` + AddApiCmdLiteral + ` ` + `-n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sChangeCmdLiteral + ` ` + K8sChangeDockerRegistryCmdLiteral

// K8sCmd represents the import command
var K8sCmd = &cobra.Command{
	Use:     K8sCmdLiteral,
	Short:   k8sCmdShortDesc,
	Long:    k8sCmdLongDesc,
	Example: k8sCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + cmd2.ImportCmdLiteral + " called")

	},
}

// init using Cobra
func init() {
	cmd2.RootCmd.AddCommand(K8sCmd)
}

