package k8s

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/operator/registry"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const K8sChangeCmdLiteral = "change"
const k8sChangeCmdShortDesc = "Change a configuration in K8s cluster resource"
const k8sChangeCmdLongDesc = "Change a configuration in K8s cluster resource"
const k8sChangeCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sChangeCmdLiteral + ` ` + K8sChangeDockerRegistryCmdLiteral

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:     K8sChangeCmdLiteral,
	Short:   k8sChangeCmdShortDesc,
	Long:    k8sChangeCmdLongDesc,
	Example: k8sChangeCmdExamples,
}

const K8sChangeDockerRegistryCmdLiteral = "registry"
const k8sChangeDockerRegistryCmdShortDesc = "Change the registry"
const k8sChangeDockerRegistryCmdLongDesc = "Change the registry to be pushed the built micro-gateway image"
const k8sChangeDockerRegistryCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sChangeCmdLiteral + ` ` + K8sChangeDockerRegistryCmdLiteral

// changeDockerRegistryCmd represents the change registry command
var changeDockerRegistryCmd = &cobra.Command{
	Use:     K8sChangeDockerRegistryCmdLiteral,
	Short:   k8sChangeDockerRegistryCmdShortDesc,
	Long:    k8sChangeDockerRegistryCmdLongDesc,
	Example: k8sChangeDockerRegistryCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(fmt.Sprintf("%s%s %s called", utils.LogPrefixInfo, K8sChangeCmdLiteral, K8sChangeDockerRegistryCmdLiteral))
		configVars := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
		if !configVars.Config.KubernetesMode {
			utils.HandleErrorAndExit("set mode to kubernetes with command: apictl set --mode kubernetes",
				errors.New("mode should be set to kubernetes"))
		}

		// check for installation mode: interactive or batch mode
		if flagBmRegistryType == "" {
			// run in interactive mode
			// read inputs for docker registry
			registry.ChooseRegistryInteractive()
			registry.ReadInputsInteractive()
		} else {
			// run in batch mode
			// set registry type
			registry.SetRegistry(flagBmRegistryType)

			flagsValues := getGivenFlagsValues()
			registry.ValidateFlags(flagsValues)       // validate flags with respect to registry type
			registry.ReadInputsFromFlags(flagsValues) // read values from flags with respect to registry type
		}

		registry.UpdateConfigsSecrets()
	},
}

func init() {
	K8sCmd.AddCommand(changeCmd)
	changeCmd.AddCommand(changeDockerRegistryCmd)

	// flags for installing api-operator in batch mode
	// only the flag "registry-type" is required and others are registry specific flags
	// same flags defined in 'installApiOperator'
	changeDockerRegistryCmd.Flags().StringVarP(&flagBmRegistryType, "registry-type", "R", "", "Registry type: DOCKER_HUB | AMAZON_ECR |GCR | HTTP")
	changeDockerRegistryCmd.Flags().StringVarP(&flagBmRepository, k8sUtils.FlagBmRepository, "r", "", "Repository name or URI")
	changeDockerRegistryCmd.Flags().StringVarP(&flagBmUsername, k8sUtils.FlagBmUsername, "u", "", "Username of the repository")
	changeDockerRegistryCmd.Flags().StringVarP(&flagBmPassword, k8sUtils.FlagBmPassword, "p", "", "Password of the given user")
	changeDockerRegistryCmd.Flags().BoolVar(&flagBmPasswordStdin, k8sUtils.FlagBmPasswordStdin, false, "Prompt for password of the given user in the stdin")
	changeDockerRegistryCmd.Flags().StringVarP(&flagBmKeyFile, k8sUtils.FlagBmKeyFile, "c", "", "Credentials file")
}
