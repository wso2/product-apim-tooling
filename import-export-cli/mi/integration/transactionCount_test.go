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

package integration

import (
	"testing"

	"github.com/wso2/product-apim-tooling/import-export-cli/mi/integration/testutils"
)

const transactionCountCmd = "transaction-counts"

func TestGetTransactions(t *testing.T) {
	testutils.ValidateTransaction(t, transactionCountCmd, config)
}

func TestGetTransactionsWithoutSettingUpEnv(t *testing.T) {
	testutils.ExecGetCommandWithoutSettingEnv(t, transactionCountCmd)
}

func TestGetTransactionsWithoutLogin(t *testing.T) {
	testutils.ExecGetCommandWithoutLogin(t, transactionCountCmd, config)
}

func TestGetTransactionsWithoutEnvFlag(t *testing.T) {
	testutils.ExecGetCommandWithoutEnvFlag(t, transactionCountCmd, config)
}

func TestGetTransactionsWithInvalidArgs(t *testing.T) {
	testutils.ExecGetTransactionCountWithInvalidArgCount(t, config, 1, transactionCountCmd, "2020")
}
