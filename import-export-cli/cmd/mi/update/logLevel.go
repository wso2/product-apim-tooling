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

package update

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var updateLogLevelCmdEnvironment string

const updateLogLevelCmdLiteral = "log-level [logger-name] [log-level]"
const updateLogLevelCmdShortDesc = "Update log level of a Logger in a Micro Integrator"

const updateLogLevelCmdLongDesc = "Update the log level of a Logger named [logger-name] to [log-level] specified by the command line arguments in a Micro Integrator in the environment specified by the flag --environment, -e"

var updateLogLevelCmdExamples = "To update the log level\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + updateCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(updateLogLevelCmdLiteral) + " org-apache-coyote DEBUG -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var updateLogLevelCmd = &cobra.Command{
	Use:     updateLogLevelCmdLiteral,
	Short:   updateLogLevelCmdShortDesc,
	Long:    updateLogLevelCmdLongDesc,
	Example: updateLogLevelCmdExamples,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		handleupdateLogLevelCmdArguments(args)
	},
}

func init() {
	UpdateCmd.AddCommand(updateLogLevelCmd)
	updateLogLevelCmd.Flags().StringVarP(&updateLogLevelCmdEnvironment, "environment", "e", "", "Environment of the micro integrator of which the logger should be updated")
	updateLogLevelCmd.MarkFlagRequired("environment")
}

func handleupdateLogLevelCmdArguments(args []string) {
	printUpdateCmdVerboseLog(miUtils.GetTrimmedCmdLiteral(updateLogLevelCmdLiteral))
	credentials.HandleMissingCredentials(updateLogLevelCmdEnvironment)
	executeUpdateLogger(args[0], args[1])
}

func executeUpdateLogger(loggerName, logLevel string) {
	resp, err := impl.UpdateMILogger(updateLogLevelCmdEnvironment, loggerName, logLevel)
	if err != nil {
		fmt.Println(utils.LogPrefixError+"updating logger [ "+loggerName+" ] ", err)
	} else {
		fmt.Println(resp)
	}
}
