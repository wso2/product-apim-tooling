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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	impl "github.com/wso2/product-apim-tooling/import-export-cli/mi/impl"
	miUtils "github.com/wso2/product-apim-tooling/import-export-cli/mi/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var getTransactionCountCmdEnvironment string
var getTransactionCountCmdFormat string

const getTransactionCountCmdLiteral = "transaction-counts [year] [month]"

const getTransactionCountCmdShortDesc = "Retrieve transaction count"
const getTransactionCountCmdLongDesc = "Retrieve transaction count based on the given year and month.\n" +
	"If year and month not provided, retrieve the count for the current year and month of Micro Integrator in the environment specified by the flag --environment, -e"

var getTransactionCountCmdExamples = "To get the transaction count for the current month\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTransactionCountCmdLiteral) + " -e dev\n" +
	"To get the transaction count for a specific month\n" +
	"  " + utils.ProjectName + " " + utils.MiCmdLiteral + " " + GetCmdLiteral + " " + miUtils.GetTrimmedCmdLiteral(getTransactionCountCmdLiteral) + " 2020 06 -e dev\n" +
	"NOTE: The flag (--environment (-e)) is mandatory"

var getTransactionCountCmd = &cobra.Command{
	Use:     getTransactionCountCmdLiteral,
	Short:   getTransactionCountCmdShortDesc,
	Long:    getTransactionCountCmdLongDesc,
	Example: getTransactionCountCmdExamples,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) == 2 {
			return nil
		}
		var errMessage = "accepts exactly 0 or 2 arg(s), received " + fmt.Sprint(len(args))
		return errors.New(errMessage)
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleGetTransactionCountCmdArguments(args)
	},
}

func init() {
	GetCmd.AddCommand(getTransactionCountCmd)
	setEnvFlag(getTransactionCountCmd, &getTransactionCountCmdEnvironment)
	setFormatFlag(getTransactionCountCmd, &getTransactionCountCmdFormat)
}

func handleGetTransactionCountCmdArguments(args []string) {
	printGetCmdVerboseLogForArtifact(miUtils.GetTrimmedCmdLiteral(getTransactionCountCmdLiteral))
	credentials.HandleMissingCredentials(getTransactionCountCmdEnvironment)
	if len(args) == 2 {
		var year = args[0]
		var month = args[1]
		executeGetTransactionCountForMonth(year, month)
	} else {
		executeGetTransactionCountForMonth()
	}
}

func executeGetTransactionCountForMonth(period ...string) {
	transactionCount, err := impl.GetTransactionCount(getTransactionCountCmdEnvironment, period)
	if err == nil {
		impl.PrintTransactionCount(transactionCount, getTransactionCountCmdFormat)
	} else {
		fmt.Println(utils.LogPrefixError+"Retrieving transactions count.", err)
	}
}
