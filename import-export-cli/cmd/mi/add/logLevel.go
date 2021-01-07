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

package add

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var addLogLevelCmdEnvironment string

const addLogLevelCmdLiteral = "log-level [logger-name] [class-name] [log-level]"
const addLogLevelCmdShortDesc = "Add new Logger to a Micro Integrator"

const addLogLevelCmdLongDesc = "Add new Logger named [logger-name] to the [class-name] with log level [log-level] specified by the command line arguments to a Micro Integrator in the environment specified by the flag --environment, -e"

var addLogLevelCmdExamples = "To add a new logger\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + addCmdLiteral + " log-level synapse-api org.apache.synapse.rest.API DEBUG -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var addLogLevelCmd = &cobra.Command{
	Use:     addLogLevelCmdLiteral,
	Short:   addLogLevelCmdShortDesc,
	Long:    addLogLevelCmdLongDesc,
	Example: addLogLevelCmdExamples,
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		handleAddLogLevelCmdArguments(args)
	},
}

func init() {
	AddCmd.AddCommand(addLogLevelCmd)
	addLogLevelCmd.Flags().StringVarP(&addLogLevelCmdEnvironment, "environment", "e", "", "Environment to be searched")
	addLogLevelCmd.MarkFlagRequired("environment")
}

func handleAddLogLevelCmdArguments(args []string) {
	printAddCmdVerboseLog("log-level")
	credentials.HandleMissingCredentials(addLogLevelCmdEnvironment)
	executeAddNewLogger(args[0], args[1], args[2])
}

func executeAddNewLogger(loggerName, logClass, logLevel string) {
	resp, err := impl.AddMILogger(addLogLevelCmdEnvironment, loggerName, logClass, logLevel)
	if err != nil {
		fmt.Println(utils.LogPrefixError+"Adding new logger [ "+loggerName+" ] ", err)
	} else {
		fmt.Println(resp)
	}
}
