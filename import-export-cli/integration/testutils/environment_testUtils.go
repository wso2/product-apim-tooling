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
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

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

func EnvironmentSetHttpRequestTimeout(t *testing.T, args *SetTestArgs) (string, error) {
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

func ValidateThatRecievingTokenTypeIsChanged(t *testing.T, args *ApiGetKeyTestArgs, expectedTokenType string) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	var err error
	_, err = GetKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())
	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Error while getting key")

	tokenType := args.Apim.GetApplication(args.Apim.GetApplicationByName(utils.DefaultApictlTestAppName).ApplicationID).TokenType
	assert.Equal(t, strings.ToUpper(expectedTokenType), tokenType, "Error getting token type of application.")

	UnsubscribeAPI(args.Apim, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)
}

func ValidateExportDirectoryIsChanged(t *testing.T, args *SetTestArgs) {
	t.Helper()
	output, _ := EnvironmentSetExportDirectory(t, args)
	base.Log(output)
	assert.Contains(t, output, "Export Directory is set to", "Export Directory change is not successful")
}

func ValidateExportApisPassed(t *testing.T, args *InitTestArgs, directoryName string) {
	t.Helper()
	time.Sleep(5 * time.Second)

	output, error := ExportApisWithOneCommand(t, args)
	assert.Nil(t, error, "Error while Exporting APIs")
	assert.Contains(t, output, "export-apis execution completed", "Error while Exporting APIs")

	//Derive exported path from output
	exportedPath := base.GetExportedPathFromOutput(strings.ReplaceAll(output, "Command: export-apis execution completed !", ""))
	count, _ := base.CountFiles(exportedPath)
	assert.Equal(t, 1, count, "Error while Exporting APIs")

	t.Cleanup(func() {
		argsDefault := &SetTestArgs{
			SrcAPIM:             args.SrcAPIM,
			ExportDirectoryFlag: utils.DefaultExportDirPath,
		}
		ValidateExportDirectoryIsChanged(t, argsDefault)
		//Remove Exported apis
		base.RemoveDir(directoryName + utils.TestMigrationDirectorySuffix)
	})
}

func ValidateETokenTypeIsChanged(t *testing.T, args *SetTestArgs) {
	t.Helper()
	output, _ := EnvironmentSetTokenType(t, args)
	base.Log(output)
	assert.Contains(t, output, "Token type is set to", "1st attempt of Token Type change is not successful")
}
