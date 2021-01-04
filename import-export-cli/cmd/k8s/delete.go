package k8s

import (
	cmd2 "github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

// Delete command related usage Info
const k8sDeleteCmdLiteral = "delete"
const k8sDeleteCmdShortDesc = "Delete resources related to kubernetes"
const k8sDeleteCmdLongDesc = `Delete resources by filenames, stdin, resources and names, or by resources and label selector in kubernetes mode`

const k8sDeleteCmdExamples = utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` petstore
` + utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` -l name=myLabel`

// k8sDeleteCmd represents the delete command
var k8sDeleteCmd = &cobra.Command{
	Use:                k8sDeleteCmdLiteral,
	Short:              k8sDeleteCmdShortDesc,
	Long:               k8sDeleteCmdLongDesc,
	Example:            k8sDeleteCmdExamples,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + k8sDeleteCmdLiteral + " called")
		k8sArgs := []string{k8sUtils.K8sDelete}
		k8sArgs = append(k8sArgs, args...)
		cmd2.ExecuteKubernetes(k8sArgs...)
	},
}

func init() {
	K8sCmd.AddCommand(k8sDeleteCmd)
}

