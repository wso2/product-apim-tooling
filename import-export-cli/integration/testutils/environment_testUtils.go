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

package testutils

import (
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"strconv"
	"testing"
)

func initProject(t *testing.T, args *initTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	output, err := base.Execute(t, "init", args.initFlag)
	return output, err
}

func initProjectWithDefinitionFlag(t *testing.T, args *initTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	output, err := base.Execute(t, "init", args.initFlag, "--definition", args.definitionFlag, "--force", strconv.FormatBool(args.forceFlag))
	return output, err
}

func initProjectWithOasFlag(t *testing.T, args *initTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.srcAPIM.GetEnvName(), args.srcAPIM.GetApimURL())
	base.Login(t, args.srcAPIM.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	output, err := base.Execute(t, "init", args.initFlag, "--oas", args.oasFlag)
	return output, err
}

func environmentSetExportDirectory(t *testing.T, args *setTestArgs) (string, error) {
	apim := args.srcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--export-directory", args.exportDirectoryFlag, "-k")
	return output, error
}

func environmentSetHttpRequestTimeout(t *testing.T, args *setTestArgs) (string, error) {
	apim := args.srcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--http-request-timeout", strconv.Itoa(args.httpRequestTimeout), "-k")
	return output, error
}

func environmentSetTokenType(t *testing.T, args *setTestArgs) (string, error) {
	apim := args.srcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--token-type", args.tokenTypeFlag, "-k")
	return output, error
}
