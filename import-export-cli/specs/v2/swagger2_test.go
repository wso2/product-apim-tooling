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

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
)

func Test_swagger2WSO2Cors(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	cors, ok, err := swagger2XWSO2Cors(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, []string{"GET", "PUT", "POST"}, cors.AccessControlAllowMethods, "should have same elements for access control")
	assert.ElementsMatch(t, []string{"test.com", "example.com"}, cors.AccessControlAllowOrigins, "should have same elements for origins")
}

func Test_swagger2Tags(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	tags := swagger2Tags(doc)
	assert.ElementsMatch(t, []string{"pet", "user", "store"}, tags, "should have same elements")
}

func Test_swagger2WSO2ProductionEndpoints(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := swagger2XWSO2ProductionEndpoints(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_swagger2WSO2SandboxEndpoints(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := swagger2XWSO2SandboxEndpoints(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func TestSwagger2Populate(t *testing.T) {
	var def APIDTODefinition
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "Swagger Petstore", def.Name, "Should return correct api name")
	assert.Equal(t, "/petstore/v1/1.0.0", def.Context)
}

func TestSwagger2PopulateWithBasePath(t *testing.T) {
	var def1, def2 APIDTODefinition

	// Basepath without {version}
	doc1, err1 := loads.Spec("testdata/petstore_with_basepath1.yaml")
	assert.Nil(t, err1, "err should be nil")
	err1 = Swagger2Populate(&def1, doc1)
	assert.Nil(t, err1, "err should be nil")

	assert.Equal(t, "/petstore/v1/1.0.0", def1.Context)
	assert.Equal(t, true, def1.IsDefaultVersion)

	// Basepath with {version}
	doc2, err2 := loads.Spec("testdata/petstore_with_basepath2.yaml")
	assert.Nil(t, err2, "err should be nil")
	err1 = Swagger2Populate(&def2, doc2)
	assert.Nil(t, err2, "err should be nil")

	assert.Equal(t, "/petstore/v1/1.0.0", def2.Context)
	assert.Equal(t, false, def2.IsDefaultVersion)
}
