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

package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"os"
)

// envsCmd related info
const EnvsCmdLiteral = "envs"
const EnvsCmdShortDesc = "Display the list of environments"

var EnvsCmdLongDesc = dedent.Dedent(`
		Display a list of environments defined in '` + utils.MainConfigFileName + `' file
	`)

var EnvsCmdExamples = dedent.Dedent(`
		` + utils.ProjectName + ` list envs
	`)

// envsCmd represents the envs command
var envsCmd = &cobra.Command{
	Use:   EnvsCmdLiteral,
	Short: EnvsCmdShortDesc,
	Long:  EnvsCmdLongDesc + EnvsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + EnvsCmdLiteral + " called")

		envs := utils.GetMainConfigFromFile(utils.MainConfigFilePath).Environments

		printEnvs(envs)
	},
}

func printEnvs(envs map[string]utils.EnvEndpoints) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Publisher Endpoint", "Registration Endpoint", "Token Endpoint"})

	var data [][]string

	for env, endpoints := range envs {
		data = append(data, []string{env, endpoints.ApiManagerEndpoint, endpoints.RegistrationEndpoint,
			endpoints.TokenEndpoint})
	}

	for _, v := range data {
		table.Append(v)
	}

	fmt.Printf("Environments available in file '%s'\n", utils.MainConfigFilePath)
	table.Render()

}

func init() {
	ListCmd.AddCommand(envsCmd)
}
