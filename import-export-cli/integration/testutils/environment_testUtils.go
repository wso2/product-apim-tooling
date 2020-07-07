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

func InitProject(t *testing.T, args *InitTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, err := base.Execute(t, "init", args.InitFlag)
	return output, err
}

func initProjectWithDefinitionFlag(t *testing.T, args *InitTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, err := base.Execute(t, "init", args.InitFlag, "--definition", args.definitionFlag, "--force", strconv.FormatBool(args.ForceFlag))
	return output, err
}

func InitProjectWithOasFlag(t *testing.T, args *InitTestArgs) (string, error) {
	//Setup Environment and login to it.
	base.SetupEnvWithoutTokenFlag(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL())
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, err := base.Execute(t, "init", args.InitFlag, "--oas", args.OasFlag)
	return output, err
}

func EnvironmentSetExportDirectory(t *testing.T, args *SetTestArgs) (string, error) {
	apim := args.SrcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--export-directory", args.ExportDirectoryFlag, "-k")
	return output, error
}

func environmentSetHttpRequestTimeout(t *testing.T, args *SetTestArgs) (string, error) {
	apim := args.SrcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--http-request-timeout", strconv.Itoa(args.httpRequestTimeout), "-k")
	return output, error
}

func EnvironmentSetTokenType(t *testing.T, args *SetTestArgs) (string, error) {
	apim := args.SrcAPIM
	base.SetupEnvWithoutTokenFlag(t, apim.GetEnvName(), apim.GetApimURL())
	output, error := base.Execute(t, "set", "--token-type", args.TokenTypeFlag, "-k")
	return output, error
}
