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

	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	defaultTransactionCountTableFormat = "table {{.Year}}\t{{.Month}}\t{{.TransactionCount}}"
)

// GetTransactionCount returns inbound transactions received by the micro integrator in a given environment
func GetTransactionCount(env string, period []string) (*artifactutils.TransactionCount, error) {

	var params map[string]string

	if len(period) == 2 {
		params = make(map[string]string)
		params["year"] = period[0]
		params["month"] = period[1]
	}

	var transactionCountResource = utils.MiManagementTransactionResource + "/" + utils.MiManagementTransactionCountResource

	resp, err := callMIManagementEndpointOfResource(transactionCountResource, params, env, &artifactutils.TransactionCount{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.TransactionCount), nil
}

// PrintTransactionCount prints the transaction count according to the given format
func PrintTransactionCount(transactionCount *artifactutils.TransactionCount, format string) {

	transactionContext := getContextWithFormat(format, defaultTransactionCountTableFormat)

	renderer := getItemRendererEndsWithNewLine(transactionCount)

	transactionCountTableHeaders := map[string]string{
		"Year":             yearHeader,
		"Month":            monthHeader,
		"TransactionCount": transactionCountHeader,
	}

	if err := transactionContext.Write(renderer, transactionCountTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
