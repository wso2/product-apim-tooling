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

package impl

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const transactionReportFilePrefix = "transaction-count-summary-"

// GetTransactionReport returns inbound transactions received by the micro integrator in a given environment as a report
func GetTransactionReport(env string, period []string) (*artifactutils.TransactionCountInfo, error) {
	params := make(map[string]string)
	params["start"] = period[0]
	params["end"] = period[1]

	var transactionReportResource = utils.MiManagementTransactionResource + "/" + utils.MiManagementTransactionReportResource
	resp, err := callMIManagementEndpointOfResource(transactionReportResource, params, env, &artifactutils.TransactionCountInfo{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.TransactionCountInfo), nil
}

// WriteTransactionReportAsCSV writes the transaction report to a csv file in the specified target directory
func WriteTransactionReportAsCSV(transactions *artifactutils.TransactionCountInfo, targetDirectory string) {
	fileName := transactionReportFilePrefix + strconv.FormatInt(time.Now().UnixNano(), 10) + ".csv"
	destinationFilePath := filepath.Join(targetDirectory, fileName)
	transactionCountLines := transactions.TransactionCounts
	err := utils.WriteLinesToCSVFile(transactionCountLines, destinationFilePath)
	if err != nil {
		fmt.Println("Error writing the transaction report", err.Error())
	} else {
		fmt.Println("Transaction Count Report created in", destinationFilePath)
	}
}
