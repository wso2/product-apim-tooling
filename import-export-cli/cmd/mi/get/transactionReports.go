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

package get

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getTransactionReportCmdEnvironment string
var transactionReportPath string

const getTransactionReportCmdLiteral = "transaction-reports [start] [end]"

const getTransactionReportCmdShortDesc = "Generate transaction count summary report"
const getTransactionReportCmdLongDesc = "Generate the transaction count summary report at the given location for the " +
	"given period of time.\nIf a location not provided, generate the report in current directory.\nIf an end date " +
	"not provided, generate the report with values upto current date of the Micro Integrator in the environment specified by the flag --environment, -e"

var getTransactionReportCmdExamples = "Example:\n" +
	"To generate transaction count report consisting data within a specified time period at a specified location\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTransactionReportCmdLiteral) + " 2020-05 2020-06 --path </dir_path> -e dev\n" +
	"To generate transaction count report with data from a given month upto the current month at a specified location\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTransactionReportCmdLiteral) + " 2020-01 -p </dir_path> -e dev\n" +
	"To generate transaction count report at the current location with data between 2020-01 and 2020-05\n" +
	"  " + utils.GetMICmdName() + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTransactionReportCmdLiteral) + " 2020-01 2020-05 -e dev\n" +
	"NOTE: The [start] argument and the flag (--environment (-e)) is mandatory"

var getTransactionReportCmd = &cobra.Command{
	Use:     getTransactionReportCmdLiteral,
	Short:   getTransactionReportCmdShortDesc,
	Long:    getTransactionReportCmdLongDesc,
	Example: getTransactionReportCmdExamples,
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		handleGetTransactionReportCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getTransactionReportCmd)
	setEnvFlag(getTransactionReportCmd, &getTransactionReportCmdEnvironment)
	getTransactionReportCmd.Flags().StringVarP(&transactionReportPath, "path", "p", "", "destination file location")
}

func handleGetTransactionReportCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getTransactionReportCmdLiteral))
	credentials.HandleMissingCredentials(getTransactionReportCmdEnvironment)
	var start = args[0]
	var end = ""
	if len(args) == 2 {
		end = args[1]
	}
	if isEmptyOrCurrentDir(transactionReportPath) {
		transactionReportPath, _ = os.Getwd()
	}
	executeGetTransactionReport(transactionReportPath, start, end)
}

func executeGetTransactionReport(targetDirectory string, period ...string) {
	transactionReport, err := impl.GetTransactionReport(getTransactionReportCmdEnvironment, period)
	if err == nil {
		impl.WriteTransactionReportAsCSV(transactionReport, targetDirectory)
	} else {
		fmt.Println(utils.LogPrefixError+"Retrieving Transaction Reports.", err)
	}
}
