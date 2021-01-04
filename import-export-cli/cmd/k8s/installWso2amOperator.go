package k8s

import (
	"fmt"

	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

const K8sInstallWso2amOperatorCmdLiteral = "wso2am-operator"
const k8sInstallWso2amOperatorCmdShortDesc = "Install WSO2AM Operator"
const k8sInstallWso2amOperatorCmdLongDesc = "Install WSO2AM Operator in the configured K8s cluster"
const k8sInstallWso2amOperatorCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallWso2amOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallWso2amOperatorCmdLiteral + ` -f path/to/operator/configs
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallWso2amOperatorCmdLiteral + ` -f path/to/operator/config/file.yaml`

// flags
var flagWso2AmOperatorFile string

// installWso2amOperatorCmd represents the 'install wso2am-operator' command
var installWso2amOperatorCmd = &cobra.Command{
	Use:     K8sInstallWso2amOperatorCmdLiteral,
	Short:   k8sInstallWso2amOperatorCmdShortDesc,
	Long:    k8sInstallWso2amOperatorCmdLongDesc,
	Example: k8sInstallWso2amOperatorCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(fmt.Sprintf("%s%s %s called", utils.LogPrefixInfo, K8sInstallCmdLiteral, K8sInstallWso2amOperatorCmdLiteral))

		// is -f or --from-file flag specified
		isLocalInstallation := flagWso2AmOperatorFile != ""
		configFile := flagWso2AmOperatorFile

		if !isLocalInstallation {
			// getting API Operator version
			operatorVersion, err := k8sUtils.GetVersion(
				"WSO2AM Operator",
				k8sUtils.Wso2AmOperatorVersionEnvVariable,
				k8sUtils.DefaultWso2AmOperatorVersion,
				k8sUtils.Wso2AmOperatorVersionValidationUrlTemplate,
				k8sUtils.Wso2AmOperatorFindVersionUrl,
			)
			if err != nil {
				utils.HandleErrorAndExit("Error in WSO2AM Operator version", err)
			}
			configFile = fmt.Sprintf(k8sUtils.Wso2AmOperatorConfigsUrlTemplate, operatorVersion)
		}

		// installing operator and configs if -f flag given
		// otherwise settings configs only
		k8sUtils.CreateControllerConfigs(configFile, 20, k8sUtils.Wso2amOpCrdApimanager)

		fmt.Println("[Setting to K8s Mode]")
		utils.SetToK8sMode()
	},
}

func init() {
	installCmd.AddCommand(installWso2amOperatorCmd)
	installWso2amOperatorCmd.Flags().StringVarP(&flagWso2AmOperatorFile, "from-file", "f", "", "Path to wso2am-operator directory")
}

