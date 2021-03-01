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
	"crypto/tls"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func GetKeys(t *testing.T, provider string, name string, version string, env string) (string, error) {
	return base.Execute(t, "get", "keys", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")
}

func InvokeAPI(t *testing.T, url string, key string, expectedCode int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)

	assert.Nil(t, err, "Error while generating GET")

	base.SetDefaultRestAPIHeaders(key, req)

	base.LogRequest("invokeAPI()", req)

	response, err := invokeWithRetry(t, client, req)

	assert.Nil(t, err, "Error while invoking API")
	assert.Equal(t, expectedCode, response.StatusCode, "API Invocation failed")
}

func invokeAPIProduct(t *testing.T, url string, key string, expectedCode int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)

	assert.Nil(t, err, "Error while generating GET")

	base.SetDefaultRestAPIHeaders(key, req)

	base.LogRequest("invokeAPIProduct()", req)

	response, err := invokeWithRetry(t, client, req)

	assert.Nil(t, err, "Error while invoking API Product")
	assert.Equal(t, expectedCode, response.StatusCode, "API Product Invocation failed")
}

func invokeWithRetry(t *testing.T, client *http.Client, req *http.Request) (*http.Response, error) {
	response, err := client.Do(req)

	attempts := 0
	for response.StatusCode == http.StatusNotFound {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		response, err = client.Do(req)

		attempts++

		t.Log("invokeWithRetry() attempts = ", attempts)

		if attempts == base.GetMaxInvocationAttempts() {
			break
		}
	}

	return response, err
}

func ValidateGetKeysFailure(t *testing.T, args *ApiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	var err error
	var result string
	if args.Api != nil {
		result, err = GetKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())
	}

	if args.ApiProduct != nil {
		result, err = GetKeys(t, args.ApiProduct.Provider, args.ApiProduct.Name, utils.DefaultApiProductVersion, args.Apim.GetEnvName())
	}

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
}

func ValidateGetKeys(t *testing.T, args *ApiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	var err error
	var result string
	if args.Api != nil {
		result, err = GetKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		InvokeAPI(t, GetResourceURL(args.Apim, args.Api), base.GetValueOfUniformResponse(result), 200)
		UnsubscribeAPI(args.Apim, args.CtlUser.Username, args.CtlUser.Password, args.Api.ID)
	}

	if args.ApiProduct != nil {
		result, err = GetKeys(t, args.ApiProduct.Provider, args.ApiProduct.Name, utils.DefaultApiProductVersion, args.Apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		invokeAPIProduct(t, getResourceURLForAPIProduct(args.Apim, args.ApiProduct), base.GetValueOfUniformResponse(result), 200)
		UnsubscribeAPI(args.Apim, args.CtlUser.Username, args.CtlUser.Password, args.ApiProduct.ID)
	}
}

func ValidateGetKeysWithoutCleanup(t *testing.T, args *ApiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.Apim.GetEnvName(), args.Apim.GetApimURL(), args.Apim.GetTokenURL())
	base.Login(t, args.Apim.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	var err error
	var result string
	if args.Api != nil {
		result, err = GetKeys(t, args.Api.Provider, args.Api.Name, args.Api.Version, args.Apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		InvokeAPI(t, GetResourceURL(args.Apim, args.Api), base.GetValueOfUniformResponse(result), 200)
	}

	if args.ApiProduct != nil {
		result, err = GetKeys(t, args.ApiProduct.Provider, args.ApiProduct.Name, utils.DefaultApiProductVersion, args.Apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		invokeAPIProduct(t, getResourceURLForAPIProduct(args.Apim, args.ApiProduct), base.GetValueOfUniformResponse(result), 200)
	}
}
