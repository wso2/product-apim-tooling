/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

	"github.com/spf13/cobra"
	"github.com/renstrom/dedent"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

// envsCmd related info
const EnvsCmdLiteral string = "envs"
const EnvsCmdShortDesc string = "Display the list of environments"
var EnvsCmdLongDesc = dedent.Dedent(`
		Display a list of environments defined in '`+utils.MainConfigFileName+`' file
	`)

var EnvsCmdExamples = dedent.Dedent(`
		`+utils.ProjectName+` list envs
	`)

// envsCmd represents the envs command
var envsCmd = &cobra.Command{
	Use:   EnvsCmdLiteral,
	Short: EnvsCmdShortDesc,
	Long: EnvsCmdLongDesc + EnvsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + EnvsCmdLiteral + " called")
	},
}

func init() {
	ListCmd.AddCommand(envsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
