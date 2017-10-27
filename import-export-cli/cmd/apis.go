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
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"github.com/renstrom/dedent"
)

// apisCmd related info
const apisCmdLiteral string = "apis"
const apisCmdShortDesc string = "Display a list of APIs in an environment"

var apisCmdLongDesc = dedent.Dedent(`
		Display a list of APIs in the environment specified by the flag --environment, -e
	`)
var apisCmdExamples = dedent.Dedent(`
	`+utils.ProjectName + ` `+ apisCmdLiteral +` `+ listCmdLiteral +` -e dev
	`+utils.ProjectName + ` `+ apisCmdLiteral +` `+ listCmdLiteral +` -e staging
	`)

// apisCmd represents the apis command
var apisCmd = &cobra.Command{
	Use:   apisCmdLiteral,
	Short: apisCmdShortDesc,
	Long:  apisCmdLongDesc + apisCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apisCmdLiteral + " called")
	},
}

func init() {
	ListCmd.AddCommand(apisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
