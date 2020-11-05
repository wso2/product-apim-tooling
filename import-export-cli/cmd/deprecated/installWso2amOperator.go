/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package deprecated

import (
	"fmt"

	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"

	"github.com/spf13/cobra"
)

const installWso2amOperatorCmdLiteral = "wso2am-operator"
const installWso2amOperatorCmdShortDesc = "Install WSO2AM Operator"
const installWso2amOperatorCmdLongDesc = "Install WSO2AM Operator in the configured K8s cluster"
const installWso2amOperatorCmdExamples = utils.ProjectName + ` ` + installCmdLiteral + ` ` + installWso2amOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + installCmdLiteral + ` ` + installWso2amOperatorCmdLiteral + ` -f path/to/operator/configs
` + utils.ProjectName + ` ` + installCmdLiteral + ` ` + installWso2amOperatorCmdLiteral + ` -f path/to/operator/config/file.yaml`

// flags
var flagWso2AmOperatorFile string

// installWso2amOperatorCmdDeprecated represents the 'install wso2am-operator' command
var installWso2amOperatorCmdDeprecated = &cobra.Command{
	Use:     installWso2amOperatorCmdLiteral,
	Short:   installWso2amOperatorCmdShortDesc,
	Long:    installWso2amOperatorCmdLongDesc,
	Example: installWso2amOperatorCmdExamples,
	Deprecated: "use \"" + cmd.K8sCmdLiteral + " " + cmd.K8sInstallCmdLiteral + " " + cmd.K8sInstallWso2amOperatorCmdLiteral +
		"\" " + "instead of \"" + installCmdLiteral + " " + installWso2amOperatorCmdLiteral + "\".",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(fmt.Sprintf("%s%s %s called", utils.LogPrefixInfo, installCmdLiteral, installWso2amOperatorCmdLiteral))

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
	installCmdDeprecated.AddCommand(installWso2amOperatorCmdDeprecated)
	installWso2amOperatorCmdDeprecated.Flags().StringVarP(&flagWso2AmOperatorFile, "from-file", "f", "", "Path to wso2am-operator directory")
}
