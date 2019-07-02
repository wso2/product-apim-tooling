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

package v2

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

var petstoreProdUrls = []string{"https://petstore.swagger.io/v2", "https://petstore.swagger.io/v2/1", "https://petstore.swagger.io/v2/2"}

func Test_oai3WSO2Basepath(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	basepath, ok, err := oai3WSO2Basepath(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.Equal(t, "/petstore/v1", basepath, "should return correct basepath")
}

func Test_oai3WSO2ProductionEndpoints(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := oai3XWSO2ProductionEndpoints(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_oai3WSO2SandboxEndpoints(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := oai3XWso2SandboxEndpoints(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_oai3Tags(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	tags := oai3Tags(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.ElementsMatch(t, []string{"pet", "user", "store"}, tags, "should have same elements")
}

func Test_oai3WSO2Cors(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	cors, ok, err := oai3XWSO2Cors(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, []string{"GET", "PUT", "POST"}, cors.AccessControlAllowMethods, "should have same elements for access control")
	assert.ElementsMatch(t, []string{"test.com", "example.com"}, cors.AccessControlAllowOrigins, "should have same elements for origins")
}
