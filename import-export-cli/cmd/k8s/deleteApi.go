package k8s

import (
	"github.com/spf13/cobra"
	cmd2 "github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// k8sDeleteAPI command related usage info
const k8sDeleteAPICmdLiteral = "api"
const k8sDeleteAPICmdShortDesc = "Delete API resources"
const k8sDeleteAPICmdLongDesc = "Delete API resources by API name or label selector in kubernetes mode"

const k8sDeleteAPICmdExamples = utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` petstore
` + "  " + utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + ` -l name=myLabel`

// k8sDeleteAPICmd represents the delete api command in kubernetes mode
var k8sDeleteAPICmd = &cobra.Command{
	Use:                utils.ProjectName + ` ` + k8sDeleteCmdLiteral + ` ` + k8sDeleteAPICmdLiteral + " (<name-of-the-api> or -l name=<name-of-the-label>)",
	Short:              k8sDeleteAPICmdShortDesc,
	Long:               k8sDeleteAPICmdLongDesc,
	Example:            k8sDeleteAPICmdExamples,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + k8sDeleteAPICmdLiteral + " called")
		k8sArgs := []string{k8sUtils.K8sDelete, k8sUtils.ApiOpCrdApi}
		k8sArgs = append(k8sArgs, args...)
		cmd2.ExecuteKubernetes(k8sArgs...)
	},
}

// Init using Cobra
func init() {
	k8sDeleteCmd.AddCommand(k8sDeleteAPICmd)
}

