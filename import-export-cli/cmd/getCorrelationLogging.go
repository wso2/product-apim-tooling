/*
*  Copyright (c) WSO2 LLC (http://www.wso2.com) All Rights Reserved.
*
*  WSO2 LLC licenses this file to you under the Apache License,
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
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getCorrelationLoggingEnvironment string
var getCorrelationLoggingCmdFormat string

const GetCorrelationLoggingCmdLiteral = "correlation-logging"
const getCorrelationLoggingCmdShortDesc = "Display a list of correlation logging components in an environment"
const getCorrelationLoggingCmdLongDesc = `Display a list of correlation logging components available in the environment specified`

var getCorrelationLoggingCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetCorrelationLoggingCmdLiteral + ` -e dev `

var getCorrelationLoggingCmd = &cobra.Command{
	Use:     GetCorrelationLoggingCmdLiteral,
	Short:   getCorrelationLoggingCmdShortDesc,
	Long:    getCorrelationLoggingCmdLongDesc,
	Example: getCorrelationLoggingCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GetCmdLiteral + " " + GetCorrelationLoggingCmdLiteral + " called")
		cred, err := GetCredentials(getCorrelationLoggingEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeGetCorrelationLoggingCmd(cred)
	},
}

func executeGetCorrelationLoggingCmd(credential credentials.Credential) {
	components, err := impl.GetCorrelationLogComponentListFromEnv(credential, getCorrelationLoggingEnvironment)
	
	if err == nil {
		impl.PrintCorrelationLoggers(components, getCorrelationLoggingCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError + "Getting list of correlation log configurations", err)
		utils.HandleErrorAndExit("Error while getting list of correlation log configurations", err)
	}
}

func init() {
	GetCmd.AddCommand(getCorrelationLoggingCmd)

	getCorrelationLoggingCmd.Flags().StringVarP(&getCorrelationLoggingEnvironment, "environment", "e",
		"", "Environment which the correlation logging components should be displayed")
	getCorrelationLoggingCmd.Flags().StringVarP(&getCorrelationLoggingCmdFormat, "format", "", "", "Pretty-print Correlation loggers "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = getCorrelationLoggingCmd.MarkFlagRequired("environment")
}
