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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

func ValidateGetAPILogLevel(t *testing.T, args *ApiLoggingTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl env
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Wait for indexing
	base.WaitForIndexing()

	// Get apictl output
	output, _ := getAPILogLevel(t, args)

	// Get REST API response
	response, _ := args.APIM.GetAPILogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.ApiId)

	// Validate output and response
	validateGetAPILogLevel(t, output, response)
}

func ValidateGetAPILogLevelError(t *testing.T, args *ApiLoggingTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl env
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Wait for indexing
	base.WaitForIndexing()

	// Get apictl output
	output, _ := getAPILogLevel(t, args)

	// Get REST API response
	_, err := args.APIM.GetAPILogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.ApiId)

	// Validate output and response
	validateInvalidPermissionError(t, output, err)
}

func ValidateSetAPILogLevel(t *testing.T, args *ApiLoggingTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl env
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Wait for indexing
	base.WaitForIndexing()

	// Get apictl output
	output, _ := setAPILogLevel(t, args)

	// Get REST API response
	response, _ := args.APIM.SetAPILogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.ApiId, args.LogLevel)

	// Validate output and response
	validateSetAPILogLevel(t, output, response)
}

func ValidateSetAPILogLevelError(t *testing.T, args *ApiLoggingTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.APIM.GetEnvName(), args.APIM.GetApimURL(), args.APIM.GetTokenURL())

	// Login to apictl env
	base.Login(t, args.APIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// Wait for indexing
	base.WaitForIndexing()

	// Get apictl output
	output, _ := setAPILogLevel(t, args)

	// Get REST API response
	_, err := args.APIM.SetAPILogLevel(args.CtlUser.Username, args.CtlUser.Password, args.TenantDomain, args.ApiId, args.LogLevel)

	// Validate output and response
	validateInvalidPermissionError(t, output, err)
}

func getAPILogLevel(t *testing.T, args *ApiLoggingTestArgs) (string, error) {
	tenant_domain := ""
	tenant_domain_flag := ""
	if args.TenantDomain != "" && args.TenantDomain != "carbon.super" {
		tenant_domain = args.TenantDomain
		tenant_domain_flag = "--tenant-domain"
	}

	api_id := ""
	api_id_flag := ""
	if args.ApiId != "" {
		api_id = args.ApiId
		api_id_flag = "--api-id"
	}

	output, err := base.Execute(t, "get", "api-logging", "-e", args.APIM.EnvName, tenant_domain_flag, tenant_domain, api_id_flag, api_id, "-k", "--verbose")
	return output, err
}

func setAPILogLevel(t *testing.T, args *ApiLoggingTestArgs) (string, error) {
	output, err := base.Execute(t, "set", "api-logging", "-e", args.APIM.EnvName, "--tenant-domain", args.TenantDomain, "--api-id", args.ApiId, "--log-level", args.LogLevel, "-k", "--verbose")
	return output, err
}

func validateGetAPILogLevel(t *testing.T, output string, response *apim.APILogLevelList) {
	for index, api := range response.Apis {
		log_line := strings.Split(output, "\n")[index+1]
		assert.Contains(t, log_line, api.ApiId, "API '"+api.ApiId+"' is not listed in the apictl output.")
		assert.Contains(t, log_line, api.LogLevel, "Log level of API '"+api.ApiId+"' is not matching.")
	}
}

func validateSetAPILogLevel(t *testing.T, output string, response *apim.APILogLevel) {
	assert.Equal(t, "Log level "+response.LogLevel+" is successfully set to the API.\n", output, "Invalid apictl output after setting API log level.")
}

func validateInvalidPermissionError(t *testing.T, output string, err error) {
	assert.Equal(t, "Exit status 1\n", output)
	assert.Contains(t, string(err.Error()), "403")
}
