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
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/cmd"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const defaulEnvsTableFormat = "table {{.Name}}\t{{.ApiManagerEndpoint}}\t{{.RegistrationEndpoint}}\t{{.TokenEndpoint}}\t{{.PublisherEndpoint}}\t{{.ApplicationEndpoint}}\t{{.AdminEndpoint}}"

var envsCmdFormat string

// envsCmd related info
const EnvsCmdLiteral = "envs"
const EnvsCmdShortDesc = "Display the list of environments"

const EnvsCmdLongDesc = `Display a list of environments defined in '` + utils.MainConfigFileName + `' file`

const EnvsCmdExamples = utils.ProjectName + " list envs"

// envsCmd represents the envs command
var envsCmdDeprecated = &cobra.Command{
	Use:        EnvsCmdLiteral,
	Short:      EnvsCmdShortDesc,
	Long:       EnvsCmdLongDesc,
	Example:    EnvsCmdExamples,
	Deprecated: "instead use \"" + cmd.GetCmdLiteral + " " + cmd.GetEnvsCmdLiteral + "\".",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + EnvsCmdLiteral + " called")
		envs := utils.GetMainConfigFromFile(utils.MainConfigFilePath).Environments
		impl.PrintEnvs(envs, envsCmdFormat, defaulEnvsTableFormat)
	},
}

func init() {
	ListCmdDeprecated.AddCommand(envsCmdDeprecated)
	envsCmdDeprecated.Flags().StringVarP(&envsCmdFormat, "format", "", defaulEnvsTableFormat, "Pretty-print "+
		"environments using go templates")
}
