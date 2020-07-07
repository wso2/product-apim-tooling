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
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"log"
	"net/http"
	"testing"
)

func getKeys(t *testing.T, provider string, name string, version string, env string) (string, error) {
	return base.Execute(t, "get-keys", "-n", name, "-v", version, "-r", provider, "-e", env, "-k")
}

func invokeAPI(t *testing.T, url string, key string, expectedCode int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)

	assert.Nil(t, err, "Error while generating GET")

	authHeader := "Bearer " + key
	req.Header.Set("Authorization", authHeader)

	t.Log("invokeAPI() url", url)

	response, err := client.Do(req)

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

	authHeader := "Bearer " + key
	req.Header.Set("Authorization", authHeader)

	t.Log("invokeAPIProduct() url", url)

	response, err := client.Do(req)

	assert.Nil(t, err, "Error while invoking API Product")
	assert.Equal(t, expectedCode, response.StatusCode, "API Product Invocation failed")
}

func validateGetKeysFailure(t *testing.T, args *apiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.apim.GetEnvName(), args.apim.GetApimURL(), args.apim.GetTokenURL())
	base.Login(t, args.apim.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	var err error
	var result string
	if args.api != nil {
		result, err = getKeys(t, args.api.Provider, args.api.Name, args.api.Version, args.apim.GetEnvName())
	}

	if args.apiProduct != nil {
		result, err = getKeys(t, args.apiProduct.Provider, args.apiProduct.Name, utils.DefaultApiProductVersion, args.apim.GetEnvName())
	}

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
}

func validateGetKeys(t *testing.T, args *apiGetKeyTestArgs) {
	t.Helper()

	base.SetupEnv(t, args.apim.GetEnvName(), args.apim.GetApimURL(), args.apim.GetTokenURL())
	base.Login(t, args.apim.GetEnvName(), args.ctlUser.Username, args.ctlUser.Password)

	var err error
	var result string
	if args.api != nil {
		result, err = getKeys(t, args.api.Provider, args.api.Name, args.api.Version, args.apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		invokeAPI(t, getResourceURL(args.apim, args.api), base.GetValueOfUniformResponse(result), 200)
		unsubscribeAPI(args.apim, args.ctlUser.Username, args.ctlUser.Password, args.api.ID)
	}

	if args.apiProduct != nil {
		result, err = getKeys(t, args.apiProduct.Provider, args.apiProduct.Name, utils.DefaultApiProductVersion, args.apim.GetEnvName())
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "Error while getting key")

		invokeAPIProduct(t, getResourceURLForAPIProduct(args.apim, args.apiProduct), base.GetValueOfUniformResponse(result), 200)
		unsubscribeAPI(args.apim, args.ctlUser.Username, args.ctlUser.Password, args.apiProduct.ID)
	}
}
