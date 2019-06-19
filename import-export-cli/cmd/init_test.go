package cmd

import (
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
	sw, _, err := loadSwagger("testdata/swaggers/swagger-2.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Simple API overview", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_loadSwagger2YAML(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-2.yaml")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Simple API overview", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_loadSwagger3(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-3.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Swagger Petstore", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_loadSwagger3YAML(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-3.yaml")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Swagger Petstore", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_APIDefinition_generateFieldsFromSwagger(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-3.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	def := &APIDefinition{}
	def.generateFieldsFromSwagger3(sw)

	assert.Equal(t, "SwaggerPetstore", def.ID.APIName, "Should correctly output name")
	assert.Equal(t, "/SwaggerPetstore/1.0.0", def.Context, "Should return correct context")
	assert.Equal(t, 14, len(def.URITemplates), "Should return correct number of uri templates")
}
