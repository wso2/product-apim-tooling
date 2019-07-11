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

package cmd

import (
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createDirectoriesWithName(t *testing.T) {
	name, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "Temp directory should be created")
	err = createDirectories(filepath.Join(name, "lorem"))
	assert.Nil(t, err, "Should be no errors when creating directory structure")
	for _, dir := range dirs {
		dirPath := filepath.Join(name, "lorem", filepath.FromSlash(dir))
		assert.DirExists(t, dirPath, "Directory "+dirPath+" should be created")
	}
	_ = os.RemoveAll(name)
}

func Test_loadSwagger2(t *testing.T) {
	doc, err := loadSwagger("testdata/swaggers/swagger-2.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Simple API overview", doc.Spec().Info.Title, "Loads correct title")
	assert.NotNil(t, doc.Spec().Paths.Paths, "Paths should not be nil")
}

func Test_loadSwagger2YAML(t *testing.T) {
	doc, err := loadSwagger("testdata/swaggers/swagger-2.yaml")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Simple API overview", doc.Spec().Info.Title, "Loads correct title")
	assert.NotNil(t, doc.Spec().Paths.Paths, "Paths should not be nil")
}

func Test_loadSwagger3(t *testing.T) {
	doc, err := loadSwagger("testdata/swaggers/swagger-3.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Swagger Petstore", doc.Spec().Info.Title, "Loads correct title")
	assert.NotNil(t, doc.Spec().Paths.Paths, "Paths should not be nil")
}

func Test_loadSwagger3YAML(t *testing.T) {
	doc, err := loadSwagger("testdata/swaggers/swagger-3.yaml")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Swagger Petstore", doc.Spec().Info.Title, "Loads correct title")
	assert.NotNil(t, doc.Spec().Paths.Paths, "Paths should not be nil")
}

func Test_APIDefinition_generateFieldsFromSwagger(t *testing.T) {
	doc, err := loadSwagger("testdata/swaggers/swagger-3.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	def := &v2.APIDefinition{}
	err = v2.Swagger2Populate(def, doc)
	assert.Nil(t, err, "Populate without errors")
	assert.Equal(t, "SwaggerPetstore", def.ID.APIName, "Should correctly output name")
	assert.Equal(t, "/SwaggerPetstore/1.0.0", def.Context, "Should return correct context")
}
